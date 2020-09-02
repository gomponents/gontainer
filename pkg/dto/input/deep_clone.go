package input

import (
	"reflect"
)

func deepClone(i interface{}) interface{} {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Slice:
		return deepCloneSlice(i)
	case reflect.Map:
		return deepCloneMap(i)
	}
	return i
}

func deepCloneSlice(i interface{}) interface{} {
	v := reflect.ValueOf(i)
	r := reflect.MakeSlice(v.Type(), 0, 0)
	for i := 0; i < v.Len(); i++ {
		curr := v.Index(i)
		cp := deepClone(curr.Interface())
		r = reflect.Append(r, reflect.ValueOf(cp))
	}
	return r.Interface()
}

func deepCloneMap(i interface{}) interface{} {
	v := reflect.ValueOf(i)
	r := reflect.MakeMap(v.Type())
	for _, k := range v.MapKeys() {
		curr := v.MapIndex(k)
		cp := deepClone(curr.Interface())
		r.SetMapIndex(k, reflect.ValueOf(cp))
	}
	return r.Interface()
}
