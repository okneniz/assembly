package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Nop — nop (no operands).
type Nop struct {
	base
}

const nopMatch = 0xD503201F

// Nop — nop (no operands, fixed encoding).
func (Builder) Nop() Instr {
	return Nop{}
}

func decodeNop(w uint32) Instr {
	return Nop{base: newBase(w)}
}

func (i Nop) ObjDump(_ disasm.ViewCtx) string {
	return "nop"
}

func (i Nop) Encode(w io.Writer) (int64, error) {
	return writeWord(w, nopMatch)
}
