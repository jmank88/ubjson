// +build go1.9

package ubjson

import "reflect"

func makeMap(typ reflect.Type, cap int) reflect.Value {
	if cap < 0 {
		cap = 0
	}
	return reflect.MakeMapWithSize(typ, cap)
}
