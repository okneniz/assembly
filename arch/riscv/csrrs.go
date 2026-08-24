package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrrs - csrrs rd, csr, rs1; pseudo: csrr (rs1=zero), frflags/frrm/frcsr/
// rdcycle/rdtime/rdinstret (rs1=zero + status CSR).
type Csrrs struct {
	base
	csrOp

	rs1 string
}

func decodeCsrrs(w uint32, addr uint64) Instr {
	return Csrrs{
		base:  newBase(addr, w),
		csrOp: newCsrOp(rvRegNames[w>>7&0x1f], int64(w>>20&0xfff)),
		rs1:   rvRegNames[w>>15&0x1f],
	}
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

func (i Csrrs) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	csr := i.csrBits()

	return writeWord(w, riscvEncodings["csrrs"][0]|regBits(i.rd)<<7|regBits(i.rs1)<<15|csr<<20)
}

func (i Csrrs) MarshalJSON() ([]byte, error) {
	return i.marshalDTO("csrrs", i.ObjDump(disasm.DefaultViewCtx()), "Zicsr",
		map[string]any{"rd": i.rd, "csr": i.text(), "rs1": i.rs1})
}

func newCsrrs(ops []Op) (Instr, error) {
	return newCsrr(ops, "csrrs")
}
