package models

type Answer struct {
	Question       Question
	SelectedOption int
	IsCorrect      bool
}
