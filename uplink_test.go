package onvif

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUplinkServiceCapabilities(t *testing.T) {
	tests := []struct {
		name               string
		handler            http.HandlerFunc
		wantErr            bool
		wantMaxUplinks     *int
		wantProtocols      string
		wantAuthModes      string
		wantStreaming       *bool
	}{
		{
			name: "successful capabilities with all fields",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tup:GetServiceCapabilitiesResponse xmlns:tup="http://www.onvif.org/ver10/uplink/wsdl">
							<tup:Capabilities MaxUplinks="4" Protocols="https wss" AuthorizationModes="mTLS AccessToken" StreamingOverUplink="true"/>
						</tup:GetServiceCapabilitiesResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:        false,
			wantMaxUplinks: func() *int { v := 4; return &v }(),
			wantProtocols:  "https wss",
			wantAuthModes:  "mTLS AccessToken",
			wantStreaming:  func() *bool { v := true; return &v }(),
		},
		{
			name: "successful capabilities with minimal fields",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tup:GetServiceCapabilitiesResponse xmlns:tup="http://www.onvif.org/ver10/uplink/wsdl">
							<tup:Capabilities Protocols="https"/>
						</tup:GetServiceCapabilitiesResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:       false,
			wantProtocols: "https",
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

			caps, err := client.GetUplinkServiceCapabilities(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUplinkServiceCapabilities() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.wantErr {
				return
			}

			if caps == nil {
				t.Fatal("Expected capabilities, got nil")
			}

			if tt.wantMaxUplinks != nil {
				if caps.MaxUplinks == nil || *caps.MaxUplinks != *tt.wantMaxUplinks {
					t.Errorf("MaxUplinks = %v, want %v", caps.MaxUplinks, tt.wantMaxUplinks)
				}
			}

			if caps.Protocols != tt.wantProtocols {
				t.Errorf("Protocols = %q, want %q", caps.Protocols, tt.wantProtocols)
			}

			if caps.AuthorizationModes != tt.wantAuthModes {
				t.Errorf("AuthorizationModes = %q, want %q", caps.AuthorizationModes, tt.wantAuthModes)
			}

			if tt.wantStreaming != nil {
				if caps.StreamingOverUplink == nil || *caps.StreamingOverUplink != *tt.wantStreaming {
					t.Errorf("StreamingOverUplink = %v, want %v", caps.StreamingOverUplink, tt.wantStreaming)
				}
			}
		})
	}
}

func TestGetUplinks(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		wantErr     bool
		wantCount   int
		checkResult func(t *testing.T, configs []*UplinkConfiguration)
	}{
		{
			name: "successful retrieval with two uplinks",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tup:GetUplinksResponse xmlns:tup="http://www.onvif.org/ver10/uplink/wsdl">
							<tup:Configuration>
								<tup:RemoteAddress>https://cloud.example.com/uplink1</tup:RemoteAddress>
								<tup:UserLevel>Administrator</tup:UserLevel>
								<tup:Status>Connected</tup:Status>
							</tup:Configuration>
							<tup:Configuration>
								<tup:RemoteAddress>wss://cloud.example.com/uplink2</tup:RemoteAddress>
								<tup:CertificateID>cert-001</tup:CertificateID>
								<tup:UserLevel>Operator</tup:UserLevel>
								<tup:Status>Offline</tup:Status>
							</tup:Configuration>
						</tup:GetUplinksResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantCount: 2,
			checkResult: func(t *testing.T, configs []*UplinkConfiguration) {
				t.Helper()

				if configs[0].RemoteAddress != "https://cloud.example.com/uplink1" {
					t.Errorf("configs[0].RemoteAddress = %q, want %q", configs[0].RemoteAddress, "https://cloud.example.com/uplink1")
				}

				if configs[0].Status == nil || *configs[0].Status != "Connected" {
					t.Errorf("configs[0].Status = %v, want Connected", configs[0].Status)
				}

				if configs[1].CertificateID == nil || *configs[1].CertificateID != "cert-001" {
					t.Errorf("configs[1].CertificateID = %v, want cert-001", configs[1].CertificateID)
				}
			},
		},
		{
			name: "empty uplinks list",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tup:GetUplinksResponse xmlns:tup="http://www.onvif.org/ver10/uplink/wsdl"/>
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
							<s:Reason><s:Text xml:lang="en">Service unavailable</s:Text></s:Reason>
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

			configs, err := client.GetUplinks(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUplinks() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.wantErr {
				return
			}

			if len(configs) != tt.wantCount {
				t.Errorf("len(configs) = %d, want %d", len(configs), tt.wantCount)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, configs)
			}
		})
	}
}

func TestSetUplink(t *testing.T) {
	tests := []struct {
		name    string
		config  UplinkConfiguration
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful set uplink",
			config: UplinkConfiguration{
				RemoteAddress: "https://cloud.example.com/uplink1",
				UserLevel:     "Administrator",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tup:SetUplinkResponse xmlns:tup="http://www.onvif.org/ver10/uplink/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name: "successful set uplink with optional fields",
			config: UplinkConfiguration{
				RemoteAddress: "wss://cloud.example.com/uplink2",
				CertificateID: func() *string { s := "cert-001"; return &s }(),
				UserLevel:     "Operator",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tup:SetUplinkResponse xmlns:tup="http://www.onvif.org/ver10/uplink/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name: "SOAP fault response",
			config: UplinkConfiguration{
				RemoteAddress: "https://invalid.example.com",
				UserLevel:     "Administrator",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Sender</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Invalid configuration</s:Text></s:Reason>
						</s:Fault>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusBadRequest)
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

			err = client.SetUplink(context.Background(), tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetUplink() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteUplink(t *testing.T) {
	tests := []struct {
		name          string
		remoteAddress string
		handler       http.HandlerFunc
		wantErr       bool
	}{
		{
			name:          "successful delete",
			remoteAddress: "https://cloud.example.com/uplink1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tup:DeleteUplinkResponse xmlns:tup="http://www.onvif.org/ver10/uplink/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:          "SOAP fault - uplink not found",
			remoteAddress: "https://nonexistent.example.com/uplink",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Sender</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">No such uplink configuration</s:Text></s:Reason>
						</s:Fault>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(response))
			},
			wantErr: true,
		},
		{
			name:          "SOAP fault - internal error",
			remoteAddress: "https://cloud.example.com/uplink2",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Internal device error</s:Text></s:Reason>
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

			err = client.DeleteUplink(context.Background(), tt.remoteAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUplink() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
