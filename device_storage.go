package onvif

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/EthanFlower1/onvif-go/internal/soap"
)

// GetStorageConfigurations retrieves storage configurations. ONVIF Specification: GetStorageConfigurations operation.
func (c *Client) GetStorageConfigurations(ctx context.Context) ([]*StorageConfiguration, error) {
	type GetStorageConfigurationsBody struct {
		XMLName xml.Name `xml:"tds:GetStorageConfigurations"`
		Xmlns   string   `xml:"xmlns:tds,attr"`
	}

	type storageConfigResp struct {
		Token string `xml:"Token"`
		Data  struct {
			Type       string `xml:"Type"`
			LocalPath  string `xml:"LocalPath"`
			StorageUri string `xml:"StorageUri"`
		} `xml:"Data"`
	}

	type GetStorageConfigurationsResponse struct {
		StorageConfigurations []storageConfigResp `xml:"StorageConfigurations"`
	}

	request := GetStorageConfigurationsBody{
		Xmlns: deviceNamespace,
	}
	var response GetStorageConfigurationsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.endpoint, "", request, &response); err != nil {
		return nil, fmt.Errorf("GetStorageConfigurations failed: %w", err)
	}

	configs := make([]*StorageConfiguration, 0, len(response.StorageConfigurations))
	for _, sc := range response.StorageConfigurations {
		configs = append(configs, &StorageConfiguration{
			Token: sc.Token,
			Data: StorageConfigurationData{
				Type:       sc.Data.Type,
				LocalPath:  sc.Data.LocalPath,
				StorageURI: sc.Data.StorageUri,
			},
		})
	}

	return configs, nil
}

// GetStorageConfiguration retrieves a storage configuration. ONVIF Specification: GetStorageConfiguration operation.
func (c *Client) GetStorageConfiguration(ctx context.Context, token string) (*StorageConfiguration, error) {
	type GetStorageConfigurationBody struct {
		XMLName xml.Name `xml:"tds:GetStorageConfiguration"`
		Xmlns   string   `xml:"xmlns:tds,attr"`
		Token   string   `xml:"tds:Token"`
	}

	type storageConfigSingleResp struct {
		Token string `xml:"Token"`
		Data  struct {
			Type       string `xml:"Type"`
			LocalPath  string `xml:"LocalPath"`
			StorageUri string `xml:"StorageUri"`
		} `xml:"Data"`
	}

	type GetStorageConfigurationResponse struct {
		StorageConfiguration storageConfigSingleResp `xml:"StorageConfiguration"`
	}

	request := GetStorageConfigurationBody{
		Xmlns: deviceNamespace,
		Token: token,
	}
	var response GetStorageConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.endpoint, "", request, &response); err != nil {
		return nil, fmt.Errorf("GetStorageConfiguration failed: %w", err)
	}

	sc := response.StorageConfiguration
	return &StorageConfiguration{
		Token: sc.Token,
		Data: StorageConfigurationData{
			Type:       sc.Data.Type,
			LocalPath:  sc.Data.LocalPath,
			StorageURI: sc.Data.StorageUri,
		},
	}, nil
}

// CreateStorageConfiguration creates a storage configuration.
// ONVIF Specification: CreateStorageConfiguration operation.
func (c *Client) CreateStorageConfiguration(ctx context.Context, config *StorageConfiguration) (string, error) {
	type CreateStorageConfigurationBody struct {
		XMLName              xml.Name              `xml:"tds:CreateStorageConfiguration"`
		Xmlns                string                `xml:"xmlns:tds,attr"`
		StorageConfiguration *StorageConfiguration `xml:"tds:StorageConfiguration"`
	}

	type CreateStorageConfigurationResponse struct {
		XMLName xml.Name `xml:"CreateStorageConfigurationResponse"`
		Token   string   `xml:"Token"`
	}

	request := CreateStorageConfigurationBody{
		Xmlns:                deviceNamespace,
		StorageConfiguration: config,
	}
	var response CreateStorageConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.endpoint, "", request, &response); err != nil {
		return "", fmt.Errorf("CreateStorageConfiguration failed: %w", err)
	}

	return response.Token, nil
}

// SetStorageConfiguration sets a storage configuration. ONVIF Specification: SetStorageConfiguration operation.
func (c *Client) SetStorageConfiguration(ctx context.Context, config *StorageConfiguration) error {
	type SetStorageConfigurationBody struct {
		XMLName              xml.Name              `xml:"tds:SetStorageConfiguration"`
		Xmlns                string                `xml:"xmlns:tds,attr"`
		StorageConfiguration *StorageConfiguration `xml:"tds:StorageConfiguration"`
	}

	type SetStorageConfigurationResponse struct {
		XMLName xml.Name `xml:"SetStorageConfigurationResponse"`
	}

	request := SetStorageConfigurationBody{
		Xmlns:                deviceNamespace,
		StorageConfiguration: config,
	}
	var response SetStorageConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.endpoint, "", request, &response); err != nil {
		return fmt.Errorf("SetStorageConfiguration failed: %w", err)
	}

	return nil
}

// DeleteStorageConfiguration deletes a storage configuration.
// ONVIF Specification: DeleteStorageConfiguration operation.
func (c *Client) DeleteStorageConfiguration(ctx context.Context, token string) error {
	type DeleteStorageConfigurationBody struct {
		XMLName xml.Name `xml:"tds:DeleteStorageConfiguration"`
		Xmlns   string   `xml:"xmlns:tds,attr"`
		Token   string   `xml:"tds:Token"`
	}

	type DeleteStorageConfigurationResponse struct {
		XMLName xml.Name `xml:"DeleteStorageConfigurationResponse"`
	}

	request := DeleteStorageConfigurationBody{
		Xmlns: deviceNamespace,
		Token: token,
	}
	var response DeleteStorageConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.endpoint, "", request, &response); err != nil {
		return fmt.Errorf("DeleteStorageConfiguration failed: %w", err)
	}

	return nil
}

// SetHashingAlgorithm sets the hashing algorithm. ONVIF Specification: SetHashingAlgorithm operation.
func (c *Client) SetHashingAlgorithm(ctx context.Context, algorithm string) error {
	type SetHashingAlgorithmBody struct {
		XMLName   xml.Name `xml:"tds:SetHashingAlgorithm"`
		Xmlns     string   `xml:"xmlns:tds,attr"`
		Algorithm string   `xml:"tds:Algorithm"`
	}

	type SetHashingAlgorithmResponse struct {
		XMLName xml.Name `xml:"SetHashingAlgorithmResponse"`
	}

	request := SetHashingAlgorithmBody{
		Xmlns:     deviceNamespace,
		Algorithm: algorithm,
	}
	var response SetHashingAlgorithmResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.endpoint, "", request, &response); err != nil {
		return fmt.Errorf("SetHashingAlgorithm failed: %w", err)
	}

	return nil
}
