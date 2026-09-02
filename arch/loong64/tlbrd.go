package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Tlbrd - tlbrd: read the indexed TLB entry into the TLB entry
// registers (no operands).
type Tlbrd struct {
	base
}

// Tlbrd - tlbrd (no operands).
func (Builder) Tlbrd() Instr {
	return Tlbrd{}
}

func decodeTlbrd(w uint32, addr uint64) Instr {
	return Tlbrd{
		base: newBase(addr, w),
	}
}

func (i Tlbrd) ObjDump(_ disasm.ViewCtx) string {
	return "tlbrd"
}

func (i Tlbrd) Encode(w io.Writer, _ uint64) (int64, error) {
	return writeWord(w, loongEncodings["tlbrd"][0])
}

func (i Tlbrd) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"tlbrd",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		nil,
	)
}
