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

// CreateTrack creates a new track within a recording.
func (c *Client) CreateTrack(ctx context.Context, recordingToken string, config *TrackConfiguration) (string, error) {
	endpoint := c.getRecordingEndpoint()

	type TrackConfigReq struct {
		TrackType   string `xml:"tt:TrackType"`
		Description string `xml:"tt:Description,omitempty"`
	}

	type CreateTrack struct {
		XMLName            xml.Name       `xml:"trc:CreateTrack"`
		Xmlns              string         `xml:"xmlns:trc,attr"`
		XmlnsTt            string         `xml:"xmlns:tt,attr"`
		RecordingToken     string         `xml:"trc:RecordingToken"`
		TrackConfiguration TrackConfigReq `xml:"trc:TrackConfiguration"`
	}

	type CreateTrackResponse struct {
		XMLName    xml.Name `xml:"CreateTrackResponse"`
		TrackToken string   `xml:"TrackToken"`
	}

	req := CreateTrack{
		Xmlns:          recordingNamespace,
		XmlnsTt:        "http://www.onvif.org/ver10/schema",
		RecordingToken: recordingToken,
		TrackConfiguration: TrackConfigReq{
			TrackType:   config.TrackType,
			Description: config.Description,
		},
	}

	var resp CreateTrackResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateTrack failed: %w", err)
	}

	return resp.TrackToken, nil
}

// DeleteTrack deletes a track from a recording.
func (c *Client) DeleteTrack(ctx context.Context, recordingToken, trackToken string) error {
	endpoint := c.getRecordingEndpoint()

	type DeleteTrack struct {
		XMLName        xml.Name `xml:"trc:DeleteTrack"`
		Xmlns          string   `xml:"xmlns:trc,attr"`
		RecordingToken string   `xml:"trc:RecordingToken"`
		TrackToken     string   `xml:"trc:TrackToken"`
	}

	req := DeleteTrack{
		Xmlns:          recordingNamespace,
		RecordingToken: recordingToken,
		TrackToken:     trackToken,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("DeleteTrack failed: %w", err)
	}

	return nil
}

// GetTrackConfiguration retrieves the configuration of a track within a recording.
func (c *Client) GetTrackConfiguration(ctx context.Context, recordingToken, trackToken string) (*TrackConfiguration, error) {
	endpoint := c.getRecordingEndpoint()

	type GetTrackConfiguration struct {
		XMLName        xml.Name `xml:"trc:GetTrackConfiguration"`
		Xmlns          string   `xml:"xmlns:trc,attr"`
		RecordingToken string   `xml:"trc:RecordingToken"`
		TrackToken     string   `xml:"trc:TrackToken"`
	}

	type GetTrackConfigurationResponse struct {
		XMLName            xml.Name `xml:"GetTrackConfigurationResponse"`
		TrackConfiguration struct {
			TrackType   string `xml:"TrackType"`
			Description string `xml:"Description"`
		} `xml:"TrackConfiguration"`
	}

	req := GetTrackConfiguration{
		Xmlns:          recordingNamespace,
		RecordingToken: recordingToken,
		TrackToken:     trackToken,
	}

	var resp GetTrackConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetTrackConfiguration failed: %w", err)
	}

	return &TrackConfiguration{
		TrackType:   resp.TrackConfiguration.TrackType,
		Description: resp.TrackConfiguration.Description,
	}, nil
}

// SetTrackConfiguration updates the configuration of a track within a recording.
func (c *Client) SetTrackConfiguration(ctx context.Context, recordingToken, trackToken string, config *TrackConfiguration) error {
	endpoint := c.getRecordingEndpoint()

	type TrackConfigReq struct {
		TrackType   string `xml:"tt:TrackType"`
		Description string `xml:"tt:Description,omitempty"`
	}

	type SetTrackConfiguration struct {
		XMLName            xml.Name       `xml:"trc:SetTrackConfiguration"`
		Xmlns              string         `xml:"xmlns:trc,attr"`
		XmlnsTt            string         `xml:"xmlns:tt,attr"`
		RecordingToken     string         `xml:"trc:RecordingToken"`
		TrackToken         string         `xml:"trc:TrackToken"`
		TrackConfiguration TrackConfigReq `xml:"trc:TrackConfiguration"`
	}

	req := SetTrackConfiguration{
		Xmlns:          recordingNamespace,
		XmlnsTt:        "http://www.onvif.org/ver10/schema",
		RecordingToken: recordingToken,
		TrackToken:     trackToken,
		TrackConfiguration: TrackConfigReq{
			TrackType:   config.TrackType,
			Description: config.Description,
		},
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("SetTrackConfiguration failed: %w", err)
	}

	return nil
}

// CreateRecordingJob creates a new recording job on the device.
func (c *Client) CreateRecordingJob(ctx context.Context, config *RecordingJobConfiguration) (string, *RecordingJobConfiguration, error) {
	endpoint := c.getRecordingEndpoint()

	type JobConfigReq struct {
		RecordingToken string `xml:"tt:RecordingToken"`
		Mode           string `xml:"tt:Mode"`
		Priority       int    `xml:"tt:Priority"`
	}

	type CreateRecordingJob struct {
		XMLName          xml.Name     `xml:"trc:CreateRecordingJob"`
		Xmlns            string       `xml:"xmlns:trc,attr"`
		XmlnsTt          string       `xml:"xmlns:tt,attr"`
		JobConfiguration JobConfigReq `xml:"trc:JobConfiguration"`
	}

	type JobConfigResp struct {
		RecordingToken string `xml:"RecordingToken"`
		Mode           string `xml:"Mode"`
		Priority       int    `xml:"Priority"`
	}

	type CreateRecordingJobResponse struct {
		XMLName          xml.Name      `xml:"CreateRecordingJobResponse"`
		JobToken         string        `xml:"JobToken"`
		JobConfiguration JobConfigResp `xml:"JobConfiguration"`
	}

	req := CreateRecordingJob{
		Xmlns:   recordingNamespace,
		XmlnsTt: "http://www.onvif.org/ver10/schema",
		JobConfiguration: JobConfigReq{
			RecordingToken: config.RecordingToken,
			Mode:           config.Mode,
			Priority:       config.Priority,
		},
	}

	var resp CreateRecordingJobResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", nil, fmt.Errorf("CreateRecordingJob failed: %w", err)
	}

	actualConfig := &RecordingJobConfiguration{
		RecordingToken: resp.JobConfiguration.RecordingToken,
		Mode:           resp.JobConfiguration.Mode,
		Priority:       resp.JobConfiguration.Priority,
	}

	return resp.JobToken, actualConfig, nil
}

// DeleteRecordingJob deletes a recording job from the device.
func (c *Client) DeleteRecordingJob(ctx context.Context, jobToken string) error {
	endpoint := c.getRecordingEndpoint()

	type DeleteRecordingJob struct {
		XMLName  xml.Name `xml:"trc:DeleteRecordingJob"`
		Xmlns    string   `xml:"xmlns:trc,attr"`
		JobToken string   `xml:"trc:JobToken"`
	}

	req := DeleteRecordingJob{
		Xmlns:    recordingNamespace,
		JobToken: jobToken,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("DeleteRecordingJob failed: %w", err)
	}

	return nil
}

// GetRecordingJobs retrieves all recording jobs on the device.
func (c *Client) GetRecordingJobs(ctx context.Context) ([]*RecordingJob, error) {
	endpoint := c.getRecordingEndpoint()

	type GetRecordingJobs struct {
		XMLName xml.Name `xml:"trc:GetRecordingJobs"`
		Xmlns   string   `xml:"xmlns:trc,attr"`
	}

	type JobConfigItem struct {
		RecordingToken string `xml:"RecordingToken"`
		Mode           string `xml:"Mode"`
		Priority       int    `xml:"Priority"`
	}

	type JobItem struct {
		JobToken         string        `xml:"JobToken"`
		JobConfiguration JobConfigItem `xml:"JobConfiguration"`
	}

	type GetRecordingJobsResponse struct {
		XMLName xml.Name  `xml:"GetRecordingJobsResponse"`
		JobItem []JobItem `xml:"JobItem"`
	}

	req := GetRecordingJobs{
		Xmlns: recordingNamespace,
	}

	var resp GetRecordingJobsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRecordingJobs failed: %w", err)
	}

	jobs := make([]*RecordingJob, 0, len(resp.JobItem))
	for i := range resp.JobItem {
		item := &resp.JobItem[i]
		jobs = append(jobs, &RecordingJob{
			Token: item.JobToken,
			Configuration: RecordingJobConfiguration{
				RecordingToken: item.JobConfiguration.RecordingToken,
				Mode:           item.JobConfiguration.Mode,
				Priority:       item.JobConfiguration.Priority,
			},
		})
	}

	return jobs, nil
}

// SetRecordingJobConfiguration updates the configuration of an existing recording job.
func (c *Client) SetRecordingJobConfiguration(
	ctx context.Context, jobToken string, config *RecordingJobConfiguration,
) (*RecordingJobConfiguration, error) {
	endpoint := c.getRecordingEndpoint()

	type JobConfigReq struct {
		RecordingToken string `xml:"tt:RecordingToken"`
		Mode           string `xml:"tt:Mode"`
		Priority       int    `xml:"tt:Priority"`
	}

	type SetRecordingJobConfiguration struct {
		XMLName          xml.Name     `xml:"trc:SetRecordingJobConfiguration"`
		Xmlns            string       `xml:"xmlns:trc,attr"`
		XmlnsTt          string       `xml:"xmlns:tt,attr"`
		JobToken         string       `xml:"trc:JobToken"`
		JobConfiguration JobConfigReq `xml:"trc:JobConfiguration"`
	}

	type JobConfigResp struct {
		RecordingToken string `xml:"RecordingToken"`
		Mode           string `xml:"Mode"`
		Priority       int    `xml:"Priority"`
	}

	type SetRecordingJobConfigurationResponse struct {
		XMLName          xml.Name      `xml:"SetRecordingJobConfigurationResponse"`
		JobConfiguration JobConfigResp `xml:"JobConfiguration"`
	}

	req := SetRecordingJobConfiguration{
		Xmlns:    recordingNamespace,
		XmlnsTt:  "http://www.onvif.org/ver10/schema",
		JobToken: jobToken,
		JobConfiguration: JobConfigReq{
			RecordingToken: config.RecordingToken,
			Mode:           config.Mode,
			Priority:       config.Priority,
		},
	}

	var resp SetRecordingJobConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("SetRecordingJobConfiguration failed: %w", err)
	}

	return &RecordingJobConfiguration{
		RecordingToken: resp.JobConfiguration.RecordingToken,
		Mode:           resp.JobConfiguration.Mode,
		Priority:       resp.JobConfiguration.Priority,
	}, nil
}

// GetRecordingJobConfiguration retrieves the configuration of a recording job.
func (c *Client) GetRecordingJobConfiguration(
	ctx context.Context, jobToken string,
) (*RecordingJobConfiguration, error) {
	endpoint := c.getRecordingEndpoint()

	type GetRecordingJobConfiguration struct {
		XMLName  xml.Name `xml:"trc:GetRecordingJobConfiguration"`
		Xmlns    string   `xml:"xmlns:trc,attr"`
		JobToken string   `xml:"trc:JobToken"`
	}

	type GetRecordingJobConfigurationResponse struct {
		XMLName          xml.Name `xml:"GetRecordingJobConfigurationResponse"`
		JobConfiguration struct {
			RecordingToken string `xml:"RecordingToken"`
			Mode           string `xml:"Mode"`
			Priority       int    `xml:"Priority"`
		} `xml:"JobConfiguration"`
	}

	req := GetRecordingJobConfiguration{
		Xmlns:    recordingNamespace,
		JobToken: jobToken,
	}

	var resp GetRecordingJobConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRecordingJobConfiguration failed: %w", err)
	}

	return &RecordingJobConfiguration{
		RecordingToken: resp.JobConfiguration.RecordingToken,
		Mode:           resp.JobConfiguration.Mode,
		Priority:       resp.JobConfiguration.Priority,
	}, nil
}

// SetRecordingJobMode sets the mode of a recording job.
func (c *Client) SetRecordingJobMode(ctx context.Context, jobToken, mode string) error {
	endpoint := c.getRecordingEndpoint()

	type SetRecordingJobMode struct {
		XMLName  xml.Name `xml:"trc:SetRecordingJobMode"`
		Xmlns    string   `xml:"xmlns:trc,attr"`
		JobToken string   `xml:"trc:JobToken"`
		Mode     string   `xml:"trc:Mode"`
	}

	req := SetRecordingJobMode{
		Xmlns:    recordingNamespace,
		JobToken: jobToken,
		Mode:     mode,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("SetRecordingJobMode failed: %w", err)
	}

	return nil
}

// GetRecordingJobState retrieves the current state of a recording job.
func (c *Client) GetRecordingJobState(ctx context.Context, jobToken string) (*RecordingJobState, error) {
	endpoint := c.getRecordingEndpoint()

	type GetRecordingJobState struct {
		XMLName  xml.Name `xml:"trc:GetRecordingJobState"`
		Xmlns    string   `xml:"xmlns:trc,attr"`
		JobToken string   `xml:"trc:JobToken"`
	}

	type SourceTokenResp struct {
		Token string `xml:"Token"`
		Type  string `xml:"Type"`
	}

	type SourceStateResp struct {
		SourceToken SourceTokenResp `xml:"SourceToken"`
		State       string          `xml:"State"`
	}

	type GetRecordingJobStateResponse struct {
		XMLName xml.Name `xml:"GetRecordingJobStateResponse"`
		State   struct {
			RecordingToken string            `xml:"RecordingToken"`
			State          string            `xml:"State"`
			Sources        []SourceStateResp `xml:"Sources"`
		} `xml:"State"`
	}

	req := GetRecordingJobState{
		Xmlns:    recordingNamespace,
		JobToken: jobToken,
	}

	var resp GetRecordingJobStateResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRecordingJobState failed: %w", err)
	}

	sources := make([]*RecordingJobSourceState, 0, len(resp.State.Sources))
	for i := range resp.State.Sources {
		s := &resp.State.Sources[i]
		src := &RecordingJobSourceState{
			State: s.State,
		}
		if s.SourceToken.Token != "" || s.SourceToken.Type != "" {
			src.SourceToken = &SourceReference{
				Token: s.SourceToken.Token,
				Type:  s.SourceToken.Type,
			}
		}
		sources = append(sources, src)
	}

	return &RecordingJobState{
		RecordingToken: resp.State.RecordingToken,
		State:          resp.State.State,
		Sources:        sources,
	}, nil
}

// ExportRecordedData initiates an export of recorded data for the given time range.
// It returns the operation token and the list of file names that will be exported.
func (c *Client) ExportRecordedData(
	ctx context.Context,
	startPoint, endPoint string,
	recordingToken string,
	fileFormat string,
	storageDestination string,
) (string, []string, error) {
	endpoint := c.getRecordingEndpoint()

	type SearchScopeReq struct {
		IncludedRecordings string `xml:"tt:IncludedRecordings"`
	}

	type StorageDestinationReq struct {
		StorageUri string `xml:"tt:StorageUri"`
	}

	type ExportRecordedData struct {
		XMLName            xml.Name              `xml:"trc:ExportRecordedData"`
		Xmlns              string                `xml:"xmlns:trc,attr"`
		XmlnsTt            string                `xml:"xmlns:tt,attr"`
		StartPoint         string                `xml:"trc:StartPoint"`
		EndPoint           string                `xml:"trc:EndPoint"`
		SearchScope        SearchScopeReq        `xml:"trc:SearchScope"`
		FileFormat         string                `xml:"trc:FileFormat"`
		StorageDestination StorageDestinationReq `xml:"trc:StorageDestination"`
	}

	type FileNamesResp struct {
		FileName []string `xml:"FileName"`
	}

	type ExportRecordedDataResponse struct {
		XMLName        xml.Name      `xml:"ExportRecordedDataResponse"`
		OperationToken string        `xml:"OperationToken"`
		FileNames      FileNamesResp `xml:"FileNames"`
	}

	req := ExportRecordedData{
		Xmlns:      recordingNamespace,
		XmlnsTt:    "http://www.onvif.org/ver10/schema",
		StartPoint: startPoint,
		EndPoint:   endPoint,
		SearchScope: SearchScopeReq{
			IncludedRecordings: recordingToken,
		},
		FileFormat: fileFormat,
		StorageDestination: StorageDestinationReq{
			StorageUri: storageDestination,
		},
	}

	var resp ExportRecordedDataResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", nil, fmt.Errorf("ExportRecordedData failed: %w", err)
	}

	return resp.OperationToken, resp.FileNames.FileName, nil
}

// StopExportRecordedData stops an ongoing export operation and returns its final state.
func (c *Client) StopExportRecordedData(ctx context.Context, operationToken string) (*ExportRecordedDataState, error) {
	endpoint := c.getRecordingEndpoint()

	type StopExportRecordedData struct {
		XMLName        xml.Name `xml:"trc:StopExportRecordedData"`
		Xmlns          string   `xml:"xmlns:trc,attr"`
		OperationToken string   `xml:"trc:OperationToken"`
	}

	type FileProgressResp struct {
		FileName string  `xml:"FileName"`
		Progress float64 `xml:"Progress"`
	}

	type FileProgressStatusResp struct {
		FileProgress []FileProgressResp `xml:"FileProgress"`
	}

	type StopExportRecordedDataResponse struct {
		XMLName            xml.Name               `xml:"StopExportRecordedDataResponse"`
		Progress           float64                `xml:"Progress"`
		FileProgressStatus FileProgressStatusResp `xml:"FileProgressStatus"`
	}

	req := StopExportRecordedData{
		Xmlns:          recordingNamespace,
		OperationToken: operationToken,
	}

	var resp StopExportRecordedDataResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("StopExportRecordedData failed: %w", err)
	}

	fileProgresses := make([]*FileProgress, 0, len(resp.FileProgressStatus.FileProgress))
	for i := range resp.FileProgressStatus.FileProgress {
		fp := &resp.FileProgressStatus.FileProgress[i]
		fileProgresses = append(fileProgresses, &FileProgress{
			FileName: fp.FileName,
			Progress: fp.Progress,
		})
	}

	return &ExportRecordedDataState{
		Progress:           resp.Progress,
		FileProgressStatus: fileProgresses,
	}, nil
}

// GetExportRecordedDataState retrieves the current state of an export operation.
func (c *Client) GetExportRecordedDataState(ctx context.Context, operationToken string) (*ExportRecordedDataState, error) {
	endpoint := c.getRecordingEndpoint()

	type GetExportRecordedDataState struct {
		XMLName        xml.Name `xml:"trc:GetExportRecordedDataState"`
		Xmlns          string   `xml:"xmlns:trc,attr"`
		OperationToken string   `xml:"trc:OperationToken"`
	}

	type FileProgressResp struct {
		FileName string  `xml:"FileName"`
		Progress float64 `xml:"Progress"`
	}

	type FileProgressStatusResp struct {
		FileProgress []FileProgressResp `xml:"FileProgress"`
	}

	type GetExportRecordedDataStateResponse struct {
		XMLName            xml.Name               `xml:"GetExportRecordedDataStateResponse"`
		Progress           float64                `xml:"Progress"`
		FileProgressStatus FileProgressStatusResp `xml:"FileProgressStatus"`
	}

	req := GetExportRecordedDataState{
		Xmlns:          recordingNamespace,
		OperationToken: operationToken,
	}

	var resp GetExportRecordedDataStateResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetExportRecordedDataState failed: %w", err)
	}

	fileProgresses := make([]*FileProgress, 0, len(resp.FileProgressStatus.FileProgress))
	for i := range resp.FileProgressStatus.FileProgress {
		fp := &resp.FileProgressStatus.FileProgress[i]
		fileProgresses = append(fileProgresses, &FileProgress{
			FileName: fp.FileName,
			Progress: fp.Progress,
		})
	}

	return &ExportRecordedDataState{
		Progress:           resp.Progress,
		FileProgressStatus: fileProgresses,
	}, nil
}

// OverrideSegmentDuration overrides the segment duration for a recording.
// targetDuration and expiration are ISO 8601 duration strings (e.g. "PT10M", "PT1H").
func (c *Client) OverrideSegmentDuration(ctx context.Context, targetDuration, expiration, recordingToken string) error {
	endpoint := c.getRecordingEndpoint()

	type OverrideSegmentDuration struct {
		XMLName               xml.Name `xml:"trc:OverrideSegmentDuration"`
		Xmlns                 string   `xml:"xmlns:trc,attr"`
		TargetSegmentDuration string   `xml:"trc:TargetSegmentDuration"`
		Expiration            string   `xml:"trc:Expiration"`
		RecordingConfiguration string  `xml:"trc:RecordingConfiguration"`
	}

	req := OverrideSegmentDuration{
		Xmlns:                 recordingNamespace,
		TargetSegmentDuration: targetDuration,
		Expiration:            expiration,
		RecordingConfiguration: recordingToken,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("OverrideSegmentDuration failed: %w", err)
	}

	return nil
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
