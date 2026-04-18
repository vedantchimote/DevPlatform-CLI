package integration

import (
	"context"
	"testing"

	"github.com/devplatform/devplatform-cli/internal/config"
	"github.com/devplatform/devplatform-cli/internal/helm"
	"github.com/devplatform/devplatform-cli/internal/logger"
	"github.com/devplatform/devplatform-cli/internal/provider"
	"github.com/devplatform/devplatform-cli/internal/terraform"
	"github.com/devplatform/devplatform-cli/test/mocks"
)

// TestContext holds common test dependencies
type TestContext struct {
	Ctx              context.Context
	Logger           *logger.Logger
	MockAWSProvider  *mocks.MockAWSProvider
	MockAzureProvider *mocks.MockAzureProvider
	MockTerraform    *terraform.MockTerraformExecutor
	MockHelm         *helm.MockHelmClient
	Config           *config.Config
}

// SetupTestContext creates a new test context with mocked dependencies
func SetupTestContext(t *testing.T) *TestContext {
	t.Helper()

	// Create logger with test mode
	log := logger.New(logger.InfoLevel, false)

	// Create mock providers
	mockAWS := mocks.NewMockAWSProvider()
	mockAzure := mocks.NewMockAzureProvider()

	// Create mock Terraform executor
	mockTF := terraform.NewMockTerraformExecutor()

	// Create mock Helm client
	mockHelm := helm.NewMockHelmClient()

	// Create default config
	cfg := &config.Config{
		AWS: config.AWSConfig{
			Region:  "us-east-1",
			Profile: "default",
		},
		Azure: config.AzureConfig{
			SubscriptionID: "12345678-1234-1234-1234-123456789012",
			TenantID:       "87654321-4321-4321-4321-210987654321",
			Location:       "eastus",
		},
		Helm: config.HelmConfig{
			ChartPath: "charts/devplatform-base",
		},
	}

	return &TestContext{
		Ctx:              context.Background(),
		Logger:           log,
		MockAWSProvider:  mockAWS,
		MockAzureProvider: mockAzure,
		MockTerraform:    mockTF,
		MockHelm:         mockHelm,
		Config:           cfg,
	}
}

// SetupMockProvider configures a mock cloud provider with default successful behavior
func SetupMockProvider(tc *TestContext, providerType string) provider.CloudProvider {
	if providerType == "aws" {
		return tc.MockAWSProvider
	}
	return tc.MockAzureProvider
}

// Cleanup performs cleanup after tests
func (tc *TestContext) Cleanup() {
	// Reset all mocks
	tc.MockAWSProvider.Reset()
	tc.MockAzureProvider.Reset()
	tc.MockTerraform.Reset()
	tc.MockHelm.Reset()
}
