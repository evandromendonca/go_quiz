package dto

type QuestionDto struct {
	Id          int      `json:"id"`
	Description string   `json:"description"`
	Options     []string `json:"options"`
}
