package handlers

import (
	"fasttrack_quiz/dto"
	"fasttrack_quiz/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type UserRepositoryInterface interface {
	RepositoryInterface[models.User]
	GetUserByName(username string) (models.User, error)
	ReadById(id int) (models.User, error)
}

type UserHandler struct {
	UserRepository UserRepositoryInterface
}

func (h UserHandler) handleGetUser(c echo.Context) error {
	// parsing input
	param := c.Param("id")
	userId, err := strconv.Atoi(param)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "UserId must be a number")
	}

	// fetching user
	user, err := h.UserRepository.ReadById(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, user)
}

func (h UserHandler) handlePostUser(c echo.Context) error {
	// parsing input
	var newUserDto dto.NewUserDto
	err := c.Bind(&newUserDto)
	if err != nil {
		return c.String(http.StatusBadRequest, "wrong model")
	}

	// validate input
	if strings.Contains(newUserDto.Username, " ") || strings.Contains(newUserDto.Password, " ") {
		return c.String(http.StatusBadRequest, "username and password cannot contain whitespaces")
	}

	// check if username is already in use
	if _, err := h.UserRepository.GetUserByName(newUserDto.Username); err == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "username already in use")
	}

	// create new user
	h.UserRepository.Create(models.User{
		Username: newUserDto.Username,
		Password: newUserDto.Password,
	})

	return c.JSON(http.StatusCreated, nil)
}
