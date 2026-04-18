package azure

import (
	"testing"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

func TestNewPricingCalculator(t *testing.T) {
	calc := NewPricingCalculator("eastus")
	testutil.AssertNotEqual(t, nil, calc)
	testutil.AssertEqual(t, "eastus", calc.location)
}

func TestPricingCalculator_CalculateVNetCost(t *testing.T) {
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
			calc := NewPricingCalculator("eastus")
			cost := calc.CalculateVNetCost(tt.envType)

			testutil.AssertTrue(t, cost >= tt.minCost, "VNet cost should be >= minimum")
			testutil.AssertTrue(t, cost <= tt.maxCost, "VNet cost should be <= maximum")
		})
	}
}

func TestPricingCalculator_CalculateAzureDatabaseCost(t *testing.T) {
	tests := []struct {
		name    string
		envType string
		minCost float64
		maxCost float64
	}{
		{
			name:    "dev_environment_b1ms",
			envType: "dev",
			minCost: 15.0,
			maxCost: 30.0,
		},
		{
			name:    "staging_environment_b2s",
			envType: "staging",
			minCost: 35.0,
			maxCost: 60.0,
		},
		{
			name:    "prod_environment_gp_standard_zone_redundant",
			envType: "prod",
			minCost: 140.0, // Zone-redundant adds 50% cost
			maxCost: 250.0,
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
			calc := NewPricingCalculator("eastus")
			cost := calc.CalculateAzureDatabaseCost(tt.envType)

			testutil.AssertTrue(t, cost >= tt.minCost, "Database cost should be >= minimum")
			testutil.AssertTrue(t, cost <= tt.maxCost, "Database cost should be <= maximum")
		})
	}
}

func TestPricingCalculator_CalculateAKSCost(t *testing.T) {
	tests := []struct {
		name    string
		envType string
		minCost float64
		maxCost float64
	}{
		{
			name:    "dev_environment_half_node",
			envType: "dev",
			minCost: 30.0,
			maxCost: 50.0,
		},
		{
			name:    "staging_environment_one_node",
			envType: "staging",
			minCost: 60.0,
			maxCost: 100.0,
		},
		{
			name:    "prod_environment_two_nodes",
			envType: "prod",
			minCost: 120.0,
			maxCost: 200.0,
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
			calc := NewPricingCalculator("eastus")
			cost := calc.CalculateAKSCost(tt.envType)

			testutil.AssertTrue(t, cost >= tt.minCost, "AKS cost should be >= minimum")
			testutil.AssertTrue(t, cost <= tt.maxCost, "AKS cost should be <= maximum")
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
			minCost: 75.0,
			maxCost: 130.0,
		},
		{
			name:    "staging_environment_total",
			envType: "staging",
			minCost: 155.0,
			maxCost: 260.0,
		},
		{
			name:    "prod_environment_total",
			envType: "prod",
			minCost: 320.0,
			maxCost: 550.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewPricingCalculator("eastus")
			costs, err := calc.CalculateTotalCost(tt.envType)

			testutil.AssertNoError(t, err)
			testutil.AssertNotEqual(t, nil, costs)
			testutil.AssertEqual(t, tt.envType, costs.Environment)

			// Verify total is sum of components
			expectedTotal := costs.VNetCost + costs.DatabaseCost + costs.AKSCost
			testutil.AssertEqual(t, expectedTotal, costs.TotalCost)

			// Verify total is within expected range
			testutil.AssertTrue(t, costs.TotalCost >= tt.minCost, "Total cost should be >= minimum")
			testutil.AssertTrue(t, costs.TotalCost <= tt.maxCost, "Total cost should be <= maximum")

			// Verify all components are positive
			testutil.AssertTrue(t, costs.VNetCost > 0, "VNet cost should be positive")
			testutil.AssertTrue(t, costs.DatabaseCost > 0, "Database cost should be positive")
			testutil.AssertTrue(t, costs.AKSCost > 0, "AKS cost should be positive")
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
		VNetCost:     50.0,
		DatabaseCost: 100.0,
		AKSCost:      150.0,
		TotalCost:    300.0,
		Environment:  "staging",
	}

	result := FormatCostBreakdown(costs)

	// Verify the breakdown contains all expected components
	testutil.AssertContains(t, result, "staging")
	testutil.AssertContains(t, result, "$50.00")
	testutil.AssertContains(t, result, "$100.00")
	testutil.AssertContains(t, result, "$150.00")
	testutil.AssertContains(t, result, "$300.00")
	testutil.AssertContains(t, result, "VNet")
	testutil.AssertContains(t, result, "Azure Database")
	testutil.AssertContains(t, result, "AKS")
	testutil.AssertContains(t, result, "Total")
}

func TestEnvironmentCosts_Fields(t *testing.T) {
	costs := &EnvironmentCosts{
		VNetCost:     50.0,
		DatabaseCost: 100.0,
		AKSCost:      150.0,
		TotalCost:    300.0,
		Environment:  "staging",
	}

	testutil.AssertEqual(t, 50.0, costs.VNetCost)
	testutil.AssertEqual(t, 100.0, costs.DatabaseCost)
	testutil.AssertEqual(t, 150.0, costs.AKSCost)
	testutil.AssertEqual(t, 300.0, costs.TotalCost)
	testutil.AssertEqual(t, "staging", costs.Environment)
}

func TestPricingCalculator_DifferentLocations(t *testing.T) {
	locations := []string{"eastus", "westus2", "northeurope", "southeastasia"}

	for _, location := range locations {
		t.Run("location_"+location, func(t *testing.T) {
			calc := NewPricingCalculator(location)
			testutil.AssertNotEqual(t, nil, calc)
			testutil.AssertEqual(t, location, calc.location)

			// Verify calculations work for all locations
			costs, err := calc.CalculateTotalCost("dev")
			testutil.AssertNoError(t, err)
			testutil.AssertTrue(t, costs.TotalCost > 0, "Total cost should be positive")
		})
	}
}
