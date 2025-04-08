package cmd

import (
	"fmt"
	"os"

	"github.com/jacobbanks/tmail/email"
	"github.com/jacobbanks/tmail/ui"
	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read emails",
	Long:  "Fetch and read emails in a terminal UI",
	Run: func(cmd *cobra.Command, args []string) {
		provider, err := email.CreateDefaultMailProvider()
		if err != nil {
			fmt.Println("Error setting up mail provider:", err)
			fmt.Println("Please run: tmail auth")
			os.Exit(1)
		}

		emails, err := provider.GetEmails(25)
		if err != nil {
			fmt.Println("Error fetching emails:", err)
			os.Exit(1)
		}

		if len(emails) == 0 {
			fmt.Println("No emails found.")
			return
		}

		reader := ui.NewEmailReader(emails, provider)
		if err := reader.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}
