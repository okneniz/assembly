package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrrwi - csrrwi rd, csr, zimm; pseudo: csrwi, fscsri/fsrmi/fsflagsi.
type Csrrwi struct {
	base
	csrOp

	zimm imm
}

// Csrrwi - csrrwi rd, csr, zimm; csr is a 12-bit CSR number
// 0..4095, zimm a 5-bit immediate 0..31.
func (Builder) Csrrwi(rd Reg, csr uint16, zimm uint8) Instr {
	return Csrrwi{
		csrOp: newCsrOp(rd.name(), int64(csr)),
		zimm:  immNum(int64(zimm)),
	}
}

func decodeCsrrwi(w uint32, addr uint64) Instr {
	return Csrrwi{
		base:  newBase(addr, w),
		csrOp: newCsrOp(rvRegNames[w>>7&0x1f], int64(w>>20&0xfff)),
		zimm:  immNum(int64(w >> 15 & 0x1f)),
	}
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

func (i Csrrwi) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	csr := i.csrBits()

	z, err := zimmBits(i.zimm)
	if err != nil {
		return 0, err
	}

	return writeWord(w, riscvEncodings["csrrwi"][0]|regBits(i.rd)<<7|z<<15|csr<<20)
}

func (i Csrrwi) MarshalJSON() ([]byte, error) {
	return i.marshalDTO("csrrwi", i.ObjDump(disasm.DefaultViewCtx()), "Zicsr",
		map[string]any{"rd": i.rd, "csr": i.text()})
}

func newCsrrwi(ops []Op) (Instr, error) {
	return newCsrI(ops, "csrrwi")
}
