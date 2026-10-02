package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Blt - blt rs1, rs2, off; pseudo: bgtz (rs1=zero), bltz (rs2=zero).
type Blt struct {
	rs1, rs2 string
	off      imm // pc-relative byte offset
}

func (i Blt) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	if i.rs1 == "zero" {
		return fmt.Sprintf("bgtz %s, %s", i.rs2, target.text())
	}

	if i.rs2 == "zero" {
		return fmt.Sprintf("bltz %s, %s", i.rs1, target.text())
	}

	return fmt.Sprintf("blt %s, %s, %s", i.rs1, i.rs2, target.text())
}

func (i Blt) Encode(w io.Writer, o EncOpts) (int64, error) {
	bits, err := encB(i.off.val)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["blt"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits

	return writeWord(w, word)
}

func newBlt(ops []Op) (Instr, error) {
	rs1, rs2, t, err := wantR2T(ops, "blt")
	if err != nil {
		return nil, err
	}

	return Blt{
		rs1: rs1,
		rs2: rs2,
		off: t,
	}, nil
}
