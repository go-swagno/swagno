package definition

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-swagno/swagno/v3/components/fields"
	"github.com/go-swagno/swagno/v3/components/http/response"
)

// Schema represents an OpenAPI 3.0 schema definition for a type.
// See: https://spec.openapis.org/oas/v3.0.3#schema-object
type Schema struct {
	Type                 string                    `json:"type,omitempty"`
	Format               string                    `json:"format,omitempty"`
	Title                string                    `json:"title,omitempty"`
	Description          string                    `json:"description,omitempty"`
	Default              interface{}               `json:"default,omitempty"`
	Example              interface{}               `json:"example,omitempty"`
	Examples             []interface{}             `json:"examples,omitempty"`
	Enum                 []interface{}             `json:"enum,omitempty"`
	Const                interface{}               `json:"const,omitempty"`
	Properties           map[string]SchemaProperty `json:"properties,omitempty"`
	AdditionalProperties interface{}               `json:"additionalProperties,omitempty"`
	Required             []string                  `json:"required,omitempty"`
	Items                *SchemaItems              `json:"items,omitempty"`
	AllOf                []*Schema                 `json:"allOf,omitempty"`
	OneOf                []*Schema                 `json:"oneOf,omitempty"`
	AnyOf                []*Schema                 `json:"anyOf,omitempty"`
	Not                  *Schema                   `json:"not,omitempty"`
	MinLength            *int64                    `json:"minLength,omitempty"`
	MaxLength            *int64                    `json:"maxLength,omitempty"`
	Pattern              string                    `json:"pattern,omitempty"`
	Minimum              *float64                  `json:"minimum,omitempty"`
	Maximum              *float64                  `json:"maximum,omitempty"`
	ExclusiveMinimum     bool                      `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum     bool                      `json:"exclusiveMaximum,omitempty"`
	MultipleOf           *float64                  `json:"multipleOf,omitempty"`
	MinItems             *int64                    `json:"minItems,omitempty"`
	MaxItems             *int64                    `json:"maxItems,omitempty"`
	UniqueItems          bool                      `json:"uniqueItems,omitempty"`
	MinProperties        *int64                    `json:"minProperties,omitempty"`
	MaxProperties        *int64                    `json:"maxProperties,omitempty"`
	Nullable             bool                      `json:"nullable,omitempty"`
	Discriminator        *Discriminator            `json:"discriminator,omitempty"`
	ReadOnly             bool                      `json:"readOnly,omitempty"`
	WriteOnly            bool                      `json:"writeOnly,omitempty"`
	XML                  *XML                      `json:"xml,omitempty"`
	ExternalDocs         *ExternalDocs             `json:"externalDocs,omitempty"`
	Deprecated           bool                      `json:"deprecated,omitempty"`
	Ref                  string                    `json:"$ref,omitempty"`
}

// SchemaProperty defines the details of a property within a Schema,
// which may include its type, format, reference to another schema, among others.
type SchemaProperty struct {
	Type        string                    `json:"type,omitempty"`
	Format      string                    `json:"format,omitempty"`
	Ref         string                    `json:"$ref,omitempty"`
	Items       *SchemaItems              `json:"items,omitempty"`
	Example     interface{}               `json:"example,omitempty"`
	Description string                    `json:"description,omitempty"`
	Default     interface{}               `json:"default,omitempty"`
	Enum        []interface{}             `json:"enum,omitempty"`
	MinLength   *int64                    `json:"minLength,omitempty"`
	MaxLength   *int64                    `json:"maxLength,omitempty"`
	Pattern     string                    `json:"pattern,omitempty"`
	Minimum     *float64                  `json:"minimum,omitempty"`
	Maximum     *float64                  `json:"maximum,omitempty"`
	MinItems    *int64                    `json:"minItems,omitempty"`
	MaxItems    *int64                    `json:"maxItems,omitempty"`
	UniqueItems bool                      `json:"uniqueItems,omitempty"`
	MultipleOf  *float64                  `json:"multipleOf,omitempty"`
	Nullable    bool                      `json:"nullable,omitempty"`
	ReadOnly    bool                      `json:"readOnly,omitempty"`
	WriteOnly   bool                      `json:"writeOnly,omitempty"`
	Deprecated  bool                      `json:"deprecated,omitempty"`
	Properties  map[string]SchemaProperty `json:"properties,omitempty"`

	// keep this info to fill Required fields later
	IsRequired bool `json:"-"`
}

// SchemaItems specifies the type or reference of array items when
// the 'type' of SchemaProperty is set to 'array'.
type SchemaItems struct {
	Type   string       `json:"type,omitempty"`
	Ref    string       `json:"$ref,omitempty"`
	Items  *SchemaItems `json:"items,omitempty"`
	Format string       `json:"format,omitempty"`
}

// Discriminator represents a discriminator object for polymorphism
type Discriminator struct {
	PropertyName string            `json:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty"`
}

// XML represents XML metadata
type XML struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Attribute bool   `json:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty"`
}

// ExternalDocs represents external documentation
type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
}

// DefinitionGenerator holds a map of Schema objects and is capable
// of adding new schemas based on reflected types.
type DefinitionGenerator struct {
	Schemas map[string]Schema
}

// NewDefinitionGenerator is a constructor function that initializes
// a DefinitionGenerator with a provided map of Schema objects.
func NewDefinitionGenerator(schemaMap map[string]Schema) *DefinitionGenerator {
	return &DefinitionGenerator{
		Schemas: schemaMap,
	}
}

// CreateDefinition analyzes the type of the provided value 't' and adds a corresponding Schema to the generator's Schemas map.
func (g DefinitionGenerator) CreateDefinition(t interface{}) {
	properties := make(map[string]SchemaProperty)
	definitionName := fmt.Sprintf("%T", t)

	reflectReturn := reflect.TypeOf(t)
	switch reflectReturn.Kind() {
	case reflect.Slice:
		reflectReturn = reflectReturn.Elem()
		if reflectReturn.Kind() == reflect.Struct {
			properties = g.createStructDefinitions(reflectReturn)
		}
		definitionName, _ = strings.CutPrefix(definitionName, "[]")
	case reflect.Struct:
		if reflectReturn == reflect.TypeOf(response.CustomResponse{}) {
			// if CustomResponseType, use Model struct in it
			g.CreateDefinition(t.(response.CustomResponse).Model)
			return
		}
		properties = g.createStructDefinitions(reflectReturn)
	}

	// merge embedded struct fields with other fields
	g.mergeEmbeddedStructFields(properties)

	// delete empty json tags
	delete(properties, "")

	g.Schemas[definitionName] = Schema{
		Type:       "object",
		Properties: properties,
		Required:   g.findRequiredFields(properties),
	}
}

func (g DefinitionGenerator) mergeEmbeddedStructFields(properties map[string]SchemaProperty) {
	for k, v := range properties {
		if k == "" && v.Ref != "" { // identify embedded structs
			embeddedModelName, _ := strings.CutPrefix(v.Ref, "#/components/schemas/")
			if schema, ok := g.Schemas[embeddedModelName]; ok {
				for propName, propValue := range schema.Properties {
					properties[propName] = propValue
				}
				delete(properties, "")
			}
		}
	}
}

func (g DefinitionGenerator) findRequiredFields(properties map[string]SchemaProperty) []string {
	requiredFields := []string{}
	for k, v := range properties {
		if v.IsRequired {
			requiredFields = append(requiredFields, k)
		}
	}
	return requiredFields
}

func (g DefinitionGenerator) createStructDefinitions(structType reflect.Type) map[string]SchemaProperty {
	properties := make(map[string]SchemaProperty)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldType := fields.Type(field.Type.Kind().String())
		fieldJsonTag := fields.JsonTag(field)

		// skip ignored tags
		if fieldJsonTag == "-" {
			continue
		}

		// skip for function and channel types
		if fieldType == "func" || fieldType == "chan" {
			continue
		}

		// if item type is array, create Schema for array element type
		switch fieldType {
		case "array":
			if field.Type.Elem().Kind() == reflect.Pointer { // []*type
				if field.Type.Elem().Elem().Kind() == reflect.Struct { // []*struct
					properties[fieldJsonTag] = SchemaProperty{
						Type: fieldType,
						Items: &SchemaItems{
							Ref: fmt.Sprintf("#/components/schemas/%s", field.Type.Elem().Elem().String()),
						},
						IsRequired:  g.isRequired(field),
						Nullable:    field.Type.Kind() == reflect.Pointer,
						Example:     fields.ExampleTag(field),
						Description: fields.DescriptionTag(field),
					}
					if structType == field.Type.Elem().Elem() {
						continue // prevent recursion
					}
					g.CreateDefinition(reflect.New(field.Type.Elem().Elem()).Elem().Interface())
				} else { // []*other
					itemType := fields.Type(field.Type.Elem().Elem().Kind().String())
					properties[fieldJsonTag] = SchemaProperty{
						Type: fieldType,
						Items: &SchemaItems{
							Type: itemType,
						},
						IsRequired:  g.isRequired(field),
						Nullable:    field.Type.Kind() == reflect.Pointer,
						Example:     fields.ExampleTag(field),
						Description: fields.DescriptionTag(field),
					}
				}
			} else if field.Type.Elem().Kind() == reflect.Struct { // []struct
				properties[fieldJsonTag] = SchemaProperty{
					Type: fieldType,
					Items: &SchemaItems{
						Ref: fmt.Sprintf("#/components/schemas/%s", field.Type.Elem().String()),
					},
					IsRequired:  g.isRequired(field),
					Example:     fields.ExampleTag(field),
					Description: fields.DescriptionTag(field),
				}
				if structType == field.Type.Elem() {
					continue // prevent recursion
				}
				g.CreateDefinition(reflect.New(field.Type.Elem()).Elem().Interface())
			} else { // []other
				properties[fieldJsonTag] = SchemaProperty{
					Type: fieldType,
					Items: &SchemaItems{
						Type: fields.Type(field.Type.Elem().Kind().String()),
					},
					IsRequired:  g.isRequired(field),
					Example:     fields.ExampleTag(field),
					Description: fields.DescriptionTag(field),
				}
			}

		case "struct":
			isRequiredField := g.isRequired(field)
			if field.Type.String() == "time.Time" {
				properties[fieldJsonTag] = g.timeProperty(field, isRequiredField)
			} else if field.Type.String() == "time.Duration" {
				properties[fieldJsonTag] = g.durationProperty(field, isRequiredField)
			} else {
				properties[fieldJsonTag] = SchemaProperty{
					Ref:         fmt.Sprintf("#/components/schemas/%s", field.Type.String()),
					IsRequired:  isRequiredField,
					Example:     fields.ExampleTag(field),
					Description: fields.DescriptionTag(field),
				}
				g.CreateDefinition(reflect.New(field.Type).Elem().Interface())
			}

		case "ptr":
			if field.Type.Elem() == structType { // prevent recursion
				properties[fieldJsonTag] = SchemaProperty{
					Example:     fmt.Sprintf("Recursive Type: %s", field.Type.Elem().String()),
					Description: fields.DescriptionTag(field),
				}
				continue
			}
			if field.Type.Elem().Kind() == reflect.Struct {
				if field.Type.Elem().String() == "time.Time" {
					properties[fieldJsonTag] = g.timeProperty(field, false)
				} else if field.Type.String() == "time.Duration" {
					properties[fieldJsonTag] = g.durationProperty(field, false)
				} else {
					properties[fieldJsonTag] = g.refProperty(field, false)
					g.CreateDefinition(reflect.New(field.Type.Elem()).Elem().Interface())
				}
			} else if field.Type.Elem().Kind() == reflect.Array || field.Type.Elem().Kind() == reflect.Slice {
				if field.Type.Elem().Elem().Kind() == reflect.Struct {
					properties[fieldJsonTag] = SchemaProperty{
						Type: fields.Type(field.Type.Elem().Kind().String()),
						Items: &SchemaItems{
							Ref: fmt.Sprintf("#/components/schemas/%s", field.Type.Elem().Elem().String()),
						},
						Nullable:    true,
						Example:     fields.ExampleTag(field),
						Description: fields.DescriptionTag(field),
					}
					if structType == field.Type.Elem().Elem() {
						continue // prevent recursion
					}
					g.CreateDefinition(reflect.New(field.Type.Elem().Elem()).Elem().Interface())
				} else {
					properties[fieldJsonTag] = SchemaProperty{
						Type: fields.Type(field.Type.Elem().Kind().String()),
						Items: &SchemaItems{
							Type: fields.Type(field.Type.Elem().Elem().Kind().String()),
						},
						Nullable:    true,
						Example:     fields.ExampleTag(field),
						Description: fields.DescriptionTag(field),
					}
				}
			} else {
				properties[fieldJsonTag] = SchemaProperty{
					Type:        fields.Type(field.Type.Elem().Kind().String()),
					Nullable:    true,
					Example:     fields.ExampleTag(field),
					Description: fields.DescriptionTag(field),
				}
			}

		case "map":
			name := fmt.Sprintf("%s.%s", structType.String(), fieldJsonTag)
			mapKeyType := field.Type.Key()
			mapValueType := field.Type.Elem()
			if mapValueType.Kind() == reflect.Ptr {
				mapValueType = mapValueType.Elem()
			}
			properties[fieldJsonTag] = SchemaProperty{
				Ref:         fmt.Sprintf("#/components/schemas/%s", name),
				Example:     fields.ExampleTag(field),
				Description: fields.DescriptionTag(field),
			}
			if mapValueType.Kind() == reflect.Struct {
				g.Schemas[name] = Schema{
					Type: "object",
					Properties: map[string]SchemaProperty{
						fields.Type(mapKeyType.String()): {
							Ref:         fmt.Sprintf("#/components/schemas/%s", mapValueType.String()),
							Example:     fields.ExampleTag(field),
							Description: fields.DescriptionTag(field),
						},
					},
				}
			} else {
				g.Schemas[name] = Schema{
					Type: "object",
					Properties: map[string]SchemaProperty{
						fields.Type(mapKeyType.String()): {
							Type:        fields.Type(mapValueType.String()),
							Example:     fields.ExampleTag(field),
							Description: fields.DescriptionTag(field),
						},
					},
				}
			}

		case "interface":
			// TODO: Find a way to get real model of interface{}
			properties[fieldJsonTag] = SchemaProperty{
				Type:        "Ambiguous Type: interface{}",
				IsRequired:  g.isRequired(field),
				Example:     fields.ExampleTag(field),
				Description: fields.DescriptionTag(field),
			}

		default:
			properties[fieldJsonTag] = g.defaultProperty(field)

		}
	}

	return properties
}

func (g DefinitionGenerator) timeProperty(field reflect.StructField, required bool) SchemaProperty {
	return SchemaProperty{
		Type:        "string",
		Format:      "date-time",
		IsRequired:  required,
		Nullable:    field.Type.Kind() == reflect.Pointer,
		Example:     fields.ExampleTag(field),
		Description: fields.DescriptionTag(field),
	}
}

func (g DefinitionGenerator) durationProperty(field reflect.StructField, required bool) SchemaProperty {
	return SchemaProperty{
		Type:        "integer",
		IsRequired:  required,
		Nullable:    field.Type.Kind() == reflect.Pointer,
		Example:     fields.ExampleTag(field),
		Description: fields.DescriptionTag(field),
	}
}

func (g DefinitionGenerator) refProperty(field reflect.StructField, required bool) SchemaProperty {
	return SchemaProperty{
		Ref:         fmt.Sprintf("#/components/schemas/%s", field.Type.Elem().String()),
		IsRequired:  required,
		Nullable:    true,
		Example:     fields.ExampleTag(field),
		Description: fields.DescriptionTag(field),
	}
}

func (g DefinitionGenerator) defaultProperty(field reflect.StructField) SchemaProperty {
	return SchemaProperty{
		Type:        fields.Type(field.Type.Kind().String()),
		IsRequired:  g.isRequired(field),
		Nullable:    field.Type.Kind() == reflect.Pointer,
		Example:     fields.ExampleTag(field),
		Description: fields.DescriptionTag(field),
	}
}

func (g DefinitionGenerator) isRequired(field reflect.StructField) bool {
	hasRequiredTag := fields.IsRequired(field)
	hasOmitemptyTag := fields.IsOmitempty(field)
	return hasRequiredTag || !hasOmitemptyTag
}
