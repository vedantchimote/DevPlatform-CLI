package azure

import (
	"context"
	"testing"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

func TestNewAzureProvider(t *testing.T) {
	tests := []struct {
		name           string
		subscriptionID string
		tenantID       string
		location       string
		resourceGroup  string
		wantErr        bool
	}{
		{
			name:           "valid_provider",
			subscriptionID: "12345678-1234-1234-1234-123456789012",
			tenantID:       "87654321-4321-4321-4321-210987654321",
			location:       "eastus",
			resourceGroup:  "devplatform-rg",
			wantErr:        false,
		},
		{
			name:           "valid_provider_different_location",
			subscriptionID: "11111111-1111-1111-1111-111111111111",
			tenantID:       "22222222-2222-2222-2222-222222222222",
			location:       "westus2",
			resourceGroup:  "test-rg",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			provider, err := NewAzureProvider(ctx, tt.subscriptionID, tt.tenantID, tt.location, tt.resourceGroup)

			if tt.wantErr {
				testutil.AssertError(t, err)
				return
			}

			testutil.AssertNoError(t, err)
			testutil.AssertNotEqual(t, nil, provider)
			testutil.AssertEqual(t, tt.subscriptionID, provider.subscriptionID)
			testutil.AssertEqual(t, tt.tenantID, provider.tenantID)
			testutil.AssertEqual(t, tt.location, provider.location)
			testutil.AssertEqual(t, tt.resourceGroup, provider.resourceGroup)
			testutil.AssertNotEqual(t, nil, provider.auth)
			testutil.AssertNotEqual(t, nil, provider.kubeconfig)
			testutil.AssertNotEqual(t, nil, provider.pricing)
		})
	}
}

func TestAzureProvider_GetProviderName(t *testing.T) {
	ctx := context.Background()
	provider, err := NewAzureProvider(ctx, "12345678-1234-1234-1234-123456789012", "87654321-4321-4321-4321-210987654321", "eastus", "devplatform-rg")
	testutil.AssertNoError(t, err)

	name := provider.GetProviderName()
	testutil.AssertEqual(t, "azure", name)
}

func TestAzureProvider_GetModulePath(t *testing.T) {
	ctx := context.Background()
	provider, err := NewAzureProvider(ctx, "12345678-1234-1234-1234-123456789012", "87654321-4321-4321-4321-210987654321", "eastus", "devplatform-rg")
	testutil.AssertNoError(t, err)

	path := provider.GetModulePath()
	testutil.AssertEqual(t, "terraform/modules/azure", path)
}

func TestAzureProvider_GetTerraformBackend(t *testing.T) {
	tests := []struct {
		name           string
		subscriptionID string
		tenantID       string
		location       string
		resourceGroup  string
		appName        string
		envType        string
	}{
		{
			name:           "dev_environment",
			subscriptionID: "12345678-1234-1234-1234-123456789012",
			tenantID:       "87654321-4321-4321-4321-210987654321",
			location:       "eastus",
			resourceGroup:  "devplatform-rg",
			appName:        "myapp",
			envType:        "dev",
		},
		{
			name:           "staging_environment",
			subscriptionID: "11111111-1111-1111-1111-111111111111",
			tenantID:       "22222222-2222-2222-2222-222222222222",
			location:       "westus2",
			resourceGroup:  "test-rg",
			appName:        "testapp",
			envType:        "staging",
		},
		{
			name:           "prod_environment",
			subscriptionID: "33333333-3333-3333-3333-333333333333",
			tenantID:       "44444444-4444-4444-4444-444444444444",
			location:       "northeurope",
			resourceGroup:  "prod-rg",
			appName:        "prodapp",
			envType:        "prod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			provider, err := NewAzureProvider(ctx, tt.subscriptionID, tt.tenantID, tt.location, tt.resourceGroup)
			testutil.AssertNoError(t, err)

			backend, err := provider.GetTerraformBackend(tt.appName, tt.envType)
			testutil.AssertNoError(t, err)
			testutil.AssertNotEqual(t, nil, backend)

			// Verify backend type
			testutil.AssertEqual(t, "azurerm", backend.Type)

			// Verify backend config
			testutil.AssertNotEqual(t, nil, backend.Config)
			testutil.AssertEqual(t, "devplatformtf"+tt.location, backend.Config["storage_account_name"])
			testutil.AssertEqual(t, "terraform-state", backend.Config["container_name"])
			testutil.AssertEqual(t, tt.appName+"/"+tt.envType+"/terraform.tfstate", backend.Config["key"])
			testutil.AssertEqual(t, tt.resourceGroup, backend.Config["resource_group_name"])
			testutil.AssertEqual(t, tt.subscriptionID, backend.Config["subscription_id"])
			testutil.AssertEqual(t, tt.tenantID, backend.Config["tenant_id"])
		})
	}
}

func TestAzureProvider_GetConnectionCommands(t *testing.T) {
	tests := []struct {
		name        string
		clusterName string
		namespace   string
	}{
		{
			name:        "dev_cluster",
			clusterName: "dev-cluster",
			namespace:   "dev-ns",
		},
		{
			name:        "prod_cluster",
			clusterName: "prod-cluster",
			namespace:   "prod-ns",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			provider, err := NewAzureProvider(ctx, "12345678-1234-1234-1234-123456789012", "87654321-4321-4321-4321-210987654321", "eastus", "devplatform-rg")
			testutil.AssertNoError(t, err)

			commands := provider.GetConnectionCommands(tt.clusterName, tt.namespace)
			testutil.AssertNotEqual(t, nil, commands)
			testutil.AssertEqual(t, 2, len(commands))

			// Verify commands contain cluster and namespace
			testutil.AssertContains(t, commands[0], tt.clusterName)
			testutil.AssertContains(t, commands[1], tt.namespace)
		})
	}
}

func TestAzureProvider_CalculateTotalCost(t *testing.T) {
	tests := []struct {
		name    string
		envType string
		wantErr bool
	}{
		{
			name:    "dev_environment",
			envType: "dev",
			wantErr: false,
		},
		{
			name:    "staging_environment",
			envType: "staging",
			wantErr: false,
		},
		{
			name:    "prod_environment",
			envType: "prod",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			provider, err := NewAzureProvider(ctx, "12345678-1234-1234-1234-123456789012", "87654321-4321-4321-4321-210987654321", "eastus", "devplatform-rg")
			testutil.AssertNoError(t, err)

			costs, err := provider.CalculateTotalCost(tt.envType)
			if tt.wantErr {
				testutil.AssertError(t, err)
				return
			}

			testutil.AssertNoError(t, err)
			testutil.AssertNotEqual(t, nil, costs)
			testutil.AssertEqual(t, tt.envType, costs.Environment)
			testutil.AssertEqual(t, "azure", costs.Provider)
			testutil.AssertTrue(t, costs.TotalCost > 0, "Total cost should be positive")
			testutil.AssertTrue(t, costs.NetworkCost > 0, "Network cost should be positive")
			testutil.AssertTrue(t, costs.DatabaseCost > 0, "Database cost should be positive")
			testutil.AssertTrue(t, costs.K8sCost > 0, "K8s cost should be positive")
		})
	}
}

func TestAzureProvider_MultipleLocations(t *testing.T) {
	locations := []string{"eastus", "westus2", "northeurope"}

	for _, location := range locations {
		t.Run("location_"+location, func(t *testing.T) {
			ctx := context.Background()
			provider, err := NewAzureProvider(ctx, "12345678-1234-1234-1234-123456789012", "87654321-4321-4321-4321-210987654321", location, "test-rg")
			testutil.AssertNoError(t, err)
			testutil.AssertNotEqual(t, nil, provider)
			testutil.AssertEqual(t, location, provider.location)
		})
	}
}
