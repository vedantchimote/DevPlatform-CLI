# DevPlatform CLI - Command Reference

Complete reference for all DevPlatform CLI commands, flags, and options.

## Table of Contents

- [Global Flags](#global-flags)
- [create Command](#create-command)
- [status Command](#status-command)
- [destroy Command](#destroy-command)
- [version Command](#version-command)
- [Error Codes](#error-codes)
- [Exit Codes](#exit-codes)
- [Environment Variables](#environment-variables)

## Global Flags

These flags are available for all commands:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--help` | `-h` | Display help information | - |
| `--verbose` | `-v` | Enable verbose output | false |
| `--debug` | | Enable debug logging | false |
| `--no-color` | | Disable colored output | false |
| `--config` | `-c` | Path to configuration file | `.devplatform.yaml` |

## create Command

Creates a complete isolated environment with infrastructure and application deployment.

### Synopsis

```bash
devplatform create --app <name> --env <type> [--provider <cloud>] [options]
```

### Description

The `create` command provisions a complete environment including:
- Network infrastructure (VPC/VNet)
- Database (RDS/Azure Database)
- Kubernetes namespace with resource quotas
- Application deployment via Helm
- Ingress for external access

### Required Flags

| Flag | Short | Type | Description | Validation |
|------|-------|------|-------------|------------|
| `--app` | `-a` | string | Application name | 3-32 chars, lowercase alphanumeric and hyphens |
| `--env` | `-e` | string | Environment type | Must be: dev, staging, or prod |

### Optional Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--provider` | `-p` | string | Cloud provider (aws or azure) | aws |
| `--dry-run` | | boolean | Preview changes without applying | false |
| `--values-file` | | string | Path to custom Helm values file | - |
| `--timeout` | | int | Operation timeout in minutes | 30 |

### Examples

#### Basic Usage

```bash
# Create dev environment on AWS (default)
devplatform create --app payment --env dev

# Create dev environment on Azure
devplatform create --app payment --env dev --provider azure
```

#### Advanced Usage

```bash
# Dry-run to preview changes
devplatform create --app payment --env staging --provider aws --dry-run

# Create with custom Helm values
devplatform create --app payment --env prod --provider azure \
  --values-file custom-values.yaml

# Create with verbose output
devplatform create --app payment --env dev --provider aws --verbose

# Create with debug logging
devplatform create --app payment --env dev --provider aws --debug
```

### Output

#### Success Output (AWS)

```
✓ Validating inputs...
✓ Checking AWS credentials...
✓ Initializing Terraform...
✓ Creating VPC... (10.0.0.0/16)
✓ Creating RDS instance... (db.t3.micro)
✓ Creating EKS namespace... (dev-payment-aws)
✓ Deploying application...
✓ Verifying pods...

Environment created successfully!

Database Endpoint: payment-dev.abc123.us-east-1.rds.amazonaws.com:5432
Ingress URL:       https://payment-dev.example.com
Namespace:         dev-payment-aws
Cloud Provider:    AWS (us-east-1)

To configure kubectl access:
  aws eks update-kubeconfig --name shared-devplatform-cluster --region us-east-1
  kubectl config set-context --current --namespace=dev-payment-aws

Estimated Monthly Cost: $72
```

#### Success Output (Azure)

```
✓ Validating inputs...
✓ Checking Azure credentials...
✓ Initializing Terraform...
✓ Creating VNet... (10.10.0.0/16)
✓ Creating Azure Database... (B_Gen5_1)
✓ Creating AKS namespace... (dev-payment-azure)
✓ Deploying application...
✓ Verifying pods...

Environment created successfully!

Database Endpoint: payment-dev.postgres.database.azure.com:5432
Ingress URL:       https://payment-dev.example.com
Namespace:         dev-payment-azure
Cloud Provider:    Azure (eastus)

To configure kubectl access:
  az aks get-credentials --name shared-devplatform-cluster --resource-group devplatform-rg
  kubectl config set-context --current --namespace=dev-payment-azure

Estimated Monthly Cost: $67
```

#### Dry-Run Output

```
Dry-run mode: No resources will be created

Planned Infrastructure Changes:
  + VPC (10.0.0.0/16)
    - 2 public subnets
    - 2 private subnets
    - 1 NAT gateway
    - Internet gateway
  
  + RDS Instance
    - Engine: PostgreSQL 14
    - Instance: db.t3.micro
    - Storage: 20 GB
    - Multi-AZ: false
  
  + EKS Namespace
    - Name: dev-payment-aws
    - Resource Quota: 2 CPU, 4Gi memory
  
  + Helm Release
    - Chart: devplatform-base
    - Replicas: 2

Estimated Monthly Cost: $72
  - VPC: $32
  - RDS: $15
  - EKS: $25
```

### Error Handling

The command automatically rolls back on failure:

```
✓ Creating VPC...
✓ Creating RDS instance...
✗ Creating EKS namespace... Failed

Error: Insufficient IAM permissions

Rolling back...
✓ Destroying RDS instance...
✓ Destroying VPC...

Rollback completed. No resources remain.

Resolution:
Grant the following IAM permissions:
  - eks:CreateNamespace
  - eks:TagResource
```

### Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error |
| 2 | Validation error |
| 3 | Authentication error |
| 4 | Terraform error |
| 5 | Helm error |

## status Command

Checks the health and status of an existing environment.

### Synopsis

```bash
devplatform status --app <name> --env <type> [--provider <cloud>] [options]
```

### Description

The `status` command queries:
- Terraform state for infrastructure resources
- Cloud provider APIs for resource health
- Kubernetes API for pod and ingress status
- Database connectivity

### Required Flags

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--app` | `-a` | string | Application name |
| `--env` | `-e` | string | Environment type |

### Optional Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--provider` | `-p` | string | Cloud provider (aws or azure) | aws |
| `--output` | `-o` | string | Output format (table, json, yaml) | table |
| `--watch` | `-w` | int | Watch mode (refresh every N seconds) | 0 (disabled) |

### Examples

```bash
# Check status (table format)
devplatform status --app payment --env dev --provider aws

# JSON output
devplatform status --app payment --env dev --provider azure --output json

# YAML output
devplatform status --app payment --env dev --provider aws --output yaml

# Watch mode (refresh every 5 seconds)
devplatform status --app payment --env dev --provider aws --watch 5
```

### Output

#### Table Format (AWS)

```
Environment Status: payment-dev (AWS)

Component       Status    Details
---------       ------    -------
VPC             OK        vpc-abc123 (10.0.0.0/16)
RDS             OK        payment-dev.abc123.us-east-1.rds.amazonaws.com
Namespace       OK        dev-payment-aws
Pods            OK        2/2 Ready
Ingress         OK        https://payment-dev.example.com

Last Updated: 2024-01-15 10:30:45 UTC
```

#### Table Format (Azure)

```
Environment Status: payment-dev (Azure)

Component       Status    Details
---------       ------    -------
VNet            OK        vnet-abc123 (10.10.0.0/16)
Azure Database  OK        payment-dev.postgres.database.azure.com
Namespace       OK        dev-payment-azure
Pods            OK        2/2 Ready
Ingress         OK        https://payment-dev.example.com

Last Updated: 2024-01-15 10:30:45 UTC
```

#### JSON Format

```json
{
  "app": "payment",
  "env": "dev",
  "provider": "aws",
  "status": "healthy",
  "components": {
    "vpc": {
      "status": "ok",
      "id": "vpc-abc123",
      "cidr": "10.0.0.0/16"
    },
    "rds": {
      "status": "ok",
      "endpoint": "payment-dev.abc123.us-east-1.rds.amazonaws.com",
      "engine": "postgres",
      "version": "14.7"
    },
    "namespace": {
      "status": "ok",
      "name": "dev-payment-aws"
    },
    "pods": {
      "status": "ok",
      "ready": 2,
      "total": 2,
      "pods": [
        {
          "name": "payment-7d8f9c5b6-abc12",
          "status": "Running",
          "ready": true
        },
        {
          "name": "payment-7d8f9c5b6-def34",
          "status": "Running",
          "ready": true
        }
      ]
    },
    "ingress": {
      "status": "ok",
      "url": "https://payment-dev.example.com"
    }
  },
  "last_updated": "2024-01-15T10:30:45Z"
}
```

#### Environment Not Found

```
Environment not found: payment-dev (AWS)

No Terraform state found for this environment.

To create this environment:
  devplatform create --app payment --env dev --provider aws
```

### Exit Codes

| Code | Description |
|------|-------------|
| 0 | Environment healthy |
| 1 | Environment degraded or not found |
| 3 | Authentication error |

## destroy Command

Destroys an existing environment and all associated resources.

### Synopsis

```bash
devplatform destroy --app <name> --env <type> [--provider <cloud>] [options]
```

### Description

The `destroy` command removes:
- Helm release and Kubernetes resources
- Database instance
- Network infrastructure
- Terraform state (unless --keep-state is specified)

**Warning:** This operation is irreversible. All data will be permanently deleted.

### Required Flags

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--app` | `-a` | string | Application name |
| `--env` | `-e` | string | Environment type |

### Optional Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--provider` | `-p` | string | Cloud provider (aws or azure) | aws |
| `--confirm` | `-y` | boolean | Skip confirmation prompt | false |
| `--force` | | boolean | Force destruction even if errors occur | false |
| `--keep-state` | | boolean | Keep Terraform state file | false |

### Examples

```bash
# Destroy with confirmation prompt
devplatform destroy --app payment --env dev --provider aws

# Destroy without confirmation
devplatform destroy --app payment --env dev --provider azure --confirm

# Force destroy (ignore errors)
devplatform destroy --app payment --env dev --provider aws --confirm --force

# Destroy but keep state file
devplatform destroy --app payment --env dev --provider aws --confirm --keep-state
```

### Output

#### With Confirmation Prompt

```
⚠ WARNING: This will destroy all resources for payment-dev on AWS

The following resources will be deleted:
  - VPC: vpc-abc123
  - RDS Instance: payment-dev
  - EKS Namespace: dev-payment-aws
  - All application pods and services

This action cannot be undone. All data will be permanently deleted.

Are you sure? (yes/no): yes

✓ Uninstalling Helm release...
✓ Deleting Kubernetes resources...
✓ Destroying RDS instance...
✓ Destroying VPC...
✓ Cleaning up Terraform state...

Environment destroyed successfully!

Estimated monthly savings: $72
```

#### Without Confirmation (--confirm flag)

```
✓ Uninstalling Helm release...
✓ Deleting Kubernetes resources...
✓ Destroying Azure Database...
✓ Destroying VNet...
✓ Cleaning up Terraform state...

Environment destroyed successfully!

Estimated monthly savings: $67
```

#### Partial Failure

```
✓ Uninstalling Helm release...
✓ Deleting Kubernetes resources...
✗ Destroying RDS instance... Failed

Error: RDS instance has deletion protection enabled

Manual cleanup required:
  1. Disable deletion protection:
     aws rds modify-db-instance --db-instance-identifier payment-dev --no-deletion-protection
  
  2. Retry destroy:
     devplatform destroy --app payment --env dev --provider aws --confirm
```

### Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error |
| 3 | Authentication error |
| 4 | Terraform error |
| 5 | Helm error |

## version Command

Displays version information for the CLI and dependencies.

### Synopsis

```bash
devplatform version [options]
```

### Description

The `version` command shows:
- CLI version number
- Git commit hash
- Build date
- Go version
- Platform/architecture
- Dependency versions (with --check-deps)

### Optional Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--short` | `-s` | boolean | Display only version number | false |
| `--check-deps` | | boolean | Check dependency versions | false |

### Examples

```bash
# Full version info
devplatform version

# Short version
devplatform version --short

# Check dependencies
devplatform version --check-deps
```

### Output

#### Full Version

```
DevPlatform CLI
Version:    1.0.0
Git Commit: a1b2c3d
Build Date: 2024-01-15T12:00:00Z
Go Version: go1.21.5
Platform:   linux/amd64
```

#### Short Version

```
1.0.0
```

#### With Dependency Check

```
DevPlatform CLI
Version:    1.0.0
Git Commit: a1b2c3d
Build Date: 2024-01-15T12:00:00Z
Go Version: go1.21.5
Platform:   linux/amd64

Dependencies:
  ✓ Terraform:  1.5.7 (minimum: 1.5.0)
  ✓ Helm:       3.12.0 (minimum: 3.0.0)
  ✓ kubectl:    1.27.3 (minimum: 1.27.0)
  ✓ AWS CLI:    2.13.0 (minimum: 2.0.0)
  ✓ Azure CLI:  2.50.0 (minimum: 2.0.0)

All dependencies satisfied!
```

#### Missing Dependencies

```
DevPlatform CLI
Version:    1.0.0

Dependencies:
  ✓ Terraform:  1.5.7 (minimum: 1.5.0)
  ✗ Helm:       Not found (minimum: 3.0.0)
  ✓ kubectl:    1.27.3 (minimum: 1.27.0)
  ✗ AWS CLI:    Not found (minimum: 2.0.0)
  ✓ Azure CLI:  2.50.0 (minimum: 2.0.0)

Missing dependencies! Please install:
  - Helm 3.0 or higher
  - AWS CLI 2.0 or higher
```

### Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success (all dependencies satisfied) |
| 1 | Missing or outdated dependencies |

## Error Codes

DevPlatform CLI uses specific error codes for different failure categories:

### Authentication Errors (1000-1099)

| Code | Description | Resolution |
|------|-------------|------------|
| 1001 | No AWS credentials found | Run `aws configure` |
| 1002 | AWS credentials expired | Refresh credentials or re-authenticate |
| 1003 | Insufficient IAM permissions | Grant required IAM permissions |
| 2001 | No Azure credentials found | Run `az login` |
| 2002 | Azure CLI not configured | Install and configure Azure CLI |
| 2003 | Azure subscription not found | Verify subscription ID |

### Validation Errors (1100-1199)

| Code | Description | Resolution |
|------|-------------|------------|
| 1101 | Invalid application name format | Use lowercase alphanumeric and hyphens (3-32 chars) |
| 1102 | Invalid environment type | Use dev, staging, or prod |
| 1103 | Missing required flag | Provide all required flags |
| 1104 | Invalid cloud provider | Use aws or azure |

### Terraform Errors (1200-1299)

| Code | Description | Resolution |
|------|-------------|------------|
| 1201 | Terraform init failed | Check Terraform configuration |
| 1202 | Terraform apply failed | Review Terraform error output |
| 1203 | State locked by another process | Wait for lock release or force unlock |
| 1204 | Backend configuration error | Verify backend settings |

### Helm Errors (1300-1399)

| Code | Description | Resolution |
|------|-------------|------------|
| 1301 | Helm chart not found | Verify chart path |
| 1302 | Helm install failed | Check Helm error output |
| 1303 | Pods not ready after timeout | Check pod logs and events |
| 1304 | Invalid values file | Verify YAML syntax |

### Network Errors (1400-1499)

| Code | Description | Resolution |
|------|-------------|------------|
| 1401 | Connection timeout | Check network connectivity |
| 1402 | DNS resolution failed | Verify DNS configuration |
| 1403 | API rate limit exceeded | Wait and retry |

### Configuration Errors (1500-1599)

| Code | Description | Resolution |
|------|-------------|------------|
| 1501 | Configuration file not found | Create .devplatform.yaml |
| 1502 | Invalid YAML syntax | Fix YAML syntax errors |
| 1503 | Invalid configuration value | Check configuration documentation |

### Azure-Specific Errors (2100-2199)

| Code | Description | Resolution |
|------|-------------|------------|
| 2101 | Azure resource creation failed | Check Azure error details |
| 2102 | Azure Storage backend error | Verify storage account configuration |
| 2103 | AKS connection failed | Check AKS cluster status |

## Exit Codes

| Code | Description | Commands |
|------|-------------|----------|
| 0 | Success | All |
| 1 | General error | All |
| 2 | Validation error | create |
| 3 | Authentication error | create, status, destroy |
| 4 | Terraform error | create, destroy |
| 5 | Helm error | create, destroy |
| 6 | Network error | All |
| 7 | Configuration error | All |

## Environment Variables

DevPlatform CLI respects the following environment variables:

### AWS

| Variable | Description | Example |
|----------|-------------|---------|
| `AWS_ACCESS_KEY_ID` | AWS access key | `AKIAIOSFODNN7EXAMPLE` |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key | `wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY` |
| `AWS_DEFAULT_REGION` | Default AWS region | `us-east-1` |
| `AWS_PROFILE` | AWS CLI profile | `devplatform` |

### Azure

| Variable | Description | Example |
|----------|-------------|---------|
| `AZURE_CLIENT_ID` | Service principal app ID | `12345678-1234-1234-1234-123456789012` |
| `AZURE_CLIENT_SECRET` | Service principal password | `secret` |
| `AZURE_TENANT_ID` | Azure tenant ID | `87654321-4321-4321-4321-210987654321` |
| `AZURE_SUBSCRIPTION_ID` | Azure subscription ID | `abcdef12-3456-7890-abcd-ef1234567890` |

### DevPlatform CLI

| Variable | Description | Example |
|----------|-------------|---------|
| `DEVPLATFORM_CONFIG` | Path to config file | `/path/to/config.yaml` |
| `DEVPLATFORM_LOG_LEVEL` | Log level | `debug` |
| `DEVPLATFORM_NO_COLOR` | Disable colored output | `true` |

## Common Use Cases

### CI/CD Integration

```bash
# Non-interactive mode with environment variables
export AWS_ACCESS_KEY_ID="..."
export AWS_SECRET_ACCESS_KEY="..."
export AWS_DEFAULT_REGION="us-east-1"

devplatform create \
  --app myapp \
  --env staging \
  --provider aws \
  --no-color \
  --confirm
```

### Automation Scripts

```bash
#!/bin/bash
set -e

# Create environment
devplatform create --app myapp --env dev --provider aws --confirm

# Wait for pods to be ready
while true; do
  STATUS=$(devplatform status --app myapp --env dev --provider aws --output json | jq -r '.components.pods.status')
  if [ "$STATUS" = "ok" ]; then
    break
  fi
  sleep 10
done

# Run tests
kubectl run test --image=mytest --namespace=dev-myapp-aws

# Cleanup
devplatform destroy --app myapp --env dev --provider aws --confirm
```

### Multi-Cloud Deployment

```bash
# Deploy to both AWS and Azure
devplatform create --app myapp --env prod --provider aws
devplatform create --app myapp --env prod --provider azure

# Check status on both
devplatform status --app myapp --env prod --provider aws
devplatform status --app myapp --env prod --provider azure
```

## See Also

- [README.md](../README.md) - Main documentation
- [Architecture](architecture.md) - System architecture
- [Troubleshooting](troubleshooting.md) - Common issues and solutions
- [AWS Guide](../aws/overview.mdx) - AWS-specific documentation
- [Azure Guide](../azure/overview.mdx) - Azure-specific documentation
