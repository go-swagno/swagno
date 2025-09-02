package security

// https://spec.openapis.org/oas/v3.0.3#security-scheme-object
type openAPISecurityScope struct {
	Name        string
	Description string
}

// Scopes is a helper function that takes multiple openAPISecurityScope objects and returns a map of scopes.
// Each scope consists of a name and description.
// Example usage: Scopes(Scope("read", "Read access"), Scope("write", "Write access"))
func Scopes(scopes ...openAPISecurityScope) map[string]string {
	scopesMap := make(map[string]string)
	for _, scope := range scopes {
		scopesMap[scope.Name] = scope.Description
	}
	return scopesMap
}

// Scope is a helper function that creates a openAPISecurityScope object with the specified name and description.
// This is used in conjunction with the Scopes function to define security scopes.
// Example usage: Scope("read", "Read access")
func Scope(name string, description string) openAPISecurityScope {
	return openAPISecurityScope{
		Name:        name,
		Description: description,
	}
}
