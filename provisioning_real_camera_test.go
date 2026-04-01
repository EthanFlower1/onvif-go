//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestProvisioning_RealCamera(t *testing.T) {
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

	if !client.HasProvisioningService() {
		t.Skip("Provisioning service not available on this camera")
	}

	t.Run("GetProvisioningServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetProvisioningServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetProvisioningServiceCapabilities not supported: %v", err)
		}

		t.Logf("Provisioning service capabilities: DefaultTimeout=%s, Sources=%d",
			caps.DefaultTimeout, len(caps.Source))
		for _, src := range caps.Source {
			t.Logf("  Source: Token=%s, MaxPanMoves=%v, MaxTiltMoves=%v, MaxZoomMoves=%v",
				src.VideoSourceToken, src.MaximumPanMoves, src.MaximumTiltMoves, src.MaximumZoomMoves)
		}
	})

	t.Run("GetProvisioningUsage", func(t *testing.T) {
		// Get a video source token from capabilities to query usage.
		caps, err := client.GetProvisioningServiceCapabilities(ctx)
		if err != nil || len(caps.Source) == 0 {
			t.Skip("GetProvisioningServiceCapabilities not supported or returned no sources")
		}

		videoSourceToken := caps.Source[0].VideoSourceToken

		usage, err := client.GetProvisioningUsage(ctx, videoSourceToken)
		if err != nil {
			t.Skipf("GetProvisioningUsage not supported: %v", err)
		}

		t.Logf("Provisioning usage for source %s: Pan=%v, Tilt=%v, Zoom=%v, Roll=%v, Focus=%v",
			videoSourceToken, usage.Pan, usage.Tilt, usage.Zoom, usage.Roll, usage.Focus)
	})
}
