package ubjson

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// The Marshal function marshals v into universal binary json.
// Types implementing Marshaler will be encoded via their MarshalUBJSON method.
func Marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	if err := NewEncoder(&b).Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// An Encoder provides methods for encoding UBJSON data types to an io.Writer.
type Encoder struct {
	// The primary writer.
	w io.Writer
	// A delegate or Writer.
	bw io.ByteWriter

	// A buffer as large as the largest fixed size type.
	buf [8]byte

	// Block notation output.
	block bool
	// Current number of indentations.
	indent int
}

// The NewEncoder functions returns a new Encoder for w.
func NewEncoder(w io.Writer) *Encoder {
	var bw io.ByteWriter
	if t, ok := w.(io.ByteWriter); ok {
		bw = t
	} else {
		bw = &byteWriter{w}
	}
	return &Encoder{w: w, bw: bw}
}

// The Block method enables block-notation encoding.
func (e *Encoder) Block() *Encoder {
	e.block = true
	return e
}

// The newLine method starts a new, indented line if block-notation is enabled,
// otherwise it is a no-op.
func (e *Encoder) newLine() error {
	if e.block {
		if err := e.writeByte('\n'); err != nil {
			return err
		}
		for i := 0; i < e.indent; i++ {
			if err := e.writeByte('\t'); err != nil {
				return err
			}
		}
	}
	return nil
}

// The writeType method writes the marker byte, respecting block notation.
func (e *Encoder) writeType(marker byte) error {
	if e.block {
		return e.write([]byte{'[', marker, ']'})
	}
	return e.writeByte(marker)
}

// The write method writes p.
func (e *Encoder) write(p []byte) error {
	_, err := e.w.Write(p)
	return err
}

// The writeByte method writes c.
func (e *Encoder) writeByte(c byte) error { return e.bw.WriteByte(c) }

// The blocked method writes s surrounded by square brackets.
func (e *Encoder) blocked(s string) error {
	_, err := fmt.Fprintf(e.w, "[%s]", s)
	return err
}

// The EncodeNull method encodes the null marker.
func (e *Encoder) EncodeNull() error {
	return e.marshal(nullMarshaler)
}

// The EncodeNoOp method encodes the NoOp marker.
func (e *Encoder) EncodeNoOp() error {
	return e.marshal(noOpMarshaler)
}

// The EncodeTrue method encodes the true Marker.
func (e *Encoder) EncodeTrue() error {
	return e.marshal(trueMarshaler)
}

// The EncodeFalse method encodes the false Marker.
func (e *Encoder) EncodeFalse() error {
	return e.marshal(falseMarshaler)
}

// The EncodeInt8 method encodes the int8 v.
func (e *Encoder) EncodeInt8(v int8) error {
	return e.marshal(int8Marshaler(v))
}

type int8Marshaler int8

func (i int8Marshaler) MarshalUBJSON() (byte, func(*Encoder) error) {
	return int8Marker, func(e *Encoder) error {
		if e.block {
			return e.blocked(strconv.FormatInt(int64(i), 10))
		}
		return e.writeByte(uint8(i))
	}
}

// The EncodeUInt8 method encodes the uint8 v.
func (e *Encoder) EncodeUInt8(v uint8) error {
	return e.marshal(uint8Marshaler(v))
}

type uint8Marshaler uint8

func (i uint8Marshaler) MarshalUBJSON() (byte, func(*Encoder) error) {
	return uInt8Marker, func(e *Encoder) error {
		if e.block {
			return e.blocked(strconv.FormatUint(uint64(i), 10))
		}
		return e.writeByte(uint8(i))
	}
}

// The EncodeInt16 method encodes the int16 v.
func (e *Encoder) EncodeInt16(v int16) error {
	return e.marshal(int16Marshaler(v))
}

type int16Marshaler int16

func (i int16Marshaler) MarshalUBJSON() (byte, func(*Encoder) error) {
	return int16Marker, func(e *Encoder) error {
		if e.block {
			return e.blocked(strconv.FormatInt(int64(i), 10))
		}
		b := e.buf[:2]
		binary.BigEndian.PutUint16(b, uint16(i))
		return e.write(b)
	}
}

// The EncodeInt32 method encodes the int32 v.
func (e *Encoder) EncodeInt32(v int32) error {
	return e.marshal(int32Marshaler(v))
}

type int32Marshaler int32

func (i int32Marshaler) MarshalUBJSON() (byte, func(*Encoder) error) {
	return int32Marker, func(e *Encoder) error {
		if e.block {
			return e.blocked(strconv.FormatInt(int64(i), 10))
		}
		b := e.buf[:4]
		binary.BigEndian.PutUint32(b, uint32(i))
		return e.write(b)
	}
}

// The EncodeInt64 method encodes the int64 v.
func (e *Encoder) EncodeInt64(v int64) error {
	return e.marshal(int64Marshaler(v))
}

type int64Marshaler int64

func (i int64Marshaler) MarshalUBJSON() (byte, func(*Encoder) error) {
	return int64Marker, func(e *Encoder) error {
		if e.block {
			return e.blocked(strconv.FormatInt(int64(i), 10))
		}
		b := e.buf[:8]
		binary.BigEndian.PutUint64(b, uint64(i))
		return e.write(b)
	}
}

// The EncodeInt method encodes the int v in the smallest possible format.
func (e *Encoder) EncodeInt(v int) error {
	return e.marshal(intMarshaler(v))
}

// The intMarshaler method returns a Marshaler for the int v, which uses the smallest possible format.
func intMarshaler(v int) Marshaler {
	switch {
	case v >= 0 && v <= 255:
		return uint8Marshaler(uint8(v))
	case v <= 127 && v >= -128:
		return int8Marshaler(int8(v))
	case v <= 32767 && v >= -32768:
		return int16Marshaler(int16(v))
	case v <= 2147483647 && v >= -2147483648:
		return int32Marshaler(int32(v))
	default:
		return int64Marshaler(int64(v))
	}
}

// The EncodeFloat32 method encodes the float32 v.
func (e *Encoder) EncodeFloat32(v float32) error {
	return e.marshal(float32Marshaler(v))
}

type float32Marshaler float32

func (f float32Marshaler) MarshalUBJSON() (byte, func(*Encoder) error) {
	return float32Marker, func(e *Encoder) error {
		if e.block {
			return e.blocked(strconv.FormatFloat(float64(f), 'g', -1, 32))
		}
		b := e.buf[:4]
		binary.BigEndian.PutUint32(b, math.Float32bits(float32(f)))
		return e.write(b)
	}
}

// The EncodeFloat64 method encodes the float64 v.
func (e *Encoder) EncodeFloat64(v float64) error {
	return e.marshal(float64Marshaler(v))
}

type float64Marshaler float64

func (f float64Marshaler) MarshalUBJSON() (byte, func(*Encoder) error) {
	return float64Marker, func(e *Encoder) error {
		if e.block {
			return e.blocked(strconv.FormatFloat(float64(f), 'g', -1, 64))
		}
		b := e.buf[:8]
		binary.BigEndian.PutUint64(b, math.Float64bits(float64(f)))
		return e.write(b)
	}
}

// The EncodeHighPrecNum method encodes the string v as a high precision number.
func (e *Encoder) EncodeHighPrecNum(v string) error {
	return e.marshal(HighPrecNumber(v))
}

// The EncodeChar method encodes the byte v.
func (e *Encoder) EncodeChar(v byte) error {
	return e.marshal(charMarshaler(v))
}

type charMarshaler byte

func (c charMarshaler) MarshalUBJSON() (byte, func(*Encoder) error) {
	return charMarker, func(e *Encoder) error {
		if e.block {
			return e.write([]byte{'[', byte(c), ']'})
		}
		return e.writeByte(byte(c))
	}
}

// The EncodeString method encodes the string v.
func (e *Encoder) EncodeString(v string) error {
	return e.marshal(stringMarshaler(v))
}

type stringMarshaler string

func (m stringMarshaler) MarshalUBJSON() (byte, func(*Encoder) error) {
	return stringMarker, func(e *Encoder) error {
		s := string(m)
		if err := e.EncodeInt(len(s)); err != nil {
			return err
		}
		if len(s) > 0 {
			if e.block {
				return e.blocked(s)
			}
			if _, err := io.WriteString(e.w, s); err != nil {
				return err
			}
		}
		return nil
	}
}

// The getType method returns the marker byte for t.
// TODO do we need complex types? [{][#][U]... []byte? something safer? TypeEncoder()?
func getType(t reflect.Type) byte {
	k := t.Kind()
	if m, ok := reflect.New(t).Interface().(Marshaler); ok {
		//TODO to be able to optimize on custom types, they must ALWAYS return the same byte type
		marker, _ := m.MarshalUBJSON()
		return marker
	}

	switch k {
	case reflect.String:
		return stringMarker
	case reflect.Bool:
		return 0 //TODO can't assume single type
		//if t {
		//	return trueMarker
		//} else {
		//	return falseMarker
		//}
	case reflect.Int:
		return 0 //TODO can't assume single type
	case reflect.Int8:
		return int8Marker
	case reflect.Uint8:
		return uInt8Marker
	case reflect.Int16:
		return int16Marker
	case reflect.Int32:
		return int32Marker
	case reflect.Int64:
		return int64Marker
	case reflect.Float32:
		return float32Marker
	case reflect.Float64:
		return float64Marker

		//TODO containers have complex types
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
		return 0
	}

	//TODO more?
	return 0
}

//TODO encodes a ubjson array (or slice) containing the elements from v.
type ArrayMarshaler struct {
	*reflect.Value
}

func (a *ArrayMarshaler) MarshalUBJSON() (byte, func(*Encoder) error) {
	return arrayStartMarker, func(e *Encoder) error {
		var elemType reflect.Type
		if a.Type().Elem().Kind() != reflect.Interface {
			elemType = a.Type().Elem()
		}
		var i int
		return e.EncodeArray(elemType, a.Len(), func() Marshaler {
			if i >= a.Len() {
				return nil
			}
			m := getMarshaler(a.Index(i).Interface())
			i++
			return m
		})
	}
}

func (e *Encoder) EncodeArray(elemType reflect.Type, count int, elems func() Marshaler) error {
	e.indent++

	// Optimize type for everything but interfaces.
	if elemType != nil {
		if err := e.writeType(typeMarker); err != nil {
			return err
		}
		//TODO this isn't correct; should be able to have strongly typed container defs in here
		m := getType(elemType)
		if m == 0 {
			//TODO can't optimize type, since it's variable (bool/int)
		}
		if err := e.writeType(m); err != nil {
			return err
		}
	}

	if count >= 0 {
		if err := e.writeType(countMarker); err != nil {
			return err
		}
		if err := e.EncodeInt(count); err != nil {
			return err
		}
	}

	for {
		elem := elems()
		if elem == nil {
			break
		}
		if err := e.newLine(); err != nil {
			return err
		}
		if elemType == nil {
			if err := e.marshal(elem); err != nil {
				return err
			}
		} else {
			_, data := elem.MarshalUBJSON()
			if err := data(e); err != nil {
				return err
			}
		}
	}

	e.indent--

	if count < 0 {
		if err := e.newLine(); err != nil {
			return err
		}

		if err := e.writeType(arrayEndMarker); err != nil {
			return err
		}
	}

	return e.newLine()
}

type MapMarshaler struct {
	*reflect.Value
}

func (m *MapMarshaler) MarshalUBJSON() (byte, func(*Encoder) error) {
	return objectStartMarker, func(e *Encoder) error {
		var elemType reflect.Type
		valKind := m.Type().Elem().Kind()
		if valKind != reflect.Interface {
			elemType = m.Type().Elem()
		}
		keys := m.MapKeys()
		var i int
		return e.EncodeObject(elemType, m.Len(), func() *MapEntry {
			if i >= len(keys) {
				return nil
			}
			key := keys[i]
			i++
			val := m.MapIndex(key)

			return &MapEntry{
				Key:   key.String(),
				Value: getMarshaler(val.Interface()),
			}
		})

	}
}

type MapEntry struct {
	Key   string
	Value Marshaler
}

//TODO doc
//TODO elemType nil means...
//TODO count -1 means no end marker
//TODo return nil entry for end
func (e *Encoder) EncodeObject(valueType reflect.Type, count int, entries func() *MapEntry) error {
	e.indent++

	if valueType != nil {
		if err := e.writeType(typeMarker); err != nil {
			return err
		}
		if err := e.writeType(getType(valueType)); err != nil {
			return err
		}
	}

	if count >= 0 {
		if err := e.writeType(countMarker); err != nil {
			return err
		}
		if err := e.EncodeInt(count); err != nil {
			return err
		}
	}

	for {
		entry := entries()
		if entry == nil {
			break
		}

		if err := e.newLine(); err != nil {
			return err
		}

		_, keyDataEncoder := stringMarshaler(entry.Key).MarshalUBJSON()
		if err := keyDataEncoder(e); err != nil {
			return err
		}

		if valueType == nil {
			if err := e.marshal(entry.Value); err != nil {
				return err
			}
		} else {
			_, dataEncoder := entry.Value.MarshalUBJSON()
			if err := dataEncoder(e); err != nil {
				return err
			}
		}
	}

	e.indent--

	if count < 0 {
		if err := e.newLine(); err != nil {
			return err
		}

		if err := e.writeType(objectEndMarker); err != nil {
			return err
		}
	}

	return e.newLine()
}

type StructMarshaler struct {
	*reflect.Value
}

//TODO recognize tags, and options
func (s *StructMarshaler) MarshalUBJSON() (byte, func(*Encoder) error) {
	return objectStartMarker, func(e *Encoder) error {
		var count int
		for i := 0; i < s.NumField(); i++ {
			f := s.Type().Field(i)
			if f.PkgPath == "" || f.Anonymous {
				count++
			}
		}
		var i int
		return e.EncodeObject(nil, count, func() *MapEntry {
			var f reflect.StructField
			for {
				if i >= s.NumField() {
					return nil
				}
				f = s.Type().Field(i)
				if f.PkgPath != "" && !f.Anonymous {
					i++
				} else {
					break
				}
			}

			jsonTag := f.Tag.Get("json")
			//TODO check ubjson tag as well eventually
			name := f.Name
			if jsonTag != "" {
				name = strings.Split(jsonTag, ",")[0]
			}

			val := s.Field(i).Interface()
			i++
			return &MapEntry{
				Key:   name,
				Value: getMarshaler(val),
			}
		})
	}
}

// The Encode methods encodes v into universal binary json.
// Types implementing Marshaler will be encoded via their MarshalUBJSON method.
func (e *Encoder) Encode(v interface{}) error {
	return e.marshal(getMarshaler(v))
}

// The marshal functions marshals m to e.
func (e *Encoder) marshal(m Marshaler) error {
	marker, dataEncoder := m.MarshalUBJSON()
	if err := e.writeType(marker); err != nil {
		return err
	}
	if dataEncoder != nil {
		return dataEncoder(e)
	}
	return nil
}

// A constMarshaler is a Marshaler for a single constant byte.
type constMarshaler byte

func (c constMarshaler) MarshalUBJSON() (byte, func(*Encoder) error) {
	return byte(c), nil
}

const (
	nullMarshaler  constMarshaler = nullMarker
	noOpMarshaler  constMarshaler = noOpMarker
	trueMarshaler  constMarshaler = trueMarker
	falseMarshaler constMarshaler = falseMarker
)

// The getMarshaler function returns a Marshaler for v, which may be itself.
func getMarshaler(v interface{}) Marshaler {
	if v == nil {
		return nullMarshaler
	}
	if m, ok := v.(Marshaler); ok {
		return m
	}
	switch t := v.(type) {
	case string:
		return stringMarshaler(t)
	case bool:
		if t {
			return trueMarshaler
		} else {
			return falseMarshaler
		}
	case int:
		return intMarshaler(t)
	case int8:
		return int8Marshaler(t)
	case uint8:
		return uint8Marshaler(t)
	case int16:
		return int16Marshaler(t)
	case int32:
		return int32Marshaler(t)
	case int64:
		return int64Marshaler(t)
	case float32:
		return float32Marshaler(t)
	case float64:
		return float64Marshaler(t)
	}

	// Containers
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		return &ArrayMarshaler{&value}
	case reflect.Map:
		return &MapMarshaler{&value}
	case reflect.Struct:
		return &StructMarshaler{&value}
	}

	panic(fmt.Sprintf("missing encoding sceme for (%T) value: %#v", v, v))
}

//TODO replace marker byte with markerEncoder? (for complex types?)
type Marshaler interface {
	//TODO note that marker must always be the same for a type (otherwise optimization breaks)
	MarshalUBJSON() (marker byte, dataEncoder func(*Encoder) error)
}

type byteWriter struct {
	io.Writer
}

func (w *byteWriter) WriteByte(c byte) error {
	_, err := w.Write([]byte{c})
	return err
}

//TODO tags
