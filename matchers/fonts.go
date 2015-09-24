package matchers

var TypeWoff = NewType("woff", "application/font-woff")
var TypeWoff2 = NewType("woff2", "application/font-woff")
var TypeTtf = NewType("ttf", "application/font-sfnt")
var TypeOtf = NewType("otf", "application/font-sfnt")

var Font = Map{
	TypeWoff:  Woff,
	TypeWoff2: Woff2,
	TypeTtf:   Ttf,
	TypeOtf:   Otf,
}

func Woff(buf []byte, length int) bool {
	return length > 7 &&
		buf[0] == 0x77 && buf[1] == 0x4F && buf[2] == 0x46 && buf[3] == 0x46 &&
		buf[4] == 0x00 && buf[5] == 0x01 && buf[6] == 0x00 && buf[7] == 0x00
}

func Woff2(buf []byte, length int) bool {
	return length > 7 &&
		buf[0] == 0x77 && buf[1] == 0x4F && buf[2] == 0x46 && buf[3] == 0x32 &&
		buf[4] == 0x00 && buf[5] == 0x01 && buf[6] == 0x00 && buf[7] == 0x00
}

func Ttf(buf []byte, length int) bool {
	return length > 4 &&
		buf[0] == 0x00 && buf[1] == 0x01 &&
		buf[2] == 0x00 && buf[3] == 0x00 &&
		buf[4] == 0x00
}

func Otf(buf []byte, length int) bool {
	return length > 4 &&
		buf[0] == 0x4F && buf[1] == 0x54 &&
		buf[2] == 0x54 && buf[3] == 0x4F &&
		buf[4] == 0x00
}
