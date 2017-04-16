package ubjson

import (
	"fmt"
	"io"
	"reflect"

	"github.com/pkg/errors"
)

// An Encoder provides methods for encoding UBJSON data types.
type Encoder struct {
	writer
	// Function to write value type markers. Normally writeMarker, but
	// overridden by containers to do validation and optimization.
	writeValType func(Marker) error
}

// The NewEncoder functions returns a new Encoder.
func NewEncoder(w io.Writer) *Encoder {
	e := &Encoder{writer: newBinaryWriter(w)}
	e.writeValType = e.writeMarker
	return e
}

// The NewEncoder functions returns a new block-notation Encoder.
func NewBlockEncoder(w io.Writer) *Encoder {
	e := &Encoder{writer: newBlockWriter(w)}
	e.writeValType = e.writeMarker
	return e
}

// The EncodeNull method encodes the null valType.
func (e *Encoder) EncodeNull() error {
	return e.encode(NullMarker, func(*Encoder) error { return nil })
}

// The EncodeNoOp method encodes the NoOp valType.
func (e *Encoder) EncodeNoOp() error {
	return e.encode(NoOpMarker, func(*Encoder) error { return nil })
}

// The EncodeBool method encodes the true or false Marker.
func (e *Encoder) EncodeBool(v bool) error {
	m := FalseMarker
	if v {
		m = TrueMarker
	}
	return e.encode(m, func(*Encoder) error { return nil })
}

// The EncodeValue method encodes the Value v, using UBJSONType and
// UnmarshalUBJSON methods.
func (e *Encoder) EncodeValue(v Value) error {
	return e.encode(v.UBJSONType(), v.MarshalUBJSON)
}

// The EncodeUInt8 method encodes the uint8 v.
func (e *Encoder) EncodeUInt8(v uint8) error {
	return e.encode(UInt8Marker, func(*Encoder) error {
		return e.writeUInt8(v)
	})
}

func (e *Encoder) encode(m Marker, fn func(*Encoder) error) error {
	// Normally actually writes, but omitted for strongly typed containers.
	if err := e.writeValType(m); err != nil {
		return err
	}
	if err := fn(e); err != nil {
		return err
	}
	return e.Flush()
}

// The EncodeInt8 method encodes the int8 v.
func (e *Encoder) EncodeInt8(v int8) error {
	return e.encode(Int8Marker, func(*Encoder) error {
		return e.writeInt8(v)
	})
}

// The EncodeInt16 method encodes the int16 v.
func (e *Encoder) EncodeInt16(v int16) error {
	return e.encode(Int16Marker, func(*Encoder) error {
		return e.writeInt16(v)
	})
}

// The EncodeInt32 method encodes the int32 v.
func (e *Encoder) EncodeInt32(v int32) error {
	return e.encode(Int32Marker, func(*Encoder) error {
		return e.writeInt32(v)
	})
}

// The EncodeInt64 method encodes the int64 v.
func (e *Encoder) EncodeInt64(v int64) error {
	return e.encode(Int64Marker, func(*Encoder) error {
		return e.writeInt64(v)
	})
}

// The EncodeInt method encodes the int v in the smallest possible format.
func (e *Encoder) EncodeInt(v int) error {
	switch smallestIntMarker(int64(v)) {
	case UInt8Marker:
		return e.EncodeUInt8(uint8(v))
	case Int8Marker:
		return e.EncodeInt8(int8(v))
	case Int16Marker:
		return e.EncodeInt16(int16(v))
	case Int32Marker:
		return e.EncodeInt32(int32(v))
	case Int64Marker:
		return e.EncodeInt64(int64(v))
	}
	return errors.New("TODO unreachable, programmere marker error")
}

// The EncodeFloat32 method encodes the float32 v.
func (e *Encoder) EncodeFloat32(v float32) error {
	return e.encode(Float32Marker, func(*Encoder) error {
		return e.writeFloat32(v)
	})
}

// The EncodeFloat64 method encodes the float64 v.
func (e *Encoder) EncodeFloat64(v float64) error {
	return e.encode(Float64Marker, func(*Encoder) error {
		return e.writeFloat64(v)
	})
}

// The EncodeHighPrecNum method encodes the string v as a high precision number.
func (e *Encoder) EncodeHighPrecNum(v string) error {
	return e.encode(HighPrecNumMarker, func(*Encoder) error {
		return e.writeString(v)
	})
}

// The EncodeChar method encodes the byte v.
func (e *Encoder) EncodeChar(v byte) error {
	return e.encode(CharMarker, func(*Encoder) error {
		return e.writeChar(v)
	})
}

// The EncodeString method encodes the string v.
func (e *Encoder) EncodeString(v string) error {
	return e.encode(StringMarker, func(*Encoder) error {
		return e.writeString(v)
	})
}

// The elementMarkerFor function returns a Marker for *strict* types which may be
// optimized away when used as container elements, otherwise it returns 0.
func elementMarkerFor(t reflect.Type) Marker {
	if t == nil {
		return 0
	}
	k := t.Kind()
	if v, ok := reflect.New(t).Interface().(Value); ok {
		m := v.UBJSONType()
		switch m {
		case TrueMarker, FalseMarker:
			return 0
		}
	}

	switch k {
	case reflect.Bool, reflect.Int:
		return 0

	case reflect.String:
		return StringMarker
	case reflect.Int8:
		return Int8Marker
	case reflect.Uint8:
		return UInt8Marker
	case reflect.Int16:
		return Int16Marker
	case reflect.Int32:
		return Int32Marker
	case reflect.Int64:
		return Int64Marker
	case reflect.Float32:
		return Int64Marker
	case reflect.Float64:
		return Int64Marker
	case reflect.Array, reflect.Slice:
		return ArrayStartMarker
	case reflect.Map, reflect.Struct:
		return ObjectStartMarker
	}
	return 0
}

// An ArrayEncoder supplements an Encoder with an End() methods, and performs
// validation and optimization of array elements. Callers must finish with a
// call to End().
type ArrayEncoder struct {
	Encoder
	elemType Marker
	len      int
	count    int
}

//The End method completes array encoding.
func (e *ArrayEncoder) End() error {
	e.decIndent()

	if e.len < 0 {
		if err := e.writeNewLine(); err != nil {
			return err
		}

		if err := e.writeMarker(arrayEndMarker); err != nil {
			return err
		}
	} else if e.len != e.count {
		return fmt.Errorf("unable to end array of length %d after %d elements", e.len, e.count)
	}

	return e.Flush()
}

// An ObjectEncoder supplements an Encoder with EncodeKey() and End() methods,
// and performs validation and optimization of object Values. Callers must
// alternate Key() and Encode*() methods for the specified number of entries
// and finish with End().
type ObjectEncoder struct {
	Encoder
	// Value type for strongly typed objects, otherwise 0.
	valType Marker
	// Number of entries, or -1 for unspecified.
	len int
	// Count of entries encoded so far.
	count int
}

func (o *ObjectEncoder) writeValType(m Marker) error {
	o.count++

	if o.len >= 0 {
		if o.count > 2*o.len {
			return errTooMany(o.len)
		}
	}
	if o.count%2 == 1 {
		return errors.New("expected key not value")
	}

	if o.valType == 0 {
		if err := o.writeMarker(m); err != nil {
			return err
		}
	} else {
		if o.valType != m {
			return errWrongTypeWrite(o.valType, m)
		}
		// Omit type marker.
	}
	return nil
}

// The EncodeKey method encodes an object key.
func (o *ObjectEncoder) EncodeKey(key string) error {
	o.count++

	if o.len >= 0 {
		if o.count > 2*o.len {
			return errTooMany(o.len)
		}
	}
	if o.count%2 == 0 {
		return errors.New("expected value not key")
	}

	if err := o.writeNewLine(); err != nil {
		return err
	}

	return o.writeString(key)
}

// The End method checks the length or writes an end maker.
func (o *ObjectEncoder) End() error {
	o.decIndent()

	if o.len < 0 {
		if err := o.writeNewLine(); err != nil {
			return err
		}

		if err := o.writeMarker(objectEndMarker); err != nil {
			return err
		}
	} else if 2*o.len != o.count {
		return fmt.Errorf("unable to end map of %d entries after %d", o.len, o.count/2)
	}

	return o.Flush()
}

// The Object method begins encoding an object container.
func (e *Encoder) Object() (*ObjectEncoder, error) {
	return e.ObjectType(0, -1)
}

// The ObjectLen method begins encoding an object container with a specified
// length.
func (e *Encoder) ObjectLen(len int) (*ObjectEncoder, error) {
	return e.ObjectType(0, len)
}

// The ObjectType method begins encoding a strongly-typed object container with
// a specified length.
func (e *Encoder) ObjectType(valType Marker, len int) (*ObjectEncoder, error) {
	e.incIndent()

	if err := e.writeContainer(valType, len); err != nil {
		return nil, err
	}

	o := &ObjectEncoder{valType: valType, len: len}
	o.Encoder.writer = e.writer
	o.Encoder.writeValType = o.writeValType
	return o, nil
}

// The Array method begins encoding an array container.
func (e *Encoder) Array() (*ArrayEncoder, error) {
	return e.ArrayType(0, -1)
}

// The ArrayLen method begins encoding an array container with a specified
// length.
func (e *Encoder) ArrayLen(len int) (*ArrayEncoder, error) {
	return e.ArrayType(0, len)
}

// The ArrayType method begins encoding a strongly-typed array container with a
// specified length.
// When encoding a single byte element type, actual elements are optimized away,
// and End() must be called immediately.
func (e *Encoder) ArrayType(elemType Marker, len int) (*ArrayEncoder, error) {
	e.incIndent()

	if err := e.writeContainer(elemType, len); err != nil {
		return nil, err
	}

	a := &ArrayEncoder{elemType: elemType, len: len}
	a.Encoder.writer = e.writer
	a.Encoder.writeValType = a.writeElemType
	return a, nil
}

func (e *Encoder) writeContainer(elemType Marker, len int) error {
	// Optimize type?
	if elemType != 0 {
		if err := e.writeMarker(typeMarker); err != nil {
			return err
		}
		if err := e.writeMarker(elemType); err != nil {
			return err
		}
	}

	// Fixed length?
	if len >= 0 {
		if err := e.writeMarker(countMarker); err != nil {
			return err
		}
		if err := writeInt(e, len); err != nil {
			return err
		}
	}
	return nil
}

func (a *ArrayEncoder) writeElemType(m Marker) error {
	if a.len >= 0 {
		a.count++
		if a.count > a.len {
			return errTooMany(a.len)
		}
	}

	if err := a.writeNewLine(); err != nil {
		return err
	}

	if a.elemType == 0 {
		if err := a.writeMarker(m); err != nil {
			return err
		}
	} else {
		if a.elemType != m {
			return errWrongTypeWrite(a.elemType, m)
		}
		// Omit type marker.
	}
	return nil
}

// The Encode methods encodes v into universal binary json.
// Types implementing Value will be encoded via their MarshalUBJSON method.
func (e *Encoder) Encode(v interface{}) error {
	if v == nil {
		return e.EncodeNull()
	}
	if val, ok := v.(Value); ok {
		return e.EncodeValue(val)
	}
	switch t := v.(type) {
	case string:
		return e.EncodeString(t)
	case bool:
		return e.EncodeBool(t)
	case int:
		return e.EncodeInt(t)
	case int8:
		return e.EncodeInt8(t)
	case uint8:
		return e.EncodeUInt8(t)
	case int16:
		return e.EncodeInt16(t)
	case int32:
		return e.EncodeInt32(t)
	case int64:
		return e.EncodeInt64(t)
	case float32:
		return e.EncodeFloat32(t)
	case float64:
		return e.EncodeFloat64(t)
	case Char:
		return e.EncodeChar(byte(t))
	case HighPrecNumber:
		return e.EncodeHighPrecNum(string(t))
	}

	// Containers
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		return e.encode(ArrayStartMarker, encodeArray(value))

	case reflect.Map:
		if value.Type().Key().Kind() != reflect.String &&
			!value.Type().Key().ConvertibleTo(stringType) {
			return fmt.Errorf("unable to encode map of type %s: key type must be string", value.Type())
		}
		return e.encode(ObjectStartMarker, encodeMap(value))

	case reflect.Struct:
		return e.encode(ObjectStartMarker, encodeStruct(value))

	case reflect.Ptr:
		if value.IsNil() {
			return e.EncodeNull()
		}
		return e.Encode(value.Elem().Interface())
	}

	return fmt.Errorf("unable to encode value: %V", v)
}

func encodeArray(arrayValue reflect.Value) func(*Encoder) error {
	return func(e *Encoder) error {
		var elemType reflect.Type
		if arrayValue.Type().Elem().Kind() != reflect.Interface {
			elemType = arrayValue.Type().Elem()
		}

		var ae *ArrayEncoder
		var err error
		if m := elementMarkerFor(elemType); m == 0 {
			ae, err = e.ArrayLen(arrayValue.Len())
		} else {
			ae, err = e.ArrayType(m, arrayValue.Len())
		}
		if err != nil {
			return err
		}

		for i := 0; i < arrayValue.Len(); i++ {
			if err := ae.Encode(arrayValue.Index(i).Interface()); err != nil {
				return err
			}
		}

		return ae.End()
	}
}

func encodeMap(mapValue reflect.Value) func(*Encoder) error {
	return func(e *Encoder) error {
		var elemType reflect.Type
		valKind := mapValue.Type().Elem().Kind()
		if valKind != reflect.Interface {
			elemType = mapValue.Type().Elem()
		}

		keys := mapKeys(mapValue)

		marker := elementMarkerFor(elemType)
		var o *ObjectEncoder
		var err error
		if marker != 0 {
			o, err = e.ObjectType(marker, len(keys))
		} else {
			o, err = e.ObjectLen(len(keys))
		}
		if err != nil {
			return err
		}

		for _, key := range keys {
			if err := o.EncodeKey(key.String()); err != nil {
				return err
			}
			o.Encode(mapValue.MapIndex(key).Interface())
		}

		return o.End()
	}
}

// Overridden for tests.
var mapKeys = func(mapValue reflect.Value) []reflect.Value {
	return mapValue.MapKeys()
}

func encodeStruct(structValue reflect.Value) func(*Encoder) error {
	return func(e *Encoder) error {
		o, err := e.Object()
		if err != nil {
			return err
		}
		for i := 0; i < structValue.NumField(); i++ {
			f := structValue.Type().Field(i)
			if f.PkgPath == "" {
				if err := o.EncodeKey(f.Name); err != nil {
					return errors.Wrapf(err, "failed to encode key %q", f.Name)
				}
				val := structValue.Field(i).Interface()
				if err := o.Encode(val); err != nil {
					return errors.Wrapf(err, "failed to encode value for key %q", f.Name)
				}
			}
		}

		return o.End()
	}
}
