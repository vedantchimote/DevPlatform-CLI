package config

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator validates configuration
type Validator struct {
	config *Config
	errors []string
}

// NewValidator creates a new configuration validator
func NewValidator(cfg *Config) *Validator {
	return &Validator{
		config: cfg,
		errors: []string{},
	}
}

// Validate performs validation on the configuration
func (v *Validator) Validate() error {
	// Validate global config
	v.validateGlobal()
	
	// Validate environment configs
	v.validateEnvironments()
	
	// Validate cloud-specific configs based on provider
	if v.config.Global.CloudProvider == "aws" {
		v.validateAWS()
	} else if v.config.Global.CloudProvider == "azure" {
		v.validateAzure()
	}
	
	// Return errors if any
	if len(v.errors) > 0 {
		return fmt.Errorf("configuration validation failed:\n  - %s", strings.Join(v.errors, "\n  - "))
	}
	
	return nil
}

func (v *Validator) validateGlobal() {
	// Validate cloud provider
	if v.config.Global.CloudProvider != "" {
		if v.config.Global.CloudProvider != "aws" && v.config.Global.CloudProvider != "azure" {
			v.errors = append(v.errors, fmt.Sprintf("invalid cloud_provider '%s': must be 'aws' or 'azure'", v.config.Global.CloudProvider))
		}
	}
	
	// Validate log level
	if v.config.Global.LogLevel != "" {
		validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
		if !validLevels[v.config.Global.LogLevel] {
			v.errors = append(v.errors, fmt.Sprintf("invalid log_level '%s': must be one of debug, info, warn, error", v.config.Global.LogLevel))
		}
	}
	
	// Validate timeout
	if v.config.Global.Timeout < 0 {
		v.errors = append(v.errors, "timeout must be a positive number")
	}
}

func (v *Validator) validateEnvironments() {
	validEnvs := map[string]bool{"dev": true, "staging": true, "prod": true}
	
	for envName, envConfig := range v.config.Environments {
		// Validate environment name
		if !validEnvs[envName] {
			v.errors = append(v.errors, fmt.Sprintf("invalid environment name '%s': must be one of dev, staging, prod", envName))
		}
		
		// Validate network CIDR
		if envConfig.NetworkCIDR != "" {
			if !isValidCIDR(envConfig.NetworkCIDR) {
				v.errors = append(v.errors, fmt.Sprintf("invalid network_cidr '%s' for environment '%s'", envConfig.NetworkCIDR, envName))
			}
		}
		
		// Validate database storage
		if envConfig.DBAllocatedStorage < 0 {
			v.errors = append(v.errors, fmt.Sprintf("db_allocated_storage must be positive for environment '%s'", envName))
		}
		
		// Validate k8s node count
		if envConfig.K8sNodeCount < 0 {
			v.errors = append(v.errors, fmt.Sprintf("k8s_node_count must be positive for environment '%s'", envName))
		}
	}
}

func (v *Validator) validateAWS() {
	// Validate AWS region format
	if v.config.AWS.Region != "" {
		if !isValidAWSRegion(v.config.AWS.Region) {
			v.errors = append(v.errors, fmt.Sprintf("invalid AWS region '%s'", v.config.AWS.Region))
		}
	}
}

func (v *Validator) validateAzure() {
	// Validate Azure location
	if v.config.Azure.Location == "" {
		v.errors = append(v.errors, "azure.location is required when cloud_provider is 'azure'")
	}
	
	// Validate subscription ID format (if provided)
	if v.config.Azure.SubscriptionID != "" {
		if !isValidUUID(v.config.Azure.SubscriptionID) {
			v.errors = append(v.errors, fmt.Sprintf("invalid Azure subscription_id format: %s", v.config.Azure.SubscriptionID))
		}
	}
	
	// Validate tenant ID format (if provided)
	if v.config.Azure.TenantID != "" {
		if !isValidUUID(v.config.Azure.TenantID) {
			v.errors = append(v.errors, fmt.Sprintf("invalid Azure tenant_id format: %s", v.config.Azure.TenantID))
		}
	}
}

// ValidateAppName validates an application name
func ValidateAppName(name string) error {
	if name == "" {
		return fmt.Errorf("app name cannot be empty")
	}
	
	if len(name) < 3 || len(name) > 32 {
		return fmt.Errorf("app name must be between 3 and 32 characters, got %d", len(name))
	}
	
	// Must contain only lowercase alphanumeric characters and hyphens
	matched, _ := regexp.MatchString("^[a-z0-9-]+$", name)
	if !matched {
		return fmt.Errorf("app name must contain only lowercase alphanumeric characters and hyphens")
	}
	
	// Cannot start or end with hyphen
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return fmt.Errorf("app name cannot start or end with a hyphen")
	}
	
	return nil
}

// ValidateEnvironmentType validates an environment type
func ValidateEnvironmentType(env string) error {
	validEnvs := map[string]bool{"dev": true, "staging": true, "prod": true}
	
	if !validEnvs[env] {
		return fmt.Errorf("invalid environment type '%s': must be one of dev, staging, prod", env)
	}
	
	return nil
}

// ValidateCloudProvider validates a cloud provider
func ValidateCloudProvider(provider string) error {
	if provider != "aws" && provider != "azure" {
		return fmt.Errorf("invalid cloud provider '%s': must be 'aws' or 'azure'", provider)
	}
	
	return nil
}

// Helper functions

func isValidCIDR(cidr string) bool {
	// Simple CIDR validation: x.x.x.x/y
	matched, _ := regexp.MatchString(`^(\d{1,3}\.){3}\d{1,3}/\d{1,2}$`, cidr)
	return matched
}

func isValidAWSRegion(region string) bool {
	// AWS region format: us-east-1, eu-west-2, etc.
	matched, _ := regexp.MatchString(`^[a-z]{2}-[a-z]+-\d+$`, region)
	return matched
}

func isValidUUID(uuid string) bool {
	// UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	matched, _ := regexp.MatchString(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, uuid)
	return matched
}
