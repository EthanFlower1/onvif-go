//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestReceiver_RealCamera(t *testing.T) {
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

	if !client.HasReceiverService() {
		t.Skip("Receiver service not available")
	}

	t.Run("GetReceiverServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetReceiverServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetReceiverServiceCapabilities not supported: %v", err)
		}

		t.Logf("Receiver service capabilities: RTPMulticast=%v, RTPTCP=%v, RTPRTSP_TCP=%v, SupportedReceivers=%d, MaximumRTSPURILength=%d",
			caps.RTPMulticast, caps.RTPTCP, caps.RTPRTSP_TCP, caps.SupportedReceivers, caps.MaximumRTSPURILength)
	})

	t.Run("GetReceivers", func(t *testing.T) {
		receivers, err := client.GetReceivers(ctx)
		if err != nil {
			t.Skipf("GetReceivers not supported: %v", err)
		}

		t.Logf("Found %d receivers", len(receivers))
		for _, rec := range receivers {
			t.Logf("  Receiver token=%s, mode=%s, mediaURI=%s",
				rec.Token, rec.Configuration.Mode, rec.Configuration.MediaURI)
		}
	})
}
