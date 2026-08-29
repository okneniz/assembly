package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Dbcl - dbcl code (I15): the debug-call breakpoint trap.
type Dbcl struct {
	base

	code imm
}

// NewDbcl - dbcl code.
func NewDbcl(code Code15) Instr {
	return Dbcl{
		code: immNum(code.Val()),
	}
}

func decodeDbcl(w uint32, addr uint64) Instr {
	return Dbcl{
		base: newBase(addr, w),
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func (i Dbcl) ObjDump(_ disasm.ViewCtx) string {
	return "dbcl " + i.code.text()
}

func (i Dbcl) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["dbcl"][0] | scatterU(i.code.val, 0, 15)

	return writeWord(w, word)
}

func (i Dbcl) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"dbcl",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		map[string]any{"code": i.code.val},
	)
}
