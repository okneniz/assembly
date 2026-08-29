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

// NewIdle - idle code.
func NewIdle(code Code15) Instr {
	return Idle{
		code: immNum(code.Val()),
	}
}

func decodeIdle(w uint32, addr uint64) Instr {
	return Idle{
		base: newBase(addr, w),
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func (i Idle) ObjDump(_ disasm.ViewCtx) string {
	return "idle " + i.code.text()
}

func (i Idle) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["idle"][0] | scatterU(i.code.val, 0, 15)

	return writeWord(w, word)
}

func (i Idle) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"idle",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		map[string]any{"code": i.code.val},
	)
}
