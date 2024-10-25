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
	TrueOrFalseAnswers    []TrueOrFalseAnswer
}

func (lists *AnswerLists) GetSingleChoiceAnswerOrNil(questionId string) *SingleChoiceAnswer {
	for _, a := range lists.SingleChoiceAnswers {
		if a.QuestionId == questionId {
			return &a
		}
	}
	return nil
}
func (lists *AnswerLists) GetMultipleChoiceAnswerOrNil(questionId string) *MultipleChoiceAnswer {
	for _, a := range lists.MultipleChoiceAnswers {
		if a.QuestionId == questionId {
			return &a
		}
	}
	return nil
}
func (lists *AnswerLists) GetTrueOrFalseAnswerOrNil(questionId string) *TrueOrFalseAnswer {
	for _, a := range lists.TrueOrFalseAnswers {
		if a.QuestionId == questionId {
			return &a
		}
	}
	return nil
}
