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
	SourceToken        *SourceReference
	AutoCreateReceiver bool
	Tracks             []*RecordingJobTrack
}

// SourceReference identifies a source by token and type.
type SourceReference struct {
	Token string
	Type  string
}

// RecordingJobTrack maps a source track to a recording track.
type RecordingJobTrack struct {
	SourceTag   string
	Destination string
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
	Spare             *int
	CompatibleSources []string
}

// RecordingTrackOptions contains track-related options.
type RecordingTrackOptions struct {
	SpareTotal    *int
	SpareVideo    *int
	SpareAudio    *int
	SpareMetadata *int
}

// RecordingServiceCapabilities represents recording service capabilities.
type RecordingServiceCapabilities struct {
	DynamicRecordings          bool
	DynamicTracks              bool
	MaxStringLength            int
	MaxRecordings              int
	MaxRecordingJobs           int
	Options                    bool
	MetadataRecording          bool
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
