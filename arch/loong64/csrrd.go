package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrrd - csrrd rd, csr (rd + a ui14 csr number): rd = CSR[csr].
type Csrrd struct {
	rd  uint8
	csr imm
}

func (i Csrrd) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["csrrd"][0] |
		uint32(i.rd) | scatterU(i.csr.val, 10, 14)

	return writeWord(w, word)
}

func (i Csrrd) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("csrrd %s, %s", laRegName(i.rd), i.csr.text())
}
