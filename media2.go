package onvif

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/0x524a/onvif-go/internal/soap"
)

// Media2 service namespace.
const media2Namespace = "http://www.onvif.org/ver20/media/wsdl"

// getMedia2Endpoint returns the media2 endpoint, falling back to the default endpoint if not set.
func (c *Client) getMedia2Endpoint() string {
	if c.media2Endpoint != "" {
		return c.media2Endpoint
	}

	return c.endpoint
}

// GetMedia2ServiceCapabilities retrieves the capabilities of the Media2 service.
func (c *Client) GetMedia2ServiceCapabilities(ctx context.Context) (*Media2ServiceCapabilities, error) {
	endpoint := c.getMedia2Endpoint()

	type GetServiceCapabilities struct {
		XMLName xml.Name `xml:"tr2:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:tr2,attr"`
	}

	type GetServiceCapabilitiesResponse struct {
		XMLName      xml.Name `xml:"GetServiceCapabilitiesResponse"`
		Capabilities struct {
			SnapshotUri     bool `xml:"SnapshotUri,attr"`
			Rotation        bool `xml:"Rotation,attr"`
			VideoSourceMode bool `xml:"VideoSourceMode,attr"`
			OSD             bool `xml:"OSD,attr"`
			Mask            bool `xml:"Mask,attr"`
			SourceMask      bool `xml:"SourceMask,attr"`
		} `xml:"Capabilities"`
	}

	req := GetServiceCapabilities{
		Xmlns: media2Namespace,
	}

	var resp GetServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMedia2ServiceCapabilities failed: %w", err)
	}

	return &Media2ServiceCapabilities{
		SnapshotUri:     resp.Capabilities.SnapshotUri,
		Rotation:        resp.Capabilities.Rotation,
		VideoSourceMode: resp.Capabilities.VideoSourceMode,
		OSD:             resp.Capabilities.OSD,
		Mask:            resp.Capabilities.Mask,
		SourceMask:      resp.Capabilities.SourceMask,
	}, nil
}

// GetMedia2Profiles retrieves media profiles from the Media2 service.
// Both token and configType are optional filters.
func (c *Client) GetMedia2Profiles(ctx context.Context, token *string, configType *string) ([]*Media2Profile, error) {
	endpoint := c.getMedia2Endpoint()

	type GetProfiles struct {
		XMLName    xml.Name `xml:"tr2:GetProfiles"`
		Xmlns      string   `xml:"xmlns:tr2,attr"`
		Token      *string  `xml:"tr2:Token,omitempty"`
		ConfigType *string  `xml:"tr2:Type,omitempty"`
	}

	type GetProfilesResponse struct {
		XMLName  xml.Name `xml:"GetProfilesResponse"`
		Profiles []struct {
			Token  string `xml:"token,attr"`
			Fixed  bool   `xml:"fixed,attr"`
			Name   string `xml:"Name"`
			Configurations *struct {
				VideoSource *struct {
					Token string `xml:"token,attr"`
					Name  string `xml:"Name"`
				} `xml:"VideoSource"`
				AudioSource *struct {
					Token string `xml:"token,attr"`
					Name  string `xml:"Name"`
				} `xml:"AudioSource"`
				VideoEncoder *struct {
					Token    string `xml:"token,attr"`
					Name     string `xml:"Name"`
					Encoding string `xml:"Encoding"`
				} `xml:"VideoEncoder"`
				AudioEncoder *struct {
					Token    string `xml:"token,attr"`
					Name     string `xml:"Name"`
					Encoding string `xml:"Encoding"`
				} `xml:"AudioEncoder"`
				PTZ *struct {
					Token string `xml:"token,attr"`
					Name  string `xml:"Name"`
				} `xml:"PTZ"`
			} `xml:"Configurations"`
		} `xml:"Profiles"`
	}

	req := GetProfiles{
		Xmlns:      media2Namespace,
		Token:      token,
		ConfigType: configType,
	}

	var resp GetProfilesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMedia2Profiles failed: %w", err)
	}

	profiles := make([]*Media2Profile, len(resp.Profiles))
	for i, p := range resp.Profiles {
		profile := &Media2Profile{
			Token: p.Token,
			Name:  p.Name,
			Fixed: p.Fixed,
		}

		if p.Configurations != nil {
			profile.Configurations = &Media2Configurations{}

			if p.Configurations.VideoSource != nil {
				profile.Configurations.VideoSource = &VideoSourceConfiguration{
					Token: p.Configurations.VideoSource.Token,
					Name:  p.Configurations.VideoSource.Name,
				}
			}

			if p.Configurations.AudioSource != nil {
				profile.Configurations.AudioSource = &AudioSourceConfiguration{
					Token: p.Configurations.AudioSource.Token,
					Name:  p.Configurations.AudioSource.Name,
				}
			}

			if p.Configurations.VideoEncoder != nil {
				profile.Configurations.VideoEncoder = &VideoEncoderConfiguration{
					Token:    p.Configurations.VideoEncoder.Token,
					Name:     p.Configurations.VideoEncoder.Name,
					Encoding: p.Configurations.VideoEncoder.Encoding,
				}
			}

			if p.Configurations.AudioEncoder != nil {
				profile.Configurations.AudioEncoder = &AudioEncoderConfiguration{
					Token:    p.Configurations.AudioEncoder.Token,
					Name:     p.Configurations.AudioEncoder.Name,
					Encoding: p.Configurations.AudioEncoder.Encoding,
				}
			}

			if p.Configurations.PTZ != nil {
				profile.Configurations.PTZ = &PTZConfiguration{
					Token: p.Configurations.PTZ.Token,
					Name:  p.Configurations.PTZ.Name,
				}
			}
		}

		profiles[i] = profile
	}

	return profiles, nil
}

// CreateMedia2Profile creates a new media profile in the Media2 service.
// Returns the token of the newly created profile.
func (c *Client) CreateMedia2Profile(ctx context.Context, name string, configurations []*Media2ConfigurationRef) (string, error) {
	endpoint := c.getMedia2Endpoint()

	type ConfigurationRef struct {
		Type  string `xml:"tr2:Type"`
		Token string `xml:"tr2:Token"`
	}

	type CreateProfile struct {
		XMLName        xml.Name           `xml:"tr2:CreateProfile"`
		Xmlns          string             `xml:"xmlns:tr2,attr"`
		Name           string             `xml:"tr2:Name"`
		Configurations []ConfigurationRef `xml:"tr2:Configuration,omitempty"`
	}

	type CreateProfileResponse struct {
		XMLName xml.Name `xml:"CreateProfileResponse"`
		Token   string   `xml:"Token"`
	}

	req := CreateProfile{
		Xmlns: media2Namespace,
		Name:  name,
	}

	for _, cfg := range configurations {
		req.Configurations = append(req.Configurations, ConfigurationRef{
			Type:  cfg.Type,
			Token: cfg.Token,
		})
	}

	var resp CreateProfileResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateMedia2Profile failed: %w", err)
	}

	return resp.Token, nil
}

// DeleteMedia2Profile deletes a media profile from the Media2 service.
func (c *Client) DeleteMedia2Profile(ctx context.Context, token string) error {
	endpoint := c.getMedia2Endpoint()

	type DeleteProfile struct {
		XMLName xml.Name `xml:"tr2:DeleteProfile"`
		Xmlns   string   `xml:"xmlns:tr2,attr"`
		Token   string   `xml:"tr2:Token"`
	}

	type DeleteProfileResponse struct {
		XMLName xml.Name `xml:"DeleteProfileResponse"`
	}

	req := DeleteProfile{
		Xmlns: media2Namespace,
		Token: token,
	}

	var resp DeleteProfileResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteMedia2Profile failed: %w", err)
	}

	return nil
}

// AddMedia2Configuration adds a configuration to a media profile in the Media2 service.
func (c *Client) AddMedia2Configuration(ctx context.Context, profileToken string, config *Media2ConfigurationRef) error {
	endpoint := c.getMedia2Endpoint()

	type ConfigurationRef struct {
		Type  string `xml:"tr2:Type"`
		Token string `xml:"tr2:Token"`
	}

	type AddConfiguration struct {
		XMLName       xml.Name         `xml:"tr2:AddConfiguration"`
		Xmlns         string           `xml:"xmlns:tr2,attr"`
		ProfileToken  string           `xml:"tr2:ProfileToken"`
		Configuration ConfigurationRef `xml:"tr2:Configuration"`
	}

	type AddConfigurationResponse struct {
		XMLName xml.Name `xml:"AddConfigurationResponse"`
	}

	req := AddConfiguration{
		Xmlns:        media2Namespace,
		ProfileToken: profileToken,
		Configuration: ConfigurationRef{
			Type:  config.Type,
			Token: config.Token,
		},
	}

	var resp AddConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("AddMedia2Configuration failed: %w", err)
	}

	return nil
}

// RemoveMedia2Configuration removes a configuration from a media profile in the Media2 service.
func (c *Client) RemoveMedia2Configuration(ctx context.Context, profileToken string, config *Media2ConfigurationRef) error {
	endpoint := c.getMedia2Endpoint()

	type ConfigurationRef struct {
		Type  string `xml:"tr2:Type"`
		Token string `xml:"tr2:Token"`
	}

	type RemoveConfiguration struct {
		XMLName       xml.Name         `xml:"tr2:RemoveConfiguration"`
		Xmlns         string           `xml:"xmlns:tr2,attr"`
		ProfileToken  string           `xml:"tr2:ProfileToken"`
		Configuration ConfigurationRef `xml:"tr2:Configuration"`
	}

	type RemoveConfigurationResponse struct {
		XMLName xml.Name `xml:"RemoveConfigurationResponse"`
	}

	req := RemoveConfiguration{
		Xmlns:        media2Namespace,
		ProfileToken: profileToken,
		Configuration: ConfigurationRef{
			Type:  config.Type,
			Token: config.Token,
		},
	}

	var resp RemoveConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("RemoveMedia2Configuration failed: %w", err)
	}

	return nil
}

// GetMedia2StreamUri retrieves the stream URI for a profile via the Media2 service.
func (c *Client) GetMedia2StreamUri(ctx context.Context, protocol, profileToken string) (string, error) {
	endpoint := c.getMedia2Endpoint()

	type GetStreamUri struct {
		XMLName      xml.Name `xml:"tr2:GetStreamUri"`
		Xmlns        string   `xml:"xmlns:tr2,attr"`
		Protocol     string   `xml:"tr2:Protocol"`
		ProfileToken string   `xml:"tr2:ProfileToken"`
	}

	type GetStreamUriResponse struct {
		XMLName xml.Name `xml:"GetStreamUriResponse"`
		Uri     string   `xml:"Uri"`
	}

	req := GetStreamUri{
		Xmlns:        media2Namespace,
		Protocol:     protocol,
		ProfileToken: profileToken,
	}

	var resp GetStreamUriResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("GetMedia2StreamUri failed: %w", err)
	}

	return resp.Uri, nil
}

// GetMedia2SnapshotUri retrieves the snapshot URI for a profile via the Media2 service.
func (c *Client) GetMedia2SnapshotUri(ctx context.Context, profileToken string) (string, error) {
	endpoint := c.getMedia2Endpoint()

	type GetSnapshotUri struct {
		XMLName      xml.Name `xml:"tr2:GetSnapshotUri"`
		Xmlns        string   `xml:"xmlns:tr2,attr"`
		ProfileToken string   `xml:"tr2:ProfileToken"`
	}

	type GetSnapshotUriResponse struct {
		XMLName xml.Name `xml:"GetSnapshotUriResponse"`
		Uri     string   `xml:"Uri"`
	}

	req := GetSnapshotUri{
		Xmlns:        media2Namespace,
		ProfileToken: profileToken,
	}

	var resp GetSnapshotUriResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("GetMedia2SnapshotUri failed: %w", err)
	}

	return resp.Uri, nil
}

// StartMedia2MulticastStreaming starts multicast streaming for a profile via the Media2 service.
func (c *Client) StartMedia2MulticastStreaming(ctx context.Context, profileToken string) error {
	endpoint := c.getMedia2Endpoint()

	type StartMulticastStreaming struct {
		XMLName      xml.Name `xml:"tr2:StartMulticastStreaming"`
		Xmlns        string   `xml:"xmlns:tr2,attr"`
		ProfileToken string   `xml:"tr2:ProfileToken"`
	}

	type StartMulticastStreamingResponse struct {
		XMLName xml.Name `xml:"StartMulticastStreamingResponse"`
	}

	req := StartMulticastStreaming{
		Xmlns:        media2Namespace,
		ProfileToken: profileToken,
	}

	var resp StartMulticastStreamingResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("StartMedia2MulticastStreaming failed: %w", err)
	}

	return nil
}

// StopMedia2MulticastStreaming stops multicast streaming for a profile via the Media2 service.
func (c *Client) StopMedia2MulticastStreaming(ctx context.Context, profileToken string) error {
	endpoint := c.getMedia2Endpoint()

	type StopMulticastStreaming struct {
		XMLName      xml.Name `xml:"tr2:StopMulticastStreaming"`
		Xmlns        string   `xml:"xmlns:tr2,attr"`
		ProfileToken string   `xml:"tr2:ProfileToken"`
	}

	type StopMulticastStreamingResponse struct {
		XMLName xml.Name `xml:"StopMulticastStreamingResponse"`
	}

	req := StopMulticastStreaming{
		Xmlns:        media2Namespace,
		ProfileToken: profileToken,
	}

	var resp StopMulticastStreamingResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("StopMedia2MulticastStreaming failed: %w", err)
	}

	return nil
}

// SetMedia2SynchronizationPoint sets a synchronization point for a profile via the Media2 service.
func (c *Client) SetMedia2SynchronizationPoint(ctx context.Context, profileToken string) error {
	endpoint := c.getMedia2Endpoint()

	type SetSynchronizationPoint struct {
		XMLName      xml.Name `xml:"tr2:SetSynchronizationPoint"`
		Xmlns        string   `xml:"xmlns:tr2,attr"`
		ProfileToken string   `xml:"tr2:ProfileToken"`
	}

	type SetSynchronizationPointResponse struct {
		XMLName xml.Name `xml:"SetSynchronizationPointResponse"`
	}

	req := SetSynchronizationPoint{
		Xmlns:        media2Namespace,
		ProfileToken: profileToken,
	}

	var resp SetSynchronizationPointResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetMedia2SynchronizationPoint failed: %w", err)
	}

	return nil
}

// GetMedia2VideoSourceConfigurations retrieves video source configurations from the Media2 service.
// Both configToken and profileToken are optional filters.
func (c *Client) GetMedia2VideoSourceConfigurations(ctx context.Context, configToken, profileToken *string) ([]*VideoSourceConfiguration, error) {
	endpoint := c.getMedia2Endpoint()

	type GetVideoSourceConfigurations struct {
		XMLName            xml.Name `xml:"tr2:GetVideoSourceConfigurations"`
		Xmlns              string   `xml:"xmlns:tr2,attr"`
		ConfigurationToken *string  `xml:"tr2:ConfigurationToken,omitempty"`
		ProfileToken       *string  `xml:"tr2:ProfileToken,omitempty"`
	}

	type GetVideoSourceConfigurationsResponse struct {
		XMLName        xml.Name `xml:"GetVideoSourceConfigurationsResponse"`
		Configurations []struct {
			Token       string `xml:"token,attr"`
			Name        string `xml:"Name"`
			UseCount    int    `xml:"UseCount"`
			SourceToken string `xml:"SourceToken"`
			Bounds      *struct {
				X      int `xml:"x,attr"`
				Y      int `xml:"y,attr"`
				Width  int `xml:"width,attr"`
				Height int `xml:"height,attr"`
			} `xml:"Bounds"`
		} `xml:"Configurations"`
	}

	req := GetVideoSourceConfigurations{
		Xmlns:              media2Namespace,
		ConfigurationToken: configToken,
		ProfileToken:       profileToken,
	}

	var resp GetVideoSourceConfigurationsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMedia2VideoSourceConfigurations failed: %w", err)
	}

	configs := make([]*VideoSourceConfiguration, len(resp.Configurations))
	for i, cfg := range resp.Configurations {
		config := &VideoSourceConfiguration{
			Token:       cfg.Token,
			Name:        cfg.Name,
			UseCount:    cfg.UseCount,
			SourceToken: cfg.SourceToken,
		}

		if cfg.Bounds != nil {
			config.Bounds = &IntRectangle{
				X:      cfg.Bounds.X,
				Y:      cfg.Bounds.Y,
				Width:  cfg.Bounds.Width,
				Height: cfg.Bounds.Height,
			}
		}

		configs[i] = config
	}

	return configs, nil
}

// GetMedia2VideoEncoderConfigurations retrieves video encoder configurations from the Media2 service.
// Both configToken and profileToken are optional filters.
func (c *Client) GetMedia2VideoEncoderConfigurations(ctx context.Context, configToken, profileToken *string) ([]*VideoEncoderConfiguration, error) {
	endpoint := c.getMedia2Endpoint()

	type GetVideoEncoderConfigurations struct {
		XMLName            xml.Name `xml:"tr2:GetVideoEncoderConfigurations"`
		Xmlns              string   `xml:"xmlns:tr2,attr"`
		ConfigurationToken *string  `xml:"tr2:ConfigurationToken,omitempty"`
		ProfileToken       *string  `xml:"tr2:ProfileToken,omitempty"`
	}

	type GetVideoEncoderConfigurationsResponse struct {
		XMLName        xml.Name `xml:"GetVideoEncoderConfigurationsResponse"`
		Configurations []struct {
			Token    string `xml:"token,attr"`
			Name     string `xml:"Name"`
			UseCount int    `xml:"UseCount"`
			Encoding string `xml:"Encoding"`
			Quality  float64 `xml:"Quality"`
			Resolution *struct {
				Width  int `xml:"Width"`
				Height int `xml:"Height"`
			} `xml:"Resolution"`
			RateControl *struct {
				FrameRateLimit   int `xml:"FrameRateLimit"`
				EncodingInterval int `xml:"EncodingInterval"`
				BitrateLimit     int `xml:"BitrateLimit"`
			} `xml:"RateControl"`
		} `xml:"Configurations"`
	}

	req := GetVideoEncoderConfigurations{
		Xmlns:              media2Namespace,
		ConfigurationToken: configToken,
		ProfileToken:       profileToken,
	}

	var resp GetVideoEncoderConfigurationsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMedia2VideoEncoderConfigurations failed: %w", err)
	}

	configs := make([]*VideoEncoderConfiguration, len(resp.Configurations))
	for i, cfg := range resp.Configurations {
		config := &VideoEncoderConfiguration{
			Token:    cfg.Token,
			Name:     cfg.Name,
			UseCount: cfg.UseCount,
			Encoding: cfg.Encoding,
			Quality:  cfg.Quality,
		}

		if cfg.Resolution != nil {
			config.Resolution = &VideoResolution{
				Width:  cfg.Resolution.Width,
				Height: cfg.Resolution.Height,
			}
		}

		if cfg.RateControl != nil {
			config.RateControl = &VideoRateControl{
				FrameRateLimit:   cfg.RateControl.FrameRateLimit,
				EncodingInterval: cfg.RateControl.EncodingInterval,
				BitrateLimit:     cfg.RateControl.BitrateLimit,
			}
		}

		configs[i] = config
	}

	return configs, nil
}

// GetMedia2AudioSourceConfigurations retrieves audio source configurations from the Media2 service.
// Both configToken and profileToken are optional filters.
func (c *Client) GetMedia2AudioSourceConfigurations(ctx context.Context, configToken, profileToken *string) ([]*AudioSourceConfiguration, error) {
	endpoint := c.getMedia2Endpoint()

	type GetAudioSourceConfigurations struct {
		XMLName            xml.Name `xml:"tr2:GetAudioSourceConfigurations"`
		Xmlns              string   `xml:"xmlns:tr2,attr"`
		ConfigurationToken *string  `xml:"tr2:ConfigurationToken,omitempty"`
		ProfileToken       *string  `xml:"tr2:ProfileToken,omitempty"`
	}

	type GetAudioSourceConfigurationsResponse struct {
		XMLName        xml.Name `xml:"GetAudioSourceConfigurationsResponse"`
		Configurations []struct {
			Token       string `xml:"token,attr"`
			Name        string `xml:"Name"`
			UseCount    int    `xml:"UseCount"`
			SourceToken string `xml:"SourceToken"`
		} `xml:"Configurations"`
	}

	req := GetAudioSourceConfigurations{
		Xmlns:              media2Namespace,
		ConfigurationToken: configToken,
		ProfileToken:       profileToken,
	}

	var resp GetAudioSourceConfigurationsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMedia2AudioSourceConfigurations failed: %w", err)
	}

	configs := make([]*AudioSourceConfiguration, len(resp.Configurations))
	for i, cfg := range resp.Configurations {
		configs[i] = &AudioSourceConfiguration{
			Token:       cfg.Token,
			Name:        cfg.Name,
			UseCount:    cfg.UseCount,
			SourceToken: cfg.SourceToken,
		}
	}

	return configs, nil
}

// GetMedia2AudioEncoderConfigurations retrieves audio encoder configurations from the Media2 service.
// Both configToken and profileToken are optional filters.
func (c *Client) GetMedia2AudioEncoderConfigurations(ctx context.Context, configToken, profileToken *string) ([]*AudioEncoderConfiguration, error) {
	endpoint := c.getMedia2Endpoint()

	type GetAudioEncoderConfigurations struct {
		XMLName            xml.Name `xml:"tr2:GetAudioEncoderConfigurations"`
		Xmlns              string   `xml:"xmlns:tr2,attr"`
		ConfigurationToken *string  `xml:"tr2:ConfigurationToken,omitempty"`
		ProfileToken       *string  `xml:"tr2:ProfileToken,omitempty"`
	}

	type GetAudioEncoderConfigurationsResponse struct {
		XMLName        xml.Name `xml:"GetAudioEncoderConfigurationsResponse"`
		Configurations []struct {
			Token      string `xml:"token,attr"`
			Name       string `xml:"Name"`
			UseCount   int    `xml:"UseCount"`
			Encoding   string `xml:"Encoding"`
			Bitrate    int    `xml:"Bitrate"`
			SampleRate int    `xml:"SampleRate"`
		} `xml:"Configurations"`
	}

	req := GetAudioEncoderConfigurations{
		Xmlns:              media2Namespace,
		ConfigurationToken: configToken,
		ProfileToken:       profileToken,
	}

	var resp GetAudioEncoderConfigurationsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMedia2AudioEncoderConfigurations failed: %w", err)
	}

	configs := make([]*AudioEncoderConfiguration, len(resp.Configurations))
	for i, cfg := range resp.Configurations {
		configs[i] = &AudioEncoderConfiguration{
			Token:      cfg.Token,
			Name:       cfg.Name,
			UseCount:   cfg.UseCount,
			Encoding:   cfg.Encoding,
			Bitrate:    cfg.Bitrate,
			SampleRate: cfg.SampleRate,
		}
	}

	return configs, nil
}
