package ubjson

import "fmt"

func errTooMany(len int) error {
	return fmt.Errorf("too many calls for container with len %d", len)
}

func errWrongTypeWrite(exp, got Marker) error {
	return fmt.Errorf("unable to write element type '%s' to container type '%s'", got, exp)
}

func errWrongTypeRead(exp, got Marker) error {
	return fmt.Errorf("tried to read type '%s' but found type '%s'", exp, got)
}
