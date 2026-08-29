package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Tlbwr - tlbwr: write the TLB entry registers to the indexed TLB
// entry without a search (no operands).
type Tlbwr struct {
	base
}

// NewTlbwr - tlbwr (no operands).
func NewTlbwr() Instr {
	return Tlbwr{}
}

func decodeTlbwr(w uint32, addr uint64) Instr {
	return Tlbwr{
		base: newBase(addr, w),
	}
}

func (i Tlbwr) ObjDump(_ disasm.ViewCtx) string {
	return "tlbwr"
}

func (i Tlbwr) Encode(w io.Writer, _ uint64) (int64, error) {
	return writeWord(w, loongEncodings["tlbwr"][0])
}

func (i Tlbwr) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"tlbwr",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		nil,
	)
}
