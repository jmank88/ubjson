package ubjson

import "io"

type HighPrecNumber string

func (h HighPrecNumber) MarshalUBJSON() (byte, func(*Encoder) error) {
	return highPrecNumMarker, func(e *Encoder) error {
		if err := e.EncodeInt(len(h)); err != nil {
			return err
		}
		if len(h) > 0 {
			if e.block {
				return e.blocked(string(h))
			}
			if _, err := io.WriteString(e.w, string(h)); err != nil {
				return err
			}
		}
		return nil
	}
}

type Char byte

func (c Char) MarshalUBJSON() (byte, func(*Encoder) error) {
	return charMarker, func(e *Encoder) error {
		if e.block {
			return e.blocked(string(c))
		}
		return e.writeByte(byte(c))
	}
}

//TODO time/date?
