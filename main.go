// tmail is a simple email client for the terminal
package main

import (
	"github.com/jacobbanks/tmail/cmd"
)

func main() {
	// Execute the root command which will handle subcommands
	cmd.Execute()
}
