package onvif

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/0x524a/onvif-go/internal/soap"
)

// Schedule service namespace.
const scheduleNamespace = "http://www.onvif.org/ver10/schedule/wsdl"

// Schedule service errors.
var (
	// ErrInvalidScheduleToken is returned when a schedule token is empty.
	ErrInvalidScheduleToken = errors.New("invalid schedule token: cannot be empty")
	// ErrScheduleNil is returned when a schedule is nil.
	ErrScheduleNil = errors.New("schedule cannot be nil")
	// ErrInvalidSpecialDayGroupToken is returned when a special day group token is empty.
	ErrInvalidSpecialDayGroupToken = errors.New("invalid special day group token: cannot be empty")
	// ErrSpecialDayGroupNil is returned when a special day group is nil.
	ErrSpecialDayGroupNil = errors.New("special day group cannot be nil")
)

// getScheduleEndpoint returns the schedule endpoint, falling back to device endpoint.
func (c *Client) getScheduleEndpoint() string {
	if c.scheduleEndpoint != "" {
		return c.scheduleEndpoint
	}

	return c.endpoint
}

// GetScheduleServiceCapabilities retrieves the capabilities of the schedule service.
func (c *Client) GetScheduleServiceCapabilities(ctx context.Context) (*ScheduleServiceCapabilities, error) {
	endpoint := c.getScheduleEndpoint()

	type GetServiceCapabilities struct {
		XMLName xml.Name `xml:"tsc:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:tsc,attr"`
	}

	type GetServiceCapabilitiesResponse struct {
		XMLName      xml.Name `xml:"GetServiceCapabilitiesResponse"`
		Capabilities struct {
			MaxLimit                     uint   `xml:"MaxLimit,attr"`
			MaxSchedules                 uint   `xml:"MaxSchedules,attr"`
			MaxTimePeriodsPerDay         uint   `xml:"MaxTimePeriodsPerDay,attr"`
			MaxSpecialDayGroups          uint   `xml:"MaxSpecialDayGroups,attr"`
			MaxDaysInSpecialDayGroup     uint   `xml:"MaxDaysInSpecialDayGroup,attr"`
			MaxSpecialDaysSchedules      uint   `xml:"MaxSpecialDaysSchedules,attr"`
			ExtendedRecurrenceSupported  bool   `xml:"ExtendedRecurrenceSupported,attr"`
			SpecialDaysSupported         bool   `xml:"SpecialDaysSupported,attr"`
			StateReportingSupported      bool   `xml:"StateReportingSupported,attr"`
			ClientSuppliedTokenSupported *bool  `xml:"ClientSuppliedTokenSupported,attr"`
		} `xml:"Capabilities"`
	}

	req := GetServiceCapabilities{Xmlns: scheduleNamespace}

	var resp GetServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetScheduleServiceCapabilities failed: %w", err)
	}

	caps := &ScheduleServiceCapabilities{
		MaxLimit:                    resp.Capabilities.MaxLimit,
		MaxSchedules:                resp.Capabilities.MaxSchedules,
		MaxTimePeriodsPerDay:        resp.Capabilities.MaxTimePeriodsPerDay,
		MaxSpecialDayGroups:         resp.Capabilities.MaxSpecialDayGroups,
		MaxDaysInSpecialDayGroup:    resp.Capabilities.MaxDaysInSpecialDayGroup,
		MaxSpecialDaysSchedules:     resp.Capabilities.MaxSpecialDaysSchedules,
		ExtendedRecurrenceSupported: resp.Capabilities.ExtendedRecurrenceSupported,
		SpecialDaysSupported:        resp.Capabilities.SpecialDaysSupported,
		StateReportingSupported:     resp.Capabilities.StateReportingSupported,
	}

	if resp.Capabilities.ClientSuppliedTokenSupported != nil {
		caps.ClientSuppliedTokenSupported = *resp.Capabilities.ClientSuppliedTokenSupported
	}

	return caps, nil
}

// GetScheduleState retrieves the current state of a schedule by token.
func (c *Client) GetScheduleState(ctx context.Context, token string) (*ScheduleState, error) {
	if token == "" {
		return nil, ErrInvalidScheduleToken
	}

	endpoint := c.getScheduleEndpoint()

	type GetScheduleState struct {
		XMLName xml.Name `xml:"tsc:GetScheduleState"`
		Xmlns   string   `xml:"xmlns:tsc,attr"`
		Token   string   `xml:"tsc:Token"`
	}

	type GetScheduleStateResponse struct {
		XMLName       xml.Name `xml:"GetScheduleStateResponse"`
		ScheduleState struct {
			Active     bool  `xml:"Active"`
			SpecialDay *bool `xml:"SpecialDay"`
		} `xml:"ScheduleState"`
	}

	req := GetScheduleState{
		Xmlns: scheduleNamespace,
		Token: token,
	}

	var resp GetScheduleStateResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetScheduleState failed: %w", err)
	}

	return &ScheduleState{
		Active:     resp.ScheduleState.Active,
		SpecialDay: resp.ScheduleState.SpecialDay,
	}, nil
}

// scheduleInfoEntry is the internal XML representation of a ScheduleInfo entry.
type scheduleInfoEntry struct {
	Token       string `xml:"token,attr"`
	Name        string `xml:"Name"`
	Description string `xml:"Description"`
}

// mapScheduleInfo maps a scheduleInfoEntry to a public ScheduleInfo.
func mapScheduleInfo(e scheduleInfoEntry) ScheduleInfo {
	return ScheduleInfo{
		Token:       e.Token,
		Name:        e.Name,
		Description: e.Description,
	}
}

// GetScheduleInfoList retrieves a paginated list of ScheduleInfo items.
func (c *Client) GetScheduleInfoList(ctx context.Context, limit *int, startReference *string) ([]*ScheduleInfo, string, error) {
	endpoint := c.getScheduleEndpoint()

	type GetScheduleInfoList struct {
		XMLName        xml.Name `xml:"tsc:GetScheduleInfoList"`
		Xmlns          string   `xml:"xmlns:tsc,attr"`
		Limit          *int     `xml:"tsc:Limit,omitempty"`
		StartReference *string  `xml:"tsc:StartReference,omitempty"`
	}

	type GetScheduleInfoListResponse struct {
		XMLName            xml.Name            `xml:"GetScheduleInfoListResponse"`
		NextStartReference string              `xml:"NextStartReference"`
		ScheduleInfo       []scheduleInfoEntry `xml:"ScheduleInfo"`
	}

	req := GetScheduleInfoList{
		Xmlns:          scheduleNamespace,
		Limit:          limit,
		StartReference: startReference,
	}

	var resp GetScheduleInfoListResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetScheduleInfoList failed: %w", err)
	}

	result := make([]*ScheduleInfo, 0, len(resp.ScheduleInfo))

	for _, e := range resp.ScheduleInfo {
		info := mapScheduleInfo(e)
		result = append(result, &info)
	}

	return result, resp.NextStartReference, nil
}

// GetScheduleInfo retrieves ScheduleInfo items by token.
func (c *Client) GetScheduleInfo(ctx context.Context, tokens []string) ([]*ScheduleInfo, error) {
	if len(tokens) == 0 {
		return nil, ErrInvalidScheduleToken
	}

	endpoint := c.getScheduleEndpoint()

	type GetScheduleInfo struct {
		XMLName xml.Name `xml:"tsc:GetScheduleInfo"`
		Xmlns   string   `xml:"xmlns:tsc,attr"`
		Token   []string `xml:"tsc:Token"`
	}

	type GetScheduleInfoResponse struct {
		XMLName      xml.Name            `xml:"GetScheduleInfoResponse"`
		ScheduleInfo []scheduleInfoEntry `xml:"ScheduleInfo"`
	}

	req := GetScheduleInfo{
		Xmlns: scheduleNamespace,
		Token: tokens,
	}

	var resp GetScheduleInfoResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetScheduleInfo failed: %w", err)
	}

	result := make([]*ScheduleInfo, 0, len(resp.ScheduleInfo))

	for _, e := range resp.ScheduleInfo {
		info := mapScheduleInfo(e)
		result = append(result, &info)
	}

	return result, nil
}

// scheduleEntry is the internal XML representation of a Schedule entry.
type scheduleEntry struct {
	Token       string `xml:"token,attr"`
	Name        string `xml:"Name"`
	Description string `xml:"Description"`
	Standard    string `xml:"Standard"`
	SpecialDays []struct {
		GroupToken string `xml:"GroupToken"`
		TimeRange  []struct {
			From  string `xml:"From"`
			Until string `xml:"Until"`
		} `xml:"TimeRange"`
	} `xml:"SpecialDays"`
}

// mapSchedule maps a scheduleEntry to a public Schedule.
func mapSchedule(e scheduleEntry) Schedule {
	s := Schedule{
		ScheduleInfo: ScheduleInfo{
			Token:       e.Token,
			Name:        e.Name,
			Description: e.Description,
		},
		Standard: e.Standard,
	}

	for _, sd := range e.SpecialDays {
		sds := SpecialDaysSchedule{
			GroupToken: sd.GroupToken,
		}

		for _, tr := range sd.TimeRange {
			sds.TimeRange = append(sds.TimeRange, TimePeriod{
				From:  tr.From,
				Until: tr.Until,
			})
		}

		s.SpecialDays = append(s.SpecialDays, sds)
	}

	return s
}

// GetScheduleList retrieves a paginated list of Schedule items.
func (c *Client) GetScheduleList(ctx context.Context, limit *int, startReference *string) ([]*Schedule, string, error) {
	endpoint := c.getScheduleEndpoint()

	type GetScheduleList struct {
		XMLName        xml.Name `xml:"tsc:GetScheduleList"`
		Xmlns          string   `xml:"xmlns:tsc,attr"`
		Limit          *int     `xml:"tsc:Limit,omitempty"`
		StartReference *string  `xml:"tsc:StartReference,omitempty"`
	}

	type GetScheduleListResponse struct {
		XMLName            xml.Name        `xml:"GetScheduleListResponse"`
		NextStartReference string          `xml:"NextStartReference"`
		Schedule           []scheduleEntry `xml:"Schedule"`
	}

	req := GetScheduleList{
		Xmlns:          scheduleNamespace,
		Limit:          limit,
		StartReference: startReference,
	}

	var resp GetScheduleListResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetScheduleList failed: %w", err)
	}

	result := make([]*Schedule, 0, len(resp.Schedule))

	for _, e := range resp.Schedule {
		sched := mapSchedule(e)
		result = append(result, &sched)
	}

	return result, resp.NextStartReference, nil
}

// GetSchedules retrieves Schedule items by token.
func (c *Client) GetSchedules(ctx context.Context, tokens []string) ([]*Schedule, error) {
	if len(tokens) == 0 {
		return nil, ErrInvalidScheduleToken
	}

	endpoint := c.getScheduleEndpoint()

	type GetSchedules struct {
		XMLName xml.Name `xml:"tsc:GetSchedules"`
		Xmlns   string   `xml:"xmlns:tsc,attr"`
		Token   []string `xml:"tsc:Token"`
	}

	type GetSchedulesResponse struct {
		XMLName  xml.Name        `xml:"GetSchedulesResponse"`
		Schedule []scheduleEntry `xml:"Schedule"`
	}

	req := GetSchedules{
		Xmlns: scheduleNamespace,
		Token: tokens,
	}

	var resp GetSchedulesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetSchedules failed: %w", err)
	}

	result := make([]*Schedule, 0, len(resp.Schedule))

	for _, e := range resp.Schedule {
		sched := mapSchedule(e)
		result = append(result, &sched)
	}

	return result, nil
}

// scheduleXMLPayload represents a Schedule in XML.
type scheduleXMLPayload struct {
	Token       string                   `xml:"token,attr,omitempty"`
	Name        string                   `xml:"tsc:Name"`
	Description string                   `xml:"tsc:Description,omitempty"`
	Standard    string                   `xml:"tsc:Standard"`
	SpecialDays []specialDaysXMLPayload  `xml:"tsc:SpecialDays,omitempty"`
}

// specialDaysXMLPayload represents SpecialDaysSchedule in XML.
type specialDaysXMLPayload struct {
	GroupToken string              `xml:"tsc:GroupToken"`
	TimeRange  []timePeriodPayload `xml:"tsc:TimeRange,omitempty"`
}

// timePeriodPayload represents a TimePeriod in XML.
type timePeriodPayload struct {
	From  string `xml:"tsc:From"`
	Until string `xml:"tsc:Until,omitempty"`
}

// buildSchedulePayload converts a Schedule to its XML payload form.
func buildSchedulePayload(s *Schedule) scheduleXMLPayload {
	p := scheduleXMLPayload{
		Token:       s.Token,
		Name:        s.Name,
		Description: s.Description,
		Standard:    s.Standard,
	}

	for _, sd := range s.SpecialDays {
		sdp := specialDaysXMLPayload{GroupToken: sd.GroupToken}

		for _, tr := range sd.TimeRange {
			sdp.TimeRange = append(sdp.TimeRange, timePeriodPayload{
				From:  tr.From,
				Until: tr.Until,
			})
		}

		p.SpecialDays = append(p.SpecialDays, sdp)
	}

	return p
}

// CreateSchedule creates a new schedule and returns its assigned token.
func (c *Client) CreateSchedule(ctx context.Context, schedule *Schedule) (string, error) {
	if schedule == nil {
		return "", ErrScheduleNil
	}

	endpoint := c.getScheduleEndpoint()

	type CreateSchedule struct {
		XMLName  xml.Name           `xml:"tsc:CreateSchedule"`
		Xmlns    string             `xml:"xmlns:tsc,attr"`
		Schedule scheduleXMLPayload `xml:"tsc:Schedule"`
	}

	type CreateScheduleResponse struct {
		XMLName xml.Name `xml:"CreateScheduleResponse"`
		Token   string   `xml:"Token"`
	}

	req := CreateSchedule{
		Xmlns:    scheduleNamespace,
		Schedule: buildSchedulePayload(schedule),
	}

	var resp CreateScheduleResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateSchedule failed: %w", err)
	}

	return resp.Token, nil
}

// ModifySchedule modifies an existing schedule.
func (c *Client) ModifySchedule(ctx context.Context, schedule *Schedule) error {
	if schedule == nil {
		return ErrScheduleNil
	}

	if schedule.Token == "" {
		return ErrInvalidScheduleToken
	}

	endpoint := c.getScheduleEndpoint()

	type ModifySchedule struct {
		XMLName  xml.Name           `xml:"tsc:ModifySchedule"`
		Xmlns    string             `xml:"xmlns:tsc,attr"`
		Schedule scheduleXMLPayload `xml:"tsc:Schedule"`
	}

	type ModifyScheduleResponse struct {
		XMLName xml.Name `xml:"ModifyScheduleResponse"`
	}

	req := ModifySchedule{
		Xmlns:    scheduleNamespace,
		Schedule: buildSchedulePayload(schedule),
	}

	var resp ModifyScheduleResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("ModifySchedule failed: %w", err)
	}

	return nil
}

// SetSchedule creates or replaces a schedule (requires ClientSuppliedTokenSupported capability).
func (c *Client) SetSchedule(ctx context.Context, schedule *Schedule) error {
	if schedule == nil {
		return ErrScheduleNil
	}

	endpoint := c.getScheduleEndpoint()

	type SetSchedule struct {
		XMLName  xml.Name           `xml:"tsc:SetSchedule"`
		Xmlns    string             `xml:"xmlns:tsc,attr"`
		Schedule scheduleXMLPayload `xml:"tsc:Schedule"`
	}

	type SetScheduleResponse struct {
		XMLName xml.Name `xml:"SetScheduleResponse"`
	}

	req := SetSchedule{
		Xmlns:    scheduleNamespace,
		Schedule: buildSchedulePayload(schedule),
	}

	var resp SetScheduleResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetSchedule failed: %w", err)
	}

	return nil
}

// DeleteSchedule deletes a schedule by token.
func (c *Client) DeleteSchedule(ctx context.Context, token string) error {
	if token == "" {
		return ErrInvalidScheduleToken
	}

	endpoint := c.getScheduleEndpoint()

	type DeleteSchedule struct {
		XMLName xml.Name `xml:"tsc:DeleteSchedule"`
		Xmlns   string   `xml:"xmlns:tsc,attr"`
		Token   string   `xml:"tsc:Token"`
	}

	type DeleteScheduleResponse struct {
		XMLName xml.Name `xml:"DeleteScheduleResponse"`
	}

	req := DeleteSchedule{
		Xmlns: scheduleNamespace,
		Token: token,
	}

	var resp DeleteScheduleResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteSchedule failed: %w", err)
	}

	return nil
}

// specialDayGroupInfoEntry is the internal XML representation of a SpecialDayGroupInfo.
type specialDayGroupInfoEntry struct {
	Token       string `xml:"token,attr"`
	Name        string `xml:"Name"`
	Description string `xml:"Description"`
}

// mapSpecialDayGroupInfo maps a specialDayGroupInfoEntry to a public SpecialDayGroupInfo.
func mapSpecialDayGroupInfo(e specialDayGroupInfoEntry) SpecialDayGroupInfo {
	return SpecialDayGroupInfo{
		Token:       e.Token,
		Name:        e.Name,
		Description: e.Description,
	}
}

// GetSpecialDayGroupInfoList retrieves a paginated list of SpecialDayGroupInfo items.
func (c *Client) GetSpecialDayGroupInfoList(ctx context.Context, limit *int, startReference *string) ([]*SpecialDayGroupInfo, string, error) {
	endpoint := c.getScheduleEndpoint()

	type GetSpecialDayGroupInfoList struct {
		XMLName        xml.Name `xml:"tsc:GetSpecialDayGroupInfoList"`
		Xmlns          string   `xml:"xmlns:tsc,attr"`
		Limit          *int     `xml:"tsc:Limit,omitempty"`
		StartReference *string  `xml:"tsc:StartReference,omitempty"`
	}

	type GetSpecialDayGroupInfoListResponse struct {
		XMLName              xml.Name                   `xml:"GetSpecialDayGroupInfoListResponse"`
		NextStartReference   string                     `xml:"NextStartReference"`
		SpecialDayGroupInfo  []specialDayGroupInfoEntry `xml:"SpecialDayGroupInfo"`
	}

	req := GetSpecialDayGroupInfoList{
		Xmlns:          scheduleNamespace,
		Limit:          limit,
		StartReference: startReference,
	}

	var resp GetSpecialDayGroupInfoListResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetSpecialDayGroupInfoList failed: %w", err)
	}

	result := make([]*SpecialDayGroupInfo, 0, len(resp.SpecialDayGroupInfo))

	for _, e := range resp.SpecialDayGroupInfo {
		info := mapSpecialDayGroupInfo(e)
		result = append(result, &info)
	}

	return result, resp.NextStartReference, nil
}

// GetSpecialDayGroupInfo retrieves SpecialDayGroupInfo items by token.
func (c *Client) GetSpecialDayGroupInfo(ctx context.Context, tokens []string) ([]*SpecialDayGroupInfo, error) {
	if len(tokens) == 0 {
		return nil, ErrInvalidSpecialDayGroupToken
	}

	endpoint := c.getScheduleEndpoint()

	type GetSpecialDayGroupInfo struct {
		XMLName xml.Name `xml:"tsc:GetSpecialDayGroupInfo"`
		Xmlns   string   `xml:"xmlns:tsc,attr"`
		Token   []string `xml:"tsc:Token"`
	}

	type GetSpecialDayGroupInfoResponse struct {
		XMLName             xml.Name                   `xml:"GetSpecialDayGroupInfoResponse"`
		SpecialDayGroupInfo []specialDayGroupInfoEntry `xml:"SpecialDayGroupInfo"`
	}

	req := GetSpecialDayGroupInfo{
		Xmlns: scheduleNamespace,
		Token: tokens,
	}

	var resp GetSpecialDayGroupInfoResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetSpecialDayGroupInfo failed: %w", err)
	}

	result := make([]*SpecialDayGroupInfo, 0, len(resp.SpecialDayGroupInfo))

	for _, e := range resp.SpecialDayGroupInfo {
		info := mapSpecialDayGroupInfo(e)
		result = append(result, &info)
	}

	return result, nil
}

// specialDayGroupEntry is the internal XML representation of a SpecialDayGroup.
type specialDayGroupEntry struct {
	Token       string `xml:"token,attr"`
	Name        string `xml:"Name"`
	Description string `xml:"Description"`
	Days        string `xml:"Days"`
}

// mapSpecialDayGroup maps a specialDayGroupEntry to a public SpecialDayGroup.
func mapSpecialDayGroup(e specialDayGroupEntry) SpecialDayGroup {
	return SpecialDayGroup{
		SpecialDayGroupInfo: SpecialDayGroupInfo{
			Token:       e.Token,
			Name:        e.Name,
			Description: e.Description,
		},
		Days: e.Days,
	}
}

// GetSpecialDayGroupList retrieves a paginated list of SpecialDayGroup items.
func (c *Client) GetSpecialDayGroupList(ctx context.Context, limit *int, startReference *string) ([]*SpecialDayGroup, string, error) {
	endpoint := c.getScheduleEndpoint()

	type GetSpecialDayGroupList struct {
		XMLName        xml.Name `xml:"tsc:GetSpecialDayGroupList"`
		Xmlns          string   `xml:"xmlns:tsc,attr"`
		Limit          *int     `xml:"tsc:Limit,omitempty"`
		StartReference *string  `xml:"tsc:StartReference,omitempty"`
	}

	type GetSpecialDayGroupListResponse struct {
		XMLName            xml.Name               `xml:"GetSpecialDayGroupListResponse"`
		NextStartReference string                 `xml:"NextStartReference"`
		SpecialDayGroup    []specialDayGroupEntry `xml:"SpecialDayGroup"`
	}

	req := GetSpecialDayGroupList{
		Xmlns:          scheduleNamespace,
		Limit:          limit,
		StartReference: startReference,
	}

	var resp GetSpecialDayGroupListResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetSpecialDayGroupList failed: %w", err)
	}

	result := make([]*SpecialDayGroup, 0, len(resp.SpecialDayGroup))

	for _, e := range resp.SpecialDayGroup {
		g := mapSpecialDayGroup(e)
		result = append(result, &g)
	}

	return result, resp.NextStartReference, nil
}

// GetSpecialDayGroups retrieves SpecialDayGroup items by token.
func (c *Client) GetSpecialDayGroups(ctx context.Context, tokens []string) ([]*SpecialDayGroup, error) {
	if len(tokens) == 0 {
		return nil, ErrInvalidSpecialDayGroupToken
	}

	endpoint := c.getScheduleEndpoint()

	type GetSpecialDayGroups struct {
		XMLName xml.Name `xml:"tsc:GetSpecialDayGroups"`
		Xmlns   string   `xml:"xmlns:tsc,attr"`
		Token   []string `xml:"tsc:Token"`
	}

	type GetSpecialDayGroupsResponse struct {
		XMLName         xml.Name               `xml:"GetSpecialDayGroupsResponse"`
		SpecialDayGroup []specialDayGroupEntry `xml:"SpecialDayGroup"`
	}

	req := GetSpecialDayGroups{
		Xmlns: scheduleNamespace,
		Token: tokens,
	}

	var resp GetSpecialDayGroupsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetSpecialDayGroups failed: %w", err)
	}

	result := make([]*SpecialDayGroup, 0, len(resp.SpecialDayGroup))

	for _, e := range resp.SpecialDayGroup {
		g := mapSpecialDayGroup(e)
		result = append(result, &g)
	}

	return result, nil
}

// specialDayGroupXMLPayload represents a SpecialDayGroup in XML.
type specialDayGroupXMLPayload struct {
	Token       string `xml:"token,attr,omitempty"`
	Name        string `xml:"tsc:Name"`
	Description string `xml:"tsc:Description,omitempty"`
	Days        string `xml:"tsc:Days,omitempty"`
}

// buildSpecialDayGroupPayload converts a SpecialDayGroup to its XML payload form.
func buildSpecialDayGroupPayload(g *SpecialDayGroup) specialDayGroupXMLPayload {
	return specialDayGroupXMLPayload{
		Token:       g.Token,
		Name:        g.Name,
		Description: g.Description,
		Days:        g.Days,
	}
}

// CreateSpecialDayGroup creates a new special day group and returns its assigned token.
func (c *Client) CreateSpecialDayGroup(ctx context.Context, group *SpecialDayGroup) (string, error) {
	if group == nil {
		return "", ErrSpecialDayGroupNil
	}

	endpoint := c.getScheduleEndpoint()

	type CreateSpecialDayGroup struct {
		XMLName         xml.Name                  `xml:"tsc:CreateSpecialDayGroup"`
		Xmlns           string                    `xml:"xmlns:tsc,attr"`
		SpecialDayGroup specialDayGroupXMLPayload `xml:"tsc:SpecialDayGroup"`
	}

	type CreateSpecialDayGroupResponse struct {
		XMLName xml.Name `xml:"CreateSpecialDayGroupResponse"`
		Token   string   `xml:"Token"`
	}

	req := CreateSpecialDayGroup{
		Xmlns:           scheduleNamespace,
		SpecialDayGroup: buildSpecialDayGroupPayload(group),
	}

	var resp CreateSpecialDayGroupResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateSpecialDayGroup failed: %w", err)
	}

	return resp.Token, nil
}

// ModifySpecialDayGroup modifies an existing special day group.
func (c *Client) ModifySpecialDayGroup(ctx context.Context, group *SpecialDayGroup) error {
	if group == nil {
		return ErrSpecialDayGroupNil
	}

	if group.Token == "" {
		return ErrInvalidSpecialDayGroupToken
	}

	endpoint := c.getScheduleEndpoint()

	type ModifySpecialDayGroup struct {
		XMLName         xml.Name                  `xml:"tsc:ModifySpecialDayGroup"`
		Xmlns           string                    `xml:"xmlns:tsc,attr"`
		SpecialDayGroup specialDayGroupXMLPayload `xml:"tsc:SpecialDayGroup"`
	}

	type ModifySpecialDayGroupResponse struct {
		XMLName xml.Name `xml:"ModifySpecialDayGroupResponse"`
	}

	req := ModifySpecialDayGroup{
		Xmlns:           scheduleNamespace,
		SpecialDayGroup: buildSpecialDayGroupPayload(group),
	}

	var resp ModifySpecialDayGroupResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("ModifySpecialDayGroup failed: %w", err)
	}

	return nil
}

// SetSpecialDayGroup creates or replaces a special day group (requires ClientSuppliedTokenSupported capability).
func (c *Client) SetSpecialDayGroup(ctx context.Context, group *SpecialDayGroup) error {
	if group == nil {
		return ErrSpecialDayGroupNil
	}

	endpoint := c.getScheduleEndpoint()

	type SetSpecialDayGroup struct {
		XMLName         xml.Name                  `xml:"tsc:SetSpecialDayGroup"`
		Xmlns           string                    `xml:"xmlns:tsc,attr"`
		SpecialDayGroup specialDayGroupXMLPayload `xml:"tsc:SpecialDayGroup"`
	}

	type SetSpecialDayGroupResponse struct {
		XMLName xml.Name `xml:"SetSpecialDayGroupResponse"`
	}

	req := SetSpecialDayGroup{
		Xmlns:           scheduleNamespace,
		SpecialDayGroup: buildSpecialDayGroupPayload(group),
	}

	var resp SetSpecialDayGroupResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetSpecialDayGroup failed: %w", err)
	}

	return nil
}

// DeleteSpecialDayGroup deletes a special day group by token.
func (c *Client) DeleteSpecialDayGroup(ctx context.Context, token string) error {
	if token == "" {
		return ErrInvalidSpecialDayGroupToken
	}

	endpoint := c.getScheduleEndpoint()

	type DeleteSpecialDayGroup struct {
		XMLName xml.Name `xml:"tsc:DeleteSpecialDayGroup"`
		Xmlns   string   `xml:"xmlns:tsc,attr"`
		Token   string   `xml:"tsc:Token"`
	}

	type DeleteSpecialDayGroupResponse struct {
		XMLName xml.Name `xml:"DeleteSpecialDayGroupResponse"`
	}

	req := DeleteSpecialDayGroup{
		Xmlns: scheduleNamespace,
		Token: token,
	}

	var resp DeleteSpecialDayGroupResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteSpecialDayGroup failed: %w", err)
	}

	return nil
}
