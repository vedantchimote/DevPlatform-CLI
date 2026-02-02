# Azure Environment Configurations

This directory contains environment-specific Terraform variable files for Azure deployments.

## Directory Structure

```
azure/
├── dev/
│   └── terraform.tfvars       # Development environment variables
├── staging/
│   └── terraform.tfvars       # Staging environment variables
├── prod/
│   └── terraform.tfvars       # Production environment variables
└── README.md                  # This file
```

## Environment Specifications

### Development (dev)
- **VNet CIDR**: 10.10.0.0/16
- **Availability Zones**: 2
- **NAT Gateways**: 1 (cost optimization)
- **Database**: B_Standard_B1ms, 32GB storage, 7-day backups
- **Purpose**: Development and testing

### Staging (staging)
- **VNet CIDR**: 10.11.0.0/16
- **Availability Zones**: 2
- **NAT Gateways**: 2 (high availability)
- **Database**: B_Standard_B2s, 64GB storage, 14-day backups
- **Purpose**: Pre-production testing and validation

### Production (prod)
- **VNet CIDR**: 10.12.0.0/16
- **Availability Zones**: 3
- **NAT Gateways**: 3 (high availability)
- **Database**: GP_Standard_D2s_v3, 128GB storage, 30-day backups
- **Purpose**: Production workloads

## Usage

These variable files are automatically loaded by the DevPlatform CLI based on the `--env` flag:

```bash
# Development deployment
devplatform create --app myapp --env dev --provider azure

# Staging deployment
devplatform create --app myapp --env staging --provider azure

# Production deployment
devplatform create --app myapp --env prod --provider azure
```

## Customization

You can override these defaults in two ways:

1. **Configuration File**: Create a `.devplatform.yaml` file in your project root
2. **CLI Flags**: Pass specific values via command-line flags

Example `.devplatform.yaml`:

```yaml
global:
  cloud_provider: azure
  region: eastus
  resource_group: myapp-rg

environments:
  dev:
    terraform:
      vnet_cidr: "10.20.0.0/16"
      zone_count: 2
    database:
      sku_name: "B_Standard_B2s"
      storage_mb: 65536
```

## Resource Sizing Guidelines

### Database SKUs

| Environment | SKU Name           | vCores | Memory | Use Case |
|-------------|--------------------|--------|--------|----------|
| dev         | B_Standard_B1ms    | 1      | 2 GB   | Light development workloads |
| staging     | B_Standard_B2s     | 2      | 4 GB   | Testing with realistic data |
| prod        | GP_Standard_D2s_v3 | 2      | 8 GB   | Production workloads |

### Network Configuration

| Environment | Zones | NAT Gateways | Cost/Month* | Availability |
|-------------|-------|--------------|-------------|--------------|
| dev         | 2     | 1            | ~$35        | Single point of failure |
| staging     | 2     | 2            | ~$70        | High availability |
| prod        | 3     | 3            | ~$105       | Maximum availability |

*NAT Gateway costs only, excludes data transfer

### Storage Sizing

| Environment | Storage | IOPS | Throughput | Use Case |
|-------------|---------|------|------------|----------|
| dev         | 32 GB   | 120  | 25 MB/s    | Development |
| staging     | 64 GB   | 240  | 50 MB/s    | Testing |
| prod        | 128 GB  | 500  | 100 MB/s   | Production |

## Azure-Specific Features

### High Availability
- **Dev/Staging**: Single-zone deployment
- **Production**: Zone-redundant deployment with automatic failover

### Backup & Recovery
- **Geo-redundant backups**: Enabled for production
- **Point-in-time restore**: Available for all environments
- **Retention**: 7 days (dev), 14 days (staging), 30 days (prod)

### Security
- **Key Vault**: Stores database credentials and connection strings
- **Private DNS**: Database accessible only within VNet
- **Network Security Groups**: Restrict database access to VNet only
- **Workload Identity**: Azure AD integration for AKS pods

### Monitoring
- **Network Watcher**: Flow logs for all environments
- **Traffic Analytics**: 10-minute intervals
- **Log Analytics**: 7-day retention for flow logs

## Tags

All resources are automatically tagged with:
- `App_Name`: Application name from `--app` flag
- `Env_Type`: Environment type (dev/staging/prod)
- `Cloud_Provider`: azure
- `ManagedBy`: devplatform-cli
- `Timestamp`: Creation timestamp

Additional tags can be specified in the configuration file.

## Cost Optimization Tips

### Development
- Use B-series (burstable) SKUs for cost savings
- Single NAT Gateway reduces costs by ~50%
- Single-zone deployment minimizes data transfer costs

### Staging
- Balance between cost and production-like environment
- Two NAT Gateways for HA testing
- B-series SKUs still cost-effective

### Production
- GP-series (general purpose) SKUs for consistent performance
- Zone-redundant deployment for 99.99% SLA
- Three NAT Gateways for maximum availability
- Geo-redundant backups for disaster recovery

## Azure Regions

Recommended regions for deployment:
- **East US** (eastus): Default, good for US-based workloads
- **West Europe** (westeurope): Good for EU-based workloads
- **Southeast Asia** (southeastasia): Good for APAC-based workloads

Choose a region close to your users for optimal latency.

## Troubleshooting

### NAT Gateway Quota
If you hit NAT Gateway quota limits:
- Request quota increase in Azure Portal
- Or reduce `nat_gateway_count` for non-prod environments

### Database SKU Availability
Some SKUs may not be available in all regions:
- Check SKU availability: `az postgres flexible-server list-skus --location <region>`
- Choose alternative SKU from same tier

### Zone Availability
Not all regions support 3 availability zones:
- Check zone support: `az account list-locations --query "[?name=='<region>'].availabilityZoneMappings"`
- Adjust `zone_count` accordingly

## Migration from AWS

Key differences when migrating from AWS:
- VPC → VNet (Virtual Network)
- RDS → PostgreSQL Flexible Server
- EKS → AKS (Azure Kubernetes Service)
- Secrets Manager → Key Vault
- IRSA → Workload Identity
- Security Groups → Network Security Groups

See the main documentation for detailed migration guide.
