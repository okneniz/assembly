package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bne - bne rs1, rs2, off; pseudo: bnez.
type Bne struct {
	base

	rs1, rs2 string
	off      imm // pc-relative byte offset
}

// Bne - bne rs1, rs2, off (the pc-relative byte offset; the absolute target is off + the instruction address).
func (Builder) Bne(rs1, rs2 Reg, off int64) Instr {
	return Bne{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off),
	}
}

func decodeBne(w uint32) Instr {
	return Bne{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(bImm(w)),
	}
}

// cBne - compressed forms (c.bnez): base - halfword, length 2.
func cBne(h uint32, rs1, rs2 string, off int64) Bne {
	return Bne{
		base: newHalfBase(h),
		rs1:  rs1,
		rs2:  rs2,
		off:  immNum(off),
	}
}

func (i Bne) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	if i.rs2 == "zero" {
		return fmt.Sprintf("bnez %s, %s", i.rs1, target.text())
	}

	return fmt.Sprintf("bne %s, %s, %s", i.rs1, i.rs2, target.text())
}

func (i Bne) Encode(w io.Writer, o EncOpts) (int64, error) {
	bits, err := encB(i.off.val)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["bne"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits
	if !o.NoRVC {
		if half, ok := cbeqz(i.rs1, i.rs2, "bne", i.off.val); ok {
			return writeHalf(w, half)
		}
	}

	return writeWord(w, word)
}

func newBne(ops []Op) (Instr, error) {
	rs1, rs2, t, err := wantR2T(ops, "bne")
	if err != nil {
		return nil, err
	}

	return Bne{
		rs1: rs1,
		rs2: rs2,
		off: t,
	}, nil
}
