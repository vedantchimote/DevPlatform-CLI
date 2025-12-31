package provider

import (
	"context"
	"fmt"

	"github.com/devplatform/devplatform-cli/internal/aws"
	"github.com/devplatform/devplatform-cli/internal/azure"
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
	ResourceGroup  string
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

// newAWSProvider creates an AWS provider
func newAWSProvider(ctx context.Context, cfg *ProviderConfig) (CloudProvider, error) {
	return aws.NewAWSProvider(ctx, cfg.Region, cfg.Profile)
}

// newAzureProvider creates an Azure provider
func newAzureProvider(ctx context.Context, cfg *ProviderConfig) (CloudProvider, error) {
	return azure.NewAzureProvider(ctx, cfg.SubscriptionID, cfg.TenantID, cfg.Location, cfg.ResourceGroup)
}
