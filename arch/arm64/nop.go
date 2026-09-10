package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Nop — nop (no operands).
type Nop struct {
	base
}

// newNop - the Nop constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newNop(b base) (Nop, error) { //nolint:unparam // uniform (Instr, error) decodeCtor type
	return Nop{
		base: b,
	}, nil
}

const nopMatch = 0xD503201F

func (i Nop) ObjDump(_ disasm.ViewCtx) string {
	return "nop"
}

func (i Nop) Encode(w io.Writer) (int64, error) {
	return writeWord(w, nopMatch)
}

// Nop — nop (no operands, fixed encoding).
func (Builder) Nop() Instr {
	in, err := newNop(base{})
	if err != nil {
		panic(err) // no operands - cannot fail
	}

	return in
}

func decodeNop(w uint32) (Instr, error) {
	in, err := newNop(newBase(w))
	if err != nil {
		return nil, err
	}

	return in, nil
}
