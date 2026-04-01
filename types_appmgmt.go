package onvif

import "time"

// AppState represents the current state of an installed application.
type AppState string

const (
	// AppStateActive indicates the application is running.
	AppStateActive AppState = "Active"
	// AppStateInactive indicates the application is installed but not running.
	AppStateInactive AppState = "Inactive"
	// AppStateInstalling indicates the application is currently being installed.
	AppStateInstalling AppState = "Installing"
	// AppStateUninstalling indicates the application is currently being uninstalled.
	AppStateUninstalling AppState = "Uninstalling"
	// AppStateRemoved indicates the application has been removed.
	AppStateRemoved AppState = "Removed"
	// AppStateInstallationFailed indicates the application installation failed.
	AppStateInstallationFailed AppState = "InstallationFailed"
)

// AppMgmtServiceCapabilities contains the capabilities of the App Management service.
type AppMgmtServiceCapabilities struct {
	// FormatsSupported lists supported app container formats that can be uploaded.
	FormatsSupported string
	// Licensing signals support for licensing of applications.
	Licensing *bool
	// UploadPath is the path part of the URI to which applications can be uploaded via HTTP POST.
	UploadPath string
	// EventTopicPrefix is an optional event topic prefix used for app events.
	EventTopicPrefix string
}

// LicenseInfo contains information about a license associated with an application.
type LicenseInfo struct {
	// Name is the textual name of the license.
	Name string
	// ValidFrom is the optional start time of license validity.
	ValidFrom *time.Time
	// ValidUntil is the optional end time of license validity.
	ValidUntil *time.Time
}

// AppInfo contains detailed information about an installed application.
type AppInfo struct {
	// AppID is the unique app identifier of the application instance.
	AppID string
	// Name is the user readable application name.
	Name string
	// Version is the version of the installed application.
	Version string
	// Licenses contains licenses associated with the application.
	Licenses []*LicenseInfo
	// Privileges lists privileges granted to the application.
	Privileges []string
	// InstallationDate is the date and time when the application was installed.
	InstallationDate time.Time
	// LastUpdate is the time of the last update to this app.
	LastUpdate time.Time
	// State is the current state of the application.
	State AppState
	// Status is supplemental information about the current state.
	Status string
	// Autostart indicates if the application starts automatically after device boot.
	Autostart bool
	// Website is a link to supplementary information about the application or vendor.
	Website string
	// OpenSource is an optional link to a list of open source licenses used by the application.
	OpenSource string
	// Configuration is an optional URI for backup and restore of the application configuration.
	Configuration string
	// InterfaceDescription contains optional references to the interface definition of the application.
	InterfaceDescription []string
}

// InstalledApp provides basic identification of an installed application.
type InstalledApp struct {
	// Name is the user readable application name.
	Name string
	// AppID is the unique app identifier.
	AppID string
}
