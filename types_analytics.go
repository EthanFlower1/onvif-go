package onvif

// AnalyticsRule represents an analytics rule configuration.
type AnalyticsRule struct {
	Name       string
	Type       string
	Parameters []*SimpleItem
}

// AnalyticsModule represents an analytics module configuration.
type AnalyticsModule struct {
	Name       string
	Type       string
	Parameters []*SimpleItem
}

// SupportedRule describes a supported analytics rule type.
type SupportedRule struct {
	Name       string
	Parameters []*SimpleItem
}

// SupportedAnalyticsModule describes a supported analytics module type.
type SupportedAnalyticsModule struct {
	Name       string
	Parameters []*SimpleItem
}

// RuleOptions represents options for configuring analytics rules.
type RuleOptions struct {
	RuleType string
	Items    []*SimpleItem
}

// AnalyticsModuleOptions represents options for configuring analytics modules.
type AnalyticsModuleOptions struct {
	ModuleType string
	Items      []*SimpleItem
}

// AnalyticsEngineModuleConfiguration contains analytics engine module settings.
// Note: AnalyticsEngineConfiguration (with different fields) already exists in types.go.
type AnalyticsEngineModuleConfiguration struct {
	AnalyticsModule []*AnalyticsModule
}

// AnalyticsEngine represents an analytics processing engine.
type AnalyticsEngine struct {
	Token                        string
	Name                         string
	AnalyticsEngineConfiguration *AnalyticsEngineModuleConfiguration
}

// AnalyticsEngineControl represents control settings for an analytics engine.
type AnalyticsEngineControl struct {
	Token              string
	Name               string
	EngineToken        string
	EngineConfigToken  string
	InputToken         []string
	MulticastOutput    bool
	SubscriptionPolicy string
	Mode               string
}

// AnalyticsEngineInput represents an input to an analytics engine.
type AnalyticsEngineInput struct {
	Token                string
	Name                 string
	SourceIdentification *SourceIdentification
	VideoInput           *AnalyticsVideoInput
	MetadataInput        *MetadataInput
}

// SourceIdentification identifies a source for analytics.
type SourceIdentification struct {
	Name  string
	Token []string
}

// AnalyticsVideoInput represents video input settings for analytics.
// Named AnalyticsVideoInput to avoid collision with any VideoInput type.
type AnalyticsVideoInput struct {
	InputToken string
	FrameRate  *int
	Resolution *VideoResolution
}

// MetadataInput represents metadata input settings.
type MetadataInput struct {
	MetadataConfig []string
	Extensions     interface{}
}

// AnalyticsState represents the current state of analytics processing.
type AnalyticsState struct {
	State string
	Error string
}

// AnalyticsServiceCapabilities represents analytics service capabilities.
type AnalyticsServiceCapabilities struct {
	RuleSupport                        bool
	AnalyticsModuleSupport             bool
	CellBasedSceneDescriptionSupported bool
}

// AnalyticsDeviceServiceCapabilities represents analytics device service capabilities.
type AnalyticsDeviceServiceCapabilities struct {
	RuleSupport bool
}

// SupportedMetadata describes metadata formats supported by the analytics service.
type SupportedMetadata struct {
	AnalyticsModules []string
}
