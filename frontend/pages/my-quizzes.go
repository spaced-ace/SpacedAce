package pages

import (
	"github.com/labstack/echo/v4"
	"spaced-ace/context"
)

type MyQuizzesPageData struct {
	Session *context.Session
}

func MyQuizzesPage(c echo.Context) error {
	cc := c.(*context.Context)

	pageData := MyQuizzesPageData{
		Session: cc.Session,
	}

	return c.Render(200, "my-quizzes", pageData)
}
