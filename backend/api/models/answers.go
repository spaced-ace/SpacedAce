package models

import (
	"fmt"
	"spaced-ace-backend/db"
	"strings"
)

type AnswerType int

const (
	SingleChoiceAnswerType AnswerType = iota
	MultipleChoiceAnswerType
	TrueOrFalseAnswerType
)

type SingleChoiceAnswer struct {
	ID         string     `json:"id"`
	SessionID  string     `json:"sessionId"`
	QuestionID string     `json:"questionId"`
	AnswerType AnswerType `json:"answerType"`
	Answer     string     `json:"answer"`
}

func MapSingleChoiceAnswer(dba db.SingleChoiceAnswer) (*SingleChoiceAnswer, error) {
	if len(dba.Answer) > 1 {
		return nil, fmt.Errorf("invalid number of answers (got: %d, expected: <= 1)", len(dba.Answer))
	} else if len(dba.Answer) == 1 && dba.Answer[0] != "A" && dba.Answer[0] != "B" && dba.Answer[0] != "C" && dba.Answer[0] != "D" {
		return nil, fmt.Errorf("invalid answer option: got: `%s`, expected: A or B or C or D", dba.Answer)
	}

	return &SingleChoiceAnswer{
		ID:         dba.ID,
		SessionID:  dba.SessionID,
		QuestionID: dba.QuestionID,
		AnswerType: SingleChoiceAnswerType,
		Answer:     dba.Answer[0],
	}, nil
}

type MultipleChoiceAnswer struct {
	ID         string     `json:"id"`
	SessionID  string     `json:"sessionId"`
	QuestionID string     `json:"questionId"`
	AnswerType AnswerType `json:"answerType"`
	Answers    string     `json:"answers"`
}

func MapMultipleChoiceAnswer(dba db.MultipleChoiceAnswer) (*MultipleChoiceAnswer, error) {
	if len(dba.Answers) > 4 {
		return nil, fmt.Errorf("invalid number of answers (got: %d, expected: >= 0 and <= 4)", len(dba.Answers))
	}

	answers := ""
	for _, a := range dba.Answers {
		if a != "A" && a != "B" && a != "C" && a != "D" {
			return nil, fmt.Errorf("invalid answer option: got: `%s`, expected: A or B or C or D", a)
		}

		if strings.Contains(a, answers) {
			return nil, fmt.Errorf("the option `%s` appears more than once in %s", a, answers)
		}

		answers += a
	}

	return &MultipleChoiceAnswer{
		ID:         dba.ID,
		SessionID:  dba.SessionID,
		QuestionID: dba.QuestionID,
		AnswerType: MultipleChoiceAnswerType,
		Answers:    answers,
	}, nil
}

type TrueOrFalseAnswer struct {
	ID         string     `json:"id"`
	SessionID  string     `json:"sessionId"`
	QuestionID string     `json:"questionId"`
	AnswerType AnswerType `json:"answerType"`
	Answer     bool       `json:"answer"`
}

func MapTrueOrFalseAnswer(dba db.TrueOrFalseAnswer) (*TrueOrFalseAnswer, error) {
	return &TrueOrFalseAnswer{
		ID:         dba.ID,
		SessionID:  dba.SessionID,
		QuestionID: dba.QuestionID,
		AnswerType: TrueOrFalseAnswerType,
		Answer:     dba.Answer.Bool,
	}, nil
}
