package mocks

import (
	"context"
	"errors"
	"testing"

	"github.com/devplatform/devplatform-cli/internal/provider/types"
)

func TestMockAWSProvider_ValidateCredentials(t *testing.T) {
	mock := NewMockAWSProvider()
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

func TestMockAWSProvider_ValidateCredentialsWithCustomFunc(t *testing.T) {
	mock := NewMockAWSProvider()
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

func TestMockAWSProvider_GetCallerIdentity(t *testing.T) {
	mock := NewMockAWSProvider()
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
	if identity.Arn == "" {
		t.Error("GetCallerIdentity() returned empty Arn")
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

func TestMockAWSProvider_GetCallerIdentityWithCustomFunc(t *testing.T) {
	mock := NewMockAWSProvider()
	ctx := context.Background()
	expectedIdentity := &types.CallerIdentity{
		Account: "999888777666",
		Arn:     "arn:aws:iam::999888777666:user/custom-user",
		UserId:  "CUSTOMUSERID",
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

func TestMockAWSProvider_UpdateKubeconfig(t *testing.T) {
	mock := NewMockAWSProvider()
	clusterName := "test-cluster"

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

func TestMockAWSProvider_GetConnectionCommands(t *testing.T) {
	mock := NewMockAWSProvider()
	clusterName := "test-cluster"
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

func TestMockAWSProvider_GetConnectionCommandsWithCustomFunc(t *testing.T) {
	mock := NewMockAWSProvider()
	clusterName := "test-cluster"
	namespace := "test-namespace"
	expectedCommands := []string{"custom command 1", "custom command 2"}

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

func TestMockAWSProvider_CalculateTotalCost(t *testing.T) {
	mock := NewMockAWSProvider()

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
			if costs.Provider != "aws" {
				t.Errorf("CalculateTotalCost().Provider = %s, want aws", costs.Provider)
			}
		})
	}

	// Verify call count
	if mock.GetCalculateTotalCostCallCount() != 3 {
		t.Errorf("GetCalculateTotalCostCallCount() = %d, want 3", mock.GetCalculateTotalCostCallCount())
	}
}

func TestMockAWSProvider_GetTerraformBackend(t *testing.T) {
	mock := NewMockAWSProvider()
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
	if backend.Type != "s3" {
		t.Errorf("GetTerraformBackend().Type = %s, want s3", backend.Type)
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

func TestMockAWSProvider_GetModulePath(t *testing.T) {
	mock := NewMockAWSProvider()

	// Test default behavior
	path := mock.GetModulePath()
	if path == "" {
		t.Error("GetModulePath() returned empty string")
	}
	if path != "terraform/modules/aws" {
		t.Errorf("GetModulePath() = %s, want terraform/modules/aws", path)
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

func TestMockAWSProvider_GetProviderName(t *testing.T) {
	mock := NewMockAWSProvider()

	// Test default behavior
	name := mock.GetProviderName()
	if name != "aws" {
		t.Errorf("GetProviderName() = %s, want aws", name)
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

func TestMockAWSProvider_Reset(t *testing.T) {
	mock := NewMockAWSProvider()
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

func TestMockAWSProvider_MultipleCalls(t *testing.T) {
	mock := NewMockAWSProvider()
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

func TestMockAWSProvider_Timestamps(t *testing.T) {
	mock := NewMockAWSProvider()
	ctx := context.Background()

	// Make a call
	_ = mock.ValidateCredentials(ctx)

	// Verify timestamp was recorded
	if mock.ValidateCredentialsCalls[0].Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}
}
