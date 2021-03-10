package filetype

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
)

func TestMatch(t *testing.T) {
	cases := []struct {
		buf []byte
		ext string
	}{
		{[]byte{0xFF, 0xD8, 0xFF}, "jpg"},
		{[]byte{0xFF, 0xD8, 0x00}, "unknown"},
		{[]byte{0x89, 0x50, 0x4E, 0x47}, "png"},
	}

	for _, test := range cases {
		match, err := Match(test.buf)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}

		if match.Extension != test.ext {
			t.Fatalf("Invalid image type: %s != %s", match.Extension, test.ext)
		}
	}
}

func TestMatchFile(t *testing.T) {
	cases := []struct {
		ext string
	}{
		{"gif"},
		{"jpg"},
		{"png"},
		{"zip"},
		{"tar"},
		{"tif"},
		{"mp4"},
		{"mkv"},
		{"webm"},
		{"docx"},
		{"pptx"},
		{"xlsx"},
		{"mov"},
		{"wasm"},
		{"dwg"},
		{"zst"},
	}

	for _, test := range cases {
		kind, _ := MatchFile("./fixtures/sample." + test.ext)
		if kind.Extension != test.ext {
			t.Fatalf("Invalid type: %s != %s", kind.Extension, test.ext)
		}
	}
	// test zstd with skippable frame
	kind, _ := MatchFile("./fixtures/sample_skippable.zst")
	if kind.Extension != "zst" {
		t.Fatalf("Invalid type: %s != %s", kind.Extension, "zst")
	}
}

func TestMatchReader(t *testing.T) {
	cases := []struct {
		buf io.Reader
		ext string
	}{
		{bytes.NewBuffer([]byte{0xFF, 0xD8, 0xFF}), "jpg"},
		{bytes.NewBuffer([]byte{0xFF, 0xD8, 0x00}), "unknown"},
		{bytes.NewBuffer([]byte{0x89, 0x50, 0x4E, 0x47}), "png"},
	}

	for _, test := range cases {
		match, err := MatchReader(test.buf)
		if err != nil {
			t.Fatalf("Error: %s", err)
		}

		if match.Extension != test.ext {
			t.Fatalf("Invalid image type: %s", match.Extension)
		}
	}
}

func TestMatches(t *testing.T) {
	cases := []struct {
		buf   []byte
		match bool
	}{
		{[]byte{0xFF, 0xD8, 0xFF}, true},
		{[]byte{0xFF, 0x0, 0x0}, false},
		{[]byte{0x89, 0x50, 0x4E, 0x47}, true},
	}

	for _, test := range cases {
		if Matches(test.buf) != test.match {
			t.Fatalf("Do not matches: %#v", test.buf)
		}
	}
}

func TestAddMatcher(t *testing.T) {
	fileType := AddType("foo", "foo/foo")

	AddMatcher(fileType, func(buf []byte) bool {
		return len(buf) == 2 && buf[0] == 0x00 && buf[1] == 0x00
	})

	if !Is([]byte{0x00, 0x00}, "foo") {
		t.Fatalf("Type cannot match")
	}

	if !IsSupported("foo") {
		t.Fatalf("Not supported extension")
	}

	if !IsMIMESupported("foo/foo") {
		t.Fatalf("Not supported MIME type")
	}
}

func TestMatchMap(t *testing.T) {
	cases := []struct {
		buf  []byte
		kind types.Type
	}{
		{[]byte{0xFF, 0xD8, 0xFF}, types.Get("jpg")},
		{[]byte{0x89, 0x50, 0x4E, 0x47}, types.Get("png")},
		{[]byte{0xFF, 0x0, 0x0}, Unknown},
	}

	for _, test := range cases {
		if kind := MatchMap(test.buf, matchers.Image); kind != test.kind {
			t.Fatalf("Do not matches: %#v", test.buf)
		}
	}
}

func TestMatchesMap(t *testing.T) {
	cases := []struct {
		buf   []byte
		match bool
	}{
		{[]byte{0xFF, 0xD8, 0xFF}, true},
		{[]byte{0x89, 0x50, 0x4E, 0x47}, true},
		{[]byte{0xFF, 0x0, 0x0}, false},
	}

	for _, test := range cases {
		if match := MatchesMap(test.buf, matchers.Image); match != test.match {
			t.Fatalf("Do not matches: %#v", test.buf)
		}
	}
}

//
// Benchmarks
//

var tarBuffer, _ = ioutil.ReadFile("./fixtures/sample.tar")
var zipBuffer, _ = ioutil.ReadFile("./fixtures/sample.zip")
var jpgBuffer, _ = ioutil.ReadFile("./fixtures/sample.jpg")
var gifBuffer, _ = ioutil.ReadFile("./fixtures/sample.gif")
var pngBuffer, _ = ioutil.ReadFile("./fixtures/sample.png")
var xlsxBuffer, _ = ioutil.ReadFile("./fixtures/sample.xlsx")
var pptxBuffer, _ = ioutil.ReadFile("./fixtures/sample.pptx")
var docxBuffer, _ = ioutil.ReadFile("./fixtures/sample.docx")
var dwgBuffer, _ = ioutil.ReadFile("./fixtures/sample.dwg")
var mkvBuffer, _ = ioutil.ReadFile("./fixtures/sample.mkv")
var webmBuffer, _ = ioutil.ReadFile("./fixtures/sample.webm")

func BenchmarkMatchTar(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Match(tarBuffer)
	}
}

func BenchmarkMatchZip(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Match(zipBuffer)
	}
}

func BenchmarkMatchJpeg(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Match(jpgBuffer)
	}
}

func BenchmarkMatchGif(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Match(gifBuffer)
	}
}

func BenchmarkMatchPng(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Match(pngBuffer)
	}
}

func BenchmarkMatchXlsx(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Match(xlsxBuffer)
	}
}

func BenchmarkMatchPptx(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Match(pptxBuffer)
	}
}

func BenchmarkMatchDocx(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Match(docxBuffer)
	}
}

func BenchmarkMatchDwg(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Match(dwgBuffer)
	}
}

func BenchmarkMatchMkv(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Match(mkvBuffer)
	}
}

func BenchmarkMatchWebm(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Match(webmBuffer)
	}
}
