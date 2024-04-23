package models

type QuestionType int

const (
	SingleChoice QuestionType = iota
	MultipleChoice
	TrueOrFalse
	OpenEnded
)

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
	Id            string   `json:"id"`
	QuizId        string   `json:"quizid"`
	Question      string   `json:"question"`
	Answers       []string `json:"answers"`
	CorrectAnswer string   `json:"correctAnswer"`
}

func NewSingleChoiceQuestion(id string, quizId string, order int, question string, options []Option) SingleChoiceQuestion {
	return SingleChoiceQuestion{
		Id:           id,
		QuizId:       quizId,
		Order:        order,
		QuestionType: SingleChoice,
		Question:     question,
		Options:      options,
	}
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
	Id             string   `json:"id"`
	QuizId         string   `json:"quizid"`
	Question       string   `json:"question"`
	Answers        []string `json:"answers"`
	CorrectAnswers []string `json:"correctAnswers"`
}

func NewMultipleChoiceQuestion(id string, quizId string, order int, question string, options []Option) MultipleChoiceQuestion {
	return MultipleChoiceQuestion{
		Id:           id,
		QuizId:       quizId,
		Order:        order,
		QuestionType: MultipleChoice,
		Question:     question,
		Options:      options,
	}
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
	Id            string `json:"id"`
	QuizId        string `json:"quizid"`
	Question      string `json:"question"`
	CorrectAnswer bool   `json:"correct_answer"`
}

func NewTrueOrFalseQuestion(id string, quizId string, order int, question string, answer bool) TrueOrFalseQuestion {
	return TrueOrFalseQuestion{
		Id:           id,
		QuizId:       quizId,
		Order:        order,
		QuestionType: TrueOrFalse,
		Question:     question,
		Answer:       answer,
	}
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

func NewOpenEndedQuestion(id string, quizId string, order int, question string, context string, answer string) OpenEndedQuestion {
	return OpenEndedQuestion{
		Id:           id,
		QuizId:       quizId,
		Order:        order,
		QuestionType: OpenEnded,
		Question:     question,
		Context:      context,
		Answer:       answer,
	}
}
