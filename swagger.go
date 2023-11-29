package swagno

import (
	"encoding/json"
	"log"
	"os"

	"github.com/domhoward14/swagno/components/definition"
	"github.com/domhoward14/swagno/components/endpoint"
	"github.com/domhoward14/swagno/components/tag"
)

// The full JSON model for swagger v2 documentation
// https://swagger.io/docs/specification/2-0/basic-structure/
type Swagger struct {
	Swagger             string                                      `json:"swagger" default:"2.0"`
	Info                Info                                        `json:"info"`
	Paths               map[string]map[string]endpoint.JsonEndPoint `json:"paths"`
	BasePath            string                                      `json:"basePath" default:"/"`
	Host                string                                      `json:"host" default:""`
	Definitions         map[string]definition.Definition            `json:"definitions"`
	Schemes             []string                                    `json:"schemes,omitempty"`
	Tags                []tag.Tag                                   `json:"tags,omitempty"`
	SecurityDefinitions map[string]securityDefinition               `json:"securityDefinitions,omitempty"`
	endpoints           []*endpoint.EndPoint
}

// Info represents the information about the API.
// https://swagger.io/specification/v2/#info-object
type Info struct {
	Title          string   `json:"title"`
	Version        string   `json:"version"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

// Contact represents the contact information for the API.
// https://swagger.io/specification/v2/#contact-object
type contact struct {
	Name  string `json:"name,omitempty"`
	Url   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// License represents the license information for the API.
// https://swagger.io/specification/v2/#license-object
type license struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

// securityDefinition represents the security definition object in Swagger.
// https://swagger.io/specification/v2/#securityDefinitionsObject
type securityDefinition struct {
	Type             string            `json:"type"`
	Description      string            `json:"description,omitempty"`
	Name             string            `json:"name,omitempty"`
	In               string            `json:"in,omitempty"`
	Flow             string            `json:"flow,omitempty"`
	AuthorizationUrl string            `json:"authorizationUrl,omitempty"`
	TokenUrl         string            `json:"tokenUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

// New creates a new swagger instance with the provided config
func New(c Config) *Swagger {
	return buildSwagger(c)
}

func (swagger *Swagger) AddTags(tags ...tag.Tag) {
	swagger.Tags = append(swagger.Tags, tags...)
}

// Add EndPoint models to Swagger endpoints
func (s *Swagger) AddEndpoints(e []*endpoint.EndPoint) {
	s.endpoints = append(s.endpoints, e...)
}

func (s *Swagger) AddEndpoint(e *endpoint.EndPoint) {
	s.endpoints = append(s.endpoints, e)
}

// Contact struct represents the contact information for Swagger documentation.
type Contact struct {
	Name string `json:"name"` // name of the contact person
}

// License struct represents the license information for Swagger documentation.
type License struct {
	Name string `json:"name"` // name of the license
	URL  string `json:"url"`  // URL for the license
}

// Config struct represents the configuration for Swagger documentation.
type Config struct {
	Title   string   // title of the Swagger documentation
	Version string   // version of the Swagger documentation
	Host    string   // host URL for the API
	Path    string   // path to the Swagger JSON file
	License *License // license information for the Swagger documentation
	Contact *Contact // contact information for the Swagger documentation
}

// buildSwagger creates a new swagger instance with the given title, version, and optional arguments.
func buildSwagger(c Config) (swagger *Swagger) {
	if c.Title == "" {
		c.Title = "Swagger API"
	}
	if c.Version == "" {
		c.Version = "1.0"
	}
	if c.Path == "" {
		c.Path = "/"
	}
	if c.Host == "" {
		c.Host = "localhost"
	}

	swagger = &Swagger{
		Swagger: "2.0",
		Info: Info{
			Title:   c.Title,
			Version: c.Version,
			License: c.License,
			Contact: c.Contact,
		},
		Paths:               make(map[string]map[string]endpoint.JsonEndPoint),
		BasePath:            c.Path,
		Host:                c.Host,
		Definitions:         make(map[string]definition.Definition),
		Schemes:             []string{"http", "https"},
		Tags:                []tag.Tag{},
		SecurityDefinitions: make(map[string]securityDefinition),
		endpoints:           []*endpoint.EndPoint{},
	}

	return
}

// ExportSwaggerDocs exports the Swagger documentation as a JSON file.
func (swagger *Swagger) ExportSwaggerDocs(out_file string) string {
	json, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		log.Println("Error while generating swagger json")
	}
	err = os.WriteFile(out_file, json, 0644)
	if err != nil {
		log.Println("Error writing swagger file")
	}
	return string(json)
}
