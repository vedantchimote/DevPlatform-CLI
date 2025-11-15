# DevPlatform CLI

A powerful command-line tool for deploying and managing development environments on AWS and Azure with automated infrastructure provisioning and Kubernetes orchestration.

## 🚀 Quick Start

```bash
# Install DevPlatform CLI
npm install -g devplatform-cli

# Create a development environment on AWS
devplatform create --app myapp --env dev --provider aws

# Create a development environment on Azure
devplatform create --app myapp --env dev --provider azure

# Check environment status
devplatform status --app myapp --env dev

# Destroy environment when done
devplatform destroy --app myapp --env dev
```

## ✨ Features

- **Multi-Cloud Support**: Deploy to AWS or Azure with a single command
- **Infrastructure as Code**: Automated Terraform provisioning
- **Kubernetes Integration**: Deploy applications to EKS (AWS) or AKS (Azure)
- **Environment Management**: Separate dev, staging, and production environments
- **State Management**: Secure state storage with locking
- **Cost Optimization**: Right-sized resources per environment

## 📚 Documentation

All documentation is now unified in the `docs/` directory!

### Quick Links
- **[Documentation Hub](docs/README.md)** - Complete documentation navigation
- **[Interactive Docs](http://localhost:3000)** - Mintlify documentation (run `npx mintlify dev` in `docs/`)
- **[Quick Start Guide](docs/quickstart.mdx)** - 5-minute quick start
- **[Installation](docs/installation.mdx)** - Installation guide

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

### Create Environment
```bash
devplatform create --app <name> --env <dev|staging|prod> [--provider <aws|azure>]
```

### Check Status
```bash
devplatform status --app <name> --env <dev|staging|prod> [--provider <aws|azure>]
```

### Destroy Environment
```bash
devplatform destroy --app <name> --env <dev|staging|prod> [--provider <aws|azure>]
```

### Version
```bash
devplatform version
```

## 🔧 Configuration

Create a `.devplatform.yaml` file in your project root:

```yaml
provider: aws  # or azure

aws:
  region: us-east-1
  vpc_cidr: 10.0.0.0/16
  database:
    instance_class: db.t3.medium
    allocated_storage: 100

azure:
  location: eastus
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

## 🔐 Authentication

### AWS
```bash
# Configure AWS credentials
aws configure

# Or use environment variables
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_DEFAULT_REGION=us-east-1
```

### Azure
```bash
# Login to Azure
az login

# Set subscription
az account set --subscription "My Subscription"
```

## 💰 Cost Estimates

| Environment | AWS (Monthly) | Azure (Monthly) |
|-------------|---------------|-----------------|
| Development | $50-75        | $45-70          |
| Staging     | $200-300      | $180-280        |
| Production  | $800-1200     | $750-1100       |

## 📱 Interactive Documentation

Run the Mintlify dev server to access interactive documentation:

```bash
cd docs
npx mintlify dev
```

Then open http://localhost:3000 in your browser.

## 🤝 Contributing

Contributions are welcome! Please read our contributing guidelines before submitting pull requests.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🔗 Links

- **Documentation**: [docs/README.md](docs/README.md)
- **Interactive Docs**: http://localhost:3000 (when dev server is running)
- **GitHub**: https://github.com/your-org/devplatform-cli
- **Issues**: https://github.com/your-org/devplatform-cli/issues
- **Slack Community**: https://slack.devplatform.io

## 📞 Support

- **Email**: support@devplatform.io
- **Documentation**: See `docs/` directory
- **Troubleshooting**: See [docs/reference/troubleshooting.md](docs/reference/troubleshooting.md)

## 🎯 Roadmap

- [x] AWS support
- [x] Azure support
- [x] Multi-environment management
- [x] Comprehensive documentation
- [x] Unified documentation structure
- [ ] GCP support
- [ ] CI/CD integration templates
- [ ] Cost optimization recommendations
- [ ] Automated backup and recovery

---

**Built with ❤️ for developers who want to focus on code, not infrastructure.**
