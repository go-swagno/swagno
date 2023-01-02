package swagno

// https://swagger.io/specification/v2/#pathsObject
type swaggerEndpoint struct {
	Description string                     `json:"description"`
	Consumes    []string                   `json:"consumes" default:"application/json"`
	Produces    []string                   `json:"produces" default:"application/json"`
	Tags        []string                   `json:"tags"`
	Summary     string                     `json:"summary"`
	OperationId string                     `json:"operationId,omitempty"`
	Parameters  []swaggerParameter         `json:"parameters"`
	Responses   map[string]swaggerResponse `json:"responses"`
	Security    []map[string][]string      `json:"security,omitempty"`
}

// https://swagger.io/specification/v2/#parameterObject
type swaggerParameter struct {
	Type              string                 `json:"type"`
	Description       string                 `json:"description"`
	Name              string                 `json:"name"`
	In                string                 `json:"in"`
	Required          bool                   `json:"required"`
	Schema            *swaggerResponseScheme `json:"schema,omitempty"`
	Format            string                 `json:"format,omitempty"`
	Items             *ParameterItems        `json:"items,omitempty"`
	Enum              []interface{}          `json:"enum,omitempty"`
	Default           interface{}            `json:"default,omitempty"`
	Min               int64                  `json:"minimum,omitempty"`
	Max               int64                  `json:"maximum,omitempty"`
	MinLen            int64                  `json:"minLength,omitempty"`
	MaxLen            int64                  `json:"maxLength,omitempty"`
	Pattern           string                 `json:"pattern,omitempty"`
	MaxItems          int64                  `json:"maxItems,omitempty"`
	MinItems          int64                  `json:"minItems,omitempty"`
	UniqueItems       bool                   `json:"uniqueItems,omitempty"`
	MultipleOf        int64                  `json:"multipleOf,omitempty"`
	CollenctionFormat string                 `json:"collectionFormat,omitempty"`
}

// https://swagger.io/specification/v2/#response-object
type swaggerResponse struct {
	Description string                `json:"description"`
	Schema      swaggerResponseScheme `json:"schema"`
}

// https://swagger.io/specification/v2/#schema-object
type swaggerResponseScheme struct {
	Ref   string                      `json:"$ref,omitempty"`
	Type  string                      `json:"type,omitempty"`
	Items *swaggerResponseSchemeItems `json:"items,omitempty"`
}


type swaggerResponseSchemeItems struct {
	Type string `json:"type,omitempty"`
	Ref  string `json:"$ref,omitempty"`
}