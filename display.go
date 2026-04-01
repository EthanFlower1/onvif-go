package onvif

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/EthanFlower1/onvif-go/internal/soap"
)

// Display service namespace.
const displayNamespace = "http://www.onvif.org/ver10/display/wsdl"

// getDisplayEndpoint returns the display endpoint, falling back to the device endpoint.
func (c *Client) getDisplayEndpoint() string {
	if c.displayEndpoint != "" {
		return c.displayEndpoint
	}

	return c.endpoint
}

// newDisplaySOAPClient creates a SOAP client for the display service.
func (c *Client) newDisplaySOAPClient() *soap.Client {
	username, password := c.GetCredentials()

	return soap.NewClient(c.httpClient, username, password)
}

// ============================================================
// Capabilities
// ============================================================

// GetDisplayServiceCapabilities returns the capabilities of the Display service.
func (c *Client) GetDisplayServiceCapabilities(ctx context.Context) (*DisplayServiceCapabilities, error) {
	endpoint := c.getDisplayEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tls:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:tls,attr"`
	}

	type capsXML struct {
		FixedLayout *bool `xml:"FixedLayout,attr"`
	}

	type Response struct {
		XMLName      xml.Name `xml:"GetServiceCapabilitiesResponse"`
		Capabilities capsXML  `xml:"Capabilities"`
	}

	req := Request{Xmlns: displayNamespace}

	var resp Response

	if err := c.newDisplaySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetDisplayServiceCapabilities failed: %w", err)
	}

	return &DisplayServiceCapabilities{
		FixedLayout: resp.Capabilities.FixedLayout,
	}, nil
}

// ============================================================
// Layout
// ============================================================

// GetLayout returns the current layout of the specified video output.
func (c *Client) GetLayout(ctx context.Context, videoOutputToken string) (*Layout, error) {
	endpoint := c.getDisplayEndpoint()

	type Request struct {
		XMLName     xml.Name `xml:"tls:GetLayout"`
		Xmlns       string   `xml:"xmlns:tls,attr"`
		VideoOutput string   `xml:"tls:VideoOutput"`
	}

	type paneXML struct {
		Pane string  `xml:"Pane"`
		Area struct {
			Bottom float64 `xml:"bottom,attr"`
			Top    float64 `xml:"top,attr"`
			Right  float64 `xml:"right,attr"`
			Left   float64 `xml:"left,attr"`
		} `xml:"Area"`
	}

	type layoutXML struct {
		Pane []paneXML `xml:"PaneLayout"`
	}

	type Response struct {
		XMLName xml.Name  `xml:"GetLayoutResponse"`
		Layout  layoutXML `xml:"Layout"`
	}

	req := Request{
		Xmlns:       displayNamespace,
		VideoOutput: videoOutputToken,
	}

	var resp Response

	if err := c.newDisplaySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetLayout failed: %w", err)
	}

	layout := &Layout{
		Pane: make([]PaneLayout, len(resp.Layout.Pane)),
	}

	for i, p := range resp.Layout.Pane {
		layout.Pane[i] = PaneLayout{
			Pane: p.Pane,
			Area: FloatRectangle{
				Bottom: p.Area.Bottom,
				Top:    p.Area.Top,
				Right:  p.Area.Right,
				Left:   p.Area.Left,
			},
		}
	}

	return layout, nil
}

// SetLayout changes the layout of the specified video output.
func (c *Client) SetLayout(ctx context.Context, videoOutputToken string, layout Layout) error {
	endpoint := c.getDisplayEndpoint()

	type areaXML struct {
		Bottom float64 `xml:"bottom,attr"`
		Top    float64 `xml:"top,attr"`
		Right  float64 `xml:"right,attr"`
		Left   float64 `xml:"left,attr"`
	}

	type paneXML struct {
		XMLName xml.Name `xml:"tls:PaneLayout"`
		Pane    string   `xml:"tls:Pane"`
		Area    areaXML  `xml:"tls:Area"`
	}

	type layoutXML struct {
		XMLName xml.Name  `xml:"tls:Layout"`
		Pane    []paneXML `xml:",omitempty"`
	}

	type Request struct {
		XMLName     xml.Name  `xml:"tls:SetLayout"`
		Xmlns       string    `xml:"xmlns:tls,attr"`
		VideoOutput string    `xml:"tls:VideoOutput"`
		Layout      layoutXML
	}

	type Response struct {
		XMLName xml.Name `xml:"SetLayoutResponse"`
	}

	panes := make([]paneXML, len(layout.Pane))

	for i, p := range layout.Pane {
		panes[i] = paneXML{
			Pane: p.Pane,
			Area: areaXML{
				Bottom: p.Area.Bottom,
				Top:    p.Area.Top,
				Right:  p.Area.Right,
				Left:   p.Area.Left,
			},
		}
	}

	req := Request{
		Xmlns:       displayNamespace,
		VideoOutput: videoOutputToken,
		Layout:      layoutXML{Pane: panes},
	}

	var resp Response

	if err := c.newDisplaySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetLayout failed: %w", err)
	}

	return nil
}

// ============================================================
// Display Options
// ============================================================

// GetDisplayOptions returns the layout options and coding capabilities for the specified video output.
func (c *Client) GetDisplayOptions(ctx context.Context, videoOutputToken string) (*DisplayOptions, error) {
	endpoint := c.getDisplayEndpoint()

	type Request struct {
		XMLName     xml.Name `xml:"tls:GetDisplayOptions"`
		Xmlns       string   `xml:"xmlns:tls,attr"`
		VideoOutput string   `xml:"tls:VideoOutput"`
	}

	type limitsXML struct {
		Max int `xml:"Max"`
	}

	type codingCapsXML struct {
		InputTokensLimits  *limitsXML `xml:"InputTokensLimits"`
		OutputTokensLimits *limitsXML `xml:"OutputTokensLimits"`
	}

	type Response struct {
		XMLName            xml.Name      `xml:"GetDisplayOptionsResponse"`
		CodingCapabilities codingCapsXML `xml:"CodingCapabilities"`
	}

	req := Request{
		Xmlns:       displayNamespace,
		VideoOutput: videoOutputToken,
	}

	var resp Response

	if err := c.newDisplaySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetDisplayOptions failed: %w", err)
	}

	opts := &DisplayOptions{
		CodingCapabilities: CodingCapabilities{},
	}

	if resp.CodingCapabilities.InputTokensLimits != nil {
		opts.CodingCapabilities.InputTokensLimits = &CodingCapabilityLimits{
			Max: resp.CodingCapabilities.InputTokensLimits.Max,
		}
	}

	if resp.CodingCapabilities.OutputTokensLimits != nil {
		opts.CodingCapabilities.OutputTokensLimits = &CodingCapabilityLimits{
			Max: resp.CodingCapabilities.OutputTokensLimits.Max,
		}
	}

	return opts, nil
}

// ============================================================
// Pane Configurations
// ============================================================

// GetPaneConfigurations returns all pane configurations for the specified video output.
func (c *Client) GetPaneConfigurations(ctx context.Context, videoOutputToken string) ([]*PaneConfiguration, error) {
	endpoint := c.getDisplayEndpoint()

	type Request struct {
		XMLName     xml.Name `xml:"tls:GetPaneConfigurations"`
		Xmlns       string   `xml:"xmlns:tls,attr"`
		VideoOutput string   `xml:"tls:VideoOutput"`
	}

	type paneConfigXML struct {
		Token            string  `xml:"token,attr"`
		PaneName         string  `xml:"PaneName"`
		AudioOutputToken *string `xml:"AudioOutputToken"`
		AudioSourceToken *string `xml:"AudioSourceToken"`
		ReceiverToken    *string `xml:"ReceiverToken"`
		MediaUri         *string `xml:"MediaUri"`
		Profile          *string `xml:"Profile"`
	}

	type Response struct {
		XMLName           xml.Name         `xml:"GetPaneConfigurationsResponse"`
		PaneConfiguration []*paneConfigXML `xml:"PaneConfiguration"`
	}

	req := Request{
		Xmlns:       displayNamespace,
		VideoOutput: videoOutputToken,
	}

	var resp Response

	if err := c.newDisplaySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetPaneConfigurations failed: %w", err)
	}

	result := make([]*PaneConfiguration, 0, len(resp.PaneConfiguration))

	for _, pc := range resp.PaneConfiguration {
		result = append(result, &PaneConfiguration{
			Token:            pc.Token,
			PaneName:         pc.PaneName,
			AudioOutputToken: pc.AudioOutputToken,
			AudioSourceToken: pc.AudioSourceToken,
			ReceiverToken:    pc.ReceiverToken,
			MediaUri:         pc.MediaUri,
			Profile:          pc.Profile,
		})
	}

	return result, nil
}

// GetPaneConfiguration returns the configuration for a specific pane.
func (c *Client) GetPaneConfiguration(ctx context.Context, videoOutputToken, paneToken string) (*PaneConfiguration, error) {
	endpoint := c.getDisplayEndpoint()

	type Request struct {
		XMLName     xml.Name `xml:"tls:GetPaneConfiguration"`
		Xmlns       string   `xml:"xmlns:tls,attr"`
		VideoOutput string   `xml:"tls:VideoOutput"`
		Pane        string   `xml:"tls:Pane"`
	}

	type paneConfigXML struct {
		Token            string  `xml:"token,attr"`
		PaneName         string  `xml:"PaneName"`
		AudioOutputToken *string `xml:"AudioOutputToken"`
		AudioSourceToken *string `xml:"AudioSourceToken"`
		ReceiverToken    *string `xml:"ReceiverToken"`
		MediaUri         *string `xml:"MediaUri"`
		Profile          *string `xml:"Profile"`
	}

	type Response struct {
		XMLName           xml.Name      `xml:"GetPaneConfigurationResponse"`
		PaneConfiguration paneConfigXML `xml:"PaneConfiguration"`
	}

	req := Request{
		Xmlns:       displayNamespace,
		VideoOutput: videoOutputToken,
		Pane:        paneToken,
	}

	var resp Response

	if err := c.newDisplaySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetPaneConfiguration failed: %w", err)
	}

	return &PaneConfiguration{
		Token:            resp.PaneConfiguration.Token,
		PaneName:         resp.PaneConfiguration.PaneName,
		AudioOutputToken: resp.PaneConfiguration.AudioOutputToken,
		AudioSourceToken: resp.PaneConfiguration.AudioSourceToken,
		ReceiverToken:    resp.PaneConfiguration.ReceiverToken,
		MediaUri:         resp.PaneConfiguration.MediaUri,
		Profile:          resp.PaneConfiguration.Profile,
	}, nil
}

// SetPaneConfigurations modifies one or more pane configurations for the specified video output.
func (c *Client) SetPaneConfigurations(ctx context.Context, videoOutputToken string, paneConfigs []PaneConfiguration) error {
	endpoint := c.getDisplayEndpoint()

	type paneConfigXML struct {
		XMLName          xml.Name `xml:"tls:PaneConfiguration"`
		Token            string   `xml:"token,attr"`
		PaneName         string   `xml:"tls:PaneName"`
		AudioOutputToken *string  `xml:"tls:AudioOutputToken,omitempty"`
		AudioSourceToken *string  `xml:"tls:AudioSourceToken,omitempty"`
		ReceiverToken    *string  `xml:"tls:ReceiverToken,omitempty"`
		MediaUri         *string  `xml:"tls:MediaUri,omitempty"`
		Profile          *string  `xml:"tls:Profile,omitempty"`
	}

	type Request struct {
		XMLName           xml.Name         `xml:"tls:SetPaneConfigurations"`
		Xmlns             string           `xml:"xmlns:tls,attr"`
		VideoOutput       string           `xml:"tls:VideoOutput"`
		PaneConfiguration []paneConfigXML  `xml:",omitempty"`
	}

	type Response struct {
		XMLName xml.Name `xml:"SetPaneConfigurationsResponse"`
	}

	xmlConfigs := make([]paneConfigXML, len(paneConfigs))

	for i, pc := range paneConfigs {
		xmlConfigs[i] = paneConfigXML{
			Token:            pc.Token,
			PaneName:         pc.PaneName,
			AudioOutputToken: pc.AudioOutputToken,
			AudioSourceToken: pc.AudioSourceToken,
			ReceiverToken:    pc.ReceiverToken,
			MediaUri:         pc.MediaUri,
			Profile:          pc.Profile,
		}
	}

	req := Request{
		Xmlns:             displayNamespace,
		VideoOutput:       videoOutputToken,
		PaneConfiguration: xmlConfigs,
	}

	var resp Response

	if err := c.newDisplaySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetPaneConfigurations failed: %w", err)
	}

	return nil
}

// SetPaneConfiguration changes the configuration of the specified pane.
func (c *Client) SetPaneConfiguration(ctx context.Context, videoOutputToken string, paneConfig PaneConfiguration) error {
	endpoint := c.getDisplayEndpoint()

	type paneConfigXML struct {
		XMLName          xml.Name `xml:"tls:PaneConfiguration"`
		Token            string   `xml:"token,attr"`
		PaneName         string   `xml:"tls:PaneName"`
		AudioOutputToken *string  `xml:"tls:AudioOutputToken,omitempty"`
		AudioSourceToken *string  `xml:"tls:AudioSourceToken,omitempty"`
		ReceiverToken    *string  `xml:"tls:ReceiverToken,omitempty"`
		MediaUri         *string  `xml:"tls:MediaUri,omitempty"`
		Profile          *string  `xml:"tls:Profile,omitempty"`
	}

	type Request struct {
		XMLName           xml.Name      `xml:"tls:SetPaneConfiguration"`
		Xmlns             string        `xml:"xmlns:tls,attr"`
		VideoOutput       string        `xml:"tls:VideoOutput"`
		PaneConfiguration paneConfigXML
	}

	type Response struct {
		XMLName xml.Name `xml:"SetPaneConfigurationResponse"`
	}

	req := Request{
		Xmlns:       displayNamespace,
		VideoOutput: videoOutputToken,
		PaneConfiguration: paneConfigXML{
			Token:            paneConfig.Token,
			PaneName:         paneConfig.PaneName,
			AudioOutputToken: paneConfig.AudioOutputToken,
			AudioSourceToken: paneConfig.AudioSourceToken,
			ReceiverToken:    paneConfig.ReceiverToken,
			MediaUri:         paneConfig.MediaUri,
			Profile:          paneConfig.Profile,
		},
	}

	var resp Response

	if err := c.newDisplaySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetPaneConfiguration failed: %w", err)
	}

	return nil
}

// CreatePaneConfiguration creates a new pane configuration for the specified video output.
func (c *Client) CreatePaneConfiguration(ctx context.Context, videoOutputToken string, paneConfig PaneConfiguration) (string, error) {
	endpoint := c.getDisplayEndpoint()

	type paneConfigXML struct {
		XMLName          xml.Name `xml:"tls:PaneConfiguration"`
		Token            string   `xml:"token,attr,omitempty"`
		PaneName         string   `xml:"tls:PaneName"`
		AudioOutputToken *string  `xml:"tls:AudioOutputToken,omitempty"`
		AudioSourceToken *string  `xml:"tls:AudioSourceToken,omitempty"`
		ReceiverToken    *string  `xml:"tls:ReceiverToken,omitempty"`
		MediaUri         *string  `xml:"tls:MediaUri,omitempty"`
		Profile          *string  `xml:"tls:Profile,omitempty"`
	}

	type Request struct {
		XMLName           xml.Name      `xml:"tls:CreatePaneConfiguration"`
		Xmlns             string        `xml:"xmlns:tls,attr"`
		VideoOutput       string        `xml:"tls:VideoOutput"`
		PaneConfiguration paneConfigXML
	}

	type Response struct {
		XMLName   xml.Name `xml:"CreatePaneConfigurationResponse"`
		PaneToken string   `xml:"PaneToken"`
	}

	req := Request{
		Xmlns:       displayNamespace,
		VideoOutput: videoOutputToken,
		PaneConfiguration: paneConfigXML{
			Token:            paneConfig.Token,
			PaneName:         paneConfig.PaneName,
			AudioOutputToken: paneConfig.AudioOutputToken,
			AudioSourceToken: paneConfig.AudioSourceToken,
			ReceiverToken:    paneConfig.ReceiverToken,
			MediaUri:         paneConfig.MediaUri,
			Profile:          paneConfig.Profile,
		},
	}

	var resp Response

	if err := c.newDisplaySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreatePaneConfiguration failed: %w", err)
	}

	return resp.PaneToken, nil
}

// DeletePaneConfiguration deletes a pane configuration from the specified video output.
func (c *Client) DeletePaneConfiguration(ctx context.Context, videoOutputToken, paneToken string) error {
	endpoint := c.getDisplayEndpoint()

	type Request struct {
		XMLName     xml.Name `xml:"tls:DeletePaneConfiguration"`
		Xmlns       string   `xml:"xmlns:tls,attr"`
		VideoOutput string   `xml:"tls:VideoOutput"`
		PaneToken   string   `xml:"tls:PaneToken"`
	}

	type Response struct {
		XMLName xml.Name `xml:"DeletePaneConfigurationResponse"`
	}

	req := Request{
		Xmlns:       displayNamespace,
		VideoOutput: videoOutputToken,
		PaneToken:   paneToken,
	}

	var resp Response

	if err := c.newDisplaySOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeletePaneConfiguration failed: %w", err)
	}

	return nil
}
