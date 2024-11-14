package request

type SubmitReviewItemQuestionForm struct {
	SingleChoiceValue   string   `form:"single-choice-value"`
	MultipleChoiceValue []string `form:"multiple-choice-value"`
	TrueOrFalseValue    bool     `form:"true-or-false-value"`
}
