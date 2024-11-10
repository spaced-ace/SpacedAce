package external

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"spaced-ace/models"
	"spaced-ace/models/business"
	"spaced-ace/utils"
)

type GenerateQuestionRequestBody struct {
	QuizId string `json:"quizId"`
	Prompt string `json:"prompt"`
}

type SingleChoiceQuestionResponseBody struct {
	Id            string              `json:"id"`
	QuizId        string              `json:"quizid"`
	QuestionType  models.QuestionType `json:"questionType"`
	Question      string              `json:"question"`
	Answers       []string            `json:"answers"`
	CorrectAnswer string              `json:"correctAnswer"`
}

func (q SingleChoiceQuestionResponseBody) MapToBusiness() (*business.SingleChoiceQuestion, error) {
	if len(q.Answers) != 4 {
		return nil, echo.NewHTTPError(500, fmt.Sprintf("Invalid number of possible answers. Expected: %d, got: %d", 4, len(q.Answers)))
	}
	if !utils.StringInArray(q.CorrectAnswer, []string{"A", "B", "C", "D"}) {
		return nil, echo.NewHTTPError(500, fmt.Sprintf("Invalid answer. Expected on of [A, B, C, D], got: %s", q.CorrectAnswer))
	}

	return &business.SingleChoiceQuestion{
		CommonQuestionProperties: business.CommonQuestionProperties{
			Id:           q.Id,
			QuizId:       q.QuizId,
			Order:        0,
			QuestionType: models.SingleChoice,
			Question:     q.Question,
		},
		Options: []business.QuestionOption{
			{Value: q.Answers[0], Correct: q.CorrectAnswer == "A"},
			{Value: q.Answers[1], Correct: q.CorrectAnswer == "B"},
			{Value: q.Answers[2], Correct: q.CorrectAnswer == "C"},
			{Value: q.Answers[3], Correct: q.CorrectAnswer == "D"},
		},
	}, nil
}

type MultipleChoiceQuestionResponseBody struct {
	Id             string              `json:"id"`
	QuizId         string              `json:"quizid"`
	QuestionType   models.QuestionType `json:"questionType"`
	Question       string              `json:"question"`
	Answers        []string            `json:"answers"`
	CorrectAnswers []string            `json:"correctAnswers"`
}

func (q MultipleChoiceQuestionResponseBody) MapToBusiness() (*business.MultipleChoiceQuestion, error) {
	if len(q.Answers) != 4 {
		return nil, echo.NewHTTPError(500, fmt.Sprintf("Invalid number of possible answers. Expected: %d, got: %d", 4, len(q.Answers)))
	}
	if len(q.CorrectAnswers) == 0 {
		return nil, echo.NewHTTPError(500, fmt.Sprintf("There are no correct answers: %+v", q.CorrectAnswers))
	}

	return &business.MultipleChoiceQuestion{
		CommonQuestionProperties: business.CommonQuestionProperties{
			Id:           q.Id,
			QuizId:       q.QuizId,
			Order:        0,
			QuestionType: models.MultipleChoice,
			Question:     q.Question,
		},
		Options: []business.QuestionOption{
			{Value: q.Answers[0], Correct: utils.StringInArray("A", q.CorrectAnswers)},
			{Value: q.Answers[1], Correct: utils.StringInArray("B", q.CorrectAnswers)},
			{Value: q.Answers[2], Correct: utils.StringInArray("C", q.CorrectAnswers)},
			{Value: q.Answers[3], Correct: utils.StringInArray("D", q.CorrectAnswers)},
		},
	}, nil
}

type TrueOrFalseQuestionResponseBody struct {
	Id            string              `json:"id"`
	QuizId        string              `json:"quizid"`
	QuestionType  models.QuestionType `json:"questionType"`
	Question      string              `json:"question"`
	CorrectAnswer bool                `json:"correct_answer"`
}

func (q TrueOrFalseQuestionResponseBody) MapToBusiness() (*business.TrueOrFalseQuestion, error) {
	return &business.TrueOrFalseQuestion{
		CommonQuestionProperties: business.CommonQuestionProperties{
			Id:           q.Id,
			QuizId:       q.QuizId,
			Order:        0,
			QuestionType: models.TrueOrFalse,
			Question:     q.Question,
		},
		Answer: q.CorrectAnswer,
	}, nil
}
