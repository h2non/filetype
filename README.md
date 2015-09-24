# filetype [![Build Status](https://travis-ci.org/h2non/filetype.png)](https://travis-ci.org/h2non/filetype) [![GoDoc](https://godoc.org/github.com/h2non/filetype?status.svg)](https://godoc.org/github.com/h2non/filetype)

Small [Go](https://golang.org) package to infer the file type checking the [magic number](https://en.wikipedia.org/wiki/Magic_number_(programming)#Magic_numbers_in_files) of a given binary buffer.

Supports a wide range of file types, including images formats, fonts, videos, audio and other common application files, and provides the proper file extension and convenient MIME code.

## Installation

```bash
go get gopkg.in/h2non/filetype.v0
```

## Usage

```go
import (
  "fmt"
  "io/ioutil"
  "gopkg.in/h2non/filetype.v0"
)

func main() {
  buf, _ := ioutil.ReadFile("sample.jpg")

  kind, unkwown := filetype.Type(buf)
  if unkwown != nil {
    fmt.Printf("Unkwown file type")
    return
  }

  fmt.Printf("File type found: %s. MIME: %s", kind.Extension, kind.MIME.Value)
}
```

## API



## License

MIT - Tomas Aparicio
