package dto

import (
	"time"
)

type GameDto struct {
	Id          string        `json:"id"`
	Username    string        `json:"username"`
	Questions   []QuestionDto `json:"questions"`
	CreatedDate time.Time     `json:"createdDate"`
}
