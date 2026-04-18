package mocks

import (
	"context"
	"sync"
	"time"

	"github.com/devplatform/devplatform-cli/internal/provider/types"
)

// MockAzureProvider is a mock implementation of CloudProvider for Azure testing
type MockAzureProvider struct {
	// Function fields for configuring mock behavior
	ValidateCredentialsFunc  func(ctx context.Context) error
	GetCallerIdentityFunc    func(ctx context.Context) (*types.CallerIdentity, error)
	UpdateKubeconfigFunc     func(clusterName string) error
	GetConnectionCommandsFunc func(clusterName string, namespace string) []string
	CalculateTotalCostFunc   func(envType string) (*types.EnvironmentCosts, error)
	GetTerraformBackendFunc  func(appName string, envType string) (*types.TerraformBackend, error)
	GetModulePathFunc        func() string
	GetProviderNameFunc      func() string

	// Call tracking
	ValidateCredentialsCalls  []MockCall
	GetCallerIdentityCalls    []MockCall
	UpdateKubeconfigCalls     []MockCall
	GetConnectionCommandsCalls []MockCall
	CalculateTotalCostCalls   []MockCall
	GetTerraformBackendCalls  []MockCall
	GetModulePathCalls        []MockCall
	GetProviderNameCalls      []MockCall

	// Mutex for thread-safe call tracking
	mu sync.Mutex
}

// NewMockAzureProvider creates a new mock Azure provider with default behavior
func NewMockAzureProvider() *MockAzureProvider {
	return &MockAzureProvider{
		ValidateCredentialsCalls:  make([]MockCall, 0),
		GetCallerIdentityCalls:    make([]MockCall, 0),
		UpdateKubeconfigCalls:     make([]MockCall, 0),
		GetConnectionCommandsCalls: make([]MockCall, 0),
		CalculateTotalCostCalls:   make([]MockCall, 0),
		GetTerraformBackendCalls:  make([]MockCall, 0),
		GetModulePathCalls:        make([]MockCall, 0),
		GetProviderNameCalls:      make([]MockCall, 0),
	}
}

// ValidateCredentials implements CloudProvider.ValidateCredentials
func (m *MockAzureProvider) ValidateCredentials(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.ValidateCredentialsCalls = append(m.ValidateCredentialsCalls, MockCall{
		Args:      []interface{}{ctx},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.ValidateCredentialsFunc != nil {
		return m.ValidateCredentialsFunc(ctx)
	}

	// Default behavior: return nil (success)
	return nil
}

// GetCallerIdentity implements CloudProvider.GetCallerIdentity
func (m *MockAzureProvider) GetCallerIdentity(ctx context.Context) (*types.CallerIdentity, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.GetCallerIdentityCalls = append(m.GetCallerIdentityCalls, MockCall{
		Args:      []interface{}{ctx},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.GetCallerIdentityFunc != nil {
		return m.GetCallerIdentityFunc(ctx)
	}

	// Default behavior: return mock identity
	return &types.CallerIdentity{
		Account: "12345678-1234-1234-1234-123456789012",
		Arn:     "00000000-0000-0000-0000-000000000000",
		UserId:  "Authenticated",
	}, nil
}

// UpdateKubeconfig implements CloudProvider.UpdateKubeconfig
func (m *MockAzureProvider) UpdateKubeconfig(clusterName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.UpdateKubeconfigCalls = append(m.UpdateKubeconfigCalls, MockCall{
		Args:      []interface{}{clusterName},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.UpdateKubeconfigFunc != nil {
		return m.UpdateKubeconfigFunc(clusterName)
	}

	// Default behavior: return nil (success)
	return nil
}

// GetConnectionCommands implements CloudProvider.GetConnectionCommands
func (m *MockAzureProvider) GetConnectionCommands(clusterName string, namespace string) []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.GetConnectionCommandsCalls = append(m.GetConnectionCommandsCalls, MockCall{
		Args:      []interface{}{clusterName, namespace},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.GetConnectionCommandsFunc != nil {
		return m.GetConnectionCommandsFunc(clusterName, namespace)
	}

	// Default behavior: return sample commands
	return []string{
		"az aks get-credentials --name " + clusterName + " --resource-group devplatform-rg",
		"kubectl config set-context --current --namespace=" + namespace,
		"kubectl get pods",
	}
}

// CalculateTotalCost implements CloudProvider.CalculateTotalCost
func (m *MockAzureProvider) CalculateTotalCost(envType string) (*types.EnvironmentCosts, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.CalculateTotalCostCalls = append(m.CalculateTotalCostCalls, MockCall{
		Args:      []interface{}{envType},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.CalculateTotalCostFunc != nil {
		return m.CalculateTotalCostFunc(envType)
	}

	// Default behavior: return mock costs
	costs := map[string]float64{
		"dev":     140.0,
		"staging": 280.0,
		"prod":    560.0,
	}

	totalCost := costs[envType]
	if totalCost == 0 {
		totalCost = 140.0 // default to dev cost
	}

	return &types.EnvironmentCosts{
		NetworkCost:  totalCost * 0.1,
		DatabaseCost: totalCost * 0.4,
		K8sCost:      totalCost * 0.5,
		TotalCost:    totalCost,
		Environment:  envType,
		Provider:     "azure",
	}, nil
}

// GetTerraformBackend implements CloudProvider.GetTerraformBackend
func (m *MockAzureProvider) GetTerraformBackend(appName string, envType string) (*types.TerraformBackend, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.GetTerraformBackendCalls = append(m.GetTerraformBackendCalls, MockCall{
		Args:      []interface{}{appName, envType},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.GetTerraformBackendFunc != nil {
		return m.GetTerraformBackendFunc(appName, envType)
	}

	// Default behavior: return Azure Storage backend config
	return &types.TerraformBackend{
		Type: "azurerm",
		Config: map[string]string{
			"storage_account_name": "devplatformtfeastus",
			"container_name":       "terraform-state",
			"key":                  appName + "/" + envType + "/terraform.tfstate",
			"resource_group_name":  "devplatform-rg",
		},
	}, nil
}

// GetModulePath implements CloudProvider.GetModulePath
func (m *MockAzureProvider) GetModulePath() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.GetModulePathCalls = append(m.GetModulePathCalls, MockCall{
		Args:      []interface{}{},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.GetModulePathFunc != nil {
		return m.GetModulePathFunc()
	}

	// Default behavior: return Azure module path
	return "terraform/modules/azure"
}

// GetProviderName implements CloudProvider.GetProviderName
func (m *MockAzureProvider) GetProviderName() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.GetProviderNameCalls = append(m.GetProviderNameCalls, MockCall{
		Args:      []interface{}{},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.GetProviderNameFunc != nil {
		return m.GetProviderNameFunc()
	}

	// Default behavior: return "azure"
	return "azure"
}

// Reset clears all recorded calls
func (m *MockAzureProvider) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ValidateCredentialsCalls = make([]MockCall, 0)
	m.GetCallerIdentityCalls = make([]MockCall, 0)
	m.UpdateKubeconfigCalls = make([]MockCall, 0)
	m.GetConnectionCommandsCalls = make([]MockCall, 0)
	m.CalculateTotalCostCalls = make([]MockCall, 0)
	m.GetTerraformBackendCalls = make([]MockCall, 0)
	m.GetModulePathCalls = make([]MockCall, 0)
	m.GetProviderNameCalls = make([]MockCall, 0)
}

// GetValidateCredentialsCallCount returns the number of times ValidateCredentials was called
func (m *MockAzureProvider) GetValidateCredentialsCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.ValidateCredentialsCalls)
}

// GetCallerIdentityCallCount returns the number of times GetCallerIdentity was called
func (m *MockAzureProvider) GetCallerIdentityCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.GetCallerIdentityCalls)
}

// GetUpdateKubeconfigCallCount returns the number of times UpdateKubeconfig was called
func (m *MockAzureProvider) GetUpdateKubeconfigCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.UpdateKubeconfigCalls)
}

// GetConnectionCommandsCallCount returns the number of times GetConnectionCommands was called
func (m *MockAzureProvider) GetConnectionCommandsCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.GetConnectionCommandsCalls)
}

// GetCalculateTotalCostCallCount returns the number of times CalculateTotalCost was called
func (m *MockAzureProvider) GetCalculateTotalCostCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.CalculateTotalCostCalls)
}

// GetTerraformBackendCallCount returns the number of times GetTerraformBackend was called
func (m *MockAzureProvider) GetTerraformBackendCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.GetTerraformBackendCalls)
}

// GetModulePathCallCount returns the number of times GetModulePath was called
func (m *MockAzureProvider) GetModulePathCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.GetModulePathCalls)
}

// GetProviderNameCallCount returns the number of times GetProviderName was called
func (m *MockAzureProvider) GetProviderNameCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.GetProviderNameCalls)
}
