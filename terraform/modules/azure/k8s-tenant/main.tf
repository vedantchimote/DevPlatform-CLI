terraform {
  required_version = ">= 1.0"
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }
}

# Get current Azure client configuration
data "azurerm_client_config" "current" {}

# Environment-specific resource quotas
locals {
  namespace = "${var.app_name}-${var.env_type}-azure"

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
      name          = local.namespace
      app           = var.app_name
      environment   = var.env_type
      cloud-provider = "azure"
      managed-by    = "devplatform-cli"
    }

    annotations = {
      "devplatform.io/app-name"       = var.app_name
      "devplatform.io/env-type"       = var.env_type
      "devplatform.io/cloud-provider" = "azure"
      "devplatform.io/timestamp"      = timestamp()
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

# Create Azure User Assigned Identity for workload identity
resource "azurerm_user_assigned_identity" "workload" {
  name                = "${var.app_name}-${var.env_type}-workload-identity"
  location            = var.location
  resource_group_name = var.resource_group_name

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-workload-identity"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Create federated identity credential for workload identity
resource "azurerm_federated_identity_credential" "workload" {
  name                = "${var.app_name}-${var.env_type}-federated-credential"
  resource_group_name = var.resource_group_name
  parent_id           = azurerm_user_assigned_identity.workload.id
  audience            = ["api://AzureADTokenExchange"]
  issuer              = var.oidc_issuer_url
  subject             = "system:serviceaccount:${local.namespace}:${var.app_name}-sa"
}

# Grant Key Vault access to workload identity
resource "azurerm_role_assignment" "keyvault_secrets_user" {
  count                = var.keyvault_id != "" ? 1 : 0
  scope                = var.keyvault_id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_user_assigned_identity.workload.principal_id
}

# Create Kubernetes service account with workload identity annotation
resource "kubernetes_service_account" "app" {
  metadata {
    name      = "${var.app_name}-sa"
    namespace = kubernetes_namespace.app.metadata[0].name

    annotations = {
      "azure.workload.identity/client-id" = azurerm_user_assigned_identity.workload.client_id
      "azure.workload.identity/tenant-id" = data.azurerm_client_config.current.tenant_id
    }

    labels = {
      app            = var.app_name
      environment    = var.env_type
      cloud-provider = "azure"
      managed-by     = "devplatform-cli"
      "azure.workload.identity/use" = "true"
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
