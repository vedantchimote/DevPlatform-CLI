package helm

import (
	"context"
	"sync"
	"time"
)

// MockHelmClient is a mock implementation of HelmClient for testing
type MockHelmClient struct {
	// Function fields for configuring mock behavior
	InstallFunc   func(ctx context.Context, opts InstallOptions) error
	UpgradeFunc   func(ctx context.Context, opts UpgradeOptions) error
	UninstallFunc func(ctx context.Context, opts UninstallOptions) error
	StatusFunc    func(ctx context.Context, releaseName, namespace string) (*ReleaseStatus, error)
	ListFunc      func(ctx context.Context, namespace string) ([]*Release, error)

	// Call tracking
	InstallCalls   []MockCall
	UpgradeCalls   []MockCall
	UninstallCalls []MockCall
	StatusCalls    []MockCall
	ListCalls      []MockCall

	// Mutex for thread-safe call tracking
	mu sync.Mutex
}

// MockCall represents a recorded method call on a mock
type MockCall struct {
	Args      []interface{}
	Timestamp time.Time
}

// NewMockHelmClient creates a new mock Helm client with default behavior
func NewMockHelmClient() *MockHelmClient {
	return &MockHelmClient{
		InstallCalls:   make([]MockCall, 0),
		UpgradeCalls:   make([]MockCall, 0),
		UninstallCalls: make([]MockCall, 0),
		StatusCalls:    make([]MockCall, 0),
		ListCalls:      make([]MockCall, 0),
	}
}

// Install implements HelmClient.Install
func (m *MockHelmClient) Install(ctx context.Context, opts InstallOptions) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.InstallCalls = append(m.InstallCalls, MockCall{
		Args:      []interface{}{ctx, opts},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.InstallFunc != nil {
		return m.InstallFunc(ctx, opts)
	}

	// Default behavior: return nil (success)
	return nil
}

// Upgrade implements HelmClient.Upgrade
func (m *MockHelmClient) Upgrade(ctx context.Context, opts UpgradeOptions) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.UpgradeCalls = append(m.UpgradeCalls, MockCall{
		Args:      []interface{}{ctx, opts},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.UpgradeFunc != nil {
		return m.UpgradeFunc(ctx, opts)
	}

	// Default behavior: return nil (success)
	return nil
}

// Uninstall implements HelmClient.Uninstall
func (m *MockHelmClient) Uninstall(ctx context.Context, opts UninstallOptions) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.UninstallCalls = append(m.UninstallCalls, MockCall{
		Args:      []interface{}{ctx, opts},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.UninstallFunc != nil {
		return m.UninstallFunc(ctx, opts)
	}

	// Default behavior: return nil (success)
	return nil
}

// Status implements HelmClient.Status
func (m *MockHelmClient) Status(ctx context.Context, releaseName, namespace string) (*ReleaseStatus, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.StatusCalls = append(m.StatusCalls, MockCall{
		Args:      []interface{}{ctx, releaseName, namespace},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.StatusFunc != nil {
		return m.StatusFunc(ctx, releaseName, namespace)
	}

	// Default behavior: return deployed status
	return &ReleaseStatus{
		Name:      releaseName,
		Namespace: namespace,
		Status:    "deployed",
		Revision:  1,
		Updated:   time.Now(),
	}, nil
}

// List implements HelmClient.List
func (m *MockHelmClient) List(ctx context.Context, namespace string) ([]*Release, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.ListCalls = append(m.ListCalls, MockCall{
		Args:      []interface{}{ctx, namespace},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.ListFunc != nil {
		return m.ListFunc(ctx, namespace)
	}

	// Default behavior: return empty list
	return []*Release{}, nil
}

// Reset clears all recorded calls
func (m *MockHelmClient) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.InstallCalls = make([]MockCall, 0)
	m.UpgradeCalls = make([]MockCall, 0)
	m.UninstallCalls = make([]MockCall, 0)
	m.StatusCalls = make([]MockCall, 0)
	m.ListCalls = make([]MockCall, 0)
}

// GetInstallCallCount returns the number of times Install was called
func (m *MockHelmClient) GetInstallCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.InstallCalls)
}

// GetUpgradeCallCount returns the number of times Upgrade was called
func (m *MockHelmClient) GetUpgradeCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.UpgradeCalls)
}

// GetUninstallCallCount returns the number of times Uninstall was called
func (m *MockHelmClient) GetUninstallCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.UninstallCalls)
}

// GetStatusCallCount returns the number of times Status was called
func (m *MockHelmClient) GetStatusCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.StatusCalls)
}

// GetListCallCount returns the number of times List was called
func (m *MockHelmClient) GetListCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.ListCalls)
}
