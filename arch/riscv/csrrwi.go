package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrrwi - csrrwi rd, csr, zimm; pseudo: csrwi, fscsri/fsrmi/fsflagsi.
type Csrrwi struct {
	csrOp

	zimm imm
}

func (i Csrrwi) Encode(w io.Writer, o EncOpts) (int64, error) {
	csr := i.csrBits()

	z, err := zimmBits(i.zimm)
	if err != nil {
		return 0, err
	}

	return writeWord(w, riscvEncodings["csrrwi"][0]|regBits(i.rd)<<7|z<<15|csr<<20)
}

func (i Csrrwi) ObjDump(_ disasm.ViewCtx) string {
	if i.rd == "zero" {
		switch i.text() {
		case "fflags":
			return "fsflagsi " + i.zimm.text()
		case "frm":
			return "fsrmi " + i.zimm.text()
		case "fcsr":
			return "fscsri " + i.zimm.text()
		}

		return fmt.Sprintf("csrwi %s, %s", i.text(), i.zimm.text())
	}

	return fmt.Sprintf("csrrwi %s, %s, %s", i.rd, i.text(), i.zimm.text())
}

func newCsrrwi(ops []Op) (Instr, error) {
	return newCsrI(ops, "csrrwi")
}
