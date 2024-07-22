/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// rankingCmd represents the ranking command
var rankingCmd = &cobra.Command{
	Use:   "ranking",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		url := fmt.Sprintf("http://localhost:8080/game/leaderboard/%d", numQuestions)

		response, err := http.Get(url)
		if err != nil {
			fmt.Println("error in request:", err)
			os.Exit(1)
		}

		resBody, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("client: could not read response body: %s\n", err)
			os.Exit(1)
		}

		fmt.Println(string(resBody))
	},
}

func init() {
	rootCmd.AddCommand(rankingCmd)
}
