package matchers

import (
	"sort"

	"gopkg.in/h2non/filetype.v1/types"
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

// MatcherTypes store sorted matcher key
var MatcherTypes = make([]*types.Type, 0)

type typs []*types.Type

func (t typs) Len() int           { return len(t) }
func (t typs) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t typs) Less(i, j int) bool { return t[i].Priority < t[j].Priority }

// Create and register a new type matcher function
func NewMatcher(kind types.Type, fn Matcher) TypeMatcher {
	matcher := func(buf []byte) types.Type {
		if fn(buf) {
			return kind
		}
		return types.Unknown
	}

	Matchers[kind] = matcher
	MatcherTypes = append(MatcherTypes, &kind)

	return matcher
}

// When iterating over a map with a range loop,
// the iteration order is not specified and is not guaranteed to be the same from one iteration to the next
// If you require a stable iteration order you must maintain a separate data structure that specifies that order
// see: https://blog.golang.org/go-maps-in-action
func register(matchers ...Map) {
	for _, m := range matchers {
		for kind, matcher := range m {
			NewMatcher(kind, matcher)
		}
	}

	sort.Sort(sort.Reverse(typs(MatcherTypes)))
}

func init() {
	// Arguments order is intentional
	register(Image, Video, Audio, Font, Document, Archive)
}
