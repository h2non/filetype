package filetype

import (
	"testing"
)

func TestKind(t *testing.T) {
	var cases = []struct {
		buf []byte
		ext string
	}{
		{[]byte{0xFF, 0xD8, 0xFF}, "jpg"},
		{[]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, "png"},
		{[]byte{0x89, 0x0, 0x0}, "unknown"},
	}

	for _, test := range cases {
		kind, _ := Image(test.buf)
		if kind.Extension != test.ext {
			t.Fatalf("Invalid match: %s != %s", kind.Extension, test.ext)
		}
	}
}

func TestIsKind(t *testing.T) {
	var cases = []struct {
		buf   []byte
		match bool
	}{
		{[]byte{0xFF, 0xD8, 0xFF}, true},
		{[]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, true},
		{[]byte{0x89, 0x0, 0x0}, false},
	}

	for _, test := range cases {
		if IsImage(test.buf) != test.match {
			t.Fatalf("Invalid match: %t", test.match)
		}
	}
}
