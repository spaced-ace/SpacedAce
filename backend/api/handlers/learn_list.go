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

func GetLearnList(c echo.Context) error {
	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	sessionUserID, err := auth.GetUserIdBySession(session.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	quizAccesses, err := quiz.GetQuizAccessesOfUser(sessionUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting quiz accesses for user with ID `%s`: %w", sessionUserID, err))
	}

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	addedQuizzes, err := sqlcQuerier.GetAddedLearnListItems(ctx, sessionUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting added quizes for user with ID `%s`: %w", sessionUserID, err))
	}

	available, selected, err := buildLearnListItems(*quizAccesses, addedQuizzes)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, models.LearnListResponseBody{
		AvailableItems: available,
		SelectedItems:  selected,
	})
}
func PostAddQuizToLearnList(c echo.Context) error {
	quizID := c.Param("quizID")
	if quizID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("missing path param quizID"))
	}

	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	sessionUserID, err := auth.GetUserIdBySession(session.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Add the quiz to the user's learn list
	err = sqlcQuerier.AddQuizToLearnList(
		ctx,
		db.AddQuizToLearnListParams{
			UserID: sessionUserID,
			QuizID: quizID,
		},
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("add quiz with ID %q to learn list for user with ID %q: %w", quizID, sessionUserID, err))
	}

	// Fetch the quizzes that user has access to
	quizAccesses, err := quiz.GetQuizAccessesOfUser(sessionUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting quiz accesses for user with ID `%s`: %w", sessionUserID, err))
	}

	// Fetch the quizzes that are already added to the user's learn list
	addedQuizzes, err := sqlcQuerier.GetAddedLearnListItems(ctx, sessionUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error getting added quizes for user with ID `%s`: %w", sessionUserID, err))
	}

	// Build the learn list
	available, selected, err := buildLearnListItems(*quizAccesses, addedQuizzes)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, models.LearnListResponseBody{
		AvailableItems: available,
		SelectedItems:  selected,
	})
}
func PostRemoveQuizFromLearnList(c echo.Context) error {
	quizID := c.Param("quizID")
	if quizID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("missing path param quizID"))
	}

	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	sessionUserID, err := auth.GetUserIdBySession(session.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Remove the quiz from the user's learn list
	err = sqlcQuerier.RemoveQuizFromLearnList(
		ctx,
		db.RemoveQuizFromLearnListParams{
			UserID: sessionUserID,
			QuizID: quizID,
		},
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("removing quiz with ID %q to learn list for user with ID %q: %w", quizID, sessionUserID, err))
	}

	// Fetch the quizzes that user has access to
	quizAccesses, err := quiz.GetQuizAccessesOfUser(sessionUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("getting quiz accesses for user with ID `%s`: %w", sessionUserID, err))
	}

	// Fetch the quizzes that are already added to the user's learn list
	addedQuizzes, err := sqlcQuerier.GetAddedLearnListItems(ctx, sessionUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("getting added quizes for user with ID `%s`: %w", sessionUserID, err))
	}

	// Build the learn list
	available, selected, err := buildLearnListItems(*quizAccesses, addedQuizzes)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, models.LearnListResponseBody{
		AvailableItems: available,
		SelectedItems:  selected,
	})
}

func buildLearnListItems(quizAccesses []quiz.DBQuizAccess, addedQuizzes []*db.LearnListAddedItem) (available, selected []models.LearnListItem, err error) {
	available = make([]models.LearnListItem, 0, len(quizAccesses)-len(addedQuizzes))
	selected = make([]models.LearnListItem, 0, len(addedQuizzes))

	addedQuizMap := make(map[string]struct{}, len(addedQuizzes))
	for _, item := range addedQuizzes {
		addedQuizMap[item.QuizID] = struct{}{}
	}

	for _, access := range quizAccesses {
		dbQuiz, err := quiz.GetQuizById(access.QuizId)
		if err != nil {
			return nil, nil, fmt.Errorf("getting quiz with ID %q: %w", access.QuizId, err)
		}

		item := models.LearnListItem{
			QuizID:   dbQuiz.Id,
			QuizName: dbQuiz.Name,
		}

		if _, isAdded := addedQuizMap[access.QuizId]; isAdded {
			selected = append(selected, item)
		} else {
			available = append(available, item)
		}
	}

	return available, selected, nil
}
