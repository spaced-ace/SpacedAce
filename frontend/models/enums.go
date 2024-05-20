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
