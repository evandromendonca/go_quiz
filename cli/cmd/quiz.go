/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"goquiz/dto"
	"goquiz/util"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// quizCmd represents the quiz command
var quizCmd = &cobra.Command{
	Use:   "quiz [username]",
	Short: "Start the GoQuiz game",
	Long:  `Start the GoQuiz game`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Username argument required")
			return
		}
		username := args[0]

		fmt.Printf("Type your password for %s:\n", username)
		password := util.ReadPassword()

		response, game := util.GetJsonAuth[dto.GameDto]("game", username, password)

		util.EnsureAuthorized(response.StatusCode)

		if !util.IsSuccessStatusCode(response.StatusCode) {
			game = createNewGame(username, password)
		} else {
			fmt.Println("Ongoing quiz found, please complete this one before starting another one")
		}

		answers := []dto.AnswerDto{}

		for _, q := range game.Questions {
			fmt.Printf("\nQuestion: %s\n", q.Description)
			for i, o := range q.Options {
				fmt.Printf("Option %d: %s\n", i+1, o)
			}
			optionSelected := util.GetInputNumber("Select your answer", 1, 4) - 1
			answers = append(answers, dto.AnswerDto{QuestionId: q.Id, SelectedOption: optionSelected})
		}

		gameResult := postGameAnswers(username, password, answers)

		fmt.Printf("\n%d%% of your answers where correct!\n", gameResult.ScorePercentage)
		fmt.Printf("\nYou were better than: %d%% of all users that answered %d questions. Your ranking position is %d!!!\n",
			gameResult.PercentileScore, len(gameResult.QuestionsResults), gameResult.RankingPosition+1)
	},
}

func init() {
	rootCmd.AddCommand(quizCmd)
}

func createNewGame(username, password string) dto.GameDto {
	numQuestions := util.GetInputNumber("No ongoing quiz found, creating a new one. Select the number of questions", 1, 15)

	response, game := util.PostJsonAuth[dto.GameDto](fmt.Sprintf("game/%d", numQuestions), username, password, nil)

	if !util.IsSuccessStatusCode(response.StatusCode) {
		body, _ := io.ReadAll(response.Body)
		fmt.Printf("\n%d: error creating game. %s\n", response.StatusCode, string(body))
		os.Exit(1)
	}

	return game
}

func postGameAnswers(username, password string, answers []dto.AnswerDto) dto.GameResultDto {
	response, gameResult := util.PostJsonAuth[dto.GameResultDto]("game/answers", username, password, answers)

	if !util.IsSuccessStatusCode(response.StatusCode) {
		body, _ := io.ReadAll(response.Body)
		fmt.Printf("\n%d: error saving answers. %s\n", response.StatusCode, string(body))
		os.Exit(1)
	}

	return gameResult
}
