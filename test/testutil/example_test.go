package testutil_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

// Example_basicAssertions demonstrates basic assertion helpers
func Example_basicAssertions() {
	t := &testing.T{} // In real tests, this comes from the test function parameter

	// Assert no error
	testutil.AssertNoError(t, nil)

	// Assert error exists
	testutil.AssertError(t, errors.New("some error"))

	// Assert values are equal
	testutil.AssertEqual(t, 42, 42)
	testutil.AssertEqual(t, "hello", "hello")

	// Assert string contains substring
	testutil.AssertContains(t, "hello world", "world")

	// Assert values are not equal
	testutil.AssertNotEqual(t, 42, 43)

	// Assert boolean conditions
	testutil.AssertTrue(t, true, "This should be true")
	testutil.AssertFalse(t, false, "This should be false")
}

// Example_propertyBasedTesting demonstrates property-based testing
func Example_propertyBasedTesting() {
	t := &testing.T{}

	// Run a property test 100 times
	testutil.PropertyTest(t, "addition_is_commutative", 100, func(t *testing.T) {
		// Generate test data (in real tests, you'd use random data)
		a, b := 5, 10

		// Test the property: a + b == b + a
		testutil.AssertEqual(t, a+b, b+a)
	})
}

// Example_propertyTestingWithContext demonstrates property testing with context
func Example_propertyTestingWithContext() {
	t := &testing.T{}

	// Run a property test with context and timeout
	testutil.PropertyTestWithContext(t, "operation_completes_within_timeout", 10, 5*time.Second,
		func(t *testing.T, ctx context.Context) {
			// Simulate an operation that respects context
			select {
			case <-time.After(100 * time.Millisecond):
				// Operation completed
			case <-ctx.Done():
				t.Error("Operation was cancelled")
			}
		})
}

// Example_createTempConfig demonstrates creating temporary config files
func Example_createTempConfigYAML() {
	t := &testing.T{}

	// Create a test configuration YAML
	yamlContent := `
global:
  cloud_provider: aws
  timeout: 3600
  log_level: info

aws:
  region: us-east-1
  profile: default
`

	// Create a temporary config file
	configPath := testutil.CreateTempConfigYAML(t, yamlContent)

	// Use the config file in your tests
	_ = configPath // The file will be automatically cleaned up after the test
}

// Example_panicAssertions demonstrates panic testing
func Example_panicAssertions() {
	t := &testing.T{}

	// Assert that a function panics
	testutil.AssertPanic(t, func() {
		panic("expected panic")
	})

	// Assert that a function does not panic
	testutil.AssertNoPanic(t, func() {
		// Safe operation
		_ = 1 + 1
	})
}

// Example_timeoutTesting demonstrates timeout testing
func Example_timeoutTesting() {
	t := &testing.T{}

	// Run a function with a timeout
	testutil.WithTimeout(t, 1*time.Second, func() {
		// This should complete within 1 second
		time.Sleep(100 * time.Millisecond)
	})
}

// Example_errorCodeValidation demonstrates error code validation
func Example_errorCodeValidation() {
	t := &testing.T{}

	// Create an error with a specific code
	err := errors.New("INVALID_CONFIG: configuration is invalid")

	// Assert the error has the expected code
	testutil.AssertErrorCode(t, err, "INVALID_CONFIG")
}

// Example_tableDrivenTest demonstrates using helpers in table-driven tests
func ExamplePropertyTest_tableDriven() {
	t := &testing.T{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase", "hello", "HELLO"},
		{"uppercase", "WORLD", "WORLD"},
		{"mixed", "HeLLo", "HELLO"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use assertion helpers in table-driven tests
			result := toUpper(tt.input)
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

// Helper function for example
func toUpper(s string) string {
	result := ""
	for _, c := range s {
		if c >= 'a' && c <= 'z' {
			result += string(c - 32)
		} else {
			result += string(c)
		}
	}
	return result
}
