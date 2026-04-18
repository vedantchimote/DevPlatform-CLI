package aws

import (
	"context"
	"testing"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

func TestNewAWSProvider(t *testing.T) {
	tests := []struct {
		name    string
		region  string
		profile string
		wantErr bool
	}{
		{
			name:    "valid_provider_with_region",
			region:  "us-east-1",
			profile: "",
			wantErr: false,
		},
		{
			name:    "empty_region",
			region:  "",
			profile: "",
			wantErr: false, // AWS SDK allows empty region
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			provider, err := NewAWSProvider(ctx, tt.region, tt.profile)

			if tt.wantErr {
				testutil.AssertError(t, err)
				return
			}

			testutil.AssertNoError(t, err)
			testutil.AssertNotEqual(t, nil, provider)
			testutil.AssertEqual(t, tt.region, provider.region)
			testutil.AssertEqual(t, tt.profile, provider.profile)
			testutil.AssertNotEqual(t, nil, provider.auth)
			testutil.AssertNotEqual(t, nil, provider.kubeconfig)
			testutil.AssertNotEqual(t, nil, provider.pricing)
		})
	}
}

func TestAWSProvider_GetProviderName(t *testing.T) {
	ctx := context.Background()
	provider, err := NewAWSProvider(ctx, "us-east-1", "")
	testutil.AssertNoError(t, err)

	name := provider.GetProviderName()
	testutil.AssertEqual(t, "aws", name)
}

func TestAWSProvider_GetModulePath(t *testing.T) {
	ctx := context.Background()
	provider, err := NewAWSProvider(ctx, "us-east-1", "")
	testutil.AssertNoError(t, err)

	path := provider.GetModulePath()
	testutil.AssertEqual(t, "terraform/modules/aws", path)
}

func TestAWSProvider_GetTerraformBackend(t *testing.T) {
	tests := []struct {
		name    string
		region  string
		appName string
		envType string
	}{
		{
			name:    "dev_environment",
			region:  "us-east-1",
			appName: "myapp",
			envType: "dev",
		},
		{
			name:    "staging_environment",
			region:  "us-west-2",
			appName: "testapp",
			envType: "staging",
		},
		{
			name:    "prod_environment",
			region:  "eu-west-1",
			appName: "prodapp",
			envType: "prod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			provider, err := NewAWSProvider(ctx, tt.region, "")
			testutil.AssertNoError(t, err)

			backend, err := provider.GetTerraformBackend(tt.appName, tt.envType)
			testutil.AssertNoError(t, err)
			testutil.AssertNotEqual(t, nil, backend)

			// Verify backend type
			testutil.AssertEqual(t, "s3", backend.Type)

			// Verify backend config
			testutil.AssertNotEqual(t, nil, backend.Config)
			testutil.AssertEqual(t, "devplatform-terraform-state-"+tt.region, backend.Config["bucket"])
			testutil.AssertEqual(t, tt.appName+"/"+tt.envType+"/terraform.tfstate", backend.Config["key"])
			testutil.AssertEqual(t, tt.region, backend.Config["region"])
			testutil.AssertEqual(t, "devplatform-terraform-locks", backend.Config["dynamodb_table"])
			testutil.AssertEqual(t, "true", backend.Config["encrypt"])
		})
	}
}

func TestAWSProvider_GetConnectionCommands(t *testing.T) {
	tests := []struct {
		name        string
		region      string
		clusterName string
		namespace   string
	}{
		{
			name:        "dev_cluster",
			region:      "us-east-1",
			clusterName: "dev-cluster",
			namespace:   "dev-ns",
		},
		{
			name:        "prod_cluster",
			region:      "us-west-2",
			clusterName: "prod-cluster",
			namespace:   "prod-ns",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			provider, err := NewAWSProvider(ctx, tt.region, "")
			testutil.AssertNoError(t, err)

			commands := provider.GetConnectionCommands(tt.clusterName, tt.namespace)
			testutil.AssertNotEqual(t, nil, commands)
			testutil.AssertEqual(t, 2, len(commands))

			// Verify commands contain cluster and namespace
			testutil.AssertContains(t, commands[0], tt.clusterName)
			testutil.AssertContains(t, commands[1], tt.namespace)
		})
	}
}

func TestAWSProvider_CalculateTotalCost(t *testing.T) {
	tests := []struct {
		name    string
		envType string
		wantErr bool
	}{
		{
			name:    "dev_environment",
			envType: "dev",
			wantErr: false,
		},
		{
			name:    "staging_environment",
			envType: "staging",
			wantErr: false,
		},
		{
			name:    "prod_environment",
			envType: "prod",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			provider, err := NewAWSProvider(ctx, "us-east-1", "")
			testutil.AssertNoError(t, err)

			costs, err := provider.CalculateTotalCost(tt.envType)
			if tt.wantErr {
				testutil.AssertError(t, err)
				return
			}

			testutil.AssertNoError(t, err)
			testutil.AssertNotEqual(t, nil, costs)
			testutil.AssertEqual(t, tt.envType, costs.Environment)
			testutil.AssertEqual(t, "aws", costs.Provider)
			testutil.AssertTrue(t, costs.TotalCost > 0, "Total cost should be positive")
			testutil.AssertTrue(t, costs.NetworkCost > 0, "Network cost should be positive")
			testutil.AssertTrue(t, costs.DatabaseCost > 0, "Database cost should be positive")
			testutil.AssertTrue(t, costs.K8sCost > 0, "K8s cost should be positive")
		})
	}
}

func TestAWSProvider_MultipleRegions(t *testing.T) {
	regions := []string{"us-east-1", "us-west-2", "eu-west-1"}

	for _, region := range regions {
		t.Run("region_"+region, func(t *testing.T) {
			ctx := context.Background()
			provider, err := NewAWSProvider(ctx, region, "")
			testutil.AssertNoError(t, err)
			testutil.AssertNotEqual(t, nil, provider)
			testutil.AssertEqual(t, region, provider.region)
		})
	}
}
