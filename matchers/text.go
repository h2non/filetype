package matchers

var (
	TypeHtml = newType("html", "application/html")
)

var Text = Map{
	TypeHtml: Html,
}

// Wasm detects a Web Assembly 1.0 filetype.
func Html(buf []byte) bool {
	// WASM has starts with `\0asm`, followed by the version.
	// http://webassembly.github.io/spec/core/binary/modules.html#binary-magic
	return len(buf) >= 14 &&
		buf[0] == 0x3C && buf[1] == 0x21 &&
		buf[2] == 0x64 && buf[3] == 0x6F &&
		buf[4] == 0x63 && buf[5] == 0x74 &&
		buf[6] == 0x79 && buf[7] == 0x70 &&
		buf[8] == 0x65 && buf[9] == 0x20 &&
		buf[10] == 0x68 && buf[11] == 0x74 &&
		buf[12] == 0x6D && buf[13] == 0x6C
}