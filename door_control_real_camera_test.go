//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestDoorControl_RealCamera(t *testing.T) {
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

	if !client.HasDoorControlService() {
		t.Skip("Door Control service not available on this camera")
	}

	t.Run("GetDoorControlServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetDoorControlServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetDoorControlServiceCapabilities not supported: %v", err)
		}

		t.Logf("Door Control service capabilities: MaxLimit=%d, MaxDoors=%d, ClientSuppliedTokenSupported=%v, DoorManagementSupported=%v",
			caps.MaxLimit, caps.MaxDoors, caps.ClientSuppliedTokenSupported, caps.DoorManagementSupported)
	})

	t.Run("GetDoorInfoList", func(t *testing.T) {
		infos, _, err := client.GetDoorInfoList(ctx, nil, nil)
		if err != nil {
			t.Skipf("GetDoorInfoList not supported: %v", err)
		}

		t.Logf("Found %d doors", len(infos))
		for _, door := range infos {
			t.Logf("  Door: Token=%s, Name=%s", door.Token, door.Name)
		}
	})

	t.Run("GetDoorState", func(t *testing.T) {
		infos, _, err := client.GetDoorInfoList(ctx, nil, nil)
		if err != nil || len(infos) == 0 {
			t.Skip("GetDoorInfoList not supported or returned no doors")
		}

		token := infos[0].Token
		state, err := client.GetDoorState(ctx, token)
		if err != nil {
			t.Skipf("GetDoorState not supported: %v", err)
		}

		t.Logf("Door state for token=%s: DoorMode=%s", token, state.DoorMode)
		if state.DoorPhysicalState != nil {
			t.Logf("  DoorPhysicalState=%s", *state.DoorPhysicalState)
		}
		if state.LockPhysicalState != nil {
			t.Logf("  LockPhysicalState=%s", *state.LockPhysicalState)
		}
	})
}
