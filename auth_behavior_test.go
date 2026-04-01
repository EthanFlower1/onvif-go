package onvif

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testAuthBehaviorXMLHeader = `<?xml version="1.0" encoding="UTF-8"?>`

func newMockAuthBehaviorServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml")

		body := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(body)
		bodyStr := string(body)

		var response string

		switch {
		case strings.Contains(bodyStr, "GetServiceCapabilities") && strings.Contains(bodyStr, "authenticationbehavior"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:GetServiceCapabilitiesResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl">
      <tab:Capabilities MaxLimit="100" MaxAuthenticationProfiles="20"
        MaxPoliciesPerAuthenticationProfile="5" MaxSecurityLevels="10"
        MaxRecognitionGroupsPerSecurityLevel="4" MaxRecognitionMethodsPerRecognitionGroup="3"
        ClientSuppliedTokenSupported="true" SupportedAuthenticationModes="pt:SingleCredential pt:DualCredential"/>
    </tab:GetServiceCapabilitiesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAuthenticationProfileInfoList"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:GetAuthenticationProfileInfoListResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl">
      <tab:NextStartReference>ref_002</tab:NextStartReference>
      <tab:AuthenticationProfileInfo token="ap_001">
        <tab:Name>Standard Profile</tab:Name>
        <tab:Description>Default authentication profile</tab:Description>
      </tab:AuthenticationProfileInfo>
      <tab:AuthenticationProfileInfo token="ap_002">
        <tab:Name>High Security Profile</tab:Name>
      </tab:AuthenticationProfileInfo>
    </tab:GetAuthenticationProfileInfoListResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAuthenticationProfileInfo"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:GetAuthenticationProfileInfoResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl">
      <tab:AuthenticationProfileInfo token="ap_001">
        <tab:Name>Standard Profile</tab:Name>
        <tab:Description>Default authentication profile</tab:Description>
      </tab:AuthenticationProfileInfo>
    </tab:GetAuthenticationProfileInfoResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAuthenticationProfileList"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:GetAuthenticationProfileListResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl">
      <tab:NextStartReference>ref_002</tab:NextStartReference>
      <tab:AuthenticationProfile token="ap_001">
        <tab:Name>Standard Profile</tab:Name>
        <tab:Description>Default authentication profile</tab:Description>
        <tab:DefaultSecurityLevelToken>sl_001</tab:DefaultSecurityLevelToken>
        <tab:AuthenticationPolicy>
          <tab:ScheduleToken>sched_001</tab:ScheduleToken>
          <tab:SecurityLevelConstraint>
            <tab:ActiveRegularSchedule>true</tab:ActiveRegularSchedule>
            <tab:ActiveSpecialDaySchedule>false</tab:ActiveSpecialDaySchedule>
            <tab:AuthenticationMode>pt:SingleCredential</tab:AuthenticationMode>
            <tab:SecurityLevelToken>sl_002</tab:SecurityLevelToken>
          </tab:SecurityLevelConstraint>
        </tab:AuthenticationPolicy>
      </tab:AuthenticationProfile>
    </tab:GetAuthenticationProfileListResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAuthenticationProfiles"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:GetAuthenticationProfilesResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl">
      <tab:AuthenticationProfile token="ap_001">
        <tab:Name>Standard Profile</tab:Name>
        <tab:DefaultSecurityLevelToken>sl_001</tab:DefaultSecurityLevelToken>
      </tab:AuthenticationProfile>
    </tab:GetAuthenticationProfilesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "CreateAuthenticationProfile"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:CreateAuthenticationProfileResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl">
      <tab:Token>ap_new_001</tab:Token>
    </tab:CreateAuthenticationProfileResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "ModifyAuthenticationProfile"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:ModifyAuthenticationProfileResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetAuthenticationProfile"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:SetAuthenticationProfileResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "DeleteAuthenticationProfile"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:DeleteAuthenticationProfileResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetSecurityLevelInfoList"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:GetSecurityLevelInfoListResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl">
      <tab:NextStartReference>ref_sl_002</tab:NextStartReference>
      <tab:SecurityLevelInfo token="sl_001">
        <tab:Name>Low Security</tab:Name>
        <tab:Priority>1</tab:Priority>
        <tab:Description>Single credential required</tab:Description>
      </tab:SecurityLevelInfo>
      <tab:SecurityLevelInfo token="sl_002">
        <tab:Name>High Security</tab:Name>
        <tab:Priority>10</tab:Priority>
      </tab:SecurityLevelInfo>
    </tab:GetSecurityLevelInfoListResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetSecurityLevelInfo"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:GetSecurityLevelInfoResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl">
      <tab:SecurityLevelInfo token="sl_001">
        <tab:Name>Low Security</tab:Name>
        <tab:Priority>1</tab:Priority>
        <tab:Description>Single credential required</tab:Description>
      </tab:SecurityLevelInfo>
    </tab:GetSecurityLevelInfoResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetSecurityLevelList"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:GetSecurityLevelListResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl">
      <tab:NextStartReference>ref_sl_002</tab:NextStartReference>
      <tab:SecurityLevel token="sl_001">
        <tab:Name>Low Security</tab:Name>
        <tab:Priority>1</tab:Priority>
        <tab:Description>Single credential required</tab:Description>
        <tab:RecognitionGroup>
          <tab:RecognitionMethod>
            <tab:RecognitionType>pt:CardNumber</tab:RecognitionType>
            <tab:Order>1</tab:Order>
          </tab:RecognitionMethod>
        </tab:RecognitionGroup>
      </tab:SecurityLevel>
    </tab:GetSecurityLevelListResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetSecurityLevels"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:GetSecurityLevelsResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl">
      <tab:SecurityLevel token="sl_001">
        <tab:Name>Low Security</tab:Name>
        <tab:Priority>1</tab:Priority>
        <tab:RecognitionGroup>
          <tab:RecognitionMethod>
            <tab:RecognitionType>pt:CardNumber</tab:RecognitionType>
            <tab:Order>1</tab:Order>
          </tab:RecognitionMethod>
        </tab:RecognitionGroup>
      </tab:SecurityLevel>
    </tab:GetSecurityLevelsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "CreateSecurityLevel"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:CreateSecurityLevelResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl">
      <tab:Token>sl_new_001</tab:Token>
    </tab:CreateSecurityLevelResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "ModifySecurityLevel"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:ModifySecurityLevelResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetSecurityLevel"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:SetSecurityLevelResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "DeleteSecurityLevel"):
			response = testAuthBehaviorXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tab:DeleteSecurityLevelResponse xmlns:tab="http://www.onvif.org/ver10/authenticationbehavior/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		default:
			http.Error(w, "unknown operation", http.StatusBadRequest)

			return
		}

		_, _ = w.Write([]byte(response))
	}))
}

func newMockAuthBehaviorFaultServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(testAuthBehaviorXMLHeader + `
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

func TestGetAuthBehaviorServiceCapabilities(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	caps, err := client.GetAuthBehaviorServiceCapabilities(context.Background())
	if err != nil {
		t.Fatalf("GetAuthBehaviorServiceCapabilities error: %v", err)
	}

	if caps.MaxLimit != 100 {
		t.Errorf("expected MaxLimit 100, got %d", caps.MaxLimit)
	}

	if caps.MaxAuthenticationProfiles != 20 {
		t.Errorf("expected MaxAuthenticationProfiles 20, got %d", caps.MaxAuthenticationProfiles)
	}

	if caps.MaxSecurityLevels != 10 {
		t.Errorf("expected MaxSecurityLevels 10, got %d", caps.MaxSecurityLevels)
	}

	if !caps.ClientSuppliedTokenSupported {
		t.Error("expected ClientSuppliedTokenSupported true")
	}
}

func TestGetAuthBehaviorServiceCapabilitiesFault(t *testing.T) {
	server := newMockAuthBehaviorFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.GetAuthBehaviorServiceCapabilities(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetAuthenticationProfileInfoList(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	infos, nextRef, err := client.GetAuthenticationProfileInfoList(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetAuthenticationProfileInfoList error: %v", err)
	}

	if len(infos) != 2 {
		t.Errorf("expected 2 items, got %d", len(infos))
	}

	if infos[0].Token != "ap_001" {
		t.Errorf("expected token ap_001, got %s", infos[0].Token)
	}

	if infos[0].Name != "Standard Profile" {
		t.Errorf("expected name Standard Profile, got %s", infos[0].Name)
	}

	if infos[0].Description != "Default authentication profile" {
		t.Errorf("expected description, got %s", infos[0].Description)
	}

	if nextRef != "ref_002" {
		t.Errorf("expected NextStartReference ref_002, got %s", nextRef)
	}
}

func TestGetAuthenticationProfileInfoListFault(t *testing.T) {
	server := newMockAuthBehaviorFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, _, err = client.GetAuthenticationProfileInfoList(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetAuthenticationProfileInfo(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	infos, err := client.GetAuthenticationProfileInfo(context.Background(), []string{"ap_001"})
	if err != nil {
		t.Fatalf("GetAuthenticationProfileInfo error: %v", err)
	}

	if len(infos) != 1 {
		t.Errorf("expected 1 item, got %d", len(infos))
	}

	if infos[0].Token != "ap_001" {
		t.Errorf("expected token ap_001, got %s", infos[0].Token)
	}
}

func TestGetAuthenticationProfileInfoEmptyTokens(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.GetAuthenticationProfileInfo(context.Background(), []string{})
	if !errors.Is(err, ErrInvalidAuthenticationProfileToken) {
		t.Errorf("expected ErrInvalidAuthenticationProfileToken, got %v", err)
	}
}

func TestGetAuthenticationProfileInfoFault(t *testing.T) {
	server := newMockAuthBehaviorFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.GetAuthenticationProfileInfo(context.Background(), []string{"ap_001"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetAuthenticationProfileList(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	profiles, nextRef, err := client.GetAuthenticationProfileList(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetAuthenticationProfileList error: %v", err)
	}

	if len(profiles) != 1 {
		t.Fatalf("expected 1 item, got %d", len(profiles))
	}

	if profiles[0].Token != "ap_001" {
		t.Errorf("expected token ap_001, got %s", profiles[0].Token)
	}

	if profiles[0].DefaultSecurityLevelToken != "sl_001" {
		t.Errorf("expected DefaultSecurityLevelToken sl_001, got %s", profiles[0].DefaultSecurityLevelToken)
	}

	if len(profiles[0].AuthenticationPolicies) != 1 {
		t.Fatalf("expected 1 policy, got %d", len(profiles[0].AuthenticationPolicies))
	}

	policy := profiles[0].AuthenticationPolicies[0]
	if policy.ScheduleToken != "sched_001" {
		t.Errorf("expected ScheduleToken sched_001, got %s", policy.ScheduleToken)
	}

	if len(policy.SecurityLevelConstraints) != 1 {
		t.Fatalf("expected 1 constraint, got %d", len(policy.SecurityLevelConstraints))
	}

	if !policy.SecurityLevelConstraints[0].ActiveRegularSchedule {
		t.Error("expected ActiveRegularSchedule true")
	}

	if nextRef != "ref_002" {
		t.Errorf("expected NextStartReference ref_002, got %s", nextRef)
	}
}

func TestGetAuthenticationProfileListFault(t *testing.T) {
	server := newMockAuthBehaviorFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, _, err = client.GetAuthenticationProfileList(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetAuthenticationProfiles(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	profiles, err := client.GetAuthenticationProfiles(context.Background(), []string{"ap_001"})
	if err != nil {
		t.Fatalf("GetAuthenticationProfiles error: %v", err)
	}

	if len(profiles) != 1 {
		t.Errorf("expected 1 item, got %d", len(profiles))
	}

	if profiles[0].Token != "ap_001" {
		t.Errorf("expected token ap_001, got %s", profiles[0].Token)
	}

	if profiles[0].DefaultSecurityLevelToken != "sl_001" {
		t.Errorf("expected DefaultSecurityLevelToken sl_001, got %s", profiles[0].DefaultSecurityLevelToken)
	}
}

func TestGetAuthenticationProfilesEmptyTokens(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.GetAuthenticationProfiles(context.Background(), []string{})
	if !errors.Is(err, ErrInvalidAuthenticationProfileToken) {
		t.Errorf("expected ErrInvalidAuthenticationProfileToken, got %v", err)
	}
}

func TestCreateAuthenticationProfile(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	profile := &AuthenticationProfile{
		AuthenticationProfileInfo: AuthenticationProfileInfo{
			Name:        "New Profile",
			Description: "Test profile",
		},
		DefaultSecurityLevelToken: "sl_001",
	}

	token, err := client.CreateAuthenticationProfile(context.Background(), profile)
	if err != nil {
		t.Fatalf("CreateAuthenticationProfile error: %v", err)
	}

	if token != "ap_new_001" {
		t.Errorf("expected token ap_new_001, got %s", token)
	}
}

func TestCreateAuthenticationProfileNil(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.CreateAuthenticationProfile(context.Background(), nil)
	if !errors.Is(err, ErrAuthenticationProfileNil) {
		t.Errorf("expected ErrAuthenticationProfileNil, got %v", err)
	}
}

func TestCreateAuthenticationProfileFault(t *testing.T) {
	server := newMockAuthBehaviorFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	profile := &AuthenticationProfile{
		AuthenticationProfileInfo: AuthenticationProfileInfo{
			Name: "New Profile",
		},
		DefaultSecurityLevelToken: "sl_001",
	}

	_, err = client.CreateAuthenticationProfile(context.Background(), profile)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestModifyAuthenticationProfile(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	profile := &AuthenticationProfile{
		AuthenticationProfileInfo: AuthenticationProfileInfo{
			Token: "ap_001",
			Name:  "Updated Profile",
		},
		DefaultSecurityLevelToken: "sl_001",
	}

	err = client.ModifyAuthenticationProfile(context.Background(), profile)
	if err != nil {
		t.Fatalf("ModifyAuthenticationProfile error: %v", err)
	}
}

func TestModifyAuthenticationProfileNil(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = client.ModifyAuthenticationProfile(context.Background(), nil)
	if !errors.Is(err, ErrAuthenticationProfileNil) {
		t.Errorf("expected ErrAuthenticationProfileNil, got %v", err)
	}
}

func TestModifyAuthenticationProfileEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	profile := &AuthenticationProfile{
		AuthenticationProfileInfo: AuthenticationProfileInfo{
			Name: "No Token",
		},
		DefaultSecurityLevelToken: "sl_001",
	}

	err = client.ModifyAuthenticationProfile(context.Background(), profile)
	if !errors.Is(err, ErrInvalidAuthenticationProfileToken) {
		t.Errorf("expected ErrInvalidAuthenticationProfileToken, got %v", err)
	}
}

func TestSetAuthenticationProfile(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	profile := &AuthenticationProfile{
		AuthenticationProfileInfo: AuthenticationProfileInfo{
			Token: "ap_001",
			Name:  "Set Profile",
		},
		DefaultSecurityLevelToken: "sl_001",
	}

	err = client.SetAuthenticationProfile(context.Background(), profile)
	if err != nil {
		t.Fatalf("SetAuthenticationProfile error: %v", err)
	}
}

func TestSetAuthenticationProfileNil(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = client.SetAuthenticationProfile(context.Background(), nil)
	if !errors.Is(err, ErrAuthenticationProfileNil) {
		t.Errorf("expected ErrAuthenticationProfileNil, got %v", err)
	}
}

func TestDeleteAuthenticationProfile(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = client.DeleteAuthenticationProfile(context.Background(), "ap_001")
	if err != nil {
		t.Fatalf("DeleteAuthenticationProfile error: %v", err)
	}
}

func TestDeleteAuthenticationProfileEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = client.DeleteAuthenticationProfile(context.Background(), "")
	if !errors.Is(err, ErrInvalidAuthenticationProfileToken) {
		t.Errorf("expected ErrInvalidAuthenticationProfileToken, got %v", err)
	}
}

func TestDeleteAuthenticationProfileFault(t *testing.T) {
	server := newMockAuthBehaviorFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = client.DeleteAuthenticationProfile(context.Background(), "ap_001")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetSecurityLevelInfoList(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	infos, nextRef, err := client.GetSecurityLevelInfoList(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetSecurityLevelInfoList error: %v", err)
	}

	if len(infos) != 2 {
		t.Errorf("expected 2 items, got %d", len(infos))
	}

	if infos[0].Token != "sl_001" {
		t.Errorf("expected token sl_001, got %s", infos[0].Token)
	}

	if infos[0].Name != "Low Security" {
		t.Errorf("expected name Low Security, got %s", infos[0].Name)
	}

	if infos[0].Priority != 1 {
		t.Errorf("expected priority 1, got %d", infos[0].Priority)
	}

	if nextRef != "ref_sl_002" {
		t.Errorf("expected NextStartReference ref_sl_002, got %s", nextRef)
	}
}

func TestGetSecurityLevelInfoListFault(t *testing.T) {
	server := newMockAuthBehaviorFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, _, err = client.GetSecurityLevelInfoList(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetSecurityLevelInfo(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	infos, err := client.GetSecurityLevelInfo(context.Background(), []string{"sl_001"})
	if err != nil {
		t.Fatalf("GetSecurityLevelInfo error: %v", err)
	}

	if len(infos) != 1 {
		t.Errorf("expected 1 item, got %d", len(infos))
	}

	if infos[0].Token != "sl_001" {
		t.Errorf("expected token sl_001, got %s", infos[0].Token)
	}

	if infos[0].Priority != 1 {
		t.Errorf("expected priority 1, got %d", infos[0].Priority)
	}
}

func TestGetSecurityLevelInfoEmptyTokens(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.GetSecurityLevelInfo(context.Background(), []string{})
	if !errors.Is(err, ErrInvalidSecurityLevelToken) {
		t.Errorf("expected ErrInvalidSecurityLevelToken, got %v", err)
	}
}

func TestGetSecurityLevelInfoFault(t *testing.T) {
	server := newMockAuthBehaviorFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.GetSecurityLevelInfo(context.Background(), []string{"sl_001"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetSecurityLevelList(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	levels, nextRef, err := client.GetSecurityLevelList(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetSecurityLevelList error: %v", err)
	}

	if len(levels) != 1 {
		t.Fatalf("expected 1 item, got %d", len(levels))
	}

	if levels[0].Token != "sl_001" {
		t.Errorf("expected token sl_001, got %s", levels[0].Token)
	}

	if levels[0].Priority != 1 {
		t.Errorf("expected priority 1, got %d", levels[0].Priority)
	}

	if len(levels[0].RecognitionGroups) != 1 {
		t.Fatalf("expected 1 recognition group, got %d", len(levels[0].RecognitionGroups))
	}

	if len(levels[0].RecognitionGroups[0].RecognitionMethods) != 1 {
		t.Fatalf("expected 1 recognition method, got %d", len(levels[0].RecognitionGroups[0].RecognitionMethods))
	}

	if levels[0].RecognitionGroups[0].RecognitionMethods[0].RecognitionType != "pt:CardNumber" {
		t.Errorf("expected RecognitionType pt:CardNumber, got %s", levels[0].RecognitionGroups[0].RecognitionMethods[0].RecognitionType)
	}

	if nextRef != "ref_sl_002" {
		t.Errorf("expected NextStartReference ref_sl_002, got %s", nextRef)
	}
}

func TestGetSecurityLevelListFault(t *testing.T) {
	server := newMockAuthBehaviorFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, _, err = client.GetSecurityLevelList(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetSecurityLevels(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	levels, err := client.GetSecurityLevels(context.Background(), []string{"sl_001"})
	if err != nil {
		t.Fatalf("GetSecurityLevels error: %v", err)
	}

	if len(levels) != 1 {
		t.Errorf("expected 1 item, got %d", len(levels))
	}

	if levels[0].Token != "sl_001" {
		t.Errorf("expected token sl_001, got %s", levels[0].Token)
	}
}

func TestGetSecurityLevelsEmptyTokens(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.GetSecurityLevels(context.Background(), []string{})
	if !errors.Is(err, ErrInvalidSecurityLevelToken) {
		t.Errorf("expected ErrInvalidSecurityLevelToken, got %v", err)
	}
}

func TestCreateSecurityLevel(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	sl := &SecurityLevel{
		SecurityLevelInfo: SecurityLevelInfo{
			Name:     "New Level",
			Priority: 5,
		},
	}

	token, err := client.CreateSecurityLevel(context.Background(), sl)
	if err != nil {
		t.Fatalf("CreateSecurityLevel error: %v", err)
	}

	if token != "sl_new_001" {
		t.Errorf("expected token sl_new_001, got %s", token)
	}
}

func TestCreateSecurityLevelNil(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.CreateSecurityLevel(context.Background(), nil)
	if !errors.Is(err, ErrSecurityLevelNil) {
		t.Errorf("expected ErrSecurityLevelNil, got %v", err)
	}
}

func TestCreateSecurityLevelFault(t *testing.T) {
	server := newMockAuthBehaviorFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	sl := &SecurityLevel{
		SecurityLevelInfo: SecurityLevelInfo{
			Name:     "New Level",
			Priority: 5,
		},
	}

	_, err = client.CreateSecurityLevel(context.Background(), sl)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestModifySecurityLevel(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	sl := &SecurityLevel{
		SecurityLevelInfo: SecurityLevelInfo{
			Token:    "sl_001",
			Name:     "Updated Level",
			Priority: 2,
		},
	}

	err = client.ModifySecurityLevel(context.Background(), sl)
	if err != nil {
		t.Fatalf("ModifySecurityLevel error: %v", err)
	}
}

func TestModifySecurityLevelNil(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = client.ModifySecurityLevel(context.Background(), nil)
	if !errors.Is(err, ErrSecurityLevelNil) {
		t.Errorf("expected ErrSecurityLevelNil, got %v", err)
	}
}

func TestModifySecurityLevelEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	sl := &SecurityLevel{
		SecurityLevelInfo: SecurityLevelInfo{
			Name:     "No Token",
			Priority: 1,
		},
	}

	err = client.ModifySecurityLevel(context.Background(), sl)
	if !errors.Is(err, ErrInvalidSecurityLevelToken) {
		t.Errorf("expected ErrInvalidSecurityLevelToken, got %v", err)
	}
}

func TestSetSecurityLevel(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	sl := &SecurityLevel{
		SecurityLevelInfo: SecurityLevelInfo{
			Token:    "sl_001",
			Name:     "Set Level",
			Priority: 5,
		},
	}

	err = client.SetSecurityLevel(context.Background(), sl)
	if err != nil {
		t.Fatalf("SetSecurityLevel error: %v", err)
	}
}

func TestSetSecurityLevelNil(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = client.SetSecurityLevel(context.Background(), nil)
	if !errors.Is(err, ErrSecurityLevelNil) {
		t.Errorf("expected ErrSecurityLevelNil, got %v", err)
	}
}

func TestDeleteSecurityLevel(t *testing.T) {
	server := newMockAuthBehaviorServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = client.DeleteSecurityLevel(context.Background(), "sl_001")
	if err != nil {
		t.Fatalf("DeleteSecurityLevel error: %v", err)
	}
}

func TestDeleteSecurityLevelEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = client.DeleteSecurityLevel(context.Background(), "")
	if !errors.Is(err, ErrInvalidSecurityLevelToken) {
		t.Errorf("expected ErrInvalidSecurityLevelToken, got %v", err)
	}
}

func TestDeleteSecurityLevelFault(t *testing.T) {
	server := newMockAuthBehaviorFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = client.DeleteSecurityLevel(context.Background(), "sl_001")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
