package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrrw - csrrw rd, csr, rs1; pseudo: csrw (rd=zero), fs*/fscsr/fsrm/
// fsflags (rd=zero + status CSR).
type Csrrw struct {
	csrOp

	rs1 string
}

func (i Csrrw) ObjDump(_ disasm.ViewCtx) string {
	if i.rd == "zero" { // write forms (rd == x0)
		switch i.text() {
		case "fflags":
			return "fsflags " + i.rs1
		case "frm":
			return "fsrm " + i.rs1
		case "fcsr":
			return "fscsr " + i.rs1
		}

		return fmt.Sprintf("csrw %s, %s", i.text(), i.rs1)
	}

	return fmt.Sprintf("csrrw %s, %s, %s", i.rd, i.text(), i.rs1)
}

func (i Csrrw) Encode(w io.Writer, o EncOpts) (int64, error) {
	csr := i.csrBits()

	return writeWord(w, riscvEncodings["csrrw"][0]|regBits(i.rd)<<7|regBits(i.rs1)<<15|csr<<20)
}

func newCsrrw(ops []Op) (Instr, error) {
	return newCsrr(ops, "csrrw")
}
