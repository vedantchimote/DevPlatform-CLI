package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	clierrors "github.com/devplatform/devplatform-cli/internal/errors"
)

// AuthValidator handles Azure credential validation
type AuthValidator struct {
	credential     *azidentity.DefaultAzureCredential
	subscriptionID string
	tenantID       string
}

// NewAuthValidator creates a new Azure auth validator
func NewAuthValidator(ctx context.Context, subscriptionID string, tenantID string) (*AuthValidator, error) {
	// Use DefaultAzureCredential which supports:
	// 1. Azure CLI authentication
	// 2. Service Principal (environment variables)
	// 3. Managed Identity
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, clierrors.NewAuthError(
			clierrors.ErrCodeAuthMissingCredentials,
			"Failed to create Azure credential",
			err,
		).WithDetails("Check Azure CLI login or service principal environment variables")
	}
	
	return &AuthValidator{
		credential:     cred,
		subscriptionID: subscriptionID,
		tenantID:       tenantID,
	}, nil
}

// ValidateCredentials validates Azure credentials by listing subscriptions
func (a *AuthValidator) ValidateCredentials(ctx context.Context) error {
	// Create subscriptions client to validate credentials
	client, err := armsubscriptions.NewClient(a.credential, nil)
	if err != nil {
		return clierrors.NewAuthError(
			clierrors.ErrCodeAuthMissingCredentials,
			"Failed to create Azure subscriptions client",
			err,
		)
	}
	
	// Try to get the subscription to validate credentials
	if a.subscriptionID != "" {
		_, err = client.Get(ctx, a.subscriptionID, nil)
		if err != nil {
			return clierrors.NewAuthError(
				clierrors.ErrCodeAuthInvalidCredentials,
				"Azure credentials are invalid or expired",
				err,
			).WithDetails(fmt.Sprintf("Subscription ID: %s", a.subscriptionID))
		}
	} else {
		// If no subscription ID provided, just list subscriptions to validate
		pager := client.NewListPager(nil)
		if !pager.More() {
			return clierrors.NewAuthError(
				clierrors.ErrCodeAuthInvalidCredentials,
				"Azure credentials are invalid: no subscriptions found",
				nil,
			)
		}
		_, err = pager.NextPage(ctx)
		if err != nil {
			return clierrors.NewAuthError(
				clierrors.ErrCodeAuthInvalidCredentials,
				"Azure credentials are invalid or expired",
				err,
			)
		}
	}
	
	return nil
}

// GetCallerIdentity returns information about the Azure credentials
func (a *AuthValidator) GetCallerIdentity(ctx context.Context) (*CallerIdentity, error) {
	// Create subscriptions client
	client, err := armsubscriptions.NewClient(a.credential, nil)
	if err != nil {
		return nil, clierrors.NewAuthError(
			clierrors.ErrCodeAuthMissingCredentials,
			"Failed to create Azure subscriptions client",
			err,
		)
	}
	
	// Get subscription details
	sub, err := client.Get(ctx, a.subscriptionID, nil)
	if err != nil {
		return nil, clierrors.NewAuthError(
			clierrors.ErrCodeAuthInvalidCredentials,
			"Failed to get Azure subscription details",
			err,
		).WithDetails(fmt.Sprintf("Subscription ID: %s", a.subscriptionID))
	}
	
	return &CallerIdentity{
		SubscriptionID:   *sub.SubscriptionID,
		SubscriptionName: *sub.DisplayName,
		TenantID:         *sub.TenantID,
		State:            string(*sub.State),
	}, nil
}

// CallerIdentity represents Azure caller identity information
type CallerIdentity struct {
	SubscriptionID   string
	SubscriptionName string
	TenantID         string
	State            string
}

// GetCredential returns the Azure credential
func (a *AuthValidator) GetCredential() *azidentity.DefaultAzureCredential {
	return a.credential
}

// GetSubscriptionID returns the subscription ID
func (a *AuthValidator) GetSubscriptionID() string {
	return a.subscriptionID
}

// GetTenantID returns the tenant ID
func (a *AuthValidator) GetTenantID() string {
	return a.tenantID
}
