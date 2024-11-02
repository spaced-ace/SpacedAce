package external

import (
	"spaced-ace/models"
	"spaced-ace/models/business"
	"strings"
)

type SingleChoiceAnswer struct {
	ID         string            `json:"id"`
	SessionID  string            `json:"sessionId"`
	QuestionID string            `json:"questionId"`
	AnswerType models.AnswerType `json:"answerType"`
	Answer     string            `json:"answer"`
}
type MultipleChoiceAnswer struct {
	ID         string            `json:"id"`
	SessionID  string            `json:"sessionId"`
	QuestionID string            `json:"questionId"`
	AnswerType models.AnswerType `json:"answerType"`
	Answers    string            `json:"answers"`
}
type TrueOrFalseAnswer struct {
	ID         string            `json:"id"`
	SessionID  string            `json:"sessionId"`
	QuestionID string            `json:"questionId"`
	AnswerType models.AnswerType `json:"answerType"`
	Answer     *bool             `json:"answer"`
}

func (a SingleChoiceAnswer) MapToBusiness() (*business.SingleChoiceAnswer, error) {
	return &business.SingleChoiceAnswer{
		CommonAnswerData: business.CommonAnswerData{
			Id:         a.ID,
			SessionId:  a.SessionID,
			QuestionId: a.QuestionID,
			AnswerType: a.AnswerType,
		},
		Answer: a.Answer,
	}, nil
}
func (a MultipleChoiceAnswer) MapToBusiness() (*business.MultipleChoiceAnswer, error) {
	return &business.MultipleChoiceAnswer{
		CommonAnswerData: business.CommonAnswerData{
			Id:         a.ID,
			SessionId:  a.SessionID,
			QuestionId: a.QuestionID,
			AnswerType: a.AnswerType,
		},
		Answers: strings.Split(a.Answers, ""),
	}, nil
}
func (a TrueOrFalseAnswer) MapToBusiness() (*business.TrueOrFalseAnswer, error) {
	return &business.TrueOrFalseAnswer{
		CommonAnswerData: business.CommonAnswerData{
			Id:         a.ID,
			SessionId:  a.SessionID,
			QuestionId: a.QuestionID,
			AnswerType: a.AnswerType,
		},
		Answer: a.Answer,
	}, nil
}

type AnswersResponse struct {
	SingleChoiceAnswers   []SingleChoiceAnswer   `json:"singleChoiceAnswers"`
	MultipleChoiceAnswers []MultipleChoiceAnswer `json:"multipleChoiceAnswers"`
	TrueOrFalseAnswer     []TrueOrFalseAnswer    `json:"trueOrFalseAnswer"`
}

func (r AnswersResponse) MapToBusiness() (*business.AnswerLists, error) {
	singleChoiceAnswers := make([]business.SingleChoiceAnswer, len(r.SingleChoiceAnswers))
	for i, a := range r.SingleChoiceAnswers {
		answer, err := a.MapToBusiness()
		if err != nil {
			return nil, err
		}
		singleChoiceAnswers[i] = *answer
	}

	multipleChoiceAnswers := make([]business.MultipleChoiceAnswer, len(r.MultipleChoiceAnswers))
	for i, a := range r.MultipleChoiceAnswers {
		answer, err := a.MapToBusiness()
		if err != nil {
			return nil, err
		}
		multipleChoiceAnswers[i] = *answer
	}

	trueOrFalseAnswers := make([]business.TrueOrFalseAnswer, len(r.TrueOrFalseAnswer))
	for i, a := range r.TrueOrFalseAnswer {
		answer, err := a.MapToBusiness()
		if err != nil {
			return nil, err
		}
		trueOrFalseAnswers[i] = *answer
	}

	return &business.AnswerLists{
		SingleChoiceAnswers:   singleChoiceAnswers,
		MultipleChoiceAnswers: multipleChoiceAnswers,
		TrueOrFalseAnswers:    trueOrFalseAnswers,
	}, nil
}

type AnswerRequestBody struct {
	QuestionId string `json:"questionId"`
	AnswerType string `json:"answerType"`
}
type SingleChoiceAnswerRequestBody struct {
	AnswerRequestBody
	Answer string `json:"answer"`
}
type MultipleChoiceAnswerRequestBody struct {
	AnswerRequestBody
	Answers []string `json:"answers"`
}
type TrueOrFalseAnswerRequestBody struct {
	AnswerRequestBody
	Answer bool `json:"answer"`
}

func NewSingleChoiceAnswerRequestBody(questionId string, answer string) *SingleChoiceAnswerRequestBody {
	return &SingleChoiceAnswerRequestBody{
		AnswerRequestBody: AnswerRequestBody{
			QuestionId: questionId,
			AnswerType: "single-choice",
		},
		Answer: answer,
	}
}
func NewMultipleChoiceAnswerRequestBody(questionId string, answers []string) *MultipleChoiceAnswerRequestBody {
	return &MultipleChoiceAnswerRequestBody{
		AnswerRequestBody: AnswerRequestBody{
			QuestionId: questionId,
			AnswerType: "multiple-choice",
		},
		Answers: answers,
	}
}
func NewTrueOrFalseAnswerRequestBody(questionId string, answer bool) *TrueOrFalseAnswerRequestBody {
	return &TrueOrFalseAnswerRequestBody{
		AnswerRequestBody: AnswerRequestBody{
			QuestionId: questionId,
			AnswerType: "true-or-false",
		},
		Answer: answer,
	}
}
