//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestUplink_RealCamera(t *testing.T) {
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

	if !client.HasUplinkService() {
		t.Skip("Uplink service not available on this camera")
	}

	t.Run("GetUplinkServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetUplinkServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetUplinkServiceCapabilities not supported: %v", err)
		}

		t.Logf("Uplink service capabilities: MaxUplinks=%v", caps.MaxUplinks)
	})

	t.Run("GetUplinks", func(t *testing.T) {
		uplinks, err := client.GetUplinks(ctx)
		if err != nil {
			t.Skipf("GetUplinks not supported: %v", err)
		}

		t.Logf("Found %d uplink configurations", len(uplinks))
		for _, u := range uplinks {
			t.Logf("  Uplink: RemoteAddress=%s, UserLevel=%s, Status=%v",
				u.RemoteAddress, u.UserLevel, u.Status)
		}
	})
}
