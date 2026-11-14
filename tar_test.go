package filetype

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/h2non/filetype/matchers"
)

func tarArchive(t testing.TB, name string, legacy, signedChecksum bool) []byte {
	t.Helper()
	var archive bytes.Buffer
	w := tar.NewWriter(&archive)
	if err := w.WriteHeader(&tar.Header{Name: "sample.txt", Mode: 0644, Size: 3, Format: tar.FormatUSTAR}); err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("abc")); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	data := archive.Bytes()
	if legacy {
		// V7 headers end after linkname; the remaining bytes are padding.
		for i := 257; i < 512; i++ {
			data[i] = 0
		}
		for i := 0; i < 100; i++ {
			data[i] = 0
		}
		copy(data[:100], name)
		data[156] = 0
		for i := 148; i < 156; i++ {
			data[i] = ' '
		}
		var checksum int
		for _, b := range data[:512] {
			if signedChecksum {
				checksum += int(int8(b))
			} else {
				checksum += int(b)
			}
		}
		format := "%06o\x00 "
		if signedChecksum {
			format = "%6o\x00 "
		}
		copy(data[148:156], fmt.Sprintf(format, checksum))
	}

	// Check the generated fixture with an independent TAR reader.
	r := tar.NewReader(bytes.NewReader(data))
	h, err := r.Next()
	if err != nil {
		t.Fatalf("archive/tar rejected fixture: %v", err)
	}
	if h.Name != name || h.Size != 3 {
		t.Fatalf("unexpected TAR header: %+v", h)
	}
	payload, err := ioutil.ReadAll(r)
	if err != nil || string(payload) != "abc" {
		t.Fatalf("unexpected TAR payload: %q, %v", payload, err)
	}
	return data
}

func TestMatchTarFormats(t *testing.T) {
	for _, test := range []struct {
		name, filename string
		legacy, signed bool
	}{
		{"ustar", "sample.txt", false, false},
		{"v7", "sample.txt", true, false},
		{"v7-unsigned-unicode", "café.txt", true, false},
		{"v7-signed-unicode", "café.txt", true, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			data := tarArchive(t, test.filename, test.legacy, test.signed)
			t.Run("MatchHeader", func(t *testing.T) {
				kind, err := Match(data[:512])
				if err != nil || kind.Extension != "tar" || kind.MIME.Value != "application/x-tar" {
					t.Fatalf("got %v, %v; want tar application/x-tar", kind, err)
				}
			})
			t.Run("Match", func(t *testing.T) {
				kind, err := Match(data)
				if err != nil || kind.Extension != "tar" || kind.MIME.Value != "application/x-tar" {
					t.Fatalf("got %v, %v; want tar application/x-tar", kind, err)
				}
			})
			t.Run("MatchReader", func(t *testing.T) {
				kind, err := MatchReader(bytes.NewReader(data))
				if err != nil || kind.Extension != "tar" || kind.MIME.Value != "application/x-tar" {
					t.Fatalf("got %v, %v; want tar application/x-tar", kind, err)
				}
			})
		})
	}

	// The USTAR signature remains sufficient without a complete header.
	ustar := tarArchive(t, "sample.txt", false, false)
	if !matchers.Tar(ustar[:262]) {
		t.Fatal("USTAR prefix was not recognized")
	}
}

func TestTarRejectsInvalidHeaders(t *testing.T) {
	valid := tarArchive(t, "sample.txt", true, false)
	withChecksum := func(value string) []byte {
		data := append([]byte(nil), valid...)
		copy(data[148:156], value)
		return data
	}
	corrupt := append([]byte(nil), valid...)
	corrupt[0]++
	random := make([]byte, 512)
	for i := range random {
		random[i] = byte(i*37 + 11)
	}
	for _, test := range []struct {
		name string
		data []byte
	}{
		{"empty", nil},
		{"truncated", valid[:511]},
		{"zero", make([]byte, 512)},
		{"random", random},
		{"corrupted", corrupt},
		{"non-octal", withChecksum("000008\x00 ")},
		{"signed-number", withChecksum("+00001\x00 ")},
		{"embedded-nul", withChecksum("00\x00001\x00 ")},
		{"embedded-space", withChecksum("00 001\x00 ")},
		{"empty-checksum", withChecksum("\x00\x00\x00\x00\x00\x00\x00\x00")},
		{"blank-checksum", withChecksum("        ")},
	} {
		t.Run(test.name, func(t *testing.T) {
			if matchers.Tar(test.data) {
				t.Fatal("invalid header matched as TAR")
			}
			kind, _ := Match(test.data)
			if kind.Extension == "tar" {
				t.Fatal("Match returned TAR for an invalid header")
			}
			if len(test.data) >= 512 {
				kind, err := MatchReader(bytes.NewReader(test.data))
				if err != nil || kind.Extension == "tar" {
					t.Fatalf("MatchReader returned %v, %v for an invalid header", kind, err)
				}
			}
		})
	}
}

var tarMatchResult bool
var tarMatcher = matchers.Tar

func BenchmarkTarHeader(b *testing.B) {
	ustar := tarArchive(b, "sample.txt", false, false)
	legacy := tarArchive(b, "sample.txt", true, false)
	for _, test := range []struct {
		name string
		data []byte
	}{
		{"ustar", ustar[:262]},
		{"v7", legacy[:512]},
		{"non-tar", bytes.Repeat([]byte{'x'}, 512)},
	} {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				tarMatchResult = tarMatcher(test.data)
			}
		})
	}
}
