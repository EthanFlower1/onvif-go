//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestAppMgmt_RealCamera(t *testing.T) {
	endpoint := os.Getenv("ONVIF_ENDPOINT")
	username := os.Getenv("ONVIF_USERNAME")
	password := os.Getenv("ONVIF_PASSWORD")
	if endpoint == "" {
		t.Skip("ONVIF_ENDPOINT not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := NewClient(endpoint, WithCredentials(username, password))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	if err := client.Initialize(ctx); err != nil {
		t.Fatalf("Initialize: %v", err)
	}

	if !client.HasAppMgmtService() {
		t.Skip("App Management service not available on this camera")
	}

	t.Run("GetAppMgmtServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetAppMgmtServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetAppMgmtServiceCapabilities not supported: %v", err)
		}

		t.Logf("App Management capabilities: FormatsSupported=%s, Licensing=%v, UploadPath=%s",
			caps.FormatsSupported, caps.Licensing, caps.UploadPath)
	})

	t.Run("GetInstalledApps", func(t *testing.T) {
		apps, err := client.GetInstalledApps(ctx)
		if err != nil {
			t.Skipf("GetInstalledApps not supported: %v", err)
		}

		t.Logf("Found %d installed apps", len(apps))
		for _, app := range apps {
			t.Logf("  App: AppID=%s, Name=%s", app.AppID, app.Name)
		}
	})
}
