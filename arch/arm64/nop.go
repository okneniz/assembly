package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Nop — nop (no operands).
type Nop struct {
}

// newNop - the Nop constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newNop() (Nop, error) {
	return Nop{}, nil
}

const nopMatch = 0xD503201F

func (i Nop) Encode(w io.Writer) (int64, error) {
	return writeWord(w, nopMatch)
}

func (i Nop) ObjDump(_ disasm.ViewCtx) string {
	return "nop"
}
