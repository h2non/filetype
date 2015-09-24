package filetype

import (
	"errors"
	"gopkg.in/h2non/filetype.v0/matchers"
	"gopkg.in/h2non/filetype.v0/types"
)

// Map of extensions and file types
var Types = types.Types

// Map of file matchers
var Matchers = matchers.Matchers

// Default types
var Empty = types.Empty
var Unknown = types.Unknown

// Predefined errors
var EmptyBufferErr = errors.New("Empty buffer")
var UnknownBufferErr = errors.New("Unknown buffer type")

func DoMatch(buf []byte) (types.Type, error) {
	return matchers.Match(buf)
}

// Infer the file type of a buffer inspecting the magic numbers
func Match(buf []byte) (types.Type, error) {
	return DoMatch(buf)
}

// Alias to Match()
func Type(buf []byte) (types.Type, error) {
	return Match(buf)
}

func doMatchMap(buf []byte, machers matchers.Map) (types.Type, error) {
	kind := matchers.MatchMap(buf, machers)
	if kind != types.Unknown {
		return kind, nil
	}
	return kind, UnknownBufferErr
}

// Match file as image type
func Image(buf []byte) (types.Type, error) {
	return doMatchMap(buf, matchers.Image)
}

// Match file as audio type
func Audio(buf []byte) (types.Type, error) {
	return doMatchMap(buf, matchers.Audio)
}

// Match file as video type
func Video(buf []byte) (types.Type, error) {
	return doMatchMap(buf, matchers.Audio)
}

// Match file as text font type
func Font(buf []byte) (types.Type, error) {
	return doMatchMap(buf, matchers.Font)
}

// Match file as generic archive type
func Archive(buf []byte) (types.Type, error) {
	return doMatchMap(buf, matchers.Archive)
}

func Is(buf []byte, ext string) bool {
	kind, ok := types.Types[ext]
	if ok {
		return IsType(buf, kind)
	}
	return false
}

func IsType(buf []byte, kind types.Type) bool {
	matcher := matchers.Matchers[kind]
	if matcher == nil {
		return false
	}

	length := len(buf)
	return matcher(buf, length) != types.Unknown
}

// Register a new matcher type
func AddMatcher(fileType types.Type, matcher matchers.Matcher) matchers.TypeMatcher {
	return matchers.NewMatcher(fileType, matcher)
}

// Register a new file type
func AddType(ext, mime string) types.Type {
	return types.NewType(ext, mime)
}

// Check if a given file extension is supported
func IsSupported(ext string) bool {
	for name, _ := range Types {
		if name == ext {
			return true
		}
	}
	return false
}

// Check if a given MIME expression is supported
func IsMIMESupported(mime string) bool {
	for _, m := range Types {
		if m.MIME.Value == mime {
			return true
		}
	}
	return false
}
