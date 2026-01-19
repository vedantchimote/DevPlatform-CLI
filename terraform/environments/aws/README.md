# AWS Environment Configurations

This directory contains environment-specific Terraform variable files for AWS deployments.

## Directory Structure

```
aws/
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
- **VPC CIDR**: 10.0.0.0/16
- **Availability Zones**: 2
- **NAT Gateways**: 1 (cost optimization)
- **Database**: db.t3.micro, 20GB storage, 7-day backups
- **Purpose**: Development and testing

### Staging (staging)
- **VPC CIDR**: 10.1.0.0/16
- **Availability Zones**: 2
- **NAT Gateways**: 2 (high availability)
- **Database**: db.t3.small, 50GB storage, 14-day backups
- **Purpose**: Pre-production testing and validation

### Production (prod)
- **VPC CIDR**: 10.2.0.0/16
- **Availability Zones**: 3
- **NAT Gateways**: 3 (high availability)
- **Database**: db.r6g.large, 100GB storage, 30-day backups
- **Purpose**: Production workloads

## Usage

These variable files are automatically loaded by the DevPlatform CLI based on the `--env` flag:

```bash
# Development deployment
devplatform create --app myapp --env dev --provider aws

# Staging deployment
devplatform create --app myapp --env staging --provider aws

# Production deployment
devplatform create --app myapp --env prod --provider aws
```

## Customization

You can override these defaults in two ways:

1. **Configuration File**: Create a `.devplatform.yaml` file in your project root
2. **CLI Flags**: Pass specific values via command-line flags

Example `.devplatform.yaml`:

```yaml
global:
  cloud_provider: aws
  region: us-east-1

environments:
  dev:
    terraform:
      vpc_cidr: "10.10.0.0/16"
      az_count: 2
    database:
      instance_class: "db.t3.small"
      allocated_storage: 30
```

## Resource Sizing Guidelines

### Database Instance Classes

| Environment | Instance Class | vCPUs | Memory | Use Case |
|-------------|---------------|-------|--------|----------|
| dev         | db.t3.micro   | 2     | 1 GB   | Light development workloads |
| staging     | db.t3.small   | 2     | 2 GB   | Testing with realistic data |
| prod        | db.r6g.large  | 2     | 16 GB  | Production workloads |

### Network Configuration

| Environment | AZs | NAT Gateways | Cost/Month* | Availability |
|-------------|-----|--------------|-------------|--------------|
| dev         | 2   | 1            | ~$32        | Single point of failure |
| staging     | 2   | 2            | ~$64        | High availability |
| prod        | 3   | 3            | ~$96        | Maximum availability |

*NAT Gateway costs only, excludes data transfer

## Tags

All resources are automatically tagged with:
- `App_Name`: Application name from `--app` flag
- `Env_Type`: Environment type (dev/staging/prod)
- `Cloud_Provider`: aws
- `ManagedBy`: devplatform-cli
- `Timestamp`: Creation timestamp

Additional tags can be specified in the configuration file.
