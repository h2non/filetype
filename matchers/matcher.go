package matchers

import "gopkg.in/h2non/filetype.v0/types"

type Map map[types.Type]Matcher

// Matcher function interface as type alias
type Matcher func([]byte, int) bool

// Type specific matcher function interface
type TypeMatcher func([]byte, int) types.Type

// Store registered file type matchers
var Matchers = make(map[types.Type]TypeMatcher)

var NewType = types.NewType

// Create and register a new type matcher function
func NewMatcher(kind types.Type, fn Matcher) TypeMatcher {
	matcher := func(buf []byte, length int) types.Type {
		if fn(buf, length) {
			return kind
		}
		return types.Unknown
	}
	Matchers[kind] = matcher
	return matcher
}

// Match the file type of a given buffer
func Match(buf []byte) (types.Type, error) {
	length := len(buf)
	if length == 0 {
		return types.Empty, nil
	}

	for _, checker := range Matchers {
		match := checker(buf, length)
		if match != types.Unknown && match.Extension != "" {
			return match, nil
		}
	}

	return types.Unknown, nil
}

func MatchType(buf []byte, kind types.Type) bool {
	return true
}

func MatchMap(buf []byte, matchers Map) types.Type {
	length := len(buf)
	for kind, matcher := range matchers {
		if matcher(buf, length) {
			return kind
		}
	}
	return types.Unknown
}

func MatchesMap(buf []byte, matchers Map) bool {
	return MatchMap(buf, matchers) != types.Unknown
}

func Matches(buf []byte) bool {
	kind, _ := Match(buf)
	return kind != types.Unknown
}

func register(matchers ...Map) {
	for _, m := range matchers {
		for kind, matcher := range m {
			NewMatcher(kind, matcher)
		}
	}
}

func init() {
	register(Image, Video, Audio, Font, Archive)
}
