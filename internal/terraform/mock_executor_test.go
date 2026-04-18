package terraform

import (
	"context"
	"errors"
	"testing"
)

func TestMockTerraformExecutor_Init(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()
	workingDir := "/test/dir"

	// Test default behavior (success)
	err := mock.Init(ctx, workingDir)
	if err != nil {
		t.Errorf("Init() unexpected error = %v", err)
	}

	// Verify call was recorded
	if len(mock.InitCalls) != 1 {
		t.Errorf("Expected 1 Init call, got %d", len(mock.InitCalls))
	}

	// Verify call count helper
	if mock.GetInitCallCount() != 1 {
		t.Errorf("GetInitCallCount() = %d, want 1", mock.GetInitCallCount())
	}

	// Verify call arguments
	if len(mock.InitCalls[0].Args) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(mock.InitCalls[0].Args))
	}
	if mock.InitCalls[0].Args[1] != workingDir {
		t.Errorf("Expected workingDir %s, got %v", workingDir, mock.InitCalls[0].Args[1])
	}
}

func TestMockTerraformExecutor_InitWithCustomFunc(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()
	workingDir := "/test/dir"
	expectedErr := errors.New("init failed")

	// Configure custom behavior
	mock.InitFunc = func(ctx context.Context, workingDir string) error {
		return expectedErr
	}

	// Test custom behavior
	err := mock.Init(ctx, workingDir)
	if err != expectedErr {
		t.Errorf("Init() error = %v, want %v", err, expectedErr)
	}

	// Verify call was still recorded
	if len(mock.InitCalls) != 1 {
		t.Errorf("Expected 1 Init call, got %d", len(mock.InitCalls))
	}
}

func TestMockTerraformExecutor_Plan(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()
	workingDir := "/test/dir"
	varFile := "vars.tfvars"

	// Test default behavior
	output, err := mock.Plan(ctx, workingDir, varFile)
	if err != nil {
		t.Errorf("Plan() unexpected error = %v", err)
	}
	if output == "" {
		t.Error("Plan() returned empty output")
	}

	// Verify call was recorded
	if len(mock.PlanCalls) != 1 {
		t.Errorf("Expected 1 Plan call, got %d", len(mock.PlanCalls))
	}

	// Verify call count helper
	if mock.GetPlanCallCount() != 1 {
		t.Errorf("GetPlanCallCount() = %d, want 1", mock.GetPlanCallCount())
	}
}

func TestMockTerraformExecutor_PlanWithCustomFunc(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()
	workingDir := "/test/dir"
	varFile := "vars.tfvars"
	expectedOutput := "custom plan output"

	// Configure custom behavior
	mock.PlanFunc = func(ctx context.Context, workingDir string, varFile string) (string, error) {
		return expectedOutput, nil
	}

	// Test custom behavior
	output, err := mock.Plan(ctx, workingDir, varFile)
	if err != nil {
		t.Errorf("Plan() unexpected error = %v", err)
	}
	if output != expectedOutput {
		t.Errorf("Plan() output = %s, want %s", output, expectedOutput)
	}
}

func TestMockTerraformExecutor_Apply(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()
	workingDir := "/test/dir"
	varFile := "vars.tfvars"
	autoApprove := true

	// Test default behavior
	err := mock.Apply(ctx, workingDir, varFile, autoApprove)
	if err != nil {
		t.Errorf("Apply() unexpected error = %v", err)
	}

	// Verify call was recorded
	if len(mock.ApplyCalls) != 1 {
		t.Errorf("Expected 1 Apply call, got %d", len(mock.ApplyCalls))
	}

	// Verify call count helper
	if mock.GetApplyCallCount() != 1 {
		t.Errorf("GetApplyCallCount() = %d, want 1", mock.GetApplyCallCount())
	}

	// Verify autoApprove argument
	if mock.ApplyCalls[0].Args[3] != autoApprove {
		t.Errorf("Expected autoApprove %v, got %v", autoApprove, mock.ApplyCalls[0].Args[3])
	}
}

func TestMockTerraformExecutor_Destroy(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()
	workingDir := "/test/dir"
	varFile := "vars.tfvars"
	autoApprove := true

	// Test default behavior
	err := mock.Destroy(ctx, workingDir, varFile, autoApprove)
	if err != nil {
		t.Errorf("Destroy() unexpected error = %v", err)
	}

	// Verify call was recorded
	if len(mock.DestroyCalls) != 1 {
		t.Errorf("Expected 1 Destroy call, got %d", len(mock.DestroyCalls))
	}

	// Verify call count helper
	if mock.GetDestroyCallCount() != 1 {
		t.Errorf("GetDestroyCallCount() = %d, want 1", mock.GetDestroyCallCount())
	}
}

func TestMockTerraformExecutor_Output(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()
	workingDir := "/test/dir"
	outputName := "cluster_endpoint"

	// Test default behavior
	output, err := mock.Output(ctx, workingDir, outputName)
	if err != nil {
		t.Errorf("Output() unexpected error = %v", err)
	}

	// Verify call was recorded
	if len(mock.OutputCalls) != 1 {
		t.Errorf("Expected 1 Output call, got %d", len(mock.OutputCalls))
	}

	// Verify call count helper
	if mock.GetOutputCallCount() != 1 {
		t.Errorf("GetOutputCallCount() = %d, want 1", mock.GetOutputCallCount())
	}

	// Verify output name argument
	if mock.OutputCalls[0].Args[2] != outputName {
		t.Errorf("Expected outputName %s, got %v", outputName, mock.OutputCalls[0].Args[2])
	}

	// Test with custom function
	expectedOutput := "https://test-cluster.example.com"
	mock.OutputFunc = func(ctx context.Context, workingDir string, outputName string) (string, error) {
		return expectedOutput, nil
	}

	output, err = mock.Output(ctx, workingDir, outputName)
	if err != nil {
		t.Errorf("Output() unexpected error = %v", err)
	}
	if output != expectedOutput {
		t.Errorf("Output() = %s, want %s", output, expectedOutput)
	}
}

func TestMockTerraformExecutor_Reset(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()

	// Make several calls
	_ = mock.Init(ctx, "/test/dir")
	_, _ = mock.Plan(ctx, "/test/dir", "vars.tfvars")
	_ = mock.Apply(ctx, "/test/dir", "vars.tfvars", true)

	// Verify calls were recorded
	if len(mock.InitCalls) != 1 {
		t.Errorf("Expected 1 Init call before reset, got %d", len(mock.InitCalls))
	}
	if len(mock.PlanCalls) != 1 {
		t.Errorf("Expected 1 Plan call before reset, got %d", len(mock.PlanCalls))
	}
	if len(mock.ApplyCalls) != 1 {
		t.Errorf("Expected 1 Apply call before reset, got %d", len(mock.ApplyCalls))
	}

	// Reset the mock
	mock.Reset()

	// Verify all calls were cleared
	if len(mock.InitCalls) != 0 {
		t.Errorf("Expected 0 Init calls after reset, got %d", len(mock.InitCalls))
	}
	if len(mock.PlanCalls) != 0 {
		t.Errorf("Expected 0 Plan calls after reset, got %d", len(mock.PlanCalls))
	}
	if len(mock.ApplyCalls) != 0 {
		t.Errorf("Expected 0 Apply calls after reset, got %d", len(mock.ApplyCalls))
	}
	if len(mock.DestroyCalls) != 0 {
		t.Errorf("Expected 0 Destroy calls after reset, got %d", len(mock.DestroyCalls))
	}
	if len(mock.OutputCalls) != 0 {
		t.Errorf("Expected 0 Output calls after reset, got %d", len(mock.OutputCalls))
	}
}

func TestMockTerraformExecutor_MultipleCalls(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()

	// Make multiple calls to the same method
	_ = mock.Init(ctx, "/test/dir1")
	_ = mock.Init(ctx, "/test/dir2")
	_ = mock.Init(ctx, "/test/dir3")

	// Verify all calls were recorded
	if len(mock.InitCalls) != 3 {
		t.Errorf("Expected 3 Init calls, got %d", len(mock.InitCalls))
	}

	// Verify call count helper
	if mock.GetInitCallCount() != 3 {
		t.Errorf("GetInitCallCount() = %d, want 3", mock.GetInitCallCount())
	}

	// Verify each call has correct arguments
	expectedDirs := []string{"/test/dir1", "/test/dir2", "/test/dir3"}
	for i, call := range mock.InitCalls {
		if call.Args[1] != expectedDirs[i] {
			t.Errorf("Call %d: expected workingDir %s, got %v", i, expectedDirs[i], call.Args[1])
		}
	}
}

func TestMockTerraformExecutor_Timestamps(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()

	// Make a call
	_ = mock.Init(ctx, "/test/dir")

	// Verify timestamp was recorded
	if mock.InitCalls[0].Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}
}
