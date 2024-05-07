package pages

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/context"
	"spaced-ace/models"
)

type SignupPageData struct {
	Session *models.Session
}

func SignupPage(c echo.Context) error {
	cc := c.(*context.AppContext)
	if cc.Session != nil {
		return c.Redirect(http.StatusFound, "/my-quizzes")
	}

	pageData := SignupPageData{
		Session: cc.Session,
	}

	return c.Render(200, "signup", pageData)
}
