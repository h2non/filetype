package filetype

import (
	"io"
	"os"

	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
)

// Matchers is an alias to matchers.Matchers
var Matchers = matchers.Matchers

// MatcherKeys is an alias to matchers.MatcherKeys
var MatcherKeys = &matchers.MatcherKeys

// NewMatcher is an alias to matchers.NewMatcher
var NewMatcher = matchers.NewMatcher

const maxBufSize = 8192

// MatcherWriter is a matcher that coplies to the writer interface
type MatcherWriter struct {
	buf []byte
}

func (mb *MatcherWriter) Write(p []byte) (n int, err error) {
	incomingSize := len(p)
	currentSize := len(mb.buf)
	if currentSize < maxBufSize {
		// write only when the current size of the buffer is less than the max buffer size
		newSize := currentSize + incomingSize
		overflow := newSize - maxBufSize
		reSlice := p
		if overflow > 0 {
			// if the maxBufSize is exceeded by the new size, we need to do some re slicing
			maxLen := incomingSize - overflow
			reSlice = p[0:maxLen]
		}
		mb.buf = append(mb.buf, reSlice...)
	}
	return incomingSize, nil
}

// Match calls the Match function with the inner buffer of the MatcherWriter
func (mb *MatcherWriter) Match() (types.Type, error) {
	if mb.buf == nil {
		return types.Unknown, ErrEmptyBuffer
	}
	return Match(mb.buf)
}

// NewMatcherWriter creates a matcher which is a io.Writer
func NewMatcherWriter() *MatcherWriter {
	return &MatcherWriter{
		buf: make([]byte, 0, maxBufSize),
	}
}

// Match infers the file type of a given buffer inspecting its magic numbers signature
func Match(buf []byte) (types.Type, error) {
	length := len(buf)
	if length == 0 {
		return types.Unknown, ErrEmptyBuffer
	}

	for _, kind := range *MatcherKeys {
		checker := Matchers[kind]
		match := checker(buf)
		if match != types.Unknown && match.Extension != "" {
			return match, nil
		}
	}

	return types.Unknown, nil
}

// Get is an alias to Match()
func Get(buf []byte) (types.Type, error) {
	return Match(buf)
}

// MatchFile infers a file type for a file
func MatchFile(filepath string) (types.Type, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return types.Unknown, err
	}
	defer file.Close()

	return MatchReader(file)
}

// MatchReader is convenient wrapper to Match() any Reader
func MatchReader(reader io.Reader) (types.Type, error) {
	buffer := make([]byte, 8192) // 8K makes msooxml tests happy and allows for expanded custom file checks

	_, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return types.Unknown, err
	}

	return Match(buffer)
}

// AddMatcher registers a new matcher type
func AddMatcher(fileType types.Type, matcher matchers.Matcher) matchers.TypeMatcher {
	return matchers.NewMatcher(fileType, matcher)
}

// Matches checks if the given buffer matches with some supported file type
func Matches(buf []byte) bool {
	kind, _ := Match(buf)
	return kind != types.Unknown
}

// MatchMap performs a file matching against a map of match functions
func MatchMap(buf []byte, matchers matchers.Map) types.Type {
	for kind, matcher := range matchers {
		if matcher(buf) {
			return kind
		}
	}
	return types.Unknown
}

// MatchesMap is an alias to Matches() but using matching against a map of match functions
func MatchesMap(buf []byte, matchers matchers.Map) bool {
	return MatchMap(buf, matchers) != types.Unknown
}
