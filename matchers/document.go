package matchers

import (
	"bytes"

	"gopkg.in/h2non/filetype.v1/util"
)

const defaultDocumentPriority = 500
const officeContentTypesFileName = "[Content_Types].xml"

var (
	TypeDoc  = newType("doc", "application/msword", defaultDocumentPriority)
	TypeDocx = newType("docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", defaultDocumentPriority)
	TypeXls  = newType("xls", "application/vnd.ms-excel", defaultDocumentPriority)
	TypeXlsx = newType("xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", defaultDocumentPriority)
	TypePpt  = newType("ppt", "application/vnd.ms-powerpoint", defaultDocumentPriority)
	TypePptx = newType("pptx", "application/vnd.openxmlformats-officedocument.presentationml.presentation", defaultDocumentPriority)
)

var Document = Map{
	TypeDoc:  Doc,
	TypeDocx: Docx,
	TypeXls:  Xls,
	TypeXlsx: Xlsx,
	TypePpt:  Ppt,
	TypePptx: Pptx,
}

func Doc(buf []byte) bool {
	return len(buf) > 7 &&
		buf[0] == 0xD0 && buf[1] == 0xCF &&
		buf[2] == 0x11 && buf[3] == 0xE0 &&
		buf[4] == 0xA1 && buf[5] == 0xB1 &&
		buf[6] == 0x1A && buf[7] == 0xE1
}

func Docx(buf []byte) bool {
	if len(buf) <= 3 {
		return false
	}
	if buf[0] != 0x50 || buf[1] != 0x4B ||
		buf[2] != 0x03 || buf[3] != 0x04 {
		return false
	}

	xml, err := util.ReadMSOfficeXMLFile(buf, officeContentTypesFileName)
	if err != nil {
		return false
	}

	return bytes.Contains(xml, []byte(TypeDocx.MIME.Value))
}

func Xls(buf []byte) bool {
	return len(buf) > 7 &&
		buf[0] == 0xD0 && buf[1] == 0xCF &&
		buf[2] == 0x11 && buf[3] == 0xE0 &&
		buf[4] == 0xA1 && buf[5] == 0xB1 &&
		buf[6] == 0x1A && buf[7] == 0xE1
}

func Xlsx(buf []byte) bool {
	if len(buf) <= 3 {
		return false
	}
	if buf[0] != 0x50 || buf[1] != 0x4B ||
		buf[2] != 0x03 || buf[3] != 0x04 {
		return false
	}

	xml, err := util.ReadMSOfficeXMLFile(buf, officeContentTypesFileName)
	if err != nil {
		return false
	}

	return bytes.Contains(xml, []byte(TypeXlsx.MIME.Value))
}

func Ppt(buf []byte) bool {
	return len(buf) > 7 &&
		buf[0] == 0xD0 && buf[1] == 0xCF &&
		buf[2] == 0x11 && buf[3] == 0xE0 &&
		buf[4] == 0xA1 && buf[5] == 0xB1 &&
		buf[6] == 0x1A && buf[7] == 0xE1
}

func Pptx(buf []byte) bool {
	if len(buf) <= 3 {
		return false
	}
	if buf[0] != 0x50 || buf[1] != 0x4B ||
		buf[2] != 0x03 || buf[3] != 0x04 {
		return false
	}

	xml, err := util.ReadMSOfficeXMLFile(buf, officeContentTypesFileName)
	if err != nil {
		return false
	}

	return bytes.Contains(xml, []byte(TypePptx.MIME.Value))
}
