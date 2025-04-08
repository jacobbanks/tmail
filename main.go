// tmail is a simple email client for the terminal
package main

import (
	"github.com/jacobbanks/tmail/cmd"
)

// Version information variables - set during build via -ldflags
var (
	Version   = "development"
	BuildDate = "unknown"
	GitCommit = "unknown"
	GitState  = "unknown"
)

func main() {
	// Pass version information to the command package
	cmd.Version = Version
	cmd.BuildDate = BuildDate
	cmd.GitCommit = GitCommit
	cmd.GitState = GitState
	
	cmd.Execute()
}
