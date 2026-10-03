package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrxchg - csrxchg rd, rj, csr: the CSR[csr] bits selected by the rj
// write mask get rd's bits; the old value to rd.
type Csrxchg struct {
	rd, rj uint8
	csr    imm
}

func (i Csrxchg) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["csrxchg"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.csr.val, 10, 14)

	return writeWord(w, word)
}

func (i Csrxchg) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("csrxchg %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.csr.text())
}
