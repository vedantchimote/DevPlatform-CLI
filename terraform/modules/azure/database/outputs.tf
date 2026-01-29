output "database_endpoint" {
  description = "PostgreSQL server FQDN"
  value       = azurerm_postgresql_flexible_server.main.fqdn
}

output "database_name" {
  description = "Database name"
  value       = azurerm_postgresql_flexible_server_database.main.name
}

output "database_port" {
  description = "Database port"
  value       = 5432
}

output "database_id" {
  description = "PostgreSQL server ID"
  value       = azurerm_postgresql_flexible_server.main.id
}

output "database_username" {
  description = "Database administrator username"
  value       = azurerm_postgresql_flexible_server.main.administrator_login
  sensitive   = true
}

output "keyvault_id" {
  description = "Key Vault ID"
  value       = azurerm_key_vault.db.id
}

output "keyvault_name" {
  description = "Key Vault name"
  value       = azurerm_key_vault.db.name
}

output "keyvault_uri" {
  description = "Key Vault URI"
  value       = azurerm_key_vault.db.vault_uri
}

output "keyvault_secret_id_password" {
  description = "Key Vault secret ID for database password"
  value       = azurerm_key_vault_secret.db_password.id
  sensitive   = true
}

output "keyvault_secret_id_connection_string" {
  description = "Key Vault secret ID for database connection string"
  value       = azurerm_key_vault_secret.db_connection_string.id
  sensitive   = true
}

output "private_dns_zone_id" {
  description = "Private DNS zone ID"
  value       = azurerm_private_dns_zone.postgres.id
}

output "nsg_id" {
  description = "Network Security Group ID for database"
  value       = azurerm_network_security_group.db.id
}
