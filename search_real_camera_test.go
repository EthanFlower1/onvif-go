//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestSearch_RealCamera(t *testing.T) {
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

	if !client.HasSearchService() {
		t.Skip("Search service not available")
	}

	t.Run("GetSearchServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetSearchServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetSearchServiceCapabilities not supported: %v", err)
		}

		t.Logf("Search service capabilities: MetadataSearch=%v", caps.MetadataSearch)
	})

	t.Run("GetRecordingSummary", func(t *testing.T) {
		summary, err := client.GetRecordingSummary(ctx)
		if err != nil {
			t.Skipf("GetRecordingSummary not supported: %v", err)
		}

		t.Logf("Recording summary: NumberRecordings=%d, DataFrom=%s, DataUntil=%s",
			summary.NumberRecordings, summary.DataFrom.Format(time.RFC3339), summary.DataUntil.Format(time.RFC3339))
	})

	// Get a recording token from the recording service to use in chained tests.
	var recordingToken string

	if client.HasRecordingService() {
		recordings, err := client.GetRecordings(ctx)
		if err == nil && len(recordings) > 0 {
			recordingToken = recordings[0].Token
		}
	}

	t.Run("GetRecordingInformation", func(t *testing.T) {
		if recordingToken == "" {
			t.Skip("No recordings available to test GetRecordingInformation")
		}

		info, err := client.GetRecordingInformation(ctx, recordingToken)
		if err != nil {
			t.Skipf("GetRecordingInformation not supported: %v", err)
		}

		t.Logf("Recording information for token=%s: status=%s, tracks=%d",
			info.RecordingToken, info.RecordingStatus, len(info.TrackInformation))
		if info.EarliestRecording != nil {
			t.Logf("  EarliestRecording=%s", info.EarliestRecording.Format(time.RFC3339))
		}
		if info.LatestRecording != nil {
			t.Logf("  LatestRecording=%s", info.LatestRecording.Format(time.RFC3339))
		}
	})

	t.Run("FindRecordings_GetRecordingSearchResults", func(t *testing.T) {
		searchToken, err := client.FindRecordings(ctx, nil, nil, "PT10S")
		if err != nil {
			t.Skipf("FindRecordings not supported: %v", err)
		}

		t.Logf("FindRecordings returned searchToken=%s", searchToken)

		maxResults := 10
		results, err := client.GetRecordingSearchResults(ctx, searchToken, nil, &maxResults, "PT5S")
		if err != nil {
			t.Skipf("GetRecordingSearchResults not supported: %v", err)
		}

		t.Logf("GetRecordingSearchResults: searchState=%s, recordings=%d",
			results.SearchState, len(results.RecordingInformation))
		for _, ri := range results.RecordingInformation {
			t.Logf("  Result token=%s, status=%s", ri.RecordingToken, ri.RecordingStatus)
		}
	})
}
