package onvif

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/0x524a/onvif-go/internal/soap"
)

// Advanced Security service namespace.
const advancedSecurityNamespace = "http://www.onvif.org/ver10/advancedsecurity/wsdl"

// getAdvancedSecurityEndpoint returns the advanced security endpoint, falling back to device endpoint.
func (c *Client) getAdvancedSecurityEndpoint() string {
	if c.advancedSecurityEndpoint != "" {
		return c.advancedSecurityEndpoint
	}

	return c.endpoint
}

// newAdvancedSecuritySOAPClient creates a SOAP client for the advanced security service.
func (c *Client) newAdvancedSecuritySOAPClient() *soap.Client {
	username, password := c.GetCredentials()

	return soap.NewClient(c.httpClient, username, password)
}

// ============================================================
// Capabilities
// ============================================================

// GetAdvancedSecurityServiceCapabilities returns the capabilities of the Advanced Security service.
func (c *Client) GetAdvancedSecurityServiceCapabilities(ctx context.Context) (*AdvancedSecurityCapabilities, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
	}

	type keystoreCapsXML struct {
		MaximumNumberOfKeys               *uint   `xml:"MaximumNumberOfKeys,attr"`
		MaximumNumberOfCertificates       *uint   `xml:"MaximumNumberOfCertificates,attr"`
		MaximumNumberOfCertificationPaths *uint   `xml:"MaximumNumberOfCertificationPaths,attr"`
		RSAKeyPairGeneration              *bool   `xml:"RSAKeyPairGeneration,attr"`
		ECCKeyPairGeneration              *bool   `xml:"ECCKeyPairGeneration,attr"`
		PKCS10                            *bool   `xml:"PKCS10,attr"`
		SelfSignedCertificateCreation     *bool   `xml:"SelfSignedCertificateCreation,attr"`
		PKCS8                             *bool   `xml:"PKCS8,attr"`
		PKCS12                            *bool   `xml:"PKCS12,attr"`
		MaximumNumberOfPassphrases        *uint   `xml:"MaximumNumberOfPassphrases,attr"`
		MaximumNumberOfCRLs               *uint   `xml:"MaximumNumberOfCRLs,attr"`
		MaximumNumberOfCertificationPathValidationPolicies *uint `xml:"MaximumNumberOfCertificationPathValidationPolicies,attr"`
	}

	type tlsCapsXML struct {
		TLSServerSupported                              string `xml:"TLSServerSupported,attr"`
		EnabledVersionsSupported                        *bool  `xml:"EnabledVersionsSupported,attr"`
		MaximumNumberOfTLSCertificationPaths            *uint  `xml:"MaximumNumberOfTLSCertificationPaths,attr"`
		TLSClientAuthSupported                          *bool  `xml:"TLSClientAuthSupported,attr"`
		CnMapsToUserSupported                           *bool  `xml:"CnMapsToUserSupported,attr"`
		MaximumNumberOfTLSCertificationPathValidationPolicies *uint `xml:"MaximumNumberOfTLSCertificationPathValidationPolicies,attr"`
	}

	type dot1xCapsXML struct {
		MaximumNumberOfDot1XConfigurations *uint  `xml:"MaximumNumberOfDot1XConfigurations,attr"`
		Dot1XMethods                       string `xml:"Dot1XMethods,attr"`
	}

	type Response struct {
		XMLName      xml.Name `xml:"GetServiceCapabilitiesResponse"`
		Capabilities struct {
			KeystoreCapabilities keystoreCapsXML `xml:"KeystoreCapabilities"`
			TLSServerCapabilities tlsCapsXML      `xml:"TLSServerCapabilities"`
			Dot1XCapabilities    *dot1xCapsXML   `xml:"Dot1XCapabilities"`
		} `xml:"Capabilities"`
	}

	req := Request{Xmlns: advancedSecurityNamespace}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAdvancedSecurityServiceCapabilities failed: %w", err)
	}

	caps := &AdvancedSecurityCapabilities{
		KeystoreCapabilities: KeystoreCapabilities{
			MaximumNumberOfKeys:               resp.Capabilities.KeystoreCapabilities.MaximumNumberOfKeys,
			MaximumNumberOfCertificates:       resp.Capabilities.KeystoreCapabilities.MaximumNumberOfCertificates,
			MaximumNumberOfCertificationPaths: resp.Capabilities.KeystoreCapabilities.MaximumNumberOfCertificationPaths,
			RSAKeyPairGeneration:              resp.Capabilities.KeystoreCapabilities.RSAKeyPairGeneration,
			ECCKeyPairGeneration:              resp.Capabilities.KeystoreCapabilities.ECCKeyPairGeneration,
			PKCS10:                            resp.Capabilities.KeystoreCapabilities.PKCS10,
			SelfSignedCertificateCreation:     resp.Capabilities.KeystoreCapabilities.SelfSignedCertificateCreation,
			PKCS8:                             resp.Capabilities.KeystoreCapabilities.PKCS8,
			PKCS12:                            resp.Capabilities.KeystoreCapabilities.PKCS12,
			MaximumNumberOfPassphrases:        resp.Capabilities.KeystoreCapabilities.MaximumNumberOfPassphrases,
			MaximumNumberOfCRLs:               resp.Capabilities.KeystoreCapabilities.MaximumNumberOfCRLs,
			MaximumNumberOfCertificationPathValidationPolicies: resp.Capabilities.KeystoreCapabilities.MaximumNumberOfCertificationPathValidationPolicies,
		},
		TLSServerCapabilities: TLSServerCapabilities{
			EnabledVersionsSupported:             resp.Capabilities.TLSServerCapabilities.EnabledVersionsSupported,
			MaximumNumberOfTLSCertificationPaths: resp.Capabilities.TLSServerCapabilities.MaximumNumberOfTLSCertificationPaths,
			TLSClientAuthSupported:               resp.Capabilities.TLSServerCapabilities.TLSClientAuthSupported,
			CnMapsToUserSupported:                resp.Capabilities.TLSServerCapabilities.CnMapsToUserSupported,
			MaximumNumberOfTLSCertificationPathValidationPolicies: resp.Capabilities.TLSServerCapabilities.MaximumNumberOfTLSCertificationPathValidationPolicies,
		},
	}

	if resp.Capabilities.TLSServerCapabilities.TLSServerSupported != "" {
		caps.TLSServerCapabilities.TLSServerSupported = strings.Fields(resp.Capabilities.TLSServerCapabilities.TLSServerSupported)
	}

	if d := resp.Capabilities.Dot1XCapabilities; d != nil {
		dot1x := &Dot1XCapabilities{
			MaximumNumberOfDot1XConfigurations: d.MaximumNumberOfDot1XConfigurations,
		}

		if d.Dot1XMethods != "" {
			dot1x.Dot1XMethods = strings.Fields(d.Dot1XMethods)
		}

		caps.Dot1XCapabilities = dot1x
	}

	return caps, nil
}

// ============================================================
// JWT operations
// ============================================================

// GetJWTConfiguration returns the JWT authorization configuration.
func (c *Client) GetJWTConfiguration(ctx context.Context) (*JWTConfiguration, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetJWTConfiguration"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
	}

	type customClaimXML struct {
		Name            string   `xml:"Name"`
		SupportedValues string   `xml:"SupportedValues"`
	}

	type Response struct {
		XMLName       xml.Name `xml:"GetJWTConfigurationResponse"`
		Configuration struct {
			Audiences      []string         `xml:"Audiences"`
			TrustedIssuers []string         `xml:"TrustedIssuers"`
			KeyID          []string         `xml:"KeyID"`
			ValidationPolicy *string        `xml:"ValidationPolicy"`
			CustomClaims   []customClaimXML `xml:"CustomClaims"`
		} `xml:"Configuration"`
	}

	req := Request{Xmlns: advancedSecurityNamespace}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetJWTConfiguration failed: %w", err)
	}

	cfg := &JWTConfiguration{
		Audiences:        resp.Configuration.Audiences,
		TrustedIssuers:   resp.Configuration.TrustedIssuers,
		KeyIDs:           resp.Configuration.KeyID,
		ValidationPolicy: resp.Configuration.ValidationPolicy,
	}

	for _, cc := range resp.Configuration.CustomClaims {
		cfg.CustomClaims = append(cfg.CustomClaims, JWTCustomClaim{
			Name:            cc.Name,
			SupportedValues: strings.Fields(cc.SupportedValues),
		})
	}

	return cfg, nil
}

// SetJWTConfiguration sets the JWT authorization configuration.
func (c *Client) SetJWTConfiguration(ctx context.Context, config JWTConfiguration) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type customClaimXML struct {
		Name            string `xml:"tas:Name"`
		SupportedValues string `xml:"tas:SupportedValues"`
	}

	type configXML struct {
		Audiences      []string         `xml:"tas:Audiences"`
		TrustedIssuers []string         `xml:"tas:TrustedIssuers"`
		KeyID          []string         `xml:"tas:KeyID"`
		ValidationPolicy *string        `xml:"tas:ValidationPolicy"`
		CustomClaims   []customClaimXML `xml:"tas:CustomClaims"`
	}

	type Request struct {
		XMLName       xml.Name  `xml:"tas:SetJWTConfiguration"`
		Xmlns         string    `xml:"xmlns:tas,attr"`
		Configuration configXML `xml:"tas:Configuration"`
	}

	type Response struct {
		XMLName xml.Name `xml:"SetJWTConfigurationResponse"`
	}

	claims := make([]customClaimXML, 0, len(config.CustomClaims))
	for _, cc := range config.CustomClaims {
		claims = append(claims, customClaimXML{
			Name:            cc.Name,
			SupportedValues: strings.Join(cc.SupportedValues, " "),
		})
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		Configuration: configXML{
			Audiences:      config.Audiences,
			TrustedIssuers: config.TrustedIssuers,
			KeyID:          config.KeyIDs,
			ValidationPolicy: config.ValidationPolicy,
			CustomClaims:   claims,
		},
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetJWTConfiguration failed: %w", err)
	}

	return nil
}

// ============================================================
// Key management operations
// ============================================================

// CreateRSAKeyPair triggers asynchronous generation of an RSA key pair.
func (c *Client) CreateRSAKeyPair(ctx context.Context, keyLength uint, alias *string) (*CreateRSAKeyPairResponse, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName   xml.Name `xml:"tas:CreateRSAKeyPair"`
		Xmlns     string   `xml:"xmlns:tas,attr"`
		KeyLength uint     `xml:"tas:KeyLength"`
		Alias     *string  `xml:"tas:Alias,omitempty"`
	}

	type Response struct {
		XMLName               xml.Name `xml:"CreateRSAKeyPairResponse"`
		KeyID                 string   `xml:"KeyID"`
		EstimatedCreationTime string   `xml:"EstimatedCreationTime"`
	}

	req := Request{
		Xmlns:     advancedSecurityNamespace,
		KeyLength: keyLength,
		Alias:     alias,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("CreateRSAKeyPair failed: %w", err)
	}

	return &CreateRSAKeyPairResponse{
		KeyID:                 resp.KeyID,
		EstimatedCreationTime: resp.EstimatedCreationTime,
	}, nil
}

// CreateECCKeyPair triggers asynchronous generation of an ECC key pair.
func (c *Client) CreateECCKeyPair(ctx context.Context, ellipticCurve string, alias *string) (*CreateECCKeyPairResponse, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName       xml.Name `xml:"tas:CreateECCKeyPair"`
		Xmlns         string   `xml:"xmlns:tas,attr"`
		EllipticCurve string   `xml:"tas:EllipticCurve"`
		Alias         *string  `xml:"tas:Alias,omitempty"`
	}

	type Response struct {
		XMLName               xml.Name `xml:"CreateECCKeyPairResponse"`
		KeyID                 string   `xml:"KeyID"`
		EstimatedCreationTime string   `xml:"EstimatedCreationTime"`
	}

	req := Request{
		Xmlns:         advancedSecurityNamespace,
		EllipticCurve: ellipticCurve,
		Alias:         alias,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("CreateECCKeyPair failed: %w", err)
	}

	return &CreateECCKeyPairResponse{
		KeyID:                 resp.KeyID,
		EstimatedCreationTime: resp.EstimatedCreationTime,
	}, nil
}

// UploadKeyPairInPKCS8 uploads a key pair in PKCS#8 format.
func (c *Client) UploadKeyPairInPKCS8(ctx context.Context, keyPair []byte, alias *string, encryptionPassphraseID *string, encryptionPassphrase *string) (string, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName                xml.Name `xml:"tas:UploadKeyPairInPKCS8"`
		Xmlns                  string   `xml:"xmlns:tas,attr"`
		KeyPair                []byte   `xml:"tas:KeyPair"`
		Alias                  *string  `xml:"tas:Alias,omitempty"`
		EncryptionPassphraseID *string  `xml:"tas:EncryptionPassphraseID,omitempty"`
		EncryptionPassphrase   *string  `xml:"tas:EncryptionPassphrase,omitempty"`
	}

	type Response struct {
		XMLName xml.Name `xml:"UploadKeyPairInPKCS8Response"`
		KeyID   string   `xml:"KeyID"`
	}

	req := Request{
		Xmlns:                  advancedSecurityNamespace,
		KeyPair:                keyPair,
		Alias:                  alias,
		EncryptionPassphraseID: encryptionPassphraseID,
		EncryptionPassphrase:   encryptionPassphrase,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("UploadKeyPairInPKCS8 failed: %w", err)
	}

	return resp.KeyID, nil
}

// UploadCertificateWithPrivateKeyInPKCS12 uploads a PKCS#12 file containing cert and private key.
func (c *Client) UploadCertificateWithPrivateKeyInPKCS12(
	ctx context.Context,
	certWithPrivateKey []byte,
	certificationPathAlias *string,
	keyAlias *string,
	ignoreAdditionalCertificates *bool,
	integrityPassphraseID *string,
	encryptionPassphraseID *string,
	passphrase *string,
) (*UploadCertificateWithPrivateKeyInPKCS12Response, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName                      xml.Name `xml:"tas:UploadCertificateWithPrivateKeyInPKCS12"`
		Xmlns                        string   `xml:"xmlns:tas,attr"`
		CertWithPrivateKey           []byte   `xml:"tas:CertWithPrivateKey"`
		CertificationPathAlias       *string  `xml:"tas:CertificationPathAlias,omitempty"`
		KeyAlias                     *string  `xml:"tas:KeyAlias,omitempty"`
		IgnoreAdditionalCertificates *bool    `xml:"tas:IgnoreAdditionalCertificates,omitempty"`
		IntegrityPassphraseID        *string  `xml:"tas:IntegrityPassphraseID,omitempty"`
		EncryptionPassphraseID       *string  `xml:"tas:EncryptionPassphraseID,omitempty"`
		Passphrase                   *string  `xml:"tas:Passphrase,omitempty"`
	}

	type Response struct {
		XMLName             xml.Name `xml:"UploadCertificateWithPrivateKeyInPKCS12Response"`
		CertificationPathID string   `xml:"CertificationPathID"`
		KeyID               string   `xml:"KeyID"`
	}

	req := Request{
		Xmlns:                        advancedSecurityNamespace,
		CertWithPrivateKey:           certWithPrivateKey,
		CertificationPathAlias:       certificationPathAlias,
		KeyAlias:                     keyAlias,
		IgnoreAdditionalCertificates: ignoreAdditionalCertificates,
		IntegrityPassphraseID:        integrityPassphraseID,
		EncryptionPassphraseID:       encryptionPassphraseID,
		Passphrase:                   passphrase,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("UploadCertificateWithPrivateKeyInPKCS12 failed: %w", err)
	}

	return &UploadCertificateWithPrivateKeyInPKCS12Response{
		CertificationPathID: resp.CertificationPathID,
		KeyID:               resp.KeyID,
	}, nil
}

// GetKeyStatus returns the status of a key.
func (c *Client) GetKeyStatus(ctx context.Context, keyID string) (string, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetKeyStatus"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
		KeyID   string   `xml:"tas:KeyID"`
	}

	type Response struct {
		XMLName   xml.Name `xml:"GetKeyStatusResponse"`
		KeyStatus string   `xml:"KeyStatus"`
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		KeyID: keyID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("GetKeyStatus failed: %w", err)
	}

	return resp.KeyStatus, nil
}

// GetPrivateKeyStatus returns whether a key pair contains a private key.
func (c *Client) GetPrivateKeyStatus(ctx context.Context, keyID string) (bool, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetPrivateKeyStatus"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
		KeyID   string   `xml:"tas:KeyID"`
	}

	type Response struct {
		XMLName       xml.Name `xml:"GetPrivateKeyStatusResponse"`
		HasPrivateKey bool     `xml:"hasPrivateKey"`
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		KeyID: keyID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return false, fmt.Errorf("GetPrivateKeyStatus failed: %w", err)
	}

	return resp.HasPrivateKey, nil
}

// GetAllKeys returns information about all keys in the keystore.
func (c *Client) GetAllKeys(ctx context.Context) ([]KeyAttribute, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetAllKeys"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
	}

	type keyAttrXML struct {
		KeyID               string  `xml:"KeyID"`
		Alias               *string `xml:"Alias"`
		HasPrivateKey       *bool   `xml:"hasPrivateKey"`
		KeyStatus           string  `xml:"KeyStatus"`
		ExternallyGenerated *bool   `xml:"externallyGenerated"`
		SecurelyStored      *bool   `xml:"securelyStored"`
	}

	type Response struct {
		XMLName      xml.Name     `xml:"GetAllKeysResponse"`
		KeyAttribute []keyAttrXML `xml:"KeyAttribute"`
	}

	req := Request{Xmlns: advancedSecurityNamespace}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAllKeys failed: %w", err)
	}

	keys := make([]KeyAttribute, 0, len(resp.KeyAttribute))
	for _, k := range resp.KeyAttribute {
		keys = append(keys, KeyAttribute{
			KeyID:               k.KeyID,
			Alias:               k.Alias,
			HasPrivateKey:       k.HasPrivateKey,
			KeyStatus:           k.KeyStatus,
			ExternallyGenerated: k.ExternallyGenerated,
			SecurelyStored:      k.SecurelyStored,
		})
	}

	return keys, nil
}

// DeleteKey deletes a key from the keystore.
func (c *Client) DeleteKey(ctx context.Context, keyID string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:DeleteKey"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
		KeyID   string   `xml:"tas:KeyID"`
	}

	type Response struct {
		XMLName xml.Name `xml:"DeleteKeyResponse"`
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		KeyID: keyID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteKey failed: %w", err)
	}

	return nil
}

// CreatePKCS10CSR generates a DER-encoded PKCS#10 certification request.
func (c *Client) CreatePKCS10CSR(ctx context.Context, subject DistinguishedName, keyID string, signatureAlgorithm AlgorithmIdentifier) ([]byte, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type sigAlgXML struct {
		Algorithm  string  `xml:"tas:algorithm"`
		Parameters []byte  `xml:"tas:parameters,omitempty"`
	}

	type subjectXML struct {
		CommonName          []string `xml:"tas:CommonName,omitempty"`
		Organization        []string `xml:"tas:Organization,omitempty"`
		OrganizationalUnit  []string `xml:"tas:OrganizationalUnit,omitempty"`
		Country             []string `xml:"tas:Country,omitempty"`
		StateOrProvinceName []string `xml:"tas:StateOrProvinceName,omitempty"`
		Locality            []string `xml:"tas:Locality,omitempty"`
	}

	type Request struct {
		XMLName            xml.Name   `xml:"tas:CreatePKCS10CSR"`
		Xmlns              string     `xml:"xmlns:tas,attr"`
		Subject            subjectXML `xml:"tas:Subject"`
		KeyID              string     `xml:"tas:KeyID"`
		SignatureAlgorithm sigAlgXML  `xml:"tas:SignatureAlgorithm"`
	}

	type Response struct {
		XMLName   xml.Name `xml:"CreatePKCS10CSRResponse"`
		PKCS10CSR []byte   `xml:"PKCS10CSR"`
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		Subject: subjectXML{
			CommonName:          subject.CommonName,
			Organization:        subject.Organization,
			OrganizationalUnit:  subject.OrganizationalUnit,
			Country:             subject.Country,
			StateOrProvinceName: subject.StateOrProvinceName,
			Locality:            subject.Locality,
		},
		KeyID: keyID,
		SignatureAlgorithm: sigAlgXML{
			Algorithm:  signatureAlgorithm.Algorithm,
			Parameters: signatureAlgorithm.Parameters,
		},
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("CreatePKCS10CSR failed: %w", err)
	}

	return resp.PKCS10CSR, nil
}

// ============================================================
// Certificate operations
// ============================================================

// CreateSelfSignedCertificate generates a self-signed X.509 certificate.
func (c *Client) CreateSelfSignedCertificate(
	ctx context.Context,
	subject DistinguishedName,
	keyID string,
	alias *string,
	notValidBefore *string,
	notValidAfter *string,
	signatureAlgorithm AlgorithmIdentifier,
	x509Version *uint,
) (string, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type sigAlgXML struct {
		Algorithm  string `xml:"tas:algorithm"`
		Parameters []byte `xml:"tas:parameters,omitempty"`
	}

	type subjectXML struct {
		CommonName          []string `xml:"tas:CommonName,omitempty"`
		Organization        []string `xml:"tas:Organization,omitempty"`
		OrganizationalUnit  []string `xml:"tas:OrganizationalUnit,omitempty"`
		Country             []string `xml:"tas:Country,omitempty"`
		StateOrProvinceName []string `xml:"tas:StateOrProvinceName,omitempty"`
		Locality            []string `xml:"tas:Locality,omitempty"`
	}

	type Request struct {
		XMLName            xml.Name   `xml:"tas:CreateSelfSignedCertificate"`
		Xmlns              string     `xml:"xmlns:tas,attr"`
		X509Version        *uint      `xml:"tas:X509Version,omitempty"`
		Subject            subjectXML `xml:"tas:Subject"`
		KeyID              string     `xml:"tas:KeyID"`
		Alias              *string    `xml:"tas:Alias,omitempty"`
		NotValidBefore     *string    `xml:"tas:notValidBefore,omitempty"`
		NotValidAfter      *string    `xml:"tas:notValidAfter,omitempty"`
		SignatureAlgorithm sigAlgXML  `xml:"tas:SignatureAlgorithm"`
	}

	type Response struct {
		XMLName       xml.Name `xml:"CreateSelfSignedCertificateResponse"`
		CertificateID string   `xml:"CertificateID"`
	}

	req := Request{
		Xmlns:       advancedSecurityNamespace,
		X509Version: x509Version,
		Subject: subjectXML{
			CommonName:          subject.CommonName,
			Organization:        subject.Organization,
			OrganizationalUnit:  subject.OrganizationalUnit,
			Country:             subject.Country,
			StateOrProvinceName: subject.StateOrProvinceName,
			Locality:            subject.Locality,
		},
		KeyID:          keyID,
		Alias:          alias,
		NotValidBefore: notValidBefore,
		NotValidAfter:  notValidAfter,
		SignatureAlgorithm: sigAlgXML{
			Algorithm:  signatureAlgorithm.Algorithm,
			Parameters: signatureAlgorithm.Parameters,
		},
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateSelfSignedCertificate failed: %w", err)
	}

	return resp.CertificateID, nil
}

// UploadCertificate uploads an X.509 certificate to the keystore.
func (c *Client) UploadCertificate(ctx context.Context, certificate []byte, alias *string, keyAlias *string, privateKeyRequired *bool) (*UploadCertificateResponse, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName            xml.Name `xml:"tas:UploadCertificate"`
		Xmlns              string   `xml:"xmlns:tas,attr"`
		Certificate        []byte   `xml:"tas:Certificate"`
		Alias              *string  `xml:"tas:Alias,omitempty"`
		KeyAlias           *string  `xml:"tas:KeyAlias,omitempty"`
		PrivateKeyRequired *bool    `xml:"tas:PrivateKeyRequired,omitempty"`
	}

	type Response struct {
		XMLName       xml.Name `xml:"UploadCertificateResponse"`
		CertificateID string   `xml:"CertificateID"`
		KeyID         string   `xml:"KeyID"`
	}

	req := Request{
		Xmlns:              advancedSecurityNamespace,
		Certificate:        certificate,
		Alias:              alias,
		KeyAlias:           keyAlias,
		PrivateKeyRequired: privateKeyRequired,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("UploadCertificate failed: %w", err)
	}

	return &UploadCertificateResponse{
		CertificateID: resp.CertificateID,
		KeyID:         resp.KeyID,
	}, nil
}

// GetCertificate returns a specific certificate from the keystore.
func (c *Client) GetCertificate(ctx context.Context, certificateID string) (*X509Certificate, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName       xml.Name `xml:"tas:GetCertificate"`
		Xmlns         string   `xml:"xmlns:tas,attr"`
		CertificateID string   `xml:"tas:CertificateID"`
	}

	type certXML struct {
		CertificateID      string  `xml:"CertificateID"`
		KeyID              string  `xml:"KeyID"`
		Alias              *string `xml:"Alias"`
		CertificateContent []byte  `xml:"CertificateContent"`
		HasPrivateKey      *bool   `xml:"HasPrivateKey"`
	}

	type Response struct {
		XMLName     xml.Name `xml:"GetCertificateResponse"`
		Certificate certXML  `xml:"Certificate"`
	}

	req := Request{
		Xmlns:         advancedSecurityNamespace,
		CertificateID: certificateID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetCertificate failed: %w", err)
	}

	return &X509Certificate{
		CertificateID:      resp.Certificate.CertificateID,
		KeyID:              resp.Certificate.KeyID,
		Alias:              resp.Certificate.Alias,
		CertificateContent: resp.Certificate.CertificateContent,
		HasPrivateKey:      resp.Certificate.HasPrivateKey,
	}, nil
}

// GetAllCertificates returns all certificates stored in the keystore.
func (c *Client) GetAllCertificates(ctx context.Context) ([]X509Certificate, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetAllCertificates"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
	}

	type certXML struct {
		CertificateID      string  `xml:"CertificateID"`
		KeyID              string  `xml:"KeyID"`
		Alias              *string `xml:"Alias"`
		CertificateContent []byte  `xml:"CertificateContent"`
		HasPrivateKey      *bool   `xml:"HasPrivateKey"`
	}

	type Response struct {
		XMLName      xml.Name  `xml:"GetAllCertificatesResponse"`
		Certificates []certXML `xml:"Certificate"`
	}

	req := Request{Xmlns: advancedSecurityNamespace}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAllCertificates failed: %w", err)
	}

	certs := make([]X509Certificate, 0, len(resp.Certificates))
	for _, cert := range resp.Certificates {
		certs = append(certs, X509Certificate{
			CertificateID:      cert.CertificateID,
			KeyID:              cert.KeyID,
			Alias:              cert.Alias,
			CertificateContent: cert.CertificateContent,
			HasPrivateKey:      cert.HasPrivateKey,
		})
	}

	return certs, nil
}

// DeleteCertificate deletes a certificate from the keystore.
func (c *Client) DeleteCertificate(ctx context.Context, certificateID string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName       xml.Name `xml:"tas:DeleteCertificate"`
		Xmlns         string   `xml:"xmlns:tas,attr"`
		CertificateID string   `xml:"tas:CertificateID"`
	}

	type Response struct {
		XMLName xml.Name `xml:"DeleteCertificateResponse"`
	}

	req := Request{
		Xmlns:         advancedSecurityNamespace,
		CertificateID: certificateID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteCertificate failed: %w", err)
	}

	return nil
}

// ============================================================
// Certification Path operations
// ============================================================

// CreateCertificationPath creates a certification path from a sequence of certificate IDs.
func (c *Client) CreateCertificationPath(ctx context.Context, certificateIDs []string, alias *string) (string, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type certIDsXML struct {
		CertificateID []string `xml:"tas:CertificateID"`
	}

	type Request struct {
		XMLName        xml.Name   `xml:"tas:CreateCertificationPath"`
		Xmlns          string     `xml:"xmlns:tas,attr"`
		CertificateIDs certIDsXML `xml:"tas:CertificateIDs"`
		Alias          *string    `xml:"tas:Alias,omitempty"`
	}

	type Response struct {
		XMLName             xml.Name `xml:"CreateCertificationPathResponse"`
		CertificationPathID string   `xml:"CertificationPathID"`
	}

	req := Request{
		Xmlns:          advancedSecurityNamespace,
		CertificateIDs: certIDsXML{CertificateID: certificateIDs},
		Alias:          alias,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateCertificationPath failed: %w", err)
	}

	return resp.CertificationPathID, nil
}

// GetCertificationPath returns a specific certification path.
func (c *Client) GetCertificationPath(ctx context.Context, certificationPathID string) (*CertificationPath, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName             xml.Name `xml:"tas:GetCertificationPath"`
		Xmlns               string   `xml:"xmlns:tas,attr"`
		CertificationPathID string   `xml:"tas:CertificationPathID"`
	}

	type pathXML struct {
		CertificateID []string `xml:"CertificateID"`
		Alias         *string  `xml:"Alias"`
	}

	type Response struct {
		XMLName           xml.Name `xml:"GetCertificationPathResponse"`
		CertificationPath pathXML  `xml:"CertificationPath"`
	}

	req := Request{
		Xmlns:               advancedSecurityNamespace,
		CertificationPathID: certificationPathID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetCertificationPath failed: %w", err)
	}

	return &CertificationPath{
		CertificateIDs: resp.CertificationPath.CertificateID,
		Alias:          resp.CertificationPath.Alias,
	}, nil
}

// GetAllCertificationPaths returns the IDs of all certification paths in the keystore.
func (c *Client) GetAllCertificationPaths(ctx context.Context) ([]string, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetAllCertificationPaths"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
	}

	type Response struct {
		XMLName              xml.Name `xml:"GetAllCertificationPathsResponse"`
		CertificationPathID []string `xml:"CertificationPathID"`
	}

	req := Request{Xmlns: advancedSecurityNamespace}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAllCertificationPaths failed: %w", err)
	}

	return resp.CertificationPathID, nil
}

// SetCertificationPath modifies an existing certification path.
func (c *Client) SetCertificationPath(ctx context.Context, certificationPathID string, path CertificationPath) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type pathXML struct {
		CertificateID []string `xml:"tas:CertificateID"`
		Alias         *string  `xml:"tas:Alias,omitempty"`
	}

	type Request struct {
		XMLName             xml.Name `xml:"tas:SetCertificationPath"`
		Xmlns               string   `xml:"xmlns:tas,attr"`
		CertificationPathID string   `xml:"tas:CertificationPathID"`
		CertificationPath   pathXML  `xml:"tas:CertificationPath"`
	}

	type Response struct {
		XMLName xml.Name `xml:"SetCertificationPathResponse"`
	}

	req := Request{
		Xmlns:               advancedSecurityNamespace,
		CertificationPathID: certificationPathID,
		CertificationPath: pathXML{
			CertificateID: path.CertificateIDs,
			Alias:         path.Alias,
		},
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetCertificationPath failed: %w", err)
	}

	return nil
}

// DeleteCertificationPath deletes a certification path from the keystore.
func (c *Client) DeleteCertificationPath(ctx context.Context, certificationPathID string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName             xml.Name `xml:"tas:DeleteCertificationPath"`
		Xmlns               string   `xml:"xmlns:tas,attr"`
		CertificationPathID string   `xml:"tas:CertificationPathID"`
	}

	type Response struct {
		XMLName xml.Name `xml:"DeleteCertificationPathResponse"`
	}

	req := Request{
		Xmlns:               advancedSecurityNamespace,
		CertificationPathID: certificationPathID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteCertificationPath failed: %w", err)
	}

	return nil
}

// ============================================================
// Passphrase operations
// ============================================================

// UploadPassphrase uploads a passphrase to the keystore.
func (c *Client) UploadPassphrase(ctx context.Context, passphrase string, alias *string) (string, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName         xml.Name `xml:"tas:UploadPassphrase"`
		Xmlns           string   `xml:"xmlns:tas,attr"`
		Passphrase      string   `xml:"tas:Passphrase"`
		PassphraseAlias *string  `xml:"tas:PassphraseAlias,omitempty"`
	}

	type Response struct {
		XMLName      xml.Name `xml:"UploadPassphraseResponse"`
		PassphraseID string   `xml:"PassphraseID"`
	}

	req := Request{
		Xmlns:           advancedSecurityNamespace,
		Passphrase:      passphrase,
		PassphraseAlias: alias,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("UploadPassphrase failed: %w", err)
	}

	return resp.PassphraseID, nil
}

// GetAllPassphrases returns information about all passphrases in the keystore.
func (c *Client) GetAllPassphrases(ctx context.Context) ([]PassphraseAttribute, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetAllPassphrases"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
	}

	type passphraseAttrXML struct {
		PassphraseID string  `xml:"PassphraseID"`
		Alias        *string `xml:"Alias"`
	}

	type Response struct {
		XMLName             xml.Name            `xml:"GetAllPassphrasesResponse"`
		PassphraseAttribute []passphraseAttrXML `xml:"PassphraseAttribute"`
	}

	req := Request{Xmlns: advancedSecurityNamespace}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAllPassphrases failed: %w", err)
	}

	passphrases := make([]PassphraseAttribute, 0, len(resp.PassphraseAttribute))
	for _, p := range resp.PassphraseAttribute {
		passphrases = append(passphrases, PassphraseAttribute{
			PassphraseID: p.PassphraseID,
			Alias:        p.Alias,
		})
	}

	return passphrases, nil
}

// DeletePassphrase deletes a passphrase from the keystore.
func (c *Client) DeletePassphrase(ctx context.Context, passphraseID string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName      xml.Name `xml:"tas:DeletePassphrase"`
		Xmlns        string   `xml:"xmlns:tas,attr"`
		PassphraseID string   `xml:"tas:PassphraseID"`
	}

	type Response struct {
		XMLName xml.Name `xml:"DeletePassphraseResponse"`
	}

	req := Request{
		Xmlns:        advancedSecurityNamespace,
		PassphraseID: passphraseID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeletePassphrase failed: %w", err)
	}

	return nil
}

// ============================================================
// TLS Server operations
// ============================================================

// AddServerCertificateAssignment assigns a certification path to the TLS server.
func (c *Client) AddServerCertificateAssignment(ctx context.Context, certificationPathID string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName             xml.Name `xml:"tas:AddServerCertificateAssignment"`
		Xmlns               string   `xml:"xmlns:tas,attr"`
		CertificationPathID string   `xml:"tas:CertificationPathID"`
	}

	type Response struct {
		XMLName xml.Name `xml:"AddServerCertificateAssignmentResponse"`
	}

	req := Request{
		Xmlns:               advancedSecurityNamespace,
		CertificationPathID: certificationPathID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("AddServerCertificateAssignment failed: %w", err)
	}

	return nil
}

// RemoveServerCertificateAssignment removes a certification path assignment from the TLS server.
func (c *Client) RemoveServerCertificateAssignment(ctx context.Context, certificationPathID string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName             xml.Name `xml:"tas:RemoveServerCertificateAssignment"`
		Xmlns               string   `xml:"xmlns:tas,attr"`
		CertificationPathID string   `xml:"tas:CertificationPathID"`
	}

	type Response struct {
		XMLName xml.Name `xml:"RemoveServerCertificateAssignmentResponse"`
	}

	req := Request{
		Xmlns:               advancedSecurityNamespace,
		CertificationPathID: certificationPathID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("RemoveServerCertificateAssignment failed: %w", err)
	}

	return nil
}

// ReplaceServerCertificateAssignment replaces a TLS server certification path assignment.
func (c *Client) ReplaceServerCertificateAssignment(ctx context.Context, oldCertificationPathID, newCertificationPathID string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName                xml.Name `xml:"tas:ReplaceServerCertificateAssignment"`
		Xmlns                  string   `xml:"xmlns:tas,attr"`
		OldCertificationPathID string   `xml:"tas:OldCertificationPathID"`
		NewCertificationPathID string   `xml:"tas:NewCertificationPathID"`
	}

	type Response struct {
		XMLName xml.Name `xml:"ReplaceServerCertificateAssignmentResponse"`
	}

	req := Request{
		Xmlns:                  advancedSecurityNamespace,
		OldCertificationPathID: oldCertificationPathID,
		NewCertificationPathID: newCertificationPathID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("ReplaceServerCertificateAssignment failed: %w", err)
	}

	return nil
}

// GetAssignedServerCertificates returns all certification paths assigned to the TLS server.
func (c *Client) GetAssignedServerCertificates(ctx context.Context) ([]string, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetAssignedServerCertificates"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
	}

	type Response struct {
		XMLName              xml.Name `xml:"GetAssignedServerCertificatesResponse"`
		CertificationPathID []string `xml:"CertificationPathID"`
	}

	req := Request{Xmlns: advancedSecurityNamespace}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAssignedServerCertificates failed: %w", err)
	}

	return resp.CertificationPathID, nil
}

// SetEnabledTLSVersions sets the list of enabled TLS versions.
func (c *Client) SetEnabledTLSVersions(ctx context.Context, versions []string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName  xml.Name `xml:"tas:SetEnabledTLSVersions"`
		Xmlns    string   `xml:"xmlns:tas,attr"`
		Versions string   `xml:"tas:Versions"`
	}

	type Response struct {
		XMLName xml.Name `xml:"SetEnabledTLSVersionsResponse"`
	}

	req := Request{
		Xmlns:    advancedSecurityNamespace,
		Versions: strings.Join(versions, " "),
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetEnabledTLSVersions failed: %w", err)
	}

	return nil
}

// GetEnabledTLSVersions returns the list of enabled TLS versions.
func (c *Client) GetEnabledTLSVersions(ctx context.Context) ([]string, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetEnabledTLSVersions"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
	}

	type Response struct {
		XMLName  xml.Name `xml:"GetEnabledTLSVersionsResponse"`
		Versions string   `xml:"Versions"`
	}

	req := Request{Xmlns: advancedSecurityNamespace}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetEnabledTLSVersions failed: %w", err)
	}

	if resp.Versions == "" {
		return nil, nil
	}

	return strings.Fields(resp.Versions), nil
}

// ============================================================
// CRL operations
// ============================================================

// UploadCRL uploads a Certificate Revocation List to the device.
func (c *Client) UploadCRL(ctx context.Context, crl []byte, alias *string) (string, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:UploadCRL"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
		Crl     []byte   `xml:"tas:Crl"`
		Alias   *string  `xml:"tas:Alias,omitempty"`
	}

	type Response struct {
		XMLName xml.Name `xml:"UploadCRLResponse"`
		CrlID   string   `xml:"CrlID"`
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		Crl:   crl,
		Alias: alias,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("UploadCRL failed: %w", err)
	}

	return resp.CrlID, nil
}

// GetCRL returns a specific CRL from the device.
func (c *Client) GetCRL(ctx context.Context, crlID string) (*CRL, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetCRL"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
		CrlID   string   `xml:"tas:CrlID"`
	}

	type crlXML struct {
		CRLID      string `xml:"CRLID"`
		Alias      string `xml:"Alias"`
		CRLContent []byte `xml:"CRLContent"`
	}

	type Response struct {
		XMLName xml.Name `xml:"GetCRLResponse"`
		Crl     crlXML   `xml:"Crl"`
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		CrlID: crlID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetCRL failed: %w", err)
	}

	return &CRL{
		CRLID:      resp.Crl.CRLID,
		Alias:      resp.Crl.Alias,
		CRLContent: resp.Crl.CRLContent,
	}, nil
}

// GetAllCRLs returns all CRLs stored on the device.
func (c *Client) GetAllCRLs(ctx context.Context) ([]CRL, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetAllCRLs"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
	}

	type crlXML struct {
		CRLID      string `xml:"CRLID"`
		Alias      string `xml:"Alias"`
		CRLContent []byte `xml:"CRLContent"`
	}

	type Response struct {
		XMLName xml.Name `xml:"GetAllCRLsResponse"`
		Crl     []crlXML `xml:"Crl"`
	}

	req := Request{Xmlns: advancedSecurityNamespace}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAllCRLs failed: %w", err)
	}

	crls := make([]CRL, 0, len(resp.Crl))
	for _, crl := range resp.Crl {
		crls = append(crls, CRL{
			CRLID:      crl.CRLID,
			Alias:      crl.Alias,
			CRLContent: crl.CRLContent,
		})
	}

	return crls, nil
}

// DeleteCRL deletes a CRL from the device.
func (c *Client) DeleteCRL(ctx context.Context, crlID string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:DeleteCRL"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
		CrlID   string   `xml:"tas:CrlID"`
	}

	type Response struct {
		XMLName xml.Name `xml:"DeleteCRLResponse"`
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		CrlID: crlID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteCRL failed: %w", err)
	}

	return nil
}

// ============================================================
// Cert Path Validation Policy operations
// ============================================================

// CreateCertPathValidationPolicy creates a certification path validation policy.
func (c *Client) CreateCertPathValidationPolicy(ctx context.Context, alias *string, parameters CertPathValidationParameters, trustAnchors []TrustAnchor) (string, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type paramsXML struct {
		RequireTLSWWWClientAuthExtendedKeyUsage *bool `xml:"tas:RequireTLSWWWClientAuthExtendedKeyUsage,omitempty"`
		UseDeltaCRLs                            *bool `xml:"tas:UseDeltaCRLs,omitempty"`
	}

	type trustAnchorXML struct {
		CertificateID string `xml:"tas:CertificateID"`
	}

	type Request struct {
		XMLName      xml.Name         `xml:"tas:CreateCertPathValidationPolicy"`
		Xmlns        string           `xml:"xmlns:tas,attr"`
		Alias        *string          `xml:"tas:Alias,omitempty"`
		Parameters   paramsXML        `xml:"tas:Parameters"`
		TrustAnchor  []trustAnchorXML `xml:"tas:TrustAnchor"`
	}

	type Response struct {
		XMLName                    xml.Name `xml:"CreateCertPathValidationPolicyResponse"`
		CertPathValidationPolicyID string   `xml:"CertPathValidationPolicyID"`
	}

	anchors := make([]trustAnchorXML, 0, len(trustAnchors))
	for _, ta := range trustAnchors {
		anchors = append(anchors, trustAnchorXML{CertificateID: ta.CertificateID})
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		Alias: alias,
		Parameters: paramsXML{
			RequireTLSWWWClientAuthExtendedKeyUsage: parameters.RequireTLSWWWClientAuthExtendedKeyUsage,
			UseDeltaCRLs:                            parameters.UseDeltaCRLs,
		},
		TrustAnchor: anchors,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateCertPathValidationPolicy failed: %w", err)
	}

	return resp.CertPathValidationPolicyID, nil
}

// GetCertPathValidationPolicy returns a specific cert path validation policy.
func (c *Client) GetCertPathValidationPolicy(ctx context.Context, certPathValidationPolicyID string) (*CertPathValidationPolicy, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName                    xml.Name `xml:"tas:GetCertPathValidationPolicy"`
		Xmlns                      string   `xml:"xmlns:tas,attr"`
		CertPathValidationPolicyID string   `xml:"tas:CertPathValidationPolicyID"`
	}

	type trustAnchorXML struct {
		CertificateID string `xml:"CertificateID"`
	}

	type policyXML struct {
		CertPathValidationPolicyID string `xml:"CertPathValidationPolicyID"`
		Alias                      *string `xml:"Alias"`
		Parameters                 struct {
			RequireTLSWWWClientAuthExtendedKeyUsage *bool `xml:"RequireTLSWWWClientAuthExtendedKeyUsage"`
			UseDeltaCRLs                            *bool `xml:"UseDeltaCRLs"`
		} `xml:"Parameters"`
		TrustAnchor []trustAnchorXML `xml:"TrustAnchor"`
	}

	type Response struct {
		XMLName                  xml.Name  `xml:"GetCertPathValidationPolicyResponse"`
		CertPathValidationPolicy policyXML `xml:"CertPathValidationPolicy"`
	}

	req := Request{
		Xmlns:                      advancedSecurityNamespace,
		CertPathValidationPolicyID: certPathValidationPolicyID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetCertPathValidationPolicy failed: %w", err)
	}

	policy := &CertPathValidationPolicy{
		CertPathValidationPolicyID: resp.CertPathValidationPolicy.CertPathValidationPolicyID,
		Alias:                      resp.CertPathValidationPolicy.Alias,
		Parameters: CertPathValidationParameters{
			RequireTLSWWWClientAuthExtendedKeyUsage: resp.CertPathValidationPolicy.Parameters.RequireTLSWWWClientAuthExtendedKeyUsage,
			UseDeltaCRLs:                            resp.CertPathValidationPolicy.Parameters.UseDeltaCRLs,
		},
	}

	for _, ta := range resp.CertPathValidationPolicy.TrustAnchor {
		policy.TrustAnchors = append(policy.TrustAnchors, TrustAnchor{CertificateID: ta.CertificateID})
	}

	return policy, nil
}

// GetAllCertPathValidationPolicies returns all cert path validation policies.
func (c *Client) GetAllCertPathValidationPolicies(ctx context.Context) ([]CertPathValidationPolicy, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetAllCertPathValidationPolicies"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
	}

	type trustAnchorXML struct {
		CertificateID string `xml:"CertificateID"`
	}

	type policyXML struct {
		CertPathValidationPolicyID string `xml:"CertPathValidationPolicyID"`
		Alias                      *string `xml:"Alias"`
		Parameters                 struct {
			RequireTLSWWWClientAuthExtendedKeyUsage *bool `xml:"RequireTLSWWWClientAuthExtendedKeyUsage"`
			UseDeltaCRLs                            *bool `xml:"UseDeltaCRLs"`
		} `xml:"Parameters"`
		TrustAnchor []trustAnchorXML `xml:"TrustAnchor"`
	}

	type Response struct {
		XMLName                    xml.Name    `xml:"GetAllCertPathValidationPoliciesResponse"`
		CertPathValidationPolicy   []policyXML `xml:"CertPathValidationPolicy"`
	}

	req := Request{Xmlns: advancedSecurityNamespace}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAllCertPathValidationPolicies failed: %w", err)
	}

	policies := make([]CertPathValidationPolicy, 0, len(resp.CertPathValidationPolicy))
	for _, p := range resp.CertPathValidationPolicy {
		policy := CertPathValidationPolicy{
			CertPathValidationPolicyID: p.CertPathValidationPolicyID,
			Alias:                      p.Alias,
			Parameters: CertPathValidationParameters{
				RequireTLSWWWClientAuthExtendedKeyUsage: p.Parameters.RequireTLSWWWClientAuthExtendedKeyUsage,
				UseDeltaCRLs:                            p.Parameters.UseDeltaCRLs,
			},
		}

		for _, ta := range p.TrustAnchor {
			policy.TrustAnchors = append(policy.TrustAnchors, TrustAnchor{CertificateID: ta.CertificateID})
		}

		policies = append(policies, policy)
	}

	return policies, nil
}

// SetCertPathValidationPolicy modifies an existing cert path validation policy.
func (c *Client) SetCertPathValidationPolicy(ctx context.Context, policy CertPathValidationPolicy) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type paramsXML struct {
		RequireTLSWWWClientAuthExtendedKeyUsage *bool `xml:"tas:RequireTLSWWWClientAuthExtendedKeyUsage,omitempty"`
		UseDeltaCRLs                            *bool `xml:"tas:UseDeltaCRLs,omitempty"`
	}

	type trustAnchorXML struct {
		CertificateID string `xml:"tas:CertificateID"`
	}

	type policyXML struct {
		CertPathValidationPolicyID string           `xml:"tas:CertPathValidationPolicyID"`
		Alias                      *string          `xml:"tas:Alias,omitempty"`
		Parameters                 paramsXML        `xml:"tas:Parameters"`
		TrustAnchor                []trustAnchorXML `xml:"tas:TrustAnchor"`
	}

	type Request struct {
		XMLName                  xml.Name  `xml:"tas:SetCertPathValidationPolicy"`
		Xmlns                    string    `xml:"xmlns:tas,attr"`
		CertPathValidationPolicy policyXML `xml:"tas:CertPathValidationPolicy"`
	}

	type Response struct {
		XMLName xml.Name `xml:"SetCertPathValidationPolicyResponse"`
	}

	anchors := make([]trustAnchorXML, 0, len(policy.TrustAnchors))
	for _, ta := range policy.TrustAnchors {
		anchors = append(anchors, trustAnchorXML{CertificateID: ta.CertificateID})
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		CertPathValidationPolicy: policyXML{
			CertPathValidationPolicyID: policy.CertPathValidationPolicyID,
			Alias:                      policy.Alias,
			Parameters: paramsXML{
				RequireTLSWWWClientAuthExtendedKeyUsage: policy.Parameters.RequireTLSWWWClientAuthExtendedKeyUsage,
				UseDeltaCRLs:                            policy.Parameters.UseDeltaCRLs,
			},
			TrustAnchor: anchors,
		},
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetCertPathValidationPolicy failed: %w", err)
	}

	return nil
}

// DeleteCertPathValidationPolicy deletes a cert path validation policy.
func (c *Client) DeleteCertPathValidationPolicy(ctx context.Context, certPathValidationPolicyID string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName                    xml.Name `xml:"tas:DeleteCertPathValidationPolicy"`
		Xmlns                      string   `xml:"xmlns:tas,attr"`
		CertPathValidationPolicyID string   `xml:"tas:CertPathValidationPolicyID"`
	}

	type Response struct {
		XMLName xml.Name `xml:"DeleteCertPathValidationPolicyResponse"`
	}

	req := Request{
		Xmlns:                      advancedSecurityNamespace,
		CertPathValidationPolicyID: certPathValidationPolicyID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteCertPathValidationPolicy failed: %w", err)
	}

	return nil
}

// ============================================================
// TLS Client Auth operations
// ============================================================

// SetClientAuthenticationRequired sets whether TLS client authentication is required.
func (c *Client) SetClientAuthenticationRequired(ctx context.Context, required bool) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName                      xml.Name `xml:"tas:SetClientAuthenticationRequired"`
		Xmlns                        string   `xml:"xmlns:tas,attr"`
		ClientAuthenticationRequired bool     `xml:"tas:clientAuthenticationRequired"`
	}

	type Response struct {
		XMLName xml.Name `xml:"SetClientAuthenticationRequiredResponse"`
	}

	req := Request{
		Xmlns:                        advancedSecurityNamespace,
		ClientAuthenticationRequired: required,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetClientAuthenticationRequired failed: %w", err)
	}

	return nil
}

// GetClientAuthenticationRequired returns whether TLS client authentication is required.
func (c *Client) GetClientAuthenticationRequired(ctx context.Context) (bool, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetClientAuthenticationRequired"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
	}

	type Response struct {
		XMLName                      xml.Name `xml:"GetClientAuthenticationRequiredResponse"`
		ClientAuthenticationRequired bool     `xml:"clientAuthenticationRequired"`
	}

	req := Request{Xmlns: advancedSecurityNamespace}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return false, fmt.Errorf("GetClientAuthenticationRequired failed: %w", err)
	}

	return resp.ClientAuthenticationRequired, nil
}

// SetCnMapsToUser sets whether CN maps to user for TLS client authentication.
func (c *Client) SetCnMapsToUser(ctx context.Context, cnMapsToUser bool) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName      xml.Name `xml:"tas:SetCnMapsToUser"`
		Xmlns        string   `xml:"xmlns:tas,attr"`
		CnMapsToUser bool     `xml:"tas:cnMapsToUser"`
	}

	type Response struct {
		XMLName xml.Name `xml:"SetCnMapsToUserResponse"`
	}

	req := Request{
		Xmlns:        advancedSecurityNamespace,
		CnMapsToUser: cnMapsToUser,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetCnMapsToUser failed: %w", err)
	}

	return nil
}

// GetCnMapsToUser returns whether CN maps to user for TLS client authentication.
func (c *Client) GetCnMapsToUser(ctx context.Context) (bool, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetCnMapsToUser"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
	}

	type Response struct {
		XMLName      xml.Name `xml:"GetCnMapsToUserResponse"`
		CnMapsToUser bool     `xml:"cnMapsToUser"`
	}

	req := Request{Xmlns: advancedSecurityNamespace}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return false, fmt.Errorf("GetCnMapsToUser failed: %w", err)
	}

	return resp.CnMapsToUser, nil
}

// AddCertPathValidationPolicyAssignment assigns a cert path validation policy to the TLS server.
func (c *Client) AddCertPathValidationPolicyAssignment(ctx context.Context, certPathValidationPolicyID string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName                    xml.Name `xml:"tas:AddCertPathValidationPolicyAssignment"`
		Xmlns                      string   `xml:"xmlns:tas,attr"`
		CertPathValidationPolicyID string   `xml:"tas:CertPathValidationPolicyID"`
	}

	type Response struct {
		XMLName xml.Name `xml:"AddCertPathValidationPolicyAssignmentResponse"`
	}

	req := Request{
		Xmlns:                      advancedSecurityNamespace,
		CertPathValidationPolicyID: certPathValidationPolicyID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("AddCertPathValidationPolicyAssignment failed: %w", err)
	}

	return nil
}

// RemoveCertPathValidationPolicyAssignment removes a cert path validation policy from the TLS server.
func (c *Client) RemoveCertPathValidationPolicyAssignment(ctx context.Context, certPathValidationPolicyID string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName                    xml.Name `xml:"tas:RemoveCertPathValidationPolicyAssignment"`
		Xmlns                      string   `xml:"xmlns:tas,attr"`
		CertPathValidationPolicyID string   `xml:"tas:CertPathValidationPolicyID"`
	}

	type Response struct {
		XMLName xml.Name `xml:"RemoveCertPathValidationPolicyAssignmentResponse"`
	}

	req := Request{
		Xmlns:                      advancedSecurityNamespace,
		CertPathValidationPolicyID: certPathValidationPolicyID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("RemoveCertPathValidationPolicyAssignment failed: %w", err)
	}

	return nil
}

// ReplaceCertPathValidationPolicyAssignment replaces a cert path validation policy assignment.
func (c *Client) ReplaceCertPathValidationPolicyAssignment(ctx context.Context, oldCertPathValidationPolicyID, newCertPathValidationPolicyID string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName                      xml.Name `xml:"tas:ReplaceCertPathValidationPolicyAssignment"`
		Xmlns                        string   `xml:"xmlns:tas,attr"`
		OldCertPathValidationPolicyID string  `xml:"tas:OldCertPathValidationPolicyID"`
		NewCertPathValidationPolicyID string  `xml:"tas:NewCertPathValidationPolicyID"`
	}

	type Response struct {
		XMLName xml.Name `xml:"ReplaceCertPathValidationPolicyAssignmentResponse"`
	}

	req := Request{
		Xmlns:                         advancedSecurityNamespace,
		OldCertPathValidationPolicyID: oldCertPathValidationPolicyID,
		NewCertPathValidationPolicyID: newCertPathValidationPolicyID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("ReplaceCertPathValidationPolicyAssignment failed: %w", err)
	}

	return nil
}

// GetAssignedCertPathValidationPolicies returns all cert path validation policies assigned to the TLS server.
func (c *Client) GetAssignedCertPathValidationPolicies(ctx context.Context) ([]string, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetAssignedCertPathValidationPolicies"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
	}

	type Response struct {
		XMLName                    xml.Name `xml:"GetAssignedCertPathValidationPoliciesResponse"`
		CertPathValidationPolicyID []string `xml:"CertPathValidationPolicyID"`
	}

	req := Request{Xmlns: advancedSecurityNamespace}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAssignedCertPathValidationPolicies failed: %w", err)
	}

	return resp.CertPathValidationPolicyID, nil
}

// ============================================================
// IEEE 802.1X operations
// ============================================================

// AddAdvSecDot1XConfiguration adds a new 802.1X configuration.
func (c *Client) AddAdvSecDot1XConfiguration(ctx context.Context, config AdvSecDot1XConfiguration) (string, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type stageXML struct {
		Method                     string    `xml:"Method,attr"`
		CertPathValidationPolicyID *string   `xml:"CertPathValidationPolicyID,attr,omitempty"`
		Identity                   *string   `xml:"tas:Identity,omitempty"`
		CertificationPathID        *string   `xml:"tas:CertificationPathID,omitempty"`
		PassphraseID               *string   `xml:"tas:PassphraseID,omitempty"`
		Inner                      *stageXML `xml:"tas:Inner,omitempty"`
	}

	var buildStageXML func(stage AdvSecDot1XStage) stageXML
	buildStageXML = func(stage AdvSecDot1XStage) stageXML {
		s := stageXML{
			Method:                     stage.Method,
			CertPathValidationPolicyID: stage.CertPathValidationPolicyID,
			Identity:                   stage.Identity,
			CertificationPathID:        stage.CertificationPathID,
			PassphraseID:               stage.PassphraseID,
		}

		if stage.Inner != nil {
			inner := buildStageXML(*stage.Inner)
			s.Inner = &inner
		}

		return s
	}

	type configXML struct {
		Dot1XID *string  `xml:"tas:Dot1XID,omitempty"`
		Alias   *string  `xml:"tas:Alias,omitempty"`
		Outer   stageXML `xml:"tas:Outer"`
	}

	type Request struct {
		XMLName              xml.Name  `xml:"tas:AddDot1XConfiguration"`
		Xmlns                string    `xml:"xmlns:tas,attr"`
		Dot1XConfiguration   configXML `xml:"tas:Dot1XConfiguration"`
	}

	type Response struct {
		XMLName xml.Name `xml:"AddDot1XConfigurationResponse"`
		Dot1XID string   `xml:"Dot1XID"`
	}

	var dot1xID *string
	if config.Dot1XID != "" {
		dot1xID = &config.Dot1XID
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		Dot1XConfiguration: configXML{
			Dot1XID: dot1xID,
			Alias:   config.Alias,
			Outer:   buildStageXML(config.Outer),
		},
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("AddAdvSecDot1XConfiguration failed: %w", err)
	}

	return resp.Dot1XID, nil
}

// GetAllAdvSecDot1XConfigurations returns all 802.1X configurations.
func (c *Client) GetAllAdvSecDot1XConfigurations(ctx context.Context) ([]AdvSecDot1XConfiguration, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetAllDot1XConfigurations"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
	}

	type stageXML struct {
		Method                     string    `xml:"Method,attr"`
		CertPathValidationPolicyID *string   `xml:"CertPathValidationPolicyID,attr"`
		Identity                   *string   `xml:"Identity"`
		CertificationPathID        *string   `xml:"CertificationPathID"`
		PassphraseID               *string   `xml:"PassphraseID"`
		Inner                      *stageXML `xml:"Inner"`
	}

	type configXML struct {
		Dot1XID *string  `xml:"Dot1XID"`
		Alias   *string  `xml:"Alias"`
		Outer   stageXML `xml:"Outer"`
	}

	type Response struct {
		XMLName       xml.Name    `xml:"GetAllDot1XConfigurationsResponse"`
		Configuration []configXML `xml:"Configuration"`
	}

	req := Request{Xmlns: advancedSecurityNamespace}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAllAdvSecDot1XConfigurations failed: %w", err)
	}

	var buildStage func(s stageXML) AdvSecDot1XStage
	buildStage = func(s stageXML) AdvSecDot1XStage {
		stage := AdvSecDot1XStage{
			Method:                     s.Method,
			CertPathValidationPolicyID: s.CertPathValidationPolicyID,
			Identity:                   s.Identity,
			CertificationPathID:        s.CertificationPathID,
			PassphraseID:               s.PassphraseID,
		}

		if s.Inner != nil {
			inner := buildStage(*s.Inner)
			stage.Inner = &inner
		}

		return stage
	}

	configs := make([]AdvSecDot1XConfiguration, 0, len(resp.Configuration))
	for _, c := range resp.Configuration {
		var dot1xID string
		if c.Dot1XID != nil {
			dot1xID = *c.Dot1XID
		}

		configs = append(configs, AdvSecDot1XConfiguration{
			Dot1XID: dot1xID,
			Alias:   c.Alias,
			Outer:   buildStage(c.Outer),
		})
	}

	return configs, nil
}

// GetAdvSecDot1XConfiguration returns a specific 802.1X configuration.
func (c *Client) GetAdvSecDot1XConfiguration(ctx context.Context, dot1xID string) (*AdvSecDot1XConfiguration, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetDot1XConfiguration"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
		Dot1XID string   `xml:"tas:Dot1XID"`
	}

	type stageXML struct {
		Method                     string    `xml:"Method,attr"`
		CertPathValidationPolicyID *string   `xml:"CertPathValidationPolicyID,attr"`
		Identity                   *string   `xml:"Identity"`
		CertificationPathID        *string   `xml:"CertificationPathID"`
		PassphraseID               *string   `xml:"PassphraseID"`
		Inner                      *stageXML `xml:"Inner"`
	}

	type configXML struct {
		Dot1XID *string  `xml:"Dot1XID"`
		Alias   *string  `xml:"Alias"`
		Outer   stageXML `xml:"Outer"`
	}

	type Response struct {
		XMLName              xml.Name  `xml:"GetDot1XConfigurationResponse"`
		Dot1XConfiguration   configXML `xml:"Dot1XConfiguration"`
	}

	req := Request{
		Xmlns:   advancedSecurityNamespace,
		Dot1XID: dot1xID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAdvSecDot1XConfiguration failed: %w", err)
	}

	var buildStage func(s stageXML) AdvSecDot1XStage
	buildStage = func(s stageXML) AdvSecDot1XStage {
		stage := AdvSecDot1XStage{
			Method:                     s.Method,
			CertPathValidationPolicyID: s.CertPathValidationPolicyID,
			Identity:                   s.Identity,
			CertificationPathID:        s.CertificationPathID,
			PassphraseID:               s.PassphraseID,
		}

		if s.Inner != nil {
			inner := buildStage(*s.Inner)
			stage.Inner = &inner
		}

		return stage
	}

	var id string
	if resp.Dot1XConfiguration.Dot1XID != nil {
		id = *resp.Dot1XConfiguration.Dot1XID
	}

	return &AdvSecDot1XConfiguration{
		Dot1XID: id,
		Alias:   resp.Dot1XConfiguration.Alias,
		Outer:   buildStage(resp.Dot1XConfiguration.Outer),
	}, nil
}

// DeleteAdvSecDot1XConfiguration deletes an 802.1X configuration.
func (c *Client) DeleteAdvSecDot1XConfiguration(ctx context.Context, dot1xID string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:DeleteDot1XConfiguration"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
		Dot1XID string   `xml:"tas:Dot1XID"`
	}

	type Response struct {
		XMLName xml.Name `xml:"DeleteDot1XConfigurationResponse"`
	}

	req := Request{
		Xmlns:   advancedSecurityNamespace,
		Dot1XID: dot1xID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteAdvSecDot1XConfiguration failed: %w", err)
	}

	return nil
}

// SetNetworkInterfaceAdvSecDot1XConfiguration assigns an 802.1X configuration to a network interface.
func (c *Client) SetNetworkInterfaceAdvSecDot1XConfiguration(ctx context.Context, token string, dot1xID string) (bool, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:SetNetworkInterfaceDot1XConfiguration"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
		Token   string   `xml:"tas:token"`
		Dot1XID string   `xml:"tas:Dot1XID"`
	}

	type Response struct {
		XMLName      xml.Name `xml:"SetNetworkInterfaceDot1XConfigurationResponse"`
		RebootNeeded bool     `xml:"RebootNeeded"`
	}

	req := Request{
		Xmlns:   advancedSecurityNamespace,
		Token:   token,
		Dot1XID: dot1xID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return false, fmt.Errorf("SetNetworkInterfaceAdvSecDot1XConfiguration failed: %w", err)
	}

	return resp.RebootNeeded, nil
}

// GetNetworkInterfaceAdvSecDot1XConfiguration returns the 802.1X configuration for a network interface.
func (c *Client) GetNetworkInterfaceAdvSecDot1XConfiguration(ctx context.Context, token string) (string, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetNetworkInterfaceDot1XConfiguration"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
		Token   string   `xml:"tas:token"`
	}

	type Response struct {
		XMLName xml.Name `xml:"GetNetworkInterfaceDot1XConfigurationResponse"`
		Dot1XID *string  `xml:"Dot1XID"`
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		Token: token,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("GetNetworkInterfaceAdvSecDot1XConfiguration failed: %w", err)
	}

	if resp.Dot1XID == nil {
		return "", nil
	}

	return *resp.Dot1XID, nil
}

// DeleteNetworkInterfaceAdvSecDot1XConfiguration removes the 802.1X configuration from a network interface.
func (c *Client) DeleteNetworkInterfaceAdvSecDot1XConfiguration(ctx context.Context, token string) (bool, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:DeleteNetworkInterfaceDot1XConfiguration"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
		Token   string   `xml:"tas:token"`
	}

	type Response struct {
		XMLName      xml.Name `xml:"DeleteNetworkInterfaceDot1XConfigurationResponse"`
		RebootNeeded bool     `xml:"RebootNeeded"`
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		Token: token,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return false, fmt.Errorf("DeleteNetworkInterfaceAdvSecDot1XConfiguration failed: %w", err)
	}

	return resp.RebootNeeded, nil
}

// ============================================================
// Media Signing operations
// ============================================================

// AddMediaSigningCertificateAssignment assigns a certification path for media signing.
func (c *Client) AddMediaSigningCertificateAssignment(ctx context.Context, certificationPathID string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName             xml.Name `xml:"tas:AddMediaSigningCertificateAssignment"`
		Xmlns               string   `xml:"xmlns:tas,attr"`
		CertificationPathID string   `xml:"tas:CertificationPathID"`
	}

	type Response struct {
		XMLName xml.Name `xml:"AddMediaSigningCertificateAssignmentResponse"`
	}

	req := Request{
		Xmlns:               advancedSecurityNamespace,
		CertificationPathID: certificationPathID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("AddMediaSigningCertificateAssignment failed: %w", err)
	}

	return nil
}

// RemoveMediaSigningCertificateAssignment removes a certification path from media signing.
func (c *Client) RemoveMediaSigningCertificateAssignment(ctx context.Context, certificationPathID string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName             xml.Name `xml:"tas:RemoveMediaSigningCertificateAssignment"`
		Xmlns               string   `xml:"xmlns:tas,attr"`
		CertificationPathID string   `xml:"tas:CertificationPathID"`
	}

	type Response struct {
		XMLName xml.Name `xml:"RemoveMediaSigningCertificateAssignmentResponse"`
	}

	req := Request{
		Xmlns:               advancedSecurityNamespace,
		CertificationPathID: certificationPathID,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("RemoveMediaSigningCertificateAssignment failed: %w", err)
	}

	return nil
}

// GetAssignedMediaSigningCertificates returns all certification paths assigned for media signing.
func (c *Client) GetAssignedMediaSigningCertificates(ctx context.Context) ([]string, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetAssignedMediaSigningCertificates"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
	}

	type Response struct {
		XMLName              xml.Name `xml:"GetAssignedMediaSigningCertificatesResponse"`
		CertificationPathID []string `xml:"CertificationPathID"`
	}

	req := Request{Xmlns: advancedSecurityNamespace}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAssignedMediaSigningCertificates failed: %w", err)
	}

	return resp.CertificationPathID, nil
}

// ============================================================
// Authorization Server operations
// ============================================================

// GetAuthorizationServerConfigurations returns authorization server configurations.
func (c *Client) GetAuthorizationServerConfigurations(ctx context.Context, token *string) ([]AuthorizationServerConfiguration, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:GetAuthorizationServerConfigurations"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
		Token   *string  `xml:"tas:Token,omitempty"`
	}

	type dataXML struct {
		Type                       string  `xml:"Type,attr"`
		ClientAuth                 *string `xml:"ClientAuth,attr"`
		ServerURI                  string  `xml:"ServerUri"`
		ClientID                   *string `xml:"ClientID"`
		ClientSecret               *string `xml:"ClientSecret"`
		Scope                      *string `xml:"Scope"`
		KeyID                      *string `xml:"KeyID"`
		CertificateID              *string `xml:"CertificateID"`
		CertPathValidationPolicyID *string `xml:"CertPathValidationPolicyID"`
	}

	type configXML struct {
		Token string  `xml:"token,attr"`
		Data  dataXML `xml:"Data"`
	}

	type Response struct {
		XMLName       xml.Name    `xml:"GetAuthorizationServerConfigurationsResponse"`
		Configuration []configXML `xml:"Configuration"`
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		Token: token,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAuthorizationServerConfigurations failed: %w", err)
	}

	configs := make([]AuthorizationServerConfiguration, 0, len(resp.Configuration))
	for _, cfg := range resp.Configuration {
		configs = append(configs, AuthorizationServerConfiguration{
			Token: cfg.Token,
			Data: AuthorizationServerConfigurationData{
				Type:                       cfg.Data.Type,
				ClientAuth:                 cfg.Data.ClientAuth,
				ServerURI:                  cfg.Data.ServerURI,
				ClientID:                   cfg.Data.ClientID,
				ClientSecret:               cfg.Data.ClientSecret,
				Scope:                      cfg.Data.Scope,
				KeyID:                      cfg.Data.KeyID,
				CertificateID:              cfg.Data.CertificateID,
				CertPathValidationPolicyID: cfg.Data.CertPathValidationPolicyID,
			},
		})
	}

	return configs, nil
}

// CreateAuthorizationServerConfiguration creates a new authorization server configuration.
func (c *Client) CreateAuthorizationServerConfiguration(ctx context.Context, config AuthorizationServerConfigurationData) (string, error) {
	endpoint := c.getAdvancedSecurityEndpoint()

	type dataXML struct {
		Type                       string  `xml:"Type,attr"`
		ClientAuth                 *string `xml:"ClientAuth,attr,omitempty"`
		ServerURI                  string  `xml:"tas:ServerUri"`
		ClientID                   *string `xml:"tas:ClientID,omitempty"`
		ClientSecret               *string `xml:"tas:ClientSecret,omitempty"`
		Scope                      *string `xml:"tas:Scope,omitempty"`
		KeyID                      *string `xml:"tas:KeyID,omitempty"`
		CertificateID              *string `xml:"tas:CertificateID,omitempty"`
		CertPathValidationPolicyID *string `xml:"tas:CertPathValidationPolicyID,omitempty"`
	}

	type Request struct {
		XMLName       xml.Name `xml:"tas:CreateAuthorizationServerConfiguration"`
		Xmlns         string   `xml:"xmlns:tas,attr"`
		Configuration dataXML  `xml:"tas:Configuration"`
	}

	type Response struct {
		XMLName xml.Name `xml:"CreateAuthorizationServerConfigurationResponse"`
		Token   string   `xml:"Token"`
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		Configuration: dataXML{
			Type:                       config.Type,
			ClientAuth:                 config.ClientAuth,
			ServerURI:                  config.ServerURI,
			ClientID:                   config.ClientID,
			ClientSecret:               config.ClientSecret,
			Scope:                      config.Scope,
			KeyID:                      config.KeyID,
			CertificateID:              config.CertificateID,
			CertPathValidationPolicyID: config.CertPathValidationPolicyID,
		},
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateAuthorizationServerConfiguration failed: %w", err)
	}

	return resp.Token, nil
}

// SetAuthorizationServerConfiguration updates an authorization server configuration.
func (c *Client) SetAuthorizationServerConfiguration(ctx context.Context, config AuthorizationServerConfiguration) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type dataXML struct {
		Type                       string  `xml:"Type,attr"`
		ClientAuth                 *string `xml:"ClientAuth,attr,omitempty"`
		ServerURI                  string  `xml:"tas:ServerUri"`
		ClientID                   *string `xml:"tas:ClientID,omitempty"`
		ClientSecret               *string `xml:"tas:ClientSecret,omitempty"`
		Scope                      *string `xml:"tas:Scope,omitempty"`
		KeyID                      *string `xml:"tas:KeyID,omitempty"`
		CertificateID              *string `xml:"tas:CertificateID,omitempty"`
		CertPathValidationPolicyID *string `xml:"tas:CertPathValidationPolicyID,omitempty"`
	}

	type configXML struct {
		Token string  `xml:"token,attr"`
		Data  dataXML `xml:"tas:Data"`
	}

	type Request struct {
		XMLName       xml.Name  `xml:"tas:SetAuthorizationServerConfiguration"`
		Xmlns         string    `xml:"xmlns:tas,attr"`
		Configuration configXML `xml:"tas:Configuration"`
	}

	type Response struct {
		XMLName xml.Name `xml:"SetAuthorizationServerConfigurationResponse"`
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		Configuration: configXML{
			Token: config.Token,
			Data: dataXML{
				Type:                       config.Data.Type,
				ClientAuth:                 config.Data.ClientAuth,
				ServerURI:                  config.Data.ServerURI,
				ClientID:                   config.Data.ClientID,
				ClientSecret:               config.Data.ClientSecret,
				Scope:                      config.Data.Scope,
				KeyID:                      config.Data.KeyID,
				CertificateID:              config.Data.CertificateID,
				CertPathValidationPolicyID: config.Data.CertPathValidationPolicyID,
			},
		},
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetAuthorizationServerConfiguration failed: %w", err)
	}

	return nil
}

// DeleteAuthorizationServerConfiguration deletes an authorization server configuration.
func (c *Client) DeleteAuthorizationServerConfiguration(ctx context.Context, token string) error {
	endpoint := c.getAdvancedSecurityEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tas:DeleteAuthorizationServerConfiguration"`
		Xmlns   string   `xml:"xmlns:tas,attr"`
		Token   string   `xml:"tas:Token"`
	}

	type Response struct {
		XMLName xml.Name `xml:"DeleteAuthorizationServerConfigurationResponse"`
	}

	req := Request{
		Xmlns: advancedSecurityNamespace,
		Token: token,
	}

	var resp Response

	if err := c.newAdvancedSecuritySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteAuthorizationServerConfiguration failed: %w", err)
	}

	return nil
}
