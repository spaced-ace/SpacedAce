package business

type LearnListItem struct {
	QuizID   string
	QuizName string
}
type LearnList struct {
	AvailableItems []LearnListItem
	SelectedItems  []LearnListItem
}
