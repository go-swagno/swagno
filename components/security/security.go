package security

// Struct definition for BasicAuth
type BasicAuth struct {
	BasicAuth []string `json:"basicAuth"` // Slice of strings representing basic authentication
}

// Method for BasicAuth struct
func (b *BasicAuth) New(username string, password string) {
	b.BasicAuth = []string{username, password}
}

// Struct definition for ApiKeyAuth
type ApiKeyAuth struct {
	Name []string `json:"-"` // Slice of strings representing the name of API key
}

// Method for ApiKeyAuth struct
func (a *ApiKeyAuth) New(apiKey string) {
	a.Name = []string{apiKey}
}

// Struct definition for OAuth
type OAuth struct {
	Name   string   `json:"-"` // Name of OAuth
	Scopes []string `json:"-"` // Slice of strings representing OAuth scopes
}

// Method for OAuth struct
func (o *OAuth) New(name string, scopes []string) {
	o.Name = name
	o.Scopes = scopes
}

// Struct definition for Security
type Security struct {
	Schemes []map[string][]string `json:"-"` // Slice of maps representing security schemes
}

// Method for Security struct
func (s *Security) New(schemes []map[string][]string) {
	s.Schemes = schemes
}
