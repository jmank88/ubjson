package ubjson_test

import (
	"fmt"

	"github.com/jmank88/ubjson"
)

// A Int8Vec3 encodes itself as a fixed length, strongly-typed array container.
type Int8Vec3 struct {
	a, b, c int8
}

func (iv *Int8Vec3) UBJSONType() ubjson.Marker {
	return ubjson.ArrayStartMarker
}

func (iv *Int8Vec3) MarshalUBJSON(e *ubjson.Encoder) error {
	a, err := e.ArrayType(ubjson.Int8Marker, 3)
	if err != nil {
		return err
	}
	if err := a.EncodeInt8(iv.a); err != nil {
		return err
	}
	if err := a.EncodeInt8(iv.b); err != nil {
		return err
	}
	if err := a.EncodeInt8(iv.c); err != nil {
		return err
	}
	return a.End()
}

func (iv *Int8Vec3) UnmarshalUBJSON(d *ubjson.Decoder) error {
	a, err := d.Array()
	if err != nil {
		return err
	}
	if i, err := a.DecodeInt8(); err != nil {
		return err
	} else {
		iv.a = i
	}
	if i, err := a.DecodeInt8(); err != nil {
		return err
	} else {
		iv.b = i
	}
	if i, err := a.DecodeInt8(); err != nil {
		return err
	} else {
		iv.c = i
	}
	return a.End()
}

func Example_pair() {
	iv := &Int8Vec3{a: 100, b: 42, c: -55}
	if b, err := ubjson.MarshalBlock(iv); err != nil {
		fmt.Println("error: " + err.Error())
	} else {
		fmt.Println(string(b))
	}

	// Output:
	// [[][$][i][#][U][3]
	// 	[100]
	// 	[42]
	// 	[-55]
}
