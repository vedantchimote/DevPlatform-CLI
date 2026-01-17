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

output "service_account_arn" {
  description = "IAM role ARN for the service account"
  value       = aws_iam_role.service_account.arn
}

output "iam_role_name" {
  description = "IAM role name for the service account"
  value       = aws_iam_role.service_account.name
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
