package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"slices"
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
		url := fmt.Sprintf("/quiz-history/%s", quizSessionId)
		c.Response().Header().Set("HX-Replace-Url", url)
		return c.Redirect(http.StatusFound, url)
	}

	quiz, err := cc.ApiService.GetQuiz(quizId)
	if err != nil {
		return err
	}

	var answerLists *business.AnswerLists
	answers, err := cc.ApiService.GetAnswers(quizSession.Id)
	if err == nil {
		answerLists = answers
	}

	viewModel := pages.TakeQuizPageViewModel{
		QuizSession: quizSession,
		Quiz:        quiz,
		AnswerLists: answerLists,
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

func handleQuizResultPage(c echo.Context) error {
	hxRequest := c.Request().Header.Get("HX-Request") == "true"
	if !hxRequest {
		return handleNonHXRequest(c)
	}

	cc := c.(*context.AppContext)

	redirectToMyQuizzes := func() error {
		url := "my-quizzes"
		c.Response().Header().Set("HX-Replace-Url", url)
		return c.Redirect(http.StatusFound, url)
	}

	quizSessionId := c.Param("quizSessionId")
	if quizSessionId == "" {
		log.Default().Print("missing quizSessionId in url params\n")
		return redirectToMyQuizzes()
	}

	quizSession, err := cc.ApiService.GetQuizSession(quizSessionId)
	if err != nil {
		log.Default().Printf("invalid quiz session ID `%s`\n", quizSessionId)
		return redirectToMyQuizzes()
	}

	quizResult, err := cc.ApiService.GetQuizResult(quizSession.Id)
	if err != nil {
		log.Default().Printf("quiz session is not finised. quiz session ID `%s`\n", quizSessionId)
		url := fmt.Sprintf("/quizzes/%s/take/%s", quizSession.QuizId, quizSession.Id)
		c.Response().Header().Set("HX-Replace-Url", url)
		return c.Redirect(http.StatusFound, url)
	}

	quiz, err := cc.ApiService.GetQuiz(quizSession.QuizId)
	if err != nil {
		log.Default().Printf("quiz not find with ID `%s`", quizSession.QuizId)
		return redirectToMyQuizzes()
	}

	var answerLists *business.AnswerLists
	answers, err := cc.ApiService.GetAnswers(quizSession.Id)
	if err == nil {
		answerLists = answers
	}

	viewModel := pages.QuizResulPageViewModel{
		QuizSession: quizSession,
		Quiz:        quiz,
		AnswerLists: answerLists,
		QuizResult:  quizResult,
	}
	return render.TemplRender(c, 200, pages.QuizResultPage(viewModel))
}

func handleQuizHistoryPage(c echo.Context) error {
	hxRequest := c.Request().Header.Get("HX-Request") == "true"
	if !hxRequest {
		return handleNonHXRequest(c)
	}

	cc := c.(*context.AppContext)
	redirectToMyQuizzes := func() error {
		url := "my-quizzes"
		c.Response().Header().Set("HX-Replace-Url", url)
		return c.Redirect(http.StatusFound, url)
	}

	quizHistoryEntries, err := cc.ApiService.GetQuizHistory(cc.Session.User.Id)
	if err != nil {
		return redirectToMyQuizzes()
	}

	// Sort the entries to show the most recent ones first
	slices.SortFunc(quizHistoryEntries, func(a, b business.QuizHistoryEntry) int {
		return b.DateTaken.Compare(a.DateTaken)
	})

	viewModel := pages.QuizHistoryPageViewModel{
		QuizHistoryEntries: quizHistoryEntries,
	}
	return render.TemplRender(c, 200, pages.QuizHistoryPage(viewModel))
}

func handleLearnPage(c echo.Context) error {
	hxRequest := c.Request().Header.Get("HX-Request") == "true"
	if !hxRequest {
		return handleNonHXRequest(c)
	}

	cc := c.(*context.AppContext)

	total, dueToReview, err := cc.ApiService.GetReviewItemCounts()
	if err != nil {
		log.Default().Print(fmt.Errorf("getting review item counts: %s\n", err))
		total = -1
		dueToReview = -1
	}

	viewModel := pages.LearnPageViewModel{
		TotalQuestions:    total,
		QuestionsToReview: dueToReview,
	}
	return render.TemplRender(c, 200, pages.LearnPage(viewModel))
}
func handleReviewPage(c echo.Context) error {
	hxRequest := c.Request().Header.Get("HX-Request") == "true"
	if !hxRequest {
		return handleNonHXRequest(c)
	}

	reviewItemID := c.Param("reviewItemID")
	if reviewItemID == "" {
		// TODO get a review item and an ID for the next one
	} else {
		// TODO get the review item for the ID
		// TODO set nextReviewItemID to ""
	}

	//cc := c.(*context.AppContext)

	viewModel := pages.QuizReviewPageViewModel{}
	return render.TemplRender(c, 200, pages.QuizReviewPage(viewModel))
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
