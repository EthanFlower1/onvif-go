package onvif

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetNodes tests GetNodes operation.
func TestGetNodes(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful get nodes",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:GetNodesResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
			<tptz:PTZNode token="NodeToken1">
				<tt:Name xmlns:tt="http://www.onvif.org/ver10/schema">PTZ Node 1</tt:Name>
				<tt:HomeSupported xmlns:tt="http://www.onvif.org/ver10/schema">true</tt:HomeSupported>
				<tt:MaximumNumberOfPresets xmlns:tt="http://www.onvif.org/ver10/schema">255</tt:MaximumNumberOfPresets>
			</tptz:PTZNode>
		</tptz:GetNodesResponse>
	</soap:Body>
</soap:Envelope>`
				w.Header().Set("Content-Type", "application/soap+xml")
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
				t.Fatalf("NewClient() failed: %v", err)
			}

			client.ptzEndpoint = server.URL

			nodes, err := client.GetNodes(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodes() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if len(nodes) != 1 {
					t.Errorf("Expected 1 node, got %d", len(nodes))
				}

				if len(nodes) > 0 && nodes[0].Token != "NodeToken1" {
					t.Errorf("Expected token NodeToken1, got %s", nodes[0].Token)
				}

				if len(nodes) > 0 && nodes[0].Name != "PTZ Node 1" {
					t.Errorf("Expected name 'PTZ Node 1', got %s", nodes[0].Name)
				}

				if len(nodes) > 0 && nodes[0].MaximumNumberOfPresets != 255 {
					t.Errorf("Expected MaximumNumberOfPresets 255, got %d", nodes[0].MaximumNumberOfPresets)
				}
			}
		})
	}
}

// TestGetNode tests GetNode operation.
func TestGetNode(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful get node",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:GetNodeResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
			<tptz:PTZNode token="NodeToken1">
				<tt:Name xmlns:tt="http://www.onvif.org/ver10/schema">PTZ Node 1</tt:Name>
				<tt:HomeSupported xmlns:tt="http://www.onvif.org/ver10/schema">true</tt:HomeSupported>
				<tt:MaximumNumberOfPresets xmlns:tt="http://www.onvif.org/ver10/schema">100</tt:MaximumNumberOfPresets>
			</tptz:PTZNode>
		</tptz:GetNodeResponse>
	</soap:Body>
</soap:Envelope>`
				w.Header().Set("Content-Type", "application/soap+xml")
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
			<s:Reason><s:Text xml:lang="en">Node not found</s:Text></s:Reason>
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
				t.Fatalf("NewClient() failed: %v", err)
			}

			client.ptzEndpoint = server.URL

			node, err := client.GetNode(context.Background(), "NodeToken1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNode() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if node == nil {
					t.Fatal("Expected node, got nil")
				}

				if node.Token != "NodeToken1" {
					t.Errorf("Expected token NodeToken1, got %s", node.Token)
				}

				if node.Name != "PTZ Node 1" {
					t.Errorf("Expected name 'PTZ Node 1', got %s", node.Name)
				}

				if node.MaximumNumberOfPresets != 100 {
					t.Errorf("Expected MaximumNumberOfPresets 100, got %d", node.MaximumNumberOfPresets)
				}
			}
		})
	}
}

// TestGetPTZConfigurationOptions tests GetPTZConfigurationOptions operation.
func TestGetPTZConfigurationOptions(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful get configuration options",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:GetConfigurationOptionsResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
			<tptz:PTZConfigurationOptions>
				<tt:PTZTimeout xmlns:tt="http://www.onvif.org/ver10/schema">
					<tt:Min>PT1S</tt:Min>
					<tt:Max>PT60S</tt:Max>
				</tt:PTZTimeout>
			</tptz:PTZConfigurationOptions>
		</tptz:GetConfigurationOptionsResponse>
	</soap:Body>
</soap:Envelope>`
				w.Header().Set("Content-Type", "application/soap+xml")
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
			<s:Reason><s:Text xml:lang="en">Configuration not found</s:Text></s:Reason>
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
				t.Fatalf("NewClient() failed: %v", err)
			}

			client.ptzEndpoint = server.URL

			opts, err := client.GetPTZConfigurationOptions(context.Background(), "ConfigToken1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPTZConfigurationOptions() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if opts == nil {
					t.Fatal("Expected options, got nil")
				}

				if opts.PTZTimeout == nil {
					t.Fatal("Expected PTZTimeout, got nil")
				}

				if opts.PTZTimeout.Min != "PT1S" {
					t.Errorf("Expected PTZTimeout.Min PT1S, got %s", opts.PTZTimeout.Min)
				}

				if opts.PTZTimeout.Max != "PT60S" {
					t.Errorf("Expected PTZTimeout.Max PT60S, got %s", opts.PTZTimeout.Max)
				}
			}
		})
	}
}

// TestSetPTZConfiguration tests SetPTZConfiguration operation.
func TestSetPTZConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful set configuration",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:SetConfigurationResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl"/>
	</soap:Body>
</soap:Envelope>`
				w.Header().Set("Content-Type", "application/soap+xml")
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
			<s:Reason><s:Text xml:lang="en">Invalid configuration</s:Text></s:Reason>
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
				t.Fatalf("NewClient() failed: %v", err)
			}

			client.ptzEndpoint = server.URL

			config := &PTZConfiguration{
				Token:     "ConfigToken1",
				Name:      "Main PTZ Config",
				NodeToken: "NodeToken1",
			}

			err = client.SetPTZConfiguration(context.Background(), config, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetPTZConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetPTZServiceCapabilities tests GetPTZServiceCapabilities operation.
func TestGetPTZServiceCapabilities(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful get service capabilities",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:GetServiceCapabilitiesResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
			<tptz:Capabilities EFlip="true" Reverse="true"/>
		</tptz:GetServiceCapabilitiesResponse>
	</soap:Body>
</soap:Envelope>`
				w.Header().Set("Content-Type", "application/soap+xml")
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
				t.Fatalf("NewClient() failed: %v", err)
			}

			client.ptzEndpoint = server.URL

			caps, err := client.GetPTZServiceCapabilities(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPTZServiceCapabilities() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if caps == nil {
					t.Fatal("Expected capabilities, got nil")
				}

				if !caps.EFlip {
					t.Error("Expected EFlip to be true")
				}

				if !caps.Reverse {
					t.Error("Expected Reverse to be true")
				}
			}
		})
	}
}

// TestGetCompatiblePTZConfigurationsForProfile tests GetCompatiblePTZConfigurationsForProfile operation.
func TestGetCompatiblePTZConfigurationsForProfile(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful get compatible configurations",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:GetCompatibleConfigurationsResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
			<tptz:PTZConfiguration token="ConfigToken1">
				<tt:Name xmlns:tt="http://www.onvif.org/ver10/schema">Main PTZ Config</tt:Name>
				<tt:NodeToken xmlns:tt="http://www.onvif.org/ver10/schema">NodeToken1</tt:NodeToken>
			</tptz:PTZConfiguration>
		</tptz:GetCompatibleConfigurationsResponse>
	</soap:Body>
</soap:Envelope>`
				w.Header().Set("Content-Type", "application/soap+xml")
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
			<s:Reason><s:Text xml:lang="en">Profile not found</s:Text></s:Reason>
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
				t.Fatalf("NewClient() failed: %v", err)
			}

			client.ptzEndpoint = server.URL

			configs, err := client.GetCompatiblePTZConfigurationsForProfile(context.Background(), "Profile1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCompatiblePTZConfigurationsForProfile() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if len(configs) != 1 {
					t.Errorf("Expected 1 configuration, got %d", len(configs))
				}

				if len(configs) > 0 && configs[0].Token != "ConfigToken1" {
					t.Errorf("Expected token ConfigToken1, got %s", configs[0].Token)
				}

				if len(configs) > 0 && configs[0].Name != "Main PTZ Config" {
					t.Errorf("Expected name 'Main PTZ Config', got %s", configs[0].Name)
				}

				if len(configs) > 0 && configs[0].NodeToken != "NodeToken1" {
					t.Errorf("Expected NodeToken NodeToken1, got %s", configs[0].NodeToken)
				}
			}
		})
	}
}
