package arm64

// Field describes one bit field of an instruction: offset, width and an
// optional transform — the name of the inverse transform for the
// assembler's legacy path (see transform_registry.go). Used by the decode
// table (Schema.Fields; isa_generated.go).
type Field struct {
	Name      string
	Offset    uint
	Width     uint
	Transform string
}

// NewField — a named bit field; the optional transform — the name of the
// inverse transform (see transform_registry.go), without it — a field
// with no conversion.
func NewField(name string, offset, width uint, transform ...string) Field {
	f := Field{
		Name:   name,
		Offset: offset,
		Width:  width,
	}
	if len(transform) > 0 {
		f.Transform = transform[0]
	}

	return f
}
