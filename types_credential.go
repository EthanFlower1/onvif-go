package onvif

import "time"

// CredentialServiceCapabilities represents the capabilities of the Credential service.
type CredentialServiceCapabilities struct {
	MaxLimit                                uint
	MaxCredentials                          uint
	MaxAccessProfilesPerCredential          uint
	CredentialValiditySupported             bool
	CredentialAccessProfileValiditySupported bool
	ValiditySupportsTimeValue               bool
	ResetAntipassbackSupported              bool
	ClientSuppliedTokenSupported            bool
	DefaultCredentialSuspensionDuration     string
	MaxWhitelistedItems                     uint
	MaxBlacklistedItems                     uint
	SupportedIdentifierTypes                []string
	SupportedExemptionTypes                 []string
}

// CredentialInfo contains basic information about a credential instance.
type CredentialInfo struct {
	Token                    string
	Description              string
	CredentialHolderReference string
	ValidFrom                *time.Time
	ValidTo                  *time.Time
}

// Credential includes all properties of CredentialInfo plus identifiers, access profiles, and attributes.
type Credential struct {
	CredentialInfo
	CredentialIdentifiers    []CredentialIdentifier
	CredentialAccessProfiles []CredentialAccessProfile
	ExtendedGrantTime        *bool
	Attributes               []CredentialAttribute
}

// CredentialIdentifier represents a credential identifier (card number, PIN, biometric, etc.).
type CredentialIdentifier struct {
	Type                      CredentialIdentifierType
	ExemptedFromAuthentication bool
	Value                     []byte // hex binary
}

// CredentialIdentifierType specifies the name and format type of a credential identifier.
type CredentialIdentifierType struct {
	Name       string
	FormatType string
}

// CredentialAccessProfile represents the association between a credential and an access profile.
type CredentialAccessProfile struct {
	AccessProfileToken string
	ValidFrom          *time.Time
	ValidTo            *time.Time
}

// CredentialAttribute represents a name/value attribute pair on a credential.
type CredentialAttribute struct {
	Name  string
	Value string
}

// CredentialState contains state information for a credential.
type CredentialState struct {
	Enabled          bool
	Reason           string
	AntipassbackState *AntipassbackState
}

// AntipassbackState contains anti-passback state information for a credential.
type AntipassbackState struct {
	AntipassbackViolated bool
}

// CredentialIdentifierFormatTypeInfo contains information about a format type.
type CredentialIdentifierFormatTypeInfo struct {
	FormatType  string
	Description string
}

// CredentialData holds a credential and its associated state (used by SetCredential).
type CredentialData struct {
	Credential      Credential
	CredentialState CredentialState
}

// CredentialIdentifierItem is a credential identifier used in whitelist/blacklist operations.
type CredentialIdentifierItem struct {
	Type  CredentialIdentifierType
	Value []byte // hex binary
}
