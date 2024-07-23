package handlers

import (
	"fasttrack_quiz/dto"
	"fasttrack_quiz/models"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type GameRepositoryInterface interface {
	RepositoryInterface[models.Game]
	GetOngoingUserGame(userId int) (models.Game, error)
	ArchiveGame(ongoingGame models.Game)
	RemoveFromOngoingGame(ongoingGame models.Game)
	GetLeaderboard(numQuestions int) ([]models.LeaderboardItem, error)
	InsertLeaderboardItem(numQuestions, position int, leaderboardItem models.LeaderboardItem)
}

type QuestionRepositoryInterface interface {
	GetQuestions() ([]models.Question, error)
}

type GameHandler struct {
	GameRepository     GameRepositoryInterface
	QuestionRepository QuestionRepositoryInterface
}

func (h GameHandler) handleGetGame(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// check if there is an open game for that user, and if so, return it
	ongoingGame, err := h.GameRepository.GetOngoingUserGame(user.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	var gameDto = dto.GameDto{}
	gameDto.FromGame(ongoingGame)

	return c.JSON(http.StatusOK, gameDto)
}

func (h GameHandler) handlePostGame(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// parsing input
	param := c.Param("numQuestions")
	numQuestions, err := strconv.Atoi(param)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "number of questions must be a number")
	}

	// check if there is an open game for that user, and if so throw
	if _, err := h.GameRepository.GetOngoingUserGame(user.Id); err == nil {
		return c.JSON(http.StatusBadRequest, "there is an ongoing game, finish the game before creating a new one")
	}

	possibleQuestions, err := h.QuestionRepository.GetQuestions()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	game := models.Game{
		User:         user,
		NumQuestions: numQuestions,
		Questions:    h.getRandomQuestions(possibleQuestions, numQuestions),
	}

	h.GameRepository.Create(game)

	var gameDto = dto.GameDto{}
	gameDto.FromGame(game)

	return c.JSON(http.StatusCreated, gameDto)
}

func (h GameHandler) handlePostGameAnswers(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	var answersDto []dto.AnswerDto
	err := c.Bind(&answersDto)
	if err != nil {
		return c.String(http.StatusBadRequest, "wrong model")
	}

	// check if there is an open game for that user
	ongoingGame, err := h.GameRepository.GetOngoingUserGame(user.Id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// update game with answers
	if err = h.updateGameAnswers(&ongoingGame, answersDto); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// calcualte and update leaderboard
	numQuestions := len(ongoingGame.Questions)
	leaderboard, _ := h.GameRepository.GetLeaderboard(numQuestions)

	position, leaderboardItem := h.rankGame(&ongoingGame, leaderboard)

	h.GameRepository.InsertLeaderboardItem(numQuestions, position, leaderboardItem)

	// archive and delete from ongoing games
	h.GameRepository.ArchiveGame(ongoingGame)
	h.GameRepository.RemoveFromOngoingGame(ongoingGame)

	gameResult := dto.GameResultDto{}
	gameResult.FromGame(ongoingGame)

	return c.JSON(http.StatusOK, gameResult)
}

func (h GameHandler) handleGetLeaderboard(c echo.Context) error {
	param := c.Param("numQuestions")

	numQuestions, err := strconv.Atoi(param)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "number of questions must be a number")
	}

	leaderboard, err := h.GameRepository.GetLeaderboard(numQuestions)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	var leaderboardDto = []dto.LeaderboardDto{}
	for _, r := range leaderboard {
		leaderboardItemDto := dto.LeaderboardDto{}
		leaderboardItemDto.FromLeaderboard(r)
		leaderboardDto = append(leaderboardDto, leaderboardItemDto)
	}

	return c.JSON(http.StatusOK, leaderboardDto)
}

func (h GameHandler) getRandomQuestions(questions []models.Question, num int) []models.Question {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })

	if num > len(questions) {
		num = len(questions)
	}

	return questions[:num]
}

func (h GameHandler) updateGameAnswers(game *models.Game, answersDto []dto.AnswerDto) error {
	if len(game.Questions) != len(answersDto) {
		return fmt.Errorf("number of answers %d is different from the number of questions %d",
			len(answersDto), len(game.Questions))
	}

	var answers []models.Answer
	for _, a := range answersDto {
		answers = append(answers, a.ToAnswer())
	}

	correctCount := 0
	for _, question := range game.Questions {
		answerFound := false
		for i := range answers {
			if question.Id == answers[i].Question.Id {
				answers[i].Question = question

				if question.CorrectOption == answers[i].SelectedOption {
					answers[i].IsCorrect = true
					correctCount++
				}

				answerFound = true
				break
			}
		}

		if !answerFound {
			return fmt.Errorf("answer not found for question %d", question.Id)
		}
	}

	game.Answers = answers
	game.ScorePercentage = float64(correctCount) / float64(len(game.Questions)) * 100
	game.CompletedDate = time.Now()

	return nil
}

func (h GameHandler) rankGame(game *models.Game, leaderboard []models.LeaderboardItem) (int, models.LeaderboardItem) {
	// find the user in the board
	userPosition := -1
	for i, r := range leaderboard {
		if r.User.Id == game.User.Id {
			userPosition = i
		}
	}

	// get a board slice without the user
	leaderboardWithoutUser := leaderboard
	if userPosition != -1 {
		// if the user is in the board, we remove them
		leaderboardWithoutUser = append(leaderboard[:userPosition], leaderboard[userPosition+1:]...)
	}

	// find the position they should be inserted in the board
	i := 0
	for ; i < len(leaderboardWithoutUser); i++ {
		if game.ScorePercentage > leaderboardWithoutUser[i].HighestScore {
			break
		}
	}

	// user should be in position i now, user was better than % of all users that played
	percentile := 100
	if len(leaderboardWithoutUser) > 0 {
		percentile = int(float64(len(leaderboardWithoutUser)-i) / float64(len(leaderboardWithoutUser)) * 100)
	}

	game.RankingPosition = i
	game.PercentileScore = percentile

	return i, models.LeaderboardItem{User: game.User, HighestScore: game.ScorePercentage}
}
