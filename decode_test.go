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

func TestUnmarshalDiscardUnknownFields(t *testing.T) {
	type val struct{ A int8 }

	exp := val{8}
	var got val

	bin := []byte{'{', 'U', 0x01, 'A', 'i', 0x08, 'U', 0x01, 'b', 'i', 0x05, '}'}

	if err := Unmarshal(bin, &got); err != nil {
		t.Fatal(err)
	} else if got != exp {
		t.Errorf("\nexpected: %T %v \nbut got:  %T %v", exp, exp, got, got)
	}

	block := "[{]\n\t[U][1][A][i][8]\n\t[U][1][B][i][5]\n[}]"

	if err := UnmarshalBlock([]byte(block), &got); err != nil {
		t.Fatal(err)
	} else if got != exp {
		t.Errorf("\nexpected: %T %v \nbut got:  %T %v", exp, exp, got, got)
	}
}
