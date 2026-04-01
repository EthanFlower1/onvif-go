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

// GetMedia2Masks retrieves privacy masks from the Media2 service.
// configToken is an optional filter.
func (c *Client) GetMedia2Masks(ctx context.Context, configToken *string) ([]*Mask, error) {
	endpoint := c.getMedia2Endpoint()

	type GetMasks struct {
		XMLName            xml.Name `xml:"tr2:GetMasks"`
		Xmlns              string   `xml:"xmlns:tr2,attr"`
		ConfigurationToken *string  `xml:"tr2:ConfigurationToken,omitempty"`
	}

	type MaskPoint struct {
		X float64 `xml:"x,attr"`
		Y float64 `xml:"y,attr"`
	}

	type MaskPolygon struct {
		Points []MaskPoint `xml:"Point"`
	}

	type MaskEntry struct {
		Token              string      `xml:"token,attr"`
		ConfigurationToken string      `xml:"ConfigurationToken,attr"`
		Polygon            MaskPolygon `xml:"Polygon"`
		Type               string      `xml:"Type"`
		Enabled            bool        `xml:"Enabled"`
	}

	type GetMasksResponse struct {
		XMLName xml.Name    `xml:"GetMasksResponse"`
		Masks   []MaskEntry `xml:"Masks"`
	}

	req := GetMasks{
		Xmlns:              media2Namespace,
		ConfigurationToken: configToken,
	}

	var resp GetMasksResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMedia2Masks failed: %w", err)
	}

	masks := make([]*Mask, len(resp.Masks))
	for i, m := range resp.Masks {
		mask := &Mask{
			Token:              m.Token,
			ConfigurationToken: m.ConfigurationToken,
			Type:               m.Type,
			Enabled:            m.Enabled,
		}

		if len(m.Polygon.Points) > 0 {
			poly := &Polygon{}
			for _, pt := range m.Polygon.Points {
				poly.Points = append(poly.Points, &Vector{X: pt.X, Y: pt.Y})
			}

			mask.Polygon = poly
		}

		masks[i] = mask
	}

	return masks, nil
}

// GetMedia2MaskOptions retrieves mask configuration options for a video source configuration.
func (c *Client) GetMedia2MaskOptions(ctx context.Context, configToken string) (*MaskOptions, error) {
	endpoint := c.getMedia2Endpoint()

	type GetMaskOptions struct {
		XMLName            xml.Name `xml:"tr2:GetMaskOptions"`
		Xmlns              string   `xml:"xmlns:tr2,attr"`
		ConfigurationToken string   `xml:"tr2:ConfigurationToken"`
	}

	type GetMaskOptionsResponse struct {
		XMLName xml.Name `xml:"GetMaskOptionsResponse"`
		Options struct {
			MaxMasks        int      `xml:"MaxMasks,attr"`
			MaxPoints       int      `xml:"MaxPoints,attr"`
			Types           []string `xml:"Types"`
			SingleColorOnly bool     `xml:"SingleColorOnly"`
		} `xml:"Options"`
	}

	req := GetMaskOptions{
		Xmlns:              media2Namespace,
		ConfigurationToken: configToken,
	}

	var resp GetMaskOptionsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMedia2MaskOptions failed: %w", err)
	}

	return &MaskOptions{
		MaxMasks:        resp.Options.MaxMasks,
		MaxPoints:       resp.Options.MaxPoints,
		Types:           resp.Options.Types,
		SingleColorOnly: resp.Options.SingleColorOnly,
	}, nil
}

// SetMedia2Mask updates an existing privacy mask.
func (c *Client) SetMedia2Mask(ctx context.Context, mask *Mask) error {
	endpoint := c.getMedia2Endpoint()

	type MaskPoint struct {
		X float64 `xml:"x,attr"`
		Y float64 `xml:"y,attr"`
	}

	type MaskPolygon struct {
		Points []MaskPoint `xml:"tr2:Point"`
	}

	type MaskElement struct {
		Token              string      `xml:"token,attr"`
		ConfigurationToken string      `xml:"ConfigurationToken,attr"`
		Polygon            MaskPolygon `xml:"tr2:Polygon"`
		Type               string      `xml:"tr2:Type"`
		Enabled            bool        `xml:"tr2:Enabled"`
	}

	type SetMask struct {
		XMLName xml.Name    `xml:"tr2:SetMask"`
		Xmlns   string      `xml:"xmlns:tr2,attr"`
		Mask    MaskElement `xml:"tr2:Mask"`
	}

	type SetMaskResponse struct {
		XMLName xml.Name `xml:"SetMaskResponse"`
	}

	reqMask := MaskElement{
		Token:              mask.Token,
		ConfigurationToken: mask.ConfigurationToken,
		Type:               mask.Type,
		Enabled:            mask.Enabled,
	}

	if mask.Polygon != nil {
		for _, pt := range mask.Polygon.Points {
			reqMask.Polygon.Points = append(reqMask.Polygon.Points, MaskPoint{X: pt.X, Y: pt.Y})
		}
	}

	req := SetMask{
		Xmlns: media2Namespace,
		Mask:  reqMask,
	}

	var resp SetMaskResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetMedia2Mask failed: %w", err)
	}

	return nil
}

// CreateMedia2Mask creates a new privacy mask and returns the assigned token.
func (c *Client) CreateMedia2Mask(ctx context.Context, mask *Mask) (string, error) {
	endpoint := c.getMedia2Endpoint()

	type MaskPoint struct {
		X float64 `xml:"x,attr"`
		Y float64 `xml:"y,attr"`
	}

	type MaskPolygon struct {
		Points []MaskPoint `xml:"tr2:Point"`
	}

	type MaskElement struct {
		Token              string      `xml:"token,attr,omitempty"`
		ConfigurationToken string      `xml:"ConfigurationToken,attr"`
		Polygon            MaskPolygon `xml:"tr2:Polygon"`
		Type               string      `xml:"tr2:Type"`
		Enabled            bool        `xml:"tr2:Enabled"`
	}

	type CreateMask struct {
		XMLName xml.Name    `xml:"tr2:CreateMask"`
		Xmlns   string      `xml:"xmlns:tr2,attr"`
		Mask    MaskElement `xml:"tr2:Mask"`
	}

	type CreateMaskResponse struct {
		XMLName xml.Name `xml:"CreateMaskResponse"`
		Token   string   `xml:"Token"`
	}

	reqMask := MaskElement{
		ConfigurationToken: mask.ConfigurationToken,
		Type:               mask.Type,
		Enabled:            mask.Enabled,
	}

	if mask.Polygon != nil {
		for _, pt := range mask.Polygon.Points {
			reqMask.Polygon.Points = append(reqMask.Polygon.Points, MaskPoint{X: pt.X, Y: pt.Y})
		}
	}

	req := CreateMask{
		Xmlns: media2Namespace,
		Mask:  reqMask,
	}

	var resp CreateMaskResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("CreateMedia2Mask failed: %w", err)
	}

	return resp.Token, nil
}

// DeleteMedia2Mask deletes a privacy mask by token.
func (c *Client) DeleteMedia2Mask(ctx context.Context, token string) error {
	endpoint := c.getMedia2Endpoint()

	type DeleteMask struct {
		XMLName xml.Name `xml:"tr2:DeleteMask"`
		Xmlns   string   `xml:"xmlns:tr2,attr"`
		Token   string   `xml:"tr2:Token"`
	}

	type DeleteMaskResponse struct {
		XMLName xml.Name `xml:"DeleteMaskResponse"`
	}

	req := DeleteMask{
		Xmlns: media2Namespace,
		Token: token,
	}

	var resp DeleteMaskResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteMedia2Mask failed: %w", err)
	}

	return nil
}

// GetMedia2AudioClips retrieves all audio clips from the Media2 service.
func (c *Client) GetMedia2AudioClips(ctx context.Context) ([]*AudioClip, error) {
	endpoint := c.getMedia2Endpoint()

	type GetAudioClips struct {
		XMLName xml.Name `xml:"tr2:GetAudioClips"`
		Xmlns   string   `xml:"xmlns:tr2,attr"`
	}

	type ClipEntry struct {
		Token    string `xml:"token,attr"`
		Name     string `xml:"Name"`
		MediaURI string `xml:"MediaUri"`
	}

	type GetAudioClipsResponse struct {
		XMLName    xml.Name    `xml:"GetAudioClipsResponse"`
		AudioClips []ClipEntry `xml:"AudioClips"`
	}

	req := GetAudioClips{Xmlns: media2Namespace}

	var resp GetAudioClipsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMedia2AudioClips failed: %w", err)
	}

	clips := make([]*AudioClip, len(resp.AudioClips))
	for i, c := range resp.AudioClips {
		clips[i] = &AudioClip{
			Token:    c.Token,
			Name:     c.Name,
			MediaURI: c.MediaURI,
		}
	}

	return clips, nil
}

// AddMedia2AudioClip uploads or registers a new audio clip and returns the assigned token.
func (c *Client) AddMedia2AudioClip(ctx context.Context, clip *AudioClip) (string, error) {
	endpoint := c.getMedia2Endpoint()

	type ClipElement struct {
		Name     string `xml:"tr2:Name"`
		MediaURI string `xml:"tr2:MediaUri"`
	}

	type AddAudioClip struct {
		XMLName xml.Name    `xml:"tr2:AddAudioClip"`
		Xmlns   string      `xml:"xmlns:tr2,attr"`
		Clip    ClipElement `xml:"tr2:AudioClip"`
	}

	type AddAudioClipResponse struct {
		XMLName xml.Name `xml:"AddAudioClipResponse"`
		Token   string   `xml:"Token"`
	}

	req := AddAudioClip{
		Xmlns: media2Namespace,
		Clip: ClipElement{
			Name:     clip.Name,
			MediaURI: clip.MediaURI,
		},
	}

	var resp AddAudioClipResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("AddMedia2AudioClip failed: %w", err)
	}

	return resp.Token, nil
}

// SetMedia2AudioClip updates an existing audio clip.
func (c *Client) SetMedia2AudioClip(ctx context.Context, clip *AudioClip) error {
	endpoint := c.getMedia2Endpoint()

	type ClipElement struct {
		Token    string `xml:"token,attr"`
		Name     string `xml:"tr2:Name"`
		MediaURI string `xml:"tr2:MediaUri"`
	}

	type SetAudioClip struct {
		XMLName xml.Name    `xml:"tr2:SetAudioClip"`
		Xmlns   string      `xml:"xmlns:tr2,attr"`
		Clip    ClipElement `xml:"tr2:AudioClip"`
	}

	type SetAudioClipResponse struct {
		XMLName xml.Name `xml:"SetAudioClipResponse"`
	}

	req := SetAudioClip{
		Xmlns: media2Namespace,
		Clip: ClipElement{
			Token:    clip.Token,
			Name:     clip.Name,
			MediaURI: clip.MediaURI,
		},
	}

	var resp SetAudioClipResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetMedia2AudioClip failed: %w", err)
	}

	return nil
}

// DeleteMedia2AudioClip deletes an audio clip by token.
func (c *Client) DeleteMedia2AudioClip(ctx context.Context, token string) error {
	endpoint := c.getMedia2Endpoint()

	type DeleteAudioClip struct {
		XMLName xml.Name `xml:"tr2:DeleteAudioClip"`
		Xmlns   string   `xml:"xmlns:tr2,attr"`
		Token   string   `xml:"tr2:Token"`
	}

	type DeleteAudioClipResponse struct {
		XMLName xml.Name `xml:"DeleteAudioClipResponse"`
	}

	req := DeleteAudioClip{
		Xmlns: media2Namespace,
		Token: token,
	}

	var resp DeleteAudioClipResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeleteMedia2AudioClip failed: %w", err)
	}

	return nil
}

// PlayMedia2AudioClip plays an audio clip on the specified profile's audio output.
func (c *Client) PlayMedia2AudioClip(ctx context.Context, token, profileToken string) error {
	endpoint := c.getMedia2Endpoint()

	type PlayAudioClip struct {
		XMLName      xml.Name `xml:"tr2:PlayAudioClip"`
		Xmlns        string   `xml:"xmlns:tr2,attr"`
		Token        string   `xml:"tr2:Token"`
		ProfileToken string   `xml:"tr2:ProfileToken"`
	}

	type PlayAudioClipResponse struct {
		XMLName xml.Name `xml:"PlayAudioClipResponse"`
	}

	req := PlayAudioClip{
		Xmlns:        media2Namespace,
		Token:        token,
		ProfileToken: profileToken,
	}

	var resp PlayAudioClipResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("PlayMedia2AudioClip failed: %w", err)
	}

	return nil
}

// GetMedia2PlayingAudioClips returns the currently playing audio clips.
func (c *Client) GetMedia2PlayingAudioClips(ctx context.Context) ([]*AudioClip, error) {
	endpoint := c.getMedia2Endpoint()

	type GetPlayingAudioClips struct {
		XMLName xml.Name `xml:"tr2:GetPlayingAudioClips"`
		Xmlns   string   `xml:"xmlns:tr2,attr"`
	}

	type ClipEntry struct {
		Token    string `xml:"token,attr"`
		Name     string `xml:"Name"`
		MediaURI string `xml:"MediaUri"`
	}

	type GetPlayingAudioClipsResponse struct {
		XMLName    xml.Name    `xml:"GetPlayingAudioClipsResponse"`
		AudioClips []ClipEntry `xml:"AudioClips"`
	}

	req := GetPlayingAudioClips{Xmlns: media2Namespace}

	var resp GetPlayingAudioClipsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMedia2PlayingAudioClips failed: %w", err)
	}

	clips := make([]*AudioClip, len(resp.AudioClips))
	for i, c := range resp.AudioClips {
		clips[i] = &AudioClip{
			Token:    c.Token,
			Name:     c.Name,
			MediaURI: c.MediaURI,
		}
	}

	return clips, nil
}

// GetMedia2WebRTCConfigurations retrieves WebRTC streaming configuration from the Media2 service.
func (c *Client) GetMedia2WebRTCConfigurations(ctx context.Context) (*WebRTCConfiguration, error) {
	endpoint := c.getMedia2Endpoint()

	type GetWebRTCConfigurations struct {
		XMLName xml.Name `xml:"tr2:GetWebRTCConfigurations"`
		Xmlns   string   `xml:"xmlns:tr2,attr"`
	}

	type GetWebRTCConfigurationsResponse struct {
		XMLName        xml.Name `xml:"GetWebRTCConfigurationsResponse"`
		Configurations struct {
			SignalingServerURI string `xml:"SignalingServerURI"`
			STUNServer        string `xml:"STUNServer"`
			TURNServer        string `xml:"TURNServer"`
		} `xml:"Configurations"`
	}

	req := GetWebRTCConfigurations{Xmlns: media2Namespace}

	var resp GetWebRTCConfigurationsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMedia2WebRTCConfigurations failed: %w", err)
	}

	return &WebRTCConfiguration{
		SignalingServerURI: resp.Configurations.SignalingServerURI,
		STUNServer:        resp.Configurations.STUNServer,
		TURNServer:        resp.Configurations.TURNServer,
	}, nil
}

// SetMedia2WebRTCConfigurations updates WebRTC streaming configuration.
func (c *Client) SetMedia2WebRTCConfigurations(ctx context.Context, config *WebRTCConfiguration) error {
	endpoint := c.getMedia2Endpoint()

	type ConfigElement struct {
		SignalingServerURI string `xml:"tr2:SignalingServerURI,omitempty"`
		STUNServer        string `xml:"tr2:STUNServer,omitempty"`
		TURNServer        string `xml:"tr2:TURNServer,omitempty"`
	}

	type SetWebRTCConfigurations struct {
		XMLName        xml.Name      `xml:"tr2:SetWebRTCConfigurations"`
		Xmlns          string        `xml:"xmlns:tr2,attr"`
		Configurations ConfigElement `xml:"tr2:Configurations"`
	}

	type SetWebRTCConfigurationsResponse struct {
		XMLName xml.Name `xml:"SetWebRTCConfigurationsResponse"`
	}

	req := SetWebRTCConfigurations{
		Xmlns: media2Namespace,
		Configurations: ConfigElement{
			SignalingServerURI: config.SignalingServerURI,
			STUNServer:        config.STUNServer,
			TURNServer:        config.TURNServer,
		},
	}

	var resp SetWebRTCConfigurationsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetMedia2WebRTCConfigurations failed: %w", err)
	}

	return nil
}

// GetMedia2MulticastAudioDecoderConfigurations retrieves multicast audio decoder configurations.
// Both configToken and profileToken are optional filters.
func (c *Client) GetMedia2MulticastAudioDecoderConfigurations(ctx context.Context, configToken, profileToken *string) ([]*MulticastAudioDecoderConfiguration, error) {
	endpoint := c.getMedia2Endpoint()

	type GetMulticastAudioDecoderConfigurations struct {
		XMLName            xml.Name `xml:"tr2:GetMulticastAudioDecoderConfigurations"`
		Xmlns              string   `xml:"xmlns:tr2,attr"`
		ConfigurationToken *string  `xml:"tr2:ConfigurationToken,omitempty"`
		ProfileToken       *string  `xml:"tr2:ProfileToken,omitempty"`
	}

	type GetMulticastAudioDecoderConfigurationsResponse struct {
		XMLName        xml.Name `xml:"GetMulticastAudioDecoderConfigurationsResponse"`
		Configurations []struct {
			Token          string `xml:"token,attr"`
			Name           string `xml:"Name"`
			UseCount       int    `xml:"UseCount"`
			SessionTimeout string `xml:"SessionTimeout"`
			Multicast      *struct {
				Address struct {
					Type        string `xml:"Type"`
					IPv4Address string `xml:"IPv4Address"`
					IPv6Address string `xml:"IPv6Address"`
				} `xml:"Address"`
				Port      int  `xml:"Port"`
				TTL       int  `xml:"TTL"`
				AutoStart bool `xml:"AutoStart"`
			} `xml:"Multicast"`
		} `xml:"Configurations"`
	}

	req := GetMulticastAudioDecoderConfigurations{
		Xmlns:              media2Namespace,
		ConfigurationToken: configToken,
		ProfileToken:       profileToken,
	}

	var resp GetMulticastAudioDecoderConfigurationsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMedia2MulticastAudioDecoderConfigurations failed: %w", err)
	}

	configs := make([]*MulticastAudioDecoderConfiguration, len(resp.Configurations))
	for i, cfg := range resp.Configurations {
		config := &MulticastAudioDecoderConfiguration{
			Token:          cfg.Token,
			Name:           cfg.Name,
			UseCount:       cfg.UseCount,
			SessionTimeout: cfg.SessionTimeout,
		}

		if cfg.Multicast != nil {
			config.Multicast = &MulticastConfiguration{
				Port:      cfg.Multicast.Port,
				TTL:       cfg.Multicast.TTL,
				AutoStart: cfg.Multicast.AutoStart,
				Address: &IPAddress{
					Type:        cfg.Multicast.Address.Type,
					IPv4Address: cfg.Multicast.Address.IPv4Address,
					IPv6Address: cfg.Multicast.Address.IPv6Address,
				},
			}
		}

		configs[i] = config
	}

	return configs, nil
}

// GetMedia2MulticastAudioDecoderConfigurationOptions retrieves multicast audio decoder configuration options.
// Both configToken and profileToken are optional filters.
func (c *Client) GetMedia2MulticastAudioDecoderConfigurationOptions(ctx context.Context, configToken, profileToken *string) (interface{}, error) {
	endpoint := c.getMedia2Endpoint()

	type GetMulticastAudioDecoderConfigurationOptions struct {
		XMLName            xml.Name `xml:"tr2:GetMulticastAudioDecoderConfigurationOptions"`
		Xmlns              string   `xml:"xmlns:tr2,attr"`
		ConfigurationToken *string  `xml:"tr2:ConfigurationToken,omitempty"`
		ProfileToken       *string  `xml:"tr2:ProfileToken,omitempty"`
	}

	type GetMulticastAudioDecoderConfigurationOptionsResponse struct {
		XMLName xml.Name   `xml:"GetMulticastAudioDecoderConfigurationOptionsResponse"`
		Options []xml.Name `xml:",any"`
	}

	req := GetMulticastAudioDecoderConfigurationOptions{
		Xmlns:              media2Namespace,
		ConfigurationToken: configToken,
		ProfileToken:       profileToken,
	}

	var resp GetMulticastAudioDecoderConfigurationOptionsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMedia2MulticastAudioDecoderConfigurationOptions failed: %w", err)
	}

	return &resp, nil
}

// SetMedia2MulticastAudioDecoderConfiguration updates a multicast audio decoder configuration.
func (c *Client) SetMedia2MulticastAudioDecoderConfiguration(ctx context.Context, config *MulticastAudioDecoderConfiguration) error {
	endpoint := c.getMedia2Endpoint()

	type MulticastAddressElement struct {
		Type        string `xml:"tr2:Type,omitempty"`
		IPv4Address string `xml:"tr2:IPv4Address,omitempty"`
		IPv6Address string `xml:"tr2:IPv6Address,omitempty"`
	}

	type MulticastElement struct {
		Address   MulticastAddressElement `xml:"tr2:Address"`
		Port      int                     `xml:"tr2:Port,omitempty"`
		TTL       int                     `xml:"tr2:TTL,omitempty"`
		AutoStart bool                    `xml:"tr2:AutoStart"`
	}

	type ConfigElement struct {
		Token          string            `xml:"token,attr"`
		Name           string            `xml:"tr2:Name"`
		UseCount       int               `xml:"tr2:UseCount,omitempty"`
		SessionTimeout string            `xml:"tr2:SessionTimeout,omitempty"`
		Multicast      *MulticastElement `xml:"tr2:Multicast,omitempty"`
	}

	type SetMulticastAudioDecoderConfiguration struct {
		XMLName xml.Name      `xml:"tr2:SetMulticastAudioDecoderConfiguration"`
		Xmlns   string        `xml:"xmlns:tr2,attr"`
		Config  ConfigElement `xml:"tr2:Configuration"`
	}

	type SetMulticastAudioDecoderConfigurationResponse struct {
		XMLName xml.Name `xml:"SetMulticastAudioDecoderConfigurationResponse"`
	}

	reqConfig := ConfigElement{
		Token:          config.Token,
		Name:           config.Name,
		UseCount:       config.UseCount,
		SessionTimeout: config.SessionTimeout,
	}

	if config.Multicast != nil {
		addrElem := MulticastAddressElement{}
		if config.Multicast.Address != nil {
			addrElem.Type = config.Multicast.Address.Type
			addrElem.IPv4Address = config.Multicast.Address.IPv4Address
			addrElem.IPv6Address = config.Multicast.Address.IPv6Address
		}

		reqConfig.Multicast = &MulticastElement{
			Port:      config.Multicast.Port,
			TTL:       config.Multicast.TTL,
			AutoStart: config.Multicast.AutoStart,
			Address:   addrElem,
		}
	}

	req := SetMulticastAudioDecoderConfiguration{
		Xmlns:  media2Namespace,
		Config: reqConfig,
	}

	var resp SetMulticastAudioDecoderConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetMedia2MulticastAudioDecoderConfiguration failed: %w", err)
	}

	return nil
}

// SetMedia2EQPreset sets the equalizer preset for an audio output.
func (c *Client) SetMedia2EQPreset(ctx context.Context, audioOutputToken, eqPresetToken string) error {
	endpoint := c.getMedia2Endpoint()

	type SetEQPreset struct {
		XMLName          xml.Name `xml:"tr2:SetEQPreset"`
		Xmlns            string   `xml:"xmlns:tr2,attr"`
		AudioOutputToken string   `xml:"tr2:AudioOutputToken"`
		EQPresetToken    string   `xml:"tr2:EQPresetToken"`
	}

	type SetEQPresetResponse struct {
		XMLName xml.Name `xml:"SetEQPresetResponse"`
	}

	req := SetEQPreset{
		Xmlns:            media2Namespace,
		AudioOutputToken: audioOutputToken,
		EQPresetToken:    eqPresetToken,
	}

	var resp SetEQPresetResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetMedia2EQPreset failed: %w", err)
	}

	return nil
}

// GetMedia2VideoEncoderInstances retrieves the number of available encoder instances for a video source configuration.
func (c *Client) GetMedia2VideoEncoderInstances(ctx context.Context, configToken string) (*VideoEncoderInstances, error) {
	endpoint := c.getMedia2Endpoint()

	type GetVideoEncoderInstances struct {
		XMLName            xml.Name `xml:"tr2:GetVideoEncoderInstances"`
		Xmlns              string   `xml:"xmlns:tr2,attr"`
		ConfigurationToken string   `xml:"tr2:ConfigurationToken"`
	}

	type GetVideoEncoderInstancesResponse struct {
		XMLName xml.Name `xml:"GetVideoEncoderInstancesResponse"`
		Info    struct {
			Total int  `xml:"Total"`
			JPEG  *int `xml:"JPEG"`
			H264  *int `xml:"H264"`
			MPEG4 *int `xml:"MPEG4"`
		} `xml:"Info"`
	}

	req := GetVideoEncoderInstances{
		Xmlns:              media2Namespace,
		ConfigurationToken: configToken,
	}

	var resp GetVideoEncoderInstancesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMedia2VideoEncoderInstances failed: %w", err)
	}

	return &VideoEncoderInstances{
		Total: resp.Info.Total,
		JPEG:  resp.Info.JPEG,
		H264:  resp.Info.H264,
		MPEG4: resp.Info.MPEG4,
	}, nil
}

// GetMedia2AnalyticsConfigurations retrieves video analytics configurations from the Media2 service.
// Both configToken and profileToken are optional filters.
func (c *Client) GetMedia2AnalyticsConfigurations(ctx context.Context, configToken, profileToken *string) ([]*VideoAnalyticsConfiguration, error) {
	endpoint := c.getMedia2Endpoint()

	type GetAnalyticsConfigurations struct {
		XMLName            xml.Name `xml:"tr2:GetAnalyticsConfigurations"`
		Xmlns              string   `xml:"xmlns:tr2,attr"`
		ConfigurationToken *string  `xml:"tr2:ConfigurationToken,omitempty"`
		ProfileToken       *string  `xml:"tr2:ProfileToken,omitempty"`
	}

	type GetAnalyticsConfigurationsResponse struct {
		XMLName        xml.Name `xml:"GetAnalyticsConfigurationsResponse"`
		Configurations []struct {
			Token    string `xml:"token,attr"`
			Name     string `xml:"Name"`
			UseCount int    `xml:"UseCount"`
		} `xml:"Configurations"`
	}

	req := GetAnalyticsConfigurations{
		Xmlns:              media2Namespace,
		ConfigurationToken: configToken,
		ProfileToken:       profileToken,
	}

	var resp GetAnalyticsConfigurationsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMedia2AnalyticsConfigurations failed: %w", err)
	}

	configs := make([]*VideoAnalyticsConfiguration, len(resp.Configurations))
	for i, cfg := range resp.Configurations {
		configs[i] = &VideoAnalyticsConfiguration{
			Token:    cfg.Token,
			Name:     cfg.Name,
			UseCount: cfg.UseCount,
		}
	}

	return configs, nil
}
