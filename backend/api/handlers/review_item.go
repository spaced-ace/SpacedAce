package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"net/http"
	"spaced-ace-backend/api/models"
	"spaced-ace-backend/auth"
	"spaced-ace-backend/constants"
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

func GetQuestionAndNextReviewItem(c echo.Context) error {
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

	nextReviewItemID := ""
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

		reviewItem, err = models.MapReviewItemFromReviewItemsRow(dbReviewItems[0])
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("mapping review item: %w\n", err))
		}

		if len(dbReviewItems) >= 2 {
			nextReviewItemID = dbReviewItems[1].ID
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

	response := models.ReviewItemPageDataResponseBody{
		CurrentReviewItemID:    reviewItem.ID,
		SingleChoiceQuestion:   singleChoiceQuestion,
		MultipleChoiceQuestion: multipleChoiceQuestion,
		TrueOrFalseQuestion:    trueOrFalseQuestion,
		NextReviewItemID:       nextReviewItemID,
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

	request := new(models.SubmitReviewItemQuestionRequestBody)
	if err = json.NewDecoder(c.Request().Body).Decode(request); err != nil {
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

	fmt.Printf("do something with %+v\n", reviewItem)
	// TODO get question, calculate the score with the answer and update the review item

	return c.JSON(http.StatusOK, reviewItem) // TODO return the updated review item
}
