package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Tlbclr - tlbclr: invalidate the TLB entry selected by the TLB index
// (no operands).
type Tlbclr struct {
	base
}

// NewTlbclr - tlbclr (no operands).
func NewTlbclr() Instr {
	return Tlbclr{}
}

func decodeTlbclr(w uint32, addr uint64) Instr {
	return Tlbclr{
		base: newBase(addr, w),
	}
}

func (i Tlbclr) ObjDump(_ disasm.ViewCtx) string {
	return "tlbclr"
}

func (i Tlbclr) Encode(w io.Writer, _ uint64) (int64, error) {
	return writeWord(w, loongEncodings["tlbclr"][0])
}

func (i Tlbclr) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"tlbclr",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		nil,
	)
}
