package config

// Merger handles merging CLI flags with configuration file values
type Merger struct {
	config *Config
}

// NewMerger creates a new configuration merger
func NewMerger(cfg *Config) *Merger {
	return &Merger{
		config: cfg,
	}
}

// MergeFlags merges command-line flags with configuration
// CLI flags take precedence over configuration file values
func (m *Merger) MergeFlags(flags map[string]interface{}) {
	// Merge global flags
	if provider, ok := flags["provider"].(string); ok && provider != "" {
		m.config.Global.CloudProvider = provider
	}
	
	if timeout, ok := flags["timeout"].(int); ok && timeout > 0 {
		m.config.Global.Timeout = timeout
	}
	
	if logLevel, ok := flags["log_level"].(string); ok && logLevel != "" {
		m.config.Global.LogLevel = logLevel
	}
	
	// Merge AWS flags
	if region, ok := flags["aws_region"].(string); ok && region != "" {
		m.config.AWS.Region = region
	}
	
	if profile, ok := flags["aws_profile"].(string); ok && profile != "" {
		m.config.AWS.Profile = profile
	}
	
	// Merge Azure flags
	if subscriptionID, ok := flags["azure_subscription_id"].(string); ok && subscriptionID != "" {
		m.config.Azure.SubscriptionID = subscriptionID
	}
	
	if location, ok := flags["azure_location"].(string); ok && location != "" {
		m.config.Azure.Location = location
	}
	
	if tenantID, ok := flags["azure_tenant_id"].(string); ok && tenantID != "" {
		m.config.Azure.TenantID = tenantID
	}
}

// GetConfig returns the merged configuration
func (m *Merger) GetConfig() *Config {
	return m.config
}

// MergeWithFlags is a helper function that loads config and merges with flags
func MergeWithFlags(configFile string, flags map[string]interface{}) (*Config, error) {
	// Load configuration
	loader := NewLoader(configFile)
	cfg, err := loader.Load()
	if err != nil {
		return nil, err
	}
	
	// Merge flags
	merger := NewMerger(cfg)
	merger.MergeFlags(flags)
	
	// Validate merged configuration
	validator := NewValidator(cfg)
	if err := validator.Validate(); err != nil {
		return nil, err
	}
	
	return cfg, nil
}
