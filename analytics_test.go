package onvif

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAnalyticsServiceCapabilities(t *testing.T) {
	tests := []struct {
		name                string
		handler             http.HandlerFunc
		wantErr             bool
		wantRuleSupport     bool
		wantModuleSupport   bool
		wantCellBasedScene  bool
	}{
		{
			name: "successful capabilities retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tan:GetServiceCapabilitiesResponse xmlns:tan="http://www.onvif.org/ver20/analytics/wsdl">
							<tan:Capabilities RuleSupport="true" AnalyticsModuleSupport="true" CellBasedSceneDescriptionSupported="false"/>
						</tan:GetServiceCapabilitiesResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:            false,
			wantRuleSupport:    true,
			wantModuleSupport:  true,
			wantCellBasedScene: false,
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

			caps, err := client.GetAnalyticsServiceCapabilities(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAnalyticsServiceCapabilities() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if caps == nil {
					t.Fatal("Expected capabilities, got nil")
				}

				if caps.RuleSupport != tt.wantRuleSupport {
					t.Errorf("RuleSupport = %v, want %v", caps.RuleSupport, tt.wantRuleSupport)
				}

				if caps.AnalyticsModuleSupport != tt.wantModuleSupport {
					t.Errorf("AnalyticsModuleSupport = %v, want %v", caps.AnalyticsModuleSupport, tt.wantModuleSupport)
				}

				if caps.CellBasedSceneDescriptionSupported != tt.wantCellBasedScene {
					t.Errorf("CellBasedSceneDescriptionSupported = %v, want %v", caps.CellBasedSceneDescriptionSupported, tt.wantCellBasedScene)
				}
			}
		})
	}
}

func TestGetSupportedRules(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantErr    bool
		wantCount  int
		checkFirst func(t *testing.T, rule *SupportedRule)
	}{
		{
			name: "successful supported rules retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tan:GetSupportedRulesResponse xmlns:tan="http://www.onvif.org/ver20/analytics/wsdl">
							<tan:SupportedRules>
								<tan:RuleDescription Name="tt:LineDetector">
									<tan:Parameters>
										<tan:SimpleItemDescription Name="Direction" Type="xs:string"/>
									</tan:Parameters>
								</tan:RuleDescription>
								<tan:RuleDescription Name="tt:FieldDetector">
									<tan:Parameters>
										<tan:SimpleItemDescription Name="ActiveCells" Type="xs:hexBinary"/>
									</tan:Parameters>
								</tan:RuleDescription>
							</tan:SupportedRules>
						</tan:GetSupportedRulesResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantCount: 2,
			checkFirst: func(t *testing.T, rule *SupportedRule) {
				t.Helper()

				if rule.Name != "tt:LineDetector" {
					t.Errorf("Name = %v, want tt:LineDetector", rule.Name)
				}

				if len(rule.Parameters) != 1 {
					t.Fatalf("Expected 1 parameter, got %d", len(rule.Parameters))
				}

				if rule.Parameters[0].Name != "Direction" {
					t.Errorf("Parameter Name = %v, want Direction", rule.Parameters[0].Name)
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
							<s:Reason><s:Text xml:lang="en">Invalid token</s:Text></s:Reason>
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

			rules, err := client.GetSupportedRules(context.Background(), "vac1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSupportedRules() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if len(rules) != tt.wantCount {
					t.Errorf("Expected %d rules, got %d", tt.wantCount, len(rules))

					return
				}

				if tt.checkFirst != nil && len(rules) > 0 {
					tt.checkFirst(t, rules[0])
				}
			}
		})
	}
}

func TestGetRules(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantErr    bool
		wantCount  int
		checkFirst func(t *testing.T, rule *AnalyticsRule)
	}{
		{
			name: "successful rules retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tan:GetRulesResponse xmlns:tan="http://www.onvif.org/ver20/analytics/wsdl">
							<tan:Rule Name="MyRule" Type="tt:LineDetector">
								<tan:Parameters>
									<tan:SimpleItem Name="Direction" Value="Left"/>
								</tan:Parameters>
							</tan:Rule>
						</tan:GetRulesResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantCount: 1,
			checkFirst: func(t *testing.T, rule *AnalyticsRule) {
				t.Helper()

				if rule.Name != "MyRule" {
					t.Errorf("Name = %v, want MyRule", rule.Name)
				}

				if rule.Type != "tt:LineDetector" {
					t.Errorf("Type = %v, want tt:LineDetector", rule.Type)
				}

				if len(rule.Parameters) != 1 {
					t.Fatalf("Expected 1 parameter, got %d", len(rule.Parameters))
				}

				if rule.Parameters[0].Name != "Direction" {
					t.Errorf("Parameter Name = %v, want Direction", rule.Parameters[0].Name)
				}

				if rule.Parameters[0].Value != "Left" {
					t.Errorf("Parameter Value = %v, want Left", rule.Parameters[0].Value)
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

			rules, err := client.GetRules(context.Background(), "vac1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRules() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if len(rules) != tt.wantCount {
					t.Errorf("Expected %d rules, got %d", tt.wantCount, len(rules))

					return
				}

				if tt.checkFirst != nil && len(rules) > 0 {
					tt.checkFirst(t, rules[0])
				}
			}
		})
	}
}

func TestCreateRules(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		rules   []*AnalyticsRule
		wantErr bool
	}{
		{
			name: "successful rule creation",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tan:CreateRulesResponse xmlns:tan="http://www.onvif.org/ver20/analytics/wsdl"/>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			rules: []*AnalyticsRule{
				{
					Name: "NewRule",
					Type: "tt:LineDetector",
					Parameters: []*SimpleItem{
						{Name: "Direction", Value: "Right"},
					},
				},
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
							<s:Reason><s:Text xml:lang="en">Rule already exists</s:Text></s:Reason>
						</s:Fault>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(response))
			},
			rules:   []*AnalyticsRule{{Name: "Dup", Type: "tt:LineDetector"}},
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

			err = client.CreateRules(context.Background(), "vac1", tt.rules)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRules() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetAnalyticsModules(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantErr    bool
		wantCount  int
		checkFirst func(t *testing.T, mod *AnalyticsModule)
	}{
		{
			name: "successful modules retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tan:GetAnalyticsModulesResponse xmlns:tan="http://www.onvif.org/ver20/analytics/wsdl">
							<tan:AnalyticsModule Name="MyModule" Type="tt:CellMotionDetector">
								<tan:Parameters>
									<tan:SimpleItem Name="Sensitivity" Value="50"/>
								</tan:Parameters>
							</tan:AnalyticsModule>
						</tan:GetAnalyticsModulesResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantCount: 1,
			checkFirst: func(t *testing.T, mod *AnalyticsModule) {
				t.Helper()

				if mod.Name != "MyModule" {
					t.Errorf("Name = %v, want MyModule", mod.Name)
				}

				if mod.Type != "tt:CellMotionDetector" {
					t.Errorf("Type = %v, want tt:CellMotionDetector", mod.Type)
				}

				if len(mod.Parameters) != 1 {
					t.Fatalf("Expected 1 parameter, got %d", len(mod.Parameters))
				}

				if mod.Parameters[0].Name != "Sensitivity" {
					t.Errorf("Parameter Name = %v, want Sensitivity", mod.Parameters[0].Name)
				}

				if mod.Parameters[0].Value != "50" {
					t.Errorf("Parameter Value = %v, want 50", mod.Parameters[0].Value)
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
							<s:Reason><s:Text xml:lang="en">Not supported</s:Text></s:Reason>
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

			modules, err := client.GetAnalyticsModules(context.Background(), "vac1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAnalyticsModules() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if len(modules) != tt.wantCount {
					t.Errorf("Expected %d modules, got %d", tt.wantCount, len(modules))

					return
				}

				if tt.checkFirst != nil && len(modules) > 0 {
					tt.checkFirst(t, modules[0])
				}
			}
		})
	}
}

func TestGetSupportedMetadata(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		wantErr   bool
		wantCount int
		wantFirst string
	}{
		{
			name: "successful metadata retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<tan:GetSupportedMetadataResponse xmlns:tan="http://www.onvif.org/ver20/analytics/wsdl">
							<tan:SupportedMetadata>
								<tan:AnalyticsModule>tt:CellMotionDetector</tan:AnalyticsModule>
								<tan:AnalyticsModule>tt:LineDetector</tan:AnalyticsModule>
							</tan:SupportedMetadata>
						</tan:GetSupportedMetadataResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantCount: 2,
			wantFirst: "tt:CellMotionDetector",
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Service error</s:Text></s:Reason>
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

			meta, err := client.GetSupportedMetadata(context.Background(), "tt:AnalyticsModule")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSupportedMetadata() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if meta == nil {
					t.Fatal("Expected metadata, got nil")
				}

				if len(meta.AnalyticsModules) != tt.wantCount {
					t.Errorf("Expected %d modules, got %d", tt.wantCount, len(meta.AnalyticsModules))

					return
				}

				if tt.wantFirst != "" && len(meta.AnalyticsModules) > 0 {
					if meta.AnalyticsModules[0] != tt.wantFirst {
						t.Errorf("First module = %v, want %v", meta.AnalyticsModules[0], tt.wantFirst)
					}
				}
			}
		})
	}
}
