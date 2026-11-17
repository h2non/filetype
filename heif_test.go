package filetype

import (
	"encoding/binary"
	"encoding/hex"
	"testing"

	"github.com/h2non/filetype/matchers"
)

func TestMatchHEIFBrands(t *testing.T) {
	ftyp := func(major string, compatible ...string) []byte {
		b := make([]byte, 16+4*len(compatible))
		binary.BigEndian.PutUint32(b, uint32(len(b)))
		copy(b[4:8], "ftyp")
		copy(b[8:12], major)
		for i, brand := range compatible {
			copy(b[16+4*i:], brand)
		}
		return b
	}
	// File type boxes from two Sony HIF samples; no image content or metadata.
	sony, err := hex.DecodeString("000000286674797068656978000000006d696631686569786d6961664d6948416a70656753484946")
	if err != nil {
		t.Fatal(err)
	}
	oriented, err := hex.DecodeString("0000001c6674797068656978000000006d696631686569786d696166")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		buf  []byte
	}{
		{"sony", sony},
		{"sony orientation", oriented},
		{"heic major", ftyp("heic", "mif1", "heic")},
		{"mif1 heic", ftyp("mif1", "mif1", "heic")},
		{"msf1 heic", ftyp("msf1", "msf1", "heic")},
		{"mif1 heix", ftyp("mif1", "mif1", "miaf", "heix")},
		{"msf1 heix", ftyp("msf1", "msf1", "heix")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			kind, err := Match(tc.buf)
			if err != nil || kind.Extension != "heif" || kind.MIME.Value != "image/heif" {
				t.Fatalf("Match() = %v, %v; want heif / image/heif", kind, err)
			}
		})
	}
	for _, tc := range []struct {
		name string
		buf  []byte
	}{
		{"mp4", ftyp("isom", "isom", "mp41")},
		{"avif", ftyp("avif", "avif", "mif1", "miaf")},
		{"generic mif1", ftyp("mif1", "mif1", "miaf")},
		{"unknown major", ftyp("test", "heix")},
		{"heix in minor version", append(append([]byte(nil), ftyp("mif1", "mif1")[:12]...), []byte("heixmif1")...)},
		{"heix outside ftyp", append(ftyp("mif1", "mif1"), []byte("heix")...)},
		{"not ftyp", []byte("\x00\x00\x00\x1cmdatheix\x00\x00\x00\x00mif1heixmiaf")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if matchers.Heif(tc.buf) {
				t.Fatal("non-HEIF header matched HEIF")
			}
			kind, err := Match(tc.buf)
			if err != nil || kind == matchers.TypeHeif {
				t.Fatalf("Match() = %v, %v", kind, err)
			}
		})
	}
	for n := 0; n < len(oriented); n++ {
		if matchers.Heif(oriented[:n]) {
			t.Fatalf("truncated header matched at length %d", n)
		}
	}
}
