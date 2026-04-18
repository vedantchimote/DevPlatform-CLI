package helm

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMockHelmClient_Install(t *testing.T) {
	mock := NewMockHelmClient()
	ctx := context.Background()
	opts := InstallOptions{
		ReleaseName: "test-release",
		Chart:       "test-chart",
		Namespace:   "test-namespace",
	}

	// Test default behavior (success)
	err := mock.Install(ctx, opts)
	if err != nil {
		t.Errorf("Install() unexpected error = %v", err)
	}

	// Verify call was recorded
	if len(mock.InstallCalls) != 1 {
		t.Errorf("Expected 1 Install call, got %d", len(mock.InstallCalls))
	}

	// Verify call count helper
	if mock.GetInstallCallCount() != 1 {
		t.Errorf("GetInstallCallCount() = %d, want 1", mock.GetInstallCallCount())
	}

	// Verify call arguments
	if len(mock.InstallCalls[0].Args) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(mock.InstallCalls[0].Args))
	}
}

func TestMockHelmClient_InstallWithCustomFunc(t *testing.T) {
	mock := NewMockHelmClient()
	ctx := context.Background()
	opts := InstallOptions{
		ReleaseName: "test-release",
		Chart:       "test-chart",
		Namespace:   "test-namespace",
	}
	expectedErr := errors.New("install failed")

	// Configure custom behavior
	mock.InstallFunc = func(ctx context.Context, opts InstallOptions) error {
		return expectedErr
	}

	// Test custom behavior
	err := mock.Install(ctx, opts)
	if err != expectedErr {
		t.Errorf("Install() error = %v, want %v", err, expectedErr)
	}

	// Verify call was still recorded
	if len(mock.InstallCalls) != 1 {
		t.Errorf("Expected 1 Install call, got %d", len(mock.InstallCalls))
	}
}

func TestMockHelmClient_Upgrade(t *testing.T) {
	mock := NewMockHelmClient()
	ctx := context.Background()
	opts := UpgradeOptions{
		ReleaseName: "test-release",
		Chart:       "test-chart",
		Namespace:   "test-namespace",
	}

	// Test default behavior
	err := mock.Upgrade(ctx, opts)
	if err != nil {
		t.Errorf("Upgrade() unexpected error = %v", err)
	}

	// Verify call was recorded
	if len(mock.UpgradeCalls) != 1 {
		t.Errorf("Expected 1 Upgrade call, got %d", len(mock.UpgradeCalls))
	}

	// Verify call count helper
	if mock.GetUpgradeCallCount() != 1 {
		t.Errorf("GetUpgradeCallCount() = %d, want 1", mock.GetUpgradeCallCount())
	}
}

func TestMockHelmClient_Uninstall(t *testing.T) {
	mock := NewMockHelmClient()
	ctx := context.Background()
	opts := UninstallOptions{
		ReleaseName: "test-release",
		Namespace:   "test-namespace",
	}

	// Test default behavior
	err := mock.Uninstall(ctx, opts)
	if err != nil {
		t.Errorf("Uninstall() unexpected error = %v", err)
	}

	// Verify call was recorded
	if len(mock.UninstallCalls) != 1 {
		t.Errorf("Expected 1 Uninstall call, got %d", len(mock.UninstallCalls))
	}

	// Verify call count helper
	if mock.GetUninstallCallCount() != 1 {
		t.Errorf("GetUninstallCallCount() = %d, want 1", mock.GetUninstallCallCount())
	}
}

func TestMockHelmClient_Status(t *testing.T) {
	mock := NewMockHelmClient()
	ctx := context.Background()
	releaseName := "test-release"
	namespace := "test-namespace"

	// Test default behavior
	status, err := mock.Status(ctx, releaseName, namespace)
	if err != nil {
		t.Errorf("Status() unexpected error = %v", err)
	}
	if status == nil {
		t.Fatal("Status() returned nil")
	}
	if status.Name != releaseName {
		t.Errorf("Status().Name = %s, want %s", status.Name, releaseName)
	}
	if status.Namespace != namespace {
		t.Errorf("Status().Namespace = %s, want %s", status.Namespace, namespace)
	}
	if status.Status != "deployed" {
		t.Errorf("Status().Status = %s, want deployed", status.Status)
	}

	// Verify call was recorded
	if len(mock.StatusCalls) != 1 {
		t.Errorf("Expected 1 Status call, got %d", len(mock.StatusCalls))
	}

	// Verify call count helper
	if mock.GetStatusCallCount() != 1 {
		t.Errorf("GetStatusCallCount() = %d, want 1", mock.GetStatusCallCount())
	}
}

func TestMockHelmClient_StatusWithCustomFunc(t *testing.T) {
	mock := NewMockHelmClient()
	ctx := context.Background()
	releaseName := "test-release"
	namespace := "test-namespace"
	expectedStatus := &ReleaseStatus{
		Name:      releaseName,
		Namespace: namespace,
		Status:    "failed",
		Revision:  3,
	}

	// Configure custom behavior
	mock.StatusFunc = func(ctx context.Context, releaseName, namespace string) (*ReleaseStatus, error) {
		return expectedStatus, nil
	}

	// Test custom behavior
	status, err := mock.Status(ctx, releaseName, namespace)
	if err != nil {
		t.Errorf("Status() unexpected error = %v", err)
	}
	if status.Status != "failed" {
		t.Errorf("Status().Status = %s, want failed", status.Status)
	}
	if status.Revision != 3 {
		t.Errorf("Status().Revision = %d, want 3", status.Revision)
	}
}

func TestMockHelmClient_List(t *testing.T) {
	mock := NewMockHelmClient()
	ctx := context.Background()
	namespace := "test-namespace"

	// Test default behavior
	releases, err := mock.List(ctx, namespace)
	if err != nil {
		t.Errorf("List() unexpected error = %v", err)
	}
	if releases == nil {
		t.Fatal("List() returned nil")
	}
	if len(releases) != 0 {
		t.Errorf("List() returned %d releases, want 0", len(releases))
	}

	// Verify call was recorded
	if len(mock.ListCalls) != 1 {
		t.Errorf("Expected 1 List call, got %d", len(mock.ListCalls))
	}

	// Verify call count helper
	if mock.GetListCallCount() != 1 {
		t.Errorf("GetListCallCount() = %d, want 1", mock.GetListCallCount())
	}
}

func TestMockHelmClient_ListWithCustomFunc(t *testing.T) {
	mock := NewMockHelmClient()
	ctx := context.Background()
	namespace := "test-namespace"
	expectedReleases := []*Release{
		{Name: "release1", Namespace: namespace, Status: "deployed"},
		{Name: "release2", Namespace: namespace, Status: "deployed"},
	}

	// Configure custom behavior
	mock.ListFunc = func(ctx context.Context, namespace string) ([]*Release, error) {
		return expectedReleases, nil
	}

	// Test custom behavior
	releases, err := mock.List(ctx, namespace)
	if err != nil {
		t.Errorf("List() unexpected error = %v", err)
	}
	if len(releases) != 2 {
		t.Errorf("List() returned %d releases, want 2", len(releases))
	}
}

func TestMockHelmClient_Reset(t *testing.T) {
	mock := NewMockHelmClient()
	ctx := context.Background()

	// Make several calls
	_ = mock.Install(ctx, InstallOptions{ReleaseName: "test"})
	_ = mock.Upgrade(ctx, UpgradeOptions{ReleaseName: "test"})
	_, _ = mock.Status(ctx, "test", "default")

	// Verify calls were recorded
	if len(mock.InstallCalls) != 1 {
		t.Errorf("Expected 1 Install call before reset, got %d", len(mock.InstallCalls))
	}
	if len(mock.UpgradeCalls) != 1 {
		t.Errorf("Expected 1 Upgrade call before reset, got %d", len(mock.UpgradeCalls))
	}
	if len(mock.StatusCalls) != 1 {
		t.Errorf("Expected 1 Status call before reset, got %d", len(mock.StatusCalls))
	}

	// Reset the mock
	mock.Reset()

	// Verify all calls were cleared
	if len(mock.InstallCalls) != 0 {
		t.Errorf("Expected 0 Install calls after reset, got %d", len(mock.InstallCalls))
	}
	if len(mock.UpgradeCalls) != 0 {
		t.Errorf("Expected 0 Upgrade calls after reset, got %d", len(mock.UpgradeCalls))
	}
	if len(mock.UninstallCalls) != 0 {
		t.Errorf("Expected 0 Uninstall calls after reset, got %d", len(mock.UninstallCalls))
	}
	if len(mock.StatusCalls) != 0 {
		t.Errorf("Expected 0 Status calls after reset, got %d", len(mock.StatusCalls))
	}
	if len(mock.ListCalls) != 0 {
		t.Errorf("Expected 0 List calls after reset, got %d", len(mock.ListCalls))
	}
}

func TestMockHelmClient_MultipleCalls(t *testing.T) {
	mock := NewMockHelmClient()
	ctx := context.Background()

	// Make multiple calls to the same method
	_ = mock.Install(ctx, InstallOptions{ReleaseName: "release1"})
	_ = mock.Install(ctx, InstallOptions{ReleaseName: "release2"})
	_ = mock.Install(ctx, InstallOptions{ReleaseName: "release3"})

	// Verify all calls were recorded
	if len(mock.InstallCalls) != 3 {
		t.Errorf("Expected 3 Install calls, got %d", len(mock.InstallCalls))
	}

	// Verify call count helper
	if mock.GetInstallCallCount() != 3 {
		t.Errorf("GetInstallCallCount() = %d, want 3", mock.GetInstallCallCount())
	}
}

func TestMockHelmClient_Timestamps(t *testing.T) {
	mock := NewMockHelmClient()
	ctx := context.Background()

	// Make a call
	_ = mock.Install(ctx, InstallOptions{ReleaseName: "test"})

	// Verify timestamp was recorded
	if mock.InstallCalls[0].Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}

	// Verify timestamp is recent (within last second)
	if time.Since(mock.InstallCalls[0].Timestamp) > time.Second {
		t.Error("Timestamp is not recent")
	}
}
