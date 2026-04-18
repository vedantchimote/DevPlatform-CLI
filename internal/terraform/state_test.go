package terraform

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/devplatform/devplatform-cli/internal/provider/types"
	"github.com/devplatform/devplatform-cli/test/testutil"
)

// TestNewStateManager tests creating a new state manager
func TestNewStateManager(t *testing.T) {
	mockExecutor := NewMockTerraformExecutor()
	sm := NewStateManager(mockExecutor)
	
	testutil.AssertTrue(t, sm != nil, "StateManager should not be nil")
	testutil.AssertTrue(t, sm.executor != nil, "Executor should be set")
}

// TestGenerateS3BackendConfig tests generating S3 backend configuration
func TestGenerateS3BackendConfig(t *testing.T) {
	mockExecutor := NewMockTerraformExecutor()
	sm := NewStateManager(mockExecutor)
	
	config := map[string]string{
		"bucket":         "my-terraform-state",
		"key":            "app/terraform.tfstate",
		"region":         "us-east-1",
		"dynamodb_table": "terraform-locks",
		"encrypt":        "true",
	}
	
	result := sm.generateS3BackendConfig(config)
	
	testutil.AssertContains(t, result, "backend \"s3\"")
	testutil.AssertContains(t, result, "my-terraform-state")
	testutil.AssertContains(t, result, "app/terraform.tfstate")
	testutil.AssertContains(t, result, "us-east-1")
	testutil.AssertContains(t, result, "terraform-locks")
	testutil.AssertContains(t, result, "true")
}

// TestGenerateAzureBackendConfig tests generating Azure backend configuration
func TestGenerateAzureBackendConfig(t *testing.T) {
	mockExecutor := NewMockTerraformExecutor()
	sm := NewStateManager(mockExecutor)
	
	config := map[string]string{
		"storage_account_name": "mystorageaccount",
		"container_name":       "tfstate",
		"key":                  "app.terraform.tfstate",
		"resource_group_name":  "my-rg",
		"subscription_id":      "sub-123",
		"tenant_id":            "tenant-456",
	}
	
	result := sm.generateAzureBackendConfig(config)
	
	testutil.AssertContains(t, result, "backend \"azurerm\"")
	testutil.AssertContains(t, result, "mystorageaccount")
	testutil.AssertContains(t, result, "tfstate")
	testutil.AssertContains(t, result, "app.terraform.tfstate")
	testutil.AssertContains(t, result, "my-rg")
	testutil.AssertContains(t, result, "sub-123")
	testutil.AssertContains(t, result, "tenant-456")
}

// TestConfigureBackendS3 tests configuring S3 backend
func TestConfigureBackendS3(t *testing.T) {
	mockExecutor := NewMockTerraformExecutor()
	sm := NewStateManager(mockExecutor)
	
	tmpDir := t.TempDir()
	
	backend := &types.TerraformBackend{
		Type: "s3",
		Config: map[string]string{
			"bucket":         "my-bucket",
			"key":            "terraform.tfstate",
			"region":         "us-west-2",
			"dynamodb_table": "locks",
			"encrypt":        "true",
		},
	}
	
	err := sm.ConfigureBackend(context.Background(), tmpDir, backend)
	testutil.AssertNoError(t, err)
	
	// Verify backend file was created
	backendFile := filepath.Join(tmpDir, "backend.tf")
	content, err := os.ReadFile(backendFile)
	testutil.AssertNoError(t, err)
	
	contentStr := string(content)
	testutil.AssertContains(t, contentStr, "backend \"s3\"")
	testutil.AssertContains(t, contentStr, "my-bucket")
}

// TestConfigureBackendAzure tests configuring Azure backend
func TestConfigureBackendAzure(t *testing.T) {
	mockExecutor := NewMockTerraformExecutor()
	sm := NewStateManager(mockExecutor)
	
	tmpDir := t.TempDir()
	
	backend := &types.TerraformBackend{
		Type: "azurerm",
		Config: map[string]string{
			"storage_account_name": "myaccount",
			"container_name":       "tfstate",
			"key":                  "terraform.tfstate",
			"resource_group_name":  "my-rg",
			"subscription_id":      "sub-123",
			"tenant_id":            "tenant-456",
		},
	}
	
	err := sm.ConfigureBackend(context.Background(), tmpDir, backend)
	testutil.AssertNoError(t, err)
	
	// Verify backend file was created
	backendFile := filepath.Join(tmpDir, "backend.tf")
	content, err := os.ReadFile(backendFile)
	testutil.AssertNoError(t, err)
	
	contentStr := string(content)
	testutil.AssertContains(t, contentStr, "backend \"azurerm\"")
	testutil.AssertContains(t, contentStr, "myaccount")
}

// TestStateExistsLocalFile tests checking if state exists locally
func TestStateExistsLocalFile(t *testing.T) {
	mockExecutor := NewMockTerraformExecutor()
	sm := NewStateManager(mockExecutor)
	
	tmpDir := t.TempDir()
	
	// Create a state file
	stateFile := filepath.Join(tmpDir, "terraform.tfstate")
	err := os.WriteFile(stateFile, []byte("{}"), 0644)
	testutil.AssertNoError(t, err)
	
	exists, err := sm.StateExists(context.Background(), tmpDir)
	testutil.AssertNoError(t, err)
	testutil.AssertTrue(t, exists, "State should exist")
}

// TestStateExistsNoState tests checking when no state exists
func TestStateExistsNoState(t *testing.T) {
	mockExecutor := NewMockTerraformExecutor()
	mockExecutor.OutputFunc = func(ctx context.Context, workingDir, outputName string) (string, error) {
		return "", &TerraformError{Message: "No state file found"}
	}
	
	sm := NewStateManager(mockExecutor)
	tmpDir := t.TempDir()
	
	exists, err := sm.StateExists(context.Background(), tmpDir)
	testutil.AssertNoError(t, err)
	testutil.AssertFalse(t, exists, "State should not exist")
}

// TestGetStateKey tests generating state keys
func TestGetStateKey(t *testing.T) {
	mockExecutor := NewMockTerraformExecutor()
	sm := NewStateManager(mockExecutor)
	
	tests := []struct {
		name          string
		appName       string
		envType       string
		cloudProvider string
		expected      string
	}{
		{
			name:          "AWS dev",
			appName:       "myapp",
			envType:       "dev",
			cloudProvider: "aws",
			expected:      "aws/myapp/dev/terraform.tfstate",
		},
		{
			name:          "Azure production",
			appName:       "webapp",
			envType:       "production",
			cloudProvider: "azure",
			expected:      "azure/webapp/production/terraform.tfstate",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sm.GetStateKey(tt.appName, tt.envType, tt.cloudProvider)
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

// TestValidateStateLock tests validating state lock
func TestValidateStateLock(t *testing.T) {
	mockExecutor := NewMockTerraformExecutor()
	mockExecutor.InitFunc = func(ctx context.Context, workingDir string) error {
		return nil
	}
	
	sm := NewStateManager(mockExecutor)
	tmpDir := t.TempDir()
	
	err := sm.ValidateStateLock(context.Background(), tmpDir)
	testutil.AssertNoError(t, err)
}

// TestValidateStateLockLocked tests validating when state is locked
func TestValidateStateLockLocked(t *testing.T) {
	mockExecutor := NewMockTerraformExecutor()
	mockExecutor.InitFunc = func(ctx context.Context, workingDir string) error {
		return &TerraformError{Message: "state lock acquired by another process"}
	}
	
	sm := NewStateManager(mockExecutor)
	tmpDir := t.TempDir()
	
	err := sm.ValidateStateLock(context.Background(), tmpDir)
	testutil.AssertError(t, err)
	testutil.AssertContains(t, err.Error(), "state is locked")
}

// TestForceUnlock tests force unlock
func TestForceUnlock(t *testing.T) {
	mockExecutor := NewMockTerraformExecutor()
	sm := NewStateManager(mockExecutor)
	
	err := sm.ForceUnlock(context.Background(), "/tmp", "lock-123")
	testutil.AssertError(t, err)
	testutil.AssertContains(t, err.Error(), "force unlock must be performed manually")
	testutil.AssertContains(t, err.Error(), "lock-123")
}

// TestGetStateLockInfo tests getting lock info from error
func TestGetStateLockInfo(t *testing.T) {
	mockExecutor := NewMockTerraformExecutor()
	sm := NewStateManager(mockExecutor)
	
	tests := []struct {
		name     string
		err      error
		contains string
	}{
		{
			name:     "error with lock info",
			err:      &TerraformError{Message: "Lock Info: ID=123, Operation=apply"},
			contains: "Lock Info:",
		},
		{
			name:     "error without lock info",
			err:      &TerraformError{Message: "some other error"},
			contains: "State is locked",
		},
		{
			name:     "nil error",
			err:      nil,
			contains: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sm.GetStateLockInfo(tt.err)
			if tt.contains != "" {
				testutil.AssertContains(t, result, tt.contains)
			} else {
				testutil.AssertEqual(t, "", result)
			}
		})
	}
}

// TestContains tests the contains helper function
func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{
			name:     "contains at start",
			s:        "hello world",
			substr:   "hello",
			expected: true,
		},
		{
			name:     "contains at end",
			s:        "hello world",
			substr:   "world",
			expected: true,
		},
		{
			name:     "contains in middle",
			s:        "hello world",
			substr:   "lo wo",
			expected: true,
		},
		{
			name:     "does not contain",
			s:        "hello world",
			substr:   "goodbye",
			expected: false,
		},
		{
			name:     "exact match",
			s:        "hello",
			substr:   "hello",
			expected: true,
		},
		{
			name:     "empty substring",
			s:        "hello",
			substr:   "",
			expected: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.s, tt.substr)
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}
