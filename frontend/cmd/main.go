package main

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"net/http"
	"spaced-ace/auth"
	"spaced-ace/constants"
	"spaced-ace/context"
	"spaced-ace/pages"
	"time"
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

type SingleChoiceQuestion struct {
	Id           string   `json:"id"`
	QuestionType int      `json:"questionType"`
	Question     string   `json:"question"`
	Options      []string `json:"options"`
	Answer       int      `json:"answer"`
}

type SingleChoiceQuestionData struct {
	Order    int
	Edit     bool
	Question SingleChoiceQuestion
}

func main() {
	e := echo.New()

	e.Renderer = newTemplate()
	e.Use(middleware.Logger())
	e.Use(context.SessionMiddleware)

	// Static files
	e.Static("/static", "static")

	public := e.Group("")

	protected := e.Group("")
	protected.Use(context.RequireSessionMiddleware)

	// Public pages
	public.GET("/", pages.IndexPage)
	public.GET("/login", pages.LoginPage)
	public.GET("/signup", pages.SignupPage)

	// My quizzes page
	protected.GET("/my-quizzes", pages.MyQuizzesPage)

	// Quiz creation page
	protected.GET("/create-new-quiz", pages.CreateNewQuizPage)
	protected.POST("/quizzes/create", pages.PostCreateQuiz)

	// Question generation
	protected.GET("/quizzes/:id/edit", pages.EditQuizPage)
	protected.POST("/generate", pages.PostGenerateQuestion)
	protected.PATCH("/quizzes/:id", pages.PatchUpdateQuiz)
	protected.DELETE("/questions/:questionId", pages.DeleteQuestion)

	// Auth endpoints
	public.POST("/login", auth.PostLogin)
	public.POST("/signup", auth.PostRegister)
	protected.POST("/logout", func(c echo.Context) error {
		cc := c.(*context.Context)
		if cc.Session != nil {
			_ = context.DeleteSession(cc.Session.Id)
		}

		sessionCookie, err := c.Cookie("session")
		if err == nil {
			sessionCookie.MaxAge = -1
			c.SetCookie(sessionCookie)
		}

		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(http.StatusOK)
	})

	// Temporary question generator endpoints
	public.POST("/generate/single-choice", func(c echo.Context) error {
		data := SingleChoiceQuestionData{
			Order: 1,
			Edit:  false,
			Question: SingleChoiceQuestion{
				Id:           "1",
				QuestionType: 1,
				Question:     "Mi Franciaország fővárosa?",
				Options:      []string{"Párizs", "London", "Berlin", "Madrid"},
				Answer:       0,
			},
		}

		time.Sleep(3 * time.Second)
		return c.Render(200, "single-choice-question", data)
	})
	protected.POST("multiple-choice-question", func(c echo.Context) error {
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

		resp, err := http.Post(constants.BACKEND_URL+"/multiple-choice", "application/json", bodyBuffer)
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
		return c.Render(200, "multiple-choice-question", question)
	})

	e.Logger.Fatal(e.Start(":" + constants.PORT))
}
