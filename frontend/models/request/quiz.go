package request

type CreateQuizRequestForm struct {
	Title       string `form:"title"`
	Description string `form:"description"`
}

type UpdateQuizRequestForm struct {
	QuizId      string `form:"quizId"`
	Title       string `form:"title"`
	Description string `form:"description"`
}
