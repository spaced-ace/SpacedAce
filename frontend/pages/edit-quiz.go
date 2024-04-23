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

var mockSingleChoiceQuestion = models.NewSingleChoiceQuestion(
	"1",
	"ae664251-9ee7-4ca6-9f16-ff072de61632",
	1,
	"What is the capital of France?",
	[]models.Option{
		{Value: "Paris", Correct: true},
		{Value: "London", Correct: false},
		{Value: "Berlin", Correct: false},
		{Value: "Madrid", Correct: false},
	})

var mockMultipleChoiceQuestion = models.NewMultipleChoiceQuestion(
	"2",
	"ae664251-9ee7-4ca6-9f16-ff072de61632",
	2,
	"Which of the following are European countries?",
	[]models.Option{
		{Value: "Canada", Correct: false},
		{Value: "France", Correct: true},
		{Value: "Germany", Correct: true},
		{Value: "Brazil", Correct: false},
	})

var mockQuiz = QuizWithMetaData{
	QuizInfo: models.QuizInfo{
		Id:          "ae664251-9ee7-4ca6-9f16-ff072de61632",
		Title:       "My QuizWithMetaData",
		Description: "This is a quiz",
		CreatorId:   "73975759-99f9-46be-b84b-cfa4d2222112",
		CreatorName: "John Doe",
	},
	QuestionsWithMetaData: []QuestionWithMetaData{
		{
			EditMode: false,
			Question: mockSingleChoiceQuestion,
		},
	},
}

var mockTrueOrFalseQuestion = models.NewTrueOrFalseQuestion(
	"3",
	"ae664251-9ee7-4ca6-9f16-ff072de61632",
	3,
	"Is the sun hot?",
	true)

var mockOpenEndedQuestion = models.NewOpenEndedQuestion(
	"4",
	"ae664251-9ee7-4ca6-9f16-ff072de61632",
	4,
	"What are the main benefits of using Go programming language?",
	"",
	"Go has a simple syntax which makes it easy to learn. It is statically typed and compiled, which helps in catching errors early. It also has built-in support for concurrent programming.")

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

	var quizResponse models.Quiz
	err = json.Unmarshal(body, &quizResponse)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error parsing response body: "+err.Error())
	}

	quizWithMetaData := QuizWithMetaData{
		QuizInfo: models.QuizInfo{
			Id:          quizResponse.Id,
			Title:       quizResponse.Title,
			Description: quizResponse.Description,
			CreatorId:   quizResponse.CreatorId,
			CreatorName: quizResponse.CreatorName,
		},
		QuestionsWithMetaData: []QuestionWithMetaData{},
	}
	// TODO get questions and map them to QuestionWithMetaData
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
	var createdQuestion models.Question
	err = json.NewDecoder(resp.Body).Decode(&createdQuestion)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	if resp.StatusCode != http.StatusOK {
		return c.String(resp.StatusCode, "Error creating question")
	}

	fmt.Println("Created question:", createdQuestion)

	if questionType == "single-choice" {
		return c.Render(200, "single-choice-question", createdQuestion)
	}
	if questionType == "multiple-choice" {
		return c.Render(200, "multiple-choice-question", createdQuestion)
	}
	if questionType == "true-or-false" {
		return c.Render(200, "true-or-false-question", createdQuestion)
	}
	if questionType == "open-ended" {
		return c.Render(200, "open-ended-question", createdQuestion)
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
