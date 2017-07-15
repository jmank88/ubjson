// Package ubjson implements encoding and decoding of UBJSON (spec 12).
// http://ubjson.org/
//
// Most types can be automatically encoded through reflection with the Marshal
// and Unmarshal functions. Encoders and Decoders additionally provide type
// specific methods. Custom encodings can be defined by implementing the Value
// interface.
//
//	b, _ := ubjson.MarshalBlock(8)
//	// [U][8]
//	b, _ = ubjson.MarshalBlock("hello")
//	// [S][U][5][hello]
//	var v interface{}
//	...
//	b, _ = ubjson.Marshal(v)
//	// ...
//
package ubjson

import "bytes"

// The Value interface defines a custom encoding.
type Value interface {
	// The type marker for this kind of value. Must always return the same value.
	UBJSONType() Marker
	// Marshals the value to an encoder using containers and primitive values.
	MarshalUBJSON(*Encoder) error
	// Unmarshals containers and primitive values from a decoder into the value.
	UnmarshalUBJSON(*Decoder) error
}

// Marshal encodes a value into UBJSON. Types implementing Value will be encoded
// with their UBJSONType and MarshalUBJSON methods.
func Marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	if err := NewEncoder(&b).Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// MarshalBlock encodes a value into UBJSON block-notation. Types implementing
// Value will be encoded with their UBJSONType and MarshalUBJSON methods.
func MarshalBlock(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	if err := NewBlockEncoder(&b).Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Unmarshal decodes a value from UBJSON. Types implementing Value will be
// decoded via their UBJSONType and UnmarshalUBJSON methods.
func Unmarshal(binary []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(binary)).Decode(v)
}

// UnmarshalBlock decodes a value from UBJSON block-notation. Types implementing
// Value will be encoded with their UBJSONType and MarshalUBJSON methods.
func UnmarshalBlock(block []byte, v interface{}) error {
	return NewBlockDecoder(bytes.NewReader(block)).Decode(v)
}

// A Char is a byte which is encoded as 'C' instead of 'U'. Must be <=127.
type Char byte

// A HighPrecNumber is a decimal string of arbitrary precision, which is encoded
// with 'H' instead of 'S'.
type HighPrecNumber string
