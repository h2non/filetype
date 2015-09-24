package matchers

var TypeJpeg = NewType("jpg", "image/jpeg")
var TypePng = NewType("png", "image/png")
var TypeGif = NewType("gif", "image/gif")
var TypeWebp = NewType("webp", "image/webp")
var TypeCR2 = NewType("cr2", "image/x-canon-cr2")
var TypeTiff = NewType("tif", "image/tiff")
var TypeBmp = NewType("bmp", "image/bmp")
var TypeJxr = NewType("jxr", "image/vnd.ms-photo")
var TypePsd = NewType("psd", "image/vnd.adobe.photoshop")
var TypeIco = NewType("ico", "image/x-icon")

var Image = Map{
	TypeJpeg: Jpeg,
	TypePng:  Png,
	TypeGif:  Gif,
	TypeWebp: Webp,
	TypeCR2:  CR2,
	TypeTiff: Tiff,
	TypeBmp:  Bmp,
	TypeJxr:  Jxr,
	TypePsd:  Psd,
	TypeIco:  Ico,
}

func Jpeg(buf []byte, length int) bool {
	return length > 2 &&
		buf[0] == 0xFF && buf[1] == 0xD8 && buf[2] == 0xFF
}

func Png(buf []byte, length int) bool {
	return length > 3 &&
		buf[0] == 0x89 && buf[1] == 0x50 &&
		buf[2] == 0x4E && buf[3] == 0x47
}

func Gif(buf []byte, length int) bool {
	return length > 2 &&
		buf[0] == 0x47 && buf[1] == 0x49 && buf[2] == 0x46
}

func Webp(buf []byte, length int) bool {
	return length > 11 &&
		buf[8] == 0x57 && buf[9] == 0x45 &&
		buf[10] == 0x42 && buf[11] == 0x50
}

func CR2(buf []byte, length int) bool {
	return length > 9 &&
		((buf[0] == 0x49 && buf[1] == 0x49 && buf[2] == 0x2A && buf[3] == 0x0) ||
			(buf[0] == 0x4D && buf[1] == 0x4D && buf[2] == 0x0 && buf[3] == 0x2A)) &&
		buf[8] == 0x43 && buf[9] == 0x52
}

func Tiff(buf []byte, length int) bool {
	return length > 3 &&
		(buf[0] == 0x49 && buf[1] == 0x49 && buf[2] == 0x2A && buf[3] == 0x0) ||
		(buf[0] == 0x4D && buf[1] == 0x4D && buf[2] == 0x0 && buf[3] == 0x2A)
}

func Bmp(buf []byte, length int) bool {
	return length > 1 &&
		buf[0] == 0x42 && buf[1] == 0x4D
}

func Jxr(buf []byte, length int) bool {
	return length > 2 &&
		buf[0] == 0x49 && buf[1] == 0x49 && buf[2] == 0xBC
}

func Psd(buf []byte, length int) bool {
	return length > 3 &&
		buf[0] == 0x38 && buf[1] == 0x42 &&
		buf[2] == 0x50 && buf[3] == 0x53
}

func Ico(buf []byte, length int) bool {
	return length > 3 &&
		buf[0] == 0x00 && buf[1] == 0x00 &&
		buf[2] == 0x01 && buf[3] == 0x00
}
