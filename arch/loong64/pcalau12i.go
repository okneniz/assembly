package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Pcalau12i - pcalau12i rd, si20 (1RI20): rd = (pc + si20<<12) with the
// low 12 bits cleared. This is an address computation, not a jump: the
// raw si20 is stored.
type Pcalau12i struct {
	base

	rd  uint8
	imm imm
}

func decodePcalau12i(w uint32) Instr {
	return Pcalau12i{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		imm:  immNum(sField(w, 5, 20)),
	}
}

func (i Pcalau12i) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("pcalau12i %s, %s", laRegName(i.rd), i.imm.text())
}

func (i Pcalau12i) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["pcalau12i"][0] | uint32(i.rd) | scatterS(i.imm.val, 5, 20)

	return writeWord(w, word)
}
