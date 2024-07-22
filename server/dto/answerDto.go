package dto

import (
	"fasttrack_quiz/models"
)

type AnswerDto struct {
	QuestionId     int `json:"questionId"`
	SelectedOption int `json:"selectedOption"`
}

func (o *AnswerDto) ToAnswer() (answer models.Answer) {
	answer = models.Answer{
		Question: models.Question{
			Id: o.QuestionId,
		},
		SelectedOption: o.SelectedOption,
	}

	return answer
}
