package pages

import "github.com/labstack/echo/v4"

func NotFoundPage(c echo.Context) error {
	return c.Render(404, "not-found.html", nil)
}
