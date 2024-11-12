package models

import (
	"fmt"
	"spaced-ace-backend/db"
)

type LearnListItem struct {
	QuizID   string `json:"quizID"`
	QuizName string `json:"quizName"`
}
type LearnListResponseBody struct {
	AvailableItems []LearnListItem `json:"availableItems"`
	SelectedItems  []LearnListItem `json:"selectedItems"`
}

type ReviewItem struct {
	ID                       string       `json:"id"`
	UserID                   string       `json:"userID"`
	QuizID                   string       `json:"quizID"`
	QuizName                 string       `json:"quizName"`
	SingleChoiceQuestionID   *string      `json:"singleChoiceQuestionID"`
	MultipleChoiceQuestionID *string      `json:"multipleChoiceQuestionID"`
	TrueOrFalseQuestionID    *string      `json:"trueOrFalseQuestionID"`
	QuestionName             string       `json:"questionName"`
	EaseFactor               float64      `json:"easeFactor"`
	Difficulty               float64      `json:"difficulty"`
	Streak                   int32        `json:"streak"`
	NextReviewDate           NullableTime `json:"nextReviewDate"`
	IntervalInMinutes        int32        `json:"intervalInMinutes"`
}
type ReviewItemResponseBody struct {
	ReviewItems              []*ReviewItem `json:"reviewItems"`
	ReviewItemCountForFilter int           `json:"reviewItemCountForFilter"`
}
type QuizOptionsResponseBody struct {
	QuizOptions []*Option `json:"quizOptions"`
}
type ReviewItemCountsResponseBody struct {
	Total       int `json:"total"`
	DueToReview int `json:"dueToReview"`
}
type ReviewItemPageDataResponseBody struct {
	CurrentReviewItemID    string                  `json:"currentReviewItemID"`
	SingleChoiceQuestion   *SingleChoiceQuestion   `json:"singleChoiceQuestion"`
	MultipleChoiceQuestion *MultipleChoiceQuestion `json:"multipleChoiceQuestion"`
	TrueOrFalseQuestion    *TrueOrFalseQuestion    `json:"trueOrFalseQuestion"`
	NextReviewItemID       string                  `json:"nextReviewItemID"`
}

func MapReviewItem(dbItem *db.GetReviewItemRow) (*ReviewItem, error) {
	if dbItem == nil {
		return nil, fmt.Errorf("nil review item")
	}
	return &ReviewItem{
		ID:                       dbItem.ID,
		UserID:                   dbItem.UserID,
		QuizID:                   dbItem.QuizID,
		QuizName:                 dbItem.QuizName,
		SingleChoiceQuestionID:   dbItem.SingleChoiceQuestionID,
		MultipleChoiceQuestionID: dbItem.MultipleChoiceQuestionID,
		TrueOrFalseQuestionID:    dbItem.TrueOrFalseQuestionID,
		QuestionName:             dbItem.QuestionName,
		EaseFactor:               dbItem.EaseFactor,
		Difficulty:               dbItem.Difficulty,
		Streak:                   dbItem.Streak,
		NextReviewDate:           newNullableTime(dbItem.NextReviewDate.Time),
		IntervalInMinutes:        dbItem.IntervalInMinutes,
	}, nil
}
func MapReviewItemFromReviewItemsRow(dbItem *db.GetReviewItemsRow) (*ReviewItem, error) {
	if dbItem == nil {
		return nil, fmt.Errorf("nil review item")
	}
	return &ReviewItem{
		ID:                       dbItem.ID,
		UserID:                   dbItem.UserID,
		QuizID:                   dbItem.QuizID,
		QuizName:                 dbItem.QuizName,
		SingleChoiceQuestionID:   dbItem.SingleChoiceQuestionID,
		MultipleChoiceQuestionID: dbItem.MultipleChoiceQuestionID,
		TrueOrFalseQuestionID:    dbItem.TrueOrFalseQuestionID,
		QuestionName:             dbItem.QuestionName,
		EaseFactor:               dbItem.EaseFactor,
		Difficulty:               dbItem.Difficulty,
		Streak:                   dbItem.Streak,
		NextReviewDate:           newNullableTime(dbItem.NextReviewDate.Time),
		IntervalInMinutes:        dbItem.IntervalInMinutes,
	}, nil
}
