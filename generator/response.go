package generator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-swagno/swagno/components/parameter"
	"github.com/go-swagno/swagno/utils"
)

type ResponseGenerator struct {
}

func NewResponseGenerator() *ResponseGenerator {
	return &ResponseGenerator{}
}

func (g ResponseGenerator) GenerateJsonResponseScheme(model any) *parameter.JsonResponseSchema {
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
					Type: getType(sliceElementKind.String()),
				},
			}
		}

	case reflect.Map:
		ref := strings.ReplaceAll(fmt.Sprintf("#/definitions/%T", model), "[]", "")
		// preventing override map definitions, generate random name for map if its not typed
		if reflect.TypeOf(model).Name() == "" {
			ref = strings.ReplaceAll(fmt.Sprintf("#/definitions/%T_%s", model, utils.GetHashOfMap(utils.ConvertInterfaceToMap(model))), "[]", "")
		}
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
