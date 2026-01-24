terraform {
  required_version = ">= 1.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }
}

# Create Virtual Network
resource "azurerm_virtual_network" "main" {
  name                = "${var.app_name}-${var.env_type}-vnet"
  location            = var.location
  resource_group_name = var.resource_group_name
  address_space       = [var.vnet_cidr]

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-vnet"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Create public subnets
resource "azurerm_subnet" "public" {
  count                = var.zone_count
  name                 = "${var.app_name}-${var.env_type}-public-${count.index + 1}"
  resource_group_name  = var.resource_group_name
  virtual_network_name = azurerm_virtual_network.main.name
  address_prefixes     = [cidrsubnet(var.vnet_cidr, 4, count.index)]
}

# Create private subnets
resource "azurerm_subnet" "private" {
  count                = var.zone_count
  name                 = "${var.app_name}-${var.env_type}-private-${count.index + 1}"
  resource_group_name  = var.resource_group_name
  virtual_network_name = azurerm_virtual_network.main.name
  address_prefixes     = [cidrsubnet(var.vnet_cidr, 4, count.index + var.zone_count)]
}

# Create public IP addresses for NAT Gateways
resource "azurerm_public_ip" "nat" {
  count               = var.nat_gateway_count
  name                = "${var.app_name}-${var.env_type}-nat-pip-${count.index + 1}"
  location            = var.location
  resource_group_name = var.resource_group_name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = var.zone_count > 1 ? [tostring((count.index % var.zone_count) + 1)] : null

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-nat-pip-${count.index + 1}"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Create NAT Gateways
resource "azurerm_nat_gateway" "main" {
  count               = var.nat_gateway_count
  name                = "${var.app_name}-${var.env_type}-nat-${count.index + 1}"
  location            = var.location
  resource_group_name = var.resource_group_name
  sku_name            = "Standard"
  zones               = var.zone_count > 1 ? [tostring((count.index % var.zone_count) + 1)] : null

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-nat-${count.index + 1}"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Associate public IPs with NAT Gateways
resource "azurerm_nat_gateway_public_ip_association" "main" {
  count                = var.nat_gateway_count
  nat_gateway_id       = azurerm_nat_gateway.main[count.index].id
  public_ip_address_id = azurerm_public_ip.nat[count.index].id
}

# Associate NAT Gateways with private subnets
resource "azurerm_subnet_nat_gateway_association" "main" {
  count          = var.zone_count
  subnet_id      = azurerm_subnet.private[count.index].id
  nat_gateway_id = azurerm_nat_gateway.main[count.index % var.nat_gateway_count].id
}

# Create Network Security Group
resource "azurerm_network_security_group" "main" {
  name                = "${var.app_name}-${var.env_type}-nsg"
  location            = var.location
  resource_group_name = var.resource_group_name

  # Allow internal VNet traffic
  security_rule {
    name                       = "AllowVNetInbound"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "*"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "VirtualNetwork"
  }

  # Allow outbound internet access
  security_rule {
    name                       = "AllowInternetOutbound"
    priority                   = 100
    direction                  = "Outbound"
    access                     = "Allow"
    protocol                   = "*"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "Internet"
  }

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-nsg"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Associate NSG with public subnets
resource "azurerm_subnet_network_security_group_association" "public" {
  count                     = var.zone_count
  subnet_id                 = azurerm_subnet.public[count.index].id
  network_security_group_id = azurerm_network_security_group.main.id
}

# Associate NSG with private subnets
resource "azurerm_subnet_network_security_group_association" "private" {
  count                     = var.zone_count
  subnet_id                 = azurerm_subnet.private[count.index].id
  network_security_group_id = azurerm_network_security_group.main.id
}

# Create Network Watcher (if not exists)
resource "azurerm_network_watcher" "main" {
  count               = var.enable_flow_logs ? 1 : 0
  name                = "NetworkWatcher_${var.location}"
  location            = var.location
  resource_group_name = var.resource_group_name

  tags = merge(
    var.tags,
    {
      Name      = "NetworkWatcher_${var.location}"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Create Storage Account for flow logs
resource "azurerm_storage_account" "flow_logs" {
  count                    = var.enable_flow_logs ? 1 : 0
  name                     = replace("${var.app_name}${var.env_type}flowlogs", "-", "")
  resource_group_name      = var.resource_group_name
  location                 = var.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  min_tls_version          = "TLS1_2"

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-flow-logs"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Create Log Analytics Workspace for flow logs
resource "azurerm_log_analytics_workspace" "flow_logs" {
  count               = var.enable_flow_logs ? 1 : 0
  name                = "${var.app_name}-${var.env_type}-flow-logs-workspace"
  location            = var.location
  resource_group_name = var.resource_group_name
  sku                 = "PerGB2018"
  retention_in_days   = 7

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-flow-logs-workspace"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Enable NSG flow logs
resource "azurerm_network_watcher_flow_log" "main" {
  count                     = var.enable_flow_logs ? 1 : 0
  name                      = "${var.app_name}-${var.env_type}-nsg-flow-log"
  network_watcher_name      = azurerm_network_watcher.main[0].name
  resource_group_name       = var.resource_group_name
  network_security_group_id = azurerm_network_security_group.main.id
  storage_account_id        = azurerm_storage_account.flow_logs[0].id
  enabled                   = true
  version                   = 2

  retention_policy {
    enabled = true
    days    = 7
  }

  traffic_analytics {
    enabled               = true
    workspace_id          = azurerm_log_analytics_workspace.flow_logs[0].workspace_id
    workspace_region      = var.location
    workspace_resource_id = azurerm_log_analytics_workspace.flow_logs[0].id
    interval_in_minutes   = 10
  }

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-nsg-flow-log"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}
