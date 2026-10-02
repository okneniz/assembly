package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Tlbsrch - tlbsrch: search the TLB for an entry matching the address
// in the TLB search registers (no operands).
type Tlbsrch struct {
}

func (i Tlbsrch) ObjDump(_ disasm.ViewCtx) string {
	return "tlbsrch"
}

func (i Tlbsrch) Encode(w io.Writer) (int64, error) {
	return writeWord(w, loongEncodings["tlbsrch"][0])
}
