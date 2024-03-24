package pages

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/context"
)

type LoginPageData struct {
	Session *context.Session
}

func Login(c echo.Context) error {
	cc := c.(*context.Context)
	if cc.Session != nil {
		return c.Redirect(http.StatusFound, "/my-quizzes")
	}

	pageData := LoginPageData{
		Session: cc.Session,
	}

	return c.Render(200, "login", pageData)
}
