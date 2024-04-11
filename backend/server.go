package main

import (
	"net/http"
	handlers "spaced-ace-backend/api/handlers"
	"spaced-ace-backend/auth"
	"spaced-ace-backend/constants"
	"spaced-ace-backend/question"
	"spaced-ace-backend/quiz"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	auth.InitDb()
	quiz.InitDb()
	question.InitDb()

	public := e.Group("")
	protected := e.Group("")
	quiz := e.Group("/quizzes")
	questions := e.Group("/questions")
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

	e.Logger.Fatal(e.Start(":" + constants.PORT))
}
