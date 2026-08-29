package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Pcaddi - pcaddi rd, si20 (1RI20): rd = pc + (si20 << 2). This is an
// address computation, not a jump: the raw si20 is stored.
type Pcaddi struct {
	base

	rd  uint8
	imm imm
}

// NewPcaddi - pcaddi rd, si20.
func NewPcaddi(rd Reg, v Imm20) Instr {
	return Pcaddi{
		rd:  rd.Num(),
		imm: immNum(v.Val()),
	}
}

func decodePcaddi(w uint32, addr uint64) Instr {
	return Pcaddi{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		imm:  immNum(sField(w, 5, 20)),
	}
}

func (i Pcaddi) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("pcaddi %s, %s", laRegName(i.rd), i.imm.text())
}

func (i Pcaddi) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["pcaddi"][0] | uint32(i.rd) | scatterS(i.imm.val, 5, 20)

	return writeWord(w, word)
}

func (i Pcaddi) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"pcaddi",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "imm": i.imm.val},
	)
}
