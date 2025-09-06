package parameter

// OpenAPI3Parameter represents a parameter in OpenAPI 3.0 format
type OpenAPI3Parameter struct {
	Name            string                 `json:"name"` // REQUIRED
	In              string                 `json:"in"`   // REQUIRED
	Description     string                 `json:"description,omitempty"`
	Required        bool                   `json:"required,omitempty"`
	Deprecated      bool                   `json:"deprecated,omitempty"`
	AllowEmptyValue bool                   `json:"allowEmptyValue,omitempty"`
	Style           string                 `json:"style,omitempty"`
	Explode         *bool                  `json:"explode,omitempty"`
	AllowReserved   bool                   `json:"allowReserved,omitempty"`
	Schema          *JsonResponseSchema    `json:"schema,omitempty"`
	Example         interface{}            `json:"example,omitempty"`
	Examples        map[string]interface{} `json:"examples,omitempty"`
	Content         map[string]interface{} `json:"content,omitempty"`
}

// AsOpenAPI3Json converts Parameter to OpenAPI 3.0 compliant format
func (p *Parameter) AsOpenAPI3Json() OpenAPI3Parameter {
	// Create schema object - validation fields are within the schema
	schema := &JsonResponseSchema{
		Type:        p.typeValue.String(),
		UniqueItems: p.uniqueItems,
	}

	// Min/max değerlerini ekle
	if p.min != 0 {
		min := float64(p.min)
		schema.Min = &min
	}
	if p.max != 0 {
		max := float64(p.max)
		schema.Max = &max
	}
	if p.minLen != 0 {
		schema.MinLen = &p.minLen
	}
	if p.maxLen != 0 {
		schema.MaxLen = &p.maxLen
	}
	if p.minItems != 0 {
		schema.MinItems = &p.minItems
	}
	if p.maxItems != 0 {
		schema.MaxItems = &p.maxItems
	}
	if p.multipleOf != 0 {
		multiple := float64(p.multipleOf)
		schema.MultipleOf = &multiple
	}

	openAPI3Param := OpenAPI3Parameter{
		Name:        p.name,
		In:          p.in.String(),
		Description: p.description,
		Required:    p.required,
		Deprecated:  p.deprecated,
		Schema:      schema,
	}

	// OpenAPI 3.0 specific features
	if p.style != "" {
		openAPI3Param.Style = p.style
	}
	if p.explode {
		explode := true
		openAPI3Param.Explode = &explode
	}
	if p.allowReserved {
		openAPI3Param.AllowReserved = p.allowReserved
	}
	if p.allowEmptyValue {
		openAPI3Param.AllowEmptyValue = p.allowEmptyValue
	}
	if p.example != nil {
		openAPI3Param.Example = p.example
	}

	return openAPI3Param
}

// NewOpenAPI3Parameter creates a new OpenAPI 3.0 compliant parameter
func NewOpenAPI3Parameter(name, in string, schema *JsonResponseSchema) *OpenAPI3Parameter {
	return &OpenAPI3Parameter{
		Name:   name,
		In:     in,
		Schema: schema,
	}
}

// SetRequired sets the required flag for the parameter
func (p *OpenAPI3Parameter) SetRequired(required bool) *OpenAPI3Parameter {
	p.Required = required
	return p
}

// SetDescription sets the description for the parameter
func (p *OpenAPI3Parameter) SetDescription(description string) *OpenAPI3Parameter {
	p.Description = description
	return p
}

// SetDeprecated marks the parameter as deprecated
func (p *OpenAPI3Parameter) SetDeprecated(deprecated bool) *OpenAPI3Parameter {
	p.Deprecated = deprecated
	return p
}

// SetExample sets an example value for the parameter
func (p *OpenAPI3Parameter) SetExample(example interface{}) *OpenAPI3Parameter {
	p.Example = example
	return p
}

// SetStyle sets the serialization style for the parameter
func (p *OpenAPI3Parameter) SetStyle(style string) *OpenAPI3Parameter {
	p.Style = style
	return p
}

// SetExplode sets the explode flag for the parameter
func (p *OpenAPI3Parameter) SetExplode(explode bool) *OpenAPI3Parameter {
	p.Explode = &explode
	return p
}

// SetAllowReserved sets the allowReserved flag for the parameter
func (p *OpenAPI3Parameter) SetAllowReserved(allowReserved bool) *OpenAPI3Parameter {
	p.AllowReserved = allowReserved
	return p
}

// SetAllowEmptyValue sets the allowEmptyValue flag for the parameter
func (p *OpenAPI3Parameter) SetAllowEmptyValue(allowEmptyValue bool) *OpenAPI3Parameter {
	p.AllowEmptyValue = allowEmptyValue
	return p
}
