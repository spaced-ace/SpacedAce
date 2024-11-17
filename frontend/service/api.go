package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"slices"
	"spaced-ace/constants"
	"spaced-ace/models"
	"spaced-ace/models/business"
	"spaced-ace/models/external"
	"spaced-ace/models/request"
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

		if questionType == models.MultipleChoice {
			var questionDto external.MultipleChoiceQuestionResponseBody
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

	// Reverse the order of questions
	slices.Reverse(questions)

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

func (a *ApiService) CreateQuizSession(userId, quizId string) (*business.QuizSession, error) {
	requestBody := &external.CreateQuizSessionRequestBody{
		UserID: userId,
		QuizID: quizId,
	}

	quizSessionDto := new(external.QuizSession)
	if err := a.getResponse("POST", "/quiz-sessions/start", requestBody, quizSessionDto); err != nil {
		return nil, err
	}

	return quizSessionDto.MapToBusiness()
}
func (a *ApiService) GetQuizSessions(userId, quizId string) ([]*business.QuizSession, error) {
	quizSessionsResponse := new(external.GetQuizSessionsResponseBody)
	if err := a.getResponse("GET", fmt.Sprintf("/quiz-sessions?userId=%s&quizId=%s", userId, quizId), nil, quizSessionsResponse); err != nil {
		return nil, err
	}

	quizSessions := make([]*business.QuizSession, quizSessionsResponse.Length)
	for i, quizSessionDto := range quizSessionsResponse.QuizSessions {
		quizSession, err := quizSessionDto.MapToBusiness()
		if err != nil {
			return nil, err
		}
		quizSessions[i] = quizSession
	}

	return quizSessions, nil
}
func (a *ApiService) GetQuizSession(quizSessionId string) (*business.QuizSession, error) {
	quizSessionResponse := new(external.QuizSession)
	if err := a.getResponse("GET", fmt.Sprintf("/quiz-sessions/%s", quizSessionId), nil, quizSessionResponse); err != nil {
		return nil, err
	}

	quizSession, err := quizSessionResponse.MapToBusiness()
	if err != nil {
		return nil, err
	}

	return quizSession, nil
}
func (a *ApiService) HasQuizSession(userId, quizId string) (bool, error) {
	if err := a.getResponse("GET", fmt.Sprintf("/quiz-sessions/has-session?userId=%s&quizId=%s", userId, quizId), nil, nil); err != nil {
		return false, nil
	}
	return true, nil
}

func (a *ApiService) SubmitQuiz(quizSessionId string) (*business.QuizResult, error) {
	quizResultDto := new(external.QuizResult)
	if err := a.getResponse("POST", fmt.Sprintf("/quiz-sessions/%s/submit", quizSessionId), nil, &quizResultDto); err != nil {
		return nil, fmt.Errorf("error submitting quiz session with ID `%s`: %w", quizSessionId, err)
	}

	quizResult, err := quizResultDto.MapToBusiness()
	if err != nil {
		return nil, fmt.Errorf("error mapping quiz result with ID `%s`: %w", quizResultDto.ID, err)
	}

	return quizResult, err
}
func (a *ApiService) GetQuizResult(quizSessionId string) (*business.QuizResult, error) {
	quizResultDto := new(external.QuizResult)
	if err := a.getResponse("GET", fmt.Sprintf("/quiz-sessions/%s/result", quizSessionId), nil, &quizResultDto); err != nil {
		return nil, fmt.Errorf("error getting quiz result for session with ID `%s`: %w", quizSessionId, err)
	}

	quizResult, err := quizResultDto.MapToBusiness()
	if err != nil {
		return nil, fmt.Errorf("error mapping quiz result with ID `%s`: %w", quizResultDto.ID, err)
	}

	return quizResult, err
}

func (a *ApiService) GetAnswers(quizSessionId string) (*business.AnswerLists, error) {
	answersResponse := new(external.AnswersResponse)
	if err := a.getResponse("GET", fmt.Sprintf("/quiz-sessions/%s/answers", quizSessionId), nil, &answersResponse); err != nil {
		return nil, err
	}

	answerLists, err := answersResponse.MapToBusiness()
	if err != nil {
		return nil, err
	}
	return answerLists, nil
}
func (a *ApiService) CreateOrUpdateSingleChoiceAnswer(quizSessionId, questionId string, answer string) (*business.SingleChoiceAnswer, error) {
	requestBody := external.NewSingleChoiceAnswerRequestBody(questionId, answer)

	responseBody := new(external.SingleChoiceAnswer)
	if err := a.getResponse("PUT", fmt.Sprintf("/quiz-sessions/%s/answers", quizSessionId), requestBody, responseBody); err != nil {
		return nil, err
	}

	singleChoiceAnswer, err := responseBody.MapToBusiness()
	if err != nil {
		return nil, err
	}
	return singleChoiceAnswer, nil
}
func (a *ApiService) CreateOrUpdateMultipleChoiceAnswer(quizSessionId, questionId string, answers []string) (*business.MultipleChoiceAnswer, error) {
	requestBody := external.NewMultipleChoiceAnswerRequestBody(questionId, answers)

	responseBody := new(external.MultipleChoiceAnswer)
	if err := a.getResponse("PUT", fmt.Sprintf("/quiz-sessions/%s/answers", quizSessionId), requestBody, responseBody); err != nil {
		return nil, err
	}

	multipleChoiceAnswer, err := responseBody.MapToBusiness()
	if err != nil {
		return nil, err
	}
	return multipleChoiceAnswer, nil
}
func (a *ApiService) CreateOrUpdateTrueOrFalseAnswer(quizSessionId, questionId string, answer bool) (*business.TrueOrFalseAnswer, error) {
	requestBody := external.NewTrueOrFalseAnswerRequestBody(questionId, answer)

	responseBody := new(external.TrueOrFalseAnswer)
	if err := a.getResponse("PUT", fmt.Sprintf("/quiz-sessions/%s/answers", quizSessionId), requestBody, responseBody); err != nil {
		return nil, err
	}

	trueOrFalseAnswer, err := responseBody.MapToBusiness()
	if err != nil {
		return nil, err
	}
	return trueOrFalseAnswer, nil
}

func (a *ApiService) GetQuizHistory(userID string) ([]business.QuizHistoryEntry, error) {
	responseBody := new(external.QuizHistoryEntriesResponseBody)
	if err := a.getResponse("GET", fmt.Sprintf("/quiz-history?userID=%s", userID), nil, responseBody); err != nil {
		return []business.QuizHistoryEntry{}, err
	}

	entries := make([]business.QuizHistoryEntry, 0, responseBody.Length)
	for _, e := range responseBody.QuizHistoryEntries {
		entry, err := e.MapToBusiness()
		if err != nil {
			return []business.QuizHistoryEntry{}, err
		}
		entries = append(entries, *entry)
	}

	return entries, nil
}

func (a *ApiService) GetLearnList() (*business.LearnList, error) {
	responseBody := new(external.LearnListResponseBody)
	if err := a.getResponse("GET", "/learn-list", nil, responseBody); err != nil {
		return nil, err
	}

	learnList, err := responseBody.MapToBusiness()
	if err != nil {
		return nil, err
	}

	return learnList, nil
}
func (a *ApiService) AddQuizToLearnList(quizID string) (*business.LearnList, error) {
	responseBody := new(external.LearnListResponseBody)
	if err := a.getResponse("POST", fmt.Sprintf("/learn-list/%s/add", quizID), nil, responseBody); err != nil {
		return nil, err
	}

	learnList, err := responseBody.MapToBusiness()
	if err != nil {
		return nil, err
	}

	return learnList, nil
}
func (a *ApiService) RemoveQuizFromLearnList(quizID string) (*business.LearnList, error) {
	responseBody := new(external.LearnListResponseBody)
	if err := a.getResponse("POST", fmt.Sprintf("/learn-list/%s/remove", quizID), nil, responseBody); err != nil {
		return nil, err
	}

	learnList, err := responseBody.MapToBusiness()
	if err != nil {
		return nil, err
	}

	return learnList, nil
}

func (a *ApiService) GetReviewItemListData(quizID, difficulty, status, query string, page int) (reviewItems []business.ReviewItem, quizOptions []business.Option, maxReviewItemCount int, err error) {
	reviewItemsRequestBody := external.ReviewItemsRequestBody{
		QuizID:     quizID,
		Difficulty: difficulty,
		Status:     status,
		Page:       page,
		Query:      query,
	}

	reviewItemsResponseBody := new(external.ReviewItemResponseBody)
	if err = a.getResponse("GET", "/review-items", reviewItemsRequestBody, reviewItemsResponseBody); err != nil {
		return nil, nil, 0, fmt.Errorf("getting review items: %w\n", err)
	}

	reviewItems = make([]business.ReviewItem, 0, len(reviewItemsResponseBody.ReviewItems))
	for _, item := range reviewItemsResponseBody.ReviewItems {
		reviewItem, err := item.MapToBusiness()
		if err != nil {
			return nil, nil, 0, fmt.Errorf("mapping review item %+v: %w\n", reviewItem, err)
		}
		reviewItems = append(reviewItems, *reviewItem)
	}

	maxReviewItemCount = reviewItemsResponseBody.ReviewItemCountForFilter

	quizOptionsResponseBody := new(external.QuizOptionsResponseBody)
	if err = a.getResponse("GET", "/review-items/quiz-options", nil, quizOptionsResponseBody); err != nil {
		return nil, nil, 0, fmt.Errorf("getting quiz options: %w\n", err)
	}

	quizOptions = make([]business.Option, 0, len(quizOptionsResponseBody.QuizOptions))
	for _, option := range quizOptionsResponseBody.QuizOptions {
		quizOption, err := option.MapToBusiness()
		if err != nil {
			return nil, nil, 0, fmt.Errorf("mapping quiz options: %w\n", err)
		}
		quizOptions = append(quizOptions, *quizOption)
	}

	return reviewItems, quizOptions, maxReviewItemCount, nil
}
func (a *ApiService) GetReviewItemCounts() (total int, dueToReview int, err error) {
	response := new(external.ReviewItemCountsResponseBody)
	if err := a.getResponse("GET", "/review-items/item-counts", nil, response); err != nil {
		return 0, 0, err
	}
	return response.Total, response.DueToReview, nil
}
func (a *ApiService) GetReviewItemQuestion(reviewItemID string) (*business.ReviewItemQuestionData, error) {
	url := "/review-items/get-question"
	if reviewItemID != "" {
		url = fmt.Sprintf("%s/%s", url, reviewItemID)
	}

	requestBody := new(external.ReviewItemQuestionResponseBody)
	if err := a.getResponse("GET", url, nil, requestBody); err != nil {
		return nil, fmt.Errorf("getting question for review item with ID %q: %w\n", reviewItemID, err)
	}

	reviewItemPageData, err := requestBody.MapToBusiness()
	if err != nil {
		return nil, fmt.Errorf("mapping review item page data: %w\n", err)
	}

	return reviewItemPageData, nil
}
func (a *ApiService) SubmitReviewItemQuestion(reviewItemID string, form request.SubmitReviewItemQuestionForm) (*business.ReviewItem, error) {
	requestBody := external.SubmitReviewItemQuestionRequestBody{
		SingleChoiceValue:   form.SingleChoiceValue,
		MultipleChoiceValue: form.MultipleChoiceValue,
		TrueOrFalseValue:    form.TrueOrFalseValue,
	}

	responseBody := new(external.ReviewItem)
	if err := a.getResponse("POST", fmt.Sprintf("/review-items/%s/submit", reviewItemID), requestBody, responseBody); err != nil {
		return nil, err
	}

	reviewItem, err := responseBody.MapToBusiness()
	if err != nil {
		return nil, err
	}

	return reviewItem, nil
}
