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

// Tlbclr - tlbclr (no operands).
func (Builder) Tlbclr() Instr {
	return Tlbclr{}
}

func decodeTlbclr(w uint32) Instr {
	return Tlbclr{
		base: newBase(w),
	}
}

func (i Tlbclr) ObjDump(_ disasm.ViewCtx) string {
	return "tlbclr"
}

func (i Tlbclr) Encode(w io.Writer) (int64, error) {
	return writeWord(w, loongEncodings["tlbclr"][0])
}
