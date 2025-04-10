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

// Execute runs the root command of the tmail CLI application.
// It handles command-line parsing and dispatching to the appropriate subcommands.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(sendCmd)
	sendCmd.Flags().BoolVar(&debugMode, "debug", false, "Enable debug mode for troubleshooting")
}
