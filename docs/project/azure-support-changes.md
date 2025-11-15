# Azure Support - Key Changes Summary

## Overview
This document summarizes the key architectural changes required to add Azure support alongside AWS in the DevPlatform CLI.

## Requirements Document Changes

### Updated Requirements
1. **Requirement 1**: Added Cloud_Provider parameter and validation
2. **Requirement 2**: Renamed to "Cloud Provider Authentication" with AWS and Azure support
3. **Requirement 3**: Added Azure Storage backend support for Terraform state
4. **Requirement 5**: Added Azure resource querying (VNet, Azure Database, AKS)
5. **Requirement 8**: Added Cloud_Provider validation (aws or azure)
6. **Requirement 9**: Added Azure Storage backend and blob lease locking
7. **Requirement 10**: Added cloud-specific Terraform modules
8. **Requirement 11**: Added --provider flag with aws as default
9. **Requirement 13**: Added Azure kubeconfig commands (az aks get-credentials)
10. **Requirement 14**: Added Cloud_Provider to resource tags
11. **Requirement 15**: Added Cloud_Provider to state key isolation
12. **Requirement 17**: Added Azure subscription to configuration
13. **Requirement 19**: Added Azure CLI (az) version checking
14. **Requirement 20**: Added Azure pricing data
15. **Requirement 22**: Added VNet provisioning for Azure
16. **Requirement 23**: Added Azure Database and Key Vault support
17. **Requirement 24**: Added Azure AD Workload Identity
18. **Requirement 25**: Added Azure CLI verification

### New Requirements
26. **Cloud Provider Abstraction**: Interface for consistent multi-cloud experience
27. **Azure-Specific Authentication**: Azure CLI, service principal, managed identity
28. **Azure Resource Provisioning**: VNet, Azure Database, AKS, NSG, Network Watcher
29. **Multi-Cloud Configuration**: Separate AWS and Azure config sections
30. **Cloud Provider Migration Support**: Resource mapping, cost comparison, migration guidance

## Design Document Changes

### Architecture Changes
1. **New Layer**: Cloud Provider Abstraction Layer between Business Logic and Infrastructure Adapters
2. **Provider Interface**: Common interface for AWS and Azure implementations
3. **Multi-Cloud State Management**: Support for both S3+DynamoDB and Azure Storage backends

### Component Changes

#### New Components
1. **internal/provider/** package
   - `provider.go`: CloudProvider interface
   - `aws/provider.go`: AWS implementation
   - `azure/provider.go`: Azure implementation
   - `factory.go`: Provider factory based on Cloud_Provider flag

2. **internal/azure/** package
   - `auth.go`: Azure authentication
   - `kubeconfig.go`: AKS kubeconfig management
   - `pricing.go`: Azure cost calculation

#### Updated Components
1. **internal/config/**
   - Add `CloudProvider` field to GlobalConfig
   - Add `AzureConfig` struct (subscription_id, location, tenant_id)
   - Add provider-specific validation

2. **internal/terraform/**
   - Support multiple backend types (S3, Azure Storage)
   - Provider-specific state key generation
   - Cloud-specific module path resolution

3. **internal/aws/** (existing)
   - Refactor to implement CloudProvider interface
   - Keep AWS-specific logic isolated

### Data Model Changes

```go
// New Provider Interface
type CloudProvider interface {
    ValidateCredentials() error
    GetCallerIdentity() (*Identity, error)
    UpdateKubeconfig(clusterName, region string) error
    CalculateNetworkCost(config NetworkConfig) float64
    CalculateDatabaseCost(config DatabaseConfig) float64
    CalculateK8sCost(config K8sConfig) float64
    GetTerraformBackend() BackendConfig
    GetModulePath(moduleType string) string
}

// Updated Configuration
type GlobalConfig struct {
    CloudProvider string      `yaml:"cloud_provider"` // aws or azure
    Region        string      `yaml:"region"`         // AWS region
    AWSProfile    string      `yaml:"aws_profile"`
    Azure         AzureConfig `yaml:"azure"`
    Timeout       int         `yaml:"timeout"`
    LogLevel      string      `yaml:"log_level"`
}

type AzureConfig struct {
    SubscriptionID string `yaml:"subscription_id"`
    Location       string `yaml:"location"`
    TenantID       string `yaml:"tenant_id"`
}

// Backend Config supports multiple types
type BackendConfig struct {
    Type          string `yaml:"type"` // s3 or azurerm
    // AWS S3 Backend
    Bucket        string `yaml:"bucket"`
    DynamoDBTable string `yaml:"dynamodb_table"`
    Region        string `yaml:"region"`
    // Azure Storage Backend
    StorageAccount string `yaml:"storage_account"`
    ContainerName  string `yaml:"container_name"`
    ResourceGroup  string `yaml:"resource_group"`
}
```

### Terraform Module Structure

```
terraform/
├── modules/
│   ├── aws/
│   │   ├── network/      # VPC module
│   │   ├── database/     # RDS module
│   │   └── k8s-tenant/   # EKS namespace module
│   └── azure/
│       ├── network/      # VNet module
│       ├── database/     # Azure Database module
│       └── k8s-tenant/   # AKS namespace module
└── environments/
    ├── aws/
    │   ├── dev/
    │   ├── staging/
    │   └── prod/
    └── azure/
        ├── dev/
        ├── staging/
        └── prod/
```

## Tasks Document Changes

### New Tasks

**After Task 1 (Project Setup):**
- 1.5: Add Azure SDK dependencies (github.com/Azure/azure-sdk-for-go)

**After Task 7 (AWS utilities):**
- 7.5: Implement cloud provider interface
- 7.6: Create provider factory
- 7.7: Implement Azure authentication
- 7.8: Implement Azure kubeconfig management
- 7.9: Implement Azure cost calculation

**After Task 9 (Terraform modules):**
- 9.6: Create Azure network module (VNet, subnets, NSG)
- 9.7: Create Azure database module (Azure Database for PostgreSQL)
- 9.8: Create Azure K8s tenant module (AKS namespace, Workload Identity)
- 9.9: Create Azure environment-specific configurations

**New Task 25:**
- 25: Implement multi-cloud testing
  - 25.1: Test AWS provisioning end-to-end
  - 25.2: Test Azure provisioning end-to-end
  - 25.3: Test switching between providers
  - 25.4: Test concurrent multi-cloud operations

### Updated Tasks

**Task 2.2 (Version command):**
- Add Azure CLI (az) version checking

**Task 3.1 (Configuration structures):**
- Add CloudProvider, AzureConfig fields

**Task 4.1 (Input validation):**
- Add Cloud_Provider validation (aws or azure)

**Task 8.3 (State management):**
- Support both S3 and Azure Storage backends
- Provider-specific state key generation

**Task 13.1 (Create command):**
- Add --provider flag with aws as default

**Task 13.2 (Create orchestration):**
- Use provider factory to get cloud-specific implementation
- Validate cloud-specific credentials

**Task 16.2 (Status checking):**
- Query cloud-specific resources based on provider

**Task 21.1 (Documentation):**
- Document multi-cloud support
- Document Azure-specific setup
- Document provider migration

## Configuration File Example

```yaml
# Global settings
global:
  cloud_provider: aws  # or azure
  timeout: 30
  log_level: info

# AWS-specific settings
aws:
  region: us-east-1
  profile: default

# Azure-specific settings
azure:
  subscription_id: "12345678-1234-1234-1234-123456789012"
  location: eastus
  tenant_id: "87654321-4321-4321-4321-210987654321"

# Environment-specific settings (same for both clouds)
environments:
  dev:
    network_cidr: 10.0.0.0/16
    database_instance_class: small
    database_allocated_storage: 20
    k8s_node_count: 2
    
  staging:
    network_cidr: 10.1.0.0/16
    database_instance_class: medium
    database_allocated_storage: 100
    k8s_node_count: 3
    
  prod:
    network_cidr: 10.2.0.0/16
    database_instance_class: large
    database_allocated_storage: 500
    database_multi_zone: true
    k8s_node_count: 5

# Terraform settings
terraform:
  backend:
    # For AWS
    type: s3
    bucket: terraform-state-bucket
    dynamodb_table: terraform-locks
    region: us-east-1
    
    # For Azure (alternative)
    # type: azurerm
    # storage_account: tfstatestorage
    # container_name: tfstate
    # resource_group: terraform-state-rg
  
  modules_path: ./terraform/modules
  
# Helm settings (same for both clouds)
helm:
  chart_path: ./charts/devplatform-base
  default_values:
    image:
      repository: nginx
      tag: latest
    resources:
      requests:
        cpu: 100m
        memory: 128Mi
```

## CLI Usage Examples

### AWS Provisioning
```bash
# Using default provider (aws)
devplatform create --app payment --env dev

# Explicitly specifying AWS
devplatform create --app payment --env dev --provider aws

# Check status
devplatform status --app payment --env dev --provider aws

# Destroy
devplatform destroy --app payment --env dev --provider aws --confirm
```

### Azure Provisioning
```bash
# Provision on Azure
devplatform create --app payment --env dev --provider azure

# Check status
devplatform status --app payment --env dev --provider azure

# Destroy
devplatform destroy --app payment --env dev --provider azure --confirm
```

### Multi-Cloud Operations
```bash
# Provision same app on both clouds
devplatform create --app payment --env dev --provider aws
devplatform create --app payment --env dev --provider azure

# Check status of both
devplatform status --app payment --env dev --provider aws
devplatform status --app payment --env dev --provider azure

# Cost comparison
devplatform create --app payment --env prod --provider aws --dry-run
devplatform create --app payment --env prod --provider azure --dry-run
```

## Implementation Priority

### Phase 1: Core Multi-Cloud Support (MVP)
1. Provider interface and factory
2. Azure authentication
3. Azure Terraform modules (network, database, k8s-tenant)
4. Configuration updates
5. Basic Azure provisioning

### Phase 2: Feature Parity
1. Azure cost calculation
2. Azure kubeconfig management
3. Azure-specific error handling
4. Multi-cloud testing

### Phase 3: Advanced Features
1. Cloud provider migration tools
2. Cost comparison features
3. Multi-cloud monitoring
4. Advanced documentation

## Testing Strategy

### Unit Tests
- Provider interface implementations
- Azure authentication
- Azure cost calculation
- Configuration validation for both clouds

### Integration Tests
- Azure Terraform module execution
- Azure resource provisioning
- AKS namespace creation
- Azure Storage backend operations

### End-to-End Tests
- Full AWS provisioning workflow
- Full Azure provisioning workflow
- Switching between providers
- Concurrent multi-cloud operations

## Documentation Updates

### README.md
- Add Azure support section
- Update prerequisites (add Azure CLI)
- Add Azure setup instructions
- Add multi-cloud examples

### New Documentation Files
- `docs/azure-setup.md`: Azure-specific setup guide
- `docs/multi-cloud.md`: Multi-cloud usage patterns
- `docs/provider-comparison.md`: AWS vs Azure resource mapping
- `docs/migration-guide.md`: Migrating between cloud providers

## Breaking Changes

### Backward Compatibility
- `--provider` flag defaults to `aws` for backward compatibility
- Existing configurations without `cloud_provider` field default to AWS
- Existing Terraform state keys remain unchanged for AWS

### Migration Path
1. Existing users continue to work without changes
2. Add `cloud_provider: aws` to configuration for explicitness
3. Azure users add Azure-specific configuration
4. No breaking changes to existing AWS deployments

## Success Criteria

1. ✅ Provision environments on AWS (existing functionality maintained)
2. ✅ Provision environments on Azure with same CLI commands
3. ✅ Consistent developer experience across both clouds
4. ✅ Cloud-specific features properly abstracted
5. ✅ Comprehensive documentation for both clouds
6. ✅ All tests passing for both AWS and Azure
7. ✅ Cost estimation working for both clouds
8. ✅ State management working for both backends
