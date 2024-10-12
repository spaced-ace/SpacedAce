package external

import (
	"spaced-ace/models"
	"spaced-ace/models/business"
)

type CreateQuizSessionRequestBody struct {
	UserID string `json:"userid"`
	QuizID string `json:"quizid"`
}

type GetQuizSessionsResponseBody struct {
	QuizSessions []QuizSession `json:"quizSessions"`
	Length       int           `json:"length"`
}

type QuizSession struct {
	ID         string              `json:"id"`
	UserID     string              `json:"userid"`
	QuizID     string              `json:"quizid"`
	StartedAt  models.NullableTime `json:"startedAt"`
	FinishedAt models.NullableTime `json:"finishedAt"`
	ClosesAt   models.NullableTime `json:"closesAt"`
}

func (qs QuizSession) MapToBusiness() (*business.QuizSession, error) {
	return &business.QuizSession{
		Id:         qs.ID,
		UserId:     qs.UserID,
		QuizId:     qs.QuizID,
		StartedAt:  qs.StartedAt.Time,
		FinishedAt: qs.FinishedAt.Time,
		ClosesAt:   qs.ClosesAt.Time,
	}, nil
}
