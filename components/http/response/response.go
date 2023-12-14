package response

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-swagno/swagno/components/definition"
	"github.com/go-swagno/swagno/components/parameter"
)

// Response is an interface for response information.
type Response interface {
	Description() string
	ReturnCode() string
}

// ResponseGenerator is a struct that provides functionality to generate response schemas.
type ResponseGenerator struct{}

// NewResponseGenerator creates a new instance of ResponseGenerator.
func NewResponseGenerator() *ResponseGenerator {
	return &ResponseGenerator{}
}

// Generate generates a JSON response schema based on the provided model.
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
					Ref: strings.ReplaceAll(fmt.Sprintf("#/definitions/%s", reflect.TypeOf(model).Elem().String()), "[]", ""),
				},
			}
		} else {
			return &parameter.JsonResponseSchema{
				Type: "array",
				Items: &parameter.JsonResponseSchemeItems{
					Type: definition.Type(sliceElementKind.String()),
				},
			}
		}

	case reflect.Map:
		ref := strings.ReplaceAll(fmt.Sprintf("#/definitions/%T", model), "[]", "")
		return &parameter.JsonResponseSchema{
			Ref: ref,
		}

	default:
		if g.hasStructFields(model) {
			return &parameter.JsonResponseSchema{
				Ref: strings.ReplaceAll(fmt.Sprintf("#/definitions/%T", model), "[]", ""),
			}
		}
	}

	return nil
}

func (g ResponseGenerator) hasStructFields(s interface{}) bool {
	rv := reflect.ValueOf(s)

	if rv.Kind() != reflect.Struct {
		return false
	}

	numFields := rv.NumField()
	return numFields > 0
}
