package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"spaced-ace-backend/api/handlers"
	"spaced-ace-backend/auth"
	"spaced-ace-backend/constants"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	public := e.Group("")
	protected := e.Group("")
	quiz := e.Group("/quizzes")
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

	public.POST("/multiple-choice", handlers.PostMultipleChoiceQuestion)

	e.Logger.Fatal(e.Start(":" + constants.PORT))
}
