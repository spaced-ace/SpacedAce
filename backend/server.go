package main

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"spaced-ace-backend/auth"
	"spaced-ace-backend/constants"
)

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

	resp, err := http.Post(constants.LLM_API_URL+"/multiple-choice/create", "application/json", bodyBuffer)
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
	e := echo.New()
	e.Use(middleware.Logger())

	public := e.Group("")
	public.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	public.POST("/authenticate-user", auth.AuthenticateUser)
	public.POST("/create-user", auth.CreateUser)
	public.POST("/multiple-choice", multipleChoiceQuestion)

	e.Logger.Fatal(e.Start(":" + constants.PORT))
}
