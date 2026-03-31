# Phase 3: Profile G (Recording Ecosystem) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the complete Recording, Receiver, Search, and Replay client services — the four services required for Profile G (NVR recording/playback) compliance.

**Architecture:** Four new service files with per-service type files. All operations follow the established callMethod pattern. Recording and Receiver are independent; Search and Replay depend on Recording types. Each operation gets unit tests + integration test scaffolding.

**Tech Stack:** Go stdlib (encoding/xml, net/http, context, time), internal/soap package, httptest for mocks.

**Depends on:** Phase 1 (Foundation) must be complete — endpoint discovery and testing infrastructure.

**Spec:** `docs/superpowers/specs/2026-03-31-onvif-full-compliance-design.md`

---

## File Structure

**New files:**
- `types_recording.go` — Recording, RecordingJob, RecordingConfiguration, TrackConfiguration types
- `types_search.go` — SearchScope, RecordingSummary, search result types
- `types_replay.go` — ReplayConfiguration type
- `types_receiver.go` — Receiver, ReceiverConfiguration types
- `recording.go` — 22 Recording service client methods
- `recording_test.go` — Unit tests for all recording operations
- `recording_real_camera_test.go` — Integration tests
- `search.go` — 14 Search service client methods
- `search_test.go` — Unit tests
- `search_real_camera_test.go` — Integration tests
- `replay.go` — 4 Replay service client methods
- `replay_test.go` — Unit tests
- `replay_real_camera_test.go` — Integration tests
- `receiver.go` — 8 Receiver service client methods
- `receiver_test.go` — Unit tests
- `receiver_real_camera_test.go` — Integration tests

---

## Task 1: Recording Service Types

**Files:**
- Create: `types_recording.go`

- [ ] **Step 1: Create types_recording.go**

```go
package onvif

// Recording represents an ONVIF recording.
type Recording struct {
	Token         string
	Configuration RecordingConfiguration
	Tracks        []*RecordingTrack
}

// RecordingConfiguration contains recording configuration settings.
type RecordingConfiguration struct {
	Source               RecordingSourceInformation
	Content              string
	MaximumRetentionTime string
}

// RecordingSourceInformation identifies the source of a recording.
type RecordingSourceInformation struct {
	SourceId    string
	Name        string
	Location    string
	Description string
	Address     string
}

// RecordingTrack represents a track within a recording.
type RecordingTrack struct {
	Token         string
	Configuration TrackConfiguration
	DataFrom      string
	DataTo        string
}

// TrackConfiguration contains track settings.
type TrackConfiguration struct {
	TrackType   string
	Description string
}

// RecordingJob represents a recording job.
type RecordingJob struct {
	Token         string
	Configuration RecordingJobConfiguration
}

// RecordingJobConfiguration contains recording job settings.
type RecordingJobConfiguration struct {
	RecordingToken string
	Mode           string
	Priority       int
	Source         []*RecordingJobSource
}

// RecordingJobSource identifies a source for recording.
type RecordingJobSource struct {
	SourceToken  *SourceReference
	AutoCreateReceiver bool
	Tracks       []*RecordingJobTrack
}

// SourceReference identifies a source by profile and endpoint.
type SourceReference struct {
	Token string
	Type  string
}

// RecordingJobTrack maps a source track to a recording track.
type RecordingJobTrack struct {
	SourceTag      string
	Destination    string
}

// RecordingJobState contains recording job state information.
type RecordingJobState struct {
	RecordingToken string
	State          string
	Sources        []*RecordingJobSourceState
}

// RecordingJobSourceState contains the state of a recording job source.
type RecordingJobSourceState struct {
	SourceToken *SourceReference
	State       string
	Tracks      []*RecordingJobTrackState
}

// RecordingJobTrackState contains the state of a recording job track.
type RecordingJobTrackState struct {
	SourceTag   string
	Destination string
	Error       string
	State       string
}

// RecordingOptions contains available recording configuration options.
type RecordingOptions struct {
	Job   *RecordingJobOptions
	Track *RecordingTrackOptions
}

// RecordingJobOptions contains job-related options.
type RecordingJobOptions struct {
	Spare       *int
	CompatibleSources []*string
}

// RecordingTrackOptions contains track-related options.
type RecordingTrackOptions struct {
	SpareTotal      *int
	SpareVideo      *int
	SpareAudio      *int
	SpareMetadata   *int
}

// RecordingServiceCapabilities represents recording service capabilities.
type RecordingServiceCapabilities struct {
	DynamicRecordings       bool
	DynamicTracks           bool
	MaxStringLength         int
	MaxRecordings           int
	MaxRecordingJobs        int
	Options                 bool
	MetadataRecording       bool
	SupportedExportFileFormats []string
}

// ExportRecordedDataState contains export progress information.
type ExportRecordedDataState struct {
	Progress           float64
	FileProgressStatus []*FileProgress
}

// FileProgress contains progress for a single export file.
type FileProgress struct {
	FileName string
	Progress float64
}
```

- [ ] **Step 2: Verify compilation**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go build ./...`
Expected: Clean build.

- [ ] **Step 3: Commit**

```bash
git add types_recording.go
git commit -m "feat: add Recording service type definitions"
```

---

## Task 2: Recording Service — Core Operations (7 ops)

**Files:**
- Create: `recording.go`
- Create: `recording_test.go`

Implements: `GetServiceCapabilities`, `CreateRecording`, `DeleteRecording`, `GetRecordings`, `SetRecordingConfiguration`, `GetRecordingConfiguration`, `GetRecordingOptions`.

- [ ] **Step 1: Write failing tests**

```go
package onvif

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetRecordings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), "trc:GetRecordings") {
			t.Error("expected trc:GetRecordings in request")
		}
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<trc:GetRecordingsResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
<trc:RecordingItem>
<trc:RecordingToken>REC_001</trc:RecordingToken>
<trc:Configuration>
<tt:Source xmlns:tt="http://www.onvif.org/ver10/schema">
<tt:SourceId>SRC_001</tt:SourceId>
<tt:Name>Camera 1</tt:Name>
<tt:Location>Building A</tt:Location>
<tt:Description>Front entrance</tt:Description>
<tt:Address>http://192.168.1.100/onvif/device_service</tt:Address>
</tt:Source>
<tt:Content xmlns:tt="http://www.onvif.org/ver10/schema">Recording from Camera 1</tt:Content>
<tt:MaximumRetentionTime xmlns:tt="http://www.onvif.org/ver10/schema">PT72H</tt:MaximumRetentionTime>
</trc:Configuration>
<trc:Tracks>
<trc:Track>
<trc:TrackToken>VIDEO001</trc:TrackToken>
<trc:Configuration>
<tt:TrackType xmlns:tt="http://www.onvif.org/ver10/schema">Video</tt:TrackType>
<tt:Description xmlns:tt="http://www.onvif.org/ver10/schema">Video track</tt:Description>
</trc:Configuration>
</trc:Track>
</trc:Tracks>
</trc:RecordingItem>
</trc:GetRecordingsResponse>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	recordings, err := client.GetRecordings(context.Background())
	if err != nil {
		t.Fatalf("GetRecordings failed: %v", err)
	}
	if len(recordings) != 1 {
		t.Fatalf("expected 1 recording, got %d", len(recordings))
	}
	if recordings[0].Token != "REC_001" {
		t.Errorf("expected token REC_001, got %s", recordings[0].Token)
	}
	if recordings[0].Configuration.Source.Name != "Camera 1" {
		t.Errorf("expected source name Camera 1, got %s", recordings[0].Configuration.Source.Name)
	}
	if len(recordings[0].Tracks) != 1 {
		t.Fatalf("expected 1 track, got %d", len(recordings[0].Tracks))
	}
	if recordings[0].Tracks[0].Token != "VIDEO001" {
		t.Errorf("expected track token VIDEO001, got %s", recordings[0].Tracks[0].Token)
	}
}

func TestGetRecordings_SOAPFault(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<soap:Fault>
<soap:Code><soap:Value>soap:Sender</soap:Value></soap:Code>
<soap:Reason><soap:Text>Not Authorized</soap:Text></soap:Reason>
</soap:Fault>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	_, err := client.GetRecordings(context.Background())
	if err == nil {
		t.Fatal("expected error for SOAP fault")
	}
}

func TestGetRecordings_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<trc:GetRecordingsResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl"/>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	recordings, err := client.GetRecordings(context.Background())
	if err != nil {
		t.Fatalf("GetRecordings failed: %v", err)
	}
	if len(recordings) != 0 {
		t.Errorf("expected 0 recordings, got %d", len(recordings))
	}
}

func TestCreateRecording(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<trc:CreateRecordingResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
<trc:RecordingToken>REC_NEW</trc:RecordingToken>
</trc:CreateRecordingResponse>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	token, err := client.CreateRecording(context.Background(), &RecordingConfiguration{
		Source: RecordingSourceInformation{Name: "Camera 1"},
		Content: "Test recording",
		MaximumRetentionTime: "PT24H",
	})
	if err != nil {
		t.Fatalf("CreateRecording failed: %v", err)
	}
	if token != "REC_NEW" {
		t.Errorf("expected token REC_NEW, got %s", token)
	}
}

func TestDeleteRecording(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<trc:DeleteRecordingResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl"/>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	err := client.DeleteRecording(context.Background(), "REC_001")
	if err != nil {
		t.Fatalf("DeleteRecording failed: %v", err)
	}
}

func TestGetRecordingServiceCapabilities(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<trc:GetServiceCapabilitiesResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
<trc:Capabilities DynamicRecordings="true" DynamicTracks="true" MaxRecordings="50" MaxRecordingJobs="50"/>
</trc:GetServiceCapabilitiesResponse>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	caps, err := client.GetRecordingServiceCapabilities(context.Background())
	if err != nil {
		t.Fatalf("GetRecordingServiceCapabilities failed: %v", err)
	}
	if !caps.DynamicRecordings {
		t.Error("expected DynamicRecordings to be true")
	}
	if caps.MaxRecordings != 50 {
		t.Errorf("expected MaxRecordings 50, got %d", caps.MaxRecordings)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go test -run "TestGetRecordings|TestCreateRecording|TestDeleteRecording|TestGetRecordingServiceCapabilities" -v`
Expected: FAIL — methods not defined.

- [ ] **Step 3: Implement recording.go with core operations**

```go
package onvif

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/0x524a/onvif-go/internal/soap"
)

const recordingNamespace = "http://www.onvif.org/ver10/recording/wsdl"

func (c *Client) getRecordingEndpoint() string {
	if c.recordingEndpoint != "" {
		return c.recordingEndpoint
	}

	return c.endpoint
}

// GetRecordingServiceCapabilities returns the recording service capabilities.
func (c *Client) GetRecordingServiceCapabilities(ctx context.Context) (*RecordingServiceCapabilities, error) {
	type getServiceCapabilitiesRequest struct {
		XMLName xml.Name `xml:"trc:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:trc,attr"`
	}

	type getServiceCapabilitiesResponse struct {
		Capabilities struct {
			DynamicRecordings bool `xml:"DynamicRecordings,attr"`
			DynamicTracks     bool `xml:"DynamicTracks,attr"`
			MaxRecordings     int  `xml:"MaxRecordings,attr"`
			MaxRecordingJobs  int  `xml:"MaxRecordingJobs,attr"`
		} `xml:"Capabilities"`
	}

	req := getServiceCapabilitiesRequest{Xmlns: recordingNamespace}
	var resp getServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getRecordingEndpoint(), "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRecordingServiceCapabilities failed: %w", err)
	}

	return &RecordingServiceCapabilities{
		DynamicRecordings: resp.Capabilities.DynamicRecordings,
		DynamicTracks:     resp.Capabilities.DynamicTracks,
		MaxRecordings:     resp.Capabilities.MaxRecordings,
		MaxRecordingJobs:  resp.Capabilities.MaxRecordingJobs,
	}, nil
}

// GetRecordings returns all recordings on the device.
func (c *Client) GetRecordings(ctx context.Context) ([]*Recording, error) {
	type getRecordingsRequest struct {
		XMLName xml.Name `xml:"trc:GetRecordings"`
		Xmlns   string   `xml:"xmlns:trc,attr"`
	}

	type getRecordingsResponse struct {
		RecordingItem []struct {
			RecordingToken string `xml:"RecordingToken"`
			Configuration  struct {
				Source struct {
					SourceId    string `xml:"SourceId"`
					Name        string `xml:"Name"`
					Location    string `xml:"Location"`
					Description string `xml:"Description"`
					Address     string `xml:"Address"`
				} `xml:"Source"`
				Content              string `xml:"Content"`
				MaximumRetentionTime string `xml:"MaximumRetentionTime"`
			} `xml:"Configuration"`
			Tracks struct {
				Track []struct {
					TrackToken    string `xml:"TrackToken"`
					Configuration struct {
						TrackType   string `xml:"TrackType"`
						Description string `xml:"Description"`
					} `xml:"Configuration"`
					DataFrom string `xml:"DataFrom"`
					DataTo   string `xml:"DataTo"`
				} `xml:"Track"`
			} `xml:"Tracks"`
		} `xml:"RecordingItem"`
	}

	req := getRecordingsRequest{Xmlns: recordingNamespace}
	var resp getRecordingsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getRecordingEndpoint(), "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRecordings failed: %w", err)
	}

	recordings := make([]*Recording, 0, len(resp.RecordingItem))
	for _, item := range resp.RecordingItem {
		rec := &Recording{
			Token: item.RecordingToken,
			Configuration: RecordingConfiguration{
				Source: RecordingSourceInformation{
					SourceId:    item.Configuration.Source.SourceId,
					Name:        item.Configuration.Source.Name,
					Location:    item.Configuration.Source.Location,
					Description: item.Configuration.Source.Description,
					Address:     item.Configuration.Source.Address,
				},
				Content:              item.Configuration.Content,
				MaximumRetentionTime: item.Configuration.MaximumRetentionTime,
			},
		}
		tracks := make([]*RecordingTrack, 0, len(item.Tracks.Track))
		for _, tr := range item.Tracks.Track {
			tracks = append(tracks, &RecordingTrack{
				Token: tr.TrackToken,
				Configuration: TrackConfiguration{
					TrackType:   tr.Configuration.TrackType,
					Description: tr.Configuration.Description,
				},
				DataFrom: tr.DataFrom,
				DataTo:   tr.DataTo,
			})
		}
		rec.Tracks = tracks
		recordings = append(recordings, rec)
	}

	return recordings, nil
}

// CreateRecording creates a new recording.
func (c *Client) CreateRecording(ctx context.Context, config *RecordingConfiguration) (string, error) {
	type createRecordingRequest struct {
		XMLName                xml.Name `xml:"trc:CreateRecording"`
		Xmlns                  string   `xml:"xmlns:trc,attr"`
		RecordingConfiguration struct {
			Source struct {
				SourceId    string `xml:"SourceId,omitempty"`
				Name        string `xml:"Name,omitempty"`
				Location    string `xml:"Location,omitempty"`
				Description string `xml:"Description,omitempty"`
				Address     string `xml:"Address,omitempty"`
			} `xml:"Source"`
			Content              string `xml:"Content"`
			MaximumRetentionTime string `xml:"MaximumRetentionTime"`
		} `xml:"trc:RecordingConfiguration"`
	}

	type createRecordingResponse struct {
		RecordingToken string `xml:"RecordingToken"`
	}

	req := createRecordingRequest{Xmlns: recordingNamespace}
	req.RecordingConfiguration.Source.SourceId = config.Source.SourceId
	req.RecordingConfiguration.Source.Name = config.Source.Name
	req.RecordingConfiguration.Source.Location = config.Source.Location
	req.RecordingConfiguration.Source.Description = config.Source.Description
	req.RecordingConfiguration.Source.Address = config.Source.Address
	req.RecordingConfiguration.Content = config.Content
	req.RecordingConfiguration.MaximumRetentionTime = config.MaximumRetentionTime

	var resp createRecordingResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getRecordingEndpoint(), "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateRecording failed: %w", err)
	}

	return resp.RecordingToken, nil
}

// DeleteRecording deletes a recording.
func (c *Client) DeleteRecording(ctx context.Context, recordingToken string) error {
	type deleteRecordingRequest struct {
		XMLName        xml.Name `xml:"trc:DeleteRecording"`
		Xmlns          string   `xml:"xmlns:trc,attr"`
		RecordingToken string   `xml:"trc:RecordingToken"`
	}

	type deleteRecordingResponse struct{}

	req := deleteRecordingRequest{Xmlns: recordingNamespace, RecordingToken: recordingToken}
	var resp deleteRecordingResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getRecordingEndpoint(), "", req, &resp); err != nil {
		return fmt.Errorf("DeleteRecording failed: %w", err)
	}

	return nil
}

// SetRecordingConfiguration sets the configuration for a recording.
func (c *Client) SetRecordingConfiguration(ctx context.Context, recordingToken string, config *RecordingConfiguration) error {
	type setRecordingConfigurationRequest struct {
		XMLName                xml.Name `xml:"trc:SetRecordingConfiguration"`
		Xmlns                  string   `xml:"xmlns:trc,attr"`
		RecordingToken         string   `xml:"trc:RecordingToken"`
		RecordingConfiguration struct {
			Source struct {
				SourceId    string `xml:"SourceId,omitempty"`
				Name        string `xml:"Name,omitempty"`
				Location    string `xml:"Location,omitempty"`
				Description string `xml:"Description,omitempty"`
				Address     string `xml:"Address,omitempty"`
			} `xml:"Source"`
			Content              string `xml:"Content"`
			MaximumRetentionTime string `xml:"MaximumRetentionTime"`
		} `xml:"trc:RecordingConfiguration"`
	}

	type setRecordingConfigurationResponse struct{}

	req := setRecordingConfigurationRequest{Xmlns: recordingNamespace, RecordingToken: recordingToken}
	req.RecordingConfiguration.Source.SourceId = config.Source.SourceId
	req.RecordingConfiguration.Source.Name = config.Source.Name
	req.RecordingConfiguration.Content = config.Content
	req.RecordingConfiguration.MaximumRetentionTime = config.MaximumRetentionTime

	var resp setRecordingConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getRecordingEndpoint(), "", req, &resp); err != nil {
		return fmt.Errorf("SetRecordingConfiguration failed: %w", err)
	}

	return nil
}

// GetRecordingConfiguration returns the configuration of a recording.
func (c *Client) GetRecordingConfiguration(ctx context.Context, recordingToken string) (*RecordingConfiguration, error) {
	type getRecordingConfigurationRequest struct {
		XMLName        xml.Name `xml:"trc:GetRecordingConfiguration"`
		Xmlns          string   `xml:"xmlns:trc,attr"`
		RecordingToken string   `xml:"trc:RecordingToken"`
	}

	type getRecordingConfigurationResponse struct {
		RecordingConfiguration struct {
			Source struct {
				SourceId    string `xml:"SourceId"`
				Name        string `xml:"Name"`
				Location    string `xml:"Location"`
				Description string `xml:"Description"`
				Address     string `xml:"Address"`
			} `xml:"Source"`
			Content              string `xml:"Content"`
			MaximumRetentionTime string `xml:"MaximumRetentionTime"`
		} `xml:"RecordingConfiguration"`
	}

	req := getRecordingConfigurationRequest{Xmlns: recordingNamespace, RecordingToken: recordingToken}
	var resp getRecordingConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getRecordingEndpoint(), "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRecordingConfiguration failed: %w", err)
	}

	return &RecordingConfiguration{
		Source: RecordingSourceInformation{
			SourceId:    resp.RecordingConfiguration.Source.SourceId,
			Name:        resp.RecordingConfiguration.Source.Name,
			Location:    resp.RecordingConfiguration.Source.Location,
			Description: resp.RecordingConfiguration.Source.Description,
			Address:     resp.RecordingConfiguration.Source.Address,
		},
		Content:              resp.RecordingConfiguration.Content,
		MaximumRetentionTime: resp.RecordingConfiguration.MaximumRetentionTime,
	}, nil
}

// GetRecordingOptions returns available options for a recording.
func (c *Client) GetRecordingOptions(ctx context.Context, recordingToken string) (*RecordingOptions, error) {
	type getRecordingOptionsRequest struct {
		XMLName        xml.Name `xml:"trc:GetRecordingOptions"`
		Xmlns          string   `xml:"xmlns:trc,attr"`
		RecordingToken string   `xml:"trc:RecordingToken"`
	}

	type getRecordingOptionsResponse struct {
		Options struct {
			Track *struct {
				SpareTotal    *int `xml:"SpareTotal"`
				SpareVideo    *int `xml:"SpareVideo"`
				SpareAudio    *int `xml:"SpareAudio"`
				SpareMetadata *int `xml:"SpareMetadata"`
			} `xml:"Track"`
		} `xml:"Options"`
	}

	req := getRecordingOptionsRequest{Xmlns: recordingNamespace, RecordingToken: recordingToken}
	var resp getRecordingOptionsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getRecordingEndpoint(), "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRecordingOptions failed: %w", err)
	}

	opts := &RecordingOptions{}
	if resp.Options.Track != nil {
		opts.Track = &RecordingTrackOptions{
			SpareTotal:    resp.Options.Track.SpareTotal,
			SpareVideo:    resp.Options.Track.SpareVideo,
			SpareAudio:    resp.Options.Track.SpareAudio,
			SpareMetadata: resp.Options.Track.SpareMetadata,
		}
	}

	return opts, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go test -run "TestGetRecordings|TestCreateRecording|TestDeleteRecording|TestGetRecordingServiceCapabilities" -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add recording.go recording_test.go types_recording.go
git commit -m "feat: add Recording service core operations (GetRecordings, CreateRecording, DeleteRecording, capabilities, configuration)"
```

---

## Task 3: Recording Service — Track Operations (4 ops)

**Files:**
- Modify: `recording.go`
- Modify: `recording_test.go`

Adds: `CreateTrack`, `DeleteTrack`, `GetTrackConfiguration`, `SetTrackConfiguration`.

All take `recordingToken` + `trackToken` parameters. Follow the identical pattern from Task 2.

- [ ] **Step 1: Write tests for all 4 track operations**

Follow the TestCreateRecording/TestDeleteRecording pattern. Each takes a `recordingToken` param.

- [ ] **Step 2: Implement all 4 methods**

Signatures:
```go
func (c *Client) CreateTrack(ctx context.Context, recordingToken string, config *TrackConfiguration) (string, error)
func (c *Client) DeleteTrack(ctx context.Context, recordingToken, trackToken string) error
func (c *Client) GetTrackConfiguration(ctx context.Context, recordingToken, trackToken string) (*TrackConfiguration, error)
func (c *Client) SetTrackConfiguration(ctx context.Context, recordingToken, trackToken string, config *TrackConfiguration) error
```

- [ ] **Step 3: Run tests, verify pass**

- [ ] **Step 4: Commit**

```bash
git add recording.go recording_test.go
git commit -m "feat: add Recording track operations (CreateTrack, DeleteTrack, GetTrackConfiguration, SetTrackConfiguration)"
```

---

## Task 4: Recording Service — Job Operations (7 ops)

**Files:**
- Modify: `recording.go`
- Modify: `recording_test.go`

Adds: `CreateRecordingJob`, `DeleteRecordingJob`, `GetRecordingJobs`, `SetRecordingJobConfiguration`, `GetRecordingJobConfiguration`, `SetRecordingJobMode`, `GetRecordingJobState`.

- [ ] **Step 1: Write tests for all 7 job operations**
- [ ] **Step 2: Implement all 7 methods**

Signatures:
```go
func (c *Client) CreateRecordingJob(ctx context.Context, config *RecordingJobConfiguration) (string, *RecordingJobConfiguration, error)
func (c *Client) DeleteRecordingJob(ctx context.Context, jobToken string) error
func (c *Client) GetRecordingJobs(ctx context.Context) ([]*RecordingJob, error)
func (c *Client) SetRecordingJobConfiguration(ctx context.Context, jobToken string, config *RecordingJobConfiguration) (*RecordingJobConfiguration, error)
func (c *Client) GetRecordingJobConfiguration(ctx context.Context, jobToken string) (*RecordingJobConfiguration, error)
func (c *Client) SetRecordingJobMode(ctx context.Context, jobToken, mode string) error
func (c *Client) GetRecordingJobState(ctx context.Context, jobToken string) (*RecordingJobState, error)
```

- [ ] **Step 3: Run tests, verify pass**
- [ ] **Step 4: Commit**

```bash
git add recording.go recording_test.go
git commit -m "feat: add Recording job operations (Create, Delete, Get, Set configuration and mode)"
```

---

## Task 5: Recording Service — Export Operations (4 ops)

**Files:**
- Modify: `recording.go`
- Modify: `recording_test.go`

Adds: `ExportRecordedData`, `StopExportRecordedData`, `GetExportRecordedDataState`, `OverrideSegmentDuration`.

- [ ] **Step 1: Write tests and implement**

Signatures:
```go
func (c *Client) ExportRecordedData(ctx context.Context, startPoint, endPoint *time.Time, searchScope *SearchScope, fileFormat string, storageDest string) (string, []string, error)
func (c *Client) StopExportRecordedData(ctx context.Context, operationToken string) (*ExportRecordedDataState, error)
func (c *Client) GetExportRecordedDataState(ctx context.Context, operationToken string) (*ExportRecordedDataState, error)
func (c *Client) OverrideSegmentDuration(ctx context.Context, targetDuration, expiration string, recordingToken string) error
```

- [ ] **Step 2: Run tests, verify pass**
- [ ] **Step 3: Commit**

```bash
git add recording.go recording_test.go
git commit -m "feat: add Recording export operations (ExportRecordedData, StopExport, GetExportState, OverrideSegmentDuration)"
```

---

## Task 6: Receiver Service (8 ops)

**Files:**
- Create: `types_receiver.go`
- Create: `receiver.go`
- Create: `receiver_test.go`

- [ ] **Step 1: Create types_receiver.go**

```go
package onvif

// Receiver represents an ONVIF media receiver.
type Receiver struct {
	Token         string
	Configuration ReceiverConfiguration
}

// ReceiverConfiguration contains receiver settings.
type ReceiverConfiguration struct {
	Mode         string
	MediaURI     string
	StreamSetup  *StreamSetup
}

// StreamSetup defines stream transport settings.
type StreamSetup struct {
	Stream    string
	Transport *Transport
}

// Transport defines the transport protocol.
type Transport struct {
	Protocol string
}

// ReceiverState contains receiver state information.
type ReceiverState struct {
	State      string
	AutoCreated bool
}

// ReceiverServiceCapabilities represents receiver service capabilities.
type ReceiverServiceCapabilities struct {
	RTP_Multicast            bool
	RTP_TCP                  bool
	RTP_RTSP_TCP             bool
	SupportedReceivers       int
	MaximumRTSPURILength     int
}
```

- [ ] **Step 2: Create receiver.go with all 8 operations**

```go
package onvif

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/0x524a/onvif-go/internal/soap"
)

const receiverNamespace = "http://www.onvif.org/ver10/receiver/wsdl"

func (c *Client) getReceiverEndpoint() string {
	if c.receiverEndpoint != "" {
		return c.receiverEndpoint
	}

	return c.endpoint
}
```

Then implement all 8 operations following the established pattern. Signatures:

```go
func (c *Client) GetReceiverServiceCapabilities(ctx context.Context) (*ReceiverServiceCapabilities, error)
func (c *Client) GetReceivers(ctx context.Context) ([]*Receiver, error)
func (c *Client) GetReceiver(ctx context.Context, receiverToken string) (*Receiver, error)
func (c *Client) CreateReceiver(ctx context.Context, config *ReceiverConfiguration) (*Receiver, error)
func (c *Client) DeleteReceiver(ctx context.Context, receiverToken string) error
func (c *Client) ConfigureReceiver(ctx context.Context, receiverToken string, config *ReceiverConfiguration) error
func (c *Client) SetReceiverMode(ctx context.Context, receiverToken, mode string) error
func (c *Client) GetReceiverState(ctx context.Context, receiverToken string) (*ReceiverState, error)
```

- [ ] **Step 3: Write tests for all 8 operations**
- [ ] **Step 4: Run tests, verify pass**
- [ ] **Step 5: Commit**

```bash
git add types_receiver.go receiver.go receiver_test.go
git commit -m "feat: add Receiver service (GetReceivers, CreateReceiver, ConfigureReceiver, state management)"
```

---

## Task 7: Search Service (14 ops)

**Files:**
- Create: `types_search.go`
- Create: `search.go`
- Create: `search_test.go`

- [ ] **Step 1: Create types_search.go**

```go
package onvif

import "time"

// SearchScope defines the scope for recording searches.
type SearchScope struct {
	IncludedSources   []*SourceReference
	IncludedRecordings []string
	RecordingInformationFilter string
	Extension          interface{}
}

// RecordingSummary contains a summary of available recordings.
type RecordingSummary struct {
	DataFrom     time.Time
	DataUntil    time.Time
	NumberRecordings int
}

// RecordingInformation contains detailed information about a recording.
type RecordingInformation struct {
	RecordingToken       string
	Source               RecordingSourceInformation
	EarliestRecording    *time.Time
	LatestRecording      *time.Time
	Content              string
	TrackInformation     []*TrackInformation
	RecordingStatus      string
}

// TrackInformation contains information about a track.
type TrackInformation struct {
	TrackToken  string
	TrackType   string
	Description string
	DataFrom    time.Time
	DataTo      time.Time
}

// FindRecordingResult contains a recording search result.
type FindRecordingResult struct {
	SearchState string
	RecordingInformation []*RecordingInformation
}

// FindEventResult contains an event search result.
type FindEventResult struct {
	SearchState string
	Events      []*EventResult
}

// EventResult represents an event found during search.
type EventResult struct {
	RecordingToken string
	TrackToken     string
	Time           time.Time
	Event          interface{}
	StartStateEvent bool
}

// FindPTZPositionResult contains a PTZ position search result.
type FindPTZPositionResult struct {
	SearchState string
	Positions   []*PTZPositionResult
}

// PTZPositionResult represents a PTZ position found during search.
type PTZPositionResult struct {
	RecordingToken string
	TrackToken     string
	Time           time.Time
	Position       *PTZVector
}

// FindMetadataResult contains a metadata search result.
type FindMetadataResult struct {
	SearchState string
	Results     []*MetadataResult
}

// MetadataResult represents metadata found during search.
type MetadataResult struct {
	RecordingToken string
	TrackToken     string
	Time           time.Time
}

// MediaAttributes contains media attributes for recordings.
type MediaAttributes struct {
	RecordingToken    string
	TrackAttributes   []*TrackAttributes
	From              *time.Time
	Until             *time.Time
}

// TrackAttributes contains media attributes for a track.
type TrackAttributes struct {
	TrackInformation *TrackInformation
	VideoAttributes  *VideoAttributes
	AudioAttributes  *AudioAttributes
}

// VideoAttributes contains video-specific attributes.
type VideoAttributes struct {
	Bitrate    *int
	Width      int
	Height     int
	Encoding   string
	Framerate  *float64
}

// AudioAttributes contains audio-specific attributes.
type AudioAttributes struct {
	Bitrate    *int
	Encoding   string
	Samplerate int
}

// SearchServiceCapabilities represents search service capabilities.
type SearchServiceCapabilities struct {
	MetadataSearch bool
}
```

- [ ] **Step 2: Create search.go with all 14 operations**

Namespace: `http://www.onvif.org/ver10/search/wsdl`, prefix: `tse:`.

Signatures:
```go
func (c *Client) GetSearchServiceCapabilities(ctx context.Context) (*SearchServiceCapabilities, error)
func (c *Client) GetRecordingSummary(ctx context.Context) (*RecordingSummary, error)
func (c *Client) GetRecordingInformation(ctx context.Context, recordingToken string) (*RecordingInformation, error)
func (c *Client) GetMediaAttributes(ctx context.Context, recordingTokens []string, time time.Time) ([]*MediaAttributes, error)
func (c *Client) FindRecordings(ctx context.Context, scope *SearchScope, maxMatches *int, keepAliveTime string) (string, error)
func (c *Client) GetRecordingSearchResults(ctx context.Context, searchToken string, minResults, maxResults *int, waitTime string) (*FindRecordingResult, error)
func (c *Client) FindEvents(ctx context.Context, startPoint time.Time, endPoint *time.Time, scope *SearchScope, searchFilter string, includeStartState bool, maxMatches *int, keepAliveTime string) (string, error)
func (c *Client) GetEventSearchResults(ctx context.Context, searchToken string, minResults, maxResults *int, waitTime string) (*FindEventResult, error)
func (c *Client) FindPTZPosition(ctx context.Context, startPoint time.Time, endPoint *time.Time, scope *SearchScope, maxMatches *int, keepAliveTime string) (string, error)
func (c *Client) GetPTZPositionSearchResults(ctx context.Context, searchToken string, minResults, maxResults *int, waitTime string) (*FindPTZPositionResult, error)
func (c *Client) FindMetadata(ctx context.Context, startPoint time.Time, endPoint *time.Time, scope *SearchScope, maxMatches *int, keepAliveTime string) (string, error)
func (c *Client) GetMetadataSearchResults(ctx context.Context, searchToken string, minResults, maxResults *int, waitTime string) (*FindMetadataResult, error)
func (c *Client) GetSearchState(ctx context.Context, searchToken string) (string, error)
func (c *Client) EndSearch(ctx context.Context, searchToken string) (*time.Time, error)
```

- [ ] **Step 3: Write tests for all 14 operations**
- [ ] **Step 4: Run tests, verify pass**
- [ ] **Step 5: Commit**

```bash
git add types_search.go search.go search_test.go
git commit -m "feat: add Search service (FindRecordings, FindEvents, FindPTZPosition, FindMetadata, search results, media attributes)"
```

---

## Task 8: Replay Service (4 ops)

**Files:**
- Create: `types_replay.go`
- Create: `replay.go`
- Create: `replay_test.go`

- [ ] **Step 1: Create types and implementation**

```go
// types_replay.go
package onvif

// ReplayConfiguration contains replay settings.
type ReplayConfiguration struct {
	SessionTimeout string
}

// ReplayServiceCapabilities represents replay service capabilities.
type ReplayServiceCapabilities struct {
	ReversePlayback bool
	SessionTimeoutRange *struct {
		Min string
		Max string
	}
	RTP_RTSP_TCP bool
}
```

```go
// replay.go
package onvif

const replayNamespace = "http://www.onvif.org/ver10/replay/wsdl"
```

Signatures:
```go
func (c *Client) GetReplayServiceCapabilities(ctx context.Context) (*ReplayServiceCapabilities, error)
func (c *Client) GetReplayUri(ctx context.Context, streamSetup *StreamSetup, recordingToken string) (string, error)
func (c *Client) GetReplayConfiguration(ctx context.Context) (*ReplayConfiguration, error)
func (c *Client) SetReplayConfiguration(ctx context.Context, config *ReplayConfiguration) error
```

- [ ] **Step 2: Write tests for all 4 operations**
- [ ] **Step 3: Run tests, verify pass**
- [ ] **Step 4: Commit**

```bash
git add types_replay.go replay.go replay_test.go
git commit -m "feat: add Replay service (GetReplayUri, GetReplayConfiguration, SetReplayConfiguration)"
```

---

## Task 9: Integration Tests for All Phase 3 Services

**Files:**
- Create: `recording_real_camera_test.go`
- Create: `search_real_camera_test.go`
- Create: `replay_real_camera_test.go`
- Create: `receiver_real_camera_test.go`

- [ ] **Step 1: Create integration test files**

Each follows the pattern from Phase 2 Task 11. Example for recording:

```go
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
		t.Logf("Recording caps: DynamicRecordings=%v, MaxRecordings=%d", caps.DynamicRecordings, caps.MaxRecordings)
	})

	t.Run("GetRecordings", func(t *testing.T) {
		recordings, err := client.GetRecordings(ctx)
		if err != nil {
			t.Skipf("GetRecordings not supported: %v", err)
		}
		t.Logf("Found %d recordings", len(recordings))
		for _, rec := range recordings {
			t.Logf("  Recording: %s (%s)", rec.Token, rec.Configuration.Source.Name)
		}
	})

	t.Run("GetRecordingJobs", func(t *testing.T) {
		jobs, err := client.GetRecordingJobs(ctx)
		if err != nil {
			t.Skipf("GetRecordingJobs not supported: %v", err)
		}
		t.Logf("Found %d recording jobs", len(jobs))
	})
}
```

- [ ] **Step 2: Verify compilation with build tag**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go vet -tags=real_camera ./...`
Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add recording_real_camera_test.go search_real_camera_test.go replay_real_camera_test.go receiver_real_camera_test.go
git commit -m "feat: add integration test scaffolding for Recording, Search, Replay, Receiver services"
```

---

## Task 10: Full Test Suite Verification

- [ ] **Step 1: Run all unit tests**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go test ./... -count=1`
Expected: All PASS.

- [ ] **Step 2: Run linter**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && make check`
Expected: Clean.

- [ ] **Step 3: Fix any lint issues and commit**

```bash
git add -A
git commit -m "fix: resolve lint issues from Phase 3 implementation"
```
