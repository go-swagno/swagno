package swagger

import (
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

type Parameter struct {
	Name        string          `json:"name"`
	Type        string          `json:"type"`
	Query       bool            `json:"query"`
	Required    bool            `json:"required"`
	Description string          `json:"description"`
	Format      string          `json:"format,omitempty"`
	Items       *ParameterItems `json:"items,omitempty"`
}

type ParameterItems struct {
	Type   string `json:"type"`
	Format string `json:"format,omitempty"`
}

type Endpoint struct {
	Method      string      `json:"method"`
	Path        string      `json:"path"`
	Params      []Parameter `json:"params"`
	Tags        []string    `json:"tags"`
	Return      interface{} `json:"return"`
	Error       interface{} `json:"error"`
	Body        interface{} `json:"body"`
	Description string      `json:"description"`
}

func EndPoint(method MethodType, path string, tags string, params []Parameter, body interface{}, ret interface{}, err interface{}, des string) Endpoint {
	removedSpace := strings.ReplaceAll(tags, " ", "")
	return Endpoint{
		Method:      string(method),
		Path:        path,
		Tags:        strings.Split(removedSpace, ","),
		Params:      params,
		Return:      ret,
		Body:        body,
		Error:       err,
		Description: des,
	}
}

func Params(params ...Parameter) (paramsArr []Parameter) {
	paramsArr = append(paramsArr, params...)
	return
}

func NewParam(name string, t string, query bool, required bool, description string, args ...string) Parameter {
	if len(args) == 0 {
		return Parameter{
			Name:        name,
			Type:        t,
			Query:       query,
			Required:    required,
			Description: description,
		}
	} else {
		return Parameter{
			Name:        name,
			Type:        t,
			Query:       query,
			Required:    required,
			Description: description,
			Format:      args[0],
		}
	}
}

func NewArrayParam(name string, t string, arr []interface{}, query bool, required bool, description string, args ...string) Parameter {
	param := NewParam(name, t, query, required, description, args...)
	param.Type = "array"
	param.Items.Type = t
	param.Items.Format = param.Format
	return param
}

// args: name, required, description, format(optional)
func IntParam(name string, required bool, description string, args ...string) Parameter {
	return NewParam(name, "int", false, required, description, args...)
}

// args: name, required, description, format(optional)
func StrParam(name string, required bool, description string, args ...string) Parameter {
	return NewParam(name, "string", false, required, description, args...)
}

// args: name, required, description, format(optional)
func BoolParam(name string, required bool, description string, args ...string) Parameter {
	return NewParam(name, "bool", false, required, description, args...)
}

// args: name, required, description, format(optional)
func IntQuery(name string, required bool, description string, args ...string) Parameter {
	param := IntParam(name, required, description, args...)
	param.Query = true
	return param
}

// args: name, required, description, format(optional)
func StrQuery(name string, required bool, description string, args ...string) Parameter {
	param := StrParam(name, required, description, args...)
	param.Query = true
	return param
}

// args: name, required, description, format(optional)
func BoolQuery(name string, required bool, description string, args ...string) Parameter {
	param := BoolParam(name, required, description, args...)
	param.Query = true
	return param
}

// args: name, array, required, description, format(optional)
func IntArrQuery(name string, arr []int64, required bool, description string, args ...string) Parameter {
	param := NewParam(name, "int", true, required, description, args...)
	param.Type = "array"
	param.Items = &ParameterItems{}
	param.Items.Type = "integer"
	param.Items.Format = param.Format
	return param
}

// args: name, array, required, description, format(optional)
func StrArrQuery(name string, arr []string, required bool, description string, args ...string) Parameter {
	param := NewParam(name, "string", true, required, description, args...)
	param.Type = "array"
	param.Items = &ParameterItems{}
	param.Items.Type = "string"
	param.Items.Format = param.Format
	return param
}
