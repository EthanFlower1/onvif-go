//go:build real_camera

package onvif

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestEvent_RealCamera(t *testing.T) {
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

	t.Run("GetEventServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetEventServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetEventServiceCapabilities not supported: %v", err)
		}

		t.Logf("Event service capabilities: WSSubscriptionPolicySupport=%v, WSPausableSubscriptionManagerInterfaceSupport=%v, MaxPullPoints=%d",
			caps.WSSubscriptionPolicySupport, caps.WSPausableSubscriptionManagerInterfaceSupport, caps.MaxPullPoints)
	})

	t.Run("CreatePullPointSubscription", func(t *testing.T) {
		subCtx, subCancel := context.WithTimeout(ctx, 15*time.Second)
		defer subCancel()

		sub, err := client.CreatePullPointSubscription(subCtx, "", nil, "")
		if err != nil {
			t.Skipf("CreatePullPointSubscription not supported: %v", err)
		}

		t.Logf("Created pull point subscription: ref=%s, terminationTime=%s",
			sub.SubscriptionReference, sub.TerminationTime)

		// Pull a batch of messages (non-blocking, short timeout)
		pullCtx, pullCancel := context.WithTimeout(ctx, 10*time.Second)
		defer pullCancel()

		messages, err := client.PullMessages(pullCtx, sub.SubscriptionReference, 5*time.Second, 10)
		if err != nil {
			t.Logf("PullMessages returned error (may be normal if no events): %v", err)
		} else {
			t.Logf("Pulled %d event messages", len(messages))
		}

		// Clean up subscription
		unsubCtx, unsubCancel := context.WithTimeout(ctx, 10*time.Second)
		defer unsubCancel()

		if err := client.Unsubscribe(unsubCtx, sub.SubscriptionReference); err != nil {
			t.Logf("Warning: Unsubscribe failed (non-fatal): %v", err)
		}
	})

	t.Run("CreateLegacyPullPoint", func(t *testing.T) {
		legacyCtx, legacyCancel := context.WithTimeout(ctx, 15*time.Second)
		defer legacyCancel()

		pullPointRef, err := client.CreateLegacyPullPoint(legacyCtx)
		if err != nil {
			t.Skipf("CreateLegacyPullPoint not supported: %v", err)
		}

		t.Logf("Created legacy pull point: ref=%s", pullPointRef)

		// Attempt to pull messages
		msgCtx, msgCancel := context.WithTimeout(ctx, 10*time.Second)
		defer msgCancel()

		messages, err := client.GetLegacyMessages(msgCtx, pullPointRef, 10)
		if err != nil {
			t.Logf("GetLegacyMessages returned error (may be normal if no events): %v", err)
		} else {
			t.Logf("Got %d legacy event messages", len(messages))
		}

		// Clean up
		destroyCtx, destroyCancel := context.WithTimeout(ctx, 10*time.Second)
		defer destroyCancel()

		if err := client.DestroyLegacyPullPoint(destroyCtx, pullPointRef); err != nil {
			t.Logf("Warning: DestroyLegacyPullPoint failed (non-fatal): %v", err)
		}
	})

	t.Run("PauseResumeSubscription", func(t *testing.T) {
		subCtx, subCancel := context.WithTimeout(ctx, 15*time.Second)
		defer subCancel()

		sub, err := client.CreatePullPointSubscription(subCtx, "", nil, "")
		if err != nil {
			t.Skipf("CreatePullPointSubscription not supported (required for PauseResume test): %v", err)
		}

		pauseCtx, pauseCancel := context.WithTimeout(ctx, 10*time.Second)
		defer pauseCancel()

		if err := client.PauseSubscription(pauseCtx, sub.SubscriptionReference); err != nil {
			t.Skipf("PauseSubscription not supported: %v", err)
		}

		t.Logf("Successfully paused subscription: %s", sub.SubscriptionReference)

		resumeCtx, resumeCancel := context.WithTimeout(ctx, 10*time.Second)
		defer resumeCancel()

		if err := client.ResumeSubscription(resumeCtx, sub.SubscriptionReference); err != nil {
			t.Skipf("ResumeSubscription not supported: %v", err)
		}

		t.Logf("Successfully resumed subscription: %s", sub.SubscriptionReference)

		// Clean up
		unsubCtx, unsubCancel := context.WithTimeout(ctx, 10*time.Second)
		defer unsubCancel()

		if err := client.Unsubscribe(unsubCtx, sub.SubscriptionReference); err != nil {
			t.Logf("Warning: Unsubscribe failed (non-fatal): %v", err)
		}
	})
}
