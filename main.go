package main

import (
	"fmt"
	"os"
)

var (
	// Version information - will be set during build
	Version   = "dev"
	GitCommit = "none"
	BuildDate = "unknown"
)

func main() {
	fmt.Println("DevPlatform CLI")
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Git Commit: %s\n", GitCommit)
	fmt.Printf("Build Date: %s\n", BuildDate)
	
	os.Exit(0)
}
