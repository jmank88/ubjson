// +build !go1.9

package ubjson

import "reflect"

func makeMap(typ reflect.Type, _ int) reflect.Value {
	return reflect.MakeMap(typ)
}
