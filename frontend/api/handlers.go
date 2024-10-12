package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/context"
	"spaced-ace/models"
	"spaced-ace/models/request"
	"spaced-ace/render"
	"spaced-ace/utils"
	"spaced-ace/views/components"
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
func handleGenerateQuestionStart(c echo.Context) error {
	errors := map[string]string{}

	var requestForm request.GenerateQuestionForm
	if err := c.Bind(&requestForm); err != nil {
		errors["other"] = "Parsing error: " + err.Error()
		return render.TemplRender(
			c,
			200,
			forms.GenerateQuestionForm(
				false,
				nil,
				requestForm,
				errors,
			),
		)
	}

	questionType := requestForm.QuestionType
	if questionType != models.SingleChoiceQuestion && questionType != models.MultipleChoiceQuestion && questionType != models.TrueOrFalseQuestion {
		if questionType == "open-ended" {
			errors["other"] = fmt.Sprintf("Currently not suppored question type: '%s'.", questionType)
		} else {
			errors["other"] = fmt.Sprintf("Invalid question type: '%s'.", questionType)
		}
		return render.TemplRender(
			c,
			200,
			forms.GenerateQuestionForm(
				false,
				nil,
				requestForm,
				errors,
			),
		)
	}

	if requestForm.QuizId == "" {
		errors["other"] = "quizId is required"
	}
	if requestForm.Context == "" {
		errors["context"] = "Context is required"
	}
	if len(errors) > 0 {
		return render.TemplRender(
			c,
			200,
			forms.GenerateQuestionForm(
				false,
				nil,
				requestForm,
				errors,
			),
		)
	}

	return render.TemplRender(
		c,
		200,
		forms.GenerateQuestionForm(
			true,
			nil,
			requestForm,
			errors,
		),
	)
}
func handleGenerateQuestion(c echo.Context) error {
	errors := map[string]string{}
	cc := c.(*context.AppContext)

	var requestForm request.GenerateQuestionForm
	if err := c.Bind(&requestForm); err != nil {
		errors["other"] = "Parsing error: " + err.Error()
		return render.TemplRender(
			c,
			200,
			forms.GenerateQuestionForm(
				false,
				nil,
				requestForm,
				errors,
			),
		)
	}

	questionType := requestForm.QuestionType
	if questionType != models.SingleChoiceQuestion && questionType != models.MultipleChoiceQuestion && questionType != models.TrueOrFalseQuestion {
		if questionType == "open-ended" {
			errors["other"] = fmt.Sprintf("Currently not suppored question type: '%s'.", questionType)
		} else {
			errors["other"] = fmt.Sprintf("Invalid question type: '%s'.", questionType)
		}
		return render.TemplRender(
			c,
			200,
			forms.GenerateQuestionForm(
				false,
				nil,
				requestForm,
				errors,
			),
		)
	}

	if requestForm.QuizId == "" {
		errors["other"] = "quizId is required"
	}
	if requestForm.Context == "" {
		errors["context"] = "Context is required"
	}

	switch questionType {
	case models.SingleChoiceQuestion:
		{
			question, err := cc.ApiService.GenerateSingleChoiceQuestion(requestForm.QuizId, requestForm.Context)
			if err != nil {
				errors["other"] = "Error generating question: " + err.Error()
				return render.TemplRender(
					c,
					200,
					forms.GenerateQuestionForm(
						false,
						nil,
						requestForm,
						errors,
					),
				)
			}

			questionComponent := components.SingleChoiceQuestion(
				question,
				true,
				true,
			)
			return render.TemplRender(
				c,
				200,
				forms.GenerateQuestionForm(
					false,
					questionComponent,
					requestForm,
					errors,
				),
			)
		}
	case models.MultipleChoiceQuestion:
		{
			question, err := cc.ApiService.GenerateMultipleChoiceQuestion(requestForm.QuizId, requestForm.Context)
			if err != nil {
				errors["other"] = "Error generating question: " + err.Error()
				return render.TemplRender(
					c,
					200,
					forms.GenerateQuestionForm(
						false,
						nil,
						requestForm,
						errors,
					),
				)
			}

			questionComponent := components.MultipleChoiceQuestion(
				question,
				true,
				true,
			)
			return render.TemplRender(
				c,
				200,
				forms.GenerateQuestionForm(
					false,
					questionComponent,
					requestForm,
					errors,
				),
			)
		}
	default:
		{
			question, err := cc.ApiService.GenerateTrueOrFalseQuestion(requestForm.QuizId, requestForm.Context)
			if err != nil {
				errors["other"] = "Error generating question: " + err.Error()
				return render.TemplRender(
					c,
					200,
					forms.GenerateQuestionForm(
						false,
						nil,
						requestForm,
						errors,
					),
				)
			}

			questionComponent := components.TrueOrFalseQuestion(
				question,
				true,
				true,
			)
			return render.TemplRender(
				c,
				200,
				forms.GenerateQuestionForm(
					false,
					questionComponent,
					requestForm,
					errors,
				),
			)
		}
	}
}

func handleUpdateQuiz(c echo.Context) error {
	errors := map[string]string{}
	messages := map[string]string{}

	cc := c.(*context.AppContext)

	var requestForm request.UpdateQuizRequestForm
	if err := c.Bind(&requestForm); err != nil {
		errors["other"] = "Parsing error: " + err.Error()
		return render.TemplRender(c, 200, forms.UpdateQuizForm(requestForm, errors, messages))
	}

	if requestForm.QuizId == "" {
		errors["other"] = "quizId is required"
	}
	if requestForm.Title == "" {
		errors["title"] = "Title is required"
	}
	if requestForm.Description == "" {
		errors["description"] = "Description is required"
	}
	if len(errors) > 0 {
		return render.TemplRender(c, 200, forms.UpdateQuizForm(requestForm, errors, messages))
	}

	_, err := cc.ApiService.UpdateQuiz(requestForm.QuizId, requestForm.Title, requestForm.Description)
	if err != nil {
		errors["other"] = fmt.Sprintf("Error updating %s, error: %s", requestForm.QuizId, err.Error())
		return render.TemplRender(c, 200, forms.UpdateQuizForm(requestForm, errors, messages))
	}

	messages["successful"] = fmt.Sprintf("Succesfuly updated '%s'!", requestForm.Title)
	return render.TemplRender(c, 200, forms.UpdateQuizForm(requestForm, errors, messages))
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

func handleQuizPreviewPopup(c echo.Context) error {
	cc := c.(*context.AppContext)

	quizId := c.Param("quizId")
	quiz, err := cc.ApiService.GetQuiz(quizId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "quizId is required")
	}

	// TODO
	//hasQuizSession, err := cc.ApiService.HasQuizSession(cc.Session.User.Id, quizId)
	_, err = cc.ApiService.HasQuizSession(cc.Session.User.Id, quizId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error getting checking open sessions: "+err.Error())
	}

	// TODO get all sessions and then find the open one

	props := components.PreviewQuizPopupProps{
		Quiz: quiz,
	}
	return render.TemplRender(c, 200, components.PreviewQuizPopup(props))
}
func handleClosePopup(c echo.Context) error {
	return c.NoContent(200)
}
