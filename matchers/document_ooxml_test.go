package matchers

import (
	"archive/zip"
	"bytes"
	"testing"
)

func zipNamed(names ...string) []byte {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	for _, name := range names {
		w, err := zw.Create(name)
		if err != nil {
			panic(err)
		}
		if _, err := w.Write([]byte("x")); err != nil {
			panic(err)
		}
	}
	if err := zw.Close(); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func TestXlsxDetectsLateXlEntry(t *testing.T) {
	// excelize/Sheets often put xl/ after several other ZIP entries.
	buf := zipNamed(
		"docProps/app.xml",
		"docProps/core.xml",
		"[Content_Types].xml",
		"_rels/.rels",
		"xl/workbook.xml",
		"xl/_rels/workbook.xml.rels",
	)
	if !Xlsx(buf) {
		t.Fatal("expected xlsx when xl/ is not the 3rd ZIP entry")
	}
	if Docx(buf) || Pptx(buf) {
		t.Fatal("xlsx must not also match docx/pptx")
	}
}

func TestDocxDetectsLateWordEntry(t *testing.T) {
	buf := zipNamed(
		"[Content_Types].xml",
		"_rels/.rels",
		"docProps/core.xml",
		"customXml/item1.xml",
		"word/document.xml",
	)
	if !Docx(buf) {
		t.Fatal("expected docx when word/ is not an early ZIP entry")
	}
}

func TestPptxDetectsLatePptEntry(t *testing.T) {
	buf := zipNamed(
		"[Content_Types].xml",
		"_rels/.rels",
		"docProps/core.xml",
		"ppt/presentation.xml",
	)
	if !Pptx(buf) {
		t.Fatal("expected pptx when ppt/ is not an early ZIP entry")
	}
}

func TestMsooxmlIgnoresPlainZip(t *testing.T) {
	buf := zipNamed("readme.txt", "data/file.bin")
	if Xlsx(buf) || Docx(buf) || Pptx(buf) {
		t.Fatal("plain zip must not match ooxml document types")
	}
}
