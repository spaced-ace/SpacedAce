package pages

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/api/models"
	"spaced-ace/constants"
	"spaced-ace/context"
)

type CreateNewQuizPageData struct {
	Session *context.Session
}

func CreateNewQuizPage(c echo.Context) error {
	cc := c.(*context.Context)

	pageData := MyQuizzesPageData{
		Session: cc.Session,
	}

	return c.Render(200, "create-new-quiz", pageData)
}

type QuizRequestBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateQuizRequestForm struct {
	Title       string `form:"title"`
	Description string `form:"description"`
}

func PostCreateQuiz(c echo.Context) error {
	cc := c.(*context.Context)

	var requestForm CreateQuizRequestForm
	if err := c.Bind(&requestForm); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	if requestForm.Title == "" {
		return c.String(http.StatusBadRequest, "Title is required")
	}
	if requestForm.Description == "" {
		return c.String(http.StatusBadRequest, "Description is required")
	}

	requestBody, _ := json.Marshal(QuizRequestBody{
		Name:        requestForm.Title,
		Description: requestForm.Description,
	})

	req, _ := http.NewRequest("POST", constants.BACKEND_URL+"/quizzes/create", bytes.NewBuffer(requestBody))
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: cc.Session.Id,
	})
	client := &http.Client{}

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	var responseBody models.QuizInfo
	_ = json.NewDecoder(resp.Body).Decode(&responseBody)

	c.Response().Header().Set("HX-Redirect", "/quizzes/"+responseBody.Id+"/edit")
	return c.NoContent(http.StatusCreated)
}
