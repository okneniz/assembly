package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Break - break code (Ud15): the debugger breakpoint trap, the code
// occupies the [14:0] field.
type Break struct {
	base

	code imm
}

// NewBreak - break code (a 15-bit code).
func NewBreak(code Code15) Instr {
	return Break{
		code: immNum(code.Val()),
	}
}

func decodeBreak(w uint32, addr uint64) Instr {
	return Break{
		base: newBase(addr, w),
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func (i Break) ObjDump(_ disasm.ViewCtx) string {
	return "break " + i.code.text()
}

func (i Break) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["break"][0] | scatterU(i.code.val, 0, 15)

	return writeWord(w, word)
}

func (i Break) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"break",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"code": i.code.val},
	)
}
