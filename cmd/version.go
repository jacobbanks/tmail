package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Build information - to be overridden during build
var (
	Version   string
	BuildDate string
	GitCommit string
	GitState  string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show tmail version information",
	Long:  `Display version, build date, git commit, and git state for tmail.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tmail - Terminal Mail Client")
		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Git State:  %s\n", GitState)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}