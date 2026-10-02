package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Syscall - syscall code (Ud15): the system-call trap.
type Syscall struct {
	base

	code imm
}

func decodeSyscall(w uint32) Instr {
	return Syscall{
		base: newBase(w),
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func (i Syscall) ObjDump(_ disasm.ViewCtx) string {
	return "syscall " + i.code.text()
}

func (i Syscall) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["syscall"][0] | scatterU(i.code.val, 0, 15)

	return writeWord(w, word)
}
