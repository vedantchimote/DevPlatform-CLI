package types

// CallerIdentity represents cloud provider identity information
type CallerIdentity struct {
	Account string
	Arn     string
	UserId  string
}

// EnvironmentCosts represents the total costs for an environment
type EnvironmentCosts struct {
	NetworkCost  float64 // VPC/VNet cost
	DatabaseCost float64 // RDS/Azure Database cost
	K8sCost      float64 // EKS/AKS tenant cost
	TotalCost    float64
	Environment  string
	Provider     string
}

// TerraformBackend represents Terraform backend configuration
type TerraformBackend struct {
	Type   string            // "s3" or "azurerm"
	Config map[string]string // Backend-specific configuration
}
