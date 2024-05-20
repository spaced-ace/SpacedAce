package business

import "spaced-ace/models"

type Option struct {
	Value   string `json:"value"`
	Correct bool   `json:"correct"`
}

type SingleChoiceQuestion struct {
	Id           string              `json:"id"`
	QuizId       string              `json:"quizId"`
	Order        int                 `json:"order"`
	QuestionType models.QuestionType `json:"questionType"`
	Question     string              `json:"question"`
	Options      []Option            `json:"options"`
}

type MultipleChoiceQuestion struct {
	Id           string              `json:"id"`
	QuizId       string              `json:"quizId"`
	Order        int                 `json:"order"`
	QuestionType models.QuestionType `json:"questionType"`
	Question     string              `json:"question"`
	Options      []Option            `json:"options"`
}

type TrueOrFalseQuestion struct {
	Id           string              `json:"id"`
	QuizId       string              `json:"quizId"`
	Order        int                 `json:"order"`
	QuestionType models.QuestionType `json:"questionType"`
	Question     string              `json:"question"`
	Answer       bool                `json:"answer"`
}

type OpenEndedQuestion struct {
	Id           string              `json:"id"`
	QuizId       string              `json:"quizId"`
	Order        int                 `json:"order"`
	QuestionType models.QuestionType `json:"questionType"`
	Question     string              `json:"question"`
	Context      string              `json:"context"`
	Answer       string              `json:"answer"`
}
