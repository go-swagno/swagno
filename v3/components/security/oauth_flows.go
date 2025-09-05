package security

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
