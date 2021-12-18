package ubjson

import "github.com/pkg/errors"

func errTooMany(len int) error {
	return errors.Errorf("too many calls for container with len %d", len)
}

func errWrongTypeWrite(exp, got Marker) error {
	return errors.Errorf("unable to write element type '%s' to container type '%s'", got, exp)
}

func errWrongTypeRead(exp, got Marker) error {
	return errors.Errorf("tried to read type '%s' but found type '%s'", exp, got)
}
