package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FaddD — fadd.d fd, fs1, fs2.
type FaddD struct {
	base

	rd, rs1, rs2 string
	rm           imm // rounding mode (not shown in text)
}

func decodeFaddD(w uint32, addr uint64) Instr {
	return FaddD{
		base: newBase(addr, w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func (i FaddD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fadd.d %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i FaddD) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	rm := i.rm.val

	if rm < 0 || rm > 7 {
		return 0, fmt.Errorf("rounding mode %d out of range", rm)
	}

	return writeWord(w, riscvEncodings["fadd_d"][0]|uint32(rm)<<12|
		fregBits(i.rd)<<7|fregBits(i.rs1)<<15|fregBits(i.rs2)<<20)
}

func (i FaddD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"fadd.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64D",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newFaddD(ops []Op) (Instr, error) {
	regs, err := wantFP(ops, "fadd.d", 3)
	if err != nil {
		return nil, err
	}

	return FaddD{
		rd:  regs[0],
		rs1: regs[1],
		rs2: regs[2],
		rm:  wantRM(ops, 3),
	}, nil
}
