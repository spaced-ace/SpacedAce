package request

type GenerateQuestionForm struct {
	QuizId       string `form:"quizId"`
	QuestionType string `form:"questionType"`
	Context      string `form:"context"`
}
