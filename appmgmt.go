package onvif

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/0x524a/onvif-go/internal/soap"
)

// App Management service namespace.
const appMgmtNamespace = "http://www.onvif.org/ver10/appmgmt/wsdl"

// getAppMgmtEndpoint returns the app management service endpoint, falling back to the device endpoint.
func (c *Client) getAppMgmtEndpoint() string {
	if c.appmgmtEndpoint != "" {
		return c.appmgmtEndpoint
	}

	return c.endpoint
}

// GetAppMgmtServiceCapabilities returns the capabilities of the app management service.
func (c *Client) GetAppMgmtServiceCapabilities(ctx context.Context) (*AppMgmtServiceCapabilities, error) {
	endpoint := c.getAppMgmtEndpoint()

	type GetServiceCapabilities struct {
		XMLName xml.Name `xml:"tap:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:tap,attr"`
	}

	type CapabilitiesEntry struct {
		FormatsSupported string  `xml:"FormatsSupported,attr"`
		Licensing        *bool   `xml:"Licensing,attr"`
		UploadPath       string  `xml:"UploadPath,attr"`
		EventTopicPrefix string  `xml:"EventTopicPrefix,attr"`
	}

	type GetServiceCapabilitiesResponse struct {
		XMLName      xml.Name          `xml:"GetServiceCapabilitiesResponse"`
		Capabilities CapabilitiesEntry `xml:"Capabilities"`
	}

	req := GetServiceCapabilities{
		Xmlns: appMgmtNamespace,
	}

	var resp GetServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAppMgmtServiceCapabilities failed: %w", err)
	}

	return &AppMgmtServiceCapabilities{
		FormatsSupported: resp.Capabilities.FormatsSupported,
		Licensing:        resp.Capabilities.Licensing,
		UploadPath:       resp.Capabilities.UploadPath,
		EventTopicPrefix: resp.Capabilities.EventTopicPrefix,
	}, nil
}

// GetInstalledApps lists installed apps on the device.
func (c *Client) GetInstalledApps(ctx context.Context) ([]*InstalledApp, error) {
	endpoint := c.getAppMgmtEndpoint()

	type GetInstalledApps struct {
		XMLName xml.Name `xml:"tap:GetInstalledApps"`
		Xmlns   string   `xml:"xmlns:tap,attr"`
	}

	type AppEntry struct {
		Name  string `xml:"Name"`
		AppID string `xml:"AppID"`
	}

	type GetInstalledAppsResponse struct {
		XMLName xml.Name   `xml:"GetInstalledAppsResponse"`
		Apps    []AppEntry `xml:"App"`
	}

	req := GetInstalledApps{
		Xmlns: appMgmtNamespace,
	}

	var resp GetInstalledAppsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetInstalledApps failed: %w", err)
	}

	apps := make([]*InstalledApp, 0, len(resp.Apps))
	for i := range resp.Apps {
		apps = append(apps, &InstalledApp{
			Name:  resp.Apps[i].Name,
			AppID: resp.Apps[i].AppID,
		})
	}

	return apps, nil
}

// GetAppsInfo retrieves detailed information about installed applications.
// If appID is empty, information for all installed applications is returned.
func (c *Client) GetAppsInfo(ctx context.Context, appID string) ([]*AppInfo, error) {
	endpoint := c.getAppMgmtEndpoint()

	type GetAppsInfo struct {
		XMLName xml.Name `xml:"tap:GetAppsInfo"`
		Xmlns   string   `xml:"xmlns:tap,attr"`
		AppID   *string  `xml:"tap:AppID,omitempty"`
	}

	type LicenseEntry struct {
		Name       string  `xml:"Name"`
		ValidFrom  *string `xml:"ValidFrom"`
		ValidUntil *string `xml:"ValidUntil"`
	}

	type AppInfoEntry struct {
		AppID                string         `xml:"AppID"`
		Name                 string         `xml:"Name"`
		Version              string         `xml:"Version"`
		Licenses             []LicenseEntry `xml:"Licenses"`
		Privileges           []string       `xml:"Privileges"`
		InstallationDate     string         `xml:"InstallationDate"`
		LastUpdate           string         `xml:"LastUpdate"`
		State                string         `xml:"State"`
		Status               string         `xml:"Status"`
		Autostart            bool           `xml:"Autostart"`
		Website              string         `xml:"Website"`
		OpenSource           string         `xml:"OpenSource"`
		Configuration        string         `xml:"Configuration"`
		InterfaceDescription []string       `xml:"InterfaceDescription"`
	}

	type GetAppsInfoResponse struct {
		XMLName xml.Name       `xml:"GetAppsInfoResponse"`
		Info    []AppInfoEntry `xml:"Info"`
	}

	req := GetAppsInfo{
		Xmlns: appMgmtNamespace,
	}

	if appID != "" {
		req.AppID = &appID
	}

	var resp GetAppsInfoResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAppsInfo failed: %w", err)
	}

	result := make([]*AppInfo, 0, len(resp.Info))

	for i := range resp.Info {
		entry := &resp.Info[i]
		info := &AppInfo{
			AppID:                entry.AppID,
			Name:                 entry.Name,
			Version:              entry.Version,
			Privileges:           entry.Privileges,
			State:                AppState(entry.State),
			Status:               entry.Status,
			Autostart:            entry.Autostart,
			Website:              entry.Website,
			OpenSource:           entry.OpenSource,
			Configuration:        entry.Configuration,
			InterfaceDescription: entry.InterfaceDescription,
		}

		if t, err := time.Parse(time.RFC3339, entry.InstallationDate); err == nil {
			info.InstallationDate = t
		}

		if t, err := time.Parse(time.RFC3339, entry.LastUpdate); err == nil {
			info.LastUpdate = t
		}

		for j := range entry.Licenses {
			lic := &LicenseInfo{Name: entry.Licenses[j].Name}

			if entry.Licenses[j].ValidFrom != nil {
				if t, err := time.Parse(time.RFC3339, *entry.Licenses[j].ValidFrom); err == nil {
					lic.ValidFrom = &t
				}
			}

			if entry.Licenses[j].ValidUntil != nil {
				if t, err := time.Parse(time.RFC3339, *entry.Licenses[j].ValidUntil); err == nil {
					lic.ValidUntil = &t
				}
			}

			info.Licenses = append(info.Licenses, lic)
		}

		result = append(result, info)
	}

	return result, nil
}

// ActivateApp starts an application identified by appID.
func (c *Client) ActivateApp(ctx context.Context, appID string) error {
	endpoint := c.getAppMgmtEndpoint()

	type Activate struct {
		XMLName xml.Name `xml:"tap:Activate"`
		Xmlns   string   `xml:"xmlns:tap,attr"`
		AppID   string   `xml:"tap:AppID"`
	}

	type ActivateResponse struct {
		XMLName xml.Name `xml:"ActivateResponse"`
	}

	req := Activate{
		Xmlns: appMgmtNamespace,
		AppID: appID,
	}

	var resp ActivateResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("ActivateApp failed: %w", err)
	}

	return nil
}

// DeactivateApp stops an application identified by appID.
func (c *Client) DeactivateApp(ctx context.Context, appID string) error {
	endpoint := c.getAppMgmtEndpoint()

	type Deactivate struct {
		XMLName xml.Name `xml:"tap:Deactivate"`
		Xmlns   string   `xml:"xmlns:tap,attr"`
		AppID   string   `xml:"tap:AppID"`
	}

	type DeactivateResponse struct {
		XMLName xml.Name `xml:"DeactivateResponse"`
	}

	req := Deactivate{
		Xmlns: appMgmtNamespace,
		AppID: appID,
	}

	var resp DeactivateResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("DeactivateApp failed: %w", err)
	}

	return nil
}

// UninstallApp removes an application from the device.
// This method returns immediately; completion or failure is delivered via an UninstallCompletion event.
func (c *Client) UninstallApp(ctx context.Context, appID string) error {
	endpoint := c.getAppMgmtEndpoint()

	type Uninstall struct {
		XMLName xml.Name `xml:"tap:Uninstall"`
		Xmlns   string   `xml:"xmlns:tap,attr"`
		AppID   string   `xml:"tap:AppID"`
	}

	type UninstallResponse struct {
		XMLName xml.Name `xml:"UninstallResponse"`
	}

	req := Uninstall{
		Xmlns: appMgmtNamespace,
		AppID: appID,
	}

	var resp UninstallResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("UninstallApp failed: %w", err)
	}

	return nil
}

// InstallLicense installs a license on the device.
// If the device requires per-app licensing, appID must be provided.
func (c *Client) InstallLicense(ctx context.Context, appID string, license string) error {
	endpoint := c.getAppMgmtEndpoint()

	type InstallLicense struct {
		XMLName xml.Name `xml:"tap:InstallLicense"`
		Xmlns   string   `xml:"xmlns:tap,attr"`
		AppID   *string  `xml:"tap:AppID,omitempty"`
		License string   `xml:"tap:License"`
	}

	type InstallLicenseResponse struct {
		XMLName xml.Name `xml:"InstallLicenseResponse"`
	}

	req := InstallLicense{
		Xmlns:   appMgmtNamespace,
		License: license,
	}

	if appID != "" {
		req.AppID = &appID
	}

	var resp InstallLicenseResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return fmt.Errorf("InstallLicense failed: %w", err)
	}

	return nil
}

// GetAppDeviceID returns the unique device ID to which licenses are issued.
func (c *Client) GetAppDeviceID(ctx context.Context) (string, error) {
	endpoint := c.getAppMgmtEndpoint()

	type GetDeviceId struct {
		XMLName xml.Name `xml:"tap:GetDeviceId"`
		Xmlns   string   `xml:"xmlns:tap,attr"`
	}

	type GetDeviceIdResponse struct {
		XMLName  xml.Name `xml:"GetDeviceIdResponse"`
		DeviceId string   `xml:"DeviceId"`
	}

	req := GetDeviceId{
		Xmlns: appMgmtNamespace,
	}

	var resp GetDeviceIdResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return "", fmt.Errorf("GetAppDeviceID failed: %w", err)
	}

	return resp.DeviceId, nil
}
