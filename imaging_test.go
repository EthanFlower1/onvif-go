package onvif

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetImagingServiceCapabilities(t *testing.T) {
	tests := []struct {
		name               string
		handler            http.HandlerFunc
		wantErr            bool
		wantStabilization  bool
		wantPresets        bool
	}{
		{
			name: "successful capabilities retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<timg:GetServiceCapabilitiesResponse xmlns:timg="http://www.onvif.org/ver20/imaging/wsdl">
							<timg:Capabilities ImageStabilization="true" Presets="true"/>
						</timg:GetServiceCapabilitiesResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:           false,
			wantStabilization: true,
			wantPresets:       true,
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Internal error</s:Text></s:Reason>
						</s:Fault>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(response))
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

			caps, err := client.GetImagingServiceCapabilities(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetImagingServiceCapabilities() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if caps == nil {
					t.Fatal("Expected capabilities, got nil")
				}

				if caps.ImageStabilization != tt.wantStabilization {
					t.Errorf("ImageStabilization = %v, want %v", caps.ImageStabilization, tt.wantStabilization)
				}

				if caps.Presets != tt.wantPresets {
					t.Errorf("Presets = %v, want %v", caps.Presets, tt.wantPresets)
				}
			}
		})
	}
}

func TestGetImagingPresets(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantErr    bool
		wantCount  int
		checkFirst func(t *testing.T, preset *ImagingPreset)
	}{
		{
			name: "successful presets retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<timg:GetPresetsResponse xmlns:timg="http://www.onvif.org/ver20/imaging/wsdl">
							<timg:Preset token="preset1" type="Custom">
								<timg:Name>My Preset</timg:Name>
							</timg:Preset>
							<timg:Preset token="preset2" type="Default">
								<timg:Name>Default Preset</timg:Name>
							</timg:Preset>
						</timg:GetPresetsResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantCount: 2,
			checkFirst: func(t *testing.T, preset *ImagingPreset) {
				t.Helper()

				if preset.Token != "preset1" {
					t.Errorf("Token = %q, want %q", preset.Token, "preset1")
				}

				if preset.Type != "Custom" {
					t.Errorf("Type = %q, want %q", preset.Type, "Custom")
				}

				if preset.Name != "My Preset" {
					t.Errorf("Name = %q, want %q", preset.Name, "My Preset")
				}
			},
		},
		{
			name: "empty presets list",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<timg:GetPresetsResponse xmlns:timg="http://www.onvif.org/ver20/imaging/wsdl">
						</timg:GetPresetsResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Service not available</s:Text></s:Reason>
						</s:Fault>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(response))
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

			presets, err := client.GetImagingPresets(context.Background(), "VideoSource_1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetImagingPresets() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if len(presets) != tt.wantCount {
					t.Errorf("len(presets) = %d, want %d", len(presets), tt.wantCount)
				}

				if tt.checkFirst != nil && len(presets) > 0 {
					tt.checkFirst(t, presets[0])
				}
			}
		})
	}
}

func TestGetCurrentImagingPreset(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		wantErr   bool
		wantNil   bool
		wantToken string
		wantName  string
	}{
		{
			name: "successful current preset retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<timg:GetCurrentPresetResponse xmlns:timg="http://www.onvif.org/ver20/imaging/wsdl">
							<timg:Preset token="preset1" type="Custom">
								<timg:Name>Current</timg:Name>
							</timg:Preset>
						</timg:GetCurrentPresetResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantNil:   false,
			wantToken: "preset1",
			wantName:  "Current",
		},
		{
			name: "no current preset set",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<timg:GetCurrentPresetResponse xmlns:timg="http://www.onvif.org/ver20/imaging/wsdl">
						</timg:GetCurrentPresetResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
			wantNil: true,
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Internal error</s:Text></s:Reason>
						</s:Fault>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(response))
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

			preset, err := client.GetCurrentImagingPreset(context.Background(), "VideoSource_1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentImagingPreset() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if tt.wantNil && preset != nil {
					t.Errorf("Expected nil preset, got %+v", preset)

					return
				}

				if !tt.wantNil {
					if preset == nil {
						t.Fatal("Expected preset, got nil")
					}

					if preset.Token != tt.wantToken {
						t.Errorf("Token = %q, want %q", preset.Token, tt.wantToken)
					}

					if preset.Name != tt.wantName {
						t.Errorf("Name = %q, want %q", preset.Name, tt.wantName)
					}
				}
			}
		})
	}
}

func TestSetCurrentImagingPreset(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful set current preset",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<timg:SetCurrentPresetResponse xmlns:timg="http://www.onvif.org/ver20/imaging/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Invalid preset token</s:Text></s:Reason>
						</s:Fault>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(response))
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

			err = client.SetCurrentImagingPreset(context.Background(), "VideoSource_1", "preset1")
			if (err != nil) != tt.wantErr {
				t.Errorf("SetCurrentImagingPreset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
