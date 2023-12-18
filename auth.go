package swagno

// https://swagger.io/specification/v2/#security-definitions-object

// SetBasicAuth sets the basic authentication security definition in the Swagger object.
// For more information, refer to: https://swagger.io/specification/v2/#basic-authentication-sample
func (s Swagger) SetBasicAuth(description ...string) {
	desc := "Basic Authentication"
	if len(description) > 0 {
		desc = description[0]
	}
	s.SecurityDefinitions["basicAuth"] = securityDefinition{
		Type:        "basic",
		Description: desc,
	}
}

// SetApiKeyAuth sets the API key authentication security definition in the Swagger object.
// For more information, refer to: https://swagger.io/specification/v2/#api-key-sample
func (s Swagger) SetApiKeyAuth(name string, in string, description ...string) {
	desc := "API Key Authentication"
	if len(description) > 0 {
		desc = description[0]
	}
	s.SecurityDefinitions[name] = securityDefinition{
		Type:        "apiKey",
		Name:        name,
		In:          in,
		Description: desc,
	}
}

// SetOAuth2Auth sets the OAuth2 authentication security definition in the Swagger object.
// For more information, refer to: https://swagger.io/specification/v2/#implicit-oauth2-sample
func (s Swagger) SetOAuth2Auth(name string, flow string, authorizationUrl string, tokenUrl string, scopes map[string]string, description ...string) {
	desc := "OAuth2 Authentication"
	if len(description) > 0 {
		desc = description[0]
	}

	definition := securityDefinition{
		Type:        "oauth2",
		Flow:        flow,
		Scopes:      scopes,
		Description: desc,
	}

	if flow == "implicit" || flow == "accessCode" {
		definition.AuthorizationUrl = authorizationUrl
	}
	if flow == "password" || flow == "accessCode" || flow == "application" {
		definition.TokenUrl = tokenUrl
	}
	s.SecurityDefinitions[name] = definition
}
