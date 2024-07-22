package cmd

import (
	"fasttrack_quiz_cli/dto"
	"fasttrack_quiz_cli/util"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [username]",
	Short: "Creates a new user for the quiz game",
	Long:  "Creates a new user for the quiz game",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Username argument required")
			return
		}
		username := args[0]

		fmt.Printf("Choose a password for: %s\n", username)
		password := util.ReadPassword()

		userDto := dto.UserDto{
			Username: username,
			Password: password,
		}

		resp := util.PostJson("user", userDto)
		defer resp.Body.Close()

		// read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		if resp.StatusCode == http.StatusCreated {
			fmt.Printf("User %s successfully created\n", username)
			return
		}

		fmt.Printf("Error creating user %s\n", username)
		fmt.Printf("%d: %s", resp.StatusCode, string(body))
	},
}

func init() {
	userCmd.AddCommand(newCmd)
}
