package security

type SecuritySchemeName string

const (
	BasicAuth     SecuritySchemeName = "basicAuth"
	BearerAuth    SecuritySchemeName = "bearerAuth"
	APIKeyAuth    SecuritySchemeName = "apiKeyAuth"
	OAuth2        SecuritySchemeName = "oauth2"
	OpenIDConnect SecuritySchemeName = "openIdConnect"
)

type SecuritySchemeType string

const (
	SecuritySchemeType_APIKey        SecuritySchemeType = "apiKey"
	SecuritySchemeType_HTTP          SecuritySchemeType = "http"
	SecuritySchemeType_OAuth2        SecuritySchemeType = "oauth2"
	SecuritySchemeType_OpenIDConnect SecuritySchemeType = "openIdConnect"
)

type SecuritySchemeIn string

const (
	Query  SecuritySchemeIn = "query"
	Header SecuritySchemeIn = "header"
	Cookie SecuritySchemeIn = "cookie"
)
