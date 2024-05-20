package external

import (
	"encoding/json"
)

type CreateQuizRequestBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateQuizRequestBody struct {
	Title       string `json:"name"`
	Description string `json:"description"`
}

type QuizInfo struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatorId   string `json:"creatorId"`
	CreatorName string `json:"creatorName"`
}

type Quiz struct {
	QuizInfo
	Questions []json.RawMessage `json:"questions"`
}

type QuizInfosResponse struct {
	Quizzes []QuizInfo `json:"quizzes"`
	Length  int        `json:"length"`
}
