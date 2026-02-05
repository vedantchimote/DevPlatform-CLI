# Azure Resource Tagging Strategy

All Azure resources created by the DevPlatform CLI follow a consistent tagging strategy for cost tracking, resource management, and compliance.

## Standard Tags

Every resource is tagged with the following standard tags:

| Tag Key | Description | Example | Source |
|---------|-------------|---------|--------|
| `Name` | Human-readable resource name | `myapp-prod-vnet` | Auto-generated |
| `App_Name` | Application identifier | `myapp` | CLI `--app` flag |
| `Env_Type` | Environment type | `prod` | CLI `--env` flag |
| `Cloud_Provider` | Cloud provider | `azure` | CLI `--provider` flag |
| `ManagedBy` | Management tool | `devplatform-cli` | Auto-set |
| `Timestamp` | Creation timestamp | `2024-01-15T10:30:00Z` | Auto-generated |

## Additional Tags

Additional tags can be specified in two ways:

### 1. Configuration File

Add custom tags in `.devplatform.yaml`:

```yaml
global:
  cloud_provider: azure
  resource_group: myapp-rg
  location: eastus
  tags:
    CostCenter: engineering
    Team: platform
    Project: internal-tools

environments:
  prod:
    tags:
      Compliance: iso27001
      BackupPolicy: daily
```

### 2. Environment-Specific tfvars

Add tags in `terraform/environments/azure/{env}/terraform.tfvars`:

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
- No additional tags (uses standard tags only)

### Database Resources
- No additional tags (uses standard tags only)

### Kubernetes Resources
- `managed-by`: `devplatform-cli` (Kubernetes label)
- `app`: Application name (Kubernetes label)
- `environment`: Environment type (Kubernetes label)
- `cloud-provider`: `azure` (Kubernetes label)

## Cost Management Tags

For Azure Cost Management and billing reports, use these tags:

1. `App_Name` - Track costs per application
2. `Env_Type` - Track costs per environment
3. `CostCenter` - Track costs per team/department
4. `ManagedBy` - Track costs by management tool

### Viewing Costs by Tags

1. Go to Azure Portal → Cost Management + Billing
2. Select "Cost analysis"
3. Add filter by tag (e.g., `App_Name = myapp`)
4. Group by tag to see cost breakdown

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
- Tag keys must be 1-512 characters
- Tag values must be 0-256 characters
- Maximum 50 tags per resource
- Tag keys are case-insensitive (Azure converts to lowercase)

## Examples

### Development Environment

```hcl
tags = {
  Name           = "myapp-dev-vnet"
  App_Name       = "myapp"
  Env_Type       = "dev"
  Cloud_Provider = "azure"
  ManagedBy      = "devplatform-cli"
  Timestamp      = "2024-01-15T10:30:00Z"
  CostCenter     = "engineering"
  Team           = "backend"
}
```

### Production Environment

```hcl
tags = {
  Name           = "myapp-prod-db"
  App_Name       = "myapp"
  Env_Type       = "prod"
  Cloud_Provider = "azure"
  ManagedBy      = "devplatform-cli"
  Timestamp      = "2024-01-15T10:30:00Z"
  CostCenter     = "engineering"
  Team           = "backend"
  Compliance     = "iso27001"
  BackupPolicy   = "daily"
  Owner          = "platform-team@company.com"
}
```

## Querying Resources by Tags

### Azure CLI

```bash
# Find all resources for an application
az resource list --tag App_Name=myapp

# Find all production resources
az resource list --tag Env_Type=prod

# Find resources managed by DevPlatform CLI
az resource list --tag ManagedBy=devplatform-cli

# Get resources with multiple tag filters
az resource list --tag App_Name=myapp --tag Env_Type=prod
```

### Azure PowerShell

```powershell
# Find all resources for an application
Get-AzResource -TagName "App_Name" -TagValue "myapp"

# Find all production resources
Get-AzResource -TagName "Env_Type" -TagValue "prod"

# Get resources with multiple tags
Get-AzResource -Tag @{App_Name="myapp"; Env_Type="prod"}
```

### Terraform

```hcl
data "azurerm_resources" "app" {
  type = "Microsoft.Network/virtualNetworks"
  
  required_tags = {
    App_Name = "myapp"
    Env_Type = "prod"
  }
}
```

## Azure Policy Integration

### Enforce Tagging with Azure Policy

Create a policy to require specific tags:

```json
{
  "mode": "Indexed",
  "policyRule": {
    "if": {
      "allOf": [
        {
          "field": "type",
          "equals": "Microsoft.Resources/subscriptions/resourceGroups"
        },
        {
          "anyOf": [
            {
              "field": "tags['App_Name']",
              "exists": "false"
            },
            {
              "field": "tags['Env_Type']",
              "exists": "false"
            },
            {
              "field": "tags['ManagedBy']",
              "exists": "false"
            }
          ]
        }
      ]
    },
    "then": {
      "effect": "deny"
    }
  }
}
```

### Inherit Tags from Resource Group

Use Azure Policy to automatically inherit tags from resource group:

```json
{
  "mode": "Indexed",
  "policyRule": {
    "if": {
      "field": "tags['CostCenter']",
      "exists": "false"
    },
    "then": {
      "effect": "modify",
      "details": {
        "roleDefinitionIds": [
          "/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c"
        ],
        "operations": [
          {
            "operation": "addOrReplace",
            "field": "tags['CostCenter']",
            "value": "[resourceGroup().tags['CostCenter']]"
          }
        ]
      }
    }
  }
}
```

## Best Practices

1. **Consistency**: Use the same tag keys across all environments
2. **Lowercase**: Azure converts tag keys to lowercase automatically
3. **Hyphens**: Use hyphens for multi-word tag keys (`cost-center`, not `CostCenter`)
4. **No PII**: Never include personally identifiable information in tags
5. **Automation**: Let the CLI manage standard tags automatically
6. **Documentation**: Document custom tags in your project README
7. **Resource Groups**: Tag resource groups with the same tags as resources

## Azure-Specific Considerations

### Tag Inheritance
- Tags are NOT automatically inherited from resource groups
- Use Azure Policy to enforce tag inheritance if needed
- DevPlatform CLI applies tags to all resources directly

### Tag Limits
- Maximum 50 tags per resource
- Tag keys: 1-512 characters
- Tag values: 0-256 characters
- Some resource types have lower limits

### Case Sensitivity
- Tag keys are case-insensitive (converted to lowercase)
- Tag values are case-sensitive
- Use consistent casing for values

### Reserved Tags
Avoid using these reserved prefixes:
- `microsoft.*`
- `azure.*`
- `windows.*`

## Troubleshooting

### Tags Not Appearing in Cost Management

- Wait 24-48 hours for tags to appear in cost reports
- Ensure resources were created after tagging policy
- Check that tags are applied correctly: `az resource show --ids <resource-id>`

### Tag Limit Exceeded

- Azure allows maximum 50 tags per resource
- Review and remove unnecessary custom tags
- Consolidate related tags into single values

### Inconsistent Tags

- Use the CLI's built-in tagging to ensure consistency
- Avoid manually tagging resources created by the CLI
- Use `terraform plan` to preview tag changes before applying

### Tag Policy Conflicts

- Check for conflicting Azure Policies
- Ensure DevPlatform CLI tags don't conflict with org policies
- Use `az policy state list` to check policy compliance

## Comparison with AWS Tagging

Key differences from AWS:

| Feature | AWS | Azure |
|---------|-----|-------|
| Tag key length | 128 chars | 512 chars |
| Tag value length | 256 chars | 256 chars |
| Max tags per resource | 50 | 50 |
| Case sensitivity | Case-sensitive | Keys: case-insensitive, Values: case-sensitive |
| Cost allocation | Activate in billing console | Automatic in Cost Management |
| Tag inheritance | Not supported | Via Azure Policy |
| Reserved prefixes | `aws:*` | `microsoft.*`, `azure.*`, `windows.*` |

## Migration from AWS

When migrating from AWS:
- Tag keys will be converted to lowercase in Azure
- Adjust any automation that relies on case-sensitive tag keys
- Update cost allocation reports to use Azure Cost Management
- Review and update any tag-based access policies

## Additional Resources

- [Azure Tagging Best Practices](https://docs.microsoft.com/azure/cloud-adoption-framework/ready/azure-best-practices/naming-and-tagging)
- [Azure Policy for Tagging](https://docs.microsoft.com/azure/governance/policy/samples/built-in-policies#tags)
- [Azure Cost Management](https://docs.microsoft.com/azure/cost-management-billing/)
- [Resource Naming Conventions](https://docs.microsoft.com/azure/cloud-adoption-framework/ready/azure-best-practices/resource-naming)
