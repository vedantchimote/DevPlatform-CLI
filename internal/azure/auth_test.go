package azure

import (
	"testing"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

func TestCallerIdentity_Fields(t *testing.T) {
	identity := &CallerIdentity{
		SubscriptionID:   "12345678-1234-1234-1234-123456789012",
		SubscriptionName: "Test Subscription",
		TenantID:         "87654321-4321-4321-4321-210987654321",
		State:            "Enabled",
	}

	testutil.AssertEqual(t, "12345678-1234-1234-1234-123456789012", identity.SubscriptionID)
	testutil.AssertEqual(t, "Test Subscription", identity.SubscriptionName)
	testutil.AssertEqual(t, "87654321-4321-4321-4321-210987654321", identity.TenantID)
	testutil.AssertEqual(t, "Enabled", identity.State)
}
