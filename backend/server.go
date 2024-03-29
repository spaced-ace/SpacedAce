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
	public.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	public.POST("/authenticate-user", auth.AuthenticateUser)
	public.GET("/authenticated", auth.Authenticated)
	public.POST("/create-user", auth.Register)
	protected.POST("/logout", auth.Logout)

	public.GET("/quiz-infos", handlers.GetQuizInfos)

	public.POST("/multiple-choice", handlers.PostMultipleChoiceQuestion)

	e.Logger.Fatal(e.Start(":" + constants.PORT))
}
