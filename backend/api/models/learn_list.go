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
	SingleChoiceQuestionID   *string      `json:"singleChoiceQuestionID"`
	MultipleChoiceQuestionID *string      `json:"multipleChoiceQuestionID"`
	TrueOrFalseQuestionID    *string      `json:"trueOrFalseQuestionID"`
	EaseFactor               float64      `json:"easeFactor"`
	Difficulty               float64      `json:"difficulty"`
	Streak                   int32        `json:"streak"`
	NextReviewDate           NullableTime `json:"nextReviewDate"`
	IntervalInMinutes        int32        `json:"intervalInMinutes"`
}
type ReviewItemListResponseBody struct {
	ReviewItems []*ReviewItem `json:"reviewItems"`
}

func MapReviewItem(dbItem *db.ReviewItem) (*ReviewItem, error) {
	if dbItem == nil {
		return nil, fmt.Errorf("nil review item")
	}
	return &ReviewItem{
		ID:                       dbItem.ID,
		UserID:                   dbItem.UserID,
		SingleChoiceQuestionID:   dbItem.SingleChoiceQuestionID,
		MultipleChoiceQuestionID: dbItem.MultipleChoiceQuestionID,
		TrueOrFalseQuestionID:    dbItem.TrueOrFalseQuestionID,
		EaseFactor:               dbItem.EaseFactor,
		Difficulty:               dbItem.Difficulty,
		Streak:                   dbItem.Streak,
		NextReviewDate:           newNullableTime(dbItem.NextReviewDate.Time),
		IntervalInMinutes:        dbItem.IntervalInMinutes,
	}, nil
}
