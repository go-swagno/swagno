package swagno

// https://swagger.io/specification/v2/#security-definitions-object

// https://swagger.io/specification/v2/#basic-authentication-sample
func (s Swagger) SetBasicAuth(description ...string) {
	desc := "Basic Authentication"
	if len(description) > 0 {
		desc = description[0]
	}
	s.SecurityDefinitions["basicAuth"] = swaggerSecurityDefinition{
		Type:        "basic",
		Description: desc,
	}
}

// https://swagger.io/specification/v2/#api-key-sample
func (s Swagger) SetApiKeyAuth(name string, in string, description ...string) {
	desc := "API Key Authentication"
	if len(description) > 0 {
		desc = description[0]
	}
	s.SecurityDefinitions[name] = swaggerSecurityDefinition{
		Type:        "apiKey",
		Name:        name,
		In:          in,
		Description: desc,
	}
}


// https://swagger.io/specification/v2/#implicit-oauth2-sample
func (s Swagger) SetOAuth2Auth(name string, flow string, authorizationUrl string, tokenUrl string, scopes map[string]string, description ...string) {
	desc := "OAuth2 Authentication"
	if len(description) > 0 {
		desc = description[0]
	}

	definition := swaggerSecurityDefinition{
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