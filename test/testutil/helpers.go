package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// LoadFixture loads a test fixture file from the fixtures directory
// and returns its contents as a byte slice.
func LoadFixture(t *testing.T, path string) []byte {
	t.Helper()
	
	// Find project root by looking for go.mod
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}
	
	// Construct full path relative to project root
	fullPath := filepath.Join(projectRoot, "test", "fixtures", path)
	
	data, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("Failed to load fixture %s: %v", path, err)
	}
	
	return data
}

// findProjectRoot finds the project root directory by looking for go.mod
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	
	// Walk up the directory tree until we find go.mod
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find go.mod in any parent directory")
		}
		dir = parent
	}
}

// CreateTempConfigYAML creates a temporary YAML configuration file for testing
// and returns the path to the file. The file will be automatically
// cleaned up when the test completes.
// The yamlContent parameter should be the complete YAML content as a string.
func CreateTempConfigYAML(t *testing.T, yamlContent string) string {
	t.Helper()
	
	// Create a temporary directory
	tempDir := t.TempDir()
	
	// Create config file path
	configPath := filepath.Join(tempDir, "config.yaml")
	
	// Write config file
	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config: %v", err)
	}
	
	return configPath
}

// AssertNoError fails the test if the error is not nil.
// It provides a clear error message indicating where the error occurred.
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
}

// AssertError fails the test if the error is nil.
// It ensures that an error was returned when one was expected.
func AssertError(t *testing.T, err error) {
	t.Helper()
	
	if err == nil {
		t.Fatal("Expected an error, but got nil")
	}
}

// AssertEqual fails the test if the expected and actual values are not equal.
// It uses deep equality comparison and provides detailed error messages.
func AssertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Values not equal:\n  Expected: %+v\n  Actual:   %+v", expected, actual)
	}
}

// AssertContains fails the test if the string does not contain the substring.
// It's useful for checking error messages and log output.
func AssertContains(t *testing.T, str, substr string) {
	t.Helper()
	
	if !strings.Contains(str, substr) {
		t.Fatalf("String does not contain expected substring:\n  String:    %q\n  Substring: %q", str, substr)
	}
}

// PropertyTest runs a test function multiple times with different iterations.
// This is useful for property-based testing where you want to verify that
// a property holds across many different inputs or scenarios.
//
// The test function is called 'iterations' times, each time with a unique
// subtest name that includes the iteration number.
//
// Example usage:
//
//	PropertyTest(t, "mock_records_calls", 100, func(t *testing.T) {
//	    // Test code that should pass for all iterations
//	    mock := NewMockExecutor()
//	    mock.Execute()
//	    AssertEqual(t, 1, len(mock.Calls))
//	})
func PropertyTest(t *testing.T, name string, iterations int, testFunc func(t *testing.T)) {
	t.Helper()
	
	for i := 0; i < iterations; i++ {
		t.Run(fmt.Sprintf("%s_iteration_%d", name, i), testFunc)
	}
}

// PropertyTestWithContext runs a property test with a context that can be cancelled.
// This is useful for testing timeout and cancellation behavior.
//
// The test function receives both a testing.T and a context.Context.
// The context will be cancelled after the specified timeout duration.
//
// Example usage:
//
//	PropertyTestWithContext(t, "operation_respects_timeout", 10, 5*time.Second, 
//	    func(t *testing.T, ctx context.Context) {
//	        err := LongRunningOperation(ctx)
//	        AssertNoError(t, err)
//	    })
func PropertyTestWithContext(t *testing.T, name string, iterations int, timeout time.Duration, testFunc func(t *testing.T, ctx context.Context)) {
	t.Helper()
	
	for i := 0; i < iterations; i++ {
		t.Run(fmt.Sprintf("%s_iteration_%d", name, i), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			
			testFunc(t, ctx)
		})
	}
}

// AssertErrorCode fails the test if the error does not have the expected error code.
// This is useful for validating specific error types from the internal/errors package.
func AssertErrorCode(t *testing.T, err error, expectedCode string) {
	t.Helper()
	
	if err == nil {
		t.Fatalf("Expected error with code %s, but got nil", expectedCode)
	}
	
	// Check if error message contains the error code
	// This is a simple implementation; adjust based on your error package structure
	errMsg := err.Error()
	if !strings.Contains(errMsg, expectedCode) {
		t.Fatalf("Error does not contain expected code:\n  Expected code: %s\n  Error: %v", expectedCode, err)
	}
}

// AssertNotEqual fails the test if the expected and actual values are equal.
// It's the inverse of AssertEqual.
func AssertNotEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	
	if reflect.DeepEqual(expected, actual) {
		t.Fatalf("Values should not be equal, but both are: %+v", expected)
	}
}

// AssertTrue fails the test if the condition is false.
func AssertTrue(t *testing.T, condition bool, message string) {
	t.Helper()
	
	if !condition {
		t.Fatalf("Assertion failed: %s", message)
	}
}

// AssertFalse fails the test if the condition is true.
func AssertFalse(t *testing.T, condition bool, message string) {
	t.Helper()
	
	if condition {
		t.Fatalf("Assertion failed: %s", message)
	}
}

// AssertPanic fails the test if the function does not panic.
// It's useful for testing that invalid inputs cause panics.
func AssertPanic(t *testing.T, f func()) {
	t.Helper()
	
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected function to panic, but it did not")
		}
	}()
	
	f()
}

// AssertNoPanic fails the test if the function panics.
// It's useful for ensuring that valid inputs don't cause panics.
func AssertNoPanic(t *testing.T, f func()) {
	t.Helper()
	
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Expected function not to panic, but it panicked with: %v", r)
		}
	}()
	
	f()
}

// WithTimeout runs a test function with a timeout.
// If the function doesn't complete within the timeout, the test fails.
func WithTimeout(t *testing.T, timeout time.Duration, f func()) {
	t.Helper()
	
	done := make(chan bool)
	
	go func() {
		f()
		done <- true
	}()
	
	select {
	case <-done:
		// Test completed successfully
	case <-time.After(timeout):
		t.Fatalf("Test timed out after %v", timeout)
	}
}
