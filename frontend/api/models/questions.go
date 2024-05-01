package models

type QuestionType int

const (
	SingleChoice QuestionType = iota
	MultipleChoice
	TrueOrFalse
	OpenEnded
	Unknown
)

func ParseFloatToQuestionType(floatValue float64) QuestionType {
	intValue := int(floatValue)
	switch intValue {
	case 0:
		return SingleChoice
	case 1:
		return MultipleChoice
	case 2:
		return TrueOrFalse
	case 3:
		return OpenEnded
	default:
		return Unknown
	}
}

type Question interface{}

type Option struct {
	Value   string `json:"value"`
	Correct bool   `json:"correct"`
}

type SingleChoiceQuestion struct {
	Id           string       `json:"id"`
	QuizId       string       `json:"quizId"`
	Order        int          `json:"order"`
	QuestionType QuestionType `json:"questionType"`
	Question     string       `json:"question"`
	Options      []Option     `json:"options"`
}

type SingleChoiceQuestionResponseBody struct {
	Id            string       `json:"id"`
	QuizId        string       `json:"quizid"`
	QuestionType  QuestionType `json:"questionType"`
	Question      string       `json:"question"`
	Answers       []string     `json:"answers"`
	CorrectAnswer string       `json:"correctAnswer"`
}

type MultipleChoiceQuestion struct {
	Id           string       `json:"id"`
	QuizId       string       `json:"quizId"`
	Order        int          `json:"order"`
	QuestionType QuestionType `json:"questionType"`
	Question     string       `json:"question"`
	Options      []Option     `json:"options"`
}

type MultipleChoiceQuestionResponseBody struct {
	Id             string       `json:"id"`
	QuizId         string       `json:"quizid"`
	QuestionType   QuestionType `json:"questionType"`
	Question       string       `json:"question"`
	Answers        []string     `json:"answers"`
	CorrectAnswers []string     `json:"correctAnswers"`
}

type TrueOrFalseQuestion struct {
	Id           string       `json:"id"`
	QuizId       string       `json:"quizId"`
	Order        int          `json:"order"`
	QuestionType QuestionType `json:"questionType"`
	Question     string       `json:"question"`
	Answer       bool         `json:"answer"`
}

type TrueOrFalseQuestionResponseBody struct {
	Id            string       `json:"id"`
	QuizId        string       `json:"quizid"`
	QuestionType  QuestionType `json:"questionType"`
	Question      string       `json:"question"`
	CorrectAnswer bool         `json:"correct_answer"`
}

type OpenEndedQuestion struct {
	Id           string       `json:"id"`
	QuizId       string       `json:"quizId"`
	Order        int          `json:"order"`
	QuestionType QuestionType `json:"questionType"`
	Question     string       `json:"question"`
	Context      string       `json:"context"`
	Answer       string       `json:"answer"`
}
