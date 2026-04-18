# Test Fixtures

This directory contains test fixtures used across the DevPlatform CLI test suite. Fixtures provide realistic sample data for testing without requiring actual cloud resources or external services.

## Directory Structure

```
test/fixtures/
├── configs/        # Configuration file fixtures
├── terraform/      # Terraform state and variable fixtures
├── helm/           # Helm values and release status fixtures
├── aws/            # AWS API response fixtures
└── azure/          # Azure API response fixtures
```

## Configuration Fixtures (`configs/`)

Sample configuration files for different environments and cloud providers.

### Valid Configurations

- `valid-aws-dev.yaml` - AWS development environment configuration
- `valid-aws-staging.yaml` - AWS staging environment configuration
- `valid-aws-prod.yaml` - AWS production environment configuration
- `valid-azure-dev.yaml` - Azure development environment configuration
- `valid-azure-staging.yaml` - Azure staging environment configuration
- `valid-azure-prod.yaml` - Azure production environment configuration

### Invalid Configurations (for error testing)

- `invalid-missing-provider.yaml` - Missing cloud_provider field
- `invalid-bad-provider.yaml` - Unsupported cloud provider value
- `invalid-malformed.yaml` - Malformed YAML syntax
- `invalid-bad-cidr.yaml` - Invalid CIDR block format

## Terraform Fixtures (`terraform/`)

Terraform variable files, state files, and output files for testing Terraform operations.

### Variable Files

- `aws-network-dev.tfvars` - AWS network module variables for dev
- `aws-network-prod.tfvars` - AWS network module variables for prod
- `aws-database-dev.tfvars` - AWS database module variables for dev

### State Files

- `state-aws-dev.json` - Sample Terraform state for AWS dev environment
- `state-azure-prod.json` - Sample Terraform state for Azure prod environment

### Output Files

- `output-aws-dev.json` - Sample Terraform output for AWS dev environment

### Invalid Files (for error testing)

- `invalid-bad-app-name.tfvars` - Invalid app name (uppercase characters)
- `invalid-bad-env.tfvars` - Invalid environment type

## Helm Fixtures (`helm/`)

Helm values files and release status responses for testing Helm operations.

### Values Files

- `values-dev.yaml` - Helm values for dev environment
- `values-staging.yaml` - Helm values for staging environment
- `values-prod.yaml` - Helm values for prod environment

### Release Status Files

- `release-status-deployed.json` - Successful deployment status
- `release-status-failed.json` - Failed deployment status

### Invalid Files (for error testing)

- `invalid-malformed.yaml` - Malformed YAML syntax
- `invalid-missing-required.yaml` - Missing required fields

## AWS Fixtures (`aws/`)

Sample AWS API responses for testing AWS provider operations.

### Resource Responses

- `eks-cluster-dev.json` - EKS cluster description response
- `rds-instance-dev.json` - RDS instance description response
- `pricing-response.json` - EC2 pricing API response
- `credentials-valid.json` - Valid STS GetCallerIdentity response

### Error Responses

- `error-invalid-credentials.json` - Invalid credentials error
- `error-access-denied.json` - Access denied error

## Azure Fixtures (`azure/`)

Sample Azure API responses for testing Azure provider operations.

### Resource Responses

- `aks-cluster-prod.json` - AKS cluster description response
- `postgres-server-prod.json` - PostgreSQL flexible server response
- `pricing-response.json` - Azure pricing API response
- `credentials-valid.json` - Valid subscription details response

### Error Responses

- `error-invalid-credentials.json` - Authentication failed error
- `error-subscription-not-found.json` - Subscription not found error
- `error-resource-not-found.json` - Resource not found error

## Usage in Tests

### Loading Configuration Fixtures

```go
import "github.com/yourusername/devplatform-cli/test/testutil"

func TestConfigLoading(t *testing.T) {
    data := testutil.LoadFixture(t, "test/fixtures/configs/valid-aws-dev.yaml")
    // Use data in your test
}
```

### Loading JSON Fixtures

```go
import (
    "encoding/json"
    "github.com/yourusername/devplatform-cli/test/testutil"
)

func TestAWSResponse(t *testing.T) {
    data := testutil.LoadFixture(t, "test/fixtures/aws/eks-cluster-dev.json")
    
    var cluster EKSCluster
    err := json.Unmarshal(data, &cluster)
    testutil.AssertNoError(t, err)
    
    // Use cluster in your test
}
```

### Using Fixtures with Mocks

```go
func TestTerraformOutput(t *testing.T) {
    // Load fixture
    outputData := testutil.LoadFixture(t, "test/fixtures/terraform/output-aws-dev.json")
    
    // Configure mock to return fixture data
    mockExecutor := &terraform.MockTerraformExecutor{
        OutputFunc: func(ctx context.Context, workingDir, outputName string) (string, error) {
            return string(outputData), nil
        },
    }
    
    // Test code that uses the mock
}
```

## Fixture Maintenance

### Adding New Fixtures

1. Create the fixture file in the appropriate subdirectory
2. Ensure the fixture is properly formatted (valid YAML/JSON)
3. Include realistic data that matches actual API responses
4. Add documentation to this README
5. Create corresponding test cases that use the fixture

### Updating Fixtures

When updating fixtures to match new API versions or schema changes:

1. Update the fixture file
2. Run all tests to ensure no breakage
3. Update this README if the fixture structure changed
4. Document any breaking changes in test code

### Fixture Validation

All fixtures should be validated for:

- **Syntax**: Valid YAML/JSON format
- **Schema**: Matches expected structure
- **Realism**: Represents actual API responses
- **Completeness**: Includes all required fields
- **Variety**: Covers success and error cases

## Best Practices

1. **Keep fixtures realistic**: Use data that closely matches actual API responses
2. **Include both valid and invalid examples**: Test both success and error paths
3. **Use consistent naming**: Follow the pattern `<resource>-<environment>.json`
4. **Document fixture purpose**: Add comments explaining what each fixture tests
5. **Avoid sensitive data**: Never include real credentials, API keys, or PII
6. **Version fixtures**: Update fixtures when API versions change
7. **Test fixture loading**: Ensure fixtures can be loaded without errors

## Related Documentation

- [Test Utilities](../testutil/README.md) - Helper functions for loading fixtures
- [Mock Framework](../mocks/README.md) - Using mocks with fixtures
- [Testing Guide](../../docs/testing/README.md) - Overall testing strategy
