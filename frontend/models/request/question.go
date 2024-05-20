package request

type GenerateQuestionForm struct {
	QuizId  string `form:"quizId"`
	Context string `form:"context"`
}
