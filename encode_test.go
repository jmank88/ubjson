package ubjson

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

func TestMarshal(t *testing.T) {
	t.Parallel()
	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			b, err := Marshal(testCase.value)
			if err != nil {
				t.Fatal("failed to marshal:", err.Error())
			}
			if diff := firstBytesDiff(testCase.binary, b); diff.index != -1 {
				t.Errorf("(%T) %v:\n %s\n expected:\n %#v\n\n  but got:\n %#v\n\n",
					testCase.value, testCase.value, diff, testCase.binary, b)
			}
		})
	}
}

func TestMarshalBlock(t *testing.T) {
	t.Parallel()
	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			if b, err := MarshalBlock(testCase.value); err != nil {
				t.Fatal("failed to marshal block:", err.Error())
			} else if diff := firstStringDiff(testCase.block, string(b)); diff != noStringDiff {
				t.Errorf("(%T) %v:\n %s\n expected:\n %q\n\n  but got:\n %q\n\n",
					testCase.value, testCase.value, diff, testCase.block, string(b))
			}
		})
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
