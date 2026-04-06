# Terraform Module Structure

Complete documentation for DevPlatform CLI Terraform modules, organization, and environment-specific configurations.

## Table of Contents

- [Overview](#overview)
- [Module Organization](#module-organization)
- [AWS Modules](#aws-modules)
- [Azure Modules](#azure-modules)
- [Environment Configurations](#environment-configurations)
- [Variables](#variables)
- [Outputs](#outputs)
- [State Management](#state-management)
- [Customization](#customization)

## Overview

DevPlatform CLI uses Terraform modules to provision infrastructure on AWS and Azure. The modules are organized by cloud provider and resource type, with environment-specific configurations for dev, staging, and production.

### Design Principles

1. **Modularity**: Each resource type (network, database, Kubernetes) is a separate module
2. **Reusability**: Modules are reusable across environments with different variables
3. **Cloud Abstraction**: Similar structure for AWS and Azure modules
4. **Environment-Specific**: Configurations tailored for dev, staging, and prod workloads
5. **Best Practices**: Security, high availability, and cost optimization built-in

## Module Organization

```
terraform/
├── modules/
│   ├── aws/
│   │   ├── network/          # VPC, subnets, NAT gateways
│   │   ├── database/         # RDS PostgreSQL
│   │   ├── eks-tenant/       # EKS namespace and RBAC
│   │   └── TAGGING.md        # Tagging strategy
│   └── azure/
│       ├── network/          # VNet, subnets, NAT gateways
│       ├── database/         # Azure Database for PostgreSQL
│       ├── k8s-tenant/       # AKS namespace and RBAC
│       └── TAGGING.md        # Tagging strategy
└── environments/
    ├── aws/
    │   ├── dev/
    │   │   └── terraform.tfvars
    │   ├── staging/
    │   │   └── terraform.tfvars
    │   └── prod/
    │       └── terraform.tfvars
    └── azure/
        ├── dev/
        │   └── terraform.tfvars
        ├── staging/
        │   └── terraform.tfvars
        └── prod/
            └── terraform.tfvars
```

## AWS Modules

### Network Module (aws/network)

Creates VPC infrastructure with public and private subnets, NAT gateways, and security groups.

**Location**: `terraform/modules/aws/network/`

**Resources Created**:
- VPC with configurable CIDR block
- Public subnets (one per availability zone)
- Private subnets (one per availability zone)
- Internet Gateway
- NAT Gateways (configurable count)
- Route tables and associations
- Security groups for database access
- VPC Flow Logs

**Key Variables**:

| Variable | Type | Description | Example |
|----------|------|-------------|---------|
| `app_name` | string | Application name | `payment` |
| `env_type` | string | Environment type | `dev` |
| `vpc_cidr` | string | VPC CIDR block | `10.0.0.0/16` |
| `availability_zones` | list(string) | AZs to use | `["us-east-1a", "us-east-1b"]` |
| `nat_gateway_count` | number | Number of NAT gateways | `1` (dev), `2` (staging), `3` (prod) |
| `enable_flow_logs` | bool | Enable VPC flow logs | `true` |
| `tags` | map(string) | Additional tags | `{}` |

**Outputs**:

| Output | Description |
|--------|-------------|
| `vpc_id` | VPC ID |
| `vpc_cidr` | VPC CIDR block |
| `public_subnet_ids` | List of public subnet IDs |
| `private_subnet_ids` | List of private subnet IDs |
| `nat_gateway_ids` | List of NAT gateway IDs |
| `db_security_group_id` | Security group ID for database |

**Example Usage**:

```hcl
module "network" {
  source = "../../modules/aws/network"

  app_name           = var.app_name
  env_type           = var.env_type
  vpc_cidr           = "10.0.0.0/16"
  availability_zones = ["us-east-1a", "us-east-1b"]
  nat_gateway_count  = 1
  enable_flow_logs   = true

  tags = {
    ManagedBy = "devplatform-cli"
  }
}
```

### Database Module (aws/database)

Creates RDS PostgreSQL instance with automated backups and security configuration.

**Location**: `terraform/modules/aws/database/`

**Resources Created**:
- RDS PostgreSQL instance
- DB subnet group
- DB parameter group
- Random password (stored in Secrets Manager)
- CloudWatch alarms for monitoring

**Key Variables**:

| Variable | Type | Description | Example |
|----------|------|-------------|---------|
| `app_name` | string | Application name | `payment` |
| `env_type` | string | Environment type | `dev` |
| `instance_class` | string | RDS instance class | `db.t3.micro` |
| `allocated_storage` | number | Storage in GB | `20` |
| `engine_version` | string | PostgreSQL version | `14.7` |
| `multi_az` | bool | Enable Multi-AZ | `false` (dev), `true` (prod) |
| `backup_retention_days` | number | Backup retention | `7` (dev), `30` (prod) |
| `vpc_id` | string | VPC ID | From network module |
| `subnet_ids` | list(string) | Subnet IDs | From network module |
| `security_group_ids` | list(string) | Security group IDs | From network module |

**Outputs**:

| Output | Description |
|--------|-------------|
| `db_endpoint` | Database endpoint |
| `db_port` | Database port |
| `db_name` | Database name |
| `db_username` | Database username |
| `secret_arn` | Secrets Manager ARN for password |

**Example Usage**:

```hcl
module "database" {
  source = "../../modules/aws/database"

  app_name              = var.app_name
  env_type              = var.env_type
  instance_class        = "db.t3.micro"
  allocated_storage     = 20
  engine_version        = "14.7"
  multi_az              = false
  backup_retention_days = 7

  vpc_id             = module.network.vpc_id
  subnet_ids         = module.network.private_subnet_ids
  security_group_ids = [module.network.db_security_group_id]
}
```

### EKS Tenant Module (aws/eks-tenant)

Creates Kubernetes namespace with resource quotas and IRSA configuration.

**Location**: `terraform/modules/aws/eks-tenant/`

**Resources Created**:
- Kubernetes namespace
- Resource quota
- Limit range
- Service account
- IAM role for IRSA
- Role binding

**Key Variables**:

| Variable | Type | Description | Example |
|----------|------|-------------|---------|
| `app_name` | string | Application name | `payment` |
| `env_type` | string | Environment type | `dev` |
| `cluster_name` | string | EKS cluster name | `shared-devplatform-cluster` |
| `cpu_quota` | string | CPU quota | `2` (dev), `10` (prod) |
| `memory_quota` | string | Memory quota | `4Gi` (dev), `20Gi` (prod) |
| `oidc_provider_arn` | string | OIDC provider ARN | From EKS cluster |

**Outputs**:

| Output | Description |
|--------|-------------|
| `namespace` | Kubernetes namespace name |
| `service_account` | Service account name |
| `iam_role_arn` | IAM role ARN for IRSA |

**Example Usage**:

```hcl
module "eks_tenant" {
  source = "../../modules/aws/eks-tenant"

  app_name     = var.app_name
  env_type     = var.env_type
  cluster_name = "shared-devplatform-cluster"
  cpu_quota    = "2"
  memory_quota = "4Gi"

  oidc_provider_arn = data.aws_eks_cluster.main.identity[0].oidc[0].issuer
}
```

## Azure Modules

### Network Module (azure/network)

Creates VNet infrastructure with public and private subnets, NAT gateways, and network security groups.

**Location**: `terraform/modules/azure/network/`

**Resources Created**:
- Virtual Network (VNet)
- Public subnets (one per availability zone)
- Private subnets (one per availability zone)
- NAT Gateways (configurable count)
- Public IPs for NAT gateways
- Network Security Groups
- Network Watcher flow logs

**Key Variables**:

| Variable | Type | Description | Example |
|----------|------|-------------|---------|
| `app_name` | string | Application name | `payment` |
| `env_type` | string | Environment type | `dev` |
| `resource_group_name` | string | Resource group name | `devplatform-rg` |
| `location` | string | Azure region | `eastus` |
| `vnet_cidr` | string | VNet CIDR block | `10.10.0.0/16` |
| `availability_zones` | list(string) | Zones to use | `["1", "2"]` |
| `nat_gateway_count` | number | Number of NAT gateways | `1` (dev), `3` (prod) |
| `enable_flow_logs` | bool | Enable flow logs | `true` |

**Outputs**:

| Output | Description |
|--------|-------------|
| `vnet_id` | VNet ID |
| `vnet_name` | VNet name |
| `public_subnet_ids` | List of public subnet IDs |
| `private_subnet_ids` | List of private subnet IDs |
| `db_nsg_id` | NSG ID for database |

**Example Usage**:

```hcl
module "network" {
  source = "../../modules/azure/network"

  app_name            = var.app_name
  env_type            = var.env_type
  resource_group_name = azurerm_resource_group.main.name
  location            = var.location
  vnet_cidr           = "10.10.0.0/16"
  availability_zones  = ["1", "2"]
  nat_gateway_count   = 1
  enable_flow_logs    = true
}
```

### Database Module (azure/database)

Creates Azure Database for PostgreSQL Flexible Server with automated backups.

**Location**: `terraform/modules/azure/database/`

**Resources Created**:
- PostgreSQL Flexible Server
- Private DNS zone
- VNet integration
- Firewall rules
- Random password (stored in Key Vault)
- Diagnostic settings

**Key Variables**:

| Variable | Type | Description | Example |
|----------|------|-------------|---------|
| `app_name` | string | Application name | `payment` |
| `env_type` | string | Environment type | `dev` |
| `resource_group_name` | string | Resource group name | `devplatform-rg` |
| `location` | string | Azure region | `eastus` |
| `sku_name` | string | Database SKU | `B_Standard_B1ms` |
| `storage_mb` | number | Storage in MB | `32768` |
| `version` | string | PostgreSQL version | `14` |
| `zone` | string | Availability zone | `1` |
| `high_availability_enabled` | bool | Enable HA | `false` (dev), `true` (prod) |
| `backup_retention_days` | number | Backup retention | `7` (dev), `30` (prod) |
| `vnet_id` | string | VNet ID | From network module |
| `subnet_id` | string | Subnet ID | From network module |

**Outputs**:

| Output | Description |
|--------|-------------|
| `db_fqdn` | Database FQDN |
| `db_name` | Database name |
| `db_username` | Database username |
| `key_vault_secret_id` | Key Vault secret ID for password |

**Example Usage**:

```hcl
module "database" {
  source = "../../modules/azure/database"

  app_name                  = var.app_name
  env_type                  = var.env_type
  resource_group_name       = azurerm_resource_group.main.name
  location                  = var.location
  sku_name                  = "B_Standard_B1ms"
  storage_mb                = 32768
  version                   = "14"
  zone                      = "1"
  high_availability_enabled = false
  backup_retention_days     = 7

  vnet_id   = module.network.vnet_id
  subnet_id = module.network.private_subnet_ids[0]
}
```

### K8s Tenant Module (azure/k8s-tenant)

Creates Kubernetes namespace with resource quotas and Workload Identity configuration.

**Location**: `terraform/modules/azure/k8s-tenant/`

**Resources Created**:
- Kubernetes namespace
- Resource quota
- Limit range
- Service account
- Azure AD application
- Federated identity credential
- Role assignment

**Key Variables**:

| Variable | Type | Description | Example |
|----------|------|-------------|---------|
| `app_name` | string | Application name | `payment` |
| `env_type` | string | Environment type | `dev` |
| `cluster_name` | string | AKS cluster name | `shared-devplatform-cluster` |
| `cpu_quota` | string | CPU quota | `2` (dev), `10` (prod) |
| `memory_quota` | string | Memory quota | `4Gi` (dev), `20Gi` (prod) |
| `oidc_issuer_url` | string | OIDC issuer URL | From AKS cluster |

**Outputs**:

| Output | Description |
|--------|-------------|
| `namespace` | Kubernetes namespace name |
| `service_account` | Service account name |
| `client_id` | Azure AD application client ID |

**Example Usage**:

```hcl
module "k8s_tenant" {
  source = "../../modules/azure/k8s-tenant"

  app_name     = var.app_name
  env_type     = var.env_type
  cluster_name = "shared-devplatform-cluster"
  cpu_quota    = "2"
  memory_quota = "4Gi"

  oidc_issuer_url = data.azurerm_kubernetes_cluster.main.oidc_issuer_url
}
```

## Environment Configurations

Environment-specific configurations are stored in `terraform/environments/{provider}/{env}/terraform.tfvars`.

### Development Environment

**Purpose**: Development and testing with minimal resources

**AWS Configuration** (`terraform/environments/aws/dev/terraform.tfvars`):

```hcl
# Network
vpc_cidr           = "10.0.0.0/16"
availability_zones = ["us-east-1a", "us-east-1b"]
nat_gateway_count  = 1  # Cost optimization

# Database
db_instance_class        = "db.t3.micro"
db_allocated_storage     = 20
db_multi_az              = false
db_backup_retention_days = 7

# Kubernetes
k8s_cpu_quota    = "2"
k8s_memory_quota = "4Gi"
```

**Azure Configuration** (`terraform/environments/azure/dev/terraform.tfvars`):

```hcl
# Network
vnet_cidr          = "10.10.0.0/16"
availability_zones = ["1", "2"]
nat_gateway_count  = 1  # Cost optimization

# Database
db_sku_name                  = "B_Standard_B1ms"
db_storage_mb                = 32768
db_high_availability_enabled = false
db_backup_retention_days     = 7

# Kubernetes
k8s_cpu_quota    = "2"
k8s_memory_quota = "4Gi"
```

### Staging Environment

**Purpose**: Pre-production testing with realistic resources

**AWS Configuration** (`terraform/environments/aws/staging/terraform.tfvars`):

```hcl
# Network
vpc_cidr           = "10.1.0.0/16"
availability_zones = ["us-east-1a", "us-east-1b"]
nat_gateway_count  = 2  # High availability

# Database
db_instance_class        = "db.t3.small"
db_allocated_storage     = 50
db_multi_az              = false
db_backup_retention_days = 14

# Kubernetes
k8s_cpu_quota    = "5"
k8s_memory_quota = "10Gi"
```

**Azure Configuration** (`terraform/environments/azure/staging/terraform.tfvars`):

```hcl
# Network
vnet_cidr          = "10.11.0.0/16"
availability_zones = ["1", "2"]
nat_gateway_count  = 2  # High availability

# Database
db_sku_name                  = "B_Standard_B2s"
db_storage_mb                = 65536
db_high_availability_enabled = false
db_backup_retention_days     = 14

# Kubernetes
k8s_cpu_quota    = "5"
k8s_memory_quota = "10Gi"
```

### Production Environment

**Purpose**: Production workloads with high availability and performance

**AWS Configuration** (`terraform/environments/aws/prod/terraform.tfvars`):

```hcl
# Network
vpc_cidr           = "10.2.0.0/16"
availability_zones = ["us-east-1a", "us-east-1b", "us-east-1c"]
nat_gateway_count  = 3  # Maximum availability

# Database
db_instance_class        = "db.r6g.large"
db_allocated_storage     = 100
db_multi_az              = true
db_backup_retention_days = 30

# Kubernetes
k8s_cpu_quota    = "10"
k8s_memory_quota = "20Gi"
```

**Azure Configuration** (`terraform/environments/azure/prod/terraform.tfvars`):

```hcl
# Network
vnet_cidr          = "10.12.0.0/16"
availability_zones = ["1", "2", "3"]
nat_gateway_count  = 3  # Maximum availability

# Database
db_sku_name                  = "GP_Standard_D2s_v3"
db_storage_mb                = 131072
db_high_availability_enabled = true
db_backup_retention_days     = 30

# Kubernetes
k8s_cpu_quota    = "10"
k8s_memory_quota = "20Gi"
```

## Variables

### Common Variables

These variables are used across all modules:

| Variable | Type | Required | Description |
|----------|------|----------|-------------|
| `app_name` | string | Yes | Application name (3-32 chars, lowercase alphanumeric and hyphens) |
| `env_type` | string | Yes | Environment type (dev, staging, prod) |
| `cloud_provider` | string | Yes | Cloud provider (aws or azure) |
| `tags` | map(string) | No | Additional resource tags |

### AWS-Specific Variables

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `region` | string | `us-east-1` | AWS region |
| `vpc_cidr` | string | - | VPC CIDR block |
| `availability_zones` | list(string) | - | List of AZs |
| `db_instance_class` | string | - | RDS instance class |
| `db_engine_version` | string | `14.7` | PostgreSQL version |

### Azure-Specific Variables

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `location` | string | `eastus` | Azure region |
| `resource_group_name` | string | - | Resource group name |
| `vnet_cidr` | string | - | VNet CIDR block |
| `db_sku_name` | string | - | Database SKU |
| `db_version` | string | `14` | PostgreSQL version |

## Outputs

### Network Module Outputs

**AWS**:
```hcl
output "vpc_id" {
  value = aws_vpc.main.id
}

output "private_subnet_ids" {
  value = aws_subnet.private[*].id
}
```

**Azure**:
```hcl
output "vnet_id" {
  value = azurerm_virtual_network.main.id
}

output "private_subnet_ids" {
  value = azurerm_subnet.private[*].id
}
```

### Database Module Outputs

**AWS**:
```hcl
output "db_endpoint" {
  value = aws_db_instance.main.endpoint
}

output "secret_arn" {
  value = aws_secretsmanager_secret.db_password.arn
}
```

**Azure**:
```hcl
output "db_fqdn" {
  value = azurerm_postgresql_flexible_server.main.fqdn
}

output "key_vault_secret_id" {
  value = azurerm_key_vault_secret.db_password.id
}
```

## State Management

### AWS Backend Configuration

```hcl
terraform {
  backend "s3" {
    bucket         = "terraform-state-bucket"
    key            = "devplatform/${app_name}/${env_type}/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "terraform-locks"
    encrypt        = true
  }
}
```

### Azure Backend Configuration

```hcl
terraform {
  backend "azurerm" {
    storage_account_name = "tfstatestorage"
    container_name       = "tfstate"
    key                  = "devplatform/${app_name}/${env_type}/terraform.tfstate"
    resource_group_name  = "terraform-state-rg"
  }
}
```

### State Locking

- **AWS**: Uses DynamoDB table for state locking
- **Azure**: Uses blob lease for state locking
- **Isolation**: Each app+env+provider combination has unique state key
- **Concurrent Operations**: State locking prevents concurrent modifications

## Customization

### Override Default Values

Create a `.devplatform.yaml` file:

```yaml
environments:
  dev:
    terraform:
      vpc_cidr: "10.20.0.0/16"
      nat_gateway_count: 2
    database:
      instance_class: "db.t3.small"
      allocated_storage: 50
```

### Add Custom Modules

1. Create module in `terraform/modules/{provider}/custom-module/`
2. Reference in environment configuration
3. Update CLI to pass variables

### Modify Existing Modules

1. Edit module files in `terraform/modules/{provider}/{module}/`
2. Update variables and outputs as needed
3. Test with `terraform plan`

## Best Practices

1. **Use Remote State**: Always use S3 or Azure Storage for state
2. **Enable State Locking**: Prevent concurrent modifications
3. **Tag Resources**: Use consistent tagging for cost tracking
4. **Version Modules**: Pin module versions in production
5. **Validate Changes**: Use `terraform plan` before applying
6. **Backup State**: Enable versioning on state storage
7. **Secure Secrets**: Use Secrets Manager or Key Vault
8. **Monitor Costs**: Review resource costs regularly

## See Also

- [README.md](../README.md) - Main documentation
- [Command Reference](COMMAND_REFERENCE.md) - CLI commands
- [Helm Charts](HELM_CHARTS.md) - Helm chart documentation
- [AWS Tagging](../terraform/modules/aws/TAGGING.md) - AWS tagging strategy
- [Azure Tagging](../terraform/modules/azure/TAGGING.md) - Azure tagging strategy
