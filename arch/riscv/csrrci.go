package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrrci - csrrci rd, csr, zimm.
type Csrrci struct {
	base
	csrOp

	zimm imm
}

// Csrrci - csrrci rd, csr, zimm; csr is a 12-bit CSR number
// 0..4095, zimm a 5-bit immediate 0..31.
func (Builder) Csrrci(rd Reg, csr uint16, zimm uint8) Instr {
	return Csrrci{
		csrOp: newCsrOp(rd.name(), int64(csr)),
		zimm:  immNum(int64(zimm)),
	}
}

func decodeCsrrci(w uint32, addr uint64) Instr {
	return Csrrci{
		base:  newBase(addr, w),
		csrOp: newCsrOp(rvRegNames[w>>7&0x1f], int64(w>>20&0xfff)),
		zimm:  immNum(int64(w >> 15 & 0x1f)),
	}
}

func (i Csrrci) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("csrrci %s, %s, %s", i.rd, i.text(), i.zimm.text())
}

func (i Csrrci) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	csr := i.csrBits()

	z, err := zimmBits(i.zimm)
	if err != nil {
		return 0, err
	}

	return writeWord(w, riscvEncodings["csrrci"][0]|regBits(i.rd)<<7|z<<15|csr<<20)
}

func (i Csrrci) MarshalJSON() ([]byte, error) {
	return i.marshalDTO("csrrci", i.ObjDump(disasm.DefaultViewCtx()), "Zicsr",
		map[string]any{"rd": i.rd, "csr": i.text()})
}

func newCsrrci(ops []Op) (Instr, error) {
	return newCsrI(ops, "csrrci")
}
