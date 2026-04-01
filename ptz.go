package onvif

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/EthanFlower1/onvif-go/internal/soap"
)

// PTZ service namespace.
const ptzNamespace = "http://www.onvif.org/ver20/ptz/wsdl"

// ptzPanTiltXML is a shared type for PTZ pan/tilt XML serialization.
type ptzPanTiltXML struct {
	X     float64 `xml:"x,attr"`
	Y     float64 `xml:"y,attr"`
	Space string  `xml:"space,attr,omitempty"`
}

// ptzZoomXML is a shared type for PTZ zoom XML serialization.
type ptzZoomXML struct {
	X     float64 `xml:"x,attr"`
	Space string  `xml:"space,attr,omitempty"`
}

// ptzVectorXML is a shared type for PTZ position/velocity XML serialization.
type ptzVectorXML struct {
	PanTilt *ptzPanTiltXML `xml:"PanTilt,omitempty"`
	Zoom    *ptzZoomXML    `xml:"Zoom,omitempty"`
}

// ptzSpeedXML is a shared type for PTZ speed XML serialization.
type ptzSpeedXML struct {
	PanTilt *ptzPanTiltXML `xml:"PanTilt,omitempty"`
	Zoom    *ptzZoomXML    `xml:"Zoom,omitempty"`
}

// convertToPTZVectorXML converts PTZVector to XML struct.
func convertToPTZVectorXML(v *PTZVector) *ptzVectorXML {
	if v == nil {
		return nil
	}
	result := &ptzVectorXML{}
	if v.PanTilt != nil {
		result.PanTilt = &ptzPanTiltXML{X: v.PanTilt.X, Y: v.PanTilt.Y, Space: v.PanTilt.Space}
	}
	if v.Zoom != nil {
		result.Zoom = &ptzZoomXML{X: v.Zoom.X, Space: v.Zoom.Space}
	}

	return result
}

// convertToPTZSpeedXML converts PTZSpeed to XML struct.
func convertToPTZSpeedXML(s *PTZSpeed) *ptzSpeedXML {
	if s == nil {
		return nil
	}
	result := &ptzSpeedXML{}
	if s.PanTilt != nil {
		result.PanTilt = &ptzPanTiltXML{X: s.PanTilt.X, Y: s.PanTilt.Y, Space: s.PanTilt.Space}
	}
	if s.Zoom != nil {
		result.Zoom = &ptzZoomXML{X: s.Zoom.X, Space: s.Zoom.Space}
	}

	return result
}

// ContinuousMove starts continuous PTZ movement.
func (c *Client) ContinuousMove(ctx context.Context, profileToken string, velocity *PTZSpeed, timeout *string) error {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return ErrServiceNotSupported
	}

	type ContinuousMove struct {
		XMLName      xml.Name     `xml:"tptz:ContinuousMove"`
		Xmlns        string       `xml:"xmlns:tptz,attr"`
		ProfileToken string       `xml:"tptz:ProfileToken"`
		Velocity     *ptzSpeedXML `xml:"tptz:Velocity"`
		Timeout      *string      `xml:"tptz:Timeout,omitempty"`
	}

	req := ContinuousMove{
		Xmlns:        ptzNamespace,
		ProfileToken: profileToken,
		Velocity:     convertToPTZSpeedXML(velocity),
		Timeout:      timeout,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("ContinuousMove failed: %w", err)
	}

	return nil
}

// AbsoluteMove moves PTZ to an absolute position.
func (c *Client) AbsoluteMove(ctx context.Context, profileToken string, position *PTZVector, speed *PTZSpeed) error {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return ErrServiceNotSupported
	}

	type AbsoluteMove struct {
		XMLName      xml.Name      `xml:"tptz:AbsoluteMove"`
		Xmlns        string        `xml:"xmlns:tptz,attr"`
		ProfileToken string        `xml:"tptz:ProfileToken"`
		Position     *ptzVectorXML `xml:"tptz:Position"`
		Speed        *ptzSpeedXML  `xml:"tptz:Speed,omitempty"`
	}

	req := AbsoluteMove{
		Xmlns:        ptzNamespace,
		ProfileToken: profileToken,
		Position:     convertToPTZVectorXML(position),
		Speed:        convertToPTZSpeedXML(speed),
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("AbsoluteMove failed: %w", err)
	}

	return nil
}

// RelativeMove moves PTZ relative to current position.
func (c *Client) RelativeMove(ctx context.Context, profileToken string, translation *PTZVector, speed *PTZSpeed) error {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return ErrServiceNotSupported
	}

	type RelativeMove struct {
		XMLName      xml.Name      `xml:"tptz:RelativeMove"`
		Xmlns        string        `xml:"xmlns:tptz,attr"`
		ProfileToken string        `xml:"tptz:ProfileToken"`
		Translation  *ptzVectorXML `xml:"tptz:Translation"`
		Speed        *ptzSpeedXML  `xml:"tptz:Speed,omitempty"`
	}

	req := RelativeMove{
		Xmlns:        ptzNamespace,
		ProfileToken: profileToken,
		Translation:  convertToPTZVectorXML(translation),
		Speed:        convertToPTZSpeedXML(speed),
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("RelativeMove failed: %w", err)
	}

	return nil
}

// Stop stops PTZ movement.
func (c *Client) Stop(ctx context.Context, profileToken string, panTilt, zoom bool) error {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return ErrServiceNotSupported
	}

	type Stop struct {
		XMLName      xml.Name `xml:"tptz:Stop"`
		Xmlns        string   `xml:"xmlns:tptz,attr"`
		ProfileToken string   `xml:"tptz:ProfileToken"`
		PanTilt      *bool    `xml:"tptz:PanTilt,omitempty"`
		Zoom         *bool    `xml:"tptz:Zoom,omitempty"`
	}

	req := Stop{
		Xmlns:        ptzNamespace,
		ProfileToken: profileToken,
	}

	if panTilt {
		req.PanTilt = &panTilt
	}
	if zoom {
		req.Zoom = &zoom
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("Stop failed: %w", err)
	}

	return nil
}

// GetStatus retrieves PTZ status.
func (c *Client) GetStatus(ctx context.Context, profileToken string) (*PTZStatus, error) {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return nil, ErrServiceNotSupported
	}

	type GetStatus struct {
		XMLName      xml.Name `xml:"tptz:GetStatus"`
		Xmlns        string   `xml:"xmlns:tptz,attr"`
		ProfileToken string   `xml:"tptz:ProfileToken"`
	}

	type GetStatusResponse struct {
		XMLName   xml.Name `xml:"GetStatusResponse"`
		PTZStatus struct {
			Position *struct {
				PanTilt *struct {
					X     float64 `xml:"x,attr"`
					Y     float64 `xml:"y,attr"`
					Space string  `xml:"space,attr,omitempty"`
				} `xml:"PanTilt"`
				Zoom *struct {
					X     float64 `xml:"x,attr"`
					Space string  `xml:"space,attr,omitempty"`
				} `xml:"Zoom"`
			} `xml:"Position"`
			MoveStatus *struct {
				PanTilt string `xml:"PanTilt"`
				Zoom    string `xml:"Zoom"`
			} `xml:"MoveStatus"`
			Error   string `xml:"Error"`
			UTCTime string `xml:"UtcTime"`
		} `xml:"PTZStatus"`
	}

	req := GetStatus{
		Xmlns:        ptzNamespace,
		ProfileToken: profileToken,
	}

	var resp GetStatusResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetStatus failed: %w", err)
	}

	status := &PTZStatus{
		Error: resp.PTZStatus.Error,
	}

	if resp.PTZStatus.Position != nil {
		status.Position = &PTZVector{}
		if resp.PTZStatus.Position.PanTilt != nil {
			status.Position.PanTilt = &Vector2D{
				X:     resp.PTZStatus.Position.PanTilt.X,
				Y:     resp.PTZStatus.Position.PanTilt.Y,
				Space: resp.PTZStatus.Position.PanTilt.Space,
			}
		}
		if resp.PTZStatus.Position.Zoom != nil {
			status.Position.Zoom = &Vector1D{
				X:     resp.PTZStatus.Position.Zoom.X,
				Space: resp.PTZStatus.Position.Zoom.Space,
			}
		}
	}

	if resp.PTZStatus.MoveStatus != nil {
		status.MoveStatus = &PTZMoveStatus{
			PanTilt: resp.PTZStatus.MoveStatus.PanTilt,
			Zoom:    resp.PTZStatus.MoveStatus.Zoom,
		}
	}

	return status, nil
}

// GetPresets retrieves PTZ presets.
func (c *Client) GetPresets(ctx context.Context, profileToken string) ([]*PTZPreset, error) {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return nil, ErrServiceNotSupported
	}

	type GetPresets struct {
		XMLName      xml.Name `xml:"tptz:GetPresets"`
		Xmlns        string   `xml:"xmlns:tptz,attr"`
		ProfileToken string   `xml:"tptz:ProfileToken"`
	}

	type GetPresetsResponse struct {
		XMLName xml.Name `xml:"GetPresetsResponse"`
		Preset  []struct {
			Token       string `xml:"token,attr"`
			Name        string `xml:"Name"`
			PTZPosition *struct {
				PanTilt *struct {
					X     float64 `xml:"x,attr"`
					Y     float64 `xml:"y,attr"`
					Space string  `xml:"space,attr,omitempty"`
				} `xml:"PanTilt"`
				Zoom *struct {
					X     float64 `xml:"x,attr"`
					Space string  `xml:"space,attr,omitempty"`
				} `xml:"Zoom"`
			} `xml:"PTZPosition"`
		} `xml:"Preset"`
	}

	req := GetPresets{
		Xmlns:        ptzNamespace,
		ProfileToken: profileToken,
	}

	var resp GetPresetsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetPresets failed: %w", err)
	}

	presets := make([]*PTZPreset, len(resp.Preset))
	for i, p := range resp.Preset {
		preset := &PTZPreset{
			Token: p.Token,
			Name:  p.Name,
		}

		if p.PTZPosition != nil {
			preset.PTZPosition = &PTZVector{}
			if p.PTZPosition.PanTilt != nil {
				preset.PTZPosition.PanTilt = &Vector2D{
					X:     p.PTZPosition.PanTilt.X,
					Y:     p.PTZPosition.PanTilt.Y,
					Space: p.PTZPosition.PanTilt.Space,
				}
			}
			if p.PTZPosition.Zoom != nil {
				preset.PTZPosition.Zoom = &Vector1D{
					X:     p.PTZPosition.Zoom.X,
					Space: p.PTZPosition.Zoom.Space,
				}
			}
		}

		presets[i] = preset
	}

	return presets, nil
}

// GotoPreset moves PTZ to a preset position.
func (c *Client) GotoPreset(ctx context.Context, profileToken, presetToken string, speed *PTZSpeed) error {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return ErrServiceNotSupported
	}

	type GotoPreset struct {
		XMLName      xml.Name     `xml:"tptz:GotoPreset"`
		Xmlns        string       `xml:"xmlns:tptz,attr"`
		ProfileToken string       `xml:"tptz:ProfileToken"`
		PresetToken  string       `xml:"tptz:PresetToken"`
		Speed        *ptzSpeedXML `xml:"tptz:Speed,omitempty"`
	}

	req := GotoPreset{
		Xmlns:        ptzNamespace,
		ProfileToken: profileToken,
		PresetToken:  presetToken,
		Speed:        convertToPTZSpeedXML(speed),
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("GotoPreset failed: %w", err)
	}

	return nil
}

// SetPreset sets a preset position.
func (c *Client) SetPreset(ctx context.Context, profileToken, presetName, presetToken string) (string, error) {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return "", ErrServiceNotSupported
	}

	type SetPreset struct {
		XMLName      xml.Name `xml:"tptz:SetPreset"`
		Xmlns        string   `xml:"xmlns:tptz,attr"`
		ProfileToken string   `xml:"tptz:ProfileToken"`
		PresetName   *string  `xml:"tptz:PresetName,omitempty"`
		PresetToken  *string  `xml:"tptz:PresetToken,omitempty"`
	}

	type SetPresetResponse struct {
		XMLName     xml.Name `xml:"SetPresetResponse"`
		PresetToken string   `xml:"PresetToken"`
	}

	req := SetPreset{
		Xmlns:        ptzNamespace,
		ProfileToken: profileToken,
	}

	if presetName != "" {
		req.PresetName = &presetName
	}
	if presetToken != "" {
		req.PresetToken = &presetToken
	}

	var resp SetPresetResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("SetPreset failed: %w", err)
	}

	return resp.PresetToken, nil
}

// RemovePreset removes a preset.
func (c *Client) RemovePreset(ctx context.Context, profileToken, presetToken string) error {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return ErrServiceNotSupported
	}

	type RemovePreset struct {
		XMLName      xml.Name `xml:"tptz:RemovePreset"`
		Xmlns        string   `xml:"xmlns:tptz,attr"`
		ProfileToken string   `xml:"tptz:ProfileToken"`
		PresetToken  string   `xml:"tptz:PresetToken"`
	}

	req := RemovePreset{
		Xmlns:        ptzNamespace,
		ProfileToken: profileToken,
		PresetToken:  presetToken,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("RemovePreset failed: %w", err)
	}

	return nil
}

// GotoHomePosition moves PTZ to home position.
func (c *Client) GotoHomePosition(ctx context.Context, profileToken string, speed *PTZSpeed) error {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return ErrServiceNotSupported
	}

	type GotoHomePosition struct {
		XMLName      xml.Name     `xml:"tptz:GotoHomePosition"`
		Xmlns        string       `xml:"xmlns:tptz,attr"`
		ProfileToken string       `xml:"tptz:ProfileToken"`
		Speed        *ptzSpeedXML `xml:"tptz:Speed,omitempty"`
	}

	req := GotoHomePosition{
		Xmlns:        ptzNamespace,
		ProfileToken: profileToken,
		Speed:        convertToPTZSpeedXML(speed),
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("GotoHomePosition failed: %w", err)
	}

	return nil
}

// SetHomePosition sets the current position as home position.
func (c *Client) SetHomePosition(ctx context.Context, profileToken string) error {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return ErrServiceNotSupported
	}

	type SetHomePosition struct {
		XMLName      xml.Name `xml:"tptz:SetHomePosition"`
		Xmlns        string   `xml:"xmlns:tptz,attr"`
		ProfileToken string   `xml:"tptz:ProfileToken"`
	}

	req := SetHomePosition{
		Xmlns:        ptzNamespace,
		ProfileToken: profileToken,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("SetHomePosition failed: %w", err)
	}

	return nil
}

// GetConfiguration retrieves PTZ configuration.
func (c *Client) GetConfiguration(ctx context.Context, configurationToken string) (*PTZConfiguration, error) {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return nil, ErrServiceNotSupported
	}

	type GetConfiguration struct {
		XMLName               xml.Name `xml:"tptz:GetConfiguration"`
		Xmlns                 string   `xml:"xmlns:tptz,attr"`
		PTZConfigurationToken string   `xml:"tptz:PTZConfigurationToken"`
	}

	type GetConfigurationResponse struct {
		XMLName          xml.Name `xml:"GetConfigurationResponse"`
		PTZConfiguration struct {
			Token     string `xml:"token,attr"`
			Name      string `xml:"Name"`
			UseCount  int    `xml:"UseCount"`
			NodeToken string `xml:"NodeToken"`
		} `xml:"PTZConfiguration"`
	}

	req := GetConfiguration{
		Xmlns:                 ptzNamespace,
		PTZConfigurationToken: configurationToken,
	}

	var resp GetConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetConfiguration failed: %w", err)
	}

	return &PTZConfiguration{
		Token:     resp.PTZConfiguration.Token,
		Name:      resp.PTZConfiguration.Name,
		UseCount:  resp.PTZConfiguration.UseCount,
		NodeToken: resp.PTZConfiguration.NodeToken,
	}, nil
}

// GetConfigurations retrieves all PTZ configurations.
func (c *Client) GetConfigurations(ctx context.Context) ([]*PTZConfiguration, error) {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return nil, ErrServiceNotSupported
	}

	type GetConfigurations struct {
		XMLName xml.Name `xml:"tptz:GetConfigurations"`
		Xmlns   string   `xml:"xmlns:tptz,attr"`
	}

	type GetConfigurationsResponse struct {
		XMLName          xml.Name `xml:"GetConfigurationsResponse"`
		PTZConfiguration []struct {
			Token     string `xml:"token,attr"`
			Name      string `xml:"Name"`
			UseCount  int    `xml:"UseCount"`
			NodeToken string `xml:"NodeToken"`
		} `xml:"PTZConfiguration"`
	}

	req := GetConfigurations{
		Xmlns: ptzNamespace,
	}

	var resp GetConfigurationsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetConfigurations failed: %w", err)
	}

	configs := make([]*PTZConfiguration, len(resp.PTZConfiguration))
	for i, cfg := range resp.PTZConfiguration {
		configs[i] = &PTZConfiguration{
			Token:     cfg.Token,
			Name:      cfg.Name,
			UseCount:  cfg.UseCount,
			NodeToken: cfg.NodeToken,
		}
	}

	return configs, nil
}

// GetNodes retrieves all PTZ nodes on the device.
func (c *Client) GetNodes(ctx context.Context) ([]*PTZNode, error) {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return nil, ErrServiceNotSupported
	}

	type GetNodes struct {
		XMLName xml.Name `xml:"tptz:GetNodes"`
		Xmlns   string   `xml:"xmlns:tptz,attr"`
	}

	type GetNodesResponse struct {
		XMLName xml.Name `xml:"GetNodesResponse"`
		PTZNode []struct {
			Token                  string   `xml:"token,attr"`
			Name                   string   `xml:"Name"`
			HomeSupported          bool     `xml:"HomeSupported"`
			MaximumNumberOfPresets int      `xml:"MaximumNumberOfPresets"`
			AuxiliaryCommands      []string `xml:"AuxiliaryCommands"`
		} `xml:"PTZNode"`
	}

	req := GetNodes{
		Xmlns: ptzNamespace,
	}

	var resp GetNodesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetNodes failed: %w", err)
	}

	nodes := make([]*PTZNode, len(resp.PTZNode))
	for i, n := range resp.PTZNode {
		nodes[i] = &PTZNode{
			Token:                  n.Token,
			Name:                   n.Name,
			HomeSupported:          n.HomeSupported,
			MaximumNumberOfPresets: n.MaximumNumberOfPresets,
			AuxiliaryCommands:      n.AuxiliaryCommands,
		}
	}

	return nodes, nil
}

// GetNode retrieves a specific PTZ node by token.
func (c *Client) GetNode(ctx context.Context, nodeToken string) (*PTZNode, error) {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return nil, ErrServiceNotSupported
	}

	type GetNode struct {
		XMLName   xml.Name `xml:"tptz:GetNode"`
		Xmlns     string   `xml:"xmlns:tptz,attr"`
		NodeToken string   `xml:"tptz:NodeToken"`
	}

	type GetNodeResponse struct {
		XMLName xml.Name `xml:"GetNodeResponse"`
		PTZNode struct {
			Token                  string   `xml:"token,attr"`
			Name                   string   `xml:"Name"`
			HomeSupported          bool     `xml:"HomeSupported"`
			MaximumNumberOfPresets int      `xml:"MaximumNumberOfPresets"`
			AuxiliaryCommands      []string `xml:"AuxiliaryCommands"`
		} `xml:"PTZNode"`
	}

	req := GetNode{
		Xmlns:     ptzNamespace,
		NodeToken: nodeToken,
	}

	var resp GetNodeResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetNode failed: %w", err)
	}

	return &PTZNode{
		Token:                  resp.PTZNode.Token,
		Name:                   resp.PTZNode.Name,
		HomeSupported:          resp.PTZNode.HomeSupported,
		MaximumNumberOfPresets: resp.PTZNode.MaximumNumberOfPresets,
		AuxiliaryCommands:      resp.PTZNode.AuxiliaryCommands,
	}, nil
}

// GetPTZConfigurationOptions retrieves PTZ configuration options for a given configuration token.
func (c *Client) GetPTZConfigurationOptions(ctx context.Context, configurationToken string) (*PTZConfigurationOptions, error) {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return nil, ErrServiceNotSupported
	}

	type GetPTZConfigurationOptions struct {
		XMLName            xml.Name `xml:"tptz:GetConfigurationOptions"`
		Xmlns              string   `xml:"xmlns:tptz,attr"`
		ConfigurationToken string   `xml:"tptz:ConfigurationToken"`
	}

	type GetPTZConfigurationOptionsResponse struct {
		XMLName              xml.Name `xml:"GetConfigurationOptionsResponse"`
		PTZConfigurationOptions struct {
			PTZTimeout *struct {
				Min string `xml:"Min"`
				Max string `xml:"Max"`
			} `xml:"PTZTimeout"`
		} `xml:"PTZConfigurationOptions"`
	}

	req := GetPTZConfigurationOptions{
		Xmlns:              ptzNamespace,
		ConfigurationToken: configurationToken,
	}

	var resp GetPTZConfigurationOptionsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetPTZConfigurationOptions failed: %w", err)
	}

	opts := &PTZConfigurationOptions{}
	if resp.PTZConfigurationOptions.PTZTimeout != nil {
		opts.PTZTimeout = &struct {
			Min string
			Max string
		}{
			Min: resp.PTZConfigurationOptions.PTZTimeout.Min,
			Max: resp.PTZConfigurationOptions.PTZTimeout.Max,
		}
	}

	return opts, nil
}

// SetPTZConfiguration sets a PTZ configuration on the device.
func (c *Client) SetPTZConfiguration(ctx context.Context, config *PTZConfiguration, forcePersistence bool) error {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return ErrServiceNotSupported
	}

	type setPTZConfigurationXML struct {
		Token     string `xml:"token,attr"`
		Name      string `xml:"tt:Name"`
		NodeToken string `xml:"tt:NodeToken"`
	}

	type SetPTZConfiguration struct {
		XMLName          xml.Name               `xml:"tptz:SetConfiguration"`
		Xmlns            string                 `xml:"xmlns:tptz,attr"`
		XmlnsTT          string                 `xml:"xmlns:tt,attr"`
		PTZConfiguration setPTZConfigurationXML `xml:"tptz:PTZConfiguration"`
		ForcePersistence bool                   `xml:"tptz:ForcePersistence"`
	}

	req := SetPTZConfiguration{
		Xmlns:   ptzNamespace,
		XmlnsTT: "http://www.onvif.org/ver10/schema",
		PTZConfiguration: setPTZConfigurationXML{
			Token:     config.Token,
			Name:      config.Name,
			NodeToken: config.NodeToken,
		},
		ForcePersistence: forcePersistence,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("SetPTZConfiguration failed: %w", err)
	}

	return nil
}

// GetPTZServiceCapabilities retrieves the capabilities of the PTZ service.
func (c *Client) GetPTZServiceCapabilities(ctx context.Context) (*PTZServiceCapabilities, error) {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return nil, ErrServiceNotSupported
	}

	type GetServiceCapabilities struct {
		XMLName xml.Name `xml:"tptz:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:tptz,attr"`
	}

	type GetServiceCapabilitiesResponse struct {
		XMLName      xml.Name `xml:"GetServiceCapabilitiesResponse"`
		Capabilities struct {
			EFlip   bool `xml:"EFlip,attr"`
			Reverse bool `xml:"Reverse,attr"`
		} `xml:"Capabilities"`
	}

	req := GetServiceCapabilities{
		Xmlns: ptzNamespace,
	}

	var resp GetServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetPTZServiceCapabilities failed: %w", err)
	}

	return &PTZServiceCapabilities{
		EFlip:   resp.Capabilities.EFlip,
		Reverse: resp.Capabilities.Reverse,
	}, nil
}

// GetPresetTours retrieves all preset tours for a profile.
func (c *Client) GetPresetTours(ctx context.Context, profileToken string) ([]*PresetTour, error) {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return nil, ErrServiceNotSupported
	}

	type GetPresetTours struct {
		XMLName      xml.Name `xml:"tptz:GetPresetTours"`
		Xmlns        string   `xml:"xmlns:tptz,attr"`
		ProfileToken string   `xml:"tptz:ProfileToken"`
	}

	type GetPresetToursResponse struct {
		XMLName    xml.Name `xml:"GetPresetToursResponse"`
		PresetTour []struct {
			Token     string `xml:"token,attr"`
			Name      string `xml:"Name"`
			Status    struct {
				State string `xml:"State"`
			} `xml:"Status"`
			AutoStart bool `xml:"AutoStart"`
		} `xml:"PresetTour"`
	}

	req := GetPresetTours{
		Xmlns:        ptzNamespace,
		ProfileToken: profileToken,
	}

	var resp GetPresetToursResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetPresetTours failed: %w", err)
	}

	tours := make([]*PresetTour, len(resp.PresetTour))
	for i, t := range resp.PresetTour {
		tours[i] = &PresetTour{
			Token:     t.Token,
			Name:      t.Name,
			Status:    t.Status.State,
			AutoStart: t.AutoStart,
		}
	}

	return tours, nil
}

// GetPresetTour retrieves a specific preset tour by token.
func (c *Client) GetPresetTour(ctx context.Context, profileToken, presetTourToken string) (*PresetTour, error) {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return nil, ErrServiceNotSupported
	}

	type GetPresetTour struct {
		XMLName         xml.Name `xml:"tptz:GetPresetTour"`
		Xmlns           string   `xml:"xmlns:tptz,attr"`
		ProfileToken    string   `xml:"tptz:ProfileToken"`
		PresetTourToken string   `xml:"tptz:PresetTourToken"`
	}

	type GetPresetTourResponse struct {
		XMLName    xml.Name `xml:"GetPresetTourResponse"`
		PresetTour struct {
			Token     string `xml:"token,attr"`
			Name      string `xml:"Name"`
			Status    struct {
				State string `xml:"State"`
			} `xml:"Status"`
			AutoStart bool `xml:"AutoStart"`
		} `xml:"PresetTour"`
	}

	req := GetPresetTour{
		Xmlns:           ptzNamespace,
		ProfileToken:    profileToken,
		PresetTourToken: presetTourToken,
	}

	var resp GetPresetTourResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetPresetTour failed: %w", err)
	}

	return &PresetTour{
		Token:     resp.PresetTour.Token,
		Name:      resp.PresetTour.Name,
		Status:    resp.PresetTour.Status.State,
		AutoStart: resp.PresetTour.AutoStart,
	}, nil
}

// GetPresetTourOptions retrieves available options for a preset tour.
func (c *Client) GetPresetTourOptions(ctx context.Context, profileToken string, presetTourToken string) (*PTZPresetTourOptions, error) {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return nil, ErrServiceNotSupported
	}

	type GetPresetTourOptions struct {
		XMLName         xml.Name `xml:"tptz:GetPresetTourOptions"`
		Xmlns           string   `xml:"xmlns:tptz,attr"`
		ProfileToken    string   `xml:"tptz:ProfileToken"`
		PresetTourToken *string  `xml:"tptz:PresetTourToken,omitempty"`
	}

	type GetPresetTourOptionsResponse struct {
		XMLName xml.Name `xml:"GetPresetTourOptionsResponse"`
		Options struct {
			AutoStart         bool `xml:"AutoStart"`
			StartingCondition *struct {
				RecurringTimeRange *struct {
					Min int `xml:"Min"`
					Max int `xml:"Max"`
				} `xml:"RecurringTimeRange"`
				RecurringDurationRange *struct {
					Min string `xml:"Min"`
					Max string `xml:"Max"`
				} `xml:"RecurringDurationRange"`
			} `xml:"StartingCondition"`
		} `xml:"Options"`
	}

	req := GetPresetTourOptions{
		Xmlns:        ptzNamespace,
		ProfileToken: profileToken,
	}
	if presetTourToken != "" {
		req.PresetTourToken = &presetTourToken
	}

	var resp GetPresetTourOptionsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetPresetTourOptions failed: %w", err)
	}

	opts := &PTZPresetTourOptions{
		AutoStart: resp.Options.AutoStart,
	}

	if resp.Options.StartingCondition != nil {
		opts.StartingCondition = &PresetTourStartingConditionOptions{}

		if resp.Options.StartingCondition.RecurringTimeRange != nil {
			opts.StartingCondition.RecurringTimeRange = &IntRange{
				Min: resp.Options.StartingCondition.RecurringTimeRange.Min,
				Max: resp.Options.StartingCondition.RecurringTimeRange.Max,
			}
		}

		if resp.Options.StartingCondition.RecurringDurationRange != nil {
			opts.StartingCondition.RecurringDurationRange = &DurationRange{
				Min: resp.Options.StartingCondition.RecurringDurationRange.Min,
				Max: resp.Options.StartingCondition.RecurringDurationRange.Max,
			}
		}
	}

	return opts, nil
}

// CreatePresetTour creates a new preset tour and returns its token.
func (c *Client) CreatePresetTour(ctx context.Context, profileToken string) (string, error) {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return "", ErrServiceNotSupported
	}

	type CreatePresetTour struct {
		XMLName      xml.Name `xml:"tptz:CreatePresetTour"`
		Xmlns        string   `xml:"xmlns:tptz,attr"`
		ProfileToken string   `xml:"tptz:ProfileToken"`
	}

	type CreatePresetTourResponse struct {
		XMLName         xml.Name `xml:"CreatePresetTourResponse"`
		PresetTourToken string   `xml:"PresetTourToken"`
	}

	req := CreatePresetTour{
		Xmlns:        ptzNamespace,
		ProfileToken: profileToken,
	}

	var resp CreatePresetTourResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreatePresetTour failed: %w", err)
	}

	return resp.PresetTourToken, nil
}

// ModifyPresetTour modifies an existing preset tour.
func (c *Client) ModifyPresetTour(ctx context.Context, profileToken string, presetTour *PresetTour) error {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return ErrServiceNotSupported
	}

	type presetTourXML struct {
		Token     string `xml:"token,attr"`
		Name      string `xml:"tt:Name,omitempty"`
		AutoStart bool   `xml:"tt:AutoStart"`
	}

	type ModifyPresetTour struct {
		XMLName      xml.Name      `xml:"tptz:ModifyPresetTour"`
		Xmlns        string        `xml:"xmlns:tptz,attr"`
		XmlnsTT      string        `xml:"xmlns:tt,attr"`
		ProfileToken string        `xml:"tptz:ProfileToken"`
		PresetTour   presetTourXML `xml:"tptz:PresetTour"`
	}

	req := ModifyPresetTour{
		Xmlns:        ptzNamespace,
		XmlnsTT:      "http://www.onvif.org/ver10/schema",
		ProfileToken: profileToken,
		PresetTour: presetTourXML{
			Token:     presetTour.Token,
			Name:      presetTour.Name,
			AutoStart: presetTour.AutoStart,
		},
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("ModifyPresetTour failed: %w", err)
	}

	return nil
}

// OperatePresetTour controls a preset tour (Start/Stop/Pause/Extended).
func (c *Client) OperatePresetTour(ctx context.Context, profileToken, presetTourToken, operation string) error {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return ErrServiceNotSupported
	}

	type OperatePresetTour struct {
		XMLName         xml.Name `xml:"tptz:OperatePresetTour"`
		Xmlns           string   `xml:"xmlns:tptz,attr"`
		ProfileToken    string   `xml:"tptz:ProfileToken"`
		PresetTourToken string   `xml:"tptz:PresetTourToken"`
		Operation       string   `xml:"tptz:Operation"`
	}

	req := OperatePresetTour{
		Xmlns:           ptzNamespace,
		ProfileToken:    profileToken,
		PresetTourToken: presetTourToken,
		Operation:       operation,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("OperatePresetTour failed: %w", err)
	}

	return nil
}

// RemovePresetTour removes a preset tour.
func (c *Client) RemovePresetTour(ctx context.Context, profileToken, presetTourToken string) error {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return ErrServiceNotSupported
	}

	type RemovePresetTour struct {
		XMLName         xml.Name `xml:"tptz:RemovePresetTour"`
		Xmlns           string   `xml:"xmlns:tptz,attr"`
		ProfileToken    string   `xml:"tptz:ProfileToken"`
		PresetTourToken string   `xml:"tptz:PresetTourToken"`
	}

	req := RemovePresetTour{
		Xmlns:           ptzNamespace,
		ProfileToken:    profileToken,
		PresetTourToken: presetTourToken,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("RemovePresetTour failed: %w", err)
	}

	return nil
}

// GetCompatiblePTZConfigurationsForProfile retrieves PTZ configurations compatible with a given profile.
func (c *Client) GetCompatiblePTZConfigurationsForProfile(ctx context.Context, profileToken string) ([]*PTZConfiguration, error) {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return nil, ErrServiceNotSupported
	}

	type GetCompatibleConfigurations struct {
		XMLName      xml.Name `xml:"tptz:GetCompatibleConfigurations"`
		Xmlns        string   `xml:"xmlns:tptz,attr"`
		ProfileToken string   `xml:"tptz:ProfileToken"`
	}

	type GetCompatibleConfigurationsResponse struct {
		XMLName          xml.Name `xml:"GetCompatibleConfigurationsResponse"`
		PTZConfiguration []struct {
			Token     string `xml:"token,attr"`
			Name      string `xml:"Name"`
			NodeToken string `xml:"NodeToken"`
		} `xml:"PTZConfiguration"`
	}

	req := GetCompatibleConfigurations{
		Xmlns:        ptzNamespace,
		ProfileToken: profileToken,
	}

	var resp GetCompatibleConfigurationsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetCompatiblePTZConfigurationsForProfile failed: %w", err)
	}

	configs := make([]*PTZConfiguration, len(resp.PTZConfiguration))
	for i, cfg := range resp.PTZConfiguration {
		configs[i] = &PTZConfiguration{
			Token:     cfg.Token,
			Name:      cfg.Name,
			NodeToken: cfg.NodeToken,
		}
	}

	return configs, nil
}

// PTZSendAuxiliaryCommand sends an auxiliary command to the PTZ node.
func (c *Client) PTZSendAuxiliaryCommand(ctx context.Context, profileToken, auxiliaryData string) (string, error) {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return "", ErrServiceNotSupported
	}

	type SendAuxiliaryCommand struct {
		XMLName       xml.Name `xml:"tptz:SendAuxiliaryCommand"`
		Xmlns         string   `xml:"xmlns:tptz,attr"`
		ProfileToken  string   `xml:"tptz:ProfileToken"`
		AuxiliaryData string   `xml:"tptz:AuxiliaryData"`
	}

	type SendAuxiliaryCommandResponse struct {
		XMLName           xml.Name `xml:"SendAuxiliaryCommandResponse"`
		AuxiliaryResponse string   `xml:"AuxiliaryResponse"`
	}

	req := SendAuxiliaryCommand{
		Xmlns:         ptzNamespace,
		ProfileToken:  profileToken,
		AuxiliaryData: auxiliaryData,
	}

	var resp SendAuxiliaryCommandResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("PTZSendAuxiliaryCommand failed: %w", err)
	}

	return resp.AuxiliaryResponse, nil
}

// geoLocationXML is a shared type for GeoLocation XML serialization.
type geoLocationXML struct {
	Lon       float64  `xml:"lon,attr"`
	Lat       float64  `xml:"lat,attr"`
	Elevation *float64 `xml:"elevation,attr,omitempty"`
}

// convertToGeoLocationXML converts GeoLocation to XML struct.
func convertToGeoLocationXML(g *GeoLocation) *geoLocationXML {
	if g == nil {
		return nil
	}

	result := &geoLocationXML{
		Lon: g.Lon,
		Lat: g.Lat,
	}
	if g.Elevation != 0 {
		elev := g.Elevation
		result.Elevation = &elev
	}

	return result
}

// GeoMove moves the PTZ unit to a geographic location.
func (c *Client) GeoMove(ctx context.Context, profileToken string, target *GeoLocation, speed *PTZSpeed, areaHeight, areaWidth *float64) error {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return ErrServiceNotSupported
	}

	type GeoMove struct {
		XMLName      xml.Name        `xml:"tptz:GeoMove"`
		Xmlns        string          `xml:"xmlns:tptz,attr"`
		ProfileToken string          `xml:"tptz:ProfileToken"`
		Target       *geoLocationXML `xml:"tptz:Target,omitempty"`
		Speed        *ptzSpeedXML    `xml:"tptz:Speed,omitempty"`
		AreaHeight   *float64        `xml:"tptz:AreaHeight,omitempty"`
		AreaWidth    *float64        `xml:"tptz:AreaWidth,omitempty"`
	}

	req := GeoMove{
		Xmlns:        ptzNamespace,
		ProfileToken: profileToken,
		Target:       convertToGeoLocationXML(target),
		Speed:        convertToPTZSpeedXML(speed),
		AreaHeight:   areaHeight,
		AreaWidth:    areaWidth,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("GeoMove failed: %w", err)
	}

	return nil
}

// MoveAndStartTracking moves to a position and starts object tracking.
func (c *Client) MoveAndStartTracking(ctx context.Context, request *MoveAndStartTrackingRequest) error {
	endpoint := c.ptzEndpoint
	if endpoint == "" {
		return ErrServiceNotSupported
	}

	type MoveAndStartTracking struct {
		XMLName        xml.Name        `xml:"tptz:MoveAndStartTracking"`
		Xmlns          string          `xml:"xmlns:tptz,attr"`
		ProfileToken   string          `xml:"tptz:ProfileToken"`
		PresetToken    *string         `xml:"tptz:PresetToken,omitempty"`
		GeoLocation    *geoLocationXML `xml:"tptz:GeoLocation,omitempty"`
		TargetPosition *ptzVectorXML   `xml:"tptz:TargetPosition,omitempty"`
		Speed          *ptzSpeedXML    `xml:"tptz:Speed,omitempty"`
		ObjectID       *int            `xml:"tptz:ObjectID,omitempty"`
	}

	req := MoveAndStartTracking{
		Xmlns:          ptzNamespace,
		ProfileToken:   request.ProfileToken,
		PresetToken:    request.PresetToken,
		GeoLocation:    convertToGeoLocationXML(request.GeoLocation),
		TargetPosition: convertToPTZVectorXML(request.TargetPosition),
		Speed:          convertToPTZSpeedXML(request.Speed),
		ObjectID:       request.ObjectID,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("MoveAndStartTracking failed: %w", err)
	}

	return nil
}
