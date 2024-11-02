package external

import (
	"spaced-ace/models/business"
	"time"
)

type QuizHistoryEntry struct {
	QuizID          string        `json:"quizID"`
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

func (q QuizHistoryEntry) MapToBusiness() (*business.QuizHistoryEntry, error) {
	return &business.QuizHistoryEntry{
		QuizId:          q.QuizID,
		QuizName:        q.QuizName,
		SessionID:       q.SessionID,
		Finished:        q.Finished,
		DateTaken:       q.DateTaken,
		TimeSpent:       q.TimeSpent,
		ScorePercentage: q.ScorePercentage,
	}, nil
}
