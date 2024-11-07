package business

import (
	"time"
)

type QuizSession struct {
	Id         string
	UserId     string
	QuizId     string
	StartedAt  time.Time
	FinishedAt time.Time
	ClosesAt   time.Time
	Finished   bool
}
