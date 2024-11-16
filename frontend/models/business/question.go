package business

import (
	"spaced-ace/models"
)

type QuestionOption struct {
	Value   string `json:"value"`
	Correct bool   `json:"correct"`
}

type CommonQuestionProperties struct {
	Id           string              `json:"id"`
	QuizId       string              `json:"quizId"`
	Order        int                 `json:"order"`
	QuestionType models.QuestionType `json:"questionType"`
	Question     string              `json:"question"`
}

type SingleChoiceQuestion struct {
	CommonQuestionProperties
	Options []QuestionOption `json:"options"`
}

type MultipleChoiceQuestion struct {
	CommonQuestionProperties
	Options []QuestionOption `json:"options"`
}

type TrueOrFalseQuestion struct {
	CommonQuestionProperties
	Answer bool `json:"answer"`
}

type OpenEndedQuestion struct {
	CommonQuestionProperties
	Context string `json:"context"`
	Answer  string `json:"answer"`
}
