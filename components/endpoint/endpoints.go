package endpoint

import (
	"fmt"
	"reflect"

	"github.com/go-swagno/swagno/components/mime"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/go-swagno/swagno/http/response"
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

// https://swagger.io/specification/v2/#response-object
type JsonResponse struct {
	Description string                        `json:"description"`
	Schema      *parameter.JsonResponseSchema `json:"schema,omitempty"`
}

// EndPoint represents an API endpoint.
type EndPoint struct {
	method            string
	path              string
	params            []*parameter.Parameter
	tags              []string
	Body              interface{}
	successfulReturns []response.Info
	errors            []response.Info
	description       string
	summary           string
	consume           []mime.MIME
	produce           []mime.MIME
	security          []map[string][]string
}

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

type EndPointOption func(e *EndPoint)

func getEndpoint() *EndPoint {
	return &EndPoint{
		method:            "",
		path:              "",
		params:            []*parameter.Parameter{},
		tags:              []string{},
		successfulReturns: []response.Info{},
		errors:            []response.Info{},
		description:       "",
		summary:           "",
		consume:           []mime.MIME{mime.JSON},
		produce:           []mime.MIME{mime.JSON},
		security:          nil,
	}
}

func (e *EndPoint) GetParams() []*parameter.Parameter {
	return e.params
}

func (e *EndPoint) GetSuccessfulReturns() []response.Info {
	return e.successfulReturns
}

func (e *EndPoint) GetErrors() []response.Info {
	return e.errors
}

func (e *EndPoint) GetMethod() string {
	return e.method
}

func (e *EndPoint) GetPath() string {
	return e.path
}

// GetBodyJsonParameter makes the body definitions and parameter for body if present. Parameters for body are described via schema
// definition so that's why it doesn't use the 'Parameter' object like the other ones.
func (e *EndPoint) GetBodyJsonParameter() *parameter.JsonParameter {
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

func WithConsume(consume []mime.MIME) EndPointOption {
	return func(e *EndPoint) {
		e.consume = consume
	}
}

func WithProduce(produce []mime.MIME) EndPointOption {
	return func(e *EndPoint) {
		e.produce = produce
	}
}

func WithMethod(method MethodType) EndPointOption {
	return func(e *EndPoint) {
		e.method = string(method)
	}
}
func WithPath(path string) EndPointOption {
	return func(e *EndPoint) {
		e.path = path
	}
}

func WithTags(tag ...string) EndPointOption {
	return func(e *EndPoint) {
		e.tags = tag
	}
}

func WithParams(params ...*parameter.Parameter) EndPointOption {
	return func(e *EndPoint) {
		e.params = params
	}
}

func WithBody(body interface{}) EndPointOption {
	return func(e *EndPoint) {
		e.Body = body
	}
}

func WithSuccessfulReturns(ret []response.Info) EndPointOption {
	return func(e *EndPoint) {
		e.successfulReturns = ret
	}
}

func WithErrors(err []response.Info) EndPointOption {
	return func(e *EndPoint) {
		e.errors = err
	}
}

func WithDescription(des string) EndPointOption {
	return func(e *EndPoint) {
		e.description = des
	}
}

func WithSummary(s string) EndPointOption {
	return func(e *EndPoint) {
		e.summary = s
	}
}

func WithSecurity(security []map[string][]string) EndPointOption {
	return func(e *EndPoint) {
		e.security = security
	}
}

// EndPoint is a function to create an API endpoint.
// args: method, path, tags, params, body, return, error, description, security, consume, produce
func New(opts ...EndPointOption) *EndPoint {
	e := getEndpoint()

	for _, opt := range opts {
		opt(e)
	}

	return e
}
