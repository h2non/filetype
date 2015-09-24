package filetype

import (
	"testing"
)

func TestMatches(t *testing.T) {
	cases := []struct {
		buf []byte
		ext string
	}{
		{[]byte{0xFF, 0xD8, 0xFF}, "jpg"},
	}

	for _, test := range cases {
		match, err := Match(test.buf)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}

		if match.Extension != test.ext {
			t.Fatalf("Invalid image type: %s", match.Extension)
		}
	}
}
