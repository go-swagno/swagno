package swagno

// https://swagger.io/specification/v2/#definitionsObject
type Definition struct {
	Type       string                              `json:"type"`
	Properties map[string]jsonDefinitionProperties `json:"properties"`
}

// https://swagger.io/specification/v2/#schemaObject
type jsonDefinitionProperties struct {
	Type    string                         `json:"type,omitempty"`
	Format  string                         `json:"format,omitempty"`
	Ref     string                         `json:"$ref,omitempty"`
	Items   *jsonDefinitionPropertiesItems `json:"items,omitempty"`
	Example interface{}                    `json:"example,omitempty"`
}

type jsonDefinitionPropertiesItems struct {
	Type string `json:"type,omitempty"`
	Ref  string `json:"$ref,omitempty"`
}
