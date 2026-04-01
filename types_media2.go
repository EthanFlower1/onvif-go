package onvif

// Media2Profile represents a media profile in the Media2 service.
type Media2Profile struct {
	Token          string
	Name           string
	Fixed          bool
	Configurations *Media2Configurations
}

// Media2Configurations contains the configurations within a Media2 profile.
type Media2Configurations struct {
	VideoSource  *VideoSourceConfiguration
	AudioSource  *AudioSourceConfiguration
	VideoEncoder *VideoEncoderConfiguration
	AudioEncoder *AudioEncoderConfiguration
	Analytics    *VideoAnalyticsConfiguration
	Metadata     *MetadataConfiguration
	AudioOutput  *AudioOutputConfiguration
	AudioDecoder *AudioDecoderConfiguration
	PTZ          *PTZConfiguration
}

// Mask represents a privacy mask.
type Mask struct {
	Token              string
	ConfigurationToken string
	Polygon            *Polygon
	Type               string
	Enabled            bool
}

// Polygon represents a polygon shape.
type Polygon struct {
	Points []*Vector
}

// Vector represents a 2D point.
type Vector struct {
	X float64
	Y float64
}

// MaskOptions represents available mask configuration options.
type MaskOptions struct {
	MaxMasks        int
	MaxPoints       int
	Types           []string
	SingleColorOnly bool
}

// AudioClip represents an audio clip.
type AudioClip struct {
	Token    string
	Name     string
	MediaURI string
}

// WebRTCConfiguration represents WebRTC streaming configuration.
type WebRTCConfiguration struct {
	SignalingServerURI string
	STUNServer        string
	TURNServer        string
}

// EQPreset represents an equalizer preset.
type EQPreset struct {
	Token string
	Name  string
}

// MulticastAudioDecoderConfiguration represents multicast audio decoder settings.
type MulticastAudioDecoderConfiguration struct {
	Token          string
	Name           string
	UseCount       int
	Multicast      *MulticastConfiguration
	SessionTimeout string
}

// Media2ServiceCapabilities represents Media2 service capabilities.
type Media2ServiceCapabilities struct {
	SnapshotUri     bool
	Rotation        bool
	VideoSourceMode bool
	OSD             bool
	Mask            bool
	SourceMask      bool
}

// Media2ConfigurationRef identifies a configuration to add/remove from a profile.
type Media2ConfigurationRef struct {
	Type  string
	Token string
}

// VideoEncoderInstances contains encoder instance information.
type VideoEncoderInstances struct {
	Total int
	JPEG  *int
	H264  *int
	MPEG4 *int
}
