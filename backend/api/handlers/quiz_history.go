package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"net/http"
	"spaced-ace-backend/api/models"
	"spaced-ace-backend/auth"
	"spaced-ace-backend/quiz"
	"spaced-ace-backend/utils"
	"time"
)

func GetQuizHistoryEntries(c echo.Context) error {
	userID := c.QueryParam("userID")
	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing query param userID")
	}

	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	sessionUserID, err := auth.GetUserIdBySession(session.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	if userID != sessionUserID {
		return echo.NewHTTPError(http.StatusForbidden, "cannot get quiz results for another user")
	}

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	quizResults, err := sqlcQuerier.GetQuizResultsByUserID(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting quiz results for user: %w", err))
	}

	quizNameMap := make(map[string]string)
	for _, result := range quizResults {
		if quizNameMap[result.QuizID] == "" {
			dbQuiz, err := quiz.GetQuizById(result.QuizID)
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("quiz not found with ID `%s`: %w", result.QuizID, err))
			}
			quizNameMap[result.QuizID] = dbQuiz.Name
		}
	}

	quizHistoryEntries := make([]models.QuizHistoryEntry, 0, len(quizResults))
	for i, result := range quizResults {
		timeSpent := time.Duration(0)
		if result.FinishedAt.Valid {
			timeSpent = result.FinishedAt.Time.Sub(result.StartedAt.Time)
		}

		quizHistoryEntries[i] = models.QuizHistoryEntry{
			QuizName:        quizNameMap[result.QuizID],
			SessionID:       result.SessionID,
			Finished:        result.FinishedAt.Valid,
			DateTaken:       result.StartedAt.Time,
			TimeSpent:       timeSpent,
			ScorePercentage: result.Score / result.MaxScore * 100,
		}
	}

	return c.JSON(http.StatusOK, quizHistoryEntries)
}
