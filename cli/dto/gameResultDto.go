package dto

import (
	"time"
)

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
