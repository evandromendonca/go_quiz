package models

type Question struct {
	Id            int
	Description   string
	Options       []string
	CorrectOption int
}
