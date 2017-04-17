package ubjson_test

import (
	"fmt"

	"github.com/jmank88/ubjson"
)

// A Nameless struct encodes itself as a fixed length, ordered Array container,
// omitting the string names to save space.
type Nameless struct {
	Field1 string
	FieldA int8
}

func (n *Nameless) UBJSONType() ubjson.Marker {
	return ubjson.ArrayStartMarker
}

func (n *Nameless) MarshalUBJSON(e *ubjson.Encoder) error {
	a, err := e.ArrayLen(2)
	if err != nil {
		return err
	}
	if err := a.EncodeString(n.Field1); err != nil {
		return err
	}
	if err := a.EncodeInt8(n.FieldA); err != nil {
		return err
	}
	return a.End()
}

func (n *Nameless) UnmarshalUBJSON(d *ubjson.Decoder) error {
	a, err := d.Array()
	if err != nil {
		return err
	}
	n.Field1, err = d.DecodeString()
	if err != nil {
		return err
	}
	n.FieldA, err = d.DecodeInt8()
	if err != nil {
		return err
	}
	return a.End()
}

func Example_structAsArray() {
	if b, err := ubjson.MarshalBlock(&Nameless{Field1: "test", FieldA: 42}); err != nil {
		fmt.Println("error: " + err.Error())
	} else {
		fmt.Println(string(b))
	}
	// Output:
	// [[][#][U][2]
	// 	[S][U][4][test]
	// 	[i][42]
}
