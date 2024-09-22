package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/context"
	"spaced-ace/models"
	"spaced-ace/models/business"
	"spaced-ace/render"
	"spaced-ace/utils"
	"spaced-ace/views/components"
	"spaced-ace/views/layout"
	"spaced-ace/views/pages"
	"strconv"
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
type TakeQuizPageData struct {
	QuizWithMetaData business.QuizWithMetaData
}
type SubmitQuizPageData struct {
	Answers []business.Answer
}

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
func handleTakeQuizPage(c echo.Context) error {
	cc := c.(*context.AppContext)

	quizId := c.Param("quizId")

	quiz, err := cc.ApiService.GetQuiz(quizId)
	if err != nil {
		return err
	}

	var questionsWithMetaData []business.QuestionWithMetaData
	for _, q := range quiz.Questions {
		questionsWithMetaData = append(
			questionsWithMetaData,
			business.QuestionWithMetaData{
				EditMode: false,
				Question: q,
			},
		)
	}

	data := NewPageTemplate(
		cc.Session,
		TakeQuizPageData{
			QuizWithMetaData: business.QuizWithMetaData{
				QuizInfo:              quiz.QuizInfo,
				QuestionsWithMetaData: questionsWithMetaData,
			},
		},
	)
	fmt.Println("handleTakeQuizPage: ", data.Data)
	return c.Render(200, "take-quiz-page", data)
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
func handleSubmitQuiz(c echo.Context) error {
	cc := c.(*context.AppContext)

	quizId := c.Param("quizId")
	if quizId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Quiz ID is required")
	}

	quiz, err := cc.ApiService.GetQuiz(quizId)
	if err != nil {
		return err
	}

	formData, err := c.FormParams()
	if err != nil {
		return err
	}

	fmt.Println("keys and values")
	for key, values := range formData {
		for _, value := range values {
			fmt.Printf("Key: %s, Value: %s\n", key, value)
		}
	}

	answers := make([]business.Answer, len(quiz.Questions))
	for i, question := range quiz.Questions {

		switch question := question.(type) {
		case *business.SingleChoiceQuestion:
			{
				answers[i] = business.Answer{
					QuestionId:   question.Id,
					QuestionText: question.Question,
					QuestionType: models.SingleChoice,
					MaxScore:     1,
					Score:        0,
					Options:      make([]business.AnswerOption, len(question.Options)),
				}

				correctOptionNum := 0
				for _, option := range question.Options {
					if option.Correct {
						correctOptionNum++
					}
				}
				correctScoreValue := float32(1) / float32(correctOptionNum)
				incorrectScoreValue := float32(1) / float32(len(question.Options)-correctOptionNum)
				fmt.Printf("correctScore: %f, incorrectScore: %f\n\n", correctScoreValue, incorrectScoreValue)

				for j, option := range question.Options {
					picked := utils.FindInFormData(formData, question.Id, string("ABCD"[j]))
					if picked {
						if option.Correct {
							answers[i].Score += correctScoreValue
						} else {
							answers[i].Score -= incorrectScoreValue
						}
					}

					answers[i].Options[j] = business.AnswerOption{
						Text:   option.Value,
						Valid:  option.Correct,
						Picked: picked,
					}
				}

				if answers[i].Score <= 0 {
					answers[i].Score = 0
				} else {
					parsedFloat, err := strconv.ParseFloat(fmt.Sprintf("%0.2f", answers[i].Score), 32)
					if err != nil {
						answers[i].Score = 0
					} else {
						answers[i].Score = float32(parsedFloat)
					}
				}
			}
		case *business.MultipleChoiceQuestion:
			{
				answers[i] = business.Answer{
					QuestionId:   question.Id,
					QuestionText: question.Question,
					QuestionType: models.MultipleChoice,
					MaxScore:     1,
					Score:        0,
					Options:      make([]business.AnswerOption, len(question.Options)),
				}

				correctOptionNum := 0
				for _, option := range question.Options {
					if option.Correct {
						correctOptionNum++
					}
				}
				correctScoreValue := float32(1) / float32(correctOptionNum)
				incorrectScoreValue := float32(1) / float32(len(question.Options)-correctOptionNum)

				for j, option := range question.Options {
					picked := utils.FindInFormData(formData, question.Id, string("ABCD"[j]))
					if picked {
						if option.Correct {
							answers[i].Score += correctScoreValue
						} else {
							answers[i].Score -= incorrectScoreValue
						}
					}

					answers[i].Options[j] = business.AnswerOption{
						Text:   option.Value,
						Valid:  option.Correct,
						Picked: picked,
					}
				}

				if answers[i].Score <= 0 {
					answers[i].Score = 0
				} else {
					parsedFloat, err := strconv.ParseFloat(fmt.Sprintf("%0.2f", answers[i].Score), 32)
					if err != nil {
						answers[i].Score = 0
					} else {
						answers[i].Score = float32(parsedFloat)
					}
				}
			}
		case *business.TrueOrFalseQuestion:
			{
				answers[i] = business.Answer{
					QuestionId:   question.Id,
					QuestionText: question.Question,
					QuestionType: models.TrueOrFalse,
					MaxScore:     1,
					Options:      make([]business.AnswerOption, 2),
				}
				pickedTrue := utils.FindInFormData(formData, question.Id, "true")
				pickedFalse := utils.FindInFormData(formData, question.Id, "false")
				answers[i].Options[0] = business.AnswerOption{
					Text:   "true",
					Valid:  question.Answer == true,
					Picked: pickedTrue,
				}
				answers[i].Options[1] = business.AnswerOption{
					Text:   "false",
					Valid:  question.Answer == false,
					Picked: pickedFalse,
				}

				if (pickedTrue && question.Answer) || (pickedFalse && !question.Answer) {
					answers[i].Score = 1
				}
			}
		default:
			panic("unhandled default case")
		}
	}

	data := NewPageTemplate(
		cc.Session,
		SubmitQuizPageData{
			Answers: answers,
		},
	)
	return c.Render(http.StatusOK, "submit-quiz-page", data)
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
