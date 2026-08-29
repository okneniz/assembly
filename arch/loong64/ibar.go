package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ibar - ibar hint (Ud15): the instruction barrier hint (serializes the
// instruction stream), the hint code occupies the [14:0] field.
type Ibar struct {
	base

	code imm
}

// NewIbar - ibar hint (a 15-bit code).
func NewIbar(code Code15) Instr {
	return Ibar{
		code: immNum(code.Val()),
	}
}

func decodeIbar(w uint32, addr uint64) Instr {
	return Ibar{
		base: newBase(addr, w),
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func (i Ibar) ObjDump(_ disasm.ViewCtx) string {
	return "ibar " + i.code.text()
}

func (i Ibar) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ibar"][0] | scatterU(i.code.val, 0, 15)

	return writeWord(w, word)
}

func (i Ibar) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ibar",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"code": i.code.val},
	)
}
