//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestDeviceIO_RealCamera(t *testing.T) {
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

	t.Run("GetDeviceIOServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetDeviceIOServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetDeviceIOServiceCapabilities not supported: %v", err)
		}

		t.Logf("DeviceIO service capabilities: VideoSources=%d, VideoOutputs=%d, AudioSources=%d, AudioOutputs=%d, RelayOutputs=%d",
			caps.VideoSources, caps.VideoOutputs, caps.AudioSources, caps.AudioOutputs, caps.RelayOutputs)
	})

	t.Run("GetDeviceIOAudioSources", func(t *testing.T) {
		sources, err := client.GetDeviceIOAudioSources(ctx)
		if err != nil {
			t.Skipf("GetDeviceIOAudioSources not supported: %v", err)
		}

		t.Logf("Found %d DeviceIO audio sources", len(sources))
		for _, s := range sources {
			t.Logf("  Audio source token: %s", s)
		}
	})

	t.Run("GetDeviceIOVideoSources", func(t *testing.T) {
		sources, err := client.GetDeviceIOVideoSources(ctx)
		if err != nil {
			t.Skipf("GetDeviceIOVideoSources not supported: %v", err)
		}

		t.Logf("Found %d DeviceIO video sources", len(sources))
		for _, s := range sources {
			t.Logf("  Video source token: %s", s)
		}
	})

	t.Run("GetDeviceIORelayOutputs", func(t *testing.T) {
		outputs, err := client.GetDeviceIORelayOutputs(ctx)
		if err != nil {
			t.Skipf("GetDeviceIORelayOutputs not supported: %v", err)
		}

		t.Logf("Found %d relay outputs", len(outputs))
		for _, o := range outputs {
			t.Logf("  Relay output: Token=%s, Mode=%s, IdleState=%s",
				o.Token, o.Properties.Mode, o.Properties.IdleState)
		}
	})
}
