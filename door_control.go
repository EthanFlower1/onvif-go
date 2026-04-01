package onvif

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/0x524a/onvif-go/internal/soap"
)

// Door Control service namespace.
const doorControlNamespace = "http://www.onvif.org/ver10/doorcontrol/wsdl"

// Door Control service errors.
var (
	// ErrInvalidDoorToken is returned when a door token is empty.
	ErrInvalidDoorToken = errors.New("invalid door token: cannot be empty")
	// ErrDoorNil is returned when a door is nil.
	ErrDoorNil = errors.New("door cannot be nil")
)

// getDoorControlEndpoint returns the door control endpoint, falling back to device endpoint.
func (c *Client) getDoorControlEndpoint() string {
	if c.doorControlEndpoint != "" {
		return c.doorControlEndpoint
	}

	return c.endpoint
}

// GetDoorControlServiceCapabilities retrieves the capabilities of the door control service.
func (c *Client) GetDoorControlServiceCapabilities(ctx context.Context) (*DoorControlServiceCapabilities, error) {
	endpoint := c.getDoorControlEndpoint()

	type GetServiceCapabilities struct {
		XMLName xml.Name `xml:"tdc:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:tdc,attr"`
	}

	type GetServiceCapabilitiesResponse struct {
		XMLName      xml.Name `xml:"GetServiceCapabilitiesResponse"`
		Capabilities struct {
			MaxLimit                    uint `xml:"MaxLimit,attr"`
			MaxDoors                    uint `xml:"MaxDoors,attr"`
			ClientSuppliedTokenSupported bool `xml:"ClientSuppliedTokenSupported,attr"`
			DoorManagementSupported     bool `xml:"DoorManagementSupported,attr"`
		} `xml:"Capabilities"`
	}

	req := GetServiceCapabilities{
		Xmlns: doorControlNamespace,
	}

	var resp GetServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetDoorControlServiceCapabilities failed: %w", err)
	}

	return &DoorControlServiceCapabilities{
		MaxLimit:                    resp.Capabilities.MaxLimit,
		MaxDoors:                    resp.Capabilities.MaxDoors,
		ClientSuppliedTokenSupported: resp.Capabilities.ClientSuppliedTokenSupported,
		DoorManagementSupported:     resp.Capabilities.DoorManagementSupported,
	}, nil
}

// doorCapabilitiesXML is used internally for XML marshaling/unmarshaling of DoorCapabilities.
type doorCapabilitiesXML struct {
	Access               *bool `xml:"Access,attr"`
	AccessTimingOverride *bool `xml:"AccessTimingOverride,attr"`
	Lock                 *bool `xml:"Lock,attr"`
	Unlock               *bool `xml:"Unlock,attr"`
	Block                *bool `xml:"Block,attr"`
	DoubleLock           *bool `xml:"DoubleLock,attr"`
	LockDown             *bool `xml:"LockDown,attr"`
	LockOpen             *bool `xml:"LockOpen,attr"`
	DoorMonitor          *bool `xml:"DoorMonitor,attr"`
	LockMonitor          *bool `xml:"LockMonitor,attr"`
	DoubleLockMonitor    *bool `xml:"DoubleLockMonitor,attr"`
	Alarm                *bool `xml:"Alarm,attr"`
	Tamper               *bool `xml:"Tamper,attr"`
	Fault                *bool `xml:"Fault,attr"`
}

func mapDoorCapabilities(caps doorCapabilitiesXML) DoorCapabilities {
	return DoorCapabilities{
		Access:               caps.Access,
		AccessTimingOverride: caps.AccessTimingOverride,
		Lock:                 caps.Lock,
		Unlock:               caps.Unlock,
		Block:                caps.Block,
		DoubleLock:           caps.DoubleLock,
		LockDown:             caps.LockDown,
		LockOpen:             caps.LockOpen,
		DoorMonitor:          caps.DoorMonitor,
		LockMonitor:          caps.LockMonitor,
		DoubleLockMonitor:    caps.DoubleLockMonitor,
		Alarm:                caps.Alarm,
		Tamper:               caps.Tamper,
		Fault:                caps.Fault,
	}
}

// doorInfoEntry is the internal XML structure for a DoorInfo response element.
type doorInfoEntry struct {
	Token        string              `xml:"token,attr"`
	Name         string              `xml:"Name"`
	Description  string              `xml:"Description"`
	Capabilities doorCapabilitiesXML `xml:"Capabilities"`
}

func mapDoorInfo(e doorInfoEntry) DoorInfo {
	return DoorInfo{
		Token:        e.Token,
		Name:         e.Name,
		Description:  e.Description,
		Capabilities: mapDoorCapabilities(e.Capabilities),
	}
}

// doorEntry is the internal XML structure for a Door response element.
type doorEntry struct {
	Token       string              `xml:"token,attr"`
	Name        string              `xml:"Name"`
	Description string              `xml:"Description"`
	Capabilities doorCapabilitiesXML `xml:"Capabilities"`
	DoorType    string              `xml:"DoorType"`
	Timings     struct {
		ReleaseTime           string `xml:"ReleaseTime"`
		OpenTime              string `xml:"OpenTime"`
		ExtendedReleaseTime   string `xml:"ExtendedReleaseTime"`
		DelayTimeBeforeRelock string `xml:"DelayTimeBeforeRelock"`
		ExtendedOpenTime      string `xml:"ExtendedOpenTime"`
		PreAlarmTime          string `xml:"PreAlarmTime"`
	} `xml:"Timings"`
}

func mapDoor(e doorEntry) Door {
	return Door{
		DoorInfo: DoorInfo{
			Token:        e.Token,
			Name:         e.Name,
			Description:  e.Description,
			Capabilities: mapDoorCapabilities(e.Capabilities),
		},
		DoorType: e.DoorType,
		Timings: DoorTimings{
			ReleaseTime:           e.Timings.ReleaseTime,
			OpenTime:              e.Timings.OpenTime,
			ExtendedReleaseTime:   e.Timings.ExtendedReleaseTime,
			DelayTimeBeforeRelock: e.Timings.DelayTimeBeforeRelock,
			ExtendedOpenTime:      e.Timings.ExtendedOpenTime,
			PreAlarmTime:          e.Timings.PreAlarmTime,
		},
	}
}

// GetDoorInfoList retrieves a paginated list of DoorInfo items.
func (c *Client) GetDoorInfoList(ctx context.Context, limit *int, startReference *string) ([]*DoorInfo, string, error) {
	endpoint := c.getDoorControlEndpoint()

	type GetDoorInfoList struct {
		XMLName        xml.Name `xml:"tdc:GetDoorInfoList"`
		Xmlns          string   `xml:"xmlns:tdc,attr"`
		Limit          *int     `xml:"tdc:Limit,omitempty"`
		StartReference *string  `xml:"tdc:StartReference,omitempty"`
	}

	type GetDoorInfoListResponse struct {
		XMLName            xml.Name        `xml:"GetDoorInfoListResponse"`
		NextStartReference string          `xml:"NextStartReference"`
		DoorInfo           []doorInfoEntry `xml:"DoorInfo"`
	}

	req := GetDoorInfoList{
		Xmlns:          doorControlNamespace,
		Limit:          limit,
		StartReference: startReference,
	}

	var resp GetDoorInfoListResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetDoorInfoList failed: %w", err)
	}

	result := make([]*DoorInfo, 0, len(resp.DoorInfo))

	for _, entry := range resp.DoorInfo {
		info := mapDoorInfo(entry)
		result = append(result, &info)
	}

	return result, resp.NextStartReference, nil
}

// GetDoorInfo retrieves DoorInfo items by token.
func (c *Client) GetDoorInfo(ctx context.Context, tokens []string) ([]*DoorInfo, error) {
	if len(tokens) == 0 {
		return nil, ErrInvalidDoorToken
	}

	endpoint := c.getDoorControlEndpoint()

	type GetDoorInfo struct {
		XMLName xml.Name `xml:"tdc:GetDoorInfo"`
		Xmlns   string   `xml:"xmlns:tdc,attr"`
		Token   []string `xml:"tdc:Token"`
	}

	type GetDoorInfoResponse struct {
		XMLName  xml.Name        `xml:"GetDoorInfoResponse"`
		DoorInfo []doorInfoEntry `xml:"DoorInfo"`
	}

	req := GetDoorInfo{
		Xmlns: doorControlNamespace,
		Token: tokens,
	}

	var resp GetDoorInfoResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetDoorInfo failed: %w", err)
	}

	result := make([]*DoorInfo, 0, len(resp.DoorInfo))

	for _, entry := range resp.DoorInfo {
		info := mapDoorInfo(entry)
		result = append(result, &info)
	}

	return result, nil
}

// GetDoorList retrieves a paginated list of Door items.
func (c *Client) GetDoorList(ctx context.Context, limit *int, startReference *string) ([]*Door, string, error) {
	endpoint := c.getDoorControlEndpoint()

	type GetDoorList struct {
		XMLName        xml.Name `xml:"tdc:GetDoorList"`
		Xmlns          string   `xml:"xmlns:tdc,attr"`
		Limit          *int     `xml:"tdc:Limit,omitempty"`
		StartReference *string  `xml:"tdc:StartReference,omitempty"`
	}

	type GetDoorListResponse struct {
		XMLName            xml.Name    `xml:"GetDoorListResponse"`
		NextStartReference string      `xml:"NextStartReference"`
		Door               []doorEntry `xml:"Door"`
	}

	req := GetDoorList{
		Xmlns:          doorControlNamespace,
		Limit:          limit,
		StartReference: startReference,
	}

	var resp GetDoorListResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, "", fmt.Errorf("GetDoorList failed: %w", err)
	}

	result := make([]*Door, 0, len(resp.Door))

	for _, entry := range resp.Door {
		door := mapDoor(entry)
		result = append(result, &door)
	}

	return result, resp.NextStartReference, nil
}

// GetDoors retrieves Door items by token.
func (c *Client) GetDoors(ctx context.Context, tokens []string) ([]*Door, error) {
	if len(tokens) == 0 {
		return nil, ErrInvalidDoorToken
	}

	endpoint := c.getDoorControlEndpoint()

	type GetDoors struct {
		XMLName xml.Name `xml:"tdc:GetDoors"`
		Xmlns   string   `xml:"xmlns:tdc,attr"`
		Token   []string `xml:"tdc:Token"`
	}

	type GetDoorsResponse struct {
		XMLName xml.Name    `xml:"GetDoorsResponse"`
		Door    []doorEntry `xml:"Door"`
	}

	req := GetDoors{
		Xmlns: doorControlNamespace,
		Token: tokens,
	}

	var resp GetDoorsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetDoors failed: %w", err)
	}

	result := make([]*Door, 0, len(resp.Door))

	for _, entry := range resp.Door {
		door := mapDoor(entry)
		result = append(result, &door)
	}

	return result, nil
}

// CreateDoor creates a new door. Returns the token of the created door.
func (c *Client) CreateDoor(ctx context.Context, door *Door) (string, error) {
	if door == nil {
		return "", ErrDoorNil
	}

	endpoint := c.getDoorControlEndpoint()

	type CapabilitiesXML struct {
		Access               *bool `xml:"Access,attr,omitempty"`
		AccessTimingOverride *bool `xml:"AccessTimingOverride,attr,omitempty"`
		Lock                 *bool `xml:"Lock,attr,omitempty"`
		Unlock               *bool `xml:"Unlock,attr,omitempty"`
		Block                *bool `xml:"Block,attr,omitempty"`
		DoubleLock           *bool `xml:"DoubleLock,attr,omitempty"`
		LockDown             *bool `xml:"LockDown,attr,omitempty"`
		LockOpen             *bool `xml:"LockOpen,attr,omitempty"`
		DoorMonitor          *bool `xml:"DoorMonitor,attr,omitempty"`
		LockMonitor          *bool `xml:"LockMonitor,attr,omitempty"`
		DoubleLockMonitor    *bool `xml:"DoubleLockMonitor,attr,omitempty"`
		Alarm                *bool `xml:"Alarm,attr,omitempty"`
		Tamper               *bool `xml:"Tamper,attr,omitempty"`
		Fault                *bool `xml:"Fault,attr,omitempty"`
	}

	type TimingsXML struct {
		ReleaseTime           string `xml:"tdc:ReleaseTime"`
		OpenTime              string `xml:"tdc:OpenTime"`
		ExtendedReleaseTime   string `xml:"tdc:ExtendedReleaseTime,omitempty"`
		DelayTimeBeforeRelock string `xml:"tdc:DelayTimeBeforeRelock,omitempty"`
		ExtendedOpenTime      string `xml:"tdc:ExtendedOpenTime,omitempty"`
		PreAlarmTime          string `xml:"tdc:PreAlarmTime,omitempty"`
	}

	type DoorXML struct {
		Token        string          `xml:"token,attr,omitempty"`
		Name         string          `xml:"tdc:Name"`
		Description  string          `xml:"tdc:Description,omitempty"`
		Capabilities CapabilitiesXML `xml:"tdc:Capabilities"`
		DoorType     string          `xml:"tdc:DoorType,omitempty"`
		Timings      TimingsXML      `xml:"tdc:Timings"`
	}

	type CreateDoor struct {
		XMLName xml.Name `xml:"tdc:CreateDoor"`
		Xmlns   string   `xml:"xmlns:tdc,attr"`
		Door    DoorXML  `xml:"tdc:Door"`
	}

	type CreateDoorResponse struct {
		XMLName xml.Name `xml:"CreateDoorResponse"`
		Token   string   `xml:"Token"`
	}

	req := CreateDoor{
		Xmlns: doorControlNamespace,
		Door: DoorXML{
			Name:        door.Name,
			Description: door.Description,
			Capabilities: CapabilitiesXML{
				Access:               door.Capabilities.Access,
				AccessTimingOverride: door.Capabilities.AccessTimingOverride,
				Lock:                 door.Capabilities.Lock,
				Unlock:               door.Capabilities.Unlock,
				Block:                door.Capabilities.Block,
				DoubleLock:           door.Capabilities.DoubleLock,
				LockDown:             door.Capabilities.LockDown,
				LockOpen:             door.Capabilities.LockOpen,
				DoorMonitor:          door.Capabilities.DoorMonitor,
				LockMonitor:          door.Capabilities.LockMonitor,
				DoubleLockMonitor:    door.Capabilities.DoubleLockMonitor,
				Alarm:                door.Capabilities.Alarm,
				Tamper:               door.Capabilities.Tamper,
				Fault:                door.Capabilities.Fault,
			},
			DoorType: door.DoorType,
			Timings: TimingsXML{
				ReleaseTime:           door.Timings.ReleaseTime,
				OpenTime:              door.Timings.OpenTime,
				ExtendedReleaseTime:   door.Timings.ExtendedReleaseTime,
				DelayTimeBeforeRelock: door.Timings.DelayTimeBeforeRelock,
				ExtendedOpenTime:      door.Timings.ExtendedOpenTime,
				PreAlarmTime:          door.Timings.PreAlarmTime,
			},
		},
	}

	var resp CreateDoorResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateDoor failed: %w", err)
	}

	return resp.Token, nil
}

// SetDoor creates or replaces a door with a client-supplied token.
func (c *Client) SetDoor(ctx context.Context, door *Door) error {
	if door == nil {
		return ErrDoorNil
	}

	if strings.TrimSpace(door.Token) == "" {
		return ErrInvalidDoorToken
	}

	endpoint := c.getDoorControlEndpoint()

	type CapabilitiesXML struct {
		Access               *bool `xml:"Access,attr,omitempty"`
		AccessTimingOverride *bool `xml:"AccessTimingOverride,attr,omitempty"`
		Lock                 *bool `xml:"Lock,attr,omitempty"`
		Unlock               *bool `xml:"Unlock,attr,omitempty"`
		Block                *bool `xml:"Block,attr,omitempty"`
		DoubleLock           *bool `xml:"DoubleLock,attr,omitempty"`
		LockDown             *bool `xml:"LockDown,attr,omitempty"`
		LockOpen             *bool `xml:"LockOpen,attr,omitempty"`
		DoorMonitor          *bool `xml:"DoorMonitor,attr,omitempty"`
		LockMonitor          *bool `xml:"LockMonitor,attr,omitempty"`
		DoubleLockMonitor    *bool `xml:"DoubleLockMonitor,attr,omitempty"`
		Alarm                *bool `xml:"Alarm,attr,omitempty"`
		Tamper               *bool `xml:"Tamper,attr,omitempty"`
		Fault                *bool `xml:"Fault,attr,omitempty"`
	}

	type TimingsXML struct {
		ReleaseTime           string `xml:"tdc:ReleaseTime"`
		OpenTime              string `xml:"tdc:OpenTime"`
		ExtendedReleaseTime   string `xml:"tdc:ExtendedReleaseTime,omitempty"`
		DelayTimeBeforeRelock string `xml:"tdc:DelayTimeBeforeRelock,omitempty"`
		ExtendedOpenTime      string `xml:"tdc:ExtendedOpenTime,omitempty"`
		PreAlarmTime          string `xml:"tdc:PreAlarmTime,omitempty"`
	}

	type DoorXML struct {
		Token        string          `xml:"token,attr"`
		Name         string          `xml:"tdc:Name"`
		Description  string          `xml:"tdc:Description,omitempty"`
		Capabilities CapabilitiesXML `xml:"tdc:Capabilities"`
		DoorType     string          `xml:"tdc:DoorType,omitempty"`
		Timings      TimingsXML      `xml:"tdc:Timings"`
	}

	type SetDoor struct {
		XMLName xml.Name `xml:"tdc:SetDoor"`
		Xmlns   string   `xml:"xmlns:tdc,attr"`
		Door    DoorXML  `xml:"tdc:Door"`
	}

	type SetDoorResponse struct {
		XMLName xml.Name `xml:"SetDoorResponse"`
	}

	req := SetDoor{
		Xmlns: doorControlNamespace,
		Door: DoorXML{
			Token:       door.Token,
			Name:        door.Name,
			Description: door.Description,
			Capabilities: CapabilitiesXML{
				Access:               door.Capabilities.Access,
				AccessTimingOverride: door.Capabilities.AccessTimingOverride,
				Lock:                 door.Capabilities.Lock,
				Unlock:               door.Capabilities.Unlock,
				Block:                door.Capabilities.Block,
				DoubleLock:           door.Capabilities.DoubleLock,
				LockDown:             door.Capabilities.LockDown,
				LockOpen:             door.Capabilities.LockOpen,
				DoorMonitor:          door.Capabilities.DoorMonitor,
				LockMonitor:          door.Capabilities.LockMonitor,
				DoubleLockMonitor:    door.Capabilities.DoubleLockMonitor,
				Alarm:                door.Capabilities.Alarm,
				Tamper:               door.Capabilities.Tamper,
				Fault:                door.Capabilities.Fault,
			},
			DoorType: door.DoorType,
			Timings: TimingsXML{
				ReleaseTime:           door.Timings.ReleaseTime,
				OpenTime:              door.Timings.OpenTime,
				ExtendedReleaseTime:   door.Timings.ExtendedReleaseTime,
				DelayTimeBeforeRelock: door.Timings.DelayTimeBeforeRelock,
				ExtendedOpenTime:      door.Timings.ExtendedOpenTime,
				PreAlarmTime:          door.Timings.PreAlarmTime,
			},
		},
	}

	var resp SetDoorResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetDoor failed: %w", err)
	}

	return nil
}

// ModifyDoor modifies an existing door.
func (c *Client) ModifyDoor(ctx context.Context, door *Door) error {
	if door == nil {
		return ErrDoorNil
	}

	if strings.TrimSpace(door.Token) == "" {
		return ErrInvalidDoorToken
	}

	endpoint := c.getDoorControlEndpoint()

	type CapabilitiesXML struct {
		Access               *bool `xml:"Access,attr,omitempty"`
		AccessTimingOverride *bool `xml:"AccessTimingOverride,attr,omitempty"`
		Lock                 *bool `xml:"Lock,attr,omitempty"`
		Unlock               *bool `xml:"Unlock,attr,omitempty"`
		Block                *bool `xml:"Block,attr,omitempty"`
		DoubleLock           *bool `xml:"DoubleLock,attr,omitempty"`
		LockDown             *bool `xml:"LockDown,attr,omitempty"`
		LockOpen             *bool `xml:"LockOpen,attr,omitempty"`
		DoorMonitor          *bool `xml:"DoorMonitor,attr,omitempty"`
		LockMonitor          *bool `xml:"LockMonitor,attr,omitempty"`
		DoubleLockMonitor    *bool `xml:"DoubleLockMonitor,attr,omitempty"`
		Alarm                *bool `xml:"Alarm,attr,omitempty"`
		Tamper               *bool `xml:"Tamper,attr,omitempty"`
		Fault                *bool `xml:"Fault,attr,omitempty"`
	}

	type TimingsXML struct {
		ReleaseTime           string `xml:"tdc:ReleaseTime"`
		OpenTime              string `xml:"tdc:OpenTime"`
		ExtendedReleaseTime   string `xml:"tdc:ExtendedReleaseTime,omitempty"`
		DelayTimeBeforeRelock string `xml:"tdc:DelayTimeBeforeRelock,omitempty"`
		ExtendedOpenTime      string `xml:"tdc:ExtendedOpenTime,omitempty"`
		PreAlarmTime          string `xml:"tdc:PreAlarmTime,omitempty"`
	}

	type DoorXML struct {
		Token        string          `xml:"token,attr"`
		Name         string          `xml:"tdc:Name"`
		Description  string          `xml:"tdc:Description,omitempty"`
		Capabilities CapabilitiesXML `xml:"tdc:Capabilities"`
		DoorType     string          `xml:"tdc:DoorType,omitempty"`
		Timings      TimingsXML      `xml:"tdc:Timings"`
	}

	type ModifyDoor struct {
		XMLName xml.Name `xml:"tdc:ModifyDoor"`
		Xmlns   string   `xml:"xmlns:tdc,attr"`
		Door    DoorXML  `xml:"tdc:Door"`
	}

	type ModifyDoorResponse struct {
		XMLName xml.Name `xml:"ModifyDoorResponse"`
	}

	req := ModifyDoor{
		Xmlns: doorControlNamespace,
		Door: DoorXML{
			Token:       door.Token,
			Name:        door.Name,
			Description: door.Description,
			Capabilities: CapabilitiesXML{
				Access:               door.Capabilities.Access,
				AccessTimingOverride: door.Capabilities.AccessTimingOverride,
				Lock:                 door.Capabilities.Lock,
				Unlock:               door.Capabilities.Unlock,
				Block:                door.Capabilities.Block,
				DoubleLock:           door.Capabilities.DoubleLock,
				LockDown:             door.Capabilities.LockDown,
				LockOpen:             door.Capabilities.LockOpen,
				DoorMonitor:          door.Capabilities.DoorMonitor,
				LockMonitor:          door.Capabilities.LockMonitor,
				DoubleLockMonitor:    door.Capabilities.DoubleLockMonitor,
				Alarm:                door.Capabilities.Alarm,
				Tamper:               door.Capabilities.Tamper,
				Fault:                door.Capabilities.Fault,
			},
			DoorType: door.DoorType,
			Timings: TimingsXML{
				ReleaseTime:           door.Timings.ReleaseTime,
				OpenTime:              door.Timings.OpenTime,
				ExtendedReleaseTime:   door.Timings.ExtendedReleaseTime,
				DelayTimeBeforeRelock: door.Timings.DelayTimeBeforeRelock,
				ExtendedOpenTime:      door.Timings.ExtendedOpenTime,
				PreAlarmTime:          door.Timings.PreAlarmTime,
			},
		},
	}

	var resp ModifyDoorResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("ModifyDoor failed: %w", err)
	}

	return nil
}

// DeleteDoor deletes a door by token.
func (c *Client) DeleteDoor(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidDoorToken
	}

	endpoint := c.getDoorControlEndpoint()

	type DeleteDoor struct {
		XMLName xml.Name `xml:"tdc:DeleteDoor"`
		Xmlns   string   `xml:"xmlns:tdc,attr"`
		Token   string   `xml:"tdc:Token"`
	}

	type DeleteDoorResponse struct {
		XMLName xml.Name `xml:"DeleteDoorResponse"`
	}

	req := DeleteDoor{
		Xmlns: doorControlNamespace,
		Token: token,
	}

	var resp DeleteDoorResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteDoor failed: %w", err)
	}

	return nil
}

// GetDoorState retrieves the current state of a door by token.
func (c *Client) GetDoorState(ctx context.Context, token string) (*DoorState, error) {
	if strings.TrimSpace(token) == "" {
		return nil, ErrInvalidDoorToken
	}

	endpoint := c.getDoorControlEndpoint()

	type GetDoorState struct {
		XMLName xml.Name `xml:"tdc:GetDoorState"`
		Xmlns   string   `xml:"xmlns:tdc,attr"`
		Token   string   `xml:"tdc:Token"`
	}

	type DoorTamperXML struct {
		Reason string `xml:"Reason"`
		State  string `xml:"State"`
	}

	type DoorFaultXML struct {
		Reason string `xml:"Reason"`
		State  string `xml:"State"`
	}

	type DoorStateXML struct {
		DoorPhysicalState       string        `xml:"DoorPhysicalState"`
		LockPhysicalState       string        `xml:"LockPhysicalState"`
		DoubleLockPhysicalState string        `xml:"DoubleLockPhysicalState"`
		Alarm                   string        `xml:"Alarm"`
		Tamper                  *DoorTamperXML `xml:"Tamper"`
		Fault                   *DoorFaultXML  `xml:"Fault"`
		DoorMode                string        `xml:"DoorMode"`
	}

	type GetDoorStateResponse struct {
		XMLName   xml.Name     `xml:"GetDoorStateResponse"`
		DoorState DoorStateXML `xml:"DoorState"`
	}

	req := GetDoorState{
		Xmlns: doorControlNamespace,
		Token: token,
	}

	var resp GetDoorStateResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetDoorState failed: %w", err)
	}

	state := &DoorState{
		DoorMode: DoorMode(resp.DoorState.DoorMode),
	}

	if resp.DoorState.DoorPhysicalState != "" {
		v := DoorPhysicalState(resp.DoorState.DoorPhysicalState)
		state.DoorPhysicalState = &v
	}

	if resp.DoorState.LockPhysicalState != "" {
		v := LockPhysicalState(resp.DoorState.LockPhysicalState)
		state.LockPhysicalState = &v
	}

	if resp.DoorState.DoubleLockPhysicalState != "" {
		v := LockPhysicalState(resp.DoorState.DoubleLockPhysicalState)
		state.DoubleLockPhysicalState = &v
	}

	if resp.DoorState.Alarm != "" {
		v := DoorAlarmState(resp.DoorState.Alarm)
		state.Alarm = &v
	}

	if resp.DoorState.Tamper != nil {
		state.Tamper = &DoorTamper{
			Reason: resp.DoorState.Tamper.Reason,
			State:  DoorTamperState(resp.DoorState.Tamper.State),
		}
	}

	if resp.DoorState.Fault != nil {
		state.Fault = &DoorFault{
			Reason: resp.DoorState.Fault.Reason,
			State:  DoorFaultState(resp.DoorState.Fault.State),
		}
	}

	return state, nil
}

// AccessDoor grants momentary access to a door.
func (c *Client) AccessDoor(ctx context.Context, token string, useExtendedTime *bool, accessDuration *string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidDoorToken
	}

	endpoint := c.getDoorControlEndpoint()

	type AccessDoor struct {
		XMLName         xml.Name `xml:"tdc:AccessDoor"`
		Xmlns           string   `xml:"xmlns:tdc,attr"`
		Token           string   `xml:"tdc:Token"`
		UseExtendedTime *bool    `xml:"tdc:UseExtendedTime,omitempty"`
		AccessDuration  *string  `xml:"tdc:AccessDuration,omitempty"`
	}

	type AccessDoorResponse struct {
		XMLName xml.Name `xml:"AccessDoorResponse"`
	}

	req := AccessDoor{
		Xmlns:           doorControlNamespace,
		Token:           token,
		UseExtendedTime: useExtendedTime,
		AccessDuration:  accessDuration,
	}

	var resp AccessDoorResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("AccessDoor failed: %w", err)
	}

	return nil
}

// LockDoor locks a door by token.
func (c *Client) LockDoor(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidDoorToken
	}

	endpoint := c.getDoorControlEndpoint()

	type LockDoor struct {
		XMLName xml.Name `xml:"tdc:LockDoor"`
		Xmlns   string   `xml:"xmlns:tdc,attr"`
		Token   string   `xml:"tdc:Token"`
	}

	type LockDoorResponse struct {
		XMLName xml.Name `xml:"LockDoorResponse"`
	}

	req := LockDoor{
		Xmlns: doorControlNamespace,
		Token: token,
	}

	var resp LockDoorResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("LockDoor failed: %w", err)
	}

	return nil
}

// UnlockDoor unlocks a door by token.
func (c *Client) UnlockDoor(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidDoorToken
	}

	endpoint := c.getDoorControlEndpoint()

	type UnlockDoor struct {
		XMLName xml.Name `xml:"tdc:UnlockDoor"`
		Xmlns   string   `xml:"xmlns:tdc,attr"`
		Token   string   `xml:"tdc:Token"`
	}

	type UnlockDoorResponse struct {
		XMLName xml.Name `xml:"UnlockDoorResponse"`
	}

	req := UnlockDoor{
		Xmlns: doorControlNamespace,
		Token: token,
	}

	var resp UnlockDoorResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("UnlockDoor failed: %w", err)
	}

	return nil
}

// BlockDoor blocks a door by token, preventing access.
func (c *Client) BlockDoor(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidDoorToken
	}

	endpoint := c.getDoorControlEndpoint()

	type BlockDoor struct {
		XMLName xml.Name `xml:"tdc:BlockDoor"`
		Xmlns   string   `xml:"xmlns:tdc,attr"`
		Token   string   `xml:"tdc:Token"`
	}

	type BlockDoorResponse struct {
		XMLName xml.Name `xml:"BlockDoorResponse"`
	}

	req := BlockDoor{
		Xmlns: doorControlNamespace,
		Token: token,
	}

	var resp BlockDoorResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("BlockDoor failed: %w", err)
	}

	return nil
}

// LockDownDoor puts a door in lockdown mode by token.
func (c *Client) LockDownDoor(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidDoorToken
	}

	endpoint := c.getDoorControlEndpoint()

	type LockDownDoor struct {
		XMLName xml.Name `xml:"tdc:LockDownDoor"`
		Xmlns   string   `xml:"xmlns:tdc,attr"`
		Token   string   `xml:"tdc:Token"`
	}

	type LockDownDoorResponse struct {
		XMLName xml.Name `xml:"LockDownDoorResponse"`
	}

	req := LockDownDoor{
		Xmlns: doorControlNamespace,
		Token: token,
	}

	var resp LockDownDoorResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("LockDownDoor failed: %w", err)
	}

	return nil
}

// LockDownReleaseDoor releases a door from lockdown mode by token.
func (c *Client) LockDownReleaseDoor(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidDoorToken
	}

	endpoint := c.getDoorControlEndpoint()

	type LockDownReleaseDoor struct {
		XMLName xml.Name `xml:"tdc:LockDownReleaseDoor"`
		Xmlns   string   `xml:"xmlns:tdc,attr"`
		Token   string   `xml:"tdc:Token"`
	}

	type LockDownReleaseDoorResponse struct {
		XMLName xml.Name `xml:"LockDownReleaseDoorResponse"`
	}

	req := LockDownReleaseDoor{
		Xmlns: doorControlNamespace,
		Token: token,
	}

	var resp LockDownReleaseDoorResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("LockDownReleaseDoor failed: %w", err)
	}

	return nil
}

// LockOpenDoor puts a door in lock-open mode (forced unlocked) by token.
func (c *Client) LockOpenDoor(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidDoorToken
	}

	endpoint := c.getDoorControlEndpoint()

	type LockOpenDoor struct {
		XMLName xml.Name `xml:"tdc:LockOpenDoor"`
		Xmlns   string   `xml:"xmlns:tdc,attr"`
		Token   string   `xml:"tdc:Token"`
	}

	type LockOpenDoorResponse struct {
		XMLName xml.Name `xml:"LockOpenDoorResponse"`
	}

	req := LockOpenDoor{
		Xmlns: doorControlNamespace,
		Token: token,
	}

	var resp LockOpenDoorResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("LockOpenDoor failed: %w", err)
	}

	return nil
}

// LockOpenReleaseDoor releases a door from lock-open mode by token.
func (c *Client) LockOpenReleaseDoor(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidDoorToken
	}

	endpoint := c.getDoorControlEndpoint()

	type LockOpenReleaseDoor struct {
		XMLName xml.Name `xml:"tdc:LockOpenReleaseDoor"`
		Xmlns   string   `xml:"xmlns:tdc,attr"`
		Token   string   `xml:"tdc:Token"`
	}

	type LockOpenReleaseDoorResponse struct {
		XMLName xml.Name `xml:"LockOpenReleaseDoorResponse"`
	}

	req := LockOpenReleaseDoor{
		Xmlns: doorControlNamespace,
		Token: token,
	}

	var resp LockOpenReleaseDoorResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("LockOpenReleaseDoor failed: %w", err)
	}

	return nil
}

// DoubleLockDoor activates double-locking on a door by token.
func (c *Client) DoubleLockDoor(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidDoorToken
	}

	endpoint := c.getDoorControlEndpoint()

	type DoubleLockDoor struct {
		XMLName xml.Name `xml:"tdc:DoubleLockDoor"`
		Xmlns   string   `xml:"xmlns:tdc,attr"`
		Token   string   `xml:"tdc:Token"`
	}

	type DoubleLockDoorResponse struct {
		XMLName xml.Name `xml:"DoubleLockDoorResponse"`
	}

	req := DoubleLockDoor{
		Xmlns: doorControlNamespace,
		Token: token,
	}

	var resp DoubleLockDoorResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DoubleLockDoor failed: %w", err)
	}

	return nil
}
