package aws

import (
	"testing"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

func TestNewKubeconfigManager(t *testing.T) {
	tests := []struct {
		name    string
		region  string
		profile string
	}{
		{
			name:    "with_region_only",
			region:  "us-east-1",
			profile: "",
		},
		{
			name:    "with_region_and_profile",
			region:  "us-west-2",
			profile: "test-profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewKubeconfigManager(tt.region, tt.profile)

			testutil.AssertNotEqual(t, nil, manager)
			testutil.AssertEqual(t, tt.region, manager.region)
			testutil.AssertEqual(t, tt.profile, manager.profile)
		})
	}
}

func TestKubeconfigManager_GetKubectlContextCommand(t *testing.T) {
	tests := []struct {
		name        string
		region      string
		clusterName string
		expected    string
	}{
		{
			name:        "us_east_1_cluster",
			region:      "us-east-1",
			clusterName: "my-cluster",
			expected:    "kubectl config use-context arn:aws:eks:us-east-1:*:cluster/my-cluster",
		},
		{
			name:        "us_west_2_cluster",
			region:      "us-west-2",
			clusterName: "prod-cluster",
			expected:    "kubectl config use-context arn:aws:eks:us-west-2:*:cluster/prod-cluster",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewKubeconfigManager(tt.region, "")
			result := manager.GetKubectlContextCommand(tt.clusterName)

			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

func TestKubeconfigManager_GetKubectlNamespaceCommand(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		expected  string
	}{
		{
			name:      "default_namespace",
			namespace: "default",
			expected:  "kubectl config set-context --current --namespace=default",
		},
		{
			name:      "custom_namespace",
			namespace: "my-app",
			expected:  "kubectl config set-context --current --namespace=my-app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewKubeconfigManager("us-east-1", "")
			result := manager.GetKubectlNamespaceCommand(tt.namespace)

			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

func TestKubeconfigManager_GetConnectionCommands(t *testing.T) {
	tests := []struct {
		name        string
		region      string
		clusterName string
		namespace   string
	}{
		{
			name:        "dev_environment",
			region:      "us-east-1",
			clusterName: "dev-cluster",
			namespace:   "dev-app",
		},
		{
			name:        "prod_environment",
			region:      "us-west-2",
			clusterName: "prod-cluster",
			namespace:   "prod-app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewKubeconfigManager(tt.region, "")
			commands := manager.GetConnectionCommands(tt.clusterName, tt.namespace)

			// Verify we get 2 commands
			testutil.AssertEqual(t, 2, len(commands))

			// Verify first command is context switch
			expectedContext := manager.GetKubectlContextCommand(tt.clusterName)
			testutil.AssertEqual(t, expectedContext, commands[0])

			// Verify second command is namespace set
			expectedNamespace := manager.GetKubectlNamespaceCommand(tt.namespace)
			testutil.AssertEqual(t, expectedNamespace, commands[1])
		})
	}
}

func TestKubeconfigManager_WithProfile(t *testing.T) {
	manager := NewKubeconfigManager("us-east-1", "my-profile")

	testutil.AssertNotEqual(t, nil, manager)
	testutil.AssertEqual(t, "us-east-1", manager.region)
	testutil.AssertEqual(t, "my-profile", manager.profile)
}

func TestKubeconfigManager_EmptyProfile(t *testing.T) {
	manager := NewKubeconfigManager("us-west-2", "")

	testutil.AssertNotEqual(t, nil, manager)
	testutil.AssertEqual(t, "us-west-2", manager.region)
	testutil.AssertEqual(t, "", manager.profile)
}
