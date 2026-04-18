package terraform

import (
	"context"
	"sync"
	"time"
)

// MockCall represents a recorded method call on a mock
type MockCall struct {
	Args      []interface{}
	Timestamp time.Time
}

// MockTerraformExecutor is a mock implementation of TerraformExecutor for testing
type MockTerraformExecutor struct {
	// Function fields for configuring mock behavior
	InitFunc    func(ctx context.Context, workingDir string) error
	PlanFunc    func(ctx context.Context, workingDir string, varFile string) (string, error)
	ApplyFunc   func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error
	DestroyFunc func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error
	OutputFunc  func(ctx context.Context, workingDir string, outputName string) (string, error)

	// Call tracking
	InitCalls    []MockCall
	PlanCalls    []MockCall
	ApplyCalls   []MockCall
	DestroyCalls []MockCall
	OutputCalls  []MockCall

	// Mutex for thread-safe call tracking
	mu sync.Mutex
}

// NewMockTerraformExecutor creates a new mock Terraform executor with default behavior
func NewMockTerraformExecutor() *MockTerraformExecutor {
	return &MockTerraformExecutor{
		InitCalls:    make([]MockCall, 0),
		PlanCalls:    make([]MockCall, 0),
		ApplyCalls:   make([]MockCall, 0),
		DestroyCalls: make([]MockCall, 0),
		OutputCalls:  make([]MockCall, 0),
	}
}

// Init implements TerraformExecutor.Init
func (m *MockTerraformExecutor) Init(ctx context.Context, workingDir string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.InitCalls = append(m.InitCalls, MockCall{
		Args:      []interface{}{ctx, workingDir},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.InitFunc != nil {
		return m.InitFunc(ctx, workingDir)
	}

	// Default behavior: return nil (success)
	return nil
}

// Plan implements TerraformExecutor.Plan
func (m *MockTerraformExecutor) Plan(ctx context.Context, workingDir string, varFile string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.PlanCalls = append(m.PlanCalls, MockCall{
		Args:      []interface{}{ctx, workingDir, varFile},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.PlanFunc != nil {
		return m.PlanFunc(ctx, workingDir, varFile)
	}

	// Default behavior: return empty plan output
	return "No changes. Infrastructure is up-to-date.", nil
}

// Apply implements TerraformExecutor.Apply
func (m *MockTerraformExecutor) Apply(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.ApplyCalls = append(m.ApplyCalls, MockCall{
		Args:      []interface{}{ctx, workingDir, varFile, autoApprove},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.ApplyFunc != nil {
		return m.ApplyFunc(ctx, workingDir, varFile, autoApprove)
	}

	// Default behavior: return nil (success)
	return nil
}

// Destroy implements TerraformExecutor.Destroy
func (m *MockTerraformExecutor) Destroy(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.DestroyCalls = append(m.DestroyCalls, MockCall{
		Args:      []interface{}{ctx, workingDir, varFile, autoApprove},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.DestroyFunc != nil {
		return m.DestroyFunc(ctx, workingDir, varFile, autoApprove)
	}

	// Default behavior: return nil (success)
	return nil
}

// Output implements TerraformExecutor.Output
func (m *MockTerraformExecutor) Output(ctx context.Context, workingDir string, outputName string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.OutputCalls = append(m.OutputCalls, MockCall{
		Args:      []interface{}{ctx, workingDir, outputName},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.OutputFunc != nil {
		return m.OutputFunc(ctx, workingDir, outputName)
	}

	// Default behavior: return empty string
	return "", nil
}

// Reset clears all recorded calls
func (m *MockTerraformExecutor) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.InitCalls = make([]MockCall, 0)
	m.PlanCalls = make([]MockCall, 0)
	m.ApplyCalls = make([]MockCall, 0)
	m.DestroyCalls = make([]MockCall, 0)
	m.OutputCalls = make([]MockCall, 0)
}

// GetInitCallCount returns the number of times Init was called
func (m *MockTerraformExecutor) GetInitCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.InitCalls)
}

// GetPlanCallCount returns the number of times Plan was called
func (m *MockTerraformExecutor) GetPlanCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.PlanCalls)
}

// GetApplyCallCount returns the number of times Apply was called
func (m *MockTerraformExecutor) GetApplyCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.ApplyCalls)
}

// GetDestroyCallCount returns the number of times Destroy was called
func (m *MockTerraformExecutor) GetDestroyCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.DestroyCalls)
}

// GetOutputCallCount returns the number of times Output was called
func (m *MockTerraformExecutor) GetOutputCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.OutputCalls)
}
