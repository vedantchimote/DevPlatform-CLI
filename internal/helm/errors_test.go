package helm

import (
	"errors"
	"testing"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

// TestNewHelmError tests creating a new HelmError
func TestNewHelmError(t *testing.T) {
	err := NewHelmError("install", "myapp", "default", "installation failed", "error output", 1)

	testutil.AssertEqual(t, "install", err.Operation)
	testutil.AssertEqual(t, "myapp", err.Release)
	testutil.AssertEqual(t, "default", err.Namespace)
	testutil.AssertEqual(t, "installation failed", err.Message)
	testutil.AssertEqual(t, "error output", err.Output)
	testutil.AssertEqual(t, 1, err.ExitCode)
}

// TestHelmErrorError tests the Error method
func TestHelmErrorError(t *testing.T) {
	tests := []struct {
		name     string
		err      *HelmError
		contains []string
	}{
		{
			name: "basic error",
			err: &HelmError{
				Operation: "install",
				Release:   "myapp",
				Namespace: "default",
				Message:   "failed to install",
			},
			contains: []string{"install", "myapp", "default", "failed to install"},
		},
		{
			name: "error with events",
			err: &HelmError{
				Operation: "upgrade",
				Release:   "myapp",
				Namespace: "production",
				Message:   "upgrade failed",
				Events:    []string{"Pod failed to start", "ImagePullBackOff"},
			},
			contains: []string{"upgrade", "myapp", "production", "Kubernetes Events", "Pod failed to start"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := tt.err.Error()
			for _, substr := range tt.contains {
				testutil.AssertContains(t, errStr, substr)
			}
		})
	}
}

// TestWithEvents tests adding events to an error
func TestWithEvents(t *testing.T) {
	err := NewHelmError("install", "myapp", "default", "failed", "", 1)
	events := []string{"Event 1", "Event 2", "Event 3"}

	err = err.WithEvents(events)

	testutil.AssertEqual(t, 3, len(err.Events))
	testutil.AssertEqual(t, "Event 1", err.Events[0])
	testutil.AssertEqual(t, "Event 2", err.Events[1])
}

// TestParseHelmError tests parsing helm error output
func TestParseHelmError(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{
			name:     "error with Error: prefix",
			output:   "Some output\nError: release already exists\nMore output",
			expected: " release already exists",
		},
		{
			name:     "error with error: prefix",
			output:   "error: chart not found",
			expected: " chart not found",
		},
		{
			name:     "error with FATAL",
			output:   "FATAL: connection refused",
			expected: "FATAL: connection refused",
		},
		{
			name:     "no specific error pattern",
			output:   "Line 1\nLine 2\nLast line",
			expected: "Last line",
		},
		{
			name:     "empty output",
			output:   "",
			expected: "Unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseHelmError(tt.output)
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

// TestIsReleaseNotFound tests release not found detection
func TestIsReleaseNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "release not found",
			err:      errors.New("Error: release: not found"),
			expected: true,
		},
		{
			name:     "not found in message",
			err:      errors.New("chart not found"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("connection refused"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsReleaseNotFound(tt.err)
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

// TestIsTimeout tests timeout detection
func TestIsTimeout(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "timeout error",
			err:      errors.New("Error: timeout waiting for pods"),
			expected: true,
		},
		{
			name:     "timed out error",
			err:      errors.New("operation timed out"),
			expected: true,
		},
		{
			name:     "deadline exceeded",
			err:      errors.New("context deadline exceeded"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("connection refused"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTimeout(tt.err)
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

// TestIsResourceConflict tests resource conflict detection
func TestIsResourceConflict(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "already exists",
			err:      errors.New("Error: resource already exists"),
			expected: true,
		},
		{
			name:     "conflict error",
			err:      errors.New("conflict: resource in use"),
			expected: true,
		},
		{
			name:     "already owned",
			err:      errors.New("resource already owned by another release"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("connection refused"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsResourceConflict(tt.err)
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

// TestIsInvalidChart tests invalid chart detection
func TestIsInvalidChart(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "chart not found",
			err:      errors.New("Error: chart not found"),
			expected: true,
		},
		{
			name:     "invalid chart",
			err:      errors.New("invalid chart format"),
			expected: true,
		},
		{
			name:     "chart.yaml error",
			err:      errors.New("failed to load Chart.yaml"),
			expected: true,
		},
		{
			name:     "no chart found",
			err:      errors.New("no chart found at path"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("connection refused"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInvalidChart(tt.err)
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

// TestIsInvalidValues tests invalid values detection
func TestIsInvalidValues(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "invalid values",
			err:      errors.New("Error: invalid values provided"),
			expected: true,
		},
		{
			name:     "values.yaml error",
			err:      errors.New("failed to load values.yaml"),
			expected: true,
		},
		{
			name:     "yaml parse error",
			err:      errors.New("yaml: parse error"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("connection refused"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInvalidValues(tt.err)
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

// TestGetErrorCategory tests error categorization
func TestGetErrorCategory(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "release not found",
			err:      errors.New("release: not found"),
			expected: "release_not_found",
		},
		{
			name:     "timeout",
			err:      errors.New("timeout waiting for pods"),
			expected: "timeout",
		},
		{
			name:     "resource conflict",
			err:      errors.New("resource already exists"),
			expected: "resource_conflict",
		},
		{
			name:     "invalid chart",
			err:      errors.New("chart not found"),
			expected: "invalid_chart",
		},
		{
			name:     "invalid values",
			err:      errors.New("invalid values"),
			expected: "invalid_values",
		},
		{
			name:     "generic helm error",
			err:      errors.New("some other error"),
			expected: "helm_error",
		},
		{
			name:     "nil error",
			err:      nil,
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetErrorCategory(tt.err)
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

// TestGetErrorResolution tests error resolution suggestions
func TestGetErrorResolution(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		contains string
	}{
		{
			name:     "release not found",
			err:      errors.New("release: not found"),
			contains: "helm list",
		},
		{
			name:     "timeout",
			err:      errors.New("timeout"),
			contains: "kubectl get pods",
		},
		{
			name:     "resource conflict",
			err:      errors.New("already exists"),
			contains: "helm upgrade",
		},
		{
			name:     "invalid chart",
			err:      errors.New("chart not found"),
			contains: "Chart.yaml",
		},
		{
			name:     "invalid values",
			err:      errors.New("invalid values"),
			contains: "YAML syntax",
		},
		{
			name:     "generic error",
			err:      errors.New("some error"),
			contains: "helm --debug",
		},
		{
			name:     "nil error",
			err:      nil,
			contains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetErrorResolution(tt.err)
			if tt.contains != "" {
				testutil.AssertContains(t, result, tt.contains)
			} else {
				testutil.AssertEqual(t, "", result)
			}
		})
	}
}

// TestHelmErrorWithMultipleEvents tests error with multiple events
func TestHelmErrorWithMultipleEvents(t *testing.T) {
	err := NewHelmError("install", "myapp", "default", "failed", "", 1)
	events := []string{
		"Pod myapp-123 failed to start",
		"ImagePullBackOff: image not found",
		"CrashLoopBackOff: container exited with code 1",
	}

	err = err.WithEvents(events)

	testutil.AssertEqual(t, 3, len(err.Events))
	
	errStr := err.Error()
	testutil.AssertContains(t, errStr, "Kubernetes Events")
	testutil.AssertContains(t, errStr, "ImagePullBackOff")
	testutil.AssertContains(t, errStr, "CrashLoopBackOff")
}

// TestParseHelmErrorWithMultipleLines tests parsing multi-line error output
func TestParseHelmErrorWithMultipleLines(t *testing.T) {
	output := `
Helm install started
Processing chart...
Error: failed to create resource
Additional context here
`
	result := ParseHelmError(output)
	testutil.AssertEqual(t, " failed to create resource", result)
}

// TestErrorCategoryPriority tests that error categories are checked in correct priority
func TestErrorCategoryPriority(t *testing.T) {
	// An error that matches multiple categories should return the first match
	err := errors.New("release: not found and timeout occurred")
	
	// Should match release_not_found first
	category := GetErrorCategory(err)
	testutil.AssertEqual(t, "release_not_found", category)
}

// TestHelmErrorImplementsError tests that HelmError implements error interface
func TestHelmErrorImplementsError(t *testing.T) {
	var _ error = &HelmError{}
	var _ error = NewHelmError("install", "myapp", "default", "failed", "", 1)
}
