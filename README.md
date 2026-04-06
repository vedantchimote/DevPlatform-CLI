# DevPlatform CLI

A powerful command-line tool for deploying and managing development environments on AWS and Azure with automated infrastructure provisioning and Kubernetes orchestration. Reduce environment provisioning time from 2 days to 3 minutes.

## 🚀 Quick Start

```bash
# Create a development environment on AWS (default)
devplatform create --app myapp --env dev

# Create a development environment on Azure
devplatform create --app myapp --env dev --provider azure

# Check environment status
devplatform status --app myapp --env dev --provider aws

# Destroy environment when done
devplatform destroy --app myapp --env dev --confirm
```

## ✨ Features

- **Multi-Cloud Support**: Deploy to AWS or Azure with a single command
- **Infrastructure as Code**: Automated Terraform provisioning (VPC/VNet, RDS/Azure Database, EKS/AKS)
- **Kubernetes Integration**: Deploy applications to EKS (AWS) or AKS (Azure) with Helm
- **Environment Management**: Separate dev, staging, and production environments with right-sized resources
- **State Management**: Secure remote state storage with locking (S3+DynamoDB or Azure Storage)
- **Cost Optimization**: Environment-specific resource sizing and cost estimation
- **Developer Self-Service**: No DevOps tickets required - provision complete environments in minutes

## 📋 Prerequisites

Before installing DevPlatform CLI, ensure you have the following tools installed:

### Required Tools

| Tool | Minimum Version | Purpose |
|------|----------------|---------|
| **Terraform** | 1.5+ | Infrastructure provisioning |
| **Helm** | 3.0+ | Kubernetes application deployment |
| **kubectl** | 1.27+ | Kubernetes cluster interaction |
| **AWS CLI** | 2.x | AWS authentication and operations (for AWS deployments) |
| **Azure CLI** | 2.x | Azure authentication and operations (for Azure deployments) |

### Installation Commands

<details>
<summary><b>macOS</b></summary>

```bash
# Install Terraform
brew tap hashicorp/tap
brew install hashicorp/tap/terraform

# Install Helm
brew install helm

# Install kubectl
brew install kubectl

# Install AWS CLI (for AWS deployments)
brew install awscli

# Install Azure CLI (for Azure deployments)
brew install azure-cli
```
</details>

<details>
<summary><b>Linux</b></summary>

```bash
# Install Terraform
wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt update && sudo apt install terraform

# Install Helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Install kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Install AWS CLI (for AWS deployments)
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

# Install Azure CLI (for Azure deployments)
curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash
```
</details>

<details>
<summary><b>Windows</b></summary>

```powershell
# Using Chocolatey
choco install terraform helm kubernetes-cli awscli azure-cli

# Or using Scoop
scoop install terraform helm kubectl aws azure-cli
```
</details>

## 📦 Installation

### Option 1: Download Binary (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/your-org/devplatform-cli/releases):

**Linux:**
```bash
# Download and install
VERSION="v1.0.0"
wget https://github.com/your-org/devplatform-cli/releases/download/${VERSION}/devplatform_linux_amd64.tar.gz
tar -xzf devplatform_linux_amd64.tar.gz
sudo mv devplatform /usr/local/bin/
sudo chmod +x /usr/local/bin/devplatform

# Verify installation
devplatform version
```

**macOS:**
```bash
# Download and install
VERSION="v1.0.0"
curl -LO https://github.com/your-org/devplatform-cli/releases/download/${VERSION}/devplatform_darwin_amd64.tar.gz
tar -xzf devplatform_darwin_amd64.tar.gz
sudo mv devplatform /usr/local/bin/
sudo chmod +x /usr/local/bin/devplatform

# Verify installation
devplatform version
```

**Windows:**
```powershell
# Download from GitHub Releases
# Extract the ZIP file
# Add the directory to your PATH
# Verify: devplatform version
```

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/your-org/devplatform-cli.git
cd devplatform-cli

# Build the binary
go build -o devplatform

# Install to system path
sudo mv devplatform /usr/local/bin/

# Verify installation
devplatform version
```

### Verify Installation

Check that all dependencies are installed:

```bash
devplatform version --check-deps
```

Expected output:
```
DevPlatform CLI: v1.0.0
✓ Terraform: v1.5.0 (minimum: v1.5.0)
✓ Helm: v3.12.0 (minimum: v3.0.0)
✓ kubectl: v1.27.0 (minimum: v1.27.0)
✓ AWS CLI: v2.13.0 (minimum: v2.0.0)
✓ Azure CLI: v2.50.0 (minimum: v2.0.0)

All dependencies satisfied!
```

## 🔐 Cloud Provider Setup

### AWS Setup

1. **Configure AWS Credentials:**

```bash
# Option 1: Using AWS CLI
aws configure

# Option 2: Using environment variables
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
export AWS_DEFAULT_REGION="us-east-1"
```

2. **Verify Credentials:**

```bash
aws sts get-caller-identity
```

3. **Required IAM Permissions:**

Your AWS user/role needs permissions for:
- VPC creation and management
- RDS instance creation and management
- EKS cluster access
- S3 bucket access (for Terraform state)
- DynamoDB table access (for state locking)
- Secrets Manager (for database credentials)

### Azure Setup

1. **Login to Azure:**

```bash
az login
```

2. **Set Subscription:**

```bash
# List subscriptions
az account list --output table

# Set active subscription
az account set --subscription "your-subscription-id"
```

3. **Verify Credentials:**

```bash
az account show
```

4. **Required Azure Permissions:**

Your Azure account needs permissions for:
- Virtual Network creation and management
- Azure Database for PostgreSQL creation and management
- AKS cluster access
- Azure Storage account access (for Terraform state)
- Key Vault access (for database credentials)
- Resource group management

## 📚 Documentation

All documentation is now unified in the `docs/` directory!

### Quick Links
- **[Documentation Hub](docs/README.md)** - Complete documentation navigation
- **[Interactive Docs](http://localhost:3000)** - Mintlify documentation (run `npx mintlify dev` in `docs/`)
- **[Quick Start Guide](docs/quickstart.mdx)** - 5-minute quick start
- **[Installation](docs/installation.mdx)** - Detailed installation guide

### Core Documentation
- **[Architecture](docs/reference/architecture.md)** - System architecture
- **[API Reference](docs/reference/api-reference.md)** - CLI commands
- **[Deployment Guide](docs/reference/deployment-guide.md)** - Deployment procedures
- **[Security Guide](docs/reference/security-guide.md)** - Security configuration
- **[Troubleshooting](docs/reference/troubleshooting.md)** - Common issues

### Cloud-Specific Guides
- **[AWS Deployment](docs/aws/overview.mdx)** - AWS deployment guide
- **[Azure Deployment](docs/azure/overview.mdx)** - Azure deployment guide

## 🏗️ Architecture

DevPlatform CLI uses a multi-cloud abstraction layer:

```
┌─────────────────────────────────────────────────────────┐
│                    DevPlatform CLI                       │
├─────────────────────────────────────────────────────────┤
│              Cloud Provider Abstraction                  │
├──────────────────────┬──────────────────────────────────┤
│    AWS Provider      │      Azure Provider              │
├──────────────────────┼──────────────────────────────────┤
│  • VPC               │  • VNet                          │
│  • RDS               │  • Azure Database                │
│  • EKS               │  • AKS                           │
│  • S3                │  • Azure Storage                 │
│  • IAM/IRSA          │  • Azure RBAC/Workload Identity  │
└──────────────────────┴──────────────────────────────────┘
```

## 🛠️ Technology Stack

- **Infrastructure**: Terraform
- **Container Orchestration**: Kubernetes (EKS/AKS)
- **Application Deployment**: Helm
- **Cloud Providers**: AWS, Azure
- **State Management**: S3/Azure Storage + DynamoDB/Blob Lease

## 📖 Documentation Structure

```
docs/                                  # All documentation in one place!
├── README.md                          # Documentation hub
├── mint.json                          # Mintlify configuration
├── introduction.mdx                   # Landing page
├── quickstart.mdx                     # Quick start
├── installation.mdx                   # Installation
│
├── reference/                         # Technical reference (markdown)
│   ├── api-reference.md
│   ├── architecture.md
│   ├── deployment-guide.md
│   ├── security-guide.md
│   ├── troubleshooting.md
│   └── workflows.md
│
├── project/                           # Project documentation
│   ├── README.md
│   ├── design.md
│   └── azure-support-changes.md
│
├── concepts/                          # Core concepts (interactive)
├── aws/                               # AWS guides (interactive)
├── azure/                             # Azure guides (interactive + reference)
├── security/                          # Security docs (interactive)
├── guides/                            # How-to guides (interactive)
├── advanced/                          # Advanced topics (interactive)
├── api-reference-interactive/         # Interactive API reference
└── mintlify/                          # Mintlify status docs
```

## 🚦 Commands

### create - Create Environment

Creates a complete isolated environment with infrastructure and application deployment.

```bash
devplatform create --app <name> --env <dev|staging|prod> [--provider <aws|azure>]
```

**Required Flags:**
- `--app, -a`: Application name (3-32 characters, lowercase alphanumeric and hyphens)
- `--env, -e`: Environment type (dev, staging, prod)

**Optional Flags:**
- `--provider, -p`: Cloud provider (aws or azure) (default: aws)
- `--dry-run`: Preview changes without applying them
- `--values-file`: Path to custom Helm values file
- `--config`: Path to configuration file (default: .devplatform.yaml)
- `--verbose, -v`: Enable verbose output
- `--debug`: Enable debug logging
- `--no-color`: Disable colored output
- `--timeout`: Operation timeout in minutes (default: 30)

**Examples:**

```bash
# Create dev environment on AWS (default provider)
devplatform create --app payment --env dev

# Create dev environment on Azure
devplatform create --app payment --env dev --provider azure

# Preview changes without creating (dry-run)
devplatform create --app payment --env staging --provider aws --dry-run

# Create with custom Helm values
devplatform create --app payment --env prod --provider azure --values-file custom-values.yaml

# Create with verbose output
devplatform create --app payment --env dev --provider aws --verbose
```

**What Gets Created:**

**AWS:**
- VPC with public and private subnets across multiple availability zones
- NAT gateways for private subnet internet access
- RDS PostgreSQL database instance
- EKS namespace with resource quotas
- Kubernetes service account with IRSA
- Application deployment via Helm
- Ingress for external access

**Azure:**
- VNet with public and private subnets across multiple availability zones
- NAT gateways for private subnet internet access
- Azure Database for PostgreSQL Flexible Server
- AKS namespace with resource quotas
- Kubernetes service account with Workload Identity
- Application deployment via Helm
- Ingress for external access

### status - Check Environment Status

Checks the health and status of an existing environment.

```bash
devplatform status --app <name> --env <dev|staging|prod> [--provider <aws|azure>]
```

**Required Flags:**
- `--app, -a`: Application name
- `--env, -e`: Environment type

**Optional Flags:**
- `--provider, -p`: Cloud provider (aws or azure) (default: aws)
- `--output, -o`: Output format (table, json, yaml) (default: table)
- `--watch, -w`: Watch mode (refresh every N seconds)
- `--verbose, -v`: Enable verbose output
- `--no-color`: Disable colored output

**Examples:**

```bash
# Check status on AWS
devplatform status --app payment --env dev

# Check status on Azure
devplatform status --app payment --env dev --provider azure

# JSON output
devplatform status --app payment --env dev --provider aws --output json

# Watch mode (refresh every 5 seconds)
devplatform status --app payment --env dev --provider azure --watch 5
```

### destroy - Destroy Environment

Destroys an existing environment and all associated resources.

```bash
devplatform destroy --app <name> --env <dev|staging|prod> [--provider <aws|azure>]
```

**Required Flags:**
- `--app, -a`: Application name
- `--env, -e`: Environment type

**Optional Flags:**
- `--provider, -p`: Cloud provider (aws or azure) (default: aws)
- `--confirm, -y`: Skip confirmation prompt
- `--force`: Force destruction even if errors occur
- `--keep-state`: Keep Terraform state file after destruction
- `--verbose, -v`: Enable verbose output
- `--no-color`: Disable colored output

**Examples:**

```bash
# Destroy with confirmation prompt
devplatform destroy --app payment --env dev --provider aws

# Destroy without confirmation
devplatform destroy --app payment --env dev --provider azure --confirm

# Force destroy (ignore errors)
devplatform destroy --app payment --env dev --provider aws --confirm --force
```

**Warning:** This command permanently deletes all resources. Use with caution, especially in production environments.

### version - Display Version

Displays version information for the CLI and dependencies.

```bash
devplatform version [options]
```

**Optional Flags:**
- `--short, -s`: Display only version number
- `--check-deps`: Check dependency versions

**Examples:**

```bash
# Full version info
devplatform version

# Short version
devplatform version --short

# Check dependencies
devplatform version --check-deps
```

## 🔧 Configuration

Create a `.devplatform.yaml` file in your project root to customize settings:

### Multi-Cloud Configuration Example

```yaml
# Global settings
global:
  cloud_provider: aws  # Default provider: aws or azure
  timeout: 30          # Operation timeout in minutes
  log_level: info      # Log level: debug, info, warn, error

# AWS-specific settings
aws:
  region: us-east-1
  profile: default     # AWS CLI profile to use

# Azure-specific settings
azure:
  subscription_id: "12345678-1234-1234-1234-123456789012"
  location: eastus
  tenant_id: "87654321-4321-4321-4321-210987654321"

# Environment-specific settings (applies to both AWS and Azure)
environments:
  dev:
    network_cidr: "10.0.0.0/16"
    db_instance_class: "db.t3.micro"      # AWS
    # db_instance_class: "B_Gen5_1"       # Azure equivalent
    db_allocated_storage: 20
    db_multi_az: false
    k8s_node_count: 2

  staging:
    network_cidr: "10.1.0.0/16"
    db_instance_class: "db.t3.small"      # AWS
    # db_instance_class: "B_Gen5_2"       # Azure equivalent
    db_allocated_storage: 50
    db_multi_az: false
    k8s_node_count: 3

  prod:
    network_cidr: "10.2.0.0/16"
    db_instance_class: "db.r6g.large"     # AWS
    # db_instance_class: "GP_Gen5_4"      # Azure equivalent
    db_allocated_storage: 100
    db_multi_az: true
    k8s_node_count: 5

# Terraform backend configuration
terraform:
  backend:
    # AWS S3 backend
    type: s3
    bucket: terraform-state-bucket
    dynamodb_table: terraform-locks
    region: us-east-1
    
    # Azure Storage backend (alternative)
    # type: azurerm
    # storage_account_name: tfstatestorage
    # container_name: tfstate
    # resource_group_name: terraform-state-rg

# Helm chart configuration
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

### AWS-Only Configuration

```yaml
global:
  cloud_provider: aws
  region: us-east-1

aws:
  region: us-east-1
  vpc_cidr: 10.0.0.0/16
  database:
    instance_class: db.t3.medium
    allocated_storage: 100

kubernetes:
  replicas:
    dev: 2
    staging: 3
    prod: 5
```

### Azure-Only Configuration

```yaml
global:
  cloud_provider: azure

azure:
  location: eastus
  subscription_id: "your-subscription-id"
  vnet_cidr: 10.0.0.0/16
  database:
    sku_name: GP_Gen5_2
    storage_mb: 102400

kubernetes:
  replicas:
    dev: 2
    staging: 3
    prod: 5
```

### Configuration Precedence

Settings are merged with the following precedence (highest to lowest):

1. **Command-line flags** - Highest priority
2. **Configuration file** (`.devplatform.yaml`)
3. **Environment variables**
4. **Default values** - Lowest priority

Example:
```bash
# Config file sets provider to 'aws', but flag overrides it
devplatform create --app myapp --env dev --provider azure
```

## 🔐 Authentication

### AWS

**Option 1: AWS CLI Configuration**
```bash
# Configure AWS credentials
aws configure

# Enter your credentials when prompted:
# AWS Access Key ID: your-access-key
# AWS Secret Access Key: your-secret-key
# Default region: us-east-1
# Default output format: json
```

**Option 2: Environment Variables**
```bash
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
export AWS_DEFAULT_REGION="us-east-1"
```

**Option 3: AWS Profile**
```bash
# Configure a named profile
aws configure --profile devplatform

# Use the profile
export AWS_PROFILE=devplatform
# Or specify in config file
```

**Verify AWS Credentials:**
```bash
aws sts get-caller-identity
```

### Azure

**Option 1: Interactive Login**
```bash
# Login to Azure (opens browser)
az login

# Set subscription
az account set --subscription "your-subscription-id"
```

**Option 2: Service Principal**
```bash
# Create service principal
az ad sp create-for-rbac --name devplatform-cli --role Contributor

# Set environment variables
export AZURE_CLIENT_ID="service-principal-app-id"
export AZURE_CLIENT_SECRET="service-principal-password"
export AZURE_TENANT_ID="your-tenant-id"
export AZURE_SUBSCRIPTION_ID="your-subscription-id"
```

**Option 3: Managed Identity**
```bash
# For Azure VMs or Azure Cloud Shell
# Managed identity is automatically detected
```

**Verify Azure Credentials:**
```bash
az account show
```

## 💰 Cost Estimates

Estimated monthly costs by environment and cloud provider:

### AWS Pricing

| Environment | VPC | RDS | EKS | Total/Month |
|-------------|-----|-----|-----|-------------|
| Development | $32 | $15 | $25 | **$72** |
| Staging     | $64 | $45 | $50 | **$159** |
| Production  | $96 | $180 | $150 | **$426** |

**Cost Breakdown:**
- **VPC**: NAT Gateway costs ($32/gateway/month)
- **RDS**: Database instance + storage + backups
- **EKS**: Namespace resources (shared cluster)

### Azure Pricing

| Environment | VNet | Azure DB | AKS | Total/Month |
|-------------|------|----------|-----|-------------|
| Development | $35 | $12 | $20 | **$67** |
| Staging     | $70 | $40 | $45 | **$155** |
| Production  | $105 | $165 | $140 | **$410** |

**Cost Breakdown:**
- **VNet**: NAT Gateway costs (~$35/gateway/month)
- **Azure DB**: PostgreSQL Flexible Server + storage + backups
- **AKS**: Namespace resources (shared cluster)

### Cost Optimization Tips

1. **Destroy unused environments**: Use `devplatform destroy` to remove dev/staging environments when not in use
2. **Right-size resources**: Start with dev, scale up only when needed
3. **Use spot instances**: Configure spot instances for non-production workloads
4. **Schedule shutdowns**: Automate environment shutdown during off-hours
5. **Monitor costs**: Use cloud provider cost management tools

### Cost Estimation Command

Preview costs before creating an environment:

```bash
# Dry-run shows estimated costs
devplatform create --app myapp --env dev --provider aws --dry-run
```

**Note:** Actual costs may vary based on:
- Data transfer charges
- Storage usage
- Backup retention
- Regional pricing differences
- Reserved instance discounts

## 📖 Usage Examples

### Complete Workflow - AWS

```bash
# 1. Create development environment
devplatform create --app payment --env dev --provider aws

# Output:
# ✓ Validating credentials...
# ✓ Creating VPC (10.0.0.0/16)...
# ✓ Creating RDS instance (db.t3.micro)...
# ✓ Creating EKS namespace (dev-payment-aws)...
# ✓ Deploying application...
# ✓ Verifying pods...
#
# Environment created successfully!
# Database Endpoint: payment-dev.abc123.us-east-1.rds.amazonaws.com:5432
# Ingress URL: https://payment-dev.example.com
# Namespace: dev-payment-aws

# 2. Configure kubectl access
aws eks update-kubeconfig --name shared-devplatform-cluster --region us-east-1
kubectl config set-context --current --namespace=dev-payment-aws

# 3. Check environment status
devplatform status --app payment --env dev --provider aws

# 4. Destroy when done
devplatform destroy --app payment --env dev --provider aws --confirm
```

### Complete Workflow - Azure

```bash
# 1. Create development environment
devplatform create --app payment --env dev --provider azure

# Output:
# ✓ Validating credentials...
# ✓ Creating VNet (10.10.0.0/16)...
# ✓ Creating Azure Database (B_Gen5_1)...
# ✓ Creating AKS namespace (dev-payment-azure)...
# ✓ Deploying application...
# ✓ Verifying pods...
#
# Environment created successfully!
# Database Endpoint: payment-dev.postgres.database.azure.com:5432
# Ingress URL: https://payment-dev.example.com
# Namespace: dev-payment-azure

# 2. Configure kubectl access
az aks get-credentials --name shared-devplatform-cluster --resource-group devplatform-rg
kubectl config set-context --current --namespace=dev-payment-azure

# 3. Check environment status
devplatform status --app payment --env dev --provider azure

# 4. Destroy when done
devplatform destroy --app payment --env dev --provider azure --confirm
```

### Multi-Environment Setup

```bash
# Create dev, staging, and prod environments
devplatform create --app payment --env dev --provider aws
devplatform create --app payment --env staging --provider aws
devplatform create --app payment --env prod --provider aws

# Check status of all environments
devplatform status --app payment --env dev --provider aws
devplatform status --app payment --env staging --provider aws
devplatform status --app payment --env prod --provider aws
```

### Custom Helm Values

```bash
# Create custom values file
cat > custom-values.yaml <<EOF
image:
  repository: myregistry/payment-service
  tag: v2.0.0

resources:
  requests:
    cpu: 500m
    memory: 1Gi
  limits:
    cpu: 2000m
    memory: 4Gi

env:
  - name: LOG_LEVEL
    value: debug
  - name: DATABASE_POOL_SIZE
    value: "20"
EOF

# Deploy with custom values
devplatform create --app payment --env staging --provider aws --values-file custom-values.yaml
```

### Dry-Run Mode

```bash
# Preview changes without creating resources
devplatform create --app payment --env prod --provider azure --dry-run

# Output shows:
# - Planned infrastructure changes
# - Estimated monthly costs
# - Resource specifications
# - No actual resources are created
```

### Watch Mode

```bash
# Monitor environment status in real-time
devplatform status --app payment --env dev --provider aws --watch 5

# Refreshes every 5 seconds
# Press Ctrl+C to exit
```

### JSON Output for Automation

```bash
# Get status in JSON format
devplatform status --app payment --env dev --provider aws --output json

# Use with jq for parsing
devplatform status --app payment --env dev --provider aws --output json | jq '.components.pods.ready'
```

## 🔍 Troubleshooting

### Common Issues

<details>
<summary><b>Authentication Failed</b></summary>

**AWS:**
```bash
# Check credentials
aws sts get-caller-identity

# If expired, reconfigure
aws configure
```

**Azure:**
```bash
# Check credentials
az account show

# If expired, re-login
az login
```
</details>

<details>
<summary><b>Terraform State Locked</b></summary>

```bash
# Check who holds the lock
devplatform status --app myapp --env dev --provider aws

# Wait for the lock to be released, or force unlock (use with caution)
cd terraform/environments/aws/dev
terraform force-unlock <lock-id>
```
</details>

<details>
<summary><b>Pods Not Ready</b></summary>

```bash
# Check pod logs
kubectl logs -n dev-myapp-aws -l app=myapp

# Describe pod for events
kubectl describe pod -n dev-myapp-aws -l app=myapp

# Check pod status
kubectl get pods -n dev-myapp-aws
```
</details>

<details>
<summary><b>Database Connection Failed</b></summary>

**AWS:**
```bash
# Get database password from Secrets Manager
aws secretsmanager get-secret-value --secret-id myapp-dev-db-password --query SecretString --output text

# Test connection
psql -h myapp-dev.abc123.us-east-1.rds.amazonaws.com -U postgres -d myapp
```

**Azure:**
```bash
# Get database password from Key Vault
az keyvault secret show --vault-name myapp-keyvault --name db-password --query value -o tsv

# Test connection
psql -h myapp-dev.postgres.database.azure.com -U postgres -d myapp
```
</details>

<details>
<summary><b>Insufficient Permissions</b></summary>

**AWS:**
Required IAM permissions:
- `ec2:*` (VPC management)
- `rds:*` (Database management)
- `eks:*` (Kubernetes access)
- `s3:*` (State storage)
- `dynamodb:*` (State locking)
- `secretsmanager:*` (Credentials)

**Azure:**
Required Azure RBAC roles:
- `Contributor` (Resource management)
- `Azure Kubernetes Service Cluster User Role` (AKS access)
- `Storage Blob Data Contributor` (State storage)
- `Key Vault Secrets User` (Credentials)
</details>

<details>
<summary><b>Command Not Found</b></summary>

```bash
# Ensure binary is in PATH
echo $PATH

# Add to PATH (Linux/macOS)
export PATH="/usr/local/bin:$PATH"

# Add to PATH permanently
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```
</details>

### Enable Debug Logging

```bash
# Run with debug flag for detailed logs
devplatform create --app myapp --env dev --provider aws --debug

# Logs are written to ~/.devplatform/logs/
# View latest log
tail -f ~/.devplatform/logs/devplatform.log
```

### Get Help

```bash
# General help
devplatform --help

# Command-specific help
devplatform create --help
devplatform status --help
devplatform destroy --help
```

## 🤝 Contributing

Contributions are welcome! Please read our contributing guidelines before submitting pull requests.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🔗 Links

- **Documentation**: [docs/README.md](docs/README.md)
- **Interactive Docs**: http://localhost:3000 (run `npx mintlify dev` in `docs/`)
- **GitHub**: https://github.com/your-org/devplatform-cli
- **Issues**: https://github.com/your-org/devplatform-cli/issues

## 📞 Support

- **Documentation**: See `docs/` directory
- **Troubleshooting**: See [docs/reference/troubleshooting.md](docs/reference/troubleshooting.md)
- **Community**: Join our discussions on GitHub

## 🎯 Roadmap

- [x] AWS support
- [x] Azure support
- [x] Multi-environment management (dev, staging, prod)
- [x] Comprehensive documentation
- [x] Multi-cloud configuration
- [x] Cost estimation
- [ ] GCP support
- [ ] CI/CD integration templates
- [ ] Automated cost optimization recommendations
- [ ] Automated backup and recovery
- [ ] Multi-region deployments
- [ ] Custom Terraform module support

## 🏗️ Project Structure

```
devplatform-cli/
├── cmd/                    # CLI commands
│   ├── create.go          # Create command
│   ├── status.go          # Status command
│   ├── destroy.go         # Destroy command
│   └── version.go         # Version command
├── internal/              # Internal packages
│   ├── config/           # Configuration management
│   ├── provider/         # Cloud provider abstraction
│   ├── terraform/        # Terraform wrapper
│   ├── helm/             # Helm wrapper
│   ├── aws/              # AWS utilities
│   ├── azure/            # Azure utilities
│   └── logger/           # Logging
├── terraform/            # Terraform modules
│   ├── modules/
│   │   ├── aws/         # AWS modules (VPC, RDS, EKS)
│   │   └── azure/       # Azure modules (VNet, Database, AKS)
│   └── environments/    # Environment configs
│       ├── aws/         # AWS environment configs
│       └── azure/       # Azure environment configs
├── charts/              # Helm charts
│   └── devplatform-base/  # Base application chart
├── docs/                # Documentation
└── README.md            # This file
```

## 📊 Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    DevPlatform CLI                       │
├─────────────────────────────────────────────────────────┤
│              Cloud Provider Abstraction                  │
├──────────────────────┬──────────────────────────────────┤
│    AWS Provider      │      Azure Provider              │
├──────────────────────┼──────────────────────────────────┤
│  • VPC               │  • VNet                          │
│  • RDS               │  • Azure Database                │
│  • EKS               │  • AKS                           │
│  • S3 + DynamoDB     │  • Azure Storage                 │
│  • IAM/IRSA          │  • Azure RBAC/Workload Identity  │
│  • Secrets Manager   │  • Key Vault                     │
└──────────────────────┴──────────────────────────────────┘
```

**Key Components:**

1. **CLI Layer**: User-facing commands (create, status, destroy, version)
2. **Provider Abstraction**: Unified interface for AWS and Azure
3. **Terraform Wrapper**: Infrastructure provisioning and state management
4. **Helm Wrapper**: Kubernetes application deployment
5. **Cloud Utilities**: Provider-specific authentication, pricing, and kubeconfig management

## 🔒 Security

### Best Practices

1. **Credentials**: Never commit credentials to version control
2. **State Files**: Store Terraform state remotely with encryption
3. **Secrets**: Use cloud provider secret management (Secrets Manager, Key Vault)
4. **RBAC**: Apply least-privilege access to cloud resources
5. **Network**: Use private subnets for databases and applications
6. **Encryption**: Enable encryption at rest and in transit

### Security Features

- **AWS**: VPC isolation, security groups, IRSA for pod authentication, encrypted RDS
- **Azure**: VNet isolation, NSGs, Workload Identity for pod authentication, encrypted databases
- **Kubernetes**: Resource quotas, network policies, pod security contexts
- **State**: Remote state with locking and encryption

## 📈 Performance

- **Provisioning Time**: ~3 minutes for complete environment
- **Concurrent Operations**: Supported via state locking
- **Resource Limits**: Configurable per environment
- **Scaling**: Horizontal pod autoscaling supported

## 🧪 Testing

```bash
# Run unit tests
go test ./...

# Run integration tests
go test ./... -tags=integration

# Run with coverage
go test ./... -cover
```

---

**Built with ❤️ for developers who want to focus on code, not infrastructure.**
