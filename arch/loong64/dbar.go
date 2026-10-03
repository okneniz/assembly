package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Dbar - dbar hint (Ud15): the data barrier hint (serializes data accesses).
type Dbar struct {
	code imm
}

func (i Dbar) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["dbar"][0] | scatterU(i.code.val, 0, 15)

	return writeWord(w, word)
}

func (i Dbar) ObjDump(_ disasm.ViewCtx) string {
	return "dbar " + i.code.text()
}
