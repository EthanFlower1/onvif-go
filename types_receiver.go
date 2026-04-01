package onvif

// Receiver represents an ONVIF media receiver.
type Receiver struct {
	Token         string
	Configuration ReceiverConfiguration
}

// ReceiverConfiguration contains receiver settings.
type ReceiverConfiguration struct {
	Mode        string
	MediaURI    string
	StreamSetup *StreamSetup
}

// ReceiverStateInformation contains receiver state information.
type ReceiverStateInformation struct {
	State       string
	AutoCreated bool
}

// ReceiverServiceCapabilities represents receiver service capabilities.
type ReceiverServiceCapabilities struct {
	RTPMulticast         bool
	RTPTCP               bool
	RTPRTSP_TCP          bool
	SupportedReceivers   int
	MaximumRTSPURILength int
}
