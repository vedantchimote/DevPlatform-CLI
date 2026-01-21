terraform {
  required_version = ">= 1.0"
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# Environment-specific resource quotas
locals {
  namespace = "${var.app_name}-${var.env_type}"

  resource_quotas = {
    dev = {
      cpu_requests    = "4"
      cpu_limits      = "8"
      memory_requests = "8Gi"
      memory_limits   = "16Gi"
      pods            = "20"
      services        = "10"
      persistentvolumeclaims = "5"
    }
    staging = {
      cpu_requests    = "8"
      cpu_limits      = "16"
      memory_requests = "16Gi"
      memory_limits   = "32Gi"
      pods            = "40"
      services        = "20"
      persistentvolumeclaims = "10"
    }
    prod = {
      cpu_requests    = "16"
      cpu_limits      = "32"
      memory_requests = "32Gi"
      memory_limits   = "64Gi"
      pods            = "100"
      services        = "50"
      persistentvolumeclaims = "20"
    }
  }

  quota = local.resource_quotas[var.env_type]
}

# Create Kubernetes namespace
resource "kubernetes_namespace" "app" {
  metadata {
    name = local.namespace

    labels = {
      name        = local.namespace
      app         = var.app_name
      environment = var.env_type
      managed-by  = "devplatform-cli"
    }

    annotations = {
      "devplatform.io/app-name"  = var.app_name
      "devplatform.io/env-type"  = var.env_type
      "devplatform.io/timestamp" = timestamp()
    }
  }
}

# Create resource quota for the namespace
resource "kubernetes_resource_quota" "app" {
  metadata {
    name      = "${local.namespace}-quota"
    namespace = kubernetes_namespace.app.metadata[0].name
  }

  spec {
    hard = {
      "requests.cpu"               = local.quota.cpu_requests
      "requests.memory"            = local.quota.memory_requests
      "limits.cpu"                 = local.quota.cpu_limits
      "limits.memory"              = local.quota.memory_limits
      "pods"                       = local.quota.pods
      "services"                   = local.quota.services
      "persistentvolumeclaims"     = local.quota.persistentvolumeclaims
    }
  }
}

# Create limit range for the namespace
resource "kubernetes_limit_range" "app" {
  metadata {
    name      = "${local.namespace}-limits"
    namespace = kubernetes_namespace.app.metadata[0].name
  }

  spec {
    limit {
      type = "Container"
      default = {
        cpu    = "500m"
        memory = "512Mi"
      }
      default_request = {
        cpu    = "100m"
        memory = "128Mi"
      }
    }

    limit {
      type = "Pod"
      max = {
        cpu    = "4"
        memory = "8Gi"
      }
    }
  }
}

# Create IAM role for service account (IRSA)
resource "aws_iam_role" "service_account" {
  name = "${var.app_name}-${var.env_type}-sa-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = var.oidc_provider_arn
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "${replace(var.oidc_provider_url, "https://", "")}:sub" = "system:serviceaccount:${local.namespace}:${var.app_name}-sa"
            "${replace(var.oidc_provider_url, "https://", "")}:aud" = "sts.amazonaws.com"
          }
        }
      }
    ]
  })

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-sa-role"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Attach policy to allow reading database secrets
resource "aws_iam_role_policy" "secrets_access" {
  count = var.db_secret_arn != "" ? 1 : 0
  name  = "${var.app_name}-${var.env_type}-secrets-policy"
  role  = aws_iam_role.service_account.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue",
          "secretsmanager:DescribeSecret"
        ]
        Resource = var.db_secret_arn
      }
    ]
  })
}

# Create Kubernetes service account with IRSA annotation
resource "kubernetes_service_account" "app" {
  metadata {
    name      = "${var.app_name}-sa"
    namespace = kubernetes_namespace.app.metadata[0].name

    annotations = {
      "eks.amazonaws.com/role-arn" = aws_iam_role.service_account.arn
    }

    labels = {
      app        = var.app_name
      environment = var.env_type
      managed-by = "devplatform-cli"
    }
  }
}

# Create network policy for namespace isolation
resource "kubernetes_network_policy" "app" {
  metadata {
    name      = "${local.namespace}-network-policy"
    namespace = kubernetes_namespace.app.metadata[0].name
  }

  spec {
    pod_selector {}

    policy_types = ["Ingress", "Egress"]

    # Allow ingress from same namespace
    ingress {
      from {
        pod_selector {}
      }
    }

    # Allow ingress from ingress controller
    ingress {
      from {
        namespace_selector {
          match_labels = {
            name = "ingress-nginx"
          }
        }
      }
    }

    # Allow all egress
    egress {}
  }
}
