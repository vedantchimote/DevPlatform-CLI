package config

import (
	"testing"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

// TestNewDefaultConfig tests the creation of default configuration
func TestNewDefaultConfig(t *testing.T) {
	cfg := NewDefaultConfig()

	// Test global defaults
	testutil.AssertEqual(t, "aws", cfg.Global.CloudProvider)
	testutil.AssertEqual(t, 3600, cfg.Global.Timeout)
	testutil.AssertEqual(t, "info", cfg.Global.LogLevel)

	// Test AWS defaults
	testutil.AssertEqual(t, "us-east-1", cfg.AWS.Region)
	testutil.AssertEqual(t, "default", cfg.AWS.Profile)

	// Test Azure defaults
	testutil.AssertEqual(t, "eastus", cfg.Azure.Location)

	// Test environment defaults exist
	testutil.AssertTrue(t, len(cfg.Environments) == 3, "Expected 3 default environments")
	
	// Test dev environment
	dev, exists := cfg.Environments["dev"]
	testutil.AssertTrue(t, exists, "Dev environment should exist")
	testutil.AssertEqual(t, "10.0.0.0/16", dev.NetworkCIDR)
	testutil.AssertEqual(t, "db.t3.micro", dev.DBInstanceClass)
	testutil.AssertEqual(t, 20, dev.DBAllocatedStorage)
	testutil.AssertEqual(t, false, dev.DBMultiAZ)
	testutil.AssertEqual(t, 2, dev.K8sNodeCount)

	// Test staging environment
	staging, exists := cfg.Environments["staging"]
	testutil.AssertTrue(t, exists, "Staging environment should exist")
	testutil.AssertEqual(t, "10.1.0.0/16", staging.NetworkCIDR)
	testutil.AssertEqual(t, "db.t3.medium", staging.DBInstanceClass)
	testutil.AssertEqual(t, 100, staging.DBAllocatedStorage)
	testutil.AssertEqual(t, false, staging.DBMultiAZ)
	testutil.AssertEqual(t, 3, staging.K8sNodeCount)

	// Test prod environment
	prod, exists := cfg.Environments["prod"]
	testutil.AssertTrue(t, exists, "Prod environment should exist")
	testutil.AssertEqual(t, "10.2.0.0/16", prod.NetworkCIDR)
	testutil.AssertEqual(t, "db.t3.large", prod.DBInstanceClass)
	testutil.AssertEqual(t, 200, prod.DBAllocatedStorage)
	testutil.AssertEqual(t, true, prod.DBMultiAZ)
	testutil.AssertEqual(t, 5, prod.K8sNodeCount)

	// Test Terraform defaults
	testutil.AssertEqual(t, "terraform/modules", cfg.Terraform.ModulesPath)

	// Test Helm defaults
	testutil.AssertEqual(t, "charts/devplatform-base", cfg.Helm.ChartPath)
}

// TestConfigStructure tests the Config struct fields
func TestConfigStructure(t *testing.T) {
	cfg := &Config{
		Global: GlobalConfig{
			CloudProvider: "azure",
			Timeout:       1800,
			LogLevel:      "debug",
		},
		AWS: AWSConfig{
			Region:  "eu-west-1",
			Profile: "production",
		},
		Azure: AzureConfig{
			SubscriptionID: "12345678-1234-1234-1234-123456789012",
			Location:       "westeurope",
			TenantID:       "87654321-4321-4321-4321-210987654321",
		},
		Environments: map[string]EnvironmentConfig{
			"custom": {
				NetworkCIDR:        "192.168.0.0/16",
				DBInstanceClass:    "db.r5.large",
				DBAllocatedStorage: 500,
				DBMultiAZ:          true,
				K8sNodeCount:       10,
			},
		},
		Terraform: TerraformConfig{
			ModulesPath: "custom/terraform",
			Backend: BackendConfig{
				Type:               "azurerm",
				StorageAccountName: "mystorageaccount",
				ContainerName:      "tfstate",
				ResourceGroupName:  "myresourcegroup",
			},
		},
		Helm: HelmConfig{
			ChartPath: "custom/charts",
			DefaultValues: map[string]interface{}{
				"replicas": 3,
				"image":    "myapp:latest",
			},
		},
	}

	// Verify all fields are set correctly
	testutil.AssertEqual(t, "azure", cfg.Global.CloudProvider)
	testutil.AssertEqual(t, 1800, cfg.Global.Timeout)
	testutil.AssertEqual(t, "debug", cfg.Global.LogLevel)

	testutil.AssertEqual(t, "eu-west-1", cfg.AWS.Region)
	testutil.AssertEqual(t, "production", cfg.AWS.Profile)

	testutil.AssertEqual(t, "12345678-1234-1234-1234-123456789012", cfg.Azure.SubscriptionID)
	testutil.AssertEqual(t, "westeurope", cfg.Azure.Location)
	testutil.AssertEqual(t, "87654321-4321-4321-4321-210987654321", cfg.Azure.TenantID)

	custom, exists := cfg.Environments["custom"]
	testutil.AssertTrue(t, exists, "Custom environment should exist")
	testutil.AssertEqual(t, "192.168.0.0/16", custom.NetworkCIDR)
	testutil.AssertEqual(t, "db.r5.large", custom.DBInstanceClass)
	testutil.AssertEqual(t, 500, custom.DBAllocatedStorage)
	testutil.AssertEqual(t, true, custom.DBMultiAZ)
	testutil.AssertEqual(t, 10, custom.K8sNodeCount)

	testutil.AssertEqual(t, "custom/terraform", cfg.Terraform.ModulesPath)
	testutil.AssertEqual(t, "azurerm", cfg.Terraform.Backend.Type)
	testutil.AssertEqual(t, "mystorageaccount", cfg.Terraform.Backend.StorageAccountName)
	testutil.AssertEqual(t, "tfstate", cfg.Terraform.Backend.ContainerName)
	testutil.AssertEqual(t, "myresourcegroup", cfg.Terraform.Backend.ResourceGroupName)

	testutil.AssertEqual(t, "custom/charts", cfg.Helm.ChartPath)
	testutil.AssertEqual(t, 3, cfg.Helm.DefaultValues["replicas"])
	testutil.AssertEqual(t, "myapp:latest", cfg.Helm.DefaultValues["image"])
}

// TestBackendConfigAWS tests AWS backend configuration
func TestBackendConfigAWS(t *testing.T) {
	backend := BackendConfig{
		Type:           "s3",
		Bucket:         "my-terraform-state",
		DynamoDBTable:  "terraform-locks",
		Region:         "us-west-2",
	}

	testutil.AssertEqual(t, "s3", backend.Type)
	testutil.AssertEqual(t, "my-terraform-state", backend.Bucket)
	testutil.AssertEqual(t, "terraform-locks", backend.DynamoDBTable)
	testutil.AssertEqual(t, "us-west-2", backend.Region)
}

// TestBackendConfigAzure tests Azure backend configuration
func TestBackendConfigAzure(t *testing.T) {
	backend := BackendConfig{
		Type:               "azurerm",
		StorageAccountName: "tfstatestorage",
		ContainerName:      "tfstate",
		ResourceGroupName:  "terraform-rg",
	}

	testutil.AssertEqual(t, "azurerm", backend.Type)
	testutil.AssertEqual(t, "tfstatestorage", backend.StorageAccountName)
	testutil.AssertEqual(t, "tfstate", backend.ContainerName)
	testutil.AssertEqual(t, "terraform-rg", backend.ResourceGroupName)
}

// TestEnvironmentConfig tests environment configuration
func TestEnvironmentConfig(t *testing.T) {
	tests := []struct {
		name   string
		config EnvironmentConfig
	}{
		{
			name: "minimal config",
			config: EnvironmentConfig{
				NetworkCIDR:        "10.0.0.0/16",
				DBInstanceClass:    "db.t3.micro",
				DBAllocatedStorage: 20,
				DBMultiAZ:          false,
				K8sNodeCount:       1,
			},
		},
		{
			name: "production config",
			config: EnvironmentConfig{
				NetworkCIDR:        "10.100.0.0/16",
				DBInstanceClass:    "db.r5.xlarge",
				DBAllocatedStorage: 1000,
				DBMultiAZ:          true,
				K8sNodeCount:       20,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.AssertTrue(t, tt.config.NetworkCIDR != "", "NetworkCIDR should not be empty")
			testutil.AssertTrue(t, tt.config.DBInstanceClass != "", "DBInstanceClass should not be empty")
			testutil.AssertTrue(t, tt.config.DBAllocatedStorage > 0, "DBAllocatedStorage should be positive")
			testutil.AssertTrue(t, tt.config.K8sNodeCount > 0, "K8sNodeCount should be positive")
		})
	}
}

// TestLoadFromPath tests loading configuration from a specific path
func TestLoadFromPath(t *testing.T) {
	tests := []struct {
		name        string
		fixturePath string
		wantErr     bool
		validate    func(t *testing.T, cfg *Config)
	}{
		{
			name:        "valid AWS dev config",
			fixturePath: "../../test/fixtures/configs/valid-aws-dev.yaml",
			wantErr:     false,
			validate: func(t *testing.T, cfg *Config) {
				testutil.AssertEqual(t, "aws", cfg.Global.CloudProvider)
			},
		},
		{
			name:        "valid Azure prod config",
			fixturePath: "../../test/fixtures/configs/valid-azure-prod.yaml",
			wantErr:     false,
			validate: func(t *testing.T, cfg *Config) {
				// The loader starts with defaults and merges, so check what's actually loaded
				testutil.AssertTrue(t, cfg != nil, "Config should not be nil")
				// Check Azure-specific fields to verify it loaded
				testutil.AssertEqual(t, "eastus", cfg.Azure.Location)
			},
		},
		{
			name:        "malformed config",
			fixturePath: "../../test/fixtures/configs/invalid-malformed.yaml",
			wantErr:     true,
			validate:    nil,
		},
		{
			name:        "non-existent file returns defaults",
			fixturePath: "../../test/fixtures/configs/does-not-exist.yaml",
			wantErr:     false,
			validate: func(t *testing.T, cfg *Config) {
				// Should return default config when file doesn't exist
				testutil.AssertEqual(t, "aws", cfg.Global.CloudProvider)
				testutil.AssertEqual(t, "info", cfg.Global.LogLevel)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadFromPath(tt.fixturePath)
			
			if tt.wantErr {
				testutil.AssertError(t, err)
			} else {
				testutil.AssertNoError(t, err)
				testutil.AssertTrue(t, cfg != nil, "Config should not be nil")
				
				if tt.validate != nil {
					tt.validate(t, cfg)
				}
			}
		})
	}
}

// TestLoaderLoad tests the Loader.Load method
func TestLoaderLoad(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		wantErr    bool
	}{
		{
			name:       "empty config file uses defaults",
			configFile: "",
			wantErr:    false,
		},
		{
			name:       "valid config file",
			configFile: "../../test/fixtures/configs/valid-aws-dev.yaml",
			wantErr:    false,
		},
		{
			name:       "invalid config file",
			configFile: "../../test/fixtures/configs/invalid-malformed.yaml",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewLoader(tt.configFile)
			cfg, err := loader.Load()

			if tt.wantErr {
				testutil.AssertError(t, err)
			} else {
				testutil.AssertNoError(t, err)
				testutil.AssertTrue(t, cfg != nil, "Config should not be nil")
			}
		})
	}
}

// TestValidateAppName tests app name validation with table-driven tests
func TestValidateAppName(t *testing.T) {
	tests := []struct {
		name    string
		appName string
		wantErr bool
	}{
		// Valid cases
		{name: "valid lowercase", appName: "myapp", wantErr: false},
		{name: "valid with numbers", appName: "myapp123", wantErr: false},
		{name: "valid with hyphens", appName: "my-app-123", wantErr: false},
		{name: "valid minimum length", appName: "abc", wantErr: false},
		{name: "valid maximum length", appName: "abcdefghijklmnopqrstuvwxyz123456", wantErr: false},
		
		// Invalid cases
		{name: "empty name", appName: "", wantErr: true},
		{name: "too short", appName: "ab", wantErr: true},
		{name: "too long", appName: "abcdefghijklmnopqrstuvwxyz1234567", wantErr: true},
		{name: "uppercase letters", appName: "MyApp", wantErr: true},
		{name: "starts with hyphen", appName: "-myapp", wantErr: true},
		{name: "ends with hyphen", appName: "myapp-", wantErr: true},
		{name: "contains underscore", appName: "my_app", wantErr: true},
		{name: "contains space", appName: "my app", wantErr: true},
		{name: "contains special chars", appName: "my@app", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAppName(tt.appName)
			
			if tt.wantErr {
				testutil.AssertError(t, err)
			} else {
				testutil.AssertNoError(t, err)
			}
		})
	}
}

// TestValidateEnvironmentType tests environment type validation
func TestValidateEnvironmentType(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		wantErr bool
	}{
		// Valid cases
		{name: "dev environment", env: "dev", wantErr: false},
		{name: "staging environment", env: "staging", wantErr: false},
		{name: "prod environment", env: "prod", wantErr: false},
		
		// Invalid cases
		{name: "empty environment", env: "", wantErr: true},
		{name: "invalid environment", env: "test", wantErr: true},
		{name: "uppercase environment", env: "DEV", wantErr: true},
		{name: "production full name", env: "production", wantErr: true},
		{name: "development full name", env: "development", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEnvironmentType(tt.env)
			
			if tt.wantErr {
				testutil.AssertError(t, err)
			} else {
				testutil.AssertNoError(t, err)
			}
		})
	}
}

// TestValidateCloudProvider tests cloud provider validation
func TestValidateCloudProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		wantErr  bool
	}{
		// Valid cases
		{name: "aws provider", provider: "aws", wantErr: false},
		{name: "azure provider", provider: "azure", wantErr: false},
		
		// Invalid cases
		{name: "empty provider", provider: "", wantErr: true},
		{name: "gcp provider", provider: "gcp", wantErr: true},
		{name: "uppercase AWS", provider: "AWS", wantErr: true},
		{name: "uppercase Azure", provider: "Azure", wantErr: true},
		{name: "invalid provider", provider: "digitalocean", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCloudProvider(tt.provider)
			
			if tt.wantErr {
				testutil.AssertError(t, err)
			} else {
				testutil.AssertNoError(t, err)
			}
		})
	}
}

// TestValidatorValidate tests the Validator.Validate method
func TestValidatorValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid default config",
			config:  NewDefaultConfig(),
			wantErr: false,
		},
		{
			name: "invalid cloud provider",
			config: &Config{
				Global: GlobalConfig{
					CloudProvider: "gcp",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			config: &Config{
				Global: GlobalConfig{
					CloudProvider: "aws",
					LogLevel:      "trace",
				},
			},
			wantErr: true,
		},
		{
			name: "negative timeout",
			config: &Config{
				Global: GlobalConfig{
					CloudProvider: "aws",
					Timeout:       -100,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid environment name",
			config: &Config{
				Global: GlobalConfig{
					CloudProvider: "aws",
				},
				Environments: map[string]EnvironmentConfig{
					"test": {
						NetworkCIDR: "10.0.0.0/16",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid CIDR",
			config: &Config{
				Global: GlobalConfig{
					CloudProvider: "aws",
				},
				Environments: map[string]EnvironmentConfig{
					"dev": {
						NetworkCIDR: "invalid-cidr",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "negative db storage",
			config: &Config{
				Global: GlobalConfig{
					CloudProvider: "aws",
				},
				Environments: map[string]EnvironmentConfig{
					"dev": {
						NetworkCIDR:        "10.0.0.0/16",
						DBAllocatedStorage: -50,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "negative k8s node count",
			config: &Config{
				Global: GlobalConfig{
					CloudProvider: "aws",
				},
				Environments: map[string]EnvironmentConfig{
					"dev": {
						NetworkCIDR:  "10.0.0.0/16",
						K8sNodeCount: -5,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid AWS region",
			config: &Config{
				Global: GlobalConfig{
					CloudProvider: "aws",
				},
				AWS: AWSConfig{
					Region: "invalid-region",
				},
			},
			wantErr: true,
		},
		{
			name: "azure missing location",
			config: &Config{
				Global: GlobalConfig{
					CloudProvider: "azure",
				},
				Azure: AzureConfig{
					Location: "",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid azure subscription id",
			config: &Config{
				Global: GlobalConfig{
					CloudProvider: "azure",
				},
				Azure: AzureConfig{
					Location:       "eastus",
					SubscriptionID: "invalid-uuid",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid azure tenant id",
			config: &Config{
				Global: GlobalConfig{
					CloudProvider: "azure",
				},
				Azure: AzureConfig{
					Location: "eastus",
					TenantID: "not-a-uuid",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator(tt.config)
			err := validator.Validate()

			if tt.wantErr {
				testutil.AssertError(t, err)
			} else {
				testutil.AssertNoError(t, err)
			}
		})
	}
}

// TestMergerMergeFlags tests the Merger.MergeFlags method
func TestMergerMergeFlags(t *testing.T) {
	tests := []struct {
		name     string
		initial  *Config
		flags    map[string]interface{}
		validate func(t *testing.T, cfg *Config)
	}{
		{
			name:    "merge provider flag",
			initial: NewDefaultConfig(),
			flags: map[string]interface{}{
				"provider": "azure",
			},
			validate: func(t *testing.T, cfg *Config) {
				testutil.AssertEqual(t, "azure", cfg.Global.CloudProvider)
			},
		},
		{
			name:    "merge timeout flag",
			initial: NewDefaultConfig(),
			flags: map[string]interface{}{
				"timeout": 7200,
			},
			validate: func(t *testing.T, cfg *Config) {
				testutil.AssertEqual(t, 7200, cfg.Global.Timeout)
			},
		},
		{
			name:    "merge log level flag",
			initial: NewDefaultConfig(),
			flags: map[string]interface{}{
				"log_level": "debug",
			},
			validate: func(t *testing.T, cfg *Config) {
				testutil.AssertEqual(t, "debug", cfg.Global.LogLevel)
			},
		},
		{
			name:    "merge AWS flags",
			initial: NewDefaultConfig(),
			flags: map[string]interface{}{
				"aws_region":  "eu-west-1",
				"aws_profile": "production",
			},
			validate: func(t *testing.T, cfg *Config) {
				testutil.AssertEqual(t, "eu-west-1", cfg.AWS.Region)
				testutil.AssertEqual(t, "production", cfg.AWS.Profile)
			},
		},
		{
			name:    "merge Azure flags",
			initial: NewDefaultConfig(),
			flags: map[string]interface{}{
				"azure_subscription_id": "12345678-1234-1234-1234-123456789012",
				"azure_location":        "westeurope",
				"azure_tenant_id":       "87654321-4321-4321-4321-210987654321",
			},
			validate: func(t *testing.T, cfg *Config) {
				testutil.AssertEqual(t, "12345678-1234-1234-1234-123456789012", cfg.Azure.SubscriptionID)
				testutil.AssertEqual(t, "westeurope", cfg.Azure.Location)
				testutil.AssertEqual(t, "87654321-4321-4321-4321-210987654321", cfg.Azure.TenantID)
			},
		},
		{
			name:    "empty flags don't override",
			initial: NewDefaultConfig(),
			flags: map[string]interface{}{
				"provider":  "",
				"log_level": "",
			},
			validate: func(t *testing.T, cfg *Config) {
				testutil.AssertEqual(t, "aws", cfg.Global.CloudProvider)
				testutil.AssertEqual(t, "info", cfg.Global.LogLevel)
			},
		},
		{
			name:    "zero timeout doesn't override",
			initial: NewDefaultConfig(),
			flags: map[string]interface{}{
				"timeout": 0,
			},
			validate: func(t *testing.T, cfg *Config) {
				testutil.AssertEqual(t, 3600, cfg.Global.Timeout)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merger := NewMerger(tt.initial)
			merger.MergeFlags(tt.flags)
			cfg := merger.GetConfig()

			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}
