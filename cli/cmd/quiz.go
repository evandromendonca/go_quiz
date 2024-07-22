/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fasttrack_quiz_cli/dto"
	"fasttrack_quiz_cli/util"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// quizCmd represents the quiz command
var quizCmd = &cobra.Command{
	Use:   "quiz",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("quiz called")

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("what's your username")
		scanner.Scan()

		username := scanner.Text()

		fmt.Println("username:", username)

		fmt.Println("type your password:")
		password := util.ReadPassword()

		fmt.Println("password:", password)

		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/game", nil)
		if err != nil {
			fmt.Println("error building request:", err)
			os.Exit(1)
		}

		req.SetBasicAuth(username, password)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("error calling GET game:", err)
			os.Exit(1)
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("client: could not read response body: %s\n", err)
			os.Exit(1)
		}

		fmt.Println(string(resBody))

		var game dto.GameDto
		json.Unmarshal(resBody, &game)

		answers := []dto.AnswerDto{}

		for _, q := range game.Questions {
			fmt.Println("Question:", q.Description)
			for i, o := range q.Options {
				fmt.Println("Option", i, ":", o)
			}

			var optionSelected int = 0
			for {
				fmt.Println("Select your answer (option number [0..3]):")
				scanner.Scan()
				optionSelected, err = strconv.Atoi(scanner.Text())
				if err != nil || optionSelected < 0 || optionSelected > 3 {
					fmt.Println("Option must be a number between 0 to 3")
					continue
				}
				fmt.Println("Your options is:", optionSelected)
				break
			}

			answers = append(answers, dto.AnswerDto{QuestionId: q.Id, SelectedOption: optionSelected})
		}

		gameResult := postGameAnswers(username, password, answers)

		fmt.Printf("%d%% of your answers where correct!\n", gameResult.ScorePercentage)
		fmt.Printf("You were better than: %d%% of all users that answered %d questions. Your ranking position is %d!!!\n",
			gameResult.PercentileScore, len(gameResult.QuestionsResults), gameResult.RankingPosition+1)
	},
}

func init() {
	rootCmd.AddCommand(quizCmd)
}

func postGameAnswers(username, password string, answers []dto.AnswerDto) dto.GameResultDto {
	postBody, _ := json.Marshal(answers)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/game/answers", bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println("error building request:", err)
		os.Exit(1)
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error calling POST game answers:", err)
		os.Exit(1)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(resBody))

	var gameResult dto.GameResultDto
	json.Unmarshal(resBody, &gameResult)

	return gameResult
}
