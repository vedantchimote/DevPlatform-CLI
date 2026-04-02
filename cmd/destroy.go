package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/devplatform/devplatform-cli/internal/config"
	clierrors "github.com/devplatform/devplatform-cli/internal/errors"
	"github.com/devplatform/devplatform-cli/internal/helm"
	"github.com/devplatform/devplatform-cli/internal/logger"
	"github.com/devplatform/devplatform-cli/internal/provider"
	"github.com/devplatform/devplatform-cli/internal/terraform"
	"github.com/spf13/cobra"
)

// DestroyOptions holds the options for the destroy command
type DestroyOptions struct {
	AppName     string
	Environment string
	Provider    string
	Confirm     bool
	Force       bool
	KeepState   bool
	ConfigFile  string
	Timeout     time.Duration
}

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy an environment and all its resources",
	Long: `Destroy an environment by uninstalling the Helm release and destroying
all infrastructure resources managed by Terraform.

This command will:
1. Prompt for confirmation (unless --confirm flag is provided)
2. Uninstall the Helm release
3. Destroy Terraform infrastructure (VPC/VNet, database, Kubernetes namespace)
4. Calculate and display cost savings

Examples:
  # Destroy a dev environment on AWS (with confirmation prompt)
  devplatform destroy --app myapp --env dev

  # Destroy a prod environment on Azure without confirmation
  devplatform destroy --app myapp --env prod --provider azure --confirm

  # Force destroy even if some resources fail to delete
  devplatform destroy --app myapp --env staging --force

  # Destroy but keep Terraform state file
  devplatform destroy --app myapp --env dev --keep-state`,
	RunE: runDestroy,
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	// Required flags
	destroyCmd.Flags().StringP("app", "a", "", "Application name (required)")
	destroyCmd.Flags().StringP("env", "e", "", "Environment type: dev, staging, or prod (required)")
	
	// Optional flags
	destroyCmd.Flags().StringP("provider", "p", "aws", "Cloud provider: aws or azure (default: aws)")
	destroyCmd.Flags().Bool("confirm", false, "Skip confirmation prompt and proceed with destruction")
	destroyCmd.Flags().Bool("force", false, "Force destruction even if some resources fail to delete")
	destroyCmd.Flags().Bool("keep-state", false, "Keep Terraform state file after destruction")
	destroyCmd.Flags().StringP("config", "c", "", "Path to configuration file (default: .devplatform.yaml)")
	destroyCmd.Flags().Duration("timeout", 30*time.Minute, "Timeout for the entire destroy operation")

	// Mark required flags
	destroyCmd.MarkFlagRequired("app")
	destroyCmd.MarkFlagRequired("env")
}

func runDestroy(cmd *cobra.Command, args []string) error {
	// Parse flags
	opts := &DestroyOptions{}
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

	opts.Confirm, err = cmd.Flags().GetBool("confirm")
	if err != nil {
		return fmt.Errorf("failed to get confirm flag: %w", err)
	}

	opts.Force, err = cmd.Flags().GetBool("force")
	if err != nil {
		return fmt.Errorf("failed to get force flag: %w", err)
	}

	opts.KeepState, err = cmd.Flags().GetBool("keep-state")
	if err != nil {
		return fmt.Errorf("failed to get keep-state flag: %w", err)
	}

	opts.ConfigFile, err = cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("failed to get config flag: %w", err)
	}

	opts.Timeout, err = cmd.Flags().GetDuration("timeout")
	if err != nil {
		return fmt.Errorf("failed to get timeout flag: %w", err)
	}

	// Execute destroy logic
	return executeDestroy(cmd, opts)
}

func executeDestroy(cmd *cobra.Command, opts *DestroyOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	// Initialize logger
	log := logger.GetDefault()
	log.Info(fmt.Sprintf("Starting destroy operation for app: %s, env: %s, provider: %s", 
		opts.AppName, opts.Environment, opts.Provider))

	// Track destruction state
	var (
		helmUninstalled      bool
		terraformDestroyed   bool
		cloudProvider        provider.CloudProvider
		cfg                  *config.Config
		destroyErrors        []error
	)

	// Step 1: Validate inputs
	if err := validateDestroyInputs(opts); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	log.Success("Input validation passed")

	// Step 2: Load configuration
	var err error
	cfg, err = loadConfiguration(opts.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	log.Success("Configuration loaded")

	// Step 3: Initialize cloud provider
	cloudProvider, err = initializeProvider(ctx, &CreateOptions{
		AppName:     opts.AppName,
		Environment: opts.Environment,
		Provider:    opts.Provider,
	}, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize cloud provider: %w", err)
	}
	log.Success(fmt.Sprintf("Cloud provider initialized: %s", opts.Provider))

	// Step 4: Check if environment exists
	workingDir := filepath.Join("terraform", "environments", opts.Provider)
	tfExecutor := terraform.NewExecutor(log)
	stateManager := terraform.NewStateManager(tfExecutor)

	stateExists, err := stateManager.StateExists(ctx, workingDir)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to check state existence: %v", err))
	}

	if !stateExists {
		log.Warn(fmt.Sprintf("⚠️  No Terraform state found for %s-%s on %s", opts.AppName, opts.Environment, opts.Provider))
		log.Info("The environment may not exist or has already been destroyed.")
		return nil
	}

	// Step 5: Calculate cost savings before destruction
	costSavings, err := calculateCostSavings(cloudProvider, opts.Environment, log)
	if err != nil {
		log.Warn(fmt.Sprintf("Failed to calculate cost savings: %v", err))
	}

	// Step 6: Prompt for confirmation if not provided
	if !opts.Confirm {
		if !promptForConfirmation(opts, costSavings, log) {
			log.Info("Destroy operation cancelled by user")
			return nil
		}
	}

	log.Warn("⚠️  Starting destruction process...")

	// Step 7: Uninstall Helm release
	if err := uninstallHelmRelease(ctx, opts, log); err != nil {
		log.Error(fmt.Sprintf("Failed to uninstall Helm release: %v", err))
		destroyErrors = append(destroyErrors, fmt.Errorf("helm uninstall failed: %w", err))
		
		if !opts.Force {
			return fmt.Errorf("helm uninstall failed (use --force to continue anyway): %w", err)
		}
		log.Warn("Continuing with Terraform destroy despite Helm failure (--force enabled)")
	} else {
		helmUninstalled = true
		log.Success("✓ Helm release uninstalled successfully")
	}

	// Step 8: Destroy Terraform infrastructure
	if err := destroyTerraformInfrastructure(ctx, opts, cfg, log); err != nil {
		log.Error(fmt.Sprintf("Failed to destroy Terraform infrastructure: %v", err))
		destroyErrors = append(destroyErrors, fmt.Errorf("terraform destroy failed: %w", err))
		
		if !opts.Force {
			displayPartialDestructionInstructions(opts, log, helmUninstalled, false)
			return fmt.Errorf("terraform destroy failed: %w", err)
		}
		log.Warn("Some Terraform resources may still exist (--force enabled)")
	} else {
		terraformDestroyed = true
		log.Success("✓ Terraform infrastructure destroyed successfully")
	}

	// Step 9: Display results
	if len(destroyErrors) > 0 {
		log.Warn(fmt.Sprintf("\n⚠️  Destruction completed with %d error(s)", len(destroyErrors)))
		for i, err := range destroyErrors {
			log.Error(fmt.Sprintf("  %d. %v", i+1, err))
		}
		displayPartialDestructionInstructions(opts, log, helmUninstalled, terraformDestroyed)
	} else {
		displayDestroySuccessMessage(opts, costSavings, log)
	}

	return nil
}

// validateDestroyInputs validates the command inputs
func validateDestroyInputs(opts *DestroyOptions) error {
	// Validate app name
	if opts.AppName == "" {
		return clierrors.NewValidationError(
			clierrors.ErrCodeValidationMissingRequired,
			"App name is required",
			nil,
		)
	}

	// Validate app name format
	if len(opts.AppName) < 3 || len(opts.AppName) > 32 {
		return clierrors.NewValidationError(
			clierrors.ErrCodeValidationInvalidAppName,
			fmt.Sprintf("App name '%s' is invalid: must be 3-32 characters", opts.AppName),
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

	return nil
}

// calculateCostSavings calculates the monthly cost savings from destroying the environment
func calculateCostSavings(cloudProvider provider.CloudProvider, environment string, log *logger.Logger) (float64, error) {
	log.Info("⏳ Calculating cost savings...")
	
	costs, err := cloudProvider.CalculateTotalCost(environment)
	if err != nil {
		return 0, err
	}

	return costs.TotalCost, nil
}

// promptForConfirmation prompts the user to confirm the destruction
func promptForConfirmation(opts *DestroyOptions, costSavings float64, log *logger.Logger) bool {
	log.Warn("\n⚠️  === DESTROY CONFIRMATION ===")
	log.Info(fmt.Sprintf("   Application:  %s", opts.AppName))
	log.Info(fmt.Sprintf("   Environment:  %s", opts.Environment))
	log.Info(fmt.Sprintf("   Provider:     %s", opts.Provider))
	
	if costSavings > 0 {
		log.Info(fmt.Sprintf("   Cost Savings: $%.2f/month", costSavings))
	}
	
	log.Warn("\n   This will permanently delete:")
	log.Warn("   • Helm release and all Kubernetes resources")
	log.Warn("   • Database and all data")
	log.Warn("   • Network infrastructure")
	log.Warn("   • Kubernetes namespace")
	
	if !opts.KeepState {
		log.Warn("   • Terraform state file")
	}
	
	log.Warn("\n   ⚠️  THIS ACTION CANNOT BE UNDONE!")
	log.Info("\n   Type 'yes' to confirm destruction: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Error(fmt.Sprintf("Failed to read input: %v", err))
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "yes"
}

// uninstallHelmRelease uninstalls the Helm release
func uninstallHelmRelease(ctx context.Context, opts *DestroyOptions, log *logger.Logger) error {
	log.Info("⏳ Uninstalling Helm release...")

	helmClient := helm.NewClient(*log)

	releaseName := fmt.Sprintf("%s-%s", opts.AppName, opts.Environment)
	namespace := fmt.Sprintf("%s-%s", opts.AppName, opts.Environment)

	uninstallOpts := helm.UninstallOptions{
		ReleaseName: releaseName,
		Namespace:   namespace,
		Wait:        true,
		Timeout:     5 * time.Minute,
	}

	if err := helmClient.Uninstall(ctx, uninstallOpts); err != nil {
		// Check if the error is because the release doesn't exist
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "release: not found") {
			log.Warn("Helm release not found, may have been already uninstalled")
			return nil
		}
		
		if cliErr, ok := err.(*clierrors.CLIError); ok {
			return cliErr
		}
		return clierrors.NewHelmError(
			clierrors.ErrCodeHelmUninstallFailed,
			"Failed to uninstall Helm release",
			err,
		)
	}

	return nil
}

// destroyTerraformInfrastructure destroys the Terraform-managed infrastructure
func destroyTerraformInfrastructure(ctx context.Context, opts *DestroyOptions, cfg *config.Config, log *logger.Logger) error {
	log.Info("⏳ Destroying Terraform infrastructure (this may take several minutes)...")

	tfExecutor := terraform.NewExecutor(log)

	workingDir := filepath.Join("terraform", "environments", opts.Provider)

	// Initialize Terraform
	log.Info("⏳ Initializing Terraform...")
	if err := tfExecutor.Init(ctx, workingDir); err != nil {
		if cliErr, ok := err.(*clierrors.CLIError); ok {
			return cliErr
		}
		return clierrors.NewTerraformError(
			clierrors.ErrCodeTerraformInitFailed,
			"Failed to initialize Terraform",
			err,
		)
	}

	log.Success("✓ Terraform initialized")
	log.Info("⏳ Destroying resources...")

	// Generate variable file path
	varFile := fmt.Sprintf("%s-%s.tfvars", opts.AppName, opts.Environment)

	// Destroy Terraform resources
	if err := tfExecutor.Destroy(ctx, workingDir, varFile, true); err != nil {
		if cliErr, ok := err.(*clierrors.CLIError); ok {
			return cliErr
		}
		return clierrors.NewTerraformError(
			clierrors.ErrCodeTerraformDestroyFailed,
			"Failed to destroy Terraform infrastructure",
			err,
		)
	}

	return nil
}

// displayDestroySuccessMessage displays the success message with cost savings
func displayDestroySuccessMessage(opts *DestroyOptions, costSavings float64, log *logger.Logger) {
	log.Success("\n🎉 === Destruction Complete ===")
	log.Info(fmt.Sprintf("   Application: %s", opts.AppName))
	log.Info(fmt.Sprintf("   Environment: %s", opts.Environment))
	log.Info(fmt.Sprintf("   Provider:    %s", opts.Provider))
	
	if costSavings > 0 {
		log.Success(fmt.Sprintf("\n💰 Estimated monthly savings: $%.2f", costSavings))
		log.Info(fmt.Sprintf("   Annual savings: $%.2f", costSavings*12))
	}
	
	log.Info("\n   All resources have been successfully destroyed.")
	log.Info("   ======================================\n")
}

// displayPartialDestructionInstructions displays instructions when destruction partially fails
func displayPartialDestructionInstructions(opts *DestroyOptions, log *logger.Logger, helmUninstalled bool, terraformDestroyed bool) {
	log.Error("\n⚠️  Partial destruction - Manual cleanup may be required\n")

	if !helmUninstalled {
		log.Info("To manually uninstall the Helm release:")
		releaseName := fmt.Sprintf("%s-%s", opts.AppName, opts.Environment)
		namespace := fmt.Sprintf("%s-%s", opts.AppName, opts.Environment)
		log.Info(fmt.Sprintf("  helm uninstall %s --namespace %s", releaseName, namespace))
		log.Info(fmt.Sprintf("  kubectl delete namespace %s", namespace))
		log.Info("")
	}

	if !terraformDestroyed {
		log.Info("To manually destroy Terraform resources:")
		workingDir := filepath.Join("terraform", "environments", opts.Provider)
		log.Info(fmt.Sprintf("  cd %s", workingDir))
		log.Info(fmt.Sprintf("  terraform init"))
		log.Info(fmt.Sprintf("  terraform destroy -var-file=%s-%s.tfvars", opts.AppName, opts.Environment))
		log.Info("")
	}

	log.Info("Alternatively, retry the destroy command with --force flag:")
	log.Info(fmt.Sprintf("  devplatform destroy --app %s --env %s --provider %s --force --confirm", 
		opts.AppName, opts.Environment, opts.Provider))

	log.Warn("\n⚠️  Please verify all resources have been deleted to avoid unexpected charges.")
	
	if opts.Provider == "aws" {
		log.Info("Check AWS Console: https://console.aws.amazon.com/")
	} else if opts.Provider == "azure" {
		log.Info("Check Azure Portal: https://portal.azure.com/")
	}
	
	log.Info("")
}
