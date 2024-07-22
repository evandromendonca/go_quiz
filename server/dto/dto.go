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

type QuestionDto struct {
	Id          int      `json:"id"`
	Description string   `json:"description"`
	Options     []string `json:"options"`
}

type AnswerDto struct {
	QuestionId     int `json:"questionId"`
	SelectedOption int `json:"selectedOption"`
}

type GameResultDto struct {
	Id               string              `json:"id"`
	Username         string              `json:"username"`
	QuestionsResults []QuestionResultDto `json:"questionsResults"`
	ScorePercentage  int                 `json:"scorePercentage"`
	CreatedDate      time.Time           `json:"createdDate"`
	CompletedDate    time.Time           `json:"completedDate"`
	PercentileScore  int                 `json:"percentileScore"`
	RankingPosition  int                 `json:"rankingPosition"`
}

type QuestionResultDto struct {
	Question       QuestionDto
	CorrectOption  int
	SelectedOption int
	IsCorrect      bool
}

type NewUserDto struct {
	Username string
	Password string
}

type LeaderboardDto struct {
	Username     string
	HighestScore int
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

func (o *AnswerDto) ToAnswer() (answer models.Answer) {
	answer = models.Answer{
		Question: models.Question{
			Id: o.QuestionId,
		},
		SelectedOption: o.SelectedOption,
	}

	return answer
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
