//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestThermal_RealCamera(t *testing.T) {
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

	if !client.HasThermalService() {
		t.Skip("Thermal service not available on this camera")
	}

	t.Run("GetThermalServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetThermalServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetThermalServiceCapabilities not supported: %v", err)
		}

		t.Logf("Thermal service capabilities: Radiometry=%v", caps.Radiometry)
	})

	t.Run("GetThermalConfigurations", func(t *testing.T) {
		cfgs, err := client.GetThermalConfigurations(ctx)
		if err != nil {
			t.Skipf("GetThermalConfigurations not supported: %v", err)
		}

		t.Logf("Found %d thermal configurations", len(cfgs))
		for _, cfg := range cfgs {
			t.Logf("  Configuration: Token=%s, Polarity=%s, ColorPalette=%s",
				cfg.Token, cfg.Configuration.Polarity, cfg.Configuration.ColorPalette.Token)
		}
	})
}
