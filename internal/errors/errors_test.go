package errors

import (
	"errors"
	"strings"
	"testing"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

// TestErrorCategoryConstants tests error category constants
func TestErrorCategoryConstants(t *testing.T) {
	tests := []struct {
		name     string
		category ErrorCategory
		expected string
	}{
		{name: "authentication category", category: CategoryAuthentication, expected: "authentication"},
		{name: "validation category", category: CategoryValidation, expected: "validation"},
		{name: "terraform category", category: CategoryTerraform, expected: "terraform"},
		{name: "helm category", category: CategoryHelm, expected: "helm"},
		{name: "network category", category: CategoryNetwork, expected: "network"},
		{name: "configuration category", category: CategoryConfiguration, expected: "configuration"},
		{name: "state category", category: CategoryState, expected: "state"},
		{name: "unknown category", category: CategoryUnknown, expected: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.AssertEqual(t, tt.expected, string(tt.category))
		})
	}
}

// TestErrorCodeConstants tests error code constants
func TestErrorCodeConstants(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		expected int
	}{
		// Authentication codes (1000-1099)
		{name: "invalid credentials", code: ErrCodeAuthInvalidCredentials, expected: 1001},
		{name: "expired credentials", code: ErrCodeAuthExpiredCredentials, expected: 1002},
		{name: "missing credentials", code: ErrCodeAuthMissingCredentials, expected: 1003},
		{name: "permission denied", code: ErrCodeAuthPermissionDenied, expected: 1004},
		
		// Validation codes (1100-1199)
		{name: "invalid app name", code: ErrCodeValidationInvalidAppName, expected: 1101},
		{name: "invalid environment", code: ErrCodeValidationInvalidEnvironment, expected: 1102},
		{name: "invalid provider", code: ErrCodeValidationInvalidProvider, expected: 1103},
		{name: "invalid config", code: ErrCodeValidationInvalidConfig, expected: 1104},
		{name: "missing required", code: ErrCodeValidationMissingRequired, expected: 1105},
		
		// Terraform codes (1200-1299)
		{name: "terraform init failed", code: ErrCodeTerraformInitFailed, expected: 1201},
		{name: "terraform plan failed", code: ErrCodeTerraformPlanFailed, expected: 1202},
		{name: "terraform apply failed", code: ErrCodeTerraformApplyFailed, expected: 1203},
		{name: "terraform destroy failed", code: ErrCodeTerraformDestroyFailed, expected: 1204},
		{name: "terraform output failed", code: ErrCodeTerraformOutputFailed, expected: 1205},
		
		// Helm codes (1300-1399)
		{name: "helm install failed", code: ErrCodeHelmInstallFailed, expected: 1301},
		{name: "helm upgrade failed", code: ErrCodeHelmUpgradeFailed, expected: 1302},
		{name: "helm uninstall failed", code: ErrCodeHelmUninstallFailed, expected: 1303},
		{name: "helm status failed", code: ErrCodeHelmStatusFailed, expected: 1304},
		{name: "helm pod not ready", code: ErrCodeHelmPodNotReady, expected: 1305},
		{name: "helm timeout", code: ErrCodeHelmTimeout, expected: 1306},
		
		// Network codes (1400-1499)
		{name: "network connection failed", code: ErrCodeNetworkConnectionFailed, expected: 1401},
		{name: "network timeout", code: ErrCodeNetworkTimeout, expected: 1402},
		{name: "network dns failed", code: ErrCodeNetworkDNSFailed, expected: 1403},
		
		// Configuration codes (1500-1599)
		{name: "config file not found", code: ErrCodeConfigFileNotFound, expected: 1501},
		{name: "config parse failed", code: ErrCodeConfigParseFailed, expected: 1502},
		{name: "config invalid format", code: ErrCodeConfigInvalidFormat, expected: 1503},
		
		// State codes (1600-1699)
		{name: "state locked", code: ErrCodeStateLocked, expected: 1601},
		{name: "state not found", code: ErrCodeStateNotFound, expected: 1602},
		{name: "state corrupted", code: ErrCodeStateCorrupted, expected: 1603},
		{name: "state access denied", code: ErrCodeStateAccessDenied, expected: 1604},
		
		// Unknown codes (9000-9999)
		{name: "unknown error", code: ErrCodeUnknown, expected: 9000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.AssertEqual(t, tt.expected, int(tt.code))
		})
	}
}

// TestNewCLIError tests creating a new CLI error
func TestNewCLIError(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewCLIError(CategoryValidation, ErrCodeValidationInvalidAppName, "Invalid app name", cause)

	testutil.AssertEqual(t, CategoryValidation, err.Category)
	testutil.AssertEqual(t, ErrCodeValidationInvalidAppName, err.Code)
	testutil.AssertEqual(t, "Invalid app name", err.Message)
	testutil.AssertEqual(t, cause, err.Cause)
	testutil.AssertEqual(t, "", err.Details)
	testutil.AssertEqual(t, "", err.Resolution)
	testutil.AssertEqual(t, "", err.LogPath)
}

// TestCLIErrorError tests the Error method
func TestCLIErrorError(t *testing.T) {
	tests := []struct {
		name     string
		err      *CLIError
		contains []string
	}{
		{
			name: "error without cause",
			err: &CLIError{
				Category: CategoryValidation,
				Code:     ErrCodeValidationInvalidAppName,
				Message:  "Invalid app name",
			},
			contains: []string{"validation", "1101", "Invalid app name"},
		},
		{
			name: "error with cause",
			err: &CLIError{
				Category: CategoryTerraform,
				Code:     ErrCodeTerraformApplyFailed,
				Message:  "Terraform apply failed",
				Cause:    errors.New("resource already exists"),
			},
			contains: []string{"terraform", "1203", "Terraform apply failed", "resource already exists"},
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

// TestCLIErrorUnwrap tests the Unwrap method
func TestCLIErrorUnwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewCLIError(CategoryValidation, ErrCodeValidationInvalidAppName, "Invalid app name", cause)

	unwrapped := err.Unwrap()
	testutil.AssertEqual(t, cause, unwrapped)
}

// TestCLIErrorUnwrapNil tests Unwrap with no cause
func TestCLIErrorUnwrapNil(t *testing.T) {
	err := NewCLIError(CategoryValidation, ErrCodeValidationInvalidAppName, "Invalid app name", nil)

	unwrapped := err.Unwrap()
	testutil.AssertTrue(t, unwrapped == nil, "Unwrap should return nil when no cause")
}

// TestWithDetails tests adding details to an error
func TestWithDetails(t *testing.T) {
	err := NewCLIError(CategoryValidation, ErrCodeValidationInvalidAppName, "Invalid app name", nil)
	err = err.WithDetails("App name must be lowercase")

	testutil.AssertEqual(t, "App name must be lowercase", err.Details)
}

// TestWithResolution tests adding resolution to an error
func TestWithResolution(t *testing.T) {
	err := NewCLIError(CategoryValidation, ErrCodeValidationInvalidAppName, "Invalid app name", nil)
	err = err.WithResolution("Use only lowercase letters")

	testutil.AssertEqual(t, "Use only lowercase letters", err.Resolution)
}

// TestWithLogPath tests adding log path to an error
func TestWithLogPath(t *testing.T) {
	err := NewCLIError(CategoryValidation, ErrCodeValidationInvalidAppName, "Invalid app name", nil)
	err = err.WithLogPath("/var/log/devplatform.log")

	testutil.AssertEqual(t, "/var/log/devplatform.log", err.LogPath)
}

// TestCLIErrorFormat tests the Format method
func TestCLIErrorFormat(t *testing.T) {
	tests := []struct {
		name     string
		err      *CLIError
		contains []string
	}{
		{
			name: "basic error",
			err: &CLIError{
				Category: CategoryValidation,
				Code:     ErrCodeValidationInvalidAppName,
				Message:  "Invalid app name",
			},
			contains: []string{"❌ Error", "validation:1101", "Message: Invalid app name"},
		},
		{
			name: "error with details",
			err: &CLIError{
				Category: CategoryValidation,
				Code:     ErrCodeValidationInvalidAppName,
				Message:  "Invalid app name",
				Details:  "App name contains uppercase letters",
			},
			contains: []string{"❌ Error", "Details:", "App name contains uppercase letters"},
		},
		{
			name: "error with resolution",
			err: &CLIError{
				Category: CategoryValidation,
				Code:     ErrCodeValidationInvalidAppName,
				Message:  "Invalid app name",
				Resolution: "Use only lowercase letters and hyphens",
			},
			contains: []string{"❌ Error", "💡 Resolution:", "Use only lowercase letters and hyphens"},
		},
		{
			name: "error with log path",
			err: &CLIError{
				Category: CategoryTerraform,
				Code:     ErrCodeTerraformApplyFailed,
				Message:  "Terraform apply failed",
				LogPath:  "/var/log/devplatform.log",
			},
			contains: []string{"❌ Error", "📝 Log file:", "/var/log/devplatform.log"},
		},
		{
			name: "complete error",
			err: &CLIError{
				Category:   CategoryHelm,
				Code:       ErrCodeHelmInstallFailed,
				Message:    "Helm install failed",
				Details:    "Pod failed to start",
				Resolution: "Check pod logs",
				LogPath:    "/var/log/devplatform.log",
			},
			contains: []string{
				"❌ Error", "helm:1301",
				"Message: Helm install failed",
				"Details:", "Pod failed to start",
				"💡 Resolution:", "Check pod logs",
				"📝 Log file:", "/var/log/devplatform.log",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := tt.err.Format()
			
			for _, substr := range tt.contains {
				testutil.AssertContains(t, formatted, substr)
			}
		})
	}
}

// TestNewAuthError tests authentication error creation
func TestNewAuthError(t *testing.T) {
	tests := []struct {
		name       string
		code       ErrorCode
		message    string
		hasResolution bool
	}{
		{
			name:       "invalid credentials",
			code:       ErrCodeAuthInvalidCredentials,
			message:    "Invalid AWS credentials",
			hasResolution: true,
		},
		{
			name:       "expired credentials",
			code:       ErrCodeAuthExpiredCredentials,
			message:    "AWS credentials expired",
			hasResolution: true,
		},
		{
			name:       "missing credentials",
			code:       ErrCodeAuthMissingCredentials,
			message:    "AWS credentials not found",
			hasResolution: true,
		},
		{
			name:       "permission denied",
			code:       ErrCodeAuthPermissionDenied,
			message:    "Insufficient permissions",
			hasResolution: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAuthError(tt.code, tt.message, nil)

			testutil.AssertEqual(t, CategoryAuthentication, err.Category)
			testutil.AssertEqual(t, tt.code, err.Code)
			testutil.AssertEqual(t, tt.message, err.Message)
			
			if tt.hasResolution {
				testutil.AssertTrue(t, err.Resolution != "", "Should have resolution")
			}
		})
	}
}

// TestNewValidationError tests validation error creation
func TestNewValidationError(t *testing.T) {
	tests := []struct {
		name       string
		code       ErrorCode
		message    string
		hasResolution bool
	}{
		{
			name:       "invalid app name",
			code:       ErrCodeValidationInvalidAppName,
			message:    "App name is invalid",
			hasResolution: true,
		},
		{
			name:       "invalid environment",
			code:       ErrCodeValidationInvalidEnvironment,
			message:    "Environment is invalid",
			hasResolution: true,
		},
		{
			name:       "invalid provider",
			code:       ErrCodeValidationInvalidProvider,
			message:    "Provider is invalid",
			hasResolution: true,
		},
		{
			name:       "invalid config",
			code:       ErrCodeValidationInvalidConfig,
			message:    "Config is invalid",
			hasResolution: true,
		},
		{
			name:       "missing required",
			code:       ErrCodeValidationMissingRequired,
			message:    "Required field missing",
			hasResolution: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationError(tt.code, tt.message, nil)

			testutil.AssertEqual(t, CategoryValidation, err.Category)
			testutil.AssertEqual(t, tt.code, err.Code)
			testutil.AssertEqual(t, tt.message, err.Message)
			
			if tt.hasResolution {
				testutil.AssertTrue(t, err.Resolution != "", "Should have resolution")
			}
		})
	}
}

// TestNewTerraformError tests Terraform error creation
func TestNewTerraformError(t *testing.T) {
	tests := []struct {
		name       string
		code       ErrorCode
		message    string
		hasResolution bool
	}{
		{
			name:       "init failed",
			code:       ErrCodeTerraformInitFailed,
			message:    "Terraform init failed",
			hasResolution: true,
		},
		{
			name:       "plan failed",
			code:       ErrCodeTerraformPlanFailed,
			message:    "Terraform plan failed",
			hasResolution: true,
		},
		{
			name:       "apply failed",
			code:       ErrCodeTerraformApplyFailed,
			message:    "Terraform apply failed",
			hasResolution: true,
		},
		{
			name:       "destroy failed",
			code:       ErrCodeTerraformDestroyFailed,
			message:    "Terraform destroy failed",
			hasResolution: true,
		},
		{
			name:       "output failed",
			code:       ErrCodeTerraformOutputFailed,
			message:    "Terraform output failed",
			hasResolution: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewTerraformError(tt.code, tt.message, nil)

			testutil.AssertEqual(t, CategoryTerraform, err.Category)
			testutil.AssertEqual(t, tt.code, err.Code)
			testutil.AssertEqual(t, tt.message, err.Message)
			
			if tt.hasResolution {
				testutil.AssertTrue(t, err.Resolution != "", "Should have resolution")
			}
		})
	}
}

// TestNewHelmError tests Helm error creation
func TestNewHelmError(t *testing.T) {
	tests := []struct {
		name       string
		code       ErrorCode
		message    string
		hasResolution bool
	}{
		{
			name:       "install failed",
			code:       ErrCodeHelmInstallFailed,
			message:    "Helm install failed",
			hasResolution: true,
		},
		{
			name:       "upgrade failed",
			code:       ErrCodeHelmUpgradeFailed,
			message:    "Helm upgrade failed",
			hasResolution: true,
		},
		{
			name:       "uninstall failed",
			code:       ErrCodeHelmUninstallFailed,
			message:    "Helm uninstall failed",
			hasResolution: true,
		},
		{
			name:       "status failed",
			code:       ErrCodeHelmStatusFailed,
			message:    "Helm status failed",
			hasResolution: true,
		},
		{
			name:       "pod not ready",
			code:       ErrCodeHelmPodNotReady,
			message:    "Pod not ready",
			hasResolution: true,
		},
		{
			name:       "timeout",
			code:       ErrCodeHelmTimeout,
			message:    "Helm operation timed out",
			hasResolution: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewHelmError(tt.code, tt.message, nil)

			testutil.AssertEqual(t, CategoryHelm, err.Category)
			testutil.AssertEqual(t, tt.code, err.Code)
			testutil.AssertEqual(t, tt.message, err.Message)
			
			if tt.hasResolution {
				testutil.AssertTrue(t, err.Resolution != "", "Should have resolution")
			}
		})
	}
}

// TestNewNetworkError tests network error creation
func TestNewNetworkError(t *testing.T) {
	tests := []struct {
		name       string
		code       ErrorCode
		message    string
		hasResolution bool
	}{
		{
			name:       "connection failed",
			code:       ErrCodeNetworkConnectionFailed,
			message:    "Connection failed",
			hasResolution: true,
		},
		{
			name:       "timeout",
			code:       ErrCodeNetworkTimeout,
			message:    "Network timeout",
			hasResolution: true,
		},
		{
			name:       "dns failed",
			code:       ErrCodeNetworkDNSFailed,
			message:    "DNS resolution failed",
			hasResolution: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewNetworkError(tt.code, tt.message, nil)

			testutil.AssertEqual(t, CategoryNetwork, err.Category)
			testutil.AssertEqual(t, tt.code, err.Code)
			testutil.AssertEqual(t, tt.message, err.Message)
			
			if tt.hasResolution {
				testutil.AssertTrue(t, err.Resolution != "", "Should have resolution")
			}
		})
	}
}

// TestNewConfigError tests configuration error creation
func TestNewConfigError(t *testing.T) {
	tests := []struct {
		name       string
		code       ErrorCode
		message    string
		hasResolution bool
	}{
		{
			name:       "file not found",
			code:       ErrCodeConfigFileNotFound,
			message:    "Config file not found",
			hasResolution: true,
		},
		{
			name:       "parse failed",
			code:       ErrCodeConfigParseFailed,
			message:    "Config parse failed",
			hasResolution: true,
		},
		{
			name:       "invalid format",
			code:       ErrCodeConfigInvalidFormat,
			message:    "Config format invalid",
			hasResolution: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewConfigError(tt.code, tt.message, nil)

			testutil.AssertEqual(t, CategoryConfiguration, err.Category)
			testutil.AssertEqual(t, tt.code, err.Code)
			testutil.AssertEqual(t, tt.message, err.Message)
			
			if tt.hasResolution {
				testutil.AssertTrue(t, err.Resolution != "", "Should have resolution")
			}
		})
	}
}

// TestNewStateError tests state error creation
func TestNewStateError(t *testing.T) {
	tests := []struct {
		name       string
		code       ErrorCode
		message    string
		hasResolution bool
	}{
		{
			name:       "state locked",
			code:       ErrCodeStateLocked,
			message:    "State is locked",
			hasResolution: true,
		},
		{
			name:       "state not found",
			code:       ErrCodeStateNotFound,
			message:    "State not found",
			hasResolution: true,
		},
		{
			name:       "state corrupted",
			code:       ErrCodeStateCorrupted,
			message:    "State is corrupted",
			hasResolution: true,
		},
		{
			name:       "access denied",
			code:       ErrCodeStateAccessDenied,
			message:    "State access denied",
			hasResolution: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewStateError(tt.code, tt.message, nil)

			testutil.AssertEqual(t, CategoryState, err.Category)
			testutil.AssertEqual(t, tt.code, err.Code)
			testutil.AssertEqual(t, tt.message, err.Message)
			
			if tt.hasResolution {
				testutil.AssertTrue(t, err.Resolution != "", "Should have resolution")
			}
		})
	}
}

// TestNewUnknownError tests unknown error creation
func TestNewUnknownError(t *testing.T) {
	cause := errors.New("unexpected error")
	err := NewUnknownError("Something went wrong", cause)

	testutil.AssertEqual(t, CategoryUnknown, err.Category)
	testutil.AssertEqual(t, ErrCodeUnknown, err.Code)
	testutil.AssertEqual(t, "Something went wrong", err.Message)
	testutil.AssertEqual(t, cause, err.Cause)
	testutil.AssertTrue(t, err.Resolution != "", "Should have resolution")
}

// TestErrorChaining tests error chaining with method calls
func TestErrorChaining(t *testing.T) {
	err := NewValidationError(
		ErrCodeValidationInvalidAppName,
		"Invalid app name",
		nil,
	).WithDetails("App name contains uppercase").WithResolution("Use lowercase only").WithLogPath("/var/log/app.log")

	testutil.AssertEqual(t, "Invalid app name", err.Message)
	testutil.AssertEqual(t, "App name contains uppercase", err.Details)
	testutil.AssertEqual(t, "Use lowercase only", err.Resolution)
	testutil.AssertEqual(t, "/var/log/app.log", err.LogPath)
}

// TestErrorCodeRanges tests that error codes are in correct ranges
func TestErrorCodeRanges(t *testing.T) {
	tests := []struct {
		name      string
		code      ErrorCode
		minRange  int
		maxRange  int
	}{
		{name: "auth codes in 1000-1099", code: ErrCodeAuthInvalidCredentials, minRange: 1000, maxRange: 1099},
		{name: "validation codes in 1100-1199", code: ErrCodeValidationInvalidAppName, minRange: 1100, maxRange: 1199},
		{name: "terraform codes in 1200-1299", code: ErrCodeTerraformInitFailed, minRange: 1200, maxRange: 1299},
		{name: "helm codes in 1300-1399", code: ErrCodeHelmInstallFailed, minRange: 1300, maxRange: 1399},
		{name: "network codes in 1400-1499", code: ErrCodeNetworkConnectionFailed, minRange: 1400, maxRange: 1499},
		{name: "config codes in 1500-1599", code: ErrCodeConfigFileNotFound, minRange: 1500, maxRange: 1599},
		{name: "state codes in 1600-1699", code: ErrCodeStateLocked, minRange: 1600, maxRange: 1699},
		{name: "unknown codes in 9000-9999", code: ErrCodeUnknown, minRange: 9000, maxRange: 9999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codeInt := int(tt.code)
			testutil.AssertTrue(t, codeInt >= tt.minRange && codeInt <= tt.maxRange,
				"Code should be in range")
		})
	}
}

// TestErrorImplementsErrorInterface tests that CLIError implements error interface
func TestErrorImplementsErrorInterface(t *testing.T) {
	var _ error = &CLIError{}
	var _ error = NewCLIError(CategoryValidation, ErrCodeValidationInvalidAppName, "test", nil)
}

// TestErrorWrapping tests error wrapping with errors.Is and errors.As
func TestErrorWrapping(t *testing.T) {
	cause := errors.New("underlying error")
	cliErr := NewCLIError(CategoryValidation, ErrCodeValidationInvalidAppName, "Invalid app name", cause)

	// Test errors.Is
	testutil.AssertTrue(t, errors.Is(cliErr, cause), "Should be able to unwrap to cause")

	// Test errors.As
	var target *CLIError
	testutil.AssertTrue(t, errors.As(cliErr, &target), "Should be able to extract CLIError")
	testutil.AssertEqual(t, cliErr, target)
}

// TestErrorMessageFormatting tests that error messages are properly formatted
func TestErrorMessageFormatting(t *testing.T) {
	tests := []struct {
		name     string
		err      *CLIError
		notContains []string
	}{
		{
			name: "no trailing newlines in message",
			err: NewCLIError(CategoryValidation, ErrCodeValidationInvalidAppName, "Invalid app name", nil),
			notContains: []string{"\n\n", "  "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := tt.err.Error()
			formatted := tt.err.Format()
			
			for _, substr := range tt.notContains {
				testutil.AssertFalse(t, strings.Contains(errStr, substr), "Error string should not contain: "+substr)
			}
			
			// Formatted output can have newlines, but check it's not empty
			testutil.AssertTrue(t, len(formatted) > 0, "Formatted output should not be empty")
		})
	}
}
