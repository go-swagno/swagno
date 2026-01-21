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
		numValue, err := strconv.ParseFloat(tagValue, 64) // for json float covers for all numbers
		if err == nil {
			return numValue
		}
		return tagValue
	}

	return nil
}

func DescriptionTag(field reflect.StructField) string {
	return field.Tag.Get("desc")
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

// Type maps a string to its corresponding OpenAPI type according to the
// OpenAPI Specification version 3.0.3 data types (https://spec.openapis.org/oas/v3.0.3#data-types).
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
