package onvif

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetReceiverServiceCapabilities(t *testing.T) {
	tests := []struct {
		name                 string
		handler              http.HandlerFunc
		wantErr              bool
		wantRTPMulticast     bool
		wantRTPTCP           bool
		wantRTPRTSP_TCP      bool
		wantSupportedRecvrs  int
		wantMaxRTSPURILength int
	}{
		{
			name: "successful capabilities retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trv:GetServiceCapabilitiesResponse xmlns:trv="http://www.onvif.org/ver10/receiver/wsdl">
							<trv:Capabilities RTP_Multicast="true" RTP_TCP="true" RTP_RTSP_TCP="true" SupportedReceivers="10" MaximumRTSPURILength="1024"/>
						</trv:GetServiceCapabilitiesResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:              false,
			wantRTPMulticast:     true,
			wantRTPTCP:           true,
			wantRTPRTSP_TCP:      true,
			wantSupportedRecvrs:  10,
			wantMaxRTSPURILength: 1024,
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

			caps, err := client.GetReceiverServiceCapabilities(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReceiverServiceCapabilities() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if caps == nil {
					t.Fatal("Expected capabilities, got nil")
				}

				if caps.RTPMulticast != tt.wantRTPMulticast {
					t.Errorf("RTPMulticast = %v, want %v", caps.RTPMulticast, tt.wantRTPMulticast)
				}

				if caps.RTPTCP != tt.wantRTPTCP {
					t.Errorf("RTPTCP = %v, want %v", caps.RTPTCP, tt.wantRTPTCP)
				}

				if caps.RTPRTSP_TCP != tt.wantRTPRTSP_TCP {
					t.Errorf("RTPRTSP_TCP = %v, want %v", caps.RTPRTSP_TCP, tt.wantRTPRTSP_TCP)
				}

				if caps.SupportedReceivers != tt.wantSupportedRecvrs {
					t.Errorf("SupportedReceivers = %d, want %d", caps.SupportedReceivers, tt.wantSupportedRecvrs)
				}

				if caps.MaximumRTSPURILength != tt.wantMaxRTSPURILength {
					t.Errorf("MaximumRTSPURILength = %d, want %d", caps.MaximumRTSPURILength, tt.wantMaxRTSPURILength)
				}
			}
		})
	}
}

func TestGetReceivers(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		wantErr   bool
		wantCount int
		checkFirst func(t *testing.T, rec *Receiver)
	}{
		{
			name: "successful receivers retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trv:GetReceiversResponse xmlns:trv="http://www.onvif.org/ver10/receiver/wsdl">
							<trv:Receivers token="rcv1">
								<trv:Configuration>
									<trv:Mode>AlwaysConnect</trv:Mode>
									<trv:MediaUri>rtsp://camera/stream</trv:MediaUri>
									<trv:StreamSetup>
										<trv:Stream>RTP-Unicast</trv:Stream>
										<trv:Transport>
											<trv:Protocol>RTSP</trv:Protocol>
										</trv:Transport>
									</trv:StreamSetup>
								</trv:Configuration>
							</trv:Receivers>
						</trv:GetReceiversResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantCount: 1,
			checkFirst: func(t *testing.T, rec *Receiver) {
				t.Helper()
				if rec.Token != "rcv1" {
					t.Errorf("Token = %q, want %q", rec.Token, "rcv1")
				}
				if rec.Configuration.Mode != "AlwaysConnect" {
					t.Errorf("Mode = %q, want %q", rec.Configuration.Mode, "AlwaysConnect")
				}
				if rec.Configuration.MediaURI != "rtsp://camera/stream" {
					t.Errorf("MediaURI = %q, want %q", rec.Configuration.MediaURI, "rtsp://camera/stream")
				}
			},
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

			receivers, err := client.GetReceivers(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReceivers() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if len(receivers) != tt.wantCount {
					t.Errorf("len(receivers) = %d, want %d", len(receivers), tt.wantCount)

					return
				}
				if tt.checkFirst != nil && len(receivers) > 0 {
					tt.checkFirst(t, receivers[0])
				}
			}
		})
	}
}

func TestGetReceiver(t *testing.T) {
	tests := []struct {
		name          string
		receiverToken string
		handler       http.HandlerFunc
		wantErr       bool
		wantToken     string
		wantMode      string
	}{
		{
			name:          "successful single receiver retrieval",
			receiverToken: "rcv1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trv:GetReceiverResponse xmlns:trv="http://www.onvif.org/ver10/receiver/wsdl">
							<trv:Receiver token="rcv1">
								<trv:Configuration>
									<trv:Mode>AlwaysConnect</trv:Mode>
									<trv:MediaUri>rtsp://camera/stream1</trv:MediaUri>
									<trv:StreamSetup>
										<trv:Stream>RTP-Unicast</trv:Stream>
										<trv:Transport>
											<trv:Protocol>RTSP</trv:Protocol>
										</trv:Transport>
									</trv:StreamSetup>
								</trv:Configuration>
							</trv:Receiver>
						</trv:GetReceiverResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantToken: "rcv1",
			wantMode:  "AlwaysConnect",
		},
		{
			name:          "SOAP fault response",
			receiverToken: "rcv_bad",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Sender</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Unknown token</s:Text></s:Reason>
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

			receiver, err := client.GetReceiver(context.Background(), tt.receiverToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReceiver() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if receiver == nil {
					t.Fatal("Expected receiver, got nil")
				}

				if receiver.Token != tt.wantToken {
					t.Errorf("Token = %q, want %q", receiver.Token, tt.wantToken)
				}

				if receiver.Configuration.Mode != tt.wantMode {
					t.Errorf("Mode = %q, want %q", receiver.Configuration.Mode, tt.wantMode)
				}
			}
		})
	}
}

func TestCreateReceiver(t *testing.T) {
	tests := []struct {
		name      string
		config    *ReceiverConfiguration
		handler   http.HandlerFunc
		wantErr   bool
		wantToken string
	}{
		{
			name: "successful receiver creation",
			config: &ReceiverConfiguration{
				Mode:     "AlwaysConnect",
				MediaURI: "rtsp://camera/newstream",
				StreamSetup: &StreamSetup{
					Stream:    "RTP-Unicast",
					Transport: &Transport{Protocol: "RTSP"},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trv:CreateReceiverResponse xmlns:trv="http://www.onvif.org/ver10/receiver/wsdl">
							<trv:Receiver token="rcv_new">
								<trv:Configuration>
									<trv:Mode>AlwaysConnect</trv:Mode>
									<trv:MediaUri>rtsp://camera/newstream</trv:MediaUri>
									<trv:StreamSetup>
										<trv:Stream>RTP-Unicast</trv:Stream>
										<trv:Transport>
											<trv:Protocol>RTSP</trv:Protocol>
										</trv:Transport>
									</trv:StreamSetup>
								</trv:Configuration>
							</trv:Receiver>
						</trv:CreateReceiverResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantToken: "rcv_new",
		},
		{
			name: "SOAP fault response",
			config: &ReceiverConfiguration{
				Mode:     "AlwaysConnect",
				MediaURI: "rtsp://camera/stream",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Max receivers reached</s:Text></s:Reason>
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

			receiver, err := client.CreateReceiver(context.Background(), tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateReceiver() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if receiver == nil {
					t.Fatal("Expected receiver, got nil")
				}

				if receiver.Token != tt.wantToken {
					t.Errorf("Token = %q, want %q", receiver.Token, tt.wantToken)
				}
			}
		})
	}
}

func TestDeleteReceiver(t *testing.T) {
	tests := []struct {
		name          string
		receiverToken string
		handler       http.HandlerFunc
		wantErr       bool
	}{
		{
			name:          "successful receiver deletion",
			receiverToken: "rcv1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trv:DeleteReceiverResponse xmlns:trv="http://www.onvif.org/ver10/receiver/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:          "SOAP fault response",
			receiverToken: "rcv_bad",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Sender</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Unknown token</s:Text></s:Reason>
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

			err = client.DeleteReceiver(context.Background(), tt.receiverToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteReceiver() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigureReceiver(t *testing.T) {
	tests := []struct {
		name          string
		receiverToken string
		config        *ReceiverConfiguration
		handler       http.HandlerFunc
		wantErr       bool
	}{
		{
			name:          "successful receiver configuration",
			receiverToken: "rcv1",
			config: &ReceiverConfiguration{
				Mode:     "NeverConnect",
				MediaURI: "rtsp://camera/updated",
				StreamSetup: &StreamSetup{
					Stream:    "RTP-Unicast",
					Transport: &Transport{Protocol: "RTSP"},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trv:ConfigureReceiverResponse xmlns:trv="http://www.onvif.org/ver10/receiver/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:          "SOAP fault response",
			receiverToken: "rcv_bad",
			config: &ReceiverConfiguration{
				Mode:     "AlwaysConnect",
				MediaURI: "rtsp://camera/stream",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Sender</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Unknown token</s:Text></s:Reason>
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

			err = client.ConfigureReceiver(context.Background(), tt.receiverToken, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigureReceiver() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetReceiverMode(t *testing.T) {
	tests := []struct {
		name          string
		receiverToken string
		mode          string
		handler       http.HandlerFunc
		wantErr       bool
	}{
		{
			name:          "successful mode change",
			receiverToken: "rcv1",
			mode:          "NeverConnect",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trv:SetReceiverModeResponse xmlns:trv="http://www.onvif.org/ver10/receiver/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:          "SOAP fault response",
			receiverToken: "rcv_bad",
			mode:          "AlwaysConnect",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Sender</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Unknown token</s:Text></s:Reason>
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

			err = client.SetReceiverMode(context.Background(), tt.receiverToken, tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetReceiverMode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetReceiverState(t *testing.T) {
	tests := []struct {
		name          string
		receiverToken string
		handler       http.HandlerFunc
		wantErr       bool
		wantState     string
		wantAutoCreated bool
	}{
		{
			name:          "successful state retrieval",
			receiverToken: "rcv1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trv:GetReceiverStateResponse xmlns:trv="http://www.onvif.org/ver10/receiver/wsdl">
							<trv:ReceiverState>
								<trv:State>Connected</trv:State>
								<trv:AutoCreated>false</trv:AutoCreated>
							</trv:ReceiverState>
						</trv:GetReceiverStateResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:         false,
			wantState:       "Connected",
			wantAutoCreated: false,
		},
		{
			name:          "SOAP fault response",
			receiverToken: "rcv_bad",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Sender</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Unknown token</s:Text></s:Reason>
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

			state, err := client.GetReceiverState(context.Background(), tt.receiverToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReceiverState() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if state == nil {
					t.Fatal("Expected state, got nil")
				}

				if state.State != tt.wantState {
					t.Errorf("State = %q, want %q", state.State, tt.wantState)
				}

				if state.AutoCreated != tt.wantAutoCreated {
					t.Errorf("AutoCreated = %v, want %v", state.AutoCreated, tt.wantAutoCreated)
				}
			}
		})
	}
}
