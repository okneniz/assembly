package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Tlbflush - tlbflush: invalidate all TLB entries (no operands).
type Tlbflush struct {
}

func (i Tlbflush) Encode(w io.Writer) (int64, error) {
	return writeWord(w, loongEncodings["tlbflush"][0])
}

func (i Tlbflush) ObjDump(_ disasm.ViewCtx) string {
	return "tlbflush"
}
