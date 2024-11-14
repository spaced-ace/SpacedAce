package pages

import (
	"spaced-ace/models/business"
	"spaced-ace/models/request"
)

type CreateNewQuizPageViewModel struct {
	Values request.CreateQuizRequestForm
	Errors map[string]string
}

type LoginPageViewModel struct {
	Errors map[string]string
}

type MyQuizzesPageViewModel struct {
	QuizInfosWithColors []business.QuizInfoWithColors
}

type SignupPageViewModel struct {
	Errors map[string]string
}

type EditQuizPageViewModel struct {
	Quiz *business.Quiz
}

type TakeQuizPageViewModel struct {
	QuizSession *business.QuizSession
	Quiz        *business.Quiz
	AnswerLists *business.AnswerLists
}

type QuizResulPageViewModel struct {
	QuizSession *business.QuizSession
	Quiz        *business.Quiz
	AnswerLists *business.AnswerLists
	QuizResult  *business.QuizResult
}

type QuizHistoryPageViewModel struct {
	QuizHistoryEntries []business.QuizHistoryEntry
}

type LearnPageViewModel struct {
	TotalQuestions    int
	QuestionsToReview int
}

type QuizReviewPageViewModel struct {
	CurrentReviewItemID       string
	SingleChoiceQuestion      *business.SingleChoiceQuestion
	MultipleChoiceQuestion    *business.MultipleChoiceQuestion
	TrueOrFalseChoiceQuestion *business.TrueOrFalseQuestion
	NextReviewItemID          string
}
