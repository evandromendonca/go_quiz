package dto

type QuestionResultDto struct {
	Question       QuestionDto `json:"question"`
	CorrectOption  int         `json:"correctionOption"`
	SelectedOption int         `json:"selectedOption"`
	IsCorrect      bool        `json:"isCorrect"`
}
