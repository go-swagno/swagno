package swagno

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
)

func (s *Swagger) generateSwaggerObject() {
	if len(s.endpoints) == 0 {
		log.Println("No endpoints found")
		return
	}

	// generate defination object of swagger json: https://swagger.io/specification/v2/#definitions-object
	generateSwaggerDefinition(s, s.endpoints)

	// convert all user EndPoint models to 'path' fields of swagger json
	// https://swagger.io/specification/v2/#paths-object
	for _, endpoint := range s.endpoints {
		path := endpoint.Path

		if s.Paths[path] == nil {
			s.Paths[path] = make(map[string]swaggerEndpoint)
		}

		method := strings.ToLower(endpoint.Method)

		consumes := []string{"application/json"}
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

		var successSchema *swaggerResponseScheme
		if endpoint.Return != nil {
			successSchema = &swaggerResponseScheme{
				Ref: fmt.Sprintf("#/definitions/%T", endpoint.Return),
			}

			definitionString := fmt.Sprintf("#/definitions/%T", endpoint.Return)
			definitionStringValidURI := strings.ReplaceAll(definitionString, "[]", "")

			if reflect.TypeOf(endpoint.Return).Kind() == reflect.Slice {
				successSchema = &swaggerResponseScheme{
					Type: "array",
					Items: &swaggerResponseSchemeItems{
						Ref: definitionStringValidURI,
					},
				}
			}
		}

		var errorSchemas []*swaggerResponseScheme

		errResponses, ok := endpoint.Error.(ErrorResponses)
		if ok {
			for _, errResponseInfo := range errResponses.GetErrors() {
				if endpoint.Error != nil {
					errorSchemas = append(errorSchemas, &swaggerResponseScheme{
						Ref: fmt.Sprintf("#/definitions/%T", errResponseInfo),
					})
				}
			}
		} else {
			if endpoint.Error != nil {
				errorSchemas = []*swaggerResponseScheme{{
					Ref: fmt.Sprintf("#/definitions/%T", endpoint.Error),
				}}
			}
		}

		var description = "OK"
		var successful_rc = "200"
		val, ok := endpoint.Return.(ResponseInfo)
		if ok {
			successful_rc = val.GetReturnCode()
			description = val.GetDescription()
		}

		responses := make(map[string]swaggerResponse)
		if successSchema != nil {
			responses[successful_rc] = swaggerResponse{
				Description: description,
				Schema:      *successSchema,
			}
		}

		if errorSchemas != nil {
			errorInfos, ok := endpoint.Error.(ErrorResponses)
			if ok {
				for i, errResponseInfo := range errorInfos.GetErrors() {
					responses[errResponseInfo.GetReturnCode()] = swaggerResponse{
						Description: errResponseInfo.GetDescription(),
						Schema:      *errorSchemas[i],
					}
				}
			} else {
				responses["404"] = swaggerResponse{
					Description: "Not Found",
					Schema:      *errorSchemas[0],
				}
			}
		}

		// add each endpoint to paths field of swagger
		s.Paths[path][method] = swaggerEndpoint{
			Description: endpoint.Description,
			Summary:     endpoint.Description,
			OperationId: method + "-" + path,
			Consumes:    consumes,
			Produces:    produces,
			Tags:        endpoint.Tags,
			Parameters:  parameters,
			Responses:   responses,
			Security:    endpoint.Security,
		}
	}
}

// Generate swagger v2 documentation as json string
func (swagger Swagger) GenerateDocs() (jsonDocs []byte) {
	// convert Swagger instance to json string and return it
	swagger.generateSwaggerObject()

	json, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		log.Println("Error while generating swagger json")
	}
	return json
}

// generate "definitions" keys from endpoints: https://swagger.io/specification/v2/#definitions-object
func generateSwaggerDefinition(swagger *Swagger, endpoints []Endpoint) {
	// create all definations for each model used in endpoint
	(*swagger).Definitions = make(map[string]swaggerDefinition)
	for _, endpoint := range endpoints {
		if endpoint.Body != nil {
			createdefinition(swagger, endpoint.Body)
		}

		if endpoint.Return != nil {
			createdefinition(swagger, endpoint.Return)
		}

		errorInfos, ok := endpoint.Error.(ErrorResponses)
		if ok {
			for _, errResponseInfo := range errorInfos.GetErrors() {
				createdefinition(swagger, errResponseInfo)
			}
		} else {
			if endpoint.Error != nil {
				createdefinition(swagger, endpoint.Error)
			}
		}
	}
}

// generate "definitions" attribute for swagger json
func createdefinition(swagger *Swagger, t interface{}) {
	reflectReturn := reflect.TypeOf(t)
	if reflectReturn.Kind() == reflect.Slice {
		reflectReturn = reflectReturn.Elem()
	}
	properties := make(map[string]swaggerDefinitionProperties)
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
				properties[getJsonTag(field)] = swaggerDefinitionProperties{
					Type: fieldType,
					Items: &swaggerDefinitionPropertiesItems{
						Ref: fmt.Sprintf("#/definitions/%s", field.Type.Elem().String()),
					},
				}
				createdefinition(swagger, reflect.New(field.Type.Elem()).Elem().Interface())
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
					createdefinition(swagger, reflect.New(field.Type).Elem().Interface())
				}
			} else if field.Type.Kind() == reflect.Pointer {
				if field.Type.Elem().Kind() == reflect.Struct {
					if field.Type.Elem().String() == "time.Time" {
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
							Ref: fmt.Sprintf("#/definitions/%s", field.Type.Elem().String()),
						}
						createdefinition(swagger, reflect.New(field.Type.Elem()).Elem().Interface())
					}
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
