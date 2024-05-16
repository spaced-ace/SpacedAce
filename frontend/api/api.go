package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/auth"
	"spaced-ace/context"
	"spaced-ace/pages"
)

func RegisterRoutes(e *echo.Echo) {
	public := e.Group("")

	protected := e.Group("")
	protected.Use(context.RequireSessionMiddleware)

	// Public pages
	public.GET("/", pages.IndexPage)
	public.GET("/login", pages.LoginPage)
	public.GET("/signup", pages.SignupPage)

	// My quizzes page
	protected.GET("/my-quizzes", pages.MyQuizzesPage)
	protected.DELETE("/quizzes/:quizId", pages.DeleteQuiz)

	// Take quiz page
	protected.GET("/quizzes/:quizId/preview", pages.QuizPreviewPage)
	//protected.GET("/quizzes/:quizId/take", pages.TakeQuizPage)

	// Quiz creation page
	protected.GET("/create-new-quiz", pages.CreateNewQuizPage)
	protected.POST("/quizzes/create", pages.PostCreateQuiz)

	// Question generation
	protected.GET("/quizzes/:id/edit", pages.EditQuizPage)
	protected.POST("/generate", pages.PostGenerateQuestion)
	protected.PATCH("/quizzes/:id", pages.PatchUpdateQuiz)
	protected.DELETE("/questions/:questionId", pages.DeleteQuestion)

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

	public.GET("/not-found", pages.NotFoundPage)
}
