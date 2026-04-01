package onvif

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/0x524a/onvif-go/internal/soap"
)

// Access Control service namespace.
const accessControlNamespace = "http://www.onvif.org/ver10/accesscontrol/wsdl"

// Access Control service errors.
var (
	// ErrInvalidAccessPointToken is returned when an access point token is empty.
	ErrInvalidAccessPointToken = errors.New("invalid access point token: cannot be empty")
	// ErrInvalidAreaToken is returned when an area token is empty.
	ErrInvalidAreaToken = errors.New("invalid area token: cannot be empty")
	// ErrAccessPointNil is returned when an access point is nil.
	ErrAccessPointNil = errors.New("access point cannot be nil")
	// ErrAreaNil is returned when an area is nil.
	ErrAreaNil = errors.New("area cannot be nil")
)

// getAccessControlEndpoint returns the access control endpoint, falling back to device endpoint.
func (c *Client) getAccessControlEndpoint() string {
	if c.accessControlEndpoint != "" {
		return c.accessControlEndpoint
	}

	return c.endpoint
}

// GetAccessControlServiceCapabilities retrieves the capabilities of the access control service.
func (c *Client) GetAccessControlServiceCapabilities(ctx context.Context) (*AccessControlServiceCapabilities, error) {
	endpoint := c.getAccessControlEndpoint()

	type GetServiceCapabilities struct {
		XMLName xml.Name `xml:"tac:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:tac,attr"`
	}

	type GetServiceCapabilitiesResponse struct {
		XMLName      xml.Name `xml:"GetServiceCapabilitiesResponse"`
		Capabilities struct {
			MaxLimit                       uint `xml:"MaxLimit,attr"`
			MaxAccessPoints                uint `xml:"MaxAccessPoints,attr"`
			MaxAreas                       uint `xml:"MaxAreas,attr"`
			ClientSuppliedTokenSupported   bool `xml:"ClientSuppliedTokenSupported,attr"`
			AccessPointManagementSupported bool `xml:"AccessPointManagementSupported,attr"`
			AreaManagementSupported        bool `xml:"AreaManagementSupported,attr"`
		} `xml:"Capabilities"`
	}

	req := GetServiceCapabilities{
		Xmlns: accessControlNamespace,
	}

	var resp GetServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAccessControlServiceCapabilities failed: %w", err)
	}

	return &AccessControlServiceCapabilities{
		MaxLimit:                       resp.Capabilities.MaxLimit,
		MaxAccessPoints:                resp.Capabilities.MaxAccessPoints,
		MaxAreas:                       resp.Capabilities.MaxAreas,
		ClientSuppliedTokenSupported:   resp.Capabilities.ClientSuppliedTokenSupported,
		AccessPointManagementSupported: resp.Capabilities.AccessPointManagementSupported,
		AreaManagementSupported:        resp.Capabilities.AreaManagementSupported,
	}, nil
}

// mapAccessPointInfo maps internal XML response to public AccessPointInfo.
func mapAccessPointInfo(token, name, description, areaFrom, areaTo, entityType, entity string,
	caps struct {
		DisableAccessPoint    bool    `xml:"DisableAccessPoint,attr"`
		Duress                *bool   `xml:"Duress,attr"`
		AnonymousAccess       *bool   `xml:"AnonymousAccess,attr"`
		AccessTaken           *bool   `xml:"AccessTaken,attr"`
		ExternalAuthorization *bool   `xml:"ExternalAuthorization,attr"`
		IdentifierAccess      *bool   `xml:"IdentifierAccess,attr"`
		SupportedSecurityLevels []struct {
			Value string `xml:",chardata"`
		} `xml:"SupportedSecurityLevels"`
	},
) AccessPointInfo {
	info := AccessPointInfo{
		Token:       token,
		Name:        name,
		Description: description,
		AreaFrom:    areaFrom,
		AreaTo:      areaTo,
		EntityType:  entityType,
		Entity:      entity,
		Capabilities: AccessPointCapabilities{
			DisableAccessPoint:    caps.DisableAccessPoint,
			Duress:                caps.Duress,
			AnonymousAccess:       caps.AnonymousAccess,
			AccessTaken:           caps.AccessTaken,
			ExternalAuthorization: caps.ExternalAuthorization,
			IdentifierAccess:      caps.IdentifierAccess,
		},
	}

	for _, sl := range caps.SupportedSecurityLevels {
		info.Capabilities.SupportedSecurityLevels = append(info.Capabilities.SupportedSecurityLevels, sl.Value)
	}

	return info
}

// GetAccessPointInfoList retrieves a paginated list of all AccessPointInfo items.
func (c *Client) GetAccessPointInfoList(ctx context.Context, limit *int, startReference *string) ([]*AccessPointInfo, string, error) {
	endpoint := c.getAccessControlEndpoint()

	type GetAccessPointInfoList struct {
		XMLName        xml.Name `xml:"tac:GetAccessPointInfoList"`
		Xmlns          string   `xml:"xmlns:tac,attr"`
		Limit          *int     `xml:"tac:Limit,omitempty"`
		StartReference *string  `xml:"tac:StartReference,omitempty"`
	}

	type AccessPointInfoEntry struct {
		Token       string `xml:"token,attr"`
		Name        string `xml:"Name"`
		Description string `xml:"Description"`
		AreaFrom    string `xml:"AreaFrom"`
		AreaTo      string `xml:"AreaTo"`
		EntityType  string `xml:"EntityType"`
		Entity      string `xml:"Entity"`
		Capabilities struct {
			DisableAccessPoint    bool  `xml:"DisableAccessPoint,attr"`
			Duress                *bool `xml:"Duress,attr"`
			AnonymousAccess       *bool `xml:"AnonymousAccess,attr"`
			AccessTaken           *bool `xml:"AccessTaken,attr"`
			ExternalAuthorization *bool `xml:"ExternalAuthorization,attr"`
			IdentifierAccess      *bool `xml:"IdentifierAccess,attr"`
			SupportedSecurityLevels []struct {
				Value string `xml:",chardata"`
			} `xml:"SupportedSecurityLevels"`
		} `xml:"Capabilities"`
	}

	type GetAccessPointInfoListResponse struct {
		XMLName          xml.Name               `xml:"GetAccessPointInfoListResponse"`
		NextStartReference string               `xml:"NextStartReference"`
		AccessPointInfo  []AccessPointInfoEntry `xml:"AccessPointInfo"`
	}

	req := GetAccessPointInfoList{
		Xmlns:          accessControlNamespace,
		Limit:          limit,
		StartReference: startReference,
	}

	var resp GetAccessPointInfoListResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetAccessPointInfoList failed: %w", err)
	}

	result := make([]*AccessPointInfo, 0, len(resp.AccessPointInfo))

	for _, entry := range resp.AccessPointInfo {
		info := mapAccessPointInfo(
			entry.Token, entry.Name, entry.Description,
			entry.AreaFrom, entry.AreaTo, entry.EntityType, entry.Entity,
			entry.Capabilities,
		)
		result = append(result, &info)
	}

	return result, resp.NextStartReference, nil
}

// GetAccessPointInfo retrieves AccessPointInfo items by token.
func (c *Client) GetAccessPointInfo(ctx context.Context, tokens []string) ([]*AccessPointInfo, error) {
	if len(tokens) == 0 {
		return nil, ErrInvalidAccessPointToken
	}

	endpoint := c.getAccessControlEndpoint()

	type GetAccessPointInfo struct {
		XMLName xml.Name `xml:"tac:GetAccessPointInfo"`
		Xmlns   string   `xml:"xmlns:tac,attr"`
		Token   []string `xml:"tac:Token"`
	}

	type AccessPointInfoEntry struct {
		Token       string `xml:"token,attr"`
		Name        string `xml:"Name"`
		Description string `xml:"Description"`
		AreaFrom    string `xml:"AreaFrom"`
		AreaTo      string `xml:"AreaTo"`
		EntityType  string `xml:"EntityType"`
		Entity      string `xml:"Entity"`
		Capabilities struct {
			DisableAccessPoint    bool  `xml:"DisableAccessPoint,attr"`
			Duress                *bool `xml:"Duress,attr"`
			AnonymousAccess       *bool `xml:"AnonymousAccess,attr"`
			AccessTaken           *bool `xml:"AccessTaken,attr"`
			ExternalAuthorization *bool `xml:"ExternalAuthorization,attr"`
			IdentifierAccess      *bool `xml:"IdentifierAccess,attr"`
			SupportedSecurityLevels []struct {
				Value string `xml:",chardata"`
			} `xml:"SupportedSecurityLevels"`
		} `xml:"Capabilities"`
	}

	type GetAccessPointInfoResponse struct {
		XMLName         xml.Name               `xml:"GetAccessPointInfoResponse"`
		AccessPointInfo []AccessPointInfoEntry `xml:"AccessPointInfo"`
	}

	req := GetAccessPointInfo{
		Xmlns: accessControlNamespace,
		Token: tokens,
	}

	var resp GetAccessPointInfoResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAccessPointInfo failed: %w", err)
	}

	result := make([]*AccessPointInfo, 0, len(resp.AccessPointInfo))

	for _, entry := range resp.AccessPointInfo {
		info := mapAccessPointInfo(
			entry.Token, entry.Name, entry.Description,
			entry.AreaFrom, entry.AreaTo, entry.EntityType, entry.Entity,
			entry.Capabilities,
		)
		result = append(result, &info)
	}

	return result, nil
}

// GetAccessPointList retrieves a paginated list of all AccessPoint items.
func (c *Client) GetAccessPointList(ctx context.Context, limit *int, startReference *string) ([]*AccessPoint, string, error) {
	endpoint := c.getAccessControlEndpoint()

	type GetAccessPointList struct {
		XMLName        xml.Name `xml:"tac:GetAccessPointList"`
		Xmlns          string   `xml:"xmlns:tac,attr"`
		Limit          *int     `xml:"tac:Limit,omitempty"`
		StartReference *string  `xml:"tac:StartReference,omitempty"`
	}

	type AccessPointEntry struct {
		Token       string `xml:"token,attr"`
		Name        string `xml:"Name"`
		Description string `xml:"Description"`
		AreaFrom    string `xml:"AreaFrom"`
		AreaTo      string `xml:"AreaTo"`
		EntityType  string `xml:"EntityType"`
		Entity      string `xml:"Entity"`
		Capabilities struct {
			DisableAccessPoint    bool  `xml:"DisableAccessPoint,attr"`
			Duress                *bool `xml:"Duress,attr"`
			AnonymousAccess       *bool `xml:"AnonymousAccess,attr"`
			AccessTaken           *bool `xml:"AccessTaken,attr"`
			ExternalAuthorization *bool `xml:"ExternalAuthorization,attr"`
			IdentifierAccess      *bool `xml:"IdentifierAccess,attr"`
			SupportedSecurityLevels []struct {
				Value string `xml:",chardata"`
			} `xml:"SupportedSecurityLevels"`
		} `xml:"Capabilities"`
		AuthenticationProfileToken string `xml:"AuthenticationProfileToken"`
	}

	type GetAccessPointListResponse struct {
		XMLName            xml.Name           `xml:"GetAccessPointListResponse"`
		NextStartReference string             `xml:"NextStartReference"`
		AccessPoint        []AccessPointEntry `xml:"AccessPoint"`
	}

	req := GetAccessPointList{
		Xmlns:          accessControlNamespace,
		Limit:          limit,
		StartReference: startReference,
	}

	var resp GetAccessPointListResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetAccessPointList failed: %w", err)
	}

	result := make([]*AccessPoint, 0, len(resp.AccessPoint))

	for _, entry := range resp.AccessPoint {
		info := mapAccessPointInfo(
			entry.Token, entry.Name, entry.Description,
			entry.AreaFrom, entry.AreaTo, entry.EntityType, entry.Entity,
			entry.Capabilities,
		)
		ap := &AccessPoint{
			AccessPointInfo:            info,
			AuthenticationProfileToken: entry.AuthenticationProfileToken,
		}
		result = append(result, ap)
	}

	return result, resp.NextStartReference, nil
}

// GetAccessPoints retrieves AccessPoint items by token.
func (c *Client) GetAccessPoints(ctx context.Context, tokens []string) ([]*AccessPoint, error) {
	if len(tokens) == 0 {
		return nil, ErrInvalidAccessPointToken
	}

	endpoint := c.getAccessControlEndpoint()

	type GetAccessPoints struct {
		XMLName xml.Name `xml:"tac:GetAccessPoints"`
		Xmlns   string   `xml:"xmlns:tac,attr"`
		Token   []string `xml:"tac:Token"`
	}

	type AccessPointEntry struct {
		Token       string `xml:"token,attr"`
		Name        string `xml:"Name"`
		Description string `xml:"Description"`
		AreaFrom    string `xml:"AreaFrom"`
		AreaTo      string `xml:"AreaTo"`
		EntityType  string `xml:"EntityType"`
		Entity      string `xml:"Entity"`
		Capabilities struct {
			DisableAccessPoint    bool  `xml:"DisableAccessPoint,attr"`
			Duress                *bool `xml:"Duress,attr"`
			AnonymousAccess       *bool `xml:"AnonymousAccess,attr"`
			AccessTaken           *bool `xml:"AccessTaken,attr"`
			ExternalAuthorization *bool `xml:"ExternalAuthorization,attr"`
			IdentifierAccess      *bool `xml:"IdentifierAccess,attr"`
			SupportedSecurityLevels []struct {
				Value string `xml:",chardata"`
			} `xml:"SupportedSecurityLevels"`
		} `xml:"Capabilities"`
		AuthenticationProfileToken string `xml:"AuthenticationProfileToken"`
	}

	type GetAccessPointsResponse struct {
		XMLName     xml.Name           `xml:"GetAccessPointsResponse"`
		AccessPoint []AccessPointEntry `xml:"AccessPoint"`
	}

	req := GetAccessPoints{
		Xmlns: accessControlNamespace,
		Token: tokens,
	}

	var resp GetAccessPointsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAccessPoints failed: %w", err)
	}

	result := make([]*AccessPoint, 0, len(resp.AccessPoint))

	for _, entry := range resp.AccessPoint {
		info := mapAccessPointInfo(
			entry.Token, entry.Name, entry.Description,
			entry.AreaFrom, entry.AreaTo, entry.EntityType, entry.Entity,
			entry.Capabilities,
		)
		ap := &AccessPoint{
			AccessPointInfo:            info,
			AuthenticationProfileToken: entry.AuthenticationProfileToken,
		}
		result = append(result, ap)
	}

	return result, nil
}

// CreateAccessPoint creates a new access point. Returns the token of the created access point.
func (c *Client) CreateAccessPoint(ctx context.Context, ap *AccessPoint) (string, error) {
	if ap == nil {
		return "", ErrAccessPointNil
	}

	endpoint := c.getAccessControlEndpoint()

	type CapabilitiesXML struct {
		DisableAccessPoint bool `xml:"DisableAccessPoint,attr"`
	}

	type AccessPointXML struct {
		Token                      string          `xml:"token,attr,omitempty"`
		Name                       string          `xml:"tac:Name"`
		Description                string          `xml:"tac:Description,omitempty"`
		AreaFrom                   string          `xml:"tac:AreaFrom,omitempty"`
		AreaTo                     string          `xml:"tac:AreaTo,omitempty"`
		Entity                     string          `xml:"tac:Entity"`
		Capabilities               CapabilitiesXML `xml:"tac:Capabilities"`
		AuthenticationProfileToken string          `xml:"tac:AuthenticationProfileToken,omitempty"`
	}

	type CreateAccessPoint struct {
		XMLName     xml.Name       `xml:"tac:CreateAccessPoint"`
		Xmlns       string         `xml:"xmlns:tac,attr"`
		AccessPoint AccessPointXML `xml:"tac:AccessPoint"`
	}

	type CreateAccessPointResponse struct {
		XMLName xml.Name `xml:"CreateAccessPointResponse"`
		Token   string   `xml:"Token"`
	}

	req := CreateAccessPoint{
		Xmlns: accessControlNamespace,
		AccessPoint: AccessPointXML{
			Name:        ap.Name,
			Description: ap.Description,
			AreaFrom:    ap.AreaFrom,
			AreaTo:      ap.AreaTo,
			Entity:      ap.Entity,
			Capabilities: CapabilitiesXML{
				DisableAccessPoint: ap.Capabilities.DisableAccessPoint,
			},
			AuthenticationProfileToken: ap.AuthenticationProfileToken,
		},
	}

	var resp CreateAccessPointResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateAccessPoint failed: %w", err)
	}

	return resp.Token, nil
}

// ModifyAccessPoint modifies an existing access point.
func (c *Client) ModifyAccessPoint(ctx context.Context, ap *AccessPoint) error {
	if ap == nil {
		return ErrAccessPointNil
	}

	if strings.TrimSpace(ap.Token) == "" {
		return ErrInvalidAccessPointToken
	}

	endpoint := c.getAccessControlEndpoint()

	type CapabilitiesXML struct {
		DisableAccessPoint bool `xml:"DisableAccessPoint,attr"`
	}

	type AccessPointXML struct {
		Token                      string          `xml:"token,attr"`
		Name                       string          `xml:"tac:Name"`
		Description                string          `xml:"tac:Description,omitempty"`
		AreaFrom                   string          `xml:"tac:AreaFrom,omitempty"`
		AreaTo                     string          `xml:"tac:AreaTo,omitempty"`
		Entity                     string          `xml:"tac:Entity"`
		Capabilities               CapabilitiesXML `xml:"tac:Capabilities"`
		AuthenticationProfileToken string          `xml:"tac:AuthenticationProfileToken,omitempty"`
	}

	type ModifyAccessPoint struct {
		XMLName     xml.Name       `xml:"tac:ModifyAccessPoint"`
		Xmlns       string         `xml:"xmlns:tac,attr"`
		AccessPoint AccessPointXML `xml:"tac:AccessPoint"`
	}

	type ModifyAccessPointResponse struct {
		XMLName xml.Name `xml:"ModifyAccessPointResponse"`
	}

	req := ModifyAccessPoint{
		Xmlns: accessControlNamespace,
		AccessPoint: AccessPointXML{
			Token:       ap.Token,
			Name:        ap.Name,
			Description: ap.Description,
			AreaFrom:    ap.AreaFrom,
			AreaTo:      ap.AreaTo,
			Entity:      ap.Entity,
			Capabilities: CapabilitiesXML{
				DisableAccessPoint: ap.Capabilities.DisableAccessPoint,
			},
			AuthenticationProfileToken: ap.AuthenticationProfileToken,
		},
	}

	var resp ModifyAccessPointResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("ModifyAccessPoint failed: %w", err)
	}

	return nil
}

// DeleteAccessPoint deletes an access point by token.
func (c *Client) DeleteAccessPoint(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidAccessPointToken
	}

	endpoint := c.getAccessControlEndpoint()

	type DeleteAccessPoint struct {
		XMLName xml.Name `xml:"tac:DeleteAccessPoint"`
		Xmlns   string   `xml:"xmlns:tac,attr"`
		Token   string   `xml:"tac:Token"`
	}

	type DeleteAccessPointResponse struct {
		XMLName xml.Name `xml:"DeleteAccessPointResponse"`
	}

	req := DeleteAccessPoint{
		Xmlns: accessControlNamespace,
		Token: token,
	}

	var resp DeleteAccessPointResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteAccessPoint failed: %w", err)
	}

	return nil
}

// SetAccessPointAuthenticationProfile sets the authentication profile for an access point.
func (c *Client) SetAccessPointAuthenticationProfile(ctx context.Context, accessPointToken, authProfileToken string) error {
	if strings.TrimSpace(accessPointToken) == "" {
		return ErrInvalidAccessPointToken
	}

	endpoint := c.getAccessControlEndpoint()

	type SetAccessPointAuthenticationProfile struct {
		XMLName                    xml.Name `xml:"tac:SetAccessPointAuthenticationProfile"`
		Xmlns                      string   `xml:"xmlns:tac,attr"`
		Token                      string   `xml:"tac:Token"`
		AuthenticationProfileToken string   `xml:"tac:AuthenticationProfileToken"`
	}

	type SetAccessPointAuthenticationProfileResponse struct {
		XMLName xml.Name `xml:"SetAccessPointAuthenticationProfileResponse"`
	}

	req := SetAccessPointAuthenticationProfile{
		Xmlns:                      accessControlNamespace,
		Token:                      accessPointToken,
		AuthenticationProfileToken: authProfileToken,
	}

	var resp SetAccessPointAuthenticationProfileResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetAccessPointAuthenticationProfile failed: %w", err)
	}

	return nil
}

// DeleteAccessPointAuthenticationProfile reverts authentication for an access point to default.
func (c *Client) DeleteAccessPointAuthenticationProfile(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidAccessPointToken
	}

	endpoint := c.getAccessControlEndpoint()

	type DeleteAccessPointAuthenticationProfile struct {
		XMLName xml.Name `xml:"tac:DeleteAccessPointAuthenticationProfile"`
		Xmlns   string   `xml:"xmlns:tac,attr"`
		Token   string   `xml:"tac:Token"`
	}

	type DeleteAccessPointAuthenticationProfileResponse struct {
		XMLName xml.Name `xml:"DeleteAccessPointAuthenticationProfileResponse"`
	}

	req := DeleteAccessPointAuthenticationProfile{
		Xmlns: accessControlNamespace,
		Token: token,
	}

	var resp DeleteAccessPointAuthenticationProfileResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteAccessPointAuthenticationProfile failed: %w", err)
	}

	return nil
}

// GetAreaInfoList retrieves a paginated list of all AreaInfo items.
func (c *Client) GetAreaInfoList(ctx context.Context, limit *int, startReference *string) ([]*AreaInfo, string, error) {
	endpoint := c.getAccessControlEndpoint()

	type GetAreaInfoList struct {
		XMLName        xml.Name `xml:"tac:GetAreaInfoList"`
		Xmlns          string   `xml:"xmlns:tac,attr"`
		Limit          *int     `xml:"tac:Limit,omitempty"`
		StartReference *string  `xml:"tac:StartReference,omitempty"`
	}

	type AreaInfoEntry struct {
		Token       string `xml:"token,attr"`
		Name        string `xml:"Name"`
		Description string `xml:"Description"`
	}

	type GetAreaInfoListResponse struct {
		XMLName            xml.Name        `xml:"GetAreaInfoListResponse"`
		NextStartReference string          `xml:"NextStartReference"`
		AreaInfo           []AreaInfoEntry `xml:"AreaInfo"`
	}

	req := GetAreaInfoList{
		Xmlns:          accessControlNamespace,
		Limit:          limit,
		StartReference: startReference,
	}

	var resp GetAreaInfoListResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetAreaInfoList failed: %w", err)
	}

	result := make([]*AreaInfo, 0, len(resp.AreaInfo))

	for _, entry := range resp.AreaInfo {
		result = append(result, &AreaInfo{
			Token:       entry.Token,
			Name:        entry.Name,
			Description: entry.Description,
		})
	}

	return result, resp.NextStartReference, nil
}

// GetAreaInfo retrieves AreaInfo items by token.
func (c *Client) GetAreaInfo(ctx context.Context, tokens []string) ([]*AreaInfo, error) {
	if len(tokens) == 0 {
		return nil, ErrInvalidAreaToken
	}

	endpoint := c.getAccessControlEndpoint()

	type GetAreaInfo struct {
		XMLName xml.Name `xml:"tac:GetAreaInfo"`
		Xmlns   string   `xml:"xmlns:tac,attr"`
		Token   []string `xml:"tac:Token"`
	}

	type AreaInfoEntry struct {
		Token       string `xml:"token,attr"`
		Name        string `xml:"Name"`
		Description string `xml:"Description"`
	}

	type GetAreaInfoResponse struct {
		XMLName  xml.Name        `xml:"GetAreaInfoResponse"`
		AreaInfo []AreaInfoEntry `xml:"AreaInfo"`
	}

	req := GetAreaInfo{
		Xmlns: accessControlNamespace,
		Token: tokens,
	}

	var resp GetAreaInfoResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAreaInfo failed: %w", err)
	}

	result := make([]*AreaInfo, 0, len(resp.AreaInfo))

	for _, entry := range resp.AreaInfo {
		result = append(result, &AreaInfo{
			Token:       entry.Token,
			Name:        entry.Name,
			Description: entry.Description,
		})
	}

	return result, nil
}

// GetAreaList retrieves a paginated list of all Area items.
func (c *Client) GetAreaList(ctx context.Context, limit *int, startReference *string) ([]*Area, string, error) {
	endpoint := c.getAccessControlEndpoint()

	type GetAreaList struct {
		XMLName        xml.Name `xml:"tac:GetAreaList"`
		Xmlns          string   `xml:"xmlns:tac,attr"`
		Limit          *int     `xml:"tac:Limit,omitempty"`
		StartReference *string  `xml:"tac:StartReference,omitempty"`
	}

	type AreaEntry struct {
		Token       string `xml:"token,attr"`
		Name        string `xml:"Name"`
		Description string `xml:"Description"`
	}

	type GetAreaListResponse struct {
		XMLName            xml.Name    `xml:"GetAreaListResponse"`
		NextStartReference string      `xml:"NextStartReference"`
		Area               []AreaEntry `xml:"Area"`
	}

	req := GetAreaList{
		Xmlns:          accessControlNamespace,
		Limit:          limit,
		StartReference: startReference,
	}

	var resp GetAreaListResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetAreaList failed: %w", err)
	}

	result := make([]*Area, 0, len(resp.Area))

	for _, entry := range resp.Area {
		result = append(result, &Area{
			AreaInfo: AreaInfo{
				Token:       entry.Token,
				Name:        entry.Name,
				Description: entry.Description,
			},
		})
	}

	return result, resp.NextStartReference, nil
}

// GetAreas retrieves Area items by token.
func (c *Client) GetAreas(ctx context.Context, tokens []string) ([]*Area, error) {
	if len(tokens) == 0 {
		return nil, ErrInvalidAreaToken
	}

	endpoint := c.getAccessControlEndpoint()

	type GetAreas struct {
		XMLName xml.Name `xml:"tac:GetAreas"`
		Xmlns   string   `xml:"xmlns:tac,attr"`
		Token   []string `xml:"tac:Token"`
	}

	type AreaEntry struct {
		Token       string `xml:"token,attr"`
		Name        string `xml:"Name"`
		Description string `xml:"Description"`
	}

	type GetAreasResponse struct {
		XMLName xml.Name    `xml:"GetAreasResponse"`
		Area    []AreaEntry `xml:"Area"`
	}

	req := GetAreas{
		Xmlns: accessControlNamespace,
		Token: tokens,
	}

	var resp GetAreasResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAreas failed: %w", err)
	}

	result := make([]*Area, 0, len(resp.Area))

	for _, entry := range resp.Area {
		result = append(result, &Area{
			AreaInfo: AreaInfo{
				Token:       entry.Token,
				Name:        entry.Name,
				Description: entry.Description,
			},
		})
	}

	return result, nil
}

// CreateArea creates a new area. Returns the token of the created area.
func (c *Client) CreateArea(ctx context.Context, area *Area) (string, error) {
	if area == nil {
		return "", ErrAreaNil
	}

	endpoint := c.getAccessControlEndpoint()

	type AreaXML struct {
		Token       string `xml:"token,attr,omitempty"`
		Name        string `xml:"tac:Name"`
		Description string `xml:"tac:Description,omitempty"`
	}

	type CreateArea struct {
		XMLName xml.Name `xml:"tac:CreateArea"`
		Xmlns   string   `xml:"xmlns:tac,attr"`
		Area    AreaXML  `xml:"tac:Area"`
	}

	type CreateAreaResponse struct {
		XMLName xml.Name `xml:"CreateAreaResponse"`
		Token   string   `xml:"Token"`
	}

	req := CreateArea{
		Xmlns: accessControlNamespace,
		Area: AreaXML{
			Name:        area.Name,
			Description: area.Description,
		},
	}

	var resp CreateAreaResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateArea failed: %w", err)
	}

	return resp.Token, nil
}

// ModifyArea modifies an existing area.
func (c *Client) ModifyArea(ctx context.Context, area *Area) error {
	if area == nil {
		return ErrAreaNil
	}

	if strings.TrimSpace(area.Token) == "" {
		return ErrInvalidAreaToken
	}

	endpoint := c.getAccessControlEndpoint()

	type AreaXML struct {
		Token       string `xml:"token,attr"`
		Name        string `xml:"tac:Name"`
		Description string `xml:"tac:Description,omitempty"`
	}

	type ModifyArea struct {
		XMLName xml.Name `xml:"tac:ModifyArea"`
		Xmlns   string   `xml:"xmlns:tac,attr"`
		Area    AreaXML  `xml:"tac:Area"`
	}

	type ModifyAreaResponse struct {
		XMLName xml.Name `xml:"ModifyAreaResponse"`
	}

	req := ModifyArea{
		Xmlns: accessControlNamespace,
		Area: AreaXML{
			Token:       area.Token,
			Name:        area.Name,
			Description: area.Description,
		},
	}

	var resp ModifyAreaResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("ModifyArea failed: %w", err)
	}

	return nil
}

// DeleteArea deletes an area by token.
func (c *Client) DeleteArea(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidAreaToken
	}

	endpoint := c.getAccessControlEndpoint()

	type DeleteArea struct {
		XMLName xml.Name `xml:"tac:DeleteArea"`
		Xmlns   string   `xml:"xmlns:tac,attr"`
		Token   string   `xml:"tac:Token"`
	}

	type DeleteAreaResponse struct {
		XMLName xml.Name `xml:"DeleteAreaResponse"`
	}

	req := DeleteArea{
		Xmlns: accessControlNamespace,
		Token: token,
	}

	var resp DeleteAreaResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteArea failed: %w", err)
	}

	return nil
}

// GetAccessPointState retrieves the state of an access point.
func (c *Client) GetAccessPointState(ctx context.Context, token string) (*AccessPointState, error) {
	if strings.TrimSpace(token) == "" {
		return nil, ErrInvalidAccessPointToken
	}

	endpoint := c.getAccessControlEndpoint()

	type GetAccessPointState struct {
		XMLName xml.Name `xml:"tac:GetAccessPointState"`
		Xmlns   string   `xml:"xmlns:tac,attr"`
		Token   string   `xml:"tac:Token"`
	}

	type GetAccessPointStateResponse struct {
		XMLName          xml.Name `xml:"GetAccessPointStateResponse"`
		AccessPointState struct {
			Enabled bool `xml:"Enabled"`
		} `xml:"AccessPointState"`
	}

	req := GetAccessPointState{
		Xmlns: accessControlNamespace,
		Token: token,
	}

	var resp GetAccessPointStateResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAccessPointState failed: %w", err)
	}

	return &AccessPointState{
		Enabled: resp.AccessPointState.Enabled,
	}, nil
}

// EnableAccessPoint enables a specific access point.
func (c *Client) EnableAccessPoint(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidAccessPointToken
	}

	endpoint := c.getAccessControlEndpoint()

	type EnableAccessPoint struct {
		XMLName xml.Name `xml:"tac:EnableAccessPoint"`
		Xmlns   string   `xml:"xmlns:tac,attr"`
		Token   string   `xml:"tac:Token"`
	}

	type EnableAccessPointResponse struct {
		XMLName xml.Name `xml:"EnableAccessPointResponse"`
	}

	req := EnableAccessPoint{
		Xmlns: accessControlNamespace,
		Token: token,
	}

	var resp EnableAccessPointResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("EnableAccessPoint failed: %w", err)
	}

	return nil
}

// DisableAccessPoint disables a specific access point.
func (c *Client) DisableAccessPoint(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidAccessPointToken
	}

	endpoint := c.getAccessControlEndpoint()

	type DisableAccessPoint struct {
		XMLName xml.Name `xml:"tac:DisableAccessPoint"`
		Xmlns   string   `xml:"xmlns:tac,attr"`
		Token   string   `xml:"tac:Token"`
	}

	type DisableAccessPointResponse struct {
		XMLName xml.Name `xml:"DisableAccessPointResponse"`
	}

	req := DisableAccessPoint{
		Xmlns: accessControlNamespace,
		Token: token,
	}

	var resp DisableAccessPointResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DisableAccessPoint failed: %w", err)
	}

	return nil
}

// ExternalAuthorization sends an access decision for an access point.
func (c *Client) ExternalAuthorization(ctx context.Context, accessPointToken, credentialToken, reason string, decision AccessControlDecision) error {
	if strings.TrimSpace(accessPointToken) == "" {
		return ErrInvalidAccessPointToken
	}

	endpoint := c.getAccessControlEndpoint()

	type ExternalAuthorization struct {
		XMLName          xml.Name              `xml:"tac:ExternalAuthorization"`
		Xmlns            string                `xml:"xmlns:tac,attr"`
		AccessPointToken string                `xml:"tac:AccessPointToken"`
		CredentialToken  string                `xml:"tac:CredentialToken,omitempty"`
		Reason           string                `xml:"tac:Reason,omitempty"`
		Decision         AccessControlDecision `xml:"tac:Decision"`
	}

	type ExternalAuthorizationResponse struct {
		XMLName xml.Name `xml:"ExternalAuthorizationResponse"`
	}

	req := ExternalAuthorization{
		Xmlns:            accessControlNamespace,
		AccessPointToken: accessPointToken,
		CredentialToken:  credentialToken,
		Reason:           reason,
		Decision:         decision,
	}

	var resp ExternalAuthorizationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("ExternalAuthorization failed: %w", err)
	}

	return nil
}
