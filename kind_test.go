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
		{[]byte{0x89, 0x50, 0x4E, 0x47}, "png"},
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
		{[]byte{0x89, 0x50, 0x4E, 0x47}, true},
		{[]byte{0x89, 0x0, 0x0}, false},
	}

	for _, test := range cases {
		if IsImage(test.buf) != test.match {
			t.Fatalf("Invalid match: %t", test.match)
		}
	}
}

func TestKindClasses(t *testing.T) {
	midi := []byte{0x4D, 0x54, 0x68, 0x64}
	flv := []byte{0x46, 0x4C, 0x56, 0x01}
	woff := []byte{0x77, 0x4F, 0x46, 0x46, 0x00, 0x01, 0x00, 0x00}
	zip := []byte{0x50, 0x4B, 0x03, 0x04}
	doc := make([]byte, 514)
	copy(doc, []byte{0xD0, 0xCF, 0x11, 0xE0})
	doc[512] = 0xEC
	doc[513] = 0xA5
	wasm := []byte{0x00, 0x61, 0x73, 0x6D, 0x01, 0x00, 0x00, 0x00}
	invalid := []byte{0x00, 0x00, 0x00, 0x00}

	// Audio
	if kind, _ := Audio(midi); kind.Extension != "mid" || !IsAudio(midi) || IsAudio(invalid) {
		t.Fatalf("Audio check failed")
	}

	// Video
	if kind, _ := Video(flv); kind.Extension != "flv" || !IsVideo(flv) || IsVideo(invalid) {
		t.Fatalf("Video check failed")
	}

	// Font
	if kind, _ := Font(woff); kind.Extension != "woff" || !IsFont(woff) || IsFont(invalid) {
		t.Fatalf("Font check failed")
	}

	// Archive
	if kind, _ := Archive(zip); kind.Extension != "zip" || !IsArchive(zip) || IsArchive(invalid) {
		t.Fatalf("Archive check failed")
	}

	// Document
	if kind, _ := Document(doc); kind.Extension != "doc" || !IsDocument(doc) || IsDocument(invalid) {
		t.Fatalf("Document check failed")
	}

	// Application
	if kind, _ := Application(wasm); kind.Extension != "wasm" || !IsApplication(wasm) || IsApplication(invalid) {
		t.Fatalf("Application check failed")
	}
}
