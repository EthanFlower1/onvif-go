# Phase 1 (Foundation) + Phase 2 (Profile S Completion) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Extend the client infrastructure to support all ONVIF services, then complete Profile S compliance by filling every missing Device, PTZ, Imaging, Event, and DeviceIO operation.

**Architecture:** Add 17 service endpoint fields to Client with dual-path discovery (GetCapabilities + GetServices). Implement ~51 missing operations in existing service files following the established callMethod pattern. Each operation gets unit tests (success + SOAP fault) and integration test scaffolding.

**Tech Stack:** Go stdlib (encoding/xml, net/http, context, sync), internal/soap package, httptest for mocks.

**Spec:** `docs/superpowers/specs/2026-03-31-onvif-full-compliance-design.md`

---

## File Structure

**Modified files (Phase 1 — Foundation):**
- `client.go` — Add 17 endpoint fields, getters, HasXxxService() methods, extend Initialize()
- `types.go` — Add capability structs for new services
- `testing/capture_types.go` — Add ServiceType constants and namespace mappings
- `testing/mock_server.go` — Add token types to tokenParams

**Modified files (Phase 2 — Profile S):**
- `device.go` — Add SetNetworkInterfaces
- `device_extended.go` — Add UpgradeSystemFirmware, GetUserRoles, SetUserRole, DeleteUserRole
- `device_security.go` — Add GetAuthFailureWarningOptions, GetPasswordComplexityOptions
- `ptz.go` — Add 16 missing operations
- `imaging.go` — Add 4 missing operations
- `event.go` — Add 8 missing operations
- `deviceio.go` — Add 15 missing operations
- All corresponding `*_test.go` files

---

## Task 1: Add Service Endpoint Fields to Client

**Files:**
- Modify: `client.go:33-46`

- [ ] **Step 1: Add endpoint fields to Client struct**

In `client.go`, replace the Client struct (lines 33-46) with:

```go
type Client struct {
	endpoint   string
	username   string
	password   string
	httpClient *http.Client
	mu         sync.RWMutex

	// Service endpoints (discovered via Initialize)
	mediaEndpoint            string
	ptzEndpoint              string
	imagingEndpoint          string
	eventEndpoint            string
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
	appmgmtEndpoint          string
}
```

- [ ] **Step 2: Add HasXxxService() public methods**

Append after the existing `Endpoint()` method in `client.go`:

```go
// HasRecordingService returns true if the recording service endpoint was discovered.
func (c *Client) HasRecordingService() bool { return c.recordingEndpoint != "" }

// HasSearchService returns true if the search service endpoint was discovered.
func (c *Client) HasSearchService() bool { return c.searchEndpoint != "" }

// HasReplayService returns true if the replay service endpoint was discovered.
func (c *Client) HasReplayService() bool { return c.replayEndpoint != "" }

// HasReceiverService returns true if the receiver service endpoint was discovered.
func (c *Client) HasReceiverService() bool { return c.receiverEndpoint != "" }

// HasAnalyticsService returns true if the analytics service endpoint was discovered.
func (c *Client) HasAnalyticsService() bool { return c.analyticsEndpoint != "" }

// HasMedia2Service returns true if the media2 service endpoint was discovered.
func (c *Client) HasMedia2Service() bool { return c.media2Endpoint != "" }

// HasAccessControlService returns true if the access control service endpoint was discovered.
func (c *Client) HasAccessControlService() bool { return c.accessControlEndpoint != "" }

// HasDoorControlService returns true if the door control service endpoint was discovered.
func (c *Client) HasDoorControlService() bool { return c.doorControlEndpoint != "" }

// HasCredentialService returns true if the credential service endpoint was discovered.
func (c *Client) HasCredentialService() bool { return c.credentialEndpoint != "" }

// HasScheduleService returns true if the schedule service endpoint was discovered.
func (c *Client) HasScheduleService() bool { return c.scheduleEndpoint != "" }

// HasAuthBehaviorService returns true if the authentication behavior service endpoint was discovered.
func (c *Client) HasAuthBehaviorService() bool { return c.authBehaviorEndpoint != "" }

// HasAdvancedSecurityService returns true if the advanced security service endpoint was discovered.
func (c *Client) HasAdvancedSecurityService() bool { return c.advancedSecurityEndpoint != "" }

// HasThermalService returns true if the thermal service endpoint was discovered.
func (c *Client) HasThermalService() bool { return c.thermalEndpoint != "" }

// HasDisplayService returns true if the display service endpoint was discovered.
func (c *Client) HasDisplayService() bool { return c.displayEndpoint != "" }

// HasProvisioningService returns true if the provisioning service endpoint was discovered.
func (c *Client) HasProvisioningService() bool { return c.provisioningEndpoint != "" }

// HasUplinkService returns true if the uplink service endpoint was discovered.
func (c *Client) HasUplinkService() bool { return c.uplinkEndpoint != "" }

// HasAppMgmtService returns true if the app management service endpoint was discovered.
func (c *Client) HasAppMgmtService() bool { return c.appmgmtEndpoint != "" }
```

- [ ] **Step 3: Build and verify compilation**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go build ./...`
Expected: Clean compilation, no errors.

- [ ] **Step 4: Commit**

```bash
git add client.go
git commit -m "feat: add service endpoint fields and HasXxxService methods to Client"
```

---

## Task 2: Extend Initialize() with GetServices() Discovery

**Files:**
- Modify: `client.go:199-222`
- Modify: `types.go` (Capabilities struct, line 14)

- [ ] **Step 1: Add new capability structs to types.go**

Append after the existing `CapabilitiesExtension` struct (line 115) in `types.go`:

```go
// RecordingCapabilities contains recording service capabilities.
type RecordingCapabilities struct {
	XAddr string
}

// SearchCapabilities contains search service capabilities.
type SearchCapabilities struct {
	XAddr string
}

// ReplayCapabilities contains replay service capabilities.
type ReplayCapabilities struct {
	XAddr string
}

// ReceiverCapabilities contains receiver service capabilities.
type ReceiverCapabilities struct {
	XAddr string
}

// DisplayCapabilities contains display service capabilities.
type DisplayCapabilities struct {
	XAddr string
}

// DeviceIOCapabilities contains device IO service capabilities.
type DeviceIOCapabilities struct {
	XAddr string
}
```

- [ ] **Step 2: Add new fields to Capabilities struct**

In `types.go`, add to the `Capabilities` struct (after line 21, before `Extension`):

```go
type Capabilities struct {
	Analytics  *AnalyticsCapabilities
	Device     *DeviceCapabilities
	Events     *EventCapabilities
	Imaging    *ImagingCapabilities
	Media      *MediaCapabilities
	PTZ        *PTZCapabilities
	Recording  *RecordingCapabilities
	Search     *SearchCapabilities
	Replay     *ReplayCapabilities
	Receiver   *ReceiverCapabilities
	Display    *DisplayCapabilities
	DeviceIO   *DeviceIOCapabilities
	Extension  *CapabilitiesExtension
}
```

- [ ] **Step 3: Add serviceNamespaceMap and extend Initialize()**

In `client.go`, add the namespace map before `Initialize()`:

```go
// serviceNamespaceMap maps ONVIF service namespace URIs to endpoint setters.
var serviceNamespaceMap = map[string]func(c *Client, xaddr string){
	"http://www.onvif.org/ver10/recording/wsdl":              func(c *Client, x string) { c.recordingEndpoint = x },
	"http://www.onvif.org/ver10/search/wsdl":                 func(c *Client, x string) { c.searchEndpoint = x },
	"http://www.onvif.org/ver10/replay/wsdl":                 func(c *Client, x string) { c.replayEndpoint = x },
	"http://www.onvif.org/ver10/receiver/wsdl":               func(c *Client, x string) { c.receiverEndpoint = x },
	"http://www.onvif.org/ver20/analytics/wsdl":              func(c *Client, x string) { c.analyticsEndpoint = x },
	"http://www.onvif.org/ver20/media/wsdl":                  func(c *Client, x string) { c.media2Endpoint = x },
	"http://www.onvif.org/ver10/accesscontrol/wsdl":          func(c *Client, x string) { c.accessControlEndpoint = x },
	"http://www.onvif.org/ver10/doorcontrol/wsdl":            func(c *Client, x string) { c.doorControlEndpoint = x },
	"http://www.onvif.org/ver10/credential/wsdl":             func(c *Client, x string) { c.credentialEndpoint = x },
	"http://www.onvif.org/ver10/schedule/wsdl":               func(c *Client, x string) { c.scheduleEndpoint = x },
	"http://www.onvif.org/ver10/authenticationbehavior/wsdl": func(c *Client, x string) { c.authBehaviorEndpoint = x },
	"http://www.onvif.org/ver10/advancedsecurity/wsdl":       func(c *Client, x string) { c.advancedSecurityEndpoint = x },
	"http://www.onvif.org/ver10/thermal/wsdl":                func(c *Client, x string) { c.thermalEndpoint = x },
	"http://www.onvif.org/ver10/display/wsdl":                func(c *Client, x string) { c.displayEndpoint = x },
	"http://www.onvif.org/ver10/provisioning/wsdl":           func(c *Client, x string) { c.provisioningEndpoint = x },
	"http://www.onvif.org/ver10/uplink/wsdl":                 func(c *Client, x string) { c.uplinkEndpoint = x },
	"http://www.onvif.org/ver10/appmgmt/wsdl":                func(c *Client, x string) { c.appmgmtEndpoint = x },
}
```

Replace `Initialize()` (lines 199-222) with:

```go
// Initialize discovers and initializes service endpoints.
func (c *Client) Initialize(ctx context.Context) error {
	// Get device capabilities (primary discovery path)
	capabilities, err := c.GetCapabilities(ctx)
	if err != nil {
		return fmt.Errorf("failed to get capabilities: %w", err)
	}

	// Extract service endpoints from capabilities
	if capabilities.Media != nil && capabilities.Media.XAddr != "" {
		c.mediaEndpoint = c.fixLocalhostURL(capabilities.Media.XAddr)
	}
	if capabilities.PTZ != nil && capabilities.PTZ.XAddr != "" {
		c.ptzEndpoint = c.fixLocalhostURL(capabilities.PTZ.XAddr)
	}
	if capabilities.Imaging != nil && capabilities.Imaging.XAddr != "" {
		c.imagingEndpoint = c.fixLocalhostURL(capabilities.Imaging.XAddr)
	}
	if capabilities.Events != nil && capabilities.Events.XAddr != "" {
		c.eventEndpoint = c.fixLocalhostURL(capabilities.Events.XAddr)
	}
	if capabilities.Recording != nil && capabilities.Recording.XAddr != "" {
		c.recordingEndpoint = c.fixLocalhostURL(capabilities.Recording.XAddr)
	}
	if capabilities.Search != nil && capabilities.Search.XAddr != "" {
		c.searchEndpoint = c.fixLocalhostURL(capabilities.Search.XAddr)
	}
	if capabilities.Replay != nil && capabilities.Replay.XAddr != "" {
		c.replayEndpoint = c.fixLocalhostURL(capabilities.Replay.XAddr)
	}
	if capabilities.Receiver != nil && capabilities.Receiver.XAddr != "" {
		c.receiverEndpoint = c.fixLocalhostURL(capabilities.Receiver.XAddr)
	}
	if capabilities.Display != nil && capabilities.Display.XAddr != "" {
		c.displayEndpoint = c.fixLocalhostURL(capabilities.Display.XAddr)
	}
	if capabilities.DeviceIO != nil && capabilities.DeviceIO.XAddr != "" {
		// Only set if not already discovered via direct endpoint
		if c.endpoint != "" {
			// DeviceIO often shares the device endpoint
		}
	}

	// Secondary discovery via GetServices (catches services not in GetCapabilities)
	services, err := c.GetServices(ctx)
	if err == nil {
		for _, svc := range services {
			if setter, ok := serviceNamespaceMap[svc.Namespace]; ok {
				xaddr := c.fixLocalhostURL(svc.XAddr)
				setter(c, xaddr)
			}
		}
	}
	// GetServices failure is non-fatal — older devices may not support it

	return nil
}
```

- [ ] **Step 4: Build and verify**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go build ./...`
Expected: Clean compilation.

- [ ] **Step 5: Write test for Initialize with GetServices**

Add to `client_test.go`:

```go
func TestInitialize_DiscoverServicesEndpoints(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		bodyStr := string(body)

		var response string
		if strings.Contains(bodyStr, "GetCapabilities") {
			callCount++
			response = `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<tds:GetCapabilitiesResponse xmlns:tds="http://www.onvif.org/ver10/device/wsdl">
<tds:Capabilities>
<tt:Media xmlns:tt="http://www.onvif.org/ver10/schema"><tt:XAddr>http://192.168.1.100/onvif/media</tt:XAddr></tt:Media>
</tds:Capabilities>
</tds:GetCapabilitiesResponse>
</soap:Body></soap:Envelope>`
		} else if strings.Contains(bodyStr, "GetServices") {
			callCount++
			response = `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<tds:GetServicesResponse xmlns:tds="http://www.onvif.org/ver10/device/wsdl">
<tds:Service>
<tds:Namespace>http://www.onvif.org/ver10/recording/wsdl</tds:Namespace>
<tds:XAddr>http://192.168.1.100/onvif/recording</tds:XAddr>
</tds:Service>
<tds:Service>
<tds:Namespace>http://www.onvif.org/ver20/analytics/wsdl</tds:Namespace>
<tds:XAddr>http://192.168.1.100/onvif/analytics</tds:XAddr>
</tds:Service>
</tds:GetServicesResponse>
</soap:Body></soap:Envelope>`
		}

		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if err := client.Initialize(context.Background()); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	if !client.HasRecordingService() {
		t.Error("expected recording service to be discovered")
	}
	if !client.HasAnalyticsService() {
		t.Error("expected analytics service to be discovered")
	}
	if callCount < 2 {
		t.Errorf("expected at least 2 SOAP calls, got %d", callCount)
	}
}
```

- [ ] **Step 6: Run tests**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go test -run TestInitialize -v`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add client.go types.go client_test.go
git commit -m "feat: extend Initialize with GetServices discovery for all ONVIF services"
```

---

## Task 3: Extend Testing Infrastructure

**Files:**
- Modify: `testing/capture_types.go:12-22` and `testing/capture_types.go:338-348`
- Modify: `testing/mock_server.go:461-493`

- [ ] **Step 1: Add ServiceType constants**

In `testing/capture_types.go`, extend the constants block (after line 22):

```go
const (
	ServiceDevice          ServiceType = "Device"
	ServiceMedia           ServiceType = "Media"
	ServicePTZ             ServiceType = "PTZ"
	ServiceImaging         ServiceType = "Imaging"
	ServiceEvent           ServiceType = "Event"
	ServiceDeviceIO        ServiceType = "DeviceIO"
	ServiceRecording       ServiceType = "Recording"
	ServiceSearch          ServiceType = "Search"
	ServiceReplay          ServiceType = "Replay"
	ServiceReceiver        ServiceType = "Receiver"
	ServiceAnalytics       ServiceType = "Analytics"
	ServiceMedia2          ServiceType = "Media2"
	ServiceAccessControl   ServiceType = "AccessControl"
	ServiceDoorControl     ServiceType = "DoorControl"
	ServiceCredential      ServiceType = "Credential"
	ServiceSchedule        ServiceType = "Schedule"
	ServiceAuthBehavior    ServiceType = "AuthBehavior"
	ServiceAdvancedSecurity ServiceType = "AdvancedSecurity"
	ServiceThermal         ServiceType = "Thermal"
	ServiceDisplay         ServiceType = "Display"
	ServiceProvisioning    ServiceType = "Provisioning"
	ServiceUplink          ServiceType = "Uplink"
	ServiceAppMgmt         ServiceType = "AppMgmt"
	ServiceUnknown         ServiceType = "Unknown"
)
```

- [ ] **Step 2: Extend serviceNamespaces map**

In `testing/capture_types.go`, extend `serviceNamespaces` (after line 348):

```go
var serviceNamespaces = map[string]ServiceType{
	"http://www.onvif.org/ver10/device/wsdl":                 ServiceDevice,
	"http://www.onvif.org/ver10/media/wsdl":                  ServiceMedia,
	"http://www.onvif.org/ver20/media/wsdl":                  ServiceMedia2,
	"http://www.onvif.org/ver20/ptz/wsdl":                    ServicePTZ,
	"http://www.onvif.org/ver10/ptz/wsdl":                    ServicePTZ,
	"http://www.onvif.org/ver20/imaging/wsdl":                ServiceImaging,
	"http://www.onvif.org/ver10/imaging/wsdl":                ServiceImaging,
	"http://www.onvif.org/ver10/events/wsdl":                 ServiceEvent,
	"http://www.onvif.org/ver10/deviceIO/wsdl":               ServiceDeviceIO,
	"http://www.onvif.org/ver10/recording/wsdl":              ServiceRecording,
	"http://www.onvif.org/ver10/search/wsdl":                 ServiceSearch,
	"http://www.onvif.org/ver10/replay/wsdl":                 ServiceReplay,
	"http://www.onvif.org/ver10/receiver/wsdl":               ServiceReceiver,
	"http://www.onvif.org/ver20/analytics/wsdl":              ServiceAnalytics,
	"http://www.onvif.org/ver10/accesscontrol/wsdl":          ServiceAccessControl,
	"http://www.onvif.org/ver10/doorcontrol/wsdl":            ServiceDoorControl,
	"http://www.onvif.org/ver10/credential/wsdl":             ServiceCredential,
	"http://www.onvif.org/ver10/schedule/wsdl":               ServiceSchedule,
	"http://www.onvif.org/ver10/authenticationbehavior/wsdl": ServiceAuthBehavior,
	"http://www.onvif.org/ver10/advancedsecurity/wsdl":       ServiceAdvancedSecurity,
	"http://www.onvif.org/ver10/thermal/wsdl":                ServiceThermal,
	"http://www.onvif.org/ver10/display/wsdl":                ServiceDisplay,
	"http://www.onvif.org/ver10/provisioning/wsdl":           ServiceProvisioning,
	"http://www.onvif.org/ver10/uplink/wsdl":                 ServiceUplink,
	"http://www.onvif.org/ver10/appmgmt/wsdl":                ServiceAppMgmt,
}
```

- [ ] **Step 3: Extend tokenParams in mock_server.go**

In `testing/mock_server.go`, add to the `tokenParams` slice (after line 493):

```go
var tokenParams = []string{
	// Core tokens
	"ProfileToken",
	"ConfigurationToken",
	"VideoSourceToken",
	"AudioSourceToken",
	"PresetToken",
	"Token",
	// Configuration tokens
	"VideoSourceConfigurationToken",
	"AudioSourceConfigurationToken",
	"VideoEncoderConfigurationToken",
	"AudioEncoderConfigurationToken",
	"MetadataConfigurationToken",
	"PTZConfigurationToken",
	// Event/subscription tokens
	"SubscriptionReference",
	// Extended tokens
	"OSDToken",
	"NodeToken",
	"RelayOutputToken",
	"VideoOutputToken",
	"DigitalInputToken",
	"SerialPortToken",
	"StorageConfigurationToken",
	"CertificateID",
	"RecordingToken",
	"RecordingJobToken",
	"AnalyticsConfigurationToken",
	"RuleToken",
	"ScheduleToken",
	"SpecialDayGroupToken",
	// New tokens for full compliance
	"TrackToken",
	"SearchToken",
	"AccessPointToken",
	"DoorToken",
	"CredentialToken",
	"AuthenticationProfileToken",
	"SecurityLevelToken",
	"KeyID",
	"CertificationPathID",
	"PassphraseID",
	"CRLID",
	"CertPathValidationPolicyID",
	"PresetTourToken",
	"AreaToken",
	"ReceiverToken",
	"AudioOutputToken",
}
```

- [ ] **Step 4: Build and run existing tests**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go build ./... && go test ./testing/ -v -count=1`
Expected: All existing tests still pass.

- [ ] **Step 5: Commit**

```bash
git add testing/capture_types.go testing/mock_server.go
git commit -m "feat: extend testing infrastructure with all ONVIF service types and tokens"
```

---

## Task 4: PTZ Service — Node and Configuration Operations (6 ops)

**Files:**
- Modify: `ptz.go`
- Modify: `ptz_test.go`

- [ ] **Step 1: Write failing tests for GetNodes, GetNode, GetConfigurationOptions, SetConfiguration, GetServiceCapabilities, GetCompatibleConfigurations**

Add to `ptz_test.go`:

```go
func TestGetNodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<tptz:GetNodesResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
<tptz:PTZNode token="NodeToken1">
<tt:Name xmlns:tt="http://www.onvif.org/ver10/schema">PTZ Node 1</tt:Name>
<tt:SupportedPTZSpaces xmlns:tt="http://www.onvif.org/ver10/schema">
<tt:AbsolutePanTiltPositionSpace>
<tt:URI>http://www.onvif.org/ver10/tptz/PanTiltSpaces/PositionGenericSpace</tt:URI>
<tt:XRange><tt:Min>-1</tt:Min><tt:Max>1</tt:Max></tt:XRange>
<tt:YRange><tt:Min>-1</tt:Min><tt:Max>1</tt:Max></tt:YRange>
</tt:AbsolutePanTiltPositionSpace>
</tt:SupportedPTZSpaces>
<tt:HomeSupported xmlns:tt="http://www.onvif.org/ver10/schema">true</tt:HomeSupported>
<tt:MaximumNumberOfPresets xmlns:tt="http://www.onvif.org/ver10/schema">255</tt:MaximumNumberOfPresets>
</tptz:PTZNode>
</tptz:GetNodesResponse>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	nodes, err := client.GetNodes(context.Background())
	if err != nil {
		t.Fatalf("GetNodes failed: %v", err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}
	if nodes[0].Token != "NodeToken1" {
		t.Errorf("expected token NodeToken1, got %s", nodes[0].Token)
	}
	if nodes[0].Name != "PTZ Node 1" {
		t.Errorf("expected name 'PTZ Node 1', got %s", nodes[0].Name)
	}
}

func TestGetNode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<tptz:GetNodeResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
<tptz:PTZNode token="NodeToken1">
<tt:Name xmlns:tt="http://www.onvif.org/ver10/schema">PTZ Node 1</tt:Name>
<tt:HomeSupported xmlns:tt="http://www.onvif.org/ver10/schema">true</tt:HomeSupported>
</tptz:PTZNode>
</tptz:GetNodeResponse>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	node, err := client.GetNode(context.Background(), "NodeToken1")
	if err != nil {
		t.Fatalf("GetNode failed: %v", err)
	}
	if node.Token != "NodeToken1" {
		t.Errorf("expected token NodeToken1, got %s", node.Token)
	}
}

func TestGetPTZConfigurationOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<tptz:GetConfigurationOptionsResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
<tptz:PTZConfigurationOptions>
<tt:PTZTimeout xmlns:tt="http://www.onvif.org/ver10/schema">
<tt:Min>PT1S</tt:Min>
<tt:Max>PT300S</tt:Max>
</tt:PTZTimeout>
</tptz:PTZConfigurationOptions>
</tptz:GetConfigurationOptionsResponse>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	opts, err := client.GetPTZConfigurationOptions(context.Background(), "ConfigToken1")
	if err != nil {
		t.Fatalf("GetPTZConfigurationOptions failed: %v", err)
	}
	if opts == nil {
		t.Fatal("expected non-nil options")
	}
}

func TestSetPTZConfiguration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<tptz:SetConfigurationResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl"/>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	err := client.SetPTZConfiguration(context.Background(), &PTZConfiguration{Token: "ConfigToken1"}, true)
	if err != nil {
		t.Fatalf("SetPTZConfiguration failed: %v", err)
	}
}

func TestGetPTZServiceCapabilities(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<tptz:GetServiceCapabilitiesResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
<tptz:Capabilities EFlip="true" Reverse="true"/>
</tptz:GetServiceCapabilitiesResponse>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	caps, err := client.GetPTZServiceCapabilities(context.Background())
	if err != nil {
		t.Fatalf("GetPTZServiceCapabilities failed: %v", err)
	}
	if caps == nil {
		t.Fatal("expected non-nil capabilities")
	}
}

func TestGetCompatiblePTZConfigurations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<tptz:GetCompatibleConfigurationsResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
<tptz:PTZConfiguration token="PTZConfig1">
<tt:Name xmlns:tt="http://www.onvif.org/ver10/schema">PTZ Config 1</tt:Name>
<tt:NodeToken xmlns:tt="http://www.onvif.org/ver10/schema">NodeToken1</tt:NodeToken>
</tptz:PTZConfiguration>
</tptz:GetCompatibleConfigurationsResponse>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	configs, err := client.GetCompatiblePTZConfigurationsForProfile(context.Background(), "Profile1")
	if err != nil {
		t.Fatalf("GetCompatiblePTZConfigurationsForProfile failed: %v", err)
	}
	if len(configs) != 1 {
		t.Fatalf("expected 1 config, got %d", len(configs))
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go test -run "TestGetNodes|TestGetNode|TestGetPTZConfigurationOptions|TestSetPTZConfiguration|TestGetPTZServiceCapabilities|TestGetCompatiblePTZConfigurations" -v`
Expected: FAIL — methods not defined.

- [ ] **Step 3: Add PTZNode and PTZConfigurationOptions types to types.go**

Append to `types.go`:

```go
// PTZNode represents a PTZ node on the device.
type PTZNode struct {
	Token                  string
	Name                   string
	HomeSupported          bool
	MaximumNumberOfPresets int
	AuxiliaryCommands      []string
}

// PTZConfigurationOptions represents available PTZ configuration options.
type PTZConfigurationOptions struct {
	PTZTimeout *struct {
		Min string
		Max string
	}
}

// PTZServiceCapabilities represents PTZ service capabilities.
type PTZServiceCapabilities struct {
	EFlip   bool
	Reverse bool
}
```

- [ ] **Step 4: Implement the 6 methods in ptz.go**

Append to `ptz.go`:

```go
// GetNodes returns all PTZ nodes on the device.
func (c *Client) GetNodes(ctx context.Context) ([]*PTZNode, error) {
	type getNodesRequest struct {
		XMLName xml.Name `xml:"tptz:GetNodes"`
		Xmlns   string   `xml:"xmlns:tptz,attr"`
	}

	type getNodesResponse struct {
		PTZNode []struct {
			Token                  string `xml:"token,attr"`
			Name                   string `xml:"Name"`
			HomeSupported          bool   `xml:"HomeSupported"`
			MaximumNumberOfPresets int    `xml:"MaximumNumberOfPresets"`
		} `xml:"PTZNode"`
	}

	req := getNodesRequest{Xmlns: ptzNamespace}
	var resp getNodesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getPTZEndpoint(), "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetNodes failed: %w", err)
	}

	nodes := make([]*PTZNode, 0, len(resp.PTZNode))
	for _, n := range resp.PTZNode {
		nodes = append(nodes, &PTZNode{
			Token:                  n.Token,
			Name:                   n.Name,
			HomeSupported:          n.HomeSupported,
			MaximumNumberOfPresets: n.MaximumNumberOfPresets,
		})
	}

	return nodes, nil
}

// GetNode returns a specific PTZ node.
func (c *Client) GetNode(ctx context.Context, nodeToken string) (*PTZNode, error) {
	type getNodeRequest struct {
		XMLName   xml.Name `xml:"tptz:GetNode"`
		Xmlns     string   `xml:"xmlns:tptz,attr"`
		NodeToken string   `xml:"tptz:NodeToken"`
	}

	type getNodeResponse struct {
		PTZNode struct {
			Token                  string `xml:"token,attr"`
			Name                   string `xml:"Name"`
			HomeSupported          bool   `xml:"HomeSupported"`
			MaximumNumberOfPresets int    `xml:"MaximumNumberOfPresets"`
		} `xml:"PTZNode"`
	}

	req := getNodeRequest{Xmlns: ptzNamespace, NodeToken: nodeToken}
	var resp getNodeResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getPTZEndpoint(), "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetNode failed: %w", err)
	}

	return &PTZNode{
		Token:                  resp.PTZNode.Token,
		Name:                   resp.PTZNode.Name,
		HomeSupported:          resp.PTZNode.HomeSupported,
		MaximumNumberOfPresets: resp.PTZNode.MaximumNumberOfPresets,
	}, nil
}

// GetPTZConfigurationOptions returns available options for a PTZ configuration.
func (c *Client) GetPTZConfigurationOptions(ctx context.Context, configurationToken string) (*PTZConfigurationOptions, error) {
	type getConfigurationOptionsRequest struct {
		XMLName            xml.Name `xml:"tptz:GetConfigurationOptions"`
		Xmlns              string   `xml:"xmlns:tptz,attr"`
		ConfigurationToken string   `xml:"tptz:ConfigurationToken"`
	}

	type getConfigurationOptionsResponse struct {
		PTZConfigurationOptions struct {
			PTZTimeout *struct {
				Min string `xml:"Min"`
				Max string `xml:"Max"`
			} `xml:"PTZTimeout"`
		} `xml:"PTZConfigurationOptions"`
	}

	req := getConfigurationOptionsRequest{Xmlns: ptzNamespace, ConfigurationToken: configurationToken}
	var resp getConfigurationOptionsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getPTZEndpoint(), "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetPTZConfigurationOptions failed: %w", err)
	}

	opts := &PTZConfigurationOptions{}
	if resp.PTZConfigurationOptions.PTZTimeout != nil {
		opts.PTZTimeout = &struct {
			Min string
			Max string
		}{
			Min: resp.PTZConfigurationOptions.PTZTimeout.Min,
			Max: resp.PTZConfigurationOptions.PTZTimeout.Max,
		}
	}

	return opts, nil
}

// SetPTZConfiguration sets a PTZ configuration.
func (c *Client) SetPTZConfiguration(ctx context.Context, config *PTZConfiguration, forcePersistence bool) error {
	type setPTZConfigurationRequest struct {
		XMLName          xml.Name `xml:"tptz:SetConfiguration"`
		Xmlns            string   `xml:"xmlns:tptz,attr"`
		PTZConfiguration struct {
			Token     string `xml:"token,attr"`
			Name      string `xml:"Name,omitempty"`
			NodeToken string `xml:"NodeToken,omitempty"`
		} `xml:"tptz:PTZConfiguration"`
		ForcePersistence bool `xml:"tptz:ForcePersistence"`
	}

	type setPTZConfigurationResponse struct{}

	req := setPTZConfigurationRequest{
		Xmlns:            ptzNamespace,
		ForcePersistence: forcePersistence,
	}
	req.PTZConfiguration.Token = config.Token
	req.PTZConfiguration.Name = config.Name
	req.PTZConfiguration.NodeToken = config.NodeToken

	var resp setPTZConfigurationResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getPTZEndpoint(), "", req, &resp); err != nil {
		return fmt.Errorf("SetPTZConfiguration failed: %w", err)
	}

	return nil
}

// GetPTZServiceCapabilities returns the PTZ service capabilities.
func (c *Client) GetPTZServiceCapabilities(ctx context.Context) (*PTZServiceCapabilities, error) {
	type getServiceCapabilitiesRequest struct {
		XMLName xml.Name `xml:"tptz:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:tptz,attr"`
	}

	type getServiceCapabilitiesResponse struct {
		Capabilities struct {
			EFlip   bool `xml:"EFlip,attr"`
			Reverse bool `xml:"Reverse,attr"`
		} `xml:"Capabilities"`
	}

	req := getServiceCapabilitiesRequest{Xmlns: ptzNamespace}
	var resp getServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getPTZEndpoint(), "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetPTZServiceCapabilities failed: %w", err)
	}

	return &PTZServiceCapabilities{
		EFlip:   resp.Capabilities.EFlip,
		Reverse: resp.Capabilities.Reverse,
	}, nil
}

// GetCompatiblePTZConfigurationsForProfile returns PTZ configurations compatible with a media profile.
func (c *Client) GetCompatiblePTZConfigurationsForProfile(ctx context.Context, profileToken string) ([]*PTZConfiguration, error) {
	type getCompatibleConfigurationsRequest struct {
		XMLName      xml.Name `xml:"tptz:GetCompatibleConfigurations"`
		Xmlns        string   `xml:"xmlns:tptz,attr"`
		ProfileToken string   `xml:"tptz:ProfileToken"`
	}

	type getCompatibleConfigurationsResponse struct {
		PTZConfiguration []struct {
			Token     string `xml:"token,attr"`
			Name      string `xml:"Name"`
			NodeToken string `xml:"NodeToken"`
		} `xml:"PTZConfiguration"`
	}

	req := getCompatibleConfigurationsRequest{Xmlns: ptzNamespace, ProfileToken: profileToken}
	var resp getCompatibleConfigurationsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getPTZEndpoint(), "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetCompatiblePTZConfigurationsForProfile failed: %w", err)
	}

	configs := make([]*PTZConfiguration, 0, len(resp.PTZConfiguration))
	for _, c := range resp.PTZConfiguration {
		configs = append(configs, &PTZConfiguration{
			Token:     c.Token,
			Name:      c.Name,
			NodeToken: c.NodeToken,
		})
	}

	return configs, nil
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go test -run "TestGetNodes|TestGetNode|TestGetPTZConfigurationOptions|TestSetPTZConfiguration|TestGetPTZServiceCapabilities|TestGetCompatiblePTZConfigurations" -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add ptz.go ptz_test.go types.go
git commit -m "feat: add PTZ node and configuration operations (GetNodes, GetNode, GetConfigurationOptions, SetConfiguration, GetServiceCapabilities, GetCompatibleConfigurations)"
```

---

## Task 5: PTZ Service — Preset Tour Operations (7 ops)

**Files:**
- Modify: `ptz.go`
- Modify: `ptz_test.go`
- Modify: `types.go`

This task adds: `GetPresetTours`, `GetPresetTour`, `GetPresetTourOptions`, `CreatePresetTour`, `ModifyPresetTour`, `OperatePresetTour`, `RemovePresetTour`.

The pattern is identical to Task 4. Each operation follows the same structure:
1. Inline request struct with `xml:"tptz:OperationName"` and `xmlns:tptz` attr
2. Inline response struct
3. Public method calling `soapClient.Call()` with `ptzNamespace`
4. Map to public types

- [ ] **Step 1: Add PresetTour types to types.go**

```go
// PresetTour represents a PTZ preset tour.
type PresetTour struct {
	Token             string
	Name              string
	Status            string
	AutoStart         bool
	StartingCondition *PresetTourStartingCondition
	TourSpot          []*PresetTourSpot
}

// PresetTourStartingCondition defines when a preset tour starts.
type PresetTourStartingCondition struct {
	RecurringTime    *int
	RecurringDuration string
	Direction        string
}

// PresetTourSpot represents a stop in a preset tour.
type PresetTourSpot struct {
	PresetDetail *PresetTourPresetDetail
	Speed        *PTZSpeed
	StayTime     string
}

// PresetTourPresetDetail specifies which preset to visit.
type PresetTourPresetDetail struct {
	PresetToken string
	Home        bool
}

// PTZPresetTourOptions represents options for preset tours.
type PTZPresetTourOptions struct {
	AutoStart bool
	StartingCondition *struct {
		RecurringTimeRange *IntRange
		RecurringDurationRange *struct {
			Min string
			Max string
		}
	}
}
```

- [ ] **Step 2: Write tests for all 7 preset tour operations**

Add to `ptz_test.go`. Each test follows the pattern from Task 4. Example for `GetPresetTours`:

```go
func TestGetPresetTours(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<tptz:GetPresetToursResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
<tptz:PresetTour token="Tour1">
<tt:Name xmlns:tt="http://www.onvif.org/ver10/schema">Default Tour</tt:Name>
<tt:Status xmlns:tt="http://www.onvif.org/ver10/schema"><tt:State>Idle</tt:State></tt:Status>
<tt:AutoStart xmlns:tt="http://www.onvif.org/ver10/schema">false</tt:AutoStart>
</tptz:PresetTour>
</tptz:GetPresetToursResponse>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	tours, err := client.GetPresetTours(context.Background(), "Profile1")
	if err != nil {
		t.Fatalf("GetPresetTours failed: %v", err)
	}
	if len(tours) != 1 {
		t.Fatalf("expected 1 tour, got %d", len(tours))
	}
	if tours[0].Token != "Tour1" {
		t.Errorf("expected token Tour1, got %s", tours[0].Token)
	}
}

func TestCreatePresetTour(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<tptz:CreatePresetTourResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
<tptz:PresetTourToken>NewTour1</tptz:PresetTourToken>
</tptz:CreatePresetTourResponse>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	token, err := client.CreatePresetTour(context.Background(), "Profile1")
	if err != nil {
		t.Fatalf("CreatePresetTour failed: %v", err)
	}
	if token != "NewTour1" {
		t.Errorf("expected token NewTour1, got %s", token)
	}
}

func TestRemovePresetTour(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<tptz:RemovePresetTourResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl"/>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	err := client.RemovePresetTour(context.Background(), "Profile1", "Tour1")
	if err != nil {
		t.Fatalf("RemovePresetTour failed: %v", err)
	}
}
```

Write similar tests for `GetPresetTour`, `GetPresetTourOptions`, `ModifyPresetTour`, `OperatePresetTour`.

- [ ] **Step 3: Implement all 7 preset tour methods in ptz.go**

Each follows the standard pattern. Example signatures:

```go
func (c *Client) GetPresetTours(ctx context.Context, profileToken string) ([]*PresetTour, error)
func (c *Client) GetPresetTour(ctx context.Context, profileToken, presetTourToken string) (*PresetTour, error)
func (c *Client) GetPresetTourOptions(ctx context.Context, profileToken string, presetTourToken string) (*PTZPresetTourOptions, error)
func (c *Client) CreatePresetTour(ctx context.Context, profileToken string) (string, error)
func (c *Client) ModifyPresetTour(ctx context.Context, profileToken string, presetTour *PresetTour) error
func (c *Client) OperatePresetTour(ctx context.Context, profileToken, presetTourToken, operation string) error
func (c *Client) RemovePresetTour(ctx context.Context, profileToken, presetTourToken string) error
```

Request structs use `tptz:` prefix. All call `c.getPTZEndpoint()`.

- [ ] **Step 4: Run tests**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go test -run "TestGetPresetTour|TestCreatePresetTour|TestRemovePresetTour" -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add ptz.go ptz_test.go types.go
git commit -m "feat: add PTZ preset tour operations (GetPresetTours, CreatePresetTour, ModifyPresetTour, OperatePresetTour, RemovePresetTour)"
```

---

## Task 6: PTZ Service — Advanced Movement Operations (3 ops)

**Files:**
- Modify: `ptz.go`
- Modify: `ptz_test.go`

Adds: `SendAuxiliaryCommand` (PTZ namespace), `GeoMove`, `MoveAndStartTracking`.

- [ ] **Step 1: Write tests**

```go
func TestPTZSendAuxiliaryCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), "tptz:SendAuxiliaryCommand") {
			t.Error("expected tptz namespace for SendAuxiliaryCommand")
		}
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<tptz:SendAuxiliaryCommandResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl">
<tptz:AuxiliaryResponse>OK</tptz:AuxiliaryResponse>
</tptz:SendAuxiliaryCommandResponse>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	resp, err := client.PTZSendAuxiliaryCommand(context.Background(), "Profile1", "tt:Wiper|On")
	if err != nil {
		t.Fatalf("PTZSendAuxiliaryCommand failed: %v", err)
	}
	if resp != "OK" {
		t.Errorf("expected OK, got %s", resp)
	}
}

func TestGeoMove(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<tptz:GeoMoveResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl"/>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	err := client.GeoMove(context.Background(), "Profile1", &GeoLocation{
		Longitude: -122.4194,
		Latitude:  37.7749,
	}, nil, nil, nil)
	if err != nil {
		t.Fatalf("GeoMove failed: %v", err)
	}
}

func TestMoveAndStartTracking(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
<soap:Body>
<tptz:MoveAndStartTrackingResponse xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl"/>
</soap:Body></soap:Envelope>`
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL)
	err := client.MoveAndStartTracking(context.Background(), &MoveAndStartTrackingRequest{
		ProfileToken: "Profile1",
		ObjectID:     intPtr(42),
	})
	if err != nil {
		t.Fatalf("MoveAndStartTracking failed: %v", err)
	}
}
```

- [ ] **Step 2: Add types and implement methods**

Add to `types.go`:

```go
// GeoLocation represents a geographic location.
type GeoLocation struct {
	Longitude float64
	Latitude  float64
	Elevation *float64
}

// MoveAndStartTrackingRequest contains parameters for MoveAndStartTracking.
type MoveAndStartTrackingRequest struct {
	ProfileToken   string
	PresetToken    *string
	GeoLocation    *GeoLocation
	TargetPosition *PTZVector
	Speed          *PTZSpeed
	ObjectID       *int
}
```

Implement in `ptz.go`:

```go
// PTZSendAuxiliaryCommand sends an auxiliary command via the PTZ service.
func (c *Client) PTZSendAuxiliaryCommand(ctx context.Context, profileToken, auxiliaryData string) (string, error) {
	type sendAuxRequest struct {
		XMLName      xml.Name `xml:"tptz:SendAuxiliaryCommand"`
		Xmlns        string   `xml:"xmlns:tptz,attr"`
		ProfileToken string   `xml:"tptz:ProfileToken"`
		AuxiliaryData string  `xml:"tptz:AuxiliaryData"`
	}

	type sendAuxResponse struct {
		AuxiliaryResponse string `xml:"AuxiliaryResponse"`
	}

	req := sendAuxRequest{Xmlns: ptzNamespace, ProfileToken: profileToken, AuxiliaryData: auxiliaryData}
	var resp sendAuxResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getPTZEndpoint(), "", req, &resp); err != nil {
		return "", fmt.Errorf("PTZSendAuxiliaryCommand failed: %w", err)
	}

	return resp.AuxiliaryResponse, nil
}

// GeoMove moves the PTZ to a geographic location.
func (c *Client) GeoMove(ctx context.Context, profileToken string, target *GeoLocation, speed *PTZSpeed, areaHeight, areaWidth *float64) error {
	type geoMoveRequest struct {
		XMLName      xml.Name `xml:"tptz:GeoMove"`
		Xmlns        string   `xml:"xmlns:tptz,attr"`
		ProfileToken string   `xml:"tptz:ProfileToken"`
		Target       struct {
			Longitude float64 `xml:"lon,attr"`
			Latitude  float64 `xml:"lat,attr"`
		} `xml:"tptz:Target"`
		AreaHeight *float64 `xml:"tptz:AreaHeight,omitempty"`
		AreaWidth  *float64 `xml:"tptz:AreaWidth,omitempty"`
	}

	type geoMoveResponse struct{}

	req := geoMoveRequest{Xmlns: ptzNamespace, ProfileToken: profileToken}
	req.Target.Longitude = target.Longitude
	req.Target.Latitude = target.Latitude
	req.AreaHeight = areaHeight
	req.AreaWidth = areaWidth

	var resp geoMoveResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getPTZEndpoint(), "", req, &resp); err != nil {
		return fmt.Errorf("GeoMove failed: %w", err)
	}

	return nil
}

// MoveAndStartTracking moves PTZ and starts tracking an object.
func (c *Client) MoveAndStartTracking(ctx context.Context, request *MoveAndStartTrackingRequest) error {
	type moveAndStartTrackingRequest struct {
		XMLName      xml.Name `xml:"tptz:MoveAndStartTracking"`
		Xmlns        string   `xml:"xmlns:tptz,attr"`
		ProfileToken string   `xml:"tptz:ProfileToken"`
		PresetToken  *string  `xml:"tptz:PresetToken,omitempty"`
		ObjectID     *int     `xml:"tptz:ObjectID,omitempty"`
	}

	type moveAndStartTrackingResponse struct{}

	req := moveAndStartTrackingRequest{
		Xmlns:        ptzNamespace,
		ProfileToken: request.ProfileToken,
		PresetToken:  request.PresetToken,
		ObjectID:     request.ObjectID,
	}

	var resp moveAndStartTrackingResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getPTZEndpoint(), "", req, &resp); err != nil {
		return fmt.Errorf("MoveAndStartTracking failed: %w", err)
	}

	return nil
}
```

Add helper to test file:

```go
func intPtr(i int) *int { return &i }
```

- [ ] **Step 3: Run tests**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go test -run "TestPTZSendAuxiliaryCommand|TestGeoMove|TestMoveAndStartTracking" -v`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add ptz.go ptz_test.go types.go
git commit -m "feat: add PTZ advanced movement operations (SendAuxiliaryCommand, GeoMove, MoveAndStartTracking)"
```

---

## Task 7: Imaging Service — Missing Operations (4 ops)

**Files:**
- Modify: `imaging.go`
- Modify: `imaging_test.go`
- Modify: `types.go`

Adds: `GetImagingServiceCapabilities`, `GetImagingPresets`, `GetCurrentImagingPreset`, `SetCurrentImagingPreset`.

- [ ] **Step 1: Add types**

```go
// ImagingServiceCapabilities represents imaging service capabilities.
type ImagingServiceCapabilities struct {
	ImageStabilization bool
	Presets            bool
}

// ImagingPreset represents an imaging preset.
type ImagingPreset struct {
	Token string
	Name  string
	Type  string
}
```

- [ ] **Step 2: Write tests and implement**

Follow the exact same pattern as Tasks 4-6. Request structs use `timg:` prefix, call `c.imagingEndpoint` with fallback.

Signatures:
```go
func (c *Client) GetImagingServiceCapabilities(ctx context.Context) (*ImagingServiceCapabilities, error)
func (c *Client) GetImagingPresets(ctx context.Context, videoSourceToken string) ([]*ImagingPreset, error)
func (c *Client) GetCurrentImagingPreset(ctx context.Context, videoSourceToken string) (*ImagingPreset, error)
func (c *Client) SetCurrentImagingPreset(ctx context.Context, videoSourceToken, presetToken string) error
```

- [ ] **Step 3: Run tests**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go test -run "TestGetImagingServiceCapabilities|TestGetImagingPresets|TestGetCurrentImagingPreset|TestSetCurrentImagingPreset" -v`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add imaging.go imaging_test.go types.go
git commit -m "feat: add imaging service capabilities and preset operations"
```

---

## Task 8: Device Service — Missing Operations (8 ops)

**Files:**
- Modify: `device.go`, `device_extended.go`, `device_security.go`
- Modify: corresponding `*_test.go` files

Adds: `SetNetworkInterfaces`, `UpgradeSystemFirmware`, `GetUserRoles`, `SetUserRole`, `DeleteUserRole`, `GetAuthFailureWarningOptions`, `GetPasswordComplexityOptions`.

All follow the established pattern in their respective files using `tds:` namespace prefix.

- [ ] **Step 1: Implement SetNetworkInterfaces in device.go**

```go
// SetNetworkInterfaces configures a network interface.
func (c *Client) SetNetworkInterfaces(ctx context.Context, interfaceToken string, config *NetworkInterfaceSetConfiguration) (bool, error) {
	type setNetworkInterfacesRequest struct {
		XMLName        xml.Name `xml:"tds:SetNetworkInterfaces"`
		Xmlns          string   `xml:"xmlns:tds,attr"`
		InterfaceToken string   `xml:"tds:InterfaceToken"`
		NetworkInterface struct {
			Enabled *bool `xml:"Enabled,omitempty"`
			MTU     *int  `xml:"MTU,omitempty"`
			IPv4    *struct {
				Enabled *bool `xml:"Enabled,omitempty"`
				DHCP    *bool `xml:"DHCP,omitempty"`
			} `xml:"IPv4,omitempty"`
		} `xml:"tds:NetworkInterface"`
	}

	type setNetworkInterfacesResponse struct {
		RebootNeeded bool `xml:"RebootNeeded"`
	}

	req := setNetworkInterfacesRequest{
		Xmlns:          deviceNamespace,
		InterfaceToken: interfaceToken,
	}
	if config != nil {
		req.NetworkInterface.Enabled = config.Enabled
		req.NetworkInterface.MTU = config.MTU
		if config.IPv4 != nil {
			req.NetworkInterface.IPv4 = &struct {
				Enabled *bool `xml:"Enabled,omitempty"`
				DHCP    *bool `xml:"DHCP,omitempty"`
			}{
				Enabled: config.IPv4.Enabled,
				DHCP:    config.IPv4.DHCP,
			}
		}
	}

	var resp setNetworkInterfacesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.endpoint, "", req, &resp); err != nil {
		return false, fmt.Errorf("SetNetworkInterfaces failed: %w", err)
	}

	return resp.RebootNeeded, nil
}
```

Add type to `types.go`:

```go
// NetworkInterfaceSetConfiguration represents network interface settings to apply.
type NetworkInterfaceSetConfiguration struct {
	Enabled *bool
	MTU     *int
	IPv4    *IPv4NetworkInterfaceSetConfiguration
}

// IPv4NetworkInterfaceSetConfiguration represents IPv4 settings.
type IPv4NetworkInterfaceSetConfiguration struct {
	Enabled *bool
	DHCP    *bool
}
```

- [ ] **Step 2: Implement remaining device operations**

Add `UpgradeSystemFirmware`, `GetUserRoles`, `SetUserRole`, `DeleteUserRole` to `device_extended.go`.
Add `GetAuthFailureWarningOptions`, `GetPasswordComplexityOptions` to `device_security.go`.

All follow the same `tds:` namespace pattern. Each returns appropriate types.

Add types:

```go
// AuthFailureWarningOptions represents auth failure warning configuration ranges.
type AuthFailureWarningOptions struct {
	MonitorPeriodRange IntRange
	AuthFailureRange   IntRange
}

// PasswordComplexityOptions represents password complexity configuration options.
type PasswordComplexityOptions struct {
	MinLenRange                      *IntRange
	UppercaseRange                   *IntRange
	NumberRange                      *IntRange
	SpecialCharsRange                *IntRange
	BlockUsernameOccurrenceSupported *bool
	PolicyConfigurationLockSupported *bool
}

// UserRole represents a user role on the device.
type UserRole struct {
	Token       string
	RoleName    string
	Description string
}
```

- [ ] **Step 3: Write tests for all operations, run, verify pass**

- [ ] **Step 4: Commit**

```bash
git add device.go device_extended.go device_security.go device_test.go device_extended_test.go device_security_test.go types.go
git commit -m "feat: add missing device operations (SetNetworkInterfaces, UpgradeSystemFirmware, user roles, security options)"
```

---

## Task 9: Event Service — Missing Operations (8 ops)

**Files:**
- Modify: `event.go`
- Modify: `event_test.go`

Adds: `Subscribe`, `GetCurrentMessage`, `CreatePullPoint`, `DestroyPullPoint`, `GetMessages`, `PauseSubscription`, `ResumeSubscription`, `Notify` parser.

These use WS-BaseNotification namespaces (`wsnt:`) rather than the ONVIF event namespace. The SOAP actions are specified (unlike other ONVIF operations which use empty action).

- [ ] **Step 1: Implement Subscribe and base notification operations**

```go
const wsBaseNotificationNamespace = "http://docs.oasis-open.org/wsn/b-2"

// Subscribe creates a push-based event subscription (WS-BaseNotification).
func (c *Client) Subscribe(ctx context.Context, consumerReference string, filter string, terminationTime *time.Duration) (string, *time.Time, error) {
	type subscribeRequest struct {
		XMLName           xml.Name `xml:"wsnt:Subscribe"`
		Xmlns             string   `xml:"xmlns:wsnt,attr"`
		ConsumerReference struct {
			Address string `xml:"wsa:Address"`
			XmlnsWSA string `xml:"xmlns:wsa,attr"`
		} `xml:"wsnt:ConsumerReference"`
		Filter *struct {
			TopicExpression struct {
				Dialect string `xml:"Dialect,attr"`
				Value   string `xml:",chardata"`
			} `xml:"wsnt:TopicExpression,omitempty"`
		} `xml:"wsnt:Filter,omitempty"`
		InitialTerminationTime string `xml:"wsnt:InitialTerminationTime,omitempty"`
	}

	type subscribeResponse struct {
		SubscriptionReference struct {
			Address string `xml:"Address"`
		} `xml:"SubscriptionReference"`
		CurrentTime     string `xml:"CurrentTime"`
		TerminationTime string `xml:"TerminationTime"`
	}

	req := subscribeRequest{
		Xmlns: wsBaseNotificationNamespace,
	}
	req.ConsumerReference.Address = consumerReference
	req.ConsumerReference.XmlnsWSA = "http://www.w3.org/2005/08/addressing"

	if filter != "" {
		req.Filter = &struct {
			TopicExpression struct {
				Dialect string `xml:"Dialect,attr"`
				Value   string `xml:",chardata"`
			} `xml:"wsnt:TopicExpression,omitempty"`
		}{}
		req.Filter.TopicExpression.Dialect = "http://www.onvif.org/ver10/tev/topicExpression/ConcreteSet"
		req.Filter.TopicExpression.Value = filter
	}

	if terminationTime != nil {
		req.InitialTerminationTime = fmt.Sprintf("PT%dS", int(terminationTime.Seconds()))
	}

	var resp subscribeResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	endpoint := c.getEventEndpoint()
	if err := soapClient.Call(ctx, endpoint, "http://docs.oasis-open.org/wsn/bw-2/NotificationProducer/SubscribeRequest", req, &resp); err != nil {
		return "", nil, fmt.Errorf("Subscribe failed: %w", err)
	}

	var termTime *time.Time
	if resp.TerminationTime != "" {
		t, err := time.Parse(time.RFC3339, resp.TerminationTime)
		if err == nil {
			termTime = &t
		}
	}

	return resp.SubscriptionReference.Address, termTime, nil
}
```

- [ ] **Step 2: Implement remaining event operations (PauseSubscription, ResumeSubscription, GetCurrentMessage, legacy CreatePullPoint, DestroyPullPoint, GetMessages)**

All follow the same WS-BaseNotification pattern with specific SOAP actions.

- [ ] **Step 3: Write tests, run, verify pass**

- [ ] **Step 4: Commit**

```bash
git add event.go event_test.go
git commit -m "feat: add WS-BaseNotification event operations (Subscribe, PauseSubscription, ResumeSubscription, legacy pull point)"
```

---

## Task 10: DeviceIO Service — Missing Operations (15 ops)

**Files:**
- Modify: `deviceio.go`
- Modify: `deviceio_test.go`

Adds all 15 missing DeviceIO operations. These use `tmd:` namespace prefix and `deviceIONamespace`.

The operations fall into 4 groups with identical patterns:
- **Token list operations** (GetAudioSources, GetAudioOutputs, GetVideoSources): Empty request, response contains `Token` list
- **Get configuration**: Takes token param, returns configuration struct
- **Get configuration options**: Takes token param, returns options struct
- **Set configuration**: Takes configuration + ForcePersistence, returns empty

- [ ] **Step 1: Implement token list operations (3 ops)**

```go
// GetDeviceIOAudioSources returns audio source tokens from the DeviceIO service.
func (c *Client) GetDeviceIOAudioSources(ctx context.Context) ([]string, error) {
	type getAudioSourcesRequest struct {
		XMLName xml.Name `xml:"tmd:GetAudioSources"`
		Xmlns   string   `xml:"xmlns:tmd,attr"`
	}

	type getAudioSourcesResponse struct {
		Token []string `xml:"Token"`
	}

	req := getAudioSourcesRequest{Xmlns: deviceIONamespace}
	var resp getAudioSourcesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, c.getDeviceIOEndpoint(), "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetDeviceIOAudioSources failed: %w", err)
	}

	return resp.Token, nil
}
```

Repeat for `GetDeviceIOAudioOutputs` and `GetDeviceIOVideoSources`.

- [ ] **Step 2: Implement configuration get/set/options operations (12 ops)**

Follow the established pattern for each of:
- Audio source: `GetDeviceIOAudioSourceConfiguration`, `GetDeviceIOAudioSourceConfigurationOptions`, `SetDeviceIOAudioSourceConfiguration`
- Audio output: `GetDeviceIOAudioOutputConfiguration`, `GetDeviceIOAudioOutputConfigurationOptions`, `SetDeviceIOAudioOutputConfiguration`
- Video source: `GetDeviceIOVideoSourceConfiguration`, `GetDeviceIOVideoSourceConfigurationOptions`, `SetDeviceIOVideoSourceConfiguration`
- Relay: `GetDeviceIORelayOutputs`, `SetDeviceIORelayOutputState`, `SetDeviceIORelayOutputSettings`

Add a `getDeviceIOEndpoint()` method if not already present.

- [ ] **Step 3: Write tests for all 15 operations**

- [ ] **Step 4: Run all tests**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go test -run "TestGetDeviceIO|TestSetDeviceIO" -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add deviceio.go deviceio_test.go
git commit -m "feat: add missing DeviceIO operations (audio/video source config, relay outputs)"
```

---

## Task 11: Integration Test Scaffolding

**Files:**
- Modify: `ptz_real_camera_test.go` (or create if not exists)
- Modify: `imaging_real_camera_test.go`
- Modify: `event_real_camera_test.go`
- Modify: `deviceio_real_camera_test.go`

- [ ] **Step 1: Add PTZ integration tests**

```go
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

	t.Run("GetPresetTours", func(t *testing.T) {
		profiles, err := client.GetProfiles(ctx)
		if err != nil || len(profiles) == 0 {
			t.Skip("No profiles available")
		}
		tours, err := client.GetPresetTours(ctx, profiles[0].Token)
		if err != nil {
			t.Skipf("GetPresetTours not supported: %v", err)
		}
		t.Logf("Found %d preset tours", len(tours))
	})

	t.Run("GetPTZServiceCapabilities", func(t *testing.T) {
		caps, err := client.GetPTZServiceCapabilities(ctx)
		if err != nil {
			t.Skipf("GetPTZServiceCapabilities not supported: %v", err)
		}
		t.Logf("PTZ capabilities: EFlip=%v, Reverse=%v", caps.EFlip, caps.Reverse)
	})
}
```

- [ ] **Step 2: Add similar integration tests for Imaging, Event, DeviceIO**

Follow the same pattern: skip on unsupported, log results, chain token-dependent tests.

- [ ] **Step 3: Verify integration tests compile**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go vet -tags=real_camera ./...`
Expected: No errors.

- [ ] **Step 4: Commit**

```bash
git add *_real_camera_test.go
git commit -m "feat: add integration test scaffolding for PTZ, Imaging, Event, DeviceIO services"
```

---

## Task 12: Run Full Test Suite and Verify

- [ ] **Step 1: Run all unit tests**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go test ./... -count=1 -v 2>&1 | tail -50`
Expected: All tests PASS. No regressions.

- [ ] **Step 2: Run linter**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && make check`
Expected: No new lint errors. Fix any that appear.

- [ ] **Step 3: Verify build for all platforms**

Run: `cd /Users/ethanflower/personal_projects/onvif-go && go build ./...`
Expected: Clean build.

- [ ] **Step 4: Final commit if any lint fixes needed**

```bash
git add -A
git commit -m "fix: resolve lint issues from Phase 1+2 implementation"
```
