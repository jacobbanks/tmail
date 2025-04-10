package cmd

import (
	"fmt"
	"strconv"

	"github.com/jacobbanks/tmail/email"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure TMail settings",
	Long: `Configure TMail settings and preferences.
Examples:
  tmail config show
  tmail config set theme blue
  tmail config set default_mails 25`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		switch args[0] {
		case "show":
			showConfig()
		case "set":
			if len(args) < 3 {
				fmt.Println("Error: Not enough arguments for set command")
				fmt.Println("Usage: tmail config set [setting] [value]")
				return
			}
			setSetting(args[1], args[2])
		default:
			fmt.Printf("Unknown config command: %s\n", args[0])
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func showConfig() {
	config, err := email.LoadUserConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return
	}

	fmt.Println("TMail Configuration:")
	fmt.Println("--------------------")
	fmt.Printf("Theme: %s\n", config.Theme)
	fmt.Printf("Default emails to fetch: %d\n", config.DefaultNumMails)
}

func setSetting(setting, value string) {
	config, err := email.LoadUserConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return
	}

	switch setting {
	case "theme":
		if value != "blue" && value != "dark" && value != "light" {
			fmt.Println("Invalid theme. Valid options: blue, dark, light")
			return
		}
		config.Theme = value
		fmt.Printf("Theme set to: %s\n", value)

	case "default_mails":
		num, err := strconv.Atoi(value)
		if err != nil || num < 1 || num > 500 {
			fmt.Println("Invalid number. Please specify a number between 1 and 500")
			return
		}
		config.DefaultNumMails = num
		fmt.Printf("Default emails to fetch set to: %d\n", num)
	default:
		fmt.Printf("Unknown setting: %s\n", setting)
		fmt.Println("Valid settings: theme, default_mails")
		return
	}

	if err := email.SaveUserConfig(config); err != nil {
		fmt.Printf("Error saving configuration: %v\n", err)
		return
	}
}
