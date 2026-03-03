package aws

import (
	"context"
	"fmt"

	"github.com/devplatform/devplatform-cli/internal/provider/types"
)

// AWSProvider implements the CloudProvider interface for AWS
type AWSProvider struct {
	auth    *AuthValidator
	kubeconfig *KubeconfigManager
	pricing *PricingCalculator
	region  string
	profile string
}

// NewAWSProvider creates a new AWS provider
func NewAWSProvider(ctx context.Context, region string, profile string) (*AWSProvider, error) {
	auth, err := NewAuthValidator(ctx, region, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS auth validator: %w", err)
	}
	
	kubeconfig := NewKubeconfigManager(region, profile)
	pricing := NewPricingCalculator(region)
	
	return &AWSProvider{
		auth:       auth,
		kubeconfig: kubeconfig,
		pricing:    pricing,
		region:     region,
		profile:    profile,
	}, nil
}

// ValidateCredentials validates AWS credentials
func (p *AWSProvider) ValidateCredentials(ctx context.Context) error {
	return p.auth.ValidateCredentials(ctx)
}

// GetCallerIdentity returns AWS caller identity information
func (p *AWSProvider) GetCallerIdentity(ctx context.Context) (*types.CallerIdentity, error) {
	identity, err := p.auth.GetCallerIdentity(ctx)
	if err != nil {
		return nil, err
	}
	
	return &types.CallerIdentity{
		Account: identity.Account,
		Arn:     identity.Arn,
		UserId:  identity.UserId,
	}, nil
}

// UpdateKubeconfig updates the kubeconfig for EKS cluster access
func (p *AWSProvider) UpdateKubeconfig(clusterName string) error {
	return p.kubeconfig.UpdateKubeconfig(clusterName)
}

// GetConnectionCommands returns kubectl commands for connecting to the cluster
func (p *AWSProvider) GetConnectionCommands(clusterName string, namespace string) []string {
	return p.kubeconfig.GetConnectionCommands(clusterName, namespace)
}

// CalculateTotalCost calculates the total monthly cost for an environment
func (p *AWSProvider) CalculateTotalCost(envType string) (*types.EnvironmentCosts, error) {
	costs, err := p.pricing.CalculateTotalCost(envType)
	if err != nil {
		return nil, err
	}
	
	return &types.EnvironmentCosts{
		NetworkCost:  costs.VPCCost,
		DatabaseCost: costs.RDSCost,
		K8sCost:      costs.EKSCost,
		TotalCost:    costs.TotalCost,
		Environment:  costs.Environment,
		Provider:     "aws",
	}, nil
}

// GetTerraformBackend returns the Terraform backend configuration for AWS
func (p *AWSProvider) GetTerraformBackend(appName string, envType string) (*types.TerraformBackend, error) {
	// S3 backend configuration
	bucket := fmt.Sprintf("devplatform-terraform-state-%s", p.region)
	key := fmt.Sprintf("%s/%s/terraform.tfstate", appName, envType)
	dynamodbTable := "devplatform-terraform-locks"
	
	return &types.TerraformBackend{
		Type: "s3",
		Config: map[string]string{
			"bucket":         bucket,
			"key":            key,
			"region":         p.region,
			"dynamodb_table": dynamodbTable,
			"encrypt":        "true",
		},
	}, nil
}

// GetModulePath returns the path to the Terraform modules for AWS
func (p *AWSProvider) GetModulePath() string {
	return "terraform/modules/aws"
}

// GetProviderName returns the name of the cloud provider
func (p *AWSProvider) GetProviderName() string {
	return "aws"
}
