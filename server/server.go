package main

import (
	"fasttrack_quiz/dal"
	"fasttrack_quiz/db"
	"fasttrack_quiz/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db := new(db.MockDatabase)
	db.LoadMockedData()

	handlers.RegisterHandlers(e,
		dal.GetUserRepository(db),
		dal.GetGameRepository(db),
		dal.GetQuestionRepository(db))

	e.Logger.Fatal(e.Start(":8080"))
}
