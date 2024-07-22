package dal

import (
	"errors"
	"fasttrack_quiz/db"
	"fasttrack_quiz/handlers"
	"fasttrack_quiz/models"
	"strings"
)

func GetUserRepository(db *db.MockDatabase) handlers.UserRepositoryInterface {
	userRepository := UserRepository{
		db: db,
	}

	return userRepository
}

type UserRepository struct {
	db *db.MockDatabase
}

func (r UserRepository) Create(user models.User) (models.User, error) {
	for _, u := range r.db.Users {
		if strings.EqualFold(u.Username, user.Username) {
			return models.User{}, errors.New("username already taken")
		}
	}

	newUser := models.User{
		Id:       len(r.db.Users), // calculate the id
		Username: user.Username,
		Password: user.Password,
	}

	r.db.Users[newUser.Id] = newUser

	return newUser, nil
}

func (r UserRepository) ReadById(id int) (models.User, error) {
	user, ok := r.db.Users[id]
	if !ok {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}

func (r UserRepository) GetUserByName(username string) (models.User, error) {
	for _, u := range r.db.Users {
		if strings.EqualFold(username, u.Username) {
			return u, nil
		}
	}

	return models.User{}, errors.New("not found")
}
