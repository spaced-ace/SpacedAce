package pages

import (
	"github.com/labstack/echo/v4"
	"spaced-ace/context"
	"spaced-ace/models"
)

type IndexPageData struct {
	Session *models.Session
}

func IndexPage(c echo.Context) error {
	cc := c.(*context.AppContext)

	data := IndexPageData{
		Session: cc.Session,
	}

	return c.Render(200, "index", data)
}
