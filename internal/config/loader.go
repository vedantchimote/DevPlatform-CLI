package config

import (
	"fmt"
	"os"
	"path/filepath"

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
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal into config struct
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
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
