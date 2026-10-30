package filetype

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"testing/iotest"
)

func TestMatchReaderShortReads(t *testing.T) {
	png := []byte{0x89, 0x50, 0x4e, 0x47}
	for _, tt := range []struct {
		name   string
		reader io.Reader
		ext    string
		err    error
	}{
		{"fragmented PNG", iotest.OneByteReader(bytes.NewReader(png)), "png", nil},
		{"fragmented JPEG", iotest.OneByteReader(bytes.NewReader([]byte{0xff, 0xd8, 0xff})), "jpg", nil},
		{"data with EOF", iotest.DataErrReader(bytes.NewReader(png)), "png", nil},
		{"unknown", iotest.OneByteReader(bytes.NewReader([]byte("unknown"))), "unknown", nil},
		{"truncated ICO", bytes.NewReader([]byte{0, 0, 1}), "unknown", nil},
		{"empty", bytes.NewReader(nil), "unknown", ErrEmptyBuffer},
	} {
		t.Run(tt.name, func(t *testing.T) {
			kind, err := MatchReader(tt.reader)
			if kind.Extension != tt.ext || err != tt.err {
				t.Fatalf("MatchReader = (%s, %v), want (%s, %v)", kind.Extension, err, tt.ext, tt.err)
			}
		})
	}
}

func TestMatchReaderReadLimit(t *testing.T) {
	data := make([]byte, 8193)
	copy(data, []byte{0x89, 0x50, 0x4e, 0x47})
	reader := bytes.NewReader(data)
	kind, err := MatchReader(iotest.HalfReader(reader))
	if kind.Extension != "png" || err != nil {
		t.Fatalf("MatchReader = (%s, %v), want (png, nil)", kind.Extension, err)
	}
	if reader.Len() != 1 {
		t.Fatalf("reader has %d bytes left, want 1", reader.Len())
	}
}

type matchReaderError struct {
	err error
}

func (r matchReaderError) Read([]byte) (int, error) { return 0, r.err }

func TestMatchReaderReadError(t *testing.T) {
	want := errors.New("read failed")
	for _, prefix := range [][]byte{nil, {0x89, 0x50, 0x4e, 0x47}} {
		reader := io.MultiReader(bytes.NewReader(prefix), matchReaderError{want})
		kind, err := MatchReader(reader)
		if kind != Unknown || err != want {
			t.Fatalf("MatchReader = (%s, %v), want (unknown, %v)", kind.Extension, err, want)
		}
	}
}
