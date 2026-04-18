package helm

import (
	"context"
	"testing"
	"time"

	"github.com/devplatform/devplatform-cli/internal/logger"
	"github.com/devplatform/devplatform-cli/test/testutil"
)

// TestNewClient tests creating a new Helm client
func TestNewClient(t *testing.T) {
	log := *logger.New(logger.InfoLevel, true)
	client := NewClient(log)

	testutil.AssertTrue(t, client != nil, "Client should not be nil")
	testutil.AssertEqual(t, "helm", client.helmBinary)
}

// TestInstallOptions tests InstallOptions structure
func TestInstallOptions(t *testing.T) {
	opts := InstallOptions{
		ReleaseName:     "myapp",
		Chart:           "stable/nginx",
		Namespace:       "default",
		ValuesFiles:     []string{"values.yaml"},
		Values:          map[string]interface{}{"replicas": 3},
		CreateNamespace: true,
		Wait:            true,
		Timeout:         5 * time.Minute,
	}

	testutil.AssertEqual(t, "myapp", opts.ReleaseName)
	testutil.AssertEqual(t, "stable/nginx", opts.Chart)
	testutil.AssertEqual(t, "default", opts.Namespace)
	testutil.AssertEqual(t, 1, len(opts.ValuesFiles))
	testutil.AssertEqual(t, true, opts.CreateNamespace)
	testutil.AssertEqual(t, true, opts.Wait)
	testutil.AssertEqual(t, 5*time.Minute, opts.Timeout)
}

// TestUpgradeOptions tests UpgradeOptions structure
func TestUpgradeOptions(t *testing.T) {
	opts := UpgradeOptions{
		ReleaseName: "myapp",
		Chart:       "stable/nginx",
		Namespace:   "default",
		ValuesFiles: []string{"values.yaml", "values-prod.yaml"},
		Values:      map[string]interface{}{"replicas": 5},
		Install:     true,
		Wait:        true,
		Timeout:     10 * time.Minute,
	}

	testutil.AssertEqual(t, "myapp", opts.ReleaseName)
	testutil.AssertEqual(t, "stable/nginx", opts.Chart)
	testutil.AssertEqual(t, 2, len(opts.ValuesFiles))
	testutil.AssertEqual(t, true, opts.Install)
	testutil.AssertEqual(t, true, opts.Wait)
}

// TestUninstallOptions tests UninstallOptions structure
func TestUninstallOptions(t *testing.T) {
	opts := UninstallOptions{
		ReleaseName: "myapp",
		Namespace:   "default",
		Wait:        true,
		Timeout:     5 * time.Minute,
	}

	testutil.AssertEqual(t, "myapp", opts.ReleaseName)
	testutil.AssertEqual(t, "default", opts.Namespace)
	testutil.AssertEqual(t, true, opts.Wait)
	testutil.AssertEqual(t, 5*time.Minute, opts.Timeout)
}

// TestReleaseStatus tests ReleaseStatus structure
func TestReleaseStatus(t *testing.T) {
	now := time.Now()
	status := &ReleaseStatus{
		Name:       "myapp",
		Namespace:  "default",
		Status:     "deployed",
		Revision:   3,
		Updated:    now,
		Chart:      "nginx-1.2.3",
		AppVersion: "1.19.0",
	}

	testutil.AssertEqual(t, "myapp", status.Name)
	testutil.AssertEqual(t, "default", status.Namespace)
	testutil.AssertEqual(t, "deployed", status.Status)
	testutil.AssertEqual(t, 3, status.Revision)
	testutil.AssertEqual(t, "nginx-1.2.3", status.Chart)
	testutil.AssertEqual(t, "1.19.0", status.AppVersion)
}

// TestRelease tests Release structure
func TestRelease(t *testing.T) {
	now := time.Now()
	release := &Release{
		Name:       "myapp",
		Namespace:  "production",
		Revision:   5,
		Updated:    now,
		Status:     "deployed",
		Chart:      "myapp-2.0.0",
		AppVersion: "2.0.0",
	}

	testutil.AssertEqual(t, "myapp", release.Name)
	testutil.AssertEqual(t, "production", release.Namespace)
	testutil.AssertEqual(t, 5, release.Revision)
	testutil.AssertEqual(t, "deployed", release.Status)
	testutil.AssertEqual(t, "myapp-2.0.0", release.Chart)
}

// TestMockClientInstall tests the mock client Install method
func TestMockClientInstall(t *testing.T) {
	mock := NewMockHelmClient()
	
	// Configure mock to succeed
	mock.InstallFunc = func(ctx context.Context, opts InstallOptions) error {
		return nil
	}

	opts := InstallOptions{
		ReleaseName: "test-release",
		Chart:       "test-chart",
		Namespace:   "test-ns",
	}

	err := mock.Install(context.Background(), opts)
	testutil.AssertNoError(t, err)

	// Verify call was recorded
	testutil.AssertEqual(t, 1, len(mock.InstallCalls))
	recordedOpts := mock.InstallCalls[0].Args[1].(InstallOptions)
	testutil.AssertEqual(t, "test-release", recordedOpts.ReleaseName)
	testutil.AssertEqual(t, "test-chart", recordedOpts.Chart)
}

// TestMockClientUpgrade tests the mock client Upgrade method
func TestMockClientUpgrade(t *testing.T) {
	mock := NewMockHelmClient()
	
	mock.UpgradeFunc = func(ctx context.Context, opts UpgradeOptions) error {
		return nil
	}

	opts := UpgradeOptions{
		ReleaseName: "test-release",
		Chart:       "test-chart",
		Namespace:   "test-ns",
		Install:     true,
	}

	err := mock.Upgrade(context.Background(), opts)
	testutil.AssertNoError(t, err)

	testutil.AssertEqual(t, 1, len(mock.UpgradeCalls))
	recordedOpts := mock.UpgradeCalls[0].Args[1].(UpgradeOptions)
	testutil.AssertEqual(t, true, recordedOpts.Install)
}

// TestMockClientUninstall tests the mock client Uninstall method
func TestMockClientUninstall(t *testing.T) {
	mock := NewMockHelmClient()
	
	mock.UninstallFunc = func(ctx context.Context, opts UninstallOptions) error {
		return nil
	}

	opts := UninstallOptions{
		ReleaseName: "test-release",
		Namespace:   "test-ns",
	}

	err := mock.Uninstall(context.Background(), opts)
	testutil.AssertNoError(t, err)

	testutil.AssertEqual(t, 1, len(mock.UninstallCalls))
	recordedOpts := mock.UninstallCalls[0].Args[1].(UninstallOptions)
	testutil.AssertEqual(t, "test-release", recordedOpts.ReleaseName)
}

// TestMockClientStatus tests the mock client Status method
func TestMockClientStatus(t *testing.T) {
	mock := NewMockHelmClient()
	
	expectedStatus := &ReleaseStatus{
		Name:      "test-release",
		Namespace: "test-ns",
		Status:    "deployed",
	}

	mock.StatusFunc = func(ctx context.Context, releaseName, namespace string) (*ReleaseStatus, error) {
		return expectedStatus, nil
	}

	status, err := mock.Status(context.Background(), "test-release", "test-ns")
	testutil.AssertNoError(t, err)
	testutil.AssertEqual(t, "deployed", status.Status)

	testutil.AssertEqual(t, 1, len(mock.StatusCalls))
}

// TestMockClientList tests the mock client List method
func TestMockClientList(t *testing.T) {
	mock := NewMockHelmClient()
	
	expectedReleases := []*Release{
		{Name: "release1", Status: "deployed"},
		{Name: "release2", Status: "failed"},
	}

	mock.ListFunc = func(ctx context.Context, namespace string) ([]*Release, error) {
		return expectedReleases, nil
	}

	releases, err := mock.List(context.Background(), "test-ns")
	testutil.AssertNoError(t, err)
	testutil.AssertEqual(t, 2, len(releases))
	testutil.AssertEqual(t, "release1", releases[0].Name)

	testutil.AssertEqual(t, 1, len(mock.ListCalls))
}

// TestMockClientReset tests the mock client Reset method
func TestMockClientReset(t *testing.T) {
	mock := NewMockHelmClient()

	// Make some calls
	_ = mock.Install(context.Background(), InstallOptions{})
	_ = mock.Upgrade(context.Background(), UpgradeOptions{})
	_ = mock.Uninstall(context.Background(), UninstallOptions{})

	testutil.AssertEqual(t, 1, len(mock.InstallCalls))
	testutil.AssertEqual(t, 1, len(mock.UpgradeCalls))
	testutil.AssertEqual(t, 1, len(mock.UninstallCalls))

	// Reset
	mock.Reset()

	testutil.AssertEqual(t, 0, len(mock.InstallCalls))
	testutil.AssertEqual(t, 0, len(mock.UpgradeCalls))
	testutil.AssertEqual(t, 0, len(mock.UninstallCalls))
}

// TestMockClientCallCounts tests the call count methods
func TestMockClientCallCounts(t *testing.T) {
	mock := NewMockHelmClient()

	testutil.AssertEqual(t, 0, mock.GetInstallCallCount())
	testutil.AssertEqual(t, 0, mock.GetUpgradeCallCount())

	_ = mock.Install(context.Background(), InstallOptions{})
	_ = mock.Install(context.Background(), InstallOptions{})
	_ = mock.Upgrade(context.Background(), UpgradeOptions{})

	testutil.AssertEqual(t, 2, mock.GetInstallCallCount())
	testutil.AssertEqual(t, 1, mock.GetUpgradeCallCount())
}

// TestInstallOptionsWithMultipleValues tests install with multiple value sources
func TestInstallOptionsWithMultipleValues(t *testing.T) {
	opts := InstallOptions{
		ReleaseName: "myapp",
		Chart:       "mychart",
		Namespace:   "default",
		ValuesFiles: []string{"base.yaml", "env.yaml", "secrets.yaml"},
		Values: map[string]interface{}{
			"replicas":    3,
			"image.tag":   "v1.2.3",
			"service.port": 8080,
		},
		CreateNamespace: true,
		Wait:            true,
		Timeout:         10 * time.Minute,
	}

	testutil.AssertEqual(t, 3, len(opts.ValuesFiles))
	testutil.AssertEqual(t, 3, len(opts.Values))
	testutil.AssertEqual(t, 3, opts.Values["replicas"])
	testutil.AssertEqual(t, "v1.2.3", opts.Values["image.tag"])
}

// TestUpgradeOptionsWithInstallFlag tests upgrade with --install flag
func TestUpgradeOptionsWithInstallFlag(t *testing.T) {
	opts := UpgradeOptions{
		ReleaseName: "myapp",
		Chart:       "mychart",
		Namespace:   "default",
		Install:     true, // This enables --install flag
		Wait:        true,
	}

	testutil.AssertEqual(t, true, opts.Install)
	testutil.AssertEqual(t, true, opts.Wait)
}

// TestContextCancellation tests that operations respect context cancellation
func TestContextCancellation(t *testing.T) {
	mock := NewMockHelmClient()

	// Configure mock to check context
	mock.InstallFunc = func(ctx context.Context, opts InstallOptions) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return nil
	}

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := mock.Install(ctx, InstallOptions{})
	testutil.AssertError(t, err)
	testutil.AssertContains(t, err.Error(), "context canceled")
}

// TestTimeoutHandling tests timeout handling
func TestTimeoutHandling(t *testing.T) {
	opts := InstallOptions{
		ReleaseName: "myapp",
		Chart:       "mychart",
		Namespace:   "default",
		Wait:        true,
		Timeout:     1 * time.Second, // Very short timeout
	}

	testutil.AssertEqual(t, 1*time.Second, opts.Timeout)
	testutil.AssertEqual(t, true, opts.Wait)
}
