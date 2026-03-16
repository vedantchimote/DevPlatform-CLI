package main

import (
	"fmt"
	"os"

	"github.com/devplatform/devplatform-cli/cmd"
	clierrors "github.com/devplatform/devplatform-cli/internal/errors"
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
	if err := cmd.Execute(); err != nil {
		// Check if it's a CLIError and format it nicely
		if cliErr, ok := err.(*clierrors.CLIError); ok {
			fmt.Fprintln(os.Stderr, cliErr.Format())
		} else {
			// For non-CLIError, just print the error
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}
