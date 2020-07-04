package matchers

import (
	"github.com/h2non/filetype/types"
)

// Store any types with registered children
var ChildMatchers = make(map[types.Type][]TypeMatcher)

// ChildMatcher creates a TypeMatcher as a child of the parent Type
func ChildMatcher(parent types.Type, kind types.Type, fn ByteMatcher) ByteMatcher {
	matcher := NewMatcher(kind, fn)
	ChildMatchers[parent] = append(ChildMatchers[parent], matcher)
	return fn
}

// AddChild creates and registers a new child for a TypeMatcher
func (m TypeMatcher) AddChild(kind types.Type, fn ByteMatcher) {
	_ = ChildMatcher(m.myType, kind, fn)
}


