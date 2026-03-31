# Phase 4: Profile T (Media2 + Analytics) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the complete Media2 and Analytics client services for Profile T (advanced streaming) compliance. ~72 new operations.

**Architecture:** Two new service files (`media2.go`, `analytics.go`) with per-service type files. All operations follow the established callMethod pattern from Phases 1-3.

**Tech Stack:** Go stdlib, internal/soap, httptest.

**Depends on:** Phase 1 (Foundation) must be complete. Independent of Phases 2, 3, 5, 6.

**Spec:** `docs/superpowers/specs/2026-03-31-onvif-full-compliance-design.md`

**WSDLs:**
- Media2: `resources/specs/wsdl/ver20/media/wsdl/media.wsdl` (namespace: `http://www.onvif.org/ver20/media/wsdl`, prefix: `tr2:`)
- Analytics v20: `resources/specs/wsdl/ver20/analytics/wsdl/analytics.wsdl` (namespace: `http://www.onvif.org/ver20/analytics/wsdl`, prefix: `tan:`)
- AnalyticsDevice v10: `resources/specs/wsdl/ver10/analyticsdevice.wsdl` (namespace: `http://www.onvif.org/ver10/analyticsdevice/wsdl`, prefix: `tad:`)

---

## File Structure

**New files:**
- `types_media2.go` — Mask, AudioClip, WebRTCConfiguration types
- `types_analytics.go` — AnalyticsModule, Rule, AnalyticsEngine types
- `media2.go` — ~41 Media2 service client methods
- `media2_test.go` — Unit tests
- `media2_real_camera_test.go` — Integration tests
- `analytics.go` — 31 Analytics service client methods (v20 + v10 AnalyticsDevice)
- `analytics_test.go` — Unit tests
- `analytics_real_camera_test.go` — Integration tests

---

## Tasks

### Task 1: Media2 Types (`types_media2.go`)
- [ ] Define types: `Mask`, `MaskOptions`, `AudioClip`, `WebRTCConfiguration`, `Media2ServiceCapabilities`, `EQPreset`
- [ ] Verify compilation
- [ ] Commit

### Task 2: Media2 — Configuration Operations (~10 ops)
- [ ] Implement: `AddConfiguration`, `RemoveConfiguration`, `GetVideoEncoderInstances`, `GetAnalyticsConfigurations`
- [ ] Implement Media2 namespace variants of shared operations: `GetProfiles`, `GetStreamUri`, `GetSnapshotUri`, `StartMulticastStreaming`, `StopMulticastStreaming`, `SetSynchronizationPoint`
- [ ] Write tests for all operations
- [ ] Commit

### Task 3: Media2 — Mask Operations (5 ops)
- [ ] Implement: `GetMasks`, `GetMaskOptions`, `SetMask`, `CreateMask`, `DeleteMask`
- [ ] Write tests
- [ ] Commit

### Task 4: Media2 — Audio Clip Operations (6 ops)
- [ ] Implement: `GetAudioClips`, `AddAudioClip`, `SetAudioClip`, `DeleteAudioClip`, `PlayAudioClip`, `GetPlayingAudioClips`
- [ ] Write tests
- [ ] Commit

### Task 5: Media2 — Remaining Operations (~6 ops)
- [ ] Implement: `GetWebRTCConfigurations`, `SetWebRTCConfigurations`, `SetEQPreset`
- [ ] Implement: `GetMulticastAudioDecoderConfigurations`, `GetMulticastAudioDecoderConfigurationOptions`, `SetMulticastAudioDecoderConfiguration`
- [ ] Write tests
- [ ] Commit

### Task 6: Analytics Types (`types_analytics.go`)
- [ ] Define types: `AnalyticsModule`, `AnalyticsModuleOptions`, `Rule`, `RuleOptions`, `AnalyticsEngine`, `AnalyticsEngineControl`, `AnalyticsEngineInput`
- [ ] Verify compilation
- [ ] Commit

### Task 7: Analytics v20 — Rule Operations (7 ops)
- [ ] Implement: `GetSupportedRules`, `GetRules`, `GetRuleOptions`, `CreateRules`, `ModifyRules`, `DeleteRules`, `GetServiceCapabilities`
- [ ] Write tests
- [ ] Commit

### Task 8: Analytics v20 — Module Operations (7 ops)
- [ ] Implement: `GetSupportedAnalyticsModules`, `GetAnalyticsModules`, `GetAnalyticsModuleOptions`, `CreateAnalyticsModules`, `ModifyAnalyticsModules`, `DeleteAnalyticsModules`, `GetSupportedMetadata`
- [ ] Write tests
- [ ] Commit

### Task 9: AnalyticsDevice v10 Operations (17 ops)
- [ ] Implement engine operations: `GetAnalyticsEngines`, `GetAnalyticsEngine`
- [ ] Implement control operations: `GetAnalyticsEngineControls`, `GetAnalyticsEngineControl`, `CreateAnalyticsEngineControl`, `SetAnalyticsEngineControl`, `DeleteAnalyticsEngineControl`
- [ ] Implement input operations: `GetAnalyticsEngineInputs`, `GetAnalyticsEngineInput`, `CreateAnalyticsEngineInputs`, `SetAnalyticsEngineInput`, `DeleteAnalyticsEngineInputs`
- [ ] Implement remaining: `GetAnalyticsDeviceStreamUri`, `GetAnalyticsState`, `GetVideoAnalyticsConfiguration`, `SetVideoAnalyticsConfiguration`, `GetServiceCapabilities`
- [ ] Write tests
- [ ] Commit

### Task 10: Integration Tests
- [ ] Create `media2_real_camera_test.go` and `analytics_real_camera_test.go`
- [ ] Verify compilation with `go vet -tags=real_camera ./...`
- [ ] Run full test suite, fix any lint issues
- [ ] Commit
