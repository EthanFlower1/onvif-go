package onvif

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
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

// TestGetPresetTours tests GetPresetTours operation.
func TestGetPresetTours(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful get preset tours",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:GetPresetToursResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
			<tptz:PresetTour token="Tour1">
				<tt:Name xmlns:tt="http://www.onvif.org/ver10/schema">Tour One</tt:Name>
				<tt:Status xmlns:tt="http://www.onvif.org/ver10/schema"><tt:State>Idle</tt:State></tt:Status>
				<tt:AutoStart xmlns:tt="http://www.onvif.org/ver10/schema">false</tt:AutoStart>
			</tptz:PresetTour>
		</tptz:GetPresetToursResponse>
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

			tours, err := client.GetPresetTours(context.Background(), "Profile1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPresetTours() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if len(tours) != 1 {
					t.Errorf("Expected 1 tour, got %d", len(tours))
				}

				if len(tours) > 0 && tours[0].Token != "Tour1" {
					t.Errorf("Expected token Tour1, got %s", tours[0].Token)
				}

				if len(tours) > 0 && tours[0].Name != "Tour One" {
					t.Errorf("Expected name 'Tour One', got %s", tours[0].Name)
				}
			}
		})
	}
}

// TestGetPresetTour tests GetPresetTour operation.
func TestGetPresetTour(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful get preset tour",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:GetPresetTourResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
			<tptz:PresetTour token="Tour1">
				<tt:Name xmlns:tt="http://www.onvif.org/ver10/schema">Tour One</tt:Name>
				<tt:Status xmlns:tt="http://www.onvif.org/ver10/schema"><tt:State>Idle</tt:State></tt:Status>
				<tt:AutoStart xmlns:tt="http://www.onvif.org/ver10/schema">true</tt:AutoStart>
			</tptz:PresetTour>
		</tptz:GetPresetTourResponse>
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
			<s:Reason><s:Text xml:lang="en">Tour not found</s:Text></s:Reason>
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

			tour, err := client.GetPresetTour(context.Background(), "Profile1", "Tour1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPresetTour() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if tour == nil {
					t.Fatal("Expected tour, got nil")
				}

				if tour.Token != "Tour1" {
					t.Errorf("Expected token Tour1, got %s", tour.Token)
				}

				if tour.Name != "Tour One" {
					t.Errorf("Expected name 'Tour One', got %s", tour.Name)
				}

				if !tour.AutoStart {
					t.Error("Expected AutoStart to be true")
				}
			}
		})
	}
}

// TestGetPresetTourOptions tests GetPresetTourOptions operation.
func TestGetPresetTourOptions(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful get preset tour options",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:GetPresetTourOptionsResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
			<tptz:Options>
				<tt:AutoStart xmlns:tt="http://www.onvif.org/ver10/schema">true</tt:AutoStart>
				<tt:StartingCondition xmlns:tt="http://www.onvif.org/ver10/schema">
					<tt:RecurringTimeRange>
						<tt:Min>1</tt:Min>
						<tt:Max>100</tt:Max>
					</tt:RecurringTimeRange>
					<tt:RecurringDurationRange>
						<tt:Min>PT10S</tt:Min>
						<tt:Max>PT3600S</tt:Max>
					</tt:RecurringDurationRange>
				</tt:StartingCondition>
			</tptz:Options>
		</tptz:GetPresetTourOptionsResponse>
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

			opts, err := client.GetPresetTourOptions(context.Background(), "Profile1", "Tour1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPresetTourOptions() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if opts == nil {
					t.Fatal("Expected options, got nil")
				}

				if !opts.AutoStart {
					t.Error("Expected AutoStart to be true")
				}

				if opts.StartingCondition == nil {
					t.Fatal("Expected StartingCondition, got nil")
				}

				if opts.StartingCondition.RecurringTimeRange == nil {
					t.Fatal("Expected RecurringTimeRange, got nil")
				}

				if opts.StartingCondition.RecurringTimeRange.Min != 1 {
					t.Errorf("Expected RecurringTimeRange.Min 1, got %d", opts.StartingCondition.RecurringTimeRange.Min)
				}

				if opts.StartingCondition.RecurringTimeRange.Max != 100 {
					t.Errorf("Expected RecurringTimeRange.Max 100, got %d", opts.StartingCondition.RecurringTimeRange.Max)
				}
			}
		})
	}
}

// TestCreatePresetTour tests CreatePresetTour operation.
func TestCreatePresetTour(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful create preset tour",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:CreatePresetTourResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
			<tptz:PresetTourToken>NewTour1</tptz:PresetTourToken>
		</tptz:CreatePresetTourResponse>
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

			token, err := client.CreatePresetTour(context.Background(), "Profile1")
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePresetTour() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && token != "NewTour1" {
				t.Errorf("Expected token NewTour1, got %s", token)
			}
		})
	}
}

// TestModifyPresetTour tests ModifyPresetTour operation.
func TestModifyPresetTour(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful modify preset tour",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:ModifyPresetTourResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl"/>
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
			<s:Reason><s:Text xml:lang="en">Tour not found</s:Text></s:Reason>
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

			tour := &PresetTour{
				Token:     "Tour1",
				Name:      "Updated Tour",
				AutoStart: true,
			}

			err = client.ModifyPresetTour(context.Background(), "Profile1", tour)
			if (err != nil) != tt.wantErr {
				t.Errorf("ModifyPresetTour() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestOperatePresetTour tests OperatePresetTour operation.
func TestOperatePresetTour(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful operate preset tour",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:OperatePresetTourResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl"/>
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
			<s:Reason><s:Text xml:lang="en">Tour not found</s:Text></s:Reason>
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

			err = client.OperatePresetTour(context.Background(), "Profile1", "Tour1", "Start")
			if (err != nil) != tt.wantErr {
				t.Errorf("OperatePresetTour() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestRemovePresetTour tests RemovePresetTour operation.
func TestRemovePresetTour(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful remove preset tour",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:RemovePresetTourResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl"/>
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
			<s:Reason><s:Text xml:lang="en">Tour not found</s:Text></s:Reason>
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

			err = client.RemovePresetTour(context.Background(), "Profile1", "Tour1")
			if (err != nil) != tt.wantErr {
				t.Errorf("RemovePresetTour() error = %v, wantErr %v", err, tt.wantErr)
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

// TestPTZSendAuxiliaryCommand tests PTZSendAuxiliaryCommand operation.
func TestPTZSendAuxiliaryCommand(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful send auxiliary command",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:SendAuxiliaryCommandResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
			<tptz:AuxiliaryResponse>OK</tptz:AuxiliaryResponse>
		</tptz:SendAuxiliaryCommandResponse>
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
			<s:Reason><s:Text xml:lang="en">Invalid auxiliary data</s:Text></s:Reason>
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

			resp, err := client.PTZSendAuxiliaryCommand(context.Background(), "Profile1", "tt:Wiper|On")
			if (err != nil) != tt.wantErr {
				t.Errorf("PTZSendAuxiliaryCommand() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && resp != "OK" {
				t.Errorf("Expected response 'OK', got %s", resp)
			}
		})
	}
}

// TestPTZSendAuxiliaryCommandUsesTPTZNamespace verifies the request uses tptz: namespace.
func TestPTZSendAuxiliaryCommandUsesTPTZNamespace(t *testing.T) {
	var capturedBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 4096)
		n, _ := r.Body.Read(buf)
		capturedBody = buf[:n]

		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:SendAuxiliaryCommandResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
			<tptz:AuxiliaryResponse>OK</tptz:AuxiliaryResponse>
		</tptz:SendAuxiliaryCommandResponse>
	</soap:Body>
</soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	client.ptzEndpoint = server.URL

	_, err = client.PTZSendAuxiliaryCommand(context.Background(), "Profile1", "tt:Wiper|On")
	if err != nil {
		t.Fatalf("PTZSendAuxiliaryCommand() failed: %v", err)
	}

	bodyStr := string(capturedBody)
	if !strings.Contains(bodyStr, "tptz:SendAuxiliaryCommand") {
		t.Errorf("Expected request body to contain 'tptz:SendAuxiliaryCommand', got: %s", bodyStr)
	}
}

// TestGeoMove tests GeoMove operation.
func TestGeoMove(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful geo move",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:GeoMoveResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
		</tptz:GeoMoveResponse>
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
			<s:Reason><s:Text xml:lang="en">GeoMove not supported</s:Text></s:Reason>
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

			target := &GeoLocation{Lon: 13.404954, Lat: 52.520008}
			err = client.GeoMove(context.Background(), "Profile1", target, nil, nil, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeoMove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGeoMoveRequestBody verifies the GeoMove request body contains tptz:GeoMove.
func TestGeoMoveRequestBody(t *testing.T) {
	var capturedBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 4096)
		n, _ := r.Body.Read(buf)
		capturedBody = buf[:n]

		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:GeoMoveResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
		</tptz:GeoMoveResponse>
	</soap:Body>
</soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	client.ptzEndpoint = server.URL

	target := &GeoLocation{Lon: 13.404954, Lat: 52.520008}
	err = client.GeoMove(context.Background(), "Profile1", target, nil, nil, nil)
	if err != nil {
		t.Fatalf("GeoMove() failed: %v", err)
	}

	bodyStr := string(capturedBody)
	if !strings.Contains(bodyStr, "tptz:GeoMove") {
		t.Errorf("Expected request body to contain 'tptz:GeoMove', got: %s", bodyStr)
	}
}

// TestMoveAndStartTracking tests MoveAndStartTracking operation.
func TestMoveAndStartTracking(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful move and start tracking",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
	<soap:Body>
		<tptz:MoveAndStartTrackingResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
		</tptz:MoveAndStartTrackingResponse>
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
			<s:Reason><s:Text xml:lang="en">Tracking not supported</s:Text></s:Reason>
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

			presetToken := "Preset1"
			req := &MoveAndStartTrackingRequest{
				ProfileToken: "Profile1",
				PresetToken:  &presetToken,
			}

			err = client.MoveAndStartTracking(context.Background(), req)
			if (err != nil) != tt.wantErr {
				t.Errorf("MoveAndStartTracking() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
