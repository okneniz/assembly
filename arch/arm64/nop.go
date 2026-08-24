package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Nop - nop (no operands).
type Nop struct {
	base
}

const nopMatch = 0xD503201F

func decodeNop(w uint32, addr uint64) Instr {
	return Nop{base: newBase(addr, w)}
}

func (i Nop) ObjDump(_ disasm.ViewCtx) string {
	return "nop"
}

func (i Nop) Encode(w io.Writer, pc uint64) (int64, error) {
	return writeWord(w, nopMatch)
}

func (i Nop) MarshalJSON() ([]byte, error) {
	return i.marshal("nop", i.ObjDump(disasm.DefaultViewCtx()), "Hint", nil)
}
