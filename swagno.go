package swagger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
)

// Create a new swagger instance
func CreateSwagger(title string, version string, args ...string) Swagger {
	return generateSwagger(title, version, args...)
}

func (swagger Swagger) GenerateDocs(endpoints []Endpoint) (jsonDocs []byte) {
	generateSwaggerDefinition(&swagger, endpoints)

	for _, endpoint := range endpoints {
		path := endpoint.Path
		for _, param := range endpoint.Params {
			if param.In == "path" {
				path = fmt.Sprintf("%s/{%s}", path, param.Name)
			}
		}
		parameters := make([]swaggerParameter, 0)
		for _, param := range endpoint.Params {
			parameters = append(parameters, swaggerParameter{
				Name:              param.Name,
				In:                param.In,
				Description:       param.Description,
				Required:          param.Required,
				Type:              param.Type,
				Format:            param.Format,
				Items:             param.Items,
				Enum:              param.Enum,
				Default:           param.Default,
				Min:               param.Min,
				Max:               param.Max,
				MinLen:            param.MinLen,
				MaxLen:            param.MaxLen,
				Pattern:           param.Pattern,
				MaxItems:          param.MaxItems,
				MinItems:          param.MinItems,
				UniqueItems:       param.UniqueItems,
				MultipleOf:        param.MultipleOf,
				CollenctionFormat: param.CollenctionFormat,
			})
		}
		if endpoint.Body != nil {
			bodySchema := swaggerResponseScheme{
				Ref: fmt.Sprintf("#/definitions/%T", endpoint.Body),
			}
			if reflect.TypeOf(endpoint.Body).Kind() == reflect.Slice {
				bodySchema = swaggerResponseScheme{
					Type: "array",
					Items: &swaggerResponseSchemeItems{
						Ref: fmt.Sprintf("#/definitions/%T", endpoint.Body),
					},
				}
			}
			parameters = append(parameters, swaggerParameter{
				Name:        "body",
				In:          "body",
				Description: "body",
				Required:    true,
				Schema:      &bodySchema,
			})
		}
		method := strings.ToLower(endpoint.Method)
		if swagger.Paths[path] == nil {
			swagger.Paths[path] = make(map[string]swaggerEndpoint)
		}

		var successSchema *swaggerResponseScheme
		if endpoint.Return != nil {
			successSchema = &swaggerResponseScheme{
				Ref: fmt.Sprintf("#/definitions/%T", endpoint.Return),
			}
			if reflect.TypeOf(endpoint.Return).Kind() == reflect.Slice {
				successSchema = &swaggerResponseScheme{
					Type: "array",
					Items: &swaggerResponseSchemeItems{
						Ref: fmt.Sprintf("#/definitions/%T", endpoint.Return),
					},
				}
			}
		}
		var errorSchema *swaggerResponseScheme
		if endpoint.Error != nil {
			errorSchema = &swaggerResponseScheme{
				Ref: fmt.Sprintf("#/definitions/%T", endpoint.Error),
			}
		}

		consumes := []string{"application/json", "application/xml"}
		produces := []string{"application/json", "application/xml"}
		for _, param := range endpoint.Params {
			if param.In == "formData" {
				consumes = append([]string{"multipart/form-data"}, consumes...)
				break
			}
		}
		if len(endpoint.Consume) > 0 {
			consumes = append(endpoint.Consume, consumes...)
		}
		if len(endpoint.Produce) > 0 {
			produces = append(endpoint.Produce, produces...)
		}

		responses := make(map[string]swaggerResponse)
		if successSchema != nil {
			responses["200"] = swaggerResponse{
				Description: "OK",
				Schema:      *successSchema,
			}
		}
		if errorSchema != nil {
			responses["404"] = swaggerResponse{
				Description: "Not Found",
				Schema:      *errorSchema,
			}
		}
		swagger.Paths[path][method] = swaggerEndpoint{
			Description: endpoint.Description,
			Summary:     endpoint.Description,
			OperationId: method + "-" + path,
			Consumes:    consumes,
			Produces:    produces,
			Tags:        endpoint.Tags,
			Parameters:  parameters,
			Responses:   responses,
		}
	}
	json, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		log.Println("Error while generating swagger json")
	}
	return json
}

// generate "definations" keys from endpoints
func generateSwaggerDefinition(swagger *Swagger, endpoints []Endpoint) {
	(*swagger).Definitions = make(map[string]swaggerDefinition)
	for _, endpoint := range endpoints {
		if endpoint.Body != nil {
			createDefination(swagger, endpoint.Body)
		}
		if endpoint.Return != nil {
			createDefination(swagger, endpoint.Return)
		}
		if endpoint.Error != nil {
			createDefination(swagger, endpoint.Error)
		}
	}
}

// generate "definations" attribute for swagger json
func createDefination(swagger *Swagger, t interface{}) {
	reflectReturn := reflect.TypeOf(t)
	if reflectReturn.Kind() == reflect.Slice {
		reflectReturn = reflectReturn.Elem()
	}
	properties := make(map[string]swaggerDefinitionProperties)
	for i := 0; i < reflectReturn.NumField(); i++ {
		field := reflectReturn.Field(i)
		fieldType := getType(field.Type.Kind().String())
		if fieldType == "func" || fieldType == "chan" {
			continue
		}
		if fieldType == "array" {
			if field.Type.Elem().Kind() == reflect.Struct {
				properties[getJsonTag(field)] = swaggerDefinitionProperties{
					Type: fieldType,
					Items: &swaggerDefinitionPropertiesItems{
						Ref: fmt.Sprintf("#/definitions/%s", field.Type.Elem().String()),
					},
				}
				createDefination(swagger, reflect.New(field.Type.Elem()).Elem().Interface())
			} else {
				properties[getJsonTag(field)] = swaggerDefinitionProperties{
					Type: fieldType,
					Items: &swaggerDefinitionPropertiesItems{
						Type: getType(field.Type.Elem().Kind().String()),
					},
				}
			}
		} else {
			if field.Type.Kind() == reflect.Struct {
				if field.Type.String() == "time.Time" {
					properties[getJsonTag(field)] = swaggerDefinitionProperties{
						Type:   "string",
						Format: "date-time",
					}
				} else if field.Type.String() == "time.Duration" {
					properties[getJsonTag(field)] = swaggerDefinitionProperties{
						Type: "integer",
					}
				} else {
					properties[getJsonTag(field)] = swaggerDefinitionProperties{
						Ref: fmt.Sprintf("#/definitions/%s", field.Type.String()),
					}
					createDefination(swagger, reflect.New(field.Type).Elem().Interface())
				}
			} else if field.Type.Kind() == reflect.Pointer {
				if field.Type.Elem().Kind() == reflect.Struct {
					properties[getJsonTag(field)] = swaggerDefinitionProperties{
						Ref: fmt.Sprintf("#/definitions/%s", field.Type.Elem().String()),
					}
					createDefination(swagger, reflect.New(field.Type.Elem()).Elem().Interface())
				} else {
					properties[getJsonTag(field)] = swaggerDefinitionProperties{
						Type: getType(field.Type.Elem().Kind().String()),
					}
				}
			} else {
				properties[getJsonTag(field)] = swaggerDefinitionProperties{
					Type: fieldType,
				}
			}

		}
	}
	(*swagger).Definitions[fmt.Sprintf("%T", t)] = swaggerDefinition{
		Type:       "object",
		Properties: properties,
	}
}

func getJsonTag(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if strings.Index(jsonTag, ",") > 0 {
		return strings.Split(jsonTag, ",")[0]
	}
	return jsonTag
}

func getType(t string) string {
	if strings.Contains(strings.ToLower(t), "int") {
		return "integer"
	} else if t == "array" || t == "slice" {
		return "array"
	} else if t == "bool" {
		return "boolean"
	} else if t == "float64" || t == "float32" {
		return "number"
	}
	return t
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
		BasePath: "/",
		Host:     "",
		Schemes:  []string{"http", "https"},
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
func (swagger *Swagger) ExportSwaggerDocs(out_file string) string {
	json, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		log.Println("Error while generating swagger json")
	}
	err = ioutil.WriteFile(out_file, json, 0644)
	if err != nil {
		log.Println("Error writing swagger file")
	}
	return string(json)
}

func (swagger *Swagger) AddTags(tags ...SwaggerTag) {
	swagger.Tags = append(swagger.Tags, tags...)
}

func Tag(name string, description string) SwaggerTag {
	return SwaggerTag{
		Name:        name,
		Description: description,
	}
}

/*
* Type definations
 */
type Swagger struct {
	Swagger     string                                `json:"swagger" default:"2.0"`
	Info        swaggerInfo                           `json:"info"`
	Paths       map[string]map[string]swaggerEndpoint `json:"paths"`
	BasePath    string                                `json:"basePath" default:"/"`
	Host        string                                `json:"host" default:""`
	Definitions map[string]swaggerDefinition          `json:"definitions"`
	Schemes     []string                              `json:"schemes,omitempty"`
	Tags        []SwaggerTag                          `json:"tags,omitempty"`
}

type SwaggerTag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type swaggerDefinition struct {
	Type       string                                 `json:"type"`
	Properties map[string]swaggerDefinitionProperties `json:"properties"`
}
type swaggerDefinitionProperties struct {
	Type   string                            `json:"type,omitempty"`
	Format string                            `json:"format,omitempty"`
	Ref    string                            `json:"$ref,omitempty"`
	Items  *swaggerDefinitionPropertiesItems `json:"items,omitempty"`
}
type swaggerDefinitionPropertiesItems struct {
	Type string `json:"type,omitempty"`
	Ref  string `json:"$ref,omitempty"`
}

type swaggerEndpoint struct {
	Description string                     `json:"description"`
	Consumes    []string                   `json:"consumes" default:"application/json"`
	Produces    []string                   `json:"produces" default:"application/json"`
	Tags        []string                   `json:"tags"`
	Summary     string                     `json:"summary"`
	OperationId string                     `json:"operationId,omitempty"`
	Parameters  []swaggerParameter         `json:"parameters"`
	Responses   map[string]swaggerResponse `json:"responses"`
}

type swaggerParameter struct {
	Type              string                 `json:"type"`
	Description       string                 `json:"description"`
	Name              string                 `json:"name"`
	In                string                 `json:"in"`
	Required          bool                   `json:"required"`
	Schema            *swaggerResponseScheme `json:"schema,omitempty"`
	Format            string                 `json:"format,omitempty"`
	Items             *ParameterItems        `json:"items,omitempty"`
	Enum              []interface{}          `json:"enum,omitempty"`
	Default           interface{}            `json:"default,omitempty"`
	Min               int64                  `json:"minimum,omitempty"`
	Max               int64                  `json:"maximum,omitempty"`
	MinLen            int64                  `json:"minLength,omitempty"`
	MaxLen            int64                  `json:"maxLength,omitempty"`
	Pattern           string                 `json:"pattern,omitempty"`
	MaxItems          int64                  `json:"maxItems,omitempty"`
	MinItems          int64                  `json:"minItems,omitempty"`
	UniqueItems       bool                   `json:"uniqueItems,omitempty"`
	MultipleOf        int64                  `json:"multipleOf,omitempty"`
	CollenctionFormat string                 `json:"collectionFormat,omitempty"`
}
type swaggerResponse struct {
	Description string                `json:"description"`
	Schema      swaggerResponseScheme `json:"schema"`
}
type swaggerResponseScheme struct {
	Ref   string                      `json:"$ref,omitempty"`
	Type  string                      `json:"type,omitempty"`
	Items *swaggerResponseSchemeItems `json:"items,omitempty"`
}
type swaggerResponseSchemeItems struct {
	Type string `json:"type,omitempty"`
	Ref  string `json:"$ref,omitempty"`
}

type swaggerInfo struct {
	Title          string         `json:"title"`
	Version        string         `json:"version"`
	TermsOfService string         `json:"termsOfService,omitempty"`
	Contact        swaggerContact `json:"contact,omitempty"`
	License        swaggerLicense `json:"license,omitempty"`
}

type swaggerContact struct {
	Name  string `json:"name,omitempty"`
	Url   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type swaggerLicense struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}
