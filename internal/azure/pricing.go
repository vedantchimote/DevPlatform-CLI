package azure

import (
	"fmt"
)

// PricingCalculator handles Azure cost estimation
type PricingCalculator struct {
	location string
}

// NewPricingCalculator creates a new pricing calculator
func NewPricingCalculator(location string) *PricingCalculator {
	return &PricingCalculator{
		location: location,
	}
}

// EnvironmentCosts represents the total costs for an environment
type EnvironmentCosts struct {
	VNetCost     float64
	DatabaseCost float64
	AKSCost      float64
	TotalCost    float64
	Environment  string
}

// CalculateTotalCost calculates the total monthly cost for an environment
func (p *PricingCalculator) CalculateTotalCost(envType string) (*EnvironmentCosts, error) {
	vnetCost := p.CalculateVNetCost(envType)
	dbCost := p.CalculateAzureDatabaseCost(envType)
	aksCost := p.CalculateAKSCost(envType)
	
	return &EnvironmentCosts{
		VNetCost:     vnetCost,
		DatabaseCost: dbCost,
		AKSCost:      aksCost,
		TotalCost:    vnetCost + dbCost + aksCost,
		Environment:  envType,
	}, nil
}

// CalculateVNetCost calculates the monthly cost for VNet resources
func (p *PricingCalculator) CalculateVNetCost(envType string) float64 {
	// VNet itself is free, but NAT Gateway has costs
	// NAT Gateway: $0.045/hour + $0.045/GB processed
	// Assuming 2 NAT Gateways for HA (one per AZ)
	
	natGateways := 2
	if envType == "dev" {
		natGateways = 1 // Single NAT Gateway for dev
	}
	
	// NAT Gateway hourly cost: $0.045/hour * 730 hours/month
	natGatewayCost := float64(natGateways) * 0.045 * 730
	
	// Data processing cost (estimated 100GB/month per NAT Gateway)
	dataProcessingCost := float64(natGateways) * 0.045 * 100
	
	return natGatewayCost + dataProcessingCost
}

// CalculateAzureDatabaseCost calculates the monthly cost for Azure Database for PostgreSQL
func (p *PricingCalculator) CalculateAzureDatabaseCost(envType string) float64 {
	// Azure Database for PostgreSQL Flexible Server pricing
	// Using East US pricing as baseline
	
	var computeCost float64
	var storageCost float64
	var zoneRedundantMultiplier float64 = 1.0
	
	switch envType {
	case "dev":
		// B1ms (1 vCore, 2GB RAM): $0.0255/hour
		computeCost = 0.0255 * 730
		// 32GB storage: $0.115/GB-month
		storageCost = 32 * 0.115
		
	case "staging":
		// B2s (2 vCore, 4GB RAM): $0.051/hour
		computeCost = 0.051 * 730
		// 64GB storage
		storageCost = 64 * 0.115
		
	case "prod":
		// GP_Standard_D2s_v3 (2 vCore, 8GB RAM): $0.192/hour
		computeCost = 0.192 * 730
		// 128GB storage
		storageCost = 128 * 0.115
		// Zone-redundant HA adds ~50% cost
		zoneRedundantMultiplier = 1.5
	}
	
	return (computeCost + storageCost) * zoneRedundantMultiplier
}

// CalculateAKSCost calculates the monthly cost for AKS tenant resources
func (p *PricingCalculator) CalculateAKSCost(envType string) float64 {
	// AKS cluster control plane: Free for standard tier (shared across tenants)
	// We only calculate the tenant namespace resource quota costs
	// This is an estimate based on the resource requests
	
	// Worker node costs (estimated based on resource quotas)
	// Assuming Standard_D2s_v3 nodes at $0.096/hour
	
	var estimatedNodeFraction float64
	
	switch envType {
	case "dev":
		// Dev: 2 CPU, 4GB RAM - approximately 0.5 of a Standard_D2s_v3 node
		estimatedNodeFraction = 0.5
		
	case "staging":
		// Staging: 4 CPU, 8GB RAM - approximately 1 Standard_D2s_v3 node
		estimatedNodeFraction = 1.0
		
	case "prod":
		// Prod: 8 CPU, 16GB RAM - approximately 2 Standard_D2s_v3 nodes
		estimatedNodeFraction = 2.0
	}
	
	// Standard_D2s_v3: $0.096/hour * 730 hours/month
	nodeCost := estimatedNodeFraction * 0.096 * 730
	
	return nodeCost
}

// FormatCost formats a cost value as a USD string
func FormatCost(cost float64) string {
	return fmt.Sprintf("$%.2f", cost)
}

// FormatCostBreakdown formats a cost breakdown as a string
func FormatCostBreakdown(costs *EnvironmentCosts) string {
	return fmt.Sprintf(`Environment: %s
VNet (NAT Gateways):        %s/month
Azure Database:             %s/month
AKS (Tenant Resources):     %s/month
---------------------------------------
Total:                      %s/month`,
		costs.Environment,
		FormatCost(costs.VNetCost),
		FormatCost(costs.DatabaseCost),
		FormatCost(costs.AKSCost),
		FormatCost(costs.TotalCost),
	)
}
