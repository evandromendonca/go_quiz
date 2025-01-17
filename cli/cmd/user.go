package cmd

import (
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage the quiz game users",
	Long:  "Manage the quiz game users",
}

func init() {
	rootCmd.AddCommand(userCmd)
}
