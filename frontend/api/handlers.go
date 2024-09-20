package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/context"
	"spaced-ace/models"
	"spaced-ace/models/business"
	"spaced-ace/models/request"
	"spaced-ace/render"
	"spaced-ace/utils"
	"spaced-ace/views/forms"
)

func handleCreateQuiz(c echo.Context) error {
	cc := c.(*context.AppContext)

	var requestForm request.CreateQuizRequestForm
	if err := c.Bind(&requestForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	errors := map[string]string{}
	if requestForm.Title == "" {
		errors["title"] = "Title is required"
	}
	if requestForm.Description == "" {
		errors["description"] = "Description is required"
	}
	if len(errors) > 0 {
		return render.TemplRender(c, 200, forms.CreateQuizForm(requestForm, errors))
	}

	quizInfo, err := cc.ApiService.CreateQuiz(requestForm.Title, requestForm.Description)
	if err != nil {
		errors["other"] = "Error creating a quiz: " + err.Error()
		return render.TemplRender(c, 200, forms.CreateQuizForm(requestForm, errors))
	}

	c.Response().Header().Set("HX-Redirect", "/quizzes/"+quizInfo.Id+"/edit")
	return c.NoContent(http.StatusCreated)
}
func handleGenerateQuestion(c echo.Context) error {
	cc := c.(*context.AppContext)

	questionType := cc.QueryParam("type")
	if questionType != "single-choice" && questionType != "multiple-choice" && questionType != "true-or-false" && questionType != "open-ended" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid question type")
	}

	var requestForm request.GenerateQuestionForm
	if err := c.Bind(&requestForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	if requestForm.QuizId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Quiz ID is required")
	}
	if requestForm.Context == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Context is required")
	}

	switch questionType {
	case "single-choice":
		{
			question, err := cc.ApiService.GenerateSingleChoiceQuestion(requestForm.QuizId, requestForm.Context)
			if err != nil {
				return err
			}
			data := NewComponentTemplate(
				cc.Session,
				business.QuestionWithMetaData{
					EditMode: true,
					Question: question,
				},
			)
			return c.Render(200, "single-choice-question", data)
		}
	case "multiple-choice":
		{
			question, err := cc.ApiService.GenerateMultipleChoiceQuestion(requestForm.QuizId, requestForm.Context)
			if err != nil {
				return err
			}
			data := NewComponentTemplate(
				cc.Session,
				business.QuestionWithMetaData{
					EditMode: true,
					Question: question,
				},
			)
			return c.Render(200, "multiple-choice-question", data)
		}
	case "true-or-false":
		{
			question, err := cc.ApiService.GenerateTrueOrFalseQuestion(requestForm.QuizId, requestForm.Context)
			if err != nil {
				return err
			}
			data := NewComponentTemplate(
				cc.Session,
				business.QuestionWithMetaData{
					EditMode: true,
					Question: question,
				},
			)
			return c.Render(200, "true-or-false-question", data)
		}
	}
	return echo.NewHTTPError(400, "Invalid question type")
}

func handleUpdateQuiz(c echo.Context) error {
	cc := c.(*context.AppContext)

	var requestForm request.UpdateQuizRequestForm
	if err := c.Bind(&requestForm); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	if requestForm.QuizId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Quiz ID is required")
	}

	updatedQuizInfo, err := cc.ApiService.UpdateQuiz(requestForm.QuizId, requestForm.Title, requestForm.Description)
	if err != nil {
		return err
	}

	if requestForm.Title != "" {
		data := NewComponentTemplate(cc.Session, updatedQuizInfo)
		return c.Render(200, "quiz-title-field", data)
	}
	if requestForm.Description != "" {
		data := NewComponentTemplate(cc.Session, updatedQuizInfo)
		return c.Render(200, "quiz-description-field", data)
	}

	return echo.NewHTTPError(http.StatusBadRequest, "Title or description is required")
}
func handleDeleteQuestion(c echo.Context) error {
	cc := c.(*context.AppContext)

	questionId := c.Param("questionId")
	if questionId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Question ID is required")
	}

	questionType := c.QueryParam("type")
	if !utils.StringInArray(questionType, []string{models.SingleChoiceQuestion, models.MultipleChoiceQuestion, models.TrueOrFalseQuestion}) {
		return echo.NewHTTPError(400, "Invalid question type: "+questionType)
	}

	quizId := c.QueryParam("quizId")
	if quizId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Quiz ID is required")
	}

	err := cc.ApiService.DeleteQuestion(questionType, quizId, questionId)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
func handleDeleteQuiz(c echo.Context) error {
	cc := c.(*context.AppContext)

	quizId := c.Param("quizId")
	if quizId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Quiz ID is required")
	}

	if err := cc.ApiService.DeleteQuiz(quizId); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
