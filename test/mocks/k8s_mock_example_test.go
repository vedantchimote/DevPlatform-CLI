package mocks_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/devplatform/devplatform-cli/test/mocks"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExampleMockKubernetesClient_basicUsage demonstrates basic mock usage
func ExampleMockKubernetesClient_basicUsage() {
	// Create a new mock
	mock := mocks.NewMockKubernetesClient()

	// Use the mock with default behavior (empty pod list)
	ctx := context.Background()
	podList, err := mock.ListPods(ctx, "default", metav1.ListOptions{})
	if err != nil {
		fmt.Printf("ListPods failed: %v\n", err)
		return
	}

	// Check how many times ListPods was called
	fmt.Printf("ListPods called %d times\n", mock.GetListPodsCallCount())
	fmt.Printf("Found %d pods\n", len(podList.Items))

	// Output:
	// ListPods called 1 times
	// Found 0 pods
}

// ExampleMockKubernetesClient_customBehavior demonstrates configuring custom behavior
func ExampleMockKubernetesClient_customBehavior() {
	// Create a new mock
	mock := mocks.NewMockKubernetesClient()

	// Configure custom behavior for ListPods
	mock.ListPodsFunc = func(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.PodList, error) {
		return &corev1.PodList{
			Items: []corev1.Pod{
				*mocks.NewMockPod("app-1", namespace, corev1.PodRunning, true),
				*mocks.NewMockPod("app-2", namespace, corev1.PodRunning, true),
				*mocks.NewMockPod("app-3", namespace, corev1.PodPending, false),
			},
		}, nil
	}

	// Use the mock
	ctx := context.Background()
	podList, err := mock.ListPods(ctx, "production", metav1.ListOptions{})
	if err != nil {
		fmt.Printf("ListPods failed: %v\n", err)
		return
	}

	fmt.Printf("Found %d pods\n", len(podList.Items))
	for _, pod := range podList.Items {
		fmt.Printf("- %s: %s\n", pod.Name, pod.Status.Phase)
	}

	// Output:
	// Found 3 pods
	// - app-1: Running
	// - app-2: Running
	// - app-3: Pending
}

// ExampleMockKubernetesClient_errorSimulation demonstrates simulating errors
func ExampleMockKubernetesClient_errorSimulation() {
	// Create a new mock
	mock := mocks.NewMockKubernetesClient()

	// Configure GetPod to return an error
	mock.GetPodFunc = func(ctx context.Context, namespace string, name string, opts metav1.GetOptions) (*corev1.Pod, error) {
		return nil, errors.New("pod not found")
	}

	// Use the mock
	ctx := context.Background()
	pod, err := mock.GetPod(ctx, "default", "missing-pod", metav1.GetOptions{})
	if err != nil {
		fmt.Printf("GetPod failed: %v\n", err)
	}
	if pod == nil {
		fmt.Println("Pod is nil")
	}

	// Output:
	// GetPod failed: pod not found
	// Pod is nil
}

// ExampleMockKubernetesClient_callTracking demonstrates call tracking
func ExampleMockKubernetesClient_callTracking() {
	// Create a new mock
	mock := mocks.NewMockKubernetesClient()

	// Make multiple calls
	ctx := context.Background()
	_ = mock.DeletePod(ctx, "default", "pod-1", metav1.DeleteOptions{})
	_ = mock.DeletePod(ctx, "default", "pod-2", metav1.DeleteOptions{})
	_ = mock.DeletePod(ctx, "default", "pod-3", metav1.DeleteOptions{})

	// Check call count
	fmt.Printf("DeletePod called %d times\n", mock.GetDeletePodCallCount())

	// Inspect call arguments
	for i, call := range mock.DeletePodCalls {
		podName := call.Args[2].(string)
		fmt.Printf("Call %d: pod=%s\n", i+1, podName)
	}

	// Output:
	// DeletePod called 3 times
	// Call 1: pod=pod-1
	// Call 2: pod=pod-2
	// Call 3: pod=pod-3
}

// ExampleMockKubernetesClient_reset demonstrates resetting the mock
func ExampleMockKubernetesClient_reset() {
	// Create a new mock
	mock := mocks.NewMockKubernetesClient()

	// Make some calls
	ctx := context.Background()
	_, _ = mock.ListPods(ctx, "default", metav1.ListOptions{})
	_, _ = mock.GetPod(ctx, "default", "test-pod", metav1.GetOptions{})

	fmt.Printf("Before reset - ListPods: %d, GetPod: %d\n",
		mock.GetListPodsCallCount(), mock.GetGetPodCallCount())

	// Reset the mock
	mock.Reset()

	fmt.Printf("After reset - ListPods: %d, GetPod: %d\n",
		mock.GetListPodsCallCount(), mock.GetGetPodCallCount())

	// Output:
	// Before reset - ListPods: 1, GetPod: 1
	// After reset - ListPods: 0, GetPod: 0
}

// ExampleMockKubernetesClient_createPod demonstrates creating a pod
func ExampleMockKubernetesClient_createPod() {
	// Create a new mock
	mock := mocks.NewMockKubernetesClient()

	// Create a pod
	ctx := context.Background()
	pod := mocks.NewMockPod("new-app", "default", corev1.PodPending, false)
	createdPod, err := mock.CreatePod(ctx, "default", pod, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("CreatePod failed: %v\n", err)
		return
	}

	fmt.Printf("Created pod: %s\n", createdPod.Name)
	fmt.Printf("Status: %s\n", createdPod.Status.Phase)
	fmt.Printf("CreatePod called %d times\n", mock.GetCreatePodCallCount())

	// Output:
	// Created pod: new-app
	// Status: Pending
	// CreatePod called 1 times
}

// ExampleMockKubernetesClient_listEvents demonstrates listing events
func ExampleMockKubernetesClient_listEvents() {
	// Create a new mock
	mock := mocks.NewMockKubernetesClient()

	// Configure custom behavior for ListEvents
	mock.ListEventsFunc = func(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.EventList, error) {
		return &corev1.EventList{
			Items: []corev1.Event{
				*mocks.NewMockEvent(namespace, "Created", "Pod created successfully", corev1.EventTypeNormal),
				*mocks.NewMockEvent(namespace, "Started", "Container started", corev1.EventTypeNormal),
				*mocks.NewMockEvent(namespace, "Failed", "Container failed to start", corev1.EventTypeWarning),
			},
		}, nil
	}

	// Use the mock
	ctx := context.Background()
	eventList, err := mock.ListEvents(ctx, "production", metav1.ListOptions{})
	if err != nil {
		fmt.Printf("ListEvents failed: %v\n", err)
		return
	}

	fmt.Printf("Found %d events\n", len(eventList.Items))
	for _, event := range eventList.Items {
		fmt.Printf("- [%s] %s: %s\n", event.Type, event.Reason, event.Message)
	}

	// Output:
	// Found 3 events
	// - [Normal] Created: Pod created successfully
	// - [Normal] Started: Container started
	// - [Warning] Failed: Container failed to start
}

// ExampleNewMockPod demonstrates creating mock pods
func ExampleNewMockPod() {
	// Create different types of mock pods
	runningPod := mocks.NewMockPod("app-1", "default", corev1.PodRunning, true)
	pendingPod := mocks.NewMockPod("app-2", "default", corev1.PodPending, false)
	failedPod := mocks.NewMockPod("app-3", "default", corev1.PodFailed, false)

	fmt.Printf("Running pod: %s - %s\n", runningPod.Name, runningPod.Status.Phase)
	fmt.Printf("Pending pod: %s - %s\n", pendingPod.Name, pendingPod.Status.Phase)
	fmt.Printf("Failed pod: %s - %s\n", failedPod.Name, failedPod.Status.Phase)

	// Output:
	// Running pod: app-1 - Running
	// Pending pod: app-2 - Pending
	// Failed pod: app-3 - Failed
}

// ExampleNewMockEvent demonstrates creating mock events
func ExampleNewMockEvent() {
	// Create different types of mock events
	normalEvent := mocks.NewMockEvent("default", "Created", "Pod created", corev1.EventTypeNormal)
	warningEvent := mocks.NewMockEvent("default", "Failed", "Pod failed", corev1.EventTypeWarning)

	fmt.Printf("Normal event: [%s] %s - %s\n", normalEvent.Type, normalEvent.Reason, normalEvent.Message)
	fmt.Printf("Warning event: [%s] %s - %s\n", warningEvent.Type, warningEvent.Reason, warningEvent.Message)

	// Output:
	// Normal event: [Normal] Created - Pod created
	// Warning event: [Warning] Failed - Pod failed
}
