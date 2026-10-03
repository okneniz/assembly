package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Beq - beq rs1, rs2, off; pseudo: beqz.
type Beq struct {
	rs1, rs2 string
	off      imm // pc-relative byte offset
}

// cBeq - compressed forms (c.beqz): base - halfword, length 2.
func cBeq(rs1, rs2 string, off int64) Beq {
	return Beq{
		rs1: rs1,
		rs2: rs2,
		off: immNum(off),
	}
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

func (i Beq) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	if i.rs2 == "zero" {
		return fmt.Sprintf("beqz %s, %s", i.rs1, target.text())
	}

	return fmt.Sprintf("beq %s, %s, %s", i.rs1, i.rs2, target.text())
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
