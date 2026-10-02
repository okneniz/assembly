package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bgeu - bgeu rs1, rs2, off.
type Bgeu struct {
	rs1, rs2 string
	off      imm // pc-relative byte offset
}

func (i Bgeu) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	return fmt.Sprintf("bgeu %s, %s, %s", i.rs1, i.rs2, target.text())
}

func (i Bgeu) Encode(w io.Writer, o EncOpts) (int64, error) {
	bits, err := encB(i.off.val)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["bgeu"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits

	return writeWord(w, word)
}

func newBgeu(ops []Op) (Instr, error) {
	rs1, rs2, t, err := wantR2T(ops, "bgeu")
	if err != nil {
		return nil, err
	}

	return Bgeu{
		rs1: rs1,
		rs2: rs2,
		off: t,
	}, nil
}
