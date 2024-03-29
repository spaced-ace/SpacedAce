package models

type QuizInfo struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	CreatorId   string `json:"creatorId"`
	CreatorName string `json:"creatorName"`
}

type Quiz struct {
	QuizInfo
	Questions []Question
}
