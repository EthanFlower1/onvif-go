//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestAnalytics_RealCamera(t *testing.T) {
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

	if !client.HasAnalyticsService() {
		t.Skip("Analytics service not available on this camera")
	}

	t.Run("GetAnalyticsServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetAnalyticsServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetAnalyticsServiceCapabilities not supported: %v", err)
		}

		t.Logf("Analytics service capabilities: RuleSupport=%v, AnalyticsModuleSupport=%v, CellBasedSceneDescriptionSupported=%v",
			caps.RuleSupport, caps.AnalyticsModuleSupport, caps.CellBasedSceneDescriptionSupported)
	})

	// getAnalyticsConfigToken returns a VideoAnalyticsConfiguration token from profiles,
	// or empty string if none is found. Many analytics operations require such a token.
	getAnalyticsConfigToken := func(t *testing.T) string {
		t.Helper()

		configs, err := client.GetVideoAnalyticsConfigurations(ctx)
		if err != nil || len(configs) == 0 {
			t.Skip("GetVideoAnalyticsConfigurations not supported or returned no configurations")
		}

		return configs[0].Token
	}

	t.Run("GetSupportedRules", func(t *testing.T) {
		configToken := getAnalyticsConfigToken(t)

		rules, err := client.GetSupportedRules(ctx, configToken)
		if err != nil {
			t.Skipf("GetSupportedRules not supported: %v", err)
		}

		t.Logf("Found %d supported rule types for config %s", len(rules), configToken)
		for _, r := range rules {
			t.Logf("  Rule: Name=%s, Parameters=%d", r.Name, len(r.Parameters))
		}
	})

	t.Run("GetRules", func(t *testing.T) {
		configToken := getAnalyticsConfigToken(t)

		rules, err := client.GetRules(ctx, configToken)
		if err != nil {
			t.Skipf("GetRules not supported: %v", err)
		}

		t.Logf("Found %d active rules for config %s", len(rules), configToken)
		for _, r := range rules {
			t.Logf("  Rule: Name=%s, Type=%s", r.Name, r.Type)
		}
	})

	t.Run("GetSupportedAnalyticsModules", func(t *testing.T) {
		configToken := getAnalyticsConfigToken(t)

		modules, err := client.GetSupportedAnalyticsModules(ctx, configToken)
		if err != nil {
			t.Skipf("GetSupportedAnalyticsModules not supported: %v", err)
		}

		t.Logf("Found %d supported analytics module types for config %s", len(modules), configToken)
		for _, m := range modules {
			t.Logf("  Module: Name=%s", m.Name)
		}
	})

	t.Run("GetAnalyticsModules", func(t *testing.T) {
		configToken := getAnalyticsConfigToken(t)

		modules, err := client.GetAnalyticsModules(ctx, configToken)
		if err != nil {
			t.Skipf("GetAnalyticsModules not supported: %v", err)
		}

		t.Logf("Found %d active analytics modules for config %s", len(modules), configToken)
		for _, m := range modules {
			t.Logf("  Module: Name=%s, Type=%s", m.Name, m.Type)
		}
	})
}
