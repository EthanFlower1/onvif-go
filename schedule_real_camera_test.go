//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestSchedule_RealCamera(t *testing.T) {
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

	if !client.HasScheduleService() {
		t.Skip("Schedule service not available on this camera")
	}

	t.Run("GetScheduleServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetScheduleServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetScheduleServiceCapabilities not supported: %v", err)
		}

		t.Logf("Schedule service capabilities: MaxLimit=%d, MaxSchedules=%d, MaxSpecialDayGroups=%d, SpecialDaysSupported=%v, StateReportingSupported=%v",
			caps.MaxLimit, caps.MaxSchedules, caps.MaxSpecialDayGroups,
			caps.SpecialDaysSupported, caps.StateReportingSupported)
	})

	t.Run("GetScheduleInfoList", func(t *testing.T) {
		infos, _, err := client.GetScheduleInfoList(ctx, nil, nil)
		if err != nil {
			t.Skipf("GetScheduleInfoList not supported: %v", err)
		}

		t.Logf("Found %d schedules", len(infos))
		for _, s := range infos {
			t.Logf("  Schedule: Token=%s, Name=%s", s.Token, s.Name)
		}
	})

	t.Run("GetScheduleState", func(t *testing.T) {
		infos, _, err := client.GetScheduleInfoList(ctx, nil, nil)
		if err != nil || len(infos) == 0 {
			t.Skip("GetScheduleInfoList not supported or returned no schedules")
		}

		token := infos[0].Token
		state, err := client.GetScheduleState(ctx, token)
		if err != nil {
			t.Skipf("GetScheduleState not supported: %v", err)
		}

		t.Logf("Schedule state for token=%s: Active=%v", token, state.Active)
		if state.SpecialDay != nil {
			t.Logf("  SpecialDay=%v", *state.SpecialDay)
		}
	})

	t.Run("GetSpecialDayGroupInfoList", func(t *testing.T) {
		infos, _, err := client.GetSpecialDayGroupInfoList(ctx, nil, nil)
		if err != nil {
			t.Skipf("GetSpecialDayGroupInfoList not supported: %v", err)
		}

		t.Logf("Found %d special day groups", len(infos))
		for _, sdg := range infos {
			t.Logf("  SpecialDayGroup: Token=%s, Name=%s", sdg.Token, sdg.Name)
		}
	})
}
