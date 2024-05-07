package pages

import (
	"github.com/labstack/echo/v4"
	"spaced-ace/context"
	"spaced-ace/models"
)

type QuizPreviewPageData struct {
	Session *models.Session
	Quiz    *models.Quiz
}

func QuizPreviewPage(c echo.Context) error {
	cc := c.(*context.AppContext)

	quizId := c.Param("quizId")

	quiz, err := cc.ApiService.GetQuiz(quizId)
	if err != nil {
		return err
	}

	pageData := QuizPreviewPageData{
		Session: cc.Session,
		Quiz:    quiz,
	}

	return c.Render(200, "index", pageData)
}
