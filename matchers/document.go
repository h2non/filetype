package matchers

import (
	"encoding/binary"
	"strings"
)

var (
	TypeDoc  = newType("doc", "application/msword")
	TypeDocx = newType("docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	TypeXls  = newType("xls", "application/vnd.ms-excel")
	TypeXlsx = newType("xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	TypePpt  = newType("ppt", "application/vnd.ms-powerpoint")
	TypePptx = newType("pptx", "application/vnd.openxmlformats-officedocument.presentationml.presentation")
	TypeOdp  = newType("odp", "application/vnd.oasis.opendocument.presentation")
	TypeOds  = newType("ods", "application/vnd.oasis.opendocument.spreadsheet")
	TypeOdt  = newType("odt", "application/vnd.oasis.opendocument.text")
)

var Document = Map{
	TypeDoc:  Doc,
	TypeDocx: Docx,
	TypeXls:  Xls,
	TypeXlsx: Xlsx,
	TypePpt:  Ppt,
	TypePptx: Pptx,
	TypeOdp:  Odp,
	TypeOds:  Ods,
	TypeOdt:  Odt,
}

type docType int

const (
	TYPE_DOC docType = iota
	TYPE_DOCX
	TYPE_XLS
	TYPE_XLSX
	TYPE_PPT
	TYPE_PPTX
	TYPE_OOXML
	TYPE_ODP
	TYPE_ODS
	TYPE_ODT
)

// reference: https://bz.apache.org/ooo/show_bug.cgi?id=111457
func Doc(buf []byte) bool {
	if len(buf) > 513 {
		return buf[0] == 0xD0 && buf[1] == 0xCF &&
			buf[2] == 0x11 && buf[3] == 0xE0 &&
			buf[512] == 0xEC && buf[513] == 0xA5
	} else {
		return len(buf) > 3 &&
			buf[0] == 0xD0 && buf[1] == 0xCF &&
			buf[2] == 0x11 && buf[3] == 0xE0
	}
}

func Docx(buf []byte) bool {
	typ, ok := msooxml(buf)
	return ok && typ == TYPE_DOCX
}

func Xls(buf []byte) bool {
	if len(buf) > 513 {
		isMSOfficeBFF := buf[0] == 0xD0 && buf[1] == 0xCF && buf[2] == 0x11 && buf[3] == 0xE0

		switch {
		case isMSOfficeBFF && buf[512] == 0x09 && buf[513] == 0x08: // BIFF5 && BIFF12(12)
			return true
		case isMSOfficeBFF && buf[512] == 0xFD && buf[513] == 0xFF: // BIFF12(11)
			return true
		}
	}

	return false
}

func Xlsx(buf []byte) bool {
	typ, ok := msooxml(buf)
	return ok && typ == TYPE_XLSX
}

func Ppt(buf []byte) bool {
	if len(buf) > 513 {
		return buf[0] == 0xD0 && buf[1] == 0xCF &&
			buf[2] == 0x11 && buf[3] == 0xE0 &&
			buf[512] == 0xA0 && buf[513] == 0x46
	} else {
		return len(buf) > 3 &&
			buf[0] == 0xD0 && buf[1] == 0xCF &&
			buf[2] == 0x11 && buf[3] == 0xE0
	}
}

func Pptx(buf []byte) bool {
	typ, ok := msooxml(buf)
	return ok && typ == TYPE_PPTX
}

func msooxml(buf []byte) (typ docType, found bool) {
	signature := []byte{'P', 'K', 0x03, 0x04}

	// start by checking for ZIP local file header signature
	if ok := compareBytes(buf, signature, 0); !ok {
		return
	}

	// Walk every local file header in the buffer. XLSX/DOCX/PPTX writers
	// (Excel, excelize, Google Sheets, LibreOffice) do not share a ZIP
	// entry order, so looking only at the 1st/3rd/4th name misses files.
	sawOOXML := false
	forEachZipLocalName(buf, func(name string) bool {
		switch {
		case strings.HasPrefix(name, "word/"):
			typ, found = TYPE_DOCX, true
			return true
		case strings.HasPrefix(name, "xl/"):
			typ, found = TYPE_XLSX, true
			return true
		case strings.HasPrefix(name, "ppt/"):
			typ, found = TYPE_PPTX, true
			return true
		case name == "[Content_Types].xml" ||
			strings.HasPrefix(name, "_rels") ||
			strings.HasPrefix(name, "docProps"):
			sawOOXML = true
		}
		return false
	})
	if found {
		return
	}
	if sawOOXML {
		return TYPE_OOXML, true
	}
	return
}

// forEachZipLocalName calls fn with each ZIP local-file name found in buf.
// fn returning true stops the walk.
func forEachZipLocalName(buf []byte, fn func(name string) bool) {
	signature := []byte{'P', 'K', 0x03, 0x04}
	offset := 0
	for offset+30 <= len(buf) {
		if !compareBytes(buf, signature, offset) {
			offset++
			continue
		}
		nameLen := int(binary.LittleEndian.Uint16(buf[offset+26 : offset+28]))
		extraLen := int(binary.LittleEndian.Uint16(buf[offset+28 : offset+30]))
		nameOff := offset + 30
		if nameLen <= 0 || nameLen > 512 || nameOff+nameLen > len(buf) {
			offset += 4
			continue
		}
		name := string(buf[nameOff : nameOff+nameLen])
		if fn(name) {
			return
		}
		next := nameOff + nameLen + extraLen
		if next <= offset {
			return
		}
		offset = next
	}
}

func compareBytes(slice, subSlice []byte, startOffset int) bool {
	sl := len(subSlice)

	if startOffset+sl > len(slice) {
		return false
	}

	s := slice[startOffset : startOffset+sl]
	for i := range s {
		if subSlice[i] != s[i] {
			return false
		}
	}

	return true
}

func Odp(buf []byte) bool {
	return checkOdf(buf, TypeOdp.MIME.Value)
}

func Ods(buf []byte) bool {
	return checkOdf(buf, TypeOds.MIME.Value)
}

func Odt(buf []byte) bool {
	return checkOdf(buf, TypeOdt.MIME.Value)
}

// https://en.wikipedia.org/wiki/OpenDocument_technical_specification
// https://en.wikipedia.org/wiki/ZIP_(file_format)
func checkOdf(buf []byte, mimetype string) bool {
	if 38+len(mimetype) >= len(buf) {
		return false
	}
	// Perform all byte checks first for better performance
	// Check ZIP start
	if buf[0] != 'P' || buf[1] != 'K' || buf[2] != 3 || buf[3] != 4 {
		return false
	}
	// Now check the first file data
	// Compression method: not compressed
	if buf[8] != 0 || buf[9] != 0 {
		return false
	}
	// Filename length must be 8 for "mimetype"
	if buf[26] != 8 || buf[27] != 0 {
		return false
	}
	// Check the file contents sizes
	if int(buf[18]) != len(mimetype) ||
		buf[19] != 0 || buf[20] != 0 || buf[21] != 0 ||
		int(buf[22]) != len(mimetype) ||
		buf[23] != 0 || buf[24] != 0 || buf[25] != 0 {
		return false
	}
	// No extra field (for data offset below)
	if buf[28] != 0 || buf[29] != 0 {
		return false
	}
	// Finally check the file name and contents
	return string(buf[30:38]) == "mimetype" &&
		string(buf[38:38+len(mimetype)]) == mimetype
}
