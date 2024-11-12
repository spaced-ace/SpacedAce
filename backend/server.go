package main

import (
	"golang.org/x/net/context"
	"log"
	"spaced-ace-backend/api/handlers"
	"spaced-ace-backend/auth"
	"spaced-ace-backend/constants"
	"spaced-ace-backend/question"
	"spaced-ace-backend/quiz"
	"spaced-ace-backend/utils"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	auth.InitDb()
	quiz.InitDb()
	question.InitDb()

	// Init and close SQLC connection gracefully
	sqlcQuerier := utils.GetQuerier()
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer closeCancel()
		if err := sqlcQuerier.Close(closeCtx); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	public := e.Group("")
	public.POST("/authenticate-user", auth.AuthenticateUser)
	public.GET("/authenticated", auth.Authenticated)
	public.POST("/create-user", auth.Register)
	public.DELETE("/delete-user/:id", auth.DeleteUserEndpoint)

	protected := e.Group("")
	protected.POST("/logout", auth.Logout)

	quizGroup := protected.Group("/quizzes")
	quizGroup.GET("/:id", handlers.GetQuizEndpoint)
	quizGroup.PATCH("/:id", handlers.UpdateQuizEndpoint)
	quizGroup.DELETE("/:id", handlers.DeleteQuizEndpoint)
	quizGroup.GET("/user/:id", handlers.GetQuizzesOfUserEndpoint)
	quizGroup.POST("/create", handlers.CreateQuizEndpoint)

	questions := protected.Group("/questions")
	questions.POST("/multiple-choice", handlers.CreateMultipleChoiceQuestionEndpoint)
	questions.GET("/multiple-choice/:id", handlers.GetMultipleChoiceEndpoint)
	questions.PATCH("/multiple-choice/:id", handlers.UpdateMultipleChoiceQuestionEndpoint)
	questions.DELETE("/multiple-choice/:quizId/:id", handlers.DeleteMultipleChoiceQuestionEndpoint)

	questions.POST("/single-choice", handlers.CreateSingleChoiceQuestionEndpoint)
	questions.GET("/single-choice/:id", handlers.GetSingleChoiceEndpoint)
	questions.PATCH("/single-choice/:id", handlers.UpdateSingleChoiceQuestionEndpoint)
	questions.DELETE("/single-choice/:quizId/:id", handlers.DeleteSingleChoiceQuestionEndpoint)

	questions.POST("/true-or-false", handlers.CreateTrueOrFalseQuestionEndpoint)
	questions.GET("/true-or-false/:id", handlers.GetTrueOrFalseEndpoint)
	questions.PATCH("/true-or-false/:id", handlers.UpdateTrueOrFalseQuestionEndpoint)
	questions.DELETE("/true-or-false/:quizId/:id", handlers.DeleteTrueOrFalseQuestionEndpoint)

	quizSessions := protected.Group("/quiz-sessions")
	quizSessions.GET("/:quizSessionId", handlers.GetQuizSession)
	quizSessions.GET("", handlers.GetQuizSessions)
	quizSessions.GET("/has-open", handlers.HasOpenQuizSession)
	quizSessions.POST("/start", handlers.StartQuizSession)
	quizSessions.POST("/:quizSessionId/submit", handlers.PostSubmitQuiz)
	quizSessions.GET("/:quizSessionId/result", handlers.GetQuizResult)
	quizSessions.GET("/:quizSessionId/answers", handlers.GetAnswers)
	quizSessions.PUT("/:quizSessionId/answers", handlers.PutCreateOrUpdateAnswer)

	quizHistory := protected.Group("/quiz-history")
	quizHistory.GET("", handlers.GetQuizHistoryEntries)

	learnList := protected.Group("/learn-list")
	learnList.GET("", handlers.GetLearnList)
	learnList.POST("/:quizID/add", handlers.PostAddQuizToLearnList)
	learnList.POST("/:quizID/remove", handlers.PostRemoveQuizFromLearnList)

	reviewItem := protected.Group("/review-items")
	reviewItem.GET("", handlers.GetReviewItems)
	reviewItem.GET("/quiz-options", handlers.GetQuizOptions)
	reviewItem.GET("/item-counts", handlers.GetReviewItemCounts)
	reviewItem.GET("/get-question-and-next-item/:reviewItemID", handlers.GetQuestionAndNextReviewItem)
	reviewItem.GET("/get-question-and-next-item", handlers.GetQuestionAndNextReviewItem)

	e.Logger.Fatal(e.Start(":" + constants.PORT))
}
