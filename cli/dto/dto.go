package dto

import "time"

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
