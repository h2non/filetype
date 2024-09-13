[English](./README.md) | 简体中文
# filetype [![GoDoc](https://godoc.org/github.com/h2non/filetype?status.svg)](https://godoc.org/github.com/h2non/filetype) [![Go Version](https://img.shields.io/badge/go-v1.0+-green.svg?style=flat)](https://github.com/h2non/gentleman)

小型且无依赖的 [Go](https://golang.org) 包，用于通过检查[魔数](https://en.wikipedia.org/wiki/Magic_number_(programming)#Magic_numbers_in_files)签名推断文件类型和 MIME 类型。

如需检查 SVG 文件类型，请参阅 [go-is-svg](https://github.com/h2non/go-is-svg) 包。Python 版本: [filetype.py](https://github.com/h2non/filetype.py)。

## 功能特性

- 支持[广泛的](#supported-types)文件类型
- 提供文件扩展名和正确的 MIME 类型
- 通过扩展名或 MIME 类型进行文件识别
- 通过类别（图像、视频、音频等）进行文件识别
- 提供大量辅助工具和文件匹配快捷方式
- [可扩展](#add-additional-file-type-matchers)：可以添加自定义的新类型和匹配器
- 简单且语义化的 API
- [速度极快](#benchmarks)，即使处理大型文件也不例外
- 只需要前 262 字节的文件头，因此你可以[传递一个切片](#file-header)
- 无依赖（仅需 Go 代码，无需 C 语言编译）
- 跨平台文件识别

## 安装

```bash
go get github.com/h2non/filetype
```

## API

请参考[Godoc](https://godoc.org/github.com/h2non/filetype) .

### 关联子包

- [`github.com/h2non/filetype/types`](https://godoc.org/github.com/h2non/filetype/types)
- [`github.com/h2non/filetype/matchers`](https://godoc.org/github.com/h2non/filetype/matchers)

## 示例

#### 简单的文件类型检查

```go
package main

import (
  "fmt"
  "io/ioutil"

  "github.com/h2non/filetype"
)

func main() {
  buf, _ := ioutil.ReadFile("sample.jpg")

  kind, _ := filetype.Match(buf)
  if kind == filetype.Unknown {
    fmt.Println("Unknown file type")
    return
  }

  fmt.Printf("File type: %s. MIME: %s\n", kind.Extension, kind.MIME.Value)
}
```

#### 检查类型类别

```go
package main

import (
  "fmt"
  "io/ioutil"

  "github.com/h2non/filetype"
)

func main() {
  buf, _ := ioutil.ReadFile("sample.jpg")

  if filetype.IsImage(buf) {
    fmt.Println("File is an image")
  } else {
    fmt.Println("Not an image")
  }
}
```

#### 类型支持

```go
package main

import (
  "fmt"

  "github.com/h2non/filetype"
)

func main() {
  // Check if file is supported by extension
  if filetype.IsSupported("jpg") {
    fmt.Println("Extension supported")
  } else {
    fmt.Println("Extension not supported")
  }

  // Check if file is supported by extension
  if filetype.IsMIMESupported("image/jpeg") {
    fmt.Println("MIME type supported")
  } else {
    fmt.Println("MIME type not supported")
  }
}
```

#### 文件头

```go
package main

import (
  "fmt"
  "os"

  "github.com/h2non/filetype"
)

func main() {
  // Open a file descriptor
  file, _ := os.Open("movie.mp4")

  // We only have to pass the file header = first 261 bytes
  head := make([]byte, 261)
  file.Read(head)

  if filetype.IsImage(head) {
    fmt.Println("File is an image")
  } else {
    fmt.Println("Not an image")
  }
}
```

#### 添加额外的文件类型匹配器

```go
package main

import (
  "fmt"

  "github.com/h2non/filetype"
)

var fooType = filetype.NewType("foo", "foo/foo")

func fooMatcher(buf []byte) bool {
  return len(buf) > 1 && buf[0] == 0x01 && buf[1] == 0x02
}

func main() {
  // Register the new matcher and its type
  filetype.AddMatcher(fooType, fooMatcher)

  // Check if the new type is supported by extension
  if filetype.IsSupported("foo") {
    fmt.Println("New supported type: foo")
  }

  // Check if the new type is supported by MIME
  if filetype.IsMIMESupported("foo/foo") {
    fmt.Println("New supported MIME type: foo/foo")
  }

  // Try to match the file
  fooFile := []byte{0x01, 0x02}
  kind, _ := filetype.Match(fooFile)
  if kind == filetype.Unknown {
    fmt.Println("Unknown file type")
  } else {
    fmt.Printf("File type matched: %s\n", kind.Extension)
  }
}
```

## 已支持的类型

#### Image

- **jpg** - `image/jpeg`
- **png** - `image/png`
- **gif** - `image/gif`
- **webp** - `image/webp`
- **cr2** - `image/x-canon-cr2`
- **tif** - `image/tiff`
- **bmp** - `image/bmp`
- **heif** - `image/heif`
- **jxr** - `image/vnd.ms-photo`
- **psd** - `image/vnd.adobe.photoshop`
- **ico** - `image/vnd.microsoft.icon`
- **dwg** - `image/vnd.dwg`
- **avif** - `image/avif`

#### Video

- **mp4** - `video/mp4`
- **m4v** - `video/x-m4v`
- **mkv** - `video/x-matroska`
- **webm** - `video/webm`
- **mov** - `video/quicktime`
- **avi** - `video/x-msvideo`
- **wmv** - `video/x-ms-wmv`
- **mpg** - `video/mpeg`
- **flv** - `video/x-flv`
- **3gp** - `video/3gpp`

#### Audio

- **mid** - `audio/midi`
- **mp3** - `audio/mpeg`
- **m4a** - `audio/mp4`
- **ogg** - `audio/ogg`
- **flac** - `audio/x-flac`
- **wav** - `audio/x-wav`
- **amr** - `audio/amr`
- **aac** - `audio/aac`
- **aiff** - `audio/x-aiff`

#### Archive

- **epub** - `application/epub+zip`
- **zip** - `application/zip`
- **tar** - `application/x-tar`
- **rar** - `application/vnd.rar`
- **gz** - `application/gzip`
- **bz2** - `application/x-bzip2`
- **7z** - `application/x-7z-compressed`
- **xz** - `application/x-xz`
- **zstd** - `application/zstd`
- **pdf** - `application/pdf`
- **exe** - `application/vnd.microsoft.portable-executable`
- **swf** - `application/x-shockwave-flash`
- **rtf** - `application/rtf`
- **iso** - `application/x-iso9660-image`
- **eot** - `application/octet-stream`
- **ps** - `application/postscript`
- **sqlite** - `application/vnd.sqlite3`
- **nes** - `application/x-nintendo-nes-rom`
- **crx** - `application/x-google-chrome-extension`
- **cab** - `application/vnd.ms-cab-compressed`
- **deb** - `application/vnd.debian.binary-package`
- **ar** - `application/x-unix-archive`
- **Z** - `application/x-compress`
- **lz** - `application/x-lzip`
- **lz4** - `application/x-lz4`
- **rpm** - `application/x-rpm`
- **elf** - `application/x-executable`
- **dcm** - `application/dicom`

#### Documents

- **doc** - `application/msword`
- **docx** - `application/vnd.openxmlformats-officedocument.wordprocessingml.document`
- **xls** - `application/vnd.ms-excel`
- **xlsx** - `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
- **ppt** - `application/vnd.ms-powerpoint`
- **pptx** - `application/vnd.openxmlformats-officedocument.presentationml.presentation`

#### Font

- **woff** - `application/font-woff`
- **woff2** - `application/font-woff`
- **ttf** - `application/font-sfnt`
- **otf** - `application/font-sfnt`

#### Application

- **wasm** - `application/wasm`
- **dex** - `application/vnd.android.dex`
- **dey** - `application/vnd.android.dey`

## 性能基准

通过 [real files](https://github.com/h2non/filetype/tree/master/fixtures) 测量.

运行环境: OSX x64 i7 2.7 Ghz

```bash
BenchmarkMatchTar-8    1000000        1083 ns/op
BenchmarkMatchZip-8    1000000        1162 ns/op
BenchmarkMatchJpeg-8   1000000        1280 ns/op
BenchmarkMatchGif-8    1000000        1315 ns/op
BenchmarkMatchPng-8    1000000        1121 ns/op
```

## License

MIT - Tomas Aparicio
