package terraform

import (
	"context"
	"encoding/json"
	"fmt"
)

// TerraformOutputs represents the structured outputs from Terraform
type TerraformOutputs struct {
	// Network outputs
	VPCID        string
	VNetID       string
	SubnetIDs    []string
	SecurityGroupIDs []string
	NSGID        string
	
	// Database outputs
	RDSEndpoint     string
	RDSPort         string
	DatabaseEndpoint string
	DatabasePort    string
	SecretARN       string
	KeyVaultSecretID string
	
	// Kubernetes outputs
	Namespace          string
	ServiceAccountName string
	ClusterName        string
	
	// Additional metadata
	Provider    string
	Environment string
	AppName     string
}

// OutputParser handles parsing of Terraform outputs
type OutputParser struct {
	executor TerraformExecutor
}

// NewOutputParser creates a new output parser
func NewOutputParser(executor TerraformExecutor) *OutputParser {
	return &OutputParser{
		executor: executor,
	}
}

// ParseOutputs parses all Terraform outputs into a structured format
func (p *OutputParser) ParseOutputs(ctx context.Context, workingDir string) (*TerraformOutputs, error) {
	// Get all outputs as JSON
	outputJSON, err := p.executor.Output(ctx, workingDir, "-json")
	if err != nil {
		return nil, fmt.Errorf("failed to get terraform outputs: %w", err)
	}
	
	// Parse JSON into map
	var outputs map[string]OutputValue
	if err := json.Unmarshal([]byte(outputJSON), &outputs); err != nil {
		return nil, fmt.Errorf("failed to parse terraform output JSON: %w", err)
	}
	
	// Extract values into structured format
	result := &TerraformOutputs{}
	
	// Network outputs
	if val, ok := outputs["vpc_id"]; ok {
		result.VPCID = getStringValue(val)
	}
	if val, ok := outputs["vnet_id"]; ok {
		result.VNetID = getStringValue(val)
	}
	if val, ok := outputs["subnet_ids"]; ok {
		result.SubnetIDs = getStringSliceValue(val)
	}
	if val, ok := outputs["security_group_ids"]; ok {
		result.SecurityGroupIDs = getStringSliceValue(val)
	}
	if val, ok := outputs["nsg_id"]; ok {
		result.NSGID = getStringValue(val)
	}
	
	// Database outputs
	if val, ok := outputs["rds_endpoint"]; ok {
		result.RDSEndpoint = getStringValue(val)
	}
	if val, ok := outputs["rds_port"]; ok {
		result.RDSPort = getStringValue(val)
	}
	if val, ok := outputs["database_endpoint"]; ok {
		result.DatabaseEndpoint = getStringValue(val)
	}
	if val, ok := outputs["database_port"]; ok {
		result.DatabasePort = getStringValue(val)
	}
	if val, ok := outputs["secret_arn"]; ok {
		result.SecretARN = getStringValue(val)
	}
	if val, ok := outputs["keyvault_secret_id"]; ok {
		result.KeyVaultSecretID = getStringValue(val)
	}
	
	// Kubernetes outputs
	if val, ok := outputs["namespace"]; ok {
		result.Namespace = getStringValue(val)
	}
	if val, ok := outputs["service_account_name"]; ok {
		result.ServiceAccountName = getStringValue(val)
	}
	if val, ok := outputs["cluster_name"]; ok {
		result.ClusterName = getStringValue(val)
	}
	
	// Metadata
	if val, ok := outputs["provider"]; ok {
		result.Provider = getStringValue(val)
	}
	if val, ok := outputs["environment"]; ok {
		result.Environment = getStringValue(val)
	}
	if val, ok := outputs["app_name"]; ok {
		result.AppName = getStringValue(val)
	}
	
	return result, nil
}

// GetOutput retrieves a single output value by name
func (p *OutputParser) GetOutput(ctx context.Context, workingDir string, outputName string) (string, error) {
	return p.executor.Output(ctx, workingDir, outputName)
}

// OutputValue represents a Terraform output value
type OutputValue struct {
	Sensitive bool        `json:"sensitive"`
	Type      interface{} `json:"type"`
	Value     interface{} `json:"value"`
}

// getStringValue extracts a string value from an OutputValue
func getStringValue(output OutputValue) string {
	if str, ok := output.Value.(string); ok {
		return str
	}
	return ""
}

// getStringSliceValue extracts a string slice value from an OutputValue
func getStringSliceValue(output OutputValue) []string {
	if slice, ok := output.Value.([]interface{}); ok {
		result := make([]string, 0, len(slice))
		for _, item := range slice {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	}
	return []string{}
}

// GetDatabaseEndpoint returns the appropriate database endpoint based on provider
func (o *TerraformOutputs) GetDatabaseEndpoint() string {
	if o.RDSEndpoint != "" {
		return o.RDSEndpoint
	}
	return o.DatabaseEndpoint
}

// GetDatabasePort returns the appropriate database port based on provider
func (o *TerraformOutputs) GetDatabasePort() string {
	if o.RDSPort != "" {
		return o.RDSPort
	}
	return o.DatabasePort
}

// GetNetworkID returns the appropriate network ID based on provider
func (o *TerraformOutputs) GetNetworkID() string {
	if o.VPCID != "" {
		return o.VPCID
	}
	return o.VNetID
}
