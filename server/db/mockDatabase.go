package db

import (
	"fasttrack_quiz/models"
	"sync"
)

type MockDatabase struct {
	Users        map[int]models.User
	PastGames    []models.Game
	OngoingGames map[int]models.Game
	Questions    []models.Question
	Leaderboards map[int][]models.LeaderboardItem
}

func (m *MockDatabase) LoadMockedData() {
	var once sync.Once

	once.Do(func() {
		user1 := models.User{
			Id:       0,
			Username: "user0",
			Password: "pass0",
		}

		user2 := models.User{
			Id:       1,
			Username: "user1",
			Password: "pass1",
		}

		user3 := models.User{
			Id:       2,
			Username: "user2",
			Password: "pass2",
		}

		m.Users = map[int]models.User{
			user1.Id: user1,
			user2.Id: user2,
			user3.Id: user3,
		}

		m.Questions = createQuizQuestions()
		m.PastGames = []models.Game{}
		m.OngoingGames = map[int]models.Game{}
		m.Leaderboards = map[int][]models.LeaderboardItem{
			4: {
				{
					User:         user1,
					HighestScore: 75,
				},
				{
					User:         user2,
					HighestScore: 50,
				},
				{
					User:         user3,
					HighestScore: 25,
				},
			},
		}
	})
}

func createQuizQuestions() []models.Question {
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
		{
			Id:            6,
			Description:   "Which game is commonly known as '21'?",
			Options:       []string{"Poker", "Baccarat", "Roulette", "Blackjack"},
			CorrectOption: 3,
		},
		{
			Id:            7,
			Description:   "What is the maximum number of paylines commonly found in modern video slots?",
			Options:       []string{"10", "25", "50", "243"},
			CorrectOption: 3,
		},
		{
			Id:            8,
			Description:   "Which of the following is a popular type of bet in roulette?",
			Options:       []string{"Straight", "Fold", "Full House", "Double Down"},
			CorrectOption: 0,
		},
		{
			Id:            9,
			Description:   "What is the name of the slot machine feature that multiplies winnings?",
			Options:       []string{"Scatter", "Wild", "Multiplier", "Bonus Round"},
			CorrectOption: 2,
		},
		{
			Id:            10,
			Description:   "In Texas Hold'em poker, what are the first three community cards called?",
			Options:       []string{"The River", "The Turn", "The Flop", "The Deal"},
			CorrectOption: 2,
		},
		{
			Id:            11,
			Description:   "Which game involves a wheel and a ball?",
			Options:       []string{"Blackjack", "Poker", "Roulette", "Baccarat"},
			CorrectOption: 2,
		},
		{
			Id:            12,
			Description:   "In slots, what is a 'wild' symbol typically used for?",
			Options:       []string{"To create paylines", "To trigger bonus rounds", "To multiply winnings", "To replace other symbols"},
			CorrectOption: 3,
		},
		{
			Id:            13,
			Description:   "What does the term 'all-in' mean in poker?",
			Options:       []string{"Folding all cards", "Raising the minimum bet", "Betting all of one's chips", "Splitting the pot"},
			CorrectOption: 2,
		},
		{
			Id:            14,
			Description:   "Which of the following is a progressive jackpot?",
			Options:       []string{"A fixed payout", "A jackpot that increases with each bet", "A bonus round reward", "A free spin reward"},
			CorrectOption: 1,
		},
		{
			Id:            15,
			Description:   "What is the house edge?",
			Options:       []string{"Player's advantage over the house", "The casino's advantage over players", "The edge of the casino table", "The payout ratio"},
			CorrectOption: 1,
		},
	}
	return questions
}
