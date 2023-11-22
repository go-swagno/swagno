package swagno

// https://swagger.io/specification/v2/#pathsObject
type jsonEndpoint struct {
	Description string                  `json:"description"`
	Consumes    []string                `json:"consumes" default:"application/json"`
	Produces    []string                `json:"produces" default:"application/json"`
	Tags        []string                `json:"tags"`
	Summary     string                  `json:"summary"`
	OperationId string                  `json:"operationId,omitempty"`
	Parameters  []jsonParameter         `json:"parameters"`
	Responses   map[string]jsonResponse `json:"responses"`
	Security    []map[string][]string   `json:"security,omitempty"`
}

// https://swagger.io/specification/v2/#parameterObject
type jsonParameter struct {
	Type              string              `json:"type"` // TODO Need to update logic to only show when using `In` == "query". I'm getting this error when running my test "Structural error at paths./product.post.parameters.0; should NOT have additional properties; additionalProperty: type"
	Description       string              `json:"description"`
	Name              string              `json:"name"`
	In                string              `json:"in"`
	Required          bool                `json:"required"`
	Schema            *jsonResponseScheme `json:"schema,omitempty"`
	Format            string              `json:"format,omitempty"`
	Enum              []interface{}       `json:"enum,omitempty"`
	Default           interface{}         `json:"default,omitempty"`
	Min               int64               `json:"minimum,omitempty"`
	Max               int64               `json:"maximum,omitempty"`
	MinLen            int64               `json:"minLength,omitempty"`
	MaxLen            int64               `json:"maxLength,omitempty"`
	Pattern           string              `json:"pattern,omitempty"`
	MaxItems          int64               `json:"maxItems,omitempty"`
	MinItems          int64               `json:"minItems,omitempty"`
	UniqueItems       bool                `json:"uniqueItems,omitempty"`
	MultipleOf        int64               `json:"multipleOf,omitempty"`
	CollenctionFormat string              `json:"collectionFormat,omitempty"`
}

// https://swagger.io/specification/v2/#response-object
type jsonResponse struct {
	Description string              `json:"description"`
	Schema      *jsonResponseScheme `json:"schema,omitempty"`
}

// https://swagger.io/specification/v2/#schema-object
type jsonResponseScheme struct {
	Ref   string                   `json:"$ref,omitempty"`
	Type  string                   `json:"type,omitempty"`
	Items *jsonResponseSchemeItems `json:"items,omitempty"`
}

type jsonResponseSchemeItems struct {
	Type string `json:"type,omitempty"`
	Ref  string `json:"$ref,omitempty"`
}
