package dto

type AnswerDto struct {
	QuestionId     int `json:"questionId"`
	SelectedOption int `json:"selectedOption"`
}
