package generator

import (
	"fmt"
	"reflect"

	"github.com/go-swagno/swagno/components/definition"
	"github.com/go-swagno/swagno/http/response"
	"github.com/go-swagno/swagno/utils"
)

type DefinitionGenerator struct {
	Definitions map[string]definition.Definition
}

func NewDefinitionGenerator(definitionMap map[string]definition.Definition) *DefinitionGenerator {
	return &DefinitionGenerator{
		Definitions: definitionMap,
	}
}

func (g DefinitionGenerator) CreateDefinition(t interface{}) {
	properties := make(map[string]definition.DefinitionProperties)
	definitionName := fmt.Sprintf("%T", t)

	reflectReturn := reflect.TypeOf(t)
	switch reflectReturn.Kind() {
	case reflect.Slice:
		reflectReturn = reflectReturn.Elem()
		if reflectReturn.Kind() == reflect.Struct {
			properties = g.createStructDefinitions(reflectReturn)
		}
	case reflect.Struct:
		if reflectReturn == reflect.TypeOf(response.CustomResponseType{}) {
			// if CustomResponseType, use Model struct in it
			g.CreateDefinition(t.(response.CustomResponseType).Model)
		}
		properties = g.createStructDefinitions(reflectReturn)
	case reflect.Map:
		if reflectReturn.Name() == "" {
			definitionName = fmt.Sprintf("%T_%s", t, utils.GetHashOfMap(utils.ConvertInterfaceToMap(t)))
		}
		properties = g.createMapDefinitions(reflect.ValueOf(t))
	}

	g.Definitions[definitionName] = definition.Definition{
		Type:       "object",
		Properties: properties,
	}
}

func (g DefinitionGenerator) createStructDefinitions(_struct reflect.Type) map[string]definition.DefinitionProperties {
	properties := make(map[string]definition.DefinitionProperties)
	for i := 0; i < _struct.NumField(); i++ {
		field := _struct.Field(i)
		fieldType := getType(field.Type.Kind().String())
		fieldJsonTag := getJsonTag(field)

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
				// TODO make a constructor function for swaggerdefinition.DefinitionProperties and create tests for all types to ensure it's extracting the tags correctly
				properties[fieldJsonTag] = definition.DefinitionProperties{
					Example: getExampleTag(field),
					Type:    fieldType,
					Items: &definition.DefinitionPropertiesItems{
						Ref: fmt.Sprintf("#/definitions/%s", field.Type.Elem().String()),
					},
				}
				if _struct == field.Type.Elem() {
					continue // prevent recursion
				}
				g.CreateDefinition(reflect.New(field.Type.Elem()).Elem().Interface())
			} else {
				properties[fieldJsonTag] = definition.DefinitionProperties{
					Example: getExampleTag(field),
					Type:    fieldType,
					Items: &definition.DefinitionPropertiesItems{
						Type: getType(field.Type.Elem().Kind().String()),
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
				properties[fieldJsonTag] = definition.DefinitionProperties{
					Example: getExampleTag(field),
					Type:    getType(field.Type.Elem().Kind().String()),
				}
			}

		case "map":
			name := fmt.Sprintf("%s.%s", _struct.String(), fieldJsonTag)
			mapKeyType := field.Type.Key()
			mapValueType := field.Type.Elem()
			if mapValueType.Kind() == reflect.Ptr {
				mapValueType = mapValueType.Elem()
			}
			properties[fieldJsonTag] = definition.DefinitionProperties{
				Ref: fmt.Sprintf("#/definitions/%s", name),
			}
			if mapValueType.Kind() == reflect.Struct {
				g.Definitions[name] = definition.Definition{
					Type: "object",
					Properties: map[string]definition.DefinitionProperties{
						getType(mapKeyType.String()): {
							Ref: fmt.Sprintf("#/definitions/%s", mapValueType.String()),
						},
					},
				}
			} else {
				g.Definitions[name] = definition.Definition{
					Type: "object",
					Properties: map[string]definition.DefinitionProperties{
						getType(mapKeyType.String()): {
							Example: getExampleTag(field),
							Type:    getType(mapValueType.String()),
						},
					},
				}
			}

		case "interface":
			// TODO: Find a way to get real model of interface{}
			properties[fieldJsonTag] = definition.DefinitionProperties{
				Example: getExampleTag(field),
				Type:    "Ambiguous Type: interface{}",
			}
		default:

			properties[fieldJsonTag] = g.defaultProperty(field)
		}
	}

	return properties
}

func (g DefinitionGenerator) createMapDefinitions(v reflect.Value) map[string]definition.DefinitionProperties {
	properties := make(map[string]definition.DefinitionProperties)

	g.walkMap(v, properties)

	return properties
}

func (g DefinitionGenerator) timeProperty(field reflect.StructField) definition.DefinitionProperties {
	return definition.DefinitionProperties{
		Example: getExampleTag(field),
		Type:    "string",
		Format:  "date-time",
	}
}

func (g DefinitionGenerator) durationProperty(field reflect.StructField) definition.DefinitionProperties {
	return definition.DefinitionProperties{
		Example: getExampleTag(field),
		Type:    "integer",
	}
}

func (g DefinitionGenerator) refProperty(field reflect.StructField) definition.DefinitionProperties {
	return definition.DefinitionProperties{
		Example: getExampleTag(field),
		Ref:     fmt.Sprintf("#/definitions/%s", field.Type.Elem().String()),
	}
}

func (g DefinitionGenerator) defaultProperty(field reflect.StructField) definition.DefinitionProperties {
	return definition.DefinitionProperties{
		Example: getExampleTag(field),
		Type:    getType(field.Type.Kind().String()),
	}
}

func (g DefinitionGenerator) walkMap(v reflect.Value, m map[string]definition.DefinitionProperties) reflect.Value {
	if v.Kind() != reflect.Map {
		return v
	}

	for _, k := range v.MapKeys() {
		val := g.walkMap(v.MapIndex(k), m)
		if val.Type().Kind() == reflect.Struct {
			m[k.String()] = definition.DefinitionProperties{
				Ref: fmt.Sprintf("#/definitions/%s", val.Type().String()),
			}
		} else {
			valueType := val.Type()
			if valueType.Kind() == reflect.Interface {
				valueType = val.Elem().Type()
			}
			m[k.String()] = definition.DefinitionProperties{
				Type: getType(valueType.String()),
			}
		}
	}

	return v
}
