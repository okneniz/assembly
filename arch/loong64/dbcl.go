package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Dbcl - dbcl code (I15): the debug-call breakpoint trap.
type Dbcl struct {
	code imm
}

func (i Dbcl) ObjDump(_ disasm.ViewCtx) string {
	return "dbcl " + i.code.text()
}

func (i Dbcl) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["dbcl"][0] | scatterU(i.code.val, 0, 15)

	return writeWord(w, word)
}
