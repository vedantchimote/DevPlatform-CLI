package config

import (
	"fmt"
	"os"
	"path/filepath"

	clierrors "github.com/devplatform/devplatform-cli/internal/errors"
	"github.com/spf13/viper"
)

// Loader handles loading configuration from files
type Loader struct {
	configFile string
}

// NewLoader creates a new configuration loader
func NewLoader(configFile string) *Loader {
	return &Loader{
		configFile: configFile,
	}
}

// Load reads and parses the configuration file
func (l *Loader) Load() (*Config, error) {
	// Start with default configuration
	cfg := NewDefaultConfig()

	// If no config file specified, look for .devplatform.yaml in current directory
	configPath := l.configFile
	if configPath == "" {
		configPath = ".devplatform.yaml"
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config file doesn't exist, return defaults
		return cfg, nil
	}

	// Set up Viper
	v := viper.New()
	
	// Set config file path
	v.SetConfigFile(configPath)
	
	// Set config type
	ext := filepath.Ext(configPath)
	if ext != "" {
		v.SetConfigType(ext[1:]) // Remove the dot
	} else {
		v.SetConfigType("yaml")
	}

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, clierrors.NewConfigError(
			clierrors.ErrCodeConfigParseFailed,
			fmt.Sprintf("Failed to read config file: %s", configPath),
			err,
		).WithDetails(fmt.Sprintf("Config path: %s", configPath))
	}

	// Unmarshal into config struct
	if err := v.Unmarshal(cfg); err != nil {
		return nil, clierrors.NewConfigError(
			clierrors.ErrCodeConfigParseFailed,
			fmt.Sprintf("Failed to parse config file: %s", configPath),
			err,
		).WithDetails("Check YAML syntax and structure")
	}

	return cfg, nil
}

// LoadFromPath loads configuration from a specific path
func LoadFromPath(path string) (*Config, error) {
	loader := NewLoader(path)
	return loader.Load()
}

// LoadDefault loads configuration from the default location (.devplatform.yaml)
func LoadDefault() (*Config, error) {
	loader := NewLoader("")
	return loader.Load()
}
