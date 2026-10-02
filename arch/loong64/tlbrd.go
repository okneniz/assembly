package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Tlbrd - tlbrd: read the indexed TLB entry into the TLB entry
// registers (no operands).
type Tlbrd struct {
	base
}

func (i Tlbrd) ObjDump(_ disasm.ViewCtx) string {
	return "tlbrd"
}

func (i Tlbrd) Encode(w io.Writer) (int64, error) {
	return writeWord(w, loongEncodings["tlbrd"][0])
}
