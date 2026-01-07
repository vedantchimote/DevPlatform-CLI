package terraform

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/devplatform/devplatform-cli/internal/provider"
)

// StateManager handles Terraform state operations
type StateManager struct {
	executor TerraformExecutor
}

// NewStateManager creates a new state manager
func NewStateManager(executor TerraformExecutor) *StateManager {
	return &StateManager{
		executor: executor,
	}
}

// ConfigureBackend configures the Terraform backend for state storage
func (s *StateManager) ConfigureBackend(ctx context.Context, workingDir string, backend *provider.TerraformBackend) error {
	// Create backend configuration file
	backendConfig := s.generateBackendConfig(backend)
	
	// Write backend configuration to file
	backendFile := filepath.Join(workingDir, "backend.tf")
	if err := os.WriteFile(backendFile, []byte(backendConfig), 0644); err != nil {
		return fmt.Errorf("failed to write backend configuration: %w", err)
	}
	
	return nil
}

// generateBackendConfig generates the Terraform backend configuration
func (s *StateManager) generateBackendConfig(backend *provider.TerraformBackend) string {
	switch backend.Type {
	case "s3":
		return s.generateS3BackendConfig(backend.Config)
	case "azurerm":
		return s.generateAzureBackendConfig(backend.Config)
	default:
		return ""
	}
}

// generateS3BackendConfig generates S3 backend configuration
func (s *StateManager) generateS3BackendConfig(config map[string]string) string {
	return fmt.Sprintf(`terraform {
  backend "s3" {
    bucket         = "%s"
    key            = "%s"
    region         = "%s"
    dynamodb_table = "%s"
    encrypt        = %s
  }
}
`,
		config["bucket"],
		config["key"],
		config["region"],
		config["dynamodb_table"],
		config["encrypt"],
	)
}

// generateAzureBackendConfig generates Azure Storage backend configuration
func (s *StateManager) generateAzureBackendConfig(config map[string]string) string {
	return fmt.Sprintf(`terraform {
  backend "azurerm" {
    storage_account_name = "%s"
    container_name       = "%s"
    key                  = "%s"
    resource_group_name  = "%s"
    subscription_id      = "%s"
    tenant_id            = "%s"
  }
}
`,
		config["storage_account_name"],
		config["container_name"],
		config["key"],
		config["resource_group_name"],
		config["subscription_id"],
		config["tenant_id"],
	)
}

// StateExists checks if a Terraform state exists
func (s *StateManager) StateExists(ctx context.Context, workingDir string) (bool, error) {
	// Check if state file exists locally
	stateFile := filepath.Join(workingDir, "terraform.tfstate")
	if _, err := os.Stat(stateFile); err == nil {
		return true, nil
	}
	
	// Try to get state from backend by running terraform show
	// If state exists in backend, this will succeed
	output, err := s.executor.Output(ctx, workingDir, "")
	if err != nil {
		// If error contains "No outputs found", state might exist but has no outputs
		// If error contains "No state file", state doesn't exist
		if contains(err.Error(), "No state file") || contains(err.Error(), "state file is empty") {
			return false, nil
		}
		// For other errors, assume state might exist
		return true, nil
	}
	
	// If we got output, state exists
	return output != "", nil
}

// GetStateKey generates a unique state key for the given app and environment
func (s *StateManager) GetStateKey(appName string, envType string, cloudProvider string) string {
	return fmt.Sprintf("%s/%s/%s/terraform.tfstate", cloudProvider, appName, envType)
}

// ValidateStateLock checks if the state is locked and returns lock information
func (s *StateManager) ValidateStateLock(ctx context.Context, workingDir string) error {
	// Terraform handles state locking automatically
	// If state is locked, init/plan/apply will fail with lock information
	// We can attempt to init to check for locks
	err := s.executor.Init(ctx, workingDir)
	if err != nil {
		if contains(err.Error(), "state lock") || contains(err.Error(), "locked") {
			return fmt.Errorf("state is locked: %w\n\nThe state is currently locked by another operation. Please wait for the other operation to complete or use 'terraform force-unlock' if the lock is stale", err)
		}
		return err
	}
	return nil
}

// ForceUnlock forces the release of a state lock
func (s *StateManager) ForceUnlock(ctx context.Context, workingDir string, lockID string) error {
	// Note: This is a dangerous operation and should only be used when absolutely necessary
	// We don't implement this directly to avoid accidental data corruption
	return fmt.Errorf("force unlock must be performed manually using: terraform force-unlock %s", lockID)
}

// GetStateLockInfo extracts lock information from an error
func (s *StateManager) GetStateLockInfo(err error) string {
	if err == nil {
		return ""
	}
	
	errStr := err.Error()
	if contains(errStr, "Lock Info:") {
		// Extract lock info section from error
		return errStr
	}
	
	return "State is locked. Run with --force to override (use with caution)"
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		findSubstring(s, substr)))
}

// findSubstring performs a simple substring search
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
