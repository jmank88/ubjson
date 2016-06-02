package ubjson

import ()

//TODO making these stricter enums could be beneficial
//TODO can use an int8, and then have an array of marker: value for lookup
const (
	// Value Types
	nullMarker = 'Z'

	noOpMarker = 'N'

	trueMarker  = 'T'
	falseMarker = 'F'

	int8Marker  = 'i'
	uInt8Marker = 'U'

	int16Marker = 'I'
	int32Marker = 'l'
	int64Marker = 'L'

	float32Marker = 'd'
	float64Marker = 'D'

	highPrecNumMarker = 'H'

	charMarker   = 'C'
	stringMarker = 'S'

	// Container Types
	arrayStartMarker = '['
	arrayEndMarker   = ']'

	objectStartMarker = '{'
	objectEndMarker   = '}'

	// Optimized Params
	countMarker = '#'
	typeMarker  = '$'
)

type Token interface {
	Marker() byte
	Length() int
	Data() []byte
}

var (
	Null  Token = &constToken{marker: nullMarker}
	NoOp  Token = &constToken{marker: noOpMarker}
	True  Token = &constToken{marker: trueMarker}
	False Token = &constToken{marker: falseMarker}
)

type constToken struct {
	marker byte
	length int
	data   []byte
}

func (t *constToken) Marker() byte {
	return t.marker
}

func (t *constToken) Length() int {
	return t.length
}

func (t *constToken) Data() []byte {
	return t.data
}

type Parser interface {
	scan() (Token, error)
	decode(into interface{}) error
}

//type parser struct {
//	scanner
//}
//
////TODO doc
//func (p *parser) scan() (Token, error) {
//	t.marker = p.marker()
//	switch t.marker {
//	case stringMarker, charMarker:
//		t.length = p.scanner.length()
//	case arrayStartMarker, objectStartMarker:
//		//TODO length optional
//		t.length = 2
//	}
//
//	switch t.marker {
//	case int8Marker, uInt8Marker, int16Marker, int32Marker, int64Marker, float32Marker, float64Marker, highPrecNumMarker, charMarker, stringMarker, arrayStartMarker, objectStartMarker:
//		//TODO data
//		t.data = p.data()
//	}
//	return t, nil
//}
//
//func (p *parser) decode(into interface{}) error {
//	t := reflect.Value(into)
//	switch t.Kind() {
//	case reflect.Struct:
//		return d.Parser.parseStruct(into)
//	case reflect.String:
//		return d.Parser.parseStr(t.String())
//		//TODO all the kinds
//
//	}
//
//	d.Parser.parseObject()
//}
