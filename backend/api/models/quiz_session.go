package models

import (
	"spaced-ace-backend/db"
)

type QuizSession struct {
	ID         string       `json:"id"`
	UserID     string       `json:"userid"`
	QuizID     string       `json:"quizid"`
	StartedAt  NullableTime `json:"startedAt"`
	FinishedAt NullableTime `json:"finishedAt"`
	ClosesAt   NullableTime `json:"closesAt"`
}

func MapQuizSession(dbo db.QuizSession) (*QuizSession, error) {
	return &QuizSession{
		ID:         dbo.ID,
		UserID:     dbo.UserID,
		QuizID:     dbo.QuizID,
		StartedAt:  newNullableTime(dbo.StartedAt.Time),
		FinishedAt: newNullableTime(dbo.FinishedAt.Time),
		ClosesAt:   newNullableTime(dbo.ClosesAt.Time),
	}, nil
}
