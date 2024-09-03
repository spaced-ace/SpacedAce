package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"spaced-ace/context"
	"spaced-ace/models"
	"spaced-ace/models/business"
	"spaced-ace/utils"
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

	data := NewPageTemplate(
		cc.Session,
		SignupPageData{},
	)
	return c.Render(200, "signup", data)
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
