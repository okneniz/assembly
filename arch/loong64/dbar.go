package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Dbar - dbar hint (Ud15): the data barrier hint (serializes data accesses).
type Dbar struct {
	base

	code imm
}

// Dbar - dbar hint (a 15-bit code).
func (Builder) Dbar(code Code15) Instr {
	return Dbar{
		code: immNum(code.Val()),
	}
}

func decodeDbar(w uint32, addr uint64) Instr {
	return Dbar{
		base: newBase(addr, w),
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func (i Dbar) ObjDump(_ disasm.ViewCtx) string {
	return "dbar " + i.code.text()
}

func (i Dbar) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["dbar"][0] | scatterU(i.code.val, 0, 15)

	return writeWord(w, word)
}

func (i Dbar) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"dbar",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"code": i.code.val},
	)
}
