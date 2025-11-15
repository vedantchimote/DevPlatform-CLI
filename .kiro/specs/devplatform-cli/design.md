# Design Document: DevPlatform CLI

## Overview

The DevPlatform CLI is an Internal Developer Platform (IDP) command-line tool written in Go that enables developers to self-service provision complete, isolated infrastructure environments on AWS or Azure. The CLI orchestrates Terraform for infrastructure management (VPC/VNet, RDS/Azure Database, EKS/AKS namespaces) and Helm for Kubernetes application deployment, reducing environment provisioning time from 2 days (via DevOps tickets) to approximately 3 minutes of automated execution.

The tool follows a declarative approach where developers specify what they want (app name, environment type, and cloud provider), and the CLI handles all the complexity of infrastructure provisioning, state management, and application deployment. It enforces cost management through easy teardown capabilities and provides clear visibility into resource status and costs across both AWS and Azure.

### Key Design Goals

1. **Developer Self-Service**: Enable developers to provision environments without DevOps intervention
2. **Multi-Cloud Support**: Provide consistent experience across AWS and Azure
3. **Infrastructure Abstraction**: Hide Terraform, Helm, and cloud-specific complexity behind simple CLI commands
4. **State Management**: Ensure safe concurrent operations through remote state locking
5. **Cost Awareness**: Provide cost estimates and savings calculations to promote FinOps practices
6. **Error Recovery**: Automatically rollback partial deployments to prevent orphaned resources
7. **Observability**: Provide clear status information and detailed logging for troubleshooting

## Architecture

### High-Level Architecture

The DevPlatform CLI follows a layered architecture with clear separation of concerns and cloud provider abstraction:

```
┌─────────────────────────────────────────────────────────────┐
│                     CLI Interface Layer                      │
│  (Cobra commands: create, status, destroy, version)         │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                   Business Logic Layer                       │
│  (Config validation, orchestration, error handling)         │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                Cloud Provider Abstraction Layer              │
│  (Provider interface, AWS impl, Azure impl)                 │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                  Infrastructure Adapters                     │
│  (Terraform wrapper, Helm wrapper, Cloud utilities)         │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                    External Dependencies                     │
│  (Terraform, Helm, AWS SDK, Azure SDK, Kubernetes API)      │
└─────────────────────────────────────────────────────────────┘
```

### Component Architecture

The system is organized into the following major components:

1. **CLI Command Layer** (`cmd/` package)
   - Root command handler with global flags
   - Subcommand implementations (create, status, destroy, version)
   - Flag parsing and validation
   - User interaction and output formatting

2. **Configuration Management** (`internal/config/` package)
   - YAML configuration file parsing
   - Schema validation
   - CLI flag merging with configuration
   - Environment-specific settings
   - Cloud provider configuration

3. **Cloud Provider Abstraction** (`internal/provider/` package)
   - CloudProvider interface defining common operations
   - AWS provider implementation
   - Azure provider implementation
   - Provider factory for instantiation

4. **Terraform Orchestration** (`internal/terraform/` package)
   - Terraform binary execution wrapper
   - State management and locking (S3+DynamoDB or Azure Storage)
   - Output parsing and extraction
   - Error handling and rollback
   - Multi-backend support

5. **Helm Orchestration** (`internal/helm/` package)
   - Helm binary execution wrapper
   - Chart loading and validation
   - Values merging and templating
   - Release status checking

6. **AWS Integration** (`internal/aws/` package)
   - AWS credential validation
   - EKS kubeconfig management
   - AWS cost calculation and estimation
   - AWS resource tagging utilities

7. **Azure Integration** (`internal/azure/` package)
   - Azure credential validation
   - AKS kubeconfig management
   - Azure cost calculation and estimation
   - Azure resource tagging utilities

8. **Logging and Observability** (`internal/logger/` package)
   - Structured logging
   - Log file rotation
   - Verbosity level management
   - Error tracking

### Data Flow

The typical data flow for environment provisioning follows this sequence:

1. **Input Phase**: User provides app name, environment type, and cloud provider via CLI flags
2. **Validation Phase**: Config validator checks input format and cloud provider credentials
3. **Configuration Phase**: Load and merge configuration from file and CLI flags
4. **Provider Selection**: Instantiate appropriate cloud provider (AWS or Azure) via factory
5. **Infrastructure Phase**: Terraform wrapper provisions network, database, and Kubernetes namespace using provider-specific modules
6. **Deployment Phase**: Helm wrapper deploys application to Kubernetes
7. **Verification Phase**: Check pod health and ingress availability
8. **Output Phase**: Display connection information and success message

### Concurrency Model

The CLI supports concurrent execution through Terraform's remote state locking mechanism:

- Each environment (app + env + provider combination) has a unique state key
- Cloud-specific locking mechanisms (DynamoDB for AWS, blob lease for Azure) provide distributed locking
- Multiple developers can provision different environments simultaneously
- Concurrent operations on the same environment are serialized by the lock
- Lock timeout and retry logic prevents indefinite blocking
- Cross-cloud operations are fully isolated (AWS and Azure states are separate)

### Error Handling Strategy

The CLI implements a comprehensive error handling strategy:

1. **Early Validation**: Catch input errors before any cloud operations
2. **Graceful Degradation**: Continue with partial operations when safe
3. **Automatic Rollback**: Clean up partial deployments on failure
4. **Clear Error Messages**: Provide actionable error messages with resolution steps
5. **Detailed Logging**: Write full error context to log files for debugging
6. **Exit Codes**: Use distinct exit codes for different error categories

## Components and Interfaces

### Cloud Provider Interface

The CloudProvider interface defines the contract that all cloud providers must implement:

```go
// CloudProvider defines the interface that all cloud providers must implement
type CloudProvider interface {
    // Authentication
    ValidateCredentials() error
    GetIdentity() (*Identity, error)
    
    // Kubernetes
    UpdateKubeconfig(clusterName string) error
    GetKubernetesClient() (kubernetes.Interface, error)
    
    // Cost Calculation
    CalculateNetworkCost(config NetworkConfig) float64
    CalculateDatabaseCost(config DatabaseConfig) float64
    CalculateKubernetesCost(config KubernetesConfig) float64
    GetTotalCost(env string) float64
    
    // Terraform Backend
    GetBackendConfig() BackendConfig
    
    // Resource Tagging
    GetDefaultTags() map[string]string
    
    // Provider-specific info
    GetProviderName() string
    GetRegion() string
}

// Identity represents cloud provider identity information
type Identity struct {
    AccountID   string
    UserID      string
    ARN         string  // AWS-specific
    TenantID    string  // Azure-specific
}

// NetworkConfig represents network configuration
type NetworkConfig struct {
    CIDR              string
    AvailabilityZones int
    NATGateways       int
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
    InstanceClass string
    Storage       int
    MultiAZ       bool
}

// KubernetesConfig represents Kubernetes configuration
type KubernetesConfig struct {
    NodeCount int
    NodeSize  string
}

// BackendConfig represents Terraform backend configuration
type BackendConfig struct {
    Type   string  // s3 or azurerm
    Config map[string]string
}
```

### Provider Factory

```go
type ProviderFactory struct{}

func (f *ProviderFactory) CreateProvider(providerType string, config map[string]string) (CloudProvider, error) {
    switch providerType {
    case "aws":
        region := config["region"]
        profile := config["profile"]
        return NewAWSProvider(region, profile)
        
    case "azure":
        subscriptionID := config["subscription_id"]
        location := config["location"]
        tenantID := config["tenant_id"]
        return NewAzureProvider(subscriptionID, location, tenantID)
        
    default:
        return nil, fmt.Errorf("unsupported provider: %s", providerType)
    }
}
```

### CLI Command Interface

The CLI exposes four primary commands through the Cobra framework:

#### create Command

```go
type CreateOptions struct {
    AppName      string
    Environment  string
    Provider     string  // aws or azure
    DryRun       bool
    ValuesFile   string
    ConfigFile   string
    Verbose      bool
    Debug        bool
    NoColor      bool
    Timeout      int
}

func RunCreate(opts CreateOptions) error
```

**Responsibilities:**
- Validate input parameters including cloud provider
- Load and merge configuration
- Instantiate appropriate cloud provider via factory
- Orchestrate Terraform infrastructure provisioning using provider-specific modules
- Orchestrate Helm application deployment
- Display success information with connection details

#### status Command

```go
type StatusOptions struct {
    AppName     string
    Environment string
    OutputFormat string  // table, json, yaml
    Watch       int      // refresh interval in seconds
    Verbose     bool
    NoColor     bool
}

func RunStatus(opts StatusOptions) error
```

**Responsibilities:**
- Check Terraform state existence
- Query AWS resources for health status
- Query Kubernetes for pod and ingress status
- Format and display status information

#### destroy Command

```go
type DestroyOptions struct {
    AppName     string
    Environment string
    Confirm     bool
    Force       bool
    KeepState   bool
    Verbose     bool
    NoColor     bool
}

func RunDestroy(opts DestroyOptions) error
```

**Responsibilities:**
- Prompt for confirmation if not provided
- Uninstall Helm release
- Destroy Terraform infrastructure
- Calculate and display cost savings

#### version Command

```go
type VersionOptions struct {
    Short      bool
    CheckDeps  bool
}

func RunVersion(opts VersionOptions) error
```

**Responsibilities:**
- Display CLI version information
- Check and display dependency versions
- Validate minimum required versions

### Configuration Management Interface

```go
type Config struct {
    Global       GlobalConfig
    Environments map[string]EnvironmentConfig
    Terraform    TerraformConfig
    Helm         HelmConfig
}

type ConfigManager interface {
    Load(path string) (*Config, error)
    Validate() error
    GetEnvironment(env string) (*EnvironmentConfig, error)
    MergeFlags(flags map[string]interface{}) error
}
```

**Responsibilities:**
- Parse YAML configuration files
- Validate configuration schema
- Merge CLI flags with file configuration (flags take precedence)
- Provide environment-specific settings

### Terraform Wrapper Interface

```go
type TerraformExecutor interface {
    Init(backend BackendConfig) error
    Plan(vars map[string]string) (string, error)
    Apply(vars map[string]string) (map[string]string, error)
    Destroy(vars map[string]string) error
    GetOutputs() (map[string]string, error)
}

type StateManager interface {
    GetState(key string) (*State, error)
    LockState(key string) error
    UnlockState(key string) error
    StateExists(key string) (bool, error)
}
```

**Responsibilities:**
- Execute Terraform binary with appropriate arguments
- Manage backend configuration for remote state
- Parse Terraform outputs
- Handle state locking and unlocking
- Implement rollback on failure

### Helm Wrapper Interface

```go
type HelmClient interface {
    Install(name string, chart string, values map[string]interface{}) error
    Upgrade(name string, chart string, values map[string]interface{}) error
    Uninstall(name string) error
    Status(name string) (*ReleaseStatus, error)
}

type ChartManager interface {
    LoadChart(path string) (*Chart, error)
    ValidateChart(chart *Chart) error
}

type ValuesMerger interface {
    MergeValues(base, override map[string]interface{}) map[string]interface{}
    LoadValuesFile(path string) (map[string]interface{}, error)
}
```

**Responsibilities:**
- Execute Helm binary for release management
- Load and validate Helm charts
- Merge default and custom values
- Verify pod readiness after deployment
- Handle uninstall during rollback

### AWS Utilities Interface

```go
type AWSAuth interface {
    ValidateCredentials() error
    GetCallerIdentity() (*Identity, error)
}

type KubeconfigManager interface {
    UpdateKubeconfig(clusterName, region string) error
    SetNamespace(namespace string) error
}

type PricingCalculator interface {
    CalculateVPCCost(config VPCConfig) float64
    CalculateRDSCost(config RDSConfig) float64
    CalculateEKSCost(config EKSConfig) float64
    GetTotalCost(env string) float64
}
```

**Responsibilities:**
- Validate AWS credentials and permissions
- Update kubeconfig for EKS cluster access
- Calculate infrastructure costs based on resource configuration
- Provide cost estimates and savings calculations

### Logger Interface

```go
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    SetLevel(level LogLevel)
}

type LogManager interface {
    RotateLogs() error
    GetLogPath() string
}
```

**Responsibilities:**
- Provide structured logging with context
- Support multiple log levels (debug, info, warn, error)
- Write logs to file with rotation
- Format console output with colors
- Track operations for audit trail

## Data Models

### Configuration Data Models

```go
type Config struct {
    Global       GlobalConfig
    AWS          AWSConfig
    Azure        AzureConfig
    Environments map[string]EnvironmentConfig
    Terraform    TerraformConfig
    Helm         HelmConfig
}

type GlobalConfig struct {
    CloudProvider string `yaml:"cloud_provider"`  // aws or azure
    Timeout       int    `yaml:"timeout"`
    LogLevel      string `yaml:"log_level"`
}

type AWSConfig struct {
    Region  string `yaml:"region"`
    Profile string `yaml:"aws_profile"`
}

type AzureConfig struct {
    SubscriptionID string `yaml:"subscription_id"`
    Location       string `yaml:"location"`
    TenantID       string `yaml:"tenant_id"`
}

type EnvironmentConfig struct {
    // Network configuration (VPC for AWS, VNet for Azure)
    NetworkCIDR          string `yaml:"network_cidr"`
    
    // Database configuration (RDS for AWS, Azure Database for Azure)
    DBInstanceClass      string `yaml:"db_instance_class"`
    DBAllocatedStorage   int    `yaml:"db_allocated_storage"`
    DBMultiAZ            bool   `yaml:"db_multi_az"`
    
    // Kubernetes configuration (EKS for AWS, AKS for Azure)
    K8sNodeCount         int    `yaml:"k8s_node_count"`
}

type TerraformConfig struct {
    Backend     BackendConfig `yaml:"backend"`
    ModulesPath string        `yaml:"modules_path"`
}

type BackendConfig struct {
    // AWS S3 backend
    Type          string `yaml:"type"`  // s3 or azurerm
    Bucket        string `yaml:"bucket"`
    DynamoDBTable string `yaml:"dynamodb_table"`
    Region        string `yaml:"region"`
    
    // Azure Storage backend
    StorageAccountName string `yaml:"storage_account_name"`
    ContainerName      string `yaml:"container_name"`
    ResourceGroupName  string `yaml:"resource_group_name"`
}

type HelmConfig struct {
    ChartPath     string                 `yaml:"chart_path"`
    DefaultValues map[string]interface{} `yaml:"default_values"`
}
```

### Terraform Data Models

```go
type TerraformState struct {
    Version   int                    `json:"version"`
    Resources []TerraformResource    `json:"resources"`
    Outputs   map[string]interface{} `json:"outputs"`
}

type TerraformResource struct {
    Type      string                 `json:"type"`
    Name      string                 `json:"name"`
    Provider  string                 `json:"provider"`
    Instances []ResourceInstance     `json:"instances"`
}

type ResourceInstance struct {
    Attributes map[string]interface{} `json:"attributes"`
}

type TerraformOutputs struct {
    VPCId        string `json:"vpc_id"`
    VPCCIDR      string `json:"vpc_cidr"`
    RDSEndpoint  string `json:"rds_endpoint"`
    RDSPort      int    `json:"rds_port"`
    Namespace    string `json:"namespace"`
    SecretARN    string `json:"db_secret_arn"`
}
```

### Helm Data Models

```go
type ReleaseStatus struct {
    Name        string
    Namespace   string
    Status      string  // deployed, failed, pending-install, etc.
    Version     int
    LastUpdated time.Time
}

type HelmValues struct {
    Image       ImageConfig                `yaml:"image"`
    Resources   ResourceRequirements       `yaml:"resources"`
    Environment map[string]string          `yaml:"environment"`
    Ingress     IngressConfig              `yaml:"ingress"`
}

type ImageConfig struct {
    Repository string `yaml:"repository"`
    Tag        string `yaml:"tag"`
    PullPolicy string `yaml:"pullPolicy"`
}

type ResourceRequirements struct {
    Requests ResourceList `yaml:"requests"`
    Limits   ResourceList `yaml:"limits"`
}

type ResourceList struct {
    CPU    string `yaml:"cpu"`
    Memory string `yaml:"memory"`
}

type IngressConfig struct {
    Enabled     bool              `yaml:"enabled"`
    Annotations map[string]string `yaml:"annotations"`
    Hosts       []string          `yaml:"hosts"`
}
```

### Status Data Models

```go
type EnvironmentStatus struct {
    AppName      string
    Environment  string
    Status       string  // healthy, degraded, failed, not_found
    Components   ComponentStatus
    LastUpdated  time.Time
}

type ComponentStatus struct {
    VPC       ResourceStatus
    RDS       ResourceStatus
    Namespace ResourceStatus
    Pods      PodStatus
    Ingress   IngressStatus
}

type ResourceStatus struct {
    Status  string  // ok, error, not_found
    ID      string
    Details map[string]string
}

type PodStatus struct {
    Status string
    Ready  int
    Total  int
    Pods   []PodInfo
}

type PodInfo struct {
    Name   string
    Status string
    Ready  bool
}

type IngressStatus struct {
    Status string
    URL    string
}
```

### Error Data Models

```go
type CLIError struct {
    Code       int
    Category   ErrorCategory
    Message    string
    Details    string
    Resolution string
    Cause      error
}

type ErrorCategory string

const (
    ErrorCategoryAuth       ErrorCategory = "authentication"
    ErrorCategoryValidation ErrorCategory = "validation"
    ErrorCategoryTerraform  ErrorCategory = "terraform"
    ErrorCategoryHelm       ErrorCategory = "helm"
    ErrorCategoryNetwork    ErrorCategory = "network"
    ErrorCategoryConfig     ErrorCategory = "configuration"
)
```

## Technology Stack

### Core Technologies

1. **Go 1.21+**: Primary programming language
   - Strong standard library for CLI development
   - Excellent concurrency support
   - Cross-platform compilation
   - Static binary distribution

2. **Cobra Framework**: CLI command structure
   - Industry-standard CLI framework (used by kubectl, helm, etc.)
   - Built-in help generation
   - Flag parsing and validation
   - Subcommand organization

3. **Viper**: Configuration management
   - YAML configuration file support
   - Environment variable binding
   - Configuration merging
   - Type-safe configuration access

### Infrastructure Technologies

1. **Terraform 1.5+**: Infrastructure as Code
   - Declarative infrastructure definition
   - State management with locking
   - Modular architecture
   - AWS and Azure provider support

2. **Helm 3.x**: Kubernetes package management
   - Templated Kubernetes manifests
   - Release management
   - Values override mechanism
   - No Tiller requirement (Helm 3)

3. **AWS SDK for Go v2**: AWS API interaction
   - Credential management
   - Service-specific clients (STS, EKS, EC2, RDS)
   - Automatic retries
   - Context support

4. **Azure SDK for Go**: Azure API interaction
   - Azure Identity for authentication
   - ARM clients for resource management
   - Subscription and resource group management
   - Context support

5. **Kubernetes client-go**: Kubernetes API interaction
   - Pod status checking
   - Namespace management
   - Resource querying
   - Watch capabilities

### Development and Build Tools

1. **golangci-lint**: Code quality and linting
2. **go test**: Unit and integration testing
3. **GitHub Actions**: CI/CD pipeline
4. **goreleaser**: Multi-platform binary builds

### External Dependencies

The CLI requires these external tools to be installed:

1. **terraform**: Version 1.5 or higher
2. **helm**: Version 3.0 or higher
3. **kubectl**: Version 1.27 or higher
4. **aws**: AWS CLI version 2.x

## Implementation Approach

### Phase 1: Core CLI Structure

1. Set up Go module and project structure
2. Implement Cobra command framework
3. Create basic flag parsing and validation
4. Implement configuration file loading
5. Add logging infrastructure

### Phase 2: Terraform Integration

1. Implement Terraform executor wrapper
2. Add state management functionality
3. Create Terraform module structure (network, database, EKS tenant)
4. Implement output parsing
5. Add rollback logic

### Phase 3: Helm Integration

1. Implement Helm client wrapper
2. Create base Helm chart
3. Add values merging logic
4. Implement pod verification
5. Add uninstall functionality

### Phase 4: AWS Integration

1. Implement credential validation
2. Add kubeconfig management
3. Create cost calculation logic
4. Implement resource tagging

### Phase 5: Error Handling and Observability

1. Implement comprehensive error handling
2. Add automatic rollback on failures
3. Create detailed logging
4. Add status checking functionality

### Phase 6: Polish and Distribution

1. Add colored output
2. Implement progress indicators
3. Create comprehensive documentation
4. Set up CI/CD pipeline
5. Create release binaries

### Repository Structure

```
devplatform-cli/
├── cmd/
│   ├── root.go           # Root command and global flags
│   ├── create.go         # Create command implementation
│   ├── status.go         # Status command implementation
│   ├── destroy.go        # Destroy command implementation
│   └── version.go        # Version command implementation
├── internal/
│   ├── config/
│   │   ├── config.go     # Configuration structures
│   │   ├── loader.go     # Configuration loading
│   │   └── validator.go  # Configuration validation
│   ├── provider/
│   │   ├── interface.go  # CloudProvider interface
│   │   ├── factory.go    # Provider factory
│   │   ├── aws.go        # AWS provider implementation
│   │   └── azure.go      # Azure provider implementation
│   ├── terraform/
│   │   ├── executor.go   # Terraform execution
│   │   ├── state.go      # State management
│   │   └── output.go     # Output parsing
│   ├── helm/
│   │   ├── client.go     # Helm client
│   │   ├── chart.go      # Chart management
│   │   └── values.go     # Values merging
│   ├── aws/
│   │   ├── auth.go       # AWS authentication
│   │   ├── kubeconfig.go # Kubeconfig management
│   │   └── pricing.go    # Cost calculation
│   ├── azure/
│   │   ├── auth.go       # Azure authentication
│   │   ├── kubeconfig.go # Kubeconfig management
│   │   └── pricing.go    # Cost calculation
│   └── logger/
│       ├── logger.go     # Logging interface
│       └── file.go       # File logging
├── terraform/
│   ├── modules/
│   │   ├── aws/
│   │   │   ├── network/      # VPC module
│   │   │   ├── database/     # RDS module
│   │   │   └── eks-tenant/   # EKS namespace module
│   │   └── azure/
│   │       ├── network/      # VNet module
│   │       ├── database/     # Azure Database module
│   │       └── aks-tenant/   # AKS namespace module
│   └── environments/
│       ├── dev/          # Dev environment config
│       ├── staging/      # Staging environment config
│       └── prod/         # Prod environment config
├── charts/
│   └── devplatform-base/ # Base Helm chart
├── docs/                 # Documentation
├── .github/
│   └── workflows/        # CI/CD workflows
├── go.mod
├── go.sum
└── README.md
```

### Testing Strategy

The DevPlatform CLI is primarily an orchestration tool that integrates external systems (Terraform, Helm, AWS, Kubernetes). The testing strategy focuses on:

1. **Unit Tests**: Test individual components in isolation
   - Configuration parsing and validation
   - Input validation logic
   - Output formatting
   - Error message generation
   - Cost calculation algorithms

2. **Integration Tests**: Test interactions with external tools
   - Terraform execution with mock state
   - Helm execution with test charts
   - AWS SDK interactions with localstack
   - Kubernetes API interactions with kind/minikube

3. **End-to-End Tests**: Test complete workflows
   - Full create/status/destroy cycle in test environment
   - Error handling and rollback scenarios
   - Concurrent execution scenarios

4. **Mock-Based Tests**: Test external tool wrappers
   - Mock Terraform binary execution
   - Mock Helm binary execution
   - Mock AWS API calls
   - Mock Kubernetes API calls

Property-based testing is not applicable for this CLI tool because:
- The CLI primarily orchestrates external tools (Terraform, Helm) rather than implementing algorithmic logic
- Most operations involve side effects (creating infrastructure, deploying applications)
- Behavior is highly dependent on external state (AWS resources, Kubernetes cluster)
- Operations are not pure functions with universal properties
- Testing focuses on integration with external systems rather than internal logic

Instead, the testing strategy emphasizes:
- Example-based unit tests for validation logic
- Integration tests with real or mocked external tools
- End-to-end tests in isolated test environments
- Snapshot tests for Terraform and Helm output

## Error Handling

### Error Categories

The CLI defines distinct error categories with specific handling strategies:

1. **Authentication Errors (1000-1099)**
   - No AWS credentials found
   - Expired credentials
   - Insufficient IAM permissions
   - Resolution: Display credential configuration instructions

2. **Validation Errors (1100-1199)**
   - Invalid app name format
   - Invalid environment type
   - Missing required flags
   - Resolution: Display validation rules and examples

3. **Terraform Errors (1200-1299)**
   - Terraform init failed
   - Terraform apply failed
   - State locked by another process
   - Resolution: Trigger rollback, display Terraform error output

4. **Helm Errors (1300-1399)**
   - Chart not found
   - Helm install failed
   - Pods not ready after timeout
   - Resolution: Trigger rollback, display Kubernetes events

5. **Network Errors (1400-1499)**
   - Connection timeout
   - DNS resolution failed
   - Resolution: Retry with exponential backoff

6. **Configuration Errors (1500-1599)**
   - Configuration file not found
   - Invalid YAML syntax
   - Resolution: Display configuration file examples

### Rollback Strategy

When errors occur during provisioning, the CLI implements automatic rollback:

1. **Helm Installation Failure**
   - Execute `helm uninstall` to remove Kubernetes resources
   - Proceed to Terraform rollback

2. **Terraform Apply Failure**
   - Execute `terraform destroy` to remove partially created infrastructure
   - Log rollback progress
   - Display manual cleanup instructions if rollback fails

3. **Pod Verification Failure**
   - Execute `helm uninstall`
   - Execute `terraform destroy`
   - Display pod logs and events for debugging

### Error Message Format

All error messages follow a consistent format:

```
❌ Error: <Category> (<Error Code>)

<User-friendly error message>

Details:
<Technical details from underlying tool>

Resolution:
<Step-by-step instructions to resolve the issue>

For more information, see the log file:
<Path to detailed log file>
```

### Logging Strategy

The CLI implements comprehensive logging:

1. **Console Output**: User-friendly messages with colors
   - Success messages in green
   - Warnings in yellow
   - Errors in red
   - Progress indicators for long operations

2. **File Logging**: Detailed logs for troubleshooting
   - Location: `~/.devplatform/logs/`
   - Format: JSON structured logs
   - Rotation: Keep last 10 log files
   - Content: All operations, API calls, and errors

3. **Verbosity Levels**
   - Default: User-friendly messages only
   - Verbose (`-v`): Include operation details
   - Debug (`--debug`): Include all API calls and responses

---

This design document provides the foundation for implementing the DevPlatform CLI. The architecture emphasizes simplicity, reliability, and developer experience while maintaining the flexibility to extend functionality in the future.
