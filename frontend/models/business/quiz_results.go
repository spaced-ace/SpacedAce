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
