package matchers

var (
	TypeWoff  = newType("woff", "application/font-woff")
	TypeWoff2 = newType("woff2", "application/font-woff")
	TypeTtf   = newType("ttf", "application/font-sfnt")
	TypeOtf   = newType("otf", "application/font-sfnt")
	TypeDoc   = newType("doc", "application/msword")
	TypeDocx  = newType("docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	TypeXls   = newType("xls", "application/vnd.ms-excel")
	TypeXlsx  = newType("xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	TypePpt   = newType("ppt", "application/vnd.ms-powerpoint")
	TypePptx  = newType("pptx", "application/vnd.openxmlformats-officedocument.presentationml.presentation")
)

var Font = Map{
	TypeWoff:  Woff,
	TypeWoff2: Woff2,
	TypeTtf:   Ttf,
	TypeOtf:   Otf,
	TypeDoc:   Doc,
	TypeDocx:  Docx,
	TypeXls:   Xls,
	TypeXlsx:  Xlsx,
	TypePpt:   Ppt,
	TypePptx:  Pptx,
}

func Woff(buf []byte) bool {
	return len(buf) > 7 &&
		buf[0] == 0x77 && buf[1] == 0x4F &&
		buf[2] == 0x46 && buf[3] == 0x46 &&
		buf[4] == 0x00 && buf[5] == 0x01 &&
		buf[6] == 0x00 && buf[7] == 0x00
}

func Woff2(buf []byte) bool {
	return len(buf) > 7 &&
		buf[0] == 0x77 && buf[1] == 0x4F &&
		buf[2] == 0x46 && buf[3] == 0x32 &&
		buf[4] == 0x00 && buf[5] == 0x01 &&
		buf[6] == 0x00 && buf[7] == 0x00
}

func Ttf(buf []byte) bool {
	return len(buf) > 4 &&
		buf[0] == 0x00 && buf[1] == 0x01 &&
		buf[2] == 0x00 && buf[3] == 0x00 &&
		buf[4] == 0x00
}

func Otf(buf []byte) bool {
	return len(buf) > 4 &&
		buf[0] == 0x4F && buf[1] == 0x54 &&
		buf[2] == 0x54 && buf[3] == 0x4F &&
		buf[4] == 0x00
}

func Doc(buf []byte) bool {
	return len(buf) > 7 &&
		buf[0] == 0xD0 && buf[1] == 0xCF &&
		buf[2] == 0x11 && buf[3] == 0xE0 &&
		buf[4] == 0xA1 && buf[5] == 0xB1 &&
		buf[6] == 0x1A && buf[7] == 0xE1
}

func Docx(buf []byte) bool {
	return len(buf) > 3 &&
		buf[0] == 0x50 && buf[1] == 0x4B &&
		buf[2] == 0x03 && buf[3] == 0x04
}

func Xls(buf []byte) bool {
	return len(buf) > 7 &&
		buf[0] == 0xD0 && buf[1] == 0xCF &&
		buf[2] == 0x11 && buf[3] == 0xE0 &&
		buf[4] == 0xA1 && buf[5] == 0xB1 &&
		buf[6] == 0x1A && buf[7] == 0xE1
}

func Xlsx(buf []byte) bool {
	return len(buf) > 3 &&
		buf[0] == 0x50 && buf[1] == 0x4B &&
		buf[2] == 0x03 && buf[3] == 0x04
}

func Ppt(buf []byte) bool {
	return len(buf) > 7 &&
		buf[0] == 0xD0 && buf[1] == 0xCF &&
		buf[2] == 0x11 && buf[3] == 0xE0 &&
		buf[4] == 0xA1 && buf[5] == 0xB1 &&
		buf[6] == 0x1A && buf[7] == 0xE1
}

func Pptx(buf []byte) bool {
	return len(buf) > 3 &&
		buf[0] == 0x50 && buf[1] == 0x4B &&
		buf[2] == 0x07 && buf[3] == 0x08
}
