package swagno

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/domhoward14/swagno/components/response"
)

func hasStructFields(s interface{}) bool {
	rv := reflect.ValueOf(s)

	if rv.Kind() != reflect.Struct {
		return false
	}

	numFields := rv.NumField()
	return numFields > 0
}

func appendResponses(sourceResponses map[string]jsonResponse, additionalResponses []response.Info) map[string]jsonResponse {
	for _, response := range additionalResponses {
		var responseSchema *jsonResponseScheme
		if hasStructFields(response) {
			responseSchema = &jsonResponseScheme{
				Ref: strings.ReplaceAll(fmt.Sprintf("#/definitions/%T", response), "[]", ""),
			}
		}

		sourceResponses[response.GetReturnCode()] = jsonResponse{
			Description: response.GetDescription(),
			Schema:      responseSchema,
		}
	}

	return sourceResponses
}

func (s *JsonSwagger) generateSwaggerJson() {
	if len(s.endpoints) == 0 {
		log.Println("No endpoints found")
		return
	}

	// generate definition object of swagger json: https://swagger.io/specification/v2/#definitions-object
	s.generateSwaggerDefinition()

	// convert all user EndPoint models to 'path' fields of swagger json
	// https://swagger.io/specification/v2/#paths-object
	for _, endpoint := range s.endpoints {
		path := endpoint.Path

		if s.Paths[path] == nil {
			s.Paths[path] = make(map[string]jsonEndpoint)
		}

		method := strings.ToLower(endpoint.Method)

		for _, param := range endpoint.Params {
			if param.In == "formData" {
				endpoint.Consume = append(endpoint.Consume, "multipart/form-data")
				break
			}
		}

		parameters := make([]jsonParameter, 0)
		for _, param := range endpoint.Params {
			parameters = append(parameters, jsonParameter{
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
			bodyRef := fmt.Sprintf("#/definitions/%T", endpoint.Body)
			bodySchema := jsonResponseScheme{
				Ref: bodyRef,
			}

			if reflect.TypeOf(endpoint.Body).Kind() == reflect.Slice {
				bodySchema = jsonResponseScheme{
					Type: "array",
					Items: &jsonResponseSchemeItems{
						Ref: fmt.Sprintf("#/definitions/%T", endpoint.Body),
					},
				}
			}
			parameters = append(parameters, jsonParameter{
				Name:        "body",
				In:          "body",
				Description: "body",
				Required:    true,
				Schema:      &bodySchema,
			})
		}

		// Creates the schema defintion for all successful return and error objects, and then links them in the responses section
		// TODO make a constructor for this
		responses := map[string]jsonResponse{}
		responses = appendResponses(responses, endpoint.SuccessfulReturns)
		responses = appendResponses(responses, endpoint.Errors)

		// add each endpoint to paths field of swagger
		s.Paths[path][method] = jsonEndpoint{
			Description: endpoint.Description,
			Summary:     endpoint.Summary, // TODO make sure this is working
			OperationId: method + "-" + path,
			Consumes:    endpoint.Consume,
			Produces:    endpoint.Produce,
			Tags:        endpoint.Tags,
			Parameters:  parameters,
			Responses:   responses,
			Security:    endpoint.Security,
		}
	}
}

// Generate swagger v2 documentation as json string
func (s JsonSwagger) GenerateDocs() (jsonDocs []byte) {
	s.generateSwaggerJson()

	json, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Printf("Error while generating swagger json: %s", err)
	}

	return json
}

// generate "definitions" keys from endpoints: https://swagger.io/specification/v2/#definitions-object
func (s *JsonSwagger) generateSwaggerDefinition() {
	for _, endpoint := range s.endpoints {
		if endpoint.Body != nil {
			s.createdefinition(endpoint.Body)
		}
		s.createDefinitions(endpoint.SuccessfulReturns)
		s.createDefinitions(endpoint.Errors)
	}
}

func (s *JsonSwagger) createDefinitions(r []response.Info) {
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

func (s *JsonSwagger) createdefinition(t interface{}) {
	reflectReturn := reflect.TypeOf(t)
	if reflectReturn.Kind() == reflect.Slice {
		reflectReturn = reflectReturn.Elem()
	}
	properties := make(map[string]jsonDefinitionProperties)
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
				// TODO make a constructor function for swaggerDefinitionProperties and create tests for all types to ensure it's extracting the tags correctly
				properties[getJsonTag(field)] = jsonDefinitionProperties{
					Example: getExampleTag(field),
					Type:    fieldType,
					Items: &jsonDefinitionPropertiesItems{
						Ref: fmt.Sprintf("#/definitions/%s", field.Type.Elem().String()),
					},
				}
				s.createdefinition(reflect.New(field.Type.Elem()).Elem().Interface())
			} else {
				properties[getJsonTag(field)] = jsonDefinitionProperties{
					Example: getExampleTag(field),
					Type:    fieldType,
					Items: &jsonDefinitionPropertiesItems{
						Type: getType(field.Type.Elem().Kind().String()),
					},
				}
			}
		} else {
			if field.Type.Kind() == reflect.Struct {
				if field.Type.String() == "time.Time" {
					properties[getJsonTag(field)] = jsonDefinitionProperties{
						Example: getExampleTag(field),
						Type:    "string",
						Format:  "date-time",
					}
				} else if field.Type.String() == "time.Duration" {
					properties[getJsonTag(field)] = jsonDefinitionProperties{
						Example: getExampleTag(field),
						Type:    "integer",
					}
				} else {
					properties[getJsonTag(field)] = jsonDefinitionProperties{
						Example: getExampleTag(field),
						Ref:     fmt.Sprintf("#/definitions/%s", field.Type.String()),
					}
					s.createdefinition(reflect.New(field.Type).Elem().Interface())
				}
			} else if field.Type.Kind() == reflect.Pointer {
				if field.Type.Elem().Kind() == reflect.Struct {
					if field.Type.Elem().String() == "time.Time" {
						properties[getJsonTag(field)] = jsonDefinitionProperties{
							Example: getExampleTag(field),
							Type:    "string",
							Format:  "date-time",
						}
					} else if field.Type.String() == "time.Duration" {
						properties[getJsonTag(field)] = jsonDefinitionProperties{
							Example: getExampleTag(field),
							Type:    "integer",
						}
					} else {
						properties[getJsonTag(field)] = jsonDefinitionProperties{
							Example: getExampleTag(field),
							Ref:     fmt.Sprintf("#/definitions/%s", field.Type.Elem().String()),
						}
						s.createdefinition(reflect.New(field.Type.Elem()).Elem().Interface())
					}
				} else {
					properties[getJsonTag(field)] = jsonDefinitionProperties{
						Example: getExampleTag(field),
						Type:    getType(field.Type.Elem().Kind().String()),
					}
				}
			} else {
				properties[getJsonTag(field)] = jsonDefinitionProperties{
					Example: getExampleTag(field),
					Type:    fieldType,
				}
			}
		}
	}
	(*s).Definitions[fmt.Sprintf("%T", t)] = Definition{
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
