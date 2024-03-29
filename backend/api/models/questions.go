package models

type QuestionType int

const (
	SingleChoice QuestionType = iota
	MultipleChoice
	TrueOrFalse
	OpenEnded
)

type Question interface{}

type SingleChoiceQuestion struct {
	Id           string       `json:"id"`
	QuestionType QuestionType `json:"questionType"`
	Question     string       `json:"question"`
	Options      []string     `json:"options"`
	Answer       int          `json:"answer"`
}

func NewSingleChoiceQuestion(id string, question string, options []string, answer int) SingleChoiceQuestion {
	return SingleChoiceQuestion{
		Id:           id,
		QuestionType: SingleChoice,
		Question:     question,
		Options:      options,
		Answer:       answer,
	}
}

type MultipleChoiceQuestion struct {
	Id           string       `json:"id"`
	QuestionType QuestionType `json:"questionType"`
	Question     string       `json:"question"`
	Options      []string     `json:"options"`
	Answers      []string     `json:"answers"`
}

func NewMultipleChoiceQuestion(id string, question string, options []string, answers []string) MultipleChoiceQuestion {
	return MultipleChoiceQuestion{
		Id:           id,
		QuestionType: MultipleChoice,
		Question:     question,
		Options:      options,
		Answers:      answers,
	}
}

type TrueOrFalseQuestion struct {
	Id           string       `json:"id"`
	QuestionType QuestionType `json:"questionType"`
	Question     string       `json:"question"`
	Answer       bool         `json:"answer"`
}

func NewTrueOrFalseQuestion(id string, question string, answer bool) TrueOrFalseQuestion {
	return TrueOrFalseQuestion{
		Id:           id,
		QuestionType: TrueOrFalse,
		Question:     question,
		Answer:       answer,
	}
}

type OpenEndedQuestion struct {
	Id           string       `json:"id"`
	QuestionType QuestionType `json:"questionType"`
	Question     string       `json:"question"`
	Context      string       `json:"context"`
	Answer       string       `json:"answer"`
}

func NewOpenEndedQuestion(id string, question string, context string, answer string) OpenEndedQuestion {
	return OpenEndedQuestion{
		Id:           id,
		QuestionType: OpenEnded,
		Question:     question,
		Context:      context,
		Answer:       answer,
	}
}
