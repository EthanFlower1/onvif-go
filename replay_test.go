package onvif

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetReplayServiceCapabilities(t *testing.T) {
	tests := []struct {
		name                string
		handler             http.HandlerFunc
		wantErr             bool
		wantReversePlayback bool
		wantRTPRTSP_TCP     bool
		wantSessionMin      string
		wantSessionMax      string
	}{
		{
			name: "successful capabilities retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trp:GetServiceCapabilitiesResponse xmlns:trp="http://www.onvif.org/ver10/replay/wsdl">
							<trp:Capabilities ReversePlayback="false" RTP_RTSP_TCP="true">
								<trp:SessionTimeoutRange>
									<trp:Min>PT60S</trp:Min>
									<trp:Max>PT600S</trp:Max>
								</trp:SessionTimeoutRange>
							</trp:Capabilities>
						</trp:GetServiceCapabilitiesResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:             false,
			wantReversePlayback: false,
			wantRTPRTSP_TCP:     true,
			wantSessionMin:      "PT60S",
			wantSessionMax:      "PT600S",
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

			caps, err := client.GetReplayServiceCapabilities(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReplayServiceCapabilities() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if caps == nil {
					t.Fatal("Expected capabilities, got nil")
				}

				if caps.ReversePlayback != tt.wantReversePlayback {
					t.Errorf("ReversePlayback = %v, want %v", caps.ReversePlayback, tt.wantReversePlayback)
				}

				if caps.RTPRTSP_TCP != tt.wantRTPRTSP_TCP {
					t.Errorf("RTPRTSP_TCP = %v, want %v", caps.RTPRTSP_TCP, tt.wantRTPRTSP_TCP)
				}

				if caps.SessionTimeoutRange == nil {
					t.Fatal("Expected SessionTimeoutRange, got nil")
				}

				if caps.SessionTimeoutRange.Min != tt.wantSessionMin {
					t.Errorf("SessionTimeoutRange.Min = %v, want %v", caps.SessionTimeoutRange.Min, tt.wantSessionMin)
				}

				if caps.SessionTimeoutRange.Max != tt.wantSessionMax {
					t.Errorf("SessionTimeoutRange.Max = %v, want %v", caps.SessionTimeoutRange.Max, tt.wantSessionMax)
				}
			}
		})
	}
}

func TestGetReplayUri(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		recordingToken string
		stream         string
		protocol       string
		wantErr        bool
		wantURI        string
	}{
		{
			name: "successful URI retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trp:GetReplayUriResponse xmlns:trp="http://www.onvif.org/ver10/replay/wsdl">
							<trp:Uri>rtsp://device/replay?token=REC1</trp:Uri>
						</trp:GetReplayUriResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			recordingToken: "REC1",
			stream:         "RTP-Unicast",
			protocol:       "RTSP",
			wantErr:        false,
			wantURI:        "rtsp://device/replay?token=REC1",
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Recording not found</s:Text></s:Reason>
						</s:Fault>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(response))
			},
			recordingToken: "INVALID",
			stream:         "RTP-Unicast",
			protocol:       "RTSP",
			wantErr:        true,
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

			uri, err := client.GetReplayUri(context.Background(), tt.recordingToken, tt.stream, tt.protocol)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReplayUri() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && uri != tt.wantURI {
				t.Errorf("GetReplayUri() = %v, want %v", uri, tt.wantURI)
			}
		})
	}
}

func TestGetReplayConfiguration(t *testing.T) {
	tests := []struct {
		name               string
		handler            http.HandlerFunc
		wantErr            bool
		wantSessionTimeout string
	}{
		{
			name: "successful configuration retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trp:GetReplayConfigurationResponse xmlns:trp="http://www.onvif.org/ver10/replay/wsdl">
							<trp:Configuration>
								<trp:SessionTimeout>PT60S</trp:SessionTimeout>
							</trp:Configuration>
						</trp:GetReplayConfigurationResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:            false,
			wantSessionTimeout: "PT60S",
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

			config, err := client.GetReplayConfiguration(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReplayConfiguration() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if config == nil {
					t.Fatal("Expected configuration, got nil")
				}

				if config.SessionTimeout != tt.wantSessionTimeout {
					t.Errorf("SessionTimeout = %v, want %v", config.SessionTimeout, tt.wantSessionTimeout)
				}
			}
		})
	}
}

func TestSetReplayConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		config  *ReplayConfiguration
		wantErr bool
	}{
		{
			name: "successful configuration update",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trp:SetReplayConfigurationResponse xmlns:trp="http://www.onvif.org/ver10/replay/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			config:  &ReplayConfiguration{SessionTimeout: "PT120S"},
			wantErr: false,
		},
		{
			name: "SOAP fault response",
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
			config:  &ReplayConfiguration{SessionTimeout: "INVALID"},
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

			err = client.SetReplayConfiguration(context.Background(), tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetReplayConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
