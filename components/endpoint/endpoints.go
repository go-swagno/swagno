package endpoint

import (
	"strings"

	"github.com/domhoward14/swagno/components/parameter"
	"github.com/domhoward14/swagno/components/response"
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

// EndPoint represents an API endpoint.
type EndPoint struct {
	Method            string                `json:"method"`
	Path              string                `json:"path"`
	Params            []parameter.Parameter `json:"params"`
	Tags              []string              `json:"tags"`
	SuccessfulReturns []response.Info       `json:"return"`
	Errors            []response.Info       `json:"error"`
	Body              interface{}           `json:"body"` // TODO Could probably make this a better defined interface. would probaly result in simpler parsing logic in the 'generateSwaggerObject()'.
	Description       string                `json:"description"`
	Summary           string                `json:"summary"`
	Consume           []string              `json:"consume"`
	Produce           []string              `json:"produce"`
	Security          []map[string][]string `json:"security"` // TODO make into struct
}

type EndPointOption func(e *EndPoint)

type EndPoints struct {
	List []EndPoint `json:"EndPoints"`
}

func getEndpoint() *EndPoint {
	return &EndPoint{
		Method:            "",
		Path:              "",
		Params:            []parameter.Parameter{},
		Tags:              []string{},
		SuccessfulReturns: []response.Info{},
		Errors:            []response.Info{},
		Body:              nil,
		Description:       "",
		Summary:           "",
		Consume:           []string{"application/json"},
		Produce:           []string{"application/json"},
		Security:          nil,
	}
}

func WithConsume(consume []string) EndPointOption {
	return func(e *EndPoint) {
		e.Consume = consume
	}
}

func WithProduce(produce []string) EndPointOption {
	return func(e *EndPoint) {
		e.Produce = produce
	}
}

func WithMethod(method MethodType) EndPointOption {
	return func(e *EndPoint) {
		e.Method = string(method)
	}
}
func WithPath(path string) EndPointOption {
	return func(e *EndPoint) {
		e.Path = path
	}
}

func WithTags(tags string) EndPointOption {
	return func(e *EndPoint) {
		e.Tags = strings.Split(strings.ReplaceAll(tags, " ", ""), ",")
	}
}

func WithParams(params []parameter.Parameter) EndPointOption {
	return func(e *EndPoint) {
		e.Params = params
	}
}

func WithBody(body interface{}) EndPointOption {
	return func(e *EndPoint) {
		e.Body = body
	}
}

func WithSuccessfulReturns(ret []response.Info) EndPointOption {
	return func(e *EndPoint) {
		e.SuccessfulReturns = ret
	}
}

func WithErrors(err []response.Info) EndPointOption {
	return func(e *EndPoint) {
		e.Errors = err
	}
}

func WithDescription(des string) EndPointOption {
	return func(e *EndPoint) {
		e.Description = des
	}
}

func WithSecurity(security []map[string][]string) EndPointOption {
	return func(e *EndPoint) {
		e.Security = security
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
