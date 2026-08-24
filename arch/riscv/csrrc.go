package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csrrc — csrrc rd, csr, rs1.
type Csrrc struct {
	base
	csrOp

	rs1 string
}

func decodeCsrrc(w uint32, addr uint64) Instr {
	return Csrrc{
		base:  newBase(addr, w),
		csrOp: newCsrOp(rvRegNames[w>>7&0x1f], int64(w>>20&0xfff)),
		rs1:   rvRegNames[w>>15&0x1f],
	}
}

func (i Csrrc) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("csrrc %s, %s, %s", i.rd, i.text(), i.rs1)
}

func (i Csrrc) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	csr := i.csrBits()

	return writeWord(w, riscvEncodings["csrrc"][0]|regBits(i.rd)<<7|regBits(i.rs1)<<15|csr<<20)
}

func (i Csrrc) MarshalJSON() ([]byte, error) {
	return i.marshalDTO("csrrc", i.ObjDump(disasm.DefaultViewCtx()), "Zicsr",
		map[string]any{"rd": i.rd, "csr": i.text(), "rs1": i.rs1})
}

func newCsrrc(ops []Op) (Instr, error) {
	return newCsrr(ops, "csrrc")
}
