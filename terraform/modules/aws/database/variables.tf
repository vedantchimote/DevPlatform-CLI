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

variable "vpc_id" {
  description = "VPC ID where database will be deployed"
  type        = string
}

variable "subnet_ids" {
  description = "List of subnet IDs for database subnet group"
  type        = list(string)
}

variable "security_group_ids" {
  description = "List of security group IDs to allow database access"
  type        = list(string)
  default     = []
}

variable "instance_class" {
  description = "RDS instance class"
  type        = string
  default     = ""
}

variable "allocated_storage" {
  description = "Allocated storage in GB"
  type        = number
  default     = 0
}

variable "engine_version" {
  description = "PostgreSQL engine version"
  type        = string
  default     = "15.4"
}

variable "backup_retention_period" {
  description = "Number of days to retain backups"
  type        = number
  default     = 7
}

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}
