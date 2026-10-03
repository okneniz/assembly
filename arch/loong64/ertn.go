package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ertn - ertn: return from exception (no operands).
type Ertn struct {
}

func (i Ertn) Encode(w io.Writer) (int64, error) {
	return writeWord(w, loongEncodings["ertn"][0])
}

func (i Ertn) ObjDump(_ disasm.ViewCtx) string {
	return "ertn"
}
