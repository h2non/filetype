package filetype

import (
	"bytes"
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/h2non/filetype/types"
)

func TestIs(t *testing.T) {
	cases := []struct {
		buf   []byte
		ext   string
		match bool
	}{
		{[]byte{0xFF, 0xD8, 0xFF}, "jpg", true},
		{[]byte{0xFF, 0xD8, 0x00}, "jpg", false},
		{[]byte{0x89, 0x50, 0x4E, 0x47}, "png", true},
	}

	for _, test := range cases {
		if Is(test.buf, test.ext) != test.match {
			t.Fatalf("Invalid match: %s", test.ext)
		}
	}
}

func TestIsType(t *testing.T) {
	cases := []struct {
		buf   []byte
		kind  types.Type
		match bool
	}{
		{[]byte{0xFF, 0xD8, 0xFF}, types.Get("jpg"), true},
		{[]byte{0xFF, 0xD8, 0x00}, types.Get("jpg"), false},
		{[]byte{0x89, 0x50, 0x4E, 0x47}, types.Get("png"), true},
	}

	for _, test := range cases {
		if IsType(test.buf, test.kind) != test.match {
			t.Fatalf("Invalid match: %s", test.kind.Extension)
		}
	}
}

func TestIsMIME(t *testing.T) {
	cases := []struct {
		buf   []byte
		mime  string
		match bool
	}{
		{[]byte{0xFF, 0xD8, 0xFF}, "image/jpeg", true},
		{[]byte{0xFF, 0xD8, 0x00}, "image/jpeg", false},
		{[]byte{0x89, 0x50, 0x4E, 0x47}, "image/png", true},
	}

	for _, test := range cases {
		if IsMIME(test.buf, test.mime) != test.match {
			t.Fatalf("Invalid match: %s", test.mime)
		}
	}
}

func TestIsSupported(t *testing.T) {
	cases := []struct {
		ext   string
		match bool
	}{
		{"jpg", true},
		{"jpeg", false},
		{"abc", false},
		{"png", true},
		{"mp4", true},
		{"", false},
	}

	for _, test := range cases {
		if IsSupported(test.ext) != test.match {
			t.Fatalf("Invalid match: %s", test.ext)
		}
	}
}

func TestIsMIMESupported(t *testing.T) {
	cases := []struct {
		mime  string
		match bool
	}{
		{"image/jpeg", true},
		{"foo/bar", false},
		{"image/png", true},
		{"video/mpeg", true},
	}

	for _, test := range cases {
		if IsMIMESupported(test.mime) != test.match {
			t.Fatalf("Invalid match: %s", test.mime)
		}
	}
}

func TestAddType(t *testing.T) {
	AddType("foo", "foo/foo")

	if !IsSupported("foo") {
		t.Fatalf("Not supported extension")
	}

	if !IsMIMESupported("foo/foo") {
		t.Fatalf("Not supported MIME type")
	}
}

func TestGetType(t *testing.T) {
	jpg := GetType("jpg")
	if jpg == types.Unknown {
		t.Fatalf("Type should be supported")
	}

	invalid := GetType("invalid")
	if invalid != Unknown {
		t.Fatalf("Type should not be supported")
	}
}

func TestMatchWriter(t *testing.T) {
	cases := []struct {
		mime   string
		reader io.Reader
		read   int64
		err    error
	}{
		{"image/jpeg", bytes.NewReader([]byte{0xFF, 0xD8, 0xFF}), 3, nil},
		{"image/png", bytes.NewReader([]byte{0x89, 0x50, 0x4E, 0x47}), 4, nil},
		{types.Unknown.MIME.Value, bytes.NewReader([]byte{}), 0, ErrEmptyBuffer},
		{types.Unknown.MIME.Value, nil, 0, ErrEmptyBuffer},
	}
	for _, test := range cases {
		var w int64
		var err error
		var mimeType types.Type
		mw := NewMatcherWriter()
		if test.reader != nil {
			w, err = io.Copy(mw, test.reader)
		}
		if err != nil {
			t.Fatalf("Error matching %s error: %v", test.mime, err)
		}
		mimeType, err = mw.Match()
		if err != test.err {
			t.Fatalf("Invalid error match: %v, expected %s", err, test.err)
		}
		if mimeType.MIME.Value != test.mime {
			t.Fatalf("Invalid mime match: %s, expected %s", mimeType.MIME.Value, test.mime)
		}
		if w != test.read {
			t.Fatalf("Invalid read match: %d, expected %d", w, test.read)
		}
	}
}

func generateRandomSlice(size int) []byte {
	slice := make([]byte, size, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = byte(rand.Intn(999))
	}
	return slice
}

func TestMatchWriterBuffer(t *testing.T) {
	randomSlice := generateRandomSlice(maxBufSize + 500)
	reader := bytes.NewReader(randomSlice)
	mw := NewMatcherWriter()
	_, err := io.Copy(mw, reader)
	if err != nil {
		t.Fatalf("error copying bytes to reader %v", err)
	}
	bufLen := len(mw.buf)
	if bufLen != maxBufSize {
		t.Fatalf("expected buffer len to be %d but is %d", maxBufSize, bufLen)
	}
	if bytes.Compare(mw.buf, randomSlice[0:maxBufSize]) != 0 {
		t.Fatalf("expected buffer to equal re sliced buffer")
	}
}
