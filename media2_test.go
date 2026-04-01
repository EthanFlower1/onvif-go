package onvif

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const soapFaultResponse = `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<s:Fault>
			<s:Code><s:Value>s:Receiver</s:Value></s:Code>
			<s:Reason><s:Text xml:lang="en">Internal error</s:Text></s:Reason>
		</s:Fault>
	</s:Body>
</s:Envelope>`

// TestGetMedia2ServiceCapabilities tests GetMedia2ServiceCapabilities operation.
func TestGetMedia2ServiceCapabilities(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		wantErr     bool
		checkResult func(t *testing.T, caps *Media2ServiceCapabilities)
	}{
		{
			name: "successful capabilities retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:GetServiceCapabilitiesResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl">
			<tr2:Capabilities SnapshotUri="true" Rotation="false" OSD="true" Mask="true" SourceMask="false" VideoSourceMode="true"/>
		</tr2:GetServiceCapabilitiesResponse>
	</s:Body>
</s:Envelope>`
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
			checkResult: func(t *testing.T, caps *Media2ServiceCapabilities) {
				t.Helper()
				if caps == nil {
					t.Fatal("Expected capabilities, got nil")
				}
				if !caps.SnapshotUri {
					t.Error("Expected SnapshotUri to be true")
				}
				if caps.Rotation {
					t.Error("Expected Rotation to be false")
				}
				if !caps.OSD {
					t.Error("Expected OSD to be true")
				}
				if !caps.Mask {
					t.Error("Expected Mask to be true")
				}
			},
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultResponse))
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
				t.Fatalf("NewClient() failed: %v", err)
			}

			caps, err := client.GetMedia2ServiceCapabilities(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMedia2ServiceCapabilities() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && tt.checkResult != nil {
				tt.checkResult(t, caps)
			}
		})
	}
}

// TestGetMedia2Profiles tests GetMedia2Profiles operation.
func TestGetMedia2Profiles(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		wantErr     bool
		checkResult func(t *testing.T, profiles []*Media2Profile)
	}{
		{
			name: "successful profiles retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:GetProfilesResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl">
			<tr2:Profiles token="prof1" fixed="true">
				<tr2:Name>Main Profile</tr2:Name>
				<tr2:Configurations>
					<tr2:VideoSource token="vs1">
						<tr2:Name>VS1</tr2:Name>
					</tr2:VideoSource>
					<tr2:VideoEncoder token="ve1">
						<tr2:Name>VE1</tr2:Name>
						<tr2:Encoding>H264</tr2:Encoding>
					</tr2:VideoEncoder>
				</tr2:Configurations>
			</tr2:Profiles>
		</tr2:GetProfilesResponse>
	</s:Body>
</s:Envelope>`
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
			checkResult: func(t *testing.T, profiles []*Media2Profile) {
				t.Helper()
				if len(profiles) != 1 {
					t.Fatalf("Expected 1 profile, got %d", len(profiles))
				}
				p := profiles[0]
				if p.Token != "prof1" {
					t.Errorf("Expected token 'prof1', got '%s'", p.Token)
				}
				if p.Name != "Main Profile" {
					t.Errorf("Expected name 'Main Profile', got '%s'", p.Name)
				}
				if !p.Fixed {
					t.Error("Expected Fixed to be true")
				}
				if p.Configurations == nil {
					t.Fatal("Expected configurations, got nil")
				}
				if p.Configurations.VideoSource == nil {
					t.Fatal("Expected VideoSource configuration, got nil")
				}
				if p.Configurations.VideoSource.Token != "vs1" {
					t.Errorf("Expected VideoSource token 'vs1', got '%s'", p.Configurations.VideoSource.Token)
				}
				if p.Configurations.VideoEncoder == nil {
					t.Fatal("Expected VideoEncoder configuration, got nil")
				}
				if p.Configurations.VideoEncoder.Encoding != "H264" {
					t.Errorf("Expected encoding 'H264', got '%s'", p.Configurations.VideoEncoder.Encoding)
				}
			},
		},
		{
			name: "profiles with optional token filter",
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Verify the request contains the token
				body := make([]byte, r.ContentLength)
				_, _ = r.Body.Read(body)
				if !strings.Contains(string(body), "prof1") {
					w.WriteHeader(http.StatusBadRequest)

					return
				}
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:GetProfilesResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl">
			<tr2:Profiles token="prof1" fixed="false">
				<tr2:Name>Profile 1</tr2:Name>
			</tr2:Profiles>
		</tr2:GetProfilesResponse>
	</s:Body>
</s:Envelope>`
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
			checkResult: func(t *testing.T, profiles []*Media2Profile) {
				t.Helper()
				if len(profiles) != 1 {
					t.Fatalf("Expected 1 profile, got %d", len(profiles))
				}
			},
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultResponse))
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
				t.Fatalf("NewClient() failed: %v", err)
			}

			var token *string
			if tt.name == "profiles with optional token filter" {
				tok := "prof1"
				token = &tok
			}

			profiles, err := client.GetMedia2Profiles(context.Background(), token, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMedia2Profiles() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && tt.checkResult != nil {
				tt.checkResult(t, profiles)
			}
		})
	}
}

// TestCreateMedia2Profile tests CreateMedia2Profile operation.
func TestCreateMedia2Profile(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
		wantToken string
	}{
		{
			name: "successful profile creation",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:CreateProfileResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl">
			<tr2:Token>newprof1</tr2:Token>
		</tr2:CreateProfileResponse>
	</s:Body>
</s:Envelope>`
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantToken: "newprof1",
		},
		{
			name: "successful creation with configurations",
			handler: func(w http.ResponseWriter, r *http.Request) {
				body := make([]byte, r.ContentLength)
				_, _ = r.Body.Read(body)
				// Verify configuration is in the request
				if !strings.Contains(string(body), "VideoSource") {
					w.WriteHeader(http.StatusBadRequest)

					return
				}
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:CreateProfileResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl">
			<tr2:Token>newprof2</tr2:Token>
		</tr2:CreateProfileResponse>
	</s:Body>
</s:Envelope>`
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantToken: "newprof2",
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultResponse))
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
				t.Fatalf("NewClient() failed: %v", err)
			}

			var configs []*Media2ConfigurationRef
			if tt.name == "successful creation with configurations" {
				configs = []*Media2ConfigurationRef{
					{Type: "VideoSource", Token: "vs1"},
				}
			}

			token, err := client.CreateMedia2Profile(context.Background(), "New Profile", configs)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMedia2Profile() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && token != tt.wantToken {
				t.Errorf("Expected token '%s', got '%s'", tt.wantToken, token)
			}
		})
	}
}

// TestDeleteMedia2Profile tests DeleteMedia2Profile operation.
func TestDeleteMedia2Profile(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful profile deletion",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:DeleteProfileResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl"/>
	</s:Body>
</s:Envelope>`
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultResponse))
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
				t.Fatalf("NewClient() failed: %v", err)
			}

			err = client.DeleteMedia2Profile(context.Background(), "prof1")
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteMedia2Profile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetMedia2StreamUri tests GetMedia2StreamUri operation.
func TestGetMedia2StreamUri(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
		wantUri string
	}{
		{
			name: "successful stream URI retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:GetStreamUriResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl">
			<tr2:Uri>rtsp://device/stream1</tr2:Uri>
		</tr2:GetStreamUriResponse>
	</s:Body>
</s:Envelope>`
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
			wantUri: "rtsp://device/stream1",
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultResponse))
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
				t.Fatalf("NewClient() failed: %v", err)
			}

			uri, err := client.GetMedia2StreamUri(context.Background(), "RtspUnicast", "prof1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMedia2StreamUri() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && uri != tt.wantUri {
				t.Errorf("Expected URI '%s', got '%s'", tt.wantUri, uri)
			}
		})
	}
}

// TestGetMedia2SnapshotUri tests GetMedia2SnapshotUri operation.
func TestGetMedia2SnapshotUri(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:GetSnapshotUriResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl">
			<tr2:Uri>http://device/snapshot?token=prof1</tr2:Uri>
		</tr2:GetSnapshotUriResponse>
	</s:Body>
</s:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	uri, err := client.GetMedia2SnapshotUri(context.Background(), "prof1")
	if err != nil {
		t.Fatalf("GetMedia2SnapshotUri() failed: %v", err)
	}

	if !strings.Contains(uri, "snapshot") {
		t.Errorf("Expected snapshot URI, got %s", uri)
	}
}

// TestAddMedia2Configuration tests AddMedia2Configuration operation.
func TestAddMedia2Configuration(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful configuration add",
			handler: func(w http.ResponseWriter, r *http.Request) {
				body := make([]byte, r.ContentLength)
				_, _ = r.Body.Read(body)
				bodyStr := string(body)
				if !strings.Contains(bodyStr, "VideoEncoder") || !strings.Contains(bodyStr, "ve1") {
					w.WriteHeader(http.StatusBadRequest)

					return
				}
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:AddConfigurationResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl"/>
	</s:Body>
</s:Envelope>`
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultResponse))
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
				t.Fatalf("NewClient() failed: %v", err)
			}

			config := &Media2ConfigurationRef{Type: "VideoEncoder", Token: "ve1"}
			err = client.AddMedia2Configuration(context.Background(), "prof1", config)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddMedia2Configuration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestRemoveMedia2Configuration tests RemoveMedia2Configuration operation.
func TestRemoveMedia2Configuration(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful configuration removal",
			handler: func(w http.ResponseWriter, r *http.Request) {
				body := make([]byte, r.ContentLength)
				_, _ = r.Body.Read(body)
				bodyStr := string(body)
				if !strings.Contains(bodyStr, "VideoEncoder") || !strings.Contains(bodyStr, "ve1") {
					w.WriteHeader(http.StatusBadRequest)

					return
				}
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:RemoveConfigurationResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl"/>
	</s:Body>
</s:Envelope>`
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultResponse))
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
				t.Fatalf("NewClient() failed: %v", err)
			}

			config := &Media2ConfigurationRef{Type: "VideoEncoder", Token: "ve1"}
			err = client.RemoveMedia2Configuration(context.Background(), "prof1", config)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveMedia2Configuration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetMedia2VideoSourceConfigurations tests GetMedia2VideoSourceConfigurations operation.
func TestGetMedia2VideoSourceConfigurations(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		wantErr     bool
		checkResult func(t *testing.T, configs []*VideoSourceConfiguration)
	}{
		{
			name: "successful video source configurations retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:GetVideoSourceConfigurationsResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl">
			<tr2:Configurations token="vs1">
				<tr2:Name>Main Video Source</tr2:Name>
				<tr2:UseCount>2</tr2:UseCount>
				<tr2:SourceToken>videosrc0</tr2:SourceToken>
				<tr2:Bounds x="0" y="0" width="1920" height="1080"/>
			</tr2:Configurations>
		</tr2:GetVideoSourceConfigurationsResponse>
	</s:Body>
</s:Envelope>`
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
			checkResult: func(t *testing.T, configs []*VideoSourceConfiguration) {
				t.Helper()
				if len(configs) != 1 {
					t.Fatalf("Expected 1 configuration, got %d", len(configs))
				}
				cfg := configs[0]
				if cfg.Token != "vs1" {
					t.Errorf("Expected token 'vs1', got '%s'", cfg.Token)
				}
				if cfg.Name != "Main Video Source" {
					t.Errorf("Expected name 'Main Video Source', got '%s'", cfg.Name)
				}
				if cfg.SourceToken != "videosrc0" {
					t.Errorf("Expected SourceToken 'videosrc0', got '%s'", cfg.SourceToken)
				}
				if cfg.Bounds == nil {
					t.Fatal("Expected Bounds, got nil")
				}
				if cfg.Bounds.Width != 1920 || cfg.Bounds.Height != 1080 {
					t.Errorf("Expected bounds 1920x1080, got %dx%d", cfg.Bounds.Width, cfg.Bounds.Height)
				}
			},
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultResponse))
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
				t.Fatalf("NewClient() failed: %v", err)
			}

			configs, err := client.GetMedia2VideoSourceConfigurations(context.Background(), nil, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMedia2VideoSourceConfigurations() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && tt.checkResult != nil {
				tt.checkResult(t, configs)
			}
		})
	}
}

// TestGetMedia2VideoEncoderConfigurations tests GetMedia2VideoEncoderConfigurations operation.
func TestGetMedia2VideoEncoderConfigurations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:GetVideoEncoderConfigurationsResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl">
			<tr2:Configurations token="ve1">
				<tr2:Name>H264 Encoder</tr2:Name>
				<tr2:UseCount>1</tr2:UseCount>
				<tr2:Encoding>H264</tr2:Encoding>
				<tr2:Quality>5.0</tr2:Quality>
				<tr2:Resolution>
					<tr2:Width>1920</tr2:Width>
					<tr2:Height>1080</tr2:Height>
				</tr2:Resolution>
				<tr2:RateControl>
					<tr2:FrameRateLimit>30</tr2:FrameRateLimit>
					<tr2:EncodingInterval>1</tr2:EncodingInterval>
					<tr2:BitrateLimit>4096</tr2:BitrateLimit>
				</tr2:RateControl>
			</tr2:Configurations>
		</tr2:GetVideoEncoderConfigurationsResponse>
	</s:Body>
</s:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	configs, err := client.GetMedia2VideoEncoderConfigurations(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetMedia2VideoEncoderConfigurations() failed: %v", err)
	}

	if len(configs) != 1 {
		t.Fatalf("Expected 1 configuration, got %d", len(configs))
	}

	cfg := configs[0]
	if cfg.Token != "ve1" {
		t.Errorf("Expected token 've1', got '%s'", cfg.Token)
	}
	if cfg.Encoding != "H264" {
		t.Errorf("Expected encoding 'H264', got '%s'", cfg.Encoding)
	}
	if cfg.Resolution == nil {
		t.Fatal("Expected Resolution, got nil")
	}
	if cfg.Resolution.Width != 1920 || cfg.Resolution.Height != 1080 {
		t.Errorf("Expected 1920x1080, got %dx%d", cfg.Resolution.Width, cfg.Resolution.Height)
	}
	if cfg.RateControl == nil {
		t.Fatal("Expected RateControl, got nil")
	}
	if cfg.RateControl.FrameRateLimit != 30 {
		t.Errorf("Expected FrameRateLimit 30, got %d", cfg.RateControl.FrameRateLimit)
	}
}

// TestGetMedia2AudioSourceConfigurations tests GetMedia2AudioSourceConfigurations operation.
func TestGetMedia2AudioSourceConfigurations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:GetAudioSourceConfigurationsResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl">
			<tr2:Configurations token="as1">
				<tr2:Name>Audio Source</tr2:Name>
				<tr2:UseCount>1</tr2:UseCount>
				<tr2:SourceToken>audiosrc0</tr2:SourceToken>
			</tr2:Configurations>
		</tr2:GetAudioSourceConfigurationsResponse>
	</s:Body>
</s:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	configs, err := client.GetMedia2AudioSourceConfigurations(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetMedia2AudioSourceConfigurations() failed: %v", err)
	}

	if len(configs) != 1 {
		t.Fatalf("Expected 1 configuration, got %d", len(configs))
	}

	cfg := configs[0]
	if cfg.Token != "as1" {
		t.Errorf("Expected token 'as1', got '%s'", cfg.Token)
	}
	if cfg.SourceToken != "audiosrc0" {
		t.Errorf("Expected SourceToken 'audiosrc0', got '%s'", cfg.SourceToken)
	}
}

// TestGetMedia2AudioEncoderConfigurations tests GetMedia2AudioEncoderConfigurations operation.
func TestGetMedia2AudioEncoderConfigurations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:GetAudioEncoderConfigurationsResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl">
			<tr2:Configurations token="ae1">
				<tr2:Name>AAC Encoder</tr2:Name>
				<tr2:UseCount>1</tr2:UseCount>
				<tr2:Encoding>AAC</tr2:Encoding>
				<tr2:Bitrate>128</tr2:Bitrate>
				<tr2:SampleRate>44100</tr2:SampleRate>
			</tr2:Configurations>
		</tr2:GetAudioEncoderConfigurationsResponse>
	</s:Body>
</s:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	configs, err := client.GetMedia2AudioEncoderConfigurations(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetMedia2AudioEncoderConfigurations() failed: %v", err)
	}

	if len(configs) != 1 {
		t.Fatalf("Expected 1 configuration, got %d", len(configs))
	}

	cfg := configs[0]
	if cfg.Token != "ae1" {
		t.Errorf("Expected token 'ae1', got '%s'", cfg.Token)
	}
	if cfg.Encoding != "AAC" {
		t.Errorf("Expected encoding 'AAC', got '%s'", cfg.Encoding)
	}
	if cfg.Bitrate != 128 {
		t.Errorf("Expected bitrate 128, got %d", cfg.Bitrate)
	}
	if cfg.SampleRate != 44100 {
		t.Errorf("Expected sample rate 44100, got %d", cfg.SampleRate)
	}
}

// TestStartMedia2MulticastStreaming tests StartMedia2MulticastStreaming operation.
func TestStartMedia2MulticastStreaming(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:StartMulticastStreamingResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl"/>
	</s:Body>
</s:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	err = client.StartMedia2MulticastStreaming(context.Background(), "prof1")
	if err != nil {
		t.Fatalf("StartMedia2MulticastStreaming() failed: %v", err)
	}
}

// TestStopMedia2MulticastStreaming tests StopMedia2MulticastStreaming operation.
func TestStopMedia2MulticastStreaming(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:StopMulticastStreamingResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl"/>
	</s:Body>
</s:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	err = client.StopMedia2MulticastStreaming(context.Background(), "prof1")
	if err != nil {
		t.Fatalf("StopMedia2MulticastStreaming() failed: %v", err)
	}
}

// TestSetMedia2SynchronizationPoint tests SetMedia2SynchronizationPoint operation.
func TestSetMedia2SynchronizationPoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
	<s:Body>
		<tr2:SetSynchronizationPointResponse xmlns:tr2="http://www.onvif.org/ver20/media/wsdl"/>
	</s:Body>
</s:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	err = client.SetMedia2SynchronizationPoint(context.Background(), "prof1")
	if err != nil {
		t.Fatalf("SetMedia2SynchronizationPoint() failed: %v", err)
	}
}
