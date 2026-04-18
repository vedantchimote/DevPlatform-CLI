# Test Utility Package

This package provides common test helper functions and utilities for the DevPlatform CLI test suite.

## Overview

The `testutil` package contains reusable test helpers that simplify writing tests across the codebase. These helpers provide:

- Assertion functions for common test validations
- Property-based testing support
- Test fixture loading
- Temporary configuration file creation
- Timeout and context management

## Installation

Import the package in your test files:

```go
import "github.com/devplatform/devplatform-cli/test/testutil"
```

## Core Functions

### Assertion Helpers

#### AssertNoError

Fails the test if an error is not nil.

```go
func TestSomething(t *testing.T) {
    err := someFunction()
    testutil.AssertNoError(t, err)
}
```

#### AssertError

Fails the test if an error is nil (when you expect an error).

```go
func TestErrorCase(t *testing.T) {
    err := functionThatShouldFail()
    testutil.AssertError(t, err)
}
```

#### AssertEqual

Fails the test if two values are not deeply equal.

```go
func TestEquality(t *testing.T) {
    result := calculate(5, 10)
    testutil.AssertEqual(t, 15, result)
    
    // Works with complex types too
    expected := &Config{Name: "test"}
    actual := loadConfig()
    testutil.AssertEqual(t, expected, actual)
}
```

#### AssertNotEqual

Fails the test if two values are equal (when they should be different).

```go
func TestDifference(t *testing.T) {
    id1 := generateID()
    id2 := generateID()
    testutil.AssertNotEqual(t, id1, id2)
}
```

#### AssertContains

Fails the test if a string doesn't contain a substring.

```go
func TestErrorMessage(t *testing.T) {
    err := validateInput("invalid")
    testutil.AssertError(t, err)
    testutil.AssertContains(t, err.Error(), "INVALID_INPUT")
}
```

#### AssertTrue / AssertFalse

Fails the test if a boolean condition is not met.

```go
func TestConditions(t *testing.T) {
    testutil.AssertTrue(t, len(items) > 0, "Items should not be empty")
    testutil.AssertFalse(t, isExpired, "Should not be expired")
}
```

#### AssertErrorCode

Fails the test if an error doesn't contain the expected error code.

```go
func TestSpecificError(t *testing.T) {
    err := validateConfig(invalidConfig)
    testutil.AssertErrorCode(t, err, "INVALID_CONFIG")
}
```

### Panic Testing

#### AssertPanic

Fails the test if a function doesn't panic.

```go
func TestInvalidInputPanics(t *testing.T) {
    testutil.AssertPanic(t, func() {
        processData(nil) // Should panic with nil input
    })
}
```

#### AssertNoPanic

Fails the test if a function panics.

```go
func TestValidInputDoesNotPanic(t *testing.T) {
    testutil.AssertNoPanic(t, func() {
        processData(validData) // Should not panic
    })
}
```

### Property-Based Testing

#### PropertyTest

Runs a test function multiple times to verify a property holds across iterations.

```go
func TestCommutativeProperty(t *testing.T) {
    testutil.PropertyTest(t, "addition_commutative", 100, func(t *testing.T) {
        // In real tests, generate random values here
        a, b := randomInt(), randomInt()
        
        // Test the property
        testutil.AssertEqual(t, add(a, b), add(b, a))
    })
}
```

**Parameters:**
- `t`: The testing.T instance
- `name`: Name of the property being tested
- `iterations`: Number of times to run the test (minimum 100 recommended)
- `testFunc`: The test function to execute

#### PropertyTestWithContext

Runs a property test with a context that has a timeout.

```go
func TestOperationTimeout(t *testing.T) {
    testutil.PropertyTestWithContext(t, "respects_timeout", 50, 5*time.Second,
        func(t *testing.T, ctx context.Context) {
            err := longRunningOperation(ctx)
            testutil.AssertNoError(t, err)
        })
}
```

**Parameters:**
- `t`: The testing.T instance
- `name`: Name of the property being tested
- `iterations`: Number of times to run the test
- `timeout`: Maximum duration for each iteration
- `testFunc`: The test function that receives a context

### Fixture Management

#### LoadFixture

Loads a test fixture file from the `test/fixtures` directory.

```go
func TestConfigParsing(t *testing.T) {
    data := testutil.LoadFixture(t, "configs/valid_config.yaml")
    
    config, err := parseConfig(data)
    testutil.AssertNoError(t, err)
    testutil.AssertEqual(t, "myapp", config.AppName)
}
```

**Parameters:**
- `t`: The testing.T instance
- `path`: Relative path from `test/fixtures/` directory

**Returns:** Byte slice containing the fixture file contents

**Note:** The test will fail immediately if the fixture file cannot be loaded.

#### CreateTempConfig

Creates a temporary configuration file for testing.

```go
func TestConfigLoading(t *testing.T) {
    cfg := &config.Config{
        AppName:     "testapp",
        Environment: "dev",
        Provider:    "aws",
        Region:      "us-east-1",
    }
    
    configPath := testutil.CreateTempConfig(t, cfg)
    
    // Use the config file
    loadedConfig, err := config.Load(configPath)
    testutil.AssertNoError(t, err)
    testutil.AssertEqual(t, cfg.AppName, loadedConfig.AppName)
    
    // File is automatically cleaned up after test
}
```

**Parameters:**
- `t`: The testing.T instance
- `cfg`: Configuration object to write to file

**Returns:** Path to the temporary configuration file

**Note:** The file is automatically cleaned up when the test completes.

### Timeout Management

#### WithTimeout

Runs a function with a timeout, failing the test if it doesn't complete in time.

```go
func TestQuickOperation(t *testing.T) {
    testutil.WithTimeout(t, 1*time.Second, func() {
        result := quickOperation()
        testutil.AssertEqual(t, "done", result)
    })
}
```

**Parameters:**
- `t`: The testing.T instance
- `timeout`: Maximum duration to wait
- `f`: The function to execute

## Usage Patterns

### Table-Driven Tests

Use assertion helpers in table-driven tests for cleaner code:

```go
func TestValidateAppName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
        errCode string
    }{
        {
            name:    "valid name",
            input:   "myapp",
            wantErr: false,
        },
        {
            name:    "invalid uppercase",
            input:   "MyApp",
            wantErr: true,
            errCode: "INVALID_APP_NAME",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateAppName(tt.input)
            
            if tt.wantErr {
                testutil.AssertError(t, err)
                if tt.errCode != "" {
                    testutil.AssertErrorCode(t, err, tt.errCode)
                }
            } else {
                testutil.AssertNoError(t, err)
            }
        })
    }
}
```

### Property-Based Testing Pattern

```go
func TestMockRecordsCallsProperty(t *testing.T) {
    // Feature: comprehensive-testing-suite, Property 2: Mock Call Recording
    testutil.PropertyTest(t, "mock_records_calls", 100, func(t *testing.T) {
        mock := NewMockExecutor()
        
        // Generate random test data
        workingDir := generateRandomPath()
        
        // Execute operation
        err := mock.Init(context.Background(), workingDir)
        testutil.AssertNoError(t, err)
        
        // Verify call was recorded
        testutil.AssertEqual(t, 1, len(mock.InitCalls))
        testutil.AssertEqual(t, workingDir, mock.InitCalls[0].Args[0])
    })
}
```

### Integration Test Pattern

```go
func TestCreateWorkflow(t *testing.T) {
    // Create test config
    cfg := &config.Config{
        AppName:     "testapp",
        Environment: "dev",
        Provider:    "aws",
        Region:      "us-east-1",
    }
    configPath := testutil.CreateTempConfig(t, cfg)
    
    // Load fixtures
    tfState := testutil.LoadFixture(t, "terraform/sample_state.json")
    
    // Run workflow
    err := runCreateWorkflow(configPath)
    testutil.AssertNoError(t, err)
    
    // Verify results
    status, err := getStatus(cfg.AppName)
    testutil.AssertNoError(t, err)
    testutil.AssertEqual(t, "running", status)
}
```

### Error Testing Pattern

```go
func TestErrorHandling(t *testing.T) {
    tests := []struct {
        name    string
        input   interface{}
        errCode string
    }{
        {"nil input", nil, "INVALID_INPUT"},
        {"empty string", "", "INVALID_INPUT"},
        {"invalid format", "###", "INVALID_FORMAT"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := process(tt.input)
            testutil.AssertError(t, err)
            testutil.AssertErrorCode(t, err, tt.errCode)
        })
    }
}
```

## Best Practices

1. **Use t.Helper()**: All helper functions call `t.Helper()` to ensure error messages point to the correct line in your test code.

2. **Descriptive Messages**: When using `AssertTrue` and `AssertFalse`, provide clear messages explaining what should be true/false.

3. **Property Test Iterations**: Use at least 100 iterations for property-based tests to ensure good coverage.

4. **Fixture Organization**: Organize fixtures in subdirectories under `test/fixtures/` by category (configs, terraform, helm, etc.).

5. **Cleanup**: Use `t.TempDir()` and `CreateTempConfig()` for automatic cleanup of temporary files.

6. **Error Validation**: Always validate both that an error occurred AND that it has the correct error code/message.

7. **Context Usage**: Use `PropertyTestWithContext` when testing operations that should respect context cancellation.

## Examples

See `example_test.go` in this package for complete working examples of all helper functions.

## Contributing

When adding new helper functions:

1. Add the function to `helpers.go`
2. Add tests to `helpers_test.go`
3. Add examples to `example_test.go`
4. Update this README with usage documentation

## Related Documentation

- [Testing Guide](../../docs/testing/README.md) - Overall testing strategy
- [Mock Framework](../mocks/README.md) - Mock implementations
- [Integration Tests](../integration/README.md) - Integration test patterns
