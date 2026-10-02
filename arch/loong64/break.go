package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Break - break code (Ud15): the debugger breakpoint trap.
type Break struct {
	code imm
}

func (i Break) ObjDump(_ disasm.ViewCtx) string {
	return "break " + i.code.text()
}

func (i Break) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["break"][0] | scatterU(i.code.val, 0, 15)

	return writeWord(w, word)
}
