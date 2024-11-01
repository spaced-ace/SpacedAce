package business

type AnswerScore struct {
	ID                     string
	QuizResultId           string
	SingleChoiceAnswerID   string
	MultipleChoiceAnswerID string
	TrueOrFalseAnswerID    string
	MaxScore               float64
	Score                  float64
}
type QuizResult struct {
	ID           string
	SessionID    string
	MaxScore     float64
	Score        float64
	AnswerScores []AnswerScore
}

func (r *QuizResult) GetAnswerScoreOrNilForSingleChoiceAnswer(answer *SingleChoiceAnswer) *AnswerScore {
	if r == nil || answer == nil {
		return nil
	}

	for _, score := range r.AnswerScores {
		if score.SingleChoiceAnswerID == answer.Id {
			return &score
		}
	}
	return nil
}
func (r *QuizResult) GetAnswerScoreOrNilForMultipleChoiceAnswer(answer *MultipleChoiceAnswer) *AnswerScore {
	if r == nil || answer == nil {
		return nil
	}

	for _, score := range r.AnswerScores {
		if score.MultipleChoiceAnswerID == answer.Id {
			return &score
		}
	}
	return nil
}
func (r *QuizResult) GetAnswerScoreOrNilForTrueOrFalseAnswer(answer *TrueOrFalseAnswer) *AnswerScore {
	if r == nil || answer == nil {
		return nil
	}

	for _, score := range r.AnswerScores {
		if score.TrueOrFalseAnswerID == answer.Id {
			return &score
		}
	}
	return nil
}
