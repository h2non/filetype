package types

var Types = make(map[string]Type)

func Add(t Type) Type {
	Types[t.Extension] = t
	return t
}
