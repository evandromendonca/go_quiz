package dto

import "fasttrack_quiz/models"

type LeaderboardDto struct {
	Username     string  `json:"username"`
	HighestScore float64 `json:"highestScore"`
}

func (l *LeaderboardDto) FromLeaderboard(leaderboardItem models.LeaderboardItem) {
	l.Username = leaderboardItem.User.Username
	l.HighestScore = leaderboardItem.HighestScore
}
