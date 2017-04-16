package ubjson_test

import (
	"fmt"

	"github.com/jmank88/ubjson"
)

// A CustomValue encodes itself as a fixed length object container.
type CustomValue struct {
	Field1 string
	FieldA int
}

func (c *CustomValue) UBJSONType() ubjson.Marker {
	return ubjson.ObjectStartMarker
}

func (c *CustomValue) MarshalUBJSON(e *ubjson.Encoder) error {
	o, err := e.ObjectLen(2)
	if err != nil {
		return err
	}
	if err := o.EncodeKey("Field1"); err != nil {
		return err
	}
	if err := o.EncodeString(c.Field1); err != nil {
		return err
	}
	if err := o.EncodeKey("FieldA"); err != nil {
		return err
	}
	if err := o.EncodeInt(c.FieldA); err != nil {
		return err
	}
	return o.End()
}

func (c *CustomValue) UnmarshalUBJSON(d *ubjson.Decoder) error {
	o, err := d.Object()
	if err != nil {
		return err
	}
	for o.NextEntry() {
		k, err := o.DecodeKey()
		if err != nil {
			return err
		}
		switch k {
		case "Field1":
			s, err := o.DecodeString()
			if err != nil {
				return err
			}
			c.Field1 = s
		case "FieldA":
			i, err := o.DecodeInt()
			if err != nil {
				return err
			}
			c.FieldA = i
		}
	}
	return o.End()
}

func Example_CustomValue() {
	v := &CustomValue{Field1: "test", FieldA: 42}
	if b, err := ubjson.MarshalBlock(v); err != nil {
		fmt.Println("error: " + err.Error())
	} else {
		fmt.Println(string(b))
	}

	// Output:
	// [{][#][U][2]
	// 	[U][6][Field1][S][U][4][test]
	// 	[U][6][FieldA][U][42]
}
