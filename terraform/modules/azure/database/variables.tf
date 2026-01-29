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

variable "resource_group_name" {
  description = "Name of the resource group"
  type        = string
}

variable "location" {
  description = "Azure region location"
  type        = string
  default     = "eastus"
}

variable "subnet_id" {
  description = "Subnet ID for database private endpoint"
  type        = string
}

variable "vnet_id" {
  description = "Virtual Network ID for private DNS zone"
  type        = string
}

variable "sku_name" {
  description = "Database SKU name"
  type        = string
  default     = ""
}

variable "storage_mb" {
  description = "Storage size in MB"
  type        = number
  default     = 0
}

variable "postgres_version" {
  description = "PostgreSQL version"
  type        = string
  default     = "15"
}

variable "backup_retention_days" {
  description = "Number of days to retain backups"
  type        = number
  default     = 7
}

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}
