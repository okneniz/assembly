package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrxchg - csrxchg rd, rj, csr: the CSR[csr] bits selected by the rj
// write mask get rd's bits; the old value to rd.
type Csrxchg struct {
	base

	rd, rj uint8
	csr    imm
}

// Csrxchg - csrxchg rd, rj, csr.
func (Builder) Csrxchg(rd, rj Reg, csr UImm14) Instr {
	return Csrxchg{
		rd:  rd.Num(),
		rj:  rj.Num(),
		csr: immNum(csr.Val()),
	}
}

func decodeCsrxchg(w uint32) Instr {
	return Csrxchg{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		csr:  immNum(int64(uField(w, 10, 14))),
	}
}

func (i Csrxchg) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("csrxchg %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.csr.text())
}

func (i Csrxchg) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["csrxchg"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterU(i.csr.val, 10, 14)

	return writeWord(w, word)
}
