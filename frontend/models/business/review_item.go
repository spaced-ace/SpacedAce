package business

import "time"

type ReviewItem struct {
	ID           string
	QuizName     string
	QuestionName string
	QuestionID   string
	NextReview   time.Time
	Difficulty   float64
	Streak       int
	NeedToReview bool
}
