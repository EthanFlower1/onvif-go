//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestAccessControl_RealCamera(t *testing.T) {
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

	if !client.HasAccessControlService() {
		t.Skip("Access Control service not available on this camera")
	}

	t.Run("GetAccessControlServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetAccessControlServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetAccessControlServiceCapabilities not supported: %v", err)
		}

		t.Logf("Access Control service capabilities: MaxLimit=%d, MaxAccessPoints=%d, MaxAreas=%d, ClientSuppliedTokenSupported=%v, AccessPointManagementSupported=%v, AreaManagementSupported=%v",
			caps.MaxLimit, caps.MaxAccessPoints, caps.MaxAreas,
			caps.ClientSuppliedTokenSupported, caps.AccessPointManagementSupported, caps.AreaManagementSupported)
	})

	t.Run("GetAccessPointInfoList", func(t *testing.T) {
		infos, _, err := client.GetAccessPointInfoList(ctx, nil, nil)
		if err != nil {
			t.Skipf("GetAccessPointInfoList not supported: %v", err)
		}

		t.Logf("Found %d access points", len(infos))
		for _, ap := range infos {
			t.Logf("  AccessPoint: Token=%s, Name=%s, AreaFrom=%s, AreaTo=%s",
				ap.Token, ap.Name, ap.AreaFrom, ap.AreaTo)
		}
	})

	t.Run("GetAreaInfoList", func(t *testing.T) {
		infos, _, err := client.GetAreaInfoList(ctx, nil, nil)
		if err != nil {
			t.Skipf("GetAreaInfoList not supported: %v", err)
		}

		t.Logf("Found %d areas", len(infos))
		for _, area := range infos {
			t.Logf("  Area: Token=%s, Name=%s", area.Token, area.Name)
		}
	})

	t.Run("GetAccessPointState", func(t *testing.T) {
		infos, _, err := client.GetAccessPointInfoList(ctx, nil, nil)
		if err != nil || len(infos) == 0 {
			t.Skip("GetAccessPointInfoList not supported or returned no access points")
		}

		token := infos[0].Token
		state, err := client.GetAccessPointState(ctx, token)
		if err != nil {
			t.Skipf("GetAccessPointState not supported: %v", err)
		}

		t.Logf("AccessPoint state for token=%s: Enabled=%v",
			token, state.Enabled)
	})
}
