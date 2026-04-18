package terraform_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/devplatform/devplatform-cli/internal/terraform"
)

// ExampleMockTerraformExecutor_basicUsage demonstrates basic mock usage
func ExampleMockTerraformExecutor_basicUsage() {
	// Create a new mock
	mock := terraform.NewMockTerraformExecutor()

	// Use the mock with default behavior (all operations succeed)
	ctx := context.Background()
	err := mock.Init(ctx, "/test/dir")
	if err != nil {
		fmt.Printf("Init failed: %v\n", err)
		return
	}

	// Check how many times Init was called
	fmt.Printf("Init called %d times\n", mock.GetInitCallCount())

	// Output:
	// Init called 1 times
}

// ExampleMockTerraformExecutor_customBehavior demonstrates configuring custom behavior
func ExampleMockTerraformExecutor_customBehavior() {
	// Create a new mock
	mock := terraform.NewMockTerraformExecutor()

	// Configure custom behavior for Plan
	mock.PlanFunc = func(ctx context.Context, workingDir string, varFile string) (string, error) {
		return "Plan: 3 to add, 0 to change, 0 to destroy", nil
	}

	// Use the mock
	ctx := context.Background()
	output, err := mock.Plan(ctx, "/test/dir", "vars.tfvars")
	if err != nil {
		fmt.Printf("Plan failed: %v\n", err)
		return
	}

	fmt.Println(output)

	// Output:
	// Plan: 3 to add, 0 to change, 0 to destroy
}

// ExampleMockTerraformExecutor_errorSimulation demonstrates simulating errors
func ExampleMockTerraformExecutor_errorSimulation() {
	// Create a new mock
	mock := terraform.NewMockTerraformExecutor()

	// Configure Apply to return an error
	mock.ApplyFunc = func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
		return errors.New("terraform apply failed: resource already exists")
	}

	// Use the mock
	ctx := context.Background()
	err := mock.Apply(ctx, "/test/dir", "vars.tfvars", true)
	if err != nil {
		fmt.Printf("Apply failed: %v\n", err)
	}

	// Output:
	// Apply failed: terraform apply failed: resource already exists
}

// ExampleMockTerraformExecutor_callTracking demonstrates call tracking
func ExampleMockTerraformExecutor_callTracking() {
	// Create a new mock
	mock := terraform.NewMockTerraformExecutor()

	// Make multiple calls
	ctx := context.Background()
	_ = mock.Init(ctx, "/test/dir1")
	_ = mock.Init(ctx, "/test/dir2")
	_ = mock.Init(ctx, "/test/dir3")

	// Check call count
	fmt.Printf("Init called %d times\n", mock.GetInitCallCount())

	// Inspect call arguments
	for i, call := range mock.InitCalls {
		workingDir := call.Args[1].(string)
		fmt.Printf("Call %d: workingDir=%s\n", i+1, workingDir)
	}

	// Output:
	// Init called 3 times
	// Call 1: workingDir=/test/dir1
	// Call 2: workingDir=/test/dir2
	// Call 3: workingDir=/test/dir3
}

// ExampleMockTerraformExecutor_reset demonstrates resetting the mock
func ExampleMockTerraformExecutor_reset() {
	// Create a new mock
	mock := terraform.NewMockTerraformExecutor()

	// Make some calls
	ctx := context.Background()
	_ = mock.Init(ctx, "/test/dir")
	_, _ = mock.Plan(ctx, "/test/dir", "vars.tfvars")

	fmt.Printf("Before reset - Init: %d, Plan: %d\n", 
		mock.GetInitCallCount(), mock.GetPlanCallCount())

	// Reset the mock
	mock.Reset()

	fmt.Printf("After reset - Init: %d, Plan: %d\n", 
		mock.GetInitCallCount(), mock.GetPlanCallCount())

	// Output:
	// Before reset - Init: 1, Plan: 1
	// After reset - Init: 0, Plan: 0
}
