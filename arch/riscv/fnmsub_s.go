package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FnmsubS — fnmsub.s fd, fs1, fs2, fs3.
type FnmsubS struct {
	base

	rd, rs1, rs2, rs3 string
	rm                imm
}

func decodeFnmsubS(w uint32, addr uint64) Instr {
	return FnmsubS{
		base: newBase(addr, w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rs3:  rvFRegNames[w>>27&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func (i FnmsubS) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fnmsub.s %s, %s, %s, %s", i.rd, i.rs1, i.rs2, i.rs3)
}

func (i FnmsubS) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	rm := i.rm.val

	if rm < 0 || rm > 7 {
		return 0, fmt.Errorf("rounding mode %d out of range", rm)
	}

	return writeWord(w, riscvEncodings["fnmsub_s"][0]|uint32(rm)<<12|
		fregBits(i.rd)<<7|fregBits(i.rs1)<<15|fregBits(i.rs2)<<20|fregBits(i.rs3)<<27)
}

func (i FnmsubS) MarshalJSON() ([]byte, error) {
	return i.marshalDTO("fnmsub.s", i.ObjDump(disasm.DefaultViewCtx()), "RV32F",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2, "rs3": i.rs3})
}

func newFnmsubS(ops []Op) (Instr, error) {
	regs, err := wantFP(ops, "fnmsub.s", 4)
	if err != nil {
		return nil, err
	}

	return FnmsubS{
		rd:  regs[0],
		rs1: regs[1],
		rs2: regs[2],
		rs3: regs[3],
		rm:  wantRM(ops, 4),
	}, nil
}
