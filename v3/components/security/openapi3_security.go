package security

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// EnhancedSecurityScheme with full OpenAPI 3.0 support
type EnhancedSecurityScheme struct {
	Type             string              `json:"type"` // REQUIRED
	Description      string              `json:"description,omitempty"`
	Name             string              `json:"name,omitempty"`             // For apiKey
	In               string              `json:"in,omitempty"`               // For apiKey
	Scheme           string              `json:"scheme,omitempty"`           // For http
	BearerFormat     string              `json:"bearerFormat,omitempty"`     // For http bearer
	Flows            *EnhancedOAuthFlows `json:"flows,omitempty"`            // For oauth2
	OpenIdConnectUrl string              `json:"openIdConnectUrl,omitempty"` // For openIdConnect
}

// EnhancedOAuthFlows with validation
type EnhancedOAuthFlows struct {
	Implicit          *EnhancedOAuthFlow `json:"implicit,omitempty"`
	Password          *EnhancedOAuthFlow `json:"password,omitempty"`
	ClientCredentials *EnhancedOAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *EnhancedOAuthFlow `json:"authorizationCode,omitempty"`
}

// EnhancedOAuthFlow with validation
type EnhancedOAuthFlow struct {
	AuthorizationUrl string            `json:"authorizationUrl,omitempty"`
	TokenUrl         string            `json:"tokenUrl,omitempty"`
	RefreshUrl       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes"` // REQUIRED
}

// SecuritySchemeType represents different types of security schemes
// (Using existing SecuritySchemeType from security.go)

// APIKeyLocation represents where the API key should be located
type APIKeyLocation string

const (
	QueryLocation  APIKeyLocation = "query"
	HeaderLocation APIKeyLocation = "header"
	CookieLocation APIKeyLocation = "cookie"
)

// HTTPScheme represents different HTTP authentication schemes
type HTTPScheme string

const (
	BasicScheme  HTTPScheme = "basic"
	BearerScheme HTTPScheme = "bearer"
	DigestScheme HTTPScheme = "digest"
)

// OAuthFlowType represents different OAuth2 flow types
type OAuthFlowType string

const (
	ImplicitFlow          OAuthFlowType = "implicit"
	PasswordFlow          OAuthFlowType = "password"
	ClientCredentialsFlow OAuthFlowType = "clientCredentials"
	AuthorizationCodeFlow OAuthFlowType = "authorizationCode"
)

// NewAPIKeySecurityScheme creates a new API key security scheme
func NewAPIKeySecurityScheme(name string, location APIKeyLocation, description string) *EnhancedSecurityScheme {
	return &EnhancedSecurityScheme{
		Type:        string(APIKey),
		Name:        name,
		In:          string(location),
		Description: description,
	}
}

// NewHTTPSecurityScheme creates a new HTTP security scheme
func NewHTTPSecurityScheme(scheme HTTPScheme, description string) *EnhancedSecurityScheme {
	return &EnhancedSecurityScheme{
		Type:        string(HTTP),
		Scheme:      string(scheme),
		Description: description,
	}
}

// NewHTTPBearerSecurityScheme creates a new HTTP Bearer security scheme
func NewHTTPBearerSecurityScheme(bearerFormat, description string) *EnhancedSecurityScheme {
	return &EnhancedSecurityScheme{
		Type:         string(HTTP),
		Scheme:       string(BearerScheme),
		BearerFormat: bearerFormat,
		Description:  description,
	}
}

// NewOAuth2SecurityScheme creates a new OAuth2 security scheme
func NewOAuth2SecurityScheme(flows *EnhancedOAuthFlows, description string) *EnhancedSecurityScheme {
	return &EnhancedSecurityScheme{
		Type:        string(OAuth2),
		Flows:       flows,
		Description: description,
	}
}

// NewOpenIDConnectSecurityScheme creates a new OpenID Connect security scheme
func NewOpenIDConnectSecurityScheme(openIdConnectUrl, description string) *EnhancedSecurityScheme {
	return &EnhancedSecurityScheme{
		Type:             string(OpenIDConnect),
		OpenIdConnectUrl: openIdConnectUrl,
		Description:      description,
	}
}

// NewEnhancedOAuthFlows creates a new OAuth flows object
func NewEnhancedOAuthFlows() *EnhancedOAuthFlows {
	return &EnhancedOAuthFlows{}
}

// SetImplicitFlow sets the implicit flow
func (flows *EnhancedOAuthFlows) SetImplicitFlow(authorizationUrl string, scopes map[string]string) *EnhancedOAuthFlows {
	flows.Implicit = &EnhancedOAuthFlow{
		AuthorizationUrl: authorizationUrl,
		Scopes:           scopes,
	}
	return flows
}

// SetPasswordFlow sets the password flow
func (flows *EnhancedOAuthFlows) SetPasswordFlow(tokenUrl string, scopes map[string]string) *EnhancedOAuthFlows {
	flows.Password = &EnhancedOAuthFlow{
		TokenUrl: tokenUrl,
		Scopes:   scopes,
	}
	return flows
}

// SetClientCredentialsFlow sets the client credentials flow
func (flows *EnhancedOAuthFlows) SetClientCredentialsFlow(tokenUrl string, scopes map[string]string) *EnhancedOAuthFlows {
	flows.ClientCredentials = &EnhancedOAuthFlow{
		TokenUrl: tokenUrl,
		Scopes:   scopes,
	}
	return flows
}

// SetAuthorizationCodeFlow sets the authorization code flow
func (flows *EnhancedOAuthFlows) SetAuthorizationCodeFlow(authorizationUrl, tokenUrl string, scopes map[string]string) *EnhancedOAuthFlows {
	flows.AuthorizationCode = &EnhancedOAuthFlow{
		AuthorizationUrl: authorizationUrl,
		TokenUrl:         tokenUrl,
		Scopes:           scopes,
	}
	return flows
}

// SetRefreshUrl sets the refresh URL for a flow
func (flow *EnhancedOAuthFlow) SetRefreshUrl(refreshUrl string) *EnhancedOAuthFlow {
	flow.RefreshUrl = refreshUrl
	return flow
}

// Validate validates the security scheme according to OpenAPI 3.0.3 rules
func (s *EnhancedSecurityScheme) Validate() error {
	switch SecuritySchemeType(s.Type) {
	case APIKey:
		return s.validateAPIKey()
	case HTTP:
		return s.validateHTTP()
	case OAuth2:
		return s.validateOAuth2()
	case OpenIDConnect:
		return s.validateOpenIDConnect()
	default:
		return fmt.Errorf("invalid security scheme type: %s", s.Type)
	}
}

// validateAPIKey validates API key security scheme
func (s *EnhancedSecurityScheme) validateAPIKey() error {
	if s.Name == "" {
		return errors.New("apiKey security scheme requires 'name' field")
	}

	if s.In == "" {
		return errors.New("apiKey security scheme requires 'in' field")
	}

	validLocations := []string{string(QueryLocation), string(HeaderLocation), string(CookieLocation)}
	isValidLocation := false
	for _, location := range validLocations {
		if s.In == location {
			isValidLocation = true
			break
		}
	}

	if !isValidLocation {
		return fmt.Errorf("apiKey 'in' must be one of: %s", strings.Join(validLocations, ", "))
	}

	return nil
}

// validateHTTP validates HTTP security scheme
func (s *EnhancedSecurityScheme) validateHTTP() error {
	if s.Scheme == "" {
		return errors.New("http security scheme requires 'scheme' field")
	}

	// Validate common HTTP schemes
	validSchemes := []string{string(BasicScheme), string(BearerScheme), string(DigestScheme)}
	isValidScheme := false
	for _, scheme := range validSchemes {
		if strings.EqualFold(s.Scheme, scheme) {
			isValidScheme = true
			break
		}
	}

	// Allow custom schemes but validate they follow RFC7235 format
	if !isValidScheme {
		if !isValidHTTPScheme(s.Scheme) {
			return fmt.Errorf("invalid HTTP scheme: %s", s.Scheme)
		}
	}

	return nil
}

// validateOAuth2 validates OAuth2 security scheme
func (s *EnhancedSecurityScheme) validateOAuth2() error {
	if s.Flows == nil {
		return errors.New("oauth2 security scheme requires 'flows' field")
	}

	return s.Flows.Validate()
}

// validateOpenIDConnect validates OpenID Connect security scheme
func (s *EnhancedSecurityScheme) validateOpenIDConnect() error {
	if s.OpenIdConnectUrl == "" {
		return errors.New("openIdConnect security scheme requires 'openIdConnectUrl' field")
	}

	// Validate URL format
	if _, err := url.Parse(s.OpenIdConnectUrl); err != nil {
		return fmt.Errorf("invalid openIdConnectUrl: %w", err)
	}

	return nil
}

// Validate validates OAuth flows
func (flows *EnhancedOAuthFlows) Validate() error {
	hasAtLeastOneFlow := flows.Implicit != nil || flows.Password != nil ||
		flows.ClientCredentials != nil || flows.AuthorizationCode != nil

	if !hasAtLeastOneFlow {
		return errors.New("OAuth2 flows must contain at least one flow")
	}

	if flows.Implicit != nil {
		if err := flows.Implicit.validateImplicitFlow(); err != nil {
			return fmt.Errorf("invalid implicit flow: %w", err)
		}
	}

	if flows.Password != nil {
		if err := flows.Password.validatePasswordFlow(); err != nil {
			return fmt.Errorf("invalid password flow: %w", err)
		}
	}

	if flows.ClientCredentials != nil {
		if err := flows.ClientCredentials.validateClientCredentialsFlow(); err != nil {
			return fmt.Errorf("invalid clientCredentials flow: %w", err)
		}
	}

	if flows.AuthorizationCode != nil {
		if err := flows.AuthorizationCode.validateAuthorizationCodeFlow(); err != nil {
			return fmt.Errorf("invalid authorizationCode flow: %w", err)
		}
	}

	return nil
}

// validateImplicitFlow validates implicit OAuth flow
func (flow *EnhancedOAuthFlow) validateImplicitFlow() error {
	if flow.AuthorizationUrl == "" {
		return errors.New("implicit flow requires authorizationUrl")
	}

	if _, err := url.Parse(flow.AuthorizationUrl); err != nil {
		return fmt.Errorf("invalid authorizationUrl: %w", err)
	}

	if len(flow.Scopes) == 0 {
		return errors.New("implicit flow requires at least one scope")
	}

	return nil
}

// validatePasswordFlow validates password OAuth flow
func (flow *EnhancedOAuthFlow) validatePasswordFlow() error {
	if flow.TokenUrl == "" {
		return errors.New("password flow requires tokenUrl")
	}

	if _, err := url.Parse(flow.TokenUrl); err != nil {
		return fmt.Errorf("invalid tokenUrl: %w", err)
	}

	if len(flow.Scopes) == 0 {
		return errors.New("password flow requires at least one scope")
	}

	return nil
}

// validateClientCredentialsFlow validates client credentials OAuth flow
func (flow *EnhancedOAuthFlow) validateClientCredentialsFlow() error {
	if flow.TokenUrl == "" {
		return errors.New("clientCredentials flow requires tokenUrl")
	}

	if _, err := url.Parse(flow.TokenUrl); err != nil {
		return fmt.Errorf("invalid tokenUrl: %w", err)
	}

	if len(flow.Scopes) == 0 {
		return errors.New("clientCredentials flow requires at least one scope")
	}

	return nil
}

// validateAuthorizationCodeFlow validates authorization code OAuth flow
func (flow *EnhancedOAuthFlow) validateAuthorizationCodeFlow() error {
	if flow.AuthorizationUrl == "" {
		return errors.New("authorizationCode flow requires authorizationUrl")
	}

	if flow.TokenUrl == "" {
		return errors.New("authorizationCode flow requires tokenUrl")
	}

	if _, err := url.Parse(flow.AuthorizationUrl); err != nil {
		return fmt.Errorf("invalid authorizationUrl: %w", err)
	}

	if _, err := url.Parse(flow.TokenUrl); err != nil {
		return fmt.Errorf("invalid tokenUrl: %w", err)
	}

	if len(flow.Scopes) == 0 {
		return errors.New("authorizationCode flow requires at least one scope")
	}

	return nil
}

// isValidHTTPScheme validates HTTP scheme according to RFC7235
func isValidHTTPScheme(scheme string) bool {
	if scheme == "" {
		return false
	}

	// RFC7235: scheme = token
	// token = 1*tchar
	// tchar = "!" / "#" / "$" / "%" / "&" / "'" / "*" / "+" / "-" / "." /
	//         "^" / "_" / "`" / "|" / "~" / DIGIT / ALPHA
	for _, char := range scheme {
		if !isTokenChar(char) {
			return false
		}
	}

	return true
}

// isTokenChar checks if a character is valid for HTTP token
func isTokenChar(char rune) bool {
	return (char >= 'A' && char <= 'Z') ||
		(char >= 'a' && char <= 'z') ||
		(char >= '0' && char <= '9') ||
		char == '!' || char == '#' || char == '$' || char == '%' ||
		char == '&' || char == '\'' || char == '*' || char == '+' ||
		char == '-' || char == '.' || char == '^' || char == '_' ||
		char == '`' || char == '|' || char == '~'
}
