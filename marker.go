package ubjson

// A Marker is a single byte UBJSON valType byte.
type Marker byte

func (m Marker) String() string {
	return string(m)
}

const (
	// Value Type Markers

	NullMarker Marker = 'Z'

	NoOpMarker Marker = 'N'

	TrueMarker  Marker = 'T'
	FalseMarker Marker = 'F'

	UInt8Marker Marker = 'U'

	Int8Marker  Marker = 'i'
	Int16Marker Marker = 'I'
	Int32Marker Marker = 'l'
	Int64Marker Marker = 'L'

	Float32Marker Marker = 'd'
	Float64Marker Marker = 'D'

	HighPrecNumMarker Marker = 'H'

	CharMarker   Marker = 'C'
	StringMarker Marker = 'S'
)

const (
	// Container Types Markers

	ArrayStartMarker  Marker = '['
	ObjectStartMarker Marker = '{'
)

const (
	// Container Meta-Markers

	arrayEndMarker  Marker = ']'
	objectEndMarker Marker = '}'
	countMarker     Marker = '#'
	typeMarker      Marker = '$'
)

// The smallestIntMarker function returns the Marker for the smallest integer
// into which v will fit.
func smallestIntMarker(v int64) Marker {
	switch {
	case v >= 0 && v <= 255:
		return UInt8Marker
	case v <= 127 && v >= -128:
		return Int8Marker
	case v <= 32767 && v >= -32768:
		return Int16Marker
	case v <= 2147483647 && v >= -2147483648:
		return Int32Marker
	default:
		return Int64Marker
	}
}
