package main

import (
	"golang.org/x/net/context"
	"log"
	"net/http"
	handlers "spaced-ace-backend/api/handlers"
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
	protected := e.Group("")
	quiz := e.Group("/quizzes")
	questions := e.Group("/questions")
	quizSessions := e.Group("/quiz-sessions")
	public.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	public.POST("/authenticate-user", auth.AuthenticateUser)
	public.GET("/authenticated", auth.Authenticated)
	public.POST("/create-user", auth.Register)
	public.DELETE("/delete-user/:id", auth.DeleteUserEndpoint)
	protected.POST("/logout", auth.Logout)

	quiz.GET("/:id", handlers.GetQuizEndpoint)
	quiz.PATCH("/:id", handlers.UpdateQuizEndpoint)
	quiz.DELETE("/:id", handlers.DeleteQuizEndpoint)
	quiz.GET("/user/:id", handlers.GetQuizzesOfUserEndpoint)
	quiz.POST("/create", handlers.CreateQuizEndpoint)

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

	quizSessions.GET("/:quizSessionId", handlers.GetQuizSession)
	quizSessions.GET("", handlers.GetQuizSessions)
	quizSessions.GET("/has-open", handlers.HasOpenQuizSession)
	quizSessions.POST("/start", handlers.StartQuizSession)
	quizSessions.POST("/:quizSessionId/stop", handlers.StopQuizSession)
	quizSessions.GET("/:quizSessionId/answers", handlers.GetAnswers)
	quizSessions.PUT("/:quizSessionId/answers", handlers.PutCreateOrUpdateAnswer)

	e.Logger.Fatal(e.Start(":" + constants.PORT))
}
