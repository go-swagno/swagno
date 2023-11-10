package swagno

// https://swagger.io/specification/v2/#scopes-object
type swaggerSecurityScope struct {
	Name        string
	Description string
}

// Scopes is a helper function that takes multiple swaggerSecurityScope objects and returns a map of scopes.
// Each scope consists of a name and description.
// Example usage: Scopes(Scope("read", "Read access"), Scope("write", "Write access"))
func Scopes(scopes ...swaggerSecurityScope) map[string]string {
	scopesMap := make(map[string]string)
	for _, scope := range scopes {
		scopesMap[scope.Name] = scope.Description
	}
	return scopesMap
}

// Scope is a helper function that creates a swaggerSecurityScope object with the specified name and description.
// This is used in conjunction with the Scopes function to define security scopes.
// Example usage: Scope("read", "Read access")
func Scope(name string, description string) swaggerSecurityScope {
	return swaggerSecurityScope{
		Name:        name,
		Description: description,
	}
}
