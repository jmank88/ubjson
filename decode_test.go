package ubjson

import (
	"reflect"
	"strings"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	for name, tc := range cases {
		t.Run(name, tc.unmarshal)
	}
}

func (tc *testCase) unmarshal(t *testing.T) {
	var expected interface{} = tc.value
	actual := reflect.New(reflect.ValueOf(tc.value).Type())

	if err := Unmarshal(tc.binary, actual.Interface()); err != nil {
		t.Fatalf("failed to unmarshal: %+v\n", err)
	}
	if !reflect.DeepEqual(actual.Elem().Interface(), expected) {
		t.Errorf("\nexpected: %T %v \nbut got:  %T %v",
			expected, expected, actual.Elem().Interface(), actual.Elem().Interface())
	}
}

func TestUnmarshalBlock(t *testing.T) {
	t.Parallel()
	for name, tc := range cases {
		t.Run(name, tc.unmarshalBlock)
	}
}

func (tc *testCase) unmarshalBlock(t *testing.T) {
	var expected interface{} = tc.value
	actual := reflect.New(reflect.ValueOf(tc.value).Type())

	if err := UnmarshalBlock([]byte(tc.block), actual.Interface()); err != nil {
		t.Fatal("failed to unmarshal block:", err.Error())
	}
	if !reflect.DeepEqual(actual.Elem().Interface(), expected) {
		t.Errorf("\nexpected: %T %#v \nbut got:  %T %#v",
			expected, expected, actual.Elem().Interface(), actual.Elem().Interface())
	}
}

func TestUnmarshalDiscardUnknownFields(t *testing.T) {
	type val struct{ A int8 }

	exp := val{8}
	var got val

	bin := []byte{'{', 'U', 0x01, 'A', 'i', 0x08, 'U', 0x01, 'b', 'i', 0x05, '}'}

	if err := Unmarshal(bin, &got); err != nil {
		t.Fatal(err)
	} else if got != exp {
		t.Errorf("\nexpected: %T %v \nbut got:  %T %v", exp, exp, got, got)
	}

	block := "[{]\n\t[U][1][A][i][8]\n\t[U][1][B][i][5]\n[}]"

	if err := UnmarshalBlock([]byte(block), &got); err != nil {
		t.Fatal(err)
	} else if got != exp {
		t.Errorf("\nexpected: %T %v \nbut got:  %T %v", exp, exp, got, got)
	}
}

func TestFuzzUnmarshalCrashers(t *testing.T) {
	for _, data := range []string{
		"[$F#i\x8a\x98b\x82ϟ6/\x9b\"\xe4\x88\xe8\xf0\xe0\f1A",
		"{l[ca[l1ca[ll[ca[l[ca[lcca[l[caP",
		"{$F#l2y2pY_9A__9y8vqOcl8Vxz9_Lu_2_wl8o4EMgH7T_3yDa8aS05Q17_YMAQHnwZfbccI_5c4",
		"[#lL00U",
		"{#l0000",
		"[l[ca[lca[l[[#l[ca[l",
		"{#ll[/\xfa\x00\x00\xfa\x80\xff\xff\xff\x01U",
		"[[#U\x01[#U\x01[#U\x01[#U\x01[#U\x01[#lca[l",
		"SlS\xfa\xb2S\xaad\xf3#",
		"Sl\u007f\x00\x00\x00",
		"SlSl\xaad\xf3#\xaad\xf3#",
		"Slintterer",
		"[#L00000000",
		"[[{I\xda0",
		"[{I\xda0",
		"[[[{I\xda0",
		"{I\x00\x00{I\x800",
		"Slen\x03\xe8r",
		"SI\x800",
		"{I\xe90",
		"{I\xfa0",
		"Sl\xff000",
	} {
		data := data
		t.Run(data, func(t *testing.T) {
			var i interface{}
			_ = Unmarshal([]byte(data), &i)
		})
	}
}

func TestFuzzUnmarshalBlockCrashers(t *testing.T) {
	for _, data := range []string{
		"[[][$][F][#][I][-7][I]4]",
		"[[][$][T][#][l][1020846876]",
		"[]",
		"[[][[][[][[][H][]",
		"[[][C][]",
		"[[][$][]",
		"[[][S][]",
		"[[][d][7][d][3][1d][7][d][3][d][7]",
		"[[][[][[][S][]",
		"[C][]",
		"[[][[][S][]",
		"[S][]",
	} {
		data := data
		t.Run(data, func(t *testing.T) {
			var i interface{}
			_ = UnmarshalBlock([]byte(data), &i)
		})
	}
}

func TestDecoder_maxCollectionAlloc(t *testing.T) {
	d := NewBlockDecoder(strings.NewReader(
		"[[][$][U][#][U][2][76][127]",
	))
	d.MaxCollectionAlloc = 1
	_, err := d.decodeInterface()
	if err == nil {
		t.Error("expected error")
	}
}

func TestFuzzUnmarshalBlock_strong_type_NoOp_array(t *testing.T) {
	var i interface{}
	if UnmarshalBlock([]byte("[[][$][N][#][I][512]"), &i) == nil {
		t.Errorf("expected failure but got: %v", i)
	}
}

func TestFuzzUnmarshalBlock_strong_type_NoOp_object(t *testing.T) {
	var i interface{}
	if UnmarshalBlock([]byte("[{][$][N][#][i][1][i][4][name]"), &i) == nil {
		t.Errorf("expected failure but got: %v", i)
	}
}
