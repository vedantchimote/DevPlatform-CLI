package mocks_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/devplatform/devplatform-cli/internal/provider/types"
	"github.com/devplatform/devplatform-cli/test/mocks"
)

// ExampleMockAzureProvider_basicUsage demonstrates basic mock usage
func ExampleMockAzureProvider_basicUsage() {
	// Create a new mock
	mock := mocks.NewMockAzureProvider()

	// Use the mock with default behavior (all operations succeed)
	ctx := context.Background()
	err := mock.ValidateCredentials(ctx)
	if err != nil {
		fmt.Printf("ValidateCredentials failed: %v\n", err)
		return
	}

	// Check how many times ValidateCredentials was called
	fmt.Printf("ValidateCredentials called %d times\n", mock.GetValidateCredentialsCallCount())

	// Output:
	// ValidateCredentials called 1 times
}

// ExampleMockAzureProvider_customBehavior demonstrates configuring custom behavior
func ExampleMockAzureProvider_customBehavior() {
	// Create a new mock
	mock := mocks.NewMockAzureProvider()

	// Configure custom behavior for GetCallerIdentity
	mock.GetCallerIdentityFunc = func(ctx context.Context) (*types.CallerIdentity, error) {
		return &types.CallerIdentity{
			Account: "aaaabbbb-cccc-dddd-eeee-ffff00001111",
			Arn:     "22223333-4444-5555-6666-777788889999",
			UserId:  "AzureAdmin",
		}, nil
	}

	// Use the mock
	ctx := context.Background()
	identity, err := mock.GetCallerIdentity(ctx)
	if err != nil {
		fmt.Printf("GetCallerIdentity failed: %v\n", err)
		return
	}

	fmt.Printf("Subscription: %s, User: %s\n", identity.Account, identity.UserId)

	// Output:
	// Subscription: aaaabbbb-cccc-dddd-eeee-ffff00001111, User: AzureAdmin
}

// ExampleMockAzureProvider_errorSimulation demonstrates simulating errors
func ExampleMockAzureProvider_errorSimulation() {
	// Create a new mock
	mock := mocks.NewMockAzureProvider()

	// Configure ValidateCredentials to return an error
	mock.ValidateCredentialsFunc = func(ctx context.Context) error {
		return errors.New("Azure credentials not configured")
	}

	// Use the mock
	ctx := context.Background()
	err := mock.ValidateCredentials(ctx)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	}

	// Output:
	// Validation failed: Azure credentials not configured
}

// ExampleMockAzureProvider_callTracking demonstrates call tracking
func ExampleMockAzureProvider_callTracking() {
	// Create a new mock
	mock := mocks.NewMockAzureProvider()

	// Make multiple calls
	_ = mock.UpdateKubeconfig("aks-cluster-1")
	_ = mock.UpdateKubeconfig("aks-cluster-2")
	_ = mock.UpdateKubeconfig("aks-cluster-3")

	// Check call count
	fmt.Printf("UpdateKubeconfig called %d times\n", mock.GetUpdateKubeconfigCallCount())

	// Inspect call arguments
	for i, call := range mock.UpdateKubeconfigCalls {
		clusterName := call.Args[0].(string)
		fmt.Printf("Call %d: cluster=%s\n", i+1, clusterName)
	}

	// Output:
	// UpdateKubeconfig called 3 times
	// Call 1: cluster=aks-cluster-1
	// Call 2: cluster=aks-cluster-2
	// Call 3: cluster=aks-cluster-3
}

// ExampleMockAzureProvider_reset demonstrates resetting the mock
func ExampleMockAzureProvider_reset() {
	// Create a new mock
	mock := mocks.NewMockAzureProvider()

	// Make some calls
	ctx := context.Background()
	_ = mock.ValidateCredentials(ctx)
	_, _ = mock.GetCallerIdentity(ctx)

	fmt.Printf("Before reset - ValidateCredentials: %d, GetCallerIdentity: %d\n",
		mock.GetValidateCredentialsCallCount(), mock.GetCallerIdentityCallCount())

	// Reset the mock
	mock.Reset()

	fmt.Printf("After reset - ValidateCredentials: %d, GetCallerIdentity: %d\n",
		mock.GetValidateCredentialsCallCount(), mock.GetCallerIdentityCallCount())

	// Output:
	// Before reset - ValidateCredentials: 1, GetCallerIdentity: 1
	// After reset - ValidateCredentials: 0, GetCallerIdentity: 0
}

// ExampleMockAzureProvider_calculateCost demonstrates mocking cost calculation
func ExampleMockAzureProvider_calculateCost() {
	// Create a new mock
	mock := mocks.NewMockAzureProvider()

	// Use default behavior
	costs, err := mock.CalculateTotalCost("prod")
	if err != nil {
		fmt.Printf("CalculateTotalCost failed: %v\n", err)
		return
	}

	fmt.Printf("Environment: %s\n", costs.Environment)
	fmt.Printf("Total Cost: $%.2f/month\n", costs.TotalCost)
	fmt.Printf("Provider: %s\n", costs.Provider)

	// Output:
	// Environment: prod
	// Total Cost: $560.00/month
	// Provider: azure
}

// ExampleMockAzureProvider_connectionCommands demonstrates getting connection commands
func ExampleMockAzureProvider_connectionCommands() {
	// Create a new mock
	mock := mocks.NewMockAzureProvider()

	// Get connection commands
	commands := mock.GetConnectionCommands("my-aks-cluster", "production")

	fmt.Printf("Connection commands:\n")
	for i, cmd := range commands {
		fmt.Printf("%d. %s\n", i+1, cmd)
	}

	// Output:
	// Connection commands:
	// 1. az aks get-credentials --name my-aks-cluster --resource-group devplatform-rg
	// 2. kubectl config set-context --current --namespace=production
	// 3. kubectl get pods
}

// ExampleMockAzureProvider_terraformBackend demonstrates getting Terraform backend config
func ExampleMockAzureProvider_terraformBackend() {
	// Create a new mock
	mock := mocks.NewMockAzureProvider()

	// Get Terraform backend configuration
	backend, err := mock.GetTerraformBackend("my-app", "staging")
	if err != nil {
		fmt.Printf("GetTerraformBackend failed: %v\n", err)
		return
	}

	fmt.Printf("Backend Type: %s\n", backend.Type)
	fmt.Printf("Storage Account: %s\n", backend.Config["storage_account_name"])
	fmt.Printf("Container: %s\n", backend.Config["container_name"])
	fmt.Printf("Key: %s\n", backend.Config["key"])

	// Output:
	// Backend Type: azurerm
	// Storage Account: devplatformtfeastus
	// Container: terraform-state
	// Key: my-app/staging/terraform.tfstate
}
