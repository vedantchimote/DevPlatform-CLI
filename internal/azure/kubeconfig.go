package azure

import (
	"fmt"
	"os/exec"
)

// KubeconfigManager handles AKS kubeconfig operations
type KubeconfigManager struct {
	subscriptionID string
	resourceGroup  string
}

// NewKubeconfigManager creates a new kubeconfig manager
func NewKubeconfigManager(subscriptionID string, resourceGroup string) *KubeconfigManager {
	return &KubeconfigManager{
		subscriptionID: subscriptionID,
		resourceGroup:  resourceGroup,
	}
}

// UpdateKubeconfig updates the kubeconfig file for AKS cluster access
func (k *KubeconfigManager) UpdateKubeconfig(clusterName string) error {
	args := []string{
		"aks",
		"get-credentials",
		"--name", clusterName,
		"--resource-group", k.resourceGroup,
		"--subscription", k.subscriptionID,
		"--overwrite-existing",
	}
	
	cmd := exec.Command("az", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update kubeconfig: %w\nOutput: %s", err, string(output))
	}
	
	return nil
}

// GetKubectlContextCommand returns the kubectl command to switch context
func (k *KubeconfigManager) GetKubectlContextCommand(clusterName string) string {
	contextName := clusterName
	return fmt.Sprintf("kubectl config use-context %s", contextName)
}

// GetKubectlNamespaceCommand returns the kubectl command to set namespace
func (k *KubeconfigManager) GetKubectlNamespaceCommand(namespace string) string {
	return fmt.Sprintf("kubectl config set-context --current --namespace=%s", namespace)
}

// GetConnectionCommands returns all commands needed to connect to the cluster
func (k *KubeconfigManager) GetConnectionCommands(clusterName string, namespace string) []string {
	commands := []string{
		k.GetKubectlContextCommand(clusterName),
		k.GetKubectlNamespaceCommand(namespace),
	}
	return commands
}
