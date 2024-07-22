package handlers

import (
	"crypto/subtle"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterHandlers(e *echo.Echo,
	userRepository UserRepositoryInterface,
	gameRepository GameRepositoryInterface,
	questionsRepository QuestionRepositoryInterface) {

	userHandler := UserHandler{
		UserRepository: userRepository,
	}

	gameHandler := GameHandler{
		GameRepository:     gameRepository,
		QuestionRepository: questionsRepository,
	}

	userGroup := e.Group("/user")
	{
		userGroup.GET("/:id", userHandler.handleGetUser)
		userGroup.POST("", userHandler.handlePostUser)
	}

	gameGroup := e.Group("/game")
	{
		gameGroup.GET("/leaderboard/:numQuestions", gameHandler.handleGetLeaderboard)
		useBasicAuth(gameGroup, userRepository)
		gameGroup.GET("", gameHandler.handleGetGame)
		gameGroup.POST("/:numQuestions", gameHandler.handlePostGame)
		gameGroup.POST("/answers", gameHandler.handlePostGameAnswers)
	}
}

func useBasicAuth(group *echo.Group, userRepository UserRepositoryInterface) {

	group.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		user, err := userRepository.GetUserByName(username)
		if err != nil {
			return false, nil
		}

		if subtle.ConstantTimeCompare([]byte(password), []byte(user.Password)) == 1 {
			c.Set("user", user)
			return true, nil
		}

		return false, nil
	}))
}
