package onvif

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testAdvancedSecurityXMLHeader = `<?xml version="1.0" encoding="UTF-8"?>`

const soapFaultAdvSec = testAdvancedSecurityXMLHeader + `
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

func newMockAdvancedSecurityServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml")

		body := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(body)
		bodyStr := string(body)

		var response string

		switch {
		// Key operations
		case strings.Contains(bodyStr, "CreateRSAKeyPair"):
			response = testAdvancedSecurityXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tas:CreateRSAKeyPairResponse xmlns:tas="http://www.onvif.org/ver10/advancedsecurity/wsdl">
      <tas:KeyID>key_rsa_001</tas:KeyID>
      <tas:EstimatedCreationTime>PT5S</tas:EstimatedCreationTime>
    </tas:CreateRSAKeyPairResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAllKeys"):
			response = testAdvancedSecurityXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tas:GetAllKeysResponse xmlns:tas="http://www.onvif.org/ver10/advancedsecurity/wsdl">
      <tas:KeyAttribute>
        <tas:KeyID>key_001</tas:KeyID>
        <tas:Alias>My RSA Key</tas:Alias>
        <tas:KeyStatus>ok</tas:KeyStatus>
      </tas:KeyAttribute>
      <tas:KeyAttribute>
        <tas:KeyID>key_002</tas:KeyID>
        <tas:KeyStatus>generating</tas:KeyStatus>
      </tas:KeyAttribute>
    </tas:GetAllKeysResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "DeleteKey"):
			response = testAdvancedSecurityXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tas:DeleteKeyResponse xmlns:tas="http://www.onvif.org/ver10/advancedsecurity/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		// Certificate operations
		case strings.Contains(bodyStr, "CreateSelfSignedCertificate"):
			response = testAdvancedSecurityXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tas:CreateSelfSignedCertificateResponse xmlns:tas="http://www.onvif.org/ver10/advancedsecurity/wsdl">
      <tas:CertificateID>cert_001</tas:CertificateID>
    </tas:CreateSelfSignedCertificateResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAllCertificates"):
			response = testAdvancedSecurityXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tas:GetAllCertificatesResponse xmlns:tas="http://www.onvif.org/ver10/advancedsecurity/wsdl">
      <tas:Certificate>
        <tas:CertificateID>cert_001</tas:CertificateID>
        <tas:KeyID>key_001</tas:KeyID>
        <tas:Alias>My Cert</tas:Alias>
        <tas:CertificateContent>AAAA</tas:CertificateContent>
      </tas:Certificate>
      <tas:Certificate>
        <tas:CertificateID>cert_002</tas:CertificateID>
        <tas:KeyID>key_002</tas:KeyID>
        <tas:CertificateContent>BBBB</tas:CertificateContent>
      </tas:Certificate>
    </tas:GetAllCertificatesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "DeleteCertificate"):
			response = testAdvancedSecurityXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tas:DeleteCertificateResponse xmlns:tas="http://www.onvif.org/ver10/advancedsecurity/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		// TLS operations
		case strings.Contains(bodyStr, "GetAssignedServerCertificates"):
			response = testAdvancedSecurityXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tas:GetAssignedServerCertificatesResponse xmlns:tas="http://www.onvif.org/ver10/advancedsecurity/wsdl">
      <tas:CertificationPathID>path_001</tas:CertificationPathID>
      <tas:CertificationPathID>path_002</tas:CertificationPathID>
    </tas:GetAssignedServerCertificatesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetEnabledTLSVersions"):
			response = testAdvancedSecurityXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tas:GetEnabledTLSVersionsResponse xmlns:tas="http://www.onvif.org/ver10/advancedsecurity/wsdl">
      <tas:Versions>TLS1.2 TLS1.3</tas:Versions>
    </tas:GetEnabledTLSVersionsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		// CRL operations
		case strings.Contains(bodyStr, "GetAllCRLs"):
			response = testAdvancedSecurityXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tas:GetAllCRLsResponse xmlns:tas="http://www.onvif.org/ver10/advancedsecurity/wsdl">
      <tas:Crl>
        <tas:CRLID>crl_001</tas:CRLID>
        <tas:Alias>My CRL</tas:Alias>
        <tas:CRLContent>CCCC</tas:CRLContent>
      </tas:Crl>
    </tas:GetAllCRLsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		// Dot1X operations
		case strings.Contains(bodyStr, "GetAllDot1XConfigurations"):
			response = testAdvancedSecurityXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tas:GetAllDot1XConfigurationsResponse xmlns:tas="http://www.onvif.org/ver10/advancedsecurity/wsdl">
      <tas:Configuration>
        <tas:Dot1XID>dot1x_001</tas:Dot1XID>
        <tas:Alias>Corp Network</tas:Alias>
        <tas:Outer Method="EAP-TLS"/>
      </tas:Configuration>
    </tas:GetAllDot1XConfigurationsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		// Passphrase operations
		case strings.Contains(bodyStr, "GetAllPassphrases"):
			response = testAdvancedSecurityXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tas:GetAllPassphrasesResponse xmlns:tas="http://www.onvif.org/ver10/advancedsecurity/wsdl">
      <tas:PassphraseAttribute>
        <tas:PassphraseID>pp_001</tas:PassphraseID>
        <tas:Alias>My Passphrase</tas:Alias>
      </tas:PassphraseAttribute>
      <tas:PassphraseAttribute>
        <tas:PassphraseID>pp_002</tas:PassphraseID>
      </tas:PassphraseAttribute>
    </tas:GetAllPassphrasesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		default:
			response = soapFaultAdvSec
		}

		_, _ = w.Write([]byte(response))
	}))
}

func newMockAdvancedSecurityFaultServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml")
		_, _ = w.Write([]byte(soapFaultAdvSec))
	}))
}

// ============================================================
// Key tests
// ============================================================

func TestCreateRSAKeyPair(t *testing.T) {
	server := newMockAdvancedSecurityServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.CreateRSAKeyPair(context.Background(), 2048, nil)
	if err != nil {
		t.Fatalf("CreateRSAKeyPair failed: %v", err)
	}

	if resp.KeyID != "key_rsa_001" {
		t.Errorf("expected KeyID %q, got %q", "key_rsa_001", resp.KeyID)
	}

	if resp.EstimatedCreationTime != "PT5S" {
		t.Errorf("expected EstimatedCreationTime %q, got %q", "PT5S", resp.EstimatedCreationTime)
	}
}

func TestCreateRSAKeyPairSOAPFault(t *testing.T) {
	server := newMockAdvancedSecurityFaultServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.CreateRSAKeyPair(context.Background(), 2048, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "CreateRSAKeyPair failed") {
		t.Errorf("expected error to contain 'CreateRSAKeyPair failed', got: %v", err)
	}
}

func TestGetAllKeys(t *testing.T) {
	server := newMockAdvancedSecurityServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	keys, err := client.GetAllKeys(context.Background())
	if err != nil {
		t.Fatalf("GetAllKeys failed: %v", err)
	}

	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}

	if keys[0].KeyID != "key_001" {
		t.Errorf("expected first KeyID %q, got %q", "key_001", keys[0].KeyID)
	}

	if keys[0].KeyStatus != "ok" {
		t.Errorf("expected KeyStatus %q, got %q", "ok", keys[0].KeyStatus)
	}

	if keys[0].Alias == nil || *keys[0].Alias != "My RSA Key" {
		t.Errorf("expected Alias %q, got %v", "My RSA Key", keys[0].Alias)
	}

	if keys[1].KeyID != "key_002" {
		t.Errorf("expected second KeyID %q, got %q", "key_002", keys[1].KeyID)
	}

	if keys[1].KeyStatus != "generating" {
		t.Errorf("expected second KeyStatus %q, got %q", "generating", keys[1].KeyStatus)
	}
}

func TestGetAllKeysSOAPFault(t *testing.T) {
	server := newMockAdvancedSecurityFaultServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.GetAllKeys(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "GetAllKeys failed") {
		t.Errorf("expected error to contain 'GetAllKeys failed', got: %v", err)
	}
}

func TestDeleteKey(t *testing.T) {
	server := newMockAdvancedSecurityServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	err = client.DeleteKey(context.Background(), "key_001")
	if err != nil {
		t.Fatalf("DeleteKey failed: %v", err)
	}
}

func TestDeleteKeySOAPFault(t *testing.T) {
	server := newMockAdvancedSecurityFaultServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	err = client.DeleteKey(context.Background(), "key_001")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "DeleteKey failed") {
		t.Errorf("expected error to contain 'DeleteKey failed', got: %v", err)
	}
}

// ============================================================
// Certificate tests
// ============================================================

func TestCreateSelfSignedCertificate(t *testing.T) {
	server := newMockAdvancedSecurityServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	subject := DistinguishedName{
		CommonName:   []string{"test.example.com"},
		Organization: []string{"Test Org"},
		Country:      []string{"US"},
	}

	sigAlg := AlgorithmIdentifier{
		Algorithm: "1.2.840.113549.1.1.5",
	}

	certID, err := client.CreateSelfSignedCertificate(
		context.Background(),
		subject,
		"key_001",
		nil,
		nil,
		nil,
		sigAlg,
		nil,
	)
	if err != nil {
		t.Fatalf("CreateSelfSignedCertificate failed: %v", err)
	}

	if certID != "cert_001" {
		t.Errorf("expected CertificateID %q, got %q", "cert_001", certID)
	}
}

func TestCreateSelfSignedCertificateSOAPFault(t *testing.T) {
	server := newMockAdvancedSecurityFaultServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	subject := DistinguishedName{CommonName: []string{"test.example.com"}}
	sigAlg := AlgorithmIdentifier{Algorithm: "1.2.840.113549.1.1.5"}

	_, err = client.CreateSelfSignedCertificate(context.Background(), subject, "key_001", nil, nil, nil, sigAlg, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "CreateSelfSignedCertificate failed") {
		t.Errorf("expected error to contain 'CreateSelfSignedCertificate failed', got: %v", err)
	}
}

func TestGetAllCertificates(t *testing.T) {
	server := newMockAdvancedSecurityServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	certs, err := client.GetAllCertificates(context.Background())
	if err != nil {
		t.Fatalf("GetAllCertificates failed: %v", err)
	}

	if len(certs) != 2 {
		t.Fatalf("expected 2 certificates, got %d", len(certs))
	}

	if certs[0].CertificateID != "cert_001" {
		t.Errorf("expected CertificateID %q, got %q", "cert_001", certs[0].CertificateID)
	}

	if certs[0].KeyID != "key_001" {
		t.Errorf("expected KeyID %q, got %q", "key_001", certs[0].KeyID)
	}

	if certs[0].Alias == nil || *certs[0].Alias != "My Cert" {
		t.Errorf("expected Alias %q, got %v", "My Cert", certs[0].Alias)
	}

	if certs[1].CertificateID != "cert_002" {
		t.Errorf("expected second CertificateID %q, got %q", "cert_002", certs[1].CertificateID)
	}
}

func TestGetAllCertificatesSOAPFault(t *testing.T) {
	server := newMockAdvancedSecurityFaultServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.GetAllCertificates(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "GetAllCertificates failed") {
		t.Errorf("expected error to contain 'GetAllCertificates failed', got: %v", err)
	}
}

func TestDeleteCertificate(t *testing.T) {
	server := newMockAdvancedSecurityServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	err = client.DeleteCertificate(context.Background(), "cert_001")
	if err != nil {
		t.Fatalf("DeleteCertificate failed: %v", err)
	}
}

func TestDeleteCertificateSOAPFault(t *testing.T) {
	server := newMockAdvancedSecurityFaultServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	err = client.DeleteCertificate(context.Background(), "cert_001")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "DeleteCertificate failed") {
		t.Errorf("expected error to contain 'DeleteCertificate failed', got: %v", err)
	}
}

// ============================================================
// TLS tests
// ============================================================

func TestGetAssignedServerCertificates(t *testing.T) {
	server := newMockAdvancedSecurityServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	paths, err := client.GetAssignedServerCertificates(context.Background())
	if err != nil {
		t.Fatalf("GetAssignedServerCertificates failed: %v", err)
	}

	if len(paths) != 2 {
		t.Fatalf("expected 2 certification paths, got %d", len(paths))
	}

	if paths[0] != "path_001" {
		t.Errorf("expected first path %q, got %q", "path_001", paths[0])
	}

	if paths[1] != "path_002" {
		t.Errorf("expected second path %q, got %q", "path_002", paths[1])
	}
}

func TestGetAssignedServerCertificatesSOAPFault(t *testing.T) {
	server := newMockAdvancedSecurityFaultServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.GetAssignedServerCertificates(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "GetAssignedServerCertificates failed") {
		t.Errorf("expected error to contain 'GetAssignedServerCertificates failed', got: %v", err)
	}
}

func TestGetEnabledTLSVersions(t *testing.T) {
	server := newMockAdvancedSecurityServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	versions, err := client.GetEnabledTLSVersions(context.Background())
	if err != nil {
		t.Fatalf("GetEnabledTLSVersions failed: %v", err)
	}

	if len(versions) != 2 {
		t.Fatalf("expected 2 TLS versions, got %d", len(versions))
	}

	if versions[0] != "TLS1.2" {
		t.Errorf("expected first version %q, got %q", "TLS1.2", versions[0])
	}

	if versions[1] != "TLS1.3" {
		t.Errorf("expected second version %q, got %q", "TLS1.3", versions[1])
	}
}

func TestGetEnabledTLSVersionsSOAPFault(t *testing.T) {
	server := newMockAdvancedSecurityFaultServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.GetEnabledTLSVersions(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "GetEnabledTLSVersions failed") {
		t.Errorf("expected error to contain 'GetEnabledTLSVersions failed', got: %v", err)
	}
}

// ============================================================
// CRL tests
// ============================================================

func TestGetAllCRLs(t *testing.T) {
	server := newMockAdvancedSecurityServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	crls, err := client.GetAllCRLs(context.Background())
	if err != nil {
		t.Fatalf("GetAllCRLs failed: %v", err)
	}

	if len(crls) != 1 {
		t.Fatalf("expected 1 CRL, got %d", len(crls))
	}

	if crls[0].CRLID != "crl_001" {
		t.Errorf("expected CRLID %q, got %q", "crl_001", crls[0].CRLID)
	}

	if crls[0].Alias != "My CRL" {
		t.Errorf("expected Alias %q, got %q", "My CRL", crls[0].Alias)
	}
}

func TestGetAllCRLsSOAPFault(t *testing.T) {
	server := newMockAdvancedSecurityFaultServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.GetAllCRLs(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "GetAllCRLs failed") {
		t.Errorf("expected error to contain 'GetAllCRLs failed', got: %v", err)
	}
}

// ============================================================
// Dot1X tests
// ============================================================

func TestGetAllAdvSecDot1XConfigurations(t *testing.T) {
	server := newMockAdvancedSecurityServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	configs, err := client.GetAllAdvSecDot1XConfigurations(context.Background())
	if err != nil {
		t.Fatalf("GetAllAdvSecDot1XConfigurations failed: %v", err)
	}

	if len(configs) != 1 {
		t.Fatalf("expected 1 configuration, got %d", len(configs))
	}

	if configs[0].Dot1XID != "dot1x_001" {
		t.Errorf("expected Dot1XID %q, got %q", "dot1x_001", configs[0].Dot1XID)
	}

	if configs[0].Alias == nil || *configs[0].Alias != "Corp Network" {
		t.Errorf("expected Alias %q, got %v", "Corp Network", configs[0].Alias)
	}

	if configs[0].Outer.Method != "EAP-TLS" {
		t.Errorf("expected Method %q, got %q", "EAP-TLS", configs[0].Outer.Method)
	}
}

func TestGetAllAdvSecDot1XConfigurationsSOAPFault(t *testing.T) {
	server := newMockAdvancedSecurityFaultServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.GetAllAdvSecDot1XConfigurations(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "GetAllAdvSecDot1XConfigurations failed") {
		t.Errorf("expected error to contain 'GetAllAdvSecDot1XConfigurations failed', got: %v", err)
	}
}

// ============================================================
// Passphrase tests
// ============================================================

func TestGetAllPassphrases(t *testing.T) {
	server := newMockAdvancedSecurityServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	passphrases, err := client.GetAllPassphrases(context.Background())
	if err != nil {
		t.Fatalf("GetAllPassphrases failed: %v", err)
	}

	if len(passphrases) != 2 {
		t.Fatalf("expected 2 passphrases, got %d", len(passphrases))
	}

	if passphrases[0].PassphraseID != "pp_001" {
		t.Errorf("expected PassphraseID %q, got %q", "pp_001", passphrases[0].PassphraseID)
	}

	if passphrases[0].Alias == nil || *passphrases[0].Alias != "My Passphrase" {
		t.Errorf("expected Alias %q, got %v", "My Passphrase", passphrases[0].Alias)
	}

	if passphrases[1].PassphraseID != "pp_002" {
		t.Errorf("expected second PassphraseID %q, got %q", "pp_002", passphrases[1].PassphraseID)
	}

	if passphrases[1].Alias != nil {
		t.Errorf("expected nil Alias for second passphrase, got %v", passphrases[1].Alias)
	}
}

func TestGetAllPassphrasesSOAPFault(t *testing.T) {
	server := newMockAdvancedSecurityFaultServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.GetAllPassphrases(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "GetAllPassphrases failed") {
		t.Errorf("expected error to contain 'GetAllPassphrases failed', got: %v", err)
	}
}

// TestAdvancedSecurityErrors verifies errors wrap properly.
func TestAdvancedSecurityErrors(t *testing.T) {
	server := newMockAdvancedSecurityFaultServer(t)
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Verify that errors are wrapped and unwrappable
	_, err = client.GetAllKeys(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// The error chain should contain the original ONVIF/SOAP error
	var onvifErr interface{ Error() string }
	if !errors.As(err, &onvifErr) {
		t.Error("expected error to be in chain")
	}
}
