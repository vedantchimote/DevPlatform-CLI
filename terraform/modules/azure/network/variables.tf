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

variable "vnet_cidr" {
  description = "CIDR block for VNet"
  type        = string
  default     = "10.0.0.0/16"
}

variable "zone_count" {
  description = "Number of availability zones to use"
  type        = number
  default     = 2

  validation {
    condition     = var.zone_count >= 1 && var.zone_count <= 3
    error_message = "Zone count must be between 1 and 3."
  }
}

variable "nat_gateway_count" {
  description = "Number of NAT gateways (1 for dev, 2+ for staging/prod)"
  type        = number
  default     = 1
}

variable "enable_flow_logs" {
  description = "Enable Network Watcher flow logs"
  type        = bool
  default     = true
}

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}
