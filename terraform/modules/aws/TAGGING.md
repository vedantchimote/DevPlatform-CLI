# AWS Resource Tagging Strategy

All AWS resources created by the DevPlatform CLI follow a consistent tagging strategy for cost tracking, resource management, and compliance.

## Standard Tags

Every resource is tagged with the following standard tags:

| Tag Key | Description | Example | Source |
|---------|-------------|---------|--------|
| `Name` | Human-readable resource name | `myapp-prod-vpc` | Auto-generated |
| `App_Name` | Application identifier | `myapp` | CLI `--app` flag |
| `Env_Type` | Environment type | `prod` | CLI `--env` flag |
| `Cloud_Provider` | Cloud provider | `aws` | CLI `--provider` flag |
| `ManagedBy` | Management tool | `devplatform-cli` | Auto-set |
| `Timestamp` | Creation timestamp | `2024-01-15T10:30:00Z` | Auto-generated |

## Additional Tags

Additional tags can be specified in two ways:

### 1. Configuration File

Add custom tags in `.devplatform.yaml`:

```yaml
global:
  cloud_provider: aws
  tags:
    CostCenter: engineering
    Team: platform
    Project: internal-tools

environments:
  prod:
    tags:
      Compliance: pci-dss
      BackupPolicy: daily
```

### 2. Environment-Specific tfvars

Add tags in `terraform/environments/aws/{env}/terraform.tfvars`:

```hcl
tags = {
  Environment = "prod"
  CostCenter  = "engineering"
  Owner       = "platform-team"
}
```

## Tag Merging

Tags are merged with the following precedence (highest to lowest):

1. **Standard tags** (Name, App_Name, Env_Type, etc.) - Always applied, cannot be overridden
2. **Configuration file tags** - From `.devplatform.yaml`
3. **Environment tfvars tags** - From `terraform.tfvars`
4. **Module default tags** - Module-specific defaults

## Resource-Specific Tags

Some resources have additional context-specific tags:

### Network Resources
- `Type`: `public` or `private` (for subnets)

### Database Resources
- `Engine`: `postgres` (for RDS instances)

### Kubernetes Resources
- `managed-by`: `devplatform-cli` (Kubernetes label)
- `app`: Application name (Kubernetes label)
- `environment`: Environment type (Kubernetes label)

## Cost Allocation Tags

For AWS Cost Explorer and billing reports, activate these cost allocation tags:

1. `App_Name` - Track costs per application
2. `Env_Type` - Track costs per environment
3. `CostCenter` - Track costs per team/department
4. `ManagedBy` - Track costs by management tool

### Activating Cost Allocation Tags

1. Go to AWS Billing Console → Cost Allocation Tags
2. Select the tags listed above
3. Click "Activate"
4. Wait 24 hours for tags to appear in Cost Explorer

## Tag Compliance

### Required Tags

All resources MUST have:
- `App_Name`
- `Env_Type`
- `ManagedBy`

### Recommended Tags

Production resources SHOULD have:
- `CostCenter`
- `Owner` or `Team`
- `Compliance` (if applicable)

### Tag Validation

The CLI validates tags before resource creation:
- Tag keys must be 1-128 characters
- Tag values must be 0-256 characters
- Maximum 50 tags per resource
- Tag keys are case-sensitive

## Examples

### Development Environment

```hcl
tags = {
  Name         = "myapp-dev-vpc"
  App_Name     = "myapp"
  Env_Type     = "dev"
  Cloud_Provider = "aws"
  ManagedBy    = "devplatform-cli"
  Timestamp    = "2024-01-15T10:30:00Z"
  CostCenter   = "engineering"
  Team         = "backend"
}
```

### Production Environment

```hcl
tags = {
  Name         = "myapp-prod-db"
  App_Name     = "myapp"
  Env_Type     = "prod"
  Cloud_Provider = "aws"
  ManagedBy    = "devplatform-cli"
  Timestamp    = "2024-01-15T10:30:00Z"
  CostCenter   = "engineering"
  Team         = "backend"
  Compliance   = "soc2"
  BackupPolicy = "daily"
  Owner        = "platform-team@company.com"
}
```

## Querying Resources by Tags

### AWS CLI

```bash
# Find all resources for an application
aws resourcegroupstaggingapi get-resources \
  --tag-filters Key=App_Name,Values=myapp

# Find all production resources
aws resourcegroupstaggingapi get-resources \
  --tag-filters Key=Env_Type,Values=prod

# Find resources managed by DevPlatform CLI
aws resourcegroupstaggingapi get-resources \
  --tag-filters Key=ManagedBy,Values=devplatform-cli
```

### Terraform

```hcl
data "aws_instances" "app" {
  filter {
    name   = "tag:App_Name"
    values = ["myapp"]
  }
  
  filter {
    name   = "tag:Env_Type"
    values = ["prod"]
  }
}
```

## Best Practices

1. **Consistency**: Use the same tag keys across all environments
2. **Lowercase**: Use lowercase for tag keys (except standard tags)
3. **Hyphens**: Use hyphens for multi-word tag keys (`cost-center`, not `CostCenter`)
4. **No PII**: Never include personally identifiable information in tags
5. **Automation**: Let the CLI manage standard tags automatically
6. **Documentation**: Document custom tags in your project README

## Troubleshooting

### Tags Not Appearing in Cost Explorer

- Wait 24 hours after activating cost allocation tags
- Ensure tags are activated in the Billing Console
- Check that resources were created after tag activation

### Tag Limit Exceeded

- AWS allows maximum 50 tags per resource
- Review and remove unnecessary custom tags
- Consolidate related tags into single values

### Inconsistent Tags

- Use the CLI's built-in tagging to ensure consistency
- Avoid manually tagging resources created by the CLI
- Use `terraform plan` to preview tag changes before applying
