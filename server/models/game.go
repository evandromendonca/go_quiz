package models

import "time"

type Game struct {
	Id              string
	User            User
	NumQuestions    int
	Questions       []Question
	Answers         []Answer
	CreatedDate     time.Time
	CompletedDate   time.Time
	ScorePercentage float64
	RankingPosition int
	PercentileScore int
}
