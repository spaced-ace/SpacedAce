package models

type LearnListItem struct {
	QuizID   string `json:"quizID"`
	QuizName string `json:"quizName"`
}
type LearnListResponseBody struct {
	AvailableItems []LearnListItem `json:"availableItems"`
	SelectedItems  []LearnListItem `json:"selectedItems"`
}
