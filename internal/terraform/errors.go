package terraform

import (
	"fmt"
	"strings"
)

// TerraformError represents a structured Terraform error
type TerraformError struct {
	Command    string
	ExitCode   int
	Message    string
	Output     string
	Category   ErrorCategory
	Suggestion string
}

// Error implements the error interface
func (e *TerraformError) Error() string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("Terraform %s failed", e.Command))
	if e.ExitCode > 0 {
		sb.WriteString(fmt.Sprintf(" (exit code %d)", e.ExitCode))
	}
	sb.WriteString(fmt.Sprintf(": %s", e.Message))
	
	if e.Suggestion != "" {
		sb.WriteString(fmt.Sprintf("\n\nSuggestion: %s", e.Suggestion))
	}
	
	if e.Output != "" {
		sb.WriteString(fmt.Sprintf("\n\nTerraform Output:\n%s", e.Output))
	}
	
	return sb.String()
}

// ErrorCategory represents the category of a Terraform error
type ErrorCategory int

const (
	ErrorCategoryUnknown ErrorCategory = iota
	ErrorCategoryConfiguration
	ErrorCategoryState
	ErrorCategoryProvider
	ErrorCategoryResource
	ErrorCategoryValidation
	ErrorCategoryPermission
	ErrorCategoryNetwork
)

// String returns the string representation of the error category
func (c ErrorCategory) String() string {
	switch c {
	case ErrorCategoryConfiguration:
		return "Configuration Error"
	case ErrorCategoryState:
		return "State Error"
	case ErrorCategoryProvider:
		return "Provider Error"
	case ErrorCategoryResource:
		return "Resource Error"
	case ErrorCategoryValidation:
		return "Validation Error"
	case ErrorCategoryPermission:
		return "Permission Error"
	case ErrorCategoryNetwork:
		return "Network Error"
	default:
		return "Unknown Error"
	}
}

// ParseTerraformError parses a Terraform error and returns a structured error
func ParseTerraformError(command string, err error, output string) *TerraformError {
	if err == nil {
		return nil
	}
	
	tfErr := &TerraformError{
		Command:  command,
		Message:  err.Error(),
		Output:   output,
		Category: ErrorCategoryUnknown,
	}
	
	// Analyze the error output to categorize and provide suggestions
	outputLower := strings.ToLower(output)
	errorLower := strings.ToLower(err.Error())
	
	// Configuration errors
	if strings.Contains(outputLower, "invalid configuration") ||
		strings.Contains(outputLower, "syntax error") ||
		strings.Contains(outputLower, "missing required argument") {
		tfErr.Category = ErrorCategoryConfiguration
		tfErr.Suggestion = "Check your Terraform configuration files for syntax errors or missing required arguments"
	}
	
	// State errors
	if strings.Contains(outputLower, "state lock") ||
		strings.Contains(outputLower, "locked") ||
		strings.Contains(errorLower, "state lock") {
		tfErr.Category = ErrorCategoryState
		tfErr.Suggestion = "The state is locked by another operation. Wait for it to complete or use 'terraform force-unlock' if the lock is stale"
	}
	
	if strings.Contains(outputLower, "no state file") ||
		strings.Contains(outputLower, "state file is empty") {
		tfErr.Category = ErrorCategoryState
		tfErr.Suggestion = "Initialize Terraform with 'terraform init' before running other commands"
	}
	
	// Provider errors
	if strings.Contains(outputLower, "provider") && 
		(strings.Contains(outputLower, "not found") || strings.Contains(outputLower, "failed to install")) {
		tfErr.Category = ErrorCategoryProvider
		tfErr.Suggestion = "Run 'terraform init' to download required providers"
	}
	
	if strings.Contains(outputLower, "authentication") ||
		strings.Contains(outputLower, "credentials") ||
		strings.Contains(outputLower, "unauthorized") {
		tfErr.Category = ErrorCategoryPermission
		tfErr.Suggestion = "Check your cloud provider credentials and ensure they have the necessary permissions"
	}
	
	// Resource errors
	if strings.Contains(outputLower, "already exists") ||
		strings.Contains(outputLower, "duplicate") {
		tfErr.Category = ErrorCategoryResource
		tfErr.Suggestion = "A resource with this name already exists. Choose a different name or import the existing resource"
	}
	
	if strings.Contains(outputLower, "not found") && 
		!strings.Contains(outputLower, "provider") {
		tfErr.Category = ErrorCategoryResource
		tfErr.Suggestion = "The resource does not exist. It may have been deleted outside of Terraform"
	}
	
	// Validation errors
	if strings.Contains(outputLower, "invalid") ||
		strings.Contains(outputLower, "validation failed") {
		tfErr.Category = ErrorCategoryValidation
		tfErr.Suggestion = "Check the input values and ensure they meet the required format and constraints"
	}
	
	// Network errors
	if strings.Contains(outputLower, "timeout") ||
		strings.Contains(outputLower, "connection refused") ||
		strings.Contains(outputLower, "network") {
		tfErr.Category = ErrorCategoryNetwork
		tfErr.Suggestion = "Check your network connection and firewall settings"
	}
	
	// Permission errors
	if strings.Contains(outputLower, "access denied") ||
		strings.Contains(outputLower, "forbidden") ||
		strings.Contains(outputLower, "insufficient permissions") {
		tfErr.Category = ErrorCategoryPermission
		tfErr.Suggestion = "Ensure your cloud provider credentials have the necessary permissions for this operation"
	}
	
	return tfErr
}

// IsStateLockError checks if an error is a state lock error
func IsStateLockError(err error) bool {
	if err == nil {
		return false
	}
	
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "state lock") || strings.Contains(errStr, "locked")
}

// IsPermissionError checks if an error is a permission error
func IsPermissionError(err error) bool {
	if err == nil {
		return false
	}
	
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "access denied") ||
		strings.Contains(errStr, "forbidden") ||
		strings.Contains(errStr, "unauthorized") ||
		strings.Contains(errStr, "insufficient permissions") ||
		strings.Contains(errStr, "authentication")
}

// IsResourceNotFoundError checks if an error is a resource not found error
func IsResourceNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "not found") && !strings.Contains(errStr, "provider")
}

// ExtractErrorMessage extracts a clean error message from Terraform output
func ExtractErrorMessage(output string) string {
	lines := strings.Split(output, "\n")
	
	var errorLines []string
	inError := false
	
	for _, line := range lines {
		lineLower := strings.ToLower(line)
		
		// Start capturing when we see "Error:"
		if strings.Contains(lineLower, "error:") {
			inError = true
		}
		
		// Capture error lines
		if inError {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				errorLines = append(errorLines, trimmed)
			}
			
			// Stop at empty line after error
			if trimmed == "" && len(errorLines) > 0 {
				break
			}
		}
	}
	
	if len(errorLines) > 0 {
		return strings.Join(errorLines, "\n")
	}
	
	// If no error section found, return last non-empty line
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed != "" {
			return trimmed
		}
	}
	
	return output
}
