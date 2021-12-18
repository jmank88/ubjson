package ubjson

import (
	"fmt"
	"reflect"
	"testing"
	"unicode/utf8"
)

func TestMarshal(t *testing.T) {
	t.Parallel()
	for name, testCase := range cases {
		t.Run(name, testCase.marshal)
	}
}

func (tc *testCase) marshal(t *testing.T) {
	if b, err := Marshal(tc.value); err != nil {
		t.Fatal("failed to marshal:", err.Error())
	} else if diff := firstBytesDiff(tc.binary, b); diff.index != -1 {
		t.Errorf("(%T) %v:\n %s\n expected:\n %#v\n\n  but got:\n %#v\n\n",
			tc.value, tc.value, diff, tc.binary, b)
	}
}

func TestMarshalBlock(t *testing.T) {
	t.Parallel()
	for name, testCase := range cases {
		t.Run(name, testCase.marshalBlock)
	}
}

func (tc *testCase) marshalBlock(t *testing.T) {
	if b, err := MarshalBlock(tc.value); err != nil {
		t.Fatal("failed to marshal block:", err.Error())
	} else if diff := firstStringDiff(tc.block, string(b)); diff != noStringDiff {
		t.Errorf("(%T) %v:\n %s\n expected:\n %q\n\n  but got:\n %q\n\n",
			tc.value, tc.value, diff, tc.block, string(b))
	}
}

var noBytesDiff = bytesDiff{index: -1}

type bytesDiff struct {
	index        int
	lens         bool
	b1           byte
	b1Pre, b1Suf []byte
	b2           byte
	b2Pre, b2Suf []byte
}

func (d bytesDiff) String() string {
	if d.index == -1 {
		return "no difference"
	} else if d.lens {
		return "different lengths"
	}
	return fmt.Sprintf(bytesDiffF, d.index, d.b1Pre, d.b1, d.b1Suf, d.b2Pre, d.b2, d.b2Suf)
}

var bytesDiffF = "first diff at index (%d);\n" +
	"%x[%x]%x\n" +
	"%x[%x]%x\n"

const diffPreSufLen = 8

func firstBytesDiff(b1, b2 []byte) bytesDiff {
	var diff bytesDiff
	var has1, has2 bool
	for {
		if has1 = diff.index < len(b1); has1 {
			diff.b1 = b1[diff.index]

			if diff.index < diffPreSufLen {
				diff.b1Pre = b1[:diff.index]
			} else {
				diff.b1Pre = b1[diff.index-diffPreSufLen : diff.index]
			}
			if diff.index+1+diffPreSufLen < len(b1) {
				diff.b1Suf = b1[diff.index+1 : diff.index+1+diffPreSufLen]
			} else {
				diff.b1Suf = b1[diff.index+1:]
			}
		}
		if has2 = diff.index < len(b2); has2 {
			diff.b2 = b2[diff.index]

			if diff.index < diffPreSufLen {
				diff.b2Pre = b2[:diff.index]
			} else {
				diff.b2Pre = b2[diff.index-diffPreSufLen : diff.index]
			}
			if diff.index+1+diffPreSufLen < len(b2) {
				diff.b2Suf = b2[diff.index+1 : diff.index+1+diffPreSufLen]
			} else {
				diff.b2Suf = b2[diff.index+1:]
			}
		}

		if !has1 && !has2 {
			return noBytesDiff
		} else if !has1 || !has2 {
			diff.lens = true
			return diff
		}
		// has1 && has2 == true
		if diff.b1 != diff.b2 {
			return diff
		}
		diff.index++
	}
}

var noStringDiff = stringDiff{index: -1}

type stringDiff struct {
	index int
	rune  int
	lens  bool
	r1    rune
	r2    rune
}

func (d stringDiff) String() string {
	if d.index == -1 {
		return "no difference"
	} else if d.lens {
		return "different lengths"
	}
	return fmt.Sprintf("first diff at index (%d) rune (%d); %q vs. %q", d.index, d.rune, d.r1, d.r2)
}

func firstStringDiff(s1, s2 string) stringDiff {
	var diff stringDiff
	var has1, has2 bool
	var l int
	for {
		if has1 = diff.index < len(s1); has1 {
			diff.r1, l = utf8.DecodeRuneInString(s1[diff.index:])
		}
		if has2 = diff.index < len(s2); has2 {
			diff.r2, _ = utf8.DecodeRuneInString(s2[diff.index:])
		}

		if !has1 && !has2 {
			return noStringDiff
		} else if !has1 || !has2 {
			diff.lens = true
			return diff
		}
		// has1 && has2 == true
		if diff.r1 != diff.r2 {
			return diff
		}
		diff.index += l
		diff.rune++
	}
}

func Test_elementMarkerFor(t *testing.T) {
	for _, tt := range []struct {
		v interface{}
		m Marker
	}{
		{uint8(1), UInt8Marker},
		{int8(2), Int8Marker},
		{int16(3), Int16Marker},
		{int32(4), Int32Marker},
		{int64(5), Int64Marker},
		{float32(1.1), Float32Marker},
		{float64(2.2), Float64Marker},
		{HighPrecNumber(""), HighPrecNumMarker},
		{Char('a'), CharMarker},
		{"", StringMarker},
	} {
		tt := tt
		t.Run(tt.m.String(), func(t *testing.T) {
			if got := elementMarkerFor(reflect.TypeOf(tt.v)); got != tt.m {
				t.Errorf("expected %s for %v but got %s", tt.m, tt.v, got)
			}
		})
	}
}
