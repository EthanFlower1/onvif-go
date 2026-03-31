//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestImaging_RealCamera(t *testing.T) {
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

	t.Run("GetImagingServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetImagingServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetImagingServiceCapabilities not supported: %v", err)
		}

		t.Logf("Imaging service capabilities: ImageStabilization=%v, Presets=%v",
			caps.ImageStabilization, caps.Presets)
	})

	t.Run("GetImagingPresets", func(t *testing.T) {
		sources, err := client.GetVideoSources(ctx)
		if err != nil || len(sources) == 0 {
			t.Skip("GetVideoSources not supported or returned no sources")
		}

		videoSourceToken := sources[0].Token
		presets, err := client.GetImagingPresets(ctx, videoSourceToken)
		if err != nil {
			t.Skipf("GetImagingPresets not supported: %v", err)
		}

		t.Logf("Found %d imaging presets for video source %s", len(presets), videoSourceToken)
		for _, p := range presets {
			t.Logf("  Preset: Token=%s, Name=%s, Type=%s", p.Token, p.Name, p.Type)
		}
	})

	t.Run("GetCurrentImagingPreset", func(t *testing.T) {
		sources, err := client.GetVideoSources(ctx)
		if err != nil || len(sources) == 0 {
			t.Skip("GetVideoSources not supported or returned no sources")
		}

		videoSourceToken := sources[0].Token
		preset, err := client.GetCurrentImagingPreset(ctx, videoSourceToken)
		if err != nil {
			t.Skipf("GetCurrentImagingPreset not supported: %v", err)
		}

		if preset == nil {
			t.Skip("No current imaging preset returned")
		}

		t.Logf("Current imaging preset: Token=%s, Name=%s, Type=%s",
			preset.Token, preset.Name, preset.Type)
	})
}
