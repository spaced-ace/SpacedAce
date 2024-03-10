package main

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"net/http"
	"os"
)

var (
	BACKEND_URL = "http://localhost:9000"
	PORT        = "42069"
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
	Name   string
	Result string
}

type Question struct {
	Question string
	Option1  string
	Option2  string
	Option3  string
	Option4  string
}

type MultipleChoiceResponse struct {
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	CorrectOption string   `json:"correctOption"`
}

func main() {
	if envBackendURL, exists := os.LookupEnv("BACKEND_URL"); exists {
		BACKEND_URL = envBackendURL
	}
	if envPort, exists := os.LookupEnv("PORT"); exists {
		PORT = envPort
	}

	e := echo.New()

	e.Renderer = newTemplate()
	e.Use(middleware.Logger())

	data := Data{Name: "HTMX", Result: "Hello World!"}

	// Static files
	e.Static("/static", "static")

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index.html", data)
	})

	e.POST("multiple-choice", func(c echo.Context) error {
		prompt := c.FormValue("prompt")
		if prompt == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "prompt is required")
		}

		bodyMap := map[string]interface{}{
			"prompt": prompt,
		}
		bodyBytes, err := json.Marshal(bodyMap)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		bodyBuffer := bytes.NewBuffer(bodyBytes)

		resp, err := http.Post(BACKEND_URL+"/multiple-choice", "application/json", bodyBuffer)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadGateway, err.Error())
		}
		defer resp.Body.Close()

		var result = MultipleChoiceResponse{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		question := Question{Question: result.Question, Option1: result.Options[0], Option2: result.Options[1], Option3: result.Options[2], Option4: result.Options[3]}
		return c.Render(200, "question.html", question)
	})

	e.Logger.Fatal(e.Start(":" + PORT))
}
