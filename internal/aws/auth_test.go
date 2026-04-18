package aws

import (
	"testing"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

func TestCallerIdentity_Fields(t *testing.T) {
	identity := &CallerIdentity{
		Account: "123456789012",
		Arn:     "arn:aws:iam::123456789012:user/test-user",
		UserId:  "AIDACKCEVSQ6C2EXAMPLE",
	}

	testutil.AssertEqual(t, "123456789012", identity.Account)
	testutil.AssertEqual(t, "arn:aws:iam::123456789012:user/test-user", identity.Arn)
	testutil.AssertEqual(t, "AIDACKCEVSQ6C2EXAMPLE", identity.UserId)
}
