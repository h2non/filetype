package matchers

import (
	"testing"
)

func TestMp3(t *testing.T) {
	cases := []struct {
		buf      []byte
		expected bool
	}{
		// ID3 container
		{[]byte{0x49, 0x44, 0x33}, true},
		{[]byte{0x49, 0x44, 0x33, 0x03, 0x00}, true},

		// MPEG-1 Layer III without CRC
		{[]byte{0xFF, 0xFB, 0x90}, true},
		// MPEG-1 Layer III with CRC
		{[]byte{0xFF, 0xFA, 0x90}, true},
		// MPEG-2 Layer III without CRC
		{[]byte{0xFF, 0xF3, 0x40}, true},
		// MPEG-2 Layer III with CRC
		{[]byte{0xFF, 0xF2, 0x40}, true},
		// MPEG-2.5 Layer III without CRC
		{[]byte{0xFF, 0xE3, 0x40}, true},
		// MPEG-2.5 Layer III with CRC
		{[]byte{0xFF, 0xE2, 0x40}, true},

		// Invalid cases
		// Too short
		{[]byte{0xFF, 0xFB}, false},
		{[]byte{0x49, 0x44}, false},
		// Invalid version (01 is reserved)
		{[]byte{0xFF, 0xEB, 0x00}, false},
		{[]byte{0xFF, 0xEA, 0x00}, false},
		// MPEG Layer II
		{[]byte{0xFF, 0xFD, 0x00}, false},
		// MPEG Layer I
		{[]byte{0xFF, 0xFF, 0x00}, false},
		// Not MPEG
		{[]byte{0x00, 0x00, 0x00}, false},
	}

	for _, test := range cases {
		if result := Mp3(test.buf); result != test.expected {
			t.Fatalf("Mp3(%#v) = %v, expected %v", test.buf, result, test.expected)
		}
	}
}
