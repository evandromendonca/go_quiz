/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fasttrack_quiz_cli/util"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new user for the quiz game",
	Long:  "Creates a new user for the quiz game",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("username required")
			return
		}

		username := args[0]

		fmt.Println("choose a password for ", username)
		password := util.ReadPassword()

		postBody, _ := json.Marshal(map[string]string{
			"username": username,
			"password": password,
		})
		responseBody := bytes.NewBuffer(postBody)

		resp, err := http.Post("http://localhost:8080/user", "application/json", responseBody)
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}
		defer resp.Body.Close()

		// read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		sb := string(body)

		log.Println(resp.Status)
		log.Println(sb)
	},
}

func init() {
	userCmd.AddCommand(newCmd)
}
