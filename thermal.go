package onvif

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/0x524a/onvif-go/internal/soap"
)

// Thermal service namespace.
const thermalNamespace = "http://www.onvif.org/ver10/thermal/wsdl"

// getThermalEndpoint returns the thermal endpoint, falling back to the device endpoint.
func (c *Client) getThermalEndpoint() string {
	if c.thermalEndpoint != "" {
		return c.thermalEndpoint
	}

	return c.endpoint
}

// newThermalSOAPClient creates a SOAP client for the thermal service.
func (c *Client) newThermalSOAPClient() *soap.Client {
	username, password := c.GetCredentials()

	return soap.NewClient(c.httpClient, username, password)
}

// ============================================================
// Capabilities
// ============================================================

// GetThermalServiceCapabilities returns the capabilities of the Thermal service.
func (c *Client) GetThermalServiceCapabilities(ctx context.Context) (*ThermalServiceCapabilities, error) {
	endpoint := c.getThermalEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tth:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:tth,attr"`
	}

	type capsXML struct {
		Radiometry *bool `xml:"Radiometry,attr"`
	}

	type Response struct {
		XMLName      xml.Name `xml:"GetServiceCapabilitiesResponse"`
		Capabilities capsXML  `xml:"Capabilities"`
	}

	req := Request{Xmlns: thermalNamespace}

	var resp Response

	if err := c.newThermalSOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetThermalServiceCapabilities failed: %w", err)
	}

	return &ThermalServiceCapabilities{
		Radiometry: resp.Capabilities.Radiometry,
	}, nil
}

// ============================================================
// Configuration
// ============================================================

// GetThermalConfigurations returns the thermal configuration for all thermal video sources.
func (c *Client) GetThermalConfigurations(ctx context.Context) ([]*ThermalConfigurations, error) {
	endpoint := c.getThermalEndpoint()

	type Request struct {
		XMLName xml.Name `xml:"tth:GetConfigurations"`
		Xmlns   string   `xml:"xmlns:tth,attr"`
	}

	type colorPaletteXML struct {
		Token string `xml:"token,attr"`
		Type  string `xml:"Type,attr"`
		Name  string `xml:"Name"`
	}

	type nucTableXML struct {
		Token           string   `xml:"token,attr"`
		LowTemperature  *float32 `xml:"LowTemperature,attr"`
		HighTemperature *float32 `xml:"HighTemperature,attr"`
		Name            string   `xml:"Name"`
	}

	type coolerXML struct {
		Enabled bool     `xml:"Enabled"`
		RunTime *float32 `xml:"RunTime"`
	}

	type configXML struct {
		ColorPalette colorPaletteXML `xml:"ColorPalette"`
		Polarity     string          `xml:"Polarity"`
		NUCTable     *nucTableXML    `xml:"NUCTable"`
		Cooler       *coolerXML      `xml:"Cooler"`
	}

	type configurationsXML struct {
		Token         string    `xml:"token,attr"`
		Configuration configXML `xml:"Configuration"`
	}

	type Response struct {
		XMLName        xml.Name             `xml:"GetConfigurationsResponse"`
		Configurations []*configurationsXML `xml:"Configurations"`
	}

	req := Request{Xmlns: thermalNamespace}

	var resp Response

	if err := c.newThermalSOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetThermalConfigurations failed: %w", err)
	}

	result := make([]*ThermalConfigurations, 0, len(resp.Configurations))

	for _, cfg := range resp.Configurations {
		tc := &ThermalConfigurations{
			Token: cfg.Token,
			Configuration: ThermalConfiguration{
				ColorPalette: ThermalColorPalette{
					Token: cfg.Configuration.ColorPalette.Token,
					Type:  cfg.Configuration.ColorPalette.Type,
					Name:  cfg.Configuration.ColorPalette.Name,
				},
				Polarity: ThermalPolarity(cfg.Configuration.Polarity),
			},
		}

		if cfg.Configuration.NUCTable != nil {
			tc.Configuration.NUCTable = &ThermalNUCTable{
				Token:           cfg.Configuration.NUCTable.Token,
				Name:            cfg.Configuration.NUCTable.Name,
				LowTemperature:  cfg.Configuration.NUCTable.LowTemperature,
				HighTemperature: cfg.Configuration.NUCTable.HighTemperature,
			}
		}

		if cfg.Configuration.Cooler != nil {
			tc.Configuration.Cooler = &ThermalCooler{
				Enabled: cfg.Configuration.Cooler.Enabled,
				RunTime: cfg.Configuration.Cooler.RunTime,
			}
		}

		result = append(result, tc)
	}

	return result, nil
}

// GetThermalConfiguration returns the thermal configuration for the specified video source.
func (c *Client) GetThermalConfiguration(ctx context.Context, videoSourceToken string) (*ThermalConfiguration, error) {
	endpoint := c.getThermalEndpoint()

	type Request struct {
		XMLName          xml.Name `xml:"tth:GetConfiguration"`
		Xmlns            string   `xml:"xmlns:tth,attr"`
		VideoSourceToken string   `xml:"tth:VideoSourceToken"`
	}

	type colorPaletteXML struct {
		Token string `xml:"token,attr"`
		Type  string `xml:"Type,attr"`
		Name  string `xml:"Name"`
	}

	type nucTableXML struct {
		Token           string   `xml:"token,attr"`
		LowTemperature  *float32 `xml:"LowTemperature,attr"`
		HighTemperature *float32 `xml:"HighTemperature,attr"`
		Name            string   `xml:"Name"`
	}

	type coolerXML struct {
		Enabled bool     `xml:"Enabled"`
		RunTime *float32 `xml:"RunTime"`
	}

	type Response struct {
		XMLName       xml.Name `xml:"GetConfigurationResponse"`
		Configuration struct {
			ColorPalette colorPaletteXML `xml:"ColorPalette"`
			Polarity     string          `xml:"Polarity"`
			NUCTable     *nucTableXML    `xml:"NUCTable"`
			Cooler       *coolerXML      `xml:"Cooler"`
		} `xml:"Configuration"`
	}

	req := Request{
		Xmlns:            thermalNamespace,
		VideoSourceToken: videoSourceToken,
	}

	var resp Response

	if err := c.newThermalSOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetThermalConfiguration failed: %w", err)
	}

	cfg := &ThermalConfiguration{
		ColorPalette: ThermalColorPalette{
			Token: resp.Configuration.ColorPalette.Token,
			Type:  resp.Configuration.ColorPalette.Type,
			Name:  resp.Configuration.ColorPalette.Name,
		},
		Polarity: ThermalPolarity(resp.Configuration.Polarity),
	}

	if resp.Configuration.NUCTable != nil {
		cfg.NUCTable = &ThermalNUCTable{
			Token:           resp.Configuration.NUCTable.Token,
			Name:            resp.Configuration.NUCTable.Name,
			LowTemperature:  resp.Configuration.NUCTable.LowTemperature,
			HighTemperature: resp.Configuration.NUCTable.HighTemperature,
		}
	}

	if resp.Configuration.Cooler != nil {
		cfg.Cooler = &ThermalCooler{
			Enabled: resp.Configuration.Cooler.Enabled,
			RunTime: resp.Configuration.Cooler.RunTime,
		}
	}

	return cfg, nil
}

// SetThermalConfiguration sets the thermal configuration for the specified video source.
func (c *Client) SetThermalConfiguration(ctx context.Context, videoSourceToken string, configuration ThermalConfiguration) error {
	endpoint := c.getThermalEndpoint()

	type colorPaletteXML struct {
		XMLName xml.Name `xml:"tth:ColorPalette"`
		Token   string   `xml:"token,attr"`
		Type    string   `xml:"Type,attr"`
		Name    string   `xml:"tth:Name"`
	}

	type nucTableXML struct {
		XMLName         xml.Name `xml:"tth:NUCTable"`
		Token           string   `xml:"token,attr"`
		LowTemperature  *float32 `xml:"LowTemperature,attr,omitempty"`
		HighTemperature *float32 `xml:"HighTemperature,attr,omitempty"`
		Name            string   `xml:"tth:Name"`
	}

	type coolerXML struct {
		XMLName xml.Name `xml:"tth:Cooler"`
		Enabled bool     `xml:"tth:Enabled"`
		RunTime *float32 `xml:"tth:RunTime,omitempty"`
	}

	type configXML struct {
		XMLName      xml.Name        `xml:"tth:Configuration"`
		ColorPalette colorPaletteXML
		Polarity     string       `xml:"tth:Polarity"`
		NUCTable     *nucTableXML `xml:",omitempty"`
		Cooler       *coolerXML   `xml:",omitempty"`
	}

	type Request struct {
		XMLName          xml.Name  `xml:"tth:SetConfiguration"`
		Xmlns            string    `xml:"xmlns:tth,attr"`
		VideoSourceToken string    `xml:"tth:VideoSourceToken"`
		Configuration    configXML
	}

	type Response struct {
		XMLName xml.Name `xml:"SetConfigurationResponse"`
	}

	cfg := configXML{
		ColorPalette: colorPaletteXML{
			Token: configuration.ColorPalette.Token,
			Type:  configuration.ColorPalette.Type,
			Name:  configuration.ColorPalette.Name,
		},
		Polarity: string(configuration.Polarity),
	}

	if configuration.NUCTable != nil {
		cfg.NUCTable = &nucTableXML{
			Token:           configuration.NUCTable.Token,
			LowTemperature:  configuration.NUCTable.LowTemperature,
			HighTemperature: configuration.NUCTable.HighTemperature,
			Name:            configuration.NUCTable.Name,
		}
	}

	if configuration.Cooler != nil {
		cfg.Cooler = &coolerXML{
			Enabled: configuration.Cooler.Enabled,
			RunTime: configuration.Cooler.RunTime,
		}
	}

	req := Request{
		Xmlns:            thermalNamespace,
		VideoSourceToken: videoSourceToken,
		Configuration:    cfg,
	}

	var resp Response

	if err := c.newThermalSOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetThermalConfiguration failed: %w", err)
	}

	return nil
}

// GetThermalConfigurationOptions returns the valid ranges for thermal configuration parameters.
func (c *Client) GetThermalConfigurationOptions(ctx context.Context, videoSourceToken string) (*ThermalConfigurationOptions, error) {
	endpoint := c.getThermalEndpoint()

	type Request struct {
		XMLName          xml.Name `xml:"tth:GetConfigurationOptions"`
		Xmlns            string   `xml:"xmlns:tth,attr"`
		VideoSourceToken string   `xml:"tth:VideoSourceToken"`
	}

	type colorPaletteXML struct {
		Token string `xml:"token,attr"`
		Type  string `xml:"Type,attr"`
		Name  string `xml:"Name"`
	}

	type nucTableXML struct {
		Token           string   `xml:"token,attr"`
		LowTemperature  *float32 `xml:"LowTemperature,attr"`
		HighTemperature *float32 `xml:"HighTemperature,attr"`
		Name            string   `xml:"Name"`
	}

	type Response struct {
		XMLName              xml.Name `xml:"GetConfigurationOptionsResponse"`
		ConfigurationOptions struct {
			ColorPalettes []*colorPaletteXML `xml:"ColorPalette"`
			NUCTables     []*nucTableXML     `xml:"NUCTable"`
			CoolerOptions *struct {
				Enabled *bool `xml:"Enabled"`
			} `xml:"CoolerOptions"`
		} `xml:"ConfigurationOptions"`
	}

	req := Request{
		Xmlns:            thermalNamespace,
		VideoSourceToken: videoSourceToken,
	}

	var resp Response

	if err := c.newThermalSOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetThermalConfigurationOptions failed: %w", err)
	}

	opts := &ThermalConfigurationOptions{}

	for _, cp := range resp.ConfigurationOptions.ColorPalettes {
		opts.ColorPalettes = append(opts.ColorPalettes, &ThermalColorPalette{
			Token: cp.Token,
			Type:  cp.Type,
			Name:  cp.Name,
		})
	}

	for _, nuc := range resp.ConfigurationOptions.NUCTables {
		opts.NUCTables = append(opts.NUCTables, &ThermalNUCTable{
			Token:           nuc.Token,
			Name:            nuc.Name,
			LowTemperature:  nuc.LowTemperature,
			HighTemperature: nuc.HighTemperature,
		})
	}

	if resp.ConfigurationOptions.CoolerOptions != nil {
		opts.CoolerOptions = &ThermalCoolerOptions{
			Enabled: resp.ConfigurationOptions.CoolerOptions.Enabled,
		}
	}

	return opts, nil
}

// ============================================================
// Radiometry
// ============================================================

// GetRadiometryConfiguration returns the radiometry configuration for the specified video source.
func (c *Client) GetRadiometryConfiguration(ctx context.Context, videoSourceToken string) (*RadiometryConfiguration, error) {
	endpoint := c.getThermalEndpoint()

	type Request struct {
		XMLName          xml.Name `xml:"tth:GetRadiometryConfiguration"`
		Xmlns            string   `xml:"xmlns:tth,attr"`
		VideoSourceToken string   `xml:"tth:VideoSourceToken"`
	}

	type globalParamsXML struct {
		ReflectedAmbientTemperature float32  `xml:"ReflectedAmbientTemperature"`
		Emissivity                  float32  `xml:"Emissivity"`
		DistanceToObject            float32  `xml:"DistanceToObject"`
		RelativeHumidity            *float32 `xml:"RelativeHumidity"`
		AtmosphericTemperature      *float32 `xml:"AtmosphericTemperature"`
		AtmosphericTransmittance    *float32 `xml:"AtmosphericTransmittance"`
		ExtOpticsTemperature        *float32 `xml:"ExtOpticsTemperature"`
		ExtOpticsTransmittance      *float32 `xml:"ExtOpticsTransmittance"`
	}

	type Response struct {
		XMLName       xml.Name `xml:"GetRadiometryConfigurationResponse"`
		Configuration struct {
			RadiometryGlobalParameters *globalParamsXML `xml:"RadiometryGlobalParameters"`
		} `xml:"Configuration"`
	}

	req := Request{
		Xmlns:            thermalNamespace,
		VideoSourceToken: videoSourceToken,
	}

	var resp Response

	if err := c.newThermalSOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRadiometryConfiguration failed: %w", err)
	}

	cfg := &RadiometryConfiguration{}

	if resp.Configuration.RadiometryGlobalParameters != nil {
		p := resp.Configuration.RadiometryGlobalParameters
		cfg.RadiometryGlobalParameters = &RadiometryGlobalParameters{
			ReflectedAmbientTemperature: p.ReflectedAmbientTemperature,
			Emissivity:                  p.Emissivity,
			DistanceToObject:            p.DistanceToObject,
			RelativeHumidity:            p.RelativeHumidity,
			AtmosphericTemperature:      p.AtmosphericTemperature,
			AtmosphericTransmittance:    p.AtmosphericTransmittance,
			ExtOpticsTemperature:        p.ExtOpticsTemperature,
			ExtOpticsTransmittance:      p.ExtOpticsTransmittance,
		}
	}

	return cfg, nil
}

// SetRadiometryConfiguration sets the radiometry configuration for the specified video source.
func (c *Client) SetRadiometryConfiguration(ctx context.Context, videoSourceToken string, configuration RadiometryConfiguration) error {
	endpoint := c.getThermalEndpoint()

	type globalParamsXML struct {
		XMLName                     xml.Name `xml:"tth:RadiometryGlobalParameters"`
		ReflectedAmbientTemperature float32  `xml:"tth:ReflectedAmbientTemperature"`
		Emissivity                  float32  `xml:"tth:Emissivity"`
		DistanceToObject            float32  `xml:"tth:DistanceToObject"`
		RelativeHumidity            *float32 `xml:"tth:RelativeHumidity,omitempty"`
		AtmosphericTemperature      *float32 `xml:"tth:AtmosphericTemperature,omitempty"`
		AtmosphericTransmittance    *float32 `xml:"tth:AtmosphericTransmittance,omitempty"`
		ExtOpticsTemperature        *float32 `xml:"tth:ExtOpticsTemperature,omitempty"`
		ExtOpticsTransmittance      *float32 `xml:"tth:ExtOpticsTransmittance,omitempty"`
	}

	type configXML struct {
		XMLName                    xml.Name         `xml:"tth:Configuration"`
		RadiometryGlobalParameters *globalParamsXML `xml:",omitempty"`
	}

	type Request struct {
		XMLName          xml.Name  `xml:"tth:SetRadiometryConfiguration"`
		Xmlns            string    `xml:"xmlns:tth,attr"`
		VideoSourceToken string    `xml:"tth:VideoSourceToken"`
		Configuration    configXML
	}

	type Response struct {
		XMLName xml.Name `xml:"SetRadiometryConfigurationResponse"`
	}

	cfg := configXML{}

	if configuration.RadiometryGlobalParameters != nil {
		p := configuration.RadiometryGlobalParameters
		cfg.RadiometryGlobalParameters = &globalParamsXML{
			ReflectedAmbientTemperature: p.ReflectedAmbientTemperature,
			Emissivity:                  p.Emissivity,
			DistanceToObject:            p.DistanceToObject,
			RelativeHumidity:            p.RelativeHumidity,
			AtmosphericTemperature:      p.AtmosphericTemperature,
			AtmosphericTransmittance:    p.AtmosphericTransmittance,
			ExtOpticsTemperature:        p.ExtOpticsTemperature,
			ExtOpticsTransmittance:      p.ExtOpticsTransmittance,
		}
	}

	req := Request{
		Xmlns:            thermalNamespace,
		VideoSourceToken: videoSourceToken,
		Configuration:    cfg,
	}

	var resp Response

	if err := c.newThermalSOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("SetRadiometryConfiguration failed: %w", err)
	}

	return nil
}

// GetRadiometryConfigurationOptions returns the valid ranges for radiometry configuration parameters.
func (c *Client) GetRadiometryConfigurationOptions(ctx context.Context, videoSourceToken string) (*RadiometryConfigurationOptions, error) {
	endpoint := c.getThermalEndpoint()

	type Request struct {
		XMLName          xml.Name `xml:"tth:GetRadiometryConfigurationOptions"`
		Xmlns            string   `xml:"xmlns:tth,attr"`
		VideoSourceToken string   `xml:"tth:VideoSourceToken"`
	}

	type floatRangeXML struct {
		Min float32 `xml:"Min"`
		Max float32 `xml:"Max"`
	}

	type globalParamOptsXML struct {
		ReflectedAmbientTemperature floatRangeXML  `xml:"ReflectedAmbientTemperature"`
		Emissivity                  floatRangeXML  `xml:"Emissivity"`
		DistanceToObject            floatRangeXML  `xml:"DistanceToObject"`
		RelativeHumidity            *floatRangeXML `xml:"RelativeHumidity"`
		AtmosphericTemperature      *floatRangeXML `xml:"AtmosphericTemperature"`
		AtmosphericTransmittance    *floatRangeXML `xml:"AtmosphericTransmittance"`
		ExtOpticsTemperature        *floatRangeXML `xml:"ExtOpticsTemperature"`
		ExtOpticsTransmittance      *floatRangeXML `xml:"ExtOpticsTransmittance"`
	}

	type Response struct {
		XMLName              xml.Name `xml:"GetRadiometryConfigurationOptionsResponse"`
		ConfigurationOptions struct {
			RadiometryGlobalParameterOptions *globalParamOptsXML `xml:"RadiometryGlobalParameterOptions"`
		} `xml:"ConfigurationOptions"`
	}

	req := Request{
		Xmlns:            thermalNamespace,
		VideoSourceToken: videoSourceToken,
	}

	var resp Response

	if err := c.newThermalSOAPClient().Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRadiometryConfigurationOptions failed: %w", err)
	}

	opts := &RadiometryConfigurationOptions{}

	if resp.ConfigurationOptions.RadiometryGlobalParameterOptions != nil {
		p := resp.ConfigurationOptions.RadiometryGlobalParameterOptions
		globalOpts := &RadiometryGlobalParameterOptions{
			ReflectedAmbientTemperature: RadiometryFloatRange{Min: p.ReflectedAmbientTemperature.Min, Max: p.ReflectedAmbientTemperature.Max},
			Emissivity:                  RadiometryFloatRange{Min: p.Emissivity.Min, Max: p.Emissivity.Max},
			DistanceToObject:            RadiometryFloatRange{Min: p.DistanceToObject.Min, Max: p.DistanceToObject.Max},
		}

		if p.RelativeHumidity != nil {
			fr := RadiometryFloatRange{Min: p.RelativeHumidity.Min, Max: p.RelativeHumidity.Max}
			globalOpts.RelativeHumidity = &fr
		}

		if p.AtmosphericTemperature != nil {
			fr := RadiometryFloatRange{Min: p.AtmosphericTemperature.Min, Max: p.AtmosphericTemperature.Max}
			globalOpts.AtmosphericTemperature = &fr
		}

		if p.AtmosphericTransmittance != nil {
			fr := RadiometryFloatRange{Min: p.AtmosphericTransmittance.Min, Max: p.AtmosphericTransmittance.Max}
			globalOpts.AtmosphericTransmittance = &fr
		}

		if p.ExtOpticsTemperature != nil {
			fr := RadiometryFloatRange{Min: p.ExtOpticsTemperature.Min, Max: p.ExtOpticsTemperature.Max}
			globalOpts.ExtOpticsTemperature = &fr
		}

		if p.ExtOpticsTransmittance != nil {
			fr := RadiometryFloatRange{Min: p.ExtOpticsTransmittance.Min, Max: p.ExtOpticsTransmittance.Max}
			globalOpts.ExtOpticsTransmittance = &fr
		}

		opts.RadiometryGlobalParameterOptions = globalOpts
	}

	return opts, nil
}
