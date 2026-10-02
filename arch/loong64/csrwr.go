package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrwr - csrwr rd, csr (rd + a ui14 csr number): CSR[csr] = rd, the
// old value to rd.
type Csrwr struct {
	base

	rd  uint8
	csr imm
}

func (i Csrwr) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("csrwr %s, %s", laRegName(i.rd), i.csr.text())
}

func (i Csrwr) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["csrwr"][0] |
		uint32(i.rd) | scatterU(i.csr.val, 10, 14)

	return writeWord(w, word)
}
