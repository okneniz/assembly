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

func decodeErtn(w uint32, addr uint64) Instr {
	return Ertn{
		base: newBase(addr, w),
	}
}

func (i Ertn) ObjDump(_ disasm.ViewCtx) string {
	return "ertn"
}

func (i Ertn) Encode(w io.Writer, _ uint64) (int64, error) {
	return writeWord(w, loongEncodings["ertn"][0])
}

func (i Ertn) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ertn",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		nil,
	)
}
