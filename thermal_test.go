package onvif

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testThermalXMLHeader = `<?xml version="1.0" encoding="UTF-8"?>`

const soapFaultThermal = testThermalXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <SOAP-ENV:Fault>
      <SOAP-ENV:Code>
        <SOAP-ENV:Value>SOAP-ENV:Sender</SOAP-ENV:Value>
      </SOAP-ENV:Code>
      <SOAP-ENV:Reason>
        <SOAP-ENV:Text xml:lang="en">Invalid argument</SOAP-ENV:Text>
      </SOAP-ENV:Reason>
    </SOAP-ENV:Fault>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`


func TestGetThermalServiceCapabilities(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		wantErr        bool
		wantRadiometry *bool
	}{
		{
			name: "successful capabilities retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				radiometry := true
				_ = radiometry
				response := testThermalXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tth:GetServiceCapabilitiesResponse xmlns:tth="http://www.onvif.org/ver10/thermal/wsdl">
      <tth:Capabilities Radiometry="true"/>
    </tth:GetServiceCapabilitiesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:        false,
			wantRadiometry: boolPtr(true),
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultThermal))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			caps, err := client.GetThermalServiceCapabilities(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThermalServiceCapabilities() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if caps == nil {
					t.Fatal("Expected capabilities, got nil")
				}

				if tt.wantRadiometry != nil {
					if caps.Radiometry == nil {
						t.Error("Expected Radiometry to be set, got nil")
					} else if *caps.Radiometry != *tt.wantRadiometry {
						t.Errorf("Radiometry = %v, want %v", *caps.Radiometry, *tt.wantRadiometry)
					}
				}
			}
		})
	}
}

func TestGetThermalConfigurations(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		wantErr   bool
		wantCount int
	}{
		{
			name: "successful configurations retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testThermalXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tth:GetConfigurationsResponse xmlns:tth="http://www.onvif.org/ver10/thermal/wsdl">
      <tth:Configurations token="VideoSource_1">
        <tth:Configuration>
          <tth:ColorPalette token="palette_iron" Type="Iron">
            <tth:Name>Iron</tth:Name>
          </tth:ColorPalette>
          <tth:Polarity>WhiteHot</tth:Polarity>
        </tth:Configuration>
      </tth:Configurations>
      <tth:Configurations token="VideoSource_2">
        <tth:Configuration>
          <tth:ColorPalette token="palette_gray" Type="Grayscale">
            <tth:Name>Grayscale</tth:Name>
          </tth:ColorPalette>
          <tth:Polarity>BlackHot</tth:Polarity>
        </tth:Configuration>
      </tth:Configurations>
    </tth:GetConfigurationsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultThermal))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			cfgs, err := client.GetThermalConfigurations(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThermalConfigurations() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if len(cfgs) != tt.wantCount {
					t.Errorf("Expected %d configurations, got %d", tt.wantCount, len(cfgs))
				}
			}
		})
	}
}

func TestGetThermalConfiguration(t *testing.T) {
	tests := []struct {
		name             string
		videoSourceToken string
		handler          http.HandlerFunc
		wantErr          bool
		wantPolarity     ThermalPolarity
		wantPaletteToken string
		wantCooler       bool
	}{
		{
			name:             "successful configuration retrieval with cooler",
			videoSourceToken: "VideoSource_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testThermalXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tth:GetConfigurationResponse xmlns:tth="http://www.onvif.org/ver10/thermal/wsdl">
      <tth:Configuration>
        <tth:ColorPalette token="palette_rainbow" Type="Rainbow">
          <tth:Name>Rainbow</tth:Name>
        </tth:ColorPalette>
        <tth:Polarity>BlackHot</tth:Polarity>
        <tth:Cooler>
          <tth:Enabled>true</tth:Enabled>
          <tth:RunTime>1024.5</tth:RunTime>
        </tth:Cooler>
      </tth:Configuration>
    </tth:GetConfigurationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:          false,
			wantPolarity:     ThermalPolarityBlackHot,
			wantPaletteToken: "palette_rainbow",
			wantCooler:       true,
		},
		{
			name:             "SOAP fault response",
			videoSourceToken: "VideoSource_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultThermal))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			cfg, err := client.GetThermalConfiguration(context.Background(), tt.videoSourceToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThermalConfiguration() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if cfg == nil {
					t.Fatal("Expected configuration, got nil")
				}

				if cfg.Polarity != tt.wantPolarity {
					t.Errorf("Polarity = %v, want %v", cfg.Polarity, tt.wantPolarity)
				}

				if cfg.ColorPalette.Token != tt.wantPaletteToken {
					t.Errorf("ColorPalette.Token = %v, want %v", cfg.ColorPalette.Token, tt.wantPaletteToken)
				}

				if tt.wantCooler && cfg.Cooler == nil {
					t.Error("Expected Cooler to be set, got nil")
				}
			}
		})
	}
}

func TestSetThermalConfiguration(t *testing.T) {
	tests := []struct {
		name             string
		videoSourceToken string
		configuration    ThermalConfiguration
		handler          http.HandlerFunc
		wantErr          bool
	}{
		{
			name:             "successful configuration set",
			videoSourceToken: "VideoSource_1",
			configuration: ThermalConfiguration{
				ColorPalette: ThermalColorPalette{Token: "palette_iron", Type: "Iron", Name: "Iron"},
				Polarity:     ThermalPolarityWhiteHot,
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testThermalXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tth:SetConfigurationResponse xmlns:tth="http://www.onvif.org/ver10/thermal/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:             "SOAP fault response",
			videoSourceToken: "VideoSource_1",
			configuration:    ThermalConfiguration{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultThermal))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			err = client.SetThermalConfiguration(context.Background(), tt.videoSourceToken, tt.configuration)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetThermalConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetThermalConfigurationOptions(t *testing.T) {
	tests := []struct {
		name             string
		videoSourceToken string
		handler          http.HandlerFunc
		wantErr          bool
		wantPaletteCount int
		wantNUCCount     int
		wantCoolerOpts   bool
	}{
		{
			name:             "successful options retrieval",
			videoSourceToken: "VideoSource_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testThermalXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tth:GetConfigurationOptionsResponse xmlns:tth="http://www.onvif.org/ver10/thermal/wsdl">
      <tth:ConfigurationOptions>
        <tth:ColorPalette token="palette_iron" Type="Iron">
          <tth:Name>Iron</tth:Name>
        </tth:ColorPalette>
        <tth:ColorPalette token="palette_rainbow" Type="Rainbow">
          <tth:Name>Rainbow</tth:Name>
        </tth:ColorPalette>
        <tth:NUCTable token="nuc_01" LowTemperature="233.15" HighTemperature="373.15">
          <tth:Name>Standard NUC</tth:Name>
        </tth:NUCTable>
        <tth:CoolerOptions>
          <tth:Enabled>true</tth:Enabled>
        </tth:CoolerOptions>
      </tth:ConfigurationOptions>
    </tth:GetConfigurationOptionsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:          false,
			wantPaletteCount: 2,
			wantNUCCount:     1,
			wantCoolerOpts:   true,
		},
		{
			name:             "SOAP fault response",
			videoSourceToken: "VideoSource_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultThermal))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			opts, err := client.GetThermalConfigurationOptions(context.Background(), tt.videoSourceToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThermalConfigurationOptions() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if opts == nil {
					t.Fatal("Expected options, got nil")
				}

				if len(opts.ColorPalettes) != tt.wantPaletteCount {
					t.Errorf("ColorPalettes count = %d, want %d", len(opts.ColorPalettes), tt.wantPaletteCount)
				}

				if len(opts.NUCTables) != tt.wantNUCCount {
					t.Errorf("NUCTables count = %d, want %d", len(opts.NUCTables), tt.wantNUCCount)
				}

				if tt.wantCoolerOpts && opts.CoolerOptions == nil {
					t.Error("Expected CoolerOptions to be set, got nil")
				}
			}
		})
	}
}

func TestGetRadiometryConfiguration(t *testing.T) {
	tests := []struct {
		name             string
		videoSourceToken string
		handler          http.HandlerFunc
		wantErr          bool
		wantEmissivity   float32
	}{
		{
			name:             "successful radiometry configuration retrieval",
			videoSourceToken: "VideoSource_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testThermalXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tth:GetRadiometryConfigurationResponse xmlns:tth="http://www.onvif.org/ver10/thermal/wsdl">
      <tth:Configuration>
        <tth:RadiometryGlobalParameters>
          <tth:ReflectedAmbientTemperature>293.15</tth:ReflectedAmbientTemperature>
          <tth:Emissivity>0.95</tth:Emissivity>
          <tth:DistanceToObject>10.0</tth:DistanceToObject>
        </tth:RadiometryGlobalParameters>
      </tth:Configuration>
    </tth:GetRadiometryConfigurationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:        false,
			wantEmissivity: 0.95,
		},
		{
			name:             "SOAP fault response",
			videoSourceToken: "VideoSource_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultThermal))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			cfg, err := client.GetRadiometryConfiguration(context.Background(), tt.videoSourceToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRadiometryConfiguration() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if cfg == nil {
					t.Fatal("Expected configuration, got nil")
				}

				if cfg.RadiometryGlobalParameters == nil {
					t.Fatal("Expected RadiometryGlobalParameters, got nil")
				}

				if cfg.RadiometryGlobalParameters.Emissivity != tt.wantEmissivity {
					t.Errorf("Emissivity = %v, want %v", cfg.RadiometryGlobalParameters.Emissivity, tt.wantEmissivity)
				}
			}
		})
	}
}

func TestSetRadiometryConfiguration(t *testing.T) {
	relHumidity := float32(45.0)

	tests := []struct {
		name             string
		videoSourceToken string
		configuration    RadiometryConfiguration
		handler          http.HandlerFunc
		wantErr          bool
	}{
		{
			name:             "successful radiometry configuration set",
			videoSourceToken: "VideoSource_1",
			configuration: RadiometryConfiguration{
				RadiometryGlobalParameters: &RadiometryGlobalParameters{
					ReflectedAmbientTemperature: 293.15,
					Emissivity:                  0.95,
					DistanceToObject:            10.0,
					RelativeHumidity:            &relHumidity,
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testThermalXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tth:SetRadiometryConfigurationResponse xmlns:tth="http://www.onvif.org/ver10/thermal/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:             "SOAP fault response",
			videoSourceToken: "VideoSource_1",
			configuration:    RadiometryConfiguration{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultThermal))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			err = client.SetRadiometryConfiguration(context.Background(), tt.videoSourceToken, tt.configuration)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetRadiometryConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetRadiometryConfigurationOptions(t *testing.T) {
	tests := []struct {
		name             string
		videoSourceToken string
		handler          http.HandlerFunc
		wantErr          bool
		wantHasGlobal    bool
	}{
		{
			name:             "successful radiometry options retrieval",
			videoSourceToken: "VideoSource_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testThermalXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tth:GetRadiometryConfigurationOptionsResponse xmlns:tth="http://www.onvif.org/ver10/thermal/wsdl">
      <tth:ConfigurationOptions>
        <tth:RadiometryGlobalParameterOptions>
          <tth:ReflectedAmbientTemperature>
            <tth:Min>233.15</tth:Min>
            <tth:Max>373.15</tth:Max>
          </tth:ReflectedAmbientTemperature>
          <tth:Emissivity>
            <tth:Min>0.1</tth:Min>
            <tth:Max>1.0</tth:Max>
          </tth:Emissivity>
          <tth:DistanceToObject>
            <tth:Min>0.5</tth:Min>
            <tth:Max>1000.0</tth:Max>
          </tth:DistanceToObject>
        </tth:RadiometryGlobalParameterOptions>
      </tth:ConfigurationOptions>
    </tth:GetRadiometryConfigurationOptionsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:       false,
			wantHasGlobal: true,
		},
		{
			name:             "SOAP fault response",
			videoSourceToken: "VideoSource_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultThermal))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			opts, err := client.GetRadiometryConfigurationOptions(context.Background(), tt.videoSourceToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRadiometryConfigurationOptions() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if opts == nil {
					t.Fatal("Expected options, got nil")
				}

				if tt.wantHasGlobal && opts.RadiometryGlobalParameterOptions == nil {
					t.Error("Expected RadiometryGlobalParameterOptions to be set, got nil")
				}
			}
		})
	}
}

// boolPtr returns a pointer to the given bool value.
func boolPtr(b bool) *bool {
	return &b
}
