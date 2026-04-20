// Package extensions implements OpenAPI Specification Extensions (x-*).
//
// The OpenAPI spec allows vendor- or tool-specific fields on most objects,
// provided their keys are prefixed with "x-". Any other key is not a valid
// extension and is dropped at serialization time.
//
// See: https://spec.openapis.org/oas/v3.0.3#specification-extensions
package extensions

import (
	"encoding/json"
	"strings"
)

// Extensions is a set of OpenAPI Specification Extension fields. Keys that do
// not start with "x-" are ignored when the host object is serialized.
type Extensions map[string]any

// Prefix is the required prefix for OpenAPI extension keys.
const Prefix = "x-"

// Merge serializes v with the standard JSON marshaler and splices the x-*
// entries of ext into the resulting object.
//
// Callers must pass a type alias of the host struct (one without its own
// MarshalJSON method) so this call does not recurse.
func Merge(v any, ext Extensions) ([]byte, error) {
	base, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if len(ext) == 0 {
		return base, nil
	}

	var m map[string]any
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}
	for k, val := range ext {
		if strings.HasPrefix(k, Prefix) {
			m[k] = val
		}
	}
	return json.Marshal(m)
}
