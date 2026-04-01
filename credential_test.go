package onvif

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testCredentialXMLHeader = `<?xml version="1.0" encoding="UTF-8"?>`

func newMockCredentialServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml")

		body := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(body)
		bodyStr := string(body)

		var response string

		switch {
		case strings.Contains(bodyStr, "GetServiceCapabilities") && strings.Contains(bodyStr, "credential"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:GetServiceCapabilitiesResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl">
      <tcr:Capabilities MaxLimit="100" MaxCredentials="500" MaxAccessProfilesPerCredential="10"
        CredentialValiditySupported="true"
        CredentialAccessProfileValiditySupported="true"
        ValiditySupportsTimeValue="true"
        ResetAntipassbackSupported="true"
        ClientSuppliedTokenSupported="false"
        DefaultCredentialSuspensionDuration="PT5M"
        MaxWhitelistedItems="200"
        MaxBlacklistedItems="200">
        <tcr:SupportedIdentifierType>pt:Card</tcr:SupportedIdentifierType>
        <tcr:SupportedIdentifierType>pt:PIN</tcr:SupportedIdentifierType>
      </tcr:Capabilities>
    </tcr:GetServiceCapabilitiesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetSupportedFormatTypes"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:GetSupportedFormatTypesResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl">
      <tcr:FormatTypeInfo>
        <tcr:FormatType>WIEGAND-26</tcr:FormatType>
        <tcr:Description>Standard 26-bit Wiegand format</tcr:Description>
      </tcr:FormatTypeInfo>
      <tcr:FormatTypeInfo>
        <tcr:FormatType>WIEGAND-34</tcr:FormatType>
        <tcr:Description>34-bit Wiegand format</tcr:Description>
      </tcr:FormatTypeInfo>
    </tcr:GetSupportedFormatTypesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetCredentialInfoList"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:GetCredentialInfoListResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl">
      <tcr:NextStartReference>ref_002</tcr:NextStartReference>
      <tcr:CredentialInfo token="cred_001">
        <tcr:Description>John Doe's card</tcr:Description>
        <tcr:CredentialHolderReference>user001</tcr:CredentialHolderReference>
      </tcr:CredentialInfo>
      <tcr:CredentialInfo token="cred_002">
        <tcr:Description>Jane Smith's card</tcr:Description>
        <tcr:CredentialHolderReference>user002</tcr:CredentialHolderReference>
      </tcr:CredentialInfo>
    </tcr:GetCredentialInfoListResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetCredentialInfo"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:GetCredentialInfoResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl">
      <tcr:CredentialInfo token="cred_001">
        <tcr:Description>John Doe's card</tcr:Description>
        <tcr:CredentialHolderReference>user001</tcr:CredentialHolderReference>
      </tcr:CredentialInfo>
    </tcr:GetCredentialInfoResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetCredentialList"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:GetCredentialListResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl">
      <tcr:NextStartReference>ref_002</tcr:NextStartReference>
      <tcr:Credential token="cred_001">
        <tcr:Description>John Doe</tcr:Description>
        <tcr:CredentialHolderReference>user001</tcr:CredentialHolderReference>
        <tcr:CredentialIdentifier>
          <tcr:Type>
            <tcr:Name>pt:Card</tcr:Name>
            <tcr:FormatType>WIEGAND-26</tcr:FormatType>
          </tcr:Type>
          <tcr:ExemptedFromAuthentication>false</tcr:ExemptedFromAuthentication>
          <tcr:Value>AABBCCDD</tcr:Value>
        </tcr:CredentialIdentifier>
        <tcr:CredentialAccessProfile>
          <tcr:AccessProfileToken>ap_001</tcr:AccessProfileToken>
        </tcr:CredentialAccessProfile>
      </tcr:Credential>
    </tcr:GetCredentialListResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetCredentials"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:GetCredentialsResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl">
      <tcr:Credential token="cred_001">
        <tcr:Description>John Doe</tcr:Description>
        <tcr:CredentialHolderReference>user001</tcr:CredentialHolderReference>
        <tcr:CredentialIdentifier>
          <tcr:Type>
            <tcr:Name>pt:Card</tcr:Name>
            <tcr:FormatType>WIEGAND-26</tcr:FormatType>
          </tcr:Type>
          <tcr:ExemptedFromAuthentication>false</tcr:ExemptedFromAuthentication>
          <tcr:Value>AABBCCDD</tcr:Value>
        </tcr:CredentialIdentifier>
      </tcr:Credential>
    </tcr:GetCredentialsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "CreateCredential"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:CreateCredentialResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl">
      <tcr:Token>cred_new_001</tcr:Token>
    </tcr:CreateCredentialResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "ModifyCredential"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:ModifyCredentialResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetCredentialIdentifier"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:SetCredentialIdentifierResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetCredentialAccessProfiles"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:SetCredentialAccessProfilesResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetCredential"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:SetCredentialResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "DeleteCredentialIdentifier"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:DeleteCredentialIdentifierResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "DeleteCredentialAccessProfiles"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:DeleteCredentialAccessProfilesResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "DeleteCredential"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:DeleteCredentialResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetCredentialState"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:GetCredentialStateResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl">
      <tcr:State>
        <tcr:Enabled>true</tcr:Enabled>
        <tcr:Reason>pt:Enabled</tcr:Reason>
        <tcr:AntipassbackState>
          <tcr:AntipassbackViolated>false</tcr:AntipassbackViolated>
        </tcr:AntipassbackState>
      </tcr:State>
    </tcr:GetCredentialStateResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "EnableCredential"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:EnableCredentialResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "DisableCredential"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:DisableCredentialResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "ResetAntipassbackViolation"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:ResetAntipassbackViolationResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetCredentialIdentifiers"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:GetCredentialIdentifiersResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl">
      <tcr:CredentialIdentifier>
        <tcr:Type>
          <tcr:Name>pt:Card</tcr:Name>
          <tcr:FormatType>WIEGAND-26</tcr:FormatType>
        </tcr:Type>
        <tcr:ExemptedFromAuthentication>false</tcr:ExemptedFromAuthentication>
        <tcr:Value>AABBCCDD</tcr:Value>
      </tcr:CredentialIdentifier>
    </tcr:GetCredentialIdentifiersResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetCredentialAccessProfiles"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:GetCredentialAccessProfilesResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl">
      <tcr:CredentialAccessProfile>
        <tcr:AccessProfileToken>ap_001</tcr:AccessProfileToken>
      </tcr:CredentialAccessProfile>
      <tcr:CredentialAccessProfile>
        <tcr:AccessProfileToken>ap_002</tcr:AccessProfileToken>
      </tcr:CredentialAccessProfile>
    </tcr:GetCredentialAccessProfilesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetWhitelist"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:GetWhitelistResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl">
      <tcr:NextStartReference>wl_ref_002</tcr:NextStartReference>
      <tcr:Identifier>
        <tcr:Type>
          <tcr:Name>pt:Card</tcr:Name>
          <tcr:FormatType>WIEGAND-26</tcr:FormatType>
        </tcr:Type>
        <tcr:Value>AABBCCDD</tcr:Value>
      </tcr:Identifier>
    </tcr:GetWhitelistResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "AddToWhitelist"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:AddToWhitelistResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "RemoveFromWhitelist"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:RemoveFromWhitelistResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "DeleteWhitelist"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:DeleteWhitelistResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetBlacklist"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:GetBlacklistResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl">
      <tcr:NextStartReference>bl_ref_002</tcr:NextStartReference>
      <tcr:Identifier>
        <tcr:Type>
          <tcr:Name>pt:Card</tcr:Name>
          <tcr:FormatType>WIEGAND-26</tcr:FormatType>
        </tcr:Type>
        <tcr:Value>DEADBEEF</tcr:Value>
      </tcr:Identifier>
    </tcr:GetBlacklistResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "AddToBlacklist"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:AddToBlacklistResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "RemoveFromBlacklist"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:RemoveFromBlacklistResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "DeleteBlacklist"):
			response = testCredentialXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tcr:DeleteBlacklistResponse xmlns:tcr="http://www.onvif.org/ver10/credential/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		default:
			http.Error(w, "unknown operation", http.StatusBadRequest)

			return
		}

		_, _ = w.Write([]byte(response))
	}))
}

func newMockCredentialFaultServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(testCredentialXMLHeader + `
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
        <SOAP-ENV:Text xml:lang="en">Invalid credential token</SOAP-ENV:Text>
      </SOAP-ENV:Reason>
    </SOAP-ENV:Fault>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
	}))
}

// --- GetCredentialServiceCapabilities ---

func TestGetCredentialServiceCapabilities(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	caps, err := client.GetCredentialServiceCapabilities(context.Background())
	if err != nil {
		t.Fatalf("GetCredentialServiceCapabilities() error = %v", err)
	}

	if caps.MaxLimit != 100 {
		t.Errorf("MaxLimit = %d, want 100", caps.MaxLimit)
	}

	if caps.MaxCredentials != 500 {
		t.Errorf("MaxCredentials = %d, want 500", caps.MaxCredentials)
	}

	if caps.MaxAccessProfilesPerCredential != 10 {
		t.Errorf("MaxAccessProfilesPerCredential = %d, want 10", caps.MaxAccessProfilesPerCredential)
	}

	if !caps.CredentialValiditySupported {
		t.Error("CredentialValiditySupported = false, want true")
	}

	if !caps.ResetAntipassbackSupported {
		t.Error("ResetAntipassbackSupported = false, want true")
	}

	if caps.MaxWhitelistedItems != 200 {
		t.Errorf("MaxWhitelistedItems = %d, want 200", caps.MaxWhitelistedItems)
	}

	if len(caps.SupportedIdentifierTypes) != 2 {
		t.Errorf("len(SupportedIdentifierTypes) = %d, want 2", len(caps.SupportedIdentifierTypes))
	}
}

func TestGetCredentialServiceCapabilitiesFault(t *testing.T) {
	server := newMockCredentialFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetCredentialServiceCapabilities(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- GetSupportedFormatTypes ---

func TestGetSupportedFormatTypes(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	items, err := client.GetSupportedFormatTypes(context.Background(), "pt:Card")
	if err != nil {
		t.Fatalf("GetSupportedFormatTypes() error = %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}

	if items[0].FormatType != "WIEGAND-26" {
		t.Errorf("FormatType = %q, want %q", items[0].FormatType, "WIEGAND-26")
	}
}

// --- GetCredentialInfo ---

func TestGetCredentialInfo(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	items, err := client.GetCredentialInfo(context.Background(), []string{"cred_001"})
	if err != nil {
		t.Fatalf("GetCredentialInfo() error = %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}

	if items[0].Token != "cred_001" {
		t.Errorf("Token = %q, want %q", items[0].Token, "cred_001")
	}

	if items[0].CredentialHolderReference != "user001" {
		t.Errorf("CredentialHolderReference = %q, want %q", items[0].CredentialHolderReference, "user001")
	}
}

func TestGetCredentialInfoEmptyTokens(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetCredentialInfo(context.Background(), []string{})
	if !errors.Is(err, ErrInvalidCredentialToken) {
		t.Errorf("expected ErrInvalidCredentialToken, got %v", err)
	}
}

func TestGetCredentialInfoFault(t *testing.T) {
	server := newMockCredentialFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetCredentialInfo(context.Background(), []string{"cred_001"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- GetCredentialInfoList ---

func TestGetCredentialInfoList(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	items, nextRef, err := client.GetCredentialInfoList(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetCredentialInfoList() error = %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}

	if nextRef != "ref_002" {
		t.Errorf("NextStartReference = %q, want %q", nextRef, "ref_002")
	}

	if items[0].Token != "cred_001" {
		t.Errorf("Token = %q, want %q", items[0].Token, "cred_001")
	}
}

// --- GetCredentialsByTokens ---

func TestGetCredentialsByTokens(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	items, err := client.GetCredentialsByTokens(context.Background(), []string{"cred_001"})
	if err != nil {
		t.Fatalf("GetCredentialsByTokens() error = %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}

	if items[0].Token != "cred_001" {
		t.Errorf("Token = %q, want %q", items[0].Token, "cred_001")
	}

	if len(items[0].CredentialIdentifiers) != 1 {
		t.Fatalf("len(CredentialIdentifiers) = %d, want 1", len(items[0].CredentialIdentifiers))
	}

	if items[0].CredentialIdentifiers[0].Type.Name != "pt:Card" {
		t.Errorf("Type.Name = %q, want %q", items[0].CredentialIdentifiers[0].Type.Name, "pt:Card")
	}
}

func TestGetCredentialsByTokensEmptyTokens(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetCredentialsByTokens(context.Background(), []string{})
	if !errors.Is(err, ErrInvalidCredentialToken) {
		t.Errorf("expected ErrInvalidCredentialToken, got %v", err)
	}
}

// --- GetCredentialList ---

func TestGetCredentialList(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	items, nextRef, err := client.GetCredentialList(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetCredentialList() error = %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}

	if nextRef != "ref_002" {
		t.Errorf("NextStartReference = %q, want %q", nextRef, "ref_002")
	}

	if items[0].Token != "cred_001" {
		t.Errorf("Token = %q, want %q", items[0].Token, "cred_001")
	}

	if len(items[0].CredentialAccessProfiles) != 1 {
		t.Fatalf("len(CredentialAccessProfiles) = %d, want 1", len(items[0].CredentialAccessProfiles))
	}

	if items[0].CredentialAccessProfiles[0].AccessProfileToken != "ap_001" {
		t.Errorf("AccessProfileToken = %q, want %q", items[0].CredentialAccessProfiles[0].AccessProfileToken, "ap_001")
	}
}

// --- CreateCredential ---

func TestCreateCredential(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	cred := Credential{
		CredentialInfo: CredentialInfo{
			Description:              "Test credential",
			CredentialHolderReference: "user_test",
		},
		CredentialIdentifiers: []CredentialIdentifier{
			{
				Type:  CredentialIdentifierType{Name: "pt:Card", FormatType: "WIEGAND-26"},
				Value: []byte{0xAA, 0xBB},
			},
		},
	}

	state := CredentialState{Enabled: true}

	token, err := client.CreateCredential(context.Background(), cred, state)
	if err != nil {
		t.Fatalf("CreateCredential() error = %v", err)
	}

	if token != "cred_new_001" {
		t.Errorf("Token = %q, want %q", token, "cred_new_001")
	}
}

func TestCreateCredentialFault(t *testing.T) {
	server := newMockCredentialFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	cred := Credential{
		CredentialInfo: CredentialInfo{CredentialHolderReference: "user_test"},
		CredentialIdentifiers: []CredentialIdentifier{
			{Type: CredentialIdentifierType{Name: "pt:Card", FormatType: "WIEGAND-26"}, Value: []byte{0xAA}},
		},
	}

	_, err = client.CreateCredential(context.Background(), cred, CredentialState{Enabled: true})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- ModifyCredential ---

func TestModifyCredential(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	cred := Credential{
		CredentialInfo: CredentialInfo{
			Token:                    "cred_001",
			CredentialHolderReference: "user001",
		},
		CredentialIdentifiers: []CredentialIdentifier{
			{Type: CredentialIdentifierType{Name: "pt:Card", FormatType: "WIEGAND-26"}, Value: []byte{0xAA}},
		},
	}

	if err := client.ModifyCredential(context.Background(), cred); err != nil {
		t.Fatalf("ModifyCredential() error = %v", err)
	}
}

func TestModifyCredentialEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.ModifyCredential(context.Background(), Credential{})
	if !errors.Is(err, ErrInvalidCredentialToken) {
		t.Errorf("expected ErrInvalidCredentialToken, got %v", err)
	}
}

// --- DeleteCredential ---

func TestDeleteCredential(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if err := client.DeleteCredential(context.Background(), "cred_001"); err != nil {
		t.Fatalf("DeleteCredential() error = %v", err)
	}
}

func TestDeleteCredentialEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.DeleteCredential(context.Background(), "")
	if !errors.Is(err, ErrInvalidCredentialToken) {
		t.Errorf("expected ErrInvalidCredentialToken, got %v", err)
	}
}

func TestDeleteCredentialFault(t *testing.T) {
	server := newMockCredentialFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.DeleteCredential(context.Background(), "cred_001")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- GetCredentialState ---

func TestGetCredentialState(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	state, err := client.GetCredentialState(context.Background(), "cred_001")
	if err != nil {
		t.Fatalf("GetCredentialState() error = %v", err)
	}

	if !state.Enabled {
		t.Error("Enabled = false, want true")
	}

	if state.Reason != "pt:Enabled" {
		t.Errorf("Reason = %q, want %q", state.Reason, "pt:Enabled")
	}

	if state.AntipassbackState == nil {
		t.Fatal("AntipassbackState is nil, want non-nil")
	}

	if state.AntipassbackState.AntipassbackViolated {
		t.Error("AntipassbackViolated = true, want false")
	}
}

func TestGetCredentialStateEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetCredentialState(context.Background(), "")
	if !errors.Is(err, ErrInvalidCredentialToken) {
		t.Errorf("expected ErrInvalidCredentialToken, got %v", err)
	}
}

func TestGetCredentialStateFault(t *testing.T) {
	server := newMockCredentialFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetCredentialState(context.Background(), "cred_001")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- EnableCredential / DisableCredential ---

func TestEnableCredential(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if err := client.EnableCredential(context.Background(), "cred_001", nil); err != nil {
		t.Fatalf("EnableCredential() error = %v", err)
	}
}

func TestEnableCredentialEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.EnableCredential(context.Background(), "", nil)
	if !errors.Is(err, ErrInvalidCredentialToken) {
		t.Errorf("expected ErrInvalidCredentialToken, got %v", err)
	}
}

func TestDisableCredential(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	reason := "pt:Suspended"

	if err := client.DisableCredential(context.Background(), "cred_001", &reason); err != nil {
		t.Fatalf("DisableCredential() error = %v", err)
	}
}

func TestDisableCredentialEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.DisableCredential(context.Background(), "", nil)
	if !errors.Is(err, ErrInvalidCredentialToken) {
		t.Errorf("expected ErrInvalidCredentialToken, got %v", err)
	}
}

// --- ResetAntipassbackViolation ---

func TestResetAntipassbackViolation(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if err := client.ResetAntipassbackViolation(context.Background(), "cred_001"); err != nil {
		t.Fatalf("ResetAntipassbackViolation() error = %v", err)
	}
}

func TestResetAntipassbackViolationEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.ResetAntipassbackViolation(context.Background(), "")
	if !errors.Is(err, ErrInvalidCredentialToken) {
		t.Errorf("expected ErrInvalidCredentialToken, got %v", err)
	}
}

// --- GetCredentialIdentifiers ---

func TestGetCredentialIdentifiers(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ids, err := client.GetCredentialIdentifiers(context.Background(), "cred_001")
	if err != nil {
		t.Fatalf("GetCredentialIdentifiers() error = %v", err)
	}

	if len(ids) != 1 {
		t.Fatalf("len(ids) = %d, want 1", len(ids))
	}

	if ids[0].Type.Name != "pt:Card" {
		t.Errorf("Type.Name = %q, want %q", ids[0].Type.Name, "pt:Card")
	}

	if ids[0].ExemptedFromAuthentication {
		t.Error("ExemptedFromAuthentication = true, want false")
	}
}

func TestGetCredentialIdentifiersFault(t *testing.T) {
	server := newMockCredentialFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetCredentialIdentifiers(context.Background(), "cred_001")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- SetCredentialIdentifier ---

func TestSetCredentialIdentifier(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	id := CredentialIdentifier{
		Type:  CredentialIdentifierType{Name: "pt:Card", FormatType: "WIEGAND-26"},
		Value: []byte{0xAA, 0xBB, 0xCC, 0xDD},
	}

	if err := client.SetCredentialIdentifier(context.Background(), "cred_001", id); err != nil {
		t.Fatalf("SetCredentialIdentifier() error = %v", err)
	}
}

// --- DeleteCredentialIdentifier ---

func TestDeleteCredentialIdentifier(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if err := client.DeleteCredentialIdentifier(context.Background(), "cred_001", "pt:Card"); err != nil {
		t.Fatalf("DeleteCredentialIdentifier() error = %v", err)
	}
}

func TestDeleteCredentialIdentifierEmptyTypeName(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.DeleteCredentialIdentifier(context.Background(), "cred_001", "")
	if !errors.Is(err, ErrInvalidCredentialIdentifierTypeName) {
		t.Errorf("expected ErrInvalidCredentialIdentifierTypeName, got %v", err)
	}
}

// --- GetCredentialAccessProfiles ---

func TestGetCredentialAccessProfiles(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	aps, err := client.GetCredentialAccessProfiles(context.Background(), "cred_001")
	if err != nil {
		t.Fatalf("GetCredentialAccessProfiles() error = %v", err)
	}

	if len(aps) != 2 {
		t.Fatalf("len(aps) = %d, want 2", len(aps))
	}

	if aps[0].AccessProfileToken != "ap_001" {
		t.Errorf("AccessProfileToken = %q, want %q", aps[0].AccessProfileToken, "ap_001")
	}
}

func TestGetCredentialAccessProfilesFault(t *testing.T) {
	server := newMockCredentialFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetCredentialAccessProfiles(context.Background(), "cred_001")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- SetCredentialAccessProfiles ---

func TestSetCredentialAccessProfiles(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	aps := []CredentialAccessProfile{
		{AccessProfileToken: "ap_001"},
		{AccessProfileToken: "ap_002"},
	}

	if err := client.SetCredentialAccessProfiles(context.Background(), "cred_001", aps); err != nil {
		t.Fatalf("SetCredentialAccessProfiles() error = %v", err)
	}
}

// --- DeleteCredentialAccessProfiles ---

func TestDeleteCredentialAccessProfiles(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if err := client.DeleteCredentialAccessProfiles(context.Background(), "cred_001", []string{"ap_001"}); err != nil {
		t.Fatalf("DeleteCredentialAccessProfiles() error = %v", err)
	}
}

func TestDeleteCredentialAccessProfilesEmptyTokens(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.DeleteCredentialAccessProfiles(context.Background(), "cred_001", []string{})
	if !errors.Is(err, ErrInvalidAccessProfileToken) {
		t.Errorf("expected ErrInvalidAccessProfileToken, got %v", err)
	}
}

// --- Whitelist operations ---

func TestGetWhitelist(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	items, nextRef, err := client.GetWhitelist(context.Background(), nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("GetWhitelist() error = %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}

	if nextRef != "wl_ref_002" {
		t.Errorf("NextStartReference = %q, want %q", nextRef, "wl_ref_002")
	}

	if items[0].Type.Name != "pt:Card" {
		t.Errorf("Type.Name = %q, want %q", items[0].Type.Name, "pt:Card")
	}
}

func TestAddToWhitelist(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ids := []CredentialIdentifierItem{
		{Type: CredentialIdentifierType{Name: "pt:Card", FormatType: "WIEGAND-26"}, Value: []byte{0xAA, 0xBB}},
	}

	if err := client.AddToWhitelist(context.Background(), ids); err != nil {
		t.Fatalf("AddToWhitelist() error = %v", err)
	}
}

func TestDeleteWhitelist(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if err := client.DeleteWhitelist(context.Background()); err != nil {
		t.Fatalf("DeleteWhitelist() error = %v", err)
	}
}

// --- Blacklist operations ---

func TestGetBlacklist(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	items, nextRef, err := client.GetBlacklist(context.Background(), nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("GetBlacklist() error = %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}

	if nextRef != "bl_ref_002" {
		t.Errorf("NextStartReference = %q, want %q", nextRef, "bl_ref_002")
	}
}

func TestAddToBlacklist(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ids := []CredentialIdentifierItem{
		{Type: CredentialIdentifierType{Name: "pt:Card", FormatType: "WIEGAND-26"}, Value: []byte{0xDE, 0xAD}},
	}

	if err := client.AddToBlacklist(context.Background(), ids); err != nil {
		t.Fatalf("AddToBlacklist() error = %v", err)
	}
}

func TestDeleteBlacklist(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if err := client.DeleteBlacklist(context.Background()); err != nil {
		t.Fatalf("DeleteBlacklist() error = %v", err)
	}
}

// --- SetCredential ---

func TestSetCredential(t *testing.T) {
	server := newMockCredentialServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	data := CredentialData{
		Credential: Credential{
			CredentialInfo: CredentialInfo{
				Token:                    "cred_001",
				CredentialHolderReference: "user001",
			},
			CredentialIdentifiers: []CredentialIdentifier{
				{Type: CredentialIdentifierType{Name: "pt:Card", FormatType: "WIEGAND-26"}, Value: []byte{0xAA}},
			},
		},
		CredentialState: CredentialState{Enabled: true},
	}

	if err := client.SetCredential(context.Background(), data); err != nil {
		t.Fatalf("SetCredential() error = %v", err)
	}
}

func TestSetCredentialFault(t *testing.T) {
	server := newMockCredentialFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	data := CredentialData{
		Credential:      Credential{CredentialInfo: CredentialInfo{Token: "cred_001", CredentialHolderReference: "user001"}},
		CredentialState: CredentialState{Enabled: true},
	}

	err = client.SetCredential(context.Background(), data)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
