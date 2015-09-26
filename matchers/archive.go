package matchers

import "bytes"

var TypeEpub = newType("epub", "application/epub+zip")
var TypeZip = newType("zip", "application/zip")
var TypeTar = newType("tar", "application/x-tar")
var TypeRar = newType("rar", "application/x-rar-compressed")
var TypeGz = newType("gz", "application/gzip")
var TypeBz2 = newType("bz2", "application/x-bzip2")
var Type7z = newType("7z", "application/x-7z-compressed")
var TypeXz = newType("xz", "application/x-xz")
var TypePdf = newType("pdf", "application/pdf")
var TypeExe = newType("exe", "application/x-msdownload")
var TypeSwf = newType("swf", "application/x-shockwave-flash")
var TypeRtf = newType("rtf", "application/rtf")
var TypeEot = newType("eot", "application/octet-stream")
var TypePs = newType("ps", "application/postscript")
var TypeSqlite = newType("sqlite", "application/x-sqlite3")

var Archive = Map{
	TypeEpub:   Epub,
	TypeZip:    Zip,
	TypeTar:    Tar,
	TypeRar:    Rar,
	TypeBz2:    Bz2,
	Type7z:     SevenZ,
	TypeXz:     Xz,
	TypePdf:    Pdf,
	TypeExe:    Exe,
	TypeSwf:    Swf,
	TypeRtf:    Rtf,
	TypeEot:    Eot,
	TypePs:     Ps,
	TypeSqlite: Sqlite,
}

var (
	magicTar = []byte{0x75, 0x73, 0x74, 0x61, 0x72}
)

func Epub(buf []byte, length int) bool {
	return length > 57 &&
		buf[0] == 0x50 && buf[1] == 0x4B && buf[2] == 0x3 && buf[3] == 0x4 &&
		buf[30] == 0x6D && buf[31] == 0x69 && buf[32] == 0x6D && buf[33] == 0x65 &&
		buf[34] == 0x74 && buf[35] == 0x79 && buf[36] == 0x70 && buf[37] == 0x65 &&
		buf[38] == 0x61 && buf[39] == 0x70 && buf[40] == 0x70 && buf[41] == 0x6C &&
		buf[42] == 0x69 && buf[43] == 0x63 && buf[44] == 0x61 && buf[45] == 0x74 &&
		buf[46] == 0x69 && buf[47] == 0x6F && buf[48] == 0x6E && buf[49] == 0x2F &&
		buf[50] == 0x65 && buf[51] == 0x70 && buf[52] == 0x75 && buf[53] == 0x62 &&
		buf[54] == 0x2B && buf[55] == 0x7A && buf[56] == 0x69 && buf[57] == 0x70
}

func Zip(buf []byte, length int) bool {
	return length > 3 &&
		buf[0] == 0x50 && buf[1] == 0x4B &&
		(buf[2] == 0x3 || buf[2] == 0x5 || buf[2] == 0x7) &&
		(buf[3] == 0x4 || buf[3] == 0x6 || buf[3] == 0x8)
}

func Tar(buf []byte, length int) bool {
	return length > 261 &&
		bytes.Equal(buf[257:257+len(magicTar)], magicTar[:])
}

func Rar(buf []byte, length int) bool {
	return length > 6 &&
		buf[0] == 0x52 && buf[1] == 0x61 && buf[2] == 0x72 &&
		buf[3] == 0x21 && buf[4] == 0x1A && buf[5] == 0x7 &&
		(buf[6] == 0x0 || buf[6] == 0x1)
}

func Gz(buf []byte, length int) bool {
	return length > 2 &&
		buf[0] == 0x1F && buf[1] == 0x8B && buf[2] == 0x8
}

func Bz2(buf []byte, length int) bool {
	return length > 2 &&
		buf[0] == 0x42 && buf[1] == 0x5A && buf[2] == 0x68
}

func SevenZ(buf []byte, length int) bool {
	return length > 5 &&
		buf[0] == 0x37 && buf[1] == 0x7A && buf[2] == 0xBC &&
		buf[3] == 0xAF && buf[4] == 0x27 && buf[5] == 0x1C
}

func Pdf(buf []byte, length int) bool {
	return length > 3 &&
		buf[0] == 0x25 && buf[1] == 0x50 &&
		buf[2] == 0x44 && buf[3] == 0x46
}

func Exe(buf []byte, length int) bool {
	return length > 1 &&
		buf[0] == 0x4D && buf[1] == 0x5A
}

func Swf(buf []byte, length int) bool {
	return length > 2 &&
		(buf[0] == 0x43 || buf[0] == 0x46) &&
		buf[1] == 0x57 && buf[2] == 0x53
}

func Rtf(buf []byte, length int) bool {
	return length > 4 &&
		buf[0] == 0x7B && buf[1] == 0x5C &&
		buf[2] == 0x72 && buf[3] == 0x74 &&
		buf[4] == 0x66
}

func Eot(buf []byte, length int) bool {
	return length > 10 &&
		buf[34] == 0x4C && buf[35] == 0x50 &&
		((buf[8] == 0x02 && buf[9] == 0x00 &&
			buf[10] == 0x01) || (buf[8] == 0x01 &&
			buf[9] == 0x00 && buf[10] == 0x00) ||
			(buf[8] == 0x02 && buf[9] == 0x00 && buf[10] == 0x02))
}

func Ps(buf []byte, length int) bool {
	return length > 1 &&
		buf[0] == 0x25 && buf[1] == 0x21
}

func Xz(buf []byte, length int) bool {
	return length > 5 &&
		buf[0] == 0xFD && buf[1] == 0x37 &&
		buf[2] == 0x7A && buf[3] == 0x58 &&
		buf[4] == 0x5A && buf[5] == 0x00
}

func Sqlite(buf []byte, length int) bool {
	return length > 3 &&
		buf[0] == 0x53 && buf[1] == 0x51 &&
		buf[2] == 0x4C && buf[3] == 0x69
}
