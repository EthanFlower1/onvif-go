package onvif

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProvisioningServiceCapabilities(t *testing.T) {
	tests := []struct {
		name               string
		handler            http.HandlerFunc
		wantErr            bool
		wantDefaultTimeout string
		wantSourceCount    int
	}{
		{
			name: "successful capabilities retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tpv:GetServiceCapabilitiesResponse xmlns:tpv="http://www.onvif.org/ver10/provisioning/wsdl">
							<tpv:Capabilities>
								<tpv:DefaultTimeout>PT10S</tpv:DefaultTimeout>
								<tpv:Source VideoSourceToken="vs1" MaximumPanMoves="1000" MaximumTiltMoves="1000" MaximumZoomMoves="500" AutoFocus="true"/>
							</tpv:Capabilities>
						</tpv:GetServiceCapabilitiesResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:            false,
			wantDefaultTimeout: "PT10S",
			wantSourceCount:    1,
		},
		{
			name: "capabilities with no sources",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tpv:GetServiceCapabilitiesResponse xmlns:tpv="http://www.onvif.org/ver10/provisioning/wsdl">
							<tpv:Capabilities>
								<tpv:DefaultTimeout>PT5S</tpv:DefaultTimeout>
							</tpv:Capabilities>
						</tpv:GetServiceCapabilitiesResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:            false,
			wantDefaultTimeout: "PT5S",
			wantSourceCount:    0,
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

			caps, err := client.GetProvisioningServiceCapabilities(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProvisioningServiceCapabilities() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if caps == nil {
					t.Fatal("Expected capabilities, got nil")
				}

				if caps.DefaultTimeout != tt.wantDefaultTimeout {
					t.Errorf("DefaultTimeout = %q, want %q", caps.DefaultTimeout, tt.wantDefaultTimeout)
				}

				if len(caps.Source) != tt.wantSourceCount {
					t.Errorf("len(Source) = %d, want %d", len(caps.Source), tt.wantSourceCount)
				}
			}
		})
	}
}

func TestPanMove(t *testing.T) {
	tests := []struct {
		name      string
		direction PanDirection
		handler   http.HandlerFunc
		wantErr   bool
	}{
		{
			name:      "successful pan left",
			direction: PanDirectionLeft,
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tpv:PanMoveResponse xmlns:tpv="http://www.onvif.org/ver10/provisioning/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:      "successful pan right",
			direction: PanDirectionRight,
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tpv:PanMoveResponse xmlns:tpv="http://www.onvif.org/ver10/provisioning/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:      "SOAP fault response",
			direction: PanDirectionLeft,
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Sender</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Invalid video source</s:Text></s:Reason>
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

			err = client.PanMove(context.Background(), "vs1", tt.direction, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("PanMove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTiltMove(t *testing.T) {
	tests := []struct {
		name      string
		direction TiltDirection
		handler   http.HandlerFunc
		wantErr   bool
	}{
		{
			name:      "successful tilt up",
			direction: TiltDirectionUp,
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tpv:TiltMoveResponse xmlns:tpv="http://www.onvif.org/ver10/provisioning/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:      "SOAP fault response",
			direction: TiltDirectionDown,
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Operation not supported</s:Text></s:Reason>
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

			err = client.TiltMove(context.Background(), "vs1", tt.direction, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("TiltMove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZoomMove(t *testing.T) {
	tests := []struct {
		name      string
		direction ZoomDirection
		handler   http.HandlerFunc
		wantErr   bool
	}{
		{
			name:      "successful zoom wide",
			direction: ZoomDirectionWide,
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tpv:ZoomMoveResponse xmlns:tpv="http://www.onvif.org/ver10/provisioning/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:      "SOAP fault response",
			direction: ZoomDirectionTelephoto,
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Zoom not supported</s:Text></s:Reason>
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

			err = client.ZoomMove(context.Background(), "vs1", tt.direction, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZoomMove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRollMove(t *testing.T) {
	tests := []struct {
		name      string
		direction RollDirection
		handler   http.HandlerFunc
		wantErr   bool
	}{
		{
			name:      "successful roll clockwise",
			direction: RollDirectionClockwise,
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tpv:RollMoveResponse xmlns:tpv="http://www.onvif.org/ver10/provisioning/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:      "successful roll auto-level",
			direction: RollDirectionAuto,
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tpv:RollMoveResponse xmlns:tpv="http://www.onvif.org/ver10/provisioning/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:      "SOAP fault response",
			direction: RollDirectionCounterclockwise,
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Roll not supported</s:Text></s:Reason>
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

			err = client.RollMove(context.Background(), "vs1", tt.direction, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("RollMove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProvisioningFocusMove(t *testing.T) {
	tests := []struct {
		name      string
		direction FocusDirection
		handler   http.HandlerFunc
		wantErr   bool
	}{
		{
			name:      "successful focus near",
			direction: FocusDirectionNear,
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tpv:FocusMoveResponse xmlns:tpv="http://www.onvif.org/ver10/provisioning/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:      "successful focus auto",
			direction: FocusDirectionAuto,
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tpv:FocusMoveResponse xmlns:tpv="http://www.onvif.org/ver10/provisioning/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:      "SOAP fault response",
			direction: FocusDirectionFar,
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Focus not supported</s:Text></s:Reason>
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

			err = client.ProvisioningFocusMove(context.Background(), "vs1", tt.direction, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProvisioningFocusMove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProvisioningStop(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful stop",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tpv:StopResponse xmlns:tpv="http://www.onvif.org/ver10/provisioning/wsdl"/>
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
							<s:Code><s:Value>s:Sender</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Invalid video source</s:Text></s:Reason>
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

			err = client.ProvisioningStop(context.Background(), "vs1")
			if (err != nil) != tt.wantErr {
				t.Errorf("ProvisioningStop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetProvisioningUsage(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		wantErr   bool
		checkUsage func(t *testing.T, usage *ProvisioningUsage)
	}{
		{
			name: "successful usage retrieval with all fields",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tpv:GetUsageResponse xmlns:tpv="http://www.onvif.org/ver10/provisioning/wsdl">
							<tpv:Usage>
								<tpv:Pan>150</tpv:Pan>
								<tpv:Tilt>200</tpv:Tilt>
								<tpv:Zoom>75</tpv:Zoom>
								<tpv:Roll>30</tpv:Roll>
								<tpv:Focus>500</tpv:Focus>
							</tpv:Usage>
						</tpv:GetUsageResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
			checkUsage: func(t *testing.T, usage *ProvisioningUsage) {
				t.Helper()

				if usage.Pan == nil || *usage.Pan != 150 {
					t.Errorf("Pan = %v, want 150", usage.Pan)
				}

				if usage.Tilt == nil || *usage.Tilt != 200 {
					t.Errorf("Tilt = %v, want 200", usage.Tilt)
				}

				if usage.Zoom == nil || *usage.Zoom != 75 {
					t.Errorf("Zoom = %v, want 75", usage.Zoom)
				}

				if usage.Roll == nil || *usage.Roll != 30 {
					t.Errorf("Roll = %v, want 30", usage.Roll)
				}

				if usage.Focus == nil || *usage.Focus != 500 {
					t.Errorf("Focus = %v, want 500", usage.Focus)
				}
			},
		},
		{
			name: "successful usage retrieval with partial fields",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tpv:GetUsageResponse xmlns:tpv="http://www.onvif.org/ver10/provisioning/wsdl">
							<tpv:Usage>
								<tpv:Pan>42</tpv:Pan>
							</tpv:Usage>
						</tpv:GetUsageResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
			checkUsage: func(t *testing.T, usage *ProvisioningUsage) {
				t.Helper()

				if usage.Pan == nil || *usage.Pan != 42 {
					t.Errorf("Pan = %v, want 42", usage.Pan)
				}

				if usage.Tilt != nil {
					t.Errorf("Tilt = %v, want nil", usage.Tilt)
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
							<s:Code><s:Value>s:Sender</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Invalid video source</s:Text></s:Reason>
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

			usage, err := client.GetProvisioningUsage(context.Background(), "vs1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProvisioningUsage() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if usage == nil {
					t.Fatal("Expected usage, got nil")
				}

				if tt.checkUsage != nil {
					tt.checkUsage(t, usage)
				}
			}
		})
	}
}
