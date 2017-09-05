package ubjson

import (
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	for name, tc := range cases {
		t.Run(name, tc.unmarshal)
	}
}

func (tc *testCase) unmarshal(t *testing.T) {
	var expected interface{} = tc.value
	actual := reflect.New(reflect.ValueOf(tc.value).Type())

	if err := Unmarshal(tc.binary, actual.Interface()); err != nil {
		t.Fatal("failed to unmarshal:", err.Error())
	}
	if !reflect.DeepEqual(actual.Elem().Interface(), expected) {
		t.Errorf("\nexpected: %T %v \nbut got:  %T %v",
			expected, expected, actual.Elem().Interface(), actual.Elem().Interface())
	}
}

func TestUnmarshalBlock(t *testing.T) {
	t.Parallel()
	for name, tc := range cases {
		t.Run(name, tc.unmarshalBlock)
	}
}

func (tc *testCase) unmarshalBlock(t *testing.T) {
	var expected interface{} = tc.value
	actual := reflect.New(reflect.ValueOf(tc.value).Type())

	if err := UnmarshalBlock([]byte(tc.block), actual.Interface()); err != nil {
		t.Fatal("failed to unmarshal block:", err.Error())
	}
	if !reflect.DeepEqual(actual.Elem().Interface(), expected) {
		t.Errorf("\nexpected: %T %#v \nbut got:  %T %#v",
			expected, expected, actual.Elem().Interface(), actual.Elem().Interface())
	}
}
