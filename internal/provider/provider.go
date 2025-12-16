package provider

import (
	"context"
)

// CloudProvider defines the interface for cloud provider operations
type CloudProvider interface {
	// ValidateCredentials validates the cloud provider credentials
	ValidateCredentials(ctx context.Context) error
	
	// GetCallerIdentity returns information about the authenticated identity
	GetCallerIdentity(ctx context.Context) (*CallerIdentity, error)
	
	// UpdateKubeconfig updates the kubeconfig for cluster access
	UpdateKubeconfig(clusterName string) error
	
	// GetConnectionCommands returns kubectl commands for connecting to the cluster
	GetConnectionCommands(clusterName string, namespace string) []string
	
	// CalculateTotalCost calculates the total monthly cost for an environment
	CalculateTotalCost(envType string) (*EnvironmentCosts, error)
	
	// GetTerraformBackend returns the Terraform backend configuration
	GetTerraformBackend(appName string, envType string) (*TerraformBackend, error)
	
	// GetModulePath returns the path to the Terraform modules for this provider
	GetModulePath() string
	
	// GetProviderName returns the name of the cloud provider
	GetProviderName() string
}

// CallerIdentity represents cloud provider identity information
type CallerIdentity struct {
	Account string
	Arn     string
	UserId  string
}

// EnvironmentCosts represents the total costs for an environment
type EnvironmentCosts struct {
	NetworkCost  float64 // VPC/VNet cost
	DatabaseCost float64 // RDS/Azure Database cost
	K8sCost      float64 // EKS/AKS tenant cost
	TotalCost    float64
	Environment  string
	Provider     string
}

// TerraformBackend represents Terraform backend configuration
type TerraformBackend struct {
	Type   string            // "s3" or "azurerm"
	Config map[string]string // Backend-specific configuration
}
