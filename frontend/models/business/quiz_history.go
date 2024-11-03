package business

import "time"

type QuizHistoryEntry struct {
	QuizId          string
	QuizName        string
	SessionID       string
	Finished        bool
	DateTaken       time.Time
	TimeSpent       time.Duration
	ScorePercentage float64
}
