package v3

import "github.com/go-swagno/swagno/v3/components/security"

// SetBasicAuth adds Basic authentication to the OpenAPI specification
func (o *OpenAPI) SetBasicAuth(description ...string) {
	desc := "Basic Authentication"
	if len(description) > 0 {
		desc = description[0]
	}

	if o.Components == nil {
		o.Components = &Components{SecuritySchemes: make(map[security.SecuritySchemeName]SecurityScheme)}
	}
	if o.Components.SecuritySchemes == nil {
		o.Components.SecuritySchemes = make(map[security.SecuritySchemeName]SecurityScheme)
	}

	o.Components.SecuritySchemes[security.BasicAuth] = SecurityScheme{
		Type:        security.SecuritySchemeType_HTTP,
		Scheme:      "basic",
		Description: desc,
	}
}

// SetBearerAuth adds Bearer authentication to the OpenAPI specification
func (o *OpenAPI) SetBearerAuth(format string, description ...string) {
	desc := "Bearer Authentication"
	if len(description) > 0 {
		desc = description[0]
	}

	if o.Components == nil {
		o.Components = &Components{SecuritySchemes: make(map[security.SecuritySchemeName]SecurityScheme)}
	}
	if o.Components.SecuritySchemes == nil {
		o.Components.SecuritySchemes = make(map[security.SecuritySchemeName]SecurityScheme)
	}

	scheme := SecurityScheme{
		Type:        security.SecuritySchemeType_HTTP,
		Scheme:      "bearer",
		Description: desc,
	}

	if format != "" {
		scheme.BearerFormat = format
	}

	o.Components.SecuritySchemes[security.BearerAuth] = scheme
}

// SetApiKeyAuth adds API Key authentication to the OpenAPI specification
func (o *OpenAPI) SetApiKeyAuth(name string, in security.SecuritySchemeIn, description ...string) {
	desc := "API Key Authentication"
	if len(description) > 0 {
		desc = description[0]
	}

	if o.Components == nil {
		o.Components = &Components{SecuritySchemes: make(map[security.SecuritySchemeName]SecurityScheme)}
	}
	if o.Components.SecuritySchemes == nil {
		o.Components.SecuritySchemes = make(map[security.SecuritySchemeName]SecurityScheme)
	}

	o.Components.SecuritySchemes[security.APIKeyAuth] = SecurityScheme{
		Type:        security.SecuritySchemeType_APIKey,
		Name:        name,
		In:          in,
		Description: desc,
	}
}

// SetOAuth2Auth adds OAuth2 authentication to the OpenAPI specification
func (o *OpenAPI) SetOAuth2Auth(flows *security.OAuthFlows, description ...string) {
	desc := "OAuth2 Authentication"
	if len(description) > 0 {
		desc = description[0]
	}

	if o.Components == nil {
		o.Components = &Components{SecuritySchemes: make(map[security.SecuritySchemeName]SecurityScheme)}
	}
	if o.Components.SecuritySchemes == nil {
		o.Components.SecuritySchemes = make(map[security.SecuritySchemeName]SecurityScheme)
	}

	o.Components.SecuritySchemes[security.OAuth2] = SecurityScheme{
		Type:        security.SecuritySchemeType_OAuth2,
		Flows:       flows,
		Description: desc,
	}
}

// SetOpenIdConnectAuth adds OpenID Connect authentication to the OpenAPI specification
func (o *OpenAPI) SetOpenIdConnectAuth(openIdConnectUrl string, description ...string) {
	desc := "OpenID Connect Authentication"
	if len(description) > 0 {
		desc = description[0]
	}

	if o.Components == nil {
		o.Components = &Components{SecuritySchemes: make(map[security.SecuritySchemeName]SecurityScheme)}
	}
	if o.Components.SecuritySchemes == nil {
		o.Components.SecuritySchemes = make(map[security.SecuritySchemeName]SecurityScheme)
	}

	o.Components.SecuritySchemes[security.OpenIDConnect] = SecurityScheme{
		Type:             security.SecuritySchemeType_OpenIDConnect,
		OpenIdConnectUrl: openIdConnectUrl,
		Description:      desc,
	}
}

// AddGlobalSecurity adds global security requirements to the OpenAPI specification
func (o *OpenAPI) AddGlobalSecurity(security map[string][]string) {
	if o.Security == nil {
		o.Security = []map[string][]string{}
	}
	o.Security = append(o.Security, security)
}
