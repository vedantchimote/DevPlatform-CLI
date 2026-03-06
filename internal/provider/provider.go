package provider

import (
	"context"

	"github.com/devplatform/devplatform-cli/internal/provider/types"
)

// CloudProvider defines the interface for cloud provider operations
type CloudProvider interface {
	// ValidateCredentials validates the cloud provider credentials
	ValidateCredentials(ctx context.Context) error
	
	// GetCallerIdentity returns information about the authenticated identity
	GetCallerIdentity(ctx context.Context) (*types.CallerIdentity, error)
	
	// UpdateKubeconfig updates the kubeconfig for cluster access
	UpdateKubeconfig(clusterName string) error
	
	// GetConnectionCommands returns kubectl commands for connecting to the cluster
	GetConnectionCommands(clusterName string, namespace string) []string
	
	// CalculateTotalCost calculates the total monthly cost for an environment
	CalculateTotalCost(envType string) (*types.EnvironmentCosts, error)
	
	// GetTerraformBackend returns the Terraform backend configuration
	GetTerraformBackend(appName string, envType string) (*types.TerraformBackend, error)
	
	// GetModulePath returns the path to the Terraform modules for this provider
	GetModulePath() string
	
	// GetProviderName returns the name of the cloud provider
	GetProviderName() string
}
