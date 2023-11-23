package definition

// https://swagger.io/specification/v2/#definitionsObject
type Definition struct {
	Type       string                          `json:"type"`
	Properties map[string]DefinitionProperties `json:"properties"`
}

// https://swagger.io/specification/v2/#schemaObject
type DefinitionProperties struct {
	Type    string                     `json:"type,omitempty"`
	Format  string                     `json:"format,omitempty"`
	Ref     string                     `json:"$ref,omitempty"`
	Items   *DefinitionPropertiesItems `json:"items,omitempty"`
	Example interface{}                `json:"example,omitempty"`
}

type DefinitionPropertiesItems struct {
	Type string `json:"type,omitempty"`
	Ref  string `json:"$ref,omitempty"`
}
