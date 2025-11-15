# Requirements Document

## Introduction

The DevPlatform CLI is an Internal Developer Platform (IDP) command-line tool that enables developers to self-service provision isolated infrastructure environments on AWS or Azure without raising DevOps tickets. The CLI orchestrates Terraform for infrastructure provisioning (VPC/VNet, RDS/Azure Database, EKS/AKS namespaces) and Helm for Kubernetes application deployment, reducing environment provisioning time from 2 days to 3 minutes while enforcing cost management through easy teardown capabilities.

## Glossary

- **CLI**: The DevPlatform command-line interface application written in Go
- **Environment**: A complete isolated infrastructure stack including network, database, Kubernetes namespace, and application deployment
- **App_Name**: The unique identifier for an application being provisioned
- **Env_Type**: The environment size category (dev, staging, prod) determining resource allocation
- **Cloud_Provider**: The cloud platform (aws or azure) where resources will be provisioned
- **Terraform_Wrapper**: The internal Go module that executes Terraform commands
- **Helm_Wrapper**: The internal Go module that executes Helm commands
- **Config_Validator**: The internal Go module that validates user inputs and configuration
- **Remote_State**: The cloud storage backend storing Terraform state (S3+DynamoDB for AWS, Azure Storage for Azure)
- **K8s_Cluster**: The managed Kubernetes cluster (EKS for AWS, AKS for Azure) hosting application deployments
- **Namespace**: A Kubernetes namespace isolating application resources
- **Ingress_URL**: The publicly accessible URL for the deployed application
- **Database_Endpoint**: The database connection string for the provisioned database (RDS for AWS, Azure Database for Azure)

## Requirements

### Requirement 1: Environment Provisioning

**User Story:** As a developer, I want to provision a complete isolated environment with a single command on my chosen cloud provider, so that I can deploy and test my application without waiting for DevOps tickets.

#### Acceptance Criteria

1. WHEN a developer executes the create command with valid App_Name, Env_Type, and Cloud_Provider, THE CLI SHALL provision network resources, database, Kubernetes namespace, and deploy the application within 5 minutes
2. WHEN provisioning completes successfully, THE CLI SHALL output the Database_Endpoint, Ingress_URL, and kubeconfig instructions to the terminal
3. IF the App_Name contains invalid characters, THEN THE Config_Validator SHALL reject the input and display an error message describing valid naming conventions
4. IF the Env_Type is not one of dev, staging, or prod, THEN THE Config_Validator SHALL reject the input and display the list of valid environment types
5. IF the Cloud_Provider is not one of aws or azure, THEN THE Config_Validator SHALL reject the input and display the list of valid cloud providers
6. WHEN provisioning fails at any step, THE CLI SHALL display a descriptive error message indicating which component failed and rollback any partially created resources

### Requirement 2: Cloud Provider Authentication

**User Story:** As a developer, I want the CLI to use my existing cloud provider credentials, so that I don't need to manage separate authentication mechanisms.

#### Acceptance Criteria

1. WHEN the CLI executes any command, THE CLI SHALL validate that cloud provider credentials are configured before proceeding
2. IF cloud provider credentials are not found or expired, THEN THE CLI SHALL display an error message with instructions to configure credentials
3. THE CLI SHALL use the active cloud profile from the developer's environment (AWS profile or Azure subscription)
4. WHEN cloud API calls fail due to insufficient permissions, THE CLI SHALL display the specific permissions required
5. WHEN Cloud_Provider is aws, THE CLI SHALL validate AWS credentials using AWS SDK
6. WHEN Cloud_Provider is azure, THE CLI SHALL validate Azure credentials using Azure SDK

### Requirement 3: Terraform Infrastructure Orchestration

**User Story:** As a developer, I want the CLI to manage Terraform execution automatically, so that I don't need to understand Terraform commands or state management.

#### Acceptance Criteria

1. WHEN provisioning an environment, THE Terraform_Wrapper SHALL execute terraform init with the Remote_State backend configuration appropriate for the Cloud_Provider
2. WHEN applying infrastructure changes, THE Terraform_Wrapper SHALL pass App_Name, Env_Type, and Cloud_Provider as Terraform variables
3. WHEN terraform apply completes, THE Terraform_Wrapper SHALL extract output values including Database_Endpoint and return them to the CLI
4. THE Terraform_Wrapper SHALL execute terraform apply with auto-approve flag to enable non-interactive execution
5. WHEN Terraform execution fails, THE Terraform_Wrapper SHALL capture and return the Terraform error output to the CLI
6. WHEN Cloud_Provider is aws, THE Terraform_Wrapper SHALL use S3 backend with DynamoDB locking
7. WHEN Cloud_Provider is azure, THE Terraform_Wrapper SHALL use Azure Storage backend with blob locking

### Requirement 4: Helm Application Deployment

**User Story:** As a developer, I want my application automatically deployed to Kubernetes after infrastructure provisioning, so that I have a complete working environment immediately.

#### Acceptance Criteria

1. WHEN infrastructure provisioning completes successfully, THE Helm_Wrapper SHALL install the application chart into the provisioned Namespace
2. WHEN installing the chart, THE Helm_Wrapper SHALL pass App_Name and Env_Type as Helm values
3. WHEN Helm installation completes, THE Helm_Wrapper SHALL verify that pods are running before returning success
4. IF Helm installation fails, THEN THE Helm_Wrapper SHALL return the Kubernetes error events to the CLI
5. THE Helm_Wrapper SHALL use helm upgrade with install flag to enable idempotent deployments

### Requirement 5: Environment Status Checking

**User Story:** As a developer, I want to check the health of my provisioned environment, so that I can verify all components are running correctly.

#### Acceptance Criteria

1. WHEN a developer executes the status command with App_Name, Env_Type, and Cloud_Provider, THE CLI SHALL display the state of infrastructure and application components
2. THE CLI SHALL check Terraform Remote_State to verify infrastructure resources exist
3. THE CLI SHALL query the K8s_Cluster to verify pod health in the Namespace
4. THE CLI SHALL display Database_Endpoint, Ingress_URL, and pod status in a formatted table
5. IF the specified environment does not exist, THEN THE CLI SHALL display a message indicating no environment was found
6. WHEN Cloud_Provider is aws, THE CLI SHALL query AWS resources (VPC, RDS, EKS)
7. WHEN Cloud_Provider is azure, THE CLI SHALL query Azure resources (VNet, Azure Database, AKS)

### Requirement 6: Environment Teardown

**User Story:** As a developer, I want to destroy my environment when I'm done testing, so that I can reduce cloud costs and avoid paying for unused resources.

#### Acceptance Criteria

1. WHEN a developer executes the destroy command with App_Name, Env_Type, Cloud_Provider, and confirm flag, THE CLI SHALL delete all provisioned resources
2. THE CLI SHALL uninstall the Helm deployment before destroying infrastructure
3. WHEN Helm uninstall completes, THE Terraform_Wrapper SHALL execute terraform destroy with auto-approve flag
4. WHEN destruction completes successfully, THE CLI SHALL display an estimated cost savings message
5. IF the confirm flag is not provided, THEN THE CLI SHALL prompt for interactive confirmation before proceeding
6. WHEN destruction fails, THE CLI SHALL display which resources failed to delete and provide manual cleanup instructions

### Requirement 7: Dry Run Mode

**User Story:** As a developer, I want to preview what resources will be created without actually provisioning them, so that I can verify the configuration before committing to resource creation.

#### Acceptance Criteria

1. WHERE the dry-run flag is provided, WHEN a developer executes the create command, THE CLI SHALL display the planned infrastructure changes without applying them
2. WHEN dry-run mode is active, THE Terraform_Wrapper SHALL execute terraform plan instead of terraform apply
3. WHEN dry-run mode is active, THE CLI SHALL not execute Helm installation
4. THE CLI SHALL clearly indicate in output that dry-run mode is active and no resources were created

### Requirement 8: Input Validation

**User Story:** As a developer, I want clear error messages when I provide invalid inputs, so that I can quickly correct my command and retry.

#### Acceptance Criteria

1. THE Config_Validator SHALL verify App_Name contains only lowercase alphanumeric characters and hyphens
2. THE Config_Validator SHALL verify App_Name length is between 3 and 32 characters
3. THE Config_Validator SHALL verify Env_Type is exactly one of: dev, staging, prod
4. THE Config_Validator SHALL verify Cloud_Provider is exactly one of: aws, azure
5. WHEN validation fails, THE Config_Validator SHALL return an error message describing the validation rule that was violated
6. THE CLI SHALL validate all required flags are provided before executing any cloud operations

### Requirement 9: Remote State Management

**User Story:** As a developer, I want Terraform state stored remotely, so that I can run CLI commands from different machines without state conflicts.

#### Acceptance Criteria

1. WHEN Cloud_Provider is aws, THE Terraform_Wrapper SHALL configure S3 backend for Terraform state storage
2. WHEN Cloud_Provider is aws, THE Terraform_Wrapper SHALL configure DynamoDB table for Terraform state locking
3. WHEN Cloud_Provider is azure, THE Terraform_Wrapper SHALL configure Azure Storage backend for Terraform state storage
4. WHEN Cloud_Provider is azure, THE Terraform_Wrapper SHALL configure blob lease locking for Terraform state
5. WHEN initializing Terraform, THE Terraform_Wrapper SHALL use a state key that includes App_Name, Env_Type, and Cloud_Provider to ensure isolation
6. IF state locking fails due to another operation in progress, THEN THE Terraform_Wrapper SHALL display a message indicating the state is locked and retry instructions

### Requirement 10: Terraform Module Structure

**User Story:** As a platform engineer, I want infrastructure organized into reusable Terraform modules for each cloud provider, so that the CLI can provision consistent environments across different applications and clouds.

#### Acceptance Criteria

1. THE CLI SHALL use separate Terraform modules for network, database, and Kubernetes tenant resources for each cloud provider
2. WHEN provisioning dev environments, THE Terraform_Wrapper SHALL use small resource sizes defined in the dev configuration
3. WHEN provisioning staging environments, THE Terraform_Wrapper SHALL use medium resource sizes defined in the staging configuration
4. WHEN provisioning prod environments, THE Terraform_Wrapper SHALL use high-availability resource sizes defined in the prod configuration
5. THE Terraform modules SHALL accept App_Name, Env_Type, and Cloud_Provider as input variables
6. WHEN Cloud_Provider is aws, THE CLI SHALL use AWS-specific modules (VPC, RDS, EKS tenant)
7. WHEN Cloud_Provider is azure, THE CLI SHALL use Azure-specific modules (VNet, Azure Database, AKS tenant)


### Requirement 11: CLI Command Structure

**User Story:** As a developer, I want intuitive CLI commands following standard conventions, so that I can use the tool without extensive documentation.

#### Acceptance Criteria

1. THE CLI SHALL implement a create subcommand that accepts app, env, and provider flags
2. THE CLI SHALL implement a status subcommand that accepts app, env, and provider flags
3. THE CLI SHALL implement a destroy subcommand that accepts app, env, provider, and confirm flags
4. WHEN a developer executes the CLI without arguments, THE CLI SHALL display usage help text
5. WHEN a developer provides the help flag, THE CLI SHALL display detailed documentation for the specified subcommand
6. THE CLI SHALL follow POSIX flag conventions using double-dash for long-form flags
7. THE provider flag SHALL default to aws if not specified for backward compatibility

### Requirement 12: Error Recovery and Rollback

**User Story:** As a developer, I want automatic cleanup of partial deployments when provisioning fails, so that I don't have orphaned resources consuming costs.

#### Acceptance Criteria

1. WHEN Terraform apply fails, THE CLI SHALL execute terraform destroy to remove partially created infrastructure
2. WHEN Helm installation fails, THE CLI SHALL execute helm uninstall to remove partially deployed Kubernetes resources
3. WHEN rollback completes, THE CLI SHALL display a summary of which resources were cleaned up
4. IF rollback fails, THEN THE CLI SHALL display manual cleanup instructions with specific resource identifiers
5. THE CLI SHALL log all operations to enable troubleshooting of failed provisioning attempts

### Requirement 13: Kubernetes Configuration

**User Story:** As a developer, I want the CLI to configure my kubectl access automatically, so that I can interact with my namespace immediately after provisioning.

#### Acceptance Criteria

1. WHEN provisioning completes successfully, THE CLI SHALL display the kubectl command to configure cluster access
2. THE CLI SHALL output the specific Namespace name that was created for the application
3. THE CLI SHALL provide instructions to set the kubectl context to the provisioned Namespace
4. THE CLI SHALL verify connectivity to the K8s_Cluster before declaring provisioning successful
5. WHEN Cloud_Provider is aws, THE CLI SHALL use aws eks update-kubeconfig command
6. WHEN Cloud_Provider is azure, THE CLI SHALL use az aks get-credentials command

### Requirement 14: Resource Tagging

**User Story:** As a platform engineer, I want all provisioned resources tagged consistently, so that I can track costs and ownership across cloud accounts.

#### Acceptance Criteria

1. THE Terraform_Wrapper SHALL apply tags to all cloud resources including App_Name, Env_Type, Cloud_Provider, and ManagedBy=devplatform-cli
2. THE Terraform_Wrapper SHALL apply a Timestamp tag indicating when the resource was created
3. THE Helm_Wrapper SHALL apply Kubernetes labels to all resources including app, environment, provider, and managed-by
4. THE CLI SHALL use tags to identify resources during status checks and teardown operations
5. WHEN Cloud_Provider is aws, THE CLI SHALL use AWS resource tags
6. WHEN Cloud_Provider is azure, THE CLI SHALL use Azure resource tags

### Requirement 15: Concurrent Execution Safety

**User Story:** As a developer, I want to provision multiple environments simultaneously, so that I can create both dev and staging environments in parallel.

#### Acceptance Criteria

1. WHEN multiple CLI instances execute concurrently with different App_Name, Env_Type, or Cloud_Provider combinations, THE CLI SHALL isolate operations using separate Terraform state keys
2. WHEN multiple CLI instances execute concurrently with the same App_Name, Env_Type, and Cloud_Provider, THE Remote_State locking SHALL prevent concurrent modifications
3. IF state locking prevents execution, THEN THE CLI SHALL display which user or process holds the lock
4. THE CLI SHALL use unique Namespace names combining App_Name, Env_Type, and Cloud_Provider to prevent Kubernetes resource conflicts

### Requirement 16: Output Formatting

**User Story:** As a developer, I want clear, readable output from CLI commands, so that I can quickly understand the results and take action.

#### Acceptance Criteria

1. THE CLI SHALL use colored output to distinguish success messages (green), errors (red), and warnings (yellow)
2. WHEN displaying status information, THE CLI SHALL format output as a table with aligned columns
3. WHEN provisioning completes, THE CLI SHALL display connection information in a copy-paste friendly format
4. THE CLI SHALL display progress indicators during long-running operations like Terraform apply
5. WHERE the no-color flag is provided, THE CLI SHALL output plain text without ANSI color codes

### Requirement 17: Configuration File Support

**User Story:** As a developer, I want to store common parameters in a configuration file, so that I don't need to specify the same flags repeatedly.

#### Acceptance Criteria

1. THE CLI SHALL read configuration from a .devplatform.yaml file in the current directory if present
2. THE Config_Validator SHALL validate the configuration file schema before using values
3. WHEN both configuration file and command-line flags are provided, THE CLI SHALL prioritize command-line flags over file values
4. IF the configuration file contains invalid YAML syntax, THEN THE Config_Validator SHALL display a parsing error with line number
5. THE CLI SHALL support configuration values for default Env_Type, Cloud_Provider, and cloud-specific settings (AWS region, Azure subscription)

### Requirement 18: Logging and Debugging

**User Story:** As a developer, I want detailed logs when troubleshooting failures, so that I can understand what went wrong and report issues effectively.

#### Acceptance Criteria

1. WHERE the verbose flag is provided, THE CLI SHALL output detailed logs including all Terraform and Helm command executions
2. WHERE the debug flag is provided, THE CLI SHALL output all API calls and responses
3. THE CLI SHALL write logs to a file in the user's home directory at ~/.devplatform/logs/
4. WHEN an error occurs, THE CLI SHALL display the log file path where detailed information was written
5. THE CLI SHALL rotate log files to prevent unbounded disk usage, keeping the most recent 10 log files

### Requirement 19: Version Management

**User Story:** As a developer, I want to know which version of the CLI I'm using, so that I can verify compatibility and report issues accurately.

#### Acceptance Criteria

1. THE CLI SHALL implement a version subcommand that displays the CLI version number
2. THE CLI SHALL display the version in semantic versioning format (major.minor.patch)
3. WHEN displaying version information, THE CLI SHALL include the Git commit hash and build date
4. THE CLI SHALL check for version compatibility with Terraform and Helm binaries before executing operations
5. IF required tool versions are not met, THEN THE CLI SHALL display the minimum required versions
6. THE CLI SHALL check for cloud CLI tools (aws CLI for AWS, az CLI for Azure) based on the Cloud_Provider

### Requirement 20: Cost Estimation

**User Story:** As a developer, I want to see estimated costs before provisioning, so that I can make informed decisions about resource usage.

#### Acceptance Criteria

1. WHEN dry-run mode is active, THE CLI SHALL display estimated monthly costs for the planned infrastructure
2. THE CLI SHALL calculate costs based on Env_Type resource sizes and current cloud provider pricing
3. WHEN destroying an environment, THE CLI SHALL display the estimated monthly savings
4. THE CLI SHALL display cost estimates in USD with a disclaimer that actual costs may vary
5. THE CLI SHALL break down costs by resource type (network, database, Kubernetes) in the estimate
6. WHEN Cloud_Provider is aws, THE CLI SHALL use AWS pricing data
7. WHEN Cloud_Provider is azure, THE CLI SHALL use Azure pricing data

### Requirement 21: Helm Chart Configuration

**User Story:** As a developer, I want to customize application deployment parameters, so that I can configure my application for different environments.

#### Acceptance Criteria

1. THE Helm_Wrapper SHALL use a base Helm chart located in the charts/devplatform-base directory
2. WHEN installing the chart, THE Helm_Wrapper SHALL pass environment-specific values based on Env_Type
3. THE CLI SHALL support a values-file flag to provide custom Helm values
4. WHERE custom values are provided, THE Helm_Wrapper SHALL merge custom values with default values, prioritizing custom values
5. THE Helm_Wrapper SHALL validate that the Helm chart exists before attempting installation

### Requirement 22: Network Configuration

**User Story:** As a platform engineer, I want networks configured with proper security and isolation, so that applications are protected and compliant with security policies.

#### Acceptance Criteria

1. WHEN provisioning a network, THE Terraform network module SHALL create public and private subnets across multiple availability zones
2. THE Terraform network module SHALL configure NAT gateways for private subnet internet access
3. THE Terraform network module SHALL create security groups restricting database access to Kubernetes worker nodes only
4. THE Terraform network module SHALL enable network flow logs for network traffic auditing
5. WHEN provisioning prod environments, THE Terraform network module SHALL create resources across at least 3 availability zones
6. WHEN Cloud_Provider is aws, THE CLI SHALL provision VPC with appropriate configuration
7. WHEN Cloud_Provider is azure, THE CLI SHALL provision VNet with appropriate configuration

### Requirement 23: Database Configuration

**User Story:** As a developer, I want databases configured appropriately for my environment type, so that I have adequate performance without over-provisioning.

#### Acceptance Criteria

1. WHEN provisioning dev environments, THE Terraform database module SHALL create single-zone database instances with minimal instance sizes
2. WHEN provisioning staging environments, THE Terraform database module SHALL create single-zone database instances with medium instance sizes
3. WHEN provisioning prod environments, THE Terraform database module SHALL create multi-zone database instances with automated backups enabled
4. THE Terraform database module SHALL generate random database passwords and store them in the cloud provider's secret management service
5. WHEN provisioning completes, THE CLI SHALL display instructions to retrieve the database password from the secret management service
6. WHEN Cloud_Provider is aws, THE CLI SHALL provision RDS with passwords in AWS Secrets Manager
7. WHEN Cloud_Provider is azure, THE CLI SHALL provision Azure Database with passwords in Azure Key Vault

### Requirement 24: Kubernetes Namespace Configuration

**User Story:** As a platform engineer, I want Kubernetes namespaces configured with resource quotas, so that applications cannot consume excessive cluster resources.

#### Acceptance Criteria

1. WHEN provisioning a Namespace, THE Terraform Kubernetes tenant module SHALL create the Namespace with a name combining App_Name, Env_Type, and Cloud_Provider
2. THE Terraform Kubernetes tenant module SHALL apply resource quotas limiting CPU and memory based on Env_Type
3. THE Terraform Kubernetes tenant module SHALL create a Kubernetes service account for the application
4. THE Terraform Kubernetes tenant module SHALL configure cloud-specific identity integration (IRSA for AWS, Workload Identity for Azure)
5. WHEN provisioning dev environments, THE Terraform Kubernetes tenant module SHALL set resource quotas to 2 CPU cores and 4GB memory
6. WHEN Cloud_Provider is aws, THE CLI SHALL configure IAM Roles for Service Accounts (IRSA)
7. WHEN Cloud_Provider is azure, THE CLI SHALL configure Azure AD Workload Identity

### Requirement 25: Binary Distribution

**User Story:** As a developer, I want to download and install the CLI easily, so that I can start using it without complex setup procedures.

#### Acceptance Criteria

1. THE CLI SHALL be distributed as a single statically-linked binary for Linux, macOS, and Windows
2. THE CLI SHALL not require Go runtime or other dependencies to be installed
3. THE CLI binary SHALL be available for download from GitHub releases
4. THE CLI SHALL provide installation instructions in the README for adding the binary to system PATH
5. THE CLI SHALL verify on startup that required external tools are available in PATH (terraform, helm, kubectl, and cloud CLI tools based on provider)
6. WHEN Cloud_Provider is aws, THE CLI SHALL verify aws CLI is available
7. WHEN Cloud_Provider is azure, THE CLI SHALL verify az CLI is available

### Requirement 26: Cloud Provider Abstraction

**User Story:** As a platform engineer, I want a consistent interface for provisioning across cloud providers, so that developers have the same experience regardless of which cloud they use.

#### Acceptance Criteria

1. THE CLI SHALL implement a cloud provider interface that abstracts AWS and Azure differences
2. THE CLI SHALL provide consistent resource naming across cloud providers
3. THE CLI SHALL map equivalent resources between clouds (VPC↔VNet, RDS↔Azure Database, EKS↔AKS)
4. THE CLI SHALL handle cloud-specific authentication mechanisms transparently
5. THE CLI SHALL provide consistent error messages regardless of cloud provider
6. THE CLI SHALL support adding new cloud providers without modifying core CLI logic

### Requirement 27: Azure-Specific Authentication

**User Story:** As a developer using Azure, I want the CLI to use my existing Azure credentials, so that I don't need to manage separate authentication mechanisms.

#### Acceptance Criteria

1. WHEN Cloud_Provider is azure, THE CLI SHALL validate Azure credentials using Azure SDK
2. THE CLI SHALL support Azure CLI authentication (az login)
3. THE CLI SHALL support Azure service principal authentication
4. THE CLI SHALL support Azure managed identity authentication
5. THE CLI SHALL use the active Azure subscription from the developer's environment
6. WHEN Azure API calls fail due to insufficient permissions, THE CLI SHALL display the specific Azure RBAC roles required

### Requirement 28: Azure Resource Provisioning

**User Story:** As a developer, I want to provision Azure resources with the same ease as AWS resources, so that I can use my preferred cloud provider.

#### Acceptance Criteria

1. WHEN Cloud_Provider is azure, THE CLI SHALL provision Azure VNet with public and private subnets
2. WHEN Cloud_Provider is azure, THE CLI SHALL provision Azure Database for PostgreSQL
3. WHEN Cloud_Provider is azure, THE CLI SHALL provision AKS namespace with resource quotas
4. WHEN Cloud_Provider is azure, THE CLI SHALL configure Azure Network Security Groups
5. WHEN Cloud_Provider is azure, THE CLI SHALL enable Azure Network Watcher flow logs
6. WHEN Cloud_Provider is azure, THE CLI SHALL configure Azure AD Workload Identity for pod authentication

### Requirement 29: Multi-Cloud Configuration

**User Story:** As a developer working with multiple clouds, I want to configure settings for both AWS and Azure in a single configuration file, so that I can easily switch between providers.

#### Acceptance Criteria

1. THE CLI SHALL support separate configuration sections for AWS and Azure in .devplatform.yaml
2. THE CLI SHALL allow specifying default Cloud_Provider in configuration file
3. THE CLI SHALL support cloud-specific settings (AWS region, Azure location, Azure subscription ID)
4. THE CLI SHALL validate cloud-specific configuration only when that provider is selected
5. THE CLI SHALL provide clear error messages when cloud-specific configuration is missing

### Requirement 30: Cloud Provider Migration Support

**User Story:** As a platform engineer, I want to understand the differences between cloud providers, so that I can help developers migrate between clouds if needed.

#### Acceptance Criteria

1. THE CLI SHALL document resource mapping between AWS and Azure
2. THE CLI SHALL provide cost comparison between AWS and Azure for equivalent environments
3. THE CLI SHALL warn users about cloud-specific features that don't have equivalents
4. THE CLI SHALL provide migration guidance in documentation
5. THE CLI SHALL support exporting environment configuration for recreation on different cloud

