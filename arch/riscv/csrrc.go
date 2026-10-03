package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrrc - csrrc rd, csr, rs1.
type Csrrc struct {
	csrOp

	rs1 string
}

func (i Csrrc) Encode(w io.Writer, o EncOpts) (int64, error) {
	csr := i.csrBits()

	return writeWord(w, riscvEncodings["csrrc"][0]|regBits(i.rd)<<7|regBits(i.rs1)<<15|csr<<20)
}

func (i Csrrc) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("csrrc %s, %s, %s", i.rd, i.text(), i.rs1)
}

func newCsrrc(ops []Op) (Instr, error) {
	return newCsrr(ops, "csrrc")
}
