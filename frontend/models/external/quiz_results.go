package external

import "spaced-ace/models/business"

type AnswerScore struct {
	ID                     string  `json:"id"`
	QuizResultId           string  `json:"quizResultId"`
	SingleChoiceAnswerID   string  `json:"singleChoiceAnswerId"`
	MultipleChoiceAnswerID string  `json:"multipleChoiceAnswerId"`
	TrueOrFalseAnswerID    string  `json:"trueOrFalseAnswerId"`
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

func (a *AnswerScore) MapToBusiness() (*business.AnswerScore, error) {
	return &business.AnswerScore{
		ID:                     a.ID,
		QuizResultId:           a.QuizResultId,
		SingleChoiceAnswerID:   a.SingleChoiceAnswerID,
		MultipleChoiceAnswerID: a.MultipleChoiceAnswerID,
		TrueOrFalseAnswerID:    a.TrueOrFalseAnswerID,
		MaxScore:               a.MaxScore,
		Score:                  a.Score,
	}, nil
}
func (result *QuizResult) MapToBusiness() (*business.QuizResult, error) {
	answerScores := make([]business.AnswerScore, 0, len(result.AnswerScores))
	for _, score := range result.AnswerScores {
		answerScore, err := score.MapToBusiness()
		if err != nil {
			return nil, err
		}
		answerScores = append(answerScores, *answerScore)
	}

	return &business.QuizResult{
		ID:           result.ID,
		SessionID:    result.SessionID,
		MaxScore:     result.MaxScore,
		Score:        result.Score,
		AnswerScores: answerScores,
	}, nil
}
