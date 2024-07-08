package response

import (
	"fmt"
	"reflect"

	"strings"

	"github.com/go-swagno/swagno/components/fields"
	"github.com/go-swagno/swagno/components/parameter"
)

// Response is an interface for response information.
type Response interface {
	Description() string
	ReturnCode() string
}

// CustomResponse represents a custom implementation of Response.
type CustomResponse struct {
	Model             any
	returnCodeString  string
	descriptionString string
}

// ResponseGenerator is a struct that provides functionality to generate response schemas.
type ResponseGenerator struct{}

// New creates a new instance of Response with the provided model return code, and description.
func New(model any, returnCode string, description string) CustomResponse {
	return CustomResponse{
		Model:             model,
		returnCodeString:  returnCode,
		descriptionString: description,
	}
}

// NewResponseGenerator creates a new instance of ResponseGenerator.
func NewResponseGenerator() *ResponseGenerator {
	return &ResponseGenerator{}
}

// Generate generates a JSON response schema based on the provided model.
// It uses reflection to determine the type of the model and constructs the appropriate JSON schema.
// This function handles different types such as slices, maps, and structures to create a detailed and accurate schema.
func (g ResponseGenerator) Generate(model any) *parameter.JsonResponseSchemaAllOf {
	switch reflect.TypeOf(model).Kind() {
	case reflect.Slice, reflect.Array:
		sliceElementKind := reflect.TypeOf(model).Elem().Kind()
		if sliceElementKind == reflect.Struct {
			return &parameter.JsonResponseSchemaAllOf{
				AllOf: []parameter.JsonResponseSchema{
					{
						Type: "array",
						Items: &parameter.JsonResponseSchemeItems{
							Ref: strings.ReplaceAll(fmt.Sprintf("#/definitions/%s", reflect.TypeOf(model).Elem().String()), "[]", ""),
						},
					},
				},
			}
		} else {
			return &parameter.JsonResponseSchemaAllOf{
				AllOf: []parameter.JsonResponseSchema{
					{
						Type: "array",
						Items: &parameter.JsonResponseSchemeItems{
							Type: fields.Type(reflect.TypeOf(model).Elem()),
						},
					},
				},
			}
		}

	case reflect.Map:
		modelType := reflect.TypeOf(model)
		if modelType.Key().Kind() == reflect.String {
			//for loop to iterate over the map and get the key and value types
			properties := make(map[string]parameter.JsonResponseSchema)
			modelValueOf := reflect.ValueOf(model)
			for _, key := range modelValueOf.MapKeys() {
				// mapIndex

				valueOfMap := modelValueOf.MapIndex(key).Elem()
				// check if value is primitive type
				var k reflect.Kind
				var v reflect.Value
				if valueOfMap.Type().Kind() == reflect.Pointer {
					k = valueOfMap.Type().Elem().Kind()
					v = valueOfMap.Elem()
				} else {
					k = valueOfMap.Type().Kind()
					v = valueOfMap
				}

				if fields.IsPrimitiveValue(k) {
					properties[key.String()] = parameter.JsonResponseSchema{
						Type:    fields.Type(valueOfMap.Type()),
						Example: v.Interface(),
					}
					continue
				} else {
					model := g.Generate(valueOfMap.Interface())
					if model != nil {
						properties[key.String()] = model.AllOf[0]

					}
				}

			}
			return &parameter.JsonResponseSchemaAllOf{
				AllOf: []parameter.JsonResponseSchema{
					{
						Type:       "object",
						Properties: properties,
					},
				},
			}
		}
		ref := strings.ReplaceAll(fmt.Sprintf("#/definitions/%T", model), "[]", "")
		return &parameter.JsonResponseSchemaAllOf{
			AllOf: []parameter.JsonResponseSchema{
				{
					Ref: ref,
				},
			},
		}

	default:
		if rv, ok := extractStructFields(model); ok {
			parameters := make([]parameter.JsonResponseSchema, 0, rv.NumField())
			for i := 0; i < rv.NumField(); i++ {
				field := rv.Field(i)
				fieldStructType := rv.Type().Field(i)

				if field.Kind() != reflect.Interface {
					continue
				}

				field = field.Elem()

				if field.Kind() == 0 {
					continue
				}

				if param := getValueFromStruct(field, fieldStructType); param != nil {
					parameters = append(parameters, *param)
				}
			}
			return &parameter.JsonResponseSchemaAllOf{
				AllOf: append(
					[]parameter.JsonResponseSchema{
						{
							Ref: strings.ReplaceAll(fmt.Sprintf("#/definitions/%T", model), "[]", ""),
						},
					},
					parameters...,
				),
			}
		}
	}

	return nil
}

func getValueFromStruct(field reflect.Value, structField reflect.StructField) *parameter.JsonResponseSchema {
	fieldJsonTag := fields.JsonTag(structField)
	if fieldJsonTag == "-" {
		return nil
	}

	switch field.Kind() {
	case reflect.Pointer:
		return getValueFromStruct(field.Elem(), structField)
	case reflect.Struct:
		if field.Type().String() == "time.Time" {
			return &parameter.JsonResponseSchema{
				Type: "object",
				Properties: map[string]parameter.JsonResponseSchema{
					fieldJsonTag: {
						Type: "string",
					},
				},
			}
		}
		if field.Type().String() == "time.Duration" {
			return &parameter.JsonResponseSchema{
				Type: "object",
				Properties: map[string]parameter.JsonResponseSchema{
					fieldJsonTag: {
						Type: "integer",
					},
				},
			}
		}

		return &parameter.JsonResponseSchema{
			Type: "object",
			Properties: map[string]parameter.JsonResponseSchema{
				fieldJsonTag: {
					Ref: fmt.Sprintf("#/definitions/%s", field.Type().String()),
				},
			},
		}
	case reflect.Slice, reflect.Array:
		var jsonSchemaItems *parameter.JsonResponseSchemeItems

		var v reflect.Value
		if field.Type().Elem().Kind() == reflect.Pointer {
			v = reflect.New(field.Type().Elem().Elem()).Elem()
		} else {
			v = reflect.New(field.Type().Elem()).Elem()
		}

		primitiveValue := getValueFromStruct(v, structField)
		if primitiveValue == nil {
			return nil
		}

		if props := primitiveValue.Properties[fieldJsonTag]; props.Ref != "" {
			jsonSchemaItems = &parameter.JsonResponseSchemeItems{
				Ref: props.Ref,
			}
		} else {
			jsonSchemaItems = &parameter.JsonResponseSchemeItems{
				Type:  props.Type,
				Items: props.Items,
			}
		}

		return &parameter.JsonResponseSchema{
			Type: "object",
			Properties: map[string]parameter.JsonResponseSchema{
				fieldJsonTag: {
					Type:  "array",
					Items: jsonSchemaItems,
				},
			},
		}
	case reflect.Map:
		properties := make(map[string]parameter.JsonResponseSchema)
		for _, key := range field.MapKeys() {
			// mapIndex

			valueOfMap := field.MapIndex(key).Elem()
			// check if value is primitive type
			if fields.IsPrimitiveValue(valueOfMap.Type().Kind()) {
				properties[key.String()] = parameter.JsonResponseSchema{
					Type:    fields.Type(valueOfMap.Type()),
					Example: valueOfMap.Interface(),
				}
				continue
			} else {
				model := getValueFromStruct(valueOfMap, structField)
				if model != nil {
					properties[key.String()] = *model
				}
			}

		}
		return &parameter.JsonResponseSchema{
			Type: "object",
			Properties: map[string]parameter.JsonResponseSchema{
				fieldJsonTag: {
					Type:       "object",
					Properties: properties,
				},
			},
		}

	default: // primitive types
		return &parameter.JsonResponseSchema{
			Type: "object",
			Properties: map[string]parameter.JsonResponseSchema{
				fieldJsonTag: {
					Type:    fields.Type(field.Type()),
					Example: fields.ExampleTag(structField),
				},
			},
		}
	}

}

// extractStructFields checks if the given interface has fields in case it's a struct.
func extractStructFields(s interface{}) (reflect.Value, bool) {
	rv := reflect.ValueOf(s)

	if rv.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}

	numFields := rv.NumField()
	return rv, numFields > 0
}

func (c CustomResponse) Description() string {
	return c.descriptionString
}
func (c CustomResponse) ReturnCode() string {
	return c.returnCodeString
}
