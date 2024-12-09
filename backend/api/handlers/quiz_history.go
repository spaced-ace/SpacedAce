package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"net/http"
	"spaced-ace-backend/api/models"
	"spaced-ace-backend/auth"
	"spaced-ace-backend/db"
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

	quizSessions, err := sqlcQuerier.GetQuizSessionsByUserId(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting quiz sessions for user: %w", err))
	}

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

	quizHistoryEntries := make([]models.QuizHistoryEntry, 0, len(quizSessions))
	for _, quizSession := range quizSessions {
		if quizSession.FinishedAt.Time.IsZero() {
			quizHistoryEntries = append(quizHistoryEntries, models.QuizHistoryEntry{
				QuizID:          quizSession.QuizID,
				QuizName:        quizNameMap[quizSession.QuizID],
				SessionID:       quizSession.ID,
				Finished:        false,
				DateTaken:       quizSession.StartedAt.Time,
				TimeSpent:       0,
				ScorePercentage: 0,
			})
		} else {
			var result *db.GetQuizResultsByUserIDRow
			for _, qr := range quizResults {
				if qr.SessionID == quizSession.ID {
					result = qr
					break
				}
			}

			timeSpent := time.Duration(0)
			percentage := 0.0

			if result == nil {
				quizResult, err := calculateAndStoreQuizResult(ctx, quizSession.ID, quizSession.QuizID)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "error calculating the quiz result: "+err.Error())
				}

				timeSpent = quizSession.FinishedAt.Time.Sub(quizSession.StartedAt.Time)
				percentage = quizResult.Score / quizResult.MaxScore * 100
			} else {
				if result.FinishedAt.Valid {
					timeSpent = result.FinishedAt.Time.Sub(result.StartedAt.Time)
				}
				if result.MaxScore != 0.0 {
					percentage = result.Score / result.MaxScore * 100
				}
			}

			quizHistoryEntries = append(quizHistoryEntries, models.QuizHistoryEntry{
				QuizID:          quizSession.QuizID,
				QuizName:        quizNameMap[quizSession.QuizID],
				SessionID:       quizSession.ID,
				Finished:        quizSession.FinishedAt.Valid,
				DateTaken:       quizSession.StartedAt.Time,
				TimeSpent:       timeSpent,
				ScorePercentage: percentage,
			})
		}
	}

	responseBody := models.QuizHistoryEntriesResponseBody{
		QuizHistoryEntries: quizHistoryEntries,
		Length:             len(quizHistoryEntries),
	}
	return c.JSON(http.StatusOK, responseBody)
}
