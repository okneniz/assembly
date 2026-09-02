package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Tlbfill - tlbfill: fill the TLB from the TLB entry registers (no
// operands).
type Tlbfill struct {
	base
}

// Tlbfill - tlbfill (no operands).
func (Builder) Tlbfill() Instr {
	return Tlbfill{}
}

func decodeTlbfill(w uint32, addr uint64) Instr {
	return Tlbfill{
		base: newBase(addr, w),
	}
}

func (i Tlbfill) ObjDump(_ disasm.ViewCtx) string {
	return "tlbfill"
}

func (i Tlbfill) Encode(w io.Writer, _ uint64) (int64, error) {
	return writeWord(w, loongEncodings["tlbfill"][0])
}

func (i Tlbfill) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"tlbfill",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		nil,
	)
}
