package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"net/http"
	"slices"
	"spaced-ace-backend/api/models"
	"spaced-ace-backend/auth"
	"spaced-ace-backend/constants"
	"spaced-ace-backend/db"
	"spaced-ace-backend/question"
	"spaced-ace-backend/quiz"
	"spaced-ace-backend/utils"
	"time"
)

type ReviewItemsRequestBody struct {
	QuizID     string `json:"quiz"`
	Difficulty string `json:"difficulty"`
	Status     string `json:"status"`
	Page       int    `json:"page"`
	Query      string `json:"query"`
}

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

	// Create review items for the quiz's questions
	_, err = createAndStoreReviewItems(ctx, sessionUserID, quizID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("creating the review items for quiz with ID %q: %w", quizID, err))
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

	// Delete the affected review items
	err = deleteReviewItems(ctx, sessionUserID, quizID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("deleting the review items for quiz with ID %q: %w", quizID, err))
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

func createAndStoreReviewItems(ctx context.Context, userID, quizID string) ([]*models.ReviewItem, error) {
	dbSingleChoiceQuestions, err := question.GetSingleChoiceQuestions(quizID)
	if err != nil {
		return nil, fmt.Errorf("getting single choice questions for quiz with ID %q\n", quizID)
	}

	reviewItems := make([]*models.ReviewItem, 0, len(dbSingleChoiceQuestions))

	for _, dbQuestion := range dbSingleChoiceQuestions {
		reviewItem, err := createSingleChoiceReviewItem(ctx, userID, dbQuestion.UUID)
		if err != nil {
			return nil, fmt.Errorf("creating review item for single choice question with ID %q: %w\n", dbQuestion.UUID, err)
		}
		reviewItems = append(reviewItems, reviewItem)
	}

	dbMultipleChoiceQuestions, err := question.GetMultipleChoiceQuestions(quizID)
	if err != nil {
		return nil, fmt.Errorf("getting multiple choice questions for quiz with ID %q\n", quizID)
	}

	for _, dbQuestion := range dbMultipleChoiceQuestions {
		reviewItem, err := createMultipleChoiceReviewItem(ctx, userID, dbQuestion.UUID)
		if err != nil {
			return nil, fmt.Errorf("creating review item for multiple choice question with ID %q: %w\n", dbQuestion.UUID, err)
		}
		reviewItems = append(reviewItems, reviewItem)
	}

	dbTrueOrFalseQuestions, err := question.GetTrueOrFalseQuestions(quizID)
	if err != nil {
		return nil, fmt.Errorf("getting true or false questions for quiz with ID %q\n", quizID)
	}

	for _, dbQuestion := range dbTrueOrFalseQuestions {
		reviewItem, err := createTrueOrFalseReviewItem(ctx, userID, dbQuestion.UUID)
		if err != nil {
			return nil, fmt.Errorf("creating review item for true or false question with ID %q: %w\n", dbQuestion.UUID, err)
		}
		reviewItems = append(reviewItems, reviewItem)
	}

	return reviewItems, nil
}
func createSingleChoiceReviewItem(ctx context.Context, userID, questionID string) (*models.ReviewItem, error) {
	sqlcQuerier := utils.GetQuerier()

	reviewItemID, err := sqlcQuerier.CreateSingleChoiceReviewItem(
		ctx,
		db.CreateSingleChoiceReviewItemParams{
			ID:                     uuid.NewString(),
			UserID:                 userID,
			SingleChoiceQuestionID: &questionID,
			EaseFactor:             constants.EASE_FACTOR_DEFAULT,
			Difficulty:             constants.REVIEW_ITEM_DIFFICULTY_DEFAULT,
			Streak:                 constants.REVIEW_ITEM_STREAK_DEFAULT,
			NextReviewDate: pgtype.Timestamptz{
				Time:             time.Now().UTC(),
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
			IntervalInMinutes: constants.REVIEW_ITEM_INTERVAL_IN_MINUTES_DEFAULT,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("creating review item for single choice question with ID %q: %w\n", questionID, err)
	}

	dbReviewItem, err := sqlcQuerier.GetReviewItem(ctx, reviewItemID)
	if err != nil {
		return nil, fmt.Errorf("getting review item with ID %q: %w\n", reviewItemID, err)
	}

	reviewItem, err := models.MapReviewItem(dbReviewItem)
	if err != nil {
		return nil, fmt.Errorf("mapping review item: %w\n", err)
	}

	return reviewItem, nil
}
func createMultipleChoiceReviewItem(ctx context.Context, userID, questionID string) (*models.ReviewItem, error) {
	sqlcQuerier := utils.GetQuerier()

	reviewItemID, err := sqlcQuerier.CreateMultipleChoiceReviewItem(
		ctx,
		db.CreateMultipleChoiceReviewItemParams{
			ID:                       uuid.NewString(),
			UserID:                   userID,
			MultipleChoiceQuestionID: &questionID,
			EaseFactor:               constants.EASE_FACTOR_DEFAULT,
			Difficulty:               constants.REVIEW_ITEM_DIFFICULTY_DEFAULT,
			Streak:                   constants.REVIEW_ITEM_STREAK_DEFAULT,
			NextReviewDate: pgtype.Timestamptz{
				Time:             time.Now().UTC(),
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
			IntervalInMinutes: constants.REVIEW_ITEM_INTERVAL_IN_MINUTES_DEFAULT,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("creating review item for multiple choice question with ID %q: %w\n", questionID, err)
	}

	dbReviewItem, err := sqlcQuerier.GetReviewItem(ctx, reviewItemID)
	if err != nil {
		return nil, fmt.Errorf("getting review item with ID %q: %w\n", reviewItemID, err)
	}

	reviewItem, err := models.MapReviewItem(dbReviewItem)
	if err != nil {
		return nil, fmt.Errorf("mapping review item: %w\n", err)
	}

	return reviewItem, nil
}
func createTrueOrFalseReviewItem(ctx context.Context, userID, questionID string) (*models.ReviewItem, error) {
	sqlcQuerier := utils.GetQuerier()

	reviewItemID, err := sqlcQuerier.CreateTrueOrFalseReviewItem(
		ctx,
		db.CreateTrueOrFalseReviewItemParams{
			ID:                    uuid.NewString(),
			UserID:                userID,
			TrueOrFalseQuestionID: &questionID,
			EaseFactor:            constants.EASE_FACTOR_DEFAULT,
			Difficulty:            constants.REVIEW_ITEM_DIFFICULTY_DEFAULT,
			Streak:                constants.REVIEW_ITEM_STREAK_DEFAULT,
			NextReviewDate: pgtype.Timestamptz{
				Time:             time.Now().UTC(),
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
			IntervalInMinutes: constants.REVIEW_ITEM_INTERVAL_IN_MINUTES_DEFAULT,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("creating review item for true or false question with ID %q: %w\n", questionID, err)
	}

	dbReviewItem, err := sqlcQuerier.GetReviewItem(ctx, reviewItemID)
	if err != nil {
		return nil, fmt.Errorf("getting review item with ID %q: %w\n", reviewItemID, err)
	}

	reviewItem, err := models.MapReviewItem(dbReviewItem)
	if err != nil {
		return nil, fmt.Errorf("mapping review item: %w\n", err)
	}

	return reviewItem, nil
}

func deleteReviewItems(ctx context.Context, userID, quizID string) error {
	sqlcQuerier := utils.GetQuerier()

	return sqlcQuerier.DeleteReviewItemsByQuizID(
		ctx,
		db.DeleteReviewItemsByQuizIDParams{
			UserID: userID,
			ID:     quizID,
		},
	)
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
			return nil, nil, fmt.Errorf("getting quiz with ID %q: %w\n", access.QuizId, err)
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

func validateReviewItemsRequest(request ReviewItemsRequestBody) (*ReviewItemsRequestBody, error) {
	if !slices.Contains([]string{"", "easy", "medium", "hard"}, request.Difficulty) {
		return nil, fmt.Errorf("invalid difficulty value %q\n", request.Difficulty)
	}

	if !slices.Contains([]string{"", "due", "not-due"}, request.Status) {
		return nil, fmt.Errorf("invalid status value %q\n", request.Status)
	}

	if request.Page < 1 {
		request.Page = 1
	}

	return &request, nil
}
