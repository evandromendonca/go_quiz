package handlers

import (
	"fasttrack_quiz/dto"
	"fasttrack_quiz/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomQuestions(t *testing.T) {
	handler := GameHandler{}
	questions := []models.Question{
		{
			Id:            1,
			Description:   "What does RTP stand for in iGaming?",
			Options:       []string{"Real-Time Play", "Return to Player", "Random Total Payout", "Risk to Profit"},
			CorrectOption: 1,
		},
		{
			Id:            2,
			Description:   "Which of the following is a popular online slot game provider?",
			Options:       []string{"Microgaming", "Sony", "Nintendo", "Hasbro"},
			CorrectOption: 0,
		},
		{
			Id:            3,
			Description:   "What is a common term for free spins or bonus rounds in slots?",
			Options:       []string{"Betting Circles", "Win Loops", "Free Rolls", "Scatter Bonuses"},
			CorrectOption: 3,
		},
		{
			Id:            4,
			Description:   "In poker, what is the term for a sequence of five cards in numerical order?",
			Options:       []string{"Full House", "Straight", "Flush", "Pair"},
			CorrectOption: 1,
		},
		{
			Id:            5,
			Description:   "What does RNG stand for in iGaming?",
			Options:       []string{"Real Number Generator", "Random Number Generator", "Risk Number Game", "Real Network Game"},
			CorrectOption: 1,
		},
	}

	// test with num less than the length of the questions slice
	num := 3
	result := handler.getRandomQuestions(questions, num)
	if len(result) != num {
		t.Errorf("expected %d questions, got %d", num, len(result))
	}

	// test with num greater than the length of the questions slice
	num = 10
	result = handler.getRandomQuestions(questions, num)
	if len(result) != len(questions) {
		t.Errorf("expected %d questions, got %d", len(questions), len(result))
	}

	// test with num 0
	num = 0
	result = handler.getRandomQuestions(questions, num)
	if len(result) != num {
		t.Errorf("expected %d questions, got %d", num, len(result))
	}
}

func TestUpdateGameAnswers(t *testing.T) {
	handler := GameHandler{}

	user := models.User{Id: 1, Username: "user1"}
	questions := []models.Question{
		{
			Id:            1,
			Description:   "What does RTP stand for in iGaming?",
			Options:       []string{"Real-Time Play", "Return to Player", "Random Total Payout", "Risk to Profit"},
			CorrectOption: 1,
		},
		{
			Id:            2,
			Description:   "Which of the following is a popular online slot game provider?",
			Options:       []string{"Microgaming", "Sony", "Nintendo", "Hasbro"},
			CorrectOption: 0,
		},
	}

	game := models.Game{
		Id:           "game1",
		User:         user,
		NumQuestions: len(questions),
		Questions:    questions,
		CreatedDate:  time.Now(),
	}

	answersDto := []dto.AnswerDto{
		{QuestionId: 1, SelectedOption: 1}, // Correct
		{QuestionId: 2, SelectedOption: 0}, // Correct
	}

	err := handler.updateGameAnswers(&game, answersDto)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(game.Answers))
	assert.Equal(t, 100.0, game.ScorePercentage)
	for _, answer := range game.Answers {
		assert.True(t, answer.IsCorrect)
	}

	// Test with incorrect number of answers
	incorrectAnswersDto := []dto.AnswerDto{
		{QuestionId: 1, SelectedOption: 2},
	}
	err = handler.updateGameAnswers(&game, incorrectAnswersDto)
	assert.Error(t, err)
	assert.Equal(t, "number of answers 1 is different from the number of questions 2", err.Error())

	// Test with missing question in answers
	missingQuestionAnswersDto := []dto.AnswerDto{
		{QuestionId: 1, SelectedOption: 2},
		{QuestionId: 3, SelectedOption: 1},
	}
	err = handler.updateGameAnswers(&game, missingQuestionAnswersDto)
	assert.Error(t, err)
	assert.Equal(t, "answer not found for question 2", err.Error())
}

func TestRankGame(t *testing.T) {
	handler := GameHandler{}

	user1 := models.User{Id: 1, Username: "user1", Password: "pass1"}
	user2 := models.User{Id: 2, Username: "user2", Password: "pass2"}

	leaderboard := []models.LeaderboardItem{
		{User: user1, HighestScore: 90.0},
		{User: user2, HighestScore: 80.0},
	}

	// test with game1 in the middle of the leaderboard with new user
	game1 := models.Game{
		Id:              "game1",
		User:            models.User{Id: 3, Username: "user3", Password: "pass3"},
		ScorePercentage: 85.0,
	}

	position, newItem := handler.rankGame(&game1, leaderboard)

	assert.Equal(t, 1, position)
	assert.Equal(t, 3, newItem.User.Id)
	assert.Equal(t, 85.0, newItem.HighestScore)
	assert.Equal(t, 1, game1.RankingPosition)
	assert.Equal(t, 50, game1.PercentileScore)

	// test with game above the leaderboard with existing user
	game2 := models.Game{
		Id:              "game2",
		User:            user2,
		ScorePercentage: 95.0,
	}

	position2, newItem2 := handler.rankGame(&game2, leaderboard)

	assert.Equal(t, 0, position2)
	assert.Equal(t, user2, newItem2.User)
	assert.Equal(t, 95.0, newItem2.HighestScore)
	assert.Equal(t, 0, game2.RankingPosition)
	assert.Equal(t, 100, game2.PercentileScore)

	// test with game below the leaderboard with existing user
	game3 := models.Game{
		Id:              "game3",
		User:            models.User{Id: 3, Username: "new_user", Password: "secret"},
		ScorePercentage: 70.0,
	}

	position3, newItem3 := handler.rankGame(&game3, leaderboard)

	assert.Equal(t, 2, position3)
	assert.Equal(t, game3.User, newItem3.User)
	assert.Equal(t, 70.0, newItem3.HighestScore)
	assert.Equal(t, 2, game3.RankingPosition)
	assert.Equal(t, 0, game3.PercentileScore)
}
