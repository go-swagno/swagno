package parameter

import (
	"fmt"
	"strings"
)

// CollectionFormat defines the format for serializing array parameters in the URL query string.
type CollectionFormat string

const (
	CSV   CollectionFormat = "csv"
	SSV   CollectionFormat = "ssv"
	TSV   CollectionFormat = "tsv"
	Pipes CollectionFormat = "pipes"
	Multi CollectionFormat = "multi"
)

// String returns the string representation of the CollectionFormat.
func (c CollectionFormat) String() string {
	return string(c)
}

// Location specifies where in the request a parameter is expected to be located.
type Location string

// String returns the string representation of the Location.
func (l Location) String() string {
	return string(l)
}

const (
	Query  Location = "query"
	Header Location = "header"
	Path   Location = "path"
	Form   Location = "formData"
)

// ParamType represents the type of a parameter in the API endpoint.
type ParamType string

// String returns the string representation of the ParamType.
func (p ParamType) String() string {
	return string(p)
}

const (
	String  ParamType = "string"
	Number  ParamType = "number"
	Integer ParamType = "integer"
	Boolean ParamType = "boolean"
	Array   ParamType = "array"
	File    ParamType = "file"
)

// JsonParameter is the JSON model version of Parameter object used for API purposes
// https://swagger.io/specification/v2/#parameterObject
type JsonParameter struct {
	Type              string              `json:"type,omitempty"`
	Description       string              `json:"description"`
	Name              string              `json:"name"`
	In                string              `json:"in,omitempty"`
	Required          bool                `json:"required"`
	Schema            *JsonResponseSchema `json:"schema,omitempty"`
	Format            string              `json:"format,omitempty"`
	Enum              []interface{}       `json:"enum,omitempty"`
	Default           interface{}         `json:"default,omitempty"`
	Min               int64               `json:"minimum,omitempty"`
	Max               int64               `json:"maximum,omitempty"`
	MinLen            int64               `json:"minLength,omitempty"`
	MaxLen            int64               `json:"maxLength,omitempty"`
	Pattern           string              `json:"pattern,omitempty"`
	MaxItems          int64               `json:"maxItems,omitempty"`
	MinItems          int64               `json:"minItems,omitempty"`
	UniqueItems       bool                `json:"uniqueItems,omitempty"`
	MultipleOf        int64               `json:"multipleOf,omitempty"`
	CollenctionFormat string              `json:"collectionFormat,omitempty"`
}

// JsonResponseSchema defines the schema for a JSON response as per the Swagger 2.0 specification.
// It is used to describe the structure and type of a response returned by an API endpoint.
// https://swagger.io/specification/v2/#schema-object
type JsonResponseSchema struct {
	Ref   string                   `json:"$ref,omitempty"`
	Type  string                   `json:"type,omitempty"`
	Items *JsonResponseSchemeItems `json:"items,omitempty"`
}

// JsonResponseSchemeItems represents the individual items in a JsonResponseSchema, especially for arrays.
// It provides the type or reference for the array items.
type JsonResponseSchemeItems struct {
	Type  string                   `json:"type,omitempty"`
	Ref   string                   `json:"$ref,omitempty"`
	Items *JsonResponseSchemeItems `json:"items,omitempty"`
}

// Parameter represents a parameter in an API endpoint.
type Parameter struct {
	name             string
	typeValue        ParamType
	in               Location
	required         bool
	description      string
	enum             []interface{}
	defaultValue     interface{}
	format           string
	min              int64
	max              int64
	minLen           int64
	maxLen           int64
	pattern          string
	maxItems         int64
	minItems         int64
	uniqueItems      bool
	multipleOf       int64
	collectionFormat CollectionFormat
}

// Location returns the location of the parameter (i.e. Query, Body, Path, and etc.)
func (p Parameter) Location() Location {
	return p.in
}

// AsJson returns the json representation of Parameter
func (p *Parameter) AsJson() JsonParameter {
	return JsonParameter{
		Name:              p.name,
		In:                p.in.String(),
		Description:       p.description,
		Required:          p.required,
		Type:              p.typeValue.String(),
		Format:            p.format,
		Enum:              p.enum,
		Default:           p.defaultValue,
		Min:               p.min,
		Max:               p.max,
		MinLen:            p.minLen,
		MaxLen:            p.maxLen,
		Pattern:           p.pattern,
		MaxItems:          p.maxItems,
		MinItems:          p.minItems,
		UniqueItems:       p.uniqueItems,
		MultipleOf:        p.multipleOf,
		CollenctionFormat: p.collectionFormat.String(),
	}
}

// NoParam is an empty slice of parameters.
var NoParam []Parameter

// Params appends parameters to an existing parameter slice.
func Params(params ...Parameter) (paramsArr []Parameter) {
	paramsArr = append(paramsArr, params...)
	return
}

// IntParam creates an integer parameter.
func IntParam(name string, l Location, opts ...Option) *Parameter {
	opts = append(opts, WithType("integer"), WithIn(l))
	return newParam(name, opts...)
}

// StrParam creates a string parameter.
func StrParam(name string, l Location, opts ...Option) *Parameter {
	opts = append(opts, WithType(String), WithIn(l))
	return newParam(name, opts...)
}

// BoolParam creates a boolean parameter.
func BoolParam(name string, l Location, opts ...Option) *Parameter {
	opts = append(opts, WithType("boolean"), WithIn(l))
	return newParam(name, opts...)
}

// FileParam creates a file parameter.
func FileParam(name string, opts ...Option) *Parameter {
	opts = append(opts, WithType("file"), WithIn(Form))
	return newParam(name, opts...)
}

// IntEnumParam creates an integer enum parameter.
func IntEnumParam(name string, l Location, arr []int64, opts ...Option) *Parameter {
	opts = append(opts, WithType(Integer), WithIn(l))
	param := newParam(name, opts...)

	if len(arr) > 0 {
		s := make([]interface{}, len(arr))
		for i, v := range arr {
			s[i] = v
		}
		param.enum = s
	}

	return param
}

// StrEnumParam creates a string enum parameter.
func StrEnumParam(name string, l Location, arr []string, opts ...Option) *Parameter {
	opts = append(opts, WithType(String), WithIn(l))
	param := newParam(name, opts...)

	if len(arr) > 0 {
		s := make([]interface{}, len(arr))
		for i, v := range arr {
			s[i] = v
		}
		param.enum = s
	}

	return param
}

// IntArrParam creates an integer array parameter.
func IntArrParam(name string, l Location, arr []int64, opts ...Option) *Parameter {
	opts = append(opts, WithType(Array), WithIn(l))
	param := newParam(name, opts...)

	if len(arr) > 0 {
		s := make([]interface{}, len(arr))
		for i, v := range arr {
			s[i] = v
		}
		param.enum = s
	}

	return param
}

// StrArrParam creates a string array parameter.
func StrArrParam(name string, l Location, arr []string, opts ...Option) *Parameter {
	opts = append(opts, WithType(Array), WithIn(Location(l)))
	param := newParam(name, opts...)

	if len(arr) > 0 {
		s := make([]interface{}, len(arr))
		for i, v := range arr {
			s[i] = v
		}
		param.enum = s
	}

	return param
}

// Option represents a function that can modify a Parameter.
type Option func(*Parameter)

// WithType sets the Type field of a Parameter.
func WithType(t ParamType) Option {
	return func(p *Parameter) {
		p.typeValue = t
	}
}

// WithIn sets the In field of a Parameter.
func WithIn(in Location) Option {
	return func(p *Parameter) {
		p.in = in
	}
}

// WithRequired sets the Required field of a Parameter.
func WithRequired() Option {
	return func(p *Parameter) {
		p.required = true
	}
}

// WithDescription sets the Description field of a Parameter.
func WithDescription(description string) Option {
	return func(p *Parameter) {
		p.description = description
	}
}

// WithDefault sets the Default field of a Parameter.
func WithDefault(defaultValue interface{}) Option {
	return func(p *Parameter) {
		p.defaultValue = defaultValue
	}
}

// WithFormat sets the Format field of a Parameter.
func WithFormat(format string) Option {
	return func(p *Parameter) {
		p.format = format
	}
}

// WithMin sets the Min field of a Parameter.
func WithMin(min int64) Option {
	return func(p *Parameter) {
		p.min = min
	}
}

// WithMax sets the Max field of a Parameter.
func WithMax(max int64) Option {
	return func(p *Parameter) {
		p.max = max
	}
}

// WithMinLen sets the MinLen field of a Parameter.
func WithMinLen(minLen int64) Option {
	return func(p *Parameter) {
		p.minLen = minLen
	}
}

// WithMaxLen sets the MaxLen field of a Parameter.
func WithMaxLen(maxLen int64) Option {
	return func(p *Parameter) {
		p.maxLen = maxLen
	}
}

// WithPattern sets the Pattern field of a Parameter.
func WithPattern(pattern string) Option {
	return func(p *Parameter) {
		p.pattern = pattern
	}
}

// WithMaxItems sets the MaxItems field of a Parameter.
func WithMaxItems(maxItems int64) Option {
	return func(p *Parameter) {
		p.maxItems = maxItems
	}
}

// WithMinItems sets the MinItems field of a Parameter.
func WithMinItems(minItems int64) Option {
	return func(p *Parameter) {
		p.minItems = minItems
	}
}

// WithUniqueItems sets the UniqueItems field of a Parameter.
func WithUniqueItems(uniqueItems bool) Option {
	return func(p *Parameter) {
		p.uniqueItems = uniqueItems
	}
}

// WithMultipleOf sets the MultipleOf field of a Parameter.
func WithMultipleOf(multipleOf int64) Option {
	return func(p *Parameter) {
		p.multipleOf = multipleOf
	}
}

// WithCollectionFormat sets the CollectionFormat field of a Parameter.
func WithCollectionFormat(c CollectionFormat) Option {
	return func(p *Parameter) {
		p.collectionFormat = c
	}
}

// newParam creates a newParam parameter with the given options.
func newParam(name string, opts ...Option) *Parameter {
	parameter := Parameter{name: name}

	for _, opt := range opts {
		opt(&parameter)
	}

	generateParamDescription(&parameter)
	return &parameter
}

// generateParamDescription generates the description for a parameter based on its properties.
func generateParamDescription(param *Parameter) {
	newDescription := ""
	if param.min != 0 {
		newDescription += "min: " + fmt.Sprint(param.min) + " "
	}
	if param.max != 0 {
		newDescription += "max: " + fmt.Sprint(param.max) + " "
	}
	if param.minLen != 0 {
		newDescription += "minLength: " + fmt.Sprint(param.minLen) + " "
	}
	if param.maxLen != 0 {
		newDescription += "maxLength: " + fmt.Sprint(param.maxLen) + " "
	}
	if param.pattern != "" {
		newDescription += "pattern: " + param.pattern + " "
	}
	if len(newDescription) > 0 {
		if len(param.description) > 0 {
			param.description += "\n"
		}
		param.description += " (" + strings.Trim(newDescription, " ") + ")"
	}
}
