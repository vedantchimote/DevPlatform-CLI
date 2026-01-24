output "vnet_id" {
  description = "Virtual Network ID"
  value       = azurerm_virtual_network.main.id
}

output "vnet_name" {
  description = "Virtual Network name"
  value       = azurerm_virtual_network.main.name
}

output "vnet_address_space" {
  description = "Virtual Network address space"
  value       = azurerm_virtual_network.main.address_space
}

output "public_subnet_ids" {
  description = "List of public subnet IDs"
  value       = azurerm_subnet.public[*].id
}

output "public_subnet_names" {
  description = "List of public subnet names"
  value       = azurerm_subnet.public[*].name
}

output "private_subnet_ids" {
  description = "List of private subnet IDs"
  value       = azurerm_subnet.private[*].id
}

output "private_subnet_names" {
  description = "List of private subnet names"
  value       = azurerm_subnet.private[*].name
}

output "nat_gateway_ids" {
  description = "List of NAT Gateway IDs"
  value       = azurerm_nat_gateway.main[*].id
}

output "nat_public_ip_addresses" {
  description = "List of NAT Gateway public IP addresses"
  value       = azurerm_public_ip.nat[*].ip_address
}

output "nsg_id" {
  description = "Network Security Group ID"
  value       = azurerm_network_security_group.main.id
}

output "nsg_name" {
  description = "Network Security Group name"
  value       = azurerm_network_security_group.main.name
}

output "network_watcher_id" {
  description = "Network Watcher ID"
  value       = var.enable_flow_logs ? azurerm_network_watcher.main[0].id : null
}

output "flow_logs_storage_account_id" {
  description = "Storage Account ID for flow logs"
  value       = var.enable_flow_logs ? azurerm_storage_account.flow_logs[0].id : null
}

output "flow_logs_workspace_id" {
  description = "Log Analytics Workspace ID for flow logs"
  value       = var.enable_flow_logs ? azurerm_log_analytics_workspace.flow_logs[0].id : null
}
