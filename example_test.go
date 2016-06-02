package ubjson_test

import (
	"github.com/jmank88/ubjson"
	"os"
)

func ExampleEncode() {
	_ = ubjson.NewEncoder(os.Stdout).Block().Encode(8)

	// Output:
	// [U][8]
}

//TODO more examples
