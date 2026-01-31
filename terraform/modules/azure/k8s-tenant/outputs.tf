output "namespace" {
  description = "Kubernetes namespace name"
  value       = kubernetes_namespace.app.metadata[0].name
}

output "namespace_id" {
  description = "Kubernetes namespace ID"
  value       = kubernetes_namespace.app.id
}

output "service_account_name" {
  description = "Kubernetes service account name"
  value       = kubernetes_service_account.app.metadata[0].name
}

output "workload_identity_client_id" {
  description = "Azure Workload Identity client ID"
  value       = azurerm_user_assigned_identity.workload.client_id
}

output "workload_identity_principal_id" {
  description = "Azure Workload Identity principal ID"
  value       = azurerm_user_assigned_identity.workload.principal_id
}

output "workload_identity_id" {
  description = "Azure Workload Identity resource ID"
  value       = azurerm_user_assigned_identity.workload.id
}

output "resource_quota_name" {
  description = "Resource quota name"
  value       = kubernetes_resource_quota.app.metadata[0].name
}

output "limit_range_name" {
  description = "Limit range name"
  value       = kubernetes_limit_range.app.metadata[0].name
}

output "network_policy_name" {
  description = "Network policy name"
  value       = kubernetes_network_policy.app.metadata[0].name
}
