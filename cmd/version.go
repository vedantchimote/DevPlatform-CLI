package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	short     bool
	checkDeps bool
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long: `Display version information for the DevPlatform CLI and optionally check 
the versions of required dependencies (terraform, helm, kubectl, aws CLI, az CLI).`,
	Run: runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	
	versionCmd.Flags().BoolVar(&short, "short", false, "display only the version number")
	versionCmd.Flags().BoolVar(&checkDeps, "check-deps", false, "check versions of required dependencies")
}

func runVersion(cmd *cobra.Command, args []string) {
	if short {
		fmt.Println(Version)
		return
	}
	
	fmt.Printf("DevPlatform CLI\n")
	fmt.Printf("Version:    %s\n", Version)
	fmt.Printf("Git Commit: %s\n", GitCommit)
	fmt.Printf("Build Date: %s\n", BuildDate)
	fmt.Printf("Go Version: %s\n", getGoVersion())
	
	if checkDeps {
		fmt.Println("\nDependency Versions:")
		checkDependency("terraform", "1.5.0")
		checkDependency("helm", "3.0.0")
		checkDependency("kubectl", "1.27.0")
		checkDependency("aws", "2.0.0")
		checkDependency("az", "2.0.0")
	}
}

func getGoVersion() string {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	// Output format: "go version go1.26.2 windows/amd64"
	parts := strings.Fields(string(output))
	if len(parts) >= 3 {
		return parts[2] // Returns "go1.26.2"
	}
	return "unknown"
}

func checkDependency(name string, minVersion string) {
	var cmd *exec.Cmd
	
	switch name {
	case "terraform":
		cmd = exec.Command("terraform", "version")
	case "helm":
		cmd = exec.Command("helm", "version", "--short")
	case "kubectl":
		cmd = exec.Command("kubectl", "version", "--client", "--short")
	case "aws":
		cmd = exec.Command("aws", "--version")
	case "az":
		cmd = exec.Command("az", "version", "--output", "json")
	default:
		fmt.Printf("  %-12s unknown tool\n", name+":")
		return
	}
	
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("  %-12s ❌ not found (minimum required: %s)\n", name+":", minVersion)
		return
	}
	
	version := parseVersion(name, string(output))
	if version == "" {
		fmt.Printf("  %-12s ⚠️  installed (version unknown, minimum required: %s)\n", name+":", minVersion)
	} else {
		fmt.Printf("  %-12s ✓ %s (minimum required: %s)\n", name+":", version, minVersion)
	}
}

func parseVersion(tool string, output string) string {
	output = strings.TrimSpace(output)
	lines := strings.Split(output, "\n")
	
	switch tool {
	case "terraform":
		// Output: "Terraform v1.14.8"
		if len(lines) > 0 {
			parts := strings.Fields(lines[0])
			if len(parts) >= 2 {
				return parts[1] // Returns "v1.14.8"
			}
		}
	case "helm":
		// Output: "v3.12.3+g3a31588"
		if len(lines) > 0 {
			parts := strings.Split(lines[0], "+")
			if len(parts) > 0 {
				return parts[0] // Returns "v3.12.3"
			}
		}
	case "kubectl":
		// Output: "Client Version: v1.28.1"
		if len(lines) > 0 {
			parts := strings.Fields(lines[0])
			if len(parts) >= 3 {
				return parts[2] // Returns "v1.28.1"
			}
		}
	case "aws":
		// Output: "aws-cli/2.x.x Python/3.x.x ..."
		parts := strings.Fields(output)
		if len(parts) > 0 {
			versionParts := strings.Split(parts[0], "/")
			if len(versionParts) >= 2 {
				return versionParts[1] // Returns "2.x.x"
			}
		}
	case "az":
		// Output is JSON, but we'll just check if it runs
		return "installed"
	}
	
	return ""
}
