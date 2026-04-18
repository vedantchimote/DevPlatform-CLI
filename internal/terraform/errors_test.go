package terraform

import (
	"errors"
	"testing"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

// TestTerraformErrorError tests the Error method
func TestTerraformErrorError(t *testing.T) {
	tests := []struct {
		name     string
		err      *TerraformError
		contains []string
	}{
		{
			name: "basic error",
			err: &TerraformError{
				Command: "apply",
				Message: "failed to apply",
			},
			contains: []string{"Terraform apply failed", "failed to apply"},
		},
		{
			name: "error with exit code",
			err: &TerraformError{
				Command:  "plan",
				ExitCode: 1,
				Message:  "planning failed",
			},
			contains: []string{"Terraform plan failed", "exit code 1", "planning failed"},
		},
		{
			name: "error with suggestion",
			err: &TerraformError{
				Command:    "init",
				Message:    "initialization failed",
				Suggestion: "Run terraform init again",
			},
			contains: []string{"Terraform init failed", "initialization failed", "Suggestion:", "Run terraform init again"},
		},
		{
			name: "error with output",
			err: &TerraformError{
				Command: "destroy",
				Message: "destroy failed",
				Output:  "Error: resource not found",
			},
			contains: []string{"Terraform destroy failed", "destroy failed", "Terraform Output:", "Error: resource not found"},
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

// TestErrorCategoryString tests the String method of ErrorCategory
func TestErrorCategoryString(t *testing.T) {
	tests := []struct {
		category ErrorCategory
		expected string
	}{
		{ErrorCategoryConfiguration, "Configuration Error"},
		{ErrorCategoryState, "State Error"},
		{ErrorCategoryProvider, "Provider Error"},
		{ErrorCategoryResource, "Resource Error"},
		{ErrorCategoryValidation, "Validation Error"},
		{ErrorCategoryPermission, "Permission Error"},
		{ErrorCategoryNetwork, "Network Error"},
		{ErrorCategoryUnknown, "Unknown Error"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.category.String()
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

// TestParseTerraformError tests parsing Terraform errors
func TestParseTerraformError(t *testing.T) {
	tests := []struct {
		name           string
		command        string
		err            error
		output         string
		expectedCat    ErrorCategory
		suggestionHas  string
	}{
		{
			name:          "configuration error",
			command:       "plan",
			err:           errors.New("invalid configuration"),
			output:        "Error: Invalid configuration\nMissing required argument",
			expectedCat:   ErrorCategoryValidation, // "invalid" triggers validation category
			suggestionHas: "input values",
		},
		{
			name:          "state lock error",
			command:       "apply",
			err:           errors.New("state lock error"),
			output:        "Error: State lock acquired",
			expectedCat:   ErrorCategoryState,
			suggestionHas: "locked",
		},
		{
			name:          "no state file error",
			command:       "output",
			err:           errors.New("no state file"),
			output:        "Error: No state file found",
			expectedCat:   ErrorCategoryState,
			suggestionHas: "terraform init",
		},
		{
			name:          "provider not found",
			command:       "init",
			err:           errors.New("provider not found"),
			output:        "Error: Provider not found",
			expectedCat:   ErrorCategoryProvider,
			suggestionHas: "terraform init",
		},
		{
			name:          "authentication error",
			command:       "apply",
			err:           errors.New("authentication failed"),
			output:        "Error: Authentication failed",
			expectedCat:   ErrorCategoryPermission,
			suggestionHas: "credentials",
		},
		{
			name:          "resource already exists",
			command:       "apply",
			err:           errors.New("resource already exists"),
			output:        "Error: Resource already exists",
			expectedCat:   ErrorCategoryResource,
			suggestionHas: "already exists",
		},
		{
			name:          "resource not found",
			command:       "destroy",
			err:           errors.New("resource not found"),
			output:        "Error: Resource not found",
			expectedCat:   ErrorCategoryResource,
			suggestionHas: "does not exist",
		},
		{
			name:          "validation error",
			command:       "plan",
			err:           errors.New("validation failed"),
			output:        "Error: Validation failed",
			expectedCat:   ErrorCategoryValidation,
			suggestionHas: "input values",
		},
		{
			name:          "network timeout",
			command:       "apply",
			err:           errors.New("connection timeout"),
			output:        "Error: Connection timeout",
			expectedCat:   ErrorCategoryNetwork,
			suggestionHas: "network connection",
		},
		{
			name:          "access denied",
			command:       "apply",
			err:           errors.New("access denied"),
			output:        "Error: Access denied",
			expectedCat:   ErrorCategoryPermission,
			suggestionHas: "permissions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseTerraformError(tt.command, tt.err, tt.output)
			testutil.AssertTrue(t, result != nil, "Result should not be nil")
			testutil.AssertEqual(t, tt.command, result.Command)
			testutil.AssertEqual(t, tt.expectedCat, result.Category)
			if tt.suggestionHas != "" {
				testutil.AssertContains(t, result.Suggestion, tt.suggestionHas)
			}
		})
	}
}

// TestParseTerraformErrorNil tests parsing nil error
func TestParseTerraformErrorNil(t *testing.T) {
	result := ParseTerraformError("apply", nil, "")
	testutil.AssertTrue(t, result == nil, "Result should be nil for nil error")
}

// TestIsStateLockError tests state lock error detection
func TestIsStateLockError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "state lock error",
			err:      errors.New("Error: state lock acquired"),
			expected: true,
		},
		{
			name:     "locked error",
			err:      errors.New("Error: resource is locked"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("Error: something else"),
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
			result := IsStateLockError(tt.err)
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

// TestIsPermissionError tests permission error detection
func TestIsPermissionError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "access denied",
			err:      errors.New("Error: access denied"),
			expected: true,
		},
		{
			name:     "forbidden",
			err:      errors.New("Error: forbidden"),
			expected: true,
		},
		{
			name:     "unauthorized",
			err:      errors.New("Error: unauthorized"),
			expected: true,
		},
		{
			name:     "insufficient permissions",
			err:      errors.New("Error: insufficient permissions"),
			expected: true,
		},
		{
			name:     "authentication failed",
			err:      errors.New("Error: authentication failed"),
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("Error: something else"),
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
			result := IsPermissionError(tt.err)
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

// TestIsResourceNotFoundError tests resource not found error detection
func TestIsResourceNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "resource not found",
			err:      errors.New("Error: resource not found"),
			expected: true,
		},
		{
			name:     "provider not found (should be false)",
			err:      errors.New("Error: provider not found"),
			expected: false,
		},
		{
			name:     "other error",
			err:      errors.New("Error: something else"),
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
			result := IsResourceNotFoundError(tt.err)
			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

// TestExtractErrorMessage tests extracting error messages
func TestExtractErrorMessage(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{
			name: "error with Error: prefix",
			output: `Some output
Error: Resource not found
The resource does not exist

More output`,
			expected: "Error: Resource not found",
		},
		{
			name:     "no error section",
			output:   "Line 1\nLine 2\nLast line",
			expected: "Last line",
		},
		{
			name:     "empty output",
			output:   "",
			expected: "",
		},
		{
			name:     "only whitespace",
			output:   "   \n  \n  ",
			expected: "   \n  \n  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractErrorMessage(tt.output)
			testutil.AssertContains(t, result, tt.expected)
		})
	}
}
