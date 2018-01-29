// +build gofuzz

package ubjson

// example run: go-fuzz -bin=fuzz-bin.zip -workdir=testdata/bin

//go:generate go-fuzz-build -func FuzzUnmarshal -o fuzz-bin.zip github.com/jmank88/ubjson

func FuzzUnmarshal(data []byte) int {
	var i interface{}
	if Unmarshal(data, &i) != nil {
		return 0
	}
	return 1
}

//go:generate go-fuzz-build -func FuzzUnmarshalBlock -o fuzz-block.zip github.com/jmank88/ubjson

func FuzzUnmarshalBlock(data []byte) int {
	var i interface{}
	if UnmarshalBlock(data, &i) != nil {
		return 0
	}
	return 1
}
