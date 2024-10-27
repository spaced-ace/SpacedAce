package models

import (
	"spaced-ace-backend/db"
)

type AnswerScore struct {
	ID                     string  `json:"id"`
	QuizResultId           string  `json:"quizResultId"`
	SingleChoiceAnswerId   string  `json:"singleChoiceAnswerId"`
	MultipleChoiceAnswerId string  `json:"multipleChoiceAnswerId"`
	TrueOrFalseAnswerId    string  `json:"trueOrFalseAnswerId"`
	MaxScore               float64 `json:"maxScore"`
	Score                  float64 `json:"score"`
}
type QuizResult struct {
	ID           string        `json:"id"`
	SessionID    string        `json:"sessionId"`
	MaxScore     float64       `json:"maxScore"`
	Score        float64       `json:"score"`
	AnswerScores []AnswerScore `json:"answerScores"`
}

func MapAnswerScore(score *db.AnswerScore) (*AnswerScore, error) {
	return &AnswerScore{
		ID:                     score.ID,
		QuizResultId:           score.QuizResultID,
		SingleChoiceAnswerId:   score.SingleChoiceAnswerID,
		MultipleChoiceAnswerId: score.MultipleChoiceAnswerID,
		TrueOrFalseAnswerId:    score.TrueOrFalseAnswerID,
		MaxScore:               score.MaxScore,
		Score:                  score.Score,
	}, nil
}
func MapQuizResult(result *db.QuizResult) (*QuizResult, error) {
	return &QuizResult{
		ID:           result.ID,
		SessionID:    result.SessionID,
		MaxScore:     result.MaxScore,
		Score:        result.Score,
		AnswerScores: []AnswerScore{},
	}, nil
}
