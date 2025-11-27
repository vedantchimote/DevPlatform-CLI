package config

// Config represents the complete configuration for DevPlatform CLI
type Config struct {
	Global       GlobalConfig                  `yaml:"global"`
	AWS          AWSConfig                     `yaml:"aws"`
	Azure        AzureConfig                   `yaml:"azure"`
	Environments map[string]EnvironmentConfig  `yaml:"environments"`
	Terraform    TerraformConfig               `yaml:"terraform"`
	Helm         HelmConfig                    `yaml:"helm"`
}

// GlobalConfig contains global settings
type GlobalConfig struct {
	CloudProvider string `yaml:"cloud_provider"` // aws or azure
	Timeout       int    `yaml:"timeout"`        // timeout in seconds
	LogLevel      string `yaml:"log_level"`      // debug, info, warn, error
}

// AWSConfig contains AWS-specific configuration
type AWSConfig struct {
	Region  string `yaml:"region"`
	Profile string `yaml:"aws_profile"`
}

// AzureConfig contains Azure-specific configuration
type AzureConfig struct {
	SubscriptionID string `yaml:"subscription_id"`
	Location       string `yaml:"location"`
	TenantID       string `yaml:"tenant_id"`
}

// EnvironmentConfig contains environment-specific settings
type EnvironmentConfig struct {
	// Network configuration (VPC for AWS, VNet for Azure)
	NetworkCIDR string `yaml:"network_cidr"`

	// Database configuration (RDS for AWS, Azure Database for Azure)
	DBInstanceClass    string `yaml:"db_instance_class"`
	DBAllocatedStorage int    `yaml:"db_allocated_storage"`
	DBMultiAZ          bool   `yaml:"db_multi_az"`

	// Kubernetes configuration (EKS for AWS, AKS for Azure)
	K8sNodeCount int `yaml:"k8s_node_count"`
}

// TerraformConfig contains Terraform-specific configuration
type TerraformConfig struct {
	Backend     BackendConfig `yaml:"backend"`
	ModulesPath string        `yaml:"modules_path"`
}

// BackendConfig contains Terraform backend configuration
type BackendConfig struct {
	// Common
	Type string `yaml:"type"` // s3 or azurerm

	// AWS S3 backend
	Bucket        string `yaml:"bucket"`
	DynamoDBTable string `yaml:"dynamodb_table"`
	Region        string `yaml:"region"`

	// Azure Storage backend
	StorageAccountName string `yaml:"storage_account_name"`
	ContainerName      string `yaml:"container_name"`
	ResourceGroupName  string `yaml:"resource_group_name"`
}

// HelmConfig contains Helm-specific configuration
type HelmConfig struct {
	ChartPath     string                 `yaml:"chart_path"`
	DefaultValues map[string]interface{} `yaml:"default_values"`
}

// NewDefaultConfig returns a Config with sensible defaults
func NewDefaultConfig() *Config {
	return &Config{
		Global: GlobalConfig{
			CloudProvider: "aws",
			Timeout:       3600,
			LogLevel:      "info",
		},
		AWS: AWSConfig{
			Region:  "us-east-1",
			Profile: "default",
		},
		Azure: AzureConfig{
			Location: "eastus",
		},
		Environments: map[string]EnvironmentConfig{
			"dev": {
				NetworkCIDR:        "10.0.0.0/16",
				DBInstanceClass:    "db.t3.micro",
				DBAllocatedStorage: 20,
				DBMultiAZ:          false,
				K8sNodeCount:       2,
			},
			"staging": {
				NetworkCIDR:        "10.1.0.0/16",
				DBInstanceClass:    "db.t3.medium",
				DBAllocatedStorage: 100,
				DBMultiAZ:          false,
				K8sNodeCount:       3,
			},
			"prod": {
				NetworkCIDR:        "10.2.0.0/16",
				DBInstanceClass:    "db.t3.large",
				DBAllocatedStorage: 200,
				DBMultiAZ:          true,
				K8sNodeCount:       5,
			},
		},
		Terraform: TerraformConfig{
			ModulesPath: "terraform/modules",
		},
		Helm: HelmConfig{
			ChartPath: "charts/devplatform-base",
		},
	}
}
