package dto

import "reflect"

func omitEmpty(values map[string]interface{}) map[string]interface{} {
	for k, v := range values {
		t := reflect.TypeOf(v)
		v := reflect.ValueOf(v)
		if t.Kind() == reflect.String && v.IsZero() {
			delete(values, k)
		}
		if t.Kind() == reflect.Ptr && v.IsNil() {
			delete(values, k)
		}
	}
	return values
}
