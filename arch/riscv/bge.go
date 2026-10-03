package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bge - bge rs1, rs2, off; pseudo: blez (rs1=zero), bgez (rs2=zero).
type Bge struct {
	rs1, rs2 string
	off      imm // pc-relative byte offset
}

func (i Bge) Encode(w io.Writer, o EncOpts) (int64, error) {
	bits, err := encB(i.off.val)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["bge"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits

	return writeWord(w, word)
}

func (i Bge) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	if i.rs1 == "zero" {
		return fmt.Sprintf("blez %s, %s", i.rs2, target.text())
	}

	if i.rs2 == "zero" {
		return fmt.Sprintf("bgez %s, %s", i.rs1, target.text())
	}

	return fmt.Sprintf("bge %s, %s, %s", i.rs1, i.rs2, target.text())
}

func newBge(ops []Op) (Instr, error) {
	rs1, rs2, t, err := wantR2T(ops, "bge")
	if err != nil {
		return nil, err
	}

	return Bge{
		rs1: rs1,
		rs2: rs2,
		off: t,
	}, nil
}
