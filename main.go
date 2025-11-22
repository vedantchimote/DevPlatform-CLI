package main

import (
	"github.com/devplatform/devplatform-cli/cmd"
)

var (
	// Version information - will be set during build
	Version   = "dev"
	GitCommit = "none"
	BuildDate = "unknown"
)

func main() {
	// Set version information in cmd package
	cmd.Version = Version
	cmd.GitCommit = GitCommit
	cmd.BuildDate = BuildDate
	
	// Execute the root command
	cmd.Execute()
}
