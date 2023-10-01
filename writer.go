package ubjson

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
)

type writer interface {
	// The writeMarker method writes a Marker byte.
	writeMarker(Marker) error
	writeUInt8(uint8) error
	writeInt8(int8) error
	writeInt16(int16) error
	writeInt32(int32) error
	writeInt64(int64) error
	writeFloat32(float32) error
	writeFloat64(float64) error
	// Writes a length-prefixed UBJSON string.
	writeString(string) error
	// Writes a UBJSON Char, which must be <=127.
	writeChar(byte) error
	Flush() error
	writeNewLine() error
	incIndent()
	decIndent()
}

// A blockWriter is a writer which writes block-notation UBJSON.
type blockWriter struct {
	*bufio.Writer
	// Current number of indentations.
	indent int
}

// The newBlockWriter function returns a new block-notation writer.
func newBlockWriter(w io.Writer) *blockWriter {
	return &blockWriter{Writer: bufio.NewWriter(w)}
}

// The writeBlocked method writes s surrounded by square brackets.
func (w *blockWriter) writeBlocked(s string) error {
	_, err := fmt.Fprintf(w, "[%s]", s)
	return err
}

func (w *blockWriter) incIndent() {
	w.indent++
}

func (w *blockWriter) decIndent() {
	w.indent--
}

// The writeNewLine method starts a new, indented line.
func (w *blockWriter) writeNewLine() error {
	if err := w.writeByte('\n'); err != nil {
		return err
	}
	for i := 0; i < w.indent; i++ {
		if err := w.writeByte('\t'); err != nil {
			return err
		}
	}
	return nil
}

func (w *blockWriter) write(p []byte) error {
	_, err := w.Write(p)
	return err
}

// The writeMarker method writes a Marker byte, respecting block notation
func (w *blockWriter) writeMarker(m Marker) error {
	return w.write([]byte{'[', byte(m), ']'})
}

func (w *blockWriter) writeByte(c byte) error { return w.WriteByte(c) }

func (w *blockWriter) writeUInt8(v uint8) error {
	return w.writeBlocked(strconv.FormatUint(uint64(v), 10))
}

func (w *blockWriter) writeInt8(v int8) error {
	return w.writeBlocked(strconv.FormatInt(int64(v), 10))
}

func (w *blockWriter) writeInt16(v int16) error {
	return w.writeBlocked(strconv.FormatInt(int64(v), 10))
}

func (w *blockWriter) writeInt32(v int32) error {
	return w.writeBlocked(strconv.FormatInt(int64(v), 10))
}

func (w *blockWriter) writeInt64(v int64) error {
	return w.writeBlocked(strconv.FormatInt(int64(v), 10))
}

func (w *blockWriter) writeFloat32(v float32) error {
	return w.writeBlocked(strconv.FormatFloat(float64(v), 'g', -1, 32))
}

func (w *blockWriter) writeFloat64(v float64) error {
	return w.writeBlocked(strconv.FormatFloat(float64(v), 'g', -1, 32))
}

func (w *blockWriter) writeChar(v byte) error {
	if v > 127 {
		return fmt.Errorf("illegal char value (%d): cannot exceed 127", v)
	}
	return w.write([]byte{'[', byte(v), ']'})
}

// The writeString method writes a length-prefixed UBJSON string.
func (w *blockWriter) writeString(s string) error {
	if err := writeInt(w, len(s)); err != nil {
		return fmt.Errorf("failed writing string lenth prefix: %w", err)
	}

	if len(s) > 0 {
		return w.writeBlocked(s)
	}
	return nil
}

// A binaryWriter is a writer which writes binary UBJSON.
type binaryWriter struct {
	*bufio.Writer

	// A buffer as large as the largest fixed size type.
	buf [8]byte
}

func newBinaryWriter(w io.Writer) *binaryWriter {
	return &binaryWriter{Writer: bufio.NewWriter(w)}
}

func (w *binaryWriter) writeNewLine() error { return nil }

func (w *binaryWriter) incIndent() {}
func (w *binaryWriter) decIndent() {}

func (w *binaryWriter) write(p []byte) error {
	_, err := w.Write(p)
	return err
}

func (w *binaryWriter) writeMarker(m Marker) error {
	return w.writeByte(byte(m))
}

func (w *binaryWriter) writeByte(c byte) error { return w.WriteByte(c) }

func (w *binaryWriter) writeUInt8(v uint8) error {
	return w.writeByte(v)
}

func (w *binaryWriter) writeInt8(v int8) error {
	return w.writeByte(uint8(v))
}

func (w *binaryWriter) writeInt16(v int16) error {
	b := w.buf[:2]
	binary.BigEndian.PutUint16(b, uint16(v))
	return w.write(b)
}

func (w *binaryWriter) writeInt32(v int32) error {
	b := w.buf[:4]
	binary.BigEndian.PutUint32(b, uint32(v))
	return w.write(b)
}

func (w *binaryWriter) writeInt64(v int64) error {
	b := w.buf[:8]
	binary.BigEndian.PutUint64(b, uint64(v))
	return w.write(b)
}

func (w *binaryWriter) writeFloat32(v float32) error {
	b := w.buf[:4]
	binary.BigEndian.PutUint32(b, math.Float32bits(float32(v)))
	return w.write(b)
}

func (w *binaryWriter) writeFloat64(v float64) error {
	b := w.buf[:8]
	binary.BigEndian.PutUint64(b, math.Float64bits(float64(v)))
	return w.write(b)
}

func (w *binaryWriter) writeChar(v byte) error {
	if v > 127 {
		return fmt.Errorf("illegal char value (%d): cannot exceed 127", v)
	}
	return w.writeByte(byte(v))
}

func (w *binaryWriter) writeString(s string) error {
	if err := writeInt(w, len(s)); err != nil {
		return fmt.Errorf("failed writing string lenth prefix: %w", err)
	}

	if len(s) > 0 {
		_, err := io.WriteString(w, s)
		return err
	}
	return nil
}

// The writeInt method writes an integer with the smallest type possible, like
// EncodeInt, but *always* writes the leading type marker. This is used for
// string length-prefix metadata, which must never have its type marker
// optimized away, as a container element or value might.
func writeInt(w writer, v int) error {
	m := smallestIntMarker(int64(v))
	if err := w.writeMarker(m); err != nil {
		return err
	}
	switch m {
	case UInt8Marker:
		return w.writeUInt8(uint8(v))
	case Int8Marker:
		return w.writeInt8(int8(v))
	case Int16Marker:
		return w.writeInt16(int16(v))
	case Int32Marker:
		return w.writeInt32(int32(v))
	case Int64Marker:
		return w.writeInt64(int64(v))
	}
	return errors.New("TODO unreachable, programmere marker error")
}
