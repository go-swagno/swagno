# Security Structure in Swagno v3

## Overview

Security features in Swagno v3 are implemented with a dedicated security component package at `components/security/` and authentication methods in the main `auth.go` file. This document outlines the actual security implementation.

## Security Package Structure

### `components/security/security.go`

Contains security scheme types, names, and location constants.

#### Security Scheme Names

```go
type SecuritySchemeName string

const (
    BasicAuth     SecuritySchemeName = "basicAuth"
    BearerAuth    SecuritySchemeName = "bearerAuth"
    APIKeyAuth    SecuritySchemeName = "apiKeyAuth"
    OAuth2        SecuritySchemeName = "oauth2"
    OpenIDConnect SecuritySchemeName = "openIdConnect"
)
```

#### Security Scheme Types

```go
type SecuritySchemeType string

const (
    SecuritySchemeType_APIKey        SecuritySchemeType = "apiKey"
    SecuritySchemeType_HTTP          SecuritySchemeType = "http"
    SecuritySchemeType_OAuth2        SecuritySchemeType = "oauth2"
    SecuritySchemeType_OpenIDConnect SecuritySchemeType = "openIdConnect"
)
```

#### Security Scheme Locations

```go
type SecuritySchemeIn string

const (
    Query  SecuritySchemeIn = "query"
    Header SecuritySchemeIn = "header"
    Cookie SecuritySchemeIn = "cookie"
)
```

### `components/security/oauth_flows.go`

Contains OAuth2 flow structures and builder methods.

#### OAuth Flow Structures

```go
type OAuthFlows struct {
    Implicit          *OAuthFlow `json:"implicit,omitempty"`
    Password          *OAuthFlow `json:"password,omitempty"`
    ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
    AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

type OAuthFlow struct {
    AuthorizationUrl string            `json:"authorizationUrl,omitempty"`
    TokenUrl         string            `json:"tokenUrl,omitempty"`
    RefreshUrl       string            `json:"refreshUrl,omitempty"`
    Scopes           map[string]string `json:"scopes"`
}
```

## Available Security Methods (in `auth.go`)

### Basic Authentication

- `SetBasicAuth(description ...string)`

### Bearer Authentication

- `SetBearerAuth(format string, description ...string)`

### API Key Authentication

- `SetApiKeyAuth(name string, in security.SecuritySchemeIn, description ...string)`

### OAuth2 Authentication

- `SetOAuth2Auth(flows *security.OAuthFlows, description ...string)`

### OpenID Connect

- `SetOpenIdConnectAuth(openIdConnectUrl string, description ...string)`

### Global Security

- `AddGlobalSecurity(security map[string][]string)`

## OAuth2 Flow Builders (in `security` package)

### NewOAuthFlows()

Creates a new OAuth flows object.

### Flow Methods

- `WithImplicit(authUrl string, scopes map[string]string) *OAuthFlows`
- `WithPassword(tokenUrl string, scopes map[string]string) *OAuthFlows`
- `WithClientCredentials(tokenUrl string, scopes map[string]string) *OAuthFlows`
- `WithAuthorizationCode(authUrl string, tokenUrl string, scopes map[string]string) *OAuthFlows`

### Flow Instance Methods

- `SetRefreshUrl(refreshUrl string)` - Sets refresh URL for individual flows

## Example Usage

```go
import "github.com/go-swagno/swagno/v3/components/security"

// Basic setup
openapi := v3.New(v3.Config{Title: "My API", Version: "1.0.0"})

// Add Bearer auth
openapi.SetBearerAuth("JWT", "Bearer authentication using JWT tokens")

// Add API Key auth
openapi.SetApiKeyAuth("X-API-Key", security.Header, "API key authentication")

// Add OAuth2 with multiple flows
flows := security.NewOAuthFlows().
    WithAuthorizationCode(
        "https://example.com/oauth/authorize",
        "https://example.com/oauth/token",
        map[string]string{
            "read":  "Read access",
            "write": "Write access",
        },
    ).
    WithClientCredentials(
        "https://example.com/oauth/token",
        map[string]string{
            "api": "API access",
        },
    )

openapi.SetOAuth2Auth(flows, "OAuth2 authentication")

// Use in endpoints
endpoint.New(
    endpoint.GET,
    "/protected",
    endpoint.WithSecurity([]map[security.SecuritySchemeName][]string{
        {security.BearerAuth: {}},
        {security.APIKeyAuth: {}},
        {security.BasicAuth: {}},
        {security.OpenIDConnect: {}},
        {security.OAuth2: {"read", "write"}}, // Multiple scopes
    }),
)
```

## Note

The security implementation now uses a proper component structure with the `components/security/` package containing type definitions and OAuth2 flow builders, while authentication setup methods remain in `auth.go` for ease of use.
