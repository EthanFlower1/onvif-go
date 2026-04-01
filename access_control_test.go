package onvif

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testAccessControlXMLHeader = `<?xml version="1.0" encoding="UTF-8"?>`

func newMockAccessControlServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml")

		body := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(body)
		bodyStr := string(body)

		var response string

		switch {
		case strings.Contains(bodyStr, "GetServiceCapabilities") && strings.Contains(bodyStr, "accesscontrol"):
			response = testAccessControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tac:GetServiceCapabilitiesResponse xmlns:tac="http://www.onvif.org/ver10/accesscontrol/wsdl">
      <tac:Capabilities MaxLimit="100" MaxAccessPoints="20" MaxAreas="50"
        ClientSuppliedTokenSupported="true"
        AccessPointManagementSupported="true"
        AreaManagementSupported="true"/>
    </tac:GetServiceCapabilitiesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAccessPointInfoList"):
			response = testAccessControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tac:GetAccessPointInfoListResponse xmlns:tac="http://www.onvif.org/ver10/accesscontrol/wsdl">
      <tac:NextStartReference>ref_002</tac:NextStartReference>
      <tac:AccessPointInfo token="ap_001">
        <tac:Name>Main Entrance</tac:Name>
        <tac:Description>Front door access point</tac:Description>
        <tac:AreaFrom>area_outside</tac:AreaFrom>
        <tac:AreaTo>area_lobby</tac:AreaTo>
        <tac:Entity>door_001</tac:Entity>
        <tac:Capabilities DisableAccessPoint="true"/>
      </tac:AccessPointInfo>
      <tac:AccessPointInfo token="ap_002">
        <tac:Name>Side Entrance</tac:Name>
        <tac:Entity>door_002</tac:Entity>
        <tac:Capabilities DisableAccessPoint="false"/>
      </tac:AccessPointInfo>
    </tac:GetAccessPointInfoListResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAccessPointInfo"):
			response = testAccessControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tac:GetAccessPointInfoResponse xmlns:tac="http://www.onvif.org/ver10/accesscontrol/wsdl">
      <tac:AccessPointInfo token="ap_001">
        <tac:Name>Main Entrance</tac:Name>
        <tac:Description>Front door access point</tac:Description>
        <tac:Entity>door_001</tac:Entity>
        <tac:Capabilities DisableAccessPoint="true"/>
      </tac:AccessPointInfo>
    </tac:GetAccessPointInfoResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAreaInfoList"):
			response = testAccessControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tac:GetAreaInfoListResponse xmlns:tac="http://www.onvif.org/ver10/accesscontrol/wsdl">
      <tac:AreaInfo token="area_001">
        <tac:Name>Lobby</tac:Name>
        <tac:Description>Main lobby area</tac:Description>
      </tac:AreaInfo>
      <tac:AreaInfo token="area_002">
        <tac:Name>Server Room</tac:Name>
        <tac:Description>Restricted server room</tac:Description>
      </tac:AreaInfo>
    </tac:GetAreaInfoListResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAreaInfo"):
			response = testAccessControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tac:GetAreaInfoResponse xmlns:tac="http://www.onvif.org/ver10/accesscontrol/wsdl">
      <tac:AreaInfo token="area_001">
        <tac:Name>Lobby</tac:Name>
        <tac:Description>Main lobby area</tac:Description>
      </tac:AreaInfo>
    </tac:GetAreaInfoResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAccessPointState"):
			response = testAccessControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tac:GetAccessPointStateResponse xmlns:tac="http://www.onvif.org/ver10/accesscontrol/wsdl">
      <tac:AccessPointState>
        <tac:Enabled>true</tac:Enabled>
      </tac:AccessPointState>
    </tac:GetAccessPointStateResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "EnableAccessPoint"):
			response = testAccessControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tac:EnableAccessPointResponse xmlns:tac="http://www.onvif.org/ver10/accesscontrol/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "DisableAccessPoint"):
			response = testAccessControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tac:DisableAccessPointResponse xmlns:tac="http://www.onvif.org/ver10/accesscontrol/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "ExternalAuthorization"):
			response = testAccessControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tac:ExternalAuthorizationResponse xmlns:tac="http://www.onvif.org/ver10/accesscontrol/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		default:
			http.Error(w, "unknown operation", http.StatusBadRequest)

			return
		}

		_, _ = w.Write([]byte(response))
	}))
}

func newMockAccessControlFaultServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(testAccessControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <SOAP-ENV:Fault>
      <SOAP-ENV:Code>
        <SOAP-ENV:Value>SOAP-ENV:Sender</SOAP-ENV:Value>
        <SOAP-ENV:Subcode>
          <SOAP-ENV:Value>ter:InvalidArgVal</SOAP-ENV:Value>
        </SOAP-ENV:Subcode>
      </SOAP-ENV:Code>
      <SOAP-ENV:Reason>
        <SOAP-ENV:Text xml:lang="en">Invalid token</SOAP-ENV:Text>
      </SOAP-ENV:Reason>
    </SOAP-ENV:Fault>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
	}))
}

func TestGetAccessControlServiceCapabilities(t *testing.T) {
	server := newMockAccessControlServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	caps, err := client.GetAccessControlServiceCapabilities(context.Background())
	if err != nil {
		t.Fatalf("GetAccessControlServiceCapabilities() error = %v", err)
	}

	if caps.MaxLimit != 100 {
		t.Errorf("MaxLimit = %d, want 100", caps.MaxLimit)
	}

	if caps.MaxAccessPoints != 20 {
		t.Errorf("MaxAccessPoints = %d, want 20", caps.MaxAccessPoints)
	}

	if caps.MaxAreas != 50 {
		t.Errorf("MaxAreas = %d, want 50", caps.MaxAreas)
	}

	if !caps.ClientSuppliedTokenSupported {
		t.Error("ClientSuppliedTokenSupported = false, want true")
	}

	if !caps.AccessPointManagementSupported {
		t.Error("AccessPointManagementSupported = false, want true")
	}

	if !caps.AreaManagementSupported {
		t.Error("AreaManagementSupported = false, want true")
	}
}

func TestGetAccessControlServiceCapabilitiesFault(t *testing.T) {
	server := newMockAccessControlFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetAccessControlServiceCapabilities(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetAccessPointInfoList(t *testing.T) {
	server := newMockAccessControlServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	items, nextRef, err := client.GetAccessPointInfoList(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetAccessPointInfoList() error = %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}

	if nextRef != "ref_002" {
		t.Errorf("NextStartReference = %q, want %q", nextRef, "ref_002")
	}

	if items[0].Token != "ap_001" {
		t.Errorf("items[0].Token = %q, want %q", items[0].Token, "ap_001")
	}

	if items[0].Name != "Main Entrance" {
		t.Errorf("items[0].Name = %q, want %q", items[0].Name, "Main Entrance")
	}

	if items[0].Description != "Front door access point" {
		t.Errorf("items[0].Description = %q, want %q", items[0].Description, "Front door access point")
	}

	if items[0].Entity != "door_001" {
		t.Errorf("items[0].Entity = %q, want %q", items[0].Entity, "door_001")
	}

	if !items[0].Capabilities.DisableAccessPoint {
		t.Error("items[0].Capabilities.DisableAccessPoint = false, want true")
	}
}

func TestGetAccessPointInfoListFault(t *testing.T) {
	server := newMockAccessControlFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, _, err = client.GetAccessPointInfoList(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetAccessPointInfo(t *testing.T) {
	server := newMockAccessControlServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	items, err := client.GetAccessPointInfo(context.Background(), []string{"ap_001"})
	if err != nil {
		t.Fatalf("GetAccessPointInfo() error = %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}

	if items[0].Token != "ap_001" {
		t.Errorf("Token = %q, want %q", items[0].Token, "ap_001")
	}
}

func TestGetAccessPointInfoEmptyTokens(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetAccessPointInfo(context.Background(), []string{})
	if !errors.Is(err, ErrInvalidAccessPointToken) {
		t.Errorf("expected ErrInvalidAccessPointToken, got %v", err)
	}
}

func TestGetAreaInfoList(t *testing.T) {
	server := newMockAccessControlServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	items, _, err := client.GetAreaInfoList(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetAreaInfoList() error = %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}

	if items[0].Token != "area_001" {
		t.Errorf("items[0].Token = %q, want %q", items[0].Token, "area_001")
	}

	if items[0].Name != "Lobby" {
		t.Errorf("items[0].Name = %q, want %q", items[0].Name, "Lobby")
	}

	if items[1].Token != "area_002" {
		t.Errorf("items[1].Token = %q, want %q", items[1].Token, "area_002")
	}
}

func TestGetAreaInfoListFault(t *testing.T) {
	server := newMockAccessControlFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, _, err = client.GetAreaInfoList(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetAreaInfo(t *testing.T) {
	server := newMockAccessControlServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	items, err := client.GetAreaInfo(context.Background(), []string{"area_001"})
	if err != nil {
		t.Fatalf("GetAreaInfo() error = %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}

	if items[0].Token != "area_001" {
		t.Errorf("Token = %q, want %q", items[0].Token, "area_001")
	}

	if items[0].Name != "Lobby" {
		t.Errorf("Name = %q, want %q", items[0].Name, "Lobby")
	}
}

func TestGetAreaInfoEmptyTokens(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetAreaInfo(context.Background(), []string{})
	if !errors.Is(err, ErrInvalidAreaToken) {
		t.Errorf("expected ErrInvalidAreaToken, got %v", err)
	}
}

func TestGetAccessPointState(t *testing.T) {
	server := newMockAccessControlServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	state, err := client.GetAccessPointState(context.Background(), "ap_001")
	if err != nil {
		t.Fatalf("GetAccessPointState() error = %v", err)
	}

	if !state.Enabled {
		t.Error("Enabled = false, want true")
	}
}

func TestGetAccessPointStateEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetAccessPointState(context.Background(), "")
	if !errors.Is(err, ErrInvalidAccessPointToken) {
		t.Errorf("expected ErrInvalidAccessPointToken, got %v", err)
	}
}

func TestGetAccessPointStateFault(t *testing.T) {
	server := newMockAccessControlFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetAccessPointState(context.Background(), "ap_001")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestEnableAccessPoint(t *testing.T) {
	server := newMockAccessControlServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.EnableAccessPoint(context.Background(), "ap_001")
	if err != nil {
		t.Fatalf("EnableAccessPoint() error = %v", err)
	}
}

func TestEnableAccessPointEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.EnableAccessPoint(context.Background(), "")
	if !errors.Is(err, ErrInvalidAccessPointToken) {
		t.Errorf("expected ErrInvalidAccessPointToken, got %v", err)
	}
}

func TestEnableAccessPointFault(t *testing.T) {
	server := newMockAccessControlFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.EnableAccessPoint(context.Background(), "ap_001")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDisableAccessPoint(t *testing.T) {
	server := newMockAccessControlServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.DisableAccessPoint(context.Background(), "ap_001")
	if err != nil {
		t.Fatalf("DisableAccessPoint() error = %v", err)
	}
}

func TestDisableAccessPointEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.DisableAccessPoint(context.Background(), "")
	if !errors.Is(err, ErrInvalidAccessPointToken) {
		t.Errorf("expected ErrInvalidAccessPointToken, got %v", err)
	}
}

func TestExternalAuthorization(t *testing.T) {
	server := newMockAccessControlServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.ExternalAuthorization(context.Background(), "ap_001", "cred_001", "authorized by operator", AccessDecisionGranted)
	if err != nil {
		t.Fatalf("ExternalAuthorization() error = %v", err)
	}
}

func TestExternalAuthorizationEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.ExternalAuthorization(context.Background(), "", "", "", AccessDecisionDenied)
	if !errors.Is(err, ErrInvalidAccessPointToken) {
		t.Errorf("expected ErrInvalidAccessPointToken, got %v", err)
	}
}

func TestExternalAuthorizationFault(t *testing.T) {
	server := newMockAccessControlFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.ExternalAuthorization(context.Background(), "ap_001", "", "", AccessDecisionGranted)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteAccessPoint(t *testing.T) {
	server := newMockAccessControlServer()
	defer server.Close()

	// Override the server to handle DeleteAccessPoint
	deleteServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml")
		_, _ = w.Write([]byte(testAccessControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tac:DeleteAccessPointResponse xmlns:tac="http://www.onvif.org/ver10/accesscontrol/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
	}))
	defer deleteServer.Close()

	client, err := NewClient(deleteServer.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.DeleteAccessPoint(context.Background(), "ap_001")
	if err != nil {
		t.Fatalf("DeleteAccessPoint() error = %v", err)
	}
}

func TestDeleteAccessPointEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.DeleteAccessPoint(context.Background(), "")
	if !errors.Is(err, ErrInvalidAccessPointToken) {
		t.Errorf("expected ErrInvalidAccessPointToken, got %v", err)
	}
}

func TestModifyAccessPointNil(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.ModifyAccessPoint(context.Background(), nil)
	if !errors.Is(err, ErrAccessPointNil) {
		t.Errorf("expected ErrAccessPointNil, got %v", err)
	}
}

func TestModifyAccessPointEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.ModifyAccessPoint(context.Background(), &AccessPoint{})
	if !errors.Is(err, ErrInvalidAccessPointToken) {
		t.Errorf("expected ErrInvalidAccessPointToken, got %v", err)
	}
}

func TestCreateAreaNil(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.CreateArea(context.Background(), nil)
	if !errors.Is(err, ErrAreaNil) {
		t.Errorf("expected ErrAreaNil, got %v", err)
	}
}

func TestModifyAreaEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.ModifyArea(context.Background(), &Area{})
	if !errors.Is(err, ErrInvalidAreaToken) {
		t.Errorf("expected ErrInvalidAreaToken, got %v", err)
	}
}

func TestDeleteAreaEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.DeleteArea(context.Background(), "")
	if !errors.Is(err, ErrInvalidAreaToken) {
		t.Errorf("expected ErrInvalidAreaToken, got %v", err)
	}
}
