package v3

import (
	"encoding/json"
	"log"
	"os"

	"github.com/go-swagno/swagno/v3/components/definition"
	"github.com/go-swagno/swagno/v3/components/endpoint"
	"github.com/go-swagno/swagno/v3/components/tag"
)

// The full JSON model for OpenAPI v3 documentation
// https://spec.openapis.org/oas/v3.0.3
type OpenAPI struct {
	OpenAPI    string                                      `json:"openapi" default:"3.0.3"`
	Info       Info                                        `json:"info"`
	Servers    []Server                                    `json:"servers,omitempty"`
	Paths      map[string]map[string]endpoint.JsonEndPoint `json:"paths"`
	Components *Components                                 `json:"components,omitempty"`
	Tags       []tag.Tag                                   `json:"tags,omitempty"`
	Security   []map[string][]string                       `json:"security,omitempty"`
	endpoints  []*endpoint.EndPoint
}

// Info represents the information about the API.
// https://spec.openapis.org/oas/v3.0.3#info-object
type Info struct {
	Title          string   `json:"title"`
	Description    string   `json:"description,omitempty"`
	Version        string   `json:"version"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

// Server represents a server object in OpenAPI 3.0
// https://spec.openapis.org/oas/v3.0.3#server-object
type Server struct {
	URL         string                    `json:"url"`
	Description string                    `json:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty"`
}

// ServerVariable represents a server variable object
// https://spec.openapis.org/oas/v3.0.3#server-variable-object
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default"`
	Description string   `json:"description,omitempty"`
}

// Components holds a set of reusable objects for different aspects of the OAS
// https://spec.openapis.org/oas/v3.0.3#components-object
type Components struct {
	Schemas         map[string]definition.Schema     `json:"schemas,omitempty"`
	Responses       map[string]endpoint.JsonResponse `json:"responses,omitempty"`
	Parameters      map[string]interface{}           `json:"parameters,omitempty"`
	Examples        map[string]interface{}           `json:"examples,omitempty"`
	RequestBodies   map[string]interface{}           `json:"requestBodies,omitempty"`
	Headers         map[string]interface{}           `json:"headers,omitempty"`
	SecuritySchemes map[string]SecurityScheme        `json:"securitySchemes,omitempty"`
	Links           map[string]interface{}           `json:"links,omitempty"`
	Callbacks       map[string]interface{}           `json:"callbacks,omitempty"`
}

// SecurityScheme represents a security scheme in OpenAPI 3.0
// https://spec.openapis.org/oas/v3.0.3#security-scheme-object
type SecurityScheme struct {
	Type             string      `json:"type"`
	Description      string      `json:"description,omitempty"`
	Name             string      `json:"name,omitempty"`
	In               string      `json:"in,omitempty"`
	Scheme           string      `json:"scheme,omitempty"`
	BearerFormat     string      `json:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `json:"flows,omitempty"`
	OpenIdConnectUrl string      `json:"openIdConnectUrl,omitempty"`
}

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

// New creates a new OpenAPI instance with the provided config
func New(c Config) *OpenAPI {
	return buildOpenAPI(c)
}

func (openapi *OpenAPI) AddTags(tags ...tag.Tag) {
	openapi.Tags = append(openapi.Tags, tags...)
}

// AddEndpoints adds endpoints to OpenAPI object
func (o *OpenAPI) AddEndpoints(e []*endpoint.EndPoint) {
	o.endpoints = append(o.endpoints, e...)
}

// AddEndpoint adds endpoint to OpenAPI object
func (o *OpenAPI) AddEndpoint(e *endpoint.EndPoint) {
	o.endpoints = append(o.endpoints, e)
}

// AddServer adds a server to the OpenAPI specification
func (o *OpenAPI) AddServer(url string, description string) {
	if o.Servers == nil {
		o.Servers = []Server{}
	}

	// Remove default server if this is the first real server being added
	if len(o.Servers) == 1 && o.Servers[0].URL == "/" && o.Servers[0].Description == "" {
		o.Servers = []Server{}
	}

	o.Servers = append(o.Servers, Server{
		URL:         url,
		Description: description,
	})
}

// Contact represents the contact information for the API.
// https://spec.openapis.org/oas/v3.0.3#contact-object
type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// License represents the license information for the API.
// https://spec.openapis.org/oas/v3.0.3#license-object
type License struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// Config struct represents the configuration for OpenAPI documentation.
type Config struct {
	Title          string   // title of the OpenAPI documentation
	Version        string   // version of the OpenAPI documentation
	Description    string   // description of the OpenAPI documentation
	Servers        []Server // servers for the API
	License        *License // license information for the OpenAPI documentation
	Contact        *Contact // contact information for the OpenAPI documentation
	TermsOfService string   // term of service information for the OpenAPI documentation
}

// buildOpenAPI creates a new OpenAPI instance with the given configuration.
func buildOpenAPI(c Config) (openapi *OpenAPI) {
	if c.Title == "" {
		c.Title = "OpenAPI API"
	}
	if c.Version == "" {
		c.Version = "1.0.0"
	}

	openapi = &OpenAPI{
		OpenAPI: "3.0.3",
		Info: Info{
			Title:          c.Title,
			Description:    c.Description,
			Version:        c.Version,
			License:        c.License,
			Contact:        c.Contact,
			TermsOfService: c.TermsOfService,
		},
		Servers: c.Servers,
		Paths:   make(map[string]map[string]endpoint.JsonEndPoint),
		Components: &Components{
			Schemas:         make(map[string]definition.Schema),
			SecuritySchemes: make(map[string]SecurityScheme),
		},
		Tags:      []tag.Tag{},
		endpoints: []*endpoint.EndPoint{},
	}

	// Set default server if none provided and none will be added later
	if len(openapi.Servers) == 0 {
		openapi.Servers = []Server{{URL: "/"}}
	}

	return
}

// ExportOpenAPIDocs exports the OpenAPI documentation as a JSON file.
func (openapi *OpenAPI) ExportOpenAPIDocs(out_file string) string {
	json, err := json.MarshalIndent(openapi, "", "  ")
	if err != nil {
		log.Println("Error while generating openapi json")
	}
	err = os.WriteFile(out_file, json, 0644)
	if err != nil {
		log.Println("Error writing openapi file")
	}
	return string(json)
}
