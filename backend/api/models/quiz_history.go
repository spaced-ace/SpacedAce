package models

import "time"

type QuizHistoryEntry struct {
	QuizName        string        `json:"quizName"`
	SessionID       string        `json:"sessionID"`
	Finished        bool          `json:"finished"`
	DateTaken       time.Time     `json:"dateTaken"`
	TimeSpent       time.Duration `json:"timeSpent"`
	ScorePercentage float64       `json:"scorePercentage"`
}

type QuizHistoryEntriesResponseBody struct {
	QuizHistoryEntries []QuizHistoryEntry `json:"quizHistoryEntries"`
	Length             int                `json:"length"`
}
