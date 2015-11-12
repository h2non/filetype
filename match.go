package filetype

import (
	"gopkg.in/h2non/filetype.v0/matchers"
	"gopkg.in/h2non/filetype.v0/types"
)

// Map of registered file matchers
var Matchers = matchers.Matchers

// Create and register a new matcher
var NewMatcher = matchers.NewMatcher

// Match infers the file type of a given buffer inspecting its magic numbers signature
func Match(buf []byte) (types.Type, error) {
	length := len(buf)
	if length == 0 {
		return types.Unknown, ErrEmptyBuffer
	}

	for _, checker := range Matchers {
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

// AddMatcher registers a new matcher type
func AddMatcher(fileType types.Type, matcher matchers.Matcher) matchers.TypeMatcher {
	return matchers.NewMatcher(fileType, matcher)
}

// Matches checks if the given buffer matches with some supported file type
func Matches(buf []byte) bool {
	kind, _ := Match(buf)
	return kind != types.Unknown
}

// MatchMap performs a file matching againts a map of match functions
func MatchMap(buf []byte, matchers matchers.Map) types.Type {
	for kind, matcher := range matchers {
		if matcher(buf) {
			return kind
		}
	}
	return types.Unknown
}

// MatchesMap is an alias to Matches() but using matching againts a map of match functions
func MatchesMap(buf []byte, matchers matchers.Map) bool {
	return MatchMap(buf, matchers) != types.Unknown
}
