//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestAuthBehavior_RealCamera(t *testing.T) {
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

	if !client.HasAuthBehaviorService() {
		t.Skip("Authentication Behavior service not available on this camera")
	}

	t.Run("GetAuthBehaviorServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetAuthBehaviorServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetAuthBehaviorServiceCapabilities not supported: %v", err)
		}

		t.Logf("Auth Behavior service capabilities: MaxLimit=%d, MaxAuthenticationProfiles=%d, MaxSecurityLevels=%d, SupportedAuthenticationModes=%s",
			caps.MaxLimit, caps.MaxAuthenticationProfiles, caps.MaxSecurityLevels, caps.SupportedAuthenticationModes)
	})

	t.Run("GetAuthenticationProfileInfoList", func(t *testing.T) {
		infos, _, err := client.GetAuthenticationProfileInfoList(ctx, nil, nil)
		if err != nil {
			t.Skipf("GetAuthenticationProfileInfoList not supported: %v", err)
		}

		t.Logf("Found %d authentication profiles", len(infos))
		for _, ap := range infos {
			t.Logf("  AuthenticationProfile: Token=%s, Name=%s", ap.Token, ap.Name)
		}
	})

	t.Run("GetSecurityLevelInfoList", func(t *testing.T) {
		infos, _, err := client.GetSecurityLevelInfoList(ctx, nil, nil)
		if err != nil {
			t.Skipf("GetSecurityLevelInfoList not supported: %v", err)
		}

		t.Logf("Found %d security levels", len(infos))
		for _, sl := range infos {
			t.Logf("  SecurityLevel: Token=%s, Name=%s", sl.Token, sl.Name)
		}
	})
}
