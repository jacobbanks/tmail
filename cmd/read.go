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
		fmt.Println("Creating mail provider...")
		// Create mail provider
		provider, err := email.CreateDefaultMailProvider()
		if err != nil {
			fmt.Println("Error setting up mail provider:", err)
			fmt.Println("Please run: tmail auth")
			os.Exit(1)
		}
		// Fetch emails
		fmt.Println("Fetching emails...")
		limit := 10 // Default to 10 emails
		emails, err := provider.GetEmails(limit)
		if err != nil {
			fmt.Println("Error fetching emails:", err)
			os.Exit(1)
		}
		fmt.Printf("Fetched %d emails\n", len(emails))

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
