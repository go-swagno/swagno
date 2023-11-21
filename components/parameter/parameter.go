package parameter

import (
	"fmt"
	"strings"
)

// Parameter represents a parameter in an API endpoint.
type Parameter struct {
	Name             string        `json:"name"`
	Type             string        `json:"type"`
	In               string        `json:"in"` // TODO these types should be it's own type
	Required         bool          `json:"required"`
	Description      string        `json:"description"`
	Enum             []interface{} `json:"enum,omitempty"`
	Default          interface{}   `json:"default,omitempty"`
	Format           string        `json:"format,omitempty"`
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

// NoParam is an empty slice of parameters.
var NoParam []Parameter

// Params appends parameters to an existing parameter slice.
func Params(params ...Parameter) (paramsArr []Parameter) {
	paramsArr = append(paramsArr, params...)
	return
}

// IntParam creates an integer parameter.
func IntParam(name string, opts ...Option) Parameter {
	opts = append(opts, WithType("integer"), WithIn("path"))
	return newParam(name, opts...)
}

// StrParam creates a string parameter.
func StrParam(name string, opts ...Option) Parameter {
	opts = append(opts, WithType("string"), WithIn("path"))
	return newParam(name, opts...)
}

// BoolParam creates a boolean parameter.
func BoolParam(name string, opts ...Option) Parameter {
	opts = append(opts, WithType("boolean"), WithIn("path"))
	return newParam(name, opts...)
}

// FileParam creates a file parameter.
func FileParam(name string, opts ...Option) Parameter {
	opts = append(opts, WithType("file"), WithIn("formData"))
	return newParam(name, opts...)
}

// IntQuery creates an integer query parameter.
func IntQuery(name string, opts ...Option) Parameter {
	opts = append(opts, WithType("integer"), WithIn("query"))
	return newParam(name, opts...)
}

// StrQuery creates a string query parameter.
func StrQuery(name string, opts ...Option) Parameter {
	param := StrParam(name, opts...)
	param.In = "query"
	return param
}

// BoolQuery creates a boolean query parameter.
func BoolQuery(name string, opts ...Option) Parameter {
	param := BoolParam(name, opts...)
	param.In = "query"
	return param
}

// IntHeader creates an integer header parameter.
func IntHeader(name string, opts ...Option) Parameter {
	param := IntParam(name, opts...)
	param.In = "header"
	return param
}

// StrHeader creates a string header parameter.
func StrHeader(name string, opts ...Option) Parameter {
	param := StrParam(name, opts...)
	param.In = "header"
	return param
}

// BoolHeader creates a boolean header parameter.
func BoolHeader(name string, opts ...Option) Parameter {
	param := BoolParam(name, opts...)
	param.In = "header"
	return param
}

// IntEnumParam creates an integer enum parameter.
func IntEnumParam(name string, arr []int64, opts ...Option) Parameter {
	opts = append(opts, WithIn("path"), WithType("integer"))
	param := newParam(name, opts...)

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
func StrEnumParam(name string, arr []string, opts ...Option) Parameter {
	opts = append(opts, WithIn("path"), WithType("string"))
	param := newParam(name, opts...)

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
func IntEnumQuery(name string, arr []int64, opts ...Option) Parameter {
	param := IntEnumParam(name, arr, opts...)
	param.In = "query"
	return param
}

// StrEnumQuery creates a string enum query parameter.
func StrEnumQuery(name string, arr []string, opts ...Option) Parameter {
	param := StrEnumParam(name, arr, opts...)
	param.In = "query"
	return param
}

// IntEnumHeader creates an integer enum header parameter.
func IntEnumHeader(name string, arr []int64, opts ...Option) Parameter {
	param := IntEnumParam(name, arr, opts...)
	param.In = "header"
	return param
}

// StrEnumHeader creates a string enum header parameter.
func StrEnumHeader(name string, arr []string, opts ...Option) Parameter {
	param := StrEnumParam(name, arr, opts...)
	param.In = "header"
	return param
}

// IntArrParam creates an integer array parameter.
func IntArrParam(name string, arr []int64, opts ...Option) Parameter {
	opts = append(opts, WithType("integer"), WithIn("path"))
	param := newParam(name, opts...)
	param.Type = "array"

	if len(arr) > 0 {
		s := make([]interface{}, len(arr))
		for i, v := range arr {
			s[i] = v
		}
		param.Enum = s
	}

	return param
}

// StrArrParam creates a string array parameter.
func StrArrParam(name string, arr []string, opts ...Option) Parameter {
	opts = append(opts, WithType("string"), WithIn("path"))
	param := newParam(name, opts...)
	param.Type = "array"

	if len(arr) > 0 {
		s := make([]interface{}, len(arr))
		for i, v := range arr {
			s[i] = v
		}
		param.Enum = s
	}

	return param
}

// IntArrQuery creates an integer array query parameter.
func IntArrQuery(name string, arr []int64, opts ...Option) Parameter {
	param := IntArrParam(name, arr, opts...)
	param.In = "query"
	return param
}

// StrArrQuery creates a string array query parameter.
func StrArrQuery(name string, arr []string, opts ...Option) Parameter {
	param := StrArrParam(name, arr, opts...)
	param.In = "query"
	return param
}

// IntArrHeader creates an integer array header parameter.
func IntArrHeader(name string, arr []int64, opts ...Option) Parameter {
	param := IntArrParam(name, arr, opts...)
	param.In = "header"
	return param
}

// StrArrHeader creates a string array header parameter.
func StrArrHeader(name string, arr []string, opts ...Option) Parameter {
	param := StrArrParam(name, arr, opts...)
	param.In = "header"
	return param
}

// Option represents a function that can modify a Parameter.
type Option func(*Parameter)

// WithType sets the Type field of a Parameter.
func WithType(t string) Option {
	return func(p *Parameter) {
		p.Type = t
	}
}

// WithIn sets the In field of a Parameter.
func WithIn(in string) Option {
	return func(p *Parameter) {
		p.In = in
	}
}

// WithRequired sets the Required field of a Parameter.
func WithRequired(required bool) Option {
	return func(p *Parameter) {
		p.Required = required
	}
}

// WithDescription sets the Description field of a Parameter.
func WithDescription(description string) Option {
	return func(p *Parameter) {
		p.Description = description
	}
}

// WithDefault sets the Default field of a Parameter.
func WithDefault(defaultValue interface{}) Option {
	return func(p *Parameter) {
		p.Default = defaultValue
	}
}

// WithFormat sets the Format field of a Parameter.
func WithFormat(format string) Option {
	return func(p *Parameter) {
		p.Format = format
	}
}

// WithMin sets the Min field of a Parameter.
func WithMin(min int64) Option {
	return func(p *Parameter) {
		p.Min = min
	}
}

// WithMax sets the Max field of a Parameter.
func WithMax(max int64) Option {
	return func(p *Parameter) {
		p.Max = max
	}
}

// WithMinLen sets the MinLen field of a Parameter.
func WithMinLen(minLen int64) Option {
	return func(p *Parameter) {
		p.MinLen = minLen
	}
}

// WithMaxLen sets the MaxLen field of a Parameter.
func WithMaxLen(maxLen int64) Option {
	return func(p *Parameter) {
		p.MaxLen = maxLen
	}
}

// WithPattern sets the Pattern field of a Parameter.
func WithPattern(pattern string) Option {
	return func(p *Parameter) {
		p.Pattern = pattern
	}
}

// WithMaxItems sets the MaxItems field of a Parameter.
func WithMaxItems(maxItems int64) Option {
	return func(p *Parameter) {
		p.MaxItems = maxItems
	}
}

// WithMinItems sets the MinItems field of a Parameter.
func WithMinItems(minItems int64) Option {
	return func(p *Parameter) {
		p.MinItems = minItems
	}
}

// WithUniqueItems sets the UniqueItems field of a Parameter.
func WithUniqueItems(uniqueItems bool) Option {
	return func(p *Parameter) {
		p.UniqueItems = uniqueItems
	}
}

// WithMultipleOf sets the MultipleOf field of a Parameter.
func WithMultipleOf(multipleOf int64) Option {
	return func(p *Parameter) {
		p.MultipleOf = multipleOf
	}
}

// WithCollectionFormat sets the CollectionFormat field of a Parameter.
func WithCollectionFormat(collectionFormat string) Option {
	return func(p *Parameter) {
		p.CollectionFormat = collectionFormat
	}
}

// newParam creates a newParam parameter with the given options.
func newParam(name string, opts ...Option) Parameter {
	parameter := Parameter{Name: name}

	for _, opt := range opts {
		opt(&parameter)
	}

	generateParamDescription(&parameter)
	return parameter
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
