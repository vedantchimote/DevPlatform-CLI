package mocks_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/devplatform/devplatform-cli/internal/provider/types"
	"github.com/devplatform/devplatform-cli/test/mocks"
)

// ExampleMockAWSProvider_basicUsage demonstrates basic mock usage
func ExampleMockAWSProvider_basicUsage() {
	// Create a new mock
	mock := mocks.NewMockAWSProvider()

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

// ExampleMockAWSProvider_customBehavior demonstrates configuring custom behavior
func ExampleMockAWSProvider_customBehavior() {
	// Create a new mock
	mock := mocks.NewMockAWSProvider()

	// Configure custom behavior for GetCallerIdentity
	mock.GetCallerIdentityFunc = func(ctx context.Context) (*types.CallerIdentity, error) {
		return &types.CallerIdentity{
			Account: "111222333444",
			Arn:     "arn:aws:iam::111222333444:user/admin",
			UserId:  "ADMINUSERID",
		}, nil
	}

	// Use the mock
	ctx := context.Background()
	identity, err := mock.GetCallerIdentity(ctx)
	if err != nil {
		fmt.Printf("GetCallerIdentity failed: %v\n", err)
		return
	}

	fmt.Printf("Account: %s, User: %s\n", identity.Account, identity.UserId)

	// Output:
	// Account: 111222333444, User: ADMINUSERID
}

// ExampleMockAWSProvider_errorSimulation demonstrates simulating errors
func ExampleMockAWSProvider_errorSimulation() {
	// Create a new mock
	mock := mocks.NewMockAWSProvider()

	// Configure ValidateCredentials to return an error
	mock.ValidateCredentialsFunc = func(ctx context.Context) error {
		return errors.New("AWS credentials not found")
	}

	// Use the mock
	ctx := context.Background()
	err := mock.ValidateCredentials(ctx)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	}

	// Output:
	// Validation failed: AWS credentials not found
}

// ExampleMockAWSProvider_callTracking demonstrates call tracking
func ExampleMockAWSProvider_callTracking() {
	// Create a new mock
	mock := mocks.NewMockAWSProvider()

	// Make multiple calls
	_ = mock.UpdateKubeconfig("cluster-1")
	_ = mock.UpdateKubeconfig("cluster-2")
	_ = mock.UpdateKubeconfig("cluster-3")

	// Check call count
	fmt.Printf("UpdateKubeconfig called %d times\n", mock.GetUpdateKubeconfigCallCount())

	// Inspect call arguments
	for i, call := range mock.UpdateKubeconfigCalls {
		clusterName := call.Args[0].(string)
		fmt.Printf("Call %d: cluster=%s\n", i+1, clusterName)
	}

	// Output:
	// UpdateKubeconfig called 3 times
	// Call 1: cluster=cluster-1
	// Call 2: cluster=cluster-2
	// Call 3: cluster=cluster-3
}

// ExampleMockAWSProvider_reset demonstrates resetting the mock
func ExampleMockAWSProvider_reset() {
	// Create a new mock
	mock := mocks.NewMockAWSProvider()

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

// ExampleMockAWSProvider_calculateCost demonstrates mocking cost calculation
func ExampleMockAWSProvider_calculateCost() {
	// Create a new mock
	mock := mocks.NewMockAWSProvider()

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
	// Total Cost: $600.00/month
	// Provider: aws
}

// ExampleMockAWSProvider_connectionCommands demonstrates getting connection commands
func ExampleMockAWSProvider_connectionCommands() {
	// Create a new mock
	mock := mocks.NewMockAWSProvider()

	// Get connection commands
	commands := mock.GetConnectionCommands("my-cluster", "production")

	fmt.Printf("Connection commands:\n")
	for i, cmd := range commands {
		fmt.Printf("%d. %s\n", i+1, cmd)
	}

	// Output:
	// Connection commands:
	// 1. aws eks update-kubeconfig --name my-cluster --region us-east-1
	// 2. kubectl config set-context --current --namespace=production
	// 3. kubectl get pods
}

// ExampleMockAWSProvider_terraformBackend demonstrates getting Terraform backend config
func ExampleMockAWSProvider_terraformBackend() {
	// Create a new mock
	mock := mocks.NewMockAWSProvider()

	// Get Terraform backend configuration
	backend, err := mock.GetTerraformBackend("my-app", "staging")
	if err != nil {
		fmt.Printf("GetTerraformBackend failed: %v\n", err)
		return
	}

	fmt.Printf("Backend Type: %s\n", backend.Type)
	fmt.Printf("Bucket: %s\n", backend.Config["bucket"])
	fmt.Printf("Key: %s\n", backend.Config["key"])

	// Output:
	// Backend Type: s3
	// Bucket: devplatform-terraform-state-us-east-1
	// Key: my-app/staging/terraform.tfstate
}
