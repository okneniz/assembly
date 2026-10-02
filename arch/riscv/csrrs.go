package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrrs - csrrs rd, csr, rs1; pseudo: csrr (rs1=zero), frflags/frrm/frcsr/
// rdcycle/rdtime/rdinstret (rs1=zero + status CSR).
type Csrrs struct {
	csrOp

	rs1 string
}

func (i Csrrs) ObjDump(_ disasm.ViewCtx) string {
	if i.rs1 == "zero" { // read forms (rs1 == x0)
		switch i.text() {
		case "fflags":
			return "frflags " + i.rd
		case "frm":
			return "frrm " + i.rd
		case "fcsr":
			return "frcsr " + i.rd
		case "cycle":
			return "rdcycle " + i.rd
		case "time":
			return "rdtime " + i.rd
		case "instret":
			return "rdinstret " + i.rd
		}

		return fmt.Sprintf("csrr %s, %s", i.rd, i.text())
	}

	return fmt.Sprintf("csrrs %s, %s, %s", i.rd, i.text(), i.rs1)
}

func (i Csrrs) Encode(w io.Writer, o EncOpts) (int64, error) {
	csr := i.csrBits()

	return writeWord(w, riscvEncodings["csrrs"][0]|regBits(i.rd)<<7|regBits(i.rs1)<<15|csr<<20)
}

func newCsrrs(ops []Op) (Instr, error) {
	return newCsrr(ops, "csrrs")
}
