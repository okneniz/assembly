package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrrci - csrrci rd, csr, zimm.
type Csrrci struct {
	csrOp

	zimm imm
}

func (i Csrrci) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("csrrci %s, %s, %s", i.rd, i.text(), i.zimm.text())
}

func (i Csrrci) Encode(w io.Writer, o EncOpts) (int64, error) {
	csr := i.csrBits()

	z, err := zimmBits(i.zimm)
	if err != nil {
		return 0, err
	}

	return writeWord(w, riscvEncodings["csrrci"][0]|regBits(i.rd)<<7|z<<15|csr<<20)
}

func newCsrrci(ops []Op) (Instr, error) {
	return newCsrI(ops, "csrrci")
}
