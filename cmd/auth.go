package cmd

import (
	"fmt"
	"os"

	"github.com/jacobbanks/tmail/auth"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Set up email authentication",
	Long:  "Configure your Gmail account credentials",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := auth.LoadUser()
		if err == nil {
			fmt.Print("Existing credentials found. Do you want to overwrite them? (y/n): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "yes" {
				fmt.Println("Setup canceled. Using existing credentials.")
				return
			}
		}

		err = auth.PromptForAuthentication()
		if err != nil {
			fmt.Printf("Authentication setup failed: %v\n", err)
			os.Exit(1)
		}
	},
}
