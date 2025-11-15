# DevPlatform CLI Documentation

Welcome to the unified documentation for DevPlatform CLI. All documentation - both markdown reference guides and interactive Mintlify pages - are now in this single directory!

## 🎯 Quick Start

### View Interactive Documentation
```bash
# From the docs directory
npx mintlify dev

# Open http://localhost:3000
```

### Browse Documentation Files
All documentation is organized in this directory by category and format.

## 📁 Unified Documentation Structure

```
docs/
├── README.md                          # This file - documentation hub
├── mint.json                          # Mintlify configuration
├── docs.json                          # Auto-generated Mintlify config
├── favicon.svg                        # Mintlify favicon
├── logo/                              # Mintlify logos
│
├── introduction.mdx                   # Landing page (interactive)
├── quickstart.mdx                     # Quick start (interactive)
├── installation.mdx                   # Installation (interactive)
│
├── reference/                         # 📘 Technical Reference (Markdown)
│   ├── api-reference.md               # CLI command reference
│   ├── architecture.md                # System architecture
│   ├── deployment-guide.md            # Deployment procedures
│   ├── security-guide.md              # Security configuration
│   ├── troubleshooting.md             # Issue resolution
│   └── workflows.md                   # Operational workflows
│
├── project/                           # 📗 Project Documentation
│   ├── README.md                      # Project overview
│   ├── design.md                      # Design document
│   └── azure-support-changes.md       # Azure integration details
│
├── concepts/                          # 🎓 Core Concepts (Interactive)
│   ├── architecture.mdx               # Architecture concepts
│   ├── multi-cloud.mdx                # Multi-cloud design
│   ├── workflows.mdx                  # Workflow details
│   └── state-management.mdx           # State management
│
├── aws/                               # ☁️ AWS Documentation (Interactive)
│   ├── overview.mdx                   # AWS deployment overview
│   ├── authentication.mdx             # AWS authentication
│   ├── networking.mdx                 # VPC configuration
│   ├── database.mdx                   # RDS configuration
│   └── kubernetes.mdx                 # EKS integration
│
├── azure/                             # ☁️ Azure Documentation (Mixed)
│   ├── overview.mdx                   # Azure deployment overview
│   ├── authentication.mdx             # Azure authentication
│   ├── networking.mdx                 # VNet configuration
│   ├── database.mdx                   # Azure Database config
│   ├── kubernetes.mdx                 # AKS integration
│   ├── documentation-updates.md       # Azure doc updates (reference)
│   └── update-guide.md                # Azure update guide (reference)
│
├── security/                          # 🔒 Security Documentation (Interactive)
│   ├── overview.mdx                   # Security architecture
│   ├── authentication.mdx             # Authentication methods
│   ├── rbac.mdx                       # RBAC configuration
│   ├── encryption.mdx                 # Encryption details
│   └── audit-logging.mdx              # Audit and compliance
│
├── api-reference-interactive/         # 📚 API Reference (Interactive)
│   ├── introduction.mdx               # API overview
│   ├── create.mdx                     # Create command
│   ├── status.mdx                     # Status command
│   ├── destroy.mdx                    # Destroy command
│   └── version.mdx                    # Version command
│
├── guides/                            # 📖 How-To Guides (Interactive)
│   ├── first-deployment.mdx           # First deployment
│   ├── multi-environment.mdx          # Multi-environment setup
│   ├── cost-optimization.mdx          # Cost optimization
│   ├── troubleshooting.mdx            # Troubleshooting guide
│   └── migration.mdx                  # Migration guide
│
├── advanced/                          # 🚀 Advanced Topics (Interactive)
│   ├── custom-modules.mdx             # Custom Terraform modules
│   ├── helm-customization.mdx         # Helm customization
│   ├── ci-cd-integration.mdx          # CI/CD integration
│   └── disaster-recovery.mdx          # Disaster recovery
│
└── mintlify/                          # 📊 Mintlify Status
    ├── README.md                      # Mintlify overview
    ├── setup.md                       # Setup guide
    ├── status.md                      # Current status
    └── progress.md                    # Progress tracking
```

## 📊 Documentation Statistics

### Total Files: 60
- **Markdown files (.md)**: 16 files
- **Interactive files (.mdx)**: 38 files
- **Configuration files**: 2 files
- **Assets**: 4 files (logos, favicon)

### By Category
| Category | Files | Format | Purpose |
|----------|-------|--------|---------|
| Reference Guides | 6 | Markdown | Technical reference |
| Project Docs | 3 | Markdown | Project information |
| Core Concepts | 4 | MDX | Interactive learning |
| AWS Guides | 5 | MDX | AWS deployment |
| Azure Guides | 5 MDX + 2 MD | Mixed | Azure deployment |
| Security | 5 | MDX | Security configuration |
| API Reference | 5 | MDX | Command reference |
| How-To Guides | 5 | MDX | Practical guides |
| Advanced Topics | 4 | MDX | Advanced usage |
| Mintlify Status | 4 | Markdown | Documentation status |

## 🎯 Documentation by Use Case

### I want to get started quickly
1. [Quick Start](quickstart.mdx) - 5-minute guide
2. [Installation](installation.mdx) - Install the CLI
3. [First Deployment](guides/first-deployment.mdx) - Deploy your first environment

### I want to understand the architecture
1. [Architecture Reference](reference/architecture.md) - Technical architecture
2. [Architecture Concepts](concepts/architecture.mdx) - Interactive concepts
3. [Design Document](project/design.md) - Design decisions
4. [Multi-Cloud Design](concepts/multi-cloud.mdx) - Multi-cloud approach

### I want to deploy on AWS
1. [AWS Overview](aws/overview.mdx) - AWS deployment overview
2. [AWS Authentication](aws/authentication.mdx) - Configure AWS credentials
3. [AWS Networking](aws/networking.mdx) - VPC configuration
4. [AWS Database](aws/database.mdx) - RDS setup
5. [AWS Kubernetes](aws/kubernetes.mdx) - EKS integration

### I want to deploy on Azure
1. [Azure Overview](azure/overview.mdx) - Azure deployment overview
2. [Azure Authentication](azure/authentication.mdx) - Configure Azure credentials
3. [Azure Networking](azure/networking.mdx) - VNet configuration
4. [Azure Database](azure/database.mdx) - Azure Database setup
5. [Azure Kubernetes](azure/kubernetes.mdx) - AKS integration

### I want to secure my deployment
1. [Security Guide](reference/security-guide.md) - Security reference
2. [Security Overview](security/overview.mdx) - Security architecture
3. [Authentication](security/authentication.mdx) - Authentication methods
4. [RBAC](security/rbac.mdx) - Role-based access control
5. [Encryption](security/encryption.mdx) - Encryption configuration

### I'm having issues
1. [Troubleshooting Reference](reference/troubleshooting.md) - Technical troubleshooting
2. [Troubleshooting Guide](guides/troubleshooting.mdx) - Interactive troubleshooting
3. [Common Issues](reference/troubleshooting.md#common-issues) - Quick fixes

### I want to optimize costs
1. [Cost Optimization Guide](guides/cost-optimization.mdx) - Cost-saving strategies
2. [Resource Sizing](aws/overview.mdx#resource-sizing-by-environment) - Right-sizing
3. [Deployment Best Practices](reference/deployment-guide.md) - Efficient deployment

### I want to set up CI/CD
1. [CI/CD Integration](advanced/ci-cd-integration.mdx) - CI/CD setup
2. [Deployment Guide](reference/deployment-guide.md) - Automation procedures

## 📖 Documentation Formats

### Markdown (.md) - Technical Reference
- **Purpose**: Detailed technical documentation, design documents
- **Audience**: Developers, architects, contributors
- **Location**: `reference/`, `project/`, `azure/` (some), `mintlify/`
- **Viewer**: Any markdown viewer, IDE, GitHub

### MDX (.mdx) - Interactive Documentation
- **Purpose**: User-facing interactive documentation
- **Audience**: End users, developers, operators
- **Location**: Root and subdirectories
- **Viewer**: Mintlify dev server (http://localhost:3000)
- **Features**: Interactive components, tabs, diagrams, code examples

## 🚀 Accessing Documentation

### Interactive Documentation (Recommended)
```bash
# Navigate to docs directory
cd docs

# Start Mintlify dev server
npx mintlify dev

# Open in browser
# http://localhost:3000
```

### Markdown Documentation
- Use any markdown viewer or IDE
- Browse files directly in `docs/reference/` and `docs/project/`
- View on GitHub

### Quick Reference
| Documentation Type | Location | Access Method |
|-------------------|----------|---------------|
| Interactive Docs | `docs/*.mdx` | Mintlify dev server |
| Technical Reference | `docs/reference/` | Markdown viewer |
| Project Info | `docs/project/` | Markdown viewer |
| Mintlify Status | `docs/mintlify/` | Markdown viewer |

## 🔍 Finding Documentation

### By Topic

**Installation & Setup**
- [installation.mdx](installation.mdx)
- [quickstart.mdx](quickstart.mdx)
- [mintlify/setup.md](mintlify/setup.md)

**Architecture**
- [reference/architecture.md](reference/architecture.md)
- [concepts/architecture.mdx](concepts/architecture.mdx)
- [project/design.md](project/design.md)

**AWS Deployment**
- [aws/overview.mdx](aws/overview.mdx)
- [aws/authentication.mdx](aws/authentication.mdx)
- [aws/networking.mdx](aws/networking.mdx)
- [aws/database.mdx](aws/database.mdx)
- [aws/kubernetes.mdx](aws/kubernetes.mdx)

**Azure Deployment**
- [azure/overview.mdx](azure/overview.mdx)
- [azure/authentication.mdx](azure/authentication.mdx)
- [azure/networking.mdx](azure/networking.mdx)
- [azure/database.mdx](azure/database.mdx)
- [azure/kubernetes.mdx](azure/kubernetes.mdx)
- [azure/documentation-updates.md](azure/documentation-updates.md)
- [azure/update-guide.md](azure/update-guide.md)

**Security**
- [reference/security-guide.md](reference/security-guide.md)
- [security/overview.mdx](security/overview.mdx)
- [security/authentication.mdx](security/authentication.mdx)
- [security/rbac.mdx](security/rbac.mdx)
- [security/encryption.mdx](security/encryption.mdx)

**API Reference**
- [reference/api-reference.md](reference/api-reference.md) - Markdown reference
- [api-reference-interactive/](api-reference-interactive/) - Interactive reference

**Troubleshooting**
- [reference/troubleshooting.md](reference/troubleshooting.md)
- [guides/troubleshooting.mdx](guides/troubleshooting.mdx)

## 🎨 Documentation Organization Principles

### 1. Single Source of Truth
All documentation is in the `docs/` directory - no more scattered files!

### 2. Format by Purpose
- **Markdown (.md)**: Technical reference, design docs, project info
- **MDX (.mdx)**: Interactive user-facing documentation

### 3. Logical Grouping
- **reference/**: Technical reference guides
- **project/**: Project-level documentation
- **concepts/**: Core concepts (interactive)
- **aws/**, **azure/**: Cloud-specific guides
- **security/**: Security documentation
- **guides/**: How-to guides
- **advanced/**: Advanced topics

### 4. Easy Navigation
- This README provides comprehensive navigation
- Mintlify provides interactive navigation
- Clear directory structure

## 📝 Contributing to Documentation

### For Markdown Documentation
1. Edit files in `docs/reference/` or `docs/project/`
2. Use standard markdown syntax
3. Test locally with any markdown viewer

### For Interactive Documentation
1. Edit MDX files in appropriate subdirectories
2. Test with Mintlify dev server: `npx mintlify dev`
3. View at http://localhost:3000
4. Verify all interactive components work

### Adding New Documentation
1. Determine format (MD for reference, MDX for interactive)
2. Place in appropriate subdirectory
3. Update this README if adding new categories
4. Update `mint.json` if adding MDX pages
5. Test locally before committing

## 🔄 Documentation Maintenance

### Regular Updates
- Update [mintlify/status.md](mintlify/status.md) when making changes
- Keep [mintlify/progress.md](mintlify/progress.md) current
- Update this README when adding new sections

### Quality Checks
- Verify all internal links work
- Test code examples
- Check Mintlify documentation renders correctly
- Ensure consistent formatting

## 🎯 Benefits of Unified Structure

### Before
- ❌ Documentation split between `docs/` and `mintlify-docs/`
- ❌ Unclear which directory to check
- ❌ Duplicate navigation files
- ❌ Confusing for contributors

### After
- ✅ All documentation in single `docs/` directory
- ✅ Clear organization by format and purpose
- ✅ Single source of truth
- ✅ Easy to find and maintain
- ✅ Unified navigation

## 📞 Support

For documentation questions:
- **GitHub Issues**: Report documentation bugs
- **Slack Community**: Ask questions
- **Email**: support@devplatform.io

## 🔗 External Resources

- **Mintlify Documentation**: https://mintlify.com/docs
- **Terraform Documentation**: https://www.terraform.io/docs
- **Kubernetes Documentation**: https://kubernetes.io/docs
- **AWS Documentation**: https://docs.aws.amazon.com
- **Azure Documentation**: https://docs.microsoft.com/azure

---

**Documentation Version**: 2.0 (Unified Structure)
**Last Updated**: 2024
**Total Pages**: 60 (16 MD + 38 MDX + 2 JSON + 4 assets)
**Status**: ✅ Complete and Unified
