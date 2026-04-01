package onvif

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/0x524a/onvif-go/internal/soap"
)

// Recording service namespace.
const recordingNamespace = "http://www.onvif.org/ver10/recording/wsdl"

// getRecordingEndpoint returns the recording service endpoint, falling back to the device endpoint.
func (c *Client) getRecordingEndpoint() string {
	if c.recordingEndpoint != "" {
		return c.recordingEndpoint
	}

	return c.endpoint
}

// GetRecordingServiceCapabilities retrieves the capabilities of the recording service.
func (c *Client) GetRecordingServiceCapabilities(ctx context.Context) (*RecordingServiceCapabilities, error) {
	endpoint := c.getRecordingEndpoint()

	type GetServiceCapabilities struct {
		XMLName xml.Name `xml:"trc:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:trc,attr"`
	}

	type GetServiceCapabilitiesResponse struct {
		XMLName      xml.Name `xml:"GetServiceCapabilitiesResponse"`
		Capabilities struct {
			DynamicRecordings bool `xml:"DynamicRecordings,attr"`
			DynamicTracks     bool `xml:"DynamicTracks,attr"`
			MaxRecordings     int  `xml:"MaxRecordings,attr"`
			MaxRecordingJobs  int  `xml:"MaxRecordingJobs,attr"`
		} `xml:"Capabilities"`
	}

	req := GetServiceCapabilities{
		Xmlns: recordingNamespace,
	}

	var resp GetServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRecordingServiceCapabilities failed: %w", err)
	}

	return &RecordingServiceCapabilities{
		DynamicRecordings: resp.Capabilities.DynamicRecordings,
		DynamicTracks:     resp.Capabilities.DynamicTracks,
		MaxRecordings:     resp.Capabilities.MaxRecordings,
		MaxRecordingJobs:  resp.Capabilities.MaxRecordingJobs,
	}, nil
}

// GetRecordings retrieves all recordings available on the device.
func (c *Client) GetRecordings(ctx context.Context) ([]*Recording, error) {
	endpoint := c.getRecordingEndpoint()

	type GetRecordings struct {
		XMLName xml.Name `xml:"trc:GetRecordings"`
		Xmlns   string   `xml:"xmlns:trc,attr"`
	}

	type TrackConfigEntry struct {
		TrackType   string `xml:"TrackType"`
		Description string `xml:"Description"`
	}

	type TrackEntry struct {
		TrackToken    string           `xml:"TrackToken"`
		Configuration TrackConfigEntry `xml:"Configuration"`
		DataFrom      string           `xml:"DataFrom"`
		DataTo        string           `xml:"DataTo"`
	}

	type TracksContainer struct {
		Track []TrackEntry `xml:"Track"`
	}

	type SourceEntry struct {
		SourceId    string `xml:"SourceId"`
		Name        string `xml:"Name"`
		Location    string `xml:"Location"`
		Description string `xml:"Description"`
		Address     string `xml:"Address"`
	}

	type ConfigEntry struct {
		Source               SourceEntry `xml:"Source"`
		Content              string      `xml:"Content"`
		MaximumRetentionTime string      `xml:"MaximumRetentionTime"`
	}

	type RecordingItem struct {
		RecordingToken string          `xml:"RecordingToken"`
		Configuration  ConfigEntry     `xml:"Configuration"`
		Tracks         TracksContainer `xml:"Tracks"`
	}

	type GetRecordingsResponse struct {
		XMLName       xml.Name        `xml:"GetRecordingsResponse"`
		RecordingItem []RecordingItem `xml:"RecordingItem"`
	}

	req := GetRecordings{
		Xmlns: recordingNamespace,
	}

	var resp GetRecordingsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRecordings failed: %w", err)
	}

	recordings := make([]*Recording, 0, len(resp.RecordingItem))
	for i := range resp.RecordingItem {
		item := &resp.RecordingItem[i]
		rec := &Recording{
			Token: item.RecordingToken,
			Configuration: RecordingConfiguration{
				Source: RecordingSourceInformation{
					SourceId:    item.Configuration.Source.SourceId,
					Name:        item.Configuration.Source.Name,
					Location:    item.Configuration.Source.Location,
					Description: item.Configuration.Source.Description,
					Address:     item.Configuration.Source.Address,
				},
				Content:              item.Configuration.Content,
				MaximumRetentionTime: item.Configuration.MaximumRetentionTime,
			},
		}

		tracks := make([]*RecordingTrack, 0, len(item.Tracks.Track))
		for j := range item.Tracks.Track {
			t := &item.Tracks.Track[j]
			tracks = append(tracks, &RecordingTrack{
				Token: t.TrackToken,
				Configuration: TrackConfiguration{
					TrackType:   t.Configuration.TrackType,
					Description: t.Configuration.Description,
				},
				DataFrom: t.DataFrom,
				DataTo:   t.DataTo,
			})
		}
		rec.Tracks = tracks

		recordings = append(recordings, rec)
	}

	return recordings, nil
}

// CreateRecording creates a new recording on the device.
func (c *Client) CreateRecording(ctx context.Context, config *RecordingConfiguration) (string, error) {
	endpoint := c.getRecordingEndpoint()

	type SourceReq struct {
		SourceId    string `xml:"tt:SourceId"`
		Name        string `xml:"tt:Name"`
		Location    string `xml:"tt:Location,omitempty"`
		Description string `xml:"tt:Description,omitempty"`
		Address     string `xml:"tt:Address,omitempty"`
	}

	type ConfigReq struct {
		Source               SourceReq `xml:"tt:Source"`
		Content              string    `xml:"tt:Content,omitempty"`
		MaximumRetentionTime string    `xml:"tt:MaximumRetentionTime,omitempty"`
	}

	type CreateRecording struct {
		XMLName               xml.Name  `xml:"trc:CreateRecording"`
		Xmlns                 string    `xml:"xmlns:trc,attr"`
		XmlnsTt               string    `xml:"xmlns:tt,attr"`
		RecordingConfiguration ConfigReq `xml:"trc:RecordingConfiguration"`
	}

	type CreateRecordingResponse struct {
		XMLName        xml.Name `xml:"CreateRecordingResponse"`
		RecordingToken string   `xml:"RecordingToken"`
	}

	req := CreateRecording{
		Xmlns:   recordingNamespace,
		XmlnsTt: "http://www.onvif.org/ver10/schema",
		RecordingConfiguration: ConfigReq{
			Source: SourceReq{
				SourceId:    config.Source.SourceId,
				Name:        config.Source.Name,
				Location:    config.Source.Location,
				Description: config.Source.Description,
				Address:     config.Source.Address,
			},
			Content:              config.Content,
			MaximumRetentionTime: config.MaximumRetentionTime,
		},
	}

	var resp CreateRecordingResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateRecording failed: %w", err)
	}

	return resp.RecordingToken, nil
}

// DeleteRecording deletes a recording from the device.
func (c *Client) DeleteRecording(ctx context.Context, recordingToken string) error {
	endpoint := c.getRecordingEndpoint()

	type DeleteRecording struct {
		XMLName        xml.Name `xml:"trc:DeleteRecording"`
		Xmlns          string   `xml:"xmlns:trc,attr"`
		RecordingToken string   `xml:"trc:RecordingToken"`
	}

	req := DeleteRecording{
		Xmlns:          recordingNamespace,
		RecordingToken: recordingToken,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("DeleteRecording failed: %w", err)
	}

	return nil
}

// SetRecordingConfiguration updates the configuration of an existing recording.
func (c *Client) SetRecordingConfiguration(
	ctx context.Context, recordingToken string, config *RecordingConfiguration,
) error {
	endpoint := c.getRecordingEndpoint()

	type SourceReq struct {
		SourceId    string `xml:"tt:SourceId"`
		Name        string `xml:"tt:Name"`
		Location    string `xml:"tt:Location,omitempty"`
		Description string `xml:"tt:Description,omitempty"`
		Address     string `xml:"tt:Address,omitempty"`
	}

	type ConfigReq struct {
		Source               SourceReq `xml:"tt:Source"`
		Content              string    `xml:"tt:Content,omitempty"`
		MaximumRetentionTime string    `xml:"tt:MaximumRetentionTime,omitempty"`
	}

	type SetRecordingConfiguration struct {
		XMLName               xml.Name  `xml:"trc:SetRecordingConfiguration"`
		Xmlns                 string    `xml:"xmlns:trc,attr"`
		XmlnsTt               string    `xml:"xmlns:tt,attr"`
		RecordingToken        string    `xml:"trc:RecordingToken"`
		RecordingConfiguration ConfigReq `xml:"trc:RecordingConfiguration"`
	}

	req := SetRecordingConfiguration{
		Xmlns:          recordingNamespace,
		XmlnsTt:        "http://www.onvif.org/ver10/schema",
		RecordingToken: recordingToken,
		RecordingConfiguration: ConfigReq{
			Source: SourceReq{
				SourceId:    config.Source.SourceId,
				Name:        config.Source.Name,
				Location:    config.Source.Location,
				Description: config.Source.Description,
				Address:     config.Source.Address,
			},
			Content:              config.Content,
			MaximumRetentionTime: config.MaximumRetentionTime,
		},
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("SetRecordingConfiguration failed: %w", err)
	}

	return nil
}

// GetRecordingConfiguration retrieves the configuration of a recording.
func (c *Client) GetRecordingConfiguration(
	ctx context.Context, recordingToken string,
) (*RecordingConfiguration, error) {
	endpoint := c.getRecordingEndpoint()

	type GetRecordingConfiguration struct {
		XMLName        xml.Name `xml:"trc:GetRecordingConfiguration"`
		Xmlns          string   `xml:"xmlns:trc,attr"`
		RecordingToken string   `xml:"trc:RecordingToken"`
	}

	type SourceResp struct {
		SourceId    string `xml:"SourceId"`
		Name        string `xml:"Name"`
		Location    string `xml:"Location"`
		Description string `xml:"Description"`
		Address     string `xml:"Address"`
	}

	type GetRecordingConfigurationResponse struct {
		XMLName               xml.Name `xml:"GetRecordingConfigurationResponse"`
		RecordingConfiguration struct {
			Source               SourceResp `xml:"Source"`
			Content              string     `xml:"Content"`
			MaximumRetentionTime string     `xml:"MaximumRetentionTime"`
		} `xml:"RecordingConfiguration"`
	}

	req := GetRecordingConfiguration{
		Xmlns:          recordingNamespace,
		RecordingToken: recordingToken,
	}

	var resp GetRecordingConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRecordingConfiguration failed: %w", err)
	}

	return &RecordingConfiguration{
		Source: RecordingSourceInformation{
			SourceId:    resp.RecordingConfiguration.Source.SourceId,
			Name:        resp.RecordingConfiguration.Source.Name,
			Location:    resp.RecordingConfiguration.Source.Location,
			Description: resp.RecordingConfiguration.Source.Description,
			Address:     resp.RecordingConfiguration.Source.Address,
		},
		Content:              resp.RecordingConfiguration.Content,
		MaximumRetentionTime: resp.RecordingConfiguration.MaximumRetentionTime,
	}, nil
}

// GetRecordingOptions retrieves the options available for a recording.
func (c *Client) GetRecordingOptions(ctx context.Context, recordingToken string) (*RecordingOptions, error) {
	endpoint := c.getRecordingEndpoint()

	type GetRecordingOptions struct {
		XMLName        xml.Name `xml:"trc:GetRecordingOptions"`
		Xmlns          string   `xml:"xmlns:trc,attr"`
		RecordingToken string   `xml:"trc:RecordingToken"`
	}

	type GetRecordingOptionsResponse struct {
		XMLName xml.Name `xml:"GetRecordingOptionsResponse"`
		Options struct {
			Track *struct {
				SpareTotal    *int `xml:"SpareTotal"`
				SpareVideo    *int `xml:"SpareVideo"`
				SpareAudio    *int `xml:"SpareAudio"`
				SpareMetadata *int `xml:"SpareMetadata"`
			} `xml:"Track"`
		} `xml:"Options"`
	}

	req := GetRecordingOptions{
		Xmlns:          recordingNamespace,
		RecordingToken: recordingToken,
	}

	var resp GetRecordingOptionsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRecordingOptions failed: %w", err)
	}

	options := &RecordingOptions{}

	if resp.Options.Track != nil {
		options.Track = &RecordingTrackOptions{
			SpareTotal:    resp.Options.Track.SpareTotal,
			SpareVideo:    resp.Options.Track.SpareVideo,
			SpareAudio:    resp.Options.Track.SpareAudio,
			SpareMetadata: resp.Options.Track.SpareMetadata,
		}
	}

	return options, nil
}
