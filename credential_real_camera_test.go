//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestCredential_RealCamera(t *testing.T) {
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

	if !client.HasCredentialService() {
		t.Skip("Credential service not available on this camera")
	}

	t.Run("GetCredentialServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetCredentialServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetCredentialServiceCapabilities not supported: %v", err)
		}

		t.Logf("Credential service capabilities: MaxLimit=%d, MaxCredentials=%d, MaxAccessProfilesPerCredential=%d, ClientSuppliedTokenSupported=%v",
			caps.MaxLimit, caps.MaxCredentials, caps.MaxAccessProfilesPerCredential, caps.ClientSuppliedTokenSupported)
		t.Logf("  SupportedIdentifierTypes=%v", caps.SupportedIdentifierTypes)
	})

	t.Run("GetSupportedFormatTypes", func(t *testing.T) {
		caps, err := client.GetCredentialServiceCapabilities(ctx)
		if err != nil || len(caps.SupportedIdentifierTypes) == 0 {
			t.Skip("GetCredentialServiceCapabilities not supported or returned no identifier types")
		}

		identifierTypeName := caps.SupportedIdentifierTypes[0]
		formats, err := client.GetSupportedFormatTypes(ctx, identifierTypeName)
		if err != nil {
			t.Skipf("GetSupportedFormatTypes not supported: %v", err)
		}

		t.Logf("Found %d format types for identifier type %s", len(formats), identifierTypeName)
		for _, f := range formats {
			t.Logf("  FormatType=%s", f.FormatType)
		}
	})

	t.Run("GetCredentialInfoList", func(t *testing.T) {
		infos, _, err := client.GetCredentialInfoList(ctx, nil, nil)
		if err != nil {
			t.Skipf("GetCredentialInfoList not supported: %v", err)
		}

		t.Logf("Found %d credentials", len(infos))
		for _, cred := range infos {
			t.Logf("  Credential: Token=%s, Description=%s", cred.Token, cred.Description)
		}
	})

	t.Run("GetCredentialState", func(t *testing.T) {
		infos, _, err := client.GetCredentialInfoList(ctx, nil, nil)
		if err != nil || len(infos) == 0 {
			t.Skip("GetCredentialInfoList not supported or returned no credentials")
		}

		token := infos[0].Token
		state, err := client.GetCredentialState(ctx, token)
		if err != nil {
			t.Skipf("GetCredentialState not supported: %v", err)
		}

		t.Logf("Credential state for token=%s: Enabled=%v, Reason=%s",
			token, state.Enabled, state.Reason)
	})
}
