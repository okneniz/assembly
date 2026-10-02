package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ibar - ibar hint (Ud15): the instruction barrier hint (serializes the instruction stream).
type Ibar struct {
	code imm
}

func (i Ibar) ObjDump(_ disasm.ViewCtx) string {
	return "ibar " + i.code.text()
}

func (i Ibar) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ibar"][0] | scatterU(i.code.val, 0, 15)

	return writeWord(w, word)
}
