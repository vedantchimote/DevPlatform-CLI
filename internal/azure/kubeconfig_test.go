package azure

import (
	"testing"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

func TestNewKubeconfigManager(t *testing.T) {
	tests := []struct {
		name           string
		subscriptionID string
		resourceGroup  string
	}{
		{
			name:           "valid_manager",
			subscriptionID: "12345678-1234-1234-1234-123456789012",
			resourceGroup:  "devplatform-rg",
		},
		{
			name:           "different_resource_group",
			subscriptionID: "11111111-1111-1111-1111-111111111111",
			resourceGroup:  "test-rg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewKubeconfigManager(tt.subscriptionID, tt.resourceGroup)

			testutil.AssertNotEqual(t, nil, manager)
			testutil.AssertEqual(t, tt.subscriptionID, manager.subscriptionID)
			testutil.AssertEqual(t, tt.resourceGroup, manager.resourceGroup)
		})
	}
}

func TestKubeconfigManager_GetKubectlContextCommand(t *testing.T) {
	tests := []struct {
		name        string
		clusterName string
		expected    string
	}{
		{
			name:        "dev_cluster",
			clusterName: "dev-cluster",
			expected:    "kubectl config use-context dev-cluster",
		},
		{
			name:        "prod_cluster",
			clusterName: "prod-cluster",
			expected:    "kubectl config use-context prod-cluster",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewKubeconfigManager("12345678-1234-1234-1234-123456789012", "devplatform-rg")
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
			manager := NewKubeconfigManager("12345678-1234-1234-1234-123456789012", "devplatform-rg")
			result := manager.GetKubectlNamespaceCommand(tt.namespace)

			testutil.AssertEqual(t, tt.expected, result)
		})
	}
}

func TestKubeconfigManager_GetConnectionCommands(t *testing.T) {
	tests := []struct {
		name        string
		clusterName string
		namespace   string
	}{
		{
			name:        "dev_environment",
			clusterName: "dev-cluster",
			namespace:   "dev-app",
		},
		{
			name:        "prod_environment",
			clusterName: "prod-cluster",
			namespace:   "prod-app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewKubeconfigManager("12345678-1234-1234-1234-123456789012", "devplatform-rg")
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

func TestKubeconfigManager_DifferentSubscriptions(t *testing.T) {
	subscriptions := []string{
		"12345678-1234-1234-1234-123456789012",
		"11111111-1111-1111-1111-111111111111",
		"22222222-2222-2222-2222-222222222222",
	}

	for _, subID := range subscriptions {
		t.Run("subscription_"+subID, func(t *testing.T) {
			manager := NewKubeconfigManager(subID, "test-rg")

			testutil.AssertNotEqual(t, nil, manager)
			testutil.AssertEqual(t, subID, manager.subscriptionID)
			testutil.AssertEqual(t, "test-rg", manager.resourceGroup)
		})
	}
}
