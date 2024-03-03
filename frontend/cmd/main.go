package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
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
	return t.tmpl.ExecuteTemplate(w, name, data)
}

type Data struct {
	Name string
}

func main() {
	e := echo.New()

	e.Renderer = newTemplate()
	e.Use(middleware.Logger())

	data := Data{Name: "HTMX"}

	// Static files
	e.Static("/static", "static")

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index.html", data)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
