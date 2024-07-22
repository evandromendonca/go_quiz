package dto

type LeaderboardDto struct {
	Username     string  `json:"username"`
	HighestScore float64 `json:"highestScore"`
}
