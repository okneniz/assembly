package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Beq - beq rs1, rs2, off; pseudo: beqz.
type Beq struct {
	base

	rs1, rs2 string
	off      imm // pc-relative byte offset
}

// Beq - beq rs1, rs2, off (the pc-relative byte offset; the absolute
// target is off + the instruction address).
func (Builder) Beq(rs1, rs2 Reg, off int64) Instr {
	return Beq{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off),
	}
}

func decodeBeq(w uint32) Instr {
	return Beq{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(bImm(w)),
	}
}

// cBeq - compressed forms (c.beqz): base - halfword, length 2.
func cBeq(h uint32, rs1, rs2 string, off int64) Beq {
	return Beq{
		base: newHalfBase(h),
		rs1:  rs1,
		rs2:  rs2,
		off:  immNum(off),
	}
}

func (i Beq) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	if i.rs2 == "zero" {
		return fmt.Sprintf("beqz %s, %s", i.rs1, target.text())
	}

	return fmt.Sprintf("beq %s, %s, %s", i.rs1, i.rs2, target.text())
}

func (i Beq) Encode(w io.Writer, o EncOpts) (int64, error) {
	bits, err := encB(i.off.val)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["beq"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits
	if !o.NoRVC {
		if half, ok := cbeqz(i.rs1, i.rs2, "beq", i.off.val); ok {
			return writeHalf(w, half)
		}
	}

	return writeWord(w, word)
}

func newBeq(ops []Op) (Instr, error) {
	rs1, rs2, off, err := wantR2T(ops, "beq")
	if err != nil {
		return nil, err
	}

	return Beq{
		rs1: rs1,
		rs2: rs2,
		off: off,
	}, nil
}
