package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	clierrors "github.com/devplatform/devplatform-cli/internal/errors"
)

// AuthValidator handles AWS credential validation
type AuthValidator struct {
	cfg aws.Config
}

// NewAuthValidator creates a new AWS auth validator
func NewAuthValidator(ctx context.Context, region string, profile string) (*AuthValidator, error) {
	// Load AWS configuration
	var cfg aws.Config
	var err error
	
	if profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithSharedConfigProfile(profile),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
		)
	}
	
	if err != nil {
		return nil, clierrors.NewAuthError(
			clierrors.ErrCodeAuthMissingCredentials,
			"Failed to load AWS configuration",
			err,
		).WithDetails(fmt.Sprintf("Region: %s, Profile: %s", region, profile))
	}
	
	return &AuthValidator{cfg: cfg}, nil
}

// ValidateCredentials validates AWS credentials by calling GetCallerIdentity
func (a *AuthValidator) ValidateCredentials(ctx context.Context) error {
	stsClient := sts.NewFromConfig(a.cfg)
	
	_, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return clierrors.NewAuthError(
			clierrors.ErrCodeAuthInvalidCredentials,
			"AWS credentials are invalid or expired",
			err,
		).WithDetails("Check AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, or AWS profile configuration")
	}
	
	return nil
}

// GetCallerIdentity returns information about the AWS credentials
func (a *AuthValidator) GetCallerIdentity(ctx context.Context) (*CallerIdentity, error) {
	stsClient := sts.NewFromConfig(a.cfg)
	
	result, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, clierrors.NewAuthError(
			clierrors.ErrCodeAuthInvalidCredentials,
			"Failed to get AWS caller identity",
			err,
		)
	}
	
	return &CallerIdentity{
		Account: aws.ToString(result.Account),
		Arn:     aws.ToString(result.Arn),
		UserId:  aws.ToString(result.UserId),
	}, nil
}

// CallerIdentity represents AWS caller identity information
type CallerIdentity struct {
	Account string
	Arn     string
	UserId  string
}

// GetConfig returns the AWS configuration
func (a *AuthValidator) GetConfig() aws.Config {
	return a.cfg
}
