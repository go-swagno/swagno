package utils

import (
	"crypto/sha256"
	"fmt"
	"reflect"
	"sort"
)

func ConvertInterfaceToMap(t interface{}) map[string]interface{} {
	m := make(map[string]interface{})

	reflectVal := reflect.ValueOf(t)
	iter := reflectVal.MapRange()
	for iter.Next() {
		m[iter.Key().String()] = iter.Value().Interface()
	}

	return m
}

func GetHashOfMap(m map[string]interface{}) string {
	h := sha256.New()

	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := m[k]
		b := sha256.Sum256([]byte(fmt.Sprintf("%v", k)))
		h.Write(b[:])
		b = sha256.Sum256([]byte(fmt.Sprintf("%v", v)))
		h.Write(b[:])
	}

	return fmt.Sprintf("%x", h.Sum(nil))[:10]
}
