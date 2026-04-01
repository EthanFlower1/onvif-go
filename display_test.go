package onvif

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testDisplayXMLHeader = `<?xml version="1.0" encoding="UTF-8"?>`

const soapFaultDisplay = testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <SOAP-ENV:Fault>
      <SOAP-ENV:Code>
        <SOAP-ENV:Value>SOAP-ENV:Sender</SOAP-ENV:Value>
      </SOAP-ENV:Code>
      <SOAP-ENV:Reason>
        <SOAP-ENV:Text xml:lang="en">Invalid argument</SOAP-ENV:Text>
      </SOAP-ENV:Reason>
    </SOAP-ENV:Fault>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

func newMockDisplayServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml")

		body := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(body)
		bodyStr := string(body)

		var response string

		switch {
		case strings.Contains(bodyStr, "GetServiceCapabilities"):
			response = testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:GetServiceCapabilitiesResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl">
      <tls:Capabilities FixedLayout="false"/>
    </tls:GetServiceCapabilitiesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetLayout"):
			response = testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:GetLayoutResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl">
      <tls:Layout>
        <tls:PaneLayout>
          <tls:Pane>Pane_1</tls:Pane>
          <tls:Area bottom="0" top="1" right="1" left="0"/>
        </tls:PaneLayout>
      </tls:Layout>
    </tls:GetLayoutResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetLayout"):
			response = testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:SetLayoutResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetDisplayOptions"):
			response = testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:GetDisplayOptionsResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl">
      <tls:CodingCapabilities>
        <tls:InputTokensLimits>
          <tls:Max>4</tls:Max>
        </tls:InputTokensLimits>
        <tls:OutputTokensLimits>
          <tls:Max>2</tls:Max>
        </tls:OutputTokensLimits>
      </tls:CodingCapabilities>
    </tls:GetDisplayOptionsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetPaneConfigurations"):
			response = testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:GetPaneConfigurationsResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl">
      <tls:PaneConfiguration token="Pane_1">
        <tls:PaneName>Main Pane</tls:PaneName>
        <tls:ReceiverToken>Receiver_1</tls:ReceiverToken>
      </tls:PaneConfiguration>
      <tls:PaneConfiguration token="Pane_2">
        <tls:PaneName>Secondary Pane</tls:PaneName>
      </tls:PaneConfiguration>
    </tls:GetPaneConfigurationsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetPaneConfiguration"):
			response = testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:GetPaneConfigurationResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl">
      <tls:PaneConfiguration token="Pane_1">
        <tls:PaneName>Main Pane</tls:PaneName>
        <tls:ReceiverToken>Receiver_1</tls:ReceiverToken>
      </tls:PaneConfiguration>
    </tls:GetPaneConfigurationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetPaneConfigurations"):
			response = testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:SetPaneConfigurationsResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetPaneConfiguration"):
			response = testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:SetPaneConfigurationResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "CreatePaneConfiguration"):
			response = testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:CreatePaneConfigurationResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl">
      <tls:PaneToken>Pane_New</tls:PaneToken>
    </tls:CreatePaneConfigurationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "DeletePaneConfiguration"):
			response = testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:DeletePaneConfigurationResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		default:
			response = soapFaultDisplay
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
}

func TestGetDisplayServiceCapabilities(t *testing.T) {
	tests := []struct {
		name            string
		handler         http.HandlerFunc
		wantErr         bool
		wantFixedLayout *bool
	}{
		{
			name: "successful capabilities retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:GetServiceCapabilitiesResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl">
      <tls:Capabilities FixedLayout="false"/>
    </tls:GetServiceCapabilitiesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:         false,
			wantFixedLayout: boolPtr(false),
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultDisplay))
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

			caps, err := client.GetDisplayServiceCapabilities(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDisplayServiceCapabilities() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if caps == nil {
					t.Fatal("Expected capabilities, got nil")
				}

				if tt.wantFixedLayout != nil {
					if caps.FixedLayout == nil {
						t.Error("Expected FixedLayout to be set, got nil")
					} else if *caps.FixedLayout != *tt.wantFixedLayout {
						t.Errorf("FixedLayout = %v, want %v", *caps.FixedLayout, *tt.wantFixedLayout)
					}
				}
			}
		})
	}
}

func TestGetLayout(t *testing.T) {
	tests := []struct {
		name             string
		videoOutputToken string
		handler          http.HandlerFunc
		wantErr          bool
		wantPaneCount    int
	}{
		{
			name:             "successful layout retrieval",
			videoOutputToken: "VideoOutput_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:GetLayoutResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl">
      <tls:Layout>
        <tls:PaneLayout>
          <tls:Pane>Pane_1</tls:Pane>
          <tls:Area bottom="0" top="1" right="1" left="0"/>
        </tls:PaneLayout>
        <tls:PaneLayout>
          <tls:Pane>Pane_2</tls:Pane>
          <tls:Area bottom="0" top="0.5" right="0.5" left="0"/>
        </tls:PaneLayout>
      </tls:Layout>
    </tls:GetLayoutResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:       false,
			wantPaneCount: 2,
		},
		{
			name:             "SOAP fault response",
			videoOutputToken: "VideoOutput_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultDisplay))
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

			layout, err := client.GetLayout(context.Background(), tt.videoOutputToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLayout() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if layout == nil {
					t.Fatal("Expected layout, got nil")
				}

				if len(layout.Pane) != tt.wantPaneCount {
					t.Errorf("Pane count = %d, want %d", len(layout.Pane), tt.wantPaneCount)
				}
			}
		})
	}
}

func TestSetLayout(t *testing.T) {
	tests := []struct {
		name             string
		videoOutputToken string
		layout           Layout
		handler          http.HandlerFunc
		wantErr          bool
	}{
		{
			name:             "successful layout set",
			videoOutputToken: "VideoOutput_1",
			layout: Layout{
				Pane: []PaneLayout{
					{
						Pane: "Pane_1",
						Area: FloatRectangle{Bottom: 0, Top: 1, Right: 1, Left: 0},
					},
				},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:SetLayoutResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:             "SOAP fault response",
			videoOutputToken: "VideoOutput_1",
			layout:           Layout{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultDisplay))
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

			err = client.SetLayout(context.Background(), tt.videoOutputToken, tt.layout)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetLayout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetDisplayOptions(t *testing.T) {
	tests := []struct {
		name             string
		videoOutputToken string
		handler          http.HandlerFunc
		wantErr          bool
		wantInputMax     int
	}{
		{
			name:             "successful options retrieval",
			videoOutputToken: "VideoOutput_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:GetDisplayOptionsResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl">
      <tls:CodingCapabilities>
        <tls:InputTokensLimits>
          <tls:Max>4</tls:Max>
        </tls:InputTokensLimits>
        <tls:OutputTokensLimits>
          <tls:Max>2</tls:Max>
        </tls:OutputTokensLimits>
      </tls:CodingCapabilities>
    </tls:GetDisplayOptionsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:      false,
			wantInputMax: 4,
		},
		{
			name:             "SOAP fault response",
			videoOutputToken: "VideoOutput_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultDisplay))
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

			opts, err := client.GetDisplayOptions(context.Background(), tt.videoOutputToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDisplayOptions() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if opts == nil {
					t.Fatal("Expected options, got nil")
				}

				if opts.CodingCapabilities.InputTokensLimits == nil {
					t.Error("Expected InputTokensLimits to be set, got nil")
				} else if opts.CodingCapabilities.InputTokensLimits.Max != tt.wantInputMax {
					t.Errorf("InputTokensLimits.Max = %d, want %d", opts.CodingCapabilities.InputTokensLimits.Max, tt.wantInputMax)
				}
			}
		})
	}
}

func TestGetPaneConfigurations(t *testing.T) {
	tests := []struct {
		name             string
		videoOutputToken string
		handler          http.HandlerFunc
		wantErr          bool
		wantCount        int
	}{
		{
			name:             "successful pane configurations retrieval",
			videoOutputToken: "VideoOutput_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:GetPaneConfigurationsResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl">
      <tls:PaneConfiguration token="Pane_1">
        <tls:PaneName>Main Pane</tls:PaneName>
        <tls:ReceiverToken>Receiver_1</tls:ReceiverToken>
      </tls:PaneConfiguration>
      <tls:PaneConfiguration token="Pane_2">
        <tls:PaneName>Secondary Pane</tls:PaneName>
      </tls:PaneConfiguration>
    </tls:GetPaneConfigurationsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:             "SOAP fault response",
			videoOutputToken: "VideoOutput_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultDisplay))
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

			panes, err := client.GetPaneConfigurations(context.Background(), tt.videoOutputToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPaneConfigurations() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if len(panes) != tt.wantCount {
					t.Errorf("Pane count = %d, want %d", len(panes), tt.wantCount)
				}
			}
		})
	}
}

func TestGetPaneConfiguration(t *testing.T) {
	tests := []struct {
		name             string
		videoOutputToken string
		paneToken        string
		handler          http.HandlerFunc
		wantErr          bool
		wantPaneName     string
	}{
		{
			name:             "successful pane configuration retrieval",
			videoOutputToken: "VideoOutput_1",
			paneToken:        "Pane_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:GetPaneConfigurationResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl">
      <tls:PaneConfiguration token="Pane_1">
        <tls:PaneName>Main Pane</tls:PaneName>
        <tls:ReceiverToken>Receiver_1</tls:ReceiverToken>
      </tls:PaneConfiguration>
    </tls:GetPaneConfigurationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:      false,
			wantPaneName: "Main Pane",
		},
		{
			name:             "SOAP fault response",
			videoOutputToken: "VideoOutput_1",
			paneToken:        "Pane_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultDisplay))
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

			pane, err := client.GetPaneConfiguration(context.Background(), tt.videoOutputToken, tt.paneToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPaneConfiguration() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if pane == nil {
					t.Fatal("Expected pane configuration, got nil")
				}

				if pane.PaneName != tt.wantPaneName {
					t.Errorf("PaneName = %q, want %q", pane.PaneName, tt.wantPaneName)
				}
			}
		})
	}
}

func TestSetPaneConfigurations(t *testing.T) {
	tests := []struct {
		name             string
		videoOutputToken string
		paneConfigs      []PaneConfiguration
		handler          http.HandlerFunc
		wantErr          bool
	}{
		{
			name:             "successful pane configurations set",
			videoOutputToken: "VideoOutput_1",
			paneConfigs: []PaneConfiguration{
				{Token: "Pane_1", PaneName: "Main Pane"},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:SetPaneConfigurationsResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:             "SOAP fault response",
			videoOutputToken: "VideoOutput_1",
			paneConfigs:      []PaneConfiguration{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultDisplay))
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

			err = client.SetPaneConfigurations(context.Background(), tt.videoOutputToken, tt.paneConfigs)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetPaneConfigurations() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetPaneConfiguration(t *testing.T) {
	tests := []struct {
		name             string
		videoOutputToken string
		paneConfig       PaneConfiguration
		handler          http.HandlerFunc
		wantErr          bool
	}{
		{
			name:             "successful pane configuration set",
			videoOutputToken: "VideoOutput_1",
			paneConfig:       PaneConfiguration{Token: "Pane_1", PaneName: "Updated Pane"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:SetPaneConfigurationResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:             "SOAP fault response",
			videoOutputToken: "VideoOutput_1",
			paneConfig:       PaneConfiguration{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultDisplay))
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

			err = client.SetPaneConfiguration(context.Background(), tt.videoOutputToken, tt.paneConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetPaneConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreatePaneConfiguration(t *testing.T) {
	tests := []struct {
		name             string
		videoOutputToken string
		paneConfig       PaneConfiguration
		handler          http.HandlerFunc
		wantErr          bool
		wantToken        string
	}{
		{
			name:             "successful pane configuration creation",
			videoOutputToken: "VideoOutput_1",
			paneConfig:       PaneConfiguration{PaneName: "New Pane"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:CreatePaneConfigurationResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl">
      <tls:PaneToken>Pane_New</tls:PaneToken>
    </tls:CreatePaneConfigurationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantToken: "Pane_New",
		},
		{
			name:             "SOAP fault response",
			videoOutputToken: "VideoOutput_1",
			paneConfig:       PaneConfiguration{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultDisplay))
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

			token, err := client.CreatePaneConfiguration(context.Background(), tt.videoOutputToken, tt.paneConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePaneConfiguration() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if token != tt.wantToken {
					t.Errorf("PaneToken = %q, want %q", token, tt.wantToken)
				}
			}
		})
	}
}

func TestDeletePaneConfiguration(t *testing.T) {
	tests := []struct {
		name             string
		videoOutputToken string
		paneToken        string
		handler          http.HandlerFunc
		wantErr          bool
	}{
		{
			name:             "successful pane configuration deletion",
			videoOutputToken: "VideoOutput_1",
			paneToken:        "Pane_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				response := testDisplayXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tls:DeletePaneConfigurationResponse xmlns:tls="http://www.onvif.org/ver10/display/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr: false,
		},
		{
			name:             "SOAP fault response",
			videoOutputToken: "VideoOutput_1",
			paneToken:        "Pane_1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/soap+xml")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(soapFaultDisplay))
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

			err = client.DeletePaneConfiguration(context.Background(), tt.videoOutputToken, tt.paneToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeletePaneConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
