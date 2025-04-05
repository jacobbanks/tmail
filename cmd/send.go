package cmd

import (
	"fmt"
	"os"

	"github.com/jacobbanks/tmail/auth"
	"github.com/jacobbanks/tmail/ui"
	"github.com/spf13/cobra"
)

// Flag for enabling debug mode
var debugMode bool

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Compose and send an email",
	Long:  "Open a TUI to compose and send an email",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := auth.LoadUser()
		if err != nil {
			fmt.Println("You need to set up your email credentials first.")
			fmt.Println("Please run: tmail auth")
			os.Exit(1)
		}

		composer := ui.NewEmailComposer(nil)
		composer.SetDebugMode(debugMode)
		if err := composer.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	// Add flags
	sendCmd.Flags().BoolVar(&debugMode, "debug", false, "Enable debug mode for troubleshooting")
}
