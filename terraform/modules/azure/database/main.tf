terraform {
  required_version = ">= 1.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

# Environment-specific defaults
locals {
  sku_defaults = {
    dev     = "B_Standard_B1ms"
    staging = "B_Standard_B2s"
    prod    = "GP_Standard_D2s_v3"
  }

  storage_defaults = {
    dev     = 32768  # 32 GB
    staging = 65536  # 64 GB
    prod    = 131072 # 128 GB
  }

  sku_name              = var.sku_name != "" ? var.sku_name : local.sku_defaults[var.env_type]
  storage_mb            = var.storage_mb > 0 ? var.storage_mb : local.storage_defaults[var.env_type]
  zone_redundant        = var.env_type == "prod" ? true : false
  backup_retention_days = var.env_type == "prod" ? 30 : var.backup_retention_days

  db_name = replace("${var.app_name}_${var.env_type}", "-", "_")
}

# Generate random password for database
resource "random_password" "db_password" {
  length  = 32
  special = true
  # Exclude characters that might cause issues
  override_special = "!#$%&*()-_=+[]{}:?"
}

# Create Key Vault for storing database credentials
resource "azurerm_key_vault" "db" {
  name                       = "${var.app_name}-${var.env_type}-kv"
  location                   = var.location
  resource_group_name        = var.resource_group_name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  purge_protection_enabled   = var.env_type == "prod" ? true : false

  # Access policy for the Terraform executor (current user/SP)
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get", "List", "Set", "Delete", "Purge", "Recover"
    ]
  }

  network_acls {
    default_action = "Allow"
    bypass         = "AzureServices"
  }

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-kv"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Get current Azure client configuration
data "azurerm_client_config" "current" {}

# Store database password in Key Vault
resource "azurerm_key_vault_secret" "db_password" {
  name         = "db-password"
  value        = random_password.db_password.result
  key_vault_id = azurerm_key_vault.db.id

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-db-password"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Store database connection string in Key Vault
resource "azurerm_key_vault_secret" "db_connection_string" {
  name  = "db-connection-string"
  value = "postgresql://dbadmin:${random_password.db_password.result}@${azurerm_postgresql_flexible_server.main.fqdn}:5432/${local.db_name}?sslmode=require"
  key_vault_id = azurerm_key_vault.db.id

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-db-connection"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Create PostgreSQL Flexible Server
resource "azurerm_postgresql_flexible_server" "main" {
  name                = "${var.app_name}-${var.env_type}-db"
  location            = var.location
  resource_group_name = var.resource_group_name

  administrator_login    = "dbadmin"
  administrator_password = random_password.db_password.result

  sku_name   = local.sku_name
  version    = var.postgres_version
  storage_mb = local.storage_mb

  backup_retention_days        = local.backup_retention_days
  geo_redundant_backup_enabled = var.env_type == "prod" ? true : false
  zone                         = local.zone_redundant ? null : "1"

  # High availability for production
  dynamic "high_availability" {
    for_each = local.zone_redundant ? [1] : []
    content {
      mode                      = "ZoneRedundant"
      standby_availability_zone = "2"
    }
  }

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-db"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Create database
resource "azurerm_postgresql_flexible_server_database" "main" {
  name      = local.db_name
  server_id = azurerm_postgresql_flexible_server.main.id
  charset   = "UTF8"
  collation = "en_US.utf8"
}

# Configure PostgreSQL server parameters
resource "azurerm_postgresql_flexible_server_configuration" "log_connections" {
  name      = "log_connections"
  server_id = azurerm_postgresql_flexible_server.main.id
  value     = "on"
}

resource "azurerm_postgresql_flexible_server_configuration" "log_disconnections" {
  name      = "log_disconnections"
  server_id = azurerm_postgresql_flexible_server.main.id
  value     = "on"
}

resource "azurerm_postgresql_flexible_server_configuration" "connection_throttling" {
  name      = "connection_throttling"
  server_id = azurerm_postgresql_flexible_server.main.id
  value     = "on"
}

# Create private DNS zone for PostgreSQL
resource "azurerm_private_dns_zone" "postgres" {
  name                = "privatelink.postgres.database.azure.com"
  resource_group_name = var.resource_group_name

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-postgres-dns"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Link private DNS zone to VNet
resource "azurerm_private_dns_zone_virtual_network_link" "postgres" {
  name                  = "${var.app_name}-${var.env_type}-postgres-dns-link"
  resource_group_name   = var.resource_group_name
  private_dns_zone_name = azurerm_private_dns_zone.postgres.name
  virtual_network_id    = var.vnet_id

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-postgres-dns-link"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Create Network Security Group for database subnet
resource "azurerm_network_security_group" "db" {
  name                = "${var.app_name}-${var.env_type}-db-nsg"
  location            = var.location
  resource_group_name = var.resource_group_name

  # Allow PostgreSQL from VNet
  security_rule {
    name                       = "AllowPostgreSQLFromVNet"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "5432"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "*"
  }

  # Deny all other inbound traffic
  security_rule {
    name                       = "DenyAllInbound"
    priority                   = 4096
    direction                  = "Inbound"
    access                     = "Deny"
    protocol                   = "*"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-db-nsg"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}
