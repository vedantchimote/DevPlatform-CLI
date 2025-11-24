# Implementation Plan: DevPlatform CLI

## Overview

This implementation plan breaks down the DevPlatform CLI into discrete coding tasks following the phased approach outlined in the design document. The CLI is a Go-based orchestration tool that wraps Terraform and Helm to provision AWS infrastructure environments. Each task builds incrementally toward a complete, working CLI tool.

## Tasks

- [x] 1. Set up Go project structure and core dependencies
  - Initialize Go module with `go mod init`
  - Add Cobra, Viper, AWS SDK v2, Azure SDK, and client-go dependencies
  - Create directory structure: `cmd/`, `internal/config/`, `internal/terraform/`, `internal/helm/`, `internal/aws/`, `internal/azure/`, `internal/provider/`, `internal/logger/`
  - Create `main.go` entry point
  - _Requirements: 25.1, 25.2, 26.1, 27.1_

- [ ] 2. Implement core CLI command structure with Cobra
  - [x] 2.1 Create root command with global flags
    - Implement `cmd/root.go` with global flags: `--verbose`, `--debug`, `--no-color`, `--config`
    - Set up command execution framework
    - Add version information constants
    - _Requirements: 11.4, 11.6, 16.4_

  - [-] 2.2 Implement version command
    - Create `cmd/version.go` with version display logic
    - Include CLI version, Git commit hash, and build date
    - Add dependency version checking for terraform, helm, kubectl, aws CLI, az CLI
    - Display minimum required versions
    - _Requirements: 19.1, 19.2, 19.3, 19.4, 19.5, 25.5, 27.1_

  - [ ]* 2.3 Write unit tests for version command
    - Test version output formatting
    - Test dependency version checking logic
    - _Requirements: 19.1, 19.2, 19.3_

- [ ] 3. Implement configuration management
  - [ ] 3.1 Create configuration data structures
    - Define `Config`, `GlobalConfig`, `EnvironmentConfig`, `TerraformConfig`, `HelmConfig`, `AzureConfig` structs in `internal/config/config.go`
    - Add YAML struct tags for parsing
    - Add CloudProvider field to GlobalConfig
    - _Requirements: 17.1, 10.2, 10.3, 10.4, 26.1, 29.1, 29.2_

  - [ ] 3.2 Implement configuration file loader
    - Create `internal/config/loader.go` with Viper integration
    - Load `.devplatform.yaml` from current directory
    - Parse YAML and populate Config structs
    - Handle file not found gracefully
    - _Requirements: 17.1, 17.4_

  - [ ] 3.3 Implement configuration validator
    - Create `internal/config/validator.go` with validation logic
    - Validate YAML schema and required fields
    - Validate environment type values (dev, staging, prod)
    - Return descriptive error messages with line numbers for YAML errors
    - _Requirements: 8.3, 17.2, 17.4_

  - [ ] 3.4 Implement CLI flag merging logic
    - Merge command-line flags with configuration file values
    - Prioritize CLI flags over file configuration
    - _Requirements: 17.3_

  - [ ]* 3.5 Write unit tests for configuration management
    - Test YAML parsing with valid and invalid files
    - Test validation rules
    - Test flag merging precedence
    - _Requirements: 17.2, 17.3, 17.4_

- [ ] 4. Implement input validation
  - [ ] 4.1 Create input validator for app name, environment, and cloud provider
    - Implement validation in `internal/config/validator.go`
    - Validate app name: lowercase alphanumeric and hyphens only, 3-32 characters
    - Validate environment type: exactly one of dev, staging, prod
    - Validate cloud provider: exactly one of aws, azure
    - Return descriptive error messages for validation failures
    - _Requirements: 1.3, 1.4, 1.5, 8.1, 8.2, 8.3, 8.4, 8.5, 26.1_

  - [ ]* 4.2 Write unit tests for input validation
    - Test app name validation with valid and invalid inputs
    - Test environment type validation
    - Test error message formatting
    - _Requirements: 8.1, 8.2, 8.3, 8.4_

- [ ] 5. Implement logging infrastructure
  - [ ] 5.1 Create logger interface and implementation
    - Define `Logger` interface in `internal/logger/logger.go`
    - Implement structured logging with Debug, Info, Warn, Error methods
    - Support log level configuration
    - Add colored console output (green for success, yellow for warnings, red for errors)
    - Respect `--no-color` flag
    - _Requirements: 16.1, 16.4, 18.1, 18.2_

  - [ ] 5.2 Implement file logging with rotation
    - Create `internal/logger/file.go` for file logging
    - Write logs to `~/.devplatform/logs/` directory
    - Implement log rotation keeping most recent 10 files
    - Use JSON format for structured logs
    - _Requirements: 18.3, 18.4, 18.5_

  - [ ]* 5.3 Write unit tests for logging
    - Test log level filtering
    - Test file rotation logic
    - Test color output formatting
    - _Requirements: 18.3, 18.5_

- [ ] 6. Checkpoint - Verify core CLI structure
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 7. Implement AWS authentication and utilities
  - [ ] 7.1 Create AWS credential validator
    - Implement `internal/aws/auth.go` with credential validation
    - Use AWS SDK v2 to validate credentials
    - Implement `GetCallerIdentity` to verify credentials work
    - Return descriptive errors for missing or expired credentials
    - _Requirements: 2.1, 2.2, 2.3, 2.4_

  - [ ] 7.2 Implement kubeconfig management
    - Create `internal/aws/kubeconfig.go` for EKS kubeconfig updates
    - Implement `UpdateKubeconfig` to configure kubectl access
    - Generate kubectl commands for namespace context switching
    - _Requirements: 13.1, 13.2, 13.3, 13.4_

  - [ ] 7.3 Implement cost calculation logic
    - Create `internal/aws/pricing.go` with cost estimation functions
    - Implement `CalculateVPCCost`, `CalculateRDSCost`, `CalculateEKSCost`
    - Calculate costs based on environment type resource sizes
    - Return monthly cost estimates in USD
    - _Requirements: 20.1, 20.2, 20.3, 20.4, 20.5_

  - [ ]* 7.4 Write unit tests for AWS utilities
    - Test cost calculation algorithms
    - Test kubeconfig command generation
    - Mock AWS SDK calls for credential validation tests
    - _Requirements: 20.2, 20.5_

- [ ] 7.5 Implement cloud provider abstraction layer
  - [ ] 7.5.1 Create CloudProvider interface
    - Define `CloudProvider` interface in `internal/provider/provider.go`
    - Define methods: ValidateCredentials, GetCallerIdentity, UpdateKubeconfig, CalculateCosts, GetTerraformBackend, GetModulePath
    - _Requirements: 26.1, 26.2_

  - [ ] 7.5.2 Create provider factory
    - Implement `NewProvider` factory function in `internal/provider/factory.go`
    - Return AWS or Azure provider based on configuration
    - _Requirements: 26.1, 26.6_

  - [ ] 7.5.3 Refactor AWS utilities to implement CloudProvider interface
    - Update AWS package to implement CloudProvider interface
    - Ensure backward compatibility
    - _Requirements: 26.1, 26.3_

- [ ] 7.6 Implement Azure authentication and utilities
  - [ ] 7.6.1 Create Azure credential validator
    - Implement `internal/azure/auth.go` with credential validation
    - Use Azure SDK to validate credentials
    - Support Azure CLI, service principal, and managed identity authentication
    - Implement `GetCallerIdentity` to verify credentials work
    - Return descriptive errors for missing or expired credentials
    - _Requirements: 27.1, 27.2, 27.3, 27.4, 27.5, 27.6_

  - [ ] 7.6.2 Implement AKS kubeconfig management
    - Create `internal/azure/kubeconfig.go` for AKS kubeconfig updates
    - Implement `UpdateKubeconfig` to configure kubectl access
    - Generate kubectl commands for namespace context switching
    - _Requirements: 28.6, 13.1, 13.2, 13.3_

  - [ ] 7.6.3 Implement Azure cost calculation logic
    - Create `internal/azure/pricing.go` with cost estimation functions
    - Implement `CalculateVNetCost`, `CalculateAzureDatabaseCost`, `CalculateAKSCost`
    - Calculate costs based on environment type resource sizes
    - Return monthly cost estimates in USD
    - _Requirements: 20.1, 20.2, 20.6, 30.2_

  - [ ]* 7.6.4 Write unit tests for Azure utilities
    - Test cost calculation algorithms
    - Test kubeconfig command generation
    - Mock Azure SDK calls for credential validation tests
    - _Requirements: 27.1, 20.2_

- [ ] 8. Implement Terraform wrapper and state management
  - [ ] 8.1 Create Terraform executor interface and implementation
    - Define `TerraformExecutor` interface in `internal/terraform/executor.go`
    - Implement `Init`, `Plan`, `Apply`, `Destroy` methods
    - Execute terraform binary with appropriate arguments
    - Capture stdout and stderr from terraform commands
    - Pass `--auto-approve` flag for non-interactive execution
    - _Requirements: 3.1, 3.2, 3.4, 3.5_

  - [ ] 8.2 Implement Terraform output parsing
    - Create `internal/terraform/output.go` for output extraction
    - Parse terraform output JSON to extract VPC ID, RDS endpoint, namespace, etc.
    - Return structured `TerraformOutputs` data
    - _Requirements: 3.3, 1.2_

  - [ ] 8.3 Implement state management with multi-backend support
    - Create `internal/terraform/state.go` for state operations
    - Configure S3 backend with bucket and DynamoDB table for AWS
    - Configure Azure Storage backend with blob lease locking for Azure
    - Generate unique state keys using app name, environment type, and cloud provider
    - Implement state existence checking
    - Handle state locking errors with descriptive messages
    - _Requirements: 9.1, 9.2, 9.3, 9.4, 9.5, 9.6, 15.1, 15.2, 15.3, 26.1_

  - [ ] 8.4 Implement Terraform error handling
    - Capture and parse terraform error output
    - Return structured errors with terraform messages
    - _Requirements: 3.5, 12.1_

  - [ ]* 8.5 Write integration tests for Terraform wrapper
    - Test terraform execution with mock state
    - Test output parsing with sample terraform output
    - Test state key generation
    - _Requirements: 3.2, 3.3, 9.3_

- [ ] 9. Create Terraform modules for infrastructure
  - [ ] 9.1 Create network module
    - Create `terraform/modules/network/` directory
    - Write `main.tf` for VPC, subnets, NAT gateways, security groups
    - Create public and private subnets across multiple AZs
    - Configure VPC flow logs
    - Add `variables.tf` for app_name, env_type, vpc_cidr
    - Add `outputs.tf` for vpc_id, subnet_ids, security_group_ids
    - _Requirements: 10.1, 10.5, 22.1, 22.2, 22.3, 22.4, 22.5_

  - [ ] 9.2 Create database module
    - Create `terraform/modules/database/` directory
    - Write `main.tf` for RDS instance, subnet group, parameter group
    - Configure single-AZ for dev/staging, multi-AZ for prod
    - Generate random password and store in Secrets Manager
    - Restrict security group access to EKS worker nodes
    - Add `variables.tf` for app_name, env_type, instance_class, allocated_storage
    - Add `outputs.tf` for rds_endpoint, rds_port, secret_arn
    - _Requirements: 10.1, 10.5, 23.1, 23.2, 23.3, 23.4, 23.5_

  - [ ] 9.3 Create EKS tenant module
    - Create `terraform/modules/eks-tenant/` directory
    - Write `main.tf` for Kubernetes namespace, resource quotas, service account, IRSA
    - Generate namespace name combining app_name and env_type
    - Configure resource quotas based on environment type
    - Add `variables.tf` for app_name, env_type, cluster_name
    - Add `outputs.tf` for namespace, service_account_name
    - _Requirements: 10.1, 10.5, 24.1, 24.2, 24.3, 24.4, 24.5, 15.4_

  - [ ] 9.4 Create environment-specific configurations
    - Create `terraform/environments/dev/`, `terraform/environments/staging/`, `terraform/environments/prod/` directories
    - Define environment-specific variable values (instance sizes, AZ counts, etc.)
    - _Requirements: 10.2, 10.3, 10.4_

  - [ ] 9.5 Implement resource tagging
    - Add tags to all AWS resources in Terraform modules
    - Include App_Name, Env_Type, Cloud_Provider, ManagedBy=devplatform-cli, Timestamp tags
    - _Requirements: 14.1, 14.2, 14.3, 14.4_

  - [ ] 9.6 Create Azure network module
    - Create `terraform/modules/azure/network/` directory
    - Write `main.tf` for VNet, subnets, NSG, NAT Gateway
    - Create public and private subnets across multiple availability zones
    - Configure Network Watcher flow logs
    - Add `variables.tf` for app_name, env_type, vnet_cidr
    - Add `outputs.tf` for vnet_id, subnet_ids, nsg_ids
    - _Requirements: 10.1, 10.5, 22.1, 22.2, 22.3, 22.4, 22.5, 28.1_

  - [ ] 9.7 Create Azure database module
    - Create `terraform/modules/azure/database/` directory
    - Write `main.tf` for Azure Database for PostgreSQL, subnet delegation
    - Configure single-zone for dev/staging, zone-redundant for prod
    - Generate random password and store in Azure Key Vault
    - Restrict NSG access to AKS worker nodes
    - Add `variables.tf` for app_name, env_type, sku, storage_mb
    - Add `outputs.tf` for database_endpoint, database_port, keyvault_secret_id
    - _Requirements: 10.1, 10.5, 23.1, 23.2, 23.3, 23.4, 23.5, 23.7, 28.2_

  - [ ] 9.8 Create Azure K8s tenant module
    - Create `terraform/modules/azure/k8s-tenant/` directory
    - Write `main.tf` for Kubernetes namespace, resource quotas, service account, Azure AD Workload Identity
    - Generate namespace name combining app_name, env_type, and cloud provider
    - Configure resource quotas based on environment type
    - Add `variables.tf` for app_name, env_type, cluster_name
    - Add `outputs.tf` for namespace, service_account_name
    - _Requirements: 10.1, 10.5, 24.1, 24.2, 24.3, 24.4, 24.7, 28.3_

  - [ ] 9.9 Create Azure environment-specific configurations
    - Create `terraform/environments/azure/dev/`, `terraform/environments/azure/staging/`, `terraform/environments/azure/prod/` directories
    - Define environment-specific variable values (SKUs, storage sizes, zone counts, etc.)
    - _Requirements: 10.2, 10.3, 10.4, 28.1, 28.2, 28.3_

  - [ ] 9.10 Implement Azure resource tagging
    - Add tags to all Azure resources in Terraform modules
    - Include App_Name, Env_Type, Cloud_Provider, ManagedBy=devplatform-cli, Timestamp tags
    - _Requirements: 14.1, 14.2, 14.3, 14.4, 28.4, 28.5_

- [ ] 10. Checkpoint - Verify Terraform integration
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 11. Implement Helm wrapper and chart management
  - [ ] 11.1 Create Helm client interface and implementation
    - Define `HelmClient` interface in `internal/helm/client.go`
    - Implement `Install`, `Upgrade`, `Uninstall`, `Status` methods
    - Execute helm binary with appropriate arguments
    - Use `helm upgrade --install` for idempotent deployments
    - Capture stdout and stderr from helm commands
    - _Requirements: 4.1, 4.4, 4.5_

  - [ ] 11.2 Implement Helm values merging
    - Create `internal/helm/values.go` for values management
    - Implement `MergeValues` to combine default and custom values
    - Implement `LoadValuesFile` to parse custom values files
    - Prioritize custom values over defaults
    - _Requirements: 21.3, 21.4_

  - [ ] 11.3 Implement pod verification
    - Add pod readiness checking in `internal/helm/client.go`
    - Use Kubernetes client-go to query pod status
    - Wait for pods to reach Running state with timeout
    - Return Kubernetes events on failure
    - _Requirements: 4.3, 4.4_

  - [ ] 11.4 Implement Helm error handling
    - Capture and parse helm error output
    - Return structured errors with Kubernetes events
    - _Requirements: 4.4_

  - [ ]* 11.5 Write integration tests for Helm wrapper
    - Test helm execution with test charts
    - Test values merging logic
    - Test pod verification with mock Kubernetes API
    - _Requirements: 4.2, 21.4_

- [ ] 12. Create base Helm chart
  - [ ] 12.1 Create Helm chart structure
    - Create `charts/devplatform-base/` directory
    - Create `Chart.yaml` with chart metadata
    - Create `values.yaml` with default values
    - Create `templates/` directory for Kubernetes manifests
    - _Requirements: 21.1, 21.5_

  - [ ] 12.2 Create Kubernetes manifest templates
    - Create `templates/deployment.yaml` for application deployment
    - Create `templates/service.yaml` for service
    - Create `templates/ingress.yaml` for ingress
    - Use Helm templating for app_name, env_type, image, resources
    - Add Kubernetes labels: app, environment, managed-by
    - _Requirements: 14.3, 21.2_

  - [ ] 12.3 Configure environment-specific values
    - Define resource requests/limits based on environment type in values.yaml
    - Configure ingress annotations and hosts
    - _Requirements: 21.2_

- [ ] 13. Implement create command
  - [ ] 13.1 Create command structure and flag parsing
    - Implement `cmd/create.go` with Cobra command definition
    - Add flags: `--app`, `--env`, `--provider`, `--dry-run`, `--values-file`, `--config`, `--timeout`
    - Set --provider default to aws for backward compatibility
    - Define `CreateOptions` struct
    - _Requirements: 11.1, 11.7, 7.1, 7.3, 26.1_

  - [ ] 13.2 Implement create command orchestration logic with multi-cloud support
    - Validate cloud provider credentials before proceeding (AWS or Azure based on --provider flag)
    - Load and validate configuration
    - Validate app name, environment type, and cloud provider inputs
    - Use provider factory to instantiate appropriate cloud provider
    - Display cost estimate in dry-run mode
    - Execute Terraform init, plan/apply with provider-specific modules
    - Parse Terraform outputs
    - Execute Helm install with merged values
    - Verify pod readiness
    - Update kubeconfig (aws eks or az aks) and display connection information
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 2.1, 7.1, 7.2, 7.3, 20.1, 26.1, 26.3, 27.1, 28.1, 28.2, 28.3_

  - [ ] 13.3 Implement dry-run mode
    - Execute terraform plan instead of apply when dry-run flag is set
    - Skip Helm installation in dry-run mode
    - Display planned changes and cost estimate
    - Clearly indicate dry-run mode in output
    - _Requirements: 7.1, 7.2, 7.3, 7.4_

  - [ ] 13.4 Implement progress indicators
    - Display progress messages during long operations
    - Show spinner or progress bar for terraform apply and helm install
    - _Requirements: 16.4_

  - [ ]* 13.5 Write integration tests for create command
    - Test full create flow with mock Terraform and Helm
    - Test dry-run mode
    - Test error handling paths
    - _Requirements: 1.1, 7.1, 7.2_

- [ ] 14. Implement error handling and rollback
  - [ ] 14.1 Create error types and categories
    - Define `CLIError` struct in `internal/errors/errors.go`
    - Define error categories: authentication, validation, terraform, helm, network, configuration
    - Assign error codes to each category (1000-1099 for auth, 1100-1199 for validation, etc.)
    - _Requirements: 1.5, 2.4, 8.4_

  - [ ] 14.2 Implement rollback logic
    - Create rollback function in create command
    - Execute helm uninstall on Helm failure
    - Execute terraform destroy on Terraform failure
    - Log rollback progress
    - Display manual cleanup instructions if rollback fails
    - _Requirements: 12.1, 12.2, 12.3, 12.4, 12.5_

  - [ ] 14.3 Implement error message formatting
    - Format errors with category, code, message, details, resolution
    - Display log file path in error messages
    - _Requirements: 8.4, 18.4_

  - [ ]* 14.4 Write unit tests for error handling
    - Test error formatting
    - Test rollback logic with mock executors
    - _Requirements: 12.1, 12.2, 12.3_

- [ ] 15. Checkpoint - Verify create command
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 16. Implement status command
  - [ ] 16.1 Create command structure and flag parsing
    - Implement `cmd/status.go` with Cobra command definition
    - Add flags: `--app`, `--env`, `--output`, `--watch`
    - Define `StatusOptions` struct
    - _Requirements: 11.2_

  - [ ] 16.2 Implement status checking logic with multi-cloud support
    - Check Terraform state existence
    - Query Terraform outputs for resource IDs
    - Query cloud provider for network and database status (VPC/VNet for AWS/Azure, RDS/Azure Database)
    - Query Kubernetes for pod and ingress status
    - Build `EnvironmentStatus` data structure
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7, 26.1, 26.3_

  - [ ] 16.3 Implement status output formatting
    - Format status as table with aligned columns
    - Support JSON and YAML output formats
    - Display component health (VPC, RDS, Namespace, Pods, Ingress)
    - Display connection information (RDS endpoint, Ingress URL)
    - _Requirements: 5.4, 16.2_

  - [ ] 16.4 Implement watch mode
    - Refresh status at specified interval when watch flag is provided
    - Clear screen and redisplay status on each refresh
    - _Requirements: 11.2_

  - [ ]* 16.5 Write integration tests for status command
    - Test status checking with mock state and APIs
    - Test output formatting
    - _Requirements: 5.1, 5.2, 5.3, 5.4_

- [ ] 17. Implement destroy command
  - [ ] 17.1 Create command structure and flag parsing
    - Implement `cmd/destroy.go` with Cobra command definition
    - Add flags: `--app`, `--env`, `--confirm`, `--force`, `--keep-state`
    - Define `DestroyOptions` struct
    - _Requirements: 11.3_

  - [ ] 17.2 Implement destroy orchestration logic
    - Prompt for confirmation if confirm flag not provided
    - Execute helm uninstall
    - Execute terraform destroy with auto-approve
    - Calculate and display cost savings
    - Handle partial deletion failures
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_

  - [ ] 17.3 Implement cost savings calculation
    - Use pricing calculator to estimate monthly savings
    - Display savings message on successful destruction
    - _Requirements: 6.4, 20.3_

  - [ ]* 17.4 Write integration tests for destroy command
    - Test destroy flow with mock Terraform and Helm
    - Test confirmation prompt
    - Test cost savings calculation
    - _Requirements: 6.1, 6.2, 6.3, 6.4_

- [ ] 18. Implement output formatting and user experience enhancements
  - [ ] 18.1 Add colored output support
    - Implement color formatting functions
    - Use green for success, yellow for warnings, red for errors
    - Respect --no-color flag
    - _Requirements: 16.1, 16.4_

  - [ ] 18.2 Implement table formatting for status output
    - Create table formatter with aligned columns
    - Format component status in readable table
    - _Requirements: 16.2_

  - [ ] 18.3 Format connection information output
    - Display RDS endpoint, Ingress URL, kubectl commands in copy-paste friendly format
    - Use code blocks or highlighted sections
    - _Requirements: 16.3_

  - [ ]* 18.4 Write unit tests for output formatting
    - Test color formatting
    - Test table alignment
    - _Requirements: 16.1, 16.2_

- [ ] 19. Implement concurrent execution safety
  - [ ] 19.1 Verify state key isolation
    - Ensure state keys include both app name and environment type
    - Test that different app/env combinations use separate state files
    - _Requirements: 15.1, 15.4_

  - [ ] 19.2 Implement state lock handling
    - Display lock holder information when state is locked
    - Provide retry instructions
    - _Requirements: 15.2, 15.3_

  - [ ]* 19.3 Write integration tests for concurrent execution
    - Test concurrent operations with different app/env combinations
    - Test state locking behavior
    - _Requirements: 15.1, 15.2_

- [ ] 20. Checkpoint - Verify all commands
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 21. Create documentation
  - [ ] 21.1 Write README with installation and usage instructions for multi-cloud
    - Document installation steps for Linux, macOS, Windows
    - Document required external tools (terraform, helm, kubectl, aws CLI, az CLI)
    - Provide usage examples for all commands with both AWS and Azure
    - Document configuration file format with multi-cloud settings
    - Add Azure setup section
    - Add multi-cloud usage patterns
    - _Requirements: 25.4, 26.1, 27.1, 29.1, 29.2, 30.1_

  - [ ] 21.2 Create command reference documentation
    - Document all commands, flags, and options
    - Provide examples for common use cases
    - Document error codes and resolutions
    - _Requirements: 11.4, 11.5_

  - [ ] 21.3 Document Terraform module structure
    - Explain module organization and variables
    - Document environment-specific configurations
    - _Requirements: 10.1, 10.2, 10.3, 10.4_

  - [ ] 21.4 Document Helm chart customization
    - Explain values.yaml structure
    - Provide examples of custom values files
    - _Requirements: 21.3, 21.4_

- [ ] 22. Set up CI/CD pipeline
  - [ ] 22.1 Create GitHub Actions workflow for testing
    - Create `.github/workflows/test.yml`
    - Run unit tests and integration tests on pull requests
    - Run linting with golangci-lint
    - _Requirements: 25.1_

  - [ ] 22.2 Create GitHub Actions workflow for releases
    - Create `.github/workflows/release.yml`
    - Use goreleaser to build binaries for Linux, macOS, Windows
    - Create GitHub releases with binaries
    - _Requirements: 25.1, 25.3_

  - [ ] 22.3 Configure goreleaser
    - Create `.goreleaser.yml` configuration
    - Configure cross-compilation for multiple platforms
    - Configure binary naming and packaging
    - _Requirements: 25.1, 25.3_

- [ ] 23. Final integration and polish
  - [ ] 23.1 Verify all commands work end-to-end
    - Test create, status, destroy workflow in test environment
    - Verify error handling and rollback
    - Verify concurrent execution
    - _Requirements: 1.1, 5.1, 6.1_

  - [ ] 23.2 Verify external tool version checking
    - Test version command with various tool versions
    - Verify minimum version enforcement
    - _Requirements: 19.4, 19.5, 25.5_

  - [ ] 23.3 Verify logging and debugging
    - Test verbose and debug modes
    - Verify log file creation and rotation
    - _Requirements: 18.1, 18.2, 18.3, 18.5_

  - [ ] 23.4 Final code cleanup and optimization
    - Remove unused code
    - Optimize performance where needed
    - Ensure consistent code style
    - _Requirements: 25.1_

- [ ] 24. Final checkpoint - Complete testing
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 25. Implement multi-cloud testing and validation
  - [ ] 25.1 Test AWS provisioning end-to-end
    - Test full create/status/destroy cycle on AWS
    - Verify all AWS resources are created correctly
    - Verify cost calculations for AWS
    - _Requirements: 1.1, 5.1, 6.1, 20.1, 20.6_

  - [ ] 25.2 Test Azure provisioning end-to-end
    - Test full create/status/destroy cycle on Azure
    - Verify all Azure resources are created correctly
    - Verify cost calculations for Azure
    - _Requirements: 28.1, 28.2, 28.3, 20.1, 20.7_

  - [ ] 25.3 Test switching between cloud providers
    - Test provisioning same app on both AWS and Azure
    - Verify state isolation between providers
    - Test status command for both providers
    - _Requirements: 26.1, 26.3, 26.4, 15.1_

  - [ ] 25.4 Test concurrent multi-cloud operations
    - Test concurrent provisioning on AWS and Azure
    - Verify no cross-cloud interference
    - Test state locking for each provider
    - _Requirements: 15.1, 15.2, 15.3, 26.1_

  - [ ] 25.5 Validate cloud provider migration documentation
    - Review resource mapping documentation
    - Test cost comparison between providers
    - Validate migration guidance
    - _Requirements: 30.1, 30.2, 30.3, 30.4_

- [ ] 26. Final checkpoint - Multi-cloud validation complete
  - Ensure all multi-cloud tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional testing tasks and can be skipped for faster MVP delivery
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation at key milestones
- The implementation follows a bottom-up approach: core infrastructure first, then orchestration, then user-facing commands
- Terraform modules and Helm charts are created alongside their respective wrapper implementations
- Error handling and rollback logic are integrated throughout rather than added at the end
- Documentation and CI/CD are included as essential deliverables, not afterthoughts
- Multi-cloud support is integrated throughout the implementation, not bolted on at the end
