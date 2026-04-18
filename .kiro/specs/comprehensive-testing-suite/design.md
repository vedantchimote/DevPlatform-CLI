# Design Document: Comprehensive Testing Suite

## Overview

This design document outlines the technical implementation for a comprehensive testing suite for the DevPlatform CLI project. The suite will include unit tests, integration tests, mock implementations, test fixtures, coverage reporting, and CI integration to achieve production-ready quality standards.

The testing suite addresses the current zero test coverage by implementing a structured approach using Go's standard testing package, table-driven test patterns, and mock-based isolation. The design ensures tests are maintainable, fast, and provide meaningful feedback to developers.

### Goals

- Achieve minimum 70% code coverage across all packages
- Enable fast, isolated unit testing through comprehensive mocking
- Provide integration tests for end-to-end workflow validation
- Integrate automated testing into CI pipeline
- Establish testing patterns and documentation for future development

### Non-Goals

- Performance benchmarking (separate effort)
- Load testing or stress testing
- UI/UX testing (CLI is text-based)
- Security penetration testing (separate security audit)

## Architecture

### Testing Layers

The testing architecture follows a three-layer approach:

```
┌─────────────────────────────────────────┐
│         Integration Tests               │
│  (End-to-end workflows with mocks)      │
└─────────────────────────────────────────┘
                  │
┌─────────────────────────────────────────┐
│          Unit Tests                     │
│  (Individual functions with mocks)      │
└─────────────────────────────────────────┘
                  │
┌─────────────────────────────────────────┐
│          Mock Layer                     │
│  (Simulated external dependencies)      │
└─────────────────────────────────────────┘
```

### Directory Structure

```
devplatform-cli/
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   ├── config_test.go          # Unit tests
│   │   └── testdata/               # Test fixtures
│   ├── terraform/
│   │   ├── executor.go
│   │   ├── executor_test.go
│   │   └── mock_executor.go        # Mock implementation
│   ├── helm/
│   │   ├── client.go
│   │   ├── client_test.go
│   │   └── mock_client.go
│   └── ...
├── test/
│   ├── integration/
│   │   ├── create_test.go          # Integration tests
│   │   ├── destroy_test.go
│   │   └── status_test.go
│   ├── fixtures/
│   │   ├── configs/                # Sample configs
│   │   ├── terraform/              # Sample TF state
│   │   └── helm/                   # Sample Helm values
│   └── mocks/
│       ├── aws_mock.go             # Shared mocks
│       ├── azure_mock.go
│       └── k8s_mock.go
└── docs/
    └── testing/
        └── README.md               # Testing documentation
```

### Test Execution Flow

1. **Unit Tests**: Run first, fast execution (< 5 seconds total)
2. **Integration Tests**: Run on PR, slower execution (< 2 minutes)
3. **Coverage Analysis**: Aggregate coverage from both test types
4. **CI Reporting**: Upload results and fail on threshold violations

## Components and Interfaces

### Mock Framework

We'll use Go's interface-based mocking approach without external frameworks to maintain simplicity and zero additional dependencies.

#### Mock Terraform Executor

```go
// MockTerraformExecutor simulates Terraform operations
type MockTerraformExecutor struct {
    InitFunc    func(ctx context.Context, workingDir string) error
    PlanFunc    func(ctx context.Context, workingDir string, varFile string) (string, error)
    ApplyFunc   func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error
    DestroyFunc func(ctx context.Context, workingDir string, varFile string, autoApprove bool) error
    OutputFunc  func(ctx context.Context, workingDir string, outputName string) (string, error)
    
    // Call tracking
    InitCalls    []MockCall
    PlanCalls    []MockCall
    ApplyCalls   []MockCall
    DestroyCalls []MockCall
    OutputCalls  []MockCall
}

type MockCall struct {
    Args      []interface{}
    Timestamp time.Time
}
```

#### Mock Helm Client

```go
// MockHelmClient simulates Helm operations
type MockHelmClient struct {
    InstallFunc   func(ctx context.Context, opts InstallOptions) error
    UpgradeFunc   func(ctx context.Context, opts UpgradeOptions) error
    UninstallFunc func(ctx context.Context, opts UninstallOptions) error
    StatusFunc    func(ctx context.Context, releaseName, namespace string) (*ReleaseStatus, error)
    ListFunc      func(ctx context.Context, namespace string) ([]*Release, error)
    
    // Call tracking
    InstallCalls   []MockCall
    UpgradeCalls   []MockCall
    UninstallCalls []MockCall
    StatusCalls    []MockCall
    ListCalls      []MockCall
}
```

#### Mock Cloud Provider Clients

```go
// MockAWSProvider simulates AWS SDK operations
type MockAWSProvider struct {
    ValidateCredentialsFunc func(ctx context.Context) error
    GetPricingFunc          func(ctx context.Context, instanceType string) (float64, error)
    GetKubeconfigFunc       func(ctx context.Context, clusterName string) ([]byte, error)
    
    // Call tracking
    ValidateCredentialsCalls []MockCall
    GetPricingCalls          []MockCall
    GetKubeconfigCalls       []MockCall
}

// MockAzureProvider simulates Azure SDK operations
type MockAzureProvider struct {
    ValidateCredentialsFunc func(ctx context.Context) error
    GetPricingFunc          func(ctx context.Context, vmSize string) (float64, error)
    GetKubeconfigFunc       func(ctx context.Context, clusterName string) ([]byte, error)
    
    // Call tracking
    ValidateCredentialsCalls []MockCall
    GetPricingCalls          []MockCall
    GetKubeconfigCalls       []MockCall
}
```

### Test Helper Functions

```go
// Test helpers for common operations
package testutil

// LoadFixture loads a test fixture file
func LoadFixture(t *testing.T, path string) []byte

// CreateTempConfig creates a temporary config file for testing
func CreateTempConfig(t *testing.T, config *Config) string

// AssertNoError fails the test if error is not nil
func AssertNoError(t *testing.T, err error)

// AssertError fails the test if error is nil
func AssertError(t *testing.T, err error)

// AssertEqual fails the test if values are not equal
func AssertEqual(t *testing.T, expected, actual interface{})

// AssertContains fails the test if string doesn't contain substring
func AssertContains(t *testing.T, str, substr string)
```

### Table-Driven Test Structure

```go
// Standard table-driven test pattern
func TestValidateAppName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
        errCode string
    }{
        {
            name:    "valid lowercase name",
            input:   "myapp",
            wantErr: false,
        },
        {
            name:    "valid with hyphens",
            input:   "my-app-123",
            wantErr: false,
        },
        {
            name:    "invalid uppercase",
            input:   "MyApp",
            wantErr: true,
            errCode: "INVALID_APP_NAME",
        },
        {
            name:    "invalid too short",
            input:   "a",
            wantErr: true,
            errCode: "INVALID_APP_NAME",
        },
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateAppName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateAppName() error = %v, wantErr %v", err, tt.wantErr)
            }
            if err != nil && tt.errCode != "" {
                // Verify error code matches
            }
        })
    }
}
```

## Data Models

### Test Fixture Models

```go
// TestConfig represents a test configuration fixture
type TestConfig struct {
    Name        string
    ConfigFile  string
    Environment string
    Provider    string
    Expected    *config.Config
}

// TestTerraformState represents a test Terraform state fixture
type TestTerraformState struct {
    Name     string
    StateFile string
    Outputs  map[string]string
}

// TestHelmRelease represents a test Helm release fixture
type TestHelmRelease struct {
    Name        string
    Namespace   string
    Chart       string
    ValuesFile  string
    Status      string
}
```

### Coverage Report Model

```go
// CoverageReport represents test coverage metrics
type CoverageReport struct {
    OverallCoverage float64
    PackageCoverage map[string]PackageCoverage
    Threshold       float64
    Passed          bool
}

// PackageCoverage represents coverage for a single package
type PackageCoverage struct {
    Package        string
    Coverage       float64
    TotalLines     int
    CoveredLines   int
    UncoveredLines []int
}
```


## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

Most of the requirements for this testing suite are structural in nature (e.g., "tests SHALL exist for package X"), which are verified through code review and manual inspection rather than automated property-based testing. However, several behavioral properties can be validated to ensure the testing infrastructure works correctly.

### Property 1: Unit Test Isolation

*For any* unit test in the test suite, when executed, it SHALL complete within 100 milliseconds, indicating proper isolation through mocks without making real external calls.

**Validates: Requirements 1.10**

### Property 2: Mock Call Recording

*For any* mock implementation, when a method is called with specific parameters, those parameters SHALL be recorded in the mock's call history for later verification.

**Validates: Requirements 2.7**

### Property 3: Mock Response Configuration

*For any* mock implementation, when configured with a predefined response, it SHALL return that exact response when the corresponding method is called.

**Validates: Requirements 2.8**

### Property 4: Table-Driven Test Error Reporting

*For any* table-driven test that fails, the error message SHALL include the specific test case name that failed, enabling quick identification of the failing scenario.

**Validates: Requirements 3.6**

### Property 5: Integration Test Isolation

*For any* integration test, it SHALL use unique resource identifiers (names, namespaces, directories) to prevent conflicts when tests run in parallel.

**Validates: Requirements 4.7**

### Property 6: Integration Test Cleanup

*For any* integration test, when it completes (whether success or failure), all test resources (files, directories, mock state) SHALL be cleaned up automatically.

**Validates: Requirements 4.8**

### Property 7: Coverage Report Completeness

*For any* test run with coverage enabled, the generated coverage report SHALL contain both overall coverage percentage and per-package coverage percentages.

**Validates: Requirements 5.3, 5.4**

### Property 8: Coverage Threshold Enforcement

*For any* test run where coverage falls below 70%, the test command SHALL exit with a non-zero status code, causing CI builds to fail.

**Validates: Requirements 5.5, 6.5**

### Property 9: Coverage Aggregation

*For any* coverage report generated, it SHALL include coverage data from both unit tests and integration tests, providing a complete picture of code coverage.

**Validates: Requirements 5.7**

### Property 10: CI Test Failure Propagation

*For any* test execution in CI where at least one test fails, the CI pipeline SHALL exit with a non-zero status code and display the test name and failure reason.

**Validates: Requirements 6.4, 6.8**

### Property 11: Fixture Validation

*For any* test fixture file, when loaded by the test suite, it SHALL be validated for correct format (valid YAML/JSON syntax, required fields present), failing fast with a clear error if corrupted.

**Validates: Requirements 7.8**

### Property 12: Error Test Completeness

*For any* test that validates error handling, it SHALL verify both the error code and error message, ensuring complete error validation.

**Validates: Requirements 8.9, 8.10**

### Property 13: Timeout Test Efficiency

*For any* test that validates timeout behavior, it SHALL use mock implementations to simulate delays and complete within 1 second, avoiding real timeout waits.

**Validates: Requirements 9.7**

## Error Handling

### Test Execution Errors

The test suite must handle various error conditions gracefully:

1. **Missing Dependencies**: If required tools (terraform, helm, kubectl) are not installed, tests should skip gracefully with clear messages rather than fail cryptically.

2. **Fixture Loading Errors**: If test fixtures are missing or corrupted, tests should fail immediately with clear indication of which fixture failed and why.

3. **Mock Configuration Errors**: If mocks are not properly configured before use, tests should fail with helpful error messages indicating what configuration is missing.

4. **Timeout Errors**: If tests exceed reasonable time limits, they should be terminated with timeout errors rather than hanging indefinitely.

### Error Reporting Standards

All test errors should follow these standards:

```go
// Good error message
t.Errorf("ValidateAppName(%q) failed: got error %v, want nil", input, err)

// Bad error message
t.Errorf("test failed")
```

Error messages should include:
- The function/method being tested
- The input values that caused the failure
- The actual result received
- The expected result

### CI Error Handling

The CI pipeline should:
- Capture and display all test output
- Highlight failed tests in the summary
- Upload test results as artifacts for debugging
- Provide links to coverage reports
- Exit with appropriate status codes (0 for success, non-zero for failure)

## Testing Strategy

### Unit Testing Approach

Unit tests will be co-located with source files following Go conventions:

```
internal/config/
├── config.go
├── config_test.go      # Unit tests for config.go
├── loader.go
├── loader_test.go      # Unit tests for loader.go
└── testdata/           # Test fixtures
    ├── valid_config.yaml
    └── invalid_config.yaml
```

**Unit Test Characteristics**:
- Fast execution (< 100ms per test)
- No external dependencies (use mocks)
- Test single functions/methods in isolation
- Use table-driven tests for multiple scenarios
- Minimum 100 iterations for property-based tests

**Unit Test Coverage Goals**:
- Config package: 80%+ coverage
- Logger package: 70%+ coverage
- Error package: 90%+ coverage
- Terraform package: 75%+ coverage
- Helm package: 75%+ coverage
- Provider packages: 70%+ coverage

### Integration Testing Approach

Integration tests will be in a separate `test/integration` directory:

```
test/
├── integration/
│   ├── create_workflow_test.go
│   ├── destroy_workflow_test.go
│   └── status_workflow_test.go
├── fixtures/
│   ├── configs/
│   ├── terraform/
│   └── helm/
└── mocks/
    ├── aws_mock.go
    ├── azure_mock.go
    └── k8s_mock.go
```

**Integration Test Characteristics**:
- Slower execution (< 2 minutes total)
- Test multiple components together
- Use mocks for external services (AWS, Azure, K8s)
- Use real Terraform/Helm with test fixtures
- Clean up all resources after each test

**Integration Test Scenarios**:
1. Create workflow: config load → validation → terraform init/plan/apply → helm install
2. Destroy workflow: config load → terraform destroy → helm uninstall
3. Status workflow: config load → terraform output → helm status → pod verification
4. Error scenarios: invalid config, terraform failure, helm failure

### Property-Based Testing Configuration

We'll use Go's standard testing package with custom property test helpers:

```go
// Property test helper
func PropertyTest(t *testing.T, name string, iterations int, testFunc func(t *testing.T)) {
    for i := 0; i < iterations; i++ {
        t.Run(fmt.Sprintf("%s_iteration_%d", name, i), testFunc)
    }
}

// Usage
func TestMockRecordsCallsProperty(t *testing.T) {
    PropertyTest(t, "mock_records_calls", 100, func(t *testing.T) {
        // Feature: comprehensive-testing-suite, Property 2: Mock Call Recording
        mock := NewMockTerraformExecutor()
        
        // Generate random test data
        workingDir := generateRandomPath()
        
        // Call mock
        _ = mock.Init(context.Background(), workingDir)
        
        // Verify call was recorded
        if len(mock.InitCalls) != 1 {
            t.Errorf("Expected 1 recorded call, got %d", len(mock.InitCalls))
        }
        if mock.InitCalls[0].Args[0] != workingDir {
            t.Errorf("Expected workingDir %s, got %v", workingDir, mock.InitCalls[0].Args[0])
        }
    })
}
```

**Property Test Configuration**:
- Minimum 100 iterations per property test
- Each property test tagged with: `Feature: comprehensive-testing-suite, Property N: <description>`
- Use randomized inputs where applicable
- Verify invariants hold across all iterations

### Coverage Reporting

Coverage will be measured using Go's built-in coverage tools:

```bash
# Run tests with coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Generate text summary
go tool cover -func=coverage.out

# Check coverage threshold
go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//' | \
  awk '{if ($1 < 70) exit 1}'
```

**Coverage Metrics**:
- Overall coverage: minimum 70%
- Per-package coverage: reported in CI artifacts
- Coverage trends: tracked over time in CI
- Uncovered lines: highlighted in HTML reports

### CI Integration

The GitHub Actions workflow will:

1. **On every push**:
   - Run all unit tests
   - Generate coverage reports
   - Upload coverage to Codecov (optional)
   - Cache Go modules for speed

2. **On pull requests to main**:
   - Run all unit tests
   - Run all integration tests
   - Generate coverage reports
   - Fail if coverage < 70%
   - Fail if any test fails
   - Upload test results as artifacts

3. **Parallel execution**:
   - Run tests across multiple packages in parallel
   - Use `go test -p 4` for parallel execution
   - Timeout after 10 minutes to prevent hanging

**CI Workflow Structure**:

```yaml
jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: go test -v -race -coverprofile=coverage.out ./...
      - run: go tool cover -func=coverage.out
      - run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$coverage < 70" | bc -l) )); then
            echo "Coverage $coverage% is below 70% threshold"
            exit 1
          fi
      - uses: actions/upload-artifact@v4
        with:
          name: coverage-report
          path: coverage.out

  integration-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - uses: hashicorp/setup-terraform@v3
      - uses: azure/setup-helm@v4
      - run: go test -v -tags=integration ./test/integration/...
```

### Test Documentation

A comprehensive testing README will be created at `docs/testing/README.md`:

**Contents**:
1. Running all tests: `go test ./...`
2. Running specific packages: `go test ./internal/config/...`
3. Running with coverage: `go test -coverprofile=coverage.out ./...`
4. Running integration tests: `go test -tags=integration ./test/integration/...`
5. Viewing coverage reports: `go tool cover -html=coverage.out`
6. Writing new tests: examples and patterns
7. Mock framework usage: creating and configuring mocks
8. Table-driven test patterns: examples
9. Property-based test patterns: examples
10. Debugging test failures: common issues and solutions

**Example Code Snippets**:

The documentation will include copy-paste examples for:
- Basic unit test structure
- Table-driven test pattern
- Mock creation and configuration
- Integration test setup and teardown
- Property-based test implementation
- Fixture loading and validation
- Error handling test patterns
- Timeout and cancellation tests

### Test Maintenance

**Guidelines for maintaining tests**:

1. **Keep tests simple**: Each test should verify one behavior
2. **Use descriptive names**: Test names should describe what they test
3. **Avoid test interdependencies**: Tests should run independently
4. **Update tests with code changes**: Tests are first-class code
5. **Review test coverage**: Regularly check for untested code paths
6. **Refactor tests**: Apply same quality standards as production code
7. **Document complex tests**: Add comments explaining non-obvious test logic

**Test Review Checklist**:
- [ ] Tests are co-located with source code
- [ ] Tests use mocks for external dependencies
- [ ] Tests clean up resources after execution
- [ ] Tests have descriptive names
- [ ] Tests verify both success and error cases
- [ ] Tests include property-based tests where applicable
- [ ] Tests are tagged with feature and property references
- [ ] Coverage meets minimum thresholds

