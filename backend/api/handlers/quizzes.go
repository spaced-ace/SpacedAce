package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace-backend/api/models"
)

type QuizzesResponseBody struct {
	Quizzes []models.QuizInfo `json:"quizzes"`
	Length  int               `json:"length"`
}

var quizzes = []models.QuizInfo{
	{
		Id:          "1",
		Title:       "Quiz 1",
		Description: "This is a quiz.",
		CreatorId:   "ccc6874f-40ab-46b1-b8e4-64ad278eaff5",
		CreatorName: "John Doe",
	},
	{
		Id:          "2",
		Title:       "Quiz 2",
		Description: "This is a quiz.",
		CreatorId:   "2c26874f-44ab-46b1-b8e4-64ad278ea55e",
		CreatorName: "Jane Doe",
	},
	{
		Id:          "3",
		Title:       "Quiz 3",
		Description: "This is a quiz.",
		CreatorId:   "2c26874f-40ab-46b1-b8e4-64ad278eaff5",
		CreatorName: "John Smith",
	},
}

func GetQuizInfos(c echo.Context) error {
	creatorId := c.QueryParam("creatorId")
	fmt.Println(creatorId)

	var filteredQuizzes []models.QuizInfo
	if creatorId != "" {
		for _, quiz := range quizzes {
			if quiz.CreatorId == creatorId {
				filteredQuizzes = append(filteredQuizzes, quiz)
			}
		}
	} else {
		filteredQuizzes = quizzes
	}

	response := QuizzesResponseBody{
		Quizzes: filteredQuizzes,
		Length:  len(filteredQuizzes),
	}

	return c.JSON(http.StatusOK, response)
}
