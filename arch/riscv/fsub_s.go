package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FsubS - fsub.s fd, fs1, fs2.
type FsubS struct {
	base

	rd, rs1, rs2 string
	rm           imm // rounding mode (not shown in text)
}

func (i FsubS) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fsub.s %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i FsubS) Encode(w io.Writer, o EncOpts) (int64, error) {
	rm := i.rm.val

	if rm < 0 || rm > 7 {
		return 0, fmt.Errorf("rounding mode %d out of range", rm)
	}

	return writeWord(w, riscvEncodings["fsub_s"][0]|uint32(rm)<<12|
		fregBits(i.rd)<<7|fregBits(i.rs1)<<15|fregBits(i.rs2)<<20)
}

func newFsubS(ops []Op) (Instr, error) {
	regs, err := wantFP(ops, "fsub.s", 3)
	if err != nil {
		return nil, err
	}

	return FsubS{
		rd:  regs[0],
		rs1: regs[1],
		rs2: regs[2],
		rm:  wantRM(ops, 3),
	}, nil
}
