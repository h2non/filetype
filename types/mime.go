package types

type MIME struct {
	Type    string
	Subtype string
	Value   string
}

func NewMIME(mime string) MIME {
	kind, subtype := splitMime(mime)
	return MIME{Type: kind, Subtype: subtype, Value: mime}
}
