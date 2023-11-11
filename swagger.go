package swagno

import (
	"encoding/json"
	"log"
	"os"
)

// The full JSON model for swagger v2 documentation
// https://swagger.io/docs/specification/2-0/basic-structure/
type Swagger struct {
	Swagger             string                                `json:"swagger" default:"2.0"`
	Info                swaggerInfo                           `json:"info"`
	Paths               map[string]map[string]swaggerEndpoint `json:"paths"`
	BasePath            string                                `json:"basePath" default:"/"`
	Host                string                                `json:"host" default:""`
	Definitions         map[string]swaggerDefinition          `json:"definitions"`
	Schemes             []string                              `json:"schemes,omitempty"`
	Tags                []SwaggerTag                          `json:"tags,omitempty"`
	SecurityDefinitions map[string]swaggerSecurityDefinition  `json:"securityDefinitions,omitempty"`
	endpoints           []Endpoint
}

// https://swagger.io/specification/v2/#info-object
type swaggerInfo struct {
	Title          string          `json:"title"`
	Version        string          `json:"version"`
	TermsOfService string          `json:"termsOfService,omitempty"`
	Contact        *SwaggerContact `json:"contact,omitempty"`
	License        *SwaggerLicense `json:"license,omitempty"`
}

// https://swagger.io/specification/v2/#contact-object
type SwaggerContact struct {
	Name  string `json:"name,omitempty"`
	Url   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// https://swagger.io/specification/v2/#license-object
type SwaggerLicense struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

// https://swagger.io/specification/v2/#securityDefinitionsObject
type swaggerSecurityDefinition struct {
	Type             string            `json:"type"`
	Description      string            `json:"description,omitempty"`
	Name             string            `json:"name,omitempty"`
	In               string            `json:"in,omitempty"`
	Flow             string            `json:"flow,omitempty"`
	AuthorizationUrl string            `json:"authorizationUrl,omitempty"`
	TokenUrl         string            `json:"tokenUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

type Contact struct {
	Name string `json:"name"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Config struct {
	Title   string
	Version string
	Host    string
	Path    string
	License *SwaggerLicense
	Contact *SwaggerContact
}

// Create a new swagger instance
func CreateNewSwagger(c Config) *Swagger {
	return generateSwagger(c)
}

// Add EndPoint models to Swagger endpoints
func (s *Swagger) AddEndpoints(e []Endpoint) {
	s.endpoints = append(s.endpoints, e...)
}

func (s *Swagger) AddEndpoint(e Endpoint) {
	s.endpoints = append(s.endpoints, e)
}

// Create a new swagger instance
// args: title, version, basePath, host
func generateSwagger(c Config) (swagger *Swagger) {
	if c.Title == "" {
		c.Title = "Swagger API"
	}
	if c.Version == "" {
		c.Version = "1.0"
	}
	swagger = &Swagger{
		Swagger: "2.0",
		Info: swaggerInfo{
			Title:   c.Title,
			Version: c.Version,
			License: c.License,
			Contact: c.Contact,
		},
		BasePath:            "/",
		Host:                "",
		Schemes:             []string{"http", "https"},
		SecurityDefinitions: make(map[string]swaggerSecurityDefinition),
	}

	swagger.BasePath = c.Path
	swagger.Host = c.Host
	swagger.Paths = make(map[string]map[string]swaggerEndpoint)
	return
}

// To export json file to an output file
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
