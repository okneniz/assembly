package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Tlbflush - tlbflush: invalidate all TLB entries (no operands).
type Tlbflush struct {
	base
}

// NewTlbflush - tlbflush (no operands).
func NewTlbflush() Instr {
	return Tlbflush{}
}

func decodeTlbflush(w uint32, addr uint64) Instr {
	return Tlbflush{
		base: newBase(addr, w),
	}
}

func (i Tlbflush) ObjDump(_ disasm.ViewCtx) string {
	return "tlbflush"
}

func (i Tlbflush) Encode(w io.Writer, _ uint64) (int64, error) {
	return writeWord(w, loongEncodings["tlbflush"][0])
}

func (i Tlbflush) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"tlbflush",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		nil,
	)
}
