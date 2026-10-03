package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FnmaddS - fnmadd.s fd, fs1, fs2, fs3.
type FnmaddS struct {
	rd, rs1, rs2, rs3 string
	rm                imm
}

func (i FnmaddS) Encode(w io.Writer, o EncOpts) (int64, error) {
	rm := i.rm.val

	if rm < 0 || rm > 7 {
		return 0, fmt.Errorf("rounding mode %d out of range", rm)
	}

	return writeWord(w, riscvEncodings["fnmadd_s"][0]|uint32(rm)<<12|
		fregBits(i.rd)<<7|fregBits(i.rs1)<<15|fregBits(i.rs2)<<20|fregBits(i.rs3)<<27)
}

func (i FnmaddS) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fnmadd.s %s, %s, %s, %s", i.rd, i.rs1, i.rs2, i.rs3)
}

func newFnmaddS(ops []Op) (Instr, error) {
	regs, err := wantFP(ops, "fnmadd.s", 4)
	if err != nil {
		return nil, err
	}

	return FnmaddS{
		rd:  regs[0],
		rs1: regs[1],
		rs2: regs[2],
		rs3: regs[3],
		rm:  wantRM(ops, 4),
	}, nil
}
