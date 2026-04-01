//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestRecording_RealCamera(t *testing.T) {
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

	if !client.HasRecordingService() {
		t.Skip("Recording service not available")
	}

	t.Run("GetRecordingServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetRecordingServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetRecordingServiceCapabilities not supported: %v", err)
		}

		t.Logf("Recording service capabilities: DynamicRecordings=%v, DynamicTracks=%v, MaxRecordings=%d, MaxRecordingJobs=%d",
			caps.DynamicRecordings, caps.DynamicTracks, caps.MaxRecordings, caps.MaxRecordingJobs)
	})

	var recordingToken string

	t.Run("GetRecordings", func(t *testing.T) {
		recordings, err := client.GetRecordings(ctx)
		if err != nil {
			t.Skipf("GetRecordings not supported: %v", err)
		}

		t.Logf("Found %d recordings", len(recordings))
		for _, rec := range recordings {
			t.Logf("  Recording token=%s, name=%s, tracks=%d",
				rec.Token, rec.Configuration.Source.Name, len(rec.Tracks))
		}

		if len(recordings) > 0 {
			recordingToken = recordings[0].Token
		}
	})

	t.Run("GetRecordingJobs", func(t *testing.T) {
		jobs, err := client.GetRecordingJobs(ctx)
		if err != nil {
			t.Skipf("GetRecordingJobs not supported: %v", err)
		}

		t.Logf("Found %d recording jobs", len(jobs))
		for _, job := range jobs {
			t.Logf("  Job token=%s, recordingToken=%s, mode=%s",
				job.Token, job.Configuration.RecordingToken, job.Configuration.Mode)
		}
	})

	t.Run("GetRecordingConfiguration", func(t *testing.T) {
		if recordingToken == "" {
			t.Skip("No recordings available to test GetRecordingConfiguration")
		}

		config, err := client.GetRecordingConfiguration(ctx, recordingToken)
		if err != nil {
			t.Skipf("GetRecordingConfiguration not supported: %v", err)
		}

		t.Logf("Recording config for token=%s: sourceName=%s, content=%s, maxRetention=%s",
			recordingToken, config.Source.Name, config.Content, config.MaximumRetentionTime)
	})
}
