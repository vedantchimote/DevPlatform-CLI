# Integration Tests

This directory contains integration tests for the DevPlatform CLI. Integration tests verify end-to-end workflows by testing multiple components together.

## Test Structure

- `create_test.go` - Tests for the create workflow
- `destroy_test.go` - Tests for the destroy workflow
- `status_test.go` - Tests for the status workflow
- `helpers.go` - Shared test helpers and utilities

## Test Coverage

### Create Workflow Tests
- `TestCreateWorkflow_Success` - Successful create workflow with all components
- `TestCreateWorkflow_DryRun` - Dry-run mode (plan only, no apply)
- `TestCreateWorkflow_InvalidInputs` - Input validation (app name, environment, provider)
- `TestCreateWorkflow_TerraformFailure` - Rollback when Terraform fails
- `TestCreateWorkflow_HelmFailure` - Rollback when Helm deployment fails
- `TestCreateWorkflow_CredentialValidationFailure` - Invalid cloud credentials
- `TestCreateWorkflow_Timeout` - Timeout handling

### Destroy Workflow Tests
- `TestDestroyWorkflow_Success` - Successful destroy workflow
- `TestDestroyWorkflow_WithConfirmation` - Confirmation prompt handling
- `TestDestroyWorkflow_WithForceFlag` - Force flag continues despite errors
- `TestDestroyWorkflow_EnvironmentNotFound` - Non-existent environment handling
- `TestDestroyWorkflow_PartialFailure` - Partial destruction with errors
- `TestDestroyWorkflow_InvalidInputs` - Input validation
- `TestDestroyWorkflow_HelmReleaseNotFound` - Missing Helm release handling
- `TestDestroyWorkflow_KeepStateFlag` - Keep state file after destroy
- `TestDestroyWorkflow_CostSavingsCalculation` - Cost savings calculation

### Status Workflow Tests
- `TestStatusWorkflow_HealthyEnvironment` - All components healthy
- `TestStatusWorkflow_DegradedEnvironment` - Some components degraded
- `TestStatusWorkflow_NonExistentEnvironment` - Environment not found
- `TestStatusWorkflow_JSONOutput` - JSON output format
- `TestStatusWorkflow_YAMLOutput` - YAML output format
- `TestStatusWorkflow_InvalidInputs` - Input validation
- `TestStatusWorkflow_NetworkStatus` - Network resource status
- `TestStatusWorkflow_DatabaseStatus` - Database resource status
- `TestStatusWorkflow_NamespaceStatus` - Kubernetes namespace status
- `TestStatusWorkflow_AzureProvider` - Azure-specific status checks
- `TestStatusWorkflow_OverallStatusDetermination` - Overall status calculation

## Running Integration Tests

```bash
# Run all integration tests
go test -v ./test/integration/...

# Run specific test
go test -v ./test/integration/... -run TestCreateWorkflow

# Run with coverage
go test -v -coverprofile=coverage.out ./test/integration/...

# Run specific workflow tests
go test -v ./test/integration/... -run TestCreateWorkflow
go test -v ./test/integration/... -run TestDestroyWorkflow
go test -v ./test/integration/... -run TestStatusWorkflow
```

## Test Characteristics

- Use mocks for external dependencies (AWS SDK, Azure SDK, Terraform, Helm, Kubernetes)
- Test complete workflows from start to finish
- Verify error handling and rollback scenarios
- Clean up all test resources after execution
- Run in isolation without affecting real infrastructure
- Fast execution (no real cloud resources created)
- Deterministic results (no flaky tests)

## Test Fixtures

Test fixtures are located in `test/fixtures/` and include:
- Sample configuration files
- Mock Terraform state files
- Mock Helm values files
- Mock cloud provider responses

## Mock Configuration

The integration tests use the following mocks:

### Cloud Provider Mocks
- `MockAWSProvider` - AWS SDK operations (credentials, identity, costs)
- `MockAzureProvider` - Azure SDK operations (credentials, identity, costs)

### Infrastructure Mocks
- `MockTerraformExecutor` - Terraform operations (init, plan, apply, destroy)
- `MockHelmClient` - Helm operations (install, upgrade, uninstall, status)

### Test Context
The `TestContext` struct provides a unified setup for all integration tests:
- Pre-configured mocks with default behavior
- Logger for test output
- Configuration with sensible defaults
- Cleanup methods to reset state between tests

## Writing New Integration Tests

1. Create a new test function following the naming convention `TestWorkflowName_Scenario`
2. Use `SetupTestContext(t)` to get a configured test context
3. Configure mock behavior for your specific test scenario
4. Execute the workflow (or verify mock configuration)
5. Assert expected outcomes using `testutil` helpers
6. Defer `tc.Cleanup()` to reset mocks after the test

Example:
```go
func TestCreateWorkflow_CustomScenario(t *testing.T) {
    tc := SetupTestContext(t)
    defer tc.Cleanup()

    // Configure mocks
    tc.MockAWSProvider.ValidateCredentialsFunc = func(ctx context.Context) error {
        return nil
    }

    // Execute workflow and assert
    testutil.AssertEqual(t, 0, tc.MockAWSProvider.GetValidateCredentialsCallCount())
}
```
