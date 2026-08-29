package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrrd - csrrd rd, csr (rd + a ui14 csr number): rd = CSR[csr].
type Csrrd struct {
	base

	rd  uint8
	csr imm
}

// NewCsrrd - csrrd rd, csr.
func NewCsrrd(rd Reg, csr UImm14) Instr {
	return Csrrd{
		rd:  rd.Num(),
		csr: immNum(csr.Val()),
	}
}

func decodeCsrrd(w uint32, addr uint64) Instr {
	return Csrrd{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		csr:  immNum(int64(uField(w, 10, 14))),
	}
}

func (i Csrrd) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("csrrd %s, %s", laRegName(i.rd), i.csr.text())
}

func (i Csrrd) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["csrrd"][0] |
		uint32(i.rd) | scatterU(i.csr.val, 10, 14)

	return writeWord(w, word)
}

func (i Csrrd) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"csrrd",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		map[string]any{"rd": laRegName(i.rd), "csr": i.csr.val},
	)
}
