package types

var Types = make(map[string]Type)

// Register a new type
func Add(t Type) Type {
	Types[t.Extension] = t
	return t
}

// Retrieve a Type by extension
func Get(ext string) Type {
	kind := Types[ext]
	if kind.Extension != "" {
		return kind
	}
	return Unknown
}
