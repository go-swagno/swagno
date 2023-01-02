package swagno

// https://swagger.io/specification/v2/#scopes-object
type swaggerSecurityScope struct {
	Name        string
	Description string
}


func Scopes(scopes ...swaggerSecurityScope) map[string]string {
	scopesMap := make(map[string]string)
	for _, scope := range scopes {
		scopesMap[scope.Name] = scope.Description
	}
	return scopesMap
}

func Scope(name string, description string) swaggerSecurityScope {
	return swaggerSecurityScope{
		Name:        name,
		Description: description,
	}
}