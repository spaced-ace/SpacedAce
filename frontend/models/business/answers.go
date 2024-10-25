package business

import "spaced-ace/models"

type CommonAnswerData struct {
	Id         string
	SessionId  string
	QuestionId string
	AnswerType models.AnswerType
}

type SingleChoiceAnswer struct {
	CommonAnswerData
	Answer string
}

type MultipleChoiceAnswer struct {
	CommonAnswerData
	Answers []string
}

type TrueOrFalseAnswer struct {
	CommonAnswerData
	Answer bool
}

type AnswerLists struct {
	SingleChoiceAnswers   []SingleChoiceAnswer
	MultipleChoiceAnswers []MultipleChoiceAnswer
	TrueOrFalseAnswer     []TrueOrFalseAnswer
}
