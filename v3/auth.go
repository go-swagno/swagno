package v3

// SetBasicAuth adds Basic authentication to the OpenAPI specification
func (o *OpenAPI) SetBasicAuth(description ...string) {
	desc := "Basic Authentication"
	if len(description) > 0 {
		desc = description[0]
	}

	if o.Components == nil {
		o.Components = &Components{SecuritySchemes: make(map[string]SecurityScheme)}
	}
	if o.Components.SecuritySchemes == nil {
		o.Components.SecuritySchemes = make(map[string]SecurityScheme)
	}

	o.Components.SecuritySchemes["basicAuth"] = SecurityScheme{
		Type:        "http",
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
		o.Components = &Components{SecuritySchemes: make(map[string]SecurityScheme)}
	}
	if o.Components.SecuritySchemes == nil {
		o.Components.SecuritySchemes = make(map[string]SecurityScheme)
	}

	scheme := SecurityScheme{
		Type:        "http",
		Scheme:      "bearer",
		Description: desc,
	}

	if format != "" {
		scheme.BearerFormat = format
	}

	o.Components.SecuritySchemes["bearerAuth"] = scheme
}

// SetApiKeyAuth adds API Key authentication to the OpenAPI specification
func (o *OpenAPI) SetApiKeyAuth(name string, in string, description ...string) {
	desc := "API Key Authentication"
	if len(description) > 0 {
		desc = description[0]
	}

	if o.Components == nil {
		o.Components = &Components{SecuritySchemes: make(map[string]SecurityScheme)}
	}
	if o.Components.SecuritySchemes == nil {
		o.Components.SecuritySchemes = make(map[string]SecurityScheme)
	}

	o.Components.SecuritySchemes["apiKeyAuth"] = SecurityScheme{
		Type:        "apiKey",
		Name:        name,
		In:          in,
		Description: desc,
	}
}

// SetOAuth2Auth adds OAuth2 authentication to the OpenAPI specification
func (o *OpenAPI) SetOAuth2Auth(flows *OAuthFlows, description ...string) {
	desc := "OAuth2 Authentication"
	if len(description) > 0 {
		desc = description[0]
	}

	if o.Components == nil {
		o.Components = &Components{SecuritySchemes: make(map[string]SecurityScheme)}
	}
	if o.Components.SecuritySchemes == nil {
		o.Components.SecuritySchemes = make(map[string]SecurityScheme)
	}

	o.Components.SecuritySchemes["oauth2"] = SecurityScheme{
		Type:        "oauth2",
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
		o.Components = &Components{SecuritySchemes: make(map[string]SecurityScheme)}
	}
	if o.Components.SecuritySchemes == nil {
		o.Components.SecuritySchemes = make(map[string]SecurityScheme)
	}

	o.Components.SecuritySchemes["openIdConnect"] = SecurityScheme{
		Type:             "openIdConnect",
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
