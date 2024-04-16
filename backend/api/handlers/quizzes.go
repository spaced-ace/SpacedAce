package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	models "spaced-ace-backend/api/models"
	"spaced-ace-backend/auth"
	"spaced-ace-backend/question"
	quiz "spaced-ace-backend/quiz"
)

type QuizzesResponse struct {
	Quizzes []models.QuizInfo `json:"quizzes"`
	Length  int               `json:"length"`
}

type QuizRequestBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Quiz struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	CreatorId   string `json:"creatorid"`
	Description string `json:"description"`
}

func init() {
}

func CreateQuizEndpoint(c echo.Context) error {
	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	uid, err := auth.GetUserIdBySession(session.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	user, err := auth.GetUserById(uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	var request = QuizRequestBody{}
	if err = json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}
	createdQuiz, err := quiz.CreateQuiz(uid, request.Name, request.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, models.QuizInfo{Id: createdQuiz.Id, Title: createdQuiz.Name, Description: createdQuiz.Description.String, CreatorName: user.Name, CreatorId: user.Id})
}

func GetQuizEndpoint(c echo.Context) error {
	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	uid, err := auth.GetUserIdBySession(session.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	quizId := c.Param("id")
	_, err = quiz.GetQuizAccess(uid, quizId)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "quiz not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	quiz, err := quiz.GetQuizById(quizId)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "quiz not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	var questions []models.Question
	singleChoiceQuestions, _ := question.GetSingleChoiceQuestions(quizId)
	for _, q := range singleChoiceQuestions {
		questions = append(questions, models.SingleChoiceQuestion{
			ID:            q.UUID,
			QuizID:        q.QuizID,
			Question:      q.Question,
			Answers:       q.Answers,
			CorrectAnswer: q.CorrectAnswer,
		})
	}
	multipleChoiceQuestions, _ := question.GetMultipleChoiceQuestions(quizId)
	for _, q := range multipleChoiceQuestions {
		questions = append(questions, models.MultipleChoiceQuestion{
			ID:             q.UUID,
			QuizID:         q.QuizID,
			Question:       q.Question,
			Answers:        q.Answers,
			CorrectAnswers: q.CorrectAnswers,
		})
	}
	trueOrFalseQuestions, _ := question.GetTrueOrFalseQuestions(quizId)
	for _, q := range trueOrFalseQuestions {
		questions = append(questions, models.TrueOrFalseQuestion{
			ID:            q.UUID,
			QuizID:        q.QuizID,
			Question:      q.Question,
			CorrectAnswer: q.CorrectAnswer,
		})
	}

	userinfo, err := auth.GetUserById(uid)
	if err != nil {
		return c.JSON(http.StatusOK, models.Quiz{
			QuizInfo: models.QuizInfo{
				Id: quiz.Id, Title: quiz.Name, Description: quiz.Description.String, CreatorName: "Deleted",
			},
			Questions: questions,
		})
	}
	return c.JSON(http.StatusOK, models.Quiz{
		QuizInfo: models.QuizInfo{
			Id: quiz.Id, Title: quiz.Name, Description: quiz.Description.String, CreatorName: userinfo.Name, CreatorId: userinfo.Id,
		},
		Questions: questions,
	})
}

func GetQuizzesOfUserEndpoint(c echo.Context) error {
	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	uid, err := auth.GetUserIdBySession(session.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	quizAccesses, err := quiz.GetQuizAccessesOfUser(uid)
	if err != nil {
		fmt.Println("failed to get quiz accesses of user")
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	var quizzes []models.QuizInfo
	for _, acc := range *quizAccesses {
		quiz, err := quiz.GetQuizById(acc.QuizId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
		creator, err := auth.GetUserById(quiz.CreatorId.String)
		if err != nil {
			fmt.Println(err)
			quizzes = append(quizzes, models.QuizInfo{Id: quiz.Id, Title: quiz.Name, Description: quiz.Description.String, CreatorName: "Deleted"})
		} else {
			quizzes = append(quizzes, models.QuizInfo{Id: quiz.Id, Title: quiz.Name, Description: quiz.Description.String, CreatorName: creator.Name})
		}
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, QuizzesResponse{Quizzes: quizzes, Length: len(quizzes)})
}

func UpdateQuizEndpoint(c echo.Context) error {
	request := QuizRequestBody{}
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}
	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	uid, err := auth.GetUserIdBySession(session.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	quizId := c.Param("id")
	role, err := quiz.GetQuizAccess(uid, quizId)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "unauthorized")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	if role != 1 {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	err = quiz.UpdateQuiz(quizId, request.Name, request.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	quiz, err := quiz.GetQuizById(quizId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	if !quiz.CreatorId.Valid {
		return c.JSON(http.StatusOK, models.QuizInfo{Id: quiz.Id, Title: quiz.Name, Description: quiz.Description.String, CreatorName: "Deleted"})
	}
	creator, err := auth.GetUserById(quiz.CreatorId.String)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, models.QuizInfo{Id: quiz.Id, Title: quiz.Name, Description: quiz.Description.String, CreatorName: "Deleted"})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, models.QuizInfo{Id: quiz.Id, Title: quiz.Name, Description: quiz.Description.String, CreatorName: creator.Name, CreatorId: creator.Id})
}

func DeleteQuizEndpoint(c echo.Context) error {
	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	uid, err := auth.GetUserIdBySession(session.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	quizId := c.Param("id")
	role, err := quiz.GetQuizAccess(uid, quizId)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "unauthorized")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	if role != 1 {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	err = quiz.DeleteQuiz(quizId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, "quiz deleted")
}
