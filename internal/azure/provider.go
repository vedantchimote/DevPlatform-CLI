package azure

import (
	"context"
	"fmt"

	"github.com/devplatform/devplatform-cli/internal/provider/types"
)

// AzureProvider implements the CloudProvider interface for Azure
type AzureProvider struct {
	auth          *AuthValidator
	kubeconfig    *KubeconfigManager
	pricing       *PricingCalculator
	subscriptionID string
	tenantID      string
	location      string
	resourceGroup string
}

// NewAzureProvider creates a new Azure provider
func NewAzureProvider(ctx context.Context, subscriptionID string, tenantID string, location string, resourceGroup string) (*AzureProvider, error) {
	auth, err := NewAuthValidator(ctx, subscriptionID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure auth validator: %w", err)
	}
	
	kubeconfig := NewKubeconfigManager(subscriptionID, resourceGroup)
	pricing := NewPricingCalculator(location)
	
	return &AzureProvider{
		auth:          auth,
		kubeconfig:    kubeconfig,
		pricing:       pricing,
		subscriptionID: subscriptionID,
		tenantID:      tenantID,
		location:      location,
		resourceGroup: resourceGroup,
	}, nil
}

// ValidateCredentials validates Azure credentials
func (p *AzureProvider) ValidateCredentials(ctx context.Context) error {
	return p.auth.ValidateCredentials(ctx)
}

// GetCallerIdentity returns Azure caller identity information
func (p *AzureProvider) GetCallerIdentity(ctx context.Context) (*types.CallerIdentity, error) {
	identity, err := p.auth.GetCallerIdentity(ctx)
	if err != nil {
		return nil, err
	}
	
	// Map Azure identity to provider identity
	// Using SubscriptionID as Account, TenantID as Arn, and State as UserId
	return &types.CallerIdentity{
		Account: identity.SubscriptionID,
		Arn:     identity.TenantID,
		UserId:  identity.State,
	}, nil
}

// UpdateKubeconfig updates the kubeconfig for AKS cluster access
func (p *AzureProvider) UpdateKubeconfig(clusterName string) error {
	return p.kubeconfig.UpdateKubeconfig(clusterName)
}

// GetConnectionCommands returns kubectl commands for connecting to the cluster
func (p *AzureProvider) GetConnectionCommands(clusterName string, namespace string) []string {
	return p.kubeconfig.GetConnectionCommands(clusterName, namespace)
}

// CalculateTotalCost calculates the total monthly cost for an environment
func (p *AzureProvider) CalculateTotalCost(envType string) (*types.EnvironmentCosts, error) {
	costs, err := p.pricing.CalculateTotalCost(envType)
	if err != nil {
		return nil, err
	}
	
	return &types.EnvironmentCosts{
		NetworkCost:  costs.VNetCost,
		DatabaseCost: costs.DatabaseCost,
		K8sCost:      costs.AKSCost,
		TotalCost:    costs.TotalCost,
		Environment:  costs.Environment,
		Provider:     "azure",
	}, nil
}

// GetTerraformBackend returns the Terraform backend configuration for Azure
func (p *AzureProvider) GetTerraformBackend(appName string, envType string) (*types.TerraformBackend, error) {
	// Azure Storage backend configuration
	storageAccountName := fmt.Sprintf("devplatformtf%s", p.location)
	containerName := "terraform-state"
	key := fmt.Sprintf("%s/%s/terraform.tfstate", appName, envType)
	
	return &types.TerraformBackend{
		Type: "azurerm",
		Config: map[string]string{
			"storage_account_name": storageAccountName,
			"container_name":       containerName,
			"key":                  key,
			"resource_group_name":  p.resourceGroup,
			"subscription_id":      p.subscriptionID,
			"tenant_id":            p.tenantID,
		},
	}, nil
}

// GetModulePath returns the path to the Terraform modules for Azure
func (p *AzureProvider) GetModulePath() string {
	return "terraform/modules/azure"
}

// GetProviderName returns the name of the cloud provider
func (p *AzureProvider) GetProviderName() string {
	return "azure"
}
