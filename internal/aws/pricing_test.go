package aws

import (
	"testing"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

func TestNewPricingCalculator(t *testing.T) {
	calc := NewPricingCalculator("us-east-1")
	testutil.AssertNotEqual(t, nil, calc)
	testutil.AssertEqual(t, "us-east-1", calc.region)
}

func TestPricingCalculator_CalculateVPCCost(t *testing.T) {
	tests := []struct {
		name    string
		envType string
		minCost float64
		maxCost float64
	}{
		{
			name:    "dev_environment_single_nat",
			envType: "dev",
			minCost: 30.0,  // 1 NAT Gateway
			maxCost: 50.0,
		},
		{
			name:    "staging_environment_dual_nat",
			envType: "staging",
			minCost: 60.0,  // 2 NAT Gateways
			maxCost: 100.0,
		},
		{
			name:    "prod_environment_dual_nat",
			envType: "prod",
			minCost: 60.0,  // 2 NAT Gateways
			maxCost: 100.0,
		},
		{
			name:    "unknown_environment_defaults_to_dual_nat",
			envType: "unknown",
			minCost: 60.0,  // Defaults to 2 NAT Gateways
			maxCost: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewPricingCalculator("us-east-1")
			cost := calc.CalculateVPCCost(tt.envType)

			testutil.AssertTrue(t, cost >= tt.minCost, "VPC cost should be >= minimum")
			testutil.AssertTrue(t, cost <= tt.maxCost, "VPC cost should be <= maximum")
		})
	}
}

func TestPricingCalculator_CalculateRDSCost(t *testing.T) {
	tests := []struct {
		name    string
		envType string
		minCost float64
		maxCost float64
	}{
		{
			name:    "dev_environment_t3_micro",
			envType: "dev",
			minCost: 10.0,
			maxCost: 20.0,
		},
		{
			name:    "staging_environment_t3_small",
			envType: "staging",
			minCost: 20.0,
			maxCost: 40.0,
		},
		{
			name:    "prod_environment_t3_medium_multi_az",
			envType: "prod",
			minCost: 100.0, // Multi-AZ doubles cost
			maxCost: 150.0,
		},
		{
			name:    "unknown_environment_defaults_to_zero",
			envType: "unknown",
			minCost: 0.0,
			maxCost: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewPricingCalculator("us-east-1")
			cost := calc.CalculateRDSCost(tt.envType)

			testutil.AssertTrue(t, cost >= tt.minCost, "RDS cost should be >= minimum")
			testutil.AssertTrue(t, cost <= tt.maxCost, "RDS cost should be <= maximum")
		})
	}
}

func TestPricingCalculator_CalculateEKSCost(t *testing.T) {
	tests := []struct {
		name    string
		envType string
		minCost float64
		maxCost float64
	}{
		{
			name:    "dev_environment_half_node",
			envType: "dev",
			minCost: 10.0,
			maxCost: 25.0,
		},
		{
			name:    "staging_environment_one_node",
			envType: "staging",
			minCost: 25.0,
			maxCost: 50.0,
		},
		{
			name:    "prod_environment_two_nodes",
			envType: "prod",
			minCost: 50.0,
			maxCost: 100.0,
		},
		{
			name:    "unknown_environment_defaults_to_zero",
			envType: "unknown",
			minCost: 0.0,
			maxCost: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewPricingCalculator("us-east-1")
			cost := calc.CalculateEKSCost(tt.envType)

			testutil.AssertTrue(t, cost >= tt.minCost, "EKS cost should be >= minimum")
			testutil.AssertTrue(t, cost <= tt.maxCost, "EKS cost should be <= maximum")
		})
	}
}

func TestPricingCalculator_CalculateTotalCost(t *testing.T) {
	tests := []struct {
		name    string
		envType string
		minCost float64
		maxCost float64
	}{
		{
			name:    "dev_environment_total",
			envType: "dev",
			minCost: 50.0,
			maxCost: 100.0,
		},
		{
			name:    "staging_environment_total",
			envType: "staging",
			minCost: 100.0,
			maxCost: 200.0,
		},
		{
			name:    "prod_environment_total",
			envType: "prod",
			minCost: 200.0,
			maxCost: 400.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewPricingCalculator("us-east-1")
			costs, err := calc.CalculateTotalCost(tt.envType)

			testutil.AssertNoError(t, err)
			testutil.AssertNotEqual(t, nil, costs)
			testutil.AssertEqual(t, tt.envType, costs.Environment)

			// Verify total is sum of components
			expectedTotal := costs.VPCCost + costs.RDSCost + costs.EKSCost
			testutil.AssertEqual(t, expectedTotal, costs.TotalCost)

			// Verify total is within expected range
			testutil.AssertTrue(t, costs.TotalCost >= tt.minCost, "Total cost should be >= minimum")
			testutil.AssertTrue(t, costs.TotalCost <= tt.maxCost, "Total cost should be <= maximum")

			// Verify all components are positive
			testutil.AssertTrue(t, costs.VPCCost > 0, "VPC cost should be positive")
			testutil.AssertTrue(t, costs.RDSCost > 0, "RDS cost should be positive")
			testutil.AssertTrue(t, costs.EKSCost > 0, "EKS cost should be positive")
		})
	}
}

func TestFormatCost(t *testing.T) {
	tests := []struct {
		name     string
		cost     float64
		expected string
	}{
		{
			name:     "zero_cost",
			cost:     0.0,
			expected: "$0.00",
		},
		{
			name:     "small_cost",
			cost:     12.34,
			expected: "$12.34",
		},
		{
			name:     "large_cost",
			cost:     1234.56,
			expected: "$1234.56",
		},
		{
			name:     "rounded_cost",
			cost:     99.999,
			expected: "$100.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCost(tt.cost)
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

func TestFormatCostBreakdown(t *testing.T) {
	costs := &EnvironmentCosts{
		VPCCost:     50.0,
		RDSCost:     100.0,
		EKSCost:     150.0,
		TotalCost:   300.0,
		Environment: "staging",
	}

	result := FormatCostBreakdown(costs)

	// Verify the breakdown contains all expected components
	testutil.AssertContains(t, result, "staging")
	testutil.AssertContains(t, result, "$50.00")
	testutil.AssertContains(t, result, "$100.00")
	testutil.AssertContains(t, result, "$150.00")
	testutil.AssertContains(t, result, "$300.00")
	testutil.AssertContains(t, result, "VPC")
	testutil.AssertContains(t, result, "RDS")
	testutil.AssertContains(t, result, "EKS")
	testutil.AssertContains(t, result, "Total")
}

func TestEnvironmentCosts_Fields(t *testing.T) {
	costs := &EnvironmentCosts{
		VPCCost:     50.0,
		RDSCost:     100.0,
		EKSCost:     150.0,
		TotalCost:   300.0,
		Environment: "staging",
	}

	testutil.AssertEqual(t, 50.0, costs.VPCCost)
	testutil.AssertEqual(t, 100.0, costs.RDSCost)
	testutil.AssertEqual(t, 150.0, costs.EKSCost)
	testutil.AssertEqual(t, 300.0, costs.TotalCost)
	testutil.AssertEqual(t, "staging", costs.Environment)
}

func TestPricingCalculator_DifferentRegions(t *testing.T) {
	regions := []string{"us-east-1", "us-west-2", "eu-west-1", "ap-southeast-1"}

	for _, region := range regions {
		t.Run("region_"+region, func(t *testing.T) {
			calc := NewPricingCalculator(region)
			testutil.AssertNotEqual(t, nil, calc)
			testutil.AssertEqual(t, region, calc.region)

			// Verify calculations work for all regions
			costs, err := calc.CalculateTotalCost("dev")
			testutil.AssertNoError(t, err)
			testutil.AssertTrue(t, costs.TotalCost > 0, "Total cost should be positive")
		})
	}
}
