package pages

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/api/models"
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

func EditQuiz(c echo.Context) error {
	cc := c.(*context.Context)

	pageData := EditQuizPageData{
		Session:          cc.Session,
		QuizWithMetaData: mockQuiz,
	}

	return c.Render(200, "edit-quiz", pageData)
}

type GenerateQuestionForm struct {
	QuizId  string `form:"quizId"`
	Context string `form:"context"`
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

	if questionType == "single-choice" {
		question := QuestionWithMetaData{
			EditMode: false,
			Question: mockSingleChoiceQuestion,
		}

		return c.Render(200, "single-choice-question", question)
	}

	if questionType == "multiple-choice" {
		question := QuestionWithMetaData{
			EditMode: false,
			Question: mockMultipleChoiceQuestion,
		}

		return c.Render(200, "multiple-choice-question", question)
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

	if requestForm.Title != "" {
		// Update the title
		return c.Render(200, "quiz-title-field", mockQuiz.QuizInfo)
	}

	if requestForm.Description != "" {
		// Update the description
		return c.Render(200, "quiz-description-field", mockQuiz.QuizInfo)
	}

	return c.NoContent(http.StatusBadRequest)
}
