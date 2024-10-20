package business

import "spaced-ace/utils"

type QuizInfo struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatorId   string `json:"creatorId"`
	CreatorName string `json:"creatorName"`
}

type Quiz struct {
	QuizInfo
	Questions []interface{}
}
type QuestionWithMetaData struct {
	EditMode bool
	Question interface{}
}

type QuizWithMetaData struct {
	QuizInfo              QuizInfo
	QuestionsWithMetaData []QuestionWithMetaData
}

type QuizInfoWithColors struct {
	QuizInfo  QuizInfo
	FromColor string
	ToColor   string
}

func NewQuizInfoWithColors(quizInfo QuizInfo) QuizInfoWithColors {
	fromColor, toColor := utils.GenerateColors(quizInfo.Title, quizInfo.Id)

	return QuizInfoWithColors{
		QuizInfo:  quizInfo,
		FromColor: fromColor,
		ToColor:   toColor,
	}
}
