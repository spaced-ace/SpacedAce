package request

type CreateOrUpdateAnswerForm struct {
	QuestionId string `form:"questionId"`
	AnswerType string `form:"answerType"`
}

type CreateOrUpdateSingleChoiceAnswerForm struct {
	CreateOrUpdateAnswerForm
	Answer string `form:"answer"`
}
type CreateOrUpdateMultipleChoiceAnswerForm struct {
	CreateOrUpdateAnswerForm
	Answers []string `form:"answer"`
}
type CreateOrUpdateTrueOrFalseAnswerForm struct {
	CreateOrUpdateAnswerForm
	Answer bool `form:"answer"`
}
