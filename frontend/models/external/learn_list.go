package external

import "spaced-ace/models/business"

type LearnListItem struct {
	QuizID   string `json:"quizID"`
	QuizName string `json:"quizName"`
}
type LearnListResponseBody struct {
	AvailableItems []LearnListItem `json:"availableItems"`
	SelectedItems  []LearnListItem `json:"selectedItems"`
}

func (r LearnListItem) MapToBusiness() (*business.LearnListItem, error) {
	return &business.LearnListItem{
		QuizID:   r.QuizID,
		QuizName: r.QuizName,
	}, nil
}
func (r LearnListResponseBody) MapToBusiness() (*business.LearnList, error) {
	availableItems := make([]business.LearnListItem, 0, len(r.AvailableItems))
	for _, item := range r.AvailableItems {
		businessItem, err := item.MapToBusiness()
		if err != nil {
			return nil, err
		}
		availableItems = append(availableItems, *businessItem)
	}

	selectedItems := make([]business.LearnListItem, 0, len(r.SelectedItems))
	for _, item := range r.SelectedItems {
		businessItem, err := item.MapToBusiness()
		if err != nil {
			return nil, err
		}
		selectedItems = append(selectedItems, *businessItem)
	}

	return &business.LearnList{
		AvailableItems: availableItems,
		SelectedItems:  selectedItems,
	}, nil
}
