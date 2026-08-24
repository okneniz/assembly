package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrrw - csrrw rd, csr, rs1; pseudo: csrw (rd=zero), fs*/fscsr/fsrm/
// fsflags (rd=zero + status CSR).
type Csrrw struct {
	base
	csrOp

	rs1 string
}

func decodeCsrrw(w uint32, addr uint64) Instr {
	return Csrrw{
		base:  newBase(addr, w),
		csrOp: newCsrOp(rvRegNames[w>>7&0x1f], int64(w>>20&0xfff)),
		rs1:   rvRegNames[w>>15&0x1f],
	}
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

func (i Csrrw) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	csr := i.csrBits()

	return writeWord(w, riscvEncodings["csrrw"][0]|regBits(i.rd)<<7|regBits(i.rs1)<<15|csr<<20)
}

func (i Csrrw) MarshalJSON() ([]byte, error) {
	return i.marshalDTO("csrrw", i.ObjDump(disasm.DefaultViewCtx()), "Zicsr",
		map[string]any{"rd": i.rd, "csr": i.text(), "rs1": i.rs1})
}

func newCsrrw(ops []Op) (Instr, error) {
	return newCsrr(ops, "csrrw")
}
