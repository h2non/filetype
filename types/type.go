package types

type Type struct {
	MIME      MIME
	Extension string
}

func NewType(ext, mime string) Type {
	t := Type{
		MIME:      NewMIME(mime),
		Extension: ext,
	}
	return Add(t)
}
