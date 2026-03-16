package helm

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	clierrors "github.com/devplatform/devplatform-cli/internal/errors"
	"github.com/devplatform/devplatform-cli/internal/logger"
)

// HelmClient defines the interface for Helm operations
type HelmClient interface {
	// Install installs a Helm chart
	Install(ctx context.Context, opts InstallOptions) error
	
	// Upgrade upgrades an existing Helm release
	Upgrade(ctx context.Context, opts UpgradeOptions) error
	
	// Uninstall uninstalls a Helm release
	Uninstall(ctx context.Context, opts UninstallOptions) error
	
	// Status gets the status of a Helm release
	Status(ctx context.Context, releaseName, namespace string) (*ReleaseStatus, error)
	
	// List lists Helm releases
	List(ctx context.Context, namespace string) ([]*Release, error)
}

// Client implements the HelmClient interface
type Client struct {
	logger     logger.Logger
	helmBinary string
}

// InstallOptions contains options for helm install
type InstallOptions struct {
	ReleaseName string
	Chart       string
	Namespace   string
	ValuesFiles []string
	Values      map[string]interface{}
	CreateNamespace bool
	Wait        bool
	Timeout     time.Duration
}

// UpgradeOptions contains options for helm upgrade
type UpgradeOptions struct {
	ReleaseName string
	Chart       string
	Namespace   string
	ValuesFiles []string
	Values      map[string]interface{}
	Install     bool // --install flag
	Wait        bool
	Timeout     time.Duration
}

// UninstallOptions contains options for helm uninstall
type UninstallOptions struct {
	ReleaseName string
	Namespace   string
	Wait        bool
	Timeout     time.Duration
}

// ReleaseStatus represents the status of a Helm release
type ReleaseStatus struct {
	Name       string
	Namespace  string
	Status     string
	Revision   int
	Updated    time.Time
	Chart      string
	AppVersion string
}

// Release represents a Helm release
type Release struct {
	Name       string
	Namespace  string
	Revision   int
	Updated    time.Time
	Status     string
	Chart      string
	AppVersion string
}

// NewClient creates a new Helm client
func NewClient(log logger.Logger) *Client {
	return &Client{
		logger:     log,
		helmBinary: "helm",
	}
}

// Install installs a Helm chart
func (c *Client) Install(ctx context.Context, opts InstallOptions) error {
	c.logger.Info(fmt.Sprintf("Installing Helm chart: %s as release: %s", opts.Chart, opts.ReleaseName))

	args := []string{
		"install",
		opts.ReleaseName,
		opts.Chart,
		"--namespace", opts.Namespace,
	}

	if opts.CreateNamespace {
		args = append(args, "--create-namespace")
	}

	// Add values files
	for _, valuesFile := range opts.ValuesFiles {
		args = append(args, "--values", valuesFile)
	}

	// Add inline values
	for key, value := range opts.Values {
		args = append(args, "--set", fmt.Sprintf("%s=%v", key, value))
	}

	if opts.Wait {
		args = append(args, "--wait")
		if opts.Timeout > 0 {
			args = append(args, "--timeout", opts.Timeout.String())
		}
	}

	return c.executeHelm(ctx, args...)
}

// Upgrade upgrades an existing Helm release
func (c *Client) Upgrade(ctx context.Context, opts UpgradeOptions) error {
	c.logger.Info(fmt.Sprintf("Upgrading Helm release: %s with chart: %s", opts.ReleaseName, opts.Chart))

	args := []string{
		"upgrade",
		opts.ReleaseName,
		opts.Chart,
		"--namespace", opts.Namespace,
	}

	if opts.Install {
		args = append(args, "--install")
	}

	// Add values files
	for _, valuesFile := range opts.ValuesFiles {
		args = append(args, "--values", valuesFile)
	}

	// Add inline values
	for key, value := range opts.Values {
		args = append(args, "--set", fmt.Sprintf("%s=%v", key, value))
	}

	if opts.Wait {
		args = append(args, "--wait")
		if opts.Timeout > 0 {
			args = append(args, "--timeout", opts.Timeout.String())
		}
	}

	return c.executeHelm(ctx, args...)
}

// Uninstall uninstalls a Helm release
func (c *Client) Uninstall(ctx context.Context, opts UninstallOptions) error {
	c.logger.Info(fmt.Sprintf("Uninstalling Helm release: %s", opts.ReleaseName))

	args := []string{
		"uninstall",
		opts.ReleaseName,
		"--namespace", opts.Namespace,
	}

	if opts.Wait {
		args = append(args, "--wait")
		if opts.Timeout > 0 {
			args = append(args, "--timeout", opts.Timeout.String())
		}
	}

	return c.executeHelm(ctx, args...)
}

// Status gets the status of a Helm release
func (c *Client) Status(ctx context.Context, releaseName, namespace string) (*ReleaseStatus, error) {
	c.logger.Debug(fmt.Sprintf("Getting status for Helm release: %s", releaseName))

	args := []string{
		"status",
		releaseName,
		"--namespace", namespace,
		"--output", "json",
	}

	output, err := c.executeHelmWithOutput(ctx, args...)
	if err != nil {
		return nil, err
	}

	// Parse JSON output (simplified - in production, use proper JSON parsing)
	status := &ReleaseStatus{
		Name:      releaseName,
		Namespace: namespace,
	}

	// Extract status from output
	if strings.Contains(output, "deployed") {
		status.Status = "deployed"
	} else if strings.Contains(output, "failed") {
		status.Status = "failed"
	} else if strings.Contains(output, "pending") {
		status.Status = "pending"
	}

	return status, nil
}

// List lists Helm releases
func (c *Client) List(ctx context.Context, namespace string) ([]*Release, error) {
	c.logger.Debug(fmt.Sprintf("Listing Helm releases in namespace: %s", namespace))

	args := []string{
		"list",
		"--namespace", namespace,
		"--output", "json",
	}

	output, err := c.executeHelmWithOutput(ctx, args...)
	if err != nil {
		return nil, err
	}

	// Parse JSON output (simplified - in production, use proper JSON parsing)
	releases := []*Release{}
	
	// If output is empty or "[]", return empty list
	if strings.TrimSpace(output) == "" || strings.TrimSpace(output) == "[]" {
		return releases, nil
	}

	return releases, nil
}

// executeHelm executes a helm command and streams output to logger
func (c *Client) executeHelm(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, c.helmBinary, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	c.logger.Debug(fmt.Sprintf("Executing: %s %s", c.helmBinary, strings.Join(args, " ")))

	err := cmd.Run()

	// Log stdout
	if stdout.Len() > 0 {
		c.logger.Info(stdout.String())
	}

	// Log stderr
	if stderr.Len() > 0 {
		if err != nil {
			c.logger.Error(stderr.String())
		} else {
			c.logger.Warn(stderr.String())
		}
	}

	if err != nil {
		return clierrors.NewHelmError(
			clierrors.ErrCodeHelmInstallFailed,
			"Helm command failed",
			err,
		).WithDetails(fmt.Sprintf("Command: %s %s\nStderr: %s", c.helmBinary, strings.Join(args, " "), stderr.String()))
	}

	return nil
}

// executeHelmWithOutput executes a helm command and returns the output
func (c *Client) executeHelmWithOutput(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, c.helmBinary, args...)

	c.logger.Debug(fmt.Sprintf("Executing: %s %s", c.helmBinary, strings.Join(args, " ")))

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.logger.Error(string(output))
		return "", clierrors.NewHelmError(
			clierrors.ErrCodeHelmStatusFailed,
			"Helm command failed",
			err,
		).WithDetails(fmt.Sprintf("Command: %s %s\nOutput: %s", c.helmBinary, strings.Join(args, " "), string(output)))
	}

	return string(output), nil
}
