package endpoint

import (
	"fmt"
	"reflect"

	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/mime"
	"github.com/go-swagno/swagno/components/parameter"
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
)

// JsonEndPoint is the JSON model version of EndPoint object used for API purposes
// https://swagger.io/specification/v2/#pathsObject
type JsonEndPoint struct {
	Description string                    `json:"description"`
	Consumes    []mime.MIME               `json:"consumes" default:"application/json"`
	Produces    []mime.MIME               `json:"produces" default:"application/json"`
	Tags        []string                  `json:"tags"`
	Summary     string                    `json:"summary"`
	OperationId string                    `json:"operationId,omitempty"`
	Parameters  []parameter.JsonParameter `json:"parameters"`
	Responses   map[string]JsonResponse   `json:"responses"`
	Security    []map[string][]string     `json:"security,omitempty"`
}

// JsonResponse represents the structure of a response in the Swagger 2.0 specification.
// It encapsulates the description and schema of a response object.
// See: https://swagger.io/specification/v2/#response-object
type JsonResponse struct {
	Description string                        `json:"description"`
	Schema      *parameter.JsonResponseSchema `json:"schema,omitempty"`
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
}

// AsJson converts an EndPoint into its JSON representation as JsonEndPoint.
// This function is typically used when the endpoint needs to be serialized
// into a format that can be used by Swagger or other API documentation tools.
func (e *EndPoint) AsJson() JsonEndPoint {
	return JsonEndPoint{
		Description: e.description,
		Summary:     e.summary,
		Consumes:    e.consume,
		Produces:    e.produce,
		Tags:        e.tags,
		Security:    e.security,
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

// BodyJsonParameter makes the body definitions and parameter for body if present. Parameters for body are described via schema
// definition so that's why it doesn't use the 'Parameter' object like the other ones.
func (e *EndPoint) BodyJsonParameter() *parameter.JsonParameter {
	if e.Body != nil {
		bodyRef := fmt.Sprintf("#/definitions/%T", e.Body)
		bodySchema := parameter.JsonResponseSchema{
			Ref: bodyRef,
		}

		if reflect.TypeOf(e.Body).Kind() == reflect.Slice {
			bodySchema = parameter.JsonResponseSchema{
				Type: "array",
				Items: &parameter.JsonResponseSchemeItems{
					Ref: fmt.Sprintf("#/definitions/%T", e.Body),
				},
			}
		}

		return &parameter.JsonParameter{
			Name:        "body",
			In:          "body",
			Description: "body",
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
