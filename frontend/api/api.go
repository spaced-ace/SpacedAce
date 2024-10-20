package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/auth"
	"spaced-ace/context"
)

func RegisterRoutes(e *echo.Echo) {
	public := e.Group("")

	protected := e.Group("")
	protected.Use(context.RequireSessionMiddleware)

	// Public pages
	public.GET("/", handleIndexPage)
	public.GET("/login", handleLoginPage)
	public.GET("/signup", handleSignupPage)
	public.GET("/close-popup", handleClosePopup)

	// My quizzes page
	protected.GET("/my-quizzes", handleMyQuizzesPage)
	protected.DELETE("/quizzes/:quizId", handleDeleteQuiz)

	// Take quiz page
	protected.GET("quizzes/:quizId/preview-popup", handleQuizPreviewPopup)
	protected.GET("/quizzes/:quizId/take", handleTakeQuizPage)
	protected.GET("/quizzes/:quizId/take/:quizSessionId", handleTakeQuizPage)

	// Quiz creation page
	protected.GET("/create-new-quiz", handleCreateNewQuizPage)
	protected.POST("/quizzes/create", handleCreateQuiz)

	// Question generation
	protected.GET("/quizzes/:id/edit", handleEditQuizPage)
	protected.POST("/generate/start", handleGenerateQuestionStart)
	protected.POST("/generate", handleGenerateQuestion)
	protected.PATCH("/quizzes/:id", handleUpdateQuiz)
	protected.DELETE("/questions/:questionId", handleDeleteQuestion)

	// Auth endpoints
	public.POST("/login", auth.PostLogin)
	public.POST("/signup", auth.PostRegister)
	protected.POST("/logout", func(c echo.Context) error {
		c.Response().Header().Set("HX-Redirect", "/")

		cc, ok := c.(*context.AppContext)
		if !ok {
			c.Logger().Printf("Cannot cast echo.Context to context.AppContext")
			return c.NoContent(http.StatusOK)
		}
		if cc.Session != nil {
			if err := cc.ApiService.DeleteSession(); err != nil {
				c.SetCookie(&http.Cookie{
					Name: "session",
				})
				return c.NoContent(http.StatusOK)
			}
		}

		sessionCookie, err := c.Cookie("session")
		if err == nil {
			sessionCookie.MaxAge = -1
			c.SetCookie(sessionCookie)
		}

		return c.NoContent(http.StatusOK)
	})
}
