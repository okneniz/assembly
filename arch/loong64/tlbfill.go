package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Tlbfill - tlbfill: fill the TLB from the TLB entry registers (no
// operands).
type Tlbfill struct {
}

func (i Tlbfill) Encode(w io.Writer) (int64, error) {
	return writeWord(w, loongEncodings["tlbfill"][0])
}

func (i Tlbfill) ObjDump(_ disasm.ViewCtx) string {
	return "tlbfill"
}
