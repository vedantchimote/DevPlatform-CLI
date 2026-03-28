package terraform

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	clierrors "github.com/devplatform/devplatform-cli/internal/errors"
	"github.com/devplatform/devplatform-cli/internal/logger"
)

// TerraformExecutor defines the interface for executing Terraform commands
type TerraformExecutor interface {
	// Init initializes a Terraform working directory
	Init(ctx context.Context, workingDir string) error
	
	// Plan generates an execution plan
	Plan(ctx context.Context, workingDir string, varFile string) (string, error)
	
	// Apply applies the changes required to reach the desired state
	Apply(ctx context.Context, workingDir string, varFile string, autoApprove bool) error
	
	// Destroy destroys all resources managed by Terraform
	Destroy(ctx context.Context, workingDir string, varFile string, autoApprove bool) error
	
	// Output retrieves the output values from the state
	Output(ctx context.Context, workingDir string, outputName string) (string, error)
}

// Executor implements the TerraformExecutor interface
type Executor struct {
	logger *logger.Logger
}

// NewExecutor creates a new Terraform executor
func NewExecutor(log *logger.Logger) *Executor {
	return &Executor{
		logger: log,
	}
}

// Init initializes a Terraform working directory
func (e *Executor) Init(ctx context.Context, workingDir string) error {
	e.logger.Info("Initializing Terraform", logger.F("workingDir", workingDir))
	
	cmd := exec.CommandContext(ctx, "terraform", "init", "-input=false")
	cmd.Dir = workingDir
	
	output, err := e.runCommand(cmd)
	if err != nil {
		return clierrors.NewTerraformError(
			clierrors.ErrCodeTerraformInitFailed,
			"Terraform init failed",
			err,
		).WithDetails(fmt.Sprintf("Working directory: %s\nOutput: %s", workingDir, output))
	}
	
	e.logger.Success("Terraform initialized successfully")
	return nil
}

// Plan generates an execution plan
func (e *Executor) Plan(ctx context.Context, workingDir string, varFile string) (string, error) {
	e.logger.Info("Generating Terraform plan", logger.F("workingDir", workingDir))
	
	args := []string{"plan", "-input=false", "-no-color"}
	if varFile != "" {
		args = append(args, "-var-file="+varFile)
	}
	
	cmd := exec.CommandContext(ctx, "terraform", args...)
	cmd.Dir = workingDir
	
	output, err := e.runCommand(cmd)
	if err != nil {
		return output, clierrors.NewTerraformError(
			clierrors.ErrCodeTerraformPlanFailed,
			"Terraform plan failed",
			err,
		).WithDetails(fmt.Sprintf("Working directory: %s\nOutput: %s", workingDir, output))
	}
	
	e.logger.Success("Terraform plan generated successfully")
	return output, nil
}

// Apply applies the changes required to reach the desired state
func (e *Executor) Apply(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
	e.logger.Info("Applying Terraform changes", logger.F("workingDir", workingDir))
	
	args := []string{"apply", "-input=false", "-no-color"}
	if varFile != "" {
		args = append(args, "-var-file="+varFile)
	}
	if autoApprove {
		args = append(args, "-auto-approve")
	}
	
	cmd := exec.CommandContext(ctx, "terraform", args...)
	cmd.Dir = workingDir
	
	output, err := e.runCommand(cmd)
	if err != nil {
		return clierrors.NewTerraformError(
			clierrors.ErrCodeTerraformApplyFailed,
			"Terraform apply failed",
			err,
		).WithDetails(fmt.Sprintf("Working directory: %s\nOutput: %s", workingDir, output))
	}
	
	e.logger.Success("Terraform apply completed successfully")
	return nil
}

// Destroy destroys all resources managed by Terraform
func (e *Executor) Destroy(ctx context.Context, workingDir string, varFile string, autoApprove bool) error {
	e.logger.Info("Destroying Terraform resources", logger.F("workingDir", workingDir))
	
	args := []string{"destroy", "-input=false", "-no-color"}
	if varFile != "" {
		args = append(args, "-var-file="+varFile)
	}
	if autoApprove {
		args = append(args, "-auto-approve")
	}
	
	cmd := exec.CommandContext(ctx, "terraform", args...)
	cmd.Dir = workingDir
	
	output, err := e.runCommand(cmd)
	if err != nil {
		return clierrors.NewTerraformError(
			clierrors.ErrCodeTerraformDestroyFailed,
			"Terraform destroy failed",
			err,
		).WithDetails(fmt.Sprintf("Working directory: %s\nOutput: %s", workingDir, output))
	}
	
	e.logger.Success("Terraform destroy completed successfully")
	return nil
}

// Output retrieves the output values from the state
func (e *Executor) Output(ctx context.Context, workingDir string, outputName string) (string, error) {
	var args []string
	if outputName == "-json" {
		// Special case for JSON output of all values
		args = []string{"output", "-json"}
	} else if outputName != "" {
		// Get a specific output value in raw format
		args = []string{"output", "-raw", outputName}
	} else {
		// Get all outputs in human-readable format
		args = []string{"output"}
	}
	
	cmd := exec.CommandContext(ctx, "terraform", args...)
	cmd.Dir = workingDir
	
	output, err := e.runCommand(cmd)
	if err != nil {
		return "", clierrors.NewTerraformError(
			clierrors.ErrCodeTerraformOutputFailed,
			"Terraform output failed",
			err,
		).WithDetails(fmt.Sprintf("Working directory: %s, Output name: %s\nOutput: %s", workingDir, outputName, output))
	}
	
	return strings.TrimSpace(output), nil
}

// runCommand executes a command and captures stdout and stderr
func (e *Executor) runCommand(cmd *exec.Cmd) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = io.MultiWriter(&stdout, &logWriter{logger: e.logger, level: logger.InfoLevel})
	cmd.Stderr = io.MultiWriter(&stderr, &logWriter{logger: e.logger, level: logger.ErrorLevel})
	
	err := cmd.Run()
	
	// Combine stdout and stderr for error reporting
	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\n" + stderr.String()
	}
	
	return output, err
}

// logWriter is an io.Writer that writes to the logger
type logWriter struct {
	logger *logger.Logger
	level  logger.LogLevel
}

// Write implements io.Writer
func (lw *logWriter) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	if msg == "" {
		return len(p), nil
	}
	
	switch lw.level {
	case logger.ErrorLevel:
		lw.logger.Error(msg)
	default:
		lw.logger.Info(msg)
	}
	
	return len(p), nil
}
