package dal

import (
	"errors"
	"fasttrack_quiz/db"
	"fasttrack_quiz/handlers"
	"fasttrack_quiz/models"
	"fmt"
	"time"
)

func GetGameRepository(db *db.MockDatabase) handlers.GameRepositoryInterface {
	gameRepository := GameRepository{
		db: db,
	}

	return gameRepository
}

type GameRepository struct {
	db *db.MockDatabase
}

func (r GameRepository) Create(game models.Game) (models.Game, error) {
	// check if the user exists and gets its reference
	user, ok := r.db.Users[game.User.Id]
	if !ok {
		return models.Game{}, errors.New("user not found to start a game")
	}

	// updating fields
	game.User = user
	game.CreatedDate = time.Now()
	game.Id = fmt.Sprintf("%d%d", user.Id, game.CreatedDate.Unix())

	r.db.OngoingGames[game.User.Id] = game

	return game, nil
}

func (r GameRepository) GetOngoingUserGame(userId int) (models.Game, error) {
	// search the user's ongoing games
	game, ok := r.db.OngoingGames[userId]
	if !ok {
		return models.Game{}, errors.New("no ongoing game found")
	}

	return game, nil
}

func (r GameRepository) RemoveFromOngoingGame(game models.Game) {
	delete(r.db.OngoingGames, game.User.Id)
}

func (r GameRepository) ArchiveGame(game models.Game) {
	r.db.PastGames = append(r.db.PastGames, game)
}

func (r GameRepository) GetLeaderboard(numQuestions int) ([]models.LeaderboardItem, error) {
	leaderboard, ok := r.db.Leaderboards[numQuestions]
	if !ok {
		return []models.LeaderboardItem{}, errors.New("leaderboard for that number of questions not found")
	}

	return leaderboard, nil
}

func (r GameRepository) InsertLeaderboardItem(numQuestions, position int, leaderboardItem models.LeaderboardItem) {
	leaderboard, ok := r.db.Leaderboards[numQuestions]
	if !ok {
		leaderboard = []models.LeaderboardItem{}
		r.db.Leaderboards[numQuestions] = leaderboard
	}

	// remove the user from the leaderboard
	leaderboardWithoutUser := []models.LeaderboardItem{}
	for i := 0; i < len(leaderboard); i++ {
		if leaderboard[i].User.Id != leaderboardItem.User.Id {
			leaderboardWithoutUser = append(leaderboardWithoutUser, leaderboard[i])
		}
	}

	// insert in the correct position of the slice
	r.db.Leaderboards[numQuestions] = append(leaderboardWithoutUser[:position],
		append([]models.LeaderboardItem{leaderboardItem}, leaderboardWithoutUser[position:]...)...)
}
