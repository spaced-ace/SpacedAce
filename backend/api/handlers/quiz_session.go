package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"spaced-ace-backend/api/models"
	"spaced-ace-backend/auth"
	"spaced-ace-backend/db"
	"spaced-ace-backend/utils"
	"time"
)

type StartQuizSessionRequestBody struct {
	QuizId string `json:"quizId"`
}
type GetQuizSessionRequestBody struct {
	QuizId string `json:"quizId"`
	UserId string `json:"userId"`
}

func StartQuizSession(c echo.Context) error {
	var request StartQuizSessionRequestBody
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	userId, err := auth.GetUserIdBySession(session.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// List quiz sessions
	openQuizSessions, err := sqlcQuerier.GetQuizSessionsByQuizIdAndUserId(
		ctx,
		db.GetQuizSessionsByQuizIdAndUserIdParams{
			QuizID: request.QuizId,
			UserID: userId,
		},
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error getting open quiz sessions: "+err.Error())
	}

	// Close the open quiz sessions
	for _, openQuizSession := range openQuizSessions {
		// Skip already closed sessions
		if openQuizSession.FinishedAt.Valid {
			continue
		}

		_, err := sqlcQuerier.UpdateQuizSessionFinishedAt(
			ctx,
			db.UpdateQuizSessionFinishedAtParams{
				ID: openQuizSession.ID,
				FinishedAt: pgtype.Timestamp{
					Time:             time.Now(),
					InfinityModifier: pgtype.Finite,
					Valid:            true,
				},
			},
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error closing quiz session: %s", openQuizSession.ID))
		}
		log.Default().Printf("quiz session closed: %s", openQuizSession.ID)
	}

	// Start a new quiz session
	dbQuizSession, err := sqlcQuerier.CreateQuizSession(
		ctx,
		db.CreateQuizSessionParams{
			ID:     uuid.NewString(),
			UserID: userId,
			QuizID: request.QuizId,
		},
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error starting a new quiz session")
	}
	log.Default().Printf("quiz session started: %s", dbQuizSession.ID)

	quizSession, err := models.MapQuizSession(dbQuizSession)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, quizSession)
}

func GetQuizSessions(c echo.Context) error {
	userId := c.QueryParam("userId")
	quizId := c.QueryParam("quizId")
	open := c.QueryParam("open")

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbQuizSessions, err := sqlcQuerier.GetQuizSessionsByQuizIdAndUserId(
		ctx,
		db.GetQuizSessionsByQuizIdAndUserIdParams{
			QuizID: quizId,
			UserID: userId,
		},
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	filter := func(s db.QuizSession) bool {
		if open == "true" {
			return !s.FinishedAt.Valid
		} else if open == "false" {
			return s.FinishedAt.Valid
		}
		return true
	}

	resultCount := 0
	for _, dbQuizSession := range dbQuizSessions {
		if filter(dbQuizSession) {
			resultCount++
		}
	}

	quizSessions := make([]models.QuizSession, resultCount)
	nextIndex := 0
	for _, dbQuizSession := range dbQuizSessions {
		if !filter(dbQuizSession) {
			continue
		}

		quizSession, err := models.MapQuizSession(dbQuizSession)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		quizSessions[nextIndex] = *quizSession
		nextIndex++
	}

	response := models.GetQuizSessionsResponseBody{
		QuizSessions: quizSessions,
		Length:       resultCount,
	}
	return c.JSON(http.StatusOK, response)
}

func GetQuizSession(c echo.Context) error {
	quizSessionId := c.Param("quizSessionId")

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbQuizSession, err := sqlcQuerier.GetQuizSession(ctx, quizSessionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	quizSession, err := models.MapQuizSession(dbQuizSession)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, quizSession)
}

func StopQuizSession(c echo.Context) error {
	quizSessionId := c.Param("quizSessionId")

	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	userId, err := auth.GetUserIdBySession(session.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	quizSession, err := sqlcQuerier.GetQuizSession(
		ctx,
		quizSessionId,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("quiz session '%s' is not found", quizSessionId))
	}

	if quizSession.FinishedAt.Valid {
		return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("quiz session '%s' is already closed", quizSessionId))
	}

	// TODO when roles (like admin and normal user) will be available, we should enable this action for admins
	if quizSession.UserID != userId {
		return echo.NewHTTPError(http.StatusForbidden, "cannot stop other user's quiz session")
	}

	_, err = sqlcQuerier.UpdateQuizSessionFinishedAt(
		ctx,
		db.UpdateQuizSessionFinishedAtParams{
			ID: quizSessionId,
			FinishedAt: pgtype.Timestamp{
				Time:             time.Now(),
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
		},
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error closing quiz session: %s", quizSessionId))
	}
	log.Default().Printf("quiz session closed: %s", quizSessionId)

	return c.NoContent(http.StatusOK)
}

func HasOpenQuizSession(c echo.Context) error {
	userId := c.QueryParam("userId")
	quizId := c.QueryParam("quizId")

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hasOpenSession, err := sqlcQuerier.HasOpenQuizSession(
		ctx,
		db.HasOpenQuizSessionParams{
			QuizID: quizId,
			UserID: userId,
		},
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if hasOpenSession {
		return c.NoContent(http.StatusOK)
	} else {
		return c.NoContent(http.StatusNotFound)
	}
}
