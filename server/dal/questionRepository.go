package dal

import (
	"errors"
	"fasttrack_quiz/db"
	"fasttrack_quiz/handlers"
	"fasttrack_quiz/models"
)

func GetQuestionRepository(db *db.MockDatabase) handlers.QuestionRepositoryInterface {
	questionRepository := QuestionRepository{
		db: db,
	}

	return questionRepository
}

type QuestionRepository struct {
	db *db.MockDatabase
}

func (r QuestionRepository) GetQuestions() ([]models.Question, error) {
	if len(r.db.Questions) == 0 {
		return []models.Question{}, errors.New("no questions registered")
	}

	return r.db.Questions, nil
}
