package mocks

import (
	"context"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MockKubernetesClient is a mock implementation of Kubernetes client operations for testing
type MockKubernetesClient struct {
	// Function fields for configuring mock behavior
	ListPodsFunc   func(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.PodList, error)
	GetPodFunc     func(ctx context.Context, namespace string, name string, opts metav1.GetOptions) (*corev1.Pod, error)
	CreatePodFunc  func(ctx context.Context, namespace string, pod *corev1.Pod, opts metav1.CreateOptions) (*corev1.Pod, error)
	DeletePodFunc  func(ctx context.Context, namespace string, name string, opts metav1.DeleteOptions) error
	ListEventsFunc func(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.EventList, error)

	// Call tracking
	ListPodsCalls   []MockCall
	GetPodCalls     []MockCall
	CreatePodCalls  []MockCall
	DeletePodCalls  []MockCall
	ListEventsCalls []MockCall

	// Mutex for thread-safe call tracking
	mu sync.Mutex
}

// NewMockKubernetesClient creates a new mock Kubernetes client with default behavior
func NewMockKubernetesClient() *MockKubernetesClient {
	return &MockKubernetesClient{
		ListPodsCalls:   make([]MockCall, 0),
		GetPodCalls:     make([]MockCall, 0),
		CreatePodCalls:  make([]MockCall, 0),
		DeletePodCalls:  make([]MockCall, 0),
		ListEventsCalls: make([]MockCall, 0),
	}
}

// ListPods lists pods in a namespace
func (m *MockKubernetesClient) ListPods(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.PodList, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.ListPodsCalls = append(m.ListPodsCalls, MockCall{
		Args:      []interface{}{ctx, namespace, opts},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.ListPodsFunc != nil {
		return m.ListPodsFunc(ctx, namespace, opts)
	}

	// Default behavior: return empty pod list
	return &corev1.PodList{
		Items: []corev1.Pod{},
	}, nil
}

// GetPod gets a specific pod
func (m *MockKubernetesClient) GetPod(ctx context.Context, namespace string, name string, opts metav1.GetOptions) (*corev1.Pod, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.GetPodCalls = append(m.GetPodCalls, MockCall{
		Args:      []interface{}{ctx, namespace, name, opts},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.GetPodFunc != nil {
		return m.GetPodFunc(ctx, namespace, name, opts)
	}

	// Default behavior: return a running pod
	now := metav1.Now()
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Namespace:         namespace,
			CreationTimestamp: now,
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
			Conditions: []corev1.PodCondition{
				{
					Type:   corev1.PodReady,
					Status: corev1.ConditionTrue,
				},
			},
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name:         "main",
					Ready:        true,
					RestartCount: 0,
				},
			},
		},
	}, nil
}

// CreatePod creates a new pod
func (m *MockKubernetesClient) CreatePod(ctx context.Context, namespace string, pod *corev1.Pod, opts metav1.CreateOptions) (*corev1.Pod, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.CreatePodCalls = append(m.CreatePodCalls, MockCall{
		Args:      []interface{}{ctx, namespace, pod, opts},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.CreatePodFunc != nil {
		return m.CreatePodFunc(ctx, namespace, pod, opts)
	}

	// Default behavior: return the created pod with status
	createdPod := pod.DeepCopy()
	createdPod.Status.Phase = corev1.PodPending
	return createdPod, nil
}

// DeletePod deletes a pod
func (m *MockKubernetesClient) DeletePod(ctx context.Context, namespace string, name string, opts metav1.DeleteOptions) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.DeletePodCalls = append(m.DeletePodCalls, MockCall{
		Args:      []interface{}{ctx, namespace, name, opts},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.DeletePodFunc != nil {
		return m.DeletePodFunc(ctx, namespace, name, opts)
	}

	// Default behavior: return nil (success)
	return nil
}

// ListEvents lists events in a namespace
func (m *MockKubernetesClient) ListEvents(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.EventList, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.ListEventsCalls = append(m.ListEventsCalls, MockCall{
		Args:      []interface{}{ctx, namespace, opts},
		Timestamp: time.Now(),
	})

	// Execute configured function if provided
	if m.ListEventsFunc != nil {
		return m.ListEventsFunc(ctx, namespace, opts)
	}

	// Default behavior: return empty event list
	return &corev1.EventList{
		Items: []corev1.Event{},
	}, nil
}

// Reset clears all recorded calls
func (m *MockKubernetesClient) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ListPodsCalls = make([]MockCall, 0)
	m.GetPodCalls = make([]MockCall, 0)
	m.CreatePodCalls = make([]MockCall, 0)
	m.DeletePodCalls = make([]MockCall, 0)
	m.ListEventsCalls = make([]MockCall, 0)
}

// GetListPodsCallCount returns the number of times ListPods was called
func (m *MockKubernetesClient) GetListPodsCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.ListPodsCalls)
}

// GetGetPodCallCount returns the number of times GetPod was called
func (m *MockKubernetesClient) GetGetPodCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.GetPodCalls)
}

// GetCreatePodCallCount returns the number of times CreatePod was called
func (m *MockKubernetesClient) GetCreatePodCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.CreatePodCalls)
}

// GetDeletePodCallCount returns the number of times DeletePod was called
func (m *MockKubernetesClient) GetDeletePodCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.DeletePodCalls)
}

// GetListEventsCallCount returns the number of times ListEvents was called
func (m *MockKubernetesClient) GetListEventsCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.ListEventsCalls)
}

// Helper functions to create test pods

// NewMockPod creates a mock pod with the given parameters
func NewMockPod(name, namespace string, phase corev1.PodPhase, ready bool) *corev1.Pod {
	now := metav1.Now()
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Namespace:         namespace,
			CreationTimestamp: now,
		},
		Status: corev1.PodStatus{
			Phase: phase,
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name:         "main",
					Ready:        ready,
					RestartCount: 0,
				},
			},
		},
	}

	if ready {
		pod.Status.Conditions = []corev1.PodCondition{
			{
				Type:   corev1.PodReady,
				Status: corev1.ConditionTrue,
			},
		}
	} else {
		pod.Status.Conditions = []corev1.PodCondition{
			{
				Type:   corev1.PodReady,
				Status: corev1.ConditionFalse,
			},
		}
	}

	return pod
}

// NewMockEvent creates a mock Kubernetes event
func NewMockEvent(namespace, reason, message, eventType string) *corev1.Event {
	now := metav1.Now()
	return &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:              reason + "-event",
			Namespace:         namespace,
			CreationTimestamp: now,
		},
		Reason:  reason,
		Message: message,
		Type:    eventType,
	}
}
