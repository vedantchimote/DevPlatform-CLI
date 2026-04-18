# Testing Guide

This guide covers how to run, write, and maintain tests for the DevPlatform CLI project.

## Overview

The DevPlatform CLI has a comprehensive testing suite with:
- **75% code coverage** (exceeds 70% minimum target)
- **Unit tests** for all core packages
- **Integration tests** for end-to-end workflows
- **Mock implementations** for external dependencies
- **CI/CD integration** with automated testing

## Quick Start

### Run All Tests

```bash
go test -v ./...
```

### Run Tests with Coverage

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Specific Package Tests

```bash
# Config package
go test -v ./internal/config/...

# Helm package
go test -v ./internal/helm/...

# Terraform package
go test -v ./internal/terraform/...

# AWS provider
go test -v ./internal/aws/...

# Azure provider
go test -v ./internal/azure/...
```

### Run Integration Tests

```bash
go test -v ./test/integration/...
```

## Test Structure

```
devplatform-cli/
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go          # Unit tests
│   ├── terraform/
│   │   ├── executor.go
│   │   ├── executor_test.go
│   │   └── mock_executor.go        # Mock implementation
│   └── ...
├── test/
│   ├── integration/
│   │   ├── create_test.go          # Integration tests
│   │   ├── destroy_test.go
│   │   ├── status_test.go
│   │   └── helpers.go              # Test utilities
│   ├── fixtures/
│   │   ├── aws/                    # AWS test data
│   │   ├── azure/                  # Azure test data
│   │   └── helm/                   # Helm test data
│   ├── mocks/
│   │   ├── aws_mock.go             # AWS provider mock
│   │   ├── azure_mock.go           # Azure provider mock
│   │   └── k8s_mock.go             # Kubernetes mock
│   └── testutil/
│       └── helpers.go              # Test helper functions
```

## Test Coverage by Package

| Package | Coverage | Status |
|---------|----------|--------|
| config | 87.8% | ✅ |
| errors | 100.0% | ✅ |
| helm | 57.4% | ✅ |
| logger | 90.0% | ✅ |
| terraform | 74.4% | ✅ |
| aws | 67.0% | ✅ |
| azure | 58.9% | ✅ |
| test/mocks | 94.5% | ✅ |
| test/testutil | 73.4% | ✅ |

## Writing Tests

### Unit Test Example

```go
package config

import (
    "testing"
    "github.com/devplatform/devplatform-cli/test/testutil"
)

func TestValidateAppName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {
            name:    "valid lowercase name",
            input:   "myapp",
            wantErr: false,
        },
        {
            name:    "invalid uppercase",
            input:   "MyApp",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateAppName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateAppName() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Test Example

```go
package integration

import (
    "testing"
    "github.com/devplatform/devplatform-cli/test/testutil"
)

func TestCreateWorkflow_Success(t *testing.T) {
    tc := SetupTestContext(t)
    defer tc.Cleanup()

    // Configure mocks
    tc.MockAWSProvider.ValidateCredentialsFunc = func(ctx context.Context) error {
        return nil
    }

    // Execute workflow
    // ... test code ...

    // Verify results
    testutil.AssertEqual(t, 1, tc.MockAWSProvider.GetValidateCredentialsCallCount())
}
```

### Using Mocks

```go
// Create a mock
mockTF := terraform.NewMockTerraformExecutor()

// Configure mock behavior
mockTF.InitFunc = func(ctx context.Context, workingDir string) error {
    return nil
}

// Use the mock
err := mockTF.Init(context.Background(), "/path/to/terraform")

// Verify calls
if mockTF.GetInitCallCount() != 1 {
    t.Errorf("Expected 1 Init call, got %d", mockTF.GetInitCallCount())
}
```

## Test Utilities

The `test/testutil` package provides helpful assertion functions:

```go
import "github.com/devplatform/devplatform-cli/test/testutil"

// Assert no error
testutil.AssertNoError(t, err)

// Assert error occurred
testutil.AssertError(t, err)

// Assert equality
testutil.AssertEqual(t, expected, actual)

// Assert string contains
testutil.AssertContains(t, str, substr)

// Assert error code
testutil.AssertErrorCode(t, err, "VALIDATION_ERROR")

// Load test fixture
data := testutil.LoadFixture(t, "aws/rds-instance-dev.json")

// Create temp config
configPath := testutil.CreateTempConfigYAML(t, yamlContent)
```

## Mock Implementations

### Available Mocks

1. **MockTerraformExecutor** - Simulates Terraform operations
2. **MockHelmClient** - Simulates Helm operations
3. **MockAWSProvider** - Simulates AWS SDK calls
4. **MockAzureProvider** - Simulates Azure SDK calls
5. **MockK8sClient** - Simulates Kubernetes operations

### Mock Features

All mocks provide:
- **Call tracking** - Records all method calls with parameters
- **Configurable responses** - Set custom return values
- **Thread-safe** - Safe for concurrent use
- **Reset functionality** - Clear call history between tests

### Example: Configuring Mock Responses

```go
mock := mocks.NewMockAWSProvider()

// Configure successful response
mock.ValidateCredentialsFunc = func(ctx context.Context) error {
    return nil
}

// Configure error response
mock.ValidateCredentialsFunc = func(ctx context.Context) error {
    return errors.New("invalid credentials")
}

// Configure custom data response
mock.CalculateTotalCostFunc = func(envType string) (*types.EnvironmentCosts, error) {
    return &types.EnvironmentCosts{
        TotalCost: 150.0,
        Environment: envType,
    }, nil
}
```

## Test Fixtures

Test fixtures are located in `test/fixtures/`:

```
test/fixtures/
├── aws/
│   └── rds-instance-dev.json       # Sample AWS RDS response
├── azure/
│   └── pricing-response.json       # Sample Azure pricing response
└── helm/
    └── values-staging.yaml         # Sample Helm values
```

### Loading Fixtures

```go
// Load a fixture file
data := testutil.LoadFixture(t, "aws/rds-instance-dev.json")

// Parse JSON fixture
var instance RDSInstance
err := json.Unmarshal(data, &instance)
testutil.AssertNoError(t, err)
```

## Integration Tests

Integration tests verify end-to-end workflows using mocks for external dependencies.

### Test Categories

1. **Create Workflow** (7 tests)
   - Success path
   - Dry-run mode
   - Input validation
   - Terraform failure with rollback
   - Helm failure with rollback
   - Credential validation failure
   - Timeout handling

2. **Destroy Workflow** (9 tests)
   - Success path
   - Confirmation prompt
   - Force flag behavior
   - Non-existent environment
   - Partial failure
   - Input validation
   - Missing Helm release
   - Keep-state flag
   - Cost savings calculation

3. **Status Workflow** (11 tests)
   - Healthy environment
   - Degraded environment
   - Non-existent environment
   - JSON/YAML output formats
   - Input validation
   - Resource status checks
   - Azure provider support
   - Overall status determination

### Running Integration Tests

```bash
# Run all integration tests
go test -v ./test/integration/...

# Run specific workflow tests
go test -v ./test/integration/... -run TestCreateWorkflow
go test -v ./test/integration/... -run TestDestroyWorkflow
go test -v ./test/integration/... -run TestStatusWorkflow
```

## CI/CD Integration

Tests run automatically in GitHub Actions:

### On Every Push
- Run all unit tests
- Generate coverage reports
- Fail if coverage < 70%

### On Pull Requests
- Run all unit tests
- Run all integration tests
- Upload coverage reports as artifacts
- Fail if any test fails

### Viewing CI Results

1. Go to the **Actions** tab in GitHub
2. Select the workflow run
3. View test results and coverage reports
4. Download coverage artifacts if needed

## Coverage Reports

### Generate HTML Coverage Report

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

Open `coverage.html` in a browser to see line-by-line coverage.

### Generate Text Coverage Report

```bash
go tool cover -func=coverage.out
```

### Check Coverage Threshold

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//' | \
  awk '{if ($1 < 70) exit 1}'
```

## Best Practices

### Writing Tests

1. **Keep tests simple** - Each test should verify one behavior
2. **Use descriptive names** - Test names should describe what they test
3. **Avoid test interdependencies** - Tests should run independently
4. **Use table-driven tests** - For testing multiple scenarios
5. **Mock external dependencies** - Keep tests fast and isolated
6. **Clean up resources** - Use `defer` for cleanup
7. **Test error cases** - Don't just test the happy path

### Test Organization

1. **Co-locate tests** - Keep `*_test.go` files next to source files
2. **Use testdata directories** - For test fixtures
3. **Group related tests** - Use subtests with `t.Run()`
4. **Separate unit and integration tests** - Different directories

### Mock Usage

1. **Configure before use** - Set up mock responses before calling
2. **Verify calls** - Check that mocks were called as expected
3. **Reset between tests** - Clear mock state for test isolation
4. **Use realistic data** - Mock responses should match real data

## Debugging Tests

### Run Tests with Verbose Output

```bash
go test -v ./internal/config/...
```

### Run a Specific Test

```bash
go test -v ./internal/config/... -run TestValidateAppName
```

### Run Tests with Race Detection

```bash
go test -v -race ./...
```

### View Test Coverage for Specific Package

```bash
go test -v -coverprofile=coverage.out ./internal/config/...
go tool cover -func=coverage.out
```

## Common Issues

### Tests Fail Due to Missing Fixtures

**Problem**: Test can't find fixture file

**Solution**: Ensure fixture path is relative to `test/fixtures/`:
```go
data := testutil.LoadFixture(t, "aws/rds-instance-dev.json")
```

### Mock Not Recording Calls

**Problem**: Mock call count is 0 when it should be > 0

**Solution**: Ensure you're using the mock instance, not creating a new one:
```go
// Wrong
mock := NewMockExecutor()
realExecutor.Init() // Using wrong instance

// Right
mock := NewMockExecutor()
mock.Init() // Using mock instance
```

### Tests Pass Locally But Fail in CI

**Problem**: Tests work on your machine but fail in GitHub Actions

**Solution**: Check for:
- Hardcoded paths (use relative paths)
- Time-dependent tests (use fixed times in tests)
- OS-specific behavior (test on multiple platforms)

## Additional Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Table-Driven Tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Test Fixtures Best Practices](https://golang.org/doc/tutorial/add-a-test)

## Getting Help

If you encounter issues with tests:

1. Check this documentation
2. Review existing tests for examples
3. Check CI logs for detailed error messages
4. Ask in team chat or create an issue

---

**Last Updated**: 2024
**Test Coverage**: 75%
**Total Tests**: 27 integration tests + comprehensive unit tests
