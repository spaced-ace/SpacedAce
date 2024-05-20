package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"spaced-ace/constants"
	"spaced-ace/models"
	"spaced-ace/models/business"
	"spaced-ace/models/external"
)

type ApiService struct {
	sessionCookie *http.Cookie
	client        *http.Client
}

func NewApiService(sessionCookie *http.Cookie) *ApiService {
	return &ApiService{
		sessionCookie: sessionCookie,
		client:        &http.Client{},
	}
}

func (a *ApiService) getResponse(method, path string, requestBody any, responseBody interface{}) error {
	data, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, constants.BACKEND_URL+path, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if a.sessionCookie != nil {
		req.AddCookie(a.sessionCookie)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var result map[string]string
		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			return err
		}

		return errors.New("error message: " + result["message"])
	}

	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return err
	}

	return nil
}

func (a *ApiService) GetSession() (*business.Session, error) {
	session := new(business.Session)
	if err := a.getResponse("GET", "/authenticated", nil, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (a *ApiService) DeleteSession() error {
	return a.getResponse("POST", "/logout", nil, nil)
}

func (a *ApiService) GetQuiz(quizId string) (*business.Quiz, error) {
	if quizId == "" {
		return nil, echo.NewHTTPError(400, fmt.Sprintf("Invalid quiz id: %s", quizId))
	}

	quizDTO := new(external.Quiz)
	if err := a.getResponse("GET", "/quizzes/"+quizId, nil, quizDTO); err != nil {
		return nil, err
	}

	var questions []interface{}
	for _, rawQuestion := range quizDTO.Questions {
		var questionDto map[string]interface{}
		if err := json.Unmarshal(rawQuestion, &questionDto); err != nil {
			continue
		}
		questionType := models.ParseFloatToQuestionType(questionDto["questionType"].(float64))

		if questionType == models.SingleChoice {
			var questionDto external.SingleChoiceQuestionResponseBody
			if err := json.Unmarshal(rawQuestion, &questionDto); err != nil {
				continue
			}
			question, err := questionDto.MapToBusiness()
			if err != nil {
				continue
			}
			questions = append(questions, question)
		}

		if questionType == models.TrueOrFalse {
			var questionDto external.TrueOrFalseQuestionResponseBody
			if err := json.Unmarshal(rawQuestion, &questionDto); err != nil {
				continue
			}
			question, err := questionDto.MapToBusiness()
			if err != nil {
				continue
			}
			questions = append(questions, question)
		}
	}

	quiz := &business.Quiz{
		QuizInfo: business.QuizInfo{
			Id:          quizDTO.Id,
			Title:       quizDTO.Title,
			Description: quizDTO.Description,
			CreatorId:   quizDTO.CreatorId,
			CreatorName: quizDTO.CreatorName,
		},
		Questions: questions,
	}

	return quiz, nil
}

func (a *ApiService) GenerateSingleChoiceQuestion(quizId, context string) (*business.SingleChoiceQuestion, error) {
	requestBody := external.GenerateQuestionRequestBody{
		QuizId: quizId,
		Prompt: context,
	}

	questionDTO := new(external.SingleChoiceQuestionResponseBody)
	if err := a.getResponse("POST", "/questions/single-choice", requestBody, questionDTO); err != nil {
		return nil, err
	}

	return questionDTO.MapToBusiness()
}
func (a *ApiService) GenerateMultipleChoiceQuestion(quizId, context string) (*business.MultipleChoiceQuestion, error) {
	requestBody := external.GenerateQuestionRequestBody{
		QuizId: quizId,
		Prompt: context,
	}

	questionDTO := new(external.MultipleChoiceQuestionResponseBody)
	if err := a.getResponse("POST", "/questions/multiple-choice", requestBody, questionDTO); err != nil {
		return nil, err
	}

	return questionDTO.MapToBusiness()
}
func (a *ApiService) GenerateTrueOrFalseQuestion(quizId, context string) (*business.TrueOrFalseQuestion, error) {
	requestBody := external.GenerateQuestionRequestBody{
		QuizId: quizId,
		Prompt: context,
	}

	questionDTO := new(external.TrueOrFalseQuestionResponseBody)
	if err := a.getResponse("POST", "/questions/true-or-false", requestBody, questionDTO); err != nil {
		return nil, err
	}

	return questionDTO.MapToBusiness()
}
func (a *ApiService) DeleteQuestion(questionType, quizId, questionId string) error {
	return a.getResponse("DELETE", fmt.Sprintf("/questions/%s/%s/%s", questionType, quizId, questionId), nil, nil)
}

func (a *ApiService) GetQuizzesInfos(userId string) ([]business.QuizInfo, error) {
	quizzesDTO := new(external.QuizInfosResponse)
	err := a.getResponse("GET", "/quizzes/user/"+userId, nil, quizzesDTO)
	if err != nil {
		return nil, err
	}

	var quizInfos []business.QuizInfo
	for _, q := range quizzesDTO.Quizzes {
		quizInfos = append(quizInfos, business.QuizInfo{
			Id:          q.Id,
			Title:       q.Title,
			Description: q.Description,
			CreatorId:   q.CreatorId,
			CreatorName: q.CreatorName,
		})
	}

	return quizInfos, err
}
func (a *ApiService) CreateQuiz(title, description string) (*business.QuizInfo, error) {
	requestBody := external.CreateQuizRequestBody{
		Name:        title,
		Description: description,
	}

	quizInfoDto := new(external.QuizInfo)
	err := a.getResponse("POST", "/quizzes/create", requestBody, quizInfoDto)
	if err != nil {
		return nil, err
	}

	return &business.QuizInfo{
		Id:          quizInfoDto.Id,
		Title:       quizInfoDto.Title,
		Description: quizInfoDto.Description,
		CreatorId:   quizInfoDto.CreatorId,
		CreatorName: quizInfoDto.CreatorName,
	}, nil
}
func (a *ApiService) UpdateQuiz(quizId, title, description string) (*business.QuizInfo, error) {
	requestBody := &external.UpdateQuizRequestBody{
		Title:       title,
		Description: description,
	}

	quizInfoDto := new(external.QuizInfo)
	if err := a.getResponse("PATCH", "/quizzes/"+quizId, requestBody, quizInfoDto); err != nil {
		return nil, err
	}

	return &business.QuizInfo{
		Id:          quizInfoDto.Id,
		Title:       quizInfoDto.Title,
		Description: quizInfoDto.Description,
		CreatorId:   quizInfoDto.CreatorId,
		CreatorName: quizInfoDto.CreatorName,
	}, nil
}
func (a *ApiService) DeleteQuiz(quizId string) error {
	return a.getResponse("DELETE", "/quizzes/"+quizId, nil, nil)
}
