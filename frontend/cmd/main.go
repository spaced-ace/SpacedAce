package main

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"net/http"
	"spaced-ace/api"
	"spaced-ace/constants"
	"spaced-ace/context"
	"strings"
)

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	fmt.Printf("Rendering `%s`, data: %+v\n", name, data)
	return t.tmpl.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	e.Renderer = newTemplate()
	e.Use(middleware.Logger())
	e.Use(context.SessionMiddleware)

	// Static files
	e.Static("/static", "static")

	api.RegisterRoutes(e)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := "internal server error"
		var he *echo.HTTPError
		if errors.As(err, &he) {
			code = he.Code
			message = extractErrorMessage(he)
		}

		err = c.Render(code, "error-message", map[string]string{
			"Message": message,
		})
		if err != nil {
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
