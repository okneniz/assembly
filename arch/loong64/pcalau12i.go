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

// Pcalau12i - pcalau12i rd, si20.
func (Builder) Pcalau12i(rd Reg, v Imm20) Instr {
	return Pcalau12i{
		rd:  rd.Num(),
		imm: immNum(v.Val()),
	}
}

func decodePcalau12i(w uint32, addr uint64) Instr {
	return Pcalau12i{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		imm:  immNum(sField(w, 5, 20)),
	}
}

func (i Pcalau12i) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("pcalau12i %s, %s", laRegName(i.rd), i.imm.text())
}

func (i Pcalau12i) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["pcalau12i"][0] | uint32(i.rd) | scatterS(i.imm.val, 5, 20)

	return writeWord(w, word)
}

func (i Pcalau12i) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"pcalau12i",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "imm": i.imm.val},
	)
}
