package main

import (
	"errors"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"spaced-ace/api"
	"spaced-ace/constants"
	"spaced-ace/views/pages"
	"strings"
)

func main() {
	e := echo.New()

	// Static files
	e.Static("/static", "static")

	e.Use(middleware.Logger())
	//e.Use(context.SessionMiddleware)
	//e.Renderer = api.NewTemplate()

	e.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, pages.IndexPage())
	})

	//api.RegisterRoutes(e)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := "internal server error"
		var he *echo.HTTPError
		if errors.As(err, &he) {
			code = he.Code
			message = extractErrorMessage(he)
		}

		data := api.NewPageTemplate(
			nil,
			map[string]string{
				"Message": message,
			},
		)

		c.Response().Header().Set("HX-Push-Url", "false")
		if err = c.Render(code, "error-message", data); err != nil {
			c.Error(err)
		}
	}

	e.Logger.Fatal(e.Start(":" + constants.PORT))
}

func extractErrorMessage(he *echo.HTTPError) string {
	parts := strings.Split(he.Error(), ", ")
	for _, part := range parts {
		if strings.HasPrefix(part, "message=") {
			return strings.TrimPrefix(part, "message=")
		}
	}
	return "unknown error"
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
