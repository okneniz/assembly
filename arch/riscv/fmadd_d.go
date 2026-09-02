package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FmaddD - fmadd.d fd, fs1, fs2, fs3.
type FmaddD struct {
	base

	rd, rs1, rs2, rs3 string
	rm                imm
}

// FmaddD - fmadd.d fd, fs1, fs2, fs3; the registers are FP registers taken
// by number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode
// 0..7 (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FmaddD(rd, rs1, rs2, rs3 Reg, rm uint8) Instr {
	return FmaddD{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rs3: fpName(rs3),
		rm:  immNum(int64(rm)),
	}
}

func decodeFmaddD(w uint32, addr uint64) Instr {
	return FmaddD{
		base: newBase(addr, w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rs3:  rvFRegNames[w>>27&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func (i FmaddD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fmadd.d %s, %s, %s, %s", i.rd, i.rs1, i.rs2, i.rs3)
}

func (i FmaddD) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	rm := i.rm.val

	if rm < 0 || rm > 7 {
		return 0, fmt.Errorf("rounding mode %d out of range", rm)
	}

	return writeWord(w, riscvEncodings["fmadd_d"][0]|uint32(rm)<<12|
		fregBits(i.rd)<<7|fregBits(i.rs1)<<15|fregBits(i.rs2)<<20|fregBits(i.rs3)<<27)
}

func (i FmaddD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO("fmadd.d", i.ObjDump(disasm.DefaultViewCtx()), "RV64D",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2, "rs3": i.rs3})
}

func newFmaddD(ops []Op) (Instr, error) {
	regs, err := wantFP(ops, "fmadd.d", 4)
	if err != nil {
		return nil, err
	}

	return FmaddD{
		rd:  regs[0],
		rs1: regs[1],
		rs2: regs[2],
		rs3: regs[3],
		rm:  wantRM(ops, 4),
	}, nil
}
