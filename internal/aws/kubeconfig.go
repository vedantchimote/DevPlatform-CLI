package aws

import (
	"fmt"
	"os/exec"
)

// KubeconfigManager handles EKS kubeconfig operations
type KubeconfigManager struct {
	region  string
	profile string
}

// NewKubeconfigManager creates a new kubeconfig manager
func NewKubeconfigManager(region string, profile string) *KubeconfigManager {
	return &KubeconfigManager{
		region:  region,
		profile: profile,
	}
}

// UpdateKubeconfig updates the kubeconfig file for EKS cluster access
func (k *KubeconfigManager) UpdateKubeconfig(clusterName string) error {
	args := []string{
		"eks",
		"update-kubeconfig",
		"--name", clusterName,
		"--region", k.region,
	}
	
	if k.profile != "" {
		args = append(args, "--profile", k.profile)
	}
	
	cmd := exec.Command("aws", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update kubeconfig: %w\nOutput: %s", err, string(output))
	}
	
	return nil
}

// GetKubectlContextCommand returns the kubectl command to switch context
func (k *KubeconfigManager) GetKubectlContextCommand(clusterName string) string {
	contextName := fmt.Sprintf("arn:aws:eks:%s:*:cluster/%s", k.region, clusterName)
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
