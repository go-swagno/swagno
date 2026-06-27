package response

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-swagno/swagno/v3/components/fields"
	"github.com/go-swagno/swagno/v3/components/parameter"
)

// Response is an interface for response information.
type Response interface {
	Description() string
	ReturnCode() string
}

// CustomResponse represents a custom implementation of Response.
type CustomResponse struct {
	Model             any
	returnCodeString  string
	descriptionString string
	example           interface{}
	examples          map[string]interface{}
}

// ResponseGenerator is a struct that provides functionality to generate response schemas.
type ResponseGenerator struct {
	// HidePackageName, when true, strips the leading package qualifier from $ref
	// values (e.g. "models.MyStruct" -> "MyStruct").
	HidePackageName bool
}

// New creates a new instance of Response with the provided model return code, and description.
func New(model any, returnCode string, description string) CustomResponse {
	return CustomResponse{
		Model:             model,
		returnCodeString:  returnCode,
		descriptionString: description,
	}
}

// WithExample sets an example for the response and clears multiple examples.
func (c CustomResponse) WithExample(example interface{}) CustomResponse {
	c.example = example
	c.examples = nil
	return c
}

// WithExamples sets multiple examples for the response and clears a single example.
func (c CustomResponse) WithExamples(examples map[string]interface{}) CustomResponse {
	c.examples = examples
	c.example = nil
	return c
}

// Example returns the example for the response.
func (c CustomResponse) Example() interface{} {
	return c.example
}

// Examples returns the examples for the response.
func (c CustomResponse) Examples() map[string]interface{} {
	return c.examples
}

// NewResponseGenerator creates a new instance of ResponseGenerator.
func NewResponseGenerator(hidePackageName bool) *ResponseGenerator {
	return &ResponseGenerator{
		HidePackageName: hidePackageName,
	}
}

// Generate generates a JSON response schema based on the provided model for OpenAPI 3.0.
// It uses reflection to determine the type of the model and constructs the appropriate JSON schema.
// This function handles different types such as slices, maps, and structures to create a detailed and accurate schema.
func (g ResponseGenerator) Generate(model any) *parameter.JsonResponseSchema {
	switch reflect.TypeOf(model).Kind() {
	case reflect.Slice:
		sliceElementKind := reflect.TypeOf(model).Elem().Kind()
		if sliceElementKind == reflect.Struct {
			return &parameter.JsonResponseSchema{
				Type: "array",
				Items: &parameter.JsonResponseSchemeItems{
					Ref: fmt.Sprintf("#/components/schemas/%s", fields.RefName(strings.ReplaceAll(reflect.TypeOf(model).Elem().String(), "[]", ""), g.HidePackageName)),
				},
			}
		} else {
			return &parameter.JsonResponseSchema{
				Type: "array",
				Items: &parameter.JsonResponseSchemeItems{
					Type: fields.Type(sliceElementKind.String()),
				},
			}
		}

	case reflect.Map:
		ref := fmt.Sprintf("#/components/schemas/%s", fields.RefName(strings.ReplaceAll(fmt.Sprintf("%T", model), "[]", ""), g.HidePackageName))
		return &parameter.JsonResponseSchema{
			Ref: ref,
		}

	default:
		if hasStructFields(model) {
			return &parameter.JsonResponseSchema{
				Ref: fmt.Sprintf("#/components/schemas/%s", fields.RefName(strings.ReplaceAll(fmt.Sprintf("%T", model), "[]", ""), g.HidePackageName)),
			}
		}
	}

	return nil
}

// hasStructFields checks if the given interface has fields in case it's a struct.
func hasStructFields(s interface{}) bool {
	rv := reflect.ValueOf(s)

	if rv.Kind() != reflect.Struct {
		return false
	}

	numFields := rv.NumField()
	return numFields > 0
}

func (c CustomResponse) Description() string {
	return c.descriptionString
}
func (c CustomResponse) ReturnCode() string {
	return c.returnCodeString
}
