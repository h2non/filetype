package matchers

import (
	"github.com/h2non/filetype/types"
)

// Internal shortcut to NewType
var newType = types.NewType

// Matcher function interface as type alias
type Matcher func([]byte) bool

// Type interface to store pairs of type with its matcher function
type Map map[types.Type]Matcher

// Type specific matcher function interface
type TypeMatcher func([]byte) types.Type

// Store registered file type matchers
var Matchers = make(map[types.Type]TypeMatcher)
var MatcherKeys []types.Type

// Create and register a new type matcher function
func NewMatcher(kind types.Type, fn Matcher) TypeMatcher {
	matcher := func(buf []byte) types.Type {
		if fn(buf) {
			return kind
		}
		return types.Unknown
	}

	Matchers[kind] = matcher
	MatcherKeys = append(MatcherKeys, kind)
	return matcher
}

func register(matchers ...Map) {
	MatcherKeys = MatcherKeys[:0]
	for _, m := range matchers {
		for kind, matcher := range m {
			NewMatcher(kind, matcher)
		}
	}
}

func init() {
	// Arguments order is intentional
	register(Image, Video, Audio, Font, Document, Archive)
}
