package dto

import (
	"fasttrack_quiz/models"
	"time"
)

type GameResultDto struct {
	Id               string              `json:"id"`
	Username         string              `json:"username"`
	QuestionsResults []QuestionResultDto `json:"questionsResults"`
	ScorePercentage  float64             `json:"scorePercentage"`
	CreatedDate      time.Time           `json:"createdDate"`
	CompletedDate    time.Time           `json:"completedDate"`
	PercentileScore  int                 `json:"percentileScore"`
	RankingPosition  int                 `json:"rankingPosition"`
}

func (o *GameResultDto) FromGame(game models.Game) {
	o.Id = game.Id
	o.Username = game.User.Username
	o.CreatedDate = game.CreatedDate
	o.CompletedDate = game.CompletedDate
	o.ScorePercentage = game.ScorePercentage
	o.QuestionsResults = []QuestionResultDto{}
	o.PercentileScore = game.PercentileScore
	o.RankingPosition = game.RankingPosition

	for _, a := range game.Answers {
		o.QuestionsResults = append(o.QuestionsResults, QuestionResultDto{
			Question: QuestionDto{
				Id:          a.Question.Id,
				Description: a.Question.Description,
				Options:     a.Question.Options,
			},
			CorrectOption:  a.Question.CorrectOption,
			SelectedOption: a.SelectedOption,
			IsCorrect:      a.IsCorrect,
		})
	}
}
