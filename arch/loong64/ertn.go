package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ertn - ertn: return from exception (no operands).
type Ertn struct {
	base
}

// Ertn - ertn (no operands).
func (Builder) Ertn() Instr {
	return Ertn{}
}

func decodeErtn(w uint32) Instr {
	return Ertn{
		base: newBase(w),
	}
}

func (i Ertn) ObjDump(_ disasm.ViewCtx) string {
	return "ertn"
}

func (i Ertn) Encode(w io.Writer) (int64, error) {
	return writeWord(w, loongEncodings["ertn"][0])
}
