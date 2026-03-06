package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/devplatform/devplatform-cli/internal/config"
	"github.com/devplatform/devplatform-cli/internal/helm"
	"github.com/devplatform/devplatform-cli/internal/logger"
	"github.com/devplatform/devplatform-cli/internal/provider"
	"github.com/devplatform/devplatform-cli/internal/terraform"
	"github.com/spf13/cobra"
)

// CreateOptions holds the options for the create command
type CreateOptions struct {
	AppName      string
	Environment  string
	Provider     string
	DryRun       bool
	ValuesFile   string
	ConfigFile   string
	Timeout      time.Duration
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new environment with infrastructure and application deployment",
	Long: `Create a new environment by provisioning cloud infrastructure using Terraform
and deploying the application using Helm.

This command will:
1. Validate cloud provider credentials
2. Provision network infrastructure (VPC/VNet, subnets, NAT gateways)
3. Provision database (RDS/Azure Database for PostgreSQL)
4. Create Kubernetes namespace with resource quotas
5. Deploy application using Helm
6. Verify pod readiness
7. Configure kubectl access

Examples:
  # Create a dev environment on AWS
  devplatform create --app myapp --env dev

  # Create a prod environment on Azure
  devplatform create --app myapp --env prod --provider azure

  # Dry-run to see what would be created
  devplatform create --app myapp --env staging --dry-run

  # Use custom Helm values file
  devplatform create --app myapp --env prod --values-file custom-values.yaml`,
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Required flags
	createCmd.Flags().StringP("app", "a", "", "Application name (required)")
	createCmd.Flags().StringP("env", "e", "", "Environment type: dev, staging, or prod (required)")
	
	// Optional flags
	createCmd.Flags().StringP("provider", "p", "aws", "Cloud provider: aws or azure (default: aws)")
	createCmd.Flags().Bool("dry-run", false, "Show what would be created without actually creating resources")
	createCmd.Flags().StringP("values-file", "f", "", "Path to custom Helm values file")
	createCmd.Flags().StringP("config", "c", "", "Path to configuration file (default: .devplatform.yaml)")
	createCmd.Flags().Duration("timeout", 30*time.Minute, "Timeout for the entire create operation")

	// Mark required flags
	createCmd.MarkFlagRequired("app")
	createCmd.MarkFlagRequired("env")
}

func runCreate(cmd *cobra.Command, args []string) error {
	// Parse flags
	opts := &CreateOptions{}
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

	opts.DryRun, err = cmd.Flags().GetBool("dry-run")
	if err != nil {
		return fmt.Errorf("failed to get dry-run flag: %w", err)
	}

	opts.ValuesFile, err = cmd.Flags().GetString("values-file")
	if err != nil {
		return fmt.Errorf("failed to get values-file flag: %w", err)
	}

	opts.ConfigFile, err = cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("failed to get config flag: %w", err)
	}

	opts.Timeout, err = cmd.Flags().GetDuration("timeout")
	if err != nil {
		return fmt.Errorf("failed to get timeout flag: %w", err)
	}

	// Execute create logic
	return executeCreate(cmd, opts)
}

func executeCreate(cmd *cobra.Command, opts *CreateOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	// Initialize logger
	log := logger.GetDefault()
	log.Info(fmt.Sprintf("Starting create operation for app: %s, env: %s, provider: %s", 
		opts.AppName, opts.Environment, opts.Provider))

	if opts.DryRun {
		log.Info("DRY-RUN MODE: No resources will be created")
	}

	// Step 1: Validate inputs
	if err := validateInputs(opts); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	log.Success("Input validation passed")

	// Step 2: Load configuration
	cfg, err := loadConfiguration(opts.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	log.Success("Configuration loaded")

	// Step 3: Initialize cloud provider
	cloudProvider, err := initializeProvider(ctx, opts, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize cloud provider: %w", err)
	}
	log.Success(fmt.Sprintf("Cloud provider initialized: %s", opts.Provider))

	// Step 4: Validate credentials and display identity
	if err := validateCredentials(ctx, cloudProvider, log); err != nil {
		return fmt.Errorf("credential validation failed: %w", err)
	}

	// Step 5: Calculate and display cost estimate
	if err := displayCostEstimate(cloudProvider, opts.Environment, log); err != nil {
		log.Warn(fmt.Sprintf("Failed to calculate cost estimate: %v", err))
	}

	if opts.DryRun {
		log.Info("Dry-run complete. No resources were created.")
		return nil
	}

	// Step 6: Provision infrastructure with Terraform
	if err := provisionInfrastructure(ctx, cloudProvider, opts, cfg, log); err != nil {
		return fmt.Errorf("infrastructure provisioning failed: %w", err)
	}

	// Step 7: Deploy application with Helm
	if err := deployApplication(ctx, opts, cfg, log); err != nil {
		return fmt.Errorf("application deployment failed: %w", err)
	}

	// Step 8: Configure kubectl access
	if err := configureKubectl(cloudProvider, opts, log); err != nil {
		log.Warn(fmt.Sprintf("Failed to configure kubectl: %v", err))
	}

	// Display success message
	displaySuccessMessage(cloudProvider, opts, log)

	return nil
}

// validateInputs validates the command inputs
func validateInputs(opts *CreateOptions) error {
	// Validate app name
	if opts.AppName == "" {
		return fmt.Errorf("app name is required")
	}

	// Validate environment
	validEnvs := map[string]bool{"dev": true, "staging": true, "prod": true}
	if !validEnvs[opts.Environment] {
		return fmt.Errorf("invalid environment: %s (must be dev, staging, or prod)", opts.Environment)
	}

	// Validate provider
	validProviders := map[string]bool{"aws": true, "azure": true}
	if !validProviders[opts.Provider] {
		return fmt.Errorf("invalid provider: %s (must be aws or azure)", opts.Provider)
	}

	return nil
}

// loadConfiguration loads the configuration from file or returns defaults
func loadConfiguration(configFile string) (*config.Config, error) {
	if configFile != "" {
		return config.LoadFromPath(configFile)
	}
	return config.LoadDefault()
}

// initializeProvider creates and initializes the cloud provider
func initializeProvider(ctx context.Context, opts *CreateOptions, cfg *config.Config) (provider.CloudProvider, error) {
	providerCfg := &provider.ProviderConfig{
		Provider: opts.Provider,
	}

	// Set provider-specific configuration
	if opts.Provider == "aws" {
		providerCfg.Region = cfg.AWS.Region
		providerCfg.Profile = cfg.AWS.Profile
	} else if opts.Provider == "azure" {
		providerCfg.SubscriptionID = cfg.Azure.SubscriptionID
		providerCfg.TenantID = cfg.Azure.TenantID
		providerCfg.Location = cfg.Azure.Location
		providerCfg.ResourceGroup = fmt.Sprintf("devplatform-%s-%s", opts.AppName, opts.Environment)
	}

	return provider.NewProvider(ctx, providerCfg)
}

// validateCredentials validates cloud provider credentials
func validateCredentials(ctx context.Context, cloudProvider provider.CloudProvider, log *logger.Logger) error {
	log.Info("Validating cloud provider credentials...")
	
	if err := cloudProvider.ValidateCredentials(ctx); err != nil {
		return err
	}

	identity, err := cloudProvider.GetCallerIdentity(ctx)
	if err != nil {
		return err
	}

	log.Success(fmt.Sprintf("Authenticated as: %s (Account: %s)", identity.UserId, identity.Account))
	return nil
}

// displayCostEstimate calculates and displays the cost estimate
func displayCostEstimate(cloudProvider provider.CloudProvider, environment string, log *logger.Logger) error {
	log.Info("Calculating cost estimate...")
	
	costs, err := cloudProvider.CalculateTotalCost(environment)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("\n=== Cost Estimate for %s environment ===", environment))
	log.Info(fmt.Sprintf("Network:  $%.2f/month", costs.NetworkCost))
	log.Info(fmt.Sprintf("Database: $%.2f/month", costs.DatabaseCost))
	log.Info(fmt.Sprintf("K8s:      $%.2f/month", costs.K8sCost))
	log.Info(fmt.Sprintf("Total:    $%.2f/month", costs.TotalCost))
	log.Info("=========================================\n")

	return nil
}

// provisionInfrastructure provisions the infrastructure using Terraform
func provisionInfrastructure(ctx context.Context, cloudProvider provider.CloudProvider, opts *CreateOptions, cfg *config.Config, log *logger.Logger) error {
	log.Info("Provisioning infrastructure with Terraform...")

	// Create Terraform executor
	tfExecutor := terraform.NewExecutor(log)

	// Determine working directory
	workingDir := filepath.Join("terraform", "environments", opts.Provider)

	// Initialize Terraform
	if err := tfExecutor.Init(ctx, workingDir); err != nil {
		return err
	}

	// Generate variable file path
	varFile := fmt.Sprintf("%s-%s.tfvars", opts.AppName, opts.Environment)

	// Apply Terraform
	if err := tfExecutor.Apply(ctx, workingDir, varFile, true); err != nil {
		return err
	}

	log.Success("Infrastructure provisioned successfully")
	return nil
}

// deployApplication deploys the application using Helm
func deployApplication(ctx context.Context, opts *CreateOptions, cfg *config.Config, log *logger.Logger) error {
	log.Info("Deploying application with Helm...")

	// Create Helm client
	helmClient := helm.NewClient(*log)

	// Prepare release name and namespace
	releaseName := fmt.Sprintf("%s-%s", opts.AppName, opts.Environment)
	namespace := fmt.Sprintf("%s-%s", opts.AppName, opts.Environment)

	// Prepare values
	values := map[string]interface{}{
		"image.repository": fmt.Sprintf("%s/%s", "registry.example.com", opts.AppName),
		"image.tag":        "latest",
		"environment":      opts.Environment,
	}

	// Prepare values files
	valuesFiles := []string{}
	if opts.ValuesFile != "" {
		valuesFiles = append(valuesFiles, opts.ValuesFile)
	}

	// Install or upgrade the release
	chartPath := cfg.Helm.ChartPath
	installOpts := helm.InstallOptions{
		ReleaseName:     releaseName,
		Chart:           chartPath,
		Namespace:       namespace,
		ValuesFiles:     valuesFiles,
		Values:          values,
		CreateNamespace: true,
		Wait:            true,
		Timeout:         5 * time.Minute,
	}

	// Try install first
	err := helmClient.Install(ctx, installOpts)
	if err != nil {
		// If install fails, try upgrade with --install flag
		upgradeOpts := helm.UpgradeOptions{
			ReleaseName: releaseName,
			Chart:       chartPath,
			Namespace:   namespace,
			ValuesFiles: valuesFiles,
			Values:      values,
			Install:     true,
			Wait:        true,
			Timeout:     5 * time.Minute,
		}
		if err := helmClient.Upgrade(ctx, upgradeOpts); err != nil {
			return err
		}
	}

	// Verify pods are ready
	verifier, err := helm.NewPodVerifier(*log)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to create pod verifier: %v", err))
	} else {
		if _, err := verifier.VerifyPods(ctx, namespace, 5*time.Minute); err != nil {
			log.Warn(fmt.Sprintf("Pod verification failed: %v", err))
		}
	}

	log.Success("Application deployed successfully")
	return nil
}

// configureKubectl configures kubectl access to the cluster
func configureKubectl(cloudProvider provider.CloudProvider, opts *CreateOptions, log *logger.Logger) error {
	log.Info("Configuring kubectl access...")

	clusterName := fmt.Sprintf("%s-%s", opts.AppName, opts.Environment)
	if err := cloudProvider.UpdateKubeconfig(clusterName); err != nil {
		return err
	}

	log.Success("kubectl configured successfully")
	return nil
}

// displaySuccessMessage displays the final success message with connection commands
func displaySuccessMessage(cloudProvider provider.CloudProvider, opts *CreateOptions, log *logger.Logger) {
	clusterName := fmt.Sprintf("%s-%s", opts.AppName, opts.Environment)
	namespace := fmt.Sprintf("%s-%s", opts.AppName, opts.Environment)

	log.Success("\n=== Deployment Complete ===")
	log.Info(fmt.Sprintf("Application: %s", opts.AppName))
	log.Info(fmt.Sprintf("Environment: %s", opts.Environment))
	log.Info(fmt.Sprintf("Provider: %s", opts.Provider))
	log.Info("\nTo connect to your cluster, run:")

	commands := cloudProvider.GetConnectionCommands(clusterName, namespace)
	for _, cmd := range commands {
		log.Info(fmt.Sprintf("  %s", cmd))
	}

	log.Info("\n===========================\n")
}
