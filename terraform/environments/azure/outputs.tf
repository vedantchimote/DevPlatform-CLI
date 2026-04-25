# =============================================================================
# Outputs matching the CLI's internal/terraform/output.go TerraformOutputs
# =============================================================================

# --- Network Outputs ---
output "vnet_id" {
  description = "Virtual Network ID"
  value       = module.network.vnet_id
}

output "subnet_ids" {
  description = "Private subnet IDs"
  value       = module.network.private_subnet_ids
}

output "nsg_id" {
  description = "Network Security Group ID"
  value       = module.network.nsg_id
}

# Database omitted

# --- Kubernetes Outputs ---
output "namespace" {
  description = "Kubernetes namespace"
  value       = module.k8s_tenant.namespace
}

output "service_account_name" {
  description = "Kubernetes service account"
  value       = module.k8s_tenant.service_account_name
}

output "cluster_name" {
  description = "AKS cluster name"
  value       = azurerm_kubernetes_cluster.main.name
}

# --- Metadata ---
output "provider" {
  description = "Cloud provider"
  value       = "azure"
}

output "environment" {
  description = "Environment type"
  value       = var.env_type
}

output "app_name" {
  description = "Application name"
  value       = var.app_name
}

output "resource_group_name" {
  description = "Resource group name"
  value       = azurerm_resource_group.main.name
}

# --- Connection Info ---
output "kube_config_command" {
  description = "Command to configure kubectl"
  value       = "az aks get-credentials --resource-group ${azurerm_resource_group.main.name} --name ${azurerm_kubernetes_cluster.main.name}"
}
