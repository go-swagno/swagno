package swagno

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var lock = &sync.Mutex{}
var swagger *Swagger
var endpoints []Endpoint

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
}

// https://swagger.io/specification/v2/#info-object
type swaggerInfo struct {
	Title          string         `json:"title"`
	Version        string         `json:"version"`
	TermsOfService string         `json:"termsOfService,omitempty"`
	Contact        swaggerContact `json:"contact,omitempty"`
	License        swaggerLicense `json:"license,omitempty"`
}

// https://swagger.io/specification/v2/#contact-object
type swaggerContact struct {
	Name  string `json:"name,omitempty"`
	Url   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// https://swagger.io/specification/v2/#license-object
type swaggerLicense struct {
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

// Create a new swagger instance
func CreateNewSwagger(title string, version string, args ...string) Swagger {
	newSwagger := generateSwagger(title, version, args...)
	swagger = &newSwagger
	return *swagger
}

// returns singleton swagger instance
func GetSwagger() Swagger {
	if swagger == nil {
		lock.Lock()
		if swagger == nil {
			// Creating single instance of swagger
			newSwagger := generateSwagger("", "")
			swagger = &newSwagger
		}
		lock.Unlock()
	}

	return *swagger
}

// Add EndPoint models to Swagger endpoints
func AddEndpoints(e []Endpoint) {
	endpoints = append(endpoints, e...)
}

func AddEndpoint(e Endpoint) {
	endpoints = append(endpoints, e)
}

// Create a new swagger instance
// args: title, version, basePath, host
func generateSwagger(title string, version string, args ...string) (swagger Swagger) {
	if title == "" {
		title = "Swagger API"
	}
	if version == "" {
		version = "1.0"
	}
	swagger = Swagger{
		Swagger: "2.0",
		Info: swaggerInfo{
			Title:   title,
			Version: version,
			License: swaggerLicense{},
			Contact: swaggerContact{},
		},
		BasePath:            "/",
		Host:                "",
		Schemes:             []string{"http", "https"},
		SecurityDefinitions: make(map[string]swaggerSecurityDefinition),
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

// To export json file to an output file
func (s *Swagger) ExportSwaggerDocs(out_file string) string {
	json, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Println("Error while generating s json")
	}
	err = os.WriteFile(out_file, json, 0644)
	if err != nil {
		log.Println("Error writing s file")
	}
	return string(json)
}
