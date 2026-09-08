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

// Dbcl - dbcl code.
func (Builder) Dbcl(code Code15) Instr {
	return Dbcl{
		code: immNum(code.Val()),
	}
}

func decodeDbcl(w uint32) Instr {
	return Dbcl{
		base: newBase(w),
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func (i Dbcl) ObjDump(_ disasm.ViewCtx) string {
	return "dbcl " + i.code.text()
}

func (i Dbcl) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["dbcl"][0] | scatterU(i.code.val, 0, 15)

	return writeWord(w, word)
}
