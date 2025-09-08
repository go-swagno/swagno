package http

import "github.com/go-swagno/swagno/v3/components/definition"

// EnhancedRequestBody with better content type support
type EnhancedRequestBody struct {
	Description string                       `json:"description,omitempty"`
	Content     map[string]EnhancedMediaType `json:"content"` // REQUIRED
	Required    bool                         `json:"required,omitempty"`
}

// EnhancedMediaType with encoding support
type EnhancedMediaType struct {
	Schema   *definition.EnhancedSchema `json:"schema,omitempty"`
	Example  interface{}                `json:"example,omitempty"`
	Examples map[string]*Example        `json:"examples,omitempty"`
	Encoding map[string]*Encoding       `json:"encoding,omitempty"`
}

// Encoding object for multipart and form data
type Encoding struct {
	ContentType   string             `json:"contentType,omitempty"`
	Headers       map[string]*Header `json:"headers,omitempty"`
	Style         string             `json:"style,omitempty"`
	Explode       bool               `json:"explode,omitempty"`
	AllowReserved bool               `json:"allowReserved,omitempty"`
}

// Example object
type Example struct {
	Summary       string      `json:"summary,omitempty"`
	Description   string      `json:"description,omitempty"`
	Value         interface{} `json:"value,omitempty"`
	ExternalValue string      `json:"externalValue,omitempty"`
}

// Header object
type Header struct {
	Description     string                        `json:"description,omitempty"`
	Required        bool                          `json:"required,omitempty"`
	Deprecated      bool                          `json:"deprecated,omitempty"`
	AllowEmptyValue bool                          `json:"allowEmptyValue,omitempty"`
	Style           string                        `json:"style,omitempty"`
	Explode         bool                          `json:"explode,omitempty"`
	AllowReserved   bool                          `json:"allowReserved,omitempty"`
	Schema          *definition.EnhancedSchema    `json:"schema,omitempty"`
	Example         interface{}                   `json:"example,omitempty"`
	Examples        map[string]*Example           `json:"examples,omitempty"`
	Content         map[string]*EnhancedMediaType `json:"content,omitempty"`
}

// EnhancedResponse extends response with enhanced features
type EnhancedResponse struct {
	Description string                       `json:"description"` // REQUIRED
	Headers     map[string]*Header           `json:"headers,omitempty"`
	Content     map[string]EnhancedMediaType `json:"content,omitempty"`
	Links       map[string]*Link             `json:"links,omitempty"`
}

// Link represents a link object in responses
type Link struct {
	OperationRef string                 `json:"operationRef,omitempty"`
	OperationId  string                 `json:"operationId,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	RequestBody  interface{}            `json:"requestBody,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Server       *Server                `json:"server,omitempty"`
}

// Server represents a server for links
type Server struct {
	URL         string                    `json:"url"`
	Description string                    `json:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty"`
}

// ServerVariable represents a server variable
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default"`
	Description string   `json:"description,omitempty"`
}

// NewEnhancedRequestBody creates a new enhanced request body
func NewEnhancedRequestBody(description string, required bool) *EnhancedRequestBody {
	return &EnhancedRequestBody{
		Description: description,
		Content:     make(map[string]EnhancedMediaType),
		Required:    required,
	}
}

// AddContent adds content to the request body
func (erb *EnhancedRequestBody) AddContent(mediaType string, schema *definition.EnhancedSchema) *EnhancedRequestBody {
	erb.Content[mediaType] = EnhancedMediaType{
		Schema: schema,
	}
	return erb
}

// AddContentWithExample adds content with example to the request body
func (erb *EnhancedRequestBody) AddContentWithExample(mediaType string, schema *definition.EnhancedSchema, example interface{}) *EnhancedRequestBody {
	erb.Content[mediaType] = EnhancedMediaType{
		Schema:  schema,
		Example: example,
	}
	return erb
}

// NewEnhancedResponse creates a new enhanced response
func NewEnhancedResponse(description string) *EnhancedResponse {
	return &EnhancedResponse{
		Description: description,
		Headers:     make(map[string]*Header),
		Content:     make(map[string]EnhancedMediaType),
		Links:       make(map[string]*Link),
	}
}

// AddHeader adds a header to the response
func (er *EnhancedResponse) AddHeader(name string, header *Header) *EnhancedResponse {
	er.Headers[name] = header
	return er
}

// AddContent adds content to the response
func (er *EnhancedResponse) AddContent(mediaType string, schema *definition.EnhancedSchema) *EnhancedResponse {
	er.Content[mediaType] = EnhancedMediaType{
		Schema: schema,
	}
	return er
}

// AddLink adds a link to the response
func (er *EnhancedResponse) AddLink(name string, link *Link) *EnhancedResponse {
	er.Links[name] = link
	return er
}

// NewExample creates a new example
func NewExample(summary, description string, value interface{}) *Example {
	return &Example{
		Summary:     summary,
		Description: description,
		Value:       value,
	}
}

// NewExampleWithExternalValue creates a new example with external value
func NewExampleWithExternalValue(summary, description, externalValue string) *Example {
	return &Example{
		Summary:       summary,
		Description:   description,
		ExternalValue: externalValue,
	}
}

// NewHeader creates a new header
func NewHeader(description string, schema *definition.EnhancedSchema) *Header {
	return &Header{
		Description: description,
		Schema:      schema,
	}
}

// SetRequired sets the required flag for the header
func (h *Header) SetRequired(required bool) *Header {
	h.Required = required
	return h
}

// SetDeprecated marks the header as deprecated
func (h *Header) SetDeprecated(deprecated bool) *Header {
	h.Deprecated = deprecated
	return h
}

// SetStyle sets the serialization style for the header
func (h *Header) SetStyle(style string) *Header {
	h.Style = style
	return h
}

// NewLink creates a new link
func NewLink(operationId, description string) *Link {
	return &Link{
		OperationId: operationId,
		Description: description,
		Parameters:  make(map[string]interface{}),
	}
}

// NewLinkWithRef creates a new link with operation reference
func NewLinkWithRef(operationRef, description string) *Link {
	return &Link{
		OperationRef: operationRef,
		Description:  description,
		Parameters:   make(map[string]interface{}),
	}
}

// AddParameter adds a parameter to the link
func (l *Link) AddParameter(name string, value interface{}) *Link {
	l.Parameters[name] = value
	return l
}

// SetRequestBody sets the request body for the link
func (l *Link) SetRequestBody(requestBody interface{}) *Link {
	l.RequestBody = requestBody
	return l
}
