package swagno

import (
	"encoding/json"
	"log"
	"os"
)

// TODO remove globals
var endpoints []Endpoint

// Swagger represents the full JSON model for swagger v2 documentation
// https://swagger.io/docs/specification/2-0/basic-structure/
type Swagger struct {
	Swagger             string                                `json:"swagger" default:"2.0"`
	Info                Info                                  `json:"info"`
	Paths               map[string]map[string]swaggerEndpoint `json:"paths"`
	BasePath            string                                `json:"basePath" default:"/"`
	Host                string                                `json:"host" default:""`
	Definitions         map[string]Definition                 `json:"definitions"`
	Schemes             []string                              `json:"schemes,omitempty"`
	Tags                []Tag                                 `json:"tags,omitempty"`
	SecurityDefinitions map[string]securityDefinition         `json:"securityDefinitions,omitempty"`
}

// Info represents the information about the API.
// https://swagger.io/specification/v2/#info-object
type Info struct {
	Title          string  `json:"title"`
	Version        string  `json:"version"`
	TermsOfService string  `json:"termsOfService,omitempty"`
	Contact        contact `json:"contact,omitempty"`
	License        license `json:"license,omitempty"`
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

// New creates a new swagger instance with the given title, version, and optional arguments.
func New(title string, version string, args ...string) *Swagger {
	return generateSwagger(title, version, args...)
}

// AddEndpoints adds EndPoint models to the Swagger endpoints.
func AddEndpoints(e []Endpoint) {
	endpoints = append(endpoints, e...)
}

// AddEndpoint adds an EndPoint model to the Swagger endpoints.
func AddEndpoint(e Endpoint) {
	endpoints = append(endpoints, e)
}

func (swagger *Swagger) AddTags(tags ...Tag) {
	swagger.Tags = append(swagger.Tags, tags...)
}

// generateSwagger creates a new swagger instance with the given title, version, and optional arguments.
// The optional arguments are title, version, basePath, host.
func generateSwagger(title string, version string, args ...string) (swagger *Swagger) {
	if title == "" {
		title = "Swagger API"
	}
	if version == "" {
		version = "1.0"
	}
	swagger = &Swagger{
		Swagger: "2.0",
		Info: Info{
			Title:   title,
			Version: version,
			License: license{},
			Contact: contact{},
		},
		BasePath:            "/",
		Host:                "",
		Schemes:             []string{"http", "https"},
		SecurityDefinitions: make(map[string]securityDefinition),
	}
	if len(args) > 0 {
		swagger.BasePath = args[0]
		if len(args) > 1 {
			swagger.Host = args[1]
		}
	}
	swagger.Paths = make(map[string]map[string]swaggerEndpoint)
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
