package filetype

import (
	"errors"
	"gopkg.in/h2non/filetype.v0/matchers"
	"gopkg.in/h2non/filetype.v0/types"
)

// Map of supported types
var Types = types.Types

// Create and register a new type
var NewType = types.NewType

// Default types
var Empty = types.Empty
var Unknown = types.Unknown

// Predefined errors
var EmptyBufferErr = errors.New("Empty buffer")
var UnknownBufferErr = errors.New("Unknown buffer type")

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
