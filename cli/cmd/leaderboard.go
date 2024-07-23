/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"goquiz/dto"
	"goquiz/util"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// leaderboardCmd represents the ranking command
var leaderboardCmd = &cobra.Command{
	Use:   "leaderboard [numQuestions]",
	Short: "Get the GoQuiz leaderboard",
	Long: `Get to GoQuiz leaderboard sorted by the better to the worst player.
To be a fair competition, leaderboards are divided by number of questions. Therefore,
you should select which leaderboard you want to see specifying how many questions the quizzes had`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("This command requires the number of questions input")
			os.Exit(1)
		}

		numQuestions, err := strconv.Atoi(args[0])
		if err != nil || numQuestions < 1 {
			fmt.Println("Number of questions must be a number greater than 0")
			os.Exit(1)
		}

		response, leaderboard := util.GetJson[[]dto.LeaderboardDto](fmt.Sprintf("game/leaderboard/%d", numQuestions))

		if !util.IsSuccessStatusCode(response.StatusCode) {
			if response.StatusCode == http.StatusNotFound {
				fmt.Printf("Leaderboard for quizzes with %d questions not found\n", numQuestions)

			} else {
				fmt.Println("Error getting leaderboad")
			}
			os.Exit(1)
		}

		fmt.Printf("\nLeaderboard for quizzes with %d questions!!!\n\n", numQuestions)
		for i, item := range leaderboard {
			fmt.Printf("Position %d: %s (%.1f%%)\n", i+1, item.Username, item.HighestScore)
		}
	},
}

func init() {
	rootCmd.AddCommand(leaderboardCmd)
}
