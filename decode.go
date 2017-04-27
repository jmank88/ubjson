package ubjson

import (
	"fmt"
	"io"
	"reflect"

	"github.com/pkg/errors"
)

// An Decoder provides methods for decoding UBJSON data types.
type Decoder struct {
	reader
	// The readValType function is called to get the next value's type
	// marker. Normally it reads the next marker, but strongly typed
	// containers will do an internal check and also manage counters.
	readValType func() (Marker, error)
}

// NewDecoder returns a new Decoder.
func NewDecoder(r io.Reader) *Decoder {
	d := &Decoder{reader: newBinaryReader(r)}
	d.readValType = d.readMarker
	return d
}

// NewBlockDecoder returns a new block-notation Decoder.
func NewBlockDecoder(r io.Reader) *Decoder {
	d := &Decoder{reader: newBlockReader(r)}
	d.readValType = d.readMarker
	return d
}

// DecodeValue decodes the next value into v.
func (d *Decoder) DecodeValue(v Value) error {
	return d.decodeValue(v.UBJSONType(), v.UnmarshalUBJSON)
}

// decodeValue asserts a value's type marker, then decodes the data.
func (d *Decoder) decodeValue(m Marker, data func(*Decoder) error) error {
	if r, err := d.readValType(); err != nil {
		return errors.Wrapf(err, "failed trying to read type '%s'", m)
	} else if r != m {
		return errWrongTypeRead(m, r)
	}
	return data(d)
}

// assertType reads the next marker and returns an error if it is not m.
func (d *Decoder) assertType(m Marker) error {
	r, err := d.readMarker()
	if err != nil {
		return errors.Wrapf(err, "failed trying to read type '%s'", m)
	}
	if r != m {
		return errWrongTypeRead(m, r)
	}
	return nil
}

// DecodeBool decodes a 'T' or 'F' marker.
func (d *Decoder) DecodeBool() (bool, error) {
	m, err := d.readValType()
	if err != nil {
		return false, err
	}
	switch m {
	case TrueMarker:
		return true, nil
	case FalseMarker:
		return false, nil
	}
	return false, errors.New("expected true or false marker")
}

// DecodeUInt8 decodes a 'U' value into a uint8.
func (d *Decoder) DecodeUInt8() (uint8, error) {
	var v uint8
	return v, d.decodeValue(UInt8Marker, func(*Decoder) error {
		var err error
		v, err = d.readUInt8()
		return err
	})
}

// DecodeInt8 decodes an 'i' value into an int8.
func (d *Decoder) DecodeInt8() (int8, error) {
	var v int8
	return v, d.decodeValue(Int8Marker, func(*Decoder) error {
		var err error
		v, err = d.readInt8()
		return err
	})
}

// DecodeInt16 decodes an 'I' value into an int16.
func (d *Decoder) DecodeInt16() (int16, error) {
	var v int16
	return v, d.decodeValue(Int16Marker, func(*Decoder) error {
		var err error
		v, err = d.readInt16()
		return err
	})
}

// DecodeInt32 decodes an 'l' value into an int32.
func (d *Decoder) DecodeInt32() (int32, error) {
	var v int32
	return v, d.decodeValue(Int32Marker, func(*Decoder) error {
		var err error
		v, err = d.readInt32()
		return err
	})
}

// DecodeInt64 decodes an 'L' value into an int64.
func (d *Decoder) DecodeInt64() (int64, error) {
	var v int64
	return v, d.decodeValue(Int64Marker, func(*Decoder) error {
		var err error
		v, err = d.readInt64()
		return err
	})
}

// DecodeInt decodes an integer value (U,i,I,l,L) into an int.
func (d *Decoder) DecodeInt() (int, error) {
	m, err := d.readValType()
	if err != nil {
		return 0, err
	}
	switch m {
	case UInt8Marker:
		u, err := d.readUInt8()
		return int(u), err
	case Int8Marker:
		i, err := d.readInt8()
		return int(i), err
	case Int16Marker:
		i, err := d.readInt16()
		return int(i), err
	case Int32Marker:
		i, err := d.readInt32()
		return int(i), err
	case Int64Marker:
		i, err := d.readInt64()
		return int(i), err
	default:
		return 0, fmt.Errorf("encountered non-int type marker: %s", m)
	}
}

// DecodeFloat32 decodes an 'f' value into a float32.
func (d *Decoder) DecodeFloat32() (float32, error) {
	var v float32
	return v, d.decodeValue(Float32Marker, func(*Decoder) error {
		var err error
		v, err = d.readFloat32()
		return err
	})
}

// DecodeFloat64 decodes an 'F' value into a float64.
func (d *Decoder) DecodeFloat64() (float64, error) {
	var v float64
	return v, d.decodeValue(Float64Marker, func(*Decoder) error {
		var err error
		v, err = d.readFloat64()
		return err
	})
}

// DecodeHighPrecNumber decodes an 'H' value into a string.
func (d *Decoder) DecodeHighPrecNumber() (string, error) {
	var v string
	return v, d.decodeValue(HighPrecNumMarker, func(*Decoder) error {
		var err error
		v, err = d.readString()
		return err
	})
}

// DecodeChar decodes a 'C' value into a byte.
func (d *Decoder) DecodeChar() (byte, error) {
	var v byte
	return v, d.decodeValue(CharMarker, func(*Decoder) error {
		var err error
		v, err = d.readChar()
		return err
	})
}

// DecodeString decodes an 'S' value into a string.
func (d *Decoder) DecodeString() (string, error) {
	var v string
	return v, d.decodeValue(StringMarker, func(*Decoder) error {
		var err error
		v, err = d.readString()
		return err
	})
}

func (d *Decoder) decodeInterface() (interface{}, error) {
	m, err := d.readValType()
	if err != nil {
		return nil, err
	}
	switch m {
	case NullMarker, NoOpMarker:
		return nil, nil

	case TrueMarker:
		return true, nil

	case FalseMarker:
		return false, nil

	case UInt8Marker:
		return d.readUInt8()

	case Int8Marker:
		return d.readInt8()

	case Int16Marker:
		return d.readInt16()

	case Int32Marker:
		return d.readInt32()

	case Int64Marker:
		return d.readInt64()

	case Float32Marker:
		return d.readFloat32()

	case Float64Marker:
		return d.readFloat64()

	case StringMarker:
		return d.readString()

	case HighPrecNumMarker:
		s, err := d.readString()
		return HighPrecNumber(s), err

	case CharMarker:
		b, err := d.readChar()
		return Char(b), err

	case ArrayStartMarker:
		a, err := d.Array()
		if err != nil {
			return nil, err
		}
		return arrayAsInterface(a)

	case ObjectStartMarker:
		o, err := d.Object()
		if err != nil {
			return nil, err
		}
		return objectAsInterface(o)

	default:
		return nil, fmt.Errorf("failed to decode: unrecgonized type marker %q", m)
	}
}

// DecodeObject decodes an object container using dataFn.
func (d *Decoder) DecodeObject(dataFn func(*ObjectDecoder) error) error {
	return d.decodeValue(ObjectStartMarker, func(d *Decoder) error {
		o, err := d.Object()
		if err != nil {
			return err
		}
		return dataFn(o)
	})
}

// Object begins decoding an object, and returns a specialized decoder for
// object entries.
func (d *Decoder) Object() (*ObjectDecoder, error) {
	m, l, err := readContainer(d)
	if err != nil {
		return nil, err
	}
	o := &ObjectDecoder{ValType: m, Len: l}

	o.Decoder.reader = d.reader
	o.Decoder.readValType = o.readValType

	return o, nil
}

// DecodeArray decodes an array container using dataFn.
func (d *Decoder) DecodeArray(dataFn func(*ArrayDecoder) error) error {
	return d.decodeValue(ArrayStartMarker, func(d *Decoder) error {
		a, err := d.Array()
		if err != nil {
			return err
		}
		return dataFn(a)
	})
}

// Array begins decoding an array, and returns a specialized decoder for array
// elements.
func (d *Decoder) Array() (*ArrayDecoder, error) {
	m, l, err := readContainer(d)
	if err != nil {
		return nil, err
	}

	a := &ArrayDecoder{ElemType: m, Len: l}

	a.Decoder.reader = d.reader
	a.Decoder.readValType = a.readElemType

	return a, nil
}

// An ObjectDecoder supplements a Decoder with NextEntry(), DecodeKey(), and
// End() methods. Callers must alternate decoding keys and values for either Len
// entries or until NextEntry() returns false, and finish with End().
type ObjectDecoder struct {
	Decoder
	// Value type, or 0 if none included.
	ValType Marker
	// Number of entries, or -1 if unspecified.
	Len int
	// Count of key + val calls (2*Len expected).
	count int
	// Deferred error to be returned by End().
	err error
}

// readValType increments and validates the count, and validate the type either
// from the stream or from o.valType.
func (o *ObjectDecoder) readValType() (Marker, error) {
	o.count++
	if o.Len >= 0 && o.count > 2*o.Len {
		return 0, errTooMany(o.Len)
	}
	if o.count%2 != 0 {
		return 0, errors.New("unable to decode value: expected key")
	}
	if o.ValType == 0 {
		return o.readMarker()
	}
	return o.ValType, nil
}

// DecodeKey reads an object key.
func (o *ObjectDecoder) DecodeKey() (string, error) {
	o.count++
	if o.Len >= 0 && o.count > 2*o.Len {
		return "", errTooMany(o.Len)
	}
	if o.count%2 == 0 {
		return "", errors.New("unable to decode key: expected value")
	}
	return o.readString()
}

// NextEntry returns true if more entries are expected, or false if the end has
// been reached or an error was encountered, in which case End() will return the
// deferred error.
func (o *ObjectDecoder) NextEntry() bool {
	if o.Len < 0 {
		m, err := o.peekMarker()
		if err != nil {
			o.err = err
			return false
		}
		return m != objectEndMarker
	}
	return o.count < 2*o.Len
}

// End completes object decoding and must be called. It may return errors (1)
// deferred from entry decoding, (2) from a missing object end marker, or (3)
// from a length vs. count mismatch.
func (o *ObjectDecoder) End() error {
	if o.err != nil {
		return o.err
	}
	if o.count%2 == 1 {
		return errors.New("cannot end an object with a key")
	}
	if o.Len < 0 {
		m, err := o.readMarker()
		if err != nil {
			return err
		}
		if m != objectEndMarker {
			return errors.New("expected end marker")
		}
	} else if 2*o.Len != o.count {
		return errors.New("len count mismatch")
	}
	return nil
}

// An ArrayDecoder is a Decoder for array container elements.
type ArrayDecoder struct {
	Decoder
	// Element type, or 0 if not present.
	ElemType Marker
	// Number of elements, or -1 if not present.
	Len int
	// Count of element calls (Len expected).
	count int
	// Deferred error.
	err error
}

// readElemType increments and validates the count, and returns the type either
// from the stream or from a.ElemType.
func (a *ArrayDecoder) readElemType() (Marker, error) {
	a.count++
	if a.Len >= 0 && a.count > a.Len {
		return 0, errTooMany(a.Len)
	}
	if a.ElemType == 0 {
		return a.readMarker()
	}
	return a.ElemType, nil
}

// NextElem returns true when there is another element to decode, or false if
// the end of the array has been reached or an error is encountered, in which
// case it will be returned by the End method.
func (a *ArrayDecoder) NextElem() bool {
	if a.Len < 0 {
		m, err := a.peekMarker()
		if err != nil {
			a.err = err
			return false
		}
		return m != arrayEndMarker
	}
	return a.count < a.Len
}

// End completes array decoding and must be called. It may return errors (1)
// deferred from element decoding, (2) from a missing array end marker, or (3)
// from a length vs. count mismatch.
func (a *ArrayDecoder) End() error {
	if a.err != nil {
		return a.err
	}
	if a.Len < 0 {
		m, err := a.readMarker()
		if err != nil {
			return err
		}
		if m != arrayEndMarker {
			return errors.New("expected end marker")
		}
	} else if a.Len != a.count {
		return errors.New("len count mismatch")
	}
	return nil
}

// Decode decodes a value into v by delegating to the appropriate type-specific
// method. Recognizes the special types Char and HighPrecNumber to distinguish
// from backing types.
func (d *Decoder) Decode(v interface{}) error {
	if v == nil {
		return errors.New("cannot decode into nil value")
	}
	if val, ok := v.(Value); ok {
		return d.DecodeValue(val)
	}
	switch t := v.(type) {
	case *interface{}:
		i, err := d.decodeInterface()
		if err == nil {
			*t = i
		}
		return err

	case *string:
		s, err := d.DecodeString()
		if err == nil {
			*t = s
		}
		return err

	case *bool:
		b, err := d.DecodeBool()
		if err == nil {
			*t = b
		}
		return err

	case *int:
		i, err := d.DecodeInt()
		if err == nil {
			*t = i
		}
		return err

	case *uint8:
		u, err := d.DecodeUInt8()
		if err == nil {
			*t = u
		}
		return err

	case *int8:
		i, err := d.DecodeInt8()
		if err == nil {
			*t = i
		}
		return err

	case *int16:
		i, err := d.DecodeInt16()
		if err == nil {
			*t = i
		}
		return err

	case *int32:
		i, err := d.DecodeInt32()
		if err == nil {
			*t = i
		}
		return err

	case *int64:
		i, err := d.DecodeInt64()
		if err == nil {
			*t = i
		}
		return err

	case *float32:
		f, err := d.DecodeFloat32()
		if err == nil {
			*t = f
		}
		return err

	case *float64:
		f, err := d.DecodeFloat64()
		if err == nil {
			*t = f
		}
		return err

	case *Char:
		b, err := d.DecodeChar()
		if err == nil {
			*t = Char(b)
		}
		return err

	case *HighPrecNumber:
		s, err := d.DecodeHighPrecNumber()
		if err == nil {
			*t = HighPrecNumber(s)
		}
		return err
	}

	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return errors.Errorf("can only decode into pointers, not: %s", value.Type())
	}
	// Containers
	switch value.Elem().Kind() {
	case reflect.Array:
		return d.decodeValue(ArrayStartMarker, arrayToArray(value))
	case reflect.Slice:
		return d.decodeValue(ArrayStartMarker, arrayToSlice(value))

	case reflect.Map:
		if value.Elem().Kind() != reflect.String &&
			!value.Elem().Type().Key().ConvertibleTo(stringType) {
			return fmt.Errorf("unable to decode map of type %s: key type must be string", value.Type())
		}
		return d.decodeValue(ObjectStartMarker, objectIntoMap(value))

	case reflect.Struct:
		return d.decodeValue(ObjectStartMarker, objectIntoStruct(value))
	}

	return fmt.Errorf("unable to decode this type of value: %T %v", v, v)
}

// arrayToArray returns a function to decode an array container into
// arrayPtr.Elem(). Returns an error if the lengths are not equal.
func arrayToArray(arrayPtr reflect.Value) func(*Decoder) error {
	return func(d *Decoder) error {
		ad, err := d.Array()
		if err != nil {
			return err
		}
		arrayValue := arrayPtr.Elem()
		elemType := arrayValue.Type().Elem()
		if ad.Len > 0 {
			if ad.Len >= 0 && ad.Len != arrayValue.Len() {
				return errors.Errorf("unable to decode data length %d into array of length %d", ad.Len, arrayValue.Len())
			}
		}

		for i := 0; i < arrayValue.Len(); i++ {
			elemPtr := reflect.New(elemType)
			if err := ad.Decode(elemPtr.Interface()); err != nil {
				return err
			}
			arrayValue.Index(i).Set(elemPtr.Elem())
		}
		return ad.End()
	}
}

// arrayToSlice returns a function to decode an array container into
// slicePtr.Elem().
func arrayToSlice(slicePtr reflect.Value) func(*Decoder) error {
	return func(d *Decoder) error {
		ad, err := d.Array()
		if err != nil {
			return err
		}
		sliceValue := slicePtr.Elem()
		elemType := sliceValue.Type().Elem()
		if ad.Len < 0 {
			sliceValue.Set(reflect.MakeSlice(sliceValue.Type(), 0, 0))
			for ad.NextElem() {
				elemPtr := reflect.New(elemType)
				if err := ad.Decode(elemPtr.Interface()); err != nil {
					return err
				}
				sliceValue = reflect.Append(sliceValue, elemPtr.Elem())
			}
		} else {
			sliceValue.Set(reflect.MakeSlice(sliceValue.Type(), ad.Len, ad.Len))

			for i := 0; i < ad.Len; i++ {
				elemPtr := reflect.New(elemType)
				if err := ad.Decode(elemPtr.Interface()); err != nil {
					return err
				}
				sliceValue.Index(i).Set(elemPtr.Elem())
			}
		}
		return ad.End()
	}
}

var zeroValue = reflect.Value{}

func objectIntoStruct(structPtr reflect.Value) func(*Decoder) error {
	return func(d *Decoder) error {
		o, err := d.Object()
		if err != nil {
			return err
		}
		for o.NextEntry() {
			k, err := o.DecodeKey()
			if err != nil {
				return errors.Wrapf(err, "failed to decode key with call #%d", o.count)
			}

			f := structPtr.Elem().FieldByName(k)
			if f == zeroValue {
				return errors.Errorf("unable to decode entry: no field named %q found", k)
			}
			if err := o.Decode(f.Addr().Interface()); err != nil {
				return errors.Wrapf(err, "failed to decode value for %q with call #%d", k, o.count)
			}
		}
		return o.End()
	}
}

func objectIntoMap(mapPtr reflect.Value) func(*Decoder) error {
	return func(d *Decoder) error {
		o, err := d.Object()
		if err != nil {
			return err
		}

		mapValue := mapPtr.Elem()
		//TODO go1.9 - MakeMapWithSize
		mapValue.Set(reflect.MakeMap(mapValue.Type()))
		elemType := mapValue.Type().Elem()

		for o.NextEntry() {
			k, err := o.DecodeKey()
			if err != nil {
				return errors.Wrapf(err, "failed to decode key #%d", o.count)
			}

			valPtr := reflect.New(elemType)

			if err := o.Decode(valPtr.Interface()); err != nil {
				return errors.Wrapf(err, "failed to decode value #%d", o.count)
			}

			mapValue.SetMapIndex(reflect.ValueOf(k), valPtr.Elem())
		}
		return o.End()
	}
}

// objectAsInterface reads an object and returns a map[string]T where T is
// either interface{} or a stricter type if the object is strongly typed.
func objectAsInterface(o *ObjectDecoder) (interface{}, error) {
	valType := elementTypeFor(o.ValType)
	mapType := reflect.MapOf(stringType, valType)

	//TODO go1.9 - MakeMapWithSize
	mapValue := reflect.MakeMap(mapType)
	for o.NextEntry() {
		k, err := o.DecodeKey()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to decode key #%d", o.count)
		}
		valPtr := reflect.New(mapType.Elem())
		if err := o.Decode(valPtr.Interface()); err != nil {
			return nil, errors.Wrapf(err, "failed to decode value #%d", o.count)
		}

		mapValue.SetMapIndex(reflect.ValueOf(k), valPtr.Elem())
	}
	return mapValue.Interface(), o.End()
}

// arrayAsInterface reads an array container into a new slice []T, where T may
// be strongly typed, or an interface{} in the general case.
func arrayAsInterface(a *ArrayDecoder) (interface{}, error) {
	var sliceValue reflect.Value
	elemType := elementTypeFor(a.ElemType)
	sliceType := reflect.SliceOf(elemType)

	if a.Len < 0 {
		sliceValue = reflect.MakeSlice(sliceType, 0, 0)
		for a.NextElem() {
			elemPtr := reflect.New(elemType)
			if err := a.Decode(elemPtr.Interface()); err != nil {
				return nil, err
			}
			sliceValue = reflect.Append(sliceValue, elemPtr.Elem())
		}
	} else {
		sliceValue = reflect.MakeSlice(sliceType, a.Len, a.Len)

		for i := 0; i < a.Len; i++ {
			elemPtr := reflect.New(elemType)
			if err := a.Decode(elemPtr.Interface()); err != nil {
				return nil, err
			}
			sliceValue.Index(i).Set(elemPtr.Elem())
		}
	}

	if err := a.End(); err != nil {
		return nil, err
	}

	return sliceValue.Interface(), nil
}

var iface interface{}
var ifaceType = reflect.TypeOf(&iface).Elem()

// elementTypeFor returns the type into which data with this marker should be
// decoded, falling back to interface{} in the general case.
func elementTypeFor(m Marker) reflect.Type {
	switch m {
	case TrueMarker, FalseMarker:
		return boolType
	case UInt8Marker:
		return uint8Type
	case Int8Marker:
		return int8Type
	case Int16Marker:
		return int16Type
	case Int32Marker:
		return int32Type
	case Int64Marker:
		return int64Type
	case Float32Marker:
		return float32Type
	case Float64Marker:
		return float64Type
	case StringMarker:
		return stringType
	case CharMarker:
		return charType
	case HighPrecNumMarker:
		return highPrecNumType
	}
	return ifaceType
}

var (
	boolType        = reflect.TypeOf(true)
	uint8Type       = reflect.TypeOf(uint8(0))
	int8Type        = reflect.TypeOf(int8(0))
	int16Type       = reflect.TypeOf(int16(0))
	int32Type       = reflect.TypeOf(int32(0))
	int64Type       = reflect.TypeOf(int64(0))
	float32Type     = reflect.TypeOf(float32(3.14))
	float64Type     = reflect.TypeOf(float64(3.14))
	stringType      = reflect.TypeOf("")
	charType        = reflect.TypeOf(Char(0))
	highPrecNumType = reflect.TypeOf(HighPrecNumber(0))
)
