package pages

import (
	"github.com/labstack/echo/v4"
	models2 "spaced-ace/api/models"
	"spaced-ace/context"
)

type EditQuizPageData struct {
	Session          *context.Session
	QuizWithMetaData QuizWithMetaData
}

type QuizWithMetaData struct {
	QuizInfo              models2.QuizInfo
	QuestionsWithMetaData []QuestionWithMetaData
}

type QuestionWithMetaData struct {
	EditMode bool
	Question models2.Question
}

func EditQuiz(c echo.Context) error {
	cc := c.(*context.Context)

	quiz := QuizWithMetaData{
		QuizInfo: models2.QuizInfo{
			Id:          "ae664251-9ee7-4ca6-9f16-ff072de61632",
			Title:       "My QuizWithMetaData",
			Description: "This is a quiz",
			CreatorId:   "73975759-99f9-46be-b84b-cfa4d2222112",
			CreatorName: "John Doe",
		},
		QuestionsWithMetaData: []QuestionWithMetaData{
			{
				EditMode: false,
				Question: models2.NewSingleChoiceQuestion(
					"1",
					"ae664251-9ee7-4ca6-9f16-ff072de61632",
					1,
					"What is the capital of France?",
					[]models2.Option{
						{Value: "Paris", Correct: true},
						{Value: "London", Correct: false},
						{Value: "Berlin", Correct: false},
						{Value: "Madrid", Correct: false},
					}),
			},
		},
	}

	pageData := EditQuizPageData{
		Session:          cc.Session,
		QuizWithMetaData: quiz,
	}

	return c.Render(200, "edit-quiz", pageData)
}
