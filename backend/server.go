package main

import (
	"bytes"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	LLM_API_URL = "http://localhost:8000"
	PORT        = "9000"
	JWT_SECRET  = "secret"
)

type LoginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type MultipleChoiceAPIRequestBody struct {
	Prompt string `json:"prompt"`
}

type MultipleChoiceAPIResponse struct {
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	CorrectOption string   `json:"correct_option"`
}

type MultipleChoiceRequestBody struct {
	Prompt string `json:"prompt"`
}

type MultipleChoiceResponse struct {
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	CorrectOption string   `json:"correctOption"`
}

func createJWT(form LoginForm, expires time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": form.Email,
		"exp":   expires.Unix(),
	})
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func login(c echo.Context) error {
	var loginForm LoginForm
	err := c.Bind(&loginForm)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if loginForm.Email != "user@email.com" && loginForm.Password != "password" {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}

	// Generate the token
	expires := time.Now().Add(1 * time.Hour)
	jwtString, err := createJWT(loginForm, expires)
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	// Set the token in a cookie
	cookie := new(http.Cookie)
	cookie.Name = "auth"
	cookie.Value = jwtString
	cookie.Expires = expires
	c.SetCookie(cookie)

	return c.String(http.StatusOK, "login successful")
}

func multipleChoiceQuestion(c echo.Context) error {
	var request = MultipleChoiceRequestBody{}
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	bodyMap := MultipleChoiceAPIRequestBody{
		Prompt: request.Prompt,
	}
	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	bodyBuffer := bytes.NewBuffer(bodyBytes)

	resp, err := http.Post(LLM_API_URL+"/multiple-choice/create", "application/json", bodyBuffer)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}
	defer resp.Body.Close()

	var result = MultipleChoiceAPIResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var response = MultipleChoiceResponse{
		Question:      result.Question,
		Options:       result.Options,
		CorrectOption: result.CorrectOption,
	}
	return c.JSON(resp.StatusCode, response)
}

func main() {
	if envLLMApiURL, exists := os.LookupEnv("LLM_API_URL"); exists {
		LLM_API_URL = envLLMApiURL
	}
	if envPort, exists := os.LookupEnv("PORT"); exists {
		PORT = envPort
	}
	if envJwtSecret, exists := os.LookupEnv("JWT_SECRET"); exists {
		JWT_SECRET = envJwtSecret
	}

	e := echo.New()
	e.Use(middleware.Logger())

	authorized := e.Group("")
	authorized.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(JWT_SECRET),
		TokenLookup: "cookie:auth",
	}))
	authorized.GET("/authorized", func(c echo.Context) error {
		return c.String(http.StatusOK, "Authorized")
	})

	public := e.Group("")
	public.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	public.POST("/login", login)
	public.POST("/multiple-choice", multipleChoiceQuestion)

	e.Logger.Fatal(e.Start(":" + PORT))
}
