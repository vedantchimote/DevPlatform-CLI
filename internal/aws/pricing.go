package aws

import (
	"fmt"
)

// PricingCalculator handles AWS cost estimation
type PricingCalculator struct {
	region string
}

// NewPricingCalculator creates a new pricing calculator
func NewPricingCalculator(region string) *PricingCalculator {
	return &PricingCalculator{
		region: region,
	}
}

// EnvironmentCosts represents the total costs for an environment
type EnvironmentCosts struct {
	VPCCost     float64
	RDSCost     float64
	EKSCost     float64
	TotalCost   float64
	Environment string
}

// CalculateTotalCost calculates the total monthly cost for an environment
func (p *PricingCalculator) CalculateTotalCost(envType string) (*EnvironmentCosts, error) {
	vpcCost := p.CalculateVPCCost(envType)
	rdsCost := p.CalculateRDSCost(envType)
	eksCost := p.CalculateEKSCost(envType)
	
	return &EnvironmentCosts{
		VPCCost:     vpcCost,
		RDSCost:     rdsCost,
		EKSCost:     eksCost,
		TotalCost:   vpcCost + rdsCost + eksCost,
		Environment: envType,
	}, nil
}

// CalculateVPCCost calculates the monthly cost for VPC resources
func (p *PricingCalculator) CalculateVPCCost(envType string) float64 {
	// VPC itself is free, but NAT Gateways have costs
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

// CalculateRDSCost calculates the monthly cost for RDS instance
func (p *PricingCalculator) CalculateRDSCost(envType string) float64 {
	// RDS pricing varies by instance type and region
	// Using us-east-1 pricing as baseline
	
	var instanceCost float64
	var storageCost float64
	var multiAZMultiplier float64 = 1.0
	
	switch envType {
	case "dev":
		// db.t3.micro: $0.017/hour
		instanceCost = 0.017 * 730
		// 20GB storage: $0.115/GB-month
		storageCost = 20 * 0.115
		
	case "staging":
		// db.t3.small: $0.034/hour
		instanceCost = 0.034 * 730
		// 50GB storage
		storageCost = 50 * 0.115
		
	case "prod":
		// db.t3.medium: $0.068/hour
		instanceCost = 0.068 * 730
		// 100GB storage
		storageCost = 100 * 0.115
		// Multi-AZ doubles the cost
		multiAZMultiplier = 2.0
	}
	
	return (instanceCost + storageCost) * multiAZMultiplier
}

// CalculateEKSCost calculates the monthly cost for EKS tenant resources
func (p *PricingCalculator) CalculateEKSCost(envType string) float64 {
	// EKS cluster control plane: $0.10/hour (shared across tenants, not included)
	// We only calculate the tenant namespace resource quota costs
	// This is an estimate based on the resource requests
	
	// Worker node costs (estimated based on resource quotas)
	// Assuming t3.medium nodes at $0.0416/hour
	
	var estimatedNodeFraction float64
	
	switch envType {
	case "dev":
		// Dev: 2 CPU, 4GB RAM - approximately 0.5 of a t3.medium node
		estimatedNodeFraction = 0.5
		
	case "staging":
		// Staging: 4 CPU, 8GB RAM - approximately 1 t3.medium node
		estimatedNodeFraction = 1.0
		
	case "prod":
		// Prod: 8 CPU, 16GB RAM - approximately 2 t3.medium nodes
		estimatedNodeFraction = 2.0
	}
	
	// t3.medium: $0.0416/hour * 730 hours/month
	nodeCost := estimatedNodeFraction * 0.0416 * 730
	
	return nodeCost
}

// FormatCost formats a cost value as a USD string
func FormatCost(cost float64) string {
	return fmt.Sprintf("$%.2f", cost)
}

// FormatCostBreakdown formats a cost breakdown as a string
func FormatCostBreakdown(costs *EnvironmentCosts) string {
	return fmt.Sprintf(`Environment: %s
VPC (NAT Gateways):     %s/month
RDS (Database):         %s/month
EKS (Tenant Resources): %s/month
-----------------------------------
Total:                  %s/month`,
		costs.Environment,
		FormatCost(costs.VPCCost),
		FormatCost(costs.RDSCost),
		FormatCost(costs.EKSCost),
		FormatCost(costs.TotalCost),
	)
}
