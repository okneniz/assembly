package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Tlbflush - tlbflush: invalidate all TLB entries (no operands).
type Tlbflush struct {
	base
}

// Tlbflush - tlbflush (no operands).
func (Builder) Tlbflush() Instr {
	return Tlbflush{}
}

func decodeTlbflush(w uint32) Instr {
	return Tlbflush{
		base: newBase(w),
	}
}

func (i Tlbflush) ObjDump(_ disasm.ViewCtx) string {
	return "tlbflush"
}

func (i Tlbflush) Encode(w io.Writer) (int64, error) {
	return writeWord(w, loongEncodings["tlbflush"][0])
}
