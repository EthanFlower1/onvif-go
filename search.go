package onvif

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/0x524a/onvif-go/internal/soap"
)

// Search service namespace.
const searchNamespace = "http://www.onvif.org/ver10/search/wsdl"

// getSearchEndpoint returns the search service endpoint, falling back to the device endpoint.
func (c *Client) getSearchEndpoint() string {
	if c.searchEndpoint != "" {
		return c.searchEndpoint
	}

	return c.endpoint
}

// GetSearchServiceCapabilities retrieves the capabilities of the search service.
func (c *Client) GetSearchServiceCapabilities(ctx context.Context) (*SearchServiceCapabilities, error) {
	endpoint := c.getSearchEndpoint()

	type GetServiceCapabilities struct {
		XMLName xml.Name `xml:"tse:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:tse,attr"`
	}

	type GetServiceCapabilitiesResponse struct {
		XMLName      xml.Name `xml:"GetServiceCapabilitiesResponse"`
		Capabilities struct {
			MetadataSearch bool `xml:"MetadataSearch,attr"`
		} `xml:"Capabilities"`
	}

	req := GetServiceCapabilities{
		Xmlns: searchNamespace,
	}

	var resp GetServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetSearchServiceCapabilities failed: %w", err)
	}

	return &SearchServiceCapabilities{
		MetadataSearch: resp.Capabilities.MetadataSearch,
	}, nil
}

// GetRecordingSummary retrieves a summary of available recordings.
func (c *Client) GetRecordingSummary(ctx context.Context) (*RecordingSummary, error) {
	endpoint := c.getSearchEndpoint()

	type GetRecordingSummary struct {
		XMLName xml.Name `xml:"tse:GetRecordingSummary"`
		Xmlns   string   `xml:"xmlns:tse,attr"`
	}

	type GetRecordingSummaryResponse struct {
		XMLName xml.Name `xml:"GetRecordingSummaryResponse"`
		Summary struct {
			DataFrom         string `xml:"DataFrom"`
			DataUntil        string `xml:"DataUntil"`
			NumberRecordings int    `xml:"NumberRecordings"`
		} `xml:"Summary"`
	}

	req := GetRecordingSummary{
		Xmlns: searchNamespace,
	}

	var resp GetRecordingSummaryResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRecordingSummary failed: %w", err)
	}

	dataFrom, _ := time.Parse(time.RFC3339, resp.Summary.DataFrom)
	dataUntil, _ := time.Parse(time.RFC3339, resp.Summary.DataUntil)

	return &RecordingSummary{
		DataFrom:         dataFrom,
		DataUntil:        dataUntil,
		NumberRecordings: resp.Summary.NumberRecordings,
	}, nil
}

// GetRecordingInformation retrieves detailed information about a specific recording.
func (c *Client) GetRecordingInformation(ctx context.Context, recordingToken string) (*RecordingInformation, error) {
	endpoint := c.getSearchEndpoint()

	type GetRecordingInformation struct {
		XMLName        xml.Name `xml:"tse:GetRecordingInformation"`
		Xmlns          string   `xml:"xmlns:tse,attr"`
		RecordingToken string   `xml:"tse:RecordingToken"`
	}

	type TrackInfoEntry struct {
		TrackToken  string `xml:"TrackToken"`
		TrackType   string `xml:"TrackType"`
		Description string `xml:"Description"`
		DataFrom    string `xml:"DataFrom"`
		DataTo      string `xml:"DataTo"`
	}

	type SourceEntry struct {
		SourceId    string `xml:"SourceId"`
		Name        string `xml:"Name"`
		Location    string `xml:"Location"`
		Description string `xml:"Description"`
		Address     string `xml:"Address"`
	}

	type RecordingInfoEntry struct {
		RecordingToken    string           `xml:"RecordingToken"`
		Source            SourceEntry      `xml:"Source"`
		EarliestRecording string           `xml:"EarliestRecording"`
		LatestRecording   string           `xml:"LatestRecording"`
		Content           string           `xml:"Content"`
		Track             []TrackInfoEntry `xml:"Track"`
		RecordingStatus   string           `xml:"RecordingStatus"`
	}

	type GetRecordingInformationResponse struct {
		XMLName            xml.Name           `xml:"GetRecordingInformationResponse"`
		RecordingInformation RecordingInfoEntry `xml:"RecordingInformation"`
	}

	req := GetRecordingInformation{
		Xmlns:          searchNamespace,
		RecordingToken: recordingToken,
	}

	var resp GetRecordingInformationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRecordingInformation failed: %w", err)
	}

	info := &RecordingInformation{
		RecordingToken: resp.RecordingInformation.RecordingToken,
		Source: RecordingSourceInformation{
			SourceId:    resp.RecordingInformation.Source.SourceId,
			Name:        resp.RecordingInformation.Source.Name,
			Location:    resp.RecordingInformation.Source.Location,
			Description: resp.RecordingInformation.Source.Description,
			Address:     resp.RecordingInformation.Source.Address,
		},
		Content:         resp.RecordingInformation.Content,
		RecordingStatus: resp.RecordingInformation.RecordingStatus,
	}

	if t, err := time.Parse(time.RFC3339, resp.RecordingInformation.EarliestRecording); err == nil {
		info.EarliestRecording = &t
	}

	if t, err := time.Parse(time.RFC3339, resp.RecordingInformation.LatestRecording); err == nil {
		info.LatestRecording = &t
	}

	tracks := make([]*TrackInformation, 0, len(resp.RecordingInformation.Track))

	for i := range resp.RecordingInformation.Track {
		tr := &resp.RecordingInformation.Track[i]
		ti := &TrackInformation{
			TrackToken:  tr.TrackToken,
			TrackType:   tr.TrackType,
			Description: tr.Description,
		}
		ti.DataFrom, _ = time.Parse(time.RFC3339, tr.DataFrom)
		ti.DataTo, _ = time.Parse(time.RFC3339, tr.DataTo)
		tracks = append(tracks, ti)
	}

	info.TrackInformation = tracks

	return info, nil
}

// GetMediaAttributes retrieves media attributes for the given recordings at a specific time.
func (c *Client) GetMediaAttributes(ctx context.Context, recordingTokens []string, atTime time.Time) ([]*MediaAttributes, error) {
	endpoint := c.getSearchEndpoint()

	type GetMediaAttributes struct {
		XMLName         xml.Name `xml:"tse:GetMediaAttributes"`
		Xmlns           string   `xml:"xmlns:tse,attr"`
		RecordingTokens []string `xml:"tse:RecordingTokens"`
		Time            string   `xml:"tse:Time"`
	}

	type VideoAttrEntry struct {
		Bitrate   *int     `xml:"Bitrate"`
		Width     int      `xml:"Width"`
		Height    int      `xml:"Height"`
		Encoding  string   `xml:"Encoding"`
		Framerate *float64 `xml:"Framerate"`
	}

	type AudioAttrEntry struct {
		Bitrate    *int   `xml:"Bitrate"`
		Encoding   string `xml:"Encoding"`
		Samplerate int    `xml:"Samplerate"`
	}

	type TrackInfoEntry struct {
		TrackToken  string `xml:"TrackToken"`
		TrackType   string `xml:"TrackType"`
		Description string `xml:"Description"`
		DataFrom    string `xml:"DataFrom"`
		DataTo      string `xml:"DataTo"`
	}

	type TrackAttrEntry struct {
		TrackInformation TrackInfoEntry  `xml:"TrackInformation"`
		VideoAttributes  *VideoAttrEntry `xml:"VideoAttributes"`
		AudioAttributes  *AudioAttrEntry `xml:"AudioAttributes"`
	}

	type MediaAttrEntry struct {
		RecordingToken  string           `xml:"RecordingToken"`
		TrackAttributes []TrackAttrEntry `xml:"TrackAttributes"`
		From            string           `xml:"From"`
		Until           string           `xml:"Until"`
	}

	type GetMediaAttributesResponse struct {
		XMLName         xml.Name         `xml:"GetMediaAttributesResponse"`
		MediaAttributes []MediaAttrEntry `xml:"MediaAttributes"`
	}

	req := GetMediaAttributes{
		Xmlns:           searchNamespace,
		RecordingTokens: recordingTokens,
		Time:            atTime.UTC().Format(time.RFC3339),
	}

	var resp GetMediaAttributesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMediaAttributes failed: %w", err)
	}

	result := make([]*MediaAttributes, 0, len(resp.MediaAttributes))

	for i := range resp.MediaAttributes {
		entry := &resp.MediaAttributes[i]
		ma := &MediaAttributes{
			RecordingToken: entry.RecordingToken,
		}

		if t, err := time.Parse(time.RFC3339, entry.From); err == nil {
			ma.From = &t
		}

		if t, err := time.Parse(time.RFC3339, entry.Until); err == nil {
			ma.Until = &t
		}

		tracks := make([]*TrackAttributes, 0, len(entry.TrackAttributes))

		for j := range entry.TrackAttributes {
			ta := &entry.TrackAttributes[j]
			attr := &TrackAttributes{
				TrackInformation: &TrackInformation{
					TrackToken:  ta.TrackInformation.TrackToken,
					TrackType:   ta.TrackInformation.TrackType,
					Description: ta.TrackInformation.Description,
				},
			}
			attr.TrackInformation.DataFrom, _ = time.Parse(time.RFC3339, ta.TrackInformation.DataFrom)
			attr.TrackInformation.DataTo, _ = time.Parse(time.RFC3339, ta.TrackInformation.DataTo)

			if ta.VideoAttributes != nil {
				attr.VideoAttributes = &VideoAttributes{
					Bitrate:   ta.VideoAttributes.Bitrate,
					Width:     ta.VideoAttributes.Width,
					Height:    ta.VideoAttributes.Height,
					Encoding:  ta.VideoAttributes.Encoding,
					Framerate: ta.VideoAttributes.Framerate,
				}
			}

			if ta.AudioAttributes != nil {
				attr.AudioAttributes = &AudioAttributes{
					Bitrate:    ta.AudioAttributes.Bitrate,
					Encoding:   ta.AudioAttributes.Encoding,
					Samplerate: ta.AudioAttributes.Samplerate,
				}
			}

			tracks = append(tracks, attr)
		}

		ma.TrackAttributes = tracks
		result = append(result, ma)
	}

	return result, nil
}

// buildScopeXML is a helper to convert SearchScope into an inline XML struct.
func buildScopeXML(scope *SearchScope) *searchScopeXML {
	if scope == nil {
		return nil
	}

	s := &searchScopeXML{}

	for _, ref := range scope.IncludedSources {
		if ref != nil {
			s.IncludedSources = append(s.IncludedSources, sourceRefXML{Token: ref.Token, Type: ref.Type})
		}
	}

	s.IncludedRecordings = scope.IncludedRecordings
	s.RecordingInformationFilter = scope.RecordingInformationFilter

	return s
}

type sourceRefXML struct {
	Token string `xml:"tt:Token"`
	Type  string `xml:"tt:Type,omitempty"`
}

type searchScopeXML struct {
	IncludedSources            []sourceRefXML `xml:"tse:IncludedSources,omitempty"`
	IncludedRecordings         []string       `xml:"tt:IncludedRecordings,omitempty"`
	RecordingInformationFilter string         `xml:"tt:RecordingInformationFilter,omitempty"`
}

// FindRecordings starts a recording search and returns a SearchToken.
func (c *Client) FindRecordings(ctx context.Context, scope *SearchScope, maxMatches *int, keepAliveTime string) (string, error) {
	endpoint := c.getSearchEndpoint()

	type FindRecordings struct {
		XMLName       xml.Name        `xml:"tse:FindRecordings"`
		Xmlns         string          `xml:"xmlns:tse,attr"`
		XmlnsTt       string          `xml:"xmlns:tt,attr"`
		Scope         *searchScopeXML `xml:"tse:Scope,omitempty"`
		MaxMatches    *int            `xml:"tse:MaxMatches,omitempty"`
		KeepAliveTime string          `xml:"tse:KeepAliveTime,omitempty"`
	}

	type FindRecordingsResponse struct {
		XMLName     xml.Name `xml:"FindRecordingsResponse"`
		SearchToken string   `xml:"SearchToken"`
	}

	req := FindRecordings{
		Xmlns:         searchNamespace,
		XmlnsTt:       "http://www.onvif.org/ver10/schema",
		Scope:         buildScopeXML(scope),
		MaxMatches:    maxMatches,
		KeepAliveTime: keepAliveTime,
	}

	var resp FindRecordingsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("FindRecordings failed: %w", err)
	}

	return resp.SearchToken, nil
}

// GetRecordingSearchResults retrieves results from a previously started recording search.
func (c *Client) GetRecordingSearchResults(ctx context.Context, searchToken string, minResults, maxResults *int, waitTime string) (*FindRecordingResult, error) {
	endpoint := c.getSearchEndpoint()

	type GetRecordingSearchResults struct {
		XMLName     xml.Name `xml:"tse:GetRecordingSearchResults"`
		Xmlns       string   `xml:"xmlns:tse,attr"`
		SearchToken string   `xml:"tse:SearchToken"`
		MinResults  *int     `xml:"tse:MinResults,omitempty"`
		MaxResults  *int     `xml:"tse:MaxResults,omitempty"`
		WaitTime    string   `xml:"tse:WaitTime,omitempty"`
	}

	type TrackInfoEntry struct {
		TrackToken  string `xml:"TrackToken"`
		TrackType   string `xml:"TrackType"`
		Description string `xml:"Description"`
		DataFrom    string `xml:"DataFrom"`
		DataTo      string `xml:"DataTo"`
	}

	type SourceEntry struct {
		SourceId    string `xml:"SourceId"`
		Name        string `xml:"Name"`
		Location    string `xml:"Location"`
		Description string `xml:"Description"`
		Address     string `xml:"Address"`
	}

	type RecordingInfoEntry struct {
		RecordingToken    string           `xml:"RecordingToken"`
		Source            SourceEntry      `xml:"Source"`
		EarliestRecording string           `xml:"EarliestRecording"`
		LatestRecording   string           `xml:"LatestRecording"`
		Content           string           `xml:"Content"`
		Track             []TrackInfoEntry `xml:"Track"`
		RecordingStatus   string           `xml:"RecordingStatus"`
	}

	type ResultList struct {
		SearchState          string               `xml:"SearchState"`
		RecordingInformation []RecordingInfoEntry `xml:"RecordingInformation"`
	}

	type GetRecordingSearchResultsResponse struct {
		XMLName    xml.Name   `xml:"GetRecordingSearchResultsResponse"`
		ResultList ResultList `xml:"ResultList"`
	}

	req := GetRecordingSearchResults{
		Xmlns:       searchNamespace,
		SearchToken: searchToken,
		MinResults:  minResults,
		MaxResults:  maxResults,
		WaitTime:    waitTime,
	}

	var resp GetRecordingSearchResultsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRecordingSearchResults failed: %w", err)
	}

	result := &FindRecordingResult{
		SearchState: resp.ResultList.SearchState,
	}

	recInfos := make([]*RecordingInformation, 0, len(resp.ResultList.RecordingInformation))

	for i := range resp.ResultList.RecordingInformation {
		entry := &resp.ResultList.RecordingInformation[i]
		ri := &RecordingInformation{
			RecordingToken: entry.RecordingToken,
			Source: RecordingSourceInformation{
				SourceId:    entry.Source.SourceId,
				Name:        entry.Source.Name,
				Location:    entry.Source.Location,
				Description: entry.Source.Description,
				Address:     entry.Source.Address,
			},
			Content:         entry.Content,
			RecordingStatus: entry.RecordingStatus,
		}

		if t, err := time.Parse(time.RFC3339, entry.EarliestRecording); err == nil {
			ri.EarliestRecording = &t
		}

		if t, err := time.Parse(time.RFC3339, entry.LatestRecording); err == nil {
			ri.LatestRecording = &t
		}

		tracks := make([]*TrackInformation, 0, len(entry.Track))

		for j := range entry.Track {
			tr := &entry.Track[j]
			ti := &TrackInformation{
				TrackToken:  tr.TrackToken,
				TrackType:   tr.TrackType,
				Description: tr.Description,
			}
			ti.DataFrom, _ = time.Parse(time.RFC3339, tr.DataFrom)
			ti.DataTo, _ = time.Parse(time.RFC3339, tr.DataTo)
			tracks = append(tracks, ti)
		}

		ri.TrackInformation = tracks
		recInfos = append(recInfos, ri)
	}

	result.RecordingInformation = recInfos

	return result, nil
}

// FindEvents starts an event search and returns a SearchToken.
func (c *Client) FindEvents(ctx context.Context, startPoint time.Time, endPoint *time.Time, scope *SearchScope, includeStartState bool, maxMatches *int, keepAliveTime string) (string, error) {
	endpoint := c.getSearchEndpoint()

	type FindEvents struct {
		XMLName           xml.Name        `xml:"tse:FindEvents"`
		Xmlns             string          `xml:"xmlns:tse,attr"`
		XmlnsTt           string          `xml:"xmlns:tt,attr"`
		StartPoint        string          `xml:"tse:StartPoint"`
		EndPoint          string          `xml:"tse:EndPoint,omitempty"`
		Scope             *searchScopeXML `xml:"tse:Scope,omitempty"`
		IncludeStartState bool            `xml:"tse:IncludeStartState"`
		MaxMatches        *int            `xml:"tse:MaxMatches,omitempty"`
		KeepAliveTime     string          `xml:"tse:KeepAliveTime,omitempty"`
	}

	type FindEventsResponse struct {
		XMLName     xml.Name `xml:"FindEventsResponse"`
		SearchToken string   `xml:"SearchToken"`
	}

	req := FindEvents{
		Xmlns:             searchNamespace,
		XmlnsTt:           "http://www.onvif.org/ver10/schema",
		StartPoint:        startPoint.UTC().Format(time.RFC3339),
		Scope:             buildScopeXML(scope),
		IncludeStartState: includeStartState,
		MaxMatches:        maxMatches,
		KeepAliveTime:     keepAliveTime,
	}

	if endPoint != nil {
		req.EndPoint = endPoint.UTC().Format(time.RFC3339)
	}

	var resp FindEventsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("FindEvents failed: %w", err)
	}

	return resp.SearchToken, nil
}

// GetEventSearchResults retrieves results from a previously started event search.
func (c *Client) GetEventSearchResults(ctx context.Context, searchToken string, minResults, maxResults *int, waitTime string) (*FindEventResult, error) {
	endpoint := c.getSearchEndpoint()

	type GetEventSearchResults struct {
		XMLName     xml.Name `xml:"tse:GetEventSearchResults"`
		Xmlns       string   `xml:"xmlns:tse,attr"`
		SearchToken string   `xml:"tse:SearchToken"`
		MinResults  *int     `xml:"tse:MinResults,omitempty"`
		MaxResults  *int     `xml:"tse:MaxResults,omitempty"`
		WaitTime    string   `xml:"tse:WaitTime,omitempty"`
	}

	type EventEntry struct {
		RecordingToken  string `xml:"RecordingToken"`
		TrackToken      string `xml:"TrackToken"`
		Time            string `xml:"Time"`
		StartStateEvent bool   `xml:"StartStateEvent"`
	}

	type ResultList struct {
		SearchState string       `xml:"SearchState"`
		Result      []EventEntry `xml:"Result"`
	}

	type GetEventSearchResultsResponse struct {
		XMLName    xml.Name   `xml:"GetEventSearchResultsResponse"`
		ResultList ResultList `xml:"ResultList"`
	}

	req := GetEventSearchResults{
		Xmlns:       searchNamespace,
		SearchToken: searchToken,
		MinResults:  minResults,
		MaxResults:  maxResults,
		WaitTime:    waitTime,
	}

	var resp GetEventSearchResultsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetEventSearchResults failed: %w", err)
	}

	result := &FindEventResult{
		SearchState: resp.ResultList.SearchState,
	}

	events := make([]*EventResult, 0, len(resp.ResultList.Result))

	for i := range resp.ResultList.Result {
		e := &resp.ResultList.Result[i]
		er := &EventResult{
			RecordingToken:  e.RecordingToken,
			TrackToken:      e.TrackToken,
			StartStateEvent: e.StartStateEvent,
		}
		er.Time, _ = time.Parse(time.RFC3339, e.Time)
		events = append(events, er)
	}

	result.Events = events

	return result, nil
}

// FindPTZPosition starts a PTZ position search and returns a SearchToken.
func (c *Client) FindPTZPosition(ctx context.Context, startPoint time.Time, endPoint *time.Time, scope *SearchScope, maxMatches *int, keepAliveTime string) (string, error) {
	endpoint := c.getSearchEndpoint()

	type FindPTZPosition struct {
		XMLName       xml.Name        `xml:"tse:FindPTZPosition"`
		Xmlns         string          `xml:"xmlns:tse,attr"`
		XmlnsTt       string          `xml:"xmlns:tt,attr"`
		StartPoint    string          `xml:"tse:StartPoint"`
		EndPoint      string          `xml:"tse:EndPoint,omitempty"`
		Scope         *searchScopeXML `xml:"tse:Scope,omitempty"`
		MaxMatches    *int            `xml:"tse:MaxMatches,omitempty"`
		KeepAliveTime string          `xml:"tse:KeepAliveTime,omitempty"`
	}

	type FindPTZPositionResponse struct {
		XMLName     xml.Name `xml:"FindPTZPositionResponse"`
		SearchToken string   `xml:"SearchToken"`
	}

	req := FindPTZPosition{
		Xmlns:         searchNamespace,
		XmlnsTt:       "http://www.onvif.org/ver10/schema",
		StartPoint:    startPoint.UTC().Format(time.RFC3339),
		Scope:         buildScopeXML(scope),
		MaxMatches:    maxMatches,
		KeepAliveTime: keepAliveTime,
	}

	if endPoint != nil {
		req.EndPoint = endPoint.UTC().Format(time.RFC3339)
	}

	var resp FindPTZPositionResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("FindPTZPosition failed: %w", err)
	}

	return resp.SearchToken, nil
}

// GetPTZPositionSearchResults retrieves results from a previously started PTZ position search.
func (c *Client) GetPTZPositionSearchResults(ctx context.Context, searchToken string, minResults, maxResults *int, waitTime string) (*FindPTZPositionResult, error) {
	endpoint := c.getSearchEndpoint()

	type GetPTZPositionSearchResults struct {
		XMLName     xml.Name `xml:"tse:GetPTZPositionSearchResults"`
		Xmlns       string   `xml:"xmlns:tse,attr"`
		SearchToken string   `xml:"tse:SearchToken"`
		MinResults  *int     `xml:"tse:MinResults,omitempty"`
		MaxResults  *int     `xml:"tse:MaxResults,omitempty"`
		WaitTime    string   `xml:"tse:WaitTime,omitempty"`
	}

	type PanTiltEntry struct {
		X    float64 `xml:"x,attr"`
		Y    float64 `xml:"y,attr"`
		Space string  `xml:"space,attr,omitempty"`
	}

	type ZoomEntry struct {
		X    float64 `xml:"x,attr"`
		Space string  `xml:"space,attr,omitempty"`
	}

	type PTZVectorEntry struct {
		PanTilt *PanTiltEntry `xml:"PanTilt"`
		Zoom    *ZoomEntry    `xml:"Zoom"`
	}

	type PTZPositionEntry struct {
		RecordingToken string         `xml:"RecordingToken"`
		TrackToken     string         `xml:"TrackToken"`
		Time           string         `xml:"Time"`
		Position       PTZVectorEntry `xml:"Position"`
	}

	type ResultList struct {
		SearchState string             `xml:"SearchState"`
		Result      []PTZPositionEntry `xml:"Result"`
	}

	type GetPTZPositionSearchResultsResponse struct {
		XMLName    xml.Name   `xml:"GetPTZPositionSearchResultsResponse"`
		ResultList ResultList `xml:"ResultList"`
	}

	req := GetPTZPositionSearchResults{
		Xmlns:       searchNamespace,
		SearchToken: searchToken,
		MinResults:  minResults,
		MaxResults:  maxResults,
		WaitTime:    waitTime,
	}

	var resp GetPTZPositionSearchResultsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetPTZPositionSearchResults failed: %w", err)
	}

	result := &FindPTZPositionResult{
		SearchState: resp.ResultList.SearchState,
	}

	positions := make([]*PTZPositionResult, 0, len(resp.ResultList.Result))

	for i := range resp.ResultList.Result {
		e := &resp.ResultList.Result[i]
		pr := &PTZPositionResult{
			RecordingToken: e.RecordingToken,
			TrackToken:     e.TrackToken,
		}
		pr.Time, _ = time.Parse(time.RFC3339, e.Time)

		vec := &PTZVector{}

		if e.Position.PanTilt != nil {
			vec.PanTilt = &Vector2D{
				X:     e.Position.PanTilt.X,
				Y:     e.Position.PanTilt.Y,
				Space: e.Position.PanTilt.Space,
			}
		}

		if e.Position.Zoom != nil {
			vec.Zoom = &Vector1D{
				X:     e.Position.Zoom.X,
				Space: e.Position.Zoom.Space,
			}
		}

		pr.Position = vec
		positions = append(positions, pr)
	}

	result.Positions = positions

	return result, nil
}

// FindMetadata starts a metadata search and returns a SearchToken.
func (c *Client) FindMetadata(ctx context.Context, startPoint time.Time, endPoint *time.Time, scope *SearchScope, maxMatches *int, keepAliveTime string) (string, error) {
	endpoint := c.getSearchEndpoint()

	type FindMetadata struct {
		XMLName       xml.Name        `xml:"tse:FindMetadata"`
		Xmlns         string          `xml:"xmlns:tse,attr"`
		XmlnsTt       string          `xml:"xmlns:tt,attr"`
		StartPoint    string          `xml:"tse:StartPoint"`
		EndPoint      string          `xml:"tse:EndPoint,omitempty"`
		Scope         *searchScopeXML `xml:"tse:Scope,omitempty"`
		MaxMatches    *int            `xml:"tse:MaxMatches,omitempty"`
		KeepAliveTime string          `xml:"tse:KeepAliveTime,omitempty"`
	}

	type FindMetadataResponse struct {
		XMLName     xml.Name `xml:"FindMetadataResponse"`
		SearchToken string   `xml:"SearchToken"`
	}

	req := FindMetadata{
		Xmlns:         searchNamespace,
		XmlnsTt:       "http://www.onvif.org/ver10/schema",
		StartPoint:    startPoint.UTC().Format(time.RFC3339),
		Scope:         buildScopeXML(scope),
		MaxMatches:    maxMatches,
		KeepAliveTime: keepAliveTime,
	}

	if endPoint != nil {
		req.EndPoint = endPoint.UTC().Format(time.RFC3339)
	}

	var resp FindMetadataResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("FindMetadata failed: %w", err)
	}

	return resp.SearchToken, nil
}

// GetMetadataSearchResults retrieves results from a previously started metadata search.
func (c *Client) GetMetadataSearchResults(ctx context.Context, searchToken string, minResults, maxResults *int, waitTime string) (*FindMetadataResult, error) {
	endpoint := c.getSearchEndpoint()

	type GetMetadataSearchResults struct {
		XMLName     xml.Name `xml:"tse:GetMetadataSearchResults"`
		Xmlns       string   `xml:"xmlns:tse,attr"`
		SearchToken string   `xml:"tse:SearchToken"`
		MinResults  *int     `xml:"tse:MinResults,omitempty"`
		MaxResults  *int     `xml:"tse:MaxResults,omitempty"`
		WaitTime    string   `xml:"tse:WaitTime,omitempty"`
	}

	type MetadataEntry struct {
		RecordingToken string `xml:"RecordingToken"`
		TrackToken     string `xml:"TrackToken"`
		Time           string `xml:"Time"`
	}

	type ResultList struct {
		SearchState string          `xml:"SearchState"`
		Result      []MetadataEntry `xml:"Result"`
	}

	type GetMetadataSearchResultsResponse struct {
		XMLName    xml.Name   `xml:"GetMetadataSearchResultsResponse"`
		ResultList ResultList `xml:"ResultList"`
	}

	req := GetMetadataSearchResults{
		Xmlns:       searchNamespace,
		SearchToken: searchToken,
		MinResults:  minResults,
		MaxResults:  maxResults,
		WaitTime:    waitTime,
	}

	var resp GetMetadataSearchResultsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetMetadataSearchResults failed: %w", err)
	}

	result := &FindMetadataResult{
		SearchState: resp.ResultList.SearchState,
	}

	results := make([]*MetadataResult, 0, len(resp.ResultList.Result))

	for i := range resp.ResultList.Result {
		e := &resp.ResultList.Result[i]
		mr := &MetadataResult{
			RecordingToken: e.RecordingToken,
			TrackToken:     e.TrackToken,
		}
		mr.Time, _ = time.Parse(time.RFC3339, e.Time)
		results = append(results, mr)
	}

	result.Results = results

	return result, nil
}

// GetSearchState retrieves the current state of a search.
func (c *Client) GetSearchState(ctx context.Context, searchToken string) (string, error) {
	endpoint := c.getSearchEndpoint()

	type GetSearchState struct {
		XMLName     xml.Name `xml:"tse:GetSearchState"`
		Xmlns       string   `xml:"xmlns:tse,attr"`
		SearchToken string   `xml:"tse:SearchToken"`
	}

	type GetSearchStateResponse struct {
		XMLName     xml.Name `xml:"GetSearchStateResponse"`
		SearchState string   `xml:"SearchState"`
	}

	req := GetSearchState{
		Xmlns:       searchNamespace,
		SearchToken: searchToken,
	}

	var resp GetSearchStateResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("GetSearchState failed: %w", err)
	}

	return resp.SearchState, nil
}

// EndSearch terminates a search session and returns the final search token.
func (c *Client) EndSearch(ctx context.Context, searchToken string) (string, error) {
	endpoint := c.getSearchEndpoint()

	type EndSearch struct {
		XMLName     xml.Name `xml:"tse:EndSearch"`
		Xmlns       string   `xml:"xmlns:tse,attr"`
		SearchToken string   `xml:"tse:SearchToken"`
	}

	type EndSearchResponse struct {
		XMLName     xml.Name `xml:"EndSearchResponse"`
		SearchToken string   `xml:"SearchToken"`
	}

	req := EndSearch{
		Xmlns:       searchNamespace,
		SearchToken: searchToken,
	}

	var resp EndSearchResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("EndSearch failed: %w", err)
	}

	return resp.SearchToken, nil
}
