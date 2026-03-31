//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestPTZ_RealCamera(t *testing.T) {
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

	t.Run("GetNodes", func(t *testing.T) {
		nodes, err := client.GetNodes(ctx)
		if err != nil {
			t.Skipf("GetNodes not supported: %v", err)
		}

		t.Logf("Found %d PTZ nodes", len(nodes))
		for _, n := range nodes {
			t.Logf("  Node: %s (%s)", n.Token, n.Name)
		}
	})

	t.Run("GetNode", func(t *testing.T) {
		nodes, err := client.GetNodes(ctx)
		if err != nil || len(nodes) == 0 {
			t.Skip("GetNodes not supported or returned no nodes")
		}

		nodeToken := nodes[0].Token
		node, err := client.GetNode(ctx, nodeToken)
		if err != nil {
			t.Skipf("GetNode not supported: %v", err)
		}

		if node.Token != nodeToken {
			t.Errorf("Expected node token %q, got %q", nodeToken, node.Token)
		}

		t.Logf("Node: Token=%s, Name=%s", node.Token, node.Name)
	})

	t.Run("GetConfigurations", func(t *testing.T) {
		configs, err := client.GetConfigurations(ctx)
		if err != nil {
			t.Skipf("GetConfigurations not supported: %v", err)
		}

		t.Logf("Found %d PTZ configurations", len(configs))
		for _, c := range configs {
			t.Logf("  Config: Token=%s, Name=%s", c.Token, c.Name)
		}
	})

	t.Run("GetPTZConfigurationOptions", func(t *testing.T) {
		configs, err := client.GetConfigurations(ctx)
		if err != nil || len(configs) == 0 {
			t.Skip("GetConfigurations not supported or returned no configurations")
		}

		configToken := configs[0].Token
		options, err := client.GetPTZConfigurationOptions(ctx, configToken)
		if err != nil {
			t.Skipf("GetPTZConfigurationOptions not supported: %v", err)
		}

		t.Logf("PTZ configuration options for token %s retrieved successfully", configToken)
		if options.PTZTimeout != nil {
			t.Logf("  PTZTimeout: Min=%s, Max=%s", options.PTZTimeout.Min, options.PTZTimeout.Max)
		}
	})

	t.Run("GetPTZServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetPTZServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetPTZServiceCapabilities not supported: %v", err)
		}

		t.Logf("PTZ service capabilities: EFlip=%v, Reverse=%v",
			caps.EFlip, caps.Reverse)
	})

	t.Run("GetCompatiblePTZConfigurationsForProfile", func(t *testing.T) {
		profiles, err := client.GetProfiles(ctx)
		if err != nil || len(profiles) == 0 {
			t.Skip("GetProfiles not supported or returned no profiles")
		}

		profileToken := profiles[0].Token
		configs, err := client.GetCompatiblePTZConfigurationsForProfile(ctx, profileToken)
		if err != nil {
			t.Skipf("GetCompatiblePTZConfigurationsForProfile not supported: %v", err)
		}

		t.Logf("Found %d compatible PTZ configurations for profile %s", len(configs), profileToken)
		for _, c := range configs {
			t.Logf("  Config: Token=%s, Name=%s", c.Token, c.Name)
		}
	})

	t.Run("GetPresetTours", func(t *testing.T) {
		profiles, err := client.GetProfiles(ctx)
		if err != nil || len(profiles) == 0 {
			t.Skip("GetProfiles not supported or returned no profiles")
		}

		profileToken := profiles[0].Token
		tours, err := client.GetPresetTours(ctx, profileToken)
		if err != nil {
			t.Skipf("GetPresetTours not supported: %v", err)
		}

		t.Logf("Found %d preset tours for profile %s", len(tours), profileToken)
		for _, tour := range tours {
			t.Logf("  Tour: Token=%s, Name=%s", tour.Token, tour.Name)
		}
	})

	// Write operations skipped: SetPTZConfiguration, CreatePresetTour
	// Side-effect operations skipped: PTZSendAuxiliaryCommand, GeoMove
}
