package onvif

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/EthanFlower1/onvif-go/internal/soap"
)

// Replay service namespace.
const replayNamespace = "http://www.onvif.org/ver10/replay/wsdl"

// getReplayEndpoint returns the replay service endpoint, falling back to the device endpoint.
func (c *Client) getReplayEndpoint() string {
	if c.replayEndpoint != "" {
		return c.replayEndpoint
	}

	return c.endpoint
}

// GetReplayServiceCapabilities retrieves the capabilities of the replay service.
func (c *Client) GetReplayServiceCapabilities(ctx context.Context) (*ReplayServiceCapabilities, error) {
	endpoint := c.getReplayEndpoint()

	type GetServiceCapabilities struct {
		XMLName xml.Name `xml:"trp:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:trp,attr"`
	}

	type SessionTimeoutRangeEntry struct {
		Min string `xml:"Min"`
		Max string `xml:"Max"`
	}

	type GetServiceCapabilitiesResponse struct {
		XMLName      xml.Name `xml:"GetServiceCapabilitiesResponse"`
		Capabilities struct {
			ReversePlayback     bool                     `xml:"ReversePlayback,attr"`
			RTPRTSP_TCP         bool                     `xml:"RTP_RTSP_TCP,attr"`
			SessionTimeoutRange SessionTimeoutRangeEntry `xml:"SessionTimeoutRange"`
		} `xml:"Capabilities"`
	}

	req := GetServiceCapabilities{
		Xmlns: replayNamespace,
	}

	var resp GetServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetReplayServiceCapabilities failed: %w", err)
	}

	caps := &ReplayServiceCapabilities{
		ReversePlayback: resp.Capabilities.ReversePlayback,
		RTPRTSP_TCP:     resp.Capabilities.RTPRTSP_TCP,
	}

	if resp.Capabilities.SessionTimeoutRange.Min != "" || resp.Capabilities.SessionTimeoutRange.Max != "" {
		caps.SessionTimeoutRange = &DurationRange{
			Min: resp.Capabilities.SessionTimeoutRange.Min,
			Max: resp.Capabilities.SessionTimeoutRange.Max,
		}
	}

	return caps, nil
}

// GetReplayUri retrieves a URI for replaying a recording.
func (c *Client) GetReplayUri(ctx context.Context, recordingToken string, stream string, protocol string) (string, error) {
	endpoint := c.getReplayEndpoint()

	type TransportEntry struct {
		Protocol string `xml:"tt:Protocol"`
	}

	type StreamSetupEntry struct {
		Stream    string         `xml:"tt:Stream"`
		Transport TransportEntry `xml:"tt:Transport"`
	}

	type GetReplayUri struct {
		XMLName        xml.Name         `xml:"trp:GetReplayUri"`
		Xmlns          string           `xml:"xmlns:trp,attr"`
		XmlnsTt        string           `xml:"xmlns:tt,attr"`
		StreamSetup    StreamSetupEntry `xml:"trp:StreamSetup"`
		RecordingToken string           `xml:"trp:RecordingToken"`
	}

	type GetReplayUriResponse struct {
		XMLName xml.Name `xml:"GetReplayUriResponse"`
		Uri     string   `xml:"Uri"`
	}

	req := GetReplayUri{
		Xmlns:   replayNamespace,
		XmlnsTt: "http://www.onvif.org/ver10/schema",
		StreamSetup: StreamSetupEntry{
			Stream: stream,
			Transport: TransportEntry{
				Protocol: protocol,
			},
		},
		RecordingToken: recordingToken,
	}

	var resp GetReplayUriResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("GetReplayUri failed: %w", err)
	}

	return resp.Uri, nil
}

// GetReplayConfiguration retrieves the current replay configuration.
func (c *Client) GetReplayConfiguration(ctx context.Context) (*ReplayConfiguration, error) {
	endpoint := c.getReplayEndpoint()

	type GetReplayConfiguration struct {
		XMLName xml.Name `xml:"trp:GetReplayConfiguration"`
		Xmlns   string   `xml:"xmlns:trp,attr"`
	}

	type GetReplayConfigurationResponse struct {
		XMLName       xml.Name `xml:"GetReplayConfigurationResponse"`
		Configuration struct {
			SessionTimeout string `xml:"SessionTimeout"`
		} `xml:"Configuration"`
	}

	req := GetReplayConfiguration{
		Xmlns: replayNamespace,
	}

	var resp GetReplayConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetReplayConfiguration failed: %w", err)
	}

	return &ReplayConfiguration{
		SessionTimeout: resp.Configuration.SessionTimeout,
	}, nil
}

// SetReplayConfiguration updates the replay configuration on the device.
func (c *Client) SetReplayConfiguration(ctx context.Context, config *ReplayConfiguration) error {
	endpoint := c.getReplayEndpoint()

	type ConfigurationEntry struct {
		SessionTimeout string `xml:"trp:SessionTimeout"`
	}

	type SetReplayConfiguration struct {
		XMLName       xml.Name           `xml:"trp:SetReplayConfiguration"`
		Xmlns         string             `xml:"xmlns:trp,attr"`
		Configuration ConfigurationEntry `xml:"trp:Configuration"`
	}

	type SetReplayConfigurationResponse struct {
		XMLName xml.Name `xml:"SetReplayConfigurationResponse"`
	}

	req := SetReplayConfiguration{
		Xmlns: replayNamespace,
		Configuration: ConfigurationEntry{
			SessionTimeout: config.SessionTimeout,
		},
	}

	var resp SetReplayConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetReplayConfiguration failed: %w", err)
	}

	return nil
}
