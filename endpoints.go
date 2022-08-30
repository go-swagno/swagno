package swagger

import (
	"strings"
)

type MethodType string

const (
	GET    MethodType = "GET"
	POST   MethodType = "POST"
	PUT    MethodType = "PUT"
	DELETE MethodType = "DELETE"
)

type Parameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Query       bool   `json:"query"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Format      string `json:"format,omitempty"`
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

func IntParam(name string, required bool, description string, args ...string) Parameter {
	return NewParam(name, "int", false, required, description, args...)
}
func StrParam(name string, required bool, description string, args ...string) Parameter {
	return NewParam(name, "string", false, required, description, args...)
}
func BoolParam(name string, required bool, description string, args ...string) Parameter {
	return NewParam(name, "bool", false, required, description, args...)
}
func IntQuery(name string, required bool, description string, args ...string) Parameter {
	param := IntParam(name, required, description, args...)
	param.Query = true
	return param
}
func StrQuery(name string, required bool, description string, args ...string) Parameter {
	param := StrParam(name, required, description, args...)
	param.Query = true
	return param
}
func BoolQuery(name string, required bool, description string, args ...string) Parameter {
	param := BoolParam(name, required, description, args...)
	param.Query = true
	return param
}
