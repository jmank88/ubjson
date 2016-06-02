package ubjson

import "io"

//TODO test this
// Pretty reads ubj from r and writes to w in human readable form.
func Block(r io.Reader, w io.Writer) error {
	p := &block{
		Reader: r,
		Writer: w,
	}
	return p.transcode()
}

type block struct {
	io.Reader
	io.Writer
}

func (p *block) transcode() error {
	b := make([]byte, 1)
	_, err := p.Read(b[:])
	if err != nil {
		return err
	}
	_, err = p.Write([]byte{'[', b[0], ']'})
	if err != nil {
		return err
	}
	switch b[0] {
	case objectStartMarker:
		return p.object()
	case arrayStartMarker:
		return p.object()
		//TODO more
	}

	return nil
}

func (p *block) object() error {
	return nil //TODO
}

func (p *block) array() error {
	return nil //TODO
}

func (p *block) int8() error {
	return nil //TODO
}
