package pages

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/api/models"
	"spaced-ace/context"
	"time"
)

type CreateNewQuizPageData struct {
	Session *context.Session
}

func CreateNewQuiz(c echo.Context) error {
	cc := c.(*context.Context)

	pageData := MyQuizzesPageData{
		Session: cc.Session,
	}

	return c.Render(200, "create-new-quiz", pageData)
}

type CreateQuizRequestForm struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	Context     string `form:"context"`
}

func PostCreateQuiz(c echo.Context) error {
	cc := c.(*context.Context)

	questionType := cc.QueryParam("type")
	if questionType != "single-choice" && questionType != "multiple-choice" && questionType != "true-or-false" && questionType != "open-ended" {
		return c.String(http.StatusBadRequest, "Invalid question type")
	}

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
	if requestForm.Context == "" {
		return c.String(http.StatusBadRequest, "Context is required")
	}

	createdQuiz := models.Quiz{
		QuizInfo: models.QuizInfo{
			Id: "ae664251-9ee7-4ca6-9f16-ff072de61632",
		},
	}

	time.Sleep(3 * time.Second)

	c.Response().Header().Set("HX-Redirect", "/quizzes/"+createdQuiz.Id+"/edit")
	return c.NoContent(http.StatusCreated)
}
