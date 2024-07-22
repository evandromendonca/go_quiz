package dto

type QuestionResultDto struct {
	Question       QuestionDto
	CorrectOption  int
	SelectedOption int
	IsCorrect      bool
}
