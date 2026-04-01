package onvif

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/0x524a/onvif-go/internal/soap"
)

// Authentication Behavior service namespace.
const authBehaviorNamespace = "http://www.onvif.org/ver10/authenticationbehavior/wsdl"

// Authentication Behavior service errors.
var (
	// ErrInvalidAuthenticationProfileToken is returned when an authentication profile token is empty.
	ErrInvalidAuthenticationProfileToken = errors.New("invalid authentication profile token: cannot be empty")
	// ErrAuthenticationProfileNil is returned when an authentication profile is nil.
	ErrAuthenticationProfileNil = errors.New("authentication profile cannot be nil")
	// ErrInvalidSecurityLevelToken is returned when a security level token is empty.
	ErrInvalidSecurityLevelToken = errors.New("invalid security level token: cannot be empty")
	// ErrSecurityLevelNil is returned when a security level is nil.
	ErrSecurityLevelNil = errors.New("security level cannot be nil")
)

// getAuthBehaviorEndpoint returns the authentication behavior endpoint, falling back to device endpoint.
func (c *Client) getAuthBehaviorEndpoint() string {
	if c.authBehaviorEndpoint != "" {
		return c.authBehaviorEndpoint
	}

	return c.endpoint
}

// GetAuthBehaviorServiceCapabilities retrieves the capabilities of the authentication behavior service.
func (c *Client) GetAuthBehaviorServiceCapabilities(ctx context.Context) (*AuthBehaviorServiceCapabilities, error) {
	endpoint := c.getAuthBehaviorEndpoint()

	type GetServiceCapabilities struct {
		XMLName xml.Name `xml:"tab:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:tab,attr"`
	}

	type GetServiceCapabilitiesResponse struct {
		XMLName      xml.Name `xml:"GetServiceCapabilitiesResponse"`
		Capabilities struct {
			MaxLimit                                uint   `xml:"MaxLimit,attr"`
			MaxAuthenticationProfiles               uint   `xml:"MaxAuthenticationProfiles,attr"`
			MaxPoliciesPerAuthenticationProfile     uint   `xml:"MaxPoliciesPerAuthenticationProfile,attr"`
			MaxSecurityLevels                       uint   `xml:"MaxSecurityLevels,attr"`
			MaxRecognitionGroupsPerSecurityLevel    uint   `xml:"MaxRecognitionGroupsPerSecurityLevel,attr"`
			MaxRecognitionMethodsPerRecognitionGroup uint  `xml:"MaxRecognitionMethodsPerRecognitionGroup,attr"`
			ClientSuppliedTokenSupported            *bool  `xml:"ClientSuppliedTokenSupported,attr"`
			SupportedAuthenticationModes            string `xml:"SupportedAuthenticationModes,attr"`
		} `xml:"Capabilities"`
	}

	req := GetServiceCapabilities{Xmlns: authBehaviorNamespace}

	var resp GetServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAuthBehaviorServiceCapabilities failed: %w", err)
	}

	caps := &AuthBehaviorServiceCapabilities{
		MaxLimit:                                resp.Capabilities.MaxLimit,
		MaxAuthenticationProfiles:               resp.Capabilities.MaxAuthenticationProfiles,
		MaxPoliciesPerAuthenticationProfile:     resp.Capabilities.MaxPoliciesPerAuthenticationProfile,
		MaxSecurityLevels:                       resp.Capabilities.MaxSecurityLevels,
		MaxRecognitionGroupsPerSecurityLevel:    resp.Capabilities.MaxRecognitionGroupsPerSecurityLevel,
		MaxRecognitionMethodsPerRecognitionGroup: resp.Capabilities.MaxRecognitionMethodsPerRecognitionGroup,
		SupportedAuthenticationModes:            resp.Capabilities.SupportedAuthenticationModes,
	}

	if resp.Capabilities.ClientSuppliedTokenSupported != nil {
		caps.ClientSuppliedTokenSupported = *resp.Capabilities.ClientSuppliedTokenSupported
	}

	return caps, nil
}

// authProfileInfoEntry is the internal XML representation of an AuthenticationProfileInfo entry.
type authProfileInfoEntry struct {
	Token       string `xml:"token,attr"`
	Name        string `xml:"Name"`
	Description string `xml:"Description"`
}

// mapAuthProfileInfo maps an authProfileInfoEntry to a public AuthenticationProfileInfo.
func mapAuthProfileInfo(e authProfileInfoEntry) AuthenticationProfileInfo {
	return AuthenticationProfileInfo{
		Token:       e.Token,
		Name:        e.Name,
		Description: e.Description,
	}
}

// GetAuthenticationProfileInfoList retrieves a paginated list of AuthenticationProfileInfo items.
func (c *Client) GetAuthenticationProfileInfoList(ctx context.Context, limit *int, startReference *string) ([]*AuthenticationProfileInfo, string, error) {
	endpoint := c.getAuthBehaviorEndpoint()

	type GetAuthenticationProfileInfoList struct {
		XMLName        xml.Name `xml:"tab:GetAuthenticationProfileInfoList"`
		Xmlns          string   `xml:"xmlns:tab,attr"`
		Limit          *int     `xml:"tab:Limit,omitempty"`
		StartReference *string  `xml:"tab:StartReference,omitempty"`
	}

	type GetAuthenticationProfileInfoListResponse struct {
		XMLName                   xml.Name               `xml:"GetAuthenticationProfileInfoListResponse"`
		NextStartReference        string                 `xml:"NextStartReference"`
		AuthenticationProfileInfo []authProfileInfoEntry `xml:"AuthenticationProfileInfo"`
	}

	req := GetAuthenticationProfileInfoList{
		Xmlns:          authBehaviorNamespace,
		Limit:          limit,
		StartReference: startReference,
	}

	var resp GetAuthenticationProfileInfoListResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetAuthenticationProfileInfoList failed: %w", err)
	}

	result := make([]*AuthenticationProfileInfo, 0, len(resp.AuthenticationProfileInfo))

	for _, e := range resp.AuthenticationProfileInfo {
		info := mapAuthProfileInfo(e)
		result = append(result, &info)
	}

	return result, resp.NextStartReference, nil
}

// GetAuthenticationProfileInfo retrieves AuthenticationProfileInfo items by token.
func (c *Client) GetAuthenticationProfileInfo(ctx context.Context, tokens []string) ([]*AuthenticationProfileInfo, error) {
	if len(tokens) == 0 {
		return nil, ErrInvalidAuthenticationProfileToken
	}

	endpoint := c.getAuthBehaviorEndpoint()

	type GetAuthenticationProfileInfo struct {
		XMLName xml.Name `xml:"tab:GetAuthenticationProfileInfo"`
		Xmlns   string   `xml:"xmlns:tab,attr"`
		Token   []string `xml:"tab:Token"`
	}

	type GetAuthenticationProfileInfoResponse struct {
		XMLName                   xml.Name               `xml:"GetAuthenticationProfileInfoResponse"`
		AuthenticationProfileInfo []authProfileInfoEntry `xml:"AuthenticationProfileInfo"`
	}

	req := GetAuthenticationProfileInfo{
		Xmlns: authBehaviorNamespace,
		Token: tokens,
	}

	var resp GetAuthenticationProfileInfoResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAuthenticationProfileInfo failed: %w", err)
	}

	result := make([]*AuthenticationProfileInfo, 0, len(resp.AuthenticationProfileInfo))

	for _, e := range resp.AuthenticationProfileInfo {
		info := mapAuthProfileInfo(e)
		result = append(result, &info)
	}

	return result, nil
}

// authProfileEntry is the internal XML representation of an AuthenticationProfile entry.
type authProfileEntry struct {
	Token                     string `xml:"token,attr"`
	Name                      string `xml:"Name"`
	Description               string `xml:"Description"`
	DefaultSecurityLevelToken string `xml:"DefaultSecurityLevelToken"`
	AuthenticationPolicy      []struct {
		ScheduleToken              string `xml:"ScheduleToken"`
		SecurityLevelConstraint    []struct {
			ActiveRegularSchedule    bool   `xml:"ActiveRegularSchedule"`
			ActiveSpecialDaySchedule bool   `xml:"ActiveSpecialDaySchedule"`
			AuthenticationMode       string `xml:"AuthenticationMode"`
			SecurityLevelToken       string `xml:"SecurityLevelToken"`
		} `xml:"SecurityLevelConstraint"`
	} `xml:"AuthenticationPolicy"`
}

// mapAuthProfile maps an authProfileEntry to a public AuthenticationProfile.
func mapAuthProfile(e authProfileEntry) AuthenticationProfile {
	ap := AuthenticationProfile{
		AuthenticationProfileInfo: AuthenticationProfileInfo{
			Token:       e.Token,
			Name:        e.Name,
			Description: e.Description,
		},
		DefaultSecurityLevelToken: e.DefaultSecurityLevelToken,
	}

	for _, p := range e.AuthenticationPolicy {
		policy := AuthenticationPolicy{
			ScheduleToken: p.ScheduleToken,
		}

		for _, sc := range p.SecurityLevelConstraint {
			policy.SecurityLevelConstraints = append(policy.SecurityLevelConstraints, SecurityLevelConstraint{
				ActiveRegularSchedule:    sc.ActiveRegularSchedule,
				ActiveSpecialDaySchedule: sc.ActiveSpecialDaySchedule,
				AuthenticationMode:       sc.AuthenticationMode,
				SecurityLevelToken:       sc.SecurityLevelToken,
			})
		}

		ap.AuthenticationPolicies = append(ap.AuthenticationPolicies, policy)
	}

	return ap
}

// GetAuthenticationProfileList retrieves a paginated list of AuthenticationProfile items.
func (c *Client) GetAuthenticationProfileList(ctx context.Context, limit *int, startReference *string) ([]*AuthenticationProfile, string, error) {
	endpoint := c.getAuthBehaviorEndpoint()

	type GetAuthenticationProfileList struct {
		XMLName        xml.Name `xml:"tab:GetAuthenticationProfileList"`
		Xmlns          string   `xml:"xmlns:tab,attr"`
		Limit          *int     `xml:"tab:Limit,omitempty"`
		StartReference *string  `xml:"tab:StartReference,omitempty"`
	}

	type GetAuthenticationProfileListResponse struct {
		XMLName               xml.Name           `xml:"GetAuthenticationProfileListResponse"`
		NextStartReference    string             `xml:"NextStartReference"`
		AuthenticationProfile []authProfileEntry `xml:"AuthenticationProfile"`
	}

	req := GetAuthenticationProfileList{
		Xmlns:          authBehaviorNamespace,
		Limit:          limit,
		StartReference: startReference,
	}

	var resp GetAuthenticationProfileListResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetAuthenticationProfileList failed: %w", err)
	}

	result := make([]*AuthenticationProfile, 0, len(resp.AuthenticationProfile))

	for _, e := range resp.AuthenticationProfile {
		p := mapAuthProfile(e)
		result = append(result, &p)
	}

	return result, resp.NextStartReference, nil
}

// GetAuthenticationProfiles retrieves AuthenticationProfile items by token.
func (c *Client) GetAuthenticationProfiles(ctx context.Context, tokens []string) ([]*AuthenticationProfile, error) {
	if len(tokens) == 0 {
		return nil, ErrInvalidAuthenticationProfileToken
	}

	endpoint := c.getAuthBehaviorEndpoint()

	type GetAuthenticationProfiles struct {
		XMLName xml.Name `xml:"tab:GetAuthenticationProfiles"`
		Xmlns   string   `xml:"xmlns:tab,attr"`
		Token   []string `xml:"tab:Token"`
	}

	type GetAuthenticationProfilesResponse struct {
		XMLName               xml.Name           `xml:"GetAuthenticationProfilesResponse"`
		AuthenticationProfile []authProfileEntry `xml:"AuthenticationProfile"`
	}

	req := GetAuthenticationProfiles{
		Xmlns: authBehaviorNamespace,
		Token: tokens,
	}

	var resp GetAuthenticationProfilesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAuthenticationProfiles failed: %w", err)
	}

	result := make([]*AuthenticationProfile, 0, len(resp.AuthenticationProfile))

	for _, e := range resp.AuthenticationProfile {
		p := mapAuthProfile(e)
		result = append(result, &p)
	}

	return result, nil
}

// authProfileXMLPayload represents an AuthenticationProfile in XML.
type authProfileXMLPayload struct {
	Token                     string                          `xml:"token,attr,omitempty"`
	Name                      string                          `xml:"tab:Name"`
	Description               string                          `xml:"tab:Description,omitempty"`
	DefaultSecurityLevelToken string                          `xml:"tab:DefaultSecurityLevelToken"`
	AuthenticationPolicy      []authPolicyXMLPayload          `xml:"tab:AuthenticationPolicy,omitempty"`
}

// authPolicyXMLPayload represents an AuthenticationPolicy in XML.
type authPolicyXMLPayload struct {
	ScheduleToken           string                              `xml:"tab:ScheduleToken"`
	SecurityLevelConstraint []securityLevelConstraintXMLPayload `xml:"tab:SecurityLevelConstraint"`
}

// securityLevelConstraintXMLPayload represents a SecurityLevelConstraint in XML.
type securityLevelConstraintXMLPayload struct {
	ActiveRegularSchedule    bool   `xml:"tab:ActiveRegularSchedule"`
	ActiveSpecialDaySchedule bool   `xml:"tab:ActiveSpecialDaySchedule"`
	AuthenticationMode       string `xml:"tab:AuthenticationMode,omitempty"`
	SecurityLevelToken       string `xml:"tab:SecurityLevelToken"`
}

// buildAuthProfilePayload converts an AuthenticationProfile to its XML payload form.
func buildAuthProfilePayload(p *AuthenticationProfile) authProfileXMLPayload {
	payload := authProfileXMLPayload{
		Token:                     p.Token,
		Name:                      p.Name,
		Description:               p.Description,
		DefaultSecurityLevelToken: p.DefaultSecurityLevelToken,
	}

	for _, policy := range p.AuthenticationPolicies {
		policyPayload := authPolicyXMLPayload{
			ScheduleToken: policy.ScheduleToken,
		}

		for _, sc := range policy.SecurityLevelConstraints {
			policyPayload.SecurityLevelConstraint = append(policyPayload.SecurityLevelConstraint, securityLevelConstraintXMLPayload{
				ActiveRegularSchedule:    sc.ActiveRegularSchedule,
				ActiveSpecialDaySchedule: sc.ActiveSpecialDaySchedule,
				AuthenticationMode:       sc.AuthenticationMode,
				SecurityLevelToken:       sc.SecurityLevelToken,
			})
		}

		payload.AuthenticationPolicy = append(payload.AuthenticationPolicy, policyPayload)
	}

	return payload
}

// CreateAuthenticationProfile creates a new authentication profile and returns its assigned token.
func (c *Client) CreateAuthenticationProfile(ctx context.Context, profile *AuthenticationProfile) (string, error) {
	if profile == nil {
		return "", ErrAuthenticationProfileNil
	}

	endpoint := c.getAuthBehaviorEndpoint()

	type CreateAuthenticationProfile struct {
		XMLName               xml.Name              `xml:"tab:CreateAuthenticationProfile"`
		Xmlns                 string                `xml:"xmlns:tab,attr"`
		AuthenticationProfile authProfileXMLPayload `xml:"tab:AuthenticationProfile"`
	}

	type CreateAuthenticationProfileResponse struct {
		XMLName xml.Name `xml:"CreateAuthenticationProfileResponse"`
		Token   string   `xml:"Token"`
	}

	req := CreateAuthenticationProfile{
		Xmlns:                 authBehaviorNamespace,
		AuthenticationProfile: buildAuthProfilePayload(profile),
	}

	var resp CreateAuthenticationProfileResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateAuthenticationProfile failed: %w", err)
	}

	return resp.Token, nil
}

// ModifyAuthenticationProfile modifies an existing authentication profile.
func (c *Client) ModifyAuthenticationProfile(ctx context.Context, profile *AuthenticationProfile) error {
	if profile == nil {
		return ErrAuthenticationProfileNil
	}

	if profile.Token == "" {
		return ErrInvalidAuthenticationProfileToken
	}

	endpoint := c.getAuthBehaviorEndpoint()

	type ModifyAuthenticationProfile struct {
		XMLName               xml.Name              `xml:"tab:ModifyAuthenticationProfile"`
		Xmlns                 string                `xml:"xmlns:tab,attr"`
		AuthenticationProfile authProfileXMLPayload `xml:"tab:AuthenticationProfile"`
	}

	type ModifyAuthenticationProfileResponse struct {
		XMLName xml.Name `xml:"ModifyAuthenticationProfileResponse"`
	}

	req := ModifyAuthenticationProfile{
		Xmlns:                 authBehaviorNamespace,
		AuthenticationProfile: buildAuthProfilePayload(profile),
	}

	var resp ModifyAuthenticationProfileResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("ModifyAuthenticationProfile failed: %w", err)
	}

	return nil
}

// SetAuthenticationProfile creates or replaces an authentication profile (requires ClientSuppliedTokenSupported capability).
func (c *Client) SetAuthenticationProfile(ctx context.Context, profile *AuthenticationProfile) error {
	if profile == nil {
		return ErrAuthenticationProfileNil
	}

	endpoint := c.getAuthBehaviorEndpoint()

	type SetAuthenticationProfile struct {
		XMLName               xml.Name              `xml:"tab:SetAuthenticationProfile"`
		Xmlns                 string                `xml:"xmlns:tab,attr"`
		AuthenticationProfile authProfileXMLPayload `xml:"tab:AuthenticationProfile"`
	}

	type SetAuthenticationProfileResponse struct {
		XMLName xml.Name `xml:"SetAuthenticationProfileResponse"`
	}

	req := SetAuthenticationProfile{
		Xmlns:                 authBehaviorNamespace,
		AuthenticationProfile: buildAuthProfilePayload(profile),
	}

	var resp SetAuthenticationProfileResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetAuthenticationProfile failed: %w", err)
	}

	return nil
}

// DeleteAuthenticationProfile deletes an authentication profile by token.
func (c *Client) DeleteAuthenticationProfile(ctx context.Context, token string) error {
	if token == "" {
		return ErrInvalidAuthenticationProfileToken
	}

	endpoint := c.getAuthBehaviorEndpoint()

	type DeleteAuthenticationProfile struct {
		XMLName xml.Name `xml:"tab:DeleteAuthenticationProfile"`
		Xmlns   string   `xml:"xmlns:tab,attr"`
		Token   string   `xml:"tab:Token"`
	}

	type DeleteAuthenticationProfileResponse struct {
		XMLName xml.Name `xml:"DeleteAuthenticationProfileResponse"`
	}

	req := DeleteAuthenticationProfile{
		Xmlns: authBehaviorNamespace,
		Token: token,
	}

	var resp DeleteAuthenticationProfileResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteAuthenticationProfile failed: %w", err)
	}

	return nil
}

// securityLevelInfoEntry is the internal XML representation of a SecurityLevelInfo entry.
type securityLevelInfoEntry struct {
	Token       string `xml:"token,attr"`
	Name        string `xml:"Name"`
	Priority    int    `xml:"Priority"`
	Description string `xml:"Description"`
}

// mapSecurityLevelInfo maps a securityLevelInfoEntry to a public SecurityLevelInfo.
func mapSecurityLevelInfo(e securityLevelInfoEntry) SecurityLevelInfo {
	return SecurityLevelInfo{
		Token:       e.Token,
		Name:        e.Name,
		Priority:    e.Priority,
		Description: e.Description,
	}
}

// GetSecurityLevelInfoList retrieves a paginated list of SecurityLevelInfo items.
func (c *Client) GetSecurityLevelInfoList(ctx context.Context, limit *int, startReference *string) ([]*SecurityLevelInfo, string, error) {
	endpoint := c.getAuthBehaviorEndpoint()

	type GetSecurityLevelInfoList struct {
		XMLName        xml.Name `xml:"tab:GetSecurityLevelInfoList"`
		Xmlns          string   `xml:"xmlns:tab,attr"`
		Limit          *int     `xml:"tab:Limit,omitempty"`
		StartReference *string  `xml:"tab:StartReference,omitempty"`
	}

	type GetSecurityLevelInfoListResponse struct {
		XMLName            xml.Name                 `xml:"GetSecurityLevelInfoListResponse"`
		NextStartReference string                   `xml:"NextStartReference"`
		SecurityLevelInfo  []securityLevelInfoEntry `xml:"SecurityLevelInfo"`
	}

	req := GetSecurityLevelInfoList{
		Xmlns:          authBehaviorNamespace,
		Limit:          limit,
		StartReference: startReference,
	}

	var resp GetSecurityLevelInfoListResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetSecurityLevelInfoList failed: %w", err)
	}

	result := make([]*SecurityLevelInfo, 0, len(resp.SecurityLevelInfo))

	for _, e := range resp.SecurityLevelInfo {
		info := mapSecurityLevelInfo(e)
		result = append(result, &info)
	}

	return result, resp.NextStartReference, nil
}

// GetSecurityLevelInfo retrieves SecurityLevelInfo items by token.
func (c *Client) GetSecurityLevelInfo(ctx context.Context, tokens []string) ([]*SecurityLevelInfo, error) {
	if len(tokens) == 0 {
		return nil, ErrInvalidSecurityLevelToken
	}

	endpoint := c.getAuthBehaviorEndpoint()

	type GetSecurityLevelInfo struct {
		XMLName xml.Name `xml:"tab:GetSecurityLevelInfo"`
		Xmlns   string   `xml:"xmlns:tab,attr"`
		Token   []string `xml:"tab:Token"`
	}

	type GetSecurityLevelInfoResponse struct {
		XMLName           xml.Name                 `xml:"GetSecurityLevelInfoResponse"`
		SecurityLevelInfo []securityLevelInfoEntry `xml:"SecurityLevelInfo"`
	}

	req := GetSecurityLevelInfo{
		Xmlns: authBehaviorNamespace,
		Token: tokens,
	}

	var resp GetSecurityLevelInfoResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetSecurityLevelInfo failed: %w", err)
	}

	result := make([]*SecurityLevelInfo, 0, len(resp.SecurityLevelInfo))

	for _, e := range resp.SecurityLevelInfo {
		info := mapSecurityLevelInfo(e)
		result = append(result, &info)
	}

	return result, nil
}

// securityLevelEntry is the internal XML representation of a SecurityLevel entry.
type securityLevelEntry struct {
	Token       string `xml:"token,attr"`
	Name        string `xml:"Name"`
	Priority    int    `xml:"Priority"`
	Description string `xml:"Description"`
	RecognitionGroup []struct {
		RecognitionMethod []struct {
			RecognitionType string `xml:"RecognitionType"`
			Order           int    `xml:"Order"`
		} `xml:"RecognitionMethod"`
	} `xml:"RecognitionGroup"`
}

// mapSecurityLevel maps a securityLevelEntry to a public SecurityLevel.
func mapSecurityLevel(e securityLevelEntry) SecurityLevel {
	sl := SecurityLevel{
		SecurityLevelInfo: SecurityLevelInfo{
			Token:       e.Token,
			Name:        e.Name,
			Priority:    e.Priority,
			Description: e.Description,
		},
	}

	for _, rg := range e.RecognitionGroup {
		group := RecognitionGroup{}

		for _, rm := range rg.RecognitionMethod {
			group.RecognitionMethods = append(group.RecognitionMethods, RecognitionMethod{
				RecognitionType: rm.RecognitionType,
				Order:           rm.Order,
			})
		}

		sl.RecognitionGroups = append(sl.RecognitionGroups, group)
	}

	return sl
}

// GetSecurityLevelList retrieves a paginated list of SecurityLevel items.
func (c *Client) GetSecurityLevelList(ctx context.Context, limit *int, startReference *string) ([]*SecurityLevel, string, error) {
	endpoint := c.getAuthBehaviorEndpoint()

	type GetSecurityLevelList struct {
		XMLName        xml.Name `xml:"tab:GetSecurityLevelList"`
		Xmlns          string   `xml:"xmlns:tab,attr"`
		Limit          *int     `xml:"tab:Limit,omitempty"`
		StartReference *string  `xml:"tab:StartReference,omitempty"`
	}

	type GetSecurityLevelListResponse struct {
		XMLName            xml.Name             `xml:"GetSecurityLevelListResponse"`
		NextStartReference string               `xml:"NextStartReference"`
		SecurityLevel      []securityLevelEntry `xml:"SecurityLevel"`
	}

	req := GetSecurityLevelList{
		Xmlns:          authBehaviorNamespace,
		Limit:          limit,
		StartReference: startReference,
	}

	var resp GetSecurityLevelListResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetSecurityLevelList failed: %w", err)
	}

	result := make([]*SecurityLevel, 0, len(resp.SecurityLevel))

	for _, e := range resp.SecurityLevel {
		sl := mapSecurityLevel(e)
		result = append(result, &sl)
	}

	return result, resp.NextStartReference, nil
}

// GetSecurityLevels retrieves SecurityLevel items by token.
func (c *Client) GetSecurityLevels(ctx context.Context, tokens []string) ([]*SecurityLevel, error) {
	if len(tokens) == 0 {
		return nil, ErrInvalidSecurityLevelToken
	}

	endpoint := c.getAuthBehaviorEndpoint()

	type GetSecurityLevels struct {
		XMLName xml.Name `xml:"tab:GetSecurityLevels"`
		Xmlns   string   `xml:"xmlns:tab,attr"`
		Token   []string `xml:"tab:Token"`
	}

	type GetSecurityLevelsResponse struct {
		XMLName       xml.Name             `xml:"GetSecurityLevelsResponse"`
		SecurityLevel []securityLevelEntry `xml:"SecurityLevel"`
	}

	req := GetSecurityLevels{
		Xmlns: authBehaviorNamespace,
		Token: tokens,
	}

	var resp GetSecurityLevelsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetSecurityLevels failed: %w", err)
	}

	result := make([]*SecurityLevel, 0, len(resp.SecurityLevel))

	for _, e := range resp.SecurityLevel {
		sl := mapSecurityLevel(e)
		result = append(result, &sl)
	}

	return result, nil
}

// securityLevelXMLPayload represents a SecurityLevel in XML.
type securityLevelXMLPayload struct {
	Token            string                      `xml:"token,attr,omitempty"`
	Name             string                      `xml:"tab:Name"`
	Priority         int                         `xml:"tab:Priority"`
	Description      string                      `xml:"tab:Description,omitempty"`
	RecognitionGroup []recognitionGroupXMLPayload `xml:"tab:RecognitionGroup,omitempty"`
}

// recognitionGroupXMLPayload represents a RecognitionGroup in XML.
type recognitionGroupXMLPayload struct {
	RecognitionMethod []recognitionMethodXMLPayload `xml:"tab:RecognitionMethod,omitempty"`
}

// recognitionMethodXMLPayload represents a RecognitionMethod in XML.
type recognitionMethodXMLPayload struct {
	RecognitionType string `xml:"tab:RecognitionType"`
	Order           int    `xml:"tab:Order"`
}

// buildSecurityLevelPayload converts a SecurityLevel to its XML payload form.
func buildSecurityLevelPayload(sl *SecurityLevel) securityLevelXMLPayload {
	payload := securityLevelXMLPayload{
		Token:       sl.Token,
		Name:        sl.Name,
		Priority:    sl.Priority,
		Description: sl.Description,
	}

	for _, rg := range sl.RecognitionGroups {
		groupPayload := recognitionGroupXMLPayload{}

		for _, rm := range rg.RecognitionMethods {
			groupPayload.RecognitionMethod = append(groupPayload.RecognitionMethod, recognitionMethodXMLPayload{
				RecognitionType: rm.RecognitionType,
				Order:           rm.Order,
			})
		}

		payload.RecognitionGroup = append(payload.RecognitionGroup, groupPayload)
	}

	return payload
}

// CreateSecurityLevel creates a new security level and returns its assigned token.
func (c *Client) CreateSecurityLevel(ctx context.Context, securityLevel *SecurityLevel) (string, error) {
	if securityLevel == nil {
		return "", ErrSecurityLevelNil
	}

	endpoint := c.getAuthBehaviorEndpoint()

	type CreateSecurityLevel struct {
		XMLName       xml.Name                `xml:"tab:CreateSecurityLevel"`
		Xmlns         string                  `xml:"xmlns:tab,attr"`
		SecurityLevel securityLevelXMLPayload `xml:"tab:SecurityLevel"`
	}

	type CreateSecurityLevelResponse struct {
		XMLName xml.Name `xml:"CreateSecurityLevelResponse"`
		Token   string   `xml:"Token"`
	}

	req := CreateSecurityLevel{
		Xmlns:         authBehaviorNamespace,
		SecurityLevel: buildSecurityLevelPayload(securityLevel),
	}

	var resp CreateSecurityLevelResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateSecurityLevel failed: %w", err)
	}

	return resp.Token, nil
}

// ModifySecurityLevel modifies an existing security level.
func (c *Client) ModifySecurityLevel(ctx context.Context, securityLevel *SecurityLevel) error {
	if securityLevel == nil {
		return ErrSecurityLevelNil
	}

	if securityLevel.Token == "" {
		return ErrInvalidSecurityLevelToken
	}

	endpoint := c.getAuthBehaviorEndpoint()

	type ModifySecurityLevel struct {
		XMLName       xml.Name                `xml:"tab:ModifySecurityLevel"`
		Xmlns         string                  `xml:"xmlns:tab,attr"`
		SecurityLevel securityLevelXMLPayload `xml:"tab:SecurityLevel"`
	}

	type ModifySecurityLevelResponse struct {
		XMLName xml.Name `xml:"ModifySecurityLevelResponse"`
	}

	req := ModifySecurityLevel{
		Xmlns:         authBehaviorNamespace,
		SecurityLevel: buildSecurityLevelPayload(securityLevel),
	}

	var resp ModifySecurityLevelResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("ModifySecurityLevel failed: %w", err)
	}

	return nil
}

// SetSecurityLevel creates or replaces a security level (requires ClientSuppliedTokenSupported capability).
func (c *Client) SetSecurityLevel(ctx context.Context, securityLevel *SecurityLevel) error {
	if securityLevel == nil {
		return ErrSecurityLevelNil
	}

	endpoint := c.getAuthBehaviorEndpoint()

	type SetSecurityLevel struct {
		XMLName       xml.Name                `xml:"tab:SetSecurityLevel"`
		Xmlns         string                  `xml:"xmlns:tab,attr"`
		SecurityLevel securityLevelXMLPayload `xml:"tab:SecurityLevel"`
	}

	type SetSecurityLevelResponse struct {
		XMLName xml.Name `xml:"SetSecurityLevelResponse"`
	}

	req := SetSecurityLevel{
		Xmlns:         authBehaviorNamespace,
		SecurityLevel: buildSecurityLevelPayload(securityLevel),
	}

	var resp SetSecurityLevelResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetSecurityLevel failed: %w", err)
	}

	return nil
}

// DeleteSecurityLevel deletes a security level by token.
func (c *Client) DeleteSecurityLevel(ctx context.Context, token string) error {
	if token == "" {
		return ErrInvalidSecurityLevelToken
	}

	endpoint := c.getAuthBehaviorEndpoint()

	type DeleteSecurityLevel struct {
		XMLName xml.Name `xml:"tab:DeleteSecurityLevel"`
		Xmlns   string   `xml:"xmlns:tab,attr"`
		Token   string   `xml:"tab:Token"`
	}

	type DeleteSecurityLevelResponse struct {
		XMLName xml.Name `xml:"DeleteSecurityLevelResponse"`
	}

	req := DeleteSecurityLevel{
		Xmlns: authBehaviorNamespace,
		Token: token,
	}

	var resp DeleteSecurityLevelResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteSecurityLevel failed: %w", err)
	}

	return nil
}
