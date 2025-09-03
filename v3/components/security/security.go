package security

// SecuritySchemeType represents the type of security scheme
type SecuritySchemeType string

const (
	APIKey        SecuritySchemeType = "apiKey"
	HTTP          SecuritySchemeType = "http"
	OAuth2        SecuritySchemeType = "oauth2"
	OpenIDConnect SecuritySchemeType = "openIdConnect"
)

// SecuritySchemeIn represents where the API key is located
type SecuritySchemeIn string

const (
	Query  SecuritySchemeIn = "query"
	Header SecuritySchemeIn = "header"
	Cookie SecuritySchemeIn = "cookie"
)

// SecurityScheme represents a security scheme in OpenAPI 3.0
// https://spec.openapis.org/oas/v3.0.3#security-scheme-object
type SecurityScheme struct {
	Type             SecuritySchemeType `json:"type"`
	Description      string             `json:"description,omitempty"`
	Name             string             `json:"name,omitempty"`             // For apiKey
	In               SecuritySchemeIn   `json:"in,omitempty"`               // For apiKey
	Scheme           string             `json:"scheme,omitempty"`           // For http
	BearerFormat     string             `json:"bearerFormat,omitempty"`     // For http bearer
	Flows            *OAuthFlows        `json:"flows,omitempty"`            // For oauth2
	OpenIdConnectUrl string             `json:"openIdConnectUrl,omitempty"` // For openIdConnect
}

// OAuthFlows represents OAuth2 flows in OpenAPI 3.0
// https://spec.openapis.org/oas/v3.0.3#oauth-flows-object
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

// OAuthFlow represents a single OAuth2 flow
// https://spec.openapis.org/oas/v3.0.3#oauth-flow-object
type OAuthFlow struct {
	AuthorizationUrl string            `json:"authorizationUrl,omitempty"`
	TokenUrl         string            `json:"tokenUrl,omitempty"`
	RefreshUrl       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes"`
}

// BasicAuth creates a basic authentication security scheme
func BasicAuth(description string) SecurityScheme {
	return SecurityScheme{
		Type:        HTTP,
		Scheme:      "basic",
		Description: description,
	}
}

// BearerAuth creates a bearer token authentication security scheme
func BearerAuth(description string, bearerFormat string) SecurityScheme {
	return SecurityScheme{
		Type:         HTTP,
		Scheme:       "bearer",
		BearerFormat: bearerFormat,
		Description:  description,
	}
}

// APIKeyAuth creates an API key authentication security scheme
func APIKeyAuth(name string, in SecuritySchemeIn, description string) SecurityScheme {
	return SecurityScheme{
		Type:        APIKey,
		Name:        name,
		In:          in,
		Description: description,
	}
}

// OAuth2Auth creates an OAuth2 authentication security scheme
func OAuth2Auth(flows *OAuthFlows, description string) SecurityScheme {
	return SecurityScheme{
		Type:        OAuth2,
		Flows:       flows,
		Description: description,
	}
}

// OpenIDConnectAuth creates an OpenID Connect authentication security scheme
func OpenIDConnectAuth(url string, description string) SecurityScheme {
	return SecurityScheme{
		Type:             OpenIDConnect,
		OpenIdConnectUrl: url,
		Description:      description,
	}
}

// NewOAuthFlows creates a new OAuth flows object
func NewOAuthFlows() *OAuthFlows {
	return &OAuthFlows{}
}

// WithImplicit adds implicit flow to OAuth flows
func (f *OAuthFlows) WithImplicit(authUrl string, scopes map[string]string) *OAuthFlows {
	f.Implicit = &OAuthFlow{
		AuthorizationUrl: authUrl,
		Scopes:           scopes,
	}
	return f
}

// WithPassword adds password flow to OAuth flows
func (f *OAuthFlows) WithPassword(tokenUrl string, scopes map[string]string) *OAuthFlows {
	f.Password = &OAuthFlow{
		TokenUrl: tokenUrl,
		Scopes:   scopes,
	}
	return f
}

// WithClientCredentials adds client credentials flow to OAuth flows
func (f *OAuthFlows) WithClientCredentials(tokenUrl string, scopes map[string]string) *OAuthFlows {
	f.ClientCredentials = &OAuthFlow{
		TokenUrl: tokenUrl,
		Scopes:   scopes,
	}
	return f
}

// WithAuthorizationCode adds authorization code flow to OAuth flows
func (f *OAuthFlows) WithAuthorizationCode(authUrl string, tokenUrl string, scopes map[string]string) *OAuthFlows {
	f.AuthorizationCode = &OAuthFlow{
		AuthorizationUrl: authUrl,
		TokenUrl:         tokenUrl,
		Scopes:           scopes,
	}
	return f
}

// SetRefreshUrl sets the refresh URL for a flow
func (flow *OAuthFlow) SetRefreshUrl(refreshUrl string) {
	flow.RefreshUrl = refreshUrl
}

// Security represents security requirements for an operation
type Security struct {
	Schemes []map[string][]string `json:"-"`
}

// NewSecurity creates a new security requirement
func NewSecurity() *Security {
	return &Security{
		Schemes: make([]map[string][]string, 0),
	}
}

// AddScheme adds a security scheme requirement
func (s *Security) AddScheme(schemeName string, scopes []string) *Security {
	scheme := map[string][]string{
		schemeName: scopes,
	}
	s.Schemes = append(s.Schemes, scheme)
	return s
}

// GetSchemes returns the security schemes
func (s *Security) GetSchemes() []map[string][]string {
	return s.Schemes
}
