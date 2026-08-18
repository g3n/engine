package obj

import (
	"strings"
	"testing"
)

// A face that references an out-of-range vertex/uv/normal index must not panic
// when the decoded object is turned into geometry; it must return an error.
func TestNewGeometryOutOfRangeIndex(t *testing.T) {
	cases := map[string]string{
		"vertex too large": "v 0 0 0\nv 1 0 0\nv 0 1 0\nf 1 2 999999\n",
		"vertex negative":  "v 0 0 0\nv 1 0 0\nv 0 1 0\nf -100 -100 -100\n",
		"uv too large":     "v 0 0 0\nv 1 0 0\nv 0 1 0\nvt 0 0\nf 1/1 2/1 3/9999\n",
		"normal too large": "v 0 0 0\nv 1 0 0\nv 0 1 0\nvn 0 0 1\nf 1//1 2//1 3//9999\n",
	}
	for name, src := range cases {
		t.Run(name, func(t *testing.T) {
			dec, err := DecodeReader(strings.NewReader(src), nil)
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if _, err := dec.NewGeometry(&dec.Objects[0]); err == nil {
				t.Fatalf("expected out-of-range error, got nil")
			}
		})
	}
}

// A well-formed OBJ must still build geometry without error.
func TestNewGeometryValid(t *testing.T) {
	src := "v 0 0 0\nv 1 0 0\nv 0 1 0\nvt 0 0\nvt 1 0\nvt 0 1\nvn 0 0 1\nf 1/1/1 2/2/1 3/3/1\n"
	dec, err := DecodeReader(strings.NewReader(src), nil)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, err := dec.NewGeometry(&dec.Objects[0]); err != nil {
		t.Fatalf("valid obj should not error: %v", err)
	}
}
