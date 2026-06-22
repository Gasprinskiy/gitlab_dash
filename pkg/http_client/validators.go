package http_client

import "reflect"

func IsStruct(v any) bool {
	t := reflect.TypeOf(v)
	return t != nil && t.Kind() == reflect.Struct
}

func IsMap(v any) bool {
	t := reflect.TypeOf(v)
	return t != nil && t.Kind() == reflect.Map
}

func IsValidDest(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return false
	}
	k := rv.Elem().Kind()
	return k == reflect.Struct || k == reflect.Slice
}
