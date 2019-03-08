package ubjson

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/pkg/errors"
)

// A reader reads UBJSON types.
type reader interface {
	// Peek at the next Marker without advancing.
	peekMarker() (Marker, error)
	// Read the next marker and advance.
	readMarker() (Marker, error)

	readUInt8() (uint8, error)
	readInt8() (int8, error)
	readInt16() (int16, error)
	readInt32() (int32, error)
	readInt64() (int64, error)

	readFloat32() (float32, error)
	readFloat64() (float64, error)

	readString(max int) (string, error)
	readChar() (byte, error)
}

// The readInt function dynamically reads an integer of unspecified size.
func readInt(r reader) (int, error) {
	m, err := r.readMarker()
	if err != nil {
		return 0, err
	}
	switch m {
	case UInt8Marker:
		u, err := r.readUInt8()
		return int(u), err
	case Int8Marker:
		i, err := r.readInt8()
		return int(i), err
	case Int16Marker:
		i, err := r.readInt16()
		return int(i), err
	case Int32Marker:
		i, err := r.readInt32()
		return int(i), err
	case Int64Marker:
		i, err := r.readInt64()
		return int(i), err
	}
	return 0, fmt.Errorf("failed to read int: expected int marker but found %q", m)
}

// The readContainer method parses and returns a container type marker and
// length, or 0 and -1 respectively when none are found.
func readContainer(r reader) (Marker, int, error) {
	m, err := r.peekMarker()
	if err != nil {
		return 0, 0, err
	}

	switch m {
	case typeMarker:
		if _, err := r.readMarker(); err != nil {
			return 0, -1, err
		}
		m, err := r.readMarker()
		if err != nil {
			return 0, 0, err
		}

		if c, err := r.readMarker(); err != nil {
			return 0, 0, err
		} else if c != countMarker {
			return 0, 0, errors.New("count marker (#) required following container type marker")
		}
		l, err := readInt(r)
		if err != nil {
			return 0, 0, err
		}
		if l < 0 {
			return 0, 0, errors.Errorf("illegal negative container length: %d", l)
		}
		return m, l, nil

	case countMarker:
		if _, err := r.readMarker(); err != nil {
			return 0, -1, err
		}
		l, err := readInt(r)
		if err != nil {
			return 0, 0, err
		}
		if l < 0 {
			return 0, 0, errors.Errorf("illegal negative container length: %d", l)
		}
		return 0, l, nil

	default:
		return 0, -1, nil
	}
}

// A binaryReader reads binary UBJSON.
type binaryReader struct {
	*bufio.Reader
	// A buffer as large the largest fixed size type.
	buf [8]byte
}

func newBinaryReader(r io.Reader) *binaryReader {
	return &binaryReader{Reader: bufio.NewReader(r)}
}

func (r *binaryReader) readMarker() (Marker, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, errors.New("failed to read marker")
	}
	return Marker(b), nil
}

func (r *binaryReader) peekMarker() (Marker, error) {
	b, err := r.Peek(1)
	if err != nil {
		return 0, errors.New("failed to peek marker")
	}
	return Marker(b[0]), nil
}

func (r *binaryReader) readUInt8() (uint8, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, errors.New("failed to read UInt8 byte")
	}
	return uint8(b), nil
}

func (r *binaryReader) readInt8() (int8, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, errors.New("failed to read Int8 byte")
	}
	return int8(b), nil
}

// The readBuf method reads len bytes into r.buf. len must not exceed 8.
func (r *binaryReader) readBuf(len int) ([]byte, error) {
	b := r.buf[:len]
	n, err := r.Read(b)
	if err != nil {
		return nil, err
	}
	if n != len {
		return nil, errors.New("TODO failed to read enough bytes")
	}
	return b, nil
}

func (r *binaryReader) readInt16() (int16, error) {
	b, err := r.readBuf(2)
	if err != nil {
		return 0, err
	}
	u := binary.BigEndian.Uint16(b)
	return int16(u), err
}

func (r *binaryReader) readInt32() (int32, error) {
	b, err := r.readBuf(4)
	if err != nil {
		return 0, err
	}
	u := binary.BigEndian.Uint32(b)
	return int32(u), err
}

func (r *binaryReader) readInt64() (int64, error) {
	b, err := r.readBuf(8)
	if err != nil {
		return 0, err
	}
	u := binary.BigEndian.Uint64(b)
	return int64(u), err
}

func (r *binaryReader) readFloat32() (float32, error) {
	b, err := r.readBuf(4)
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(binary.BigEndian.Uint32(b)), err
}

func (r *binaryReader) readFloat64() (float64, error) {
	b, err := r.readBuf(8)
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(binary.BigEndian.Uint64(b)), err
}

func (r *binaryReader) readString(max int) (string, error) {
	l, err := readInt(r)
	if err != nil {
		return "", errors.Wrap(err, "failed to read string length prefix")
	}
	if l < 0 {
		return "", errors.Errorf("illegal string length prefix: %d", l)
	}
	if l > max {
		return "", errors.Errorf("string length prefix exceeds max allocation limit of %d: %d", max, l)
	}
	b := make([]byte, l)
	n, err := r.Read(b)
	if err != nil {
		return "", errors.Wrap(err, "failed to read string bytes")
	}
	if n != l {
		return "", fmt.Errorf("failed to read full string length (%d), instead got: %s", l, string(b[:n]))
	}
	return string(b), nil
}

func (r *binaryReader) readChar() (byte, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, errors.New("failed to read marker")
	}
	if b > 127 {
		return 0, fmt.Errorf("illegal Char value %d: must not exceed 127", b)
	}
	return b, nil
}

// A blockReader reads block-notation UBJSON.
type blockReader struct {
	*bufio.Reader
	next string
}

func newBlockReader(r io.Reader) *blockReader {
	return &blockReader{Reader: bufio.NewReader(r)}
}

// The nextBlock method returns the next block, which may be cached or read
// fresh.
func (r *blockReader) nextBlock() (string, error) {
	if r.next != "" {
		n := r.next
		r.next = ""
		return n, nil
	}
	return r.readBlock()
}

// The readBlock method reads the next block.
func (r *blockReader) readBlock() (string, error) {
	if _, err := r.ReadBytes('['); err != nil {
		return "", errors.Wrap(err, "failed to read block start")
	}
	s, err := r.ReadString(']')
	if err != nil {
		return "", errors.New("failed to read through block end")
	}
	return s[:len(s)-1], nil
}

// The peekBlock method returns the next block, but caches it for the next read.
func (r *blockReader) peekBlock() (string, error) {
	if r.next != "" {
		return r.next, nil
	}
	n, err := r.readBlock()
	if err == nil {
		if n == "" {
			return "", errors.New("empty block")
		}
		r.next = n
	}
	return n, err
}

func (r *blockReader) readMarker() (Marker, error) {
	s, err := r.nextBlock()
	if err != nil {
		return 0, err
	}
	if len(s) != 1 {
		return 0, fmt.Errorf("expected single byte marker, but found: %q", s)
	}
	return Marker(s[0]), nil
}

func (r *blockReader) peekMarker() (Marker, error) {
	s, err := r.peekBlock()
	if err != nil {
		return 0, err
	}
	if len(s) > 1 {
		return 0, fmt.Errorf("expected single byte marker, but found %q", s)
	}
	return Marker(s[0]), nil
}

func (r *blockReader) readUInt8() (uint8, error) {
	s, err := r.nextBlock()
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to parse uint8")
	}
	return uint8(i), nil
}

func (r *blockReader) readInt8() (int8, error) {
	i, err := r.readBlockedInt(8)
	return int8(i), err
}

func (r *blockReader) readInt16() (int16, error) {
	i, err := r.readBlockedInt(16)
	return int16(i), err
}

func (r *blockReader) readBlockedInt(bitSize int) (int64, error) {
	s, err := r.nextBlock()
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(s, 10, bitSize)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to parse int%d", bitSize)
	}
	return i, nil
}

func (r *blockReader) readInt32() (int32, error) {
	i, err := r.readBlockedInt(32)
	return int32(i), err
}

func (r *blockReader) readInt64() (int64, error) {
	i, err := r.readBlockedInt(64)
	return i, err
}

func (r *blockReader) readFloat32() (float32, error) {
	s, err := r.nextBlock()
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(s, 32)
	return float32(f), err
}

func (r *blockReader) readFloat64() (float64, error) {
	s, err := r.nextBlock()
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(s, 64)
	return f, err
}

func (r *blockReader) readString(max int) (string, error) {
	l, err := readInt(r)
	if err != nil {
		return "", errors.Wrap(err, "failed to read string length prefix")
	}
	if l < 1 {
		return "", errors.Errorf("illegal string length prefix: %d", l)
	}
	if l > max {
		return "", errors.Errorf("string length prefix exceeds max allocation limit of %d: %d", max, l)
	}
	s, err := r.nextBlock()
	if err != nil {
		return "", errors.Wrap(err, "failed to read string block")
	}
	if len(s) != l {
		return "", fmt.Errorf("string length %d but prefix %d", len(s), l)
	}
	return s, nil
}

func (r *blockReader) readChar() (byte, error) {
	s, err := r.nextBlock()
	if err != nil {
		return 0, err
	}
	if len(s) != 1 {
		return 0, fmt.Errorf("expected single byte Char, but found: %q", s)
	}
	b := s[0]
	if b > 127 {
		return 0, fmt.Errorf("illegal Char value %d: must not exceed 127", b)
	}
	return b, nil
}
