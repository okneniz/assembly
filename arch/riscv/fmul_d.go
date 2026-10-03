package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FmulD - fmul.d fd, fs1, fs2.
type FmulD struct {
	rd, rs1, rs2 string
	rm           imm // rounding mode (not shown in text)
}

func (i FmulD) Encode(w io.Writer, o EncOpts) (int64, error) {
	rm := i.rm.val

	if rm < 0 || rm > 7 {
		return 0, fmt.Errorf("rounding mode %d out of range", rm)
	}

	return writeWord(w, riscvEncodings["fmul_d"][0]|uint32(rm)<<12|
		fregBits(i.rd)<<7|fregBits(i.rs1)<<15|fregBits(i.rs2)<<20)
}

func (i FmulD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fmul.d %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func newFmulD(ops []Op) (Instr, error) {
	regs, err := wantFP(ops, "fmul.d", 3)
	if err != nil {
		return nil, err
	}

	return FmulD{
		rd:  regs[0],
		rs1: regs[1],
		rs2: regs[2],
		rm:  wantRM(ops, 3),
	}, nil
}
