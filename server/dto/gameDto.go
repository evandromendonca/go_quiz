package dto

import (
	"fasttrack_quiz/models"
	"time"
)

type GameDto struct {
	Id          string        `json:"id"`
	Username    string        `json:"username"`
	Questions   []QuestionDto `json:"questions"`
	CreatedDate time.Time     `json:"createdDate"`
}

func (o *GameDto) FromGame(game models.Game) {
	o.Id = game.Id
	o.Username = game.User.Username
	o.CreatedDate = game.CreatedDate
	o.Questions = []QuestionDto{}

	for _, q := range game.Questions {
		o.Questions = append(o.Questions, QuestionDto{
			Id:          q.Id,
			Description: q.Description,
			Options:     q.Options,
		})
	}
}
