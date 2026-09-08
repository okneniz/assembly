package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bltu - bltu rs1, rs2, off.
type Bltu struct {
	base

	rs1, rs2 string
	off      imm // pc-relative byte offset
}

// Bltu - bltu rs1, rs2, off (the pc-relative byte offset; the absolute target is off + the instruction address).
func (Builder) Bltu(rs1, rs2 Reg, off int64) Instr {
	return Bltu{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off),
	}
}

func decodeBltu(w uint32) Instr {
	return Bltu{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(bImm(w)),
	}
}

func (i Bltu) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	return fmt.Sprintf("bltu %s, %s, %s", i.rs1, i.rs2, target.text())
}

func (i Bltu) Encode(w io.Writer, o EncOpts) (int64, error) {
	bits, err := encB(i.off.val)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["bltu"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits

	return writeWord(w, word)
}

func newBltu(ops []Op) (Instr, error) {
	rs1, rs2, t, err := wantR2T(ops, "bltu")
	if err != nil {
		return nil, err
	}

	return Bltu{
		rs1: rs1,
		rs2: rs2,
		off: t,
	}, nil
}
