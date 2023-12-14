package definition

import (
	"fmt"
	"reflect"
)

// Definition represents a Swagger 2.0 schema definition for a type.
// See: https://swagger.io/specification/v2/#definitionsObject
type Definition struct {
	Type       string                          `json:"type"`
	Properties map[string]DefinitionProperties `json:"properties"`
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

// CreateDefinition analyzes the type of the provided value 't'
// and adds a corresponding Definition to the generator's Definitions map.
func (g DefinitionGenerator) CreateDefinition(t interface{}) {
	properties := make(map[string]DefinitionProperties)
	definitionName := fmt.Sprintf("%T", t)

	reflectReturn := reflect.TypeOf(t)
	switch reflectReturn.Kind() {
	case reflect.Slice:
		reflectReturn = reflectReturn.Elem()
		if reflectReturn.Kind() == reflect.Struct {
			properties = g.createStructDefinitions(reflectReturn)
		}
	case reflect.Struct:
		properties = g.createStructDefinitions(reflectReturn)
	}

	g.Definitions[definitionName] = Definition{
		Type:       "object",
		Properties: properties,
	}
}

func (g DefinitionGenerator) createStructDefinitions(_struct reflect.Type) map[string]DefinitionProperties {
	properties := make(map[string]DefinitionProperties)
	for i := 0; i < _struct.NumField(); i++ {
		field := _struct.Field(i)
		fieldType := Type(field.Type.Kind().String())
		fieldJsonTag := JsonTag(field)

		// skip ignored tags
		if fieldJsonTag == "-" {
			continue
		}

		// skip for function and channel types
		if fieldType == "func" || fieldType == "chan" {
			continue
		}

		// if item type is array, create Definition for array element type
		switch fieldType {
		case "array":
			if field.Type.Elem().Kind() == reflect.Struct {
				// TODO make a constructor function for swaggerDefinitionProperties and create tests for all types to ensure it's extracting the tags correctly
				properties[fieldJsonTag] = DefinitionProperties{
					Example: ExampleTag(field),
					Type:    fieldType,
					Items: &DefinitionPropertiesItems{
						Ref: fmt.Sprintf("#/definitions/%s", field.Type.Elem().String()),
					},
				}
				if _struct == field.Type.Elem() {
					continue // prevent recursion
				}
				g.CreateDefinition(reflect.New(field.Type.Elem()).Elem().Interface())
			} else {
				properties[fieldJsonTag] = DefinitionProperties{
					Example: ExampleTag(field),
					Type:    fieldType,
					Items: &DefinitionPropertiesItems{
						Type: Type(field.Type.Elem().Kind().String()),
					},
				}
			}

		case "struct":
			if field.Type.String() == "time.Time" {
				properties[fieldJsonTag] = g.timeProperty(field)
			} else if field.Type.String() == "time.Duration" {
				properties[fieldJsonTag] = g.durationProperty(field)
			} else {
				properties[fieldJsonTag] = g.refProperty(field)
				g.CreateDefinition(reflect.New(field.Type).Elem().Interface())
			}

		case "ptr":
			if field.Type.Elem().Kind() == reflect.Struct {
				if field.Type.Elem().String() == "time.Time" {
					properties[fieldJsonTag] = g.timeProperty(field)
				} else if field.Type.String() == "time.Duration" {
					properties[fieldJsonTag] = g.durationProperty(field)
				} else {
					properties[fieldJsonTag] = g.refProperty(field)
					g.CreateDefinition(reflect.New(field.Type.Elem()).Elem().Interface())
				}
			} else {
				properties[fieldJsonTag] = DefinitionProperties{
					Example: ExampleTag(field),
					Type:    Type(field.Type.Elem().Kind().String()),
				}
			}

		case "map":
			name := fmt.Sprintf("%s.%s", _struct.String(), fieldJsonTag)
			mapKeyType := field.Type.Key()
			mapValueType := field.Type.Elem()
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
						Type(mapKeyType.String()): {
							Ref: fmt.Sprintf("#/definitions/%s", mapValueType.String()),
						},
					},
				}
			} else {
				g.Definitions[name] = Definition{
					Type: "object",
					Properties: map[string]DefinitionProperties{
						Type(mapKeyType.String()): {
							Example: ExampleTag(field),
							Type:    Type(mapValueType.String()),
						},
					},
				}
			}

		case "interface":
			// TODO: Find a way to get real model of interface{}
			properties[fieldJsonTag] = DefinitionProperties{
				Example: ExampleTag(field),
				Type:    "Ambiguous Type: interface{}",
			}
		default:

			properties[fieldJsonTag] = g.defaultProperty(field)
		}
	}

	return properties
}

func (g DefinitionGenerator) timeProperty(field reflect.StructField) DefinitionProperties {
	return DefinitionProperties{
		Example: ExampleTag(field),
		Type:    "string",
		Format:  "date-time",
	}
}

func (g DefinitionGenerator) durationProperty(field reflect.StructField) DefinitionProperties {
	return DefinitionProperties{
		Example: ExampleTag(field),
		Type:    "integer",
	}
}

func (g DefinitionGenerator) refProperty(field reflect.StructField) DefinitionProperties {
	return DefinitionProperties{
		Example: ExampleTag(field),
		Ref:     fmt.Sprintf("#/definitions/%s", field.Type.Elem().String()),
	}
}

func (g DefinitionGenerator) defaultProperty(field reflect.StructField) DefinitionProperties {
	return DefinitionProperties{
		Example: ExampleTag(field),
		Type:    Type(field.Type.Kind().String()),
	}
}
