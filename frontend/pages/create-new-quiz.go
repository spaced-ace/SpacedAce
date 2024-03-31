package pages

import (
	"github.com/labstack/echo/v4"
	"spaced-ace/context"
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
