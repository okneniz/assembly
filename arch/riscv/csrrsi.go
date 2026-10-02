package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrrsi - csrrsi rd, csr, zimm.
type Csrrsi struct {
	base
	csrOp

	zimm imm
}

func decodeCsrrsi(w uint32) Instr {
	return Csrrsi{
		base:  newBase(w),
		csrOp: newCsrOp(rvRegNames[w>>7&0x1f], int64(w>>20&0xfff)),
		zimm:  immNum(int64(w >> 15 & 0x1f)),
	}
}

func (i Csrrsi) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("csrrsi %s, %s, %s", i.rd, i.text(), i.zimm.text())
}

func (i Csrrsi) Encode(w io.Writer, o EncOpts) (int64, error) {
	csr := i.csrBits()

	z, err := zimmBits(i.zimm)
	if err != nil {
		return 0, err
	}

	return writeWord(w, riscvEncodings["csrrsi"][0]|regBits(i.rd)<<7|z<<15|csr<<20)
}

func newCsrrsi(ops []Op) (Instr, error) {
	return newCsrI(ops, "csrrsi")
}
