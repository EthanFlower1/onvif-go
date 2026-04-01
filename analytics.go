package onvif

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/0x524a/onvif-go/internal/soap"
)

// Analytics service namespace.
const analyticsNamespace = "http://www.onvif.org/ver20/analytics/wsdl"

func (c *Client) getAnalyticsEndpoint() string {
	if c.analyticsEndpoint != "" {
		return c.analyticsEndpoint
	}

	return c.endpoint
}

// GetAnalyticsServiceCapabilities retrieves the capabilities of the analytics service.
func (c *Client) GetAnalyticsServiceCapabilities(ctx context.Context) (*AnalyticsServiceCapabilities, error) {
	endpoint := c.getAnalyticsEndpoint()

	type GetServiceCapabilities struct {
		XMLName xml.Name `xml:"tan:GetServiceCapabilities"`
		Xmlns   string   `xml:"xmlns:tan,attr"`
	}

	type GetServiceCapabilitiesResponse struct {
		XMLName      xml.Name `xml:"GetServiceCapabilitiesResponse"`
		Capabilities struct {
			RuleSupport                        bool `xml:"RuleSupport,attr"`
			AnalyticsModuleSupport             bool `xml:"AnalyticsModuleSupport,attr"`
			CellBasedSceneDescriptionSupported bool `xml:"CellBasedSceneDescriptionSupported,attr"`
		} `xml:"Capabilities"`
	}

	req := GetServiceCapabilities{
		Xmlns: analyticsNamespace,
	}

	var resp GetServiceCapabilitiesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAnalyticsServiceCapabilities failed: %w", err)
	}

	return &AnalyticsServiceCapabilities{
		RuleSupport:                        resp.Capabilities.RuleSupport,
		AnalyticsModuleSupport:             resp.Capabilities.AnalyticsModuleSupport,
		CellBasedSceneDescriptionSupported: resp.Capabilities.CellBasedSceneDescriptionSupported,
	}, nil
}

// GetSupportedRules retrieves the supported analytics rules for a configuration token.
func (c *Client) GetSupportedRules(ctx context.Context, configToken string) ([]*SupportedRule, error) {
	endpoint := c.getAnalyticsEndpoint()

	type GetSupportedRules struct {
		XMLName            xml.Name `xml:"tan:GetSupportedRules"`
		Xmlns              string   `xml:"xmlns:tan,attr"`
		ConfigurationToken string   `xml:"tan:ConfigurationToken"`
	}

	type SimpleItemDescription struct {
		Name string `xml:"Name,attr"`
		Type string `xml:"Type,attr"`
	}

	type RuleDescription struct {
		Name       string                  `xml:"Name,attr"`
		Parameters []SimpleItemDescription `xml:"Parameters>SimpleItemDescription"`
	}

	type GetSupportedRulesResponse struct {
		XMLName        xml.Name          `xml:"GetSupportedRulesResponse"`
		RuleDescription []RuleDescription `xml:"SupportedRules>RuleDescription"`
	}

	req := GetSupportedRules{
		Xmlns:              analyticsNamespace,
		ConfigurationToken: configToken,
	}

	var resp GetSupportedRulesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetSupportedRules failed: %w", err)
	}

	rules := make([]*SupportedRule, 0, len(resp.RuleDescription))
	for i := range resp.RuleDescription {
		rule := &SupportedRule{
			Name: resp.RuleDescription[i].Name,
		}
		for j := range resp.RuleDescription[i].Parameters {
			rule.Parameters = append(rule.Parameters, &SimpleItem{
				Name:  resp.RuleDescription[i].Parameters[j].Name,
				Value: resp.RuleDescription[i].Parameters[j].Type,
			})
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

// GetRules retrieves the analytics rules for a configuration token.
func (c *Client) GetRules(ctx context.Context, configToken string) ([]*AnalyticsRule, error) {
	endpoint := c.getAnalyticsEndpoint()

	type GetRules struct {
		XMLName            xml.Name `xml:"tan:GetRules"`
		Xmlns              string   `xml:"xmlns:tan,attr"`
		ConfigurationToken string   `xml:"tan:ConfigurationToken"`
	}

	type SimpleItemEntry struct {
		Name  string `xml:"Name,attr"`
		Value string `xml:"Value,attr"`
	}

	type RuleEntry struct {
		Name       string            `xml:"Name,attr"`
		Type       string            `xml:"Type,attr"`
		Parameters []SimpleItemEntry `xml:"Parameters>SimpleItem"`
	}

	type GetRulesResponse struct {
		XMLName xml.Name    `xml:"GetRulesResponse"`
		Rule    []RuleEntry `xml:"Rule"`
	}

	req := GetRules{
		Xmlns:              analyticsNamespace,
		ConfigurationToken: configToken,
	}

	var resp GetRulesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRules failed: %w", err)
	}

	rules := make([]*AnalyticsRule, 0, len(resp.Rule))
	for i := range resp.Rule {
		rule := &AnalyticsRule{
			Name: resp.Rule[i].Name,
			Type: resp.Rule[i].Type,
		}
		for j := range resp.Rule[i].Parameters {
			rule.Parameters = append(rule.Parameters, &SimpleItem{
				Name:  resp.Rule[i].Parameters[j].Name,
				Value: resp.Rule[i].Parameters[j].Value,
			})
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

// GetRuleOptions retrieves rule options for a configuration token and optional rule type filter.
func (c *Client) GetRuleOptions(ctx context.Context, configToken string, ruleType *string) ([]*RuleOptions, error) {
	endpoint := c.getAnalyticsEndpoint()

	type GetRuleOptions struct {
		XMLName            xml.Name `xml:"tan:GetRuleOptions"`
		Xmlns              string   `xml:"xmlns:tan,attr"`
		ConfigurationToken string   `xml:"tan:ConfigurationToken"`
		RuleType           *string  `xml:"tan:RuleType,omitempty"`
	}

	type SimpleItemEntry struct {
		Name  string `xml:"Name,attr"`
		Value string `xml:"Value,attr"`
	}

	type RuleOptionsEntry struct {
		RuleType string            `xml:"RuleType,attr"`
		Items    []SimpleItemEntry `xml:"SimpleItem"`
	}

	type GetRuleOptionsResponse struct {
		XMLName     xml.Name           `xml:"GetRuleOptionsResponse"`
		RuleOptions []RuleOptionsEntry `xml:"RuleOptions"`
	}

	req := GetRuleOptions{
		Xmlns:              analyticsNamespace,
		ConfigurationToken: configToken,
		RuleType:           ruleType,
	}

	var resp GetRuleOptionsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetRuleOptions failed: %w", err)
	}

	options := make([]*RuleOptions, 0, len(resp.RuleOptions))
	for i := range resp.RuleOptions {
		opt := &RuleOptions{
			RuleType: resp.RuleOptions[i].RuleType,
		}
		for j := range resp.RuleOptions[i].Items {
			opt.Items = append(opt.Items, &SimpleItem{
				Name:  resp.RuleOptions[i].Items[j].Name,
				Value: resp.RuleOptions[i].Items[j].Value,
			})
		}
		options = append(options, opt)
	}

	return options, nil
}

// CreateRules creates analytics rules for a configuration token.
func (c *Client) CreateRules(ctx context.Context, configToken string, rules []*AnalyticsRule) error {
	endpoint := c.getAnalyticsEndpoint()

	type SimpleItemEntry struct {
		XMLName xml.Name `xml:"SimpleItem"`
		Name    string   `xml:"Name,attr"`
		Value   string   `xml:"Value,attr"`
	}

	type ParametersWrapper struct {
		SimpleItems []SimpleItemEntry
	}

	type RuleEntry struct {
		XMLName    xml.Name          `xml:"tan:Rule"`
		Name       string            `xml:"Name,attr"`
		Type       string            `xml:"Type,attr"`
		Parameters ParametersWrapper `xml:"Parameters"`
	}

	type CreateRules struct {
		XMLName            xml.Name    `xml:"tan:CreateRules"`
		Xmlns              string      `xml:"xmlns:tan,attr"`
		ConfigurationToken string      `xml:"tan:ConfigurationToken"`
		Rules              []RuleEntry `xml:"tan:Rule"`
	}

	reqRules := make([]RuleEntry, 0, len(rules))
	for _, r := range rules {
		entry := RuleEntry{
			Name: r.Name,
			Type: r.Type,
		}
		for _, p := range r.Parameters {
			entry.Parameters.SimpleItems = append(entry.Parameters.SimpleItems, SimpleItemEntry{
				Name:  p.Name,
				Value: p.Value,
			})
		}
		reqRules = append(reqRules, entry)
	}

	req := CreateRules{
		Xmlns:              analyticsNamespace,
		ConfigurationToken: configToken,
		Rules:              reqRules,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("CreateRules failed: %w", err)
	}

	return nil
}

// ModifyRules modifies existing analytics rules for a configuration token.
func (c *Client) ModifyRules(ctx context.Context, configToken string, rules []*AnalyticsRule) error {
	endpoint := c.getAnalyticsEndpoint()

	type SimpleItemEntry struct {
		XMLName xml.Name `xml:"SimpleItem"`
		Name    string   `xml:"Name,attr"`
		Value   string   `xml:"Value,attr"`
	}

	type ParametersWrapper struct {
		SimpleItems []SimpleItemEntry
	}

	type RuleEntry struct {
		XMLName    xml.Name          `xml:"tan:Rule"`
		Name       string            `xml:"Name,attr"`
		Type       string            `xml:"Type,attr"`
		Parameters ParametersWrapper `xml:"Parameters"`
	}

	type ModifyRules struct {
		XMLName            xml.Name    `xml:"tan:ModifyRules"`
		Xmlns              string      `xml:"xmlns:tan,attr"`
		ConfigurationToken string      `xml:"tan:ConfigurationToken"`
		Rules              []RuleEntry `xml:"tan:Rule"`
	}

	reqRules := make([]RuleEntry, 0, len(rules))
	for _, r := range rules {
		entry := RuleEntry{
			Name: r.Name,
			Type: r.Type,
		}
		for _, p := range r.Parameters {
			entry.Parameters.SimpleItems = append(entry.Parameters.SimpleItems, SimpleItemEntry{
				Name:  p.Name,
				Value: p.Value,
			})
		}
		reqRules = append(reqRules, entry)
	}

	req := ModifyRules{
		Xmlns:              analyticsNamespace,
		ConfigurationToken: configToken,
		Rules:              reqRules,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("ModifyRules failed: %w", err)
	}

	return nil
}

// DeleteRules deletes analytics rules by name for a configuration token.
func (c *Client) DeleteRules(ctx context.Context, configToken string, ruleNames []string) error {
	endpoint := c.getAnalyticsEndpoint()

	type DeleteRules struct {
		XMLName            xml.Name `xml:"tan:DeleteRules"`
		Xmlns              string   `xml:"xmlns:tan,attr"`
		ConfigurationToken string   `xml:"tan:ConfigurationToken"`
		RuleNames          []string `xml:"tan:RuleName"`
	}

	req := DeleteRules{
		Xmlns:              analyticsNamespace,
		ConfigurationToken: configToken,
		RuleNames:          ruleNames,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("DeleteRules failed: %w", err)
	}

	return nil
}

// GetSupportedAnalyticsModules retrieves the supported analytics modules for a configuration token.
func (c *Client) GetSupportedAnalyticsModules(ctx context.Context, configToken string) ([]*SupportedAnalyticsModule, error) {
	endpoint := c.getAnalyticsEndpoint()

	type GetSupportedAnalyticsModules struct {
		XMLName            xml.Name `xml:"tan:GetSupportedAnalyticsModules"`
		Xmlns              string   `xml:"xmlns:tan,attr"`
		ConfigurationToken string   `xml:"tan:ConfigurationToken"`
	}

	type SimpleItemDescription struct {
		Name string `xml:"Name,attr"`
		Type string `xml:"Type,attr"`
	}

	type ModuleDescription struct {
		Name       string                  `xml:"Name,attr"`
		Parameters []SimpleItemDescription `xml:"Parameters>SimpleItemDescription"`
	}

	type GetSupportedAnalyticsModulesResponse struct {
		XMLName           xml.Name            `xml:"GetSupportedAnalyticsModulesResponse"`
		ModuleDescription []ModuleDescription `xml:"SupportedAnalyticsModules>AnalyticsModuleDescription"`
	}

	req := GetSupportedAnalyticsModules{
		Xmlns:              analyticsNamespace,
		ConfigurationToken: configToken,
	}

	var resp GetSupportedAnalyticsModulesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetSupportedAnalyticsModules failed: %w", err)
	}

	modules := make([]*SupportedAnalyticsModule, 0, len(resp.ModuleDescription))
	for i := range resp.ModuleDescription {
		mod := &SupportedAnalyticsModule{
			Name: resp.ModuleDescription[i].Name,
		}
		for j := range resp.ModuleDescription[i].Parameters {
			mod.Parameters = append(mod.Parameters, &SimpleItem{
				Name:  resp.ModuleDescription[i].Parameters[j].Name,
				Value: resp.ModuleDescription[i].Parameters[j].Type,
			})
		}
		modules = append(modules, mod)
	}

	return modules, nil
}

// GetAnalyticsModules retrieves the analytics modules for a configuration token.
func (c *Client) GetAnalyticsModules(ctx context.Context, configToken string) ([]*AnalyticsModule, error) {
	endpoint := c.getAnalyticsEndpoint()

	type GetAnalyticsModules struct {
		XMLName            xml.Name `xml:"tan:GetAnalyticsModules"`
		Xmlns              string   `xml:"xmlns:tan,attr"`
		ConfigurationToken string   `xml:"tan:ConfigurationToken"`
	}

	type SimpleItemEntry struct {
		Name  string `xml:"Name,attr"`
		Value string `xml:"Value,attr"`
	}

	type ModuleEntry struct {
		Name       string            `xml:"Name,attr"`
		Type       string            `xml:"Type,attr"`
		Parameters []SimpleItemEntry `xml:"Parameters>SimpleItem"`
	}

	type GetAnalyticsModulesResponse struct {
		XMLName         xml.Name      `xml:"GetAnalyticsModulesResponse"`
		AnalyticsModule []ModuleEntry `xml:"AnalyticsModule"`
	}

	req := GetAnalyticsModules{
		Xmlns:              analyticsNamespace,
		ConfigurationToken: configToken,
	}

	var resp GetAnalyticsModulesResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAnalyticsModules failed: %w", err)
	}

	modules := make([]*AnalyticsModule, 0, len(resp.AnalyticsModule))
	for i := range resp.AnalyticsModule {
		mod := &AnalyticsModule{
			Name: resp.AnalyticsModule[i].Name,
			Type: resp.AnalyticsModule[i].Type,
		}
		for j := range resp.AnalyticsModule[i].Parameters {
			mod.Parameters = append(mod.Parameters, &SimpleItem{
				Name:  resp.AnalyticsModule[i].Parameters[j].Name,
				Value: resp.AnalyticsModule[i].Parameters[j].Value,
			})
		}
		modules = append(modules, mod)
	}

	return modules, nil
}

// GetAnalyticsModuleOptions retrieves module options for a configuration token and optional module type filter.
func (c *Client) GetAnalyticsModuleOptions(ctx context.Context, configToken string, moduleType *string) ([]*AnalyticsModuleOptions, error) {
	endpoint := c.getAnalyticsEndpoint()

	type GetAnalyticsModuleOptions struct {
		XMLName            xml.Name `xml:"tan:GetAnalyticsModuleOptions"`
		Xmlns              string   `xml:"xmlns:tan,attr"`
		ConfigurationToken string   `xml:"tan:ConfigurationToken"`
		Type               *string  `xml:"tan:Type,omitempty"`
	}

	type SimpleItemEntry struct {
		Name  string `xml:"Name,attr"`
		Value string `xml:"Value,attr"`
	}

	type ModuleOptionsEntry struct {
		ModuleType string            `xml:"Type,attr"`
		Items      []SimpleItemEntry `xml:"SimpleItem"`
	}

	type GetAnalyticsModuleOptionsResponse struct {
		XMLName xml.Name             `xml:"GetAnalyticsModuleOptionsResponse"`
		Options []ModuleOptionsEntry `xml:"Options"`
	}

	req := GetAnalyticsModuleOptions{
		Xmlns:              analyticsNamespace,
		ConfigurationToken: configToken,
		Type:               moduleType,
	}

	var resp GetAnalyticsModuleOptionsResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetAnalyticsModuleOptions failed: %w", err)
	}

	options := make([]*AnalyticsModuleOptions, 0, len(resp.Options))
	for i := range resp.Options {
		opt := &AnalyticsModuleOptions{
			ModuleType: resp.Options[i].ModuleType,
		}
		for j := range resp.Options[i].Items {
			opt.Items = append(opt.Items, &SimpleItem{
				Name:  resp.Options[i].Items[j].Name,
				Value: resp.Options[i].Items[j].Value,
			})
		}
		options = append(options, opt)
	}

	return options, nil
}

// CreateAnalyticsModules creates analytics modules for a configuration token.
func (c *Client) CreateAnalyticsModules(ctx context.Context, configToken string, modules []*AnalyticsModule) error {
	endpoint := c.getAnalyticsEndpoint()

	type SimpleItemEntry struct {
		XMLName xml.Name `xml:"SimpleItem"`
		Name    string   `xml:"Name,attr"`
		Value   string   `xml:"Value,attr"`
	}

	type ParametersWrapper struct {
		SimpleItems []SimpleItemEntry
	}

	type ModuleEntry struct {
		XMLName    xml.Name          `xml:"tan:AnalyticsModule"`
		Name       string            `xml:"Name,attr"`
		Type       string            `xml:"Type,attr"`
		Parameters ParametersWrapper `xml:"Parameters"`
	}

	type CreateAnalyticsModules struct {
		XMLName            xml.Name      `xml:"tan:CreateAnalyticsModules"`
		Xmlns              string        `xml:"xmlns:tan,attr"`
		ConfigurationToken string        `xml:"tan:ConfigurationToken"`
		Modules            []ModuleEntry `xml:"tan:AnalyticsModule"`
	}

	reqModules := make([]ModuleEntry, 0, len(modules))
	for _, m := range modules {
		entry := ModuleEntry{
			Name: m.Name,
			Type: m.Type,
		}
		for _, p := range m.Parameters {
			entry.Parameters.SimpleItems = append(entry.Parameters.SimpleItems, SimpleItemEntry{
				Name:  p.Name,
				Value: p.Value,
			})
		}
		reqModules = append(reqModules, entry)
	}

	req := CreateAnalyticsModules{
		Xmlns:              analyticsNamespace,
		ConfigurationToken: configToken,
		Modules:            reqModules,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("CreateAnalyticsModules failed: %w", err)
	}

	return nil
}

// ModifyAnalyticsModules modifies existing analytics modules for a configuration token.
func (c *Client) ModifyAnalyticsModules(ctx context.Context, configToken string, modules []*AnalyticsModule) error {
	endpoint := c.getAnalyticsEndpoint()

	type SimpleItemEntry struct {
		XMLName xml.Name `xml:"SimpleItem"`
		Name    string   `xml:"Name,attr"`
		Value   string   `xml:"Value,attr"`
	}

	type ParametersWrapper struct {
		SimpleItems []SimpleItemEntry
	}

	type ModuleEntry struct {
		XMLName    xml.Name          `xml:"tan:AnalyticsModule"`
		Name       string            `xml:"Name,attr"`
		Type       string            `xml:"Type,attr"`
		Parameters ParametersWrapper `xml:"Parameters"`
	}

	type ModifyAnalyticsModules struct {
		XMLName            xml.Name      `xml:"tan:ModifyAnalyticsModules"`
		Xmlns              string        `xml:"xmlns:tan,attr"`
		ConfigurationToken string        `xml:"tan:ConfigurationToken"`
		Modules            []ModuleEntry `xml:"tan:AnalyticsModule"`
	}

	reqModules := make([]ModuleEntry, 0, len(modules))
	for _, m := range modules {
		entry := ModuleEntry{
			Name: m.Name,
			Type: m.Type,
		}
		for _, p := range m.Parameters {
			entry.Parameters.SimpleItems = append(entry.Parameters.SimpleItems, SimpleItemEntry{
				Name:  p.Name,
				Value: p.Value,
			})
		}
		reqModules = append(reqModules, entry)
	}

	req := ModifyAnalyticsModules{
		Xmlns:              analyticsNamespace,
		ConfigurationToken: configToken,
		Modules:            reqModules,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("ModifyAnalyticsModules failed: %w", err)
	}

	return nil
}

// DeleteAnalyticsModules deletes analytics modules by name for a configuration token.
func (c *Client) DeleteAnalyticsModules(ctx context.Context, configToken string, moduleNames []string) error {
	endpoint := c.getAnalyticsEndpoint()

	type DeleteAnalyticsModules struct {
		XMLName            xml.Name `xml:"tan:DeleteAnalyticsModules"`
		Xmlns              string   `xml:"xmlns:tan,attr"`
		ConfigurationToken string   `xml:"tan:ConfigurationToken"`
		ModuleNames        []string `xml:"tan:AnalyticsModuleName"`
	}

	req := DeleteAnalyticsModules{
		Xmlns:              analyticsNamespace,
		ConfigurationToken: configToken,
		ModuleNames:        moduleNames,
	}

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, nil); err != nil {
		return fmt.Errorf("DeleteAnalyticsModules failed: %w", err)
	}

	return nil
}

// GetSupportedMetadata retrieves supported metadata types for an analytics configuration.
// The configToken parameter is used as a type filter (e.g. "tt:AnalyticsModule").
func (c *Client) GetSupportedMetadata(ctx context.Context, configToken string) (*SupportedMetadata, error) {
	endpoint := c.getAnalyticsEndpoint()

	type GetSupportedMetadata struct {
		XMLName xml.Name `xml:"tan:GetSupportedMetadata"`
		Xmlns   string   `xml:"xmlns:tan,attr"`
		Type    string   `xml:"tan:Type"`
	}

	type GetSupportedMetadataResponse struct {
		XMLName          xml.Name `xml:"GetSupportedMetadataResponse"`
		SupportedMetadata struct {
			AnalyticsModule []string `xml:"AnalyticsModule"`
		} `xml:"SupportedMetadata"`
	}

	req := GetSupportedMetadata{
		Xmlns: analyticsNamespace,
		Type:  configToken,
	}

	var resp GetSupportedMetadataResponse

	username, password := c.GetCredentials()
	soapClient := soap.NewClient(c.httpClient, username, password)

	if err := soapClient.Call(ctx, endpoint, "", req, &resp); err != nil {
		return nil, fmt.Errorf("GetSupportedMetadata failed: %w", err)
	}

	return &SupportedMetadata{
		AnalyticsModules: resp.SupportedMetadata.AnalyticsModule,
	}, nil
}
