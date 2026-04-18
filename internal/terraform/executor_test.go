package terraform

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

// TestExecutorInit tests the Init method of the Executor
func TestExecutorInit(t *testing.T) {
	tests := []struct {
		name       string
		workingDir string
		mockFunc   func(ctx context.Context, workingDir string) error
		wantErr    bool
		errContains string
	}{
		{
			name:       "successful init",
			workingDir: "/tmp/terraform",
			mockFunc: func(ctx context.Context, workingDir string) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:       "init with empty working directory",
			workingDir: "",
			mockFunc: func(ctx context.Context, workingDir string) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:       "init failure - provider not found",
			workingDir: "/tmp/terraform",
			mockFunc: func(ctx context.Context, workingDir string) error {
				return errors.New("provider not found")
			},
			wantErr:     true,
			errContains: "provider not found",
		},
		{
			name:       "init failure - invalid configuration",
			workingDir: "/tmp/terraform",
			mockFunc: func(ctx context.Context, workingDir string) error {
				return errors.New("invalid configuration")
			},
			wantErr:     true,
			errContains: "invalid configuration",
		},
		{
			name:       "init with context cancellation",
			workingDir: "/tmp/terraform",
			mockFunc: func(ctx context.Context, workingDir string) error {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				return nil
			},
			wantErr:     true,
			errContains: "context canceled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockTerraformExecutor()
			mock.InitFunc = tt.mockFunc

			ctx := context.Background()
			if tt.errContains == "context canceled" {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel() // Cancel immediately
			}

			err := mock.Init(ctx, tt.workingDir)

			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				testutil.AssertContains(t, err.Error(), tt.errContains)
			}

			// Verify call was recorded
			testutil.AssertEqual(t, 1, len(mock.InitCalls))
			testutil.AssertEqual(t, tt.workingDir, mock.InitCalls[0].Args[1])
		})
	}
}

// TestExecutorInitCallRecording tests that Init calls are properly recorded
func TestExecutorInitCallRecording(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()

	// Make multiple calls
	_ = mock.Init(ctx, "/tmp/dir1")
	_ = mock.Init(ctx, "/tmp/dir2")
	_ = mock.Init(ctx, "/tmp/dir3")

	// Verify all calls were recorded
	testutil.AssertEqual(t, 3, len(mock.InitCalls))
	testutil.AssertEqual(t, "/tmp/dir1", mock.InitCalls[0].Args[1])
	testutil.AssertEqual(t, "/tmp/dir2", mock.InitCalls[1].Args[1])
	testutil.AssertEqual(t, "/tmp/dir3", mock.InitCalls[2].Args[1])
}

// TestExecutorInitTimeout tests Init with timeout context
func TestExecutorInitTimeout(t *testing.T) {
	mock := NewMockTerraformExecutor()
	mock.InitFunc = func(ctx context.Context, workingDir string) error {
		// Simulate long-running operation
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
			return nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := mock.Init(ctx, "/tmp/terraform")
	testutil.AssertError(t, err)
	testutil.AssertContains(t, err.Error(), "deadline exceeded")
}

// TestExecutorPlan tests the Plan method of the Executor
func TestExecutorPlan(t *testing.T) {
	tests := []struct {
		name        string
		workingDir  string
		varFile     string
		mockFunc    func(ctx context.Context, workingDir string, varFile string) (string, error)
		wantOutput  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "successful plan with no changes",
			workingDir: "/tmp/terraform",
			varFile:    "dev.tfvars",
			mockFunc: func(ctx context.Context, workingDir string, varFile string) (string, error) {
				return "No changes. Infrastructure is up-to-date.", nil
			},
			wantOutput: "No changes. Infrastructure is up-to-date.",
			wantErr:    false,
		},
		{
			name:       "successful plan with changes",
			workingDir: "/tmp/terraform",
			varFile:    "dev.tfvars",
			mockFunc: func(ctx context.Context, workingDir string, varFile string) (string, error) {
				return "Plan: 3 to add, 0 to change, 0 to destroy.", nil
			},
			wantOutput: "Plan: 3 to add, 0 to change, 0 to destroy.",
			wantErr:    false,
		},
		{
			name:       "plan without var file",
			workingDir: "/tmp/terraform",
			varFile:    "",
			mockFunc: func(ctx context.Context, workingDir string, varFile string) (string, error) {
				return "No changes. Infrastructure is up-to-date.", nil
			},
			wantOutput: "No changes. Infrastructure is up-to-date.",
			wantErr:    false,
		},
		{
			name:       "plan failure - invalid configuration",
			workingDir: "/tmp/terraform",
			varFile:    "dev.tfvars",
			mockFunc: func(ctx context.Context, workingDir string, varFile string) (string, error) {
				return "", errors.New("invalid configuration")
			},
			wantErr:     true,
			errContains: "invalid configuration",
		},
		{
			name:       "plan failure - state locked",
			workingDir: "/tmp/terraform",
			varFile:    "dev.tfvars",
			mockFunc: func(ctx context.Context, workingDir string, varFile string) (string, error) {
				return "", errors.New("state is locked")
			},
			wantErr:     true,
			errContains: "state is locked",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockTerraformExecutor()
			mock.PlanFunc = tt.mockFunc

			ctx := context.Background()
			output, err := mock.Plan(ctx, tt.workingDir, tt.varFile)

			if (err != nil) != tt.wantErr {
				t.Errorf("Plan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				testutil.AssertEqual(t, tt.wantOutput, output)
			}

			if tt.wantErr && tt.errContains != "" {
				testutil.AssertContains(t, err.Error(), tt.errContains)
			}

			// Verify call was recorded
			testutil.AssertEqual(t, 1, len(mock.PlanCalls))
			testutil.AssertEqual(t, tt.workingDir, mock.PlanCalls[0].Args[1])
			testutil.AssertEqual(t, tt.varFile, mock.PlanCalls[0].Args[2])
		})
	}
}

// TestExecutorPlanCallRecording tests that Plan calls are properly recorded
func TestExecutorPlanCallRecording(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()

	// Make multiple calls with different parameters
	_, _ = mock.Plan(ctx, "/tmp/dir1", "dev.tfvars")
	_, _ = mock.Plan(ctx, "/tmp/dir2", "staging.tfvars")
	_, _ = mock.Plan(ctx, "/tmp/dir3", "")

	// Verify all calls were recorded
	testutil.AssertEqual(t, 3, len(mock.PlanCalls))
	testutil.AssertEqual(t, "/tmp/dir1", mock.PlanCalls[0].Args[1])
	testutil.AssertEqual(t, "dev.tfvars", mock.PlanCalls[0].Args[2])
	testutil.AssertEqual(t, "/tmp/dir2", mock.PlanCalls[1].Args[1])
	testutil.AssertEqual(t, "staging.tfvars", mock.PlanCalls[1].Args[2])
	testutil.AssertEqual(t, "", mock.PlanCalls[2].Args[2])
}

// TestExecutorApply tests the Apply method of the Executor
func TestExecutorApply(t *testing.T) {
	tests := []struct {
		name        string
		workingDir  string
		varFile     string
		autoApprove bool
		mockFunc    func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error
		wantErr     bool
		errContains string
	}{
		{
			name:        "successful apply with auto-approve",
			workingDir:  "/tmp/terraform",
			varFile:     "dev.tfvars",
			autoApprove: true,
			mockFunc: func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:        "successful apply without auto-approve",
			workingDir:  "/tmp/terraform",
			varFile:     "dev.tfvars",
			autoApprove: false,
			mockFunc: func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:        "apply without var file",
			workingDir:  "/tmp/terraform",
			varFile:     "",
			autoApprove: true,
			mockFunc: func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:        "apply failure - resource already exists",
			workingDir:  "/tmp/terraform",
			varFile:     "dev.tfvars",
			autoApprove: true,
			mockFunc: func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
				return errors.New("resource already exists")
			},
			wantErr:     true,
			errContains: "resource already exists",
		},
		{
			name:        "apply failure - permission denied",
			workingDir:  "/tmp/terraform",
			varFile:     "dev.tfvars",
			autoApprove: true,
			mockFunc: func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
				return errors.New("access denied")
			},
			wantErr:     true,
			errContains: "access denied",
		},
		{
			name:        "apply failure - state locked",
			workingDir:  "/tmp/terraform",
			varFile:     "dev.tfvars",
			autoApprove: true,
			mockFunc: func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
				return errors.New("state is locked")
			},
			wantErr:     true,
			errContains: "state is locked",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockTerraformExecutor()
			mock.ApplyFunc = tt.mockFunc

			ctx := context.Background()
			err := mock.Apply(ctx, tt.workingDir, tt.varFile, tt.autoApprove)

			if (err != nil) != tt.wantErr {
				t.Errorf("Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				testutil.AssertContains(t, err.Error(), tt.errContains)
			}

			// Verify call was recorded
			testutil.AssertEqual(t, 1, len(mock.ApplyCalls))
			testutil.AssertEqual(t, tt.workingDir, mock.ApplyCalls[0].Args[1])
			testutil.AssertEqual(t, tt.varFile, mock.ApplyCalls[0].Args[2])
			testutil.AssertEqual(t, tt.autoApprove, mock.ApplyCalls[0].Args[3])
		})
	}
}

// TestExecutorApplyCallRecording tests that Apply calls are properly recorded
func TestExecutorApplyCallRecording(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()

	// Make multiple calls with different parameters
	_ = mock.Apply(ctx, "/tmp/dir1", "dev.tfvars", true)
	_ = mock.Apply(ctx, "/tmp/dir2", "staging.tfvars", false)
	_ = mock.Apply(ctx, "/tmp/dir3", "", true)

	// Verify all calls were recorded
	testutil.AssertEqual(t, 3, len(mock.ApplyCalls))
	testutil.AssertEqual(t, "/tmp/dir1", mock.ApplyCalls[0].Args[1])
	testutil.AssertEqual(t, true, mock.ApplyCalls[0].Args[3])
	testutil.AssertEqual(t, "/tmp/dir2", mock.ApplyCalls[1].Args[1])
	testutil.AssertEqual(t, false, mock.ApplyCalls[1].Args[3])
}

// TestExecutorDestroy tests the Destroy method of the Executor
func TestExecutorDestroy(t *testing.T) {
	tests := []struct {
		name        string
		workingDir  string
		varFile     string
		autoApprove bool
		mockFunc    func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error
		wantErr     bool
		errContains string
	}{
		{
			name:        "successful destroy with auto-approve",
			workingDir:  "/tmp/terraform",
			varFile:     "dev.tfvars",
			autoApprove: true,
			mockFunc: func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:        "successful destroy without auto-approve",
			workingDir:  "/tmp/terraform",
			varFile:     "dev.tfvars",
			autoApprove: false,
			mockFunc: func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:        "destroy without var file",
			workingDir:  "/tmp/terraform",
			varFile:     "",
			autoApprove: true,
			mockFunc: func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
				return nil
			},
			wantErr: false,
		},
		{
			name:        "destroy failure - resource not found",
			workingDir:  "/tmp/terraform",
			varFile:     "dev.tfvars",
			autoApprove: true,
			mockFunc: func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
				return errors.New("resource not found")
			},
			wantErr:     true,
			errContains: "resource not found",
		},
		{
			name:        "destroy failure - permission denied",
			workingDir:  "/tmp/terraform",
			varFile:     "dev.tfvars",
			autoApprove: true,
			mockFunc: func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
				return errors.New("access denied")
			},
			wantErr:     true,
			errContains: "access denied",
		},
		{
			name:        "destroy failure - state locked",
			workingDir:  "/tmp/terraform",
			varFile:     "dev.tfvars",
			autoApprove: true,
			mockFunc: func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
				return errors.New("state is locked")
			},
			wantErr:     true,
			errContains: "state is locked",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockTerraformExecutor()
			mock.DestroyFunc = tt.mockFunc

			ctx := context.Background()
			err := mock.Destroy(ctx, tt.workingDir, tt.varFile, tt.autoApprove)

			if (err != nil) != tt.wantErr {
				t.Errorf("Destroy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				testutil.AssertContains(t, err.Error(), tt.errContains)
			}

			// Verify call was recorded
			testutil.AssertEqual(t, 1, len(mock.DestroyCalls))
			testutil.AssertEqual(t, tt.workingDir, mock.DestroyCalls[0].Args[1])
			testutil.AssertEqual(t, tt.varFile, mock.DestroyCalls[0].Args[2])
			testutil.AssertEqual(t, tt.autoApprove, mock.DestroyCalls[0].Args[3])
		})
	}
}

// TestExecutorDestroyCallRecording tests that Destroy calls are properly recorded
func TestExecutorDestroyCallRecording(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()

	// Make multiple calls with different parameters
	_ = mock.Destroy(ctx, "/tmp/dir1", "dev.tfvars", true)
	_ = mock.Destroy(ctx, "/tmp/dir2", "staging.tfvars", false)
	_ = mock.Destroy(ctx, "/tmp/dir3", "", true)

	// Verify all calls were recorded
	testutil.AssertEqual(t, 3, len(mock.DestroyCalls))
	testutil.AssertEqual(t, "/tmp/dir1", mock.DestroyCalls[0].Args[1])
	testutil.AssertEqual(t, true, mock.DestroyCalls[0].Args[3])
	testutil.AssertEqual(t, "/tmp/dir2", mock.DestroyCalls[1].Args[1])
	testutil.AssertEqual(t, false, mock.DestroyCalls[1].Args[3])
}

// TestExecutorOutput tests the Output method of the Executor
func TestExecutorOutput(t *testing.T) {
	tests := []struct {
		name        string
		workingDir  string
		outputName  string
		mockFunc    func(ctx context.Context, workingDir string, outputName string) (string, error)
		wantOutput  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "get specific output value",
			workingDir: "/tmp/terraform",
			outputName: "vpc_id",
			mockFunc: func(ctx context.Context, workingDir string, outputName string) (string, error) {
				return "vpc-0123456789abcdef0", nil
			},
			wantOutput: "vpc-0123456789abcdef0",
			wantErr:    false,
		},
		{
			name:       "get all outputs as JSON",
			workingDir: "/tmp/terraform",
			outputName: "-json",
			mockFunc: func(ctx context.Context, workingDir string, outputName string) (string, error) {
				return `{"vpc_id":{"value":"vpc-123"}}`, nil
			},
			wantOutput: `{"vpc_id":{"value":"vpc-123"}}`,
			wantErr:    false,
		},
		{
			name:       "get all outputs in human-readable format",
			workingDir: "/tmp/terraform",
			outputName: "",
			mockFunc: func(ctx context.Context, workingDir string, outputName string) (string, error) {
				return "vpc_id = vpc-123\nsubnet_ids = [subnet-1, subnet-2]", nil
			},
			wantOutput: "vpc_id = vpc-123\nsubnet_ids = [subnet-1, subnet-2]",
			wantErr:    false,
		},
		{
			name:       "output not found",
			workingDir: "/tmp/terraform",
			outputName: "nonexistent",
			mockFunc: func(ctx context.Context, workingDir string, outputName string) (string, error) {
				return "", errors.New("output not found")
			},
			wantErr:     true,
			errContains: "output not found",
		},
		{
			name:       "no state file",
			workingDir: "/tmp/terraform",
			outputName: "vpc_id",
			mockFunc: func(ctx context.Context, workingDir string, outputName string) (string, error) {
				return "", errors.New("no state file")
			},
			wantErr:     true,
			errContains: "no state file",
		},
		{
			name:       "empty output value",
			workingDir: "/tmp/terraform",
			outputName: "empty_output",
			mockFunc: func(ctx context.Context, workingDir string, outputName string) (string, error) {
				return "", nil
			},
			wantOutput: "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockTerraformExecutor()
			mock.OutputFunc = tt.mockFunc

			ctx := context.Background()
			output, err := mock.Output(ctx, tt.workingDir, tt.outputName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Output() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				testutil.AssertEqual(t, tt.wantOutput, output)
			}

			if tt.wantErr && tt.errContains != "" {
				testutil.AssertContains(t, err.Error(), tt.errContains)
			}

			// Verify call was recorded
			testutil.AssertEqual(t, 1, len(mock.OutputCalls))
			testutil.AssertEqual(t, tt.workingDir, mock.OutputCalls[0].Args[1])
			testutil.AssertEqual(t, tt.outputName, mock.OutputCalls[0].Args[2])
		})
	}
}

// TestExecutorOutputCallRecording tests that Output calls are properly recorded
func TestExecutorOutputCallRecording(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()

	// Make multiple calls with different parameters
	_, _ = mock.Output(ctx, "/tmp/dir1", "vpc_id")
	_, _ = mock.Output(ctx, "/tmp/dir2", "-json")
	_, _ = mock.Output(ctx, "/tmp/dir3", "")

	// Verify all calls were recorded
	testutil.AssertEqual(t, 3, len(mock.OutputCalls))
	testutil.AssertEqual(t, "/tmp/dir1", mock.OutputCalls[0].Args[1])
	testutil.AssertEqual(t, "vpc_id", mock.OutputCalls[0].Args[2])
	testutil.AssertEqual(t, "/tmp/dir2", mock.OutputCalls[1].Args[1])
	testutil.AssertEqual(t, "-json", mock.OutputCalls[1].Args[2])
	testutil.AssertEqual(t, "", mock.OutputCalls[2].Args[2])
}

// TestExecutorErrorHandling tests error handling across all executor methods
func TestExecutorErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		mockSetup   func(*MockTerraformExecutor)
		errContains string
	}{
		{
			name:      "init - network timeout",
			operation: "init",
			mockSetup: func(m *MockTerraformExecutor) {
				m.InitFunc = func(ctx context.Context, workingDir string) error {
					return errors.New("connection timeout")
				}
			},
			errContains: "connection timeout",
		},
		{
			name:      "plan - state locked",
			operation: "plan",
			mockSetup: func(m *MockTerraformExecutor) {
				m.PlanFunc = func(ctx context.Context, workingDir string, varFile string) (string, error) {
					return "", errors.New("state is locked by another operation")
				}
			},
			errContains: "state is locked",
		},
		{
			name:      "apply - authentication failure",
			operation: "apply",
			mockSetup: func(m *MockTerraformExecutor) {
				m.ApplyFunc = func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
					return errors.New("authentication failed")
				}
			},
			errContains: "authentication failed",
		},
		{
			name:      "destroy - resource in use",
			operation: "destroy",
			mockSetup: func(m *MockTerraformExecutor) {
				m.DestroyFunc = func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
					return errors.New("resource is in use and cannot be deleted")
				}
			},
			errContains: "resource is in use",
		},
		{
			name:      "output - state file corrupted",
			operation: "output",
			mockSetup: func(m *MockTerraformExecutor) {
				m.OutputFunc = func(ctx context.Context, workingDir string, outputName string) (string, error) {
					return "", errors.New("state file is corrupted")
				}
			},
			errContains: "state file is corrupted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockTerraformExecutor()
			tt.mockSetup(mock)

			ctx := context.Background()
			var err error

			switch tt.operation {
			case "init":
				err = mock.Init(ctx, "/tmp/terraform")
			case "plan":
				_, err = mock.Plan(ctx, "/tmp/terraform", "dev.tfvars")
			case "apply":
				err = mock.Apply(ctx, "/tmp/terraform", "dev.tfvars", true)
			case "destroy":
				err = mock.Destroy(ctx, "/tmp/terraform", "dev.tfvars", true)
			case "output":
				_, err = mock.Output(ctx, "/tmp/terraform", "vpc_id")
			}

			testutil.AssertError(t, err)
			testutil.AssertContains(t, err.Error(), tt.errContains)
		})
	}
}

// TestExecutorContextCancellation tests that all operations respect context cancellation
func TestExecutorContextCancellation(t *testing.T) {
	operations := []struct {
		name string
		exec func(context.Context, *MockTerraformExecutor) error
	}{
		{
			name: "init",
			exec: func(ctx context.Context, m *MockTerraformExecutor) error {
				return m.Init(ctx, "/tmp/terraform")
			},
		},
		{
			name: "plan",
			exec: func(ctx context.Context, m *MockTerraformExecutor) error {
				_, err := m.Plan(ctx, "/tmp/terraform", "dev.tfvars")
				return err
			},
		},
		{
			name: "apply",
			exec: func(ctx context.Context, m *MockTerraformExecutor) error {
				return m.Apply(ctx, "/tmp/terraform", "dev.tfvars", true)
			},
		},
		{
			name: "destroy",
			exec: func(ctx context.Context, m *MockTerraformExecutor) error {
				return m.Destroy(ctx, "/tmp/terraform", "dev.tfvars", true)
			},
		},
		{
			name: "output",
			exec: func(ctx context.Context, m *MockTerraformExecutor) error {
				_, err := m.Output(ctx, "/tmp/terraform", "vpc_id")
				return err
			},
		},
	}

	for _, op := range operations {
		t.Run(op.name, func(t *testing.T) {
			mock := NewMockTerraformExecutor()

			// Configure mock to check for context cancellation
			checkCtx := func(ctx context.Context) error {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				return nil
			}

			mock.InitFunc = func(ctx context.Context, workingDir string) error {
				return checkCtx(ctx)
			}
			mock.PlanFunc = func(ctx context.Context, workingDir string, varFile string) (string, error) {
				return "", checkCtx(ctx)
			}
			mock.ApplyFunc = func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
				return checkCtx(ctx)
			}
			mock.DestroyFunc = func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
				return checkCtx(ctx)
			}
			mock.OutputFunc = func(ctx context.Context, workingDir string, outputName string) (string, error) {
				return "", checkCtx(ctx)
			}

			// Create cancelled context
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			// Execute operation
			err := op.exec(ctx, mock)

			// Verify context cancellation was detected
			testutil.AssertError(t, err)
			testutil.AssertContains(t, err.Error(), "context canceled")
		})
	}
}

// TestMockReset tests that Reset clears all recorded calls
func TestMockReset(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()

	// Make some calls
	_ = mock.Init(ctx, "/tmp/dir1")
	_, _ = mock.Plan(ctx, "/tmp/dir1", "dev.tfvars")
	_ = mock.Apply(ctx, "/tmp/dir1", "dev.tfvars", true)
	_, _ = mock.Output(ctx, "/tmp/dir1", "vpc_id")
	_ = mock.Destroy(ctx, "/tmp/dir1", "dev.tfvars", true)

	// Verify calls were recorded
	testutil.AssertEqual(t, 1, len(mock.InitCalls))
	testutil.AssertEqual(t, 1, len(mock.PlanCalls))
	testutil.AssertEqual(t, 1, len(mock.ApplyCalls))
	testutil.AssertEqual(t, 1, len(mock.OutputCalls))
	testutil.AssertEqual(t, 1, len(mock.DestroyCalls))

	// Reset the mock
	mock.Reset()

	// Verify all calls were cleared
	testutil.AssertEqual(t, 0, len(mock.InitCalls))
	testutil.AssertEqual(t, 0, len(mock.PlanCalls))
	testutil.AssertEqual(t, 0, len(mock.ApplyCalls))
	testutil.AssertEqual(t, 0, len(mock.OutputCalls))
	testutil.AssertEqual(t, 0, len(mock.DestroyCalls))
}

// TestMockCallCountMethods tests the call count helper methods
func TestMockCallCountMethods(t *testing.T) {
	mock := NewMockTerraformExecutor()
	ctx := context.Background()

	// Initially all counts should be zero
	testutil.AssertEqual(t, 0, mock.GetInitCallCount())
	testutil.AssertEqual(t, 0, mock.GetPlanCallCount())
	testutil.AssertEqual(t, 0, mock.GetApplyCallCount())
	testutil.AssertEqual(t, 0, mock.GetDestroyCallCount())
	testutil.AssertEqual(t, 0, mock.GetOutputCallCount())

	// Make calls
	_ = mock.Init(ctx, "/tmp/dir")
	_ = mock.Init(ctx, "/tmp/dir")
	_, _ = mock.Plan(ctx, "/tmp/dir", "dev.tfvars")
	_ = mock.Apply(ctx, "/tmp/dir", "dev.tfvars", true)
	_ = mock.Apply(ctx, "/tmp/dir", "dev.tfvars", true)
	_ = mock.Apply(ctx, "/tmp/dir", "dev.tfvars", true)

	// Verify counts
	testutil.AssertEqual(t, 2, mock.GetInitCallCount())
	testutil.AssertEqual(t, 1, mock.GetPlanCallCount())
	testutil.AssertEqual(t, 3, mock.GetApplyCallCount())
	testutil.AssertEqual(t, 0, mock.GetDestroyCallCount())
	testutil.AssertEqual(t, 0, mock.GetOutputCallCount())
}
