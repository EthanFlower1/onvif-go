package onvif

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetSearchServiceCapabilities(t *testing.T) {
	tests := []struct {
		name                string
		handler             http.HandlerFunc
		wantErr             bool
		wantMetadataSearch  bool
	}{
		{
			name: "successful capabilities retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <tse:GetServiceCapabilitiesResponse xmlns:tse="http://www.onvif.org/ver10/search/wsdl">
      <tse:Capabilities MetadataSearch="true"/>
    </tse:GetServiceCapabilitiesResponse>
  </s:Body>
</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:            false,
			wantMetadataSearch: true,
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

			caps, err := client.GetSearchServiceCapabilities(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSearchServiceCapabilities() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if caps == nil {
					t.Fatal("Expected capabilities, got nil")
				}

				if caps.MetadataSearch != tt.wantMetadataSearch {
					t.Errorf("MetadataSearch = %v, want %v", caps.MetadataSearch, tt.wantMetadataSearch)
				}
			}
		})
	}
}

func TestGetRecordingSummary(t *testing.T) {
	tests := []struct {
		name                 string
		handler              http.HandlerFunc
		wantErr              bool
		wantNumberRecordings int
	}{
		{
			name: "successful summary retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <tse:GetRecordingSummaryResponse xmlns:tse="http://www.onvif.org/ver10/search/wsdl">
      <tse:Summary>
        <DataFrom>2024-01-01T00:00:00Z</DataFrom>
        <DataUntil>2024-12-31T23:59:59Z</DataUntil>
        <NumberRecordings>5</NumberRecordings>
      </tse:Summary>
    </tse:GetRecordingSummaryResponse>
  </s:Body>
</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:              false,
			wantNumberRecordings: 5,
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <s:Fault>
      <s:Code><s:Value>s:Receiver</s:Value></s:Code>
      <s:Reason><s:Text xml:lang="en">Service not available</s:Text></s:Reason>
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

			summary, err := client.GetRecordingSummary(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecordingSummary() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if summary == nil {
					t.Fatal("Expected summary, got nil")
				}

				if summary.NumberRecordings != tt.wantNumberRecordings {
					t.Errorf("NumberRecordings = %d, want %d", summary.NumberRecordings, tt.wantNumberRecordings)
				}

				if summary.DataFrom.IsZero() {
					t.Error("DataFrom should not be zero")
				}

				if summary.DataUntil.IsZero() {
					t.Error("DataUntil should not be zero")
				}
			}
		})
	}
}

func TestFindRecordings(t *testing.T) {
	tests := []struct {
		name            string
		handler         http.HandlerFunc
		scope           *SearchScope
		maxMatches      *int
		keepAliveTime   string
		wantErr         bool
		wantSearchToken string
	}{
		{
			name: "successful find recordings",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <tse:FindRecordingsResponse xmlns:tse="http://www.onvif.org/ver10/search/wsdl">
      <tse:SearchToken>SEARCH_001</tse:SearchToken>
    </tse:FindRecordingsResponse>
  </s:Body>
</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			scope: &SearchScope{
				IncludedRecordings: []string{"REC1"},
			},
			keepAliveTime:   "PT60S",
			wantErr:         false,
			wantSearchToken: "SEARCH_001",
		},
		{
			name: "find recordings with nil scope",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <tse:FindRecordingsResponse xmlns:tse="http://www.onvif.org/ver10/search/wsdl">
      <tse:SearchToken>SEARCH_002</tse:SearchToken>
    </tse:FindRecordingsResponse>
  </s:Body>
</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			scope:           nil,
			keepAliveTime:   "PT30S",
			wantErr:         false,
			wantSearchToken: "SEARCH_002",
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <s:Fault>
      <s:Code><s:Value>s:Receiver</s:Value></s:Code>
      <s:Reason><s:Text xml:lang="en">Search failed</s:Text></s:Reason>
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

			token, err := client.FindRecordings(context.Background(), tt.scope, tt.maxMatches, tt.keepAliveTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindRecordings() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if token != tt.wantSearchToken {
					t.Errorf("SearchToken = %q, want %q", token, tt.wantSearchToken)
				}
			}
		})
	}
}

func TestGetRecordingSearchResults(t *testing.T) {
	tests := []struct {
		name             string
		handler          http.HandlerFunc
		wantErr          bool
		wantSearchState  string
		wantRecordCount  int
	}{
		{
			name: "successful results retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <tse:GetRecordingSearchResultsResponse xmlns:tse="http://www.onvif.org/ver10/search/wsdl">
      <tse:ResultList>
        <SearchState>Completed</SearchState>
        <RecordingInformation>
          <RecordingToken>REC1</RecordingToken>
          <Content>Recording 1</Content>
          <RecordingStatus>Complete</RecordingStatus>
        </RecordingInformation>
        <RecordingInformation>
          <RecordingToken>REC2</RecordingToken>
          <Content>Recording 2</Content>
          <RecordingStatus>Active</RecordingStatus>
        </RecordingInformation>
      </tse:ResultList>
    </tse:GetRecordingSearchResultsResponse>
  </s:Body>
</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:         false,
			wantSearchState: "Completed",
			wantRecordCount: 2,
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <s:Fault>
      <s:Code><s:Value>s:Receiver</s:Value></s:Code>
      <s:Reason><s:Text xml:lang="en">Invalid search token</s:Text></s:Reason>
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

			maxResults := 10
			result, err := client.GetRecordingSearchResults(context.Background(), "SEARCH_001", nil, &maxResults, "PT5S")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecordingSearchResults() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Fatal("Expected result, got nil")
				}

				if result.SearchState != tt.wantSearchState {
					t.Errorf("SearchState = %q, want %q", result.SearchState, tt.wantSearchState)
				}

				if len(result.RecordingInformation) != tt.wantRecordCount {
					t.Errorf("RecordingInformation count = %d, want %d", len(result.RecordingInformation), tt.wantRecordCount)
				}

				if tt.wantRecordCount > 0 && result.RecordingInformation[0].RecordingToken != "REC1" {
					t.Errorf("First recording token = %q, want %q", result.RecordingInformation[0].RecordingToken, "REC1")
				}
			}
		})
	}
}

func TestGetSearchState(t *testing.T) {
	tests := []struct {
		name            string
		handler         http.HandlerFunc
		wantErr         bool
		wantSearchState string
	}{
		{
			name: "successful state retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <tse:GetSearchStateResponse xmlns:tse="http://www.onvif.org/ver10/search/wsdl">
      <tse:SearchState>Searching</tse:SearchState>
    </tse:GetSearchStateResponse>
  </s:Body>
</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:         false,
			wantSearchState: "Searching",
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <s:Fault>
      <s:Code><s:Value>s:Receiver</s:Value></s:Code>
      <s:Reason><s:Text xml:lang="en">Invalid token</s:Text></s:Reason>
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

			state, err := client.GetSearchState(context.Background(), "SEARCH_001")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSearchState() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if state != tt.wantSearchState {
					t.Errorf("SearchState = %q, want %q", state, tt.wantSearchState)
				}
			}
		})
	}
}

func TestEndSearch(t *testing.T) {
	tests := []struct {
		name            string
		handler         http.HandlerFunc
		wantErr         bool
		wantSearchToken string
	}{
		{
			name: "successful end search",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <tse:EndSearchResponse xmlns:tse="http://www.onvif.org/ver10/search/wsdl">
      <tse:SearchToken>SEARCH_001</tse:SearchToken>
    </tse:EndSearchResponse>
  </s:Body>
</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:         false,
			wantSearchToken: "SEARCH_001",
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <s:Fault>
      <s:Code><s:Value>s:Receiver</s:Value></s:Code>
      <s:Reason><s:Text xml:lang="en">Search token not found</s:Text></s:Reason>
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

			token, err := client.EndSearch(context.Background(), "SEARCH_001")
			if (err != nil) != tt.wantErr {
				t.Errorf("EndSearch() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if token != tt.wantSearchToken {
					t.Errorf("SearchToken = %q, want %q", token, tt.wantSearchToken)
				}
			}
		})
	}
}

func TestFindEvents(t *testing.T) {
	tests := []struct {
		name            string
		handler         http.HandlerFunc
		wantErr         bool
		wantSearchToken string
	}{
		{
			name: "successful find events",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <tse:FindEventsResponse xmlns:tse="http://www.onvif.org/ver10/search/wsdl">
      <tse:SearchToken>EVT_SEARCH_001</tse:SearchToken>
    </tse:FindEventsResponse>
  </s:Body>
</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:         false,
			wantSearchToken: "EVT_SEARCH_001",
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <s:Fault>
      <s:Code><s:Value>s:Receiver</s:Value></s:Code>
      <s:Reason><s:Text xml:lang="en">Event search failed</s:Text></s:Reason>
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

			start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
			token, err := client.FindEvents(context.Background(), start, nil, nil, false, nil, "PT60S")
			if (err != nil) != tt.wantErr {
				t.Errorf("FindEvents() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if token != tt.wantSearchToken {
					t.Errorf("SearchToken = %q, want %q", token, tt.wantSearchToken)
				}
			}
		})
	}
}

func TestGetEventSearchResults(t *testing.T) {
	tests := []struct {
		name            string
		handler         http.HandlerFunc
		wantErr         bool
		wantSearchState string
		wantEventCount  int
	}{
		{
			name: "successful event results",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <tse:GetEventSearchResultsResponse xmlns:tse="http://www.onvif.org/ver10/search/wsdl">
      <tse:ResultList>
        <SearchState>Completed</SearchState>
        <Result>
          <RecordingToken>REC1</RecordingToken>
          <TrackToken>TRK1</TrackToken>
          <Time>2024-01-01T00:00:00Z</Time>
          <StartStateEvent>false</StartStateEvent>
        </Result>
      </tse:ResultList>
    </tse:GetEventSearchResultsResponse>
  </s:Body>
</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:         false,
			wantSearchState: "Completed",
			wantEventCount:  1,
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Body>
    <s:Fault>
      <s:Code><s:Value>s:Receiver</s:Value></s:Code>
      <s:Reason><s:Text xml:lang="en">Invalid token</s:Text></s:Reason>
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

			result, err := client.GetEventSearchResults(context.Background(), "EVT_SEARCH_001", nil, nil, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventSearchResults() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Fatal("Expected result, got nil")
				}

				if result.SearchState != tt.wantSearchState {
					t.Errorf("SearchState = %q, want %q", result.SearchState, tt.wantSearchState)
				}

				if len(result.Events) != tt.wantEventCount {
					t.Errorf("Events count = %d, want %d", len(result.Events), tt.wantEventCount)
				}
			}
		})
	}
}
