package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/context"
	"spaced-ace/models/business"
	"spaced-ace/utils"
)

type quiz struct {
	business.QuizInfo
	FromColor string
	ToColor   string
}

type IndexPageData struct{}
type CreateNewQuizPageData struct{}
type EditQuizPageData struct {
	QuizWithMetaData business.QuizWithMetaData
}
type LoginPageData struct{}
type MyQuizzesPageData struct {
	Quizzes []quiz
}
type QuizPreviewPageData struct {
	Quiz *business.Quiz
}
type SignupPageData struct{}

func handleIndexPage(c echo.Context) error {
	data := NewPageTemplate(
		c.(*context.AppContext).Session,
		IndexPageData{},
	)
	return c.Render(200, "index", data)
}
func handleCreateNewQuizPage(c echo.Context) error {
	data := NewPageTemplate(
		c.(*context.AppContext).Session,
		CreateNewQuizPageData{},
	)
	return c.Render(200, "create-new-quiz", data)
}
func handleEditQuizPage(c echo.Context) error {
	cc := c.(*context.AppContext)
	quizId := c.Param("id")

	quiz, err := cc.ApiService.GetQuiz(quizId)
	if err != nil {
		return c.Redirect(http.StatusFound, "/not-found")
	}

	var questionsWithMetaData []business.QuestionWithMetaData
	for _, q := range quiz.Questions {
		questionsWithMetaData = append(
			questionsWithMetaData,
			business.QuestionWithMetaData{
				EditMode: true,
				Question: q,
			},
		)
	}

	data := NewPageTemplate(
		c.(*context.AppContext).Session,
		EditQuizPageData{
			QuizWithMetaData: business.QuizWithMetaData{
				QuizInfo:              quiz.QuizInfo,
				QuestionsWithMetaData: questionsWithMetaData,
			},
		},
	)
	return c.Render(200, "edit-quiz", data)
}
func handleLoginPage(c echo.Context) error {
	cc := c.(*context.AppContext)
	if cc.Session != nil {
		return c.Redirect(http.StatusFound, "/my-quizzes")
	}

	data := NewPageTemplate(
		cc.Session,
		LoginPageData{},
	)
	return c.Render(200, "login", data)
}
func handleMyQuizzesPage(c echo.Context) error {
	cc := c.(*context.AppContext)
	userId := cc.Session.User.Id

	quizInfos, err := cc.ApiService.GetQuizzesInfos(userId)
	if err != nil {
		return err
	}

	var quizzes []quiz
	for _, q := range quizInfos {
		fromColor, toColor := utils.GenerateColors(q.Title, q.Id)

		quizzes = append(quizzes, quiz{
			QuizInfo:  q,
			FromColor: fromColor,
			ToColor:   toColor,
		})
	}

	data := NewPageTemplate(
		cc.Session,
		MyQuizzesPageData{
			Quizzes: quizzes,
		},
	)
	return c.Render(200, "my-quizzes", data)
}
func handleNotFoundPage(c echo.Context) error {
	data := NewPageTemplate(
		c.(*context.AppContext).Session,
		nil,
	)
	return c.Render(404, "not-found.html", data)
}
func handleQuizPreviewPage(c echo.Context) error {
	cc := c.(*context.AppContext)

	quizId := c.Param("quizId")

	quiz, err := cc.ApiService.GetQuiz(quizId)
	if err != nil {
		return err
	}
	fmt.Println("quiz ", quiz)

	data := NewPageTemplate(
		cc.Session,
		nil,
	)
	return c.Render(200, "index", data)
}
func handleSignupPage(c echo.Context) error {
	cc := c.(*context.AppContext)
	if cc.Session != nil {
		return c.Redirect(http.StatusFound, "/my-quizzes")
	}

	data := NewPageTemplate(
		cc.Session,
		SignupPageData{},
	)
	return c.Render(200, "signup", data)
}
