package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Syscall - syscall code (Ud15): the system-call trap.
type Syscall struct {
	code imm
}

func (i Syscall) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["syscall"][0] | scatterU(i.code.val, 0, 15)

	return writeWord(w, word)
}

func (i Syscall) ObjDump(_ disasm.ViewCtx) string {
	return "syscall " + i.code.text()
}
