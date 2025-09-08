package definition

import "fmt"

// EnhancedSchema extends the Schema with all OpenAPI 3.0.3 features
type EnhancedSchema struct {
	Schema // Embed existing schema

	// Additional OpenAPI 3.0.3 fields
	Summary string `json:"summary,omitempty"`

	// JSON Schema 2020-12 compatibility
	If   *EnhancedSchema `json:"if,omitempty"`
	Then *EnhancedSchema `json:"then,omitempty"`
	Else *EnhancedSchema `json:"else,omitempty"`

	// Content validation
	ContentMediaType string `json:"contentMediaType,omitempty"`
	ContentEncoding  string `json:"contentEncoding,omitempty"`

	// Extended validation
	Contains              *EnhancedSchema `json:"contains,omitempty"`
	MinContains           *int64          `json:"minContains,omitempty"`
	MaxContains           *int64          `json:"maxContains,omitempty"`
	UnevaluatedItems      *EnhancedSchema `json:"unevaluatedItems,omitempty"`
	UnevaluatedProperties *EnhancedSchema `json:"unevaluatedProperties,omitempty"`

	// Additional string validation
	MinDate *string `json:"minDate,omitempty"`
	MaxDate *string `json:"maxDate,omitempty"`

	// Additional number validation
	ExclusiveMinimumValue *float64 `json:"exclusiveMinimum,omitempty"` // JSON Schema Draft 7 style
	ExclusiveMaximumValue *float64 `json:"exclusiveMaximum,omitempty"` // JSON Schema Draft 7 style
}

// Enhanced discriminator with mapping support
type AdvancedDiscriminator struct {
	PropertyName string            `json:"propertyName"` // REQUIRED
	Mapping      map[string]string `json:"mapping,omitempty"`
}

// Enhanced XML object
type AdvancedXML struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Attribute bool   `json:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty"`
}

// NewEnhancedSchema creates a new enhanced schema
func NewEnhancedSchema(schemaType string) *EnhancedSchema {
	return &EnhancedSchema{
		Schema: Schema{
			Type: schemaType,
		},
	}
}

// NewEnhancedSchemaFromSchema creates an enhanced schema from existing schema
func NewEnhancedSchemaFromSchema(schema Schema) *EnhancedSchema {
	return &EnhancedSchema{
		Schema: schema,
	}
}

// SetSummary sets the summary for the schema
func (es *EnhancedSchema) SetSummary(summary string) *EnhancedSchema {
	es.Summary = summary
	return es
}

// SetContentMediaType sets the content media type
func (es *EnhancedSchema) SetContentMediaType(mediaType string) *EnhancedSchema {
	es.ContentMediaType = mediaType
	return es
}

// SetContentEncoding sets the content encoding
func (es *EnhancedSchema) SetContentEncoding(encoding string) *EnhancedSchema {
	es.ContentEncoding = encoding
	return es
}

// SetIfThenElse sets conditional schema validation
func (es *EnhancedSchema) SetIfThenElse(ifSchema, thenSchema, elseSchema *EnhancedSchema) *EnhancedSchema {
	es.If = ifSchema
	es.Then = thenSchema
	es.Else = elseSchema
	return es
}

// SetContains sets the contains schema for arrays
func (es *EnhancedSchema) SetContains(contains *EnhancedSchema, min, max *int64) *EnhancedSchema {
	es.Contains = contains
	es.MinContains = min
	es.MaxContains = max
	return es
}

// SetUnevaluatedItems sets the unevaluated items schema
func (es *EnhancedSchema) SetUnevaluatedItems(schema *EnhancedSchema) *EnhancedSchema {
	es.UnevaluatedItems = schema
	return es
}

// SetUnevaluatedProperties sets the unevaluated properties schema
func (es *EnhancedSchema) SetUnevaluatedProperties(schema *EnhancedSchema) *EnhancedSchema {
	es.UnevaluatedProperties = schema
	return es
}

// SetExclusiveMinimumValue sets exclusive minimum value (JSON Schema Draft 7 style)
func (es *EnhancedSchema) SetExclusiveMinimumValue(min float64) *EnhancedSchema {
	es.ExclusiveMinimumValue = &min
	return es
}

// SetExclusiveMaximumValue sets exclusive maximum value (JSON Schema Draft 7 style)
func (es *EnhancedSchema) SetExclusiveMaximumValue(max float64) *EnhancedSchema {
	es.ExclusiveMaximumValue = &max
	return es
}

// SetDateRange sets date range validation
func (es *EnhancedSchema) SetDateRange(minDate, maxDate *string) *EnhancedSchema {
	es.MinDate = minDate
	es.MaxDate = maxDate
	return es
}

// Validate validates the enhanced schema for OpenAPI 3.0.3 compliance
func (es *EnhancedSchema) Validate() error {
	// Enhanced schema specific validations
	if es.MinContains != nil && es.MaxContains != nil {
		if *es.MinContains > *es.MaxContains {
			return fmt.Errorf("minContains (%d) cannot be greater than maxContains (%d)", *es.MinContains, *es.MaxContains)
		}
	}

	if es.ExclusiveMinimumValue != nil && es.ExclusiveMaximumValue != nil {
		if *es.ExclusiveMinimumValue >= *es.ExclusiveMaximumValue {
			return fmt.Errorf("exclusiveMinimum (%f) must be less than exclusiveMaximum (%f)", *es.ExclusiveMinimumValue, *es.ExclusiveMaximumValue)
		}
	}

	return nil
}
