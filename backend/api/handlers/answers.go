package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"io"
	"log"
	"net/http"
	"slices"
	"spaced-ace-backend/api/models"
	"spaced-ace-backend/db"
	"spaced-ace-backend/utils"
	"time"
)

type AnswerRequestBody struct {
	QuestionId string `json:"questionId"`
	AnswerType string `json:"answerType"`
}
type SingleChoiceAnswerRequestBody struct {
	AnswerRequestBody
	Answer string `json:"answer"`
}
type MultipleChoiceAnswerRequestBody struct {
	AnswerRequestBody
	Answers []string `json:"answers"`
}
type TrueOrFalseAnswerRequestBody struct {
	AnswerRequestBody
	Answer bool `json:"answer"`
}

func PutCreateOrUpdateAnswer(c echo.Context) error {
	quizSessionId := c.Param("quizSessionId")
	if quizSessionId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing path param quizSessionId")
	}

	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Default().Println(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad request: %s", err.Error()))
	}
	c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var answerRequestBody AnswerRequestBody
	err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&answerRequestBody)
	if err != nil {
		log.Default().Println(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad request: %s", err.Error()))
	}

	if answerRequestBody.AnswerType == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing body param answerType")
	} else if !slices.Contains([]string{"single-choice", "multiple-choice", "true-or-false"}, answerRequestBody.AnswerType) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid questionType: `%s`", answerRequestBody.AnswerType))
	}
	if answerRequestBody.QuestionId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing body param questionId")
	}

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	quizSession, err := sqlcQuerier.GetQuizSession(ctx, quizSessionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("quizSession not found: %s", err))
	}
	if quizSession.FinishedAt.Valid {
		return echo.NewHTTPError(http.StatusForbidden, "modifying answer for a submitted quiz is not allowed")
	}

	switch answerRequestBody.AnswerType {
	case "single-choice":
		var requestBody SingleChoiceAnswerRequestBody
		if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&requestBody); err != nil {
			log.Default().Println(err.Error())
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error parsing single-choice answer: %s", err.Error()))
		}

		if !slices.Contains([]string{"A", "B", "C", "D"}, requestBody.Answer) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid answer: `%s` for single-choice question", requestBody.Answer))
		}

		oldAnswer, err := sqlcQuerier.GetSingleChoiceAnswerBySessionAndQuestionId(
			ctx,
			db.GetSingleChoiceAnswerBySessionAndQuestionIdParams{
				SessionID:  quizSessionId,
				QuestionID: requestBody.QuestionId,
			},
		)

		var answer *db.SingleChoiceAnswer
		var dbError error

		if err == nil {
			answer, dbError = sqlcQuerier.UpdateSingleChoiceAnswerBySessionAndQuestionId(
				ctx,
				db.UpdateSingleChoiceAnswerBySessionAndQuestionIdParams{
					SessionID:  oldAnswer.SessionID,
					QuestionID: oldAnswer.QuestionID,
					Answer:     []string{requestBody.Answer},
				},
			)
		} else {
			answer, dbError = sqlcQuerier.CreateSingleChoiceAnswer(
				ctx,
				db.CreateSingleChoiceAnswerParams{
					ID:         uuid.NewString(),
					SessionID:  quizSessionId,
					QuestionID: requestBody.QuestionId,
					Answer:     []string{requestBody.Answer},
				},
			)
		}

		if dbError != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("db error: %s", dbError))
		}

		result, err := models.MapSingleChoiceAnswer(answer)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, result)
	case "multiple-choice":
		var requestBody MultipleChoiceAnswerRequestBody
		if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&requestBody); err != nil {
			log.Default().Println(err.Error())
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error parsing multiple-choice answer: %s", err.Error()))
		}

		seen := make(map[string]bool)
		for _, answer := range requestBody.Answers {
			if !slices.Contains([]string{"A", "B", "C", "D"}, answer) {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid answer: `%s` for multiple-choice question", answer))
			}
			if seen[answer] {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("duplicated value: ``%s", answer))
			}
			seen[answer] = true
		}

		oldAnswer, err := sqlcQuerier.GetMultipleChoiceAnswerBySessionAndQuestionId(
			ctx,
			db.GetMultipleChoiceAnswerBySessionAndQuestionIdParams{
				SessionID:  quizSessionId,
				QuestionID: requestBody.QuestionId,
			},
		)

		var answer *db.MultipleChoiceAnswer
		var dbError error

		if err == nil {
			answer, dbError = sqlcQuerier.UpdateMultipleChoiceAnswerBySessionAndQuestionId(
				ctx,
				db.UpdateMultipleChoiceAnswerBySessionAndQuestionIdParams{
					SessionID:  oldAnswer.SessionID,
					QuestionID: oldAnswer.QuestionID,
					Answers:    requestBody.Answers,
				},
			)
		} else {
			answer, dbError = sqlcQuerier.CreateMultipleChoiceAnswer(
				ctx,
				db.CreateMultipleChoiceAnswerParams{
					ID:         uuid.NewString(),
					SessionID:  quizSessionId,
					QuestionID: requestBody.QuestionId,
					Answers:    requestBody.Answers,
				},
			)
		}

		if dbError != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("db error: %s", dbError))
		}

		result, err := models.MapMultipleChoiceAnswer(answer)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, result)

	case "true-or-false":
		var requestBody TrueOrFalseAnswerRequestBody
		if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&requestBody); err != nil {
			log.Default().Println(err.Error())
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error parsing true-or-false answer: %s", err.Error()))
		}

		oldAnswer, err := sqlcQuerier.GetTrueOrFalseAnswerBySessionAndQuestionId(
			ctx,
			db.GetTrueOrFalseAnswerBySessionAndQuestionIdParams{
				SessionID:  quizSessionId,
				QuestionID: requestBody.QuestionId,
			},
		)

		var answer *db.TrueOrFalseAnswer
		var dbError error

		if err == nil {
			answer, dbError = sqlcQuerier.UpdateTrueOrFalseAnswerBySessionAndQuestionId(
				ctx,
				db.UpdateTrueOrFalseAnswerBySessionAndQuestionIdParams{
					SessionID:  oldAnswer.SessionID,
					QuestionID: oldAnswer.QuestionID,
					Answer: pgtype.Bool{
						Bool:  requestBody.Answer,
						Valid: true,
					},
				},
			)
		} else {
			answer, dbError = sqlcQuerier.CreateTrueOrFalseAnswer(
				ctx,
				db.CreateTrueOrFalseAnswerParams{
					ID:         uuid.NewString(),
					SessionID:  quizSessionId,
					QuestionID: requestBody.QuestionId,
					Answer: pgtype.Bool{
						Bool:  requestBody.Answer,
						Valid: true,
					},
				},
			)
		}

		if dbError != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("db error: %s", dbError))
		}

		result, err := models.MapTrueOrFalseAnswer(answer)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, result)
	}

	return echo.NewHTTPError(http.StatusInternalServerError, "unreachable code")
}

func GetAnswers(c echo.Context) error {
	quizSessionId := c.Param("quizSessionId")
	if quizSessionId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing path param quizSessionId")
	}

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbSingleChoiceAnswers, err := sqlcQuerier.GetSingleChoiceAnswers(
		ctx,
		quizSessionId,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("error getting single choice answers: %s", err))
	}

	singleChoiceAnswers := make([]models.SingleChoiceAnswer, len(dbSingleChoiceAnswers))
	for i, dbAnswer := range dbSingleChoiceAnswers {
		answer, err := models.MapSingleChoiceAnswer(dbAnswer)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error parsing a single choice answer: %s", err))
		}
		singleChoiceAnswers[i] = *answer
	}

	dbMultipleChoiceAnswers, err := sqlcQuerier.GetMultipleChoiceAnswers(
		ctx,
		quizSessionId,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("error getting multiple choice answers: %s", err))
	}

	multipleChoiceAnswers := make([]models.MultipleChoiceAnswer, len(dbMultipleChoiceAnswers))
	for i, dbAnswer := range dbMultipleChoiceAnswers {
		answer, err := models.MapMultipleChoiceAnswer(dbAnswer)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error parsing a multiple choice answer: %s", err))
		}
		multipleChoiceAnswers[i] = *answer
	}

	dbTrueOrFalseAnswers, err := sqlcQuerier.GetTrueOrFalseAnswers(
		ctx,
		quizSessionId,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("error getting true or false answers: %s", err))
	}

	trueOrFalseAnswers := make([]models.TrueOrFalseAnswer, len(dbTrueOrFalseAnswers))
	for i, dbAnswer := range dbTrueOrFalseAnswers {
		answer, err := models.MapTrueOrFalseAnswer(dbAnswer)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error parsing a true or false answer: %s", err))
		}
		trueOrFalseAnswers[i] = *answer
	}

	response := models.AnswersResponse{
		SingleChoiceAnswers:   singleChoiceAnswers,
		MultipleChoiceAnswers: multipleChoiceAnswers,
		TrueOrFalseAnswer:     trueOrFalseAnswers,
	}
	return c.JSON(http.StatusOK, response)
}
