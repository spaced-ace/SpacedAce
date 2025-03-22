package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"log"
	"math"
	"net/http"
	"slices"
	"spaced-ace-backend/api/models"
	"spaced-ace-backend/auth"
	"spaced-ace-backend/constants"
	"spaced-ace-backend/db"
	"spaced-ace-backend/question"
	"spaced-ace-backend/utils"
	"strings"
	"time"
)

func GetReviewItems(c echo.Context) error {
	session, err := c.Cookie("session")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	sessionUserID, err := auth.GetUserIdBySession(session.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	var request = ReviewItemsRequestBody{}
	if err = json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("parsing request body: %w\n", err))
	}

	filter, err := validateReviewItemsRequest(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("validationg the request body: %w\n", err))
	}

	sqlcQuerier := utils.GetQuerier()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbReviewItems, err := sqlcQuerier.GetReviewItems(ctx, sessionUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("getting review item for user with ID %q: %w\n", sessionUserID, err))
	}

	reviewItems := make([]*models.ReviewItem, 0, len(dbReviewItems))
	for _, dbReviewItem := range dbReviewItems {
		reviewItem, err := models.MapReviewItemFromReviewItemsRow(dbReviewItem)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("mapping review item: %w\n", err))
		}
		reviewItems = append(reviewItems, reviewItem)
	}

	filteredReviewItem := make([]*models.ReviewItem, 0, len(reviewItems))
	for _, item := range reviewItems {
		if filter.QuizID != "" && item.QuizID != filter.QuizID {
			continue
		}

		if filter.Status != "" {
			if filter.Status == "due" && !item.NextReviewDate.Time.Before(time.Now()) {
				continue
			}
			if filter.Status == "not-due" && item.NextReviewDate.Time.Before(time.Now()) {
				continue
			}
		}

		if filter.Difficulty != "" {
			if filter.Difficulty == "easy" && item.Difficulty > 1.5 {
				continue
			}
			if filter.Difficulty == "medium" && item.Difficulty <= 1.5 || item.Difficulty > 3.5 {
				continue
			}
			if filter.Difficulty == "hard" && item.Difficulty <= 3.5 {
				continue
			}
		}

		if filter.Query != "" && !strings.Contains(strings.ToLower(item.QuestionName), strings.ToLower(filter.Query)) {
			continue
		}

		filteredReviewItem = append(filteredReviewItem, item)
	}

	reviewItemCountForFilter := len(filteredReviewItem)

	lowerIndex := (filter.Page - 1) * constants.REVIEW_ITEM_PAGE_SIZE
	upperIndex := filter.Page * constants.REVIEW_ITEM_PAGE_SIZE
	reviewItemsOnPage := make([]*models.ReviewItem, 0, constants.REVIEW_ITEM_PAGE_SIZE)
	for i, item := range filteredReviewItem {
		if i >= lowerIndex && i <= upperIndex {
			reviewItemsOnPage = append(reviewItemsOnPage, item)
		}
	}

	response := models.ReviewItemResponseBody{
		ReviewItems:              reviewItemsOnPage,
		ReviewItemCountForFilter: reviewItemCountForFilter,
	}
	return c.JSON(http.StatusOK, response)
}
func GetQuizOptions(c echo.Context) error {
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

	dbQuizOptions, err := sqlcQuerier.GetQuizOptions(ctx, sessionUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("getting quiz options for user with ID %q: %w\n", sessionUserID, err))
	}

	quizOptions := make([]*models.Option, 0, len(dbQuizOptions))
	for _, dbQuizOption := range dbQuizOptions {
		quizOptions = append(
			quizOptions, &models.Option{
				Name:  dbQuizOption.QuizName,
				Value: dbQuizOption.QuizID,
			},
		)
	}

	return c.JSON(http.StatusOK, models.QuizOptionsResponseBody{QuizOptions: quizOptions})
}
func GetReviewItemCounts(c echo.Context) error {
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

	dbCounts, err := sqlcQuerier.GetReviewItemCounts(ctx, sessionUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("getting item counts for user with ID %q: %w\n", sessionUserID, err))
	}

	return c.JSON(
		http.StatusOK,
		models.ReviewItemCountsResponseBody{
			Total:       int(dbCounts.Total),
			DueToReview: int(dbCounts.DueToReview),
		},
	)
}

func GetReviewItemQuestion(c echo.Context) error {
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

	reviewItemID := c.Param("reviewItemID")

	var reviewItem *models.ReviewItem
	if reviewItemID != "" {
		dbReviewItem, err := sqlcQuerier.GetReviewItem(ctx, reviewItemID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("getting review item with ID %q: %w\n", reviewItemID, err))
		}

		reviewItem, err = models.MapReviewItem(dbReviewItem)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("mapping review item: %w\n", err))
		}
	} else {
		dbReviewItems, err := sqlcQuerier.GetReviewItems(ctx, sessionUserID)
		if err != nil || len(dbReviewItems) < 1 {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("getting enough review items for user with ID %q: %w\n", sessionUserID, err))
		}

		var reviewItemToDue *db.GetReviewItemsRow

		now := time.Now().UTC()
		for _, item := range dbReviewItems {
			if item.NextReviewDate.Time.UTC().Before(now) {
				reviewItemToDue = item
				break
			}
		}

		if reviewItemToDue == nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("getting review item to due"))
		}

		reviewItem, err = models.MapReviewItemFromReviewItemsRow(reviewItemToDue)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("mapping review item: %w\n", err))
		}
	}

	var singleChoiceQuestion *models.SingleChoiceQuestion
	if reviewItem.SingleChoiceQuestionID != nil {
		dbQuestion, err := question.GetSingleChoiceQuestion(*reviewItem.SingleChoiceQuestionID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("getting single choice question with ID %q: %w\n", reviewItem.SingleChoiceQuestionID, err))
		}

		singleChoiceQuestion = &models.SingleChoiceQuestion{
			ID:            dbQuestion.UUID,
			QuizID:        dbQuestion.QuizID,
			QuestionType:  models.SingleChoice,
			Question:      dbQuestion.Question,
			Answers:       dbQuestion.Answers,
			CorrectAnswer: dbQuestion.CorrectAnswer,
		}
	}

	var multipleChoiceQuestion *models.MultipleChoiceQuestion
	if reviewItem.MultipleChoiceQuestionID != nil {
		dbQuestion, err := question.GetMultipleChoiceQuestion(*reviewItem.MultipleChoiceQuestionID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("getting multiple choice question with ID %q: %w\n", reviewItem.MultipleChoiceQuestionID, err))
		}

		multipleChoiceQuestion = &models.MultipleChoiceQuestion{
			ID:             dbQuestion.UUID,
			QuizID:         dbQuestion.QuizID,
			QuestionType:   models.MultipleChoice,
			Question:       dbQuestion.Question,
			Answers:        dbQuestion.Answers,
			CorrectAnswers: dbQuestion.CorrectAnswers,
		}
	}

	var trueOrFalseQuestion *models.TrueOrFalseQuestion
	if reviewItem.TrueOrFalseQuestionID != nil {
		dbQuestion, err := question.GetTrueOrFalseQuestion(*reviewItem.TrueOrFalseQuestionID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("getting true or false question with ID %q: %w\n", reviewItem.TrueOrFalseQuestionID, err))
		}

		trueOrFalseQuestion = &models.TrueOrFalseQuestion{
			ID:            dbQuestion.UUID,
			QuizID:        dbQuestion.QuizID,
			QuestionType:  models.TrueOrFalse,
			Question:      dbQuestion.Question,
			CorrectAnswer: dbQuestion.CorrectAnswer,
		}
	}

	response := models.ReviewItemQuestionResponseBody{
		CurrentReviewItemID:    reviewItem.ID,
		SingleChoiceQuestion:   singleChoiceQuestion,
		MultipleChoiceQuestion: multipleChoiceQuestion,
		TrueOrFalseQuestion:    trueOrFalseQuestion,
	}
	return c.JSON(200, response)
}
func PostSubmitReviewItemQuestion(c echo.Context) error {
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

	reviewItemID := c.Param("reviewItemID")

	answers := new(models.SubmitReviewItemQuestionRequestBody)
	if err = json.NewDecoder(c.Request().Body).Decode(answers); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("parsing submit review item question request body: %w\n", err))
	}

	dbReviewItem, err := sqlcQuerier.GetReviewItem(ctx, reviewItemID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("getting review item with ID %q: %w\n", reviewItemID, err))
	}

	if dbReviewItem.UserID != sessionUserID {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("cannot submit other user's review item question"))
	}

	reviewItem, err := models.MapReviewItem(dbReviewItem)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("mapping review item with ID %q: %w\n", reviewItemID, err))
	}

	score, err := calculateReviewItemScore(reviewItem, answers)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("calculating score for review item with ID %q: %w\n", reviewItemID, err))
	}

	updatedReviewItem, err := applySpacedRepetitionAndStore(ctx, reviewItem, score)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("applying spaced repetition on review item with ID %q: %w\n", reviewItemID, err))
	}

	return c.JSON(http.StatusOK, updatedReviewItem)
}

func calculateReviewItemScore(reviewItem *models.ReviewItem, answers *models.SubmitReviewItemQuestionRequestBody) (float64, error) {
	if reviewItem.SingleChoiceQuestionID != nil {
		dbQuestion, err := question.GetSingleChoiceQuestion(*reviewItem.SingleChoiceQuestionID)
		if err != nil {
			return 0, fmt.Errorf("getting single choice question with ID %q: %w\n", *reviewItem.SingleChoiceQuestionID, err)
		}
		modelQuestion := dbQuestion.MapToModel()

		score, err := calculateReviewItemSingleChoiceQuestionScore(*modelQuestion, answers.SingleChoiceValue)
		if err != nil {
			return 0, fmt.Errorf("calculating single choice question score for review item with ID %q: %w\n", reviewItem.ID, err)
		}
		return score, nil
	}
	if reviewItem.MultipleChoiceQuestionID != nil {
		dbQuestion, err := question.GetMultipleChoiceQuestion(*reviewItem.MultipleChoiceQuestionID)
		if err != nil {
			return 0, fmt.Errorf("getting multiple choice question with ID %q: %w\n", *reviewItem.MultipleChoiceQuestionID, err)
		}
		modelQuestion := dbQuestion.MapToModel()

		score, err := calculateReviewItemMultipleChoiceQuestionScore(*modelQuestion, answers.MultipleChoiceValue)
		if err != nil {
			return 0, fmt.Errorf("calculating multiple choice question score for review item with ID %q: %w\n", reviewItem.ID, err)
		}
		return score, nil
	}
	if reviewItem.TrueOrFalseQuestionID != nil {
		dbQuestion, err := question.GetTrueOrFalseQuestion(*reviewItem.TrueOrFalseQuestionID)
		if err != nil {
			return 0, fmt.Errorf("getting true or false question with ID %q: %w\n", *reviewItem.TrueOrFalseQuestionID, err)
		}
		modelQuestion := dbQuestion.MapToModel()

		score, err := calculateReviewItemTrueOrFalseQuestionScore(*modelQuestion, answers.TrueOrFalseValue)
		if err != nil {
			return 0, fmt.Errorf("calculating true or false question score for review item with ID %q: %w\n", reviewItem.ID, err)
		}
		return score, nil
	}

	log.Default().Printf("none of the questions were actually answered, additional check may be needed")
	return 0, nil
}
func calculateReviewItemSingleChoiceQuestionScore(singleChoiceQuestion models.SingleChoiceQuestion, answer string) (float64, error) {
	if singleChoiceQuestion.CorrectAnswer == answer {
		return 1.0, nil
	}
	return 0, nil
}
func calculateReviewItemMultipleChoiceQuestionScore(multipleChoiceQuestion models.MultipleChoiceQuestion, answers []string) (float64, error) {
	positiveScore := 1.0 / float64(len(multipleChoiceQuestion.CorrectAnswers))
	negativeScore := 1.0 / float64(4-len(multipleChoiceQuestion.CorrectAnswers))

	score := 0.0
	for _, correctAnswer := range multipleChoiceQuestion.CorrectAnswers {
		if slices.Contains(answers, correctAnswer) {
			score += positiveScore
		} else {
			score -= negativeScore
		}
	}
	if score <= 0.0 {
		score = 0
	}
	return score, nil
}
func calculateReviewItemTrueOrFalseQuestionScore(trueOrFalseQuestion models.TrueOrFalseQuestion, answer bool) (float64, error) {
	if trueOrFalseQuestion.CorrectAnswer == answer {
		return 1.0, nil
	}
	return 0, nil
}

func applySpacedRepetitionAndStore(ctx context.Context, reviewItem *models.ReviewItem, percentage float64) (*models.ReviewItem, error) {
	// score is the user's performance rating, where:
	// 5 = perfect recall, 4 = correct with minor hesitation, 3 = correct but difficult,
	// 2 = incorrect, but partially remembered, 1 = completely incorrect
	score := 5 * percentage

	if score >= 3 {
		reviewItem.Difficulty = math.Max(1, reviewItem.Difficulty-1)

		if reviewItem.Streak == 0 {
			reviewItem.IntervalInMinutes = 120
		} else if reviewItem.Streak == 1 {
			reviewItem.IntervalInMinutes = 3 * 120
		} else {
			reviewItem.IntervalInMinutes = int32(float64(reviewItem.IntervalInMinutes) * reviewItem.EaseFactor * (1 / reviewItem.Difficulty))
		}

		reviewItem.EaseFactor = reviewItem.EaseFactor + (0.1 - (5-score)*0.08)
		reviewItem.Streak = reviewItem.Streak + 1
	} else {
		reviewItem.Difficulty = math.Min(5, reviewItem.Difficulty+1.5)

		reviewItem.IntervalInMinutes = 60

		reviewItem.EaseFactor = math.Max(1.3, reviewItem.EaseFactor-0.2)
		reviewItem.Streak = 0
	}

	reviewItem.NextReviewDate = models.NullableTime{
		Time: time.Now().Add(time.Duration(reviewItem.IntervalInMinutes) * time.Minute),
	}

	sqlcQuerier := utils.GetQuerier()
	err := sqlcQuerier.UpdateReviewItem(
		ctx,
		db.UpdateReviewItemParams{
			ID:         reviewItem.ID,
			EaseFactor: reviewItem.EaseFactor,
			Difficulty: reviewItem.Difficulty,
			Streak:     reviewItem.Streak,
			NextReviewDate: pgtype.Timestamptz{
				Time:             reviewItem.NextReviewDate.Time,
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
			IntervalInMinutes: reviewItem.IntervalInMinutes,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("storing updated review item with ID %q: %w\n", reviewItem.ID, err)
	}

	return reviewItem, nil
}
