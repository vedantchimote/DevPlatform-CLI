package mocks

import (
	"context"
	"errors"
	"testing"

	"github.com/devplatform/devplatform-cli/internal/provider/types"
)

func TestMockAzureProvider_ValidateCredentials(t *testing.T) {
	mock := NewMockAzureProvider()
	ctx := context.Background()

	// Test default behavior (success)
	err := mock.ValidateCredentials(ctx)
	if err != nil {
		t.Errorf("ValidateCredentials() unexpected error = %v", err)
	}

	// Verify call was recorded
	if len(mock.ValidateCredentialsCalls) != 1 {
		t.Errorf("Expected 1 ValidateCredentials call, got %d", len(mock.ValidateCredentialsCalls))
	}

	// Verify call count helper
	if mock.GetValidateCredentialsCallCount() != 1 {
		t.Errorf("GetValidateCredentialsCallCount() = %d, want 1", mock.GetValidateCredentialsCallCount())
	}
}

func TestMockAzureProvider_ValidateCredentialsWithCustomFunc(t *testing.T) {
	mock := NewMockAzureProvider()
	ctx := context.Background()
	expectedErr := errors.New("invalid credentials")

	// Configure custom behavior
	mock.ValidateCredentialsFunc = func(ctx context.Context) error {
		return expectedErr
	}

	// Test custom behavior
	err := mock.ValidateCredentials(ctx)
	if err != expectedErr {
		t.Errorf("ValidateCredentials() error = %v, want %v", err, expectedErr)
	}
}

func TestMockAzureProvider_GetCallerIdentity(t *testing.T) {
	mock := NewMockAzureProvider()
	ctx := context.Background()

	// Test default behavior
	identity, err := mock.GetCallerIdentity(ctx)
	if err != nil {
		t.Errorf("GetCallerIdentity() unexpected error = %v", err)
	}
	if identity == nil {
		t.Fatal("GetCallerIdentity() returned nil")
	}
	if identity.Account == "" {
		t.Error("GetCallerIdentity() returned empty Account")
	}

	// Verify call was recorded
	if len(mock.GetCallerIdentityCalls) != 1 {
		t.Errorf("Expected 1 GetCallerIdentity call, got %d", len(mock.GetCallerIdentityCalls))
	}

	// Verify call count helper
	if mock.GetCallerIdentityCallCount() != 1 {
		t.Errorf("GetCallerIdentityCallCount() = %d, want 1", mock.GetCallerIdentityCallCount())
	}
}

func TestMockAzureProvider_GetCallerIdentityWithCustomFunc(t *testing.T) {
	mock := NewMockAzureProvider()
	ctx := context.Background()
	expectedIdentity := &types.CallerIdentity{
		Account: "87654321-4321-4321-4321-210987654321",
		Arn:     "11111111-1111-1111-1111-111111111111",
		UserId:  "CustomUser",
	}

	// Configure custom behavior
	mock.GetCallerIdentityFunc = func(ctx context.Context) (*types.CallerIdentity, error) {
		return expectedIdentity, nil
	}

	// Test custom behavior
	identity, err := mock.GetCallerIdentity(ctx)
	if err != nil {
		t.Errorf("GetCallerIdentity() unexpected error = %v", err)
	}
	if identity.Account != expectedIdentity.Account {
		t.Errorf("GetCallerIdentity().Account = %s, want %s", identity.Account, expectedIdentity.Account)
	}
}

func TestMockAzureProvider_UpdateKubeconfig(t *testing.T) {
	mock := NewMockAzureProvider()
	clusterName := "test-aks-cluster"

	// Test default behavior
	err := mock.UpdateKubeconfig(clusterName)
	if err != nil {
		t.Errorf("UpdateKubeconfig() unexpected error = %v", err)
	}

	// Verify call was recorded
	if len(mock.UpdateKubeconfigCalls) != 1 {
		t.Errorf("Expected 1 UpdateKubeconfig call, got %d", len(mock.UpdateKubeconfigCalls))
	}

	// Verify call count helper
	if mock.GetUpdateKubeconfigCallCount() != 1 {
		t.Errorf("GetUpdateKubeconfigCallCount() = %d, want 1", mock.GetUpdateKubeconfigCallCount())
	}

	// Verify call arguments
	if mock.UpdateKubeconfigCalls[0].Args[0] != clusterName {
		t.Errorf("Expected clusterName %s, got %v", clusterName, mock.UpdateKubeconfigCalls[0].Args[0])
	}
}

func TestMockAzureProvider_GetConnectionCommands(t *testing.T) {
	mock := NewMockAzureProvider()
	clusterName := "test-aks-cluster"
	namespace := "test-namespace"

	// Test default behavior
	commands := mock.GetConnectionCommands(clusterName, namespace)
	if commands == nil {
		t.Fatal("GetConnectionCommands() returned nil")
	}
	if len(commands) == 0 {
		t.Error("GetConnectionCommands() returned empty slice")
	}

	// Verify call was recorded
	if len(mock.GetConnectionCommandsCalls) != 1 {
		t.Errorf("Expected 1 GetConnectionCommands call, got %d", len(mock.GetConnectionCommandsCalls))
	}

	// Verify call count helper
	if mock.GetConnectionCommandsCallCount() != 1 {
		t.Errorf("GetConnectionCommandsCallCount() = %d, want 1", mock.GetConnectionCommandsCallCount())
	}
}

func TestMockAzureProvider_GetConnectionCommandsWithCustomFunc(t *testing.T) {
	mock := NewMockAzureProvider()
	clusterName := "test-aks-cluster"
	namespace := "test-namespace"
	expectedCommands := []string{"custom azure command 1", "custom azure command 2"}

	// Configure custom behavior
	mock.GetConnectionCommandsFunc = func(clusterName string, namespace string) []string {
		return expectedCommands
	}

	// Test custom behavior
	commands := mock.GetConnectionCommands(clusterName, namespace)
	if len(commands) != len(expectedCommands) {
		t.Errorf("GetConnectionCommands() returned %d commands, want %d", len(commands), len(expectedCommands))
	}
}

func TestMockAzureProvider_CalculateTotalCost(t *testing.T) {
	mock := NewMockAzureProvider()

	tests := []struct {
		name    string
		envType string
	}{
		{"dev environment", "dev"},
		{"staging environment", "staging"},
		{"prod environment", "prod"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			costs, err := mock.CalculateTotalCost(tt.envType)
			if err != nil {
				t.Errorf("CalculateTotalCost() unexpected error = %v", err)
			}
			if costs == nil {
				t.Fatal("CalculateTotalCost() returned nil")
			}
			if costs.TotalCost <= 0 {
				t.Error("CalculateTotalCost() returned zero or negative cost")
			}
			if costs.Provider != "azure" {
				t.Errorf("CalculateTotalCost().Provider = %s, want azure", costs.Provider)
			}
		})
	}

	// Verify call count
	if mock.GetCalculateTotalCostCallCount() != 3 {
		t.Errorf("GetCalculateTotalCostCallCount() = %d, want 3", mock.GetCalculateTotalCostCallCount())
	}
}

func TestMockAzureProvider_GetTerraformBackend(t *testing.T) {
	mock := NewMockAzureProvider()
	appName := "test-app"
	envType := "dev"

	// Test default behavior
	backend, err := mock.GetTerraformBackend(appName, envType)
	if err != nil {
		t.Errorf("GetTerraformBackend() unexpected error = %v", err)
	}
	if backend == nil {
		t.Fatal("GetTerraformBackend() returned nil")
	}
	if backend.Type != "azurerm" {
		t.Errorf("GetTerraformBackend().Type = %s, want azurerm", backend.Type)
	}
	if backend.Config == nil {
		t.Fatal("GetTerraformBackend().Config is nil")
	}

	// Verify call was recorded
	if len(mock.GetTerraformBackendCalls) != 1 {
		t.Errorf("Expected 1 GetTerraformBackend call, got %d", len(mock.GetTerraformBackendCalls))
	}

	// Verify call count helper
	if mock.GetTerraformBackendCallCount() != 1 {
		t.Errorf("GetTerraformBackendCallCount() = %d, want 1", mock.GetTerraformBackendCallCount())
	}
}

func TestMockAzureProvider_GetModulePath(t *testing.T) {
	mock := NewMockAzureProvider()

	// Test default behavior
	path := mock.GetModulePath()
	if path == "" {
		t.Error("GetModulePath() returned empty string")
	}
	if path != "terraform/modules/azure" {
		t.Errorf("GetModulePath() = %s, want terraform/modules/azure", path)
	}

	// Verify call was recorded
	if len(mock.GetModulePathCalls) != 1 {
		t.Errorf("Expected 1 GetModulePath call, got %d", len(mock.GetModulePathCalls))
	}

	// Verify call count helper
	if mock.GetModulePathCallCount() != 1 {
		t.Errorf("GetModulePathCallCount() = %d, want 1", mock.GetModulePathCallCount())
	}
}

func TestMockAzureProvider_GetProviderName(t *testing.T) {
	mock := NewMockAzureProvider()

	// Test default behavior
	name := mock.GetProviderName()
	if name != "azure" {
		t.Errorf("GetProviderName() = %s, want azure", name)
	}

	// Verify call was recorded
	if len(mock.GetProviderNameCalls) != 1 {
		t.Errorf("Expected 1 GetProviderName call, got %d", len(mock.GetProviderNameCalls))
	}

	// Verify call count helper
	if mock.GetProviderNameCallCount() != 1 {
		t.Errorf("GetProviderNameCallCount() = %d, want 1", mock.GetProviderNameCallCount())
	}
}

func TestMockAzureProvider_Reset(t *testing.T) {
	mock := NewMockAzureProvider()
	ctx := context.Background()

	// Make several calls
	_ = mock.ValidateCredentials(ctx)
	_, _ = mock.GetCallerIdentity(ctx)
	_ = mock.UpdateKubeconfig("test-cluster")
	_ = mock.GetConnectionCommands("test-cluster", "default")

	// Verify calls were recorded
	if len(mock.ValidateCredentialsCalls) != 1 {
		t.Errorf("Expected 1 ValidateCredentials call before reset, got %d", len(mock.ValidateCredentialsCalls))
	}
	if len(mock.GetCallerIdentityCalls) != 1 {
		t.Errorf("Expected 1 GetCallerIdentity call before reset, got %d", len(mock.GetCallerIdentityCalls))
	}

	// Reset the mock
	mock.Reset()

	// Verify all calls were cleared
	if len(mock.ValidateCredentialsCalls) != 0 {
		t.Errorf("Expected 0 ValidateCredentials calls after reset, got %d", len(mock.ValidateCredentialsCalls))
	}
	if len(mock.GetCallerIdentityCalls) != 0 {
		t.Errorf("Expected 0 GetCallerIdentity calls after reset, got %d", len(mock.GetCallerIdentityCalls))
	}
	if len(mock.UpdateKubeconfigCalls) != 0 {
		t.Errorf("Expected 0 UpdateKubeconfig calls after reset, got %d", len(mock.UpdateKubeconfigCalls))
	}
	if len(mock.GetConnectionCommandsCalls) != 0 {
		t.Errorf("Expected 0 GetConnectionCommands calls after reset, got %d", len(mock.GetConnectionCommandsCalls))
	}
}

func TestMockAzureProvider_MultipleCalls(t *testing.T) {
	mock := NewMockAzureProvider()
	ctx := context.Background()

	// Make multiple calls to the same method
	_ = mock.ValidateCredentials(ctx)
	_ = mock.ValidateCredentials(ctx)
	_ = mock.ValidateCredentials(ctx)

	// Verify all calls were recorded
	if len(mock.ValidateCredentialsCalls) != 3 {
		t.Errorf("Expected 3 ValidateCredentials calls, got %d", len(mock.ValidateCredentialsCalls))
	}

	// Verify call count helper
	if mock.GetValidateCredentialsCallCount() != 3 {
		t.Errorf("GetValidateCredentialsCallCount() = %d, want 3", mock.GetValidateCredentialsCallCount())
	}
}

func TestMockAzureProvider_Timestamps(t *testing.T) {
	mock := NewMockAzureProvider()
	ctx := context.Background()

	// Make a call
	_ = mock.ValidateCredentials(ctx)

	// Verify timestamp was recorded
	if mock.ValidateCredentialsCalls[0].Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}
}
