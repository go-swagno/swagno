package fields

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ExampleTag retrieves the 'example' struct tag's value and converts it to an integer if possible.
// If conversion is successful, it returns the uint64 value, else it returns the raw string.
func ExampleTag(field reflect.StructField) interface{} {
	tagValue := field.Tag.Get("example")

	if tagValue != "" {
		numValue, err := strconv.ParseUint(tagValue, 10, 64)
		if err == nil {
			return numValue
		}
		return tagValue
	}

	return nil
}

// JsonTag extracts the 'json' struct tag's value of a struct field and returns it as a string.
// If the tag contains options (comma-separated), only the name part before the comma is returned.
func JsonTag(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if strings.Index(jsonTag, ",") > 0 {
		return strings.Split(jsonTag, ",")[0]
	}
	return jsonTag
}

// IsOmitempty extracts the 'json' struct tag's value of a struct field and returns if it has omitempty.
func IsOmitempty(field reflect.StructField) bool {
	jsonTag := field.Tag.Get("json")
	for _, part := range strings.Split(jsonTag, ",") {
		if strings.TrimSpace(part) == "omitempty" {
			return true
		}
	}
	return false
}

// IsRequired extracts the 'required' struct tag's value of a struct field and returns true if required is true.
func IsRequired(field reflect.StructField) bool {
	tagValue := field.Tag.Get("required")
	return tagValue == "true"
}

// Type maps a string to its corresponding Swagger type according to the
// Swagger Specification version 2 data types (https://swagger.io/specification/v2/#data-types).
func Type(t reflect.Type) string {

	kind := t.Kind()
	switch {
	case kind == reflect.Pointer:
		return Type(t.Elem())
	case kind == reflect.Interface:
		return "interface"
	case isInteger(kind):
		return "integer"
	case kind == reflect.Slice, kind == reflect.Array:
		return "array"
	case kind == reflect.Bool:
		return "boolean"
	case kind == reflect.Float64, kind == reflect.Float32:
		return "number"
	case kind == reflect.String:
		return "string"
	default:
		panic(fmt.Sprintf("unsupported type: %s", kind))
	}
}

func IsPrimitiveValue(valueKind reflect.Kind) bool {
	return valueKind == reflect.Bool ||
		valueKind == reflect.Float64 ||
		valueKind == reflect.Float32 ||
		valueKind == reflect.String ||
		isInteger(valueKind)
}

func isInteger(valueKind reflect.Kind) bool {
	return containsKind([]reflect.Kind{
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	}, valueKind)
}

func containsKind(s []reflect.Kind, k reflect.Kind) bool {
	for _, v := range s {
		if v == k {
			return true
		}
	}
	return false
}
