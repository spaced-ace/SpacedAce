package models

type Question interface{}

type QuizInfo struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatorId   string `json:"creatorId"`
	CreatorName string `json:"creatorName"`
}

type Quiz struct {
	QuizInfo
	Questions []Question
}
