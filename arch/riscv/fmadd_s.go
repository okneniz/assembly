package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FmaddS - fmadd.s fd, fs1, fs2, fs3.
type FmaddS struct {
	base

	rd, rs1, rs2, rs3 string
	rm                imm
}

func (i FmaddS) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fmadd.s %s, %s, %s, %s", i.rd, i.rs1, i.rs2, i.rs3)
}

func (i FmaddS) Encode(w io.Writer, o EncOpts) (int64, error) {
	rm := i.rm.val

	if rm < 0 || rm > 7 {
		return 0, fmt.Errorf("rounding mode %d out of range", rm)
	}

	return writeWord(w, riscvEncodings["fmadd_s"][0]|uint32(rm)<<12|
		fregBits(i.rd)<<7|fregBits(i.rs1)<<15|fregBits(i.rs2)<<20|fregBits(i.rs3)<<27)
}

func newFmaddS(ops []Op) (Instr, error) {
	regs, err := wantFP(ops, "fmadd.s", 4)
	if err != nil {
		return nil, err
	}

	return FmaddS{
		rd:  regs[0],
		rs1: regs[1],
		rs2: regs[2],
		rs3: regs[3],
		rm:  wantRM(ops, 4),
	}, nil
}
