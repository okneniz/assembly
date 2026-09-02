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

// Csrwr - csrwr rd, csr.
func (Builder) Csrwr(rd Reg, csr UImm14) Instr {
	return Csrwr{
		rd:  rd.Num(),
		csr: immNum(csr.Val()),
	}
}

func decodeCsrwr(w uint32, addr uint64) Instr {
	return Csrwr{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		csr:  immNum(int64(uField(w, 10, 14))),
	}
}

func (i Csrwr) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("csrwr %s, %s", laRegName(i.rd), i.csr.text())
}

func (i Csrwr) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["csrwr"][0] |
		uint32(i.rd) | scatterU(i.csr.val, 10, 14)

	return writeWord(w, word)
}

func (i Csrwr) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"csrwr",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		map[string]any{"rd": laRegName(i.rd), "csr": i.csr.val},
	)
}
