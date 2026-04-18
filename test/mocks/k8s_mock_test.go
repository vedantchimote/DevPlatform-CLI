package mocks

import (
	"context"
	"errors"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMockKubernetesClient_ListPods(t *testing.T) {
	mock := NewMockKubernetesClient()
	ctx := context.Background()
	namespace := "test-namespace"
	opts := metav1.ListOptions{}

	// Test default behavior (empty list)
	podList, err := mock.ListPods(ctx, namespace, opts)
	if err != nil {
		t.Errorf("ListPods() unexpected error = %v", err)
	}
	if podList == nil {
		t.Fatal("ListPods() returned nil")
	}
	if len(podList.Items) != 0 {
		t.Errorf("ListPods() returned %d pods, want 0", len(podList.Items))
	}

	// Verify call was recorded
	if len(mock.ListPodsCalls) != 1 {
		t.Errorf("Expected 1 ListPods call, got %d", len(mock.ListPodsCalls))
	}

	// Verify call count helper
	if mock.GetListPodsCallCount() != 1 {
		t.Errorf("GetListPodsCallCount() = %d, want 1", mock.GetListPodsCallCount())
	}
}

func TestMockKubernetesClient_ListPodsWithCustomFunc(t *testing.T) {
	mock := NewMockKubernetesClient()
	ctx := context.Background()
	namespace := "test-namespace"
	opts := metav1.ListOptions{}

	// Configure custom behavior
	expectedPods := []corev1.Pod{
		*NewMockPod("pod1", namespace, corev1.PodRunning, true),
		*NewMockPod("pod2", namespace, corev1.PodRunning, true),
	}
	mock.ListPodsFunc = func(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.PodList, error) {
		return &corev1.PodList{Items: expectedPods}, nil
	}

	// Test custom behavior
	podList, err := mock.ListPods(ctx, namespace, opts)
	if err != nil {
		t.Errorf("ListPods() unexpected error = %v", err)
	}
	if len(podList.Items) != 2 {
		t.Errorf("ListPods() returned %d pods, want 2", len(podList.Items))
	}
}

func TestMockKubernetesClient_GetPod(t *testing.T) {
	mock := NewMockKubernetesClient()
	ctx := context.Background()
	namespace := "test-namespace"
	podName := "test-pod"
	opts := metav1.GetOptions{}

	// Test default behavior
	pod, err := mock.GetPod(ctx, namespace, podName, opts)
	if err != nil {
		t.Errorf("GetPod() unexpected error = %v", err)
	}
	if pod == nil {
		t.Fatal("GetPod() returned nil")
	}
	if pod.Name != podName {
		t.Errorf("GetPod().Name = %s, want %s", pod.Name, podName)
	}
	if pod.Namespace != namespace {
		t.Errorf("GetPod().Namespace = %s, want %s", pod.Namespace, namespace)
	}
	if pod.Status.Phase != corev1.PodRunning {
		t.Errorf("GetPod().Status.Phase = %s, want Running", pod.Status.Phase)
	}

	// Verify call was recorded
	if len(mock.GetPodCalls) != 1 {
		t.Errorf("Expected 1 GetPod call, got %d", len(mock.GetPodCalls))
	}

	// Verify call count helper
	if mock.GetGetPodCallCount() != 1 {
		t.Errorf("GetGetPodCallCount() = %d, want 1", mock.GetGetPodCallCount())
	}
}

func TestMockKubernetesClient_GetPodWithCustomFunc(t *testing.T) {
	mock := NewMockKubernetesClient()
	ctx := context.Background()
	namespace := "test-namespace"
	podName := "test-pod"
	opts := metav1.GetOptions{}

	// Configure custom behavior to return error
	expectedErr := errors.New("pod not found")
	mock.GetPodFunc = func(ctx context.Context, namespace string, name string, opts metav1.GetOptions) (*corev1.Pod, error) {
		return nil, expectedErr
	}

	// Test custom behavior
	pod, err := mock.GetPod(ctx, namespace, podName, opts)
	if err != expectedErr {
		t.Errorf("GetPod() error = %v, want %v", err, expectedErr)
	}
	if pod != nil {
		t.Error("GetPod() should return nil pod on error")
	}
}

func TestMockKubernetesClient_CreatePod(t *testing.T) {
	mock := NewMockKubernetesClient()
	ctx := context.Background()
	namespace := "test-namespace"
	pod := NewMockPod("new-pod", namespace, corev1.PodPending, false)
	opts := metav1.CreateOptions{}

	// Test default behavior
	createdPod, err := mock.CreatePod(ctx, namespace, pod, opts)
	if err != nil {
		t.Errorf("CreatePod() unexpected error = %v", err)
	}
	if createdPod == nil {
		t.Fatal("CreatePod() returned nil")
	}
	if createdPod.Name != pod.Name {
		t.Errorf("CreatePod().Name = %s, want %s", createdPod.Name, pod.Name)
	}

	// Verify call was recorded
	if len(mock.CreatePodCalls) != 1 {
		t.Errorf("Expected 1 CreatePod call, got %d", len(mock.CreatePodCalls))
	}

	// Verify call count helper
	if mock.GetCreatePodCallCount() != 1 {
		t.Errorf("GetCreatePodCallCount() = %d, want 1", mock.GetCreatePodCallCount())
	}
}

func TestMockKubernetesClient_DeletePod(t *testing.T) {
	mock := NewMockKubernetesClient()
	ctx := context.Background()
	namespace := "test-namespace"
	podName := "test-pod"
	opts := metav1.DeleteOptions{}

	// Test default behavior
	err := mock.DeletePod(ctx, namespace, podName, opts)
	if err != nil {
		t.Errorf("DeletePod() unexpected error = %v", err)
	}

	// Verify call was recorded
	if len(mock.DeletePodCalls) != 1 {
		t.Errorf("Expected 1 DeletePod call, got %d", len(mock.DeletePodCalls))
	}

	// Verify call count helper
	if mock.GetDeletePodCallCount() != 1 {
		t.Errorf("GetDeletePodCallCount() = %d, want 1", mock.GetDeletePodCallCount())
	}

	// Verify call arguments
	if mock.DeletePodCalls[0].Args[1] != namespace {
		t.Errorf("Expected namespace %s, got %v", namespace, mock.DeletePodCalls[0].Args[1])
	}
	if mock.DeletePodCalls[0].Args[2] != podName {
		t.Errorf("Expected podName %s, got %v", podName, mock.DeletePodCalls[0].Args[2])
	}
}

func TestMockKubernetesClient_DeletePodWithCustomFunc(t *testing.T) {
	mock := NewMockKubernetesClient()
	ctx := context.Background()
	namespace := "test-namespace"
	podName := "test-pod"
	opts := metav1.DeleteOptions{}
	expectedErr := errors.New("delete failed")

	// Configure custom behavior
	mock.DeletePodFunc = func(ctx context.Context, namespace string, name string, opts metav1.DeleteOptions) error {
		return expectedErr
	}

	// Test custom behavior
	err := mock.DeletePod(ctx, namespace, podName, opts)
	if err != expectedErr {
		t.Errorf("DeletePod() error = %v, want %v", err, expectedErr)
	}
}

func TestMockKubernetesClient_ListEvents(t *testing.T) {
	mock := NewMockKubernetesClient()
	ctx := context.Background()
	namespace := "test-namespace"
	opts := metav1.ListOptions{}

	// Test default behavior (empty list)
	eventList, err := mock.ListEvents(ctx, namespace, opts)
	if err != nil {
		t.Errorf("ListEvents() unexpected error = %v", err)
	}
	if eventList == nil {
		t.Fatal("ListEvents() returned nil")
	}
	if len(eventList.Items) != 0 {
		t.Errorf("ListEvents() returned %d events, want 0", len(eventList.Items))
	}

	// Verify call was recorded
	if len(mock.ListEventsCalls) != 1 {
		t.Errorf("Expected 1 ListEvents call, got %d", len(mock.ListEventsCalls))
	}

	// Verify call count helper
	if mock.GetListEventsCallCount() != 1 {
		t.Errorf("GetListEventsCallCount() = %d, want 1", mock.GetListEventsCallCount())
	}
}

func TestMockKubernetesClient_ListEventsWithCustomFunc(t *testing.T) {
	mock := NewMockKubernetesClient()
	ctx := context.Background()
	namespace := "test-namespace"
	opts := metav1.ListOptions{}

	// Configure custom behavior
	expectedEvents := []corev1.Event{
		*NewMockEvent(namespace, "Created", "Pod created", corev1.EventTypeNormal),
		*NewMockEvent(namespace, "Failed", "Pod failed to start", corev1.EventTypeWarning),
	}
	mock.ListEventsFunc = func(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.EventList, error) {
		return &corev1.EventList{Items: expectedEvents}, nil
	}

	// Test custom behavior
	eventList, err := mock.ListEvents(ctx, namespace, opts)
	if err != nil {
		t.Errorf("ListEvents() unexpected error = %v", err)
	}
	if len(eventList.Items) != 2 {
		t.Errorf("ListEvents() returned %d events, want 2", len(eventList.Items))
	}
}

func TestMockKubernetesClient_Reset(t *testing.T) {
	mock := NewMockKubernetesClient()
	ctx := context.Background()
	namespace := "test-namespace"

	// Make several calls
	_, _ = mock.ListPods(ctx, namespace, metav1.ListOptions{})
	_, _ = mock.GetPod(ctx, namespace, "test-pod", metav1.GetOptions{})
	_ = mock.DeletePod(ctx, namespace, "test-pod", metav1.DeleteOptions{})

	// Verify calls were recorded
	if len(mock.ListPodsCalls) != 1 {
		t.Errorf("Expected 1 ListPods call before reset, got %d", len(mock.ListPodsCalls))
	}
	if len(mock.GetPodCalls) != 1 {
		t.Errorf("Expected 1 GetPod call before reset, got %d", len(mock.GetPodCalls))
	}
	if len(mock.DeletePodCalls) != 1 {
		t.Errorf("Expected 1 DeletePod call before reset, got %d", len(mock.DeletePodCalls))
	}

	// Reset the mock
	mock.Reset()

	// Verify all calls were cleared
	if len(mock.ListPodsCalls) != 0 {
		t.Errorf("Expected 0 ListPods calls after reset, got %d", len(mock.ListPodsCalls))
	}
	if len(mock.GetPodCalls) != 0 {
		t.Errorf("Expected 0 GetPod calls after reset, got %d", len(mock.GetPodCalls))
	}
	if len(mock.CreatePodCalls) != 0 {
		t.Errorf("Expected 0 CreatePod calls after reset, got %d", len(mock.CreatePodCalls))
	}
	if len(mock.DeletePodCalls) != 0 {
		t.Errorf("Expected 0 DeletePod calls after reset, got %d", len(mock.DeletePodCalls))
	}
	if len(mock.ListEventsCalls) != 0 {
		t.Errorf("Expected 0 ListEvents calls after reset, got %d", len(mock.ListEventsCalls))
	}
}

func TestMockKubernetesClient_MultipleCalls(t *testing.T) {
	mock := NewMockKubernetesClient()
	ctx := context.Background()
	namespace := "test-namespace"

	// Make multiple calls to the same method
	_, _ = mock.ListPods(ctx, namespace, metav1.ListOptions{})
	_, _ = mock.ListPods(ctx, namespace, metav1.ListOptions{})
	_, _ = mock.ListPods(ctx, namespace, metav1.ListOptions{})

	// Verify all calls were recorded
	if len(mock.ListPodsCalls) != 3 {
		t.Errorf("Expected 3 ListPods calls, got %d", len(mock.ListPodsCalls))
	}

	// Verify call count helper
	if mock.GetListPodsCallCount() != 3 {
		t.Errorf("GetListPodsCallCount() = %d, want 3", mock.GetListPodsCallCount())
	}
}

func TestMockKubernetesClient_Timestamps(t *testing.T) {
	mock := NewMockKubernetesClient()
	ctx := context.Background()
	namespace := "test-namespace"

	// Make a call
	_, _ = mock.ListPods(ctx, namespace, metav1.ListOptions{})

	// Verify timestamp was recorded
	if mock.ListPodsCalls[0].Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}
}

func TestNewMockPod(t *testing.T) {
	pod := NewMockPod("test-pod", "test-namespace", corev1.PodRunning, true)

	if pod.Name != "test-pod" {
		t.Errorf("NewMockPod().Name = %s, want test-pod", pod.Name)
	}
	if pod.Namespace != "test-namespace" {
		t.Errorf("NewMockPod().Namespace = %s, want test-namespace", pod.Namespace)
	}
	if pod.Status.Phase != corev1.PodRunning {
		t.Errorf("NewMockPod().Status.Phase = %s, want Running", pod.Status.Phase)
	}

	// Check ready condition
	hasReadyCondition := false
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
			hasReadyCondition = true
			break
		}
	}
	if !hasReadyCondition {
		t.Error("NewMockPod() should have Ready condition set to True")
	}
}

func TestNewMockEvent(t *testing.T) {
	event := NewMockEvent("test-namespace", "TestReason", "Test message", corev1.EventTypeNormal)

	if event.Namespace != "test-namespace" {
		t.Errorf("NewMockEvent().Namespace = %s, want test-namespace", event.Namespace)
	}
	if event.Reason != "TestReason" {
		t.Errorf("NewMockEvent().Reason = %s, want TestReason", event.Reason)
	}
	if event.Message != "Test message" {
		t.Errorf("NewMockEvent().Message = %s, want Test message", event.Message)
	}
	if event.Type != corev1.EventTypeNormal {
		t.Errorf("NewMockEvent().Type = %s, want Normal", event.Type)
	}
}
