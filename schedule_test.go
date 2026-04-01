package onvif

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testScheduleXMLHeader = `<?xml version="1.0" encoding="UTF-8"?>`

func newMockScheduleServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml")

		body := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(body)
		bodyStr := string(body)

		var response string

		switch {
		case strings.Contains(bodyStr, "GetServiceCapabilities") && strings.Contains(bodyStr, "schedule"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:GetServiceCapabilitiesResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl">
      <tsc:Capabilities MaxLimit="100" MaxSchedules="10" MaxTimePeriodsPerDay="5"
        MaxSpecialDayGroups="5" MaxDaysInSpecialDayGroup="365" MaxSpecialDaysSchedules="3"
        ExtendedRecurrenceSupported="true" SpecialDaysSupported="true"
        StateReportingSupported="true" ClientSuppliedTokenSupported="true"/>
    </tsc:GetServiceCapabilitiesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetScheduleState"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:GetScheduleStateResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl">
      <tsc:ScheduleState>
        <tsc:Active>true</tsc:Active>
        <tsc:SpecialDay>false</tsc:SpecialDay>
      </tsc:ScheduleState>
    </tsc:GetScheduleStateResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetScheduleInfoList"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:GetScheduleInfoListResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl">
      <tsc:NextStartReference>ref_002</tsc:NextStartReference>
      <tsc:ScheduleInfo token="sched_001">
        <tsc:Name>Business Hours</tsc:Name>
        <tsc:Description>Monday to Friday 9am-5pm</tsc:Description>
      </tsc:ScheduleInfo>
      <tsc:ScheduleInfo token="sched_002">
        <tsc:Name>After Hours</tsc:Name>
      </tsc:ScheduleInfo>
    </tsc:GetScheduleInfoListResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetScheduleInfo"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:GetScheduleInfoResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl">
      <tsc:ScheduleInfo token="sched_001">
        <tsc:Name>Business Hours</tsc:Name>
        <tsc:Description>Monday to Friday 9am-5pm</tsc:Description>
      </tsc:ScheduleInfo>
    </tsc:GetScheduleInfoResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetScheduleList"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:GetScheduleListResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl">
      <tsc:NextStartReference>ref_002</tsc:NextStartReference>
      <tsc:Schedule token="sched_001">
        <tsc:Name>Business Hours</tsc:Name>
        <tsc:Description>Monday to Friday 9am-5pm</tsc:Description>
        <tsc:Standard>BEGIN:VCALENDAR&#xA;END:VCALENDAR</tsc:Standard>
      </tsc:Schedule>
    </tsc:GetScheduleListResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetSchedules"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:GetSchedulesResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl">
      <tsc:Schedule token="sched_001">
        <tsc:Name>Business Hours</tsc:Name>
        <tsc:Standard>BEGIN:VCALENDAR&#xA;END:VCALENDAR</tsc:Standard>
        <tsc:SpecialDays>
          <tsc:GroupToken>sdg_001</tsc:GroupToken>
          <tsc:TimeRange>
            <tsc:From>09:00:00</tsc:From>
            <tsc:Until>12:00:00</tsc:Until>
          </tsc:TimeRange>
        </tsc:SpecialDays>
      </tsc:Schedule>
    </tsc:GetSchedulesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "CreateSchedule"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:CreateScheduleResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl">
      <tsc:Token>sched_new_001</tsc:Token>
    </tsc:CreateScheduleResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "ModifySchedule"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:ModifyScheduleResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetSchedule"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:SetScheduleResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "DeleteSchedule"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:DeleteScheduleResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetSpecialDayGroupInfoList"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:GetSpecialDayGroupInfoListResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl">
      <tsc:NextStartReference>ref_sdg_002</tsc:NextStartReference>
      <tsc:SpecialDayGroupInfo token="sdg_001">
        <tsc:Name>Public Holidays</tsc:Name>
        <tsc:Description>National public holidays</tsc:Description>
      </tsc:SpecialDayGroupInfo>
    </tsc:GetSpecialDayGroupInfoListResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetSpecialDayGroupInfo"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:GetSpecialDayGroupInfoResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl">
      <tsc:SpecialDayGroupInfo token="sdg_001">
        <tsc:Name>Public Holidays</tsc:Name>
        <tsc:Description>National public holidays</tsc:Description>
      </tsc:SpecialDayGroupInfo>
    </tsc:GetSpecialDayGroupInfoResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetSpecialDayGroupList"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:GetSpecialDayGroupListResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl">
      <tsc:NextStartReference>ref_sdg_002</tsc:NextStartReference>
      <tsc:SpecialDayGroup token="sdg_001">
        <tsc:Name>Public Holidays</tsc:Name>
        <tsc:Days>BEGIN:VCALENDAR&#xA;END:VCALENDAR</tsc:Days>
      </tsc:SpecialDayGroup>
    </tsc:GetSpecialDayGroupListResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetSpecialDayGroups"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:GetSpecialDayGroupsResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl">
      <tsc:SpecialDayGroup token="sdg_001">
        <tsc:Name>Public Holidays</tsc:Name>
        <tsc:Description>National public holidays</tsc:Description>
        <tsc:Days>BEGIN:VCALENDAR&#xA;END:VCALENDAR</tsc:Days>
      </tsc:SpecialDayGroup>
    </tsc:GetSpecialDayGroupsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "CreateSpecialDayGroup"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:CreateSpecialDayGroupResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl">
      <tsc:Token>sdg_new_001</tsc:Token>
    </tsc:CreateSpecialDayGroupResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "ModifySpecialDayGroup"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:ModifySpecialDayGroupResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetSpecialDayGroup"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:SetSpecialDayGroupResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "DeleteSpecialDayGroup"):
			response = testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tsc:DeleteSpecialDayGroupResponse xmlns:tsc="http://www.onvif.org/ver10/schedule/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		default:
			http.Error(w, "unknown operation", http.StatusBadRequest)

			return
		}

		_, _ = w.Write([]byte(response))
	}))
}

func newMockScheduleFaultServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(testScheduleXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <SOAP-ENV:Fault>
      <SOAP-ENV:Code>
        <SOAP-ENV:Value>SOAP-ENV:Sender</SOAP-ENV:Value>
        <SOAP-ENV:Subcode>
          <SOAP-ENV:Value>ter:InvalidArgVal</SOAP-ENV:Value>
        </SOAP-ENV:Subcode>
      </SOAP-ENV:Code>
      <SOAP-ENV:Reason>
        <SOAP-ENV:Text xml:lang="en">Invalid token</SOAP-ENV:Text>
      </SOAP-ENV:Reason>
    </SOAP-ENV:Fault>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
	}))
}

func TestGetScheduleServiceCapabilities(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	caps, err := client.GetScheduleServiceCapabilities(context.Background())
	if err != nil {
		t.Fatalf("GetScheduleServiceCapabilities error: %v", err)
	}

	if caps.MaxLimit != 100 {
		t.Errorf("expected MaxLimit 100, got %d", caps.MaxLimit)
	}

	if caps.MaxSchedules != 10 {
		t.Errorf("expected MaxSchedules 10, got %d", caps.MaxSchedules)
	}

	if !caps.ExtendedRecurrenceSupported {
		t.Error("expected ExtendedRecurrenceSupported true")
	}

	if !caps.SpecialDaysSupported {
		t.Error("expected SpecialDaysSupported true")
	}

	if !caps.StateReportingSupported {
		t.Error("expected StateReportingSupported true")
	}

	if !caps.ClientSuppliedTokenSupported {
		t.Error("expected ClientSuppliedTokenSupported true")
	}
}

func TestGetScheduleServiceCapabilitiesFault(t *testing.T) {
	server := newMockScheduleFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.GetScheduleServiceCapabilities(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetScheduleState(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	state, err := client.GetScheduleState(context.Background(), "sched_001")
	if err != nil {
		t.Fatalf("GetScheduleState error: %v", err)
	}

	if !state.Active {
		t.Error("expected Active true")
	}

	if state.SpecialDay == nil || *state.SpecialDay {
		t.Error("expected SpecialDay false")
	}
}

func TestGetScheduleStateEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.GetScheduleState(context.Background(), "")
	if !errors.Is(err, ErrInvalidScheduleToken) {
		t.Errorf("expected ErrInvalidScheduleToken, got %v", err)
	}
}

func TestGetScheduleStateFault(t *testing.T) {
	server := newMockScheduleFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.GetScheduleState(context.Background(), "sched_001")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetScheduleInfoList(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	infos, nextRef, err := client.GetScheduleInfoList(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetScheduleInfoList error: %v", err)
	}

	if len(infos) != 2 {
		t.Errorf("expected 2 items, got %d", len(infos))
	}

	if infos[0].Token != "sched_001" {
		t.Errorf("expected token sched_001, got %s", infos[0].Token)
	}

	if infos[0].Name != "Business Hours" {
		t.Errorf("expected name Business Hours, got %s", infos[0].Name)
	}

	if nextRef != "ref_002" {
		t.Errorf("expected NextStartReference ref_002, got %s", nextRef)
	}
}

func TestGetScheduleInfo(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	infos, err := client.GetScheduleInfo(context.Background(), []string{"sched_001"})
	if err != nil {
		t.Fatalf("GetScheduleInfo error: %v", err)
	}

	if len(infos) != 1 {
		t.Errorf("expected 1 item, got %d", len(infos))
	}

	if infos[0].Token != "sched_001" {
		t.Errorf("expected token sched_001, got %s", infos[0].Token)
	}
}

func TestGetScheduleInfoEmptyTokens(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.GetScheduleInfo(context.Background(), []string{})
	if !errors.Is(err, ErrInvalidScheduleToken) {
		t.Errorf("expected ErrInvalidScheduleToken, got %v", err)
	}
}

func TestGetScheduleList(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	schedules, nextRef, err := client.GetScheduleList(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetScheduleList error: %v", err)
	}

	if len(schedules) != 1 {
		t.Errorf("expected 1 item, got %d", len(schedules))
	}

	if schedules[0].Token != "sched_001" {
		t.Errorf("expected token sched_001, got %s", schedules[0].Token)
	}

	if nextRef != "ref_002" {
		t.Errorf("expected NextStartReference ref_002, got %s", nextRef)
	}
}

func TestGetSchedules(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	schedules, err := client.GetSchedules(context.Background(), []string{"sched_001"})
	if err != nil {
		t.Fatalf("GetSchedules error: %v", err)
	}

	if len(schedules) != 1 {
		t.Errorf("expected 1 item, got %d", len(schedules))
	}

	if len(schedules[0].SpecialDays) != 1 {
		t.Errorf("expected 1 SpecialDays entry, got %d", len(schedules[0].SpecialDays))
	}

	if schedules[0].SpecialDays[0].GroupToken != "sdg_001" {
		t.Errorf("expected GroupToken sdg_001, got %s", schedules[0].SpecialDays[0].GroupToken)
	}
}

func TestGetSchedulesEmptyTokens(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.GetSchedules(context.Background(), []string{})
	if !errors.Is(err, ErrInvalidScheduleToken) {
		t.Errorf("expected ErrInvalidScheduleToken, got %v", err)
	}
}

func TestCreateSchedule(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	sched := &Schedule{
		ScheduleInfo: ScheduleInfo{Name: "New Schedule"},
		Standard:     "BEGIN:VCALENDAR\nEND:VCALENDAR",
	}

	token, err := client.CreateSchedule(context.Background(), sched)
	if err != nil {
		t.Fatalf("CreateSchedule error: %v", err)
	}

	if token != "sched_new_001" {
		t.Errorf("expected token sched_new_001, got %s", token)
	}
}

func TestCreateScheduleNil(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.CreateSchedule(context.Background(), nil)
	if !errors.Is(err, ErrScheduleNil) {
		t.Errorf("expected ErrScheduleNil, got %v", err)
	}
}

func TestCreateScheduleFault(t *testing.T) {
	server := newMockScheduleFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	sched := &Schedule{
		ScheduleInfo: ScheduleInfo{Name: "Test"},
		Standard:     "BEGIN:VCALENDAR\nEND:VCALENDAR",
	}

	_, err = client.CreateSchedule(context.Background(), sched)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestModifySchedule(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	sched := &Schedule{
		ScheduleInfo: ScheduleInfo{Token: "sched_001", Name: "Updated Schedule"},
		Standard:     "BEGIN:VCALENDAR\nEND:VCALENDAR",
	}

	if err := client.ModifySchedule(context.Background(), sched); err != nil {
		t.Fatalf("ModifySchedule error: %v", err)
	}
}

func TestModifyScheduleNil(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	if err := client.ModifySchedule(context.Background(), nil); !errors.Is(err, ErrScheduleNil) {
		t.Errorf("expected ErrScheduleNil, got %v", err)
	}
}

func TestModifyScheduleEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	sched := &Schedule{
		ScheduleInfo: ScheduleInfo{Name: "No Token"},
		Standard:     "BEGIN:VCALENDAR\nEND:VCALENDAR",
	}

	if err := client.ModifySchedule(context.Background(), sched); !errors.Is(err, ErrInvalidScheduleToken) {
		t.Errorf("expected ErrInvalidScheduleToken, got %v", err)
	}
}

func TestSetSchedule(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	sched := &Schedule{
		ScheduleInfo: ScheduleInfo{Token: "sched_001", Name: "Business Hours"},
		Standard:     "BEGIN:VCALENDAR\nEND:VCALENDAR",
	}

	if err := client.SetSchedule(context.Background(), sched); err != nil {
		t.Fatalf("SetSchedule error: %v", err)
	}
}

func TestSetScheduleNil(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	if err := client.SetSchedule(context.Background(), nil); !errors.Is(err, ErrScheduleNil) {
		t.Errorf("expected ErrScheduleNil, got %v", err)
	}
}

func TestDeleteSchedule(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	if err := client.DeleteSchedule(context.Background(), "sched_001"); err != nil {
		t.Fatalf("DeleteSchedule error: %v", err)
	}
}

func TestDeleteScheduleEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	if err := client.DeleteSchedule(context.Background(), ""); !errors.Is(err, ErrInvalidScheduleToken) {
		t.Errorf("expected ErrInvalidScheduleToken, got %v", err)
	}
}

func TestDeleteScheduleFault(t *testing.T) {
	server := newMockScheduleFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = client.DeleteSchedule(context.Background(), "sched_001")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetSpecialDayGroupInfoList(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	infos, nextRef, err := client.GetSpecialDayGroupInfoList(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetSpecialDayGroupInfoList error: %v", err)
	}

	if len(infos) != 1 {
		t.Errorf("expected 1 item, got %d", len(infos))
	}

	if infos[0].Token != "sdg_001" {
		t.Errorf("expected token sdg_001, got %s", infos[0].Token)
	}

	if nextRef != "ref_sdg_002" {
		t.Errorf("expected NextStartReference ref_sdg_002, got %s", nextRef)
	}
}

func TestGetSpecialDayGroupInfo(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	infos, err := client.GetSpecialDayGroupInfo(context.Background(), []string{"sdg_001"})
	if err != nil {
		t.Fatalf("GetSpecialDayGroupInfo error: %v", err)
	}

	if len(infos) != 1 {
		t.Errorf("expected 1 item, got %d", len(infos))
	}

	if infos[0].Name != "Public Holidays" {
		t.Errorf("expected name Public Holidays, got %s", infos[0].Name)
	}
}

func TestGetSpecialDayGroupInfoEmptyTokens(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.GetSpecialDayGroupInfo(context.Background(), []string{})
	if !errors.Is(err, ErrInvalidSpecialDayGroupToken) {
		t.Errorf("expected ErrInvalidSpecialDayGroupToken, got %v", err)
	}
}

func TestGetSpecialDayGroupList(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	groups, nextRef, err := client.GetSpecialDayGroupList(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("GetSpecialDayGroupList error: %v", err)
	}

	if len(groups) != 1 {
		t.Errorf("expected 1 item, got %d", len(groups))
	}

	if groups[0].Token != "sdg_001" {
		t.Errorf("expected token sdg_001, got %s", groups[0].Token)
	}

	if nextRef != "ref_sdg_002" {
		t.Errorf("expected NextStartReference ref_sdg_002, got %s", nextRef)
	}
}

func TestGetSpecialDayGroups(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	groups, err := client.GetSpecialDayGroups(context.Background(), []string{"sdg_001"})
	if err != nil {
		t.Fatalf("GetSpecialDayGroups error: %v", err)
	}

	if len(groups) != 1 {
		t.Errorf("expected 1 item, got %d", len(groups))
	}

	if groups[0].Description != "National public holidays" {
		t.Errorf("expected description, got %s", groups[0].Description)
	}
}

func TestGetSpecialDayGroupsEmptyTokens(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.GetSpecialDayGroups(context.Background(), []string{})
	if !errors.Is(err, ErrInvalidSpecialDayGroupToken) {
		t.Errorf("expected ErrInvalidSpecialDayGroupToken, got %v", err)
	}
}

func TestCreateSpecialDayGroup(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	group := &SpecialDayGroup{
		SpecialDayGroupInfo: SpecialDayGroupInfo{Name: "New Holidays"},
		Days:                "BEGIN:VCALENDAR\nEND:VCALENDAR",
	}

	token, err := client.CreateSpecialDayGroup(context.Background(), group)
	if err != nil {
		t.Fatalf("CreateSpecialDayGroup error: %v", err)
	}

	if token != "sdg_new_001" {
		t.Errorf("expected token sdg_new_001, got %s", token)
	}
}

func TestCreateSpecialDayGroupNil(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.CreateSpecialDayGroup(context.Background(), nil)
	if !errors.Is(err, ErrSpecialDayGroupNil) {
		t.Errorf("expected ErrSpecialDayGroupNil, got %v", err)
	}
}

func TestCreateSpecialDayGroupFault(t *testing.T) {
	server := newMockScheduleFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	group := &SpecialDayGroup{
		SpecialDayGroupInfo: SpecialDayGroupInfo{Name: "Test"},
	}

	_, err = client.CreateSpecialDayGroup(context.Background(), group)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestModifySpecialDayGroup(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	group := &SpecialDayGroup{
		SpecialDayGroupInfo: SpecialDayGroupInfo{Token: "sdg_001", Name: "Updated Holidays"},
	}

	if err := client.ModifySpecialDayGroup(context.Background(), group); err != nil {
		t.Fatalf("ModifySpecialDayGroup error: %v", err)
	}
}

func TestModifySpecialDayGroupNil(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	if err := client.ModifySpecialDayGroup(context.Background(), nil); !errors.Is(err, ErrSpecialDayGroupNil) {
		t.Errorf("expected ErrSpecialDayGroupNil, got %v", err)
	}
}

func TestModifySpecialDayGroupEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	group := &SpecialDayGroup{
		SpecialDayGroupInfo: SpecialDayGroupInfo{Name: "No Token"},
	}

	if err := client.ModifySpecialDayGroup(context.Background(), group); !errors.Is(err, ErrInvalidSpecialDayGroupToken) {
		t.Errorf("expected ErrInvalidSpecialDayGroupToken, got %v", err)
	}
}

func TestSetSpecialDayGroup(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	group := &SpecialDayGroup{
		SpecialDayGroupInfo: SpecialDayGroupInfo{Token: "sdg_001", Name: "Holidays"},
	}

	if err := client.SetSpecialDayGroup(context.Background(), group); err != nil {
		t.Fatalf("SetSpecialDayGroup error: %v", err)
	}
}

func TestSetSpecialDayGroupNil(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	if err := client.SetSpecialDayGroup(context.Background(), nil); !errors.Is(err, ErrSpecialDayGroupNil) {
		t.Errorf("expected ErrSpecialDayGroupNil, got %v", err)
	}
}

func TestDeleteSpecialDayGroup(t *testing.T) {
	server := newMockScheduleServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	if err := client.DeleteSpecialDayGroup(context.Background(), "sdg_001"); err != nil {
		t.Fatalf("DeleteSpecialDayGroup error: %v", err)
	}
}

func TestDeleteSpecialDayGroupEmptyToken(t *testing.T) {
	client, err := NewClient("http://localhost", WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	if err := client.DeleteSpecialDayGroup(context.Background(), ""); !errors.Is(err, ErrInvalidSpecialDayGroupToken) {
		t.Errorf("expected ErrInvalidSpecialDayGroupToken, got %v", err)
	}
}

func TestDeleteSpecialDayGroupFault(t *testing.T) {
	server := newMockScheduleFaultServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = client.DeleteSpecialDayGroup(context.Background(), "sdg_001")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
