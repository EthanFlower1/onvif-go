# Phase 6: Remaining Services Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the remaining ONVIF services: Advanced Security, Thermal, Display, Provisioning, Uplink, and App Management. ~174 new operations.

**Architecture:** Six new service files with per-service type files. All services are fully independent and can be implemented in any order or in parallel.

**Tech Stack:** Go stdlib, internal/soap, httptest.

**Depends on:** Phase 1 (Foundation) must be complete. Independent of all other phases.

**Spec:** `docs/superpowers/specs/2026-03-31-onvif-full-compliance-design.md`

**WSDLs:**
- Advanced Security: `resources/specs/wsdl/ver10/advancedsecurity/wsdl/advancedsecurity.wsdl` (namespace: `http://www.onvif.org/ver10/advancedsecurity/wsdl`, prefix: `tas:`)
- Thermal: `resources/specs/wsdl/ver10/thermal/wsdl/thermal.wsdl` (namespace: `http://www.onvif.org/ver10/thermal/wsdl`, prefix: `tth:`)
- Display: `resources/specs/wsdl/ver10/display/wsdl/display.wsdl` (namespace: `http://www.onvif.org/ver10/display/wsdl`, prefix: `tls:`)
- Provisioning: `resources/specs/wsdl/ver10/provisioning/wsdl/provisioning.wsdl` (namespace: `http://www.onvif.org/ver10/provisioning/wsdl`, prefix: `tpv:`)
- Uplink: `resources/specs/wsdl/ver10/uplink/wsdl/uplink.wsdl` (namespace: `http://www.onvif.org/ver10/uplink/wsdl`, prefix: `tup:`)
- App Management: `resources/specs/wsdl/ver10/appmgmt/wsdl/appmgmt.wsdl` (namespace: `http://www.onvif.org/ver10/appmgmt/wsdl`, prefix: `tap:`)

---

## File Structure

**New files (6 triplets + 6 type files = 24 files):**
- `types_security.go` + `advanced_security.go` + `advanced_security_test.go` + `advanced_security_real_camera_test.go`
- `types_thermal.go` + `thermal.go` + `thermal_test.go` + `thermal_real_camera_test.go`
- `types_display.go` + `display.go` + `display_test.go` + `display_real_camera_test.go`
- `types_provisioning.go` + `provisioning.go` + `provisioning_test.go` + `provisioning_real_camera_test.go`
- `types_uplink.go` + `uplink.go` + `uplink_test.go` + `uplink_real_camera_test.go`
- `types_appmgmt.go` + `appmgmt.go` + `appmgmt_test.go` + `appmgmt_real_camera_test.go`

---

## Tasks

### Task 1: Advanced Security Service (~124 ops) — LARGEST SERVICE

This is the largest single service. Break into sub-tasks:

- [ ] **1a: Types** — Create `types_security.go` with KeyInfo, CertificateInfo, TLSConfiguration, CRL, CertificationPath types
- [ ] **1b: Key management** (~20 ops) — CreateRSAKeyPair, GetKeyStatus, GetPrivateKeyStatus, DeleteKey, GetAllKeys, etc.
- [ ] **1c: Certificate management** (~25 ops) — CreateSelfSignedCertificate, UploadCertificate, GetCertificate, GetAllCertificates, DeleteCertificate, CreatePKCS10CSR, etc.
- [ ] **1d: TLS configuration** (~15 ops) — GetTLSServerSupported, SetTLSServer, GetTLSClientSupported, SetTLSClient, etc.
- [ ] **1e: CRL management** (~10 ops) — UploadCRL, GetCRL, GetAllCRLs, DeleteCRL, etc.
- [ ] **1f: Certification path** (~15 ops) — CreateCertificationPath, GetCertificationPath, GetAllCertificationPaths, DeleteCertificationPath, etc.
- [ ] **1g: IEEE 802.1X** (~10 ops) — AddDot1XConfiguration, GetDot1XConfiguration, DeleteDot1XConfiguration, etc.
- [ ] **1h: Passphrase/keystore** (~15 ops) — UploadPassphrase, GetAllPassphrases, DeletePassphrase, etc.
- [ ] **1i: Remaining** (~14 ops) — GetServiceCapabilities, plus any remaining operations
- [ ] Write tests for all operations
- [ ] Commit

### Task 2: Thermal Service (~10 ops)
- [ ] Create `types_thermal.go` with ThermalConfiguration, Radiometry types
- [ ] Create `thermal.go` — GetServiceCapabilities, GetConfigurations, SetConfiguration, radiometry operations
- [ ] Write unit tests
- [ ] Commit

### Task 3: Display Service (~10 ops)
- [ ] Create `types_display.go` with Layout, Pane types
- [ ] Create `display.go` — GetServiceCapabilities, GetLayout, SetLayout, pane operations
- [ ] Write unit tests
- [ ] Commit

### Task 4: Provisioning Service (~10 ops)
- [ ] Create `types_provisioning.go` with provisioning-specific types
- [ ] Create `provisioning.go` — GetServiceCapabilities, PanMove, TiltMove, ZoomMove, RollMove, FocusMove, Stop, GetUsage
- [ ] Write unit tests
- [ ] Commit

### Task 5: Uplink Service (~10 ops)
- [ ] Create `types_uplink.go` with UplinkConnection types
- [ ] Create `uplink.go` — GetServiceCapabilities, GetUplinks, SetUplink, DeleteUplink
- [ ] Write unit tests
- [ ] Commit

### Task 6: App Management Service (~10 ops)
- [ ] Create `types_appmgmt.go` with App, AppInfo types
- [ ] Create `appmgmt.go` — GetServiceCapabilities, GetInstalledApps, InstallApp, UninstallApp, ActivateApp, DeactivateApp, GetAppStatus
- [ ] Write unit tests
- [ ] Commit

### Task 7: Integration Tests
- [ ] Create `*_real_camera_test.go` for all 6 services
- [ ] Verify compilation, run full test suite
- [ ] Commit
