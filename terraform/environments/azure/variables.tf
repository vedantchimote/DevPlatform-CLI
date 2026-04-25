variable "subscription_id" {
  description = "Azure Subscription ID"
  type        = string
}

variable "tenant_id" {
  description = "Azure Tenant ID"
  type        = string
}

variable "app_name" {
  description = "Application name"
  type        = string

  validation {
    condition     = can(regex("^[a-z0-9-]{3,32}$", var.app_name))
    error_message = "App name must be 3-32 characters, lowercase alphanumeric and hyphens only."
  }
}

variable "env_type" {
  description = "Environment type (dev, staging, prod)"
  type        = string

  validation {
    condition     = contains(["dev", "staging", "prod"], var.env_type)
    error_message = "Environment type must be dev, staging, or prod."
  }
}

variable "location" {
  description = "Azure region"
  type        = string
  default     = "eastus"
}

variable "vnet_cidr" {
  description = "VNet CIDR block"
  type        = string
  default     = "10.10.0.0/16"
}

variable "zone_count" {
  description = "Number of availability zones"
  type        = number
  default     = 1
}

variable "nat_gateway_count" {
  description = "Number of NAT gateways"
  type        = number
  default     = 1
}

variable "enable_flow_logs" {
  description = "Enable network flow logs"
  type        = bool
  default     = false
}

variable "db_sku_name" {
  description = "PostgreSQL Flexible Server SKU"
  type        = string
  default     = "B_Standard_B1ms"
}

variable "db_storage_mb" {
  description = "PostgreSQL storage in MB"
  type        = number
  default     = 32768
}

variable "db_postgres_version" {
  description = "PostgreSQL version"
  type        = string
  default     = "15"
}

variable "db_backup_retention_days" {
  description = "Database backup retention days"
  type        = number
  default     = 7
}

variable "aks_node_count" {
  description = "Number of AKS nodes"
  type        = number
  default     = 1
}

variable "aks_vm_size" {
  description = "AKS node VM size"
  type        = string
  default     = "Standard_B2s"
}

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}
