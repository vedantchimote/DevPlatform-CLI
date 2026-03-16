package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Global flags
	cfgFile  string
	verbose  bool
	debug    bool
	noColor  bool
	
	// Version information - set during build
	Version   = "dev"
	GitCommit = "none"
	BuildDate = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "devplatform",
	Short: "DevPlatform CLI - Self-service infrastructure provisioning",
	Long: `DevPlatform CLI is an Internal Developer Platform (IDP) tool that enables 
developers to self-service provision complete, isolated infrastructure 
environments on AWS or Azure.

Provision VPC/VNet, RDS/Azure Database, EKS/AKS namespaces, and deploy 
applications with a single command - reducing provisioning time from 
2 days to 3 minutes.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .devplatform.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug output (includes API calls)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")
}
