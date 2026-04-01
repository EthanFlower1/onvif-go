package onvif

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testDoorControlXMLHeader = `<?xml version="1.0" encoding="UTF-8"?>`

func newMockDoorControlServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml")

		body := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(body)
		bodyStr := string(body)

		var response string

		switch {
		case strings.Contains(bodyStr, "GetServiceCapabilities") && strings.Contains(bodyStr, "doorcontrol"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:GetServiceCapabilitiesResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl">
      <tdc:Capabilities MaxLimit="100" MaxDoors="20"
        ClientSuppliedTokenSupported="true"
        DoorManagementSupported="true"/>
    </tdc:GetServiceCapabilitiesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetDoorInfoList"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:GetDoorInfoListResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl">
      <tdc:NextStartReference>ref_002</tdc:NextStartReference>
      <tdc:DoorInfo token="door_001">
        <tdc:Name>Main Entrance</tdc:Name>
        <tdc:Description>Front door</tdc:Description>
        <tdc:Capabilities Access="true" Lock="true" Unlock="true" DoorMonitor="true"/>
      </tdc:DoorInfo>
      <tdc:DoorInfo token="door_002">
        <tdc:Name>Side Entrance</tdc:Name>
        <tdc:Capabilities Access="false" Lock="true"/>
      </tdc:DoorInfo>
    </tdc:GetDoorInfoListResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetDoorInfo"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:GetDoorInfoResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl">
      <tdc:DoorInfo token="door_001">
        <tdc:Name>Main Entrance</tdc:Name>
        <tdc:Description>Front door</tdc:Description>
        <tdc:Capabilities Access="true" Lock="true" Unlock="true"/>
      </tdc:DoorInfo>
    </tdc:GetDoorInfoResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetDoorList"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:GetDoorListResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl">
      <tdc:NextStartReference>ref_002</tdc:NextStartReference>
      <tdc:Door token="door_001">
        <tdc:Name>Main Entrance</tdc:Name>
        <tdc:Capabilities Access="true" Lock="true"/>
        <tdc:DoorType>pt:Door</tdc:DoorType>
        <tdc:Timings>
          <tdc:ReleaseTime>PT5S</tdc:ReleaseTime>
          <tdc:OpenTime>PT30S</tdc:OpenTime>
        </tdc:Timings>
      </tdc:Door>
    </tdc:GetDoorListResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetDoors"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:GetDoorsResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl">
      <tdc:Door token="door_001">
        <tdc:Name>Main Entrance</tdc:Name>
        <tdc:Capabilities Access="true" Lock="true" Unlock="true"/>
        <tdc:DoorType>pt:Door</tdc:DoorType>
        <tdc:Timings>
          <tdc:ReleaseTime>PT5S</tdc:ReleaseTime>
          <tdc:OpenTime>PT30S</tdc:OpenTime>
        </tdc:Timings>
      </tdc:Door>
    </tdc:GetDoorsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "CreateDoor"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:CreateDoorResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl">
      <tdc:Token>door_new</tdc:Token>
    </tdc:CreateDoorResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetDoor"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:SetDoorResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "ModifyDoor"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:ModifyDoorResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "DeleteDoor"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:DeleteDoorResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetDoorState"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:GetDoorStateResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl">
      <tdc:DoorState>
        <tdc:DoorPhysicalState>Closed</tdc:DoorPhysicalState>
        <tdc:LockPhysicalState>Locked</tdc:LockPhysicalState>
        <tdc:Alarm>Normal</tdc:Alarm>
        <tdc:DoorMode>Locked</tdc:DoorMode>
      </tdc:DoorState>
    </tdc:GetDoorStateResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "AccessDoor"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:AccessDoorResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "LockDownReleaseDoor"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:LockDownReleaseDoorResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "LockDownDoor"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:LockDownDoorResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "LockOpenReleaseDoor"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:LockOpenReleaseDoorResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "LockOpenDoor"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:LockOpenDoorResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "DoubleLockDoor"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:DoubleLockDoorResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "BlockDoor"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:BlockDoorResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "UnlockDoor"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:UnlockDoorResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "LockDoor"):
			response = testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tdc:LockDoorResponse xmlns:tdc="http://www.onvif.org/ver10/doorcontrol/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		default:
			http.Error(w, "unknown operation", http.StatusBadRequest)

			return
		}

		_, _ = w.Write([]byte(response))
	}))
}

func newMockDoorControlFaultServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(testDoorControlXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <SOAP-ENV:Fault>
      <SOAP-ENV:Code>
        <SOAP-ENV:Value>SOAP-ENV:Sender</SOAP-ENV:Value>
      </SOAP-ENV:Code>
      <SOAP-ENV:Reason>
        <SOAP-ENV:Text>ter:InvalidArgVal</SOAP-ENV:Text>
      </SOAP-ENV:Reason>
    </SOAP-ENV:Fault>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
	}))
}

func TestGetDoorControlServiceCapabilities(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	caps, err := client.GetDoorControlServiceCapabilities(context.Background())
	if err != nil {
		t.Fatalf("GetDoorControlServiceCapabilities failed: %v", err)
	}

	if caps.MaxLimit != 100 {
		t.Errorf("expected MaxLimit=100, got %d", caps.MaxLimit)
	}

	if caps.MaxDoors != 20 {
		t.Errorf("expected MaxDoors=20, got %d", caps.MaxDoors)
	}

	if !caps.ClientSuppliedTokenSupported {
		t.Error("expected ClientSuppliedTokenSupported=true")
	}

	if !caps.DoorManagementSupported {
		t.Error("expected DoorManagementSupported=true")
	}
}

func TestGetDoorControlServiceCapabilitiesFault(t *testing.T) {
	server := newMockDoorControlFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	_, err = client.GetDoorControlServiceCapabilities(context.Background())
	if err == nil {
		t.Fatal("expected error for SOAP fault, got nil")
	}
}

func TestGetDoorInfoList(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	doors, nextRef, err := client.GetDoorInfoList(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetDoorInfoList failed: %v", err)
	}

	if len(doors) != 2 {
		t.Fatalf("expected 2 door info entries, got %d", len(doors))
	}

	if doors[0].Token != "door_001" {
		t.Errorf("expected token door_001, got %s", doors[0].Token)
	}

	if doors[0].Name != "Main Entrance" {
		t.Errorf("expected name 'Main Entrance', got %s", doors[0].Name)
	}

	if nextRef != "ref_002" {
		t.Errorf("expected NextStartReference ref_002, got %s", nextRef)
	}
}

func TestGetDoorInfoListFault(t *testing.T) {
	server := newMockDoorControlFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	_, _, err = client.GetDoorInfoList(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("expected error for SOAP fault, got nil")
	}
}

func TestGetDoorInfo(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	doors, err := client.GetDoorInfo(context.Background(), []string{"door_001"})
	if err != nil {
		t.Fatalf("GetDoorInfo failed: %v", err)
	}

	if len(doors) != 1 {
		t.Fatalf("expected 1 door info entry, got %d", len(doors))
	}

	if doors[0].Token != "door_001" {
		t.Errorf("expected token door_001, got %s", doors[0].Token)
	}
}

func TestGetDoorInfoEmptyTokens(t *testing.T) {
	client, err := NewClient("http://localhost:8080")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	_, err = client.GetDoorInfo(context.Background(), []string{})
	if !errors.Is(err, ErrInvalidDoorToken) {
		t.Errorf("expected ErrInvalidDoorToken, got %v", err)
	}
}

func TestGetDoorList(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	doors, nextRef, err := client.GetDoorList(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetDoorList failed: %v", err)
	}

	if len(doors) != 1 {
		t.Fatalf("expected 1 door entry, got %d", len(doors))
	}

	if doors[0].Token != "door_001" {
		t.Errorf("expected token door_001, got %s", doors[0].Token)
	}

	if doors[0].DoorType != "pt:Door" {
		t.Errorf("expected DoorType pt:Door, got %s", doors[0].DoorType)
	}

	if doors[0].Timings.ReleaseTime != "PT5S" {
		t.Errorf("expected ReleaseTime PT5S, got %s", doors[0].Timings.ReleaseTime)
	}

	if nextRef != "ref_002" {
		t.Errorf("expected NextStartReference ref_002, got %s", nextRef)
	}
}

func TestGetDoors(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	doors, err := client.GetDoors(context.Background(), []string{"door_001"})
	if err != nil {
		t.Fatalf("GetDoors failed: %v", err)
	}

	if len(doors) != 1 {
		t.Fatalf("expected 1 door, got %d", len(doors))
	}

	if doors[0].Token != "door_001" {
		t.Errorf("expected token door_001, got %s", doors[0].Token)
	}
}

func TestGetDoorsEmptyTokens(t *testing.T) {
	client, err := NewClient("http://localhost:8080")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	_, err = client.GetDoors(context.Background(), []string{})
	if !errors.Is(err, ErrInvalidDoorToken) {
		t.Errorf("expected ErrInvalidDoorToken, got %v", err)
	}
}

func TestCreateDoor(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	door := &Door{
		DoorInfo: DoorInfo{
			Name:        "New Door",
			Description: "A new door",
		},
		DoorType: "pt:Door",
		Timings: DoorTimings{
			ReleaseTime: "PT5S",
			OpenTime:    "PT30S",
		},
	}

	token, err := client.CreateDoor(context.Background(), door)
	if err != nil {
		t.Fatalf("CreateDoor failed: %v", err)
	}

	if token != "door_new" {
		t.Errorf("expected token door_new, got %s", token)
	}
}

func TestCreateDoorNil(t *testing.T) {
	client, err := NewClient("http://localhost:8080")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	_, err = client.CreateDoor(context.Background(), nil)
	if !errors.Is(err, ErrDoorNil) {
		t.Errorf("expected ErrDoorNil, got %v", err)
	}
}

func TestSetDoor(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	door := &Door{
		DoorInfo: DoorInfo{
			Token: "door_001",
			Name:  "Updated Door",
		},
		Timings: DoorTimings{
			ReleaseTime: "PT5S",
			OpenTime:    "PT30S",
		},
	}

	err = client.SetDoor(context.Background(), door)
	if err != nil {
		t.Fatalf("SetDoor failed: %v", err)
	}
}

func TestSetDoorEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost:8080")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	door := &Door{
		DoorInfo: DoorInfo{Token: ""},
	}

	err = client.SetDoor(context.Background(), door)
	if !errors.Is(err, ErrInvalidDoorToken) {
		t.Errorf("expected ErrInvalidDoorToken, got %v", err)
	}
}

func TestModifyDoor(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	door := &Door{
		DoorInfo: DoorInfo{
			Token: "door_001",
			Name:  "Modified Door",
		},
		Timings: DoorTimings{
			ReleaseTime: "PT5S",
			OpenTime:    "PT30S",
		},
	}

	err = client.ModifyDoor(context.Background(), door)
	if err != nil {
		t.Fatalf("ModifyDoor failed: %v", err)
	}
}

func TestModifyDoorNil(t *testing.T) {
	client, err := NewClient("http://localhost:8080")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	err = client.ModifyDoor(context.Background(), nil)
	if !errors.Is(err, ErrDoorNil) {
		t.Errorf("expected ErrDoorNil, got %v", err)
	}
}

func TestDeleteDoor(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	err = client.DeleteDoor(context.Background(), "door_001")
	if err != nil {
		t.Fatalf("DeleteDoor failed: %v", err)
	}
}

func TestDeleteDoorEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost:8080")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	err = client.DeleteDoor(context.Background(), "")
	if !errors.Is(err, ErrInvalidDoorToken) {
		t.Errorf("expected ErrInvalidDoorToken, got %v", err)
	}
}

func TestGetDoorState(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	state, err := client.GetDoorState(context.Background(), "door_001")
	if err != nil {
		t.Fatalf("GetDoorState failed: %v", err)
	}

	if state.DoorMode != DoorModeLocked {
		t.Errorf("expected DoorMode Locked, got %s", state.DoorMode)
	}

	if state.DoorPhysicalState == nil {
		t.Fatal("expected DoorPhysicalState to be set")
	}

	if *state.DoorPhysicalState != DoorPhysicalStateClosed {
		t.Errorf("expected DoorPhysicalState Closed, got %s", *state.DoorPhysicalState)
	}

	if state.LockPhysicalState == nil {
		t.Fatal("expected LockPhysicalState to be set")
	}

	if *state.LockPhysicalState != LockPhysicalStateLocked {
		t.Errorf("expected LockPhysicalState Locked, got %s", *state.LockPhysicalState)
	}

	if state.Alarm == nil {
		t.Fatal("expected Alarm to be set")
	}

	if *state.Alarm != DoorAlarmStateNormal {
		t.Errorf("expected Alarm Normal, got %s", *state.Alarm)
	}
}

func TestGetDoorStateFault(t *testing.T) {
	server := newMockDoorControlFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	_, err = client.GetDoorState(context.Background(), "door_001")
	if err == nil {
		t.Fatal("expected error for SOAP fault, got nil")
	}
}

func TestGetDoorStateEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost:8080")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	_, err = client.GetDoorState(context.Background(), "")
	if !errors.Is(err, ErrInvalidDoorToken) {
		t.Errorf("expected ErrInvalidDoorToken, got %v", err)
	}
}

func TestAccessDoor(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	err = client.AccessDoor(context.Background(), "door_001", nil, nil)
	if err != nil {
		t.Fatalf("AccessDoor failed: %v", err)
	}
}

func TestAccessDoorFault(t *testing.T) {
	server := newMockDoorControlFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	err = client.AccessDoor(context.Background(), "door_001", nil, nil)
	if err == nil {
		t.Fatal("expected error for SOAP fault, got nil")
	}
}

func TestAccessDoorEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost:8080")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	err = client.AccessDoor(context.Background(), "", nil, nil)
	if !errors.Is(err, ErrInvalidDoorToken) {
		t.Errorf("expected ErrInvalidDoorToken, got %v", err)
	}
}

func TestLockDoor(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	err = client.LockDoor(context.Background(), "door_001")
	if err != nil {
		t.Fatalf("LockDoor failed: %v", err)
	}
}

func TestLockDoorFault(t *testing.T) {
	server := newMockDoorControlFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	err = client.LockDoor(context.Background(), "door_001")
	if err == nil {
		t.Fatal("expected error for SOAP fault, got nil")
	}
}

func TestUnlockDoor(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	err = client.UnlockDoor(context.Background(), "door_001")
	if err != nil {
		t.Fatalf("UnlockDoor failed: %v", err)
	}
}

func TestBlockDoor(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	err = client.BlockDoor(context.Background(), "door_001")
	if err != nil {
		t.Fatalf("BlockDoor failed: %v", err)
	}
}

func TestLockDownDoor(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	err = client.LockDownDoor(context.Background(), "door_001")
	if err != nil {
		t.Fatalf("LockDownDoor failed: %v", err)
	}
}

func TestLockDownReleaseDoor(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	err = client.LockDownReleaseDoor(context.Background(), "door_001")
	if err != nil {
		t.Fatalf("LockDownReleaseDoor failed: %v", err)
	}
}

func TestLockOpenDoor(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	err = client.LockOpenDoor(context.Background(), "door_001")
	if err != nil {
		t.Fatalf("LockOpenDoor failed: %v", err)
	}
}

func TestLockOpenReleaseDoor(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	err = client.LockOpenReleaseDoor(context.Background(), "door_001")
	if err != nil {
		t.Fatalf("LockOpenReleaseDoor failed: %v", err)
	}
}

func TestDoubleLockDoor(t *testing.T) {
	server := newMockDoorControlServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	err = client.DoubleLockDoor(context.Background(), "door_001")
	if err != nil {
		t.Fatalf("DoubleLockDoor failed: %v", err)
	}
}

func TestDoubleLockDoorFault(t *testing.T) {
	server := newMockDoorControlFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.doorControlEndpoint = server.URL

	err = client.DoubleLockDoor(context.Background(), "door_001")
	if err == nil {
		t.Fatal("expected error for SOAP fault, got nil")
	}
}
