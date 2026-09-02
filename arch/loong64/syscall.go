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

// Syscall - syscall code (a 15-bit code).
func (Builder) Syscall(code Code15) Instr {
	return Syscall{
		code: immNum(code.Val()),
	}
}

func decodeSyscall(w uint32, addr uint64) Instr {
	return Syscall{
		base: newBase(addr, w),
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func (i Syscall) ObjDump(_ disasm.ViewCtx) string {
	return "syscall " + i.code.text()
}

func (i Syscall) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["syscall"][0] | scatterU(i.code.val, 0, 15)

	return writeWord(w, word)
}

func (i Syscall) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"syscall",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"code": i.code.val},
	)
}
