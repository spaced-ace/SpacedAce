package business

import "spaced-ace/models"

type AnswerOption struct {
	Text   string
	Valid  bool
	Picked bool
}

type Answer struct {
	QuestionId   string
	QuestionText string
	QuestionType models.QuestionType
	Score        float32
	MaxScore     float32
	Options      []AnswerOption
}
