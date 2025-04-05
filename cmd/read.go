package cmd

import (
	"fmt"
	"os"

	"github.com/jacobbanks/tmail/auth"
	"github.com/jacobbanks/tmail/email"
	"github.com/jacobbanks/tmail/ui"
	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read emails",
	Long:  "Fetch and read emails in a terminal UI",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := auth.LoadUser()
		if err != nil {
			fmt.Println("You need to set up your email credentials first.")
			fmt.Println("Please run: tmail auth")
			os.Exit(1)
		}
		
		emails := email.FetchEmails(10) // Fetch the last 10 emails

		if len(emails) == 0 {
			fmt.Println("No emails found.")
			return
		}

		reader := ui.NewEmailReader(emails)
		if err := reader.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}
