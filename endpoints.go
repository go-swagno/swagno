package swagno

import (
	"fmt"
	"strings"
)

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

const (
	ContentTypeApplicationJSON = "application/json"
	ContentTypeApplicationXML  = "application/xml"
)

type Parameter struct {
	Name             string          `json:"name"`
	Type             string          `json:"type"`
	In               string          `json:"in"`
	Required         bool            `json:"required"`
	Description      string          `json:"description"`
	Enum             []interface{}   `json:"enum,omitempty"`
	Items            *ParameterItems `json:"items,omitempty"`
	Default          interface{}     `json:"default,omitempty"`
	Format           string          `json:"format,omitempty"`
	Min              int64           `json:"minimum,omitempty"`
	Max              int64           `json:"maximum,omitempty"`
	MinLen           int64           `json:"minLength,omitempty"`
	MaxLen           int64           `json:"maxLength,omitempty"`
	Pattern          string          `json:"pattern,omitempty"`
	MaxItems         int64           `json:"maxItems,omitempty"`
	MinItems         int64           `json:"minItems,omitempty"`
	UniqueItems      bool            `json:"uniqueItems,omitempty"`
	MultipleOf       int64           `json:"multipleOf,omitempty"`
	CollectionFormat string          `json:"collectionFormat,omitempty"`
}

type ParameterItems struct {
	Type             string        `json:"type"`
	Format           string        `json:"format,omitempty"`
	Enum             []interface{} `json:"enum,omitempty"`
	Default          interface{}   `json:"default,omitempty"`
	Min              int64         `json:"minimum,omitempty"`
	Max              int64         `json:"maximum,omitempty"`
	MinLen           int64         `json:"minLength,omitempty"`
	MaxLen           int64         `json:"maxLength,omitempty"`
	Pattern          string        `json:"pattern,omitempty"`
	MaxItems         int64         `json:"maxItems,omitempty"`
	MinItems         int64         `json:"minItems,omitempty"`
	UniqueItems      bool          `json:"uniqueItems,omitempty"`
	MultipleOf       int64         `json:"multipleOf,omitempty"`
	CollectionFormat string        `json:"collectionFormat,omitempty"`
}

type Fields struct {
	Default          interface{} `json:"default,omitempty"`
	Format           string      `json:"format,omitempty"`
	Min              int64       `json:"minimum,omitempty"`
	Max              int64       `json:"maximum,omitempty"`
	MinLen           int64       `json:"minLength,omitempty"`
	MaxLen           int64       `json:"maxLength,omitempty"`
	Pattern          string      `json:"pattern,omitempty"`
	MaxItems         int64       `json:"maxItems,omitempty"`
	MinItems         int64       `json:"minItems,omitempty"`
	UniqueItems      bool        `json:"uniqueItems,omitempty"`
	MultipleOf       int64       `json:"multipleOf,omitempty"`
	CollectionFormat string      `json:"collectionFormat,omitempty"`
}

type Endpoint struct {
	Method      string                `json:"method"`
	Path        string                `json:"path"`
	Params      []Parameter           `json:"params"`
	Tags        []string              `json:"tags"`
	Responses   []Response            `json:"responses"`
	Body        interface{}           `json:"body"`
	Description string                `json:"description"`
	Consume     []string              `json:"consume"`
	Produce     []string              `json:"produce"`
	Security    []map[string][]string `json:"security"`
}

type Response struct {
	Code        string      `json:"method"`
	Description string      `json:"description"`
	Body        interface{} `json:"return"`
}

// EndPoint args: method, path, tags, params, body, okResponse, notFoundResponse, description, security, consume, produce
func EndPoint(method MethodType, path string, tags string, params []Parameter, body interface{}, okResponse interface{}, notFoundResponse interface{}, des string, security []map[string][]string, args ...string) Endpoint {
	removedSpace := strings.ReplaceAll(tags, " ", "")
	var responses []Response

	if okResponse != nil {
		responses = append(responses, Response{
			Code:        "200",
			Description: "OK",
			Body:        okResponse,
		},
		)
	}

	if notFoundResponse != nil {
		responses = append(responses, Response{
			Code:        "400",
			Description: "Not Found",
			Body:        notFoundResponse,
		},
		)
	}

	endpoint := Endpoint{
		Method:      string(method),
		Path:        path,
		Tags:        strings.Split(removedSpace, ","),
		Params:      params,
		Responses:   responses,
		Body:        body,
		Description: des,
		Security:    security,
	}
	if len(args) > 0 && len(args[0]) > 0 {
		endpoint.Consume = strings.Split(args[0], ",")
	}
	if len(args) > 1 && len(args[1]) > 0 {
		endpoint.Produce = strings.Split(args[1], ",")
	}
	return endpoint
}

func (e Endpoint) AddResponses(responses ...Response) Endpoint {
	e.Responses = append(e.Responses, responses...)
	return e
}

func NewResponse(code string, description string, body interface{}) Response {
	return Response{
		Code:        code,
		Description: description,
		Body:        body,
	}
}

var NoParam []Parameter

func Params(params ...Parameter) (paramsArr []Parameter) {
	paramsArr = append(paramsArr, params...)
	return
}

func newParam(name string, t string, in string, required bool, description string, args ...Fields) (parameter Parameter) {
	if len(args) == 0 {
		parameter = Parameter{
			Name:        name,
			Type:        t,
			In:          in,
			Required:    required,
			Description: description,
		}
	} else {
		paramArgs := args[0]
		parameter = Parameter{
			Name:             name,
			Type:             t,
			In:               in,
			Required:         required,
			Description:      description,
			Format:           paramArgs.Format,
			Default:          paramArgs.Default,
			Min:              paramArgs.Min,
			Max:              paramArgs.Max,
			MinLen:           paramArgs.MinLen,
			MaxLen:           paramArgs.MaxLen,
			Pattern:          paramArgs.Pattern,
			MaxItems:         paramArgs.MaxItems,
			MinItems:         paramArgs.MinItems,
			UniqueItems:      paramArgs.UniqueItems,
			MultipleOf:       paramArgs.MultipleOf,
			CollectionFormat: paramArgs.CollectionFormat,
		}
	}
	generateParamDescription(&parameter)
	return
}

func generateParamDescription(param *Parameter) {
	newDescription := ""
	if param.Min != 0 {
		newDescription += "min: " + fmt.Sprint(param.Min) + " "
	}
	if param.Max != 0 {
		newDescription += "max: " + fmt.Sprint(param.Max) + " "
	}
	if param.MinLen != 0 {
		newDescription += "minLength: " + fmt.Sprint(param.MinLen) + " "
	}
	if param.MaxLen != 0 {
		newDescription += "maxLength: " + fmt.Sprint(param.MaxLen) + " "
	}
	if param.Pattern != "" {
		newDescription += "pattern: " + param.Pattern + " "
	}
	if len(newDescription) > 0 {
		if len(param.Description) > 0 {
			param.Description += "\n"
		}
		param.Description += " (" + strings.Trim(newDescription, " ") + ")"
	}
}

// args: name, required, description, format(optional)
func IntParam(name string, required bool, description string, args ...Fields) Parameter {
	return newParam(name, "integer", "path", required, description, args...)
}

// args: name, required, description, format(optional)
func StrParam(name string, required bool, description string, args ...Fields) Parameter {
	return newParam(name, "string", "path", required, description, args...)
}

// args: name, required, description, format(optional)
func BoolParam(name string, required bool, description string, args ...Fields) Parameter {
	return newParam(name, "boolean", "path", required, description, args...)
}

// args: name, required, description, format(optional)
func FileParam(name string, required bool, description string, args ...Fields) Parameter {
	return newParam(name, "file", "formData", required, description, args...)
}

// args: name, required, description, format(optional)
func IntQuery(name string, required bool, description string, args ...Fields) Parameter {
	param := IntParam(name, required, description, args...)
	param.In = "query"
	return param
}

// args: name, required, description, format(optional)
func StrQuery(name string, required bool, description string, args ...Fields) Parameter {
	param := StrParam(name, required, description, args...)
	param.In = "query"
	return param
}

// args: name, required, description, format(optional)
func BoolQuery(name string, required bool, description string, args ...Fields) Parameter {
	param := BoolParam(name, required, description, args...)
	param.In = "query"
	return param
}

// args: name, required, description, format(optional)
func IntHeader(name string, required bool, description string, args ...Fields) Parameter {
	param := IntParam(name, required, description, args...)
	param.In = "header"
	return param
}

// args: name, required, description, format(optional)
func StrHeader(name string, required bool, description string, args ...Fields) Parameter {
	param := StrParam(name, required, description, args...)
	param.In = "header"
	return param
}

// args: name, required, description, format(optional)
func BoolHeader(name string, required bool, description string, args ...Fields) Parameter {
	param := BoolParam(name, required, description, args...)
	param.In = "header"
	return param
}

// args: name, array, required, description, format(optional)
func IntEnumParam(name string, arr []int64, required bool, description string, args ...Fields) Parameter {
	param := newParam(name, "integer", "path", required, description, args...)
	param.Type = "integer"
	if len(arr) > 0 {
		s := make([]interface{}, len(arr))
		for i, v := range arr {
			s[i] = v
		}
		param.Enum = s
	}
	return param
}

// args: name, array, required, description, format(optional)
func StrEnumParam(name string, arr []string, required bool, description string, args ...Fields) Parameter {
	param := newParam(name, "string", "path", required, description, args...)
	param.Type = "string"
	if len(arr) > 0 {
		s := make([]interface{}, len(arr))
		for i, v := range arr {
			s[i] = v
		}
		param.Enum = s
	}
	return param
}

// args: name, array, required, description, format(optional)
func IntEnumQuery(name string, arr []int64, required bool, description string, args ...Fields) Parameter {
	param := IntEnumParam(name, arr, required, description, args...)
	param.In = "query"
	return param
}

// args: name, array, required, description, format(optional)
func StrEnumQuery(name string, arr []string, required bool, description string, args ...Fields) Parameter {
	param := StrEnumParam(name, arr, required, description, args...)
	param.In = "query"
	return param
}

// args: name, array, required, description, format(optional)
func IntEnumHeader(name string, arr []int64, required bool, description string, args ...Fields) Parameter {
	param := IntEnumParam(name, arr, required, description, args...)
	param.In = "header"
	return param
}

// args: name, array, required, description, format(optional)
func StrEnumHeader(name string, arr []string, required bool, description string, args ...Fields) Parameter {
	param := StrEnumParam(name, arr, required, description, args...)
	param.In = "header"
	return param
}

// args: name, array, required, description, format(optional)
func IntArrParam(name string, arr []int64, required bool, description string, args ...Fields) Parameter {
	param := newParam(name, "integer", "path", required, description, args...)
	param.Type = "array"
	param.Items = &ParameterItems{}
	param.Items.Type = "integer"
	param.Items.Format = param.Format
	if len(arr) > 0 {
		s := make([]interface{}, len(arr))
		for i, v := range arr {
			s[i] = v
		}
		param.Items.Enum = s
	}
	fillItemParams(&param)
	return param
}

// args: name, array, required, description, format(optional)
func StrArrParam(name string, arr []string, required bool, description string, args ...Fields) Parameter {
	param := newParam(name, "string", "path", required, description, args...)
	param.Type = "array"
	param.Items = &ParameterItems{}
	param.Items.Type = "string"
	param.Items.Format = param.Format
	if len(arr) > 0 {
		s := make([]interface{}, len(arr))
		for i, v := range arr {
			s[i] = v
		}
		param.Items.Enum = s
	}
	fillItemParams(&param)
	return param
}

// args: name, array, required, description, format(optional)
func IntArrQuery(name string, arr []int64, required bool, description string, args ...Fields) Parameter {
	param := IntArrParam(name, arr, required, description, args...)
	param.In = "query"
	return param
}

// args: name, array, required, description, format(optional)
func StrArrQuery(name string, arr []string, required bool, description string, args ...Fields) Parameter {
	param := StrArrParam(name, arr, required, description, args...)
	param.In = "query"
	return param
}

// args: name, array, required, description, format(optional)
func IntArrHeader(name string, arr []int64, required bool, description string, args ...Fields) Parameter {
	param := IntArrParam(name, arr, required, description, args...)
	param.In = "header"
	return param
}

// args: name, array, required, description, format(optional)
func StrArrHeader(name string, arr []string, required bool, description string, args ...Fields) Parameter {
	param := StrArrParam(name, arr, required, description, args...)
	param.In = "header"
	return param
}

func fillItemParams(param *Parameter) {
	param.Items.CollectionFormat = param.CollectionFormat
	param.Items.Default = param.Default
	param.Items.Format = param.Format
	param.Items.Max = param.Max
	param.Items.Min = param.Min
	param.Items.MaxLen = param.MaxLen
	param.Items.MinLen = param.MinLen
	param.Items.MultipleOf = param.MultipleOf
	param.Items.Pattern = param.Pattern
	param.Items.UniqueItems = param.UniqueItems
}

// Security

func BasicAuth() []map[string][]string {
	return []map[string][]string{
		{
			"basicAuth": []string{},
		},
	}
}

func ApiKeyAuth(name string) []map[string][]string {
	return []map[string][]string{
		{
			name: []string{},
		},
	}
}

func OAuth(name string, scopes ...string) []map[string][]string {
	return []map[string][]string{
		{
			name: scopes,
		},
	}
}

func Security(schemes ...[]map[string][]string) []map[string][]string {
	m := make([]map[string][]string, 0)
	for _, scheme := range schemes {
		m = append(m, scheme...)
	}
	return m
}
