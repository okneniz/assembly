package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Tlbsrch - tlbsrch: search the TLB for an entry matching the address
// in the TLB search registers (no operands).
type Tlbsrch struct {
	base
}

// Tlbsrch - tlbsrch (no operands).
func (Builder) Tlbsrch() Instr {
	return Tlbsrch{}
}

func decodeTlbsrch(w uint32, addr uint64) Instr {
	return Tlbsrch{
		base: newBase(addr, w),
	}
}

func (i Tlbsrch) ObjDump(_ disasm.ViewCtx) string {
	return "tlbsrch"
}

func (i Tlbsrch) Encode(w io.Writer, _ uint64) (int64, error) {
	return writeWord(w, loongEncodings["tlbsrch"][0])
}

func (i Tlbsrch) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"tlbsrch",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		nil,
	)
}
