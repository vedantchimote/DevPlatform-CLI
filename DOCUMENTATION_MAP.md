# DevPlatform CLI - Unified Documentation Map

This file provides a complete map of all documentation in the project. All documentation is now unified in the `docs/` directory!

## 📂 Unified Documentation Structure

```
DevPlatform-CLI/
│
├── README.md                                    # Main project README
├── DOCUMENTATION_MAP.md                         # This file - complete documentation map
│
├── docs/                                        # 📚 ALL DOCUMENTATION (UNIFIED!)
│   ├── README.md                                # Documentation hub and navigation
│   ├── mint.json                                # Mintlify configuration
│   ├── docs.json                                # Auto-generated Mintlify config
│   ├── favicon.svg                              # Mintlify favicon
│   │
│   ├── introduction.mdx                         # Landing page (interactive)
│   ├── quickstart.mdx                           # Quick start guide (interactive)
│   ├── installation.mdx                         # Installation guide (interactive)
│   │
│   ├── logo/                                    # Mintlify assets
│   │   ├── light.svg
│   │   └── dark.svg
│   │
│   ├── reference/                               # 📘 Technical Reference (Markdown)
│   │   ├── api-reference.md                     # CLI command reference
│   │   ├── architecture.md                      # System architecture
│   │   ├── deployment-guide.md                  # Deployment procedures
│   │   ├── security-guide.md                    # Security configuration
│   │   ├── troubleshooting.md                   # Issue resolution
│   │   └── workflows.md                         # Operational workflows
│   │
│   ├── project/                                 # 📗 Project Documentation
│   │   ├── README.md                            # Project overview
│   │   ├── design.md                            # Design document
│   │   └── azure-support-changes.md             # Azure integration details
│   │
│   ├── concepts/                                # 🎓 Core Concepts (Interactive - 4 pages)
│   │   ├── architecture.mdx
│   │   ├── multi-cloud.mdx
│   │   ├── workflows.mdx
│   │   └── state-management.mdx
│   │
│   ├── aws/                                     # ☁️ AWS Documentation (Interactive - 5 pages)
│   │   ├── overview.mdx
│   │   ├── authentication.mdx
│   │   ├── networking.mdx
│   │   ├── database.mdx
│   │   └── kubernetes.mdx
│   │
│   ├── azure/                                   # ☁️ Azure Documentation (Mixed - 5 MDX + 2 MD)
│   │   ├── overview.mdx                         # Interactive
│   │   ├── authentication.mdx                   # Interactive
│   │   ├── networking.mdx                       # Interactive
│   │   ├── database.mdx                         # Interactive
│   │   ├── kubernetes.mdx                       # Interactive
│   │   ├── documentation-updates.md             # Reference
│   │   └── update-guide.md                      # Reference
│   │
│   ├── security/                                # 🔒 Security Documentation (Interactive - 5 pages)
│   │   ├── overview.mdx
│   │   ├── authentication.mdx
│   │   ├── rbac.mdx
│   │   ├── encryption.mdx
│   │   └── audit-logging.mdx
│   │
│   ├── api-reference-interactive/               # 📚 API Reference (Interactive - 5 pages)
│   │   ├── introduction.mdx
│   │   ├── create.mdx
│   │   ├── status.mdx
│   │   ├── destroy.mdx
│   │   └── version.mdx
│   │
│   ├── guides/                                  # 📖 How-To Guides (Interactive - 5 pages)
│   │   ├── first-deployment.mdx
│   │   ├── multi-environment.mdx
│   │   ├── cost-optimization.mdx
│   │   ├── troubleshooting.mdx
│   │   └── migration.mdx
│   │
│   ├── advanced/                                # 🚀 Advanced Topics (Interactive - 4 pages)
│   │   ├── custom-modules.mdx
│   │   ├── helm-customization.mdx
│   │   ├── ci-cd-integration.mdx
│   │   └── disaster-recovery.mdx
│   │
│   └── mintlify/                                # 📊 Mintlify Status (Markdown - 4 files)
│       ├── README.md                            # Mintlify overview
│       ├── setup.md                             # Setup guide
│       ├── status.md                            # Current status
│       └── progress.md                          # Progress tracking
│
└── .kiro/specs/devplatform-cli/                 # 📋 Specifications
    ├── .config.kiro                             # Spec configuration
    ├── requirements.md                          # Functional requirements
    ├── design.md                                # Design specifications
    └── tasks.md                                 # Implementation tasks
```

## 📊 Documentation Statistics

### Total Documentation Files: 60
- **Markdown files (.md)**: 16 files
- **MDX files (.mdx)**: 38 files
- **Configuration files**: 2 files (mint.json, docs.json)
- **Assets**: 4 files (2 logos, 1 favicon, 1 README)

### By Category
| Category | Files | Format | Location |
|----------|-------|--------|----------|
| Getting Started | 3 | MDX | `docs/` root |
| Reference Guides | 6 | Markdown | `docs/reference/` |
| Project Docs | 3 | Markdown | `docs/project/` |
| Core Concepts | 4 | MDX | `docs/concepts/` |
| AWS Guides | 5 | MDX | `docs/aws/` |
| Azure Guides | 7 | 5 MDX + 2 MD | `docs/azure/` |
| Security | 5 | MDX | `docs/security/` |
| API Reference | 5 | MDX | `docs/api-reference-interactive/` |
| How-To Guides | 5 | MDX | `docs/guides/` |
| Advanced Topics | 4 | MDX | `docs/advanced/` |
| Mintlify Status | 4 | Markdown | `docs/mintlify/` |
| Specifications | 3 | Markdown | `.kiro/specs/devplatform-cli/` |

### Interactive Documentation Breakdown
| Section | Pages | Status |
|---------|-------|--------|
| Getting Started | 3 | ✅ Complete |
| Core Concepts | 4 | ✅ Complete |
| AWS Deployment | 5 | ✅ Complete |
| Azure Deployment | 5 | ✅ Complete |
| Security | 5 | ✅ Complete |
| API Reference | 5 | ✅ Complete |
| Guides | 5 | ✅ Complete |
| Advanced Topics | 4 | ✅ Complete |
| **Total** | **38** | **✅ 100%** |

## 🎯 Documentation by Audience

### For End Users
**Start here**: `README.md` → `docs/quickstart.mdx`

1. **Getting Started**
   - `README.md` - Project overview
   - `docs/installation.mdx` - Installation
   - `docs/quickstart.mdx` - Quick start
   - `docs/guides/first-deployment.mdx` - First deployment

2. **Cloud-Specific Guides**
   - `docs/aws/overview.mdx` - AWS deployment
   - `docs/azure/overview.mdx` - Azure deployment

3. **How-To Guides**
   - `docs/guides/multi-environment.mdx` - Multi-environment setup
   - `docs/guides/cost-optimization.mdx` - Cost optimization
   - `docs/guides/troubleshooting.mdx` - Troubleshooting

### For Developers
**Start here**: `docs/project/README.md` → `docs/reference/architecture.md`

1. **Architecture & Design**
   - `docs/project/design.md` - Design document
   - `docs/reference/architecture.md` - System architecture
   - `docs/concepts/architecture.mdx` - Architecture concepts
   - `docs/concepts/multi-cloud.mdx` - Multi-cloud design

2. **Implementation Details**
   - `.kiro/specs/devplatform-cli/requirements.md` - Requirements
   - `.kiro/specs/devplatform-cli/design.md` - Design specs
   - `.kiro/specs/devplatform-cli/tasks.md` - Implementation tasks
   - `docs/project/azure-support-changes.md` - Azure integration

3. **Workflows & Processes**
   - `docs/reference/workflows.md` - Operational workflows
   - `docs/concepts/workflows.mdx` - Workflow details

### For Operations/DevOps
**Start here**: `docs/reference/deployment-guide.md` → `docs/reference/security-guide.md`

1. **Deployment**
   - `docs/reference/deployment-guide.md` - Deployment procedures
   - `docs/guides/first-deployment.mdx` - Step-by-step guide
   - `docs/guides/multi-environment.mdx` - Environment management

2. **Security**
   - `docs/reference/security-guide.md` - Security guide
   - `docs/security/overview.mdx` - Security overview
   - `docs/security/authentication.mdx` - Authentication
   - `docs/security/rbac.mdx` - RBAC configuration

3. **Operations**
   - `docs/reference/troubleshooting.md` - Troubleshooting
   - `docs/guides/troubleshooting.mdx` - Common issues
   - `docs/advanced/disaster-recovery.mdx` - DR procedures

### For Documentation Contributors
**Start here**: `docs/README.md` → `docs/mintlify/setup.md`

1. **Documentation Setup**
   - `docs/README.md` - Documentation overview
   - `docs/mintlify/setup.md` - Mintlify setup
   - `docs/mintlify/README.md` - Mintlify guide

2. **Status & Progress**
   - `docs/mintlify/status.md` - Current status
   - `docs/mintlify/progress.md` - Progress tracking
   - `DOCUMENTATION_MAP.md` - This file

## 🔍 Finding Documentation

### By Topic

#### Installation & Setup
- `docs/installation.mdx`
- `docs/quickstart.mdx`
- `docs/mintlify/setup.md`

#### Architecture
- `docs/reference/architecture.md`
- `docs/project/design.md`
- `docs/concepts/architecture.mdx`
- `.kiro/specs/devplatform-cli/design.md`

#### AWS Deployment
- `docs/aws/overview.mdx`
- `docs/aws/authentication.mdx`
- `docs/aws/networking.mdx`
- `docs/aws/database.mdx`
- `docs/aws/kubernetes.mdx`

#### Azure Deployment
- `docs/azure/overview.mdx`
- `docs/azure/authentication.mdx`
- `docs/azure/networking.mdx`
- `docs/azure/database.mdx`
- `docs/azure/kubernetes.mdx`
- `docs/azure/documentation-updates.md`
- `docs/azure/update-guide.md`
- `docs/project/azure-support-changes.md`

#### Security
- `docs/reference/security-guide.md`
- `docs/security/overview.mdx`
- `docs/security/authentication.mdx`
- `docs/security/rbac.mdx`
- `docs/security/encryption.mdx`
- `docs/security/audit-logging.mdx`

#### API Reference
- `docs/reference/api-reference.md` (Markdown)
- `docs/api-reference-interactive/introduction.mdx` (Interactive)
- `docs/api-reference-interactive/create.mdx`
- `docs/api-reference-interactive/status.mdx`
- `docs/api-reference-interactive/destroy.mdx`
- `docs/api-reference-interactive/version.mdx`

#### Troubleshooting
- `docs/reference/troubleshooting.md`
- `docs/guides/troubleshooting.mdx`

#### Workflows
- `docs/reference/workflows.md`
- `docs/concepts/workflows.mdx`

#### Advanced Topics
- `docs/advanced/custom-modules.mdx`
- `docs/advanced/helm-customization.mdx`
- `docs/advanced/ci-cd-integration.mdx`
- `docs/advanced/disaster-recovery.mdx`

## 🚀 Accessing Documentation

### Interactive Documentation (Recommended)
```bash
# Navigate to docs directory
cd docs

# Start Mintlify dev server
npx mintlify dev

# Open http://localhost:3000
```

### Markdown Documentation
- Use any markdown viewer or IDE
- Browse files in `docs/reference/` and `docs/project/`
- View on GitHub

### File Locations
- **Root README**: `./README.md`
- **Documentation Hub**: `./docs/README.md`
- **All Documentation**: `./docs/` directory
- **Specifications**: `./.kiro/specs/devplatform-cli/`

## 📝 Documentation Formats

### Markdown (.md)
- **Purpose**: Technical documentation, guides, references
- **Location**: `docs/reference/`, `docs/project/`, `docs/mintlify/`, `docs/azure/` (some)
- **Viewers**: Any markdown viewer, IDE, GitHub

### MDX (.mdx)
- **Purpose**: Interactive documentation with components
- **Location**: `docs/` root and subdirectories
- **Viewer**: Mintlify dev server (http://localhost:3000)
- **Features**: Interactive components, tabs, diagrams, code examples

### Configuration Files
- **mint.json**: Mintlify configuration
- **docs.json**: Auto-generated Mintlify config
- **.config.kiro**: Spec configuration

## 🎨 Organization Principles

### 1. Single Source of Truth
All documentation is in the `docs/` directory - no more scattered files or separate directories!

### 2. Format by Purpose
- **Markdown (.md)**: Technical reference, design docs, project info
- **MDX (.mdx)**: Interactive user-facing documentation

### 3. Logical Grouping
Files are grouped by purpose and audience:
- **reference/**: Technical reference guides
- **project/**: Project-level information
- **concepts/**, **aws/**, **azure/**, **security/**: Topic-specific guides
- **guides/**: How-to guides
- **advanced/**: Advanced topics
- **mintlify/**: Documentation platform status

### 4. Easy Navigation
- Root README points to all documentation
- `docs/README.md` provides comprehensive navigation
- `DOCUMENTATION_MAP.md` shows complete structure
- Mintlify provides interactive navigation

## 🔄 Documentation Maintenance

### Regular Updates
- Update `docs/mintlify/status.md` when making changes
- Update `docs/mintlify/progress.md` for tracking
- Keep `DOCUMENTATION_MAP.md` (this file) current
- Update `docs/README.md` when adding new sections

### Version Control
- All documentation is version controlled
- Use meaningful commit messages
- Review changes before committing

### Quality Checks
- Verify all internal links work
- Check code examples are correct
- Ensure consistent formatting
- Test Mintlify documentation locally

## ✨ Benefits of Unified Structure

### Before
- ❌ Documentation split between `docs/` and `mintlify-docs/`
- ❌ Unclear which directory to check
- ❌ Duplicate navigation files
- ❌ Confusing for contributors
- ❌ Separate dev servers needed

### After
- ✅ All documentation in single `docs/` directory
- ✅ Clear organization by format and purpose
- ✅ Single source of truth
- ✅ Easy to find and maintain
- ✅ Unified navigation
- ✅ Single dev server for all interactive docs

## 📞 Support

For documentation questions or issues:
- **GitHub Issues**: Report documentation bugs
- **Slack Community**: Ask questions
- **Email**: support@devplatform.io

---

**Documentation Version**: 2.0 (Unified Structure)
**Last Updated**: 2024
**Total Pages**: 60 (16 MD + 38 MDX + 2 JSON + 4 assets)
**Status**: ✅ Complete and Unified
**Location**: All in `docs/` directory
