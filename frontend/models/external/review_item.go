package external

import (
	"fmt"
	"spaced-ace/models"
	"spaced-ace/models/business"
	"time"
)

type ReviewItem struct {
	ID                       string              `json:"id"`
	UserID                   string              `json:"userID"`
	QuizID                   string              `json:"quizID"`
	QuizName                 string              `json:"quizName"`
	SingleChoiceQuestionID   *string             `json:"singleChoiceQuestionID"`
	MultipleChoiceQuestionID *string             `json:"multipleChoiceQuestionID"`
	TrueOrFalseQuestionID    *string             `json:"trueOrFalseQuestionID"`
	QuestionName             string              `json:"questionName"`
	EaseFactor               float64             `json:"easeFactor"`
	Difficulty               float64             `json:"difficulty"`
	Streak                   int32               `json:"streak"`
	NextReviewDate           models.NullableTime `json:"nextReviewDate"`
	IntervalInMinutes        int32               `json:"intervalInMinutes"`
}
type ReviewItemsRequestBody struct {
	QuizID     string `json:"quiz"`
	Difficulty string `json:"difficulty"`
	Status     string `json:"status"`
	Page       int    `json:"page"`
	Query      string `json:"query"`
}
type ReviewItemResponseBody struct {
	ReviewItems              []*ReviewItem `json:"reviewItems"`
	ReviewItemCountForFilter int           `json:"reviewItemCountForFilter"`
}
type ReviewItemCountsResponseBody struct {
	Total       int `json:"total"`
	DueToReview int `json:"dueToReview"`
}
type ReviewItemPageDataResponseBody struct {
	CurrentReviewItemID    string                              `json:"currentReviewItemID"`
	SingleChoiceQuestion   *SingleChoiceQuestionResponseBody   `json:"singleChoiceQuestion"`
	MultipleChoiceQuestion *MultipleChoiceQuestionResponseBody `json:"multipleChoiceQuestion"`
	TrueOrFalseQuestion    *TrueOrFalseQuestionResponseBody    `json:"trueOrFalseQuestion"`
	NextReviewItemID       string                              `json:"nextReviewItemID"`
}
type SubmitReviewItemQuestionRequestBody struct {
	SingleChoiceValue   string   `json:"singleChoiceValue"`
	MultipleChoiceValue []string `json:"multipleChoiceValue"`
	TrueOrFalseValue    bool     `json:"trueOrFalseValue"`
}

func (r *ReviewItem) MapToBusiness() (*business.ReviewItem, error) {
	if r == nil {
		return nil, fmt.Errorf("nil review item")
	}

	var questionID string
	if r.SingleChoiceQuestionID != nil {
		questionID = *r.SingleChoiceQuestionID
	} else if r.MultipleChoiceQuestionID != nil {
		questionID = *r.MultipleChoiceQuestionID
	} else if r.TrueOrFalseQuestionID != nil {
		questionID = *r.TrueOrFalseQuestionID
	}
	if questionID == "" {
		return nil, fmt.Errorf("nil question ID")
	}

	return &business.ReviewItem{
		ID:           r.ID,
		QuizName:     r.QuizName,
		QuestionName: r.QuestionName,
		QuestionID:   questionID,
		NextReview:   r.NextReviewDate.Time,
		Difficulty:   r.Difficulty,
		Streak:       int(r.Streak),
		NeedToReview: r.NextReviewDate.Time.Before(time.Now()),
	}, nil
}
func (r *ReviewItemPageDataResponseBody) MapToBusiness() (*business.ReviewItemPageData, error) {
	var singleChoiceQuestion *business.SingleChoiceQuestion
	if r.SingleChoiceQuestion != nil {
		var err error
		singleChoiceQuestion, err = r.SingleChoiceQuestion.MapToBusiness()
		if err != nil {
			return nil, err
		}
	}

	var multipleChoiceQuestion *business.MultipleChoiceQuestion
	if r.MultipleChoiceQuestion != nil {
		var err error
		multipleChoiceQuestion, err = r.MultipleChoiceQuestion.MapToBusiness()
		if err != nil {
			return nil, err
		}
	}

	var trueOrFalseQuestion *business.TrueOrFalseQuestion
	if r.TrueOrFalseQuestion != nil {
		var err error
		trueOrFalseQuestion, err = r.TrueOrFalseQuestion.MapToBusiness()
		if err != nil {
			return nil, err
		}
	}

	return &business.ReviewItemPageData{
		CurrentReviewItemID:    r.CurrentReviewItemID,
		SingleChoiceQuestion:   singleChoiceQuestion,
		MultipleChoiceQuestion: multipleChoiceQuestion,
		TrueOrFalseQuestion:    trueOrFalseQuestion,
		NextReviewItemID:       r.NextReviewItemID,
	}, nil
}
