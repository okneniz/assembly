package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Idle - idle code: wait for an interrupt; the ui15 code selects the
// idle level (0 in the base ISA).
type Idle struct {
	base

	code imm
}

func decodeIdle(w uint32) Instr {
	return Idle{
		base: newBase(w),
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func (i Idle) ObjDump(_ disasm.ViewCtx) string {
	return "idle " + i.code.text()
}

func (i Idle) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["idle"][0] | scatterU(i.code.val, 0, 15)

	return writeWord(w, word)
}
