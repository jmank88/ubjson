package ubjson_test

import (
	"fmt"

	"github.com/jmank88/ubjson"
)

// A TaggedStruct has fields with 'ubjson' tags.
type TaggedStruct struct {
	Field1 string `ubjson:"field1"`
	FieldA int    `json:"ignored" ubjson:"fieldA"`
}

func Example_taggedStruct() {
	v := &TaggedStruct{Field1: "test", FieldA: 42}
	if b, err := ubjson.MarshalBlock(v); err != nil {
		fmt.Println("error: " + err.Error())
	} else {
		fmt.Println(string(b))
	}

	// Output:
	// [{]
	// 	[U][6][field1][S][U][4][test]
	// 	[U][6][fieldA][U][42]
	// [}]
}
