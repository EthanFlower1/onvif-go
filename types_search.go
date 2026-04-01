package onvif

import "time"

// SearchScope defines the scope for recording searches.
type SearchScope struct {
	IncludedSources            []*SourceReference
	IncludedRecordings         []string
	RecordingInformationFilter string
}

// RecordingSummary contains a summary of available recordings.
type RecordingSummary struct {
	DataFrom         time.Time
	DataUntil        time.Time
	NumberRecordings int
}

// RecordingInformation contains detailed information about a recording.
type RecordingInformation struct {
	RecordingToken    string
	Source            RecordingSourceInformation
	EarliestRecording *time.Time
	LatestRecording   *time.Time
	Content           string
	TrackInformation  []*TrackInformation
	RecordingStatus   string
}

// TrackInformation contains information about a track.
type TrackInformation struct {
	TrackToken  string
	TrackType   string
	Description string
	DataFrom    time.Time
	DataTo      time.Time
}

// FindRecordingResult contains recording search results.
type FindRecordingResult struct {
	SearchState          string
	RecordingInformation []*RecordingInformation
}

// FindEventResult contains event search results.
type FindEventResult struct {
	SearchState string
	Events      []*EventResult
}

// EventResult represents an event found during search.
type EventResult struct {
	RecordingToken  string
	TrackToken      string
	Time            time.Time
	StartStateEvent bool
}

// FindPTZPositionResult contains PTZ position search results.
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

// FindMetadataResult contains metadata search results.
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
	RecordingToken  string
	TrackAttributes []*TrackAttributes
	From            *time.Time
	Until           *time.Time
}

// TrackAttributes contains media attributes for a track.
type TrackAttributes struct {
	TrackInformation *TrackInformation
	VideoAttributes  *VideoAttributes
	AudioAttributes  *AudioAttributes
}

// VideoAttributes contains video-specific attributes.
type VideoAttributes struct {
	Bitrate   *int
	Width     int
	Height    int
	Encoding  string
	Framerate *float64
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
