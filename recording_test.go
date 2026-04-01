package onvif

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRecordingServiceCapabilities(t *testing.T) {
	tests := []struct {
		name                  string
		handler               http.HandlerFunc
		wantErr               bool
		wantDynamicRecordings bool
		wantDynamicTracks     bool
		wantMaxRecordings     int
		wantMaxRecordingJobs  int
	}{
		{
			name: "successful capabilities retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:GetServiceCapabilitiesResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
							<trc:Capabilities DynamicRecordings="true" DynamicTracks="true" MaxRecordings="100" MaxRecordingJobs="10"/>
						</trc:GetServiceCapabilitiesResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:               false,
			wantDynamicRecordings: true,
			wantDynamicTracks:     true,
			wantMaxRecordings:     100,
			wantMaxRecordingJobs:  10,
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

			caps, err := client.GetRecordingServiceCapabilities(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecordingServiceCapabilities() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if caps == nil {
					t.Fatal("Expected capabilities, got nil")
				}

				if caps.DynamicRecordings != tt.wantDynamicRecordings {
					t.Errorf("DynamicRecordings = %v, want %v", caps.DynamicRecordings, tt.wantDynamicRecordings)
				}

				if caps.DynamicTracks != tt.wantDynamicTracks {
					t.Errorf("DynamicTracks = %v, want %v", caps.DynamicTracks, tt.wantDynamicTracks)
				}

				if caps.MaxRecordings != tt.wantMaxRecordings {
					t.Errorf("MaxRecordings = %d, want %d", caps.MaxRecordings, tt.wantMaxRecordings)
				}

				if caps.MaxRecordingJobs != tt.wantMaxRecordingJobs {
					t.Errorf("MaxRecordingJobs = %d, want %d", caps.MaxRecordingJobs, tt.wantMaxRecordingJobs)
				}
			}
		})
	}
}

func TestGetRecordings(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantErr    bool
		wantCount  int
		checkFirst func(t *testing.T, rec *Recording)
	}{
		{
			name: "successful recordings retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:GetRecordingsResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
							<trc:RecordingItem>
								<trc:RecordingToken>REC_001</trc:RecordingToken>
								<trc:Configuration>
									<tt:Source xmlns:tt="http://www.onvif.org/ver10/schema">
										<tt:SourceId>SRC_001</tt:SourceId>
										<tt:Name>Camera 1</tt:Name>
										<tt:Location>Building A</tt:Location>
										<tt:Description>Front entrance</tt:Description>
										<tt:Address>http://192.168.1.100/onvif/device_service</tt:Address>
									</tt:Source>
									<tt:Content xmlns:tt="http://www.onvif.org/ver10/schema">Recording</tt:Content>
									<tt:MaximumRetentionTime xmlns:tt="http://www.onvif.org/ver10/schema">PT72H</tt:MaximumRetentionTime>
								</trc:Configuration>
								<trc:Tracks>
									<trc:Track>
										<trc:TrackToken>VIDEO001</trc:TrackToken>
										<trc:Configuration>
											<tt:TrackType xmlns:tt="http://www.onvif.org/ver10/schema">Video</tt:TrackType>
											<tt:Description xmlns:tt="http://www.onvif.org/ver10/schema">Video track</tt:Description>
										</trc:Configuration>
									</trc:Track>
								</trc:Tracks>
							</trc:RecordingItem>
						</trc:GetRecordingsResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantCount: 1,
			checkFirst: func(t *testing.T, rec *Recording) {
				t.Helper()

				if rec.Token != "REC_001" {
					t.Errorf("Token = %q, want %q", rec.Token, "REC_001")
				}

				if rec.Configuration.Source.Name != "Camera 1" {
					t.Errorf("Source.Name = %q, want %q", rec.Configuration.Source.Name, "Camera 1")
				}

				if rec.Configuration.Source.SourceId != "SRC_001" {
					t.Errorf("Source.SourceId = %q, want %q", rec.Configuration.Source.SourceId, "SRC_001")
				}

				if rec.Configuration.Content != "Recording" {
					t.Errorf("Content = %q, want %q", rec.Configuration.Content, "Recording")
				}

				if rec.Configuration.MaximumRetentionTime != "PT72H" {
					t.Errorf("MaximumRetentionTime = %q, want %q", rec.Configuration.MaximumRetentionTime, "PT72H")
				}

				if len(rec.Tracks) != 1 {
					t.Fatalf("len(Tracks) = %d, want 1", len(rec.Tracks))
				}

				if rec.Tracks[0].Token != "VIDEO001" {
					t.Errorf("Track.Token = %q, want %q", rec.Tracks[0].Token, "VIDEO001")
				}

				if rec.Tracks[0].Configuration.TrackType != "Video" {
					t.Errorf("Track.TrackType = %q, want %q", rec.Tracks[0].Configuration.TrackType, "Video")
				}
			},
		},
		{
			name: "empty recordings list",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:GetRecordingsResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
						</trc:GetRecordingsResponse>
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

			recordings, err := client.GetRecordings(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecordings() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if len(recordings) != tt.wantCount {
					t.Errorf("len(recordings) = %d, want %d", len(recordings), tt.wantCount)
				}

				if tt.checkFirst != nil && len(recordings) > 0 {
					tt.checkFirst(t, recordings[0])
				}
			}
		})
	}
}

func TestCreateRecording(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		wantErr   bool
		wantToken string
	}{
		{
			name: "successful recording creation",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:CreateRecordingResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
							<trc:RecordingToken>REC_NEW_001</trc:RecordingToken>
						</trc:CreateRecordingResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantToken: "REC_NEW_001",
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Maximum recordings reached</s:Text></s:Reason>
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

			config := &RecordingConfiguration{
				Source: RecordingSourceInformation{
					SourceId: "SRC_001",
					Name:     "Camera 1",
				},
				Content:              "Recording",
				MaximumRetentionTime: "PT72H",
			}

			token, err := client.CreateRecording(context.Background(), config)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRecording() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if token != tt.wantToken {
					t.Errorf("RecordingToken = %q, want %q", token, tt.wantToken)
				}
			}
		})
	}
}

func TestDeleteRecording(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful recording deletion",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:DeleteRecordingResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl"/>
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
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Recording not found</s:Text></s:Reason>
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

			err = client.DeleteRecording(context.Background(), "REC_001")
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteRecording() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetRecordingConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful configuration update",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:SetRecordingConfigurationResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl"/>
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
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Invalid recording token</s:Text></s:Reason>
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

			config := &RecordingConfiguration{
				Source: RecordingSourceInformation{
					SourceId: "SRC_001",
					Name:     "Updated Camera",
				},
				Content:              "Updated Recording",
				MaximumRetentionTime: "PT48H",
			}

			err = client.SetRecordingConfiguration(context.Background(), "REC_001", config)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetRecordingConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetRecordingConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		wantErr     bool
		wantContent string
		wantSource  string
	}{
		{
			name: "successful configuration retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:GetRecordingConfigurationResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
							<trc:RecordingConfiguration>
								<tt:Source xmlns:tt="http://www.onvif.org/ver10/schema">
									<tt:SourceId>SRC_001</tt:SourceId>
									<tt:Name>Camera 1</tt:Name>
									<tt:Location>Building A</tt:Location>
									<tt:Description>Main entrance</tt:Description>
									<tt:Address>http://192.168.1.100/onvif/device_service</tt:Address>
								</tt:Source>
								<tt:Content xmlns:tt="http://www.onvif.org/ver10/schema">My Recording</tt:Content>
								<tt:MaximumRetentionTime xmlns:tt="http://www.onvif.org/ver10/schema">PT24H</tt:MaximumRetentionTime>
							</trc:RecordingConfiguration>
						</trc:GetRecordingConfigurationResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:     false,
			wantContent: "My Recording",
			wantSource:  "Camera 1",
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

			config, err := client.GetRecordingConfiguration(context.Background(), "REC_001")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecordingConfiguration() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if config == nil {
					t.Fatal("Expected configuration, got nil")
				}

				if config.Content != tt.wantContent {
					t.Errorf("Content = %q, want %q", config.Content, tt.wantContent)
				}

				if config.Source.Name != tt.wantSource {
					t.Errorf("Source.Name = %q, want %q", config.Source.Name, tt.wantSource)
				}
			}
		})
	}
}

func TestCreateTrack(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		wantErr   bool
		wantToken string
	}{
		{
			name: "successful track creation",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:CreateTrackResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
							<trc:TrackToken>TRACK_NEW</trc:TrackToken>
						</trc:CreateTrackResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantToken: "TRACK_NEW",
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Maximum tracks reached</s:Text></s:Reason>
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

			config := &TrackConfiguration{
				TrackType:   "Video",
				Description: "Video track",
			}

			token, err := client.CreateTrack(context.Background(), "REC1", config)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTrack() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if token != tt.wantToken {
					t.Errorf("TrackToken = %q, want %q", token, tt.wantToken)
				}
			}
		})
	}
}

func TestDeleteTrack(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful track deletion",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:DeleteTrackResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl"/>
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
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Track not found</s:Text></s:Reason>
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

			err = client.DeleteTrack(context.Background(), "REC1", "TRACK1")
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTrack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTrackConfiguration(t *testing.T) {
	tests := []struct {
		name          string
		handler       http.HandlerFunc
		wantErr       bool
		wantTrackType string
		wantDesc      string
	}{
		{
			name: "successful track configuration retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:GetTrackConfigurationResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
							<trc:TrackConfiguration>
								<tt:TrackType xmlns:tt="http://www.onvif.org/ver10/schema">Video</tt:TrackType>
								<tt:Description xmlns:tt="http://www.onvif.org/ver10/schema">Video track</tt:Description>
							</trc:TrackConfiguration>
						</trc:GetTrackConfigurationResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:       false,
			wantTrackType: "Video",
			wantDesc:      "Video track",
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Track not found</s:Text></s:Reason>
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

			config, err := client.GetTrackConfiguration(context.Background(), "REC1", "TRACK1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTrackConfiguration() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if config == nil {
					t.Fatal("Expected track configuration, got nil")
				}

				if config.TrackType != tt.wantTrackType {
					t.Errorf("TrackType = %q, want %q", config.TrackType, tt.wantTrackType)
				}

				if config.Description != tt.wantDesc {
					t.Errorf("Description = %q, want %q", config.Description, tt.wantDesc)
				}
			}
		})
	}
}

func TestSetTrackConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful track configuration update",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:SetTrackConfigurationResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl"/>
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
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Invalid track token</s:Text></s:Reason>
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

			config := &TrackConfiguration{
				TrackType:   "Audio",
				Description: "Updated audio track",
			}

			err = client.SetTrackConfiguration(context.Background(), "REC1", "TRACK1", config)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetTrackConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetRecordingOptions(t *testing.T) {
	spareTotal := 5
	spareVideo := 2
	spareAudio := 1
	spareMetadata := 3

	tests := []struct {
		name          string
		handler       http.HandlerFunc
		wantErr       bool
		wantTrack     bool
		wantSpareTotal *int
		wantSpareVideo *int
	}{
		{
			name: "successful options retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:GetRecordingOptionsResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
							<trc:Options>
								<trc:Track>
									<trc:SpareTotal>5</trc:SpareTotal>
									<trc:SpareVideo>2</trc:SpareVideo>
									<trc:SpareAudio>1</trc:SpareAudio>
									<trc:SpareMetadata>3</trc:SpareMetadata>
								</trc:Track>
							</trc:Options>
						</trc:GetRecordingOptionsResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:        false,
			wantTrack:      true,
			wantSpareTotal: &spareTotal,
			wantSpareVideo: &spareVideo,
		},
		{
			name: "options without track element",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:GetRecordingOptionsResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
							<trc:Options>
							</trc:Options>
						</trc:GetRecordingOptionsResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantTrack: false,
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
			wantErr: true,
		},
	}

	// Silence the unused variable warnings by referencing them
	_ = spareAudio
	_ = spareMetadata

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			options, err := client.GetRecordingOptions(context.Background(), "REC_001")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecordingOptions() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if options == nil {
					t.Fatal("Expected options, got nil")
				}

				if tt.wantTrack && options.Track == nil {
					t.Fatal("Expected Track options, got nil")
				}

				if !tt.wantTrack && options.Track != nil {
					t.Errorf("Expected no Track options, got %+v", options.Track)
				}

				if tt.wantTrack && options.Track != nil {
					if tt.wantSpareTotal != nil {
						if options.Track.SpareTotal == nil {
							t.Error("Expected SpareTotal, got nil")
						} else if *options.Track.SpareTotal != *tt.wantSpareTotal {
							t.Errorf("SpareTotal = %d, want %d", *options.Track.SpareTotal, *tt.wantSpareTotal)
						}
					}

					if tt.wantSpareVideo != nil {
						if options.Track.SpareVideo == nil {
							t.Error("Expected SpareVideo, got nil")
						} else if *options.Track.SpareVideo != *tt.wantSpareVideo {
							t.Errorf("SpareVideo = %d, want %d", *options.Track.SpareVideo, *tt.wantSpareVideo)
						}
					}
				}
			}
		})
	}
}

func TestCreateRecordingJob(t *testing.T) {
	tests := []struct {
		name         string
		handler      http.HandlerFunc
		wantErr      bool
		wantJobToken string
		wantMode     string
	}{
		{
			name: "successful job creation",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:CreateRecordingJobResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
							<trc:JobToken>JOB_NEW</trc:JobToken>
							<trc:JobConfiguration>
								<trc:RecordingToken>REC1</trc:RecordingToken>
								<trc:Mode>Active</trc:Mode>
								<trc:Priority>1</trc:Priority>
							</trc:JobConfiguration>
						</trc:CreateRecordingJobResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:      false,
			wantJobToken: "JOB_NEW",
			wantMode:     "Active",
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Maximum jobs reached</s:Text></s:Reason>
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

			config := &RecordingJobConfiguration{
				RecordingToken: "REC1",
				Mode:           "Active",
				Priority:       1,
			}

			jobToken, actualConfig, err := client.CreateRecordingJob(context.Background(), config)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRecordingJob() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if jobToken != tt.wantJobToken {
					t.Errorf("JobToken = %q, want %q", jobToken, tt.wantJobToken)
				}

				if actualConfig == nil {
					t.Fatal("Expected job configuration, got nil")
				}

				if actualConfig.Mode != tt.wantMode {
					t.Errorf("Mode = %q, want %q", actualConfig.Mode, tt.wantMode)
				}
			}
		})
	}
}

func TestDeleteRecordingJob(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful job deletion",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:DeleteRecordingJobResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl"/>
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
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Job not found</s:Text></s:Reason>
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

			err = client.DeleteRecordingJob(context.Background(), "JOB1")
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteRecordingJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetRecordingJobs(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		wantErr   bool
		wantCount int
		checkFirst func(t *testing.T, job *RecordingJob)
	}{
		{
			name: "successful jobs retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:GetRecordingJobsResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
							<trc:JobItem>
								<trc:JobToken>JOB1</trc:JobToken>
								<trc:JobConfiguration>
									<trc:RecordingToken>REC1</trc:RecordingToken>
									<trc:Mode>Active</trc:Mode>
									<trc:Priority>1</trc:Priority>
								</trc:JobConfiguration>
							</trc:JobItem>
						</trc:GetRecordingJobsResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:   false,
			wantCount: 1,
			checkFirst: func(t *testing.T, job *RecordingJob) {
				t.Helper()

				if job.Token != "JOB1" {
					t.Errorf("Token = %q, want %q", job.Token, "JOB1")
				}

				if job.Configuration.RecordingToken != "REC1" {
					t.Errorf("RecordingToken = %q, want %q", job.Configuration.RecordingToken, "REC1")
				}

				if job.Configuration.Mode != "Active" {
					t.Errorf("Mode = %q, want %q", job.Configuration.Mode, "Active")
				}

				if job.Configuration.Priority != 1 {
					t.Errorf("Priority = %d, want 1", job.Configuration.Priority)
				}
			},
		},
		{
			name: "empty jobs list",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:GetRecordingJobsResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
						</trc:GetRecordingJobsResponse>
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

			jobs, err := client.GetRecordingJobs(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecordingJobs() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if len(jobs) != tt.wantCount {
					t.Errorf("len(jobs) = %d, want %d", len(jobs), tt.wantCount)
				}

				if tt.checkFirst != nil && len(jobs) > 0 {
					tt.checkFirst(t, jobs[0])
				}
			}
		})
	}
}

func TestSetRecordingJobConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		handler  http.HandlerFunc
		wantErr  bool
		wantMode string
	}{
		{
			name: "successful job configuration update",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:SetRecordingJobConfigurationResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
							<trc:JobConfiguration>
								<trc:RecordingToken>REC1</trc:RecordingToken>
								<trc:Mode>Idle</trc:Mode>
								<trc:Priority>2</trc:Priority>
							</trc:JobConfiguration>
						</trc:SetRecordingJobConfigurationResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:  false,
			wantMode: "Idle",
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Job not found</s:Text></s:Reason>
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

			config := &RecordingJobConfiguration{
				RecordingToken: "REC1",
				Mode:           "Idle",
				Priority:       2,
			}

			actualConfig, err := client.SetRecordingJobConfiguration(context.Background(), "JOB1", config)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetRecordingJobConfiguration() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if actualConfig == nil {
					t.Fatal("Expected job configuration, got nil")
				}

				if actualConfig.Mode != tt.wantMode {
					t.Errorf("Mode = %q, want %q", actualConfig.Mode, tt.wantMode)
				}
			}
		})
	}
}

func TestGetRecordingJobConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		handler  http.HandlerFunc
		wantErr  bool
		wantMode string
		wantRec  string
	}{
		{
			name: "successful job configuration retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:GetRecordingJobConfigurationResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
							<trc:JobConfiguration>
								<trc:RecordingToken>REC1</trc:RecordingToken>
								<trc:Mode>Active</trc:Mode>
								<trc:Priority>1</trc:Priority>
							</trc:JobConfiguration>
						</trc:GetRecordingJobConfigurationResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:  false,
			wantMode: "Active",
			wantRec:  "REC1",
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Job not found</s:Text></s:Reason>
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

			config, err := client.GetRecordingJobConfiguration(context.Background(), "JOB1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecordingJobConfiguration() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if config == nil {
					t.Fatal("Expected job configuration, got nil")
				}

				if config.Mode != tt.wantMode {
					t.Errorf("Mode = %q, want %q", config.Mode, tt.wantMode)
				}

				if config.RecordingToken != tt.wantRec {
					t.Errorf("RecordingToken = %q, want %q", config.RecordingToken, tt.wantRec)
				}
			}
		})
	}
}

func TestSetRecordingJobMode(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful mode update",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:SetRecordingJobModeResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl"/>
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
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Job not found</s:Text></s:Reason>
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

			err = client.SetRecordingJobMode(context.Background(), "JOB1", "Idle")
			if (err != nil) != tt.wantErr {
				t.Errorf("SetRecordingJobMode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetRecordingJobState(t *testing.T) {
	tests := []struct {
		name          string
		handler       http.HandlerFunc
		wantErr       bool
		wantState     string
		wantRecToken  string
		wantSrcCount  int
	}{
		{
			name: "successful job state retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:GetRecordingJobStateResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
							<trc:State>
								<trc:RecordingToken>REC1</trc:RecordingToken>
								<trc:State>Active</trc:State>
								<trc:Sources>
									<trc:SourceToken>
										<trc:Token>SRC1</trc:Token>
										<trc:Type>http://www.onvif.org/ver10/schema/Profile</trc:Type>
									</trc:SourceToken>
									<trc:State>Recording</trc:State>
								</trc:Sources>
							</trc:State>
						</trc:GetRecordingJobStateResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:      false,
			wantState:    "Active",
			wantRecToken: "REC1",
			wantSrcCount: 1,
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Job not found</s:Text></s:Reason>
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

			state, err := client.GetRecordingJobState(context.Background(), "JOB1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecordingJobState() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if state == nil {
					t.Fatal("Expected job state, got nil")
				}

				if state.State != tt.wantState {
					t.Errorf("State = %q, want %q", state.State, tt.wantState)
				}

				if state.RecordingToken != tt.wantRecToken {
					t.Errorf("RecordingToken = %q, want %q", state.RecordingToken, tt.wantRecToken)
				}

				if len(state.Sources) != tt.wantSrcCount {
					t.Errorf("len(Sources) = %d, want %d", len(state.Sources), tt.wantSrcCount)
				}
			}
		})
	}
}

func TestExportRecordedData(t *testing.T) {
	tests := []struct {
		name          string
		handler       http.HandlerFunc
		wantErr       bool
		wantOpToken   string
		wantFileCount int
		wantFirstFile string
	}{
		{
			name: "successful export initiation",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:ExportRecordedDataResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
							<trc:OperationToken>EXPORT1</trc:OperationToken>
							<trc:FileNames>
								<trc:FileName>export1.zip</trc:FileName>
							</trc:FileNames>
						</trc:ExportRecordedDataResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:       false,
			wantOpToken:   "EXPORT1",
			wantFileCount: 1,
			wantFirstFile: "export1.zip",
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Export failed</s:Text></s:Reason>
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

			opToken, fileNames, err := client.ExportRecordedData(
				context.Background(),
				"2024-01-01T00:00:00Z",
				"2024-01-02T00:00:00Z",
				"REC1",
				"ONVIF",
				"ftp://backup/exports",
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExportRecordedData() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if opToken != tt.wantOpToken {
					t.Errorf("OperationToken = %q, want %q", opToken, tt.wantOpToken)
				}

				if len(fileNames) != tt.wantFileCount {
					t.Errorf("len(FileNames) = %d, want %d", len(fileNames), tt.wantFileCount)
				}

				if tt.wantFileCount > 0 && len(fileNames) > 0 && fileNames[0] != tt.wantFirstFile {
					t.Errorf("FileNames[0] = %q, want %q", fileNames[0], tt.wantFirstFile)
				}
			}
		})
	}
}

func TestStopExportRecordedData(t *testing.T) {
	tests := []struct {
		name         string
		handler      http.HandlerFunc
		wantErr      bool
		wantProgress float64
		wantFiles    int
	}{
		{
			name: "successful stop export",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:StopExportRecordedDataResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
							<trc:Progress>0.75</trc:Progress>
							<trc:FileProgressStatus>
								<trc:FileProgress>
									<trc:FileName>export1.zip</trc:FileName>
									<trc:Progress>0.75</trc:Progress>
								</trc:FileProgress>
							</trc:FileProgressStatus>
						</trc:StopExportRecordedDataResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:      false,
			wantProgress: 0.75,
			wantFiles:    1,
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Invalid operation token</s:Text></s:Reason>
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

			state, err := client.StopExportRecordedData(context.Background(), "EXPORT1")
			if (err != nil) != tt.wantErr {
				t.Errorf("StopExportRecordedData() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if state == nil {
					t.Fatal("Expected export state, got nil")
				}

				if state.Progress != tt.wantProgress {
					t.Errorf("Progress = %v, want %v", state.Progress, tt.wantProgress)
				}

				if len(state.FileProgressStatus) != tt.wantFiles {
					t.Errorf("len(FileProgressStatus) = %d, want %d", len(state.FileProgressStatus), tt.wantFiles)
				}
			}
		})
	}
}

func TestGetExportRecordedDataState(t *testing.T) {
	tests := []struct {
		name         string
		handler      http.HandlerFunc
		wantErr      bool
		wantProgress float64
		wantFiles    int
		wantFileName string
	}{
		{
			name: "successful state retrieval",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:GetExportRecordedDataStateResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl">
							<trc:Progress>0.5</trc:Progress>
							<trc:FileProgressStatus>
								<trc:FileProgress>
									<trc:FileName>export1.zip</trc:FileName>
									<trc:Progress>0.5</trc:Progress>
								</trc:FileProgress>
							</trc:FileProgressStatus>
						</trc:GetExportRecordedDataStateResponse>
					</s:Body>
				</s:Envelope>`
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			},
			wantErr:      false,
			wantProgress: 0.5,
			wantFiles:    1,
			wantFileName: "export1.zip",
		},
		{
			name: "SOAP fault response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<s:Fault>
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Operation not found</s:Text></s:Reason>
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

			state, err := client.GetExportRecordedDataState(context.Background(), "EXPORT1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExportRecordedDataState() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if state == nil {
					t.Fatal("Expected export state, got nil")
				}

				if state.Progress != tt.wantProgress {
					t.Errorf("Progress = %v, want %v", state.Progress, tt.wantProgress)
				}

				if len(state.FileProgressStatus) != tt.wantFiles {
					t.Errorf("len(FileProgressStatus) = %d, want %d", len(state.FileProgressStatus), tt.wantFiles)
				}

				if tt.wantFiles > 0 && len(state.FileProgressStatus) > 0 {
					if state.FileProgressStatus[0].FileName != tt.wantFileName {
						t.Errorf("FileProgressStatus[0].FileName = %q, want %q", state.FileProgressStatus[0].FileName, tt.wantFileName)
					}
				}
			}
		})
	}
}

func TestOverrideSegmentDuration(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "successful override",
			handler: func(w http.ResponseWriter, r *http.Request) {
				response := `<?xml version="1.0" encoding="UTF-8"?>
				<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
					<s:Body>
						<trc:OverrideSegmentDurationResponse xmlns:trc="http://www.onvif.org/ver10/recording/wsdl"/>
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
							<s:Code><s:Value>s:Receiver</s:Value></s:Code>
							<s:Reason><s:Text xml:lang="en">Invalid recording token</s:Text></s:Reason>
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

			err = client.OverrideSegmentDuration(context.Background(), "PT10M", "PT1H", "REC1")
			if (err != nil) != tt.wantErr {
				t.Errorf("OverrideSegmentDuration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
