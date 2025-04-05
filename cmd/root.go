package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tmail",
	Short: "A simple CLI for email",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to tmail. Use 'tmail help' for commands.")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(sendCmd)
	// Auth command is registered in auth.go
}
