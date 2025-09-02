package endpoint

import (
	"fmt"
	"reflect"

	"github.com/go-swagno/swagno/v3/components/http/response"
	"github.com/go-swagno/swagno/v3/components/mime"
	"github.com/go-swagno/swagno/v3/components/parameter"
)

// MethodType represents HTTP request methods.
type MethodType string

const (
	GET     MethodType = "GET"
	POST    MethodType = "POST"
	PUT     MethodType = "PUT"
	DELETE  MethodType = "DELETE"
	PATCH   MethodType = "PATCH"
	OPTIONS MethodType = "OPTIONS"
	HEAD    MethodType = "HEAD"
	TRACE   MethodType = "TRACE" // New in OpenAPI 3.0
)

// MediaType represents a media type object in OpenAPI 3.0
// https://spec.openapis.org/oas/v3.0.3#media-type-object
type MediaType struct {
	Schema   *parameter.JsonResponseSchema `json:"schema,omitempty"`
	Example  interface{}                   `json:"example,omitempty"`
	Examples map[string]interface{}        `json:"examples,omitempty"`
	Encoding map[string]interface{}        `json:"encoding,omitempty"`
}

// RequestBody represents a request body in OpenAPI 3.0
// https://spec.openapis.org/oas/v3.0.3#request-body-object
type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Content     map[string]MediaType `json:"content"`
	Required    bool                 `json:"required,omitempty"`
}

// JsonEndPoint is the JSON model version of EndPoint object used for API purposes
// https://spec.openapis.org/oas/v3.0.3#operation-object
type JsonEndPoint struct {
	Tags         []string                  `json:"tags,omitempty"`
	Summary      string                    `json:"summary,omitempty"`
	Description  string                    `json:"description,omitempty"`
	ExternalDocs interface{}               `json:"externalDocs,omitempty"`
	OperationId  string                    `json:"operationId,omitempty"`
	Parameters   []parameter.JsonParameter `json:"parameters,omitempty"`
	RequestBody  *RequestBody              `json:"requestBody,omitempty"`
	Responses    map[string]JsonResponse   `json:"responses"`
	Callbacks    map[string]interface{}    `json:"callbacks,omitempty"`
	Deprecated   bool                      `json:"deprecated,omitempty"`
	Security     []map[string][]string     `json:"security,omitempty"`
	Servers      []interface{}             `json:"servers,omitempty"`

	// Legacy fields for compatibility
	Consumes []mime.MIME `json:"consumes,omitempty"`
	Produces []mime.MIME `json:"produces,omitempty"`
}

// JsonResponse represents the structure of a response in the OpenAPI 3.0 specification.
// It encapsulates the description, content, headers, and links of a response object.
// See: https://spec.openapis.org/oas/v3.0.3#response-object
type JsonResponse struct {
	Description string                 `json:"description"`
	Headers     map[string]interface{} `json:"headers,omitempty"`
	Content     map[string]MediaType   `json:"content,omitempty"`
	Links       map[string]interface{} `json:"links,omitempty"`

	// Legacy field for compatibility
	Schema *parameter.JsonResponseSchema `json:"schema,omitempty"`
}

// EndPoint holds the details of an API endpoint, including HTTP method, path, parameters,
// request body, responses, and metadata such as tags and security requirements.
type EndPoint struct {
	method            MethodType
	path              string
	params            []*parameter.Parameter
	tags              []string
	Body              interface{}
	successfulReturns []response.Response
	errors            []response.Response
	description       string
	summary           string
	consume           []mime.MIME
	produce           []mime.MIME
	security          []map[string][]string
	deprecated        bool
	callbacks         map[string]interface{}
	servers           []interface{}
}

// AsJson converts an EndPoint into its JSON representation as JsonEndPoint.
// This function is typically used when the endpoint needs to be serialized
// into a format that can be used by OpenAPI or other API documentation tools.
func (e *EndPoint) AsJson() JsonEndPoint {
	return JsonEndPoint{
		Tags:        e.tags,
		Summary:     e.summary,
		Description: e.description,
		Security:    e.security,
		Deprecated:  e.deprecated,
		Callbacks:   e.callbacks,
		Servers:     e.servers,
		// Legacy compatibility
		Consumes: e.consume,
		Produces: e.produce,
	}
}

// EndPointOption is a function type that takes an EndPoint pointer allowing
// for modular and configurable setting of various EndPoint properties.
type EndPointOption func(e *EndPoint)

func endpoint() *EndPoint {
	return &EndPoint{
		method:            "",
		path:              "",
		params:            []*parameter.Parameter{},
		tags:              []string{},
		successfulReturns: []response.Response{},
		errors:            []response.Response{},
		description:       "",
		summary:           "",
		consume:           []mime.MIME{mime.JSON},
		produce:           []mime.MIME{mime.JSON},
		security:          nil,
		deprecated:        false,
		callbacks:         nil,
		servers:           nil,
	}
}

// Params returns the parameters associated with the EndPoint.
func (e *EndPoint) Params() []*parameter.Parameter {
	return e.params
}

// SuccessfulReturns lists the possible successful response types that the EndPoint can return.
func (e *EndPoint) SuccessfulReturns() []response.Response {
	return e.successfulReturns
}

// Errors lists the possible error responses that the EndPoint can return.
func (e *EndPoint) Errors() []response.Response {
	return e.errors
}

// Method returns the HTTP method type (e.g., GET, POST) associated with the EndPoint.
func (e *EndPoint) Method() MethodType {
	return e.method
}

// Path returns the URL path of the EndPoint.
func (e *EndPoint) Path() string {
	return e.path
}

// BodyJsonParameter creates the request body parameter for OpenAPI 3.0.
// In OpenAPI 3.0, request bodies are handled differently than in Swagger 2.0
func (e *EndPoint) BodyJsonParameter() *parameter.JsonParameter {
	if e.Body != nil {
		bodyRef := fmt.Sprintf("#/components/schemas/%T", e.Body)
		bodySchema := parameter.JsonResponseSchema{
			Ref: bodyRef,
		}

		if reflect.TypeOf(e.Body).Kind() == reflect.Slice {
			bodySchema = parameter.JsonResponseSchema{
				Type: "array",
				Items: &parameter.JsonResponseSchemeItems{
					Ref: fmt.Sprintf("#/components/schemas/%T", e.Body),
				},
			}
		}

		return &parameter.JsonParameter{
			Name:        "body",
			In:          "body",
			Description: "Request body",
			Required:    true,
			Schema:      &bodySchema,
		}
	}

	return nil
}

// WithConsume sets the MIME types that the EndPoint can consume.
// This is typically used for specifying the expected request body formats.
func WithConsume(consume []mime.MIME) EndPointOption {
	return func(e *EndPoint) {
		e.consume = consume
	}
}

// WithProduce sets the MIME types that the EndPoint can produce.
// This is typically used for specifying the response body formats.
func WithProduce(produce []mime.MIME) EndPointOption {
	return func(e *EndPoint) {
		e.produce = produce
	}
}

// WithTags assigns a set of tags to the EndPoint, which can be used for organizing and categorizing endpoints.
func WithTags(tag ...string) EndPointOption {
	return func(e *EndPoint) {
		e.tags = append(e.tags, tag...)
	}
}

// WithParams sets the parameters for the EndPoint, defining what data can be accepted by the endpoint.
func WithParams(params ...*parameter.Parameter) EndPointOption {
	return func(e *EndPoint) {
		e.params = append(e.params, params...)
	}
}

// WithBody specifies the data structure that the EndPoint expects in the request body.
func WithBody(body interface{}) EndPointOption {
	return func(e *EndPoint) {
		e.Body = body
	}
}

// WithSuccessfulReturns sets the possible successful response types that the EndPoint can return.
func WithSuccessfulReturns(ret []response.Response) EndPointOption {
	return func(e *EndPoint) {
		e.successfulReturns = ret
	}
}

// WithErrors defines the error responses that the EndPoint can return, allowing for detailed error handling.
func WithErrors(err []response.Response) EndPointOption {
	return func(e *EndPoint) {
		e.errors = err
	}
}

// WithDescription sets a descriptive text for the EndPoint, providing context or information about its purpose.
func WithDescription(des string) EndPointOption {
	return func(e *EndPoint) {
		e.description = des
	}
}

// WithSummary provides a brief summary of what the EndPoint does, which can be useful for quick reference.
func WithSummary(s string) EndPointOption {
	return func(e *EndPoint) {
		e.summary = s
	}
}

// WithSecurity defines the security requirements for the EndPoint, such as authentication or authorization details.
func WithSecurity(security []map[string][]string) EndPointOption {
	return func(e *EndPoint) {
		e.security = security
	}
}

// WithDeprecated marks the endpoint as deprecated.
func WithDeprecated() EndPointOption {
	return func(e *EndPoint) {
		e.deprecated = true
	}
}

// WithCallbacks adds callbacks to the endpoint.
func WithCallbacks(callbacks map[string]interface{}) EndPointOption {
	return func(e *EndPoint) {
		e.callbacks = callbacks
	}
}

// WithServers adds servers specific to this endpoint.
func WithServers(servers []interface{}) EndPointOption {
	return func(e *EndPoint) {
		e.servers = servers
	}
}

// New creates a new EndPoint with the specified HTTP method and path, and applies any provided EndPointOptions.
// This is the primary constructor function for creating a new EndPoint instance.
func New(m MethodType, path string, opts ...EndPointOption) *EndPoint {
	e := endpoint()
	e.method = m
	e.path = path

	for _, opt := range opts {
		opt(e)
	}

	return e
}
