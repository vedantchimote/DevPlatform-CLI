package helm

import (
	"context"
	"fmt"
	"time"

	"github.com/devplatform/devplatform-cli/internal/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// PodVerifier handles pod readiness verification
type PodVerifier struct {
	clientset *kubernetes.Clientset
	logger    logger.Logger
}

// PodStatus represents the status of pods in a namespace
type PodStatus struct {
	TotalPods    int
	ReadyPods    int
	PendingPods  int
	FailedPods   int
	Pods         []PodInfo
	Events       []string
}

// PodInfo contains information about a single pod
type PodInfo struct {
	Name      string
	Namespace string
	Status    string
	Ready     bool
	Restarts  int32
	Age       time.Duration
}

// NewPodVerifier creates a new pod verifier using the default kubeconfig
func NewPodVerifier(log logger.Logger) (*PodVerifier, error) {
	// Use the default kubeconfig path
	kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}
	
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}
	
	return &PodVerifier{
		clientset: clientset,
		logger:    log,
	}, nil
}

// VerifyPods waits for pods in a namespace to reach Running state
func (pv *PodVerifier) VerifyPods(ctx context.Context, namespace string, timeout time.Duration) (*PodStatus, error) {
	pv.logger.Info(fmt.Sprintf("Verifying pods in namespace: %s", namespace))
	
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			status, err := pv.GetPodStatus(ctx, namespace)
			if err != nil {
				return nil, err
			}
			
			// Check if all pods are ready
			if status.TotalPods > 0 && status.ReadyPods == status.TotalPods {
				pv.logger.Success(fmt.Sprintf("All %d pods are ready", status.TotalPods))
				return status, nil
			}
			
			// Check if any pods failed
			if status.FailedPods > 0 {
				events, err := pv.GetEvents(ctx, namespace)
				if err != nil {
					pv.logger.Warn(fmt.Sprintf("Failed to get events: %v", err))
				} else {
					status.Events = events
				}
				return status, fmt.Errorf("%d pods failed", status.FailedPods)
			}
			
			// Check timeout
			if time.Now().After(deadline) {
				events, err := pv.GetEvents(ctx, namespace)
				if err != nil {
					pv.logger.Warn(fmt.Sprintf("Failed to get events: %v", err))
				} else {
					status.Events = events
				}
				return status, fmt.Errorf("timeout waiting for pods to be ready: %d/%d ready", status.ReadyPods, status.TotalPods)
			}
			
			pv.logger.Info(fmt.Sprintf("Waiting for pods: %d/%d ready", status.ReadyPods, status.TotalPods))
		}
	}
}

// GetPodStatus gets the current status of pods in a namespace
func (pv *PodVerifier) GetPodStatus(ctx context.Context, namespace string) (*PodStatus, error) {
	pods, err := pv.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}
	
	status := &PodStatus{
		TotalPods: len(pods.Items),
		Pods:      make([]PodInfo, 0, len(pods.Items)),
	}
	
	for _, pod := range pods.Items {
		podInfo := PodInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    string(pod.Status.Phase),
			Age:       time.Since(pod.CreationTimestamp.Time),
		}
		
		// Count restarts
		for _, containerStatus := range pod.Status.ContainerStatuses {
			podInfo.Restarts += containerStatus.RestartCount
		}
		
		// Check if pod is ready
		podInfo.Ready = isPodReady(&pod)
		
		// Update counters
		switch pod.Status.Phase {
		case corev1.PodRunning:
			if podInfo.Ready {
				status.ReadyPods++
			}
		case corev1.PodPending:
			status.PendingPods++
		case corev1.PodFailed:
			status.FailedPods++
		}
		
		status.Pods = append(status.Pods, podInfo)
	}
	
	return status, nil
}

// GetEvents retrieves recent events from a namespace
func (pv *PodVerifier) GetEvents(ctx context.Context, namespace string) ([]string, error) {
	events, err := pv.clientset.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{
		Limit: 10,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}
	
	var eventMessages []string
	for _, event := range events.Items {
		if event.Type == corev1.EventTypeWarning || event.Type == corev1.EventTypeNormal {
			msg := fmt.Sprintf("[%s] %s: %s", event.Type, event.Reason, event.Message)
			eventMessages = append(eventMessages, msg)
		}
	}
	
	return eventMessages, nil
}

// isPodReady checks if a pod is ready
func isPodReady(pod *corev1.Pod) bool {
	if pod.Status.Phase != corev1.PodRunning {
		return false
	}
	
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	
	return false
}
