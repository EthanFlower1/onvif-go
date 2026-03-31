package onvif

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testDeviceIOXMLHeader = `<?xml version="1.0" encoding="UTF-8"?>`

func newMockDeviceIOServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml")

		body := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(body)
		bodyStr := string(body)

		var response string

		switch {
		case strings.Contains(bodyStr, "GetServiceCapabilities") && strings.Contains(bodyStr, "deviceIO"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetServiceCapabilitiesResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:Capabilities 
        VideoSources="4"
        VideoOutputs="2"
        AudioSources="2"
        AudioOutputs="2"
        RelayOutputs="4"
        SerialPorts="2"
        DigitalInputs="8"
        DigitalInputOptions="true"
        SerialPortConfiguration="true"/>
    </tmd:GetServiceCapabilitiesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetDigitalInputConfigurationOptions"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetDigitalInputConfigurationOptionsResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:DigitalInputConfigurationOptions>
        <tmd:IdleState>open</tmd:IdleState>
        <tmd:IdleState>closed</tmd:IdleState>
      </tmd:DigitalInputConfigurationOptions>
    </tmd:GetDigitalInputConfigurationOptionsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetDigitalInputs"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetDigitalInputsResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:DigitalInputs token="input_001" IdleState="open"/>
      <tmd:DigitalInputs token="input_002" IdleState="closed"/>
    </tmd:GetDigitalInputsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetDigitalInputConfigurations"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:SetDigitalInputConfigurationsResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetVideoOutputs"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetVideoOutputsResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:VideoOutputs token="video_out_001">
        <tmd:Layout>
          <tt:Pane xmlns:tt="http://www.onvif.org/ver10/schema" Pane="main">
            <tt:Area bottom="1.0" top="0.0" right="1.0" left="0.0"/>
          </tt:Pane>
        </tmd:Layout>
        <tmd:Resolution>
          <tmd:Width>1920</tmd:Width>
          <tmd:Height>1080</tmd:Height>
        </tmd:Resolution>
        <tmd:RefreshRate>60.0</tmd:RefreshRate>
        <tmd:AspectRatio>16:9</tmd:AspectRatio>
      </tmd:VideoOutputs>
    </tmd:GetVideoOutputsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetSerialPortConfigurationOptions"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetSerialPortConfigurationOptionsResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:SerialPortConfigurationOptions token="serial_001">
        <tmd:BaudRateList><tmd:Items>9600</tmd:Items><tmd:Items>19200</tmd:Items><tmd:Items>38400</tmd:Items></tmd:BaudRateList>
        <tmd:ParityBitList><tmd:Items>None</tmd:Items><tmd:Items>Odd</tmd:Items><tmd:Items>Even</tmd:Items></tmd:ParityBitList>
        <tmd:CharacterLengthList><tmd:Items>7</tmd:Items><tmd:Items>8</tmd:Items></tmd:CharacterLengthList>
        <tmd:StopBitList><tmd:Items>1</tmd:Items><tmd:Items>2</tmd:Items></tmd:StopBitList>
      </tmd:SerialPortConfigurationOptions>
    </tmd:GetSerialPortConfigurationOptionsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetSerialPortConfiguration"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetSerialPortConfigurationResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:SerialPortConfiguration token="serial_001">
        <tmd:Type>RS232</tmd:Type>
        <tmd:BaudRate>9600</tmd:BaudRate>
        <tmd:ParityBit>None</tmd:ParityBit>
        <tmd:CharacterLength>8</tmd:CharacterLength>
        <tmd:StopBit>1</tmd:StopBit>
      </tmd:SerialPortConfiguration>
    </tmd:GetSerialPortConfigurationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetSerialPorts"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetSerialPortsResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:SerialPorts token="serial_001">
        <tmd:Type>RS232</tmd:Type>
      </tmd:SerialPorts>
      <tmd:SerialPorts token="serial_002">
        <tmd:Type>RS485</tmd:Type>
      </tmd:SerialPorts>
    </tmd:GetSerialPortsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetSerialPortConfiguration"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:SetSerialPortConfigurationResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SendReceiveSerialCommand"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:SendReceiveSerialCommandResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:SerialData>
        <tt:Binary xmlns:tt="http://www.onvif.org/ver10/schema">OK</tt:Binary>
      </tmd:SerialData>
    </tmd:SendReceiveSerialCommandResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetVideoOutputConfigurationOptions"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetVideoOutputConfigurationOptionsResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:VideoOutputConfigurationOptions>
        <tmd:Name Min="1" Max="64"/>
        <tmd:OutputTokensAvailable>video_out_001</tmd:OutputTokensAvailable>
        <tmd:OutputTokensAvailable>video_out_002</tmd:OutputTokensAvailable>
      </tmd:VideoOutputConfigurationOptions>
    </tmd:GetVideoOutputConfigurationOptionsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetVideoOutputConfiguration"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetVideoOutputConfigurationResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:VideoOutputConfiguration token="config_001">
        <tmd:Name>Main Output</tmd:Name>
        <tmd:UseCount>2</tmd:UseCount>
        <tmd:OutputToken>video_out_001</tmd:OutputToken>
      </tmd:VideoOutputConfiguration>
    </tmd:GetVideoOutputConfigurationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetVideoOutputConfiguration"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:SetVideoOutputConfigurationResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetRelayOutputOptions"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetRelayOutputOptionsResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:RelayOutputOptions token="relay_001">
        <tmd:Mode>Monostable</tmd:Mode>
        <tmd:Mode>Bistable</tmd:Mode>
        <tmd:DelayTimes>PT1S</tmd:DelayTimes>
        <tmd:DelayTimes>PT5S</tmd:DelayTimes>
        <tmd:DelayTimes>PT10S</tmd:DelayTimes>
        <tmd:Discrete>true</tmd:Discrete>
      </tmd:RelayOutputOptions>
    </tmd:GetRelayOutputOptionsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAudioSources"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetAudioSourcesResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:Token>AudioSrc1</tmd:Token>
      <tmd:Token>AudioSrc2</tmd:Token>
    </tmd:GetAudioSourcesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAudioOutputs"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetAudioOutputsResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:Token>AudioOut1</tmd:Token>
    </tmd:GetAudioOutputsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetVideoSources"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetVideoSourcesResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:Token>VideoSrc1</tmd:Token>
      <tmd:Token>VideoSrc2</tmd:Token>
      <tmd:Token>VideoSrc3</tmd:Token>
    </tmd:GetVideoSourcesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAudioSourceConfigurationOptions"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetAudioSourceConfigurationOptionsResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:AudioSourceOptions>
        <tmd:InputTokensAvailable>AudioSrc1</tmd:InputTokensAvailable>
      </tmd:AudioSourceOptions>
    </tmd:GetAudioSourceConfigurationOptionsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAudioSourceConfiguration"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetAudioSourceConfigurationResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:AudioSourceConfiguration token="cfg_audio_001">
        <tmd:Name>Microphone</tmd:Name>
        <tmd:UseCount>1</tmd:UseCount>
        <tmd:SourceToken>AudioSrc1</tmd:SourceToken>
      </tmd:AudioSourceConfiguration>
    </tmd:GetAudioSourceConfigurationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetAudioSourceConfiguration"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:SetAudioSourceConfigurationResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAudioOutputConfigurationOptions"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetAudioOutputConfigurationOptionsResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:AudioOutputOptions>
        <tmd:OutputTokensAvailable>AudioOut1</tmd:OutputTokensAvailable>
      </tmd:AudioOutputOptions>
    </tmd:GetAudioOutputConfigurationOptionsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetAudioOutputConfiguration"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetAudioOutputConfigurationResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:AudioOutputConfiguration token="cfg_out_001">
        <tmd:Name>Speaker</tmd:Name>
        <tmd:UseCount>2</tmd:UseCount>
        <tmd:OutputToken>AudioOut1</tmd:OutputToken>
      </tmd:AudioOutputConfiguration>
    </tmd:GetAudioOutputConfigurationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetAudioOutputConfiguration"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:SetAudioOutputConfigurationResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetVideoSourceConfigurationOptions"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetVideoSourceConfigurationOptionsResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:VideoSourceConfigurationOptions>
        <tmd:BoundsRange>
          <tmd:XRange><tmd:Min>0</tmd:Min><tmd:Max>1920</tmd:Max></tmd:XRange>
          <tmd:YRange><tmd:Min>0</tmd:Min><tmd:Max>1080</tmd:Max></tmd:YRange>
          <tmd:WidthRange><tmd:Min>320</tmd:Min><tmd:Max>1920</tmd:Max></tmd:WidthRange>
          <tmd:HeightRange><tmd:Min>240</tmd:Min><tmd:Max>1080</tmd:Max></tmd:HeightRange>
        </tmd:BoundsRange>
        <tmd:VideoSourceTokensAvailable>VideoSrc1</tmd:VideoSourceTokensAvailable>
      </tmd:VideoSourceConfigurationOptions>
    </tmd:GetVideoSourceConfigurationOptionsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetVideoSourceConfiguration"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetVideoSourceConfigurationResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:VideoSourceConfiguration token="cfg_vid_001">
        <tmd:Name>Main Camera</tmd:Name>
        <tmd:UseCount>3</tmd:UseCount>
        <tmd:SourceToken>VideoSrc1</tmd:SourceToken>
        <tmd:Bounds x="0" y="0" width="1920" height="1080"/>
      </tmd:VideoSourceConfiguration>
    </tmd:GetVideoSourceConfigurationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetVideoSourceConfiguration"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:SetVideoSourceConfigurationResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "GetRelayOutputs"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetRelayOutputsResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:RelayOutputs token="relay_001">
        <tmd:Properties>
          <tmd:Mode>Monostable</tmd:Mode>
          <tmd:DelayTime>PT5S</tmd:DelayTime>
          <tmd:IdleState>open</tmd:IdleState>
        </tmd:Properties>
      </tmd:RelayOutputs>
      <tmd:RelayOutputs token="relay_002">
        <tmd:Properties>
          <tmd:Mode>Bistable</tmd:Mode>
          <tmd:IdleState>closed</tmd:IdleState>
        </tmd:Properties>
      </tmd:RelayOutputs>
    </tmd:GetRelayOutputsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetRelayOutputState"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:SetRelayOutputStateResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		case strings.Contains(bodyStr, "SetRelayOutputSettings"):
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:SetRelayOutputSettingsResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

		default:
			response = testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <SOAP-ENV:Fault>
      <SOAP-ENV:Code><SOAP-ENV:Value>SOAP-ENV:Receiver</SOAP-ENV:Value></SOAP-ENV:Code>
      <SOAP-ENV:Reason><SOAP-ENV:Text>Unknown action</SOAP-ENV:Text></SOAP-ENV:Reason>
    </SOAP-ENV:Fault>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`
		}

		_, _ = w.Write([]byte(response))
	}))
}

func TestGetDeviceIOServiceCapabilities(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	caps, err := client.GetDeviceIOServiceCapabilities(ctx)
	if err != nil {
		t.Fatalf("GetDeviceIOServiceCapabilities failed: %v", err)
	}

	if caps.VideoSources != 4 {
		t.Errorf("Expected VideoSources to be 4, got %d", caps.VideoSources)
	}

	if caps.VideoOutputs != 2 {
		t.Errorf("Expected VideoOutputs to be 2, got %d", caps.VideoOutputs)
	}

	if caps.AudioSources != 2 {
		t.Errorf("Expected AudioSources to be 2, got %d", caps.AudioSources)
	}

	if caps.AudioOutputs != 2 {
		t.Errorf("Expected AudioOutputs to be 2, got %d", caps.AudioOutputs)
	}

	if caps.RelayOutputs != 4 {
		t.Errorf("Expected RelayOutputs to be 4, got %d", caps.RelayOutputs)
	}

	if caps.SerialPorts != 2 {
		t.Errorf("Expected SerialPorts to be 2, got %d", caps.SerialPorts)
	}

	if caps.DigitalInputs != 8 {
		t.Errorf("Expected DigitalInputs to be 8, got %d", caps.DigitalInputs)
	}

	if !caps.DigitalInputOptions {
		t.Error("Expected DigitalInputOptions to be true")
	}

	if !caps.SerialPortConfiguration {
		t.Error("Expected SerialPortConfiguration to be true")
	}
}

func TestGetDigitalInputs(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	inputs, err := client.GetDigitalInputs(ctx)
	if err != nil {
		t.Fatalf("GetDigitalInputs failed: %v", err)
	}

	if len(inputs) != 2 {
		t.Fatalf("Expected 2 digital inputs, got %d", len(inputs))
	}

	if inputs[0].Token != "input_001" {
		t.Errorf("Expected first input token 'input_001', got '%s'", inputs[0].Token)
	}

	if inputs[0].IdleState != DigitalIdleOpen {
		t.Errorf("Expected first input idle state 'open', got '%s'", inputs[0].IdleState)
	}

	if inputs[1].IdleState != DigitalIdleClosed {
		t.Errorf("Expected second input idle state 'closed', got '%s'", inputs[1].IdleState)
	}
}

func TestGetDigitalInputConfigurationOptions(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	options, err := client.GetDigitalInputConfigurationOptions(ctx, "input_001")
	if err != nil {
		t.Fatalf("GetDigitalInputConfigurationOptions failed: %v", err)
	}

	if len(options.IdleStateOptions) != 2 {
		t.Errorf("Expected 2 idle state options, got %d", len(options.IdleStateOptions))
	}
}

func TestGetDigitalInputConfigurationOptionsInvalidToken(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	_, err = client.GetDigitalInputConfigurationOptions(ctx, "")
	if !errors.Is(err, ErrInvalidDigitalInputToken) {
		t.Errorf("Expected ErrInvalidDigitalInputToken, got %v", err)
	}
}

func TestSetDigitalInputConfigurations(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	inputs := []*DigitalInput{
		{Token: "input_001", IdleState: DigitalIdleOpen},
		{Token: "input_002", IdleState: DigitalIdleClosed},
	}

	err = client.SetDigitalInputConfigurations(ctx, inputs)
	if err != nil {
		t.Fatalf("SetDigitalInputConfigurations failed: %v", err)
	}
}

func TestSetDigitalInputConfigurationsValidation(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test empty inputs.
	err = client.SetDigitalInputConfigurations(ctx, []*DigitalInput{})
	if !errors.Is(err, ErrDigitalInputConfigNil) {
		t.Errorf("Expected ErrDigitalInputConfigNil, got %v", err)
	}

	// Test input with empty token.
	inputs := []*DigitalInput{{Token: "", IdleState: DigitalIdleOpen}}
	err = client.SetDigitalInputConfigurations(ctx, inputs)
	if !errors.Is(err, ErrInvalidDigitalInputToken) {
		t.Errorf("Expected ErrInvalidDigitalInputToken, got %v", err)
	}
}

func TestGetVideoOutputs(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	outputs, err := client.GetVideoOutputs(ctx)
	if err != nil {
		t.Fatalf("GetVideoOutputs failed: %v", err)
	}

	if len(outputs) != 1 {
		t.Fatalf("Expected 1 video output, got %d", len(outputs))
	}

	if outputs[0].Token != "video_out_001" {
		t.Errorf("Expected video output token 'video_out_001', got '%s'", outputs[0].Token)
	}

	if outputs[0].Resolution == nil {
		t.Fatal("Expected Resolution to be set")
	}

	if outputs[0].Resolution.Width != 1920 {
		t.Errorf("Expected resolution width 1920, got %d", outputs[0].Resolution.Width)
	}

	if outputs[0].Resolution.Height != 1080 {
		t.Errorf("Expected resolution height 1080, got %d", outputs[0].Resolution.Height)
	}

	if outputs[0].RefreshRate != 60.0 {
		t.Errorf("Expected refresh rate 60.0, got %f", outputs[0].RefreshRate)
	}

	if outputs[0].AspectRatio != "16:9" {
		t.Errorf("Expected aspect ratio '16:9', got '%s'", outputs[0].AspectRatio)
	}
}

func TestGetSerialPorts(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	ports, err := client.GetSerialPorts(ctx)
	if err != nil {
		t.Fatalf("GetSerialPorts failed: %v", err)
	}

	if len(ports) != 2 {
		t.Fatalf("Expected 2 serial ports, got %d", len(ports))
	}

	if ports[0].Token != "serial_001" {
		t.Errorf("Expected first serial port token 'serial_001', got '%s'", ports[0].Token)
	}

	if ports[0].Type != SerialPortTypeRS232 {
		t.Errorf("Expected first serial port type RS232, got '%s'", ports[0].Type)
	}

	if ports[1].Type != SerialPortTypeRS485 {
		t.Errorf("Expected second serial port type RS485, got '%s'", ports[1].Type)
	}
}

func TestGetSerialPortConfiguration(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	config, err := client.GetSerialPortConfiguration(ctx, "serial_001")
	if err != nil {
		t.Fatalf("GetSerialPortConfiguration failed: %v", err)
	}

	if config.Token != "serial_001" {
		t.Errorf("Expected token 'serial_001', got '%s'", config.Token)
	}

	if config.Type != SerialPortTypeRS232 {
		t.Errorf("Expected type RS232, got '%s'", config.Type)
	}

	if config.BaudRate != 9600 {
		t.Errorf("Expected baud rate 9600, got %d", config.BaudRate)
	}

	if config.ParityBit != ParityNone {
		t.Errorf("Expected parity None, got '%s'", config.ParityBit)
	}

	if config.CharacterLength != 8 {
		t.Errorf("Expected character length 8, got %d", config.CharacterLength)
	}

	if config.StopBit != 1 {
		t.Errorf("Expected stop bit 1, got %f", config.StopBit)
	}
}

func TestGetSerialPortConfigurationInvalidToken(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	_, err = client.GetSerialPortConfiguration(ctx, "")
	if !errors.Is(err, ErrInvalidSerialPortToken) {
		t.Errorf("Expected ErrInvalidSerialPortToken, got %v", err)
	}
}

func TestGetSerialPortConfigurationOptions(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	options, err := client.GetSerialPortConfigurationOptions(ctx, "serial_001")
	if err != nil {
		t.Fatalf("GetSerialPortConfigurationOptions failed: %v", err)
	}

	if len(options.BaudRateList) != 3 {
		t.Errorf("Expected 3 baud rate options, got %d", len(options.BaudRateList))
	}

	if len(options.ParityBitList) != 3 {
		t.Errorf("Expected 3 parity bit options, got %d", len(options.ParityBitList))
	}

	if len(options.CharacterLengthList) != 2 {
		t.Errorf("Expected 2 character length options, got %d", len(options.CharacterLengthList))
	}

	if len(options.StopBitList) != 2 {
		t.Errorf("Expected 2 stop bit options, got %d", len(options.StopBitList))
	}
}

func TestGetSerialPortConfigurationOptionsInvalidToken(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	_, err = client.GetSerialPortConfigurationOptions(ctx, "")
	if !errors.Is(err, ErrInvalidSerialPortToken) {
		t.Errorf("Expected ErrInvalidSerialPortToken, got %v", err)
	}
}

func TestSetSerialPortConfiguration(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	config := &SerialPortConfiguration{
		Token:           "serial_001",
		Type:            SerialPortTypeRS232,
		BaudRate:        19200,
		ParityBit:       ParityNone,
		CharacterLength: 8,
		StopBit:         1,
	}

	err = client.SetSerialPortConfiguration(ctx, config)
	if err != nil {
		t.Fatalf("SetSerialPortConfiguration failed: %v", err)
	}
}

func TestSetSerialPortConfigurationValidation(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test nil config.
	err = client.SetSerialPortConfiguration(ctx, nil)
	if !errors.Is(err, ErrSerialPortConfigNil) {
		t.Errorf("Expected ErrSerialPortConfigNil, got %v", err)
	}

	// Test empty token.
	config := &SerialPortConfiguration{Token: ""}
	err = client.SetSerialPortConfiguration(ctx, config)
	if !errors.Is(err, ErrInvalidSerialPortToken) {
		t.Errorf("Expected ErrInvalidSerialPortToken, got %v", err)
	}
}

func TestSendReceiveSerialCommand(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	response, err := client.SendReceiveSerialCommand(ctx, "serial_001", []byte("HELLO"), 5, 10)
	if err != nil {
		t.Fatalf("SendReceiveSerialCommand failed: %v", err)
	}

	if string(response) != "OK" {
		t.Errorf("Expected response 'OK', got '%s'", string(response))
	}
}

func TestSendReceiveSerialCommandValidation(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test empty token.
	_, err = client.SendReceiveSerialCommand(ctx, "", []byte("HELLO"), 5, 10)
	if !errors.Is(err, ErrInvalidSerialPortToken) {
		t.Errorf("Expected ErrInvalidSerialPortToken, got %v", err)
	}

	// Test empty data.
	_, err = client.SendReceiveSerialCommand(ctx, "serial_001", []byte{}, 5, 10)
	if !errors.Is(err, ErrInvalidSerialData) {
		t.Errorf("Expected ErrInvalidSerialData, got %v", err)
	}
}

func TestDigitalIdleStateConstants(t *testing.T) {
	if DigitalIdleOpen != "open" {
		t.Errorf("DigitalIdleOpen should be 'open'")
	}

	if DigitalIdleClosed != "closed" {
		t.Errorf("DigitalIdleClosed should be 'closed'")
	}
}

func TestSerialPortTypeConstants(t *testing.T) {
	if SerialPortTypeRS232 != "RS232" {
		t.Errorf("SerialPortTypeRS232 should be 'RS232'")
	}

	if SerialPortTypeRS422 != "RS422" {
		t.Errorf("SerialPortTypeRS422 should be 'RS422'")
	}

	if SerialPortTypeRS485 != "RS485" {
		t.Errorf("SerialPortTypeRS485 should be 'RS485'")
	}

	if SerialPortTypeGeneric != "Generic" {
		t.Errorf("SerialPortTypeGeneric should be 'Generic'")
	}
}

func TestParityBitConstants(t *testing.T) {
	if ParityNone != "None" {
		t.Errorf("ParityNone should be 'None'")
	}

	if ParityOdd != "Odd" {
		t.Errorf("ParityOdd should be 'Odd'")
	}

	if ParityEven != "Even" {
		t.Errorf("ParityEven should be 'Even'")
	}

	if ParityMark != "Mark" {
		t.Errorf("ParityMark should be 'Mark'")
	}

	if ParitySpace != "Space" {
		t.Errorf("ParitySpace should be 'Space'")
	}
}

func TestGetVideoOutputConfiguration(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	config, err := client.GetVideoOutputConfiguration(ctx, "video_out_001")
	if err != nil {
		t.Fatalf("GetVideoOutputConfiguration failed: %v", err)
	}

	if config.Token != "config_001" {
		t.Errorf("Expected token 'config_001', got '%s'", config.Token)
	}

	if config.Name != "Main Output" {
		t.Errorf("Expected name 'Main Output', got '%s'", config.Name)
	}

	if config.UseCount != 2 {
		t.Errorf("Expected use count 2, got %d", config.UseCount)
	}

	if config.OutputToken != "video_out_001" {
		t.Errorf("Expected output token 'video_out_001', got '%s'", config.OutputToken)
	}
}

func TestGetVideoOutputConfigurationInvalidToken(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	_, err = client.GetVideoOutputConfiguration(ctx, "")
	if !errors.Is(err, ErrInvalidVideoOutputToken) {
		t.Errorf("Expected ErrInvalidVideoOutputToken, got %v", err)
	}
}

func TestGetVideoOutputConfigurationOptions(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	options, err := client.GetVideoOutputConfigurationOptions(ctx, "video_out_001")
	if err != nil {
		t.Fatalf("GetVideoOutputConfigurationOptions failed: %v", err)
	}

	if options.Name.Min != 1 {
		t.Errorf("Expected Name.Min to be 1, got %d", options.Name.Min)
	}

	if options.Name.Max != 64 {
		t.Errorf("Expected Name.Max to be 64, got %d", options.Name.Max)
	}

	if len(options.OutputTokensAvailable) != 2 {
		t.Errorf("Expected 2 output tokens available, got %d", len(options.OutputTokensAvailable))
	}
}

func TestGetVideoOutputConfigurationOptionsInvalidToken(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	_, err = client.GetVideoOutputConfigurationOptions(ctx, "")
	if !errors.Is(err, ErrInvalidVideoOutputToken) {
		t.Errorf("Expected ErrInvalidVideoOutputToken, got %v", err)
	}
}

func TestSetVideoOutputConfiguration(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	config := &VideoOutputConfiguration{
		Token:            "config_001",
		Name:             "Main Output",
		UseCount:         2,
		OutputToken:      "video_out_001",
		ForcePersistence: true,
	}

	err = client.SetVideoOutputConfiguration(ctx, config)
	if err != nil {
		t.Fatalf("SetVideoOutputConfiguration failed: %v", err)
	}
}

func TestSetVideoOutputConfigurationValidation(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test nil config.
	err = client.SetVideoOutputConfiguration(ctx, nil)
	if !errors.Is(err, ErrVideoOutputConfigNil) {
		t.Errorf("Expected ErrVideoOutputConfigNil, got %v", err)
	}

	// Test empty token.
	config := &VideoOutputConfiguration{Token: ""}
	err = client.SetVideoOutputConfiguration(ctx, config)
	if !errors.Is(err, ErrInvalidVideoOutputToken) {
		t.Errorf("Expected ErrInvalidVideoOutputToken, got %v", err)
	}
}

func TestGetRelayOutputOptions(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	options, err := client.GetRelayOutputOptions(ctx, "relay_001")
	if err != nil {
		t.Fatalf("GetRelayOutputOptions failed: %v", err)
	}

	if options.Token != "relay_001" {
		t.Errorf("Expected token 'relay_001', got '%s'", options.Token)
	}

	if len(options.Mode) != 2 {
		t.Errorf("Expected 2 modes, got %d", len(options.Mode))
	}

	if options.Mode[0] != RelayModeMonostable {
		t.Errorf("Expected first mode to be Monostable, got '%s'", options.Mode[0])
	}

	if options.Mode[1] != RelayModeBistable {
		t.Errorf("Expected second mode to be Bistable, got '%s'", options.Mode[1])
	}

	if len(options.DelayTimes) != 3 {
		t.Errorf("Expected 3 delay times, got %d", len(options.DelayTimes))
	}

	if !options.Discrete {
		t.Error("Expected Discrete to be true")
	}
}

func TestGetRelayOutputOptionsInvalidToken(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	_, err = client.GetRelayOutputOptions(ctx, "")
	if !errors.Is(err, ErrInvalidRelayOutputToken) {
		t.Errorf("Expected ErrInvalidRelayOutputToken, got %v", err)
	}
}

func TestGetDeviceIOAudioSources(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
		wantLen int
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetAudioSourcesResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:Token>AudioSrc1</tmd:Token>
      <tmd:Token>AudioSrc2</tmd:Token>
    </tmd:GetAudioSourcesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "SOAP fault",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <SOAP-ENV:Fault>
      <SOAP-ENV:Code><SOAP-ENV:Value>SOAP-ENV:Receiver</SOAP-ENV:Value></SOAP-ENV:Code>
      <SOAP-ENV:Reason><SOAP-ENV:Text>Internal error</SOAP-ENV:Text></SOAP-ENV:Reason>
    </SOAP-ENV:Fault>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL, WithCredentials("admin", "password"))
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			tokens, err := client.GetDeviceIOAudioSources(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDeviceIOAudioSources() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && len(tokens) != tt.wantLen {
				t.Errorf("Expected %d tokens, got %d", tt.wantLen, len(tokens))
			}
		})
	}
}

func TestGetDeviceIOAudioOutputs(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tokens, err := client.GetDeviceIOAudioOutputs(context.Background())
	if err != nil {
		t.Fatalf("GetDeviceIOAudioOutputs failed: %v", err)
	}

	if len(tokens) != 1 {
		t.Fatalf("Expected 1 token, got %d", len(tokens))
	}

	if tokens[0] != "AudioOut1" {
		t.Errorf("Expected token 'AudioOut1', got '%s'", tokens[0])
	}
}

func TestGetDeviceIOVideoSources(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tokens, err := client.GetDeviceIOVideoSources(context.Background())
	if err != nil {
		t.Fatalf("GetDeviceIOVideoSources failed: %v", err)
	}

	if len(tokens) != 3 {
		t.Fatalf("Expected 3 tokens, got %d", len(tokens))
	}
}

func TestGetDeviceIOAudioSourceConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetAudioSourceConfigurationResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:AudioSourceConfiguration token="cfg_audio_001">
        <tmd:Name>Microphone</tmd:Name>
        <tmd:UseCount>1</tmd:UseCount>
        <tmd:SourceToken>AudioSrc1</tmd:SourceToken>
      </tmd:AudioSourceConfiguration>
    </tmd:GetAudioSourceConfigurationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
			},
			wantErr: false,
		},
		{
			name: "SOAP fault",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <SOAP-ENV:Fault>
      <SOAP-ENV:Code><SOAP-ENV:Value>SOAP-ENV:Receiver</SOAP-ENV:Value></SOAP-ENV:Code>
      <SOAP-ENV:Reason><SOAP-ENV:Text>Not found</SOAP-ENV:Text></SOAP-ENV:Reason>
    </SOAP-ENV:Fault>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL, WithCredentials("admin", "password"))
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			cfg, err := client.GetDeviceIOAudioSourceConfiguration(context.Background(), "AudioSrc1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDeviceIOAudioSourceConfiguration() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if cfg == nil {
					t.Fatal("Expected config, got nil")
				}

				if cfg.Token != "cfg_audio_001" {
					t.Errorf("Expected token 'cfg_audio_001', got '%s'", cfg.Token)
				}

				if cfg.Name != "Microphone" {
					t.Errorf("Expected name 'Microphone', got '%s'", cfg.Name)
				}

				if cfg.SourceToken != "AudioSrc1" {
					t.Errorf("Expected source token 'AudioSrc1', got '%s'", cfg.SourceToken)
				}
			}
		})
	}
}

func TestSetDeviceIOAudioSourceConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:SetAudioSourceConfigurationResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
			},
			wantErr: false,
		},
		{
			name: "SOAP fault",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <SOAP-ENV:Fault>
      <SOAP-ENV:Code><SOAP-ENV:Value>SOAP-ENV:Receiver</SOAP-ENV:Value></SOAP-ENV:Code>
      <SOAP-ENV:Reason><SOAP-ENV:Text>Permission denied</SOAP-ENV:Text></SOAP-ENV:Reason>
    </SOAP-ENV:Fault>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL, WithCredentials("admin", "password"))
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			cfg := &AudioSourceConfiguration{
				Token:       "cfg_audio_001",
				Name:        "Microphone",
				UseCount:    1,
				SourceToken: "AudioSrc1",
			}

			err = client.SetDeviceIOAudioSourceConfiguration(context.Background(), cfg, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetDeviceIOAudioSourceConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetDeviceIOAudioOutputConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetAudioOutputConfigurationResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:AudioOutputConfiguration token="cfg_out_001">
        <tmd:Name>Speaker</tmd:Name>
        <tmd:UseCount>2</tmd:UseCount>
        <tmd:OutputToken>AudioOut1</tmd:OutputToken>
      </tmd:AudioOutputConfiguration>
    </tmd:GetAudioOutputConfigurationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
			},
			wantErr: false,
		},
		{
			name: "SOAP fault",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <SOAP-ENV:Fault>
      <SOAP-ENV:Code><SOAP-ENV:Value>SOAP-ENV:Receiver</SOAP-ENV:Value></SOAP-ENV:Code>
      <SOAP-ENV:Reason><SOAP-ENV:Text>Not found</SOAP-ENV:Text></SOAP-ENV:Reason>
    </SOAP-ENV:Fault>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL, WithCredentials("admin", "password"))
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			cfg, err := client.GetDeviceIOAudioOutputConfiguration(context.Background(), "AudioOut1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDeviceIOAudioOutputConfiguration() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if cfg == nil {
					t.Fatal("Expected config, got nil")
				}

				if cfg.Token != "cfg_out_001" {
					t.Errorf("Expected token 'cfg_out_001', got '%s'", cfg.Token)
				}

				if cfg.Name != "Speaker" {
					t.Errorf("Expected name 'Speaker', got '%s'", cfg.Name)
				}

				if cfg.OutputToken != "AudioOut1" {
					t.Errorf("Expected output token 'AudioOut1', got '%s'", cfg.OutputToken)
				}
			}
		})
	}
}

func TestSetDeviceIOAudioOutputConfiguration(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	cfg := &AudioOutputConfiguration{
		Token:       "cfg_out_001",
		Name:        "Speaker",
		UseCount:    2,
		OutputToken: "AudioOut1",
	}

	err = client.SetDeviceIOAudioOutputConfiguration(context.Background(), cfg, false)
	if err != nil {
		t.Fatalf("SetDeviceIOAudioOutputConfiguration failed: %v", err)
	}
}

func TestGetDeviceIOVideoSourceConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetVideoSourceConfigurationResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:VideoSourceConfiguration token="cfg_vid_001">
        <tmd:Name>Main Camera</tmd:Name>
        <tmd:UseCount>3</tmd:UseCount>
        <tmd:SourceToken>VideoSrc1</tmd:SourceToken>
        <tmd:Bounds x="0" y="0" width="1920" height="1080"/>
      </tmd:VideoSourceConfiguration>
    </tmd:GetVideoSourceConfigurationResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
			},
			wantErr: false,
		},
		{
			name: "SOAP fault",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <SOAP-ENV:Fault>
      <SOAP-ENV:Code><SOAP-ENV:Value>SOAP-ENV:Receiver</SOAP-ENV:Value></SOAP-ENV:Code>
      <SOAP-ENV:Reason><SOAP-ENV:Text>Not found</SOAP-ENV:Text></SOAP-ENV:Reason>
    </SOAP-ENV:Fault>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL, WithCredentials("admin", "password"))
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			cfg, err := client.GetDeviceIOVideoSourceConfiguration(context.Background(), "VideoSrc1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDeviceIOVideoSourceConfiguration() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if cfg == nil {
					t.Fatal("Expected config, got nil")
				}

				if cfg.Token != "cfg_vid_001" {
					t.Errorf("Expected token 'cfg_vid_001', got '%s'", cfg.Token)
				}

				if cfg.SourceToken != "VideoSrc1" {
					t.Errorf("Expected source token 'VideoSrc1', got '%s'", cfg.SourceToken)
				}

				if cfg.Bounds == nil {
					t.Fatal("Expected Bounds to be set")
				}

				if cfg.Bounds.Width != 1920 || cfg.Bounds.Height != 1080 {
					t.Errorf("Expected bounds 1920x1080, got %dx%d", cfg.Bounds.Width, cfg.Bounds.Height)
				}
			}
		})
	}
}

func TestSetDeviceIOVideoSourceConfiguration(t *testing.T) {
	server := newMockDeviceIOServer()
	defer server.Close()

	client, err := NewClient(server.URL, WithCredentials("admin", "password"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	cfg := &VideoSourceConfiguration{
		Token:       "cfg_vid_001",
		Name:        "Main Camera",
		UseCount:    3,
		SourceToken: "VideoSrc1",
		Bounds:      &IntRectangle{X: 0, Y: 0, Width: 1920, Height: 1080},
	}

	err = client.SetDeviceIOVideoSourceConfiguration(context.Background(), cfg, true)
	if err != nil {
		t.Fatalf("SetDeviceIOVideoSourceConfiguration failed: %v", err)
	}
}

func TestGetDeviceIORelayOutputs(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
		wantLen int
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:GetRelayOutputsResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl">
      <tmd:RelayOutputs token="relay_001">
        <tmd:Properties>
          <tmd:Mode>Monostable</tmd:Mode>
          <tmd:DelayTime>PT5S</tmd:DelayTime>
          <tmd:IdleState>open</tmd:IdleState>
        </tmd:Properties>
      </tmd:RelayOutputs>
      <tmd:RelayOutputs token="relay_002">
        <tmd:Properties>
          <tmd:Mode>Bistable</tmd:Mode>
          <tmd:IdleState>closed</tmd:IdleState>
        </tmd:Properties>
      </tmd:RelayOutputs>
    </tmd:GetRelayOutputsResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "SOAP fault",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <SOAP-ENV:Fault>
      <SOAP-ENV:Code><SOAP-ENV:Value>SOAP-ENV:Receiver</SOAP-ENV:Value></SOAP-ENV:Code>
      <SOAP-ENV:Reason><SOAP-ENV:Text>Service unavailable</SOAP-ENV:Text></SOAP-ENV:Reason>
    </SOAP-ENV:Fault>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL, WithCredentials("admin", "password"))
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			outputs, err := client.GetDeviceIORelayOutputs(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDeviceIORelayOutputs() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				if len(outputs) != tt.wantLen {
					t.Fatalf("Expected %d relay outputs, got %d", tt.wantLen, len(outputs))
				}

				if outputs[0].Token != "relay_001" {
					t.Errorf("Expected token 'relay_001', got '%s'", outputs[0].Token)
				}

				if outputs[0].Properties.Mode != RelayModeMonostable {
					t.Errorf("Expected mode Monostable, got '%s'", outputs[0].Properties.Mode)
				}

				if outputs[0].Properties.IdleState != RelayIdleStateOpen {
					t.Errorf("Expected idle state open, got '%s'", outputs[0].Properties.IdleState)
				}

				if outputs[1].Properties.Mode != RelayModeBistable {
					t.Errorf("Expected mode Bistable, got '%s'", outputs[1].Properties.Mode)
				}
			}
		})
	}
}

func TestSetDeviceIORelayOutputState(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:SetRelayOutputStateResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
			},
			wantErr: false,
		},
		{
			name: "SOAP fault",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <SOAP-ENV:Fault>
      <SOAP-ENV:Code><SOAP-ENV:Value>SOAP-ENV:Receiver</SOAP-ENV:Value></SOAP-ENV:Code>
      <SOAP-ENV:Reason><SOAP-ENV:Text>Invalid token</SOAP-ENV:Text></SOAP-ENV:Reason>
    </SOAP-ENV:Fault>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL, WithCredentials("admin", "password"))
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			err = client.SetDeviceIORelayOutputState(context.Background(), "relay_001", string(RelayLogicalStateActive))
			if (err != nil) != tt.wantErr {
				t.Errorf("SetDeviceIORelayOutputState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetDeviceIORelayOutputSettings(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <tmd:SetRelayOutputSettingsResponse xmlns:tmd="http://www.onvif.org/ver10/deviceIO/wsdl"/>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
			},
			wantErr: false,
		},
		{
			name: "SOAP fault",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(testDeviceIOXMLHeader + `
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://www.w3.org/2003/05/soap-envelope">
  <SOAP-ENV:Body>
    <SOAP-ENV:Fault>
      <SOAP-ENV:Code><SOAP-ENV:Value>SOAP-ENV:Receiver</SOAP-ENV:Value></SOAP-ENV:Code>
      <SOAP-ENV:Reason><SOAP-ENV:Text>Invalid token</SOAP-ENV:Text></SOAP-ENV:Reason>
    </SOAP-ENV:Fault>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client, err := NewClient(server.URL, WithCredentials("admin", "password"))
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			ro := &RelayOutput{
				Token: "relay_001",
				Properties: RelayOutputSettings{
					Mode:      RelayModeMonostable,
					IdleState: RelayIdleStateOpen,
				},
			}

			err = client.SetDeviceIORelayOutputSettings(context.Background(), ro)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetDeviceIORelayOutputSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
