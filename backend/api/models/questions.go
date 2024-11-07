package models

type QuestionType int

const (
	SingleChoice QuestionType = iota
	MultipleChoice
	TrueOrFalse
	OpenEnded
)

type MultipleChoiceQuestion struct {
	ID             string       `json:"id"`
	QuizID         string       `json:"quizid"`
	QuestionType   QuestionType `json:"questionType"`
	Question       string       `json:"question"`
	Answers        []string     `json:"answers"`
	CorrectAnswers []string     `json:"correctAnswers"`
}
type SingleChoiceQuestion struct {
	ID            string       `json:"id"`
	QuizID        string       `json:"quizid"`
	QuestionType  QuestionType `json:"questionType"`
	Question      string       `json:"question"`
	Answers       []string     `json:"answers"`
	CorrectAnswer string       `json:"correctAnswer"`
}
type TrueOrFalseQuestion struct {
	ID            string       `json:"id"`
	QuizID        string       `json:"quizid"`
	QuestionType  QuestionType `json:"questionType"`
	Question      string       `json:"question"`
	CorrectAnswer bool         `json:"correct_answer"`
}

type SingleChoiceUpdateRequestBody struct {
	QuizId        string   `json:"quizId"`
	Question      string   `json:"question"`
	Answers       []string `json:"answers"`
	CorrectAnswer string   `json:"correctAnswer"`
}

type MultipleChoiceUpdateRequestBody struct {
	QuizId         string   `json:"quizId"`
	Question       string   `json:"question"`
	Answers        []string `json:"answers"`
	CorrectAnswers []string `json:"correctAnswers"`
}

type TrueOrFalseUpdateRequestBody struct {
	QuizId        string `json:"quizId"`
	Question      string `json:"question"`
	CorrectAnswer bool   `json:"correctAnswer"`
}

type QuestionCreationRequestBody struct {
	QuizId string `json:"quizId"`
	Prompt string `json:"prompt"`
}
