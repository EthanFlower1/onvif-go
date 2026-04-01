//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestAdvancedSecurity_RealCamera(t *testing.T) {
	endpoint := os.Getenv("ONVIF_ENDPOINT")
	username := os.Getenv("ONVIF_USERNAME")
	password := os.Getenv("ONVIF_PASSWORD")
	if endpoint == "" {
		t.Skip("ONVIF_ENDPOINT not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := NewClient(endpoint, WithCredentials(username, password))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	if err := client.Initialize(ctx); err != nil {
		t.Fatalf("Initialize: %v", err)
	}

	if !client.HasAdvancedSecurityService() {
		t.Skip("Advanced Security service not available on this camera")
	}

	t.Run("GetAdvSecServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetAdvancedSecurityServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetAdvancedSecurityServiceCapabilities not supported: %v", err)
		}

		t.Logf("Advanced Security capabilities: RSAKeyPairGeneration=%v, ECCKeyPairGeneration=%v, PKCS10=%v, SelfSignedCertificateCreation=%v",
			caps.KeystoreCapabilities.RSAKeyPairGeneration,
			caps.KeystoreCapabilities.ECCKeyPairGeneration,
			caps.KeystoreCapabilities.PKCS10,
			caps.KeystoreCapabilities.SelfSignedCertificateCreation)
	})

	t.Run("GetAllKeys", func(t *testing.T) {
		keys, err := client.GetAllKeys(ctx)
		if err != nil {
			t.Skipf("GetAllKeys not supported: %v", err)
		}

		t.Logf("Found %d keys", len(keys))
		for _, k := range keys {
			t.Logf("  Key: ID=%s, Status=%s", k.KeyID, k.KeyStatus)
		}
	})

	t.Run("GetAllCertificates", func(t *testing.T) {
		certs, err := client.GetAllCertificates(ctx)
		if err != nil {
			t.Skipf("GetAllCertificates not supported: %v", err)
		}

		t.Logf("Found %d certificates", len(certs))
		for _, c := range certs {
			t.Logf("  Certificate: ID=%s", c.CertificateID)
		}
	})

	t.Run("GetEnabledTLSVersions", func(t *testing.T) {
		versions, err := client.GetEnabledTLSVersions(ctx)
		if err != nil {
			t.Skipf("GetEnabledTLSVersions not supported: %v", err)
		}

		t.Logf("Enabled TLS versions: %v", versions)
	})
}
