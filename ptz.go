package onvif

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/0x524a/onvif-go/internal/soap"
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
