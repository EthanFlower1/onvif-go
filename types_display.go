package onvif

// DisplayServiceCapabilities represents the capabilities of the Display service.
type DisplayServiceCapabilities struct {
	// FixedLayout indicates that the SetLayout command supports only predefined layouts.
	FixedLayout *bool
}

// PaneConfiguration describes a pane configuration for a video output.
type PaneConfiguration struct {
	// Token is the unique identifier for this pane configuration.
	Token string
	// PaneName is the descriptive name of the pane.
	PaneName string
	// AudioOutputToken refers to the audio output associated with this pane.
	AudioOutputToken *string
	// AudioSourceToken refers to the audio source associated with this pane.
	AudioSourceToken *string
	// ReceiverToken refers to the receiver providing media for this pane.
	ReceiverToken *string
	// MediaUri is the URI of the media stream displayed in this pane.
	MediaUri *string
	// Profile is the media profile token for the stream displayed in this pane.
	Profile *string
}

// LayoutOption describes a predefined layout option for a video output.
type LayoutOption struct {
	// PaneLayoutOptions lists the supported pane layouts.
	PaneLayoutOptions []PaneLayoutOption
}

// PaneLayoutOption describes valid pane area values for a given layout.
type PaneLayoutOption struct {
	// Area is the valid region for a pane in this layout.
	Area []FloatRectangle
}

// LayoutOptions describe the fixed and predefined layouts of a device.
type LayoutOptions struct {
	// FixedLayout lists predefined layouts available on the device.
	FixedLayout []LayoutOption
}

// CodingCapabilities describes the decoding and encoding capabilities of a video output.
type CodingCapabilities struct {
	// InputTokensLimits specifies the maximum number of connected receivers.
	InputTokensLimits *CodingCapabilityLimits
	// OutputTokensLimits specifies the maximum number of output streams.
	OutputTokensLimits *CodingCapabilityLimits
}

// CodingCapabilityLimits describes token limits for coding capabilities.
type CodingCapabilityLimits struct {
	// Max is the maximum number of tokens.
	Max int
}

// DisplayOptions contains the layout options and coding capabilities for a video output.
type DisplayOptions struct {
	// LayoutOptions describes the supported layouts.
	LayoutOptions *LayoutOptions
	// CodingCapabilities describes the decoding and encoding capabilities.
	CodingCapabilities CodingCapabilities
}
