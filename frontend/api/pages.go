package api

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"spaced-ace/constants"
	"spaced-ace/context"
	"spaced-ace/models"
	"spaced-ace/utils"
)

// TODO move to models
type QuizWithMetaData struct {
	QuizInfo              models.QuizInfo
	QuestionsWithMetaData []QuestionWithMetaData
}
type QuestionWithMetaData struct {
	EditMode bool
	Question models.Question
}
type quizResponse struct {
	models.QuizInfo
	Questions []map[string]interface{}
}
type quiz struct {
	models.QuizInfo
	FromColor string
	ToColor   string
}
type quizzesResponse struct {
	Quizzes []models.QuizInfo `json:"quizzes"`
	Length  int               `json:"length"`
}

// TODO move to api
func getQuiz(quizId string, sessionId string) (*QuizWithMetaData, error) {
	if quizId == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Quiz ID is required")
	}

	req, err := http.NewRequest("GET", constants.BACKEND_URL+"/quizzes/"+quizId, nil)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error creating request")
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: sessionId,
	})
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error reading response body: "+err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, echo.NewHTTPError(resp.StatusCode, string(body))
	}

	var response quizResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error parsing response body: "+err.Error())
	}

	questionsWithMetaData := make([]QuestionWithMetaData, len(response.Questions))
	for i, question := range response.Questions {
		questionsWithMetaData[i] = QuestionWithMetaData{
			EditMode: false,
		}

		jsonData, err := json.Marshal(question)
		if err != nil {
			continue
		}
		questionType := models.ParseFloatToQuestionType(question["questionType"].(float64))

		if questionType == models.SingleChoice {
			var singleChoiceQuestion models.SingleChoiceQuestionResponseBody
			err = json.Unmarshal(jsonData, &singleChoiceQuestion)
			if err != nil {
				fmt.Println("Error unmarshalling single choice question: ", err)
				continue
			}
			questionsWithMetaData[i].Question = models.SingleChoiceQuestion{
				Id:           singleChoiceQuestion.Id,
				QuizId:       singleChoiceQuestion.QuizId,
				QuestionType: models.SingleChoice,
				Question:     singleChoiceQuestion.Question,
				Options: []models.Option{
					{Value: singleChoiceQuestion.Answers[0], Correct: singleChoiceQuestion.CorrectAnswer == "A"},
					{Value: singleChoiceQuestion.Answers[1], Correct: singleChoiceQuestion.CorrectAnswer == "B"},
					{Value: singleChoiceQuestion.Answers[2], Correct: singleChoiceQuestion.CorrectAnswer == "C"},
					{Value: singleChoiceQuestion.Answers[3], Correct: singleChoiceQuestion.CorrectAnswer == "D"},
				},
			}
		}
		if questionType == models.MultipleChoice {
			var multipleChoiceQuestion models.MultipleChoiceQuestionResponseBody
			err = json.Unmarshal(jsonData, &multipleChoiceQuestion)
			if err != nil {
				fmt.Println("Error unmarshalling multiple choice question: ", err)
				continue
			}
			questionsWithMetaData[i].Question = models.MultipleChoiceQuestion{
				Id:           multipleChoiceQuestion.Id,
				QuizId:       multipleChoiceQuestion.QuizId,
				QuestionType: models.MultipleChoice,
				Question:     multipleChoiceQuestion.Question,
				Options: []models.Option{
					{Value: multipleChoiceQuestion.Answers[0], Correct: utils.StringInArray("A", multipleChoiceQuestion.CorrectAnswers)},
					{Value: multipleChoiceQuestion.Answers[1], Correct: utils.StringInArray("B", multipleChoiceQuestion.CorrectAnswers)},
					{Value: multipleChoiceQuestion.Answers[2], Correct: utils.StringInArray("C", multipleChoiceQuestion.CorrectAnswers)},
					{Value: multipleChoiceQuestion.Answers[3], Correct: utils.StringInArray("D", multipleChoiceQuestion.CorrectAnswers)},
				},
			}
		}
		if questionType == models.TrueOrFalse {
			var trueOrFalseQuestion models.TrueOrFalseQuestionResponseBody
			err = json.Unmarshal(jsonData, &trueOrFalseQuestion)
			if err != nil {
				fmt.Println("Error unmarshalling multiple choice question: ", err)
				continue
			}
			questionsWithMetaData[i].Question = models.TrueOrFalseQuestion{
				Id:           trueOrFalseQuestion.Id,
				QuizId:       trueOrFalseQuestion.QuizId,
				QuestionType: models.TrueOrFalse,
				Question:     trueOrFalseQuestion.Question,
				Answer:       trueOrFalseQuestion.CorrectAnswer,
			}
		}
	}

	quizWithMetaData := QuizWithMetaData{
		QuizInfo: models.QuizInfo{
			Id:          response.Id,
			Title:       response.Title,
			Description: response.Description,
			CreatorId:   response.CreatorId,
			CreatorName: response.CreatorName,
		},
		QuestionsWithMetaData: questionsWithMetaData,
	}
	return &quizWithMetaData, nil
}

type IndexPageData struct{}
type CreateNewQuizPageData struct{}
type EditQuizPageData struct {
	QuizWithMetaData QuizWithMetaData
}
type LoginPageData struct{}
type MyQuizzesPageData struct {
	Quizzes []quiz
}
type QuizPreviewPageData struct {
	Quiz *models.Quiz
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

	quizWithMetaData, err := getQuiz(quizId, cc.Session.Id)
	if err != nil {
		return c.Redirect(http.StatusFound, "/not-found")
	}

	data := NewPageTemplate(
		c.(*context.AppContext).Session,
		EditQuizPageData{
			QuizWithMetaData: *quizWithMetaData,
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

	req, _ := http.NewRequest("GET", constants.BACKEND_URL+"/quizzes/user/"+userId, nil)
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: cc.Session.Id,
	})
	client := &http.Client{}

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	var responseBody quizzesResponse
	_ = json.NewDecoder(resp.Body).Decode(&responseBody)

	var quizzes []quiz
	for _, q := range responseBody.Quizzes {
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
