package pages

import (
	"github.com/labstack/echo/v4"
	"spaced-ace/context"
)

type IndexPageData struct {
	Session *context.Session
}

func IndexPage(c echo.Context) error {
	cc := c.(*context.Context)

	data := IndexPageData{
		Session: cc.Session,
	}

	return c.Render(200, "index", data)
}
