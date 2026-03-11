package errors

import (
	"fmt"
)

// ErrorCategory represents the category of an error
type ErrorCategory string

const (
	// Authentication errors (1000-1099)
	CategoryAuthentication ErrorCategory = "authentication"
	
	// Validation errors (1100-1199)
	CategoryValidation ErrorCategory = "validation"
	
	// Terraform errors (1200-1299)
	CategoryTerraform ErrorCategory = "terraform"
	
	// Helm errors (1300-1399)
	CategoryHelm ErrorCategory = "helm"
	
	// Network errors (1400-1499)
	CategoryNetwork ErrorCategory = "network"
	
	// Configuration errors (1500-1599)
	CategoryConfiguration ErrorCategory = "configuration"
	
	// State errors (1600-1699)
	CategoryState ErrorCategory = "state"
	
	// Unknown errors (9000-9999)
	CategoryUnknown ErrorCategory = "unknown"
)

// ErrorCode represents a specific error code
type ErrorCode int

const (
	// Authentication error codes (1000-1099)
	ErrCodeAuthInvalidCredentials ErrorCode = 1001
	ErrCodeAuthExpiredCredentials ErrorCode = 1002
	ErrCodeAuthMissingCredentials ErrorCode = 1003
	ErrCodeAuthPermissionDenied   ErrorCode = 1004
	
	// Validation error codes (1100-1199)
	ErrCodeValidationInvalidAppName     ErrorCode = 1101
	ErrCodeValidationInvalidEnvironment ErrorCode = 1102
	ErrCodeValidationInvalidProvider    ErrorCode = 1103
	ErrCodeValidationInvalidConfig      ErrorCode = 1104
	ErrCodeValidationMissingRequired    ErrorCode = 1105
	
	// Terraform error codes (1200-1299)
	ErrCodeTerraformInitFailed    ErrorCode = 1201
	ErrCodeTerraformPlanFailed    ErrorCode = 1202
	ErrCodeTerraformApplyFailed   ErrorCode = 1203
	ErrCodeTerraformDestroyFailed ErrorCode = 1204
	ErrCodeTerraformOutputFailed  ErrorCode = 1205
	
	// Helm error codes (1300-1399)
	ErrCodeHelmInstallFailed   ErrorCode = 1301
	ErrCodeHelmUpgradeFailed   ErrorCode = 1302
	ErrCodeHelmUninstallFailed ErrorCode = 1303
	ErrCodeHelmStatusFailed    ErrorCode = 1304
	ErrCodeHelmPodNotReady     ErrorCode = 1305
	ErrCodeHelmTimeout         ErrorCode = 1306
	
	// Network error codes (1400-1499)
	ErrCodeNetworkConnectionFailed ErrorCode = 1401
	ErrCodeNetworkTimeout          ErrorCode = 1402
	ErrCodeNetworkDNSFailed        ErrorCode = 1403
	
	// Configuration error codes (1500-1599)
	ErrCodeConfigFileNotFound  ErrorCode = 1501
	ErrCodeConfigParseFailed   ErrorCode = 1502
	ErrCodeConfigInvalidFormat ErrorCode = 1503
	
	// State error codes (1600-1699)
	ErrCodeStateLocked        ErrorCode = 1601
	ErrCodeStateNotFound      ErrorCode = 1602
	ErrCodeStateCorrupted     ErrorCode = 1603
	ErrCodeStateAccessDenied  ErrorCode = 1604
	
	// Unknown error codes (9000-9999)
	ErrCodeUnknown ErrorCode = 9000
)

// CLIError represents a structured error with category, code, and resolution
type CLIError struct {
	Category   ErrorCategory
	Code       ErrorCode
	Message    string
	Details    string
	Resolution string
	Cause      error
	LogPath    string
}

// Error implements the error interface
func (e *CLIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s:%d] %s: %v", e.Category, e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s:%d] %s", e.Category, e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *CLIError) Unwrap() error {
	return e.Cause
}

// Format returns a formatted error message for display
func (e *CLIError) Format() string {
	output := fmt.Sprintf("\n❌ Error [%s:%d]\n", e.Category, e.Code)
	output += fmt.Sprintf("Message: %s\n", e.Message)
	
	if e.Details != "" {
		output += fmt.Sprintf("\nDetails:\n%s\n", e.Details)
	}
	
	if e.Resolution != "" {
		output += fmt.Sprintf("\n💡 Resolution:\n%s\n", e.Resolution)
	}
	
	if e.LogPath != "" {
		output += fmt.Sprintf("\n📝 Log file: %s\n", e.LogPath)
	}
	
	return output
}

// NewCLIError creates a new CLI error
func NewCLIError(category ErrorCategory, code ErrorCode, message string, cause error) *CLIError {
	return &CLIError{
		Category: category,
		Code:     code,
		Message:  message,
		Cause:    cause,
	}
}

// WithDetails adds details to the error
func (e *CLIError) WithDetails(details string) *CLIError {
	e.Details = details
	return e
}

// WithResolution adds resolution guidance to the error
func (e *CLIError) WithResolution(resolution string) *CLIError {
	e.Resolution = resolution
	return e
}

// WithLogPath adds the log file path to the error
func (e *CLIError) WithLogPath(logPath string) *CLIError {
	e.LogPath = logPath
	return e
}

// Authentication errors
func NewAuthError(code ErrorCode, message string, cause error) *CLIError {
	err := NewCLIError(CategoryAuthentication, code, message, cause)
	
	switch code {
	case ErrCodeAuthInvalidCredentials:
		err.WithResolution("Verify your cloud provider credentials are correct and have not expired.")
	case ErrCodeAuthExpiredCredentials:
		err.WithResolution("Refresh your cloud provider credentials and try again.")
	case ErrCodeAuthMissingCredentials:
		err.WithResolution("Configure your cloud provider credentials using the CLI or environment variables.")
	case ErrCodeAuthPermissionDenied:
		err.WithResolution("Ensure your credentials have the necessary permissions for this operation.")
	}
	
	return err
}

// Validation errors
func NewValidationError(code ErrorCode, message string, cause error) *CLIError {
	err := NewCLIError(CategoryValidation, code, message, cause)
	
	switch code {
	case ErrCodeValidationInvalidAppName:
		err.WithResolution("App name must be 3-32 characters, lowercase alphanumeric and hyphens only.")
	case ErrCodeValidationInvalidEnvironment:
		err.WithResolution("Environment must be one of: dev, staging, prod")
	case ErrCodeValidationInvalidProvider:
		err.WithResolution("Provider must be one of: aws, azure")
	case ErrCodeValidationInvalidConfig:
		err.WithResolution("Check your configuration file for syntax errors and required fields.")
	case ErrCodeValidationMissingRequired:
		err.WithResolution("Provide all required flags or configuration values.")
	}
	
	return err
}

// Terraform errors
func NewTerraformError(code ErrorCode, message string, cause error) *CLIError {
	err := NewCLIError(CategoryTerraform, code, message, cause)
	
	switch code {
	case ErrCodeTerraformInitFailed:
		err.WithResolution("Check your Terraform configuration and backend settings. Ensure Terraform is installed.")
	case ErrCodeTerraformPlanFailed:
		err.WithResolution("Review the Terraform plan output for configuration errors.")
	case ErrCodeTerraformApplyFailed:
		err.WithResolution("Check the Terraform output for resource provisioning errors. You may need to run 'terraform destroy' to clean up partial resources.")
	case ErrCodeTerraformDestroyFailed:
		err.WithResolution("Some resources may still exist. Check the Terraform output and manually delete resources if necessary.")
	case ErrCodeTerraformOutputFailed:
		err.WithResolution("Ensure Terraform state exists and contains the requested output values.")
	}
	
	return err
}

// Helm errors
func NewHelmError(code ErrorCode, message string, cause error) *CLIError {
	err := NewCLIError(CategoryHelm, code, message, cause)
	
	switch code {
	case ErrCodeHelmInstallFailed:
		err.WithResolution("Check the Helm chart configuration and Kubernetes cluster connectivity. Review pod logs for application errors.")
	case ErrCodeHelmUpgradeFailed:
		err.WithResolution("Review the Helm upgrade output. You may need to rollback using 'helm rollback'.")
	case ErrCodeHelmUninstallFailed:
		err.WithResolution("Some Kubernetes resources may still exist. Check with 'kubectl get all -n <namespace>' and delete manually if needed.")
	case ErrCodeHelmStatusFailed:
		err.WithResolution("Ensure the Helm release exists and Kubernetes cluster is accessible.")
	case ErrCodeHelmPodNotReady:
		err.WithResolution("Check pod logs with 'kubectl logs <pod-name> -n <namespace>' and Kubernetes events with 'kubectl get events -n <namespace>'.")
	case ErrCodeHelmTimeout:
		err.WithResolution("Increase the timeout value or check for resource constraints in the cluster.")
	}
	
	return err
}

// Network errors
func NewNetworkError(code ErrorCode, message string, cause error) *CLIError {
	err := NewCLIError(CategoryNetwork, code, message, cause)
	
	switch code {
	case ErrCodeNetworkConnectionFailed:
		err.WithResolution("Check your network connectivity and firewall settings.")
	case ErrCodeNetworkTimeout:
		err.WithResolution("Increase the timeout value or check your network connection.")
	case ErrCodeNetworkDNSFailed:
		err.WithResolution("Verify DNS settings and network configuration.")
	}
	
	return err
}

// Configuration errors
func NewConfigError(code ErrorCode, message string, cause error) *CLIError {
	err := NewCLIError(CategoryConfiguration, code, message, cause)
	
	switch code {
	case ErrCodeConfigFileNotFound:
		err.WithResolution("Create a .devplatform.yaml configuration file or specify the path with --config flag.")
	case ErrCodeConfigParseFailed:
		err.WithResolution("Check your configuration file for YAML syntax errors.")
	case ErrCodeConfigInvalidFormat:
		err.WithResolution("Ensure your configuration file follows the correct schema.")
	}
	
	return err
}

// State errors
func NewStateError(code ErrorCode, message string, cause error) *CLIError {
	err := NewCLIError(CategoryState, code, message, cause)
	
	switch code {
	case ErrCodeStateLocked:
		err.WithResolution("Wait for the other operation to complete or use 'terraform force-unlock <lock-id>' if the lock is stale.")
	case ErrCodeStateNotFound:
		err.WithResolution("Ensure the environment has been created with 'devplatform create' command.")
	case ErrCodeStateCorrupted:
		err.WithResolution("Restore from a backup or recreate the environment.")
	case ErrCodeStateAccessDenied:
		err.WithResolution("Ensure you have permissions to access the state storage (S3 bucket or Azure Storage).")
	}
	
	return err
}

// Unknown errors
func NewUnknownError(message string, cause error) *CLIError {
	return NewCLIError(CategoryUnknown, ErrCodeUnknown, message, cause).
		WithResolution("Check the log file for more details. If the issue persists, please report it.")
}
