package ubjson

import "bytes"

// Types which implement Value perform their own data marshaling and unmarshaling.
type Value interface {
	// The type marker for this kind of value. Must always return the same value.
	UBJSONType() Marker
	// Marshals the value to an encoder using containers and primitive values.
	MarshalUBJSON(*Encoder) error
	// Unmarshals containers and primitive values from a decoder into the value.
	UnmarshalUBJSON(*Decoder) error
}

// The Marshal function marshals v. Types implementing Value will be encoded
// with their UBJSONType and MarshalUBJSON methods.
func Marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	if err := NewEncoder(&b).Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// The MarshalBlock function is like Marshal but produces human-readable
// block-notation, rather than binary.
func MarshalBlock(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	if err := NewBlockEncoder(&b).Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// The Unmarshal function unmarshals universal binary json into v.
// Types implementing Value will be decoded via their UBJSONType and
// UnmarshalUBJSON methods.
func Unmarshal(binary []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(binary)).Decode(v)
}

// The UnmarshalBlock function is like Unmarshal but parses human-readable
// block-notation, rather than binary.
func UnmarshalBlock(block []byte, v interface{}) error {
	return NewBlockDecoder(bytes.NewReader(block)).Decode(v)
}

// A Char is a byte which implements value to use 'C' instead of 'U'.
// Must be <127.
type Char byte

// A HighPrecNumber is a decimal string of arbitrary precision, which is encoded
// with 'H'
type HighPrecNumber string
