package ubjson

import (
	"io"
)

type Decoder struct {
	parser Parser
}

func (d *Decoder) Decode(into interface{}) error {
	return d.parser.decode(into)
}

func BytesDecoder(b []byte) *Decoder {
	return &Decoder{}
}

func StringDecoder(s string) *Decoder {
	return &Decoder{}
}

func ReaderDecoder(r io.Reader) *Decoder {
	return &Decoder{}
}

func (d *Decoder) DecodeInt(v *int) error {
	return nil //TODO route to type specific by peeking type
}

//TODO mirror the rest of encode methods
//fill in the scanner interface as we go
