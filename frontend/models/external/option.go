package external

import "spaced-ace/models/business"

type Option struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type QuizOptionsResponseBody struct {
	QuizOptions []*Option `json:"quizOptions"`
}

func (o *Option) MapToBusiness() (option *business.Option, err error) {
	return &business.Option{
		Name:  o.Name,
		Value: o.Value,
	}, nil
}
