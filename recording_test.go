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
