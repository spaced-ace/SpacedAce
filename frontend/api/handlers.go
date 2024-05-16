package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"spaced-ace/constants"
	"spaced-ace/context"
	"spaced-ace/models"
	"spaced-ace/utils"
)

// TODO move to api
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

type QuizRequestBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type CreateQuizRequestForm struct {
	Title       string `form:"title"`
	Description string `form:"description"`
}
type UpdateQuizRequestBody struct {
	Title       string `json:"name"`
	Description string `json:"description"`
}
type GenerateQuestionForm struct {
	QuizId  string `form:"quizId"`
	Context string `form:"context"`
}
type QuestionCreationRequestBody struct {
	QuizId string `json:"quizId"`
	Prompt string `json:"prompt"`
}
type UpdateQuizRequestForm struct {
	QuizId      string `form:"quizId"`
	Title       string `form:"title"`
	Description string `form:"description"`
}

func handleCreateQuiz(c echo.Context) error {
	cc := c.(*context.AppContext)

	var requestForm CreateQuizRequestForm
	if err := c.Bind(&requestForm); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	if requestForm.Title == "" {
		return c.String(http.StatusBadRequest, "Title is required")
	}
	if requestForm.Description == "" {
		return c.String(http.StatusBadRequest, "Description is required")
	}

	requestBody, _ := json.Marshal(QuizRequestBody{
		Name:        requestForm.Title,
		Description: requestForm.Description,
	})

	req, _ := http.NewRequest("POST", constants.BACKEND_URL+"/quizzes/create", bytes.NewBuffer(requestBody))
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: cc.Session.Id,
	})
	client := &http.Client{}

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	var responseBody models.QuizInfo
	_ = json.NewDecoder(resp.Body).Decode(&responseBody)

	c.Response().Header().Set("HX-Redirect", "/quizzes/"+responseBody.Id+"/edit")
	return c.NoContent(http.StatusCreated)
}
func handleGenerateQuestion(c echo.Context) error {
	cc := c.(*context.AppContext)

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
					{Value: response.Answers[0], Correct: utils.StringInArray("A", response.CorrectAnswers)},
					{Value: response.Answers[1], Correct: utils.StringInArray("B", response.CorrectAnswers)},
					{Value: response.Answers[2], Correct: utils.StringInArray("C", response.CorrectAnswers)},
					{Value: response.Answers[3], Correct: utils.StringInArray("D", response.CorrectAnswers)},
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
func handleUpdateQuiz(c echo.Context) error {
	cc := c.(*context.AppContext)

	var requestForm UpdateQuizRequestForm
	if err := c.Bind(&requestForm); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	if requestForm.QuizId == "" {
		return c.String(http.StatusBadRequest, "Quiz ID is required")
	}

	updatedQuizInfo, err := updateQuiz(requestForm.QuizId, c.(*context.AppContext).Session.Id, requestForm.Title, requestForm.Description)
	if err != nil {
		return c.String(err.(*echo.HTTPError).Code, err.Error())
	}

	if requestForm.Title != "" {
		data := NewComponentTemplate(cc.Session, updatedQuizInfo)
		return c.Render(200, "quiz-title-field", data)
	}
	if requestForm.Description != "" {
		data := NewComponentTemplate(cc.Session, updatedQuizInfo)
		return c.Render(200, "quiz-description-field", data)
	}

	return c.String(http.StatusBadRequest, "Title or description is required")
}
func handleDeleteQuestion(c echo.Context) error {
	questionId := c.Param("questionId")
	if questionId == "" {
		return c.String(http.StatusBadRequest, "Question ID is required")
	}

	return c.NoContent(http.StatusOK)
}
func handleDeleteQuiz(c echo.Context) error {
	cc := c.(*context.AppContext)
	quizId := c.Param("quizId")

	req, _ := http.NewRequest("DELETE", constants.BACKEND_URL+"/quizzes/"+quizId, nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: cc.Session.Id,
	})
	client := &http.Client{}

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	return c.NoContent(http.StatusOK)
}
