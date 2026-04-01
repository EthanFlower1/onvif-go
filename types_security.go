package onvif

// Advanced Security service types.

// KeyAttribute represents the attributes of a key in the keystore.
type KeyAttribute struct {
	KeyID              string
	Alias              *string
	HasPrivateKey      *bool
	KeyStatus          string
	ExternallyGenerated *bool
	SecurelyStored     *bool
}

// X509Certificate represents an X.509 certificate.
type X509Certificate struct {
	CertificateID      string
	KeyID              string
	Alias              *string
	CertificateContent []byte
	HasPrivateKey      *bool
}

// CertificationPath represents an X.509 certification path.
type CertificationPath struct {
	CertificateIDs []string
	Alias          *string
}

// PassphraseAttribute represents a passphrase entry in the keystore.
type PassphraseAttribute struct {
	PassphraseID string
	Alias        *string
}

// CRLID is a unique identifier for CRLs.
type CRLID = string

// CRL represents a Certificate Revocation List.
type CRL struct {
	CRLID      string
	Alias      string
	CRLContent []byte
}

// CertPathValidationPolicyID is a unique identifier for cert path validation policies.
type CertPathValidationPolicyID = string

// CertPathValidationParameters holds parameters for a cert path validation policy.
type CertPathValidationParameters struct {
	RequireTLSWWWClientAuthExtendedKeyUsage *bool
	UseDeltaCRLs                            *bool
}

// TrustAnchor represents a trust anchor certificate.
type TrustAnchor struct {
	CertificateID string
}

// CertPathValidationPolicy represents a certification path validation policy.
type CertPathValidationPolicy struct {
	CertPathValidationPolicyID string
	Alias                      *string
	Parameters                 CertPathValidationParameters
	TrustAnchors               []TrustAnchor
}

// AdvSecDot1XStage represents one stage of 802.1X authentication (Advanced Security service).
type AdvSecDot1XStage struct {
	Method                     string
	Identity                   *string
	CertificationPathID        *string
	PassphraseID               *string
	CertPathValidationPolicyID *string
	Inner                      *AdvSecDot1XStage
}

// AdvSecDot1XConfiguration represents an IEEE 802.1X configuration (Advanced Security service).
type AdvSecDot1XConfiguration struct {
	Dot1XID string
	Alias   *string
	Outer   AdvSecDot1XStage
}

// AlgorithmIdentifier identifies a cryptographic algorithm.
type AlgorithmIdentifier struct {
	Algorithm  string
	Parameters []byte
}

// DistinguishedName represents an X.500 distinguished name.
type DistinguishedName struct {
	Country                    []string
	Organization               []string
	OrganizationalUnit         []string
	DistinguishedNameQualifier []string
	StateOrProvinceName        []string
	CommonName                 []string
	SerialNumber               []string
	Locality                   []string
	Title                      []string
	Surname                    []string
	GivenName                  []string
	Initials                   []string
	Pseudonym                  []string
	GenerationQualifier        []string
}

// KeystoreCapabilities represents the capabilities of a keystore implementation.
type KeystoreCapabilities struct {
	MaximumNumberOfKeys               *uint
	MaximumNumberOfCertificates       *uint
	MaximumNumberOfCertificationPaths *uint
	RSAKeyPairGeneration              *bool
	ECCKeyPairGeneration              *bool
	PKCS10                            *bool
	SelfSignedCertificateCreation     *bool
	PKCS8                             *bool
	PKCS12                            *bool
	MaximumNumberOfPassphrases        *uint
	MaximumNumberOfCRLs               *uint
	MaximumNumberOfCertificationPathValidationPolicies *uint
}

// TLSServerCapabilities represents the capabilities of the TLS server implementation.
type TLSServerCapabilities struct {
	TLSServerSupported                              []string
	EnabledVersionsSupported                        *bool
	MaximumNumberOfTLSCertificationPaths            *uint
	TLSClientAuthSupported                          *bool
	CnMapsToUserSupported                           *bool
	MaximumNumberOfTLSCertificationPathValidationPolicies *uint
}

// Dot1XCapabilities represents the capabilities of the 802.1X implementation.
type Dot1XCapabilities struct {
	MaximumNumberOfDot1XConfigurations *uint
	Dot1XMethods                       []string
}

// AuthorizationServerConfigurationCapabilities represents capabilities for auth server.
type AuthorizationServerConfigurationCapabilities struct {
	MaxConfigurations                  *int
	ConfigurationTypesSupported        []string
	ClientAuthenticationMethodsSupported []string
}

// MediaSigningCapabilities represents capabilities for media signing.
type MediaSigningCapabilities struct {
	MediaSigningSupported         *bool
	UserMediaSigningKeySupported  *bool
}

// AdvancedSecurityCapabilities represents capabilities of the Advanced Security service.
type AdvancedSecurityCapabilities struct {
	KeystoreCapabilities                KeystoreCapabilities
	TLSServerCapabilities               TLSServerCapabilities
	Dot1XCapabilities                   *Dot1XCapabilities
	AuthorizationServerCapabilities     *AuthorizationServerConfigurationCapabilities
	MediaSigningCapabilities            *MediaSigningCapabilities
}

// AuthorizationServerConfigurationData holds configuration for an authorization server.
type AuthorizationServerConfigurationData struct {
	Type                       string
	ClientAuth                 *string
	ServerURI                  string
	ClientID                   *string
	ClientSecret               *string
	Scope                      *string
	KeyID                      *string
	CertificateID              *string
	CertPathValidationPolicyID *string
}

// AuthorizationServerConfiguration includes the token and configuration data.
type AuthorizationServerConfiguration struct {
	Token string
	Data  AuthorizationServerConfigurationData
}

// JWTCustomClaim represents a custom JWT claim.
type JWTCustomClaim struct {
	Name            string
	SupportedValues []string
}

// JWTConfiguration holds the JWT authorization configuration.
type JWTConfiguration struct {
	Audiences        []string
	TrustedIssuers   []string
	KeyIDs           []string
	ValidationPolicy *string
	CustomClaims     []JWTCustomClaim
}

// CreateRSAKeyPairResponse holds the result of CreateRSAKeyPair.
type CreateRSAKeyPairResponse struct {
	KeyID                 string
	EstimatedCreationTime string
}

// CreateECCKeyPairResponse holds the result of CreateECCKeyPair.
type CreateECCKeyPairResponse struct {
	KeyID                 string
	EstimatedCreationTime string
}

// UploadCertificateWithPrivateKeyInPKCS12Response holds result of that upload operation.
type UploadCertificateWithPrivateKeyInPKCS12Response struct {
	CertificationPathID string
	KeyID               string
}

// UploadCertificateResponse holds the result of UploadCertificate.
type UploadCertificateResponse struct {
	CertificateID string
	KeyID         string
}
