//go:build go1.18
// +build go1.18

package ubjson

import (
	"embed"
	"os"
	"testing"
)

//go:embed testdata/**/*
var testdata embed.FS

func FuzzBinary(f *testing.F) {
	for _, c := range cases {
		f.Add(c.binary)
	}
	const corpus = "testdata/bin/corpus"
	if dir, err := testdata.ReadDir(corpus); err != nil {
		f.Fatal("failed to read corpus dir:", err)
	} else {
		for _, c := range dir {
			if c.IsDir() {
				continue
			}
			if b, err := os.ReadFile(corpus + "/" + c.Name()); err != nil {
				f.Error("failed to open corpus file:", err)
			} else {
				f.Add(b)
			}
		}
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		var i interface{}
		if Unmarshal(data, &i) != nil {
			t.Skip()
		}
		if _, err := Marshal(i); err != nil {
			t.Errorf("failed to re-marshal: %+v\n", err)
			t.Logf("original: 0x%x\n", data)
			t.Logf("value: %#v\n", i)
			bl, err := MarshalBlock(i)
			if err != nil {
				t.Error("block: failed to marshal:", err)
			} else {
				t.Log("block:", bl)
			}
		}
	})
}

func FuzzBlock(f *testing.F) {
	for _, c := range cases {
		f.Add(c.block)
	}
	const corpus = "testdata/block/corpus"
	if dir, err := testdata.ReadDir(corpus); err != nil {
		f.Fatal("failed to read corpus dir:", err)
	} else {
		for _, c := range dir {
			if c.IsDir() {
				continue
			}
			if b, err := os.ReadFile(corpus + "/" + c.Name()); err != nil {
				f.Error("failed to open corpus file:", err)
			} else {
				f.Add(string(b))
			}
		}
	}
	f.Fuzz(func(t *testing.T, data string) {
		var i interface{}
		if UnmarshalBlock([]byte(data), &i) != nil {
			t.Skip()
		}
		if _, err := MarshalBlock(i); err != nil {
			t.Errorf("failed to re-marshal: %+v\n", err)
			t.Log("original:", data)
			t.Logf("value: %#v\n", i)
		}
	})
}
