package onvif

// PanDirection is the direction for PanMove to move the device.
type PanDirection string

const (
	// PanDirectionLeft moves left in relation to the video source image.
	PanDirectionLeft PanDirection = "Left"
	// PanDirectionRight moves right in relation to the video source image.
	PanDirectionRight PanDirection = "Right"
)

// TiltDirection is the direction for TiltMove to move the device.
type TiltDirection string

const (
	// TiltDirectionUp moves up in relation to the video source image.
	TiltDirectionUp TiltDirection = "Up"
	// TiltDirectionDown moves down in relation to the video source image.
	TiltDirectionDown TiltDirection = "Down"
)

// ZoomDirection is the direction for ZoomMove to change the focal length.
type ZoomDirection string

const (
	// ZoomDirectionWide moves video source lens toward a wider field of view.
	ZoomDirectionWide ZoomDirection = "Wide"
	// ZoomDirectionTelephoto moves video source lens toward a narrower field of view.
	ZoomDirectionTelephoto ZoomDirection = "Telephoto"
)

// RollDirection is the direction for RollMove to move the device.
type RollDirection string

const (
	// RollDirectionClockwise moves clockwise in relation to the video source image.
	RollDirectionClockwise RollDirection = "Clockwise"
	// RollDirectionCounterclockwise moves counterclockwise in relation to the video source image.
	RollDirectionCounterclockwise RollDirection = "Counterclockwise"
	// RollDirectionAuto automatically levels the device.
	RollDirectionAuto RollDirection = "Auto"
)

// FocusDirection is the direction for FocusMove to move the focal plane.
type FocusDirection string

const (
	// FocusDirectionNear moves to focus on close objects.
	FocusDirectionNear FocusDirection = "Near"
	// FocusDirectionFar moves to focus on distant objects.
	FocusDirectionFar FocusDirection = "Far"
	// FocusDirectionAuto automatically focuses for the sharpest video source image.
	FocusDirectionAuto FocusDirection = "Auto"
)

// ProvisioningUsage contains the quantity of movement events over the lifetime of the device.
type ProvisioningUsage struct {
	Pan   *int
	Tilt  *int
	Zoom  *int
	Roll  *int
	Focus *int
}

// ProvisioningSourceCapabilities contains the provisioning capabilities of a video source.
type ProvisioningSourceCapabilities struct {
	VideoSourceToken  string
	MaximumPanMoves   *int
	MaximumTiltMoves  *int
	MaximumZoomMoves  *int
	MaximumRollMoves  *int
	AutoLevel         *bool
	MaximumFocusMoves *int
	AutoFocus         *bool
}

// ProvisioningServiceCapabilities contains the capabilities of the Provisioning Service.
type ProvisioningServiceCapabilities struct {
	DefaultTimeout string
	Source         []*ProvisioningSourceCapabilities
}
