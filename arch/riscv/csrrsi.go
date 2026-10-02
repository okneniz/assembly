package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrrsi - csrrsi rd, csr, zimm.
type Csrrsi struct {
	csrOp

	zimm imm
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
