package onvif

// AccessControlServiceCapabilities represents the capabilities of the Access Control service.
type AccessControlServiceCapabilities struct {
	MaxLimit                        uint
	MaxAccessPoints                 uint
	MaxAreas                        uint
	ClientSuppliedTokenSupported    bool
	AccessPointManagementSupported  bool
	AreaManagementSupported         bool
}

// AccessPointCapabilities represents capabilities of a specific access point.
type AccessPointCapabilities struct {
	DisableAccessPoint      bool
	Duress                  *bool
	AnonymousAccess         *bool
	AccessTaken             *bool
	ExternalAuthorization   *bool
	IdentifierAccess        *bool
	SupportedSecurityLevels []string
	SupportedRecognitionTypes []string
	SupportedFeedbackTypes  []string
}

// AccessPointInfo contains basic information about an access point instance.
type AccessPointInfo struct {
	Token        string
	Name         string
	Description  string
	AreaFrom     string
	AreaTo       string
	EntityType   string
	Entity       string
	Capabilities AccessPointCapabilities
}

// AccessPoint includes all properties of AccessPointInfo plus authentication profile.
type AccessPoint struct {
	AccessPointInfo
	AuthenticationProfileToken string
}

// AreaInfo contains basic information about an area.
type AreaInfo struct {
	Token       string
	Name        string
	Description string
}

// Area includes all properties of AreaInfo.
type Area struct {
	AreaInfo
}

// AccessPointState contains state information for an access point.
type AccessPointState struct {
	Enabled bool
}

// AccessControlDecision represents the access decision.
type AccessControlDecision string

// Access control decision constants.
const (
	AccessDecisionGranted AccessControlDecision = "Granted"
	AccessDecisionDenied  AccessControlDecision = "Denied"
)

// AccessDenyReason represents the reason for denying access.
type AccessDenyReason string

// Deny reason constants.
const (
	DenyReasonCredentialNotEnabled    AccessDenyReason = "CredentialNotEnabled"
	DenyReasonCredentialNotActive     AccessDenyReason = "CredentialNotActive"
	DenyReasonCredentialExpired       AccessDenyReason = "CredentialExpired"
	DenyReasonInvalidPIN              AccessDenyReason = "InvalidPIN"
	DenyReasonNotPermittedAtThisTime  AccessDenyReason = "NotPermittedAtThisTime"
	DenyReasonUnauthorized            AccessDenyReason = "Unauthorized"
	DenyReasonOther                   AccessDenyReason = "Other"
)
