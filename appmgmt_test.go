package onvif

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAppMgmtServiceCapabilities(t *testing.T) {
	tests := []struct {
		name             string
		handler          http.HandlerFunc
		wantErr          bool
		wantFormats      string
		wantLicensing    *bool
		wantUploadPath   string
		wantEventPrefix  string
	}{
		{
			name: "successful capabilities with all fields",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tap:GetServiceCapabilitiesResponse xmlns:tap="http://www.onvif.org/ver10/appmgmt/wsdl">
							<tap:Capabilities FormatsSupported="docker" Licensing="true" UploadPath="/onvif/appmgmt/upload" EventTopicPrefix="tns1:AppMgmt"/>
						</tap:GetServiceCapabilitiesResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:         false,
			wantFormats:     "docker",
			wantLicensing:   func() *bool { v := true; return &v }(),
			wantUploadPath:  "/onvif/appmgmt/upload",
			wantEventPrefix: "tns1:AppMgmt",
		},
		{
			name: "successful capabilities with minimal fields",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tap:GetServiceCapabilitiesResponse xmlns:tap="http://www.onvif.org/ver10/appmgmt/wsdl">
							<tap:Capabilities FormatsSupported="docker"/>
						</tap:GetServiceCapabilitiesResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:     false,
			wantFormats: "docker",
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

			caps, err := client.GetAppMgmtServiceCapabilities(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAppMgmtServiceCapabilities() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if caps == nil {
					t.Fatal("Expected capabilities, got nil")
				}

				if caps.FormatsSupported != tt.wantFormats {
					t.Errorf("FormatsSupported = %q, want %q", caps.FormatsSupported, tt.wantFormats)
				}

				if tt.wantLicensing != nil {
					if caps.Licensing == nil || *caps.Licensing != *tt.wantLicensing {
						t.Errorf("Licensing = %v, want %v", caps.Licensing, tt.wantLicensing)
					}
				}

				if caps.UploadPath != tt.wantUploadPath {
					t.Errorf("UploadPath = %q, want %q", caps.UploadPath, tt.wantUploadPath)
				}

				if caps.EventTopicPrefix != tt.wantEventPrefix {
					t.Errorf("EventTopicPrefix = %q, want %q", caps.EventTopicPrefix, tt.wantEventPrefix)
				}
			}
		})
	}
}

func TestGetInstalledApps(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		wantErr   bool
		wantCount int
		wantName  string
		wantAppID string
	}{
		{
			name: "successful list with apps",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tap:GetInstalledAppsResponse xmlns:tap="http://www.onvif.org/ver10/appmgmt/wsdl">
							<tap:App>
								<tap:Name>MyApp</tap:Name>
								<tap:AppID>app-001</tap:AppID>
							</tap:App>
							<tap:App>
								<tap:Name>OtherApp</tap:Name>
								<tap:AppID>app-002</tap:AppID>
							</tap:App>
						</tap:GetInstalledAppsResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantCount: 2,
			wantName:  "MyApp",
			wantAppID: "app-001",
		},
		{
			name: "successful list with no apps",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tap:GetInstalledAppsResponse xmlns:tap="http://www.onvif.org/ver10/appmgmt/wsdl">
						</tap:GetInstalledAppsResponse>
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

			apps, err := client.GetInstalledApps(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInstalledApps() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if len(apps) != tt.wantCount {
					t.Errorf("Got %d apps, want %d", len(apps), tt.wantCount)
				}

				if tt.wantCount > 0 {
					if apps[0].Name != tt.wantName {
						t.Errorf("App[0].Name = %q, want %q", apps[0].Name, tt.wantName)
					}

					if apps[0].AppID != tt.wantAppID {
						t.Errorf("App[0].AppID = %q, want %q", apps[0].AppID, tt.wantAppID)
					}
				}
			}
		})
	}
}

func TestGetAppsInfo(t *testing.T) {
	tests := []struct {
		name      string
		appID     string
		handler   http.HandlerFunc
		wantErr   bool
		wantCount int
		wantAppID string
		wantState AppState
	}{
		{
			name:  "successful info for all apps",
			appID: "",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tap:GetAppsInfoResponse xmlns:tap="http://www.onvif.org/ver10/appmgmt/wsdl">
							<tap:Info>
								<tap:AppID>app-001</tap:AppID>
								<tap:Name>MyApp</tap:Name>
								<tap:Version>1.2.3</tap:Version>
								<tap:InstallationDate>2024-01-15T10:00:00Z</tap:InstallationDate>
								<tap:LastUpdate>2024-06-01T08:30:00Z</tap:LastUpdate>
								<tap:State>Active</tap:State>
								<tap:Status>Running normally</tap:Status>
								<tap:Autostart>true</tap:Autostart>
								<tap:Website>https://example.com/myapp</tap:Website>
							</tap:Info>
						</tap:GetAppsInfoResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantCount: 1,
			wantAppID: "app-001",
			wantState: AppStateActive,
		},
		{
			name:  "successful info for specific app",
			appID: "app-002",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tap:GetAppsInfoResponse xmlns:tap="http://www.onvif.org/ver10/appmgmt/wsdl">
							<tap:Info>
								<tap:AppID>app-002</tap:AppID>
								<tap:Name>OtherApp</tap:Name>
								<tap:Version>2.0.0</tap:Version>
								<tap:InstallationDate>2024-03-10T12:00:00Z</tap:InstallationDate>
								<tap:LastUpdate>2024-03-10T12:00:00Z</tap:LastUpdate>
								<tap:State>Inactive</tap:State>
								<tap:Status>Stopped</tap:Status>
								<tap:Autostart>false</tap:Autostart>
								<tap:Website>https://example.com/otherapp</tap:Website>
							</tap:Info>
						</tap:GetAppsInfoResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantCount: 1,
			wantAppID: "app-002",
			wantState: AppStateInactive,
		},
		{
			name:  "SOAP fault response",
			appID: "",
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

			infos, err := client.GetAppsInfo(context.Background(), tt.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAppsInfo() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if len(infos) != tt.wantCount {
					t.Errorf("Got %d infos, want %d", len(infos), tt.wantCount)
				}

				if tt.wantCount > 0 {
					if infos[0].AppID != tt.wantAppID {
						t.Errorf("Info[0].AppID = %q, want %q", infos[0].AppID, tt.wantAppID)
					}

					if infos[0].State != tt.wantState {
						t.Errorf("Info[0].State = %q, want %q", infos[0].State, tt.wantState)
					}
				}
			}
		})
	}
}

func TestActivateApp(t *testing.T) {
	tests := []struct {
		name    string
		appID   string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name:  "successful activation",
			appID: "app-001",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tap:ActivateResponse xmlns:tap="http://www.onvif.org/ver10/appmgmt/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:  "SOAP fault on activation",
			appID: "app-999",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Sender</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">App not found</s:Text></s:Reason>
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

			err = client.ActivateApp(context.Background(), tt.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ActivateApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeactivateApp(t *testing.T) {
	tests := []struct {
		name    string
		appID   string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name:  "successful deactivation",
			appID: "app-001",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tap:DeactivateResponse xmlns:tap="http://www.onvif.org/ver10/appmgmt/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:  "SOAP fault on deactivation",
			appID: "app-999",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Sender</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">App not found</s:Text></s:Reason>
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

			err = client.DeactivateApp(context.Background(), tt.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeactivateApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUninstallApp(t *testing.T) {
	tests := []struct {
		name    string
		appID   string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name:  "successful uninstall",
			appID: "app-001",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tap:UninstallResponse xmlns:tap="http://www.onvif.org/ver10/appmgmt/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:  "SOAP fault on uninstall",
			appID: "app-999",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Sender</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">App not found</s:Text></s:Reason>
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

			err = client.UninstallApp(context.Background(), tt.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UninstallApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstallLicense(t *testing.T) {
	tests := []struct {
		name    string
		appID   string
		license string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name:    "successful global license install",
			appID:   "",
			license: "LICENSE-TOKEN-GLOBAL-12345",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tap:InstallLicenseResponse xmlns:tap="http://www.onvif.org/ver10/appmgmt/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:    "successful per-app license install",
			appID:   "app-001",
			license: "LICENSE-TOKEN-APP-67890",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tap:InstallLicenseResponse xmlns:tap="http://www.onvif.org/ver10/appmgmt/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:    "SOAP fault on invalid license",
			appID:   "",
			license: "INVALID-LICENSE",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Sender</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Invalid license</s:Text></s:Reason>
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

			err = client.InstallLicense(context.Background(), tt.appID, tt.license)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstallLicense() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetAppDeviceID(t *testing.T) {
	tests := []struct {
		name         string
		handler      http.HandlerFunc
		wantErr      bool
		wantDeviceID string
	}{
		{
			name: "successful device ID retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tap:GetDeviceIdResponse xmlns:tap="http://www.onvif.org/ver10/appmgmt/wsdl">
							<tap:DeviceId>DEVICE-SERIAL-ABCDEF123456</tap:DeviceId>
						</tap:GetDeviceIdResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:      false,
			wantDeviceID: "DEVICE-SERIAL-ABCDEF123456",
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

			deviceID, err := client.GetAppDeviceID(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAppDeviceID() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && deviceID != tt.wantDeviceID {
				t.Errorf("GetAppDeviceID() = %q, want %q", deviceID, tt.wantDeviceID)
			}
		})
	}
}
