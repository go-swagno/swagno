package generator

import (
	"reflect"
	"strconv"
	"strings"
)

func getExampleTag(field reflect.StructField) interface{} {
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

// get struct json tag as string of a struct field
func getJsonTag(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if strings.Index(jsonTag, ",") > 0 {
		return strings.Split(jsonTag, ",")[0]
	}
	return jsonTag
}

// get swagger type from reflection type
// https://swagger.io/specification/v2/#data-types
func getType(t string) string {
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
