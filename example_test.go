package ubjson_test

import (
	"fmt"

	"github.com/jmank88/ubjson"
)

func ExampleMarshalBlock() {
	if b, err := ubjson.MarshalBlock(8); err != nil {
		fmt.Println("error: " + err.Error())
	} else {
		fmt.Println(string(b))
	}

	// Output:
	// [U][8]
}

func ExampleMarshalBlock_Ints() {
	if b, err := ubjson.MarshalBlock(8); err != nil {
		fmt.Println("error: " + err.Error())
	} else {
		fmt.Println(string(b))
	}
	if b, err := ubjson.MarshalBlock(-42); err != nil {
		fmt.Println("error: " + err.Error())
	} else {
		fmt.Println(string(b))
	}
	if b, err := ubjson.MarshalBlock(256); err != nil {
		fmt.Println("error: " + err.Error())
	} else {
		fmt.Println(string(b))
	}

	// Output:
	// [U][8]
	// [i][-42]
	// [I][256]
}

func ExampleMarshalBlock_Char() {
	if b, err := ubjson.MarshalBlock(byte('c')); err != nil {
		fmt.Println("error: " + err.Error())
	} else {
		fmt.Println(string(b))
	}
	if b, err := ubjson.MarshalBlock(ubjson.Char('c')); err != nil {
		fmt.Println("error: " + err.Error())
	} else {
		fmt.Println(string(b))
	}

	// Output:
	// [U][99]
	// [C][c]
}

func ExampleMarshalBlock_HighPrecNumber() {
	number := "1234567890.657483921"
	if b, err := ubjson.MarshalBlock(number); err != nil {
		fmt.Println("error: " + err.Error())
	} else {
		fmt.Println(string(b))
	}

	hNumber := ubjson.HighPrecNumber("1234567890.657483921")
	if b, err := ubjson.MarshalBlock(hNumber); err != nil {
		fmt.Println("error: " + err.Error())
	} else {
		fmt.Println(string(b))
	}

	// Output:
	// [S][U][20][1234567890.657483921]
	// [H][U][20][1234567890.657483921]
}

func ExampleMarshalBlock_Array() {
	if b, err := ubjson.MarshalBlock([]byte("testbytes")); err != nil {
		fmt.Println("error: " + err.Error())
	} else {
		fmt.Println(string(b))
	}

	// Output:
	// [[][$][U][#][U][9]
	//	[116]
	//	[101]
	//	[115]
	//	[116]
	//	[98]
	//	[121]
	//	[116]
	//	[101]
	//	[115]
}

func ExampleMarshalBlock_Object() {
	type object struct {
		Str   string
		Int   int64
		Bytes []byte
	}
	o := &object{Str: "str", Int: 45678, Bytes: []byte("test")}
	if b, err := ubjson.MarshalBlock(o); err != nil {
		fmt.Println("error: " + err.Error())
	} else {
		fmt.Println(string(b))
	}

	// Output:
	// [{]
	//	[U][3][Str][S][U][3][str]
	//	[U][3][Int][L][45678]
	//	[U][5][Bytes][[][$][U][#][U][4]
	//		[116]
	//		[101]
	//		[115]
	//		[116]
	// [}]
}
