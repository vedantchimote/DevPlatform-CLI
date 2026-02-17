package helm

import (
	"fmt"
	"strings"
)

// HelmError represents a structured Helm error
type HelmError struct {
	Operation string   // The Helm operation that failed (install, upgrade, uninstall, etc.)
	Release   string   // The release name
	Namespace string   // The namespace
	Message   string   // The error message
	Output    string   // The full command output
	Events    []string // Kubernetes events (if available)
	ExitCode  int      // The exit code from helm command
}

// Error implements the error interface
func (e *HelmError) Error() string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("Helm %s failed for release '%s' in namespace '%s'", 
		e.Operation, e.Release, e.Namespace))
	
	if e.Message != "" {
		sb.WriteString(fmt.Sprintf(": %s", e.Message))
	}
	
	if len(e.Events) > 0 {
		sb.WriteString("\n\nKubernetes Events:\n")
		for _, event := range e.Events {
			sb.WriteString(fmt.Sprintf("  - %s\n", event))
		}
	}
	
	return sb.String()
}

// NewHelmError creates a new HelmError
func NewHelmError(operation, release, namespace, message, output string, exitCode int) *HelmError {
	return &HelmError{
		Operation: operation,
		Release:   release,
		Namespace: namespace,
		Message:   message,
		Output:    output,
		ExitCode:  exitCode,
	}
}

// WithEvents adds Kubernetes events to the error
func (e *HelmError) WithEvents(events []string) *HelmError {
	e.Events = events
	return e
}

// ParseHelmError parses helm command output to extract meaningful error information
func ParseHelmError(output string) string {
	lines := strings.Split(output, "\n")
	
	// Look for common error patterns
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines
		if line == "" {
			continue
		}
		
		// Look for "Error:" prefix
		if strings.HasPrefix(line, "Error:") {
			return strings.TrimPrefix(line, "Error:")
		}
		
		// Look for "error:" prefix (lowercase)
		if strings.HasPrefix(line, "error:") {
			return strings.TrimPrefix(line, "error:")
		}
		
		// Look for "FATAL" prefix
		if strings.Contains(strings.ToUpper(line), "FATAL") {
			return line
		}
	}
	
	// If no specific error pattern found, return the last non-empty line
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			return line
		}
	}
	
	return "Unknown error"
}

// IsReleaseNotFound checks if the error is due to release not found
func IsReleaseNotFound(err error) bool {
	if err == nil {
		return false
	}
	
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "not found") || 
		   strings.Contains(errMsg, "release: not found")
}

// IsTimeout checks if the error is due to timeout
func IsTimeout(err error) bool {
	if err == nil {
		return false
	}
	
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "timeout") || 
		   strings.Contains(errMsg, "timed out") ||
		   strings.Contains(errMsg, "deadline exceeded")
}

// IsResourceConflict checks if the error is due to resource conflict
func IsResourceConflict(err error) bool {
	if err == nil {
		return false
	}
	
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "already exists") || 
		   strings.Contains(errMsg, "conflict") ||
		   strings.Contains(errMsg, "already owned")
}

// IsInvalidChart checks if the error is due to invalid chart
func IsInvalidChart(err error) bool {
	if err == nil {
		return false
	}
	
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "chart not found") || 
		   strings.Contains(errMsg, "invalid chart") ||
		   strings.Contains(errMsg, "chart.yaml") ||
		   strings.Contains(errMsg, "no chart found")
}

// IsInvalidValues checks if the error is due to invalid values
func IsInvalidValues(err error) bool {
	if err == nil {
		return false
	}
	
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "invalid values") || 
		   strings.Contains(errMsg, "values.yaml") ||
		   strings.Contains(errMsg, "yaml") ||
		   strings.Contains(errMsg, "parse error")
}

// GetErrorCategory returns a category for the error
func GetErrorCategory(err error) string {
	if err == nil {
		return "unknown"
	}
	
	if IsReleaseNotFound(err) {
		return "release_not_found"
	}
	
	if IsTimeout(err) {
		return "timeout"
	}
	
	if IsResourceConflict(err) {
		return "resource_conflict"
	}
	
	if IsInvalidChart(err) {
		return "invalid_chart"
	}
	
	if IsInvalidValues(err) {
		return "invalid_values"
	}
	
	return "helm_error"
}

// GetErrorResolution returns a suggested resolution for the error
func GetErrorResolution(err error) string {
	if err == nil {
		return ""
	}
	
	category := GetErrorCategory(err)
	
	switch category {
	case "release_not_found":
		return "The Helm release does not exist. Use 'helm list' to check existing releases."
	case "timeout":
		return "The operation timed out. Check pod status with 'kubectl get pods' and increase timeout if needed."
	case "resource_conflict":
		return "A resource already exists. Use 'helm upgrade' instead of 'helm install', or delete the existing resources."
	case "invalid_chart":
		return "The chart is invalid or not found. Check the chart path and ensure Chart.yaml exists."
	case "invalid_values":
		return "The values file is invalid. Check YAML syntax and ensure all required values are provided."
	default:
		return "Check the error message above for details. Use 'helm --debug' for more information."
	}
}
