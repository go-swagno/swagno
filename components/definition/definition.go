package definition

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-swagno/swagno/components/fields"
	"github.com/go-swagno/swagno/components/http/response"
)

// Definition represents a Swagger 2.0 schema definition for a type.
// See: https://swagger.io/specification/v2/#definitionsObject
type Definition struct {
	Type       string                          `json:"type"`
	Properties map[string]DefinitionProperties `json:"properties"`
	Required   []string                        `json:"required"`
}

// DefinitionProperties defines the details of a property within a Definition,
// which may include its type, format, reference to another definition, among others.
// See: https://swagger.io/specification/v2/#schemaObject
type DefinitionProperties struct {
	Type    string                     `json:"type,omitempty"`
	Format  string                     `json:"format,omitempty"`
	Ref     string                     `json:"$ref,omitempty"`
	Items   *DefinitionPropertiesItems `json:"items,omitempty"`
	Example interface{}                `json:"example,omitempty"`

	// keep this info to fill Required fields later
	IsRequired bool `json:"-"`
}

// DefinitionPropertiesItems specifies the type or reference of array items when
// the 'type' of DefinitionProperties is set to 'array'.
type DefinitionPropertiesItems struct {
	Type string `json:"type,omitempty"`
	Ref  string `json:"$ref,omitempty"`
}

// DefinitionGenerator holds a map of Definition objects and is capable
// of adding new definitions based on reflected types.
type DefinitionGenerator struct {
	Definitions map[string]Definition
}

// NewDefinitionGenerator is a constructor function that initializes
// a DefinitionGenerator with a provided map of Definition objects.
func NewDefinitionGenerator(definitionMap map[string]Definition) *DefinitionGenerator {
	return &DefinitionGenerator{
		Definitions: definitionMap,
	}
}

// CreateDefinition analyzes the type of the provided value 't' and adds a corresponding Definition to the generator's Definitions map.
func (g DefinitionGenerator) CreateDefinition(t interface{}) {
	properties := make(map[string]DefinitionProperties)
	definitionName := fmt.Sprintf("%T", t)

	if strings.HasPrefix(definitionName, "map[string]") {
		return
	}

	reflectValue := reflect.ValueOf(t)
	switch reflectValue.Kind() {
	case reflect.Slice:

		if reflect.TypeOf(t).Elem().Kind() == reflect.Struct {
			properties = g.createStructDefinitions(reflectValue)
		}
		definitionName, _ = strings.CutPrefix(definitionName, "[]")
	case reflect.Struct:
		if reflectValue.Type() == reflect.TypeOf(response.CustomResponse{}) {
			// if CustomResponseType, use Model struct in it
			g.CreateDefinition(t.(response.CustomResponse).Model)
			return
		}
		properties = g.createStructDefinitions(reflectValue)
	}

	// merge embedded struct fields with other fields
	g.mergeEmbeddedStructFields(properties)

	g.Definitions[definitionName] = Definition{
		Type:       "object",
		Properties: properties,
		Required:   g.findRequiredFields(properties),
	}
}

func (g DefinitionGenerator) mergeEmbeddedStructFields(properties map[string]DefinitionProperties) {
	for k, v := range properties {
		if k == "" && v.Ref != "" { // identify embedded structs
			embeddedModelName, _ := strings.CutPrefix(v.Ref, "#/definitions/")
			if def, ok := g.Definitions[embeddedModelName]; ok {
				for propName, propValue := range def.Properties {
					properties[propName] = propValue
				}
				delete(properties, "")
			}
		}
	}
}

func (g DefinitionGenerator) findRequiredFields(properties map[string]DefinitionProperties) []string {
	requiredFields := []string{}
	for k, v := range properties {
		if v.IsRequired {
			requiredFields = append(requiredFields, k)
		}
	}
	return requiredFields
}

func (g DefinitionGenerator) createStructDefinitions(structValue reflect.Value) map[string]DefinitionProperties {
	if structValue.Kind() == reflect.Slice {
		structValue = reflect.New(structValue.Type().Elem()).Elem()
	}

	properties := make(map[string]DefinitionProperties)
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := field.Type()
		fieldKind := field.Type().Kind()
		fieldStructType := structValue.Type().Field(i)
		fieldJsonTag := fields.JsonTag(fieldStructType)

		// skip ignored tags
		if fieldJsonTag == "-" {
			continue
		}

		// skip for function and channel types
		if fieldKind == reflect.Func || fieldKind == reflect.Chan {
			continue
		}

		// if item type is array, create Definition for array element type
		switch fieldKind {
		case reflect.Array, reflect.Slice:
			if field.Elem().Kind() == reflect.Pointer { // []*T
				if field.Elem().Elem().Kind() == reflect.Struct { // []*struct
					properties[fieldJsonTag] = DefinitionProperties{
						Example: fields.ExampleTag(fieldStructType),
						Type:    "array",
						Items: &DefinitionPropertiesItems{
							Ref: fmt.Sprintf("#/definitions/%s", fieldType.Elem().Elem().String()),
						},
						IsRequired: g.isRequired(fieldStructType),
					}
					if structValue.Type() == fieldType.Elem() {
						continue // prevent recursion
					}
					g.CreateDefinition(reflect.New(fieldType.Elem().Elem()).Elem().Interface())
				} else { // []*primitve_type
					itemType := fields.Type(fieldType.Elem().Elem())
					properties[fieldJsonTag] = DefinitionProperties{
						Example: fields.ExampleTag(fieldStructType),
						Type:    "array",
						Items: &DefinitionPropertiesItems{
							Type: itemType,
						},
						IsRequired: g.isRequired(fieldStructType),
					}
				}
			} else if fieldType.Elem().Kind() == reflect.Struct { // []struct
				properties[fieldJsonTag] = DefinitionProperties{
					Example: fields.ExampleTag(fieldStructType),
					Type:    "array",
					Items: &DefinitionPropertiesItems{
						Ref: fmt.Sprintf("#/definitions/%s", fieldType.Elem().String()),
					},
					IsRequired: g.isRequired(fieldStructType),
				}
				if structValue.Type() == fieldType.Elem() {
					continue // prevent recursion
				}
				g.CreateDefinition(reflect.New(fieldType.Elem()).Elem().Interface())
			} else { // []other
				properties[fieldJsonTag] = DefinitionProperties{
					Example: fields.ExampleTag(fieldStructType),
					Type:    "array",
					Items: &DefinitionPropertiesItems{
						Type: fields.Type(fieldType.Elem()),
					},
					IsRequired: g.isRequired(fieldStructType),
				}
			}

		case reflect.Struct:
			isRequiredField := g.isRequired(fieldStructType)
			if fieldType.String() == "time.Time" {
				properties[fieldJsonTag] = g.timeProperty(fieldStructType, isRequiredField)
			} else if fieldType.String() == "time.Duration" {
				properties[fieldJsonTag] = g.durationProperty(fieldStructType, isRequiredField)
			} else {
				properties[fieldJsonTag] = DefinitionProperties{
					Example:    fields.ExampleTag(fieldStructType),
					Ref:        fmt.Sprintf("#/definitions/%s", fieldType.String()),
					IsRequired: isRequiredField,
				}
				g.CreateDefinition(reflect.New(fieldType).Elem().Interface())
			}

		case reflect.Pointer:
			if fieldType.Elem() == structValue.Type() { // prevent recursion
				properties[fieldJsonTag] = DefinitionProperties{
					Example: fmt.Sprintf("Recursive Type: %s", fieldType.Elem().String()),
				}
				continue
			}
			if fieldType.Elem().Kind() == reflect.Struct { // *struct
				if fieldType.Elem().String() == "time.Time" {
					properties[fieldJsonTag] = g.timeProperty(fieldStructType, false)
				} else if fieldType.String() == "time.Duration" {
					properties[fieldJsonTag] = g.durationProperty(fieldStructType, false)
				} else {
					properties[fieldJsonTag] = g.refProperty(fieldStructType, false)
					g.CreateDefinition(reflect.New(fieldType.Elem()).Elem().Interface())
				}
			} else if fieldType.Elem().Kind() == reflect.Array || fieldType.Elem().Kind() == reflect.Slice { // *[]T
				if fieldType.Elem().Elem().Kind() == reflect.Struct {
					properties[fieldJsonTag] = DefinitionProperties{
						Example: fields.ExampleTag(fieldStructType),
						Type:    fields.Type(fieldType.Elem()),
						Items: &DefinitionPropertiesItems{
							Ref: fmt.Sprintf("#/definitions/%s", fieldType.Elem().Elem().String()),
						},
					}
					if structValue.Type() == fieldType.Elem().Elem() {
						continue // prevent recursion
					}
					g.CreateDefinition(reflect.New(fieldType.Elem().Elem()).Elem().Interface())
				} else {
					properties[fieldJsonTag] = DefinitionProperties{
						Example: fields.ExampleTag(fieldStructType),
						Type:    fields.Type(fieldType.Elem()),
						Items: &DefinitionPropertiesItems{
							Type: fields.Type(fieldType.Elem().Elem()),
						},
					}
				}
			} else {
				properties[fieldJsonTag] = DefinitionProperties{
					Example: fields.ExampleTag(fieldStructType),
					Type:    fields.Type(fieldType.Elem()),
				}
			}

		case reflect.Map:
			name := fmt.Sprintf("%s.%s", structValue.String(), fieldJsonTag)
			mapKeyType := fieldType.Key()
			mapValueType := fieldType.Elem()
			if mapValueType.Kind() == reflect.Ptr {
				mapValueType = mapValueType.Elem()
			}
			properties[fieldJsonTag] = DefinitionProperties{
				Ref: fmt.Sprintf("#/definitions/%s", name),
			}
			if mapValueType.Kind() == reflect.Struct {
				g.Definitions[name] = Definition{
					Type: "object",
					Properties: map[string]DefinitionProperties{
						fields.Type(mapKeyType): {
							Ref: fmt.Sprintf("#/definitions/%s", mapValueType.String()),
						},
					},
				}
			} else {
				g.Definitions[name] = Definition{
					Type: "object",
					Properties: map[string]DefinitionProperties{
						fields.Type(mapKeyType): {
							Example: fields.ExampleTag(fieldStructType),
							Type:    fields.Type(mapValueType),
						},
					},
				}
			}

		case reflect.Interface:
			fek := field.Elem().Kind()
			if fek == 0 {
				continue
			}

			if fek == reflect.Pointer {
				fek = field.Elem().Elem().Kind()
			}

			if fields.IsPrimitiveValue(fek) {
				continue
			}

			if fek == reflect.Slice || fek == reflect.Array {
				var k reflect.Kind
				//[]*T
				if field.Elem().Type().Elem().Kind() == reflect.Pointer {
					k = field.Elem().Type().Elem().Elem().Kind()

				} else { //[]T
					k = field.Elem().Type().Elem().Kind()
				}

				if fields.IsPrimitiveValue(k) {
					continue
				}
			}

			g.CreateDefinition(
				reflect.New(
					field.Elem().Type(),
				).Elem().Interface(),
			)

		default:
			properties[fieldJsonTag] = g.defaultProperty(fieldStructType)

		}
	}

	return properties
}

func (g DefinitionGenerator) timeProperty(field reflect.StructField, required bool) DefinitionProperties {
	return DefinitionProperties{
		Example:    fields.ExampleTag(field),
		Type:       "string",
		Format:     "date-time",
		IsRequired: required,
	}
}

func (g DefinitionGenerator) durationProperty(field reflect.StructField, required bool) DefinitionProperties {
	return DefinitionProperties{
		Example:    fields.ExampleTag(field),
		Type:       "integer",
		IsRequired: required,
	}
}

func (g DefinitionGenerator) refProperty(field reflect.StructField, required bool) DefinitionProperties {
	return DefinitionProperties{
		Example:    fields.ExampleTag(field),
		Ref:        fmt.Sprintf("#/definitions/%s", field.Type.Elem().String()),
		IsRequired: required,
	}
}

func (g DefinitionGenerator) defaultProperty(field reflect.StructField) DefinitionProperties {
	return DefinitionProperties{
		Example:    fields.ExampleTag(field),
		Type:       fields.Type(field.Type),
		IsRequired: g.isRequired(field),
	}
}

func (g DefinitionGenerator) isRequired(field reflect.StructField) bool {
	hasRequiredTag := fields.IsRequired(field)
	hasOmitemptyTag := fields.IsOmitempty(field)
	return hasRequiredTag || !hasOmitemptyTag
}
