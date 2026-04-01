//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestReplay_RealCamera(t *testing.T) {
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

	if !client.HasReplayService() {
		t.Skip("Replay service not available")
	}

	t.Run("GetReplayServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetReplayServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetReplayServiceCapabilities not supported: %v", err)
		}

		t.Logf("Replay service capabilities: ReversePlayback=%v, RTPRTSP_TCP=%v",
			caps.ReversePlayback, caps.RTPRTSP_TCP)
		if caps.SessionTimeoutRange != nil {
			t.Logf("  SessionTimeoutRange: min=%s, max=%s",
				caps.SessionTimeoutRange.Min, caps.SessionTimeoutRange.Max)
		}
	})

	t.Run("GetReplayConfiguration", func(t *testing.T) {
		config, err := client.GetReplayConfiguration(ctx)
		if err != nil {
			t.Skipf("GetReplayConfiguration not supported: %v", err)
		}

		t.Logf("Replay configuration: SessionTimeout=%s", config.SessionTimeout)
	})

	t.Run("GetReplayUri", func(t *testing.T) {
		if !client.HasRecordingService() {
			t.Skip("Recording service not available to obtain recording token")
		}

		recordings, err := client.GetRecordings(ctx)
		if err != nil || len(recordings) == 0 {
			t.Skip("No recordings available to test GetReplayUri")
		}

		recordingToken := recordings[0].Token

		uri, err := client.GetReplayUri(ctx, recordingToken, "RTP-Unicast", "RTSP")
		if err != nil {
			t.Skipf("GetReplayUri not supported: %v", err)
		}

		t.Logf("Replay URI for token=%s: %s", recordingToken, uri)
	})

}
