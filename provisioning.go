package onvif

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/EthanFlower1/onvif-go/internal/soap"
)

// Provisioning service namespace.
const provisioningNamespace = "http://www.onvif.org/ver10/provisioning/wsdl"

// getProvisioningEndpoint returns the provisioning service endpoint, falling back to the device endpoint.
func (c *Client) getProvisioningEndpoint() string {
	if c.provisioningEndpoint != "" {
		return c.provisioningEndpoint
	}

	return c.endpoint
}

// GetProvisioningServiceCapabilities returns the capabilities of the provisioning service.
func (c *Client) GetProvisioningServiceCapabilities(ctx context.Context) (*ProvisioningServiceCapabilities, error) {
	endpoint := c.getProvisioningEndpoint()

	type GetServiceCapabilities struct {
		XMLName xml.Name `xml:"tpv:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:tpv,attr"`
	}

	type SourceEntry struct {
		VideoSourceToken  string `xml:"VideoSourceToken,attr"`
		MaximumPanMoves   *int   `xml:"MaximumPanMoves,attr"`
		MaximumTiltMoves  *int   `xml:"MaximumTiltMoves,attr"`
		MaximumZoomMoves  *int   `xml:"MaximumZoomMoves,attr"`
		MaximumRollMoves  *int   `xml:"MaximumRollMoves,attr"`
		AutoLevel         *bool  `xml:"AutoLevel,attr"`
		MaximumFocusMoves *int   `xml:"MaximumFocusMoves,attr"`
		AutoFocus         *bool  `xml:"AutoFocus,attr"`
	}

	type CapabilitiesEntry struct {
		DefaultTimeout string        `xml:"DefaultTimeout"`
		Source         []SourceEntry `xml:"Source"`
	}

	type GetServiceCapabilitiesResponse struct {
		XMLName      xml.Name          `xml:"GetServiceCapabilitiesResponse"`
		Capabilities CapabilitiesEntry `xml:"Capabilities"`
	}

	req := GetServiceCapabilities{
		Xmlns: provisioningNamespace,
	}

	var resp GetServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetProvisioningServiceCapabilities failed: %w", err)
	}

	caps := &ProvisioningServiceCapabilities{
		DefaultTimeout: resp.Capabilities.DefaultTimeout,
	}

	for i := range resp.Capabilities.Source {
		s := &resp.Capabilities.Source[i]
		src := &ProvisioningSourceCapabilities{
			VideoSourceToken:  s.VideoSourceToken,
			MaximumPanMoves:   s.MaximumPanMoves,
			MaximumTiltMoves:  s.MaximumTiltMoves,
			MaximumZoomMoves:  s.MaximumZoomMoves,
			MaximumRollMoves:  s.MaximumRollMoves,
			AutoLevel:         s.AutoLevel,
			MaximumFocusMoves: s.MaximumFocusMoves,
			AutoFocus:         s.AutoFocus,
		}
		caps.Source = append(caps.Source, src)
	}

	return caps, nil
}

// PanMove moves the device on the pan axis.
func (c *Client) PanMove(ctx context.Context, videoSource string, direction PanDirection, timeout *string) error {
	endpoint := c.getProvisioningEndpoint()

	type PanMove struct {
		XMLName     xml.Name  `xml:"tpv:PanMove"`
		Xmlns       string    `xml:"xmlns:tpv,attr"`
		VideoSource string    `xml:"tpv:VideoSource"`
		Direction   string    `xml:"tpv:Direction"`
		Timeout     *string   `xml:"tpv:Timeout,omitempty"`
	}

	req := PanMove{
		Xmlns:       provisioningNamespace,
		VideoSource: videoSource,
		Direction:   string(direction),
		Timeout:     timeout,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("PanMove failed: %w", err)
	}

	return nil
}

// TiltMove moves the device on the tilt axis.
func (c *Client) TiltMove(ctx context.Context, videoSource string, direction TiltDirection, timeout *string) error {
	endpoint := c.getProvisioningEndpoint()

	type TiltMove struct {
		XMLName     xml.Name `xml:"tpv:TiltMove"`
		Xmlns       string   `xml:"xmlns:tpv,attr"`
		VideoSource string   `xml:"tpv:VideoSource"`
		Direction   string   `xml:"tpv:Direction"`
		Timeout     *string  `xml:"tpv:Timeout,omitempty"`
	}

	req := TiltMove{
		Xmlns:       provisioningNamespace,
		VideoSource: videoSource,
		Direction:   string(direction),
		Timeout:     timeout,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("TiltMove failed: %w", err)
	}

	return nil
}

// ZoomMove moves the device on the zoom axis.
func (c *Client) ZoomMove(ctx context.Context, videoSource string, direction ZoomDirection, timeout *string) error {
	endpoint := c.getProvisioningEndpoint()

	type ZoomMove struct {
		XMLName     xml.Name `xml:"tpv:ZoomMove"`
		Xmlns       string   `xml:"xmlns:tpv,attr"`
		VideoSource string   `xml:"tpv:VideoSource"`
		Direction   string   `xml:"tpv:Direction"`
		Timeout     *string  `xml:"tpv:Timeout,omitempty"`
	}

	req := ZoomMove{
		Xmlns:       provisioningNamespace,
		VideoSource: videoSource,
		Direction:   string(direction),
		Timeout:     timeout,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("ZoomMove failed: %w", err)
	}

	return nil
}

// RollMove moves the device on the roll axis.
func (c *Client) RollMove(ctx context.Context, videoSource string, direction RollDirection, timeout *string) error {
	endpoint := c.getProvisioningEndpoint()

	type RollMove struct {
		XMLName     xml.Name `xml:"tpv:RollMove"`
		Xmlns       string   `xml:"xmlns:tpv,attr"`
		VideoSource string   `xml:"tpv:VideoSource"`
		Direction   string   `xml:"tpv:Direction"`
		Timeout     *string  `xml:"tpv:Timeout,omitempty"`
	}

	req := RollMove{
		Xmlns:       provisioningNamespace,
		VideoSource: videoSource,
		Direction:   string(direction),
		Timeout:     timeout,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("RollMove failed: %w", err)
	}

	return nil
}

// ProvisioningFocusMove moves the device on the focus axis.
func (c *Client) ProvisioningFocusMove(ctx context.Context, videoSource string, direction FocusDirection, timeout *string) error {
	endpoint := c.getProvisioningEndpoint()

	type FocusMoveReq struct {
		XMLName     xml.Name `xml:"tpv:FocusMove"`
		Xmlns       string   `xml:"xmlns:tpv,attr"`
		VideoSource string   `xml:"tpv:VideoSource"`
		Direction   string   `xml:"tpv:Direction"`
		Timeout     *string  `xml:"tpv:Timeout,omitempty"`
	}

	req := FocusMoveReq{
		Xmlns:       provisioningNamespace,
		VideoSource: videoSource,
		Direction:   string(direction),
		Timeout:     timeout,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("ProvisioningFocusMove failed: %w", err)
	}

	return nil
}

// ProvisioningStop stops device motion on all axes.
func (c *Client) ProvisioningStop(ctx context.Context, videoSource string) error {
	endpoint := c.getProvisioningEndpoint()

	type StopReq struct {
		XMLName     xml.Name `xml:"tpv:Stop"`
		Xmlns       string   `xml:"xmlns:tpv,attr"`
		VideoSource string   `xml:"tpv:VideoSource"`
	}

	req := StopReq{
		Xmlns:       provisioningNamespace,
		VideoSource: videoSource,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("ProvisioningStop failed: %w", err)
	}

	return nil
}

// GetProvisioningUsage returns the lifetime move counts for a video source.
func (c *Client) GetProvisioningUsage(ctx context.Context, videoSource string) (*ProvisioningUsage, error) {
	endpoint := c.getProvisioningEndpoint()

	type GetUsage struct {
		XMLName     xml.Name `xml:"tpv:GetUsage"`
		Xmlns       string   `xml:"xmlns:tpv,attr"`
		VideoSource string   `xml:"tpv:VideoSource"`
	}

	type UsageEntry struct {
		Pan   *int `xml:"Pan"`
		Tilt  *int `xml:"Tilt"`
		Zoom  *int `xml:"Zoom"`
		Roll  *int `xml:"Roll"`
		Focus *int `xml:"Focus"`
	}

	type GetUsageResponse struct {
		XMLName xml.Name   `xml:"GetUsageResponse"`
		Usage   UsageEntry `xml:"Usage"`
	}

	req := GetUsage{
		Xmlns:       provisioningNamespace,
		VideoSource: videoSource,
	}

	var resp GetUsageResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetProvisioningUsage failed: %w", err)
	}

	return &ProvisioningUsage{
		Pan:   resp.Usage.Pan,
		Tilt:  resp.Usage.Tilt,
		Zoom:  resp.Usage.Zoom,
		Roll:  resp.Usage.Roll,
		Focus: resp.Usage.Focus,
	}, nil
}
