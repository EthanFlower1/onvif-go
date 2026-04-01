//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestMedia2_RealCamera(t *testing.T) {
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

	if !client.HasMedia2Service() {
		t.Skip("Media2 service not available on this camera")
	}

	t.Run("GetMedia2ServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetMedia2ServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetMedia2ServiceCapabilities not supported: %v", err)
		}

		t.Logf("Media2 service capabilities: SnapshotUri=%v, Rotation=%v, VideoSourceMode=%v, OSD=%v, Mask=%v, SourceMask=%v",
			caps.SnapshotUri, caps.Rotation, caps.VideoSourceMode, caps.OSD, caps.Mask, caps.SourceMask)
	})

	t.Run("GetMedia2Profiles", func(t *testing.T) {
		profiles, err := client.GetMedia2Profiles(ctx, nil, nil)
		if err != nil {
			t.Skipf("GetMedia2Profiles not supported: %v", err)
		}

		t.Logf("Found %d Media2 profiles", len(profiles))
		for _, p := range profiles {
			t.Logf("  Profile: Token=%s, Name=%s, Fixed=%v", p.Token, p.Name, p.Fixed)
		}
	})

	t.Run("GetMedia2StreamUri", func(t *testing.T) {
		profiles, err := client.GetMedia2Profiles(ctx, nil, nil)
		if err != nil || len(profiles) == 0 {
			t.Skip("GetMedia2Profiles not supported or returned no profiles")
		}

		uri, err := client.GetMedia2StreamUri(ctx, "RtspUnicast", profiles[0].Token)
		if err != nil {
			t.Skipf("GetMedia2StreamUri not supported: %v", err)
		}

		t.Logf("Stream URI: %s", uri)
	})

	t.Run("GetMedia2Masks", func(t *testing.T) {
		masks, err := client.GetMedia2Masks(ctx, nil)
		if err != nil {
			t.Skipf("GetMedia2Masks not supported: %v", err)
		}

		t.Logf("Found %d masks", len(masks))
		for _, m := range masks {
			t.Logf("  Mask: Token=%s, ConfigurationToken=%s, Enabled=%v",
				m.Token, m.ConfigurationToken, m.Enabled)
		}
	})

	t.Run("GetMedia2VideoEncoderInstances", func(t *testing.T) {
		configs, err := client.GetMedia2VideoSourceConfigurations(ctx, nil, nil)
		if err != nil || len(configs) == 0 {
			t.Skip("GetMedia2VideoSourceConfigurations not supported or returned no configs")
		}

		instances, err := client.GetMedia2VideoEncoderInstances(ctx, configs[0].Token)
		if err != nil {
			t.Skipf("GetMedia2VideoEncoderInstances not supported: %v", err)
		}

		h264 := 0
		if instances.H264 != nil {
			h264 = *instances.H264
		}

		jpeg := 0
		if instances.JPEG != nil {
			jpeg = *instances.JPEG
		}

		t.Logf("Video encoder instances: Total=%d, H264=%d, JPEG=%d",
			instances.Total, h264, jpeg)
	})
}
