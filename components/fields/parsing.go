package fields

import (
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

// Type maps a string to its corresponding Swagger type according to the
// Swagger Specification version 2 data types (https://swagger.io/specification/v2/#data-types).
func Type(t string) string {
	if t == "interface" {
		return "interface"
	} else if strings.Contains(strings.ToLower(t), "int") {
		return "integer"
	} else if t == "array" || t == "slice" {
		return "array"
	} else if t == "bool" {
		return "boolean"
	} else if t == "float64" || t == "float32" {
		return "number"
	}
	return t
}
