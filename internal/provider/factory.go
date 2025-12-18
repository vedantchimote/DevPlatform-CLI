package provider

import (
	"context"
	"fmt"
)

// ProviderConfig contains configuration for creating a cloud provider
type ProviderConfig struct {
	Provider string
	Region   string
	Profile  string
	
	// Azure-specific
	SubscriptionID string
	TenantID       string
	Location       string
}

// NewProvider creates a new cloud provider based on the configuration
func NewProvider(ctx context.Context, cfg *ProviderConfig) (CloudProvider, error) {
	switch cfg.Provider {
	case "aws":
		return newAWSProvider(ctx, cfg)
	case "azure":
		return newAzureProvider(ctx, cfg)
	default:
		return nil, fmt.Errorf("unsupported cloud provider: %s (supported: aws, azure)", cfg.Provider)
	}
}

// newAWSProvider creates an AWS provider (placeholder for now)
func newAWSProvider(ctx context.Context, cfg *ProviderConfig) (CloudProvider, error) {
	// This will be implemented in Task 7.5.3
	return nil, fmt.Errorf("AWS provider implementation pending")
}

// newAzureProvider creates an Azure provider (placeholder for now)
func newAzureProvider(ctx context.Context, cfg *ProviderConfig) (CloudProvider, error) {
	// This will be implemented in Task 7.6
	return nil, fmt.Errorf("Azure provider implementation pending")
}
