package swagno

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-swagno/swagno/components/definition"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/mime"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/go-swagno/swagno/http/response"
)

func hasStructFields(s interface{}) bool {
	rv := reflect.ValueOf(s)

	if rv.Kind() != reflect.Struct {
		return false
	}

	numFields := rv.NumField()
	return numFields > 0
}

func appendResponses(sourceResponses map[string]endpoint.JsonResponse, additionalResponses []response.Info) map[string]endpoint.JsonResponse {
	for _, response := range additionalResponses {
		var responseSchema *parameter.JsonResponseSchema
		if hasStructFields(response) {
			responseSchema = &parameter.JsonResponseSchema{
				Ref: strings.ReplaceAll(fmt.Sprintf("#/definitions/%T", response), "[]", ""),
			}
		}

		sourceResponses[response.GetReturnCode()] = endpoint.JsonResponse{
			Description: response.GetDescription(),
			Schema:      responseSchema,
		}
	}

	return sourceResponses
}

func (s *Swagger) generateSwaggerJson() {
	if len(s.endpoints) == 0 {
		log.Println("No endpoints found")
		return
	}

	// generate definition object of swagger json: https://swagger.io/specification/v2/#definitions-object
	s.generateSwaggerDefinition()

	// convert all user EndPoint models to 'path' fields of swagger json
	// https://swagger.io/specification/v2/#paths-object
	for _, e := range s.endpoints {
		path := e.GetPath()

		if s.Paths[path] == nil {
			s.Paths[path] = make(map[string]endpoint.JsonEndPoint)
		}

		method := strings.ToLower(e.GetMethod())

		for _, param := range e.GetParams() {
			if param.GetLocation() == parameter.Form {
				endpoint.WithConsume([]mime.MIME{mime.MULTIFORM})(e)
				break
			}
		}

		parameters := make([]parameter.JsonParameter, 0)
		for _, param := range e.GetParams() {
			pj := param.AsJson()
			if pj.In != parameter.Query.String() {
				pj.Type = ""
			}
			parameters = append(parameters, param.AsJson())
		}

		if bjp := e.GetBodyJsonParameter(); bjp != nil {
			parameters = append(parameters, *bjp)
		}

		// Creates the schema defintion for all successful return and error objects, and then links them in the responses section
		responses := map[string]endpoint.JsonResponse{}
		responses = appendResponses(responses, e.GetSuccessfulReturns())
		responses = appendResponses(responses, e.GetErrors())

		// add each endpoint to paths field of swagger
		je := e.AsJson()
		je.OperationId = method + "-" + path
		je.Parameters = parameters
		je.Responses = responses
		s.Paths[path][method] = je
	}
}

// Generate swagger v2 documentation as json string
func (s Swagger) GenerateDocs() (jsonDocs []byte) {
	s.generateSwaggerJson()

	json, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Printf("Error while generating swagger json: %s", err)
	}

	return json
}

// generate "definitions" keys from endpoints: https://swagger.io/specification/v2/#definitions-object
func (s *Swagger) generateSwaggerDefinition() {
	for _, endpoint := range s.endpoints {
		if endpoint.Body != nil {
			s.createdefinition(endpoint.Body)
		}
		s.createDefinitions(endpoint.GetSuccessfulReturns())
		s.createDefinitions(endpoint.GetErrors())
	}
}

func (s *Swagger) createDefinitions(r []response.Info) {
	for _, obj := range r {
		s.createdefinition(obj)
	}
}

func getExampleTag(field reflect.StructField) interface{} {
	tagValue := field.Tag.Get("example")

	if tagValue != "" {
		numValue, err := strconv.ParseUint(tagValue, 10, 64)
		if err == nil {
			return numValue
		}
	}

	return tagValue
}

func (s *Swagger) createdefinition(t interface{}) {
	reflectReturn := reflect.TypeOf(t)
	if reflectReturn.Kind() == reflect.Slice {
		reflectReturn = reflectReturn.Elem()
	}
	properties := make(map[string]definition.DefinitionProperties)
	for i := 0; i < reflectReturn.NumField(); i++ {
		field := reflectReturn.Field(i)
		fieldType := getType(field.Type.Kind().String())

		// skip for function and channel types
		if fieldType == "func" || fieldType == "chan" {
			continue
		}

		// if item type is array, create defination for array element type
		if fieldType == "array" {
			if field.Type.Elem().Kind() == reflect.Struct {
				// TODO make a constructor function for swaggerdefinition.DefinitionProperties and create tests for all types to ensure it's extracting the tags correctly
				properties[getJsonTag(field)] = definition.DefinitionProperties{
					Example: getExampleTag(field),
					Type:    fieldType,
					Items: &definition.DefinitionPropertiesItems{
						Ref: fmt.Sprintf("#/definitions/%s", field.Type.Elem().String()),
					},
				}
				s.createdefinition(reflect.New(field.Type.Elem()).Elem().Interface())
			} else {
				properties[getJsonTag(field)] = definition.DefinitionProperties{
					Example: getExampleTag(field),
					Type:    fieldType,
					Items: &definition.DefinitionPropertiesItems{
						Type: getType(field.Type.Elem().Kind().String()),
					},
				}
			}
		} else {
			if field.Type.Kind() == reflect.Struct {
				if field.Type.String() == "time.Time" {
					properties[getJsonTag(field)] = definition.DefinitionProperties{
						Example: getExampleTag(field),
						Type:    "string",
						Format:  "date-time",
					}
				} else if field.Type.String() == "time.Duration" {
					properties[getJsonTag(field)] = definition.DefinitionProperties{
						Example: getExampleTag(field),
						Type:    "integer",
					}
				} else {
					properties[getJsonTag(field)] = definition.DefinitionProperties{
						Example: getExampleTag(field),
						Ref:     fmt.Sprintf("#/definitions/%s", field.Type.String()),
					}
					s.createdefinition(reflect.New(field.Type).Elem().Interface())
				}
			} else if field.Type.Kind() == reflect.Pointer {
				if field.Type.Elem().Kind() == reflect.Struct {
					if field.Type.Elem().String() == "time.Time" {
						properties[getJsonTag(field)] = definition.DefinitionProperties{
							Example: getExampleTag(field),
							Type:    "string",
							Format:  "date-time",
						}
					} else if field.Type.String() == "time.Duration" {
						properties[getJsonTag(field)] = definition.DefinitionProperties{
							Example: getExampleTag(field),
							Type:    "integer",
						}
					} else {
						properties[getJsonTag(field)] = definition.DefinitionProperties{
							Example: getExampleTag(field),
							Ref:     fmt.Sprintf("#/definitions/%s", field.Type.Elem().String()),
						}
						s.createdefinition(reflect.New(field.Type.Elem()).Elem().Interface())
					}
				} else {
					properties[getJsonTag(field)] = definition.DefinitionProperties{
						Example: getExampleTag(field),
						Type:    getType(field.Type.Elem().Kind().String()),
					}
				}
			} else {
				properties[getJsonTag(field)] = definition.DefinitionProperties{
					Example: getExampleTag(field),
					Type:    fieldType,
				}
			}
		}
	}
	(*s).Definitions[fmt.Sprintf("%T", t)] = definition.Definition{
		Type:       "object",
		Properties: properties,
	}
}

// get struct json tag as string of a struct field
func getJsonTag(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if strings.Index(jsonTag, ",") > 0 {
		return strings.Split(jsonTag, ",")[0]
	}
	return jsonTag
}

// get swagger type from reflection type
// https://swagger.io/specification/v2/#data-types
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
