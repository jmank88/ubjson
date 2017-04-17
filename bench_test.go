package ubjson

import (
	"bytes"
	"testing"
	"fmt"
)

//TODO >2X speedup on structs; do all the primitives too (EncodeInt vs Encode_int)

func BenchmarkEncoder_Encode_struct(b *testing.B) {
	var buf []byte
	e := NewEncoder(bytes.NewBuffer(buf))
	b.ResetTimer()
	for i:=0;i<b.N;i++ {
		if err := e.Encode(&bs); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncoder_EncodeValue_struct(b *testing.B) {
	var buf []byte
	e := NewEncoder(bytes.NewBuffer(buf))
	v := benchValue(bs)
	b.ResetTimer()
	for i:=0;i<b.N;i++ {
		if err := e.EncodeValue(&v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecoder_Decode_struct(b *testing.B) {
	var buf bytes.Buffer
	e := NewEncoder(&buf)
	if err := e.Encode(&bs); err != nil {
		b.Fatal(err)
	}
	bin := buf.Bytes()
	r := bytes.NewReader(bin)
	d := NewDecoder(r)
	v := benchStruct{}
	b.ResetTimer()

	for i:=0;i<b.N;i++ {
		r.Reset(bin)
		if err := d.Decode(&v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecoder_DecodeValue_struct(b *testing.B) {
	var buf bytes.Buffer
	e := NewEncoder(&buf)
	v := benchValue(bs)
	if err := e.EncodeValue(&v); err != nil {
		b.Fatal(err)
	}
	bin := buf.Bytes()
	r := bytes.NewReader(bin)
	d := NewDecoder(r)
	v = benchValue{}
	b.ResetTimer()

	for i:=0;i<b.N;i++ {
		r.Reset(bin)
		if err := d.DecodeValue(&v); err != nil {
			b.Fatal(err)
		}
	}
}

var bs = benchStruct{
	UInt8: 100,
	Int8: -42,
	Int16: 500,
	Int32: -1000,
	Int64: 100000,
	String: "test",
	Bytes: []byte("test"),
}

type benchStruct struct {
	UInt8 uint8
	Int8 int8
	Int16 int16
	Int32 int32
	Int64 int64
	Float32 float32
	Float64 float64
	Char Char
	HighPrec HighPrecNumber
	String string
	Bytes []byte
}

type benchValue benchStruct

func (v *benchValue) UBJSONType() Marker {
	return ObjectStartMarker
}

func (v *benchValue) MarshalUBJSON(e *Encoder) error {
	o, err := e.Object()
	if err != nil {
		return err
	}
	if err := o.EncodeKey("UInt8"); err != nil {
		return err
	}
	if err := o.EncodeUInt8(v.UInt8); err != nil {
		return err
	}
	if err := o.EncodeKey("Int8"); err != nil {
		return err
	}
	if err := o.EncodeInt8(v.Int8); err != nil {
		return err
	}
	if err := o.EncodeKey("Int16"); err != nil {
		return err
	}
	if err := o.EncodeInt16(v.Int16); err != nil {
		return err
	}
	if err := o.EncodeKey("Int32"); err != nil {
		return err
	}
	if err := o.EncodeInt32(v.Int32); err != nil {
		return err
	}
	if err := o.EncodeKey("Int64"); err != nil {
		return err
	}
	if err := o.EncodeInt64(v.Int64); err != nil {
		return err
	}
	if err := o.EncodeKey("Float32"); err != nil {
		return err
	}
	if err := o.EncodeFloat32(v.Float32); err != nil {
		return err
	}
	if err := o.EncodeKey("Float64"); err != nil {
		return err
	}
	if err := o.EncodeFloat64(v.Float64); err != nil {
		return err
	}
	if err := o.EncodeKey("Char"); err != nil {
		return err
	}
	if err := o.EncodeChar(byte(v.Char)); err != nil {
		return err
	}
	if err := o.EncodeKey("HighPrec"); err != nil {
		return err
	}
	if err := o.EncodeHighPrecNum(string(v.HighPrec)); err != nil {
		return err
	}
	if err := o.EncodeKey("String"); err != nil {
		return err
	}
	if err := o.EncodeString(v.String); err != nil {
		return err
	}
	if err := o.EncodeKey("Bytes"); err != nil {
		return err
	}
	{
		err := o.EncodeArray(func(e *Encoder) error {
			a, err := o.ArrayType(UInt8Marker, len(v.Bytes))
			if err != nil {
				return err
			}
			for i := range v.Bytes {
				if err := a.EncodeUInt8(v.Bytes[i]); err != nil {
					return err
				}
			}
			return a.End()
		})
		if err != nil {
			return err
		}
	}

	return o.End()
}

func (v *benchValue) UnmarshalUBJSON(d *Decoder) error {
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
		case "UInt8":
			v.UInt8, err = o.DecodeUInt8()
			if err != nil {
				return err
			}
		case "Int8":
			v.Int8, err = o.DecodeInt8()
			if err != nil {
				return err
			}
		case "Int16":
			v.Int16, err = o.DecodeInt16()
			if err != nil {
				return err
			}
		case "Int32":
			v.Int32, err = o.DecodeInt32()
			if err != nil {
				return err
			}
		case "Int64":
			v.Int64, err = o.DecodeInt64()
			if err != nil {
				return err
			}
		case "Float32":
			v.Float32, err = o.DecodeFloat32()
			if err != nil {
				return err
			}
		case "Float64":
			v.Float64, err = o.DecodeFloat64()
			if err != nil {
				return err
			}
		case "Char":
			b, err := o.DecodeChar()
			if err != nil {
				return err
			}
			v.Char = Char(b)
		case "HighPrec":
			s, err := o.DecodeHighPrecNumber()
			if err != nil {
				return err
			}
			v.HighPrec = HighPrecNumber(s)
		case "String":
			v.String, err = o.DecodeString()
			if err != nil {
				return err
			}
		case "Bytes":
			err := o.DecodeArray(func(a *ArrayDecoder) error {
				if a.Len < 0 {
					v.Bytes = make([]byte, 0)
				} else {
					v.Bytes = make([]byte, a.Len)
				}

				for a.NextElem() {
					e, err := a.DecodeUInt8()
					if err != nil {
						return err
					}
					v.Bytes = append(v.Bytes, e)
				}
				return a.End()
			})
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unrecognized object key: %s", k)
		}
	}

	return o.End()
}
