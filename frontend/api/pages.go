package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/context"
	"spaced-ace/models/business"
	"spaced-ace/render"
	"spaced-ace/views/components"
	"spaced-ace/views/layout"
	"spaced-ace/views/pages"
)

func handleIndexPage(c echo.Context) error {
	return render.TemplRender(c, 200, pages.IndexPage())
}
func handleCreateNewQuizPage(c echo.Context) error {
	hxRequest := c.Request().Header.Get("HX-Request") == "true"
	if !hxRequest {
		return handleNonHXRequest(c)
	}

	return render.TemplRender(c, 200, pages.CreateNewQuizPage(pages.CreateNewQuizPageViewModel{}))
}
func handleEditQuizPage(c echo.Context) error {
	hxRequest := c.Request().Header.Get("HX-Request") == "true"
	if !hxRequest {
		return handleNonHXRequest(c)
	}

	cc := c.(*context.AppContext)
	quizId := c.Param("id")

	quiz, err := cc.ApiService.GetQuiz(quizId)
	if err != nil {
		return c.Redirect(http.StatusFound, "/not-found")
	}

	viewModel := pages.EditQuizPageViewModel{
		Quiz: quiz,
	}
	return render.TemplRender(c, 200, pages.EditQuizPage(viewModel))
}
func handleLoginPage(c echo.Context) error {
	cc := c.(*context.AppContext)
	if cc.Session != nil {
		return c.Redirect(http.StatusFound, "/my-quizzes")
	}

	viewModel := pages.LoginPageViewModel{
		Errors: map[string]string{},
	}
	return render.TemplRender(c, 200, pages.LoginPage(viewModel))
}
func handleMyQuizzesPage(c echo.Context) error {
	hxRequest := c.Request().Header.Get("HX-Request") == "true"
	if !hxRequest {
		return handleNonHXRequest(c)
	}

	cc := c.(*context.AppContext)
	userId := cc.Session.User.Id

	quizInfos, err := cc.ApiService.GetQuizzesInfos(userId)
	if err != nil {
		return err
	}

	var quizInfosWithColors []business.QuizInfoWithColors
	for _, q := range quizInfos {
		quizInfosWithColors = append(
			quizInfosWithColors,
			business.NewQuizInfoWithColors(q),
		)
	}

	viewModel := pages.MyQuizzesPageViewModel{
		QuizInfosWithColors: quizInfosWithColors,
	}
	return render.TemplRender(c, 200, pages.MyQuizzesPage(viewModel))
}
func handleTakeQuizPage(c echo.Context) error {
	hxRequest := c.Request().Header.Get("HX-Request") == "true"
	if !hxRequest {
		return handleNonHXRequest(c)
	}

	cc := c.(*context.AppContext)

	quizId := c.Param("quizId")

	quizSessionId := c.Param("quizSessionId")
	if quizSessionId == "" {
		quizSession, err := cc.ApiService.CreateQuizSession(cc.Session.User.Id, quizId)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("/quizzes/%s/take/%s", quizId, quizSession.Id)
		c.Response().Header().Set("HX-Replace-Url", url)
		return c.Redirect(http.StatusFound, url)
	}

	quizSession, err := cc.ApiService.GetQuizSession(quizSessionId)
	if err != nil {
		return err
	}
	if quizSession.Finished {
		// TODO thus the quiz is finished, redirect the user to the results page
	}

	quiz, err := cc.ApiService.GetQuiz(quizId)
	if err != nil {
		return err
	}

	viewModel := pages.TakeQuizPageViewModel{
		QuizSession: quizSession,
		Quiz:        quiz,
	}
	return render.TemplRender(c, 200, pages.TakeQuizPage(viewModel))
}
func handleSignupPage(c echo.Context) error {
	cc := c.(*context.AppContext)
	if cc.Session != nil {
		return c.Redirect(http.StatusFound, "/my-quizzes")
	}

	viewModel := pages.SignupPageViewModel{
		Errors: map[string]string{},
	}
	return render.TemplRender(c, 200, pages.SignupPage(viewModel))
}

func handleNonHXRequest(c echo.Context) error {
	activeUrl := c.Request().URL.Path
	sideBarProps, err := createSideBarProps(c, activeUrl)
	if err != nil {
		return err
	}

	props := layout.AuthenticatedLayoutProps{
		SideBarProps: *sideBarProps,
	}
	return render.TemplRender(c, 200, layout.AuthenticatedLayout(props))
}

func createSideBarProps(c echo.Context, activeUrl string) (*components.SidebarProps, error) {
	cc := c.(*context.AppContext)
	userId := cc.Session.User.Id

	quizInfos, err := cc.ApiService.GetQuizzesInfos(userId)
	if err != nil {
		return nil, err
	}

	return &components.SidebarProps{
		Username:  cc.Session.User.Name,
		ActiveUrl: activeUrl,
		QuizInfos: quizInfos,
	}, nil
}
