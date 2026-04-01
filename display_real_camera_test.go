//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestDisplay_RealCamera(t *testing.T) {
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

	if !client.HasDisplayService() {
		t.Skip("Display service not available on this camera")
	}

	t.Run("GetDisplayServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetDisplayServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetDisplayServiceCapabilities not supported: %v", err)
		}

		t.Logf("Display service capabilities: FixedLayout=%v", caps.FixedLayout)
	})

	// getVideoOutputToken returns the first video output token from the device,
	// or skips the test if none are available.
	getVideoOutputToken := func(t *testing.T) string {
		t.Helper()

		outputs, err := client.GetVideoOutputs(ctx)
		if err != nil || len(outputs) == 0 {
			t.Skip("GetVideoOutputs not supported or returned no outputs")
		}

		return outputs[0].Token
	}

	t.Run("GetLayout", func(t *testing.T) {
		token := getVideoOutputToken(t)

		layout, err := client.GetLayout(ctx, token)
		if err != nil {
			t.Skipf("GetLayout not supported: %v", err)
		}

		t.Logf("Layout for output %s: %d pane layouts", token, len(layout.Pane))
	})

	t.Run("GetPaneConfigurations", func(t *testing.T) {
		token := getVideoOutputToken(t)

		panes, err := client.GetPaneConfigurations(ctx, token)
		if err != nil {
			t.Skipf("GetPaneConfigurations not supported: %v", err)
		}

		t.Logf("Found %d pane configurations for output %s", len(panes), token)
		for _, p := range panes {
			t.Logf("  Pane: Token=%s, Name=%s", p.Token, p.PaneName)
		}
	})
}
