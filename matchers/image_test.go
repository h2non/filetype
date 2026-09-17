package matchers

import (
	"encoding/binary"
	"testing"
)

func ftypBox(major string, brands ...string) []byte {
	n := 16 + 4*len(brands)
	buf := make([]byte, n)
	binary.BigEndian.PutUint32(buf[0:4], uint32(n))
	copy(buf[4:8], "ftyp")
	copy(buf[8:12], major)
	for i, b := range brands {
		copy(buf[16+4*i:], b)
	}
	return buf
}

func TestAvifBrands(t *testing.T) {
	cases := []struct {
		name string
		buf  []byte
		want bool
	}{
		{"major avif", ftypBox("avif", "avif", "mif1", "miaf"), true},
		{"mif1 + avif", ftypBox("mif1", "mif1", "avif"), true},
		{"msf1 + avif", ftypBox("msf1", "msf1", "avif"), true},
		{"miaf + avif", ftypBox("miaf", "miaf", "avif", "mif1"), true},
		{"major avis", ftypBox("avis", "avis", "avif"), true},
		{"miaf + avis", ftypBox("miaf", "miaf", "avis"), true},
		{"miaf without avif", ftypBox("miaf", "miaf", "mif1"), false},
		{"heic", ftypBox("heic", "mif1", "heic"), false},
		{"mif1 + heic", ftypBox("mif1", "mif1", "heic"), false},
		{"not isobmff", []byte("not a media box"), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Avif(tc.buf); got != tc.want {
				t.Fatalf("Avif() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHeifNotAvif(t *testing.T) {
	buf := ftypBox("miaf", "miaf", "avif", "mif1")
	if !Avif(buf) {
		t.Fatal("miaf+avif should match Avif")
	}
	if Heif(buf) {
		t.Fatal("miaf+avif should not match Heif")
	}
}
