package terraform

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

// TestOutputParserParseOutputs tests parsing Terraform outputs
func TestOutputParserParseOutputs(t *testing.T) {
	tests := []struct {
		name        string
		outputJSON  string
		wantOutputs *TerraformOutputs
		wantErr     bool
		errContains string
	}{
		{
			name: "parse AWS outputs",
			outputJSON: `{
				"vpc_id": {"sensitive": false, "type": "string", "value": "vpc-123"},
				"subnet_ids": {"sensitive": false, "type": ["list", "string"], "value": ["subnet-1", "subnet-2"]},
				"rds_endpoint": {"sensitive": false, "type": "string", "value": "db.example.com"},
				"rds_port": {"sensitive": false, "type": "string", "value": "5432"},
				"cluster_name": {"sensitive": false, "type": "string", "value": "my-cluster"},
				"namespace": {"sensitive": false, "type": "string", "value": "default"},
				"provider": {"sensitive": false, "type": "string", "value": "aws"},
				"environment": {"sensitive": false, "type": "string", "value": "dev"},
				"app_name": {"sensitive": false, "type": "string", "value": "myapp"}
			}`,
			wantOutputs: &TerraformOutputs{
				VPCID:       "vpc-123",
				SubnetIDs:   []string{"subnet-1", "subnet-2"},
				RDSEndpoint: "db.example.com",
				RDSPort:     "5432",
				ClusterName: "my-cluster",
				Namespace:   "default",
				Provider:    "aws",
				Environment: "dev",
				AppName:     "myapp",
			},
			wantErr: false,
		},
		{
			name: "parse Azure outputs",
			outputJSON: `{
				"vnet_id": {"sensitive": false, "type": "string", "value": "vnet-456"},
				"subnet_ids": {"sensitive": false, "type": ["list", "string"], "value": ["subnet-3", "subnet-4"]},
				"nsg_id": {"sensitive": false, "type": "string", "value": "nsg-789"},
				"database_endpoint": {"sensitive": false, "type": "string", "value": "postgres.azure.com"},
				"database_port": {"sensitive": false, "type": "string", "value": "5432"},
				"keyvault_secret_id": {"sensitive": false, "type": "string", "value": "secret-id-123"},
				"cluster_name": {"sensitive": false, "type": "string", "value": "aks-cluster"},
				"provider": {"sensitive": false, "type": "string", "value": "azure"},
				"environment": {"sensitive": false, "type": "string", "value": "prod"}
			}`,
			wantOutputs: &TerraformOutputs{
				VNetID:           "vnet-456",
				SubnetIDs:        []string{"subnet-3", "subnet-4"},
				NSGID:            "nsg-789",
				DatabaseEndpoint: "postgres.azure.com",
				DatabasePort:     "5432",
				KeyVaultSecretID: "secret-id-123",
				ClusterName:      "aks-cluster",
				Provider:         "azure",
				Environment:      "prod",
			},
			wantErr: false,
		},
		{
			name:       "parse empty outputs",
			outputJSON: `{}`,
			wantOutputs: &TerraformOutputs{
				SubnetIDs:        []string{},
				SecurityGroupIDs: []string{},
			},
			wantErr: false,
		},
		{
			name:        "invalid JSON",
			outputJSON:  `{invalid json}`,
			wantErr:     true,
			errContains: "failed to parse",
		},
		{
			name: "sensitive outputs",
			outputJSON: `{
				"vpc_id": {"sensitive": true, "type": "string", "value": "vpc-secret"},
				"secret_arn": {"sensitive": true, "type": "string", "value": "arn:aws:secretsmanager:us-east-1:123456789012:secret:mysecret"}
			}`,
			wantOutputs: &TerraformOutputs{
				VPCID:     "vpc-secret",
				SecretARN: "arn:aws:secretsmanager:us-east-1:123456789012:secret:mysecret",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockTerraformExecutor()
			mock.OutputFunc = func(ctx context.Context, workingDir string, outputName string) (string, error) {
				if outputName == "-json" {
					return tt.outputJSON, nil
				}
				return "", errors.New("unexpected output name")
			}

			parser := NewOutputParser(mock)
			outputs, err := parser.ParseOutputs(context.Background(), "/tmp/terraform")

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOutputs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				testutil.AssertContains(t, err.Error(), tt.errContains)
				return
			}

			if !tt.wantErr {
				// Verify specific fields
				testutil.AssertEqual(t, tt.wantOutputs.VPCID, outputs.VPCID)
				testutil.AssertEqual(t, tt.wantOutputs.VNetID, outputs.VNetID)
				testutil.AssertEqual(t, tt.wantOutputs.RDSEndpoint, outputs.RDSEndpoint)
				testutil.AssertEqual(t, tt.wantOutputs.DatabaseEndpoint, outputs.DatabaseEndpoint)
				testutil.AssertEqual(t, tt.wantOutputs.ClusterName, outputs.ClusterName)
				testutil.AssertEqual(t, tt.wantOutputs.Provider, outputs.Provider)
				testutil.AssertEqual(t, tt.wantOutputs.Environment, outputs.Environment)
				testutil.AssertEqual(t, tt.wantOutputs.AppName, outputs.AppName)
			}
		})
	}
}

// TestOutputParserGetOutput tests getting a single output value
func TestOutputParserGetOutput(t *testing.T) {
	tests := []struct {
		name        string
		outputName  string
		mockReturn  string
		mockErr     error
		wantOutput  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "get vpc_id",
			outputName: "vpc_id",
			mockReturn: "vpc-123",
			mockErr:    nil,
			wantOutput: "vpc-123",
			wantErr:    false,
		},
		{
			name:       "get cluster_name",
			outputName: "cluster_name",
			mockReturn: "my-cluster",
			mockErr:    nil,
			wantOutput: "my-cluster",
			wantErr:    false,
		},
		{
			name:        "output not found",
			outputName:  "nonexistent",
			mockReturn:  "",
			mockErr:     errors.New("output not found"),
			wantErr:     true,
			errContains: "output not found",
		},
		{
			name:       "empty output value",
			outputName: "empty",
			mockReturn: "",
			mockErr:    nil,
			wantOutput: "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockTerraformExecutor()
			mock.OutputFunc = func(ctx context.Context, workingDir string, outputName string) (string, error) {
				return tt.mockReturn, tt.mockErr
			}

			parser := NewOutputParser(mock)
			output, err := parser.GetOutput(context.Background(), "/tmp/terraform", tt.outputName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				testutil.AssertEqual(t, tt.wantOutput, output)
			}

			if tt.wantErr && tt.errContains != "" {
				testutil.AssertContains(t, err.Error(), tt.errContains)
			}
		})
	}
}

// TestTerraformOutputsHelperMethods tests the helper methods on TerraformOutputs
func TestTerraformOutputsHelperMethods(t *testing.T) {
	tests := []struct {
		name     string
		outputs  *TerraformOutputs
		testFunc func(*testing.T, *TerraformOutputs)
	}{
		{
			name: "GetDatabaseEndpoint - AWS",
			outputs: &TerraformOutputs{
				RDSEndpoint:      "rds.aws.com",
				DatabaseEndpoint: "",
			},
			testFunc: func(t *testing.T, o *TerraformOutputs) {
				testutil.AssertEqual(t, "rds.aws.com", o.GetDatabaseEndpoint())
			},
		},
		{
			name: "GetDatabaseEndpoint - Azure",
			outputs: &TerraformOutputs{
				RDSEndpoint:      "",
				DatabaseEndpoint: "postgres.azure.com",
			},
			testFunc: func(t *testing.T, o *TerraformOutputs) {
				testutil.AssertEqual(t, "postgres.azure.com", o.GetDatabaseEndpoint())
			},
		},
		{
			name: "GetDatabasePort - AWS",
			outputs: &TerraformOutputs{
				RDSPort:      "5432",
				DatabasePort: "",
			},
			testFunc: func(t *testing.T, o *TerraformOutputs) {
				testutil.AssertEqual(t, "5432", o.GetDatabasePort())
			},
		},
		{
			name: "GetDatabasePort - Azure",
			outputs: &TerraformOutputs{
				RDSPort:      "",
				DatabasePort: "5432",
			},
			testFunc: func(t *testing.T, o *TerraformOutputs) {
				testutil.AssertEqual(t, "5432", o.GetDatabasePort())
			},
		},
		{
			name: "GetNetworkID - AWS",
			outputs: &TerraformOutputs{
				VPCID:  "vpc-123",
				VNetID: "",
			},
			testFunc: func(t *testing.T, o *TerraformOutputs) {
				testutil.AssertEqual(t, "vpc-123", o.GetNetworkID())
			},
		},
		{
			name: "GetNetworkID - Azure",
			outputs: &TerraformOutputs{
				VPCID:  "",
				VNetID: "vnet-456",
			},
			testFunc: func(t *testing.T, o *TerraformOutputs) {
				testutil.AssertEqual(t, "vnet-456", o.GetNetworkID())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, tt.outputs)
		})
	}
}

// TestGetStringValue tests the getStringValue helper function
func TestGetStringValue(t *testing.T) {
	tests := []struct {
		name   string
		output OutputValue
		want   string
	}{
		{
			name: "string value",
			output: OutputValue{
				Value: "test-value",
			},
			want: "test-value",
		},
		{
			name: "non-string value",
			output: OutputValue{
				Value: 123,
			},
			want: "",
		},
		{
			name: "nil value",
			output: OutputValue{
				Value: nil,
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStringValue(tt.output)
			testutil.AssertEqual(t, tt.want, got)
		})
	}
}

// TestGetStringSliceValue tests the getStringSliceValue helper function
func TestGetStringSliceValue(t *testing.T) {
	tests := []struct {
		name   string
		output OutputValue
		want   []string
	}{
		{
			name: "string slice",
			output: OutputValue{
				Value: []interface{}{"value1", "value2", "value3"},
			},
			want: []string{"value1", "value2", "value3"},
		},
		{
			name: "empty slice",
			output: OutputValue{
				Value: []interface{}{},
			},
			want: []string{},
		},
		{
			name: "mixed types in slice",
			output: OutputValue{
				Value: []interface{}{"value1", 123, "value2"},
			},
			want: []string{"value1", "value2"},
		},
		{
			name: "non-slice value",
			output: OutputValue{
				Value: "not-a-slice",
			},
			want: []string{},
		},
		{
			name: "nil value",
			output: OutputValue{
				Value: nil,
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStringSliceValue(tt.output)
			testutil.AssertEqual(t, len(tt.want), len(got))
			for i := range tt.want {
				testutil.AssertEqual(t, tt.want[i], got[i])
			}
		})
	}
}

// TestOutputParserWithFixture tests parsing outputs using test fixtures
func TestOutputParserWithFixture(t *testing.T) {
	// Load fixture
	fixtureData := testutil.LoadFixture(t, "terraform/output-aws-dev.json")

	mock := NewMockTerraformExecutor()
	mock.OutputFunc = func(ctx context.Context, workingDir string, outputName string) (string, error) {
		if outputName == "-json" {
			return string(fixtureData), nil
		}
		return "", errors.New("unexpected output name")
	}

	parser := NewOutputParser(mock)
	outputs, err := parser.ParseOutputs(context.Background(), "/tmp/terraform")

	testutil.AssertNoError(t, err)
	testutil.AssertEqual(t, "vpc-0123456789abcdef0", outputs.VPCID)
	testutil.AssertEqual(t, 2, len(outputs.SubnetIDs))
	testutil.AssertEqual(t, "subnet-0123456789abcdef0", outputs.SubnetIDs[0])
}

// TestOutputValueJSONMarshaling tests JSON marshaling/unmarshaling of OutputValue
func TestOutputValueJSONMarshaling(t *testing.T) {
	original := map[string]OutputValue{
		"vpc_id": {
			Sensitive: false,
			Type:      "string",
			Value:     "vpc-123",
		},
		"subnet_ids": {
			Sensitive: false,
			Type:      []interface{}{"list", "string"},
			Value:     []interface{}{"subnet-1", "subnet-2"},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	testutil.AssertNoError(t, err)

	// Unmarshal back
	var parsed map[string]OutputValue
	err = json.Unmarshal(jsonData, &parsed)
	testutil.AssertNoError(t, err)

	// Verify values
	testutil.AssertEqual(t, "vpc-123", getStringValue(parsed["vpc_id"]))
	testutil.AssertEqual(t, 2, len(getStringSliceValue(parsed["subnet_ids"])))
}
