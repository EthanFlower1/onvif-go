package onvif

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"time"

	"github.com/EthanFlower1/onvif-go/internal/soap"
)

// Credential service namespace.
const credentialNamespace = "http://www.onvif.org/ver10/credential/wsdl"

// Credential service errors.
var (
	// ErrInvalidCredentialToken is returned when a credential token is empty.
	ErrInvalidCredentialToken = errors.New("invalid credential token: cannot be empty")
	// ErrCredentialNil is returned when a credential is nil.
	ErrCredentialNil = errors.New("credential cannot be nil")
	// ErrInvalidCredentialIdentifierTypeName is returned when an identifier type name is empty.
	ErrInvalidCredentialIdentifierTypeName = errors.New("credential identifier type name cannot be empty")
	// ErrInvalidAccessProfileToken is returned when an access profile token is empty.
	ErrInvalidAccessProfileToken = errors.New("invalid access profile token: cannot be empty")
)

// getCredentialEndpoint returns the credential endpoint, falling back to device endpoint.
func (c *Client) getCredentialEndpoint() string {
	if c.credentialEndpoint != "" {
		return c.credentialEndpoint
	}

	return c.endpoint
}

// GetCredentialServiceCapabilities retrieves the capabilities of the credential service.
func (c *Client) GetCredentialServiceCapabilities(ctx context.Context) (*CredentialServiceCapabilities, error) {
	endpoint := c.getCredentialEndpoint()

	type GetServiceCapabilities struct {
		XMLName xml.Name `xml:"tcr:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:tcr,attr"`
	}

	type GetServiceCapabilitiesResponse struct {
		XMLName      xml.Name `xml:"GetServiceCapabilitiesResponse"`
		Capabilities struct {
			MaxLimit                                uint   `xml:"MaxLimit,attr"`
			MaxCredentials                          uint   `xml:"MaxCredentials,attr"`
			MaxAccessProfilesPerCredential          uint   `xml:"MaxAccessProfilesPerCredential,attr"`
			CredentialValiditySupported             bool   `xml:"CredentialValiditySupported,attr"`
			CredentialAccessProfileValiditySupported bool   `xml:"CredentialAccessProfileValiditySupported,attr"`
			ValiditySupportsTimeValue               bool   `xml:"ValiditySupportsTimeValue,attr"`
			ResetAntipassbackSupported              bool   `xml:"ResetAntipassbackSupported,attr"`
			ClientSuppliedTokenSupported            bool   `xml:"ClientSuppliedTokenSupported,attr"`
			DefaultCredentialSuspensionDuration     string `xml:"DefaultCredentialSuspensionDuration,attr"`
			MaxWhitelistedItems                     uint   `xml:"MaxWhitelistedItems,attr"`
			MaxBlacklistedItems                     uint   `xml:"MaxBlacklistedItems,attr"`
			SupportedIdentifierType                 []struct {
				Value string `xml:",chardata"`
			} `xml:"SupportedIdentifierType"`
			Extension struct {
				SupportedExemptionType []struct {
					Value string `xml:",chardata"`
				} `xml:"SupportedExemptionType"`
			} `xml:"Extension"`
		} `xml:"Capabilities"`
	}

	req := GetServiceCapabilities{Xmlns: credentialNamespace}

	var resp GetServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetCredentialServiceCapabilities failed: %w", err)
	}

	caps := &CredentialServiceCapabilities{
		MaxLimit:                                resp.Capabilities.MaxLimit,
		MaxCredentials:                          resp.Capabilities.MaxCredentials,
		MaxAccessProfilesPerCredential:          resp.Capabilities.MaxAccessProfilesPerCredential,
		CredentialValiditySupported:             resp.Capabilities.CredentialValiditySupported,
		CredentialAccessProfileValiditySupported: resp.Capabilities.CredentialAccessProfileValiditySupported,
		ValiditySupportsTimeValue:               resp.Capabilities.ValiditySupportsTimeValue,
		ResetAntipassbackSupported:              resp.Capabilities.ResetAntipassbackSupported,
		ClientSuppliedTokenSupported:            resp.Capabilities.ClientSuppliedTokenSupported,
		DefaultCredentialSuspensionDuration:     resp.Capabilities.DefaultCredentialSuspensionDuration,
		MaxWhitelistedItems:                     resp.Capabilities.MaxWhitelistedItems,
		MaxBlacklistedItems:                     resp.Capabilities.MaxBlacklistedItems,
	}

	for _, sit := range resp.Capabilities.SupportedIdentifierType {
		caps.SupportedIdentifierTypes = append(caps.SupportedIdentifierTypes, sit.Value)
	}

	for _, set := range resp.Capabilities.Extension.SupportedExemptionType {
		caps.SupportedExemptionTypes = append(caps.SupportedExemptionTypes, set.Value)
	}

	return caps, nil
}

// GetSupportedFormatTypes returns all supported format types for a specified identifier type.
func (c *Client) GetSupportedFormatTypes(ctx context.Context, credentialIdentifierTypeName string) ([]*CredentialIdentifierFormatTypeInfo, error) {
	endpoint := c.getCredentialEndpoint()

	type GetSupportedFormatTypes struct {
		XMLName                      xml.Name `xml:"tcr:GetSupportedFormatTypes"`
		Xmlns                        string   `xml:"xmlns:tcr,attr"`
		CredentialIdentifierTypeName string   `xml:"tcr:CredentialIdentifierTypeName"`
	}

	type FormatTypeInfoEntry struct {
		FormatType  string `xml:"FormatType"`
		Description string `xml:"Description"`
	}

	type GetSupportedFormatTypesResponse struct {
		XMLName        xml.Name              `xml:"GetSupportedFormatTypesResponse"`
		FormatTypeInfo []FormatTypeInfoEntry `xml:"FormatTypeInfo"`
	}

	req := GetSupportedFormatTypes{
		Xmlns:                        credentialNamespace,
		CredentialIdentifierTypeName: credentialIdentifierTypeName,
	}

	var resp GetSupportedFormatTypesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetSupportedFormatTypes failed: %w", err)
	}

	result := make([]*CredentialIdentifierFormatTypeInfo, 0, len(resp.FormatTypeInfo))

	for _, entry := range resp.FormatTypeInfo {
		result = append(result, &CredentialIdentifierFormatTypeInfo{
			FormatType:  entry.FormatType,
			Description: entry.Description,
		})
	}

	return result, nil
}

// credentialInfoEntryToPublic converts an internal XML struct to public CredentialInfo.
func credentialInfoEntryToPublic(token, description, holderRef string, validFrom, validTo string) CredentialInfo {
	info := CredentialInfo{
		Token:                    token,
		Description:              description,
		CredentialHolderReference: holderRef,
	}

	if validFrom != "" {
		if t, err := time.Parse(time.RFC3339, validFrom); err == nil {
			info.ValidFrom = &t
		}
	}

	if validTo != "" {
		if t, err := time.Parse(time.RFC3339, validTo); err == nil {
			info.ValidTo = &t
		}
	}

	return info
}

// GetCredentialInfo retrieves CredentialInfo items by token.
func (c *Client) GetCredentialInfo(ctx context.Context, tokens []string) ([]*CredentialInfo, error) {
	if len(tokens) == 0 {
		return nil, ErrInvalidCredentialToken
	}

	endpoint := c.getCredentialEndpoint()

	type GetCredentialInfo struct {
		XMLName xml.Name `xml:"tcr:GetCredentialInfo"`
		Xmlns   string   `xml:"xmlns:tcr,attr"`
		Token   []string `xml:"tcr:Token"`
	}

	type CredentialInfoEntry struct {
		Token                    string `xml:"token,attr"`
		Description              string `xml:"Description"`
		CredentialHolderReference string `xml:"CredentialHolderReference"`
		ValidFrom                string `xml:"ValidFrom"`
		ValidTo                  string `xml:"ValidTo"`
	}

	type GetCredentialInfoResponse struct {
		XMLName        xml.Name              `xml:"GetCredentialInfoResponse"`
		CredentialInfo []CredentialInfoEntry `xml:"CredentialInfo"`
	}

	req := GetCredentialInfo{
		Xmlns: credentialNamespace,
		Token: tokens,
	}

	var resp GetCredentialInfoResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetCredentialInfo failed: %w", err)
	}

	result := make([]*CredentialInfo, 0, len(resp.CredentialInfo))

	for _, entry := range resp.CredentialInfo {
		info := credentialInfoEntryToPublic(entry.Token, entry.Description, entry.CredentialHolderReference, entry.ValidFrom, entry.ValidTo)
		result = append(result, &info)
	}

	return result, nil
}

// GetCredentialInfoList retrieves a paginated list of all CredentialInfo items.
func (c *Client) GetCredentialInfoList(ctx context.Context, limit *int, startReference *string) ([]*CredentialInfo, string, error) {
	endpoint := c.getCredentialEndpoint()

	type GetCredentialInfoList struct {
		XMLName        xml.Name `xml:"tcr:GetCredentialInfoList"`
		Xmlns          string   `xml:"xmlns:tcr,attr"`
		Limit          *int     `xml:"tcr:Limit,omitempty"`
		StartReference *string  `xml:"tcr:StartReference,omitempty"`
	}

	type CredentialInfoEntry struct {
		Token                    string `xml:"token,attr"`
		Description              string `xml:"Description"`
		CredentialHolderReference string `xml:"CredentialHolderReference"`
		ValidFrom                string `xml:"ValidFrom"`
		ValidTo                  string `xml:"ValidTo"`
	}

	type GetCredentialInfoListResponse struct {
		XMLName            xml.Name              `xml:"GetCredentialInfoListResponse"`
		NextStartReference string                `xml:"NextStartReference"`
		CredentialInfo     []CredentialInfoEntry `xml:"CredentialInfo"`
	}

	req := GetCredentialInfoList{
		Xmlns:          credentialNamespace,
		Limit:          limit,
		StartReference: startReference,
	}

	var resp GetCredentialInfoListResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetCredentialInfoList failed: %w", err)
	}

	result := make([]*CredentialInfo, 0, len(resp.CredentialInfo))

	for _, entry := range resp.CredentialInfo {
		info := credentialInfoEntryToPublic(entry.Token, entry.Description, entry.CredentialHolderReference, entry.ValidFrom, entry.ValidTo)
		result = append(result, &info)
	}

	return result, resp.NextStartReference, nil
}

// credentialEntryToPublic maps an internal XML credential entry to a public Credential.
func credentialEntryToPublic(
	token, description, holderRef, validFrom, validTo string,
	identifiers []struct {
		Type struct {
			Name       string `xml:"Name"`
			FormatType string `xml:"FormatType"`
		} `xml:"Type"`
		ExemptedFromAuthentication bool   `xml:"ExemptedFromAuthentication"`
		Value                      []byte `xml:"Value"`
	},
	accessProfiles []struct {
		AccessProfileToken string `xml:"AccessProfileToken"`
		ValidFrom          string `xml:"ValidFrom"`
		ValidTo            string `xml:"ValidTo"`
	},
	extendedGrantTime *bool,
	attributes []struct {
		Name  string `xml:"Name,attr"`
		Value string `xml:"Value,attr"`
	},
) Credential {
	info := credentialInfoEntryToPublic(token, description, holderRef, validFrom, validTo)
	cred := Credential{CredentialInfo: info, ExtendedGrantTime: extendedGrantTime}

	for _, id := range identifiers {
		cred.CredentialIdentifiers = append(cred.CredentialIdentifiers, CredentialIdentifier{
			Type: CredentialIdentifierType{
				Name:       id.Type.Name,
				FormatType: id.Type.FormatType,
			},
			ExemptedFromAuthentication: id.ExemptedFromAuthentication,
			Value:                      id.Value,
		})
	}

	for _, ap := range accessProfiles {
		capEntry := CredentialAccessProfile{AccessProfileToken: ap.AccessProfileToken}

		if ap.ValidFrom != "" {
			if t, err := time.Parse(time.RFC3339, ap.ValidFrom); err == nil {
				capEntry.ValidFrom = &t
			}
		}

		if ap.ValidTo != "" {
			if t, err := time.Parse(time.RFC3339, ap.ValidTo); err == nil {
				capEntry.ValidTo = &t
			}
		}

		cred.CredentialAccessProfiles = append(cred.CredentialAccessProfiles, capEntry)
	}

	for _, attr := range attributes {
		cred.Attributes = append(cred.Attributes, CredentialAttribute{
			Name:  attr.Name,
			Value: attr.Value,
		})
	}

	return cred
}

// GetCredentialsByTokens retrieves Credential items by their tokens.
func (c *Client) GetCredentialsByTokens(ctx context.Context, tokens []string) ([]*Credential, error) {
	if len(tokens) == 0 {
		return nil, ErrInvalidCredentialToken
	}

	endpoint := c.getCredentialEndpoint()

	type GetCredentials struct {
		XMLName xml.Name `xml:"tcr:GetCredentials"`
		Xmlns   string   `xml:"xmlns:tcr,attr"`
		Token   []string `xml:"tcr:Token"`
	}

	type CredentialEntry struct {
		Token                    string `xml:"token,attr"`
		Description              string `xml:"Description"`
		CredentialHolderReference string `xml:"CredentialHolderReference"`
		ValidFrom                string `xml:"ValidFrom"`
		ValidTo                  string `xml:"ValidTo"`
		CredentialIdentifier     []struct {
			Type struct {
				Name       string `xml:"Name"`
				FormatType string `xml:"FormatType"`
			} `xml:"Type"`
			ExemptedFromAuthentication bool   `xml:"ExemptedFromAuthentication"`
			Value                      []byte `xml:"Value"`
		} `xml:"CredentialIdentifier"`
		CredentialAccessProfile []struct {
			AccessProfileToken string `xml:"AccessProfileToken"`
			ValidFrom          string `xml:"ValidFrom"`
			ValidTo            string `xml:"ValidTo"`
		} `xml:"CredentialAccessProfile"`
		ExtendedGrantTime *bool `xml:"ExtendedGrantTime"`
		Attribute         []struct {
			Name  string `xml:"Name,attr"`
			Value string `xml:"Value,attr"`
		} `xml:"Attribute"`
	}

	type GetCredentialsResponse struct {
		XMLName    xml.Name          `xml:"GetCredentialsResponse"`
		Credential []CredentialEntry `xml:"Credential"`
	}

	req := GetCredentials{
		Xmlns: credentialNamespace,
		Token: tokens,
	}

	var resp GetCredentialsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetCredentials failed: %w", err)
	}

	result := make([]*Credential, 0, len(resp.Credential))

	for _, entry := range resp.Credential {
		cred := credentialEntryToPublic(
			entry.Token, entry.Description, entry.CredentialHolderReference,
			entry.ValidFrom, entry.ValidTo,
			entry.CredentialIdentifier,
			entry.CredentialAccessProfile,
			entry.ExtendedGrantTime,
			entry.Attribute,
		)
		result = append(result, &cred)
	}

	return result, nil
}

// GetCredentialList retrieves a paginated list of all Credential items.
func (c *Client) GetCredentialList(ctx context.Context, limit *int, startReference *string) ([]*Credential, string, error) {
	endpoint := c.getCredentialEndpoint()

	type GetCredentialList struct {
		XMLName        xml.Name `xml:"tcr:GetCredentialList"`
		Xmlns          string   `xml:"xmlns:tcr,attr"`
		Limit          *int     `xml:"tcr:Limit,omitempty"`
		StartReference *string  `xml:"tcr:StartReference,omitempty"`
	}

	type CredentialEntry struct {
		Token                    string `xml:"token,attr"`
		Description              string `xml:"Description"`
		CredentialHolderReference string `xml:"CredentialHolderReference"`
		ValidFrom                string `xml:"ValidFrom"`
		ValidTo                  string `xml:"ValidTo"`
		CredentialIdentifier     []struct {
			Type struct {
				Name       string `xml:"Name"`
				FormatType string `xml:"FormatType"`
			} `xml:"Type"`
			ExemptedFromAuthentication bool   `xml:"ExemptedFromAuthentication"`
			Value                      []byte `xml:"Value"`
		} `xml:"CredentialIdentifier"`
		CredentialAccessProfile []struct {
			AccessProfileToken string `xml:"AccessProfileToken"`
			ValidFrom          string `xml:"ValidFrom"`
			ValidTo            string `xml:"ValidTo"`
		} `xml:"CredentialAccessProfile"`
		ExtendedGrantTime *bool `xml:"ExtendedGrantTime"`
		Attribute         []struct {
			Name  string `xml:"Name,attr"`
			Value string `xml:"Value,attr"`
		} `xml:"Attribute"`
	}

	type GetCredentialListResponse struct {
		XMLName            xml.Name          `xml:"GetCredentialListResponse"`
		NextStartReference string            `xml:"NextStartReference"`
		Credential         []CredentialEntry `xml:"Credential"`
	}

	req := GetCredentialList{
		Xmlns:          credentialNamespace,
		Limit:          limit,
		StartReference: startReference,
	}

	var resp GetCredentialListResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetCredentialList failed: %w", err)
	}

	result := make([]*Credential, 0, len(resp.Credential))

	for _, entry := range resp.Credential {
		cred := credentialEntryToPublic(
			entry.Token, entry.Description, entry.CredentialHolderReference,
			entry.ValidFrom, entry.ValidTo,
			entry.CredentialIdentifier,
			entry.CredentialAccessProfile,
			entry.ExtendedGrantTime,
			entry.Attribute,
		)
		result = append(result, &cred)
	}

	return result, resp.NextStartReference, nil
}

// CreateCredential creates a new credential and returns the allocated token.
func (c *Client) CreateCredential(ctx context.Context, credential Credential, state CredentialState) (string, error) {
	endpoint := c.getCredentialEndpoint()

	type CredentialIdentifierXML struct {
		Type struct {
			Name       string `xml:"tcr:Name"`
			FormatType string `xml:"tcr:FormatType"`
		} `xml:"tcr:Type"`
		ExemptedFromAuthentication bool   `xml:"tcr:ExemptedFromAuthentication"`
		Value                      []byte `xml:"tcr:Value"`
	}

	type CredentialAccessProfileXML struct {
		AccessProfileToken string `xml:"tcr:AccessProfileToken"`
	}

	type CredentialXML struct {
		Token                    string                       `xml:"token,attr,omitempty"`
		Description              string                       `xml:"tcr:Description,omitempty"`
		CredentialHolderReference string                      `xml:"tcr:CredentialHolderReference"`
		CredentialIdentifier     []CredentialIdentifierXML   `xml:"tcr:CredentialIdentifier"`
		CredentialAccessProfile  []CredentialAccessProfileXML `xml:"tcr:CredentialAccessProfile,omitempty"`
	}

	type CredentialStateXML struct {
		Enabled bool   `xml:"tcr:Enabled"`
		Reason  string `xml:"tcr:Reason,omitempty"`
	}

	type CreateCredential struct {
		XMLName    xml.Name           `xml:"tcr:CreateCredential"`
		Xmlns      string             `xml:"xmlns:tcr,attr"`
		Credential CredentialXML      `xml:"tcr:Credential"`
		State      CredentialStateXML `xml:"tcr:State"`
	}

	type CreateCredentialResponse struct {
		XMLName xml.Name `xml:"CreateCredentialResponse"`
		Token   string   `xml:"Token"`
	}

	credXML := CredentialXML{
		Description:              credential.Description,
		CredentialHolderReference: credential.CredentialHolderReference,
	}

	for _, id := range credential.CredentialIdentifiers {
		ci := CredentialIdentifierXML{Value: id.Value, ExemptedFromAuthentication: id.ExemptedFromAuthentication}
		ci.Type.Name = id.Type.Name
		ci.Type.FormatType = id.Type.FormatType
		credXML.CredentialIdentifier = append(credXML.CredentialIdentifier, ci)
	}

	for _, ap := range credential.CredentialAccessProfiles {
		credXML.CredentialAccessProfile = append(credXML.CredentialAccessProfile, CredentialAccessProfileXML{
			AccessProfileToken: ap.AccessProfileToken,
		})
	}

	stateXML := CredentialStateXML{
		Enabled: state.Enabled,
		Reason:  state.Reason,
	}

	req := CreateCredential{
		Xmlns:      credentialNamespace,
		Credential: credXML,
		State:      stateXML,
	}

	var resp CreateCredentialResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateCredential failed: %w", err)
	}

	return resp.Token, nil
}

// ModifyCredential modifies an existing credential.
func (c *Client) ModifyCredential(ctx context.Context, credential Credential) error {
	if credential.Token == "" {
		return ErrInvalidCredentialToken
	}

	endpoint := c.getCredentialEndpoint()

	type CredentialIdentifierXML struct {
		Type struct {
			Name       string `xml:"tcr:Name"`
			FormatType string `xml:"tcr:FormatType"`
		} `xml:"tcr:Type"`
		ExemptedFromAuthentication bool   `xml:"tcr:ExemptedFromAuthentication"`
		Value                      []byte `xml:"tcr:Value"`
	}

	type CredentialAccessProfileXML struct {
		AccessProfileToken string `xml:"tcr:AccessProfileToken"`
	}

	type CredentialXML struct {
		Token                    string                       `xml:"token,attr"`
		Description              string                       `xml:"tcr:Description,omitempty"`
		CredentialHolderReference string                      `xml:"tcr:CredentialHolderReference"`
		CredentialIdentifier     []CredentialIdentifierXML   `xml:"tcr:CredentialIdentifier"`
		CredentialAccessProfile  []CredentialAccessProfileXML `xml:"tcr:CredentialAccessProfile,omitempty"`
	}

	type ModifyCredential struct {
		XMLName    xml.Name      `xml:"tcr:ModifyCredential"`
		Xmlns      string        `xml:"xmlns:tcr,attr"`
		Credential CredentialXML `xml:"tcr:Credential"`
	}

	credXML := CredentialXML{
		Token:                    credential.Token,
		Description:              credential.Description,
		CredentialHolderReference: credential.CredentialHolderReference,
	}

	for _, id := range credential.CredentialIdentifiers {
		ci := CredentialIdentifierXML{Value: id.Value, ExemptedFromAuthentication: id.ExemptedFromAuthentication}
		ci.Type.Name = id.Type.Name
		ci.Type.FormatType = id.Type.FormatType
		credXML.CredentialIdentifier = append(credXML.CredentialIdentifier, ci)
	}

	for _, ap := range credential.CredentialAccessProfiles {
		credXML.CredentialAccessProfile = append(credXML.CredentialAccessProfile, CredentialAccessProfileXML{
			AccessProfileToken: ap.AccessProfileToken,
		})
	}

	req := ModifyCredential{
		Xmlns:      credentialNamespace,
		Credential: credXML,
	}

	var resp struct {
		XMLName xml.Name `xml:"ModifyCredentialResponse"`
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("ModifyCredential failed: %w", err)
	}

	return nil
}

// SetCredential synchronizes a credential (upsert with state).
func (c *Client) SetCredential(ctx context.Context, data CredentialData) error {
	endpoint := c.getCredentialEndpoint()

	type CredentialIdentifierXML struct {
		Type struct {
			Name       string `xml:"tcr:Name"`
			FormatType string `xml:"tcr:FormatType"`
		} `xml:"tcr:Type"`
		ExemptedFromAuthentication bool   `xml:"tcr:ExemptedFromAuthentication"`
		Value                      []byte `xml:"tcr:Value"`
	}

	type CredentialAccessProfileXML struct {
		AccessProfileToken string `xml:"tcr:AccessProfileToken"`
	}

	type CredentialXML struct {
		Token                    string                       `xml:"token,attr,omitempty"`
		Description              string                       `xml:"tcr:Description,omitempty"`
		CredentialHolderReference string                      `xml:"tcr:CredentialHolderReference"`
		CredentialIdentifier     []CredentialIdentifierXML   `xml:"tcr:CredentialIdentifier"`
		CredentialAccessProfile  []CredentialAccessProfileXML `xml:"tcr:CredentialAccessProfile,omitempty"`
	}

	type CredentialStateXML struct {
		Enabled bool   `xml:"tcr:Enabled"`
		Reason  string `xml:"tcr:Reason,omitempty"`
	}

	type CredentialDataXML struct {
		Credential      CredentialXML      `xml:"tcr:Credential"`
		CredentialState CredentialStateXML `xml:"tcr:CredentialState"`
	}

	type SetCredential struct {
		XMLName        xml.Name          `xml:"tcr:SetCredential"`
		Xmlns          string            `xml:"xmlns:tcr,attr"`
		CredentialData CredentialDataXML `xml:"tcr:CredentialData"`
	}

	credXML := CredentialXML{
		Token:                    data.Credential.Token,
		Description:              data.Credential.Description,
		CredentialHolderReference: data.Credential.CredentialHolderReference,
	}

	for _, id := range data.Credential.CredentialIdentifiers {
		ci := CredentialIdentifierXML{Value: id.Value, ExemptedFromAuthentication: id.ExemptedFromAuthentication}
		ci.Type.Name = id.Type.Name
		ci.Type.FormatType = id.Type.FormatType
		credXML.CredentialIdentifier = append(credXML.CredentialIdentifier, ci)
	}

	for _, ap := range data.Credential.CredentialAccessProfiles {
		credXML.CredentialAccessProfile = append(credXML.CredentialAccessProfile, CredentialAccessProfileXML{
			AccessProfileToken: ap.AccessProfileToken,
		})
	}

	req := SetCredential{
		Xmlns: credentialNamespace,
		CredentialData: CredentialDataXML{
			Credential: credXML,
			CredentialState: CredentialStateXML{
				Enabled: data.CredentialState.Enabled,
				Reason:  data.CredentialState.Reason,
			},
		},
	}

	var resp struct {
		XMLName xml.Name `xml:"SetCredentialResponse"`
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetCredential failed: %w", err)
	}

	return nil
}

// DeleteCredential deletes the credential with the given token.
func (c *Client) DeleteCredential(ctx context.Context, token string) error {
	if token == "" {
		return ErrInvalidCredentialToken
	}

	endpoint := c.getCredentialEndpoint()

	type DeleteCredential struct {
		XMLName xml.Name `xml:"tcr:DeleteCredential"`
		Xmlns   string   `xml:"xmlns:tcr,attr"`
		Token   string   `xml:"tcr:Token"`
	}

	req := DeleteCredential{
		Xmlns: credentialNamespace,
		Token: token,
	}

	var resp struct {
		XMLName xml.Name `xml:"DeleteCredentialResponse"`
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteCredential failed: %w", err)
	}

	return nil
}

// GetCredentialState returns the state of the specified credential.
func (c *Client) GetCredentialState(ctx context.Context, token string) (*CredentialState, error) {
	if token == "" {
		return nil, ErrInvalidCredentialToken
	}

	endpoint := c.getCredentialEndpoint()

	type GetCredentialState struct {
		XMLName xml.Name `xml:"tcr:GetCredentialState"`
		Xmlns   string   `xml:"xmlns:tcr,attr"`
		Token   string   `xml:"tcr:Token"`
	}

	type GetCredentialStateResponse struct {
		XMLName xml.Name `xml:"GetCredentialStateResponse"`
		State   struct {
			Enabled          bool   `xml:"Enabled"`
			Reason           string `xml:"Reason"`
			AntipassbackState *struct {
				AntipassbackViolated bool `xml:"AntipassbackViolated"`
			} `xml:"AntipassbackState"`
		} `xml:"State"`
	}

	req := GetCredentialState{
		Xmlns: credentialNamespace,
		Token: token,
	}

	var resp GetCredentialStateResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetCredentialState failed: %w", err)
	}

	state := &CredentialState{
		Enabled: resp.State.Enabled,
		Reason:  resp.State.Reason,
	}

	if resp.State.AntipassbackState != nil {
		state.AntipassbackState = &AntipassbackState{
			AntipassbackViolated: resp.State.AntipassbackState.AntipassbackViolated,
		}
	}

	return state, nil
}

// EnableCredential enables the specified credential.
func (c *Client) EnableCredential(ctx context.Context, token string, reason *string) error {
	if token == "" {
		return ErrInvalidCredentialToken
	}

	endpoint := c.getCredentialEndpoint()

	type EnableCredential struct {
		XMLName xml.Name `xml:"tcr:EnableCredential"`
		Xmlns   string   `xml:"xmlns:tcr,attr"`
		Token   string   `xml:"tcr:Token"`
		Reason  *string  `xml:"tcr:Reason,omitempty"`
	}

	req := EnableCredential{
		Xmlns:  credentialNamespace,
		Token:  token,
		Reason: reason,
	}

	var resp struct {
		XMLName xml.Name `xml:"EnableCredentialResponse"`
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("EnableCredential failed: %w", err)
	}

	return nil
}

// DisableCredential disables the specified credential.
func (c *Client) DisableCredential(ctx context.Context, token string, reason *string) error {
	if token == "" {
		return ErrInvalidCredentialToken
	}

	endpoint := c.getCredentialEndpoint()

	type DisableCredential struct {
		XMLName xml.Name `xml:"tcr:DisableCredential"`
		Xmlns   string   `xml:"xmlns:tcr,attr"`
		Token   string   `xml:"tcr:Token"`
		Reason  *string  `xml:"tcr:Reason,omitempty"`
	}

	req := DisableCredential{
		Xmlns:  credentialNamespace,
		Token:  token,
		Reason: reason,
	}

	var resp struct {
		XMLName xml.Name `xml:"DisableCredentialResponse"`
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DisableCredential failed: %w", err)
	}

	return nil
}

// ResetAntipassbackViolation resets anti-passback violations for the specified credential.
func (c *Client) ResetAntipassbackViolation(ctx context.Context, credentialToken string) error {
	if credentialToken == "" {
		return ErrInvalidCredentialToken
	}

	endpoint := c.getCredentialEndpoint()

	type ResetAntipassbackViolation struct {
		XMLName         xml.Name `xml:"tcr:ResetAntipassbackViolation"`
		Xmlns           string   `xml:"xmlns:tcr,attr"`
		CredentialToken string   `xml:"tcr:CredentialToken"`
	}

	req := ResetAntipassbackViolation{
		Xmlns:           credentialNamespace,
		CredentialToken: credentialToken,
	}

	var resp struct {
		XMLName xml.Name `xml:"ResetAntipassbackViolationResponse"`
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("ResetAntipassbackViolation failed: %w", err)
	}

	return nil
}

// GetCredentialIdentifiers returns all credential identifiers for the specified credential.
func (c *Client) GetCredentialIdentifiers(ctx context.Context, credentialToken string) ([]*CredentialIdentifier, error) {
	if credentialToken == "" {
		return nil, ErrInvalidCredentialToken
	}

	endpoint := c.getCredentialEndpoint()

	type GetCredentialIdentifiers struct {
		XMLName         xml.Name `xml:"tcr:GetCredentialIdentifiers"`
		Xmlns           string   `xml:"xmlns:tcr,attr"`
		CredentialToken string   `xml:"tcr:CredentialToken"`
	}

	type CredentialIdentifierEntry struct {
		Type struct {
			Name       string `xml:"Name"`
			FormatType string `xml:"FormatType"`
		} `xml:"Type"`
		ExemptedFromAuthentication bool   `xml:"ExemptedFromAuthentication"`
		Value                      []byte `xml:"Value"`
	}

	type GetCredentialIdentifiersResponse struct {
		XMLName              xml.Name                    `xml:"GetCredentialIdentifiersResponse"`
		CredentialIdentifier []CredentialIdentifierEntry `xml:"CredentialIdentifier"`
	}

	req := GetCredentialIdentifiers{
		Xmlns:           credentialNamespace,
		CredentialToken: credentialToken,
	}

	var resp GetCredentialIdentifiersResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetCredentialIdentifiers failed: %w", err)
	}

	result := make([]*CredentialIdentifier, 0, len(resp.CredentialIdentifier))

	for _, entry := range resp.CredentialIdentifier {
		result = append(result, &CredentialIdentifier{
			Type: CredentialIdentifierType{
				Name:       entry.Type.Name,
				FormatType: entry.Type.FormatType,
			},
			ExemptedFromAuthentication: entry.ExemptedFromAuthentication,
			Value:                      entry.Value,
		})
	}

	return result, nil
}

// SetCredentialIdentifier creates or updates a credential identifier for a credential.
func (c *Client) SetCredentialIdentifier(ctx context.Context, credentialToken string, identifier CredentialIdentifier) error {
	if credentialToken == "" {
		return ErrInvalidCredentialToken
	}

	endpoint := c.getCredentialEndpoint()

	type CredentialIdentifierXML struct {
		Type struct {
			Name       string `xml:"tcr:Name"`
			FormatType string `xml:"tcr:FormatType"`
		} `xml:"tcr:Type"`
		ExemptedFromAuthentication bool   `xml:"tcr:ExemptedFromAuthentication"`
		Value                      []byte `xml:"tcr:Value"`
	}

	type SetCredentialIdentifier struct {
		XMLName              xml.Name                `xml:"tcr:SetCredentialIdentifier"`
		Xmlns                string                  `xml:"xmlns:tcr,attr"`
		CredentialToken      string                  `xml:"tcr:CredentialToken"`
		CredentialIdentifier CredentialIdentifierXML `xml:"tcr:CredentialIdentifier"`
	}

	ci := CredentialIdentifierXML{
		ExemptedFromAuthentication: identifier.ExemptedFromAuthentication,
		Value:                      identifier.Value,
	}
	ci.Type.Name = identifier.Type.Name
	ci.Type.FormatType = identifier.Type.FormatType

	req := SetCredentialIdentifier{
		Xmlns:                credentialNamespace,
		CredentialToken:      credentialToken,
		CredentialIdentifier: ci,
	}

	var resp struct {
		XMLName xml.Name `xml:"SetCredentialIdentifierResponse"`
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetCredentialIdentifier failed: %w", err)
	}

	return nil
}

// DeleteCredentialIdentifier deletes all identifier values for the specified type name.
func (c *Client) DeleteCredentialIdentifier(ctx context.Context, credentialToken string, identifierTypeName string) error {
	if credentialToken == "" {
		return ErrInvalidCredentialToken
	}

	if identifierTypeName == "" {
		return ErrInvalidCredentialIdentifierTypeName
	}

	endpoint := c.getCredentialEndpoint()

	type DeleteCredentialIdentifier struct {
		XMLName                      xml.Name `xml:"tcr:DeleteCredentialIdentifier"`
		Xmlns                        string   `xml:"xmlns:tcr,attr"`
		CredentialToken              string   `xml:"tcr:CredentialToken"`
		CredentialIdentifierTypeName string   `xml:"tcr:CredentialIdentifierTypeName"`
	}

	req := DeleteCredentialIdentifier{
		Xmlns:                        credentialNamespace,
		CredentialToken:              credentialToken,
		CredentialIdentifierTypeName: identifierTypeName,
	}

	var resp struct {
		XMLName xml.Name `xml:"DeleteCredentialIdentifierResponse"`
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteCredentialIdentifier failed: %w", err)
	}

	return nil
}

// GetCredentialAccessProfiles returns all access profiles associated with the specified credential.
func (c *Client) GetCredentialAccessProfiles(ctx context.Context, credentialToken string) ([]*CredentialAccessProfile, error) {
	if credentialToken == "" {
		return nil, ErrInvalidCredentialToken
	}

	endpoint := c.getCredentialEndpoint()

	type GetCredentialAccessProfiles struct {
		XMLName         xml.Name `xml:"tcr:GetCredentialAccessProfiles"`
		Xmlns           string   `xml:"xmlns:tcr,attr"`
		CredentialToken string   `xml:"tcr:CredentialToken"`
	}

	type CredentialAccessProfileEntry struct {
		AccessProfileToken string `xml:"AccessProfileToken"`
		ValidFrom          string `xml:"ValidFrom"`
		ValidTo            string `xml:"ValidTo"`
	}

	type GetCredentialAccessProfilesResponse struct {
		XMLName                 xml.Name                       `xml:"GetCredentialAccessProfilesResponse"`
		CredentialAccessProfile []CredentialAccessProfileEntry `xml:"CredentialAccessProfile"`
	}

	req := GetCredentialAccessProfiles{
		Xmlns:           credentialNamespace,
		CredentialToken: credentialToken,
	}

	var resp GetCredentialAccessProfilesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetCredentialAccessProfiles failed: %w", err)
	}

	result := make([]*CredentialAccessProfile, 0, len(resp.CredentialAccessProfile))

	for _, entry := range resp.CredentialAccessProfile {
		cap := &CredentialAccessProfile{AccessProfileToken: entry.AccessProfileToken}

		if entry.ValidFrom != "" {
			if t, err := time.Parse(time.RFC3339, entry.ValidFrom); err == nil {
				cap.ValidFrom = &t
			}
		}

		if entry.ValidTo != "" {
			if t, err := time.Parse(time.RFC3339, entry.ValidTo); err == nil {
				cap.ValidTo = &t
			}
		}

		result = append(result, cap)
	}

	return result, nil
}

// SetCredentialAccessProfiles adds or updates access profiles for the specified credential.
func (c *Client) SetCredentialAccessProfiles(ctx context.Context, credentialToken string, accessProfiles []CredentialAccessProfile) error {
	if credentialToken == "" {
		return ErrInvalidCredentialToken
	}

	if len(accessProfiles) == 0 {
		return ErrInvalidAccessProfileToken
	}

	endpoint := c.getCredentialEndpoint()

	type CredentialAccessProfileXML struct {
		AccessProfileToken string `xml:"tcr:AccessProfileToken"`
	}

	type SetCredentialAccessProfiles struct {
		XMLName                 xml.Name                     `xml:"tcr:SetCredentialAccessProfiles"`
		Xmlns                   string                       `xml:"xmlns:tcr,attr"`
		CredentialToken         string                       `xml:"tcr:CredentialToken"`
		CredentialAccessProfile []CredentialAccessProfileXML `xml:"tcr:CredentialAccessProfile"`
	}

	apXMLList := make([]CredentialAccessProfileXML, 0, len(accessProfiles))

	for _, ap := range accessProfiles {
		apXMLList = append(apXMLList, CredentialAccessProfileXML{AccessProfileToken: ap.AccessProfileToken})
	}

	req := SetCredentialAccessProfiles{
		Xmlns:                   credentialNamespace,
		CredentialToken:         credentialToken,
		CredentialAccessProfile: apXMLList,
	}

	var resp struct {
		XMLName xml.Name `xml:"SetCredentialAccessProfilesResponse"`
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetCredentialAccessProfiles failed: %w", err)
	}

	return nil
}

// DeleteCredentialAccessProfiles deletes the access profiles with the given tokens from the credential.
func (c *Client) DeleteCredentialAccessProfiles(ctx context.Context, credentialToken string, accessProfileTokens []string) error {
	if credentialToken == "" {
		return ErrInvalidCredentialToken
	}

	if len(accessProfileTokens) == 0 {
		return ErrInvalidAccessProfileToken
	}

	endpoint := c.getCredentialEndpoint()

	type DeleteCredentialAccessProfiles struct {
		XMLName             xml.Name `xml:"tcr:DeleteCredentialAccessProfiles"`
		Xmlns               string   `xml:"xmlns:tcr,attr"`
		CredentialToken     string   `xml:"tcr:CredentialToken"`
		AccessProfileToken  []string `xml:"tcr:AccessProfileToken"`
	}

	req := DeleteCredentialAccessProfiles{
		Xmlns:              credentialNamespace,
		CredentialToken:    credentialToken,
		AccessProfileToken: accessProfileTokens,
	}

	var resp struct {
		XMLName xml.Name `xml:"DeleteCredentialAccessProfilesResponse"`
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteCredentialAccessProfiles failed: %w", err)
	}

	return nil
}

// GetWhitelist retrieves a paginated list of whitelisted credential identifiers.
func (c *Client) GetWhitelist(ctx context.Context, limit *int, startReference *string, identifierType, formatType *string, value []byte) ([]*CredentialIdentifierItem, string, error) {
	endpoint := c.getCredentialEndpoint()

	type GetWhitelist struct {
		XMLName        xml.Name `xml:"tcr:GetWhitelist"`
		Xmlns          string   `xml:"xmlns:tcr,attr"`
		Limit          *int     `xml:"tcr:Limit,omitempty"`
		StartReference *string  `xml:"tcr:StartReference,omitempty"`
		IdentifierType *string  `xml:"tcr:IdentifierType,omitempty"`
		FormatType     *string  `xml:"tcr:FormatType,omitempty"`
		Value          []byte   `xml:"tcr:Value,omitempty"`
	}

	type CredentialIdentifierItemEntry struct {
		Type struct {
			Name       string `xml:"Name"`
			FormatType string `xml:"FormatType"`
		} `xml:"Type"`
		Value []byte `xml:"Value"`
	}

	type GetWhitelistResponse struct {
		XMLName            xml.Name                        `xml:"GetWhitelistResponse"`
		NextStartReference string                          `xml:"NextStartReference"`
		Identifier         []CredentialIdentifierItemEntry `xml:"Identifier"`
	}

	req := GetWhitelist{
		Xmlns:          credentialNamespace,
		Limit:          limit,
		StartReference: startReference,
		IdentifierType: identifierType,
		FormatType:     formatType,
		Value:          value,
	}

	var resp GetWhitelistResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetWhitelist failed: %w", err)
	}

	result := make([]*CredentialIdentifierItem, 0, len(resp.Identifier))

	for _, entry := range resp.Identifier {
		result = append(result, &CredentialIdentifierItem{
			Type: CredentialIdentifierType{
				Name:       entry.Type.Name,
				FormatType: entry.Type.FormatType,
			},
			Value: entry.Value,
		})
	}

	return result, resp.NextStartReference, nil
}

// AddToWhitelist adds credential identifiers to the whitelist.
func (c *Client) AddToWhitelist(ctx context.Context, identifiers []CredentialIdentifierItem) error {
	if len(identifiers) == 0 {
		return errors.New("at least one identifier is required")
	}

	endpoint := c.getCredentialEndpoint()

	type CredentialIdentifierItemXML struct {
		Type struct {
			Name       string `xml:"tcr:Name"`
			FormatType string `xml:"tcr:FormatType"`
		} `xml:"tcr:Type"`
		Value []byte `xml:"tcr:Value"`
	}

	type AddToWhitelist struct {
		XMLName    xml.Name                      `xml:"tcr:AddToWhitelist"`
		Xmlns      string                        `xml:"xmlns:tcr,attr"`
		Identifier []CredentialIdentifierItemXML `xml:"tcr:Identifier"`
	}

	items := make([]CredentialIdentifierItemXML, 0, len(identifiers))

	for _, id := range identifiers {
		item := CredentialIdentifierItemXML{Value: id.Value}
		item.Type.Name = id.Type.Name
		item.Type.FormatType = id.Type.FormatType
		items = append(items, item)
	}

	req := AddToWhitelist{
		Xmlns:      credentialNamespace,
		Identifier: items,
	}

	var resp struct {
		XMLName xml.Name `xml:"AddToWhitelistResponse"`
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("AddToWhitelist failed: %w", err)
	}

	return nil
}

// RemoveFromWhitelist removes credential identifiers from the whitelist.
func (c *Client) RemoveFromWhitelist(ctx context.Context, identifiers []CredentialIdentifierItem) error {
	if len(identifiers) == 0 {
		return errors.New("at least one identifier is required")
	}

	endpoint := c.getCredentialEndpoint()

	type CredentialIdentifierItemXML struct {
		Type struct {
			Name       string `xml:"tcr:Name"`
			FormatType string `xml:"tcr:FormatType"`
		} `xml:"tcr:Type"`
		Value []byte `xml:"tcr:Value"`
	}

	type RemoveFromWhitelist struct {
		XMLName    xml.Name                      `xml:"tcr:RemoveFromWhitelist"`
		Xmlns      string                        `xml:"xmlns:tcr,attr"`
		Identifier []CredentialIdentifierItemXML `xml:"tcr:Identifier"`
	}

	items := make([]CredentialIdentifierItemXML, 0, len(identifiers))

	for _, id := range identifiers {
		item := CredentialIdentifierItemXML{Value: id.Value}
		item.Type.Name = id.Type.Name
		item.Type.FormatType = id.Type.FormatType
		items = append(items, item)
	}

	req := RemoveFromWhitelist{
		Xmlns:      credentialNamespace,
		Identifier: items,
	}

	var resp struct {
		XMLName xml.Name `xml:"RemoveFromWhitelistResponse"`
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("RemoveFromWhitelist failed: %w", err)
	}

	return nil
}

// DeleteWhitelist clears the entire whitelist.
func (c *Client) DeleteWhitelist(ctx context.Context) error {
	endpoint := c.getCredentialEndpoint()

	type DeleteWhitelist struct {
		XMLName xml.Name `xml:"tcr:DeleteWhitelist"`
		Xmlns   string   `xml:"xmlns:tcr,attr"`
	}

	req := DeleteWhitelist{Xmlns: credentialNamespace}

	var resp struct {
		XMLName xml.Name `xml:"DeleteWhitelistResponse"`
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteWhitelist failed: %w", err)
	}

	return nil
}

// GetBlacklist retrieves a paginated list of blacklisted credential identifiers.
func (c *Client) GetBlacklist(ctx context.Context, limit *int, startReference *string, identifierType, formatType *string, value []byte) ([]*CredentialIdentifierItem, string, error) {
	endpoint := c.getCredentialEndpoint()

	type GetBlacklist struct {
		XMLName        xml.Name `xml:"tcr:GetBlacklist"`
		Xmlns          string   `xml:"xmlns:tcr,attr"`
		Limit          *int     `xml:"tcr:Limit,omitempty"`
		StartReference *string  `xml:"tcr:StartReference,omitempty"`
		IdentifierType *string  `xml:"tcr:IdentifierType,omitempty"`
		FormatType     *string  `xml:"tcr:FormatType,omitempty"`
		Value          []byte   `xml:"tcr:Value,omitempty"`
	}

	type CredentialIdentifierItemEntry struct {
		Type struct {
			Name       string `xml:"Name"`
			FormatType string `xml:"FormatType"`
		} `xml:"Type"`
		Value []byte `xml:"Value"`
	}

	type GetBlacklistResponse struct {
		XMLName            xml.Name                        `xml:"GetBlacklistResponse"`
		NextStartReference string                          `xml:"NextStartReference"`
		Identifier         []CredentialIdentifierItemEntry `xml:"Identifier"`
	}

	req := GetBlacklist{
		Xmlns:          credentialNamespace,
		Limit:          limit,
		StartReference: startReference,
		IdentifierType: identifierType,
		FormatType:     formatType,
		Value:          value,
	}

	var resp GetBlacklistResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetBlacklist failed: %w", err)
	}

	result := make([]*CredentialIdentifierItem, 0, len(resp.Identifier))

	for _, entry := range resp.Identifier {
		result = append(result, &CredentialIdentifierItem{
			Type: CredentialIdentifierType{
				Name:       entry.Type.Name,
				FormatType: entry.Type.FormatType,
			},
			Value: entry.Value,
		})
	}

	return result, resp.NextStartReference, nil
}

// AddToBlacklist adds credential identifiers to the blacklist.
func (c *Client) AddToBlacklist(ctx context.Context, identifiers []CredentialIdentifierItem) error {
	if len(identifiers) == 0 {
		return errors.New("at least one identifier is required")
	}

	endpoint := c.getCredentialEndpoint()

	type CredentialIdentifierItemXML struct {
		Type struct {
			Name       string `xml:"tcr:Name"`
			FormatType string `xml:"tcr:FormatType"`
		} `xml:"tcr:Type"`
		Value []byte `xml:"tcr:Value"`
	}

	type AddToBlacklist struct {
		XMLName    xml.Name                      `xml:"tcr:AddToBlacklist"`
		Xmlns      string                        `xml:"xmlns:tcr,attr"`
		Identifier []CredentialIdentifierItemXML `xml:"tcr:Identifier"`
	}

	items := make([]CredentialIdentifierItemXML, 0, len(identifiers))

	for _, id := range identifiers {
		item := CredentialIdentifierItemXML{Value: id.Value}
		item.Type.Name = id.Type.Name
		item.Type.FormatType = id.Type.FormatType
		items = append(items, item)
	}

	req := AddToBlacklist{
		Xmlns:      credentialNamespace,
		Identifier: items,
	}

	var resp struct {
		XMLName xml.Name `xml:"AddToBlacklistResponse"`
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("AddToBlacklist failed: %w", err)
	}

	return nil
}

// RemoveFromBlacklist removes credential identifiers from the blacklist.
func (c *Client) RemoveFromBlacklist(ctx context.Context, identifiers []CredentialIdentifierItem) error {
	if len(identifiers) == 0 {
		return errors.New("at least one identifier is required")
	}

	endpoint := c.getCredentialEndpoint()

	type CredentialIdentifierItemXML struct {
		Type struct {
			Name       string `xml:"tcr:Name"`
			FormatType string `xml:"tcr:FormatType"`
		} `xml:"tcr:Type"`
		Value []byte `xml:"tcr:Value"`
	}

	type RemoveFromBlacklist struct {
		XMLName    xml.Name                      `xml:"tcr:RemoveFromBlacklist"`
		Xmlns      string                        `xml:"xmlns:tcr,attr"`
		Identifier []CredentialIdentifierItemXML `xml:"tcr:Identifier"`
	}

	items := make([]CredentialIdentifierItemXML, 0, len(identifiers))

	for _, id := range identifiers {
		item := CredentialIdentifierItemXML{Value: id.Value}
		item.Type.Name = id.Type.Name
		item.Type.FormatType = id.Type.FormatType
		items = append(items, item)
	}

	req := RemoveFromBlacklist{
		Xmlns:      credentialNamespace,
		Identifier: items,
	}

	var resp struct {
		XMLName xml.Name `xml:"RemoveFromBlacklistResponse"`
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("RemoveFromBlacklist failed: %w", err)
	}

	return nil
}

// DeleteBlacklist clears the entire blacklist.
func (c *Client) DeleteBlacklist(ctx context.Context) error {
	endpoint := c.getCredentialEndpoint()

	type DeleteBlacklist struct {
		XMLName xml.Name `xml:"tcr:DeleteBlacklist"`
		Xmlns   string   `xml:"xmlns:tcr,attr"`
	}

	req := DeleteBlacklist{Xmlns: credentialNamespace}

	var resp struct {
		XMLName xml.Name `xml:"DeleteBlacklistResponse"`
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteBlacklist failed: %w", err)
	}

	return nil
}
