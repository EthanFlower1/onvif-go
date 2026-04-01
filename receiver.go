package onvif

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/0x524a/onvif-go/internal/soap"
)

// Receiver service namespace.
const receiverNamespace = "http://www.onvif.org/ver10/receiver/wsdl"

// getReceiverEndpoint returns the receiver service endpoint, falling back to the device endpoint.
func (c *Client) getReceiverEndpoint() string {
	if c.receiverEndpoint != "" {
		return c.receiverEndpoint
	}

	return c.endpoint
}

// GetReceiverServiceCapabilities retrieves the capabilities of the receiver service.
func (c *Client) GetReceiverServiceCapabilities(ctx context.Context) (*ReceiverServiceCapabilities, error) {
	endpoint := c.getReceiverEndpoint()

	type GetServiceCapabilities struct {
		XMLName xml.Name `xml:"trv:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:trv,attr"`
	}

	type GetServiceCapabilitiesResponse struct {
		XMLName      xml.Name `xml:"GetServiceCapabilitiesResponse"`
		Capabilities struct {
			RTPMulticast         bool `xml:"RTP_Multicast,attr"`
			RTPTCP               bool `xml:"RTP_TCP,attr"`
			RTPRTSP_TCP          bool `xml:"RTP_RTSP_TCP,attr"`
			SupportedReceivers   int  `xml:"SupportedReceivers,attr"`
			MaximumRTSPURILength int  `xml:"MaximumRTSPURILength,attr"`
		} `xml:"Capabilities"`
	}

	req := GetServiceCapabilities{
		Xmlns: receiverNamespace,
	}

	var resp GetServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetReceiverServiceCapabilities failed: %w", err)
	}

	return &ReceiverServiceCapabilities{
		RTPMulticast:         resp.Capabilities.RTPMulticast,
		RTPTCP:               resp.Capabilities.RTPTCP,
		RTPRTSP_TCP:          resp.Capabilities.RTPRTSP_TCP,
		SupportedReceivers:   resp.Capabilities.SupportedReceivers,
		MaximumRTSPURILength: resp.Capabilities.MaximumRTSPURILength,
	}, nil
}

// GetReceivers retrieves all receivers configured on the device.
func (c *Client) GetReceivers(ctx context.Context) ([]*Receiver, error) {
	endpoint := c.getReceiverEndpoint()

	type GetReceivers struct {
		XMLName xml.Name `xml:"trv:GetReceivers"`
		Xmlns   string   `xml:"xmlns:trv,attr"`
	}

	type TransportEntry struct {
		Protocol string `xml:"Protocol"`
	}

	type StreamSetupEntry struct {
		Stream    string         `xml:"Stream"`
		Transport TransportEntry `xml:"Transport"`
	}

	type ConfigEntry struct {
		Mode        string           `xml:"Mode"`
		MediaUri    string           `xml:"MediaUri"`
		StreamSetup StreamSetupEntry `xml:"StreamSetup"`
	}

	type ReceiverEntry struct {
		Token         string      `xml:"token,attr"`
		Configuration ConfigEntry `xml:"Configuration"`
	}

	type GetReceiversResponse struct {
		XMLName   xml.Name        `xml:"GetReceiversResponse"`
		Receivers []ReceiverEntry `xml:"Receivers"`
	}

	req := GetReceivers{
		Xmlns: receiverNamespace,
	}

	var resp GetReceiversResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetReceivers failed: %w", err)
	}

	receivers := make([]*Receiver, 0, len(resp.Receivers))
	for i := range resp.Receivers {
		entry := &resp.Receivers[i]
		rec := &Receiver{
			Token: entry.Token,
			Configuration: ReceiverConfiguration{
				Mode:     entry.Configuration.Mode,
				MediaURI: entry.Configuration.MediaUri,
				StreamSetup: &StreamSetup{
					Stream: entry.Configuration.StreamSetup.Stream,
					Transport: &Transport{
						Protocol: entry.Configuration.StreamSetup.Transport.Protocol,
					},
				},
			},
		}
		receivers = append(receivers, rec)
	}

	return receivers, nil
}

// GetReceiver retrieves a single receiver by its token.
func (c *Client) GetReceiver(ctx context.Context, receiverToken string) (*Receiver, error) {
	endpoint := c.getReceiverEndpoint()

	type GetReceiver struct {
		XMLName       xml.Name `xml:"trv:GetReceiver"`
		Xmlns         string   `xml:"xmlns:trv,attr"`
		ReceiverToken string   `xml:"trv:ReceiverToken"`
	}

	type TransportEntry struct {
		Protocol string `xml:"Protocol"`
	}

	type StreamSetupEntry struct {
		Stream    string         `xml:"Stream"`
		Transport TransportEntry `xml:"Transport"`
	}

	type ConfigEntry struct {
		Mode        string           `xml:"Mode"`
		MediaUri    string           `xml:"MediaUri"`
		StreamSetup StreamSetupEntry `xml:"StreamSetup"`
	}

	type ReceiverEntry struct {
		Token         string      `xml:"token,attr"`
		Configuration ConfigEntry `xml:"Configuration"`
	}

	type GetReceiverResponse struct {
		XMLName  xml.Name      `xml:"GetReceiverResponse"`
		Receiver ReceiverEntry `xml:"Receiver"`
	}

	req := GetReceiver{
		Xmlns:         receiverNamespace,
		ReceiverToken: receiverToken,
	}

	var resp GetReceiverResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetReceiver failed: %w", err)
	}

	return &Receiver{
		Token: resp.Receiver.Token,
		Configuration: ReceiverConfiguration{
			Mode:     resp.Receiver.Configuration.Mode,
			MediaURI: resp.Receiver.Configuration.MediaUri,
			StreamSetup: &StreamSetup{
				Stream: resp.Receiver.Configuration.StreamSetup.Stream,
				Transport: &Transport{
					Protocol: resp.Receiver.Configuration.StreamSetup.Transport.Protocol,
				},
			},
		},
	}, nil
}

// CreateReceiver creates a new receiver on the device.
func (c *Client) CreateReceiver(ctx context.Context, config *ReceiverConfiguration) (*Receiver, error) {
	endpoint := c.getReceiverEndpoint()

	type TransportReq struct {
		Protocol string `xml:"tt:Protocol"`
	}

	type StreamSetupReq struct {
		Stream    string       `xml:"tt:Stream"`
		Transport TransportReq `xml:"tt:Transport"`
	}

	type ConfigReq struct {
		Mode        string          `xml:"tt:Mode"`
		MediaUri    string          `xml:"tt:MediaUri"`
		StreamSetup *StreamSetupReq `xml:"tt:StreamSetup,omitempty"`
	}

	type CreateReceiver struct {
		XMLName       xml.Name  `xml:"trv:CreateReceiver"`
		Xmlns         string    `xml:"xmlns:trv,attr"`
		XmlnsTt       string    `xml:"xmlns:tt,attr"`
		Configuration ConfigReq `xml:"trv:Configuration"`
	}

	type TransportEntry struct {
		Protocol string `xml:"Protocol"`
	}

	type StreamSetupEntry struct {
		Stream    string         `xml:"Stream"`
		Transport TransportEntry `xml:"Transport"`
	}

	type ConfigEntry struct {
		Mode        string           `xml:"Mode"`
		MediaUri    string           `xml:"MediaUri"`
		StreamSetup StreamSetupEntry `xml:"StreamSetup"`
	}

	type ReceiverEntry struct {
		Token         string      `xml:"token,attr"`
		Configuration ConfigEntry `xml:"Configuration"`
	}

	type CreateReceiverResponse struct {
		XMLName  xml.Name      `xml:"CreateReceiverResponse"`
		Receiver ReceiverEntry `xml:"Receiver"`
	}

	configReq := ConfigReq{
		Mode:     config.Mode,
		MediaUri: config.MediaURI,
	}
	if config.StreamSetup != nil {
		ss := &StreamSetupReq{
			Stream: config.StreamSetup.Stream,
		}
		if config.StreamSetup.Transport != nil {
			ss.Transport = TransportReq{Protocol: config.StreamSetup.Transport.Protocol}
		}
		configReq.StreamSetup = ss
	}

	req := CreateReceiver{
		Xmlns:         receiverNamespace,
		XmlnsTt:       "http://www.onvif.org/ver10/schema",
		Configuration: configReq,
	}

	var resp CreateReceiverResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("CreateReceiver failed: %w", err)
	}

	return &Receiver{
		Token: resp.Receiver.Token,
		Configuration: ReceiverConfiguration{
			Mode:     resp.Receiver.Configuration.Mode,
			MediaURI: resp.Receiver.Configuration.MediaUri,
			StreamSetup: &StreamSetup{
				Stream: resp.Receiver.Configuration.StreamSetup.Stream,
				Transport: &Transport{
					Protocol: resp.Receiver.Configuration.StreamSetup.Transport.Protocol,
				},
			},
		},
	}, nil
}

// DeleteReceiver deletes a receiver from the device.
func (c *Client) DeleteReceiver(ctx context.Context, receiverToken string) error {
	endpoint := c.getReceiverEndpoint()

	type DeleteReceiver struct {
		XMLName       xml.Name `xml:"trv:DeleteReceiver"`
		Xmlns         string   `xml:"xmlns:trv,attr"`
		ReceiverToken string   `xml:"trv:ReceiverToken"`
	}

	type DeleteReceiverResponse struct {
		XMLName xml.Name `xml:"DeleteReceiverResponse"`
	}

	req := DeleteReceiver{
		Xmlns:         receiverNamespace,
		ReceiverToken: receiverToken,
	}

	var resp DeleteReceiverResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteReceiver failed: %w", err)
	}

	return nil
}

// ConfigureReceiver updates the configuration of an existing receiver.
func (c *Client) ConfigureReceiver(ctx context.Context, receiverToken string, config *ReceiverConfiguration) error {
	endpoint := c.getReceiverEndpoint()

	type TransportReq struct {
		Protocol string `xml:"tt:Protocol"`
	}

	type StreamSetupReq struct {
		Stream    string       `xml:"tt:Stream"`
		Transport TransportReq `xml:"tt:Transport"`
	}

	type ConfigReq struct {
		Mode        string          `xml:"tt:Mode"`
		MediaUri    string          `xml:"tt:MediaUri"`
		StreamSetup *StreamSetupReq `xml:"tt:StreamSetup,omitempty"`
	}

	type ConfigureReceiver struct {
		XMLName       xml.Name  `xml:"trv:ConfigureReceiver"`
		Xmlns         string    `xml:"xmlns:trv,attr"`
		XmlnsTt       string    `xml:"xmlns:tt,attr"`
		ReceiverToken string    `xml:"trv:ReceiverToken"`
		Configuration ConfigReq `xml:"trv:Configuration"`
	}

	type ConfigureReceiverResponse struct {
		XMLName xml.Name `xml:"ConfigureReceiverResponse"`
	}

	configReq := ConfigReq{
		Mode:     config.Mode,
		MediaUri: config.MediaURI,
	}
	if config.StreamSetup != nil {
		ss := &StreamSetupReq{
			Stream: config.StreamSetup.Stream,
		}
		if config.StreamSetup.Transport != nil {
			ss.Transport = TransportReq{Protocol: config.StreamSetup.Transport.Protocol}
		}
		configReq.StreamSetup = ss
	}

	req := ConfigureReceiver{
		Xmlns:         receiverNamespace,
		XmlnsTt:       "http://www.onvif.org/ver10/schema",
		ReceiverToken: receiverToken,
		Configuration: configReq,
	}

	var resp ConfigureReceiverResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("ConfigureReceiver failed: %w", err)
	}

	return nil
}

// SetReceiverMode sets the operating mode of a receiver.
func (c *Client) SetReceiverMode(ctx context.Context, receiverToken, mode string) error {
	endpoint := c.getReceiverEndpoint()

	type SetReceiverMode struct {
		XMLName       xml.Name `xml:"trv:SetReceiverMode"`
		Xmlns         string   `xml:"xmlns:trv,attr"`
		ReceiverToken string   `xml:"trv:ReceiverToken"`
		Mode          string   `xml:"trv:Mode"`
	}

	type SetReceiverModeResponse struct {
		XMLName xml.Name `xml:"SetReceiverModeResponse"`
	}

	req := SetReceiverMode{
		Xmlns:         receiverNamespace,
		ReceiverToken: receiverToken,
		Mode:          mode,
	}

	var resp SetReceiverModeResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetReceiverMode failed: %w", err)
	}

	return nil
}

// GetReceiverState retrieves the current state of a receiver.
func (c *Client) GetReceiverState(ctx context.Context, receiverToken string) (*ReceiverStateInformation, error) {
	endpoint := c.getReceiverEndpoint()

	type GetReceiverState struct {
		XMLName       xml.Name `xml:"trv:GetReceiverState"`
		Xmlns         string   `xml:"xmlns:trv,attr"`
		ReceiverToken string   `xml:"trv:ReceiverToken"`
	}

	type ReceiverStateEntry struct {
		State       string `xml:"State"`
		AutoCreated bool   `xml:"AutoCreated"`
	}

	type GetReceiverStateResponse struct {
		XMLName       xml.Name           `xml:"GetReceiverStateResponse"`
		ReceiverState ReceiverStateEntry `xml:"ReceiverState"`
	}

	req := GetReceiverState{
		Xmlns:         receiverNamespace,
		ReceiverToken: receiverToken,
	}

	var resp GetReceiverStateResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetReceiverState failed: %w", err)
	}

	return &ReceiverStateInformation{
		State:       resp.ReceiverState.State,
		AutoCreated: resp.ReceiverState.AutoCreated,
	}, nil
}
