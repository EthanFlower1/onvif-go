package onvif

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/EthanFlower1/onvif-go/internal/soap"
)

// Uplink service namespace.
const uplinkNamespace = "http://www.onvif.org/ver10/uplink/wsdl"

// getUplinkEndpoint returns the uplink service endpoint, falling back to the device endpoint.
func (c *Client) getUplinkEndpoint() string {
	if c.uplinkEndpoint != "" {
		return c.uplinkEndpoint
	}

	return c.endpoint
}

// GetUplinkServiceCapabilities returns the capabilities of the uplink service.
func (c *Client) GetUplinkServiceCapabilities(ctx context.Context) (*UplinkServiceCapabilities, error) {
	endpoint := c.getUplinkEndpoint()

	type GetServiceCapabilities struct {
		XMLName xml.Name `xml:"tup:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:tup,attr"`
	}

	type CapabilitiesEntry struct {
		MaxUplinks          *int    `xml:"MaxUplinks,attr"`
		Protocols           string  `xml:"Protocols,attr"`
		AuthorizationModes  string  `xml:"AuthorizationModes,attr"`
		StreamingOverUplink *bool   `xml:"StreamingOverUplink,attr"`
	}

	type GetServiceCapabilitiesResponse struct {
		XMLName      xml.Name          `xml:"GetServiceCapabilitiesResponse"`
		Capabilities CapabilitiesEntry `xml:"Capabilities"`
	}

	req := GetServiceCapabilities{
		Xmlns: uplinkNamespace,
	}

	var resp GetServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetUplinkServiceCapabilities failed: %w", err)
	}

	return &UplinkServiceCapabilities{
		MaxUplinks:          resp.Capabilities.MaxUplinks,
		Protocols:           resp.Capabilities.Protocols,
		AuthorizationModes:  resp.Capabilities.AuthorizationModes,
		StreamingOverUplink: resp.Capabilities.StreamingOverUplink,
	}, nil
}

// GetUplinks retrieves all configured uplink connections.
func (c *Client) GetUplinks(ctx context.Context) ([]*UplinkConfiguration, error) {
	endpoint := c.getUplinkEndpoint()

	type GetUplinks struct {
		XMLName xml.Name `xml:"tup:GetUplinks"`
		Xmlns   string   `xml:"xmlns:tup,attr"`
	}

	type ConfigurationEntry struct {
		RemoteAddress              string  `xml:"RemoteAddress"`
		CertificateID              *string `xml:"CertificateID"`
		UserLevel                  string  `xml:"UserLevel"`
		Status                     *string `xml:"Status"`
		CertPathValidationPolicyID *string `xml:"CertPathValidationPolicyID"`
		AuthorizationServer        *string `xml:"AuthorizationServer"`
		Error                      *string `xml:"Error"`
	}

	type GetUplinksResponse struct {
		XMLName       xml.Name             `xml:"GetUplinksResponse"`
		Configuration []ConfigurationEntry `xml:"Configuration"`
	}

	req := GetUplinks{
		Xmlns: uplinkNamespace,
	}

	var resp GetUplinksResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetUplinks failed: %w", err)
	}

	configs := make([]*UplinkConfiguration, 0, len(resp.Configuration))
	for i := range resp.Configuration {
		e := &resp.Configuration[i]
		configs = append(configs, &UplinkConfiguration{
			RemoteAddress:              e.RemoteAddress,
			CertificateID:              e.CertificateID,
			UserLevel:                  e.UserLevel,
			Status:                     e.Status,
			CertPathValidationPolicyID: e.CertPathValidationPolicyID,
			AuthorizationServer:        e.AuthorizationServer,
			Error:                      e.Error,
		})
	}

	return configs, nil
}

// SetUplink adds or modifies an uplink configuration.
// The device uses RemoteAddress to determine whether to update an existing entry or create a new one.
func (c *Client) SetUplink(ctx context.Context, config UplinkConfiguration) error {
	endpoint := c.getUplinkEndpoint()

	type ConfigurationReq struct {
		RemoteAddress              string  `xml:"tup:RemoteAddress"`
		CertificateID              *string `xml:"tup:CertificateID,omitempty"`
		UserLevel                  string  `xml:"tup:UserLevel"`
		CertPathValidationPolicyID *string `xml:"tup:CertPathValidationPolicyID,omitempty"`
		AuthorizationServer        *string `xml:"tup:AuthorizationServer,omitempty"`
	}

	type SetUplink struct {
		XMLName       xml.Name         `xml:"tup:SetUplink"`
		Xmlns         string           `xml:"xmlns:tup,attr"`
		Configuration ConfigurationReq `xml:"tup:Configuration"`
	}

	req := SetUplink{
		Xmlns: uplinkNamespace,
		Configuration: ConfigurationReq{
			RemoteAddress:              config.RemoteAddress,
			CertificateID:              config.CertificateID,
			UserLevel:                  config.UserLevel,
			CertPathValidationPolicyID: config.CertPathValidationPolicyID,
			AuthorizationServer:        config.AuthorizationServer,
		},
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("SetUplink failed: %w", err)
	}

	return nil
}

// DeleteUplink removes an uplink configuration identified by its RemoteAddress.
func (c *Client) DeleteUplink(ctx context.Context, remoteAddress string) error {
	endpoint := c.getUplinkEndpoint()

	type DeleteUplink struct {
		XMLName       xml.Name `xml:"tup:DeleteUplink"`
		Xmlns         string   `xml:"xmlns:tup,attr"`
		RemoteAddress string   `xml:"tup:RemoteAddress"`
	}

	req := DeleteUplink{
		Xmlns:         uplinkNamespace,
		RemoteAddress: remoteAddress,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("DeleteUplink failed: %w", err)
	}

	return nil
}
