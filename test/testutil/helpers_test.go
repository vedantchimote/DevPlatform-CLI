package testutil

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestLoadFixture tests the LoadFixture helper function
func TestLoadFixture(t *testing.T) {
	// Create a temporary fixture file
	fixtureDir := filepath.Join("../../test", "fixtures", "testdata")
	err := os.MkdirAll(fixtureDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create fixture directory: %v", err)
	}
	defer os.RemoveAll(filepath.Join("../../test", "fixtures", "testdata"))

	testContent := []byte("test fixture content")
	fixturePath := filepath.Join(fixtureDir, "test.txt")
	err = os.WriteFile(fixturePath, testContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create test fixture: %v", err)
	}

	// Test loading the fixture
	data := LoadFixture(t, "testdata/test.txt")
	if string(data) != string(testContent) {
		t.Errorf("LoadFixture() = %q, want %q", string(data), string(testContent))
	}
}

// TestCreateTempConfigYAML tests the CreateTempConfigYAML helper function
func TestCreateTempConfigYAML(t *testing.T) {
	yamlContent := `
global:
  cloud_provider: aws
  timeout: 3600
  log_level: info

aws:
  region: us-east-1
  profile: default

environments:
  dev:
    network_cidr: "10.0.0.0/16"
    db_instance_class: "db.t3.micro"
    db_allocated_storage: 20
    db_multi_az: false
    k8s_node_count: 2
`

	configPath := CreateTempConfigYAML(t, yamlContent)

	// Verify the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("CreateTempConfigYAML() did not create file at %s", configPath)
	}

	// Verify the file contains expected content
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read created config: %v", err)
	}

	content := string(data)
	if !contains(content, "cloud_provider: aws") {
		t.Errorf("Config file does not contain cloud_provider")
	}
	if !contains(content, "us-east-1") {
		t.Errorf("Config file does not contain region")
	}
	if !contains(content, "dev") {
		t.Errorf("Config file does not contain environment")
	}
	if !contains(content, "aws") {
		t.Errorf("Config file does not contain provider")
	}
}

// TestAssertNoError tests the AssertNoError helper function
func TestAssertNoError(t *testing.T) {
	// Test with nil error (should pass)
	AssertNoError(t, nil)

	// Test with non-nil error (should fail)
	// We can't directly test this without causing the test to fail,
	// so we'll just verify the function exists and compiles
}

// TestAssertError tests the AssertError helper function
func TestAssertError(t *testing.T) {
	// Test with non-nil error (should pass)
	AssertError(t, errors.New("test error"))

	// Test with nil error would fail the test, so we skip it
}

// TestAssertEqual tests the AssertEqual helper function
func TestAssertEqual(t *testing.T) {
	// Test with equal values
	AssertEqual(t, 42, 42)
	AssertEqual(t, "hello", "hello")
	AssertEqual(t, []int{1, 2, 3}, []int{1, 2, 3})

	// Test with unequal values would fail, so we skip it
}

// TestAssertContains tests the AssertContains helper function
func TestAssertContains(t *testing.T) {
	// Test with string that contains substring
	AssertContains(t, "hello world", "world")
	AssertContains(t, "error: invalid input", "invalid")

	// Test with string that doesn't contain substring would fail
}

// TestPropertyTest tests the PropertyTest helper function
func TestPropertyTest(t *testing.T) {
	callCount := 0

	PropertyTest(t, "test_property", 10, func(t *testing.T) {
		callCount++
		// Simple test that should always pass
		AssertTrue(t, true, "This should always be true")
	})

	// Verify the test function was called the correct number of times
	if callCount != 10 {
		t.Errorf("PropertyTest() called test function %d times, want 10", callCount)
	}
}

// TestPropertyTestWithContext tests the PropertyTestWithContext helper function
func TestPropertyTestWithContext(t *testing.T) {
	callCount := 0

	PropertyTestWithContext(t, "test_with_context", 5, 1*time.Second, func(t *testing.T, ctx context.Context) {
		callCount++
		
		// Verify context is not nil
		if ctx == nil {
			t.Error("Context should not be nil")
		}

		// Verify context has a deadline
		_, ok := ctx.Deadline()
		if !ok {
			t.Error("Context should have a deadline")
		}
	})

	// Verify the test function was called the correct number of times
	if callCount != 5 {
		t.Errorf("PropertyTestWithContext() called test function %d times, want 5", callCount)
	}
}

// TestAssertTrue tests the AssertTrue helper function
func TestAssertTrue(t *testing.T) {
	// Test with true condition (should pass)
	AssertTrue(t, true, "This is true")
	AssertTrue(t, 1 == 1, "One equals one")

	// Test with false condition would fail
}

// TestAssertFalse tests the AssertFalse helper function
func TestAssertFalse(t *testing.T) {
	// Test with false condition (should pass)
	AssertFalse(t, false, "This is false")
	AssertFalse(t, 1 == 2, "One does not equal two")

	// Test with true condition would fail
}

// TestAssertNotEqual tests the AssertNotEqual helper function
func TestAssertNotEqual(t *testing.T) {
	// Test with unequal values (should pass)
	AssertNotEqual(t, 42, 43)
	AssertNotEqual(t, "hello", "world")

	// Test with equal values would fail
}

// TestAssertPanic tests the AssertPanic helper function
func TestAssertPanic(t *testing.T) {
	// Test with function that panics (should pass)
	AssertPanic(t, func() {
		panic("test panic")
	})

	// Test with function that doesn't panic would fail
}

// TestAssertNoPanic tests the AssertNoPanic helper function
func TestAssertNoPanic(t *testing.T) {
	// Test with function that doesn't panic (should pass)
	AssertNoPanic(t, func() {
		// Do nothing
	})

	// Test with function that panics would fail
}

// TestWithTimeout tests the WithTimeout helper function
func TestWithTimeout(t *testing.T) {
	// Test with function that completes quickly (should pass)
	WithTimeout(t, 1*time.Second, func() {
		time.Sleep(10 * time.Millisecond)
	})

	// Test with function that times out would fail
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
