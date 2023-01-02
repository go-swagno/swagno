package swagno

// https://swagger.io/specification/v2/#definitionsObject
type swaggerDefinition struct {
	Type       string                                 `json:"type"`
	Properties map[string]swaggerDefinitionProperties `json:"properties"`
}

// https://swagger.io/specification/v2/#schemaObject
type swaggerDefinitionProperties struct {
	Type   string                            `json:"type,omitempty"`
	Format string                            `json:"format,omitempty"`
	Ref    string                            `json:"$ref,omitempty"`
	Items  *swaggerDefinitionPropertiesItems `json:"items,omitempty"`
}

type swaggerDefinitionPropertiesItems struct {
	Type string `json:"type,omitempty"`
	Ref  string `json:"$ref,omitempty"`
}