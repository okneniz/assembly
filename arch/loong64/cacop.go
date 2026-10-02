package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Cacop - cacop op, rj, si12: the cache operation op on the block at
// rj + si12 (an unscaled byte offset).
type Cacop struct {
	base

	op  imm
	rj  uint8
	off imm
}

func decodeCacop(w uint32) Instr {
	return Cacop{
		base: newBase(w),
		op:   immNum(int64(uField(w, 0, 5))),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 12)),
	}
}

func (i Cacop) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("cacop %s, %s, %s", i.op.text(), laRegName(i.rj), i.off.text())
}

func (i Cacop) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["cacop"][0] |
		scatterU(i.op.val, 0, 5) | uint32(i.rj)<<5 | scatterS(i.off.val, 10, 12)

	return writeWord(w, word)
}
