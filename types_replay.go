package onvif

// ReplayConfiguration contains replay settings.
type ReplayConfiguration struct {
	SessionTimeout string
}

// ReplayServiceCapabilities represents replay service capabilities.
type ReplayServiceCapabilities struct {
	ReversePlayback     bool
	SessionTimeoutRange *DurationRange
	RTPRTSP_TCP         bool
}
