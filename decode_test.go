package ubjson

import (
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			var expected interface{} = testCase.value
			actual := reflect.New(reflect.ValueOf(testCase.value).Type())

			err := Unmarshal(testCase.binary, actual.Interface())
			if err != nil {
				t.Fatal("failed to unmarshal:", err.Error())
			}
			if !reflect.DeepEqual(actual.Elem().Interface(), expected) {
				t.Errorf("\nexpected: %T %v \nbut got:  %T %v",
					expected, expected, actual.Elem().Interface(), actual.Elem().Interface())
			}
		})
	}
}

func TestUnmarshalBlock(t *testing.T) {
	t.Parallel()
	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			var expected interface{} = testCase.value
			actual := reflect.New(reflect.ValueOf(testCase.value).Type())

			if err := UnmarshalBlock([]byte(testCase.block), actual.Interface()); err != nil {
				t.Fatal("failed to unmarshal block:", err.Error())
			}
			if !reflect.DeepEqual(actual.Elem().Interface(), expected) {
				t.Errorf("\nexpected: %T %#v \nbut got:  %T %#v",
					expected, expected, actual.Elem().Interface(), actual.Elem().Interface())
			}
		})
	}
}
