package helm_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/devplatform/devplatform-cli/internal/helm"
)

// ExampleMockHelmClient_basicUsage demonstrates basic mock usage
func ExampleMockHelmClient_basicUsage() {
	// Create a new mock
	mock := helm.NewMockHelmClient()

	// Use the mock with default behavior (all operations succeed)
	ctx := context.Background()
	opts := helm.InstallOptions{
		ReleaseName: "my-app",
		Chart:       "charts/my-app",
		Namespace:   "default",
	}

	err := mock.Install(ctx, opts)
	if err != nil {
		fmt.Printf("Install failed: %v\n", err)
		return
	}

	// Check how many times Install was called
	fmt.Printf("Install called %d times\n", mock.GetInstallCallCount())

	// Output:
	// Install called 1 times
}

// ExampleMockHelmClient_customBehavior demonstrates configuring custom behavior
func ExampleMockHelmClient_customBehavior() {
	// Create a new mock
	mock := helm.NewMockHelmClient()

	// Configure custom behavior for Status
	mock.StatusFunc = func(ctx context.Context, releaseName, namespace string) (*helm.ReleaseStatus, error) {
		return &helm.ReleaseStatus{
			Name:      releaseName,
			Namespace: namespace,
			Status:    "deployed",
			Revision:  5,
		}, nil
	}

	// Use the mock
	ctx := context.Background()
	status, err := mock.Status(ctx, "my-app", "default")
	if err != nil {
		fmt.Printf("Status failed: %v\n", err)
		return
	}

	fmt.Printf("Release: %s, Status: %s, Revision: %d\n", 
		status.Name, status.Status, status.Revision)

	// Output:
	// Release: my-app, Status: deployed, Revision: 5
}

// ExampleMockHelmClient_errorSimulation demonstrates simulating errors
func ExampleMockHelmClient_errorSimulation() {
	// Create a new mock
	mock := helm.NewMockHelmClient()

	// Configure Uninstall to return an error
	mock.UninstallFunc = func(ctx context.Context, opts helm.UninstallOptions) error {
		return errors.New("release not found")
	}

	// Use the mock
	ctx := context.Background()
	opts := helm.UninstallOptions{
		ReleaseName: "my-app",
		Namespace:   "default",
	}

	err := mock.Uninstall(ctx, opts)
	if err != nil {
		fmt.Printf("Uninstall failed: %v\n", err)
	}

	// Output:
	// Uninstall failed: release not found
}

// ExampleMockHelmClient_callTracking demonstrates call tracking
func ExampleMockHelmClient_callTracking() {
	// Create a new mock
	mock := helm.NewMockHelmClient()

	// Make multiple calls
	ctx := context.Background()
	_ = mock.Install(ctx, helm.InstallOptions{ReleaseName: "app1", Namespace: "ns1"})
	_ = mock.Install(ctx, helm.InstallOptions{ReleaseName: "app2", Namespace: "ns2"})
	_ = mock.Install(ctx, helm.InstallOptions{ReleaseName: "app3", Namespace: "ns3"})

	// Check call count
	fmt.Printf("Install called %d times\n", mock.GetInstallCallCount())

	// Inspect call arguments
	for i, call := range mock.InstallCalls {
		opts := call.Args[1].(helm.InstallOptions)
		fmt.Printf("Call %d: release=%s, namespace=%s\n", i+1, opts.ReleaseName, opts.Namespace)
	}

	// Output:
	// Install called 3 times
	// Call 1: release=app1, namespace=ns1
	// Call 2: release=app2, namespace=ns2
	// Call 3: release=app3, namespace=ns3
}

// ExampleMockHelmClient_reset demonstrates resetting the mock
func ExampleMockHelmClient_reset() {
	// Create a new mock
	mock := helm.NewMockHelmClient()

	// Make some calls
	ctx := context.Background()
	_ = mock.Install(ctx, helm.InstallOptions{ReleaseName: "test"})
	_ = mock.Upgrade(ctx, helm.UpgradeOptions{ReleaseName: "test"})

	fmt.Printf("Before reset - Install: %d, Upgrade: %d\n", 
		mock.GetInstallCallCount(), mock.GetUpgradeCallCount())

	// Reset the mock
	mock.Reset()

	fmt.Printf("After reset - Install: %d, Upgrade: %d\n", 
		mock.GetInstallCallCount(), mock.GetUpgradeCallCount())

	// Output:
	// Before reset - Install: 1, Upgrade: 1
	// After reset - Install: 0, Upgrade: 0
}

// ExampleMockHelmClient_listReleases demonstrates mocking List operation
func ExampleMockHelmClient_listReleases() {
	// Create a new mock
	mock := helm.NewMockHelmClient()

	// Configure custom behavior for List
	mock.ListFunc = func(ctx context.Context, namespace string) ([]*helm.Release, error) {
		return []*helm.Release{
			{Name: "app1", Namespace: namespace, Status: "deployed"},
			{Name: "app2", Namespace: namespace, Status: "deployed"},
			{Name: "app3", Namespace: namespace, Status: "failed"},
		}, nil
	}

	// Use the mock
	ctx := context.Background()
	releases, err := mock.List(ctx, "production")
	if err != nil {
		fmt.Printf("List failed: %v\n", err)
		return
	}

	fmt.Printf("Found %d releases\n", len(releases))
	for _, release := range releases {
		fmt.Printf("- %s: %s\n", release.Name, release.Status)
	}

	// Output:
	// Found 3 releases
	// - app1: deployed
	// - app2: deployed
	// - app3: failed
}
