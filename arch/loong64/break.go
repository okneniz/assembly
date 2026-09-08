package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Break - break code (Ud15): the debugger breakpoint trap.
type Break struct {
	base

	code imm
}

// Break - break code (a 15-bit code).
func (Builder) Break(code Code15) Instr {
	return Break{
		code: immNum(code.Val()),
	}
}

func decodeBreak(w uint32) Instr {
	return Break{
		base: newBase(w),
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func (i Break) ObjDump(_ disasm.ViewCtx) string {
	return "break " + i.code.text()
}

func (i Break) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["break"][0] | scatterU(i.code.val, 0, 15)

	return writeWord(w, word)
}
