package pages

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/context"
)

type SignupPageData struct {
	Session *context.Session
}

func Signup(c echo.Context) error {
	cc := c.(*context.Context)
	if cc.Session != nil {
		return c.Redirect(http.StatusFound, "/my-quizzes")
	}

	pageData := SignupPageData{
		Session: cc.Session,
	}

	return c.Render(200, "signup", pageData)
}
