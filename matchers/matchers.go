package matchers

import (
	"github.com/h2non/filetype/types"
)

// Internal shortcut to NewType
var newType = types.NewType

type Matcher interface {
	Match([]byte) bool
}

type Typer interface {
	Type([]byte) types.Type
}

// ByteMatcher function interface as type alias
type ByteMatcher func([]byte) bool

// Implement Matcher interface for ByteMatcher
func (b ByteMatcher) Match(buf []byte) bool {
	return b(buf)
}

// A TypeMatcher is both a Typer and Matcher
type TypeMatcher struct {
	myType  types.Type
	matcher ByteMatcher
}

// Implement Matcher interface for TypeMatcher
func (m TypeMatcher) Match(buf []byte) bool {
	// Check any matchers for child types
	for _, c := range ChildMatchers[m.myType] {
		if c.matcher.Match(buf) {
			return true
		}
	}
	return m.matcher.Match(buf)
}

// Implement Typer interface for TypeMatcher
func (m TypeMatcher) Type(buf []byte) types.Type {
	// Return a matching child type, if any
	for _, c := range ChildMatchers[m.myType] {
		if c.matcher.Match(buf) {
			return c.myType
		}
	}
	if m.matcher.Match(buf) {
		return m.myType
	}
	return types.Unknown
}

// Store registered file type matchers
var Matchers = make(map[types.Type]TypeMatcher)
var MatcherKeys []types.Type

// Type interface to store pairs of type with its matcher
type Map map[types.Type]ByteMatcher

// Create and register a new type matcher
func NewMatcher(kind types.Type, fn ByteMatcher) TypeMatcher {
	m := TypeMatcher{kind, fn}
	Matchers[kind] = m
	// prepend here so any user defined matchers get added first
	MatcherKeys = append([]types.Type{kind}, MatcherKeys...)
	return m
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
	// Archive files will be checked last due to prepend above in func NewMatcher
	register(Archive, Document, Font, Audio, Video, Image, Application)
}
