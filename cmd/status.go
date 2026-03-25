package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/devplatform/devplatform-cli/internal/config"
	clierrors "github.com/devplatform/devplatform-cli/internal/errors"
	"github.com/devplatform/devplatform-cli/internal/helm"
	"github.com/devplatform/devplatform-cli/internal/logger"
	"github.com/devplatform/devplatform-cli/internal/terraform"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// StatusOptions holds the options for the status command
type StatusOptions struct {
	AppName      string
	Environment  string
	Provider     string
	OutputFormat string
	Watch        int
	ConfigFile   string
}

// EnvironmentStatus represents the complete status of an environment
type EnvironmentStatus struct {
	AppName     string          `json:"app_name" yaml:"app_name"`
	Environment string          `json:"environment" yaml:"environment"`
	Provider    string          `json:"provider" yaml:"provider"`
	Status      string          `json:"status" yaml:"status"` // healthy, degraded, failed, not_found
	Components  ComponentStatus `json:"components" yaml:"components"`
	LastUpdated time.Time       `json:"last_updated" yaml:"last_updated"`
}

// ComponentStatus represents the status of all components
type ComponentStatus struct {
	Network   ResourceStatus `json:"network" yaml:"network"`
	Database  ResourceStatus `json:"database" yaml:"database"`
	Namespace ResourceStatus `json:"namespace" yaml:"namespace"`
	Pods      PodStatus      `json:"pods" yaml:"pods"`
	Ingress   IngressStatus  `json:"ingress" yaml:"ingress"`
}

// ResourceStatus represents the status of a single resource
type ResourceStatus struct {
	Status  string            `json:"status" yaml:"status"` // ok, error, not_found
	ID      string            `json:"id" yaml:"id"`
	Details map[string]string `json:"details,omitempty" yaml:"details,omitempty"`
}

// PodStatus represents the status of pods
type PodStatus struct {
	Status string     `json:"status" yaml:"status"`
	Ready  int        `json:"ready" yaml:"ready"`
	Total  int        `json:"total" yaml:"total"`
	Pods   []PodInfo  `json:"pods,omitempty" yaml:"pods,omitempty"`
}

// PodInfo contains information about a single pod
type PodInfo struct {
	Name     string `json:"name" yaml:"name"`
	Status   string `json:"status" yaml:"status"`
	Ready    bool   `json:"ready" yaml:"ready"`
	Restarts int32  `json:"restarts" yaml:"restarts"`
}

// IngressStatus represents the status of ingress
type IngressStatus struct {
	Status string `json:"status" yaml:"status"`
	URL    string `json:"url,omitempty" yaml:"url,omitempty"`
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of an environment",
	Long: `Check the health and status of a provisioned environment.

This command will:
1. Check if Terraform state exists
2. Query cloud provider resources (VPC/VNet, RDS/Azure Database)
3. Query Kubernetes for pod and ingress status
4. Display status information in the requested format

Examples:
  # Check status of a dev environment on AWS
  devplatform status --app myapp --env dev

  # Check status with JSON output
  devplatform status --app myapp --env prod --output json

  # Watch status with auto-refresh every 5 seconds
  devplatform status --app myapp --env staging --watch 5`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Required flags
	statusCmd.Flags().StringP("app", "a", "", "Application name (required)")
	statusCmd.Flags().StringP("env", "e", "", "Environment type: dev, staging, or prod (required)")
	
	// Optional flags
	statusCmd.Flags().StringP("provider", "p", "aws", "Cloud provider: aws or azure (default: aws)")
	statusCmd.Flags().StringP("output", "o", "table", "Output format: table, json, or yaml (default: table)")
	statusCmd.Flags().IntP("watch", "w", 0, "Watch mode: refresh interval in seconds (0 = disabled)")
	statusCmd.Flags().StringP("config", "c", "", "Path to configuration file (default: .devplatform.yaml)")

	// Mark required flags
	statusCmd.MarkFlagRequired("app")
	statusCmd.MarkFlagRequired("env")
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Parse flags
	opts := &StatusOptions{}
	var err error

	opts.AppName, err = cmd.Flags().GetString("app")
	if err != nil {
		return fmt.Errorf("failed to get app flag: %w", err)
	}

	opts.Environment, err = cmd.Flags().GetString("env")
	if err != nil {
		return fmt.Errorf("failed to get env flag: %w", err)
	}

	opts.Provider, err = cmd.Flags().GetString("provider")
	if err != nil {
		return fmt.Errorf("failed to get provider flag: %w", err)
	}

	opts.OutputFormat, err = cmd.Flags().GetString("output")
	if err != nil {
		return fmt.Errorf("failed to get output flag: %w", err)
	}

	opts.Watch, err = cmd.Flags().GetInt("watch")
	if err != nil {
		return fmt.Errorf("failed to get watch flag: %w", err)
	}

	opts.ConfigFile, err = cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("failed to get config flag: %w", err)
	}

	// Execute status logic
	return executeStatus(cmd, opts)
}

func executeStatus(cmd *cobra.Command, opts *StatusOptions) error {
	ctx := context.Background()

	// Initialize logger
	log := logger.GetDefault()

	// Validate inputs
	if err := validateStatusInputs(opts); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Load configuration
	cfg, err := loadConfiguration(opts.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Watch mode
	if opts.Watch > 0 {
		return watchStatus(ctx, opts, cfg, log)
	}

	// Single status check
	status, err := checkEnvironmentStatus(ctx, opts, cfg, log)
	if err != nil {
		return err
	}

	// Display status
	return displayStatus(status, opts.OutputFormat, log)
}

// validateStatusInputs validates the command inputs
func validateStatusInputs(opts *StatusOptions) error {
	// Validate app name
	if opts.AppName == "" {
		return clierrors.NewValidationError(
			clierrors.ErrCodeValidationMissingRequired,
			"App name is required",
			nil,
		)
	}

	// Validate environment
	validEnvs := map[string]bool{"dev": true, "staging": true, "prod": true}
	if !validEnvs[opts.Environment] {
		return clierrors.NewValidationError(
			clierrors.ErrCodeValidationInvalidEnvironment,
			fmt.Sprintf("Invalid environment: %s", opts.Environment),
			nil,
		).WithDetails(fmt.Sprintf("Provided: %s, Valid options: dev, staging, prod", opts.Environment))
	}

	// Validate provider
	validProviders := map[string]bool{"aws": true, "azure": true}
	if !validProviders[opts.Provider] {
		return clierrors.NewValidationError(
			clierrors.ErrCodeValidationInvalidProvider,
			fmt.Sprintf("Invalid provider: %s", opts.Provider),
			nil,
		).WithDetails(fmt.Sprintf("Provided: %s, Valid options: aws, azure", opts.Provider))
	}

	// Validate output format
	validFormats := map[string]bool{"table": true, "json": true, "yaml": true}
	if !validFormats[opts.OutputFormat] {
		return clierrors.NewValidationError(
			clierrors.ErrCodeValidationInvalidConfig,
			fmt.Sprintf("Invalid output format: %s", opts.OutputFormat),
			nil,
		).WithDetails(fmt.Sprintf("Provided: %s, Valid options: table, json, yaml", opts.OutputFormat))
	}

	return nil
}

// checkEnvironmentStatus checks the status of all environment components
func checkEnvironmentStatus(ctx context.Context, opts *StatusOptions, cfg *config.Config, log *logger.Logger) (*EnvironmentStatus, error) {
	status := &EnvironmentStatus{
		AppName:     opts.AppName,
		Environment: opts.Environment,
		Provider:    opts.Provider,
		Status:      "unknown",
		LastUpdated: time.Now(),
	}

	// Check Terraform state existence
	workingDir := filepath.Join("terraform", "environments", opts.Provider)
	tfExecutor := terraform.NewExecutor(log)
	stateManager := terraform.NewStateManager(tfExecutor)

	stateExists, err := stateManager.StateExists(ctx, workingDir)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to check state existence: %v", err))
	}

	if !stateExists {
		status.Status = "not_found"
		status.Components = ComponentStatus{
			Network:   ResourceStatus{Status: "not_found"},
			Database:  ResourceStatus{Status: "not_found"},
			Namespace: ResourceStatus{Status: "not_found"},
			Pods:      PodStatus{Status: "not_found"},
			Ingress:   IngressStatus{Status: "not_found"},
		}
		return status, nil
	}

	// Get Terraform outputs
	outputParser := terraform.NewOutputParser(tfExecutor)
	outputs, err := outputParser.ParseOutputs(ctx, workingDir)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to parse Terraform outputs: %v", err))
		status.Status = "error"
		return status, nil
	}

	// Check network status
	status.Components.Network = checkNetworkStatus(outputs, opts.Provider)

	// Check database status
	status.Components.Database = checkDatabaseStatus(outputs, opts.Provider)

	// Check namespace status
	status.Components.Namespace = checkNamespaceStatus(outputs)

	// Check pod status
	status.Components.Pods = checkPodStatus(ctx, outputs.Namespace, log)

	// Check ingress status
	status.Components.Ingress = checkIngressStatus(ctx, outputs.Namespace, log)

	// Determine overall status
	status.Status = determineOverallStatus(&status.Components)

	return status, nil
}

// checkNetworkStatus checks the status of network resources
func checkNetworkStatus(outputs *terraform.TerraformOutputs, provider string) ResourceStatus {
	networkID := outputs.GetNetworkID()
	if networkID == "" {
		return ResourceStatus{Status: "not_found"}
	}

	details := make(map[string]string)
	if provider == "aws" {
		details["vpc_id"] = outputs.VPCID
		if len(outputs.SubnetIDs) > 0 {
			details["subnets"] = fmt.Sprintf("%d subnets", len(outputs.SubnetIDs))
		}
	} else if provider == "azure" {
		details["vnet_id"] = outputs.VNetID
	}

	return ResourceStatus{
		Status:  "ok",
		ID:      networkID,
		Details: details,
	}
}

// checkDatabaseStatus checks the status of database resources
func checkDatabaseStatus(outputs *terraform.TerraformOutputs, provider string) ResourceStatus {
	dbEndpoint := outputs.GetDatabaseEndpoint()
	if dbEndpoint == "" {
		return ResourceStatus{Status: "not_found"}
	}

	details := make(map[string]string)
	details["endpoint"] = dbEndpoint
	
	dbPort := outputs.GetDatabasePort()
	if dbPort != "" {
		details["port"] = dbPort
	}

	return ResourceStatus{
		Status:  "ok",
		ID:      dbEndpoint,
		Details: details,
	}
}

// checkNamespaceStatus checks the status of Kubernetes namespace
func checkNamespaceStatus(outputs *terraform.TerraformOutputs) ResourceStatus {
	if outputs.Namespace == "" {
		return ResourceStatus{Status: "not_found"}
	}

	return ResourceStatus{
		Status: "ok",
		ID:     outputs.Namespace,
	}
}

// checkPodStatus checks the status of pods in the namespace
func checkPodStatus(ctx context.Context, namespace string, log *logger.Logger) PodStatus {
	if namespace == "" {
		return PodStatus{Status: "not_found"}
	}

	verifier, err := helm.NewPodVerifier(*log)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to create pod verifier: %v", err))
		return PodStatus{Status: "error"}
	}

	helmPodStatus, err := verifier.GetPodStatus(ctx, namespace)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to get pod status: %v", err))
		return PodStatus{Status: "error"}
	}

	// Convert helm.PodStatus to our PodStatus
	pods := make([]PodInfo, 0, len(helmPodStatus.Pods))
	for _, pod := range helmPodStatus.Pods {
		pods = append(pods, PodInfo{
			Name:     pod.Name,
			Status:   pod.Status,
			Ready:    pod.Ready,
			Restarts: pod.Restarts,
		})
	}

	status := "ok"
	if helmPodStatus.FailedPods > 0 {
		status = "error"
	} else if helmPodStatus.PendingPods > 0 {
		status = "pending"
	} else if helmPodStatus.ReadyPods < helmPodStatus.TotalPods {
		status = "degraded"
	}

	return PodStatus{
		Status: status,
		Ready:  helmPodStatus.ReadyPods,
		Total:  helmPodStatus.TotalPods,
		Pods:   pods,
	}
}

// checkIngressStatus checks the status of ingress resources
func checkIngressStatus(ctx context.Context, namespace string, log *logger.Logger) IngressStatus {
	if namespace == "" {
		return IngressStatus{Status: "not_found"}
	}

	// For now, we'll return a basic status
	// In a full implementation, we would query the Kubernetes API for ingress resources
	return IngressStatus{
		Status: "ok",
		URL:    "", // Would be populated from actual ingress query
	}
}

// determineOverallStatus determines the overall environment status
func determineOverallStatus(components *ComponentStatus) string {
	// If any component is not found, environment doesn't exist
	if components.Network.Status == "not_found" ||
		components.Database.Status == "not_found" ||
		components.Namespace.Status == "not_found" {
		return "not_found"
	}

	// If any component has error, environment is failed
	if components.Network.Status == "error" ||
		components.Database.Status == "error" ||
		components.Pods.Status == "error" {
		return "failed"
	}

	// If pods are not all ready, environment is degraded
	if components.Pods.Status == "degraded" || components.Pods.Status == "pending" {
		return "degraded"
	}

	// All components are ok
	return "healthy"
}

// displayStatus displays the status in the requested format
func displayStatus(status *EnvironmentStatus, format string, log *logger.Logger) error {
	switch format {
	case "json":
		return displayStatusJSON(status)
	case "yaml":
		return displayStatusYAML(status)
	default:
		return displayStatusTable(status, log)
	}
}

// displayStatusJSON displays status in JSON format
func displayStatusJSON(status *EnvironmentStatus) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(status)
}

// displayStatusYAML displays status in YAML format
func displayStatusYAML(status *EnvironmentStatus) error {
	encoder := yaml.NewEncoder(os.Stdout)
	encoder.SetIndent(2)
	return encoder.Encode(status)
}

// displayStatusTable displays status in table format
func displayStatusTable(status *EnvironmentStatus, log *logger.Logger) error {
	// Header
	log.Info(fmt.Sprintf("\n📊 === Environment Status: %s ===", status.AppName))
	log.Info(fmt.Sprintf("   Environment: %s", status.Environment))
	log.Info(fmt.Sprintf("   Provider:    %s", status.Provider))
	
	// Overall status with color
	statusIcon := getStatusIcon(status.Status)
	log.Info(fmt.Sprintf("   Status:      %s %s", statusIcon, status.Status))
	log.Info(fmt.Sprintf("   Updated:     %s", status.LastUpdated.Format("2006-01-02 15:04:05")))
	log.Info("")

	// If environment not found, display message and return
	if status.Status == "not_found" {
		log.Warn("   ⚠️  Environment not found. Has it been created?")
		log.Info("   Run 'devplatform create' to provision this environment.")
		log.Info("")
		return nil
	}

	// Component status table
	log.Info("   Component Status:")
	log.Info("   ┌─────────────────┬──────────┬────────────────────────────────┐")
	log.Info("   │ Component       │ Status   │ Details                        │")
	log.Info("   ├─────────────────┼──────────┼────────────────────────────────┤")

	// Network
	displayComponentRow("Network", status.Components.Network, log)

	// Database
	displayComponentRow("Database", status.Components.Database, log)

	// Namespace
	displayComponentRow("Namespace", status.Components.Namespace, log)

	// Pods
	podDetails := fmt.Sprintf("%d/%d ready", status.Components.Pods.Ready, status.Components.Pods.Total)
	displayComponentRowWithDetails("Pods", status.Components.Pods.Status, podDetails, log)

	// Ingress
	ingressDetails := ""
	if status.Components.Ingress.URL != "" {
		ingressDetails = status.Components.Ingress.URL
	}
	displayComponentRowWithDetails("Ingress", status.Components.Ingress.Status, ingressDetails, log)

	log.Info("   └─────────────────┴──────────┴────────────────────────────────┘")
	log.Info("")

	// Pod details if available
	if len(status.Components.Pods.Pods) > 0 {
		log.Info("   Pod Details:")
		for _, pod := range status.Components.Pods.Pods {
			readyIcon := "✓"
			if !pod.Ready {
				readyIcon = "✗"
			}
			log.Info(fmt.Sprintf("   %s %s (%s, restarts: %d)", readyIcon, pod.Name, pod.Status, pod.Restarts))
		}
		log.Info("")
	}

	// Connection information
	if status.Components.Database.Status == "ok" {
		log.Info("   Connection Information:")
		if endpoint, ok := status.Components.Database.Details["endpoint"]; ok {
			log.Info(fmt.Sprintf("   Database: %s", endpoint))
		}
		if status.Components.Ingress.URL != "" {
			log.Info(fmt.Sprintf("   Ingress:  %s", status.Components.Ingress.URL))
		}
		log.Info("")
	}

	log.Info("   ======================================\n")
	return nil
}

// displayComponentRow displays a component row in the table
func displayComponentRow(name string, resource ResourceStatus, log *logger.Logger) {
	statusIcon := getStatusIcon(resource.Status)
	details := resource.ID
	if len(resource.Details) > 0 {
		// Show first detail
		for k, v := range resource.Details {
			details = fmt.Sprintf("%s: %s", k, v)
			break
		}
	}
	
	// Truncate details if too long
	if len(details) > 30 {
		details = details[:27] + "..."
	}
	
	log.Info(fmt.Sprintf("   │ %-15s │ %s %-6s │ %-30s │", name, statusIcon, resource.Status, details))
}

// displayComponentRowWithDetails displays a component row with custom details
func displayComponentRowWithDetails(name string, status string, details string, log *logger.Logger) {
	statusIcon := getStatusIcon(status)
	
	// Truncate details if too long
	if len(details) > 30 {
		details = details[:27] + "..."
	}
	
	log.Info(fmt.Sprintf("   │ %-15s │ %s %-6s │ %-30s │", name, statusIcon, status, details))
}

// getStatusIcon returns an icon for the status
func getStatusIcon(status string) string {
	switch status {
	case "ok", "healthy":
		return "✓"
	case "error", "failed":
		return "✗"
	case "degraded", "pending":
		return "⚠"
	case "not_found":
		return "○"
	default:
		return "?"
	}
}

// watchStatus continuously monitors and displays environment status
func watchStatus(ctx context.Context, opts *StatusOptions, cfg *config.Config, log *logger.Logger) error {
	ticker := time.NewTicker(time.Duration(opts.Watch) * time.Second)
	defer ticker.Stop()

	// Display initial status
	status, err := checkEnvironmentStatus(ctx, opts, cfg, log)
	if err != nil {
		return err
	}

	if err := displayStatus(status, opts.OutputFormat, log); err != nil {
		return err
	}

	// Watch loop
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// Clear screen (works on most terminals)
			if opts.OutputFormat == "table" {
				clearScreen()
			}

			// Check status again
			status, err := checkEnvironmentStatus(ctx, opts, cfg, log)
			if err != nil {
				log.Error(fmt.Sprintf("Failed to check status: %v", err))
				continue
			}

			// Display status
			if err := displayStatus(status, opts.OutputFormat, log); err != nil {
				log.Error(fmt.Sprintf("Failed to display status: %v", err))
			}
		}
	}
}

// clearScreen clears the terminal screen
func clearScreen() {
	// ANSI escape code to clear screen and move cursor to top-left
	fmt.Print("\033[2J\033[H")
}
