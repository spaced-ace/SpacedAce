package pages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"spaced-ace/api/models"
	"spaced-ace/constants"
	"spaced-ace/context"
)

type EditQuizPageData struct {
	Session          *context.Session
	QuizWithMetaData QuizWithMetaData
}

type QuizWithMetaData struct {
	QuizInfo              models.QuizInfo
	QuestionsWithMetaData []QuestionWithMetaData
}

type QuestionWithMetaData struct {
	EditMode bool
	Question models.Question
}

type quizResponse struct {
	models.QuizInfo
	Questions []map[string]interface{}
}

func stringInArray(target string, arr []string) bool {
	for _, item := range arr {
		if item == target {
			return true
		}
	}
	return false
}

func getQuiz(quizId string, sessionId string) (*QuizWithMetaData, error) {
	if quizId == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Quiz ID is required")
	}

	req, err := http.NewRequest("GET", constants.BACKEND_URL+"/quizzes/"+quizId, nil)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error creating request")
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: sessionId,
	})
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error reading response body: "+err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, echo.NewHTTPError(resp.StatusCode, string(body))
	}

	var response quizResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error parsing response body: "+err.Error())
	}

	questionsWithMetaData := make([]QuestionWithMetaData, len(response.Questions))
	for i, question := range response.Questions {
		questionsWithMetaData[i] = QuestionWithMetaData{
			EditMode: false,
		}

		jsonData, err := json.Marshal(question)
		if err != nil {
			continue
		}
		questionType := models.ParseFloatToQuestionType(question["questionType"].(float64))

		if questionType == models.SingleChoice {
			var singleChoiceQuestion models.SingleChoiceQuestionResponseBody
			err = json.Unmarshal(jsonData, &singleChoiceQuestion)
			if err != nil {
				fmt.Println("Error unmarshalling single choice question: ", err)
				continue
			}
			questionsWithMetaData[i].Question = models.SingleChoiceQuestion{
				Id:           singleChoiceQuestion.Id,
				QuizId:       singleChoiceQuestion.QuizId,
				QuestionType: models.SingleChoice,
				Question:     singleChoiceQuestion.Question,
				Options: []models.Option{
					{Value: singleChoiceQuestion.Answers[0], Correct: singleChoiceQuestion.CorrectAnswer == "A"},
					{Value: singleChoiceQuestion.Answers[1], Correct: singleChoiceQuestion.CorrectAnswer == "B"},
					{Value: singleChoiceQuestion.Answers[2], Correct: singleChoiceQuestion.CorrectAnswer == "C"},
					{Value: singleChoiceQuestion.Answers[3], Correct: singleChoiceQuestion.CorrectAnswer == "D"},
				},
			}
		}
		if questionType == models.MultipleChoice {
			var multipleChoiceQuestion models.MultipleChoiceQuestionResponseBody
			err = json.Unmarshal(jsonData, &multipleChoiceQuestion)
			if err != nil {
				fmt.Println("Error unmarshalling multiple choice question: ", err)
				continue
			}
			questionsWithMetaData[i].Question = models.MultipleChoiceQuestion{
				Id:           multipleChoiceQuestion.Id,
				QuizId:       multipleChoiceQuestion.QuizId,
				QuestionType: models.MultipleChoice,
				Question:     multipleChoiceQuestion.Question,
				Options: []models.Option{
					{Value: multipleChoiceQuestion.Answers[0], Correct: stringInArray("A", multipleChoiceQuestion.CorrectAnswers)},
					{Value: multipleChoiceQuestion.Answers[1], Correct: stringInArray("B", multipleChoiceQuestion.CorrectAnswers)},
					{Value: multipleChoiceQuestion.Answers[2], Correct: stringInArray("C", multipleChoiceQuestion.CorrectAnswers)},
					{Value: multipleChoiceQuestion.Answers[3], Correct: stringInArray("D", multipleChoiceQuestion.CorrectAnswers)},
				},
			}
		}
		if questionType == models.TrueOrFalse {
			var trueOrFalseQuestion models.TrueOrFalseQuestionResponseBody
			err = json.Unmarshal(jsonData, &trueOrFalseQuestion)
			if err != nil {
				fmt.Println("Error unmarshalling multiple choice question: ", err)
				continue
			}
			questionsWithMetaData[i].Question = models.TrueOrFalseQuestion{
				Id:           trueOrFalseQuestion.Id,
				QuizId:       trueOrFalseQuestion.QuizId,
				QuestionType: models.TrueOrFalse,
				Question:     trueOrFalseQuestion.Question,
				Answer:       trueOrFalseQuestion.CorrectAnswer,
			}
		}
	}

	quizWithMetaData := QuizWithMetaData{
		QuizInfo: models.QuizInfo{
			Id:          response.Id,
			Title:       response.Title,
			Description: response.Description,
			CreatorId:   response.CreatorId,
			CreatorName: response.CreatorName,
		},
		QuestionsWithMetaData: questionsWithMetaData,
	}
	return &quizWithMetaData, nil
}

type UpdateQuizRequestBody struct {
	Title       string `json:"name"`
	Description string `json:"description"`
}

func updateQuiz(quizId string, sessionId string, title string, description string) (*models.QuizInfo, error) {
	if quizId == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Quiz ID is required")
	}

	requestBody, err := json.Marshal(UpdateQuizRequestBody{
		Title:       title,
		Description: description,
	})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error marshalling request body")
	}

	req, err := http.NewRequest("PATCH", constants.BACKEND_URL+"/quizzes/"+quizId, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error creating request")
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: sessionId,
	})
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error reading response body: "+err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, echo.NewHTTPError(resp.StatusCode, string(body))
	}

	var quizInfo models.QuizInfo
	err = json.Unmarshal(body, &quizInfo)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error parsing response body: "+err.Error())
	}

	return &quizInfo, nil
}

func EditQuizPage(c echo.Context) error {
	cc := c.(*context.Context)
	quizId := c.Param("id")

	quizWithMetaData, err := getQuiz(quizId, cc.Session.Id)
	if err != nil {
		return c.Redirect(http.StatusFound, "/not-found")
	}

	pageData := EditQuizPageData{
		Session:          cc.Session,
		QuizWithMetaData: *quizWithMetaData,
	}
	return c.Render(200, "edit-quiz", pageData)
}

type GenerateQuestionForm struct {
	QuizId  string `form:"quizId"`
	Context string `form:"context"`
}

type QuestionCreationRequestBody struct {
	QuizId string `json:"quizId"`
	Prompt string `json:"prompt"`
}

func PostGenerateQuestion(c echo.Context) error {
	cc := c.(*context.Context)

	questionType := cc.QueryParam("type")
	if questionType != "single-choice" && questionType != "multiple-choice" && questionType != "true-or-false" && questionType != "open-ended" {
		return c.String(http.StatusBadRequest, "Invalid question type")
	}

	var requestForm GenerateQuestionForm
	if err := c.Bind(&requestForm); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	if requestForm.QuizId == "" {
		return c.String(http.StatusBadRequest, "Quiz ID is required")
	}
	if requestForm.Context == "" {
		return c.String(http.StatusBadRequest, "Context is required")
	}

	requestBody, err := json.Marshal(QuestionCreationRequestBody{QuizId: requestForm.QuizId, Prompt: requestForm.Context})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	var req *http.Request
	if questionType == "single-choice" {
		req, _ = http.NewRequest("POST", constants.BACKEND_URL+"/questions/single-choice", bytes.NewBuffer(requestBody))
	}
	if questionType == "multiple-choice" {
		req, _ = http.NewRequest("POST", constants.BACKEND_URL+"/questions/multiple-choice", bytes.NewBuffer(requestBody))
	}
	if questionType == "true-or-false" {
		req, _ = http.NewRequest("POST", constants.BACKEND_URL+"/questions/true-or-false", bytes.NewBuffer(requestBody))
	}
	if questionType == "open-ended" {
		req, _ = http.NewRequest("POST", constants.BACKEND_URL+"/questions/open-ended", bytes.NewBuffer(requestBody))
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: cc.Session.Id,
	})
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.String(resp.StatusCode, "Error creating question")
	}

	if questionType == "single-choice" {
		var response models.SingleChoiceQuestionResponseBody
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error parsing response body")
		}
		question := QuestionWithMetaData{
			EditMode: true,
			Question: models.SingleChoiceQuestion{
				Id:           response.Id,
				QuizId:       response.QuizId,
				QuestionType: models.SingleChoice,
				Question:     response.Question,
				Options: []models.Option{
					{Value: response.Answers[0], Correct: response.CorrectAnswer == "A"},
					{Value: response.Answers[1], Correct: response.CorrectAnswer == "B"},
					{Value: response.Answers[2], Correct: response.CorrectAnswer == "C"},
					{Value: response.Answers[3], Correct: response.CorrectAnswer == "D"},
				},
			},
		}
		return c.Render(200, "single-choice-question", question)
	}
	if questionType == "multiple-choice" {
		var response models.MultipleChoiceQuestionResponseBody
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error parsing response body")
		}
		question := QuestionWithMetaData{
			EditMode: true,
			Question: models.MultipleChoiceQuestion{
				Id:           response.Id,
				QuizId:       response.QuizId,
				QuestionType: models.MultipleChoice,
				Question:     response.Question,
				Options: []models.Option{
					{Value: response.Answers[0], Correct: stringInArray("A", response.CorrectAnswers)},
					{Value: response.Answers[1], Correct: stringInArray("B", response.CorrectAnswers)},
					{Value: response.Answers[2], Correct: stringInArray("C", response.CorrectAnswers)},
					{Value: response.Answers[3], Correct: stringInArray("D", response.CorrectAnswers)},
				},
			},
		}
		return c.Render(200, "multiple-choice-question", question)
	}
	if questionType == "true-or-false" {
		var response models.TrueOrFalseQuestionResponseBody
		err = json.NewDecoder(resp.Body).Decode(&response)
		fmt.Println("Response: ", response)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error parsing response body")
		}
		question := QuestionWithMetaData{
			EditMode: true,
			Question: models.TrueOrFalseQuestion{
				Id:           response.Id,
				QuizId:       response.QuizId,
				QuestionType: models.TrueOrFalse,
				Question:     response.Question,
				Answer:       response.CorrectAnswer,
			},
		}
		fmt.Println("Question: ", question)
		return c.Render(200, "true-or-false-question", question)
	}

	return c.NoContent(http.StatusTeapot)
}

type UpdateQuizRequestForm struct {
	QuizId      string `form:"quizId"`
	Title       string `form:"title"`
	Description string `form:"description"`
}

func PatchUpdateQuiz(c echo.Context) error {
	var requestForm UpdateQuizRequestForm
	if err := c.Bind(&requestForm); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	if requestForm.QuizId == "" {
		return c.String(http.StatusBadRequest, "Quiz ID is required")
	}

	updatedQuizInfo, err := updateQuiz(requestForm.QuizId, c.(*context.Context).Session.Id, requestForm.Title, requestForm.Description)
	if err != nil {
		return c.String(err.(*echo.HTTPError).Code, err.Error())
	}

	if requestForm.Title != "" {
		return c.Render(200, "quiz-title-field", updatedQuizInfo)
	}
	if requestForm.Description != "" {
		return c.Render(200, "quiz-description-field", updatedQuizInfo)
	}

	return c.String(http.StatusBadRequest, "Title or description is required")
}

func DeleteQuestion(c echo.Context) error {
	questionId := c.Param("questionId")
	if questionId == "" {
		return c.String(http.StatusBadRequest, "Question ID is required")
	}

	return c.NoContent(http.StatusOK)
}
