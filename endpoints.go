package swagno

import (
	"fmt"
	"strings"
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

// Parameter represents a parameter in an API endpoint.
type Parameter struct {
	Name              string          `json:"name"`
	Type              string          `json:"type"`
	In                string          `json:"in"`
	Required          bool            `json:"required"`
	Description       string          `json:"description"`
	Enum              []interface{}   `json:"enum,omitempty"`
	Items             *ParameterItems `json:"items,omitempty"`
	Default           interface{}     `json:"default,omitempty"`
	Format            string          `json:"format,omitempty"`
	Min               int64           `json:"minimum,omitempty"`
	Max               int64           `json:"maximum,omitempty"`
	MinLen            int64           `json:"minLength,omitempty"`
	MaxLen            int64           `json:"maxLength,omitempty"`
	Pattern           string          `json:"pattern,omitempty"`
	MaxItems          int64           `json:"maxItems,omitempty"`
	MinItems          int64           `json:"minItems,omitempty"`
	UniqueItems       bool            `json:"uniqueItems,omitempty"`
	MultipleOf        int64           `json:"multipleOf,omitempty"`
	CollenctionFormat string          `json:"collectionFormat,omitempty"`
}

// ParameterItems represents items within a parameter (used for array types).
type ParameterItems struct {
	Type              string        `json:"type"`
	Format            string        `json:"format,omitempty"`
	Enum              []interface{} `json:"enum,omitempty"`
	Default           interface{}   `json:"default,omitempty"`
	Min               int64         `json:"minimum,omitempty"`
	Max               int64         `json:"maximum,omitempty"`
	MinLen            int64         `json:"minLength,omitempty"`
	MaxLen            int64         `json:"maxLength,omitempty"`
	Pattern           string        `json:"pattern,omitempty"`
	MaxItems          int64         `json:"maxItems,omitempty"`
	MinItems          int64         `json:"minItems,omitempty"`
	UniqueItems       bool          `json:"uniqueItems,omitempty"`
	MultipleOf        int64         `json:"multipleOf,omitempty"`
	CollenctionFormat string        `json:"collectionFormat,omitempty"`
}

// Fields represents fields within a parameter or response object.
type Fields struct {
	Default           interface{} `json:"default,omitempty"`
	Format            string      `json:"format,omitempty"`
	Min               int64       `json:"minimum,omitempty"`
	Max               int64       `json:"maximum,omitempty"`
	MinLen            int64       `json:"minLength,omitempty"`
	MaxLen            int64       `json:"maxLength,omitempty"`
	Pattern           string      `json:"pattern,omitempty"`
	MaxItems          int64       `json:"maxItems,omitempty"`
	MinItems          int64       `json:"minItems,omitempty"`
	UniqueItems       bool        `json:"uniqueItems,omitempty"`
	MultipleOf        int64       `json:"multipleOf,omitempty"`
	CollenctionFormat string      `json:"collectionFormat,omitempty"`
}

// Endpoint represents an API endpoint.
type Endpoint struct {
	Method      string                `json:"method"`
	Path        string                `json:"path"`
	Params      []Parameter           `json:"params"`
	Tags        []string              `json:"tags"`
	Return      interface{}           `json:"return"` // TODO after learning more about these types should model them out to be a struct or a better defined interfaced
	Error       interface{}           `json:"error"`  // TODO same as above
	Body        interface{}           `json:"body"`   // TODO same as above
	Description string                `json:"description"`
	Summary     string                `json:"summary"`
	Consume     []string              `json:"consume"`
	Produce     []string              `json:"produce"`
	Security    []map[string][]string `json:"security"`
}

type EndPoints struct {
	List []Endpoint `json:"Endpoints"`
}

// ResponseInfo is an interface for response information.
type ResponseInfo interface {
	GetDescription() string
	GetReturnCode() string
}

// ErrorResponses is an interface for error responses.
type ErrorResponses interface {
	GetErrors() []ResponseInfo
}

// TODO use a better creation pattern to make these objects
// EndPoint is a function to create an API endpoint.
// args: method, path, tags, params, body, return, error, description, security, consume, produce
func EndPoint(method MethodType, path string, tags string, params []Parameter, body interface{}, ret interface{}, err interface{}, des string, security []map[string][]string, args ...string) Endpoint {
	removedSpace := strings.ReplaceAll(tags, " ", "")
	endpoint := Endpoint{
		Method:      string(method),
		Path:        path,
		Tags:        strings.Split(removedSpace, ","),
		Params:      params,
		Return:      ret,
		Body:        body,
		Error:       err,
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

// NoParam is an empty slice of parameters.
var NoParam []Parameter

// Params appends parameters to an existing parameter slice.
func Params(params ...Parameter) (paramsArr []Parameter) {
	paramsArr = append(paramsArr, params...)
	return
}

// IntParam creates an integer parameter.
// args: name, required, description, format(optional)
func IntParam(name string, required bool, description string, args ...Fields) Parameter {
	return newParam(name, "integer", "path", required, description, args...)
}

// StrParam creates a string parameter.
// args: name, required, description, format(optional)
func StrParam(name string, required bool, description string, args ...Fields) Parameter {
	return newParam(name, "string", "path", required, description, args...)
}

// BoolParam creates a boolean parameter.
// args: name, required, description, format(optional)
func BoolParam(name string, required bool, description string, args ...Fields) Parameter {
	return newParam(name, "boolean", "path", required, description, args...)
}

// FileParam creates a file parameter.
// args: name, required, description, format(optional)
func FileParam(name string, required bool, description string, args ...Fields) Parameter {
	return newParam(name, "file", "formData", required, description, args...)
}

// IntQuery creates an integer query parameter.
// args: name, required, description, format(optional)
func IntQuery(name string, required bool, description string, args ...Fields) Parameter {
	param := IntParam(name, required, description, args...)
	param.In = "query"
	return param
}

// StrQuery creates a string query parameter.
// args: name, required, description, format(optional)
func StrQuery(name string, required bool, description string, args ...Fields) Parameter {
	param := StrParam(name, required, description, args...)
	param.In = "query"
	return param
}

// BoolQuery creates a boolean query parameter.
// args: name, required, description, format(optional)
func BoolQuery(name string, required bool, description string, args ...Fields) Parameter {
	param := BoolParam(name, required, description, args...)
	param.In = "query"
	return param
}

// IntHeader creates an integer header parameter.
// args: name, required, description, format(optional)
func IntHeader(name string, required bool, description string, args ...Fields) Parameter {
	param := IntParam(name, required, description, args...)
	param.In = "header"
	return param
}

// StrHeader creates a string header parameter.
// args: name, required, description, format(optional)
func StrHeader(name string, required bool, description string, args ...Fields) Parameter {
	param := StrParam(name, required, description, args...)
	param.In = "header"
	return param
}

// BoolHeader creates a boolean header parameter.
// args: name, required, description, format(optional)
func BoolHeader(name string, required bool, description string, args ...Fields) Parameter {
	param := BoolParam(name, required, description, args...)
	param.In = "header"
	return param
}

// IntEnumParam creates an integer enum parameter.
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

// StrEnumParam creates a string enum parameter.
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

// IntEnumQuery creates an integer enum query parameter.
// args: name, array, required, description, format(optional)
func IntEnumQuery(name string, arr []int64, required bool, description string, args ...Fields) Parameter {
	param := IntEnumParam(name, arr, required, description, args...)
	param.In = "query"
	return param
}

// StrEnumQuery creates a string enum query parameter.
// args: name, array, required, description, format(optional)
func StrEnumQuery(name string, arr []string, required bool, description string, args ...Fields) Parameter {
	param := StrEnumParam(name, arr, required, description, args...)
	param.In = "query"
	return param
}

// IntEnumHeader creates an integer enum header parameter.
// args: name, array, required, description, format(optional)
func IntEnumHeader(name string, arr []int64, required bool, description string, args ...Fields) Parameter {
	param := IntEnumParam(name, arr, required, description, args...)
	param.In = "header"
	return param
}

// StrEnumHeader creates a string enum header parameter.
// args: name, array, required, description, format(optional)
func StrEnumHeader(name string, arr []string, required bool, description string, args ...Fields) Parameter {
	param := StrEnumParam(name, arr, required, description, args...)
	param.In = "header"
	return param
}

// IntArrParam creates an integer array parameter.
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

// StrArrParam creates a string array parameter.
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

// IntArrQuery creates an integer array query parameter.
// args: name, array, required, description, format(optional)
func IntArrQuery(name string, arr []int64, required bool, description string, args ...Fields) Parameter {
	param := IntArrParam(name, arr, required, description, args...)
	param.In = "query"
	return param
}

// StrArrQuery creates a string array query parameter.
// args: name, array, required, description, format(optional)
func StrArrQuery(name string, arr []string, required bool, description string, args ...Fields) Parameter {
	param := StrArrParam(name, arr, required, description, args...)
	param.In = "query"
	return param
}

// IntArrHeader creates an integer array header parameter.
// args: name, array, required, description, format(optional)
func IntArrHeader(name string, arr []int64, required bool, description string, args ...Fields) Parameter {
	param := IntArrParam(name, arr, required, description, args...)
	param.In = "header"
	return param
}

// StrArrHeader creates a string array header parameter.
// args: name, array, required, description, format(optional)
func StrArrHeader(name string, arr []string, required bool, description string, args ...Fields) Parameter {
	param := StrArrParam(name, arr, required, description, args...)
	param.In = "header"
	return param
}

// newParam creates a new parameter with the given arguments.
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
			Name:              name,
			Type:              t,
			In:                in,
			Required:          required,
			Description:       description,
			Format:            paramArgs.Format,
			Default:           paramArgs.Default,
			Min:               paramArgs.Min,
			Max:               paramArgs.Max,
			MinLen:            paramArgs.MinLen,
			MaxLen:            paramArgs.MaxLen,
			Pattern:           paramArgs.Pattern,
			MaxItems:          paramArgs.MaxItems,
			MinItems:          paramArgs.MinItems,
			UniqueItems:       paramArgs.UniqueItems,
			MultipleOf:        paramArgs.MultipleOf,
			CollenctionFormat: paramArgs.CollenctionFormat,
		}
	}
	generateParamDescription(&parameter)
	return
}

// generateParamDescription generates the description for a parameter based on its properties.
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

// fillItemParams sets item properties for array parameters.
func fillItemParams(param *Parameter) {
	param.Items.CollenctionFormat = param.CollenctionFormat
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
