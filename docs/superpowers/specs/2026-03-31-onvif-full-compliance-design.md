# ONVIF Full Compliance Design Spec

**Date:** 2026-03-31
**Goal:** Implement all missing ONVIF client operations across all service specifications so the SDK can be used to build a fully ONVIF-compliant NVR.

## Scope

- Client-only (no server implementation)
- All ONVIF profiles: S, G, T, plus PACS and remaining services
- ~556 new operations across 20+ services
- Unit tests + integration test scaffolding for every operation
- Per-service type files for new types

## Decisions

- **Approach:** Hybrid foundation + profile-ordered vertical slices
- **Type organization:** Per-service type files (`types_recording.go`, etc.). Existing `types.go` untouched for backward compatibility.
- **Testing:** Unit tests with `httptest.NewServer` mocks for every operation (success + SOAP fault + edge cases). Integration tests with `//go:build real_camera` tag using `t.Skip` for unsupported services.
- **No server code** — client methods only.

---

## Architecture

### File Organization

New files follow the existing convention of one file per service:

```
onvif-go/
# Type files (one per new service)
types_recording.go          types_search.go
types_replay.go             types_receiver.go
types_analytics.go          types_media2.go
types_access_control.go     types_door_control.go
types_credential.go         types_schedule.go
types_auth_behavior.go      types_security.go
types_thermal.go            types_display.go
types_provisioning.go       types_uplink.go
types_appmgmt.go

# Service implementation + tests (one triplet per new service)
recording.go                recording_test.go              recording_real_camera_test.go
search.go                   search_test.go                 search_real_camera_test.go
replay.go                   replay_test.go                 replay_real_camera_test.go
receiver.go                 receiver_test.go               receiver_real_camera_test.go
analytics.go                analytics_test.go              analytics_real_camera_test.go
media2.go                   media2_test.go                 media2_real_camera_test.go
access_control.go           access_control_test.go         access_control_real_camera_test.go
door_control.go             door_control_test.go           door_control_real_camera_test.go
credential.go               credential_test.go             credential_real_camera_test.go
schedule.go                 schedule_test.go               schedule_real_camera_test.go
auth_behavior.go            auth_behavior_test.go          auth_behavior_real_camera_test.go
advanced_security.go        advanced_security_test.go      advanced_security_real_camera_test.go
thermal.go                  thermal_test.go                thermal_real_camera_test.go
display.go                  display_test.go                display_real_camera_test.go
provisioning.go             provisioning_test.go           provisioning_real_camera_test.go
uplink.go                   uplink_test.go                 uplink_real_camera_test.go
appmgmt.go                  appmgmt_test.go                appmgmt_real_camera_test.go
```

Existing files extended in Phase 2: `device.go`, `device_extended.go`, `device_security.go`, `ptz.go`, `imaging.go`, `event.go`, `deviceio.go` and their corresponding `*_test.go` files.

### Client Struct Changes

New service endpoint fields in `Client` struct (`client.go`):

```go
recordingEndpoint        string
searchEndpoint           string
replayEndpoint           string
receiverEndpoint         string
analyticsEndpoint        string
media2Endpoint           string
accessControlEndpoint    string
doorControlEndpoint      string
credentialEndpoint       string
scheduleEndpoint         string
authBehaviorEndpoint     string
advancedSecurityEndpoint string
thermalEndpoint          string
displayEndpoint          string
provisioningEndpoint     string
uplinkEndpoint           string
appmgmtEndpoint         string
```

Each gets a `getXxxEndpoint()` getter (returns endpoint if discovered, falls back to device endpoint).

Each gets a `HasXxxService() bool` public method (returns whether endpoint was discovered).

### Namespace Constants

Each service file defines its namespace constant matching the WSDL `targetNamespace`:

| Service | Namespace |
|---|---|
| Recording | `http://www.onvif.org/ver10/recording/wsdl` |
| Search | `http://www.onvif.org/ver10/search/wsdl` |
| Replay | `http://www.onvif.org/ver10/replay/wsdl` |
| Receiver | `http://www.onvif.org/ver10/receiver/wsdl` |
| Analytics (v20) | `http://www.onvif.org/ver20/analytics/wsdl` |
| AnalyticsDevice (v10) | `http://www.onvif.org/ver10/analyticsdevice/wsdl` |
| Media2 | `http://www.onvif.org/ver20/media/wsdl` |
| Access Control | `http://www.onvif.org/ver10/accesscontrol/wsdl` |
| Door Control | `http://www.onvif.org/ver10/doorcontrol/wsdl` |
| Credential | `http://www.onvif.org/ver10/credential/wsdl` |
| Schedule | `http://www.onvif.org/ver10/schedule/wsdl` |
| Auth Behavior | `http://www.onvif.org/ver10/authenticationbehavior/wsdl` |
| Advanced Security | `http://www.onvif.org/ver10/advancedsecurity/wsdl` |
| Thermal | `http://www.onvif.org/ver10/thermal/wsdl` |
| Display | `http://www.onvif.org/ver10/display/wsdl` |
| Provisioning | `http://www.onvif.org/ver10/provisioning/wsdl` |
| Uplink | `http://www.onvif.org/ver10/uplink/wsdl` |
| App Management | `http://www.onvif.org/ver10/appmgmt/wsdl` |

---

## Endpoint Discovery

### Dual-Path Discovery

`Initialize()` uses two complementary discovery methods:

1. **`GetCapabilities()`** — primary path. Covers core services: Device, Media, PTZ, Imaging, Events, Recording, Search, Replay, Receiver, Display, DeviceIO.
2. **`GetServices()`** — secondary path. Returns all services by namespace URI. Catches services not in `GetCapabilities()`: Analytics, Media2, PACS, Advanced Security, Thermal, Provisioning, Uplink, App Management.

`GetServices()` failures are non-fatal (older devices may not support it). Only sets endpoints not already discovered by `GetCapabilities()`. All discovered endpoints pass through `fixLocalhostURL()`.

### Namespace-to-Endpoint Mapping

A `serviceNamespaceMap` maps namespace URIs from `GetServices()` responses to client endpoint fields:

```go
var serviceNamespaceMap = map[string]func(c *Client, xaddr string){
    "http://www.onvif.org/ver10/recording/wsdl":              func(c *Client, x string) { c.recordingEndpoint = x },
    "http://www.onvif.org/ver10/search/wsdl":                 func(c *Client, x string) { c.searchEndpoint = x },
    // ... one entry per service
}
```

### Capability Struct Extension

New fields on the existing `Capabilities` struct in `types.go`:

```go
Recording  *RecordingCapabilities
Search     *SearchCapabilities
Replay     *ReplayCapabilities
Receiver   *ReceiverCapabilities
Display    *DisplayCapabilities
DeviceIO   *DeviceIOCapabilities
```

---

## Method Implementation Pattern

Every operation follows this canonical pattern:

### Inline SOAP Structs (unexported, in service file)

```go
type getRecordingsRequest struct {
    XMLName xml.Name `xml:"trc:GetRecordings"`
    Xmlns   string   `xml:"xmlns:trc,attr"`
}

type getRecordingsResponse struct {
    RecordingItem []struct {
        RecordingToken string `xml:"RecordingToken"`
        // ... fields matching WSDL response schema
    } `xml:"RecordingItem"`
}
```

### Public Method

```go
func (c *Client) GetRecordings(ctx context.Context) ([]*Recording, error) {
    req := getRecordingsRequest{Xmlns: recordingNamespace}
    var resp getRecordingsResponse

    username, password := c.GetCredentials()
    soapClient := soap.NewClient(c.httpClient, username, password)

    if err := soapClient.Call(ctx, c.getRecordingEndpoint(), "", req, &resp); err != nil {
        return nil, fmt.Errorf("GetRecordings failed: %w", err)
    }

    // Map inline SOAP struct to public types
    recordings := make([]*Recording, 0, len(resp.RecordingItem))
    for _, item := range resp.RecordingItem {
        recordings = append(recordings, &Recording{
            Token: item.RecordingToken,
            // ...
        })
    }
    return recordings, nil
}
```

### Public Types (in `types_*.go`)

```go
type Recording struct {
    Token         string
    Configuration RecordingConfiguration
    Tracks        []*RecordingTrack
}
```

### Rules

- Inline SOAP structs are unexported and live in the service `.go` file
- Public types are exported and live in `types_*.go` with no XML tags
- Methods always map from SOAP response to public types before returning
- Error wrapping: `fmt.Errorf("OperationName failed: %w", err)`
- Pointer fields for optional ONVIF values
- Slices never nil (return empty slice)
- Context is always the first parameter

---

## Testing Strategy

### Unit Tests

Every operation gets at minimum:

1. **Success case** — mock server returns valid SOAP response, verify parsed result
2. **SOAP fault case** — mock server returns SOAP fault, verify error returned

Additional cases where applicable:
3. **Empty response** — operations returning slices handle zero results
4. **Partial response** — optional fields missing, verify nil pointers
5. **Request validation** — verify request body contains correct XML element name and token parameters

Pattern:
```go
func TestGetRecordings(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Return SOAP XML derived from WSDL response schema
    }))
    defer server.Close()
    client, _ := NewClient(server.URL)
    recordings, err := client.GetRecordings(context.Background())
    // Assert results
}
```

### Integration Tests

Each service gets a `*_real_camera_test.go` file with `//go:build real_camera` tag:

- Reads `ONVIF_ENDPOINT`, `ONVIF_USERNAME`, `ONVIF_PASSWORD` from environment
- Calls `client.Initialize(ctx)` to discover endpoints
- Each operation is a subtest using `t.Run`
- Uses `t.Skip` for unsupported services (not `t.Fatal`)
- Chains tests: discover tokens in parent, use in subtests
- 30-second timeout per test context
- Logs responses for debugging without asserting camera-specific values

### Mock Server Extensions

Extend `testing/mock_server.go` `tokenParams` with new token types:
`RecordingToken`, `RecordingJobToken`, `TrackToken`, `SearchToken`, `AccessPointToken`, `DoorToken`, `CredentialToken`, `ScheduleToken`.

Add new `ServiceType` constants to `testing/operations.go`:
`ServiceRecording`, `ServiceSearch`, `ServiceReplay`, `ServiceReceiver`, `ServiceAnalytics`, `ServiceMedia2`, `ServiceAccessControl`, `ServiceDoorControl`, `ServiceCredential`, `ServiceSchedule`, `ServiceAuthBehavior`, `ServiceAdvancedSecurity`, `ServiceThermal`, `ServiceDisplay`, `ServiceProvisioning`, `ServiceUplink`, `ServiceAppMgmt`.

---

## Implementation Phases

### Phase 1: Foundation

Extend `client.go`:
- Add 17 new endpoint fields to `Client` struct
- Add `getXxxEndpoint()` getter for each (fallback to device endpoint)
- Add `HasXxxService() bool` for each
- Extend `Initialize()` to call `GetServices()` after `GetCapabilities()`
- Implement `serviceNamespaceMap` for namespace-to-endpoint routing
- Apply `fixLocalhostURL()` to all new endpoints

Extend `types.go`:
- Add new capability struct fields to `Capabilities`

Extend `testing/mock_server.go`:
- Add new token types to `tokenParams`

Extend `testing/operations.go`:
- Add new `ServiceType` constants

**New ops: 0 | Modified files: 4**

### Phase 2: Profile S Completion

All changes to existing files. No new files created.

**Device Service gaps (8 ops):**
- `SetNetworkInterfaces` — `device.go`
- `UpgradeFirmware`, `UpgradeSystemFirmware` — `device_extended.go`
- `GetUserRoles`, `SetUserRole`, `DeleteUserRole` — `device_extended.go`
- `GetAuthFailureWarningOptions`, `GetPasswordComplexityOptions` — `device_security.go`

**PTZ Service gaps (16 ops) — `ptz.go`:**
- `GetNodes`, `GetNode`
- `GetConfigurationOptions`, `SetConfiguration`
- `GetServiceCapabilities`, `GetCompatibleConfigurations`
- `SendAuxiliaryCommand` (PTZ namespace `tptz:`)
- `GeoMove`, `MoveAndStartTracking`
- `GetPresetTours`, `GetPresetTour`, `GetPresetTourOptions`, `CreatePresetTour`, `ModifyPresetTour`, `OperatePresetTour`, `RemovePresetTour`

**Imaging Service gaps (4 ops) — `imaging.go`:**
- `GetServiceCapabilities`
- `GetPresets`, `GetCurrentPreset`, `SetCurrentPreset`

**Event Service gaps (8 ops) — `event.go`:**
- `Subscribe`, `GetCurrentMessage` (WS-BaseNotification)
- `CreatePullPoint`, `DestroyPullPoint`, `GetMessages` (legacy)
- `PauseSubscription`, `ResumeSubscription`
- `Notify` (client-side message parser for push notification payloads received via Subscribe callback)

**DeviceIO gaps (15 ops) — `deviceio.go`:**
- Audio: `GetAudioSources`, `GetAudioSourceConfiguration`, `GetAudioSourceConfigurationOptions`, `SetAudioSourceConfiguration`, `GetAudioOutputs`, `GetAudioOutputConfiguration`, `GetAudioOutputConfigurationOptions`, `SetAudioOutputConfiguration`
- Video: `GetVideoSources`, `GetVideoSourceConfiguration`, `GetVideoSourceConfigurationOptions`, `SetVideoSourceConfiguration`
- Relay: `GetRelayOutputs`, `SetRelayOutputState`, `SetRelayOutputSettings`

**New ops: ~51 | Modified files: ~12 (existing service + test files)**

All 5 service gaps can be implemented in parallel.

### Phase 3: Profile G (Recording Ecosystem)

**Recording Service — `recording.go` (22 ops):**
- `GetServiceCapabilities`
- `CreateRecording`, `DeleteRecording`, `GetRecordings`
- `SetRecordingConfiguration`, `GetRecordingConfiguration`, `GetRecordingOptions`
- `CreateTrack`, `DeleteTrack`, `GetTrackConfiguration`, `SetTrackConfiguration`
- `CreateRecordingJob`, `DeleteRecordingJob`, `GetRecordingJobs`, `SetRecordingJobConfiguration`, `GetRecordingJobConfiguration`, `SetRecordingJobMode`, `GetRecordingJobState`
- `ExportRecordedData`, `StopExportRecordedData`, `GetExportRecordedDataState`
- `OverrideSegmentDuration`

**Receiver Service — `receiver.go` (8 ops):**
- `GetServiceCapabilities`
- `GetReceivers`, `GetReceiver`
- `CreateReceiver`, `DeleteReceiver`, `ConfigureReceiver`
- `SetReceiverMode`, `GetReceiverState`

**Search Service — `search.go` (14 ops):**
- `GetServiceCapabilities`
- `GetRecordingSummary`, `GetRecordingInformation`, `GetMediaAttributes`
- `FindRecordings`, `GetRecordingSearchResults`
- `FindEvents`, `GetEventSearchResults`
- `FindPTZPosition`, `GetPTZPositionSearchResults`
- `FindMetadata`, `GetMetadataSearchResults`
- `GetSearchState`, `EndSearch`

**Replay Service — `replay.go` (4 ops):**
- `GetServiceCapabilities`
- `GetReplayUri`, `GetReplayConfiguration`, `SetReplayConfiguration`

**New ops: 48 | New files: 4 type files + 4 service files + 8 test files = 16**

Ordering: Recording types first, then Search + Replay in parallel. Receiver is independent.

### Phase 4: Profile T (Media2 + Analytics)

**Media2 — `media2.go` (~41 new ops):**
- `AddConfiguration`, `RemoveConfiguration`
- `GetVideoEncoderInstances`
- `GetAnalyticsConfigurations`
- Masks: `GetMasks`, `GetMaskOptions`, `SetMask`, `CreateMask`, `DeleteMask`
- WebRTC: `GetWebRTCConfigurations`, `SetWebRTCConfigurations`
- Audio clips: `GetAudioClips`, `AddAudioClip`, `SetAudioClip`, `DeleteAudioClip`, `PlayAudioClip`, `GetPlayingAudioClips`
- Multicast audio: `GetMulticastAudioDecoderConfigurations`, `GetMulticastAudioDecoderConfigurationOptions`, `SetMulticastAudioDecoderConfiguration`
- `SetEQPreset`
- Media2 namespace variants of shared operations (GetProfiles, GetStreamUri, etc.)

**Analytics — `analytics.go` (31 ops):**
- Ver20 Analytics: `GetServiceCapabilities`, `GetSupportedRules`, `GetRules`, `GetRuleOptions`, `CreateRules`, `ModifyRules`, `DeleteRules`, `GetSupportedAnalyticsModules`, `GetAnalyticsModules`, `GetAnalyticsModuleOptions`, `CreateAnalyticsModules`, `ModifyAnalyticsModules`, `DeleteAnalyticsModules`, `GetSupportedMetadata`
- Ver10 AnalyticsDevice: `GetAnalyticsEngines`, `GetAnalyticsEngine`, `GetAnalyticsEngineControls`, `GetAnalyticsEngineControl`, `CreateAnalyticsEngineControl`, `SetAnalyticsEngineControl`, `DeleteAnalyticsEngineControl`, `GetAnalyticsEngineInputs`, `GetAnalyticsEngineInput`, `CreateAnalyticsEngineInputs`, `SetAnalyticsEngineInput`, `DeleteAnalyticsEngineInputs`, `GetAnalyticsDeviceStreamUri`, `GetAnalyticsState`, `GetVideoAnalyticsConfiguration`, `SetVideoAnalyticsConfiguration`

**New ops: ~72 | New files: 2 type files + 2 service files + 4 test files = 8**

Media2 and Analytics can be implemented in parallel.

### Phase 5: PACS (Physical Access Control)

**Access Control — `access_control.go` (~48 ops):**
- Service capabilities, access point info/list/state
- Area info/list, enable/disable access points
- External authorization decisions

**Door Control — `door_control.go` (~38 ops):**
- Door info/list/state, door mode
- Lock/unlock/double-lock/block operations
- Door alarm and tamper handling

**Credential — `credential.go` (~56 ops):**
- Credential CRUD, state management
- Credential identifiers, access profiles
- Anti-passback configuration

**Schedule — `schedule.go` (~36 ops):**
- Schedule CRUD, special day group CRUD
- Schedule state queries

**Auth Behavior — `auth_behavior.go` (~34 ops):**
- Authentication behavior profiles and policies
- Security level management

**New ops: ~212 | New files: 5 type files + 5 service files + 10 test files = 20**

Ordering: Access Control first (shared types), then Door Control + Schedule in parallel, then Credential, then Auth Behavior.

### Phase 6: Remaining Services

**Advanced Security — `advanced_security.go` (~124 ops):**
- TLS server/client configuration
- Certificate management (PKCS#10, PKCS#12)
- Key management, CRL management
- IEEE 802.1X, certification path management

**Thermal — `thermal.go` (~10 ops)**
**Display — `display.go` (~10 ops)**
**Provisioning — `provisioning.go` (~10 ops)**
**Uplink — `uplink.go` (~10 ops)**
**App Management — `appmgmt.go` (~10 ops)**

**New ops: ~174 | New files: 6 type files + 6 service files + 12 test files = 24**

All 6 services are fully independent and can be implemented in parallel.

---

## Phase Dependencies

```
Phase 1: Foundation
    |
    +---> Phase 2: Profile S Completion (5 parallel sub-tasks)
    +---> Phase 3: Profile G (Recording first, then Search + Replay parallel)
    +---> Phase 4: Profile T (Media2 + Analytics parallel)
    +---> Phase 5: PACS (Access Control -> Door Control + Schedule -> Credential -> Auth Behavior)
    +---> Phase 6: Remaining (6 parallel sub-tasks)
```

Phases 2-6 depend only on Phase 1, not on each other.

---

## Summary

| Phase | New Ops | New Type Files | New Service Files | New Test Files |
|---|---|---|---|---|
| 1: Foundation | 0 | 0 | 0 | 0 |
| 2: Profile S | ~51 | 0 | 0 (extend existing) | 0 (extend existing) |
| 3: Profile G | 48 | 4 | 4 | 8 |
| 4: Profile T | ~72 | 2 | 2 | 4 |
| 5: PACS | ~212 | 5 | 5 | 10 |
| 6: Remaining | ~174 | 6 | 6 | 12 |
| **Total** | **~557** | **17** | **17** | **34** |

WSDL source files for all operations are in `resources/specs/wsdl/`. Every operation's request/response XML structure is derived from the corresponding WSDL, not invented.
