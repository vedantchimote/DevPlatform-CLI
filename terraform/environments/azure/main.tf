# =============================================================================
# DevPlatform-CLI - Azure Root Module
# Creates all infrastructure for a complete environment
# =============================================================================

locals {
  resource_group_name = "devplatform-${var.app_name}-${var.env_type}"
  cluster_name        = "${var.app_name}-${var.env_type}-aks"
  common_tags = merge(var.tags, {
    App_Name       = var.app_name
    Env_Type       = var.env_type
    Cloud_Provider = "azure"
    ManagedBy      = "devplatform-cli"
  })
}

# =============================================================================
# Resource Group
# =============================================================================
resource "azurerm_resource_group" "main" {
  name     = local.resource_group_name
  location = var.location
  tags     = local.common_tags
}

# =============================================================================
# Network Module (VNet, Subnets, NAT Gateway, NSG)
# =============================================================================
module "network" {
  source = "../../modules/azure/network"

  app_name            = var.app_name
  env_type            = var.env_type
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location
  vnet_cidr           = var.vnet_cidr
  zone_count          = var.zone_count
  nat_gateway_count   = var.nat_gateway_count
  enable_flow_logs    = var.enable_flow_logs
  tags                = local.common_tags
}

# =============================================================================
# AKS Subnet (separate from network module to avoid NSG conflicts)
# =============================================================================
resource "azurerm_subnet" "aks" {
  name                 = "${var.app_name}-${var.env_type}-aks-subnet"
  resource_group_name  = azurerm_resource_group.main.name
  virtual_network_name = module.network.vnet_name
  address_prefixes     = [cidrsubnet(var.vnet_cidr, 4, 4)] # e.g. 10.10.64.0/20
}

# =============================================================================
# AKS Cluster
# =============================================================================
resource "azurerm_kubernetes_cluster" "main" {
  name                = local.cluster_name
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  dns_prefix          = "${var.app_name}-${var.env_type}"

  default_node_pool {
    name           = "default"
    node_count     = var.aks_node_count
    vm_size        = var.aks_vm_size
    vnet_subnet_id = azurerm_subnet.aks.id

    # Cost optimisation for dev
    os_disk_size_gb = 30
  }

  identity {
    type = "SystemAssigned"
  }

  # Enable OIDC & Workload Identity for k8s-tenant module
  oidc_issuer_enabled       = true
  workload_identity_enabled = true

  network_profile {
    network_plugin = "azure"
    service_cidr   = "10.200.0.0/16"
    dns_service_ip = "10.200.0.10"
  }

  tags = local.common_tags
}

# Grant AKS identity Network Contributor on its subnet
resource "azurerm_role_assignment" "aks_network" {
  scope                = azurerm_subnet.aks.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_kubernetes_cluster.main.identity[0].principal_id
}

# Database module omitted for Azure Student subscription compatibility

# =============================================================================
# K8s Tenant Module (Namespace, Quotas, Workload Identity, Service Account)
# =============================================================================
module "k8s_tenant" {
  source = "../../modules/azure/k8s-tenant"

  app_name            = var.app_name
  env_type            = var.env_type
  cluster_name        = azurerm_kubernetes_cluster.main.name
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location
  oidc_issuer_url     = azurerm_kubernetes_cluster.main.oidc_issuer_url
  enable_keyvault_access = false
  tags                = local.common_tags

  depends_on = [azurerm_kubernetes_cluster.main]
}
