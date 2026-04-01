package onvif

// UplinkProtocol represents the protocol used for an uplink connection.
type UplinkProtocol string

const (
	// UplinkProtocolHTTPS is the native h2c-reverse protocol.
	UplinkProtocolHTTPS UplinkProtocol = "https"
	// UplinkProtocolWSS is the h2c-reverse over WebSocket protocol.
	UplinkProtocolWSS UplinkProtocol = "wss"
)

// UplinkAuthorizationMode represents the authorization mode for an uplink connection.
type UplinkAuthorizationMode string

const (
	// UplinkAuthorizationModeMTLS uses TLS with a client certificate.
	UplinkAuthorizationModeMTLS UplinkAuthorizationMode = "mTLS"
	// UplinkAuthorizationModeAccessToken uses an access token obtained from an authorization server.
	UplinkAuthorizationModeAccessToken UplinkAuthorizationMode = "AccessToken"
)

// UplinkConnectionStatus represents the current connection status of an uplink.
type UplinkConnectionStatus string

const (
	// UplinkConnectionStatusOffline indicates the uplink is not connected.
	UplinkConnectionStatusOffline UplinkConnectionStatus = "Offline"
	// UplinkConnectionStatusConnecting indicates the uplink is establishing a connection.
	UplinkConnectionStatusConnecting UplinkConnectionStatus = "Connecting"
	// UplinkConnectionStatusConnected indicates the uplink is connected.
	UplinkConnectionStatusConnected UplinkConnectionStatus = "Connected"
)

// UplinkServiceCapabilities contains the capabilities of the Uplink service.
type UplinkServiceCapabilities struct {
	// MaxUplinks is the maximum number of uplink connections that can be configured.
	MaxUplinks *int
	// Protocols lists the protocols supported by the device (e.g. "https wss").
	Protocols string
	// AuthorizationModes lists the supported authorization modes (e.g. "mTLS AccessToken").
	AuthorizationModes string
	// StreamingOverUplink signals support for media streaming over uplink.
	StreamingOverUplink *bool
}

// UplinkConfiguration holds the configuration for a single uplink connection.
type UplinkConfiguration struct {
	// RemoteAddress is the URI of the remote uplink server.
	RemoteAddress string
	// CertificateID is the optional ID of the certificate used for client authentication.
	CertificateID *string
	// UserLevel lists the authorization levels/roles restricting commands accepted over the uplink.
	UserLevel string
	// Status is the current connection status (readonly).
	Status *string
	// CertPathValidationPolicyID is the optional policy ID used to validate the server certificate.
	CertPathValidationPolicyID *string
	// AuthorizationServer is the optional token referring to the authorization server.
	AuthorizationServer *string
	// Error contains optional user-readable error information (readonly).
	Error *string
}
