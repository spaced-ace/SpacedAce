package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"net/http"
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

type LoginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupForm struct {
	Email         string `form:"email"`
	Name          string `form:"name"`
	Password      string `form:"password"`
	PasswordAgain string `form:"passwordAgain"`
}

type SignupRequestBody struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	Password      string `json:"password"`
	PasswordAgain string `json:"passwordAgain"`
}

type MultipleChoiceResponse struct {
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	CorrectOption string   `json:"correctOption"`
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
	public.GET("/", pages.Index)
	public.GET("/login", pages.Login)
	public.GET("/signup", pages.Signup)

	// Protected pages
	protected.GET("/my-quizzes", pages.MyQuizzes)
	protected.GET("/generate", func(c echo.Context) error {
		cc := c.(*context.Context)

		pageData := struct {
			Session *context.Session
		}{
			Session: cc.Session,
		}

		return c.Render(200, "generate", pageData)
	})

	// API endpoints
	public.POST("/login", func(c echo.Context) error {
		var loginForm = LoginForm{}
		if err := c.Bind(&loginForm); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		bodyMap := LoginRequestBody{
			Email:    loginForm.Email,
			Password: loginForm.Password,
		}
		bodyBytes, err := json.Marshal(bodyMap)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		bodyBuffer := bytes.NewBuffer(bodyBytes)

		resp, err := http.Post(constants.BACKEND_URL+"/authenticate-user", "application/json", bodyBuffer)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadGateway, err.Error())
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return echo.NewHTTPError(resp.StatusCode, resp.Status)
		}

		// parse the response
		var user context.User
		err = json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		userSession := context.CreateSession(user)
		cookie := new(http.Cookie)
		cookie.Name = "session"
		cookie.Value = userSession.Id
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.Expires = userSession.Expires

		c.SetCookie(cookie)
		c.Response().Header().Set("HX-Redirect", "/my-quizzes")
		return c.String(http.StatusOK, "login successful")
	})
	public.POST("/signup", func(c echo.Context) error {
		var signupForm = SignupForm{}
		if err := c.Bind(&signupForm); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		bodyMap := SignupRequestBody{
			Name:          signupForm.Name,
			Email:         signupForm.Email,
			Password:      signupForm.Password,
			PasswordAgain: signupForm.PasswordAgain,
		}
		bodyBytes, err := json.Marshal(bodyMap)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		bodyBuffer := bytes.NewBuffer(bodyBytes)

		resp, err := http.Post(constants.BACKEND_URL+"/create-user", "application/json", bodyBuffer)
		fmt.Println(err)
		fmt.Println(resp.Status)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadGateway, err.Error())
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return echo.NewHTTPError(resp.StatusCode, resp.Status)
		}

		// parse the response
		var user context.User
		err = json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		userSession := context.CreateSession(user)
		cookie := new(http.Cookie)
		cookie.Name = "session"
		cookie.Value = userSession.Id
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.Expires = userSession.Expires

		c.SetCookie(cookie)
		c.Response().Header().Set("HX-Redirect", "/my-quizzes")
		return c.String(http.StatusCreated, "signup successful")
	})
	protected.POST("logout", func(c echo.Context) error {
		cc := c.(*context.Context)
		if cc.Session != nil {
			context.DeleteSession(cc.Session.Id)
		}

		cookie := new(http.Cookie)
		cookie.Name = "session"
		cookie.Value = ""
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.Expires = time.Now().Add(-1 * time.Hour)

		c.SetCookie(cookie)
		c.Response().Header().Set("HX-Redirect", "/")
		//return c.Render(http.StatusOK, "index", nil)
		return c.NoContent(http.StatusOK)
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
